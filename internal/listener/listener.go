package listener

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pzabolotniy/logging/pkg/logging"
	"go.opentelemetry.io/otel"

	"github.com/pzabolotniy/listen-notify/internal/conf"
	"github.com/pzabolotniy/listen-notify/internal/db"
)

var ErrOneOfWorkersFailed = errors.New("one of the workers failed")

func Serve(ctx context.Context, logger logging.Logger, dbService *db.DBService, config *conf.Events) error {
	return workers(ctx, logger, dbService, config.ChannelName, config.WorkersCount)
}

func workers(ctx context.Context, logger logging.Logger, dbService *db.DBService,
	channelName string, nWorkers int,
) error {
	wg := new(sync.WaitGroup)
	var err error
	for i := 0; i < nWorkers; i++ {
		wg.Add(1)
		go func(wNum int) {
			logger = logging.FromContext(ctx, logger)
			workerErr := worker(ctx, logger, wNum, dbService, channelName)
			if workerErr != nil {
				logger.WithError(workerErr).Error("worker failed")
				err = ErrOneOfWorkersFailed
			}
			wg.Done()
		}(i)
	}
	wg.Wait()

	return err
}

func worker(ctx context.Context, logger logging.Logger, wNum int, dbService *db.DBService, channelName string) error {
	logger = logging.FromContext(ctx, logger)
	ctx = logging.ReplaceFieldsInContext(ctx, logging.Fields{"worker_num": wNum})
	dbConn, err := dbService.DbConn.Acquire(ctx)
	if err != nil {
		logger.WithError(err).Error("acquire connection from the pool failed")

		return fmt.Errorf("can not acquire connection: %w", err)
	}
	_, err = dbConn.Exec(ctx, fmt.Sprintf("listen %s", channelName))
	if err != nil {
		logger.
			WithError(err).
			WithField("channel_name", channelName).
			Error("pg listen channel failed")

		return fmt.Errorf("pg listen channel failed: %w", err)
	}

	for {
		dbNotification, waitErr := dbConn.Conn().WaitForNotification(ctx)
		if waitErr != nil {
			logger.WithError(waitErr).Error("wait notification failed")

			return fmt.Errorf("wait notification failed: %w", waitErr)
		}
		_ = processPgNotification(ctx, logger, dbService, dbNotification)
	}
}

func processPgNotification(ctx context.Context,
	logger logging.Logger, dbService *db.DBService,
	dbNotification *pgconn.Notification,
) error {
	ctx, span := otel.Tracer("listener").Start(ctx, "process_postgres_notification")
	defer span.End()

	logger = logging.FromContext(ctx, logger)
	logger.WithFields(logging.Fields{
		"pid":          dbNotification.PID,
		"payload":      dbNotification.Payload,
		"channel_name": dbNotification.Channel,
	}).Trace("received pg notification")

	notifyPayload := new(db.NotifyPayload)
	decodeErr := json.NewDecoder(strings.NewReader(dbNotification.Payload)).Decode(notifyPayload)
	if decodeErr != nil {
		logger.
			WithError(decodeErr).
			WithField("raw_payload", dbNotification.Payload).
			Error("decode payload failed")

		return decodeErr
	}

	eventRepo := dbService.NewEventRepository()
	dbEvent, fetchErr := eventRepo.FetchAndLockEvent(ctx, dbService.DbConn, notifyPayload.ID)
	if fetchErr != nil {
		logger.
			WithError(fetchErr).
			WithField("event_id", notifyPayload.ID).
			Error("fetch and lock event failed")

		return fetchErr
	}

	logger.WithFields(logging.Fields{
		"event_id":      dbEvent.ID,
		"event_payload": dbEvent.Payload,
	}).Trace("event payload fetched")

	return nil
}
