package listener

import (
	"context"
	"errors"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/pzabolotniy/logging/pkg/logging"
)

var ErrOneOfWorkersFailed = errors.New("one of the workers failed")

func workers(ctx context.Context, dbConn *pgx.Conn, //nolint:unused // will be used later
	channelName string, nWorkers int,
) error {
	wg := new(sync.WaitGroup)
	var err error
	for i := 0; i < nWorkers; i++ {
		wg.Add(1)
		go func() {
			logger := logging.FromContext(ctx)
			workerErr := worker(ctx, dbConn, channelName)
			defer wg.Done()
			if workerErr != nil {
				logger.WithError(workerErr).Error("worker failed")
				err = ErrOneOfWorkersFailed
			}
		}()
	}
	wg.Wait()

	return err
}

func worker(ctx context.Context, dbConn *pgx.Conn, _ string) error { //nolint:unused // will be used later
	logger := logging.FromContext(ctx)
	for {
		dbNotification, err := dbConn.WaitForNotification(ctx)
		if err != nil {
			logger.WithError(err).Error("wait notification failed")

			return err
		}
		logger.WithFields(logging.Fields{
			"pid":          dbNotification.PID,
			"payload":      dbNotification.Payload,
			"channel_name": dbNotification.Channel,
		}).Trace("received pg notification")
	}
}
