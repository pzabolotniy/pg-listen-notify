package db

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"strings"

	"github.com/google/uuid"
	"github.com/pzabolotniy/logging/pkg/logging"
)

type NotifyRepository struct {
	Logger logging.Logger
}

func (s *DBService) NewNotifyRepository() *NotifyRepository {
	return &NotifyRepository{Logger: s.Logger}
}

type NotifyPayload struct {
	ID uuid.UUID `json:"id"`
}

func (ep *NotifyPayload) Value() (driver.Value, error) {
	return json.Marshal(ep)
}

func (nr *NotifyRepository) NotifyEventCh(ctx context.Context,
	dbConn Execer,
	channelName string, payload *NotifyPayload,
) error {
	logger := logging.FromContext(ctx, nr.Logger)
	query := `SELECT pg_notify($1, $2)`
	_, err := dbConn.Exec(ctx, query, channelName, payload)
	if err != nil {
		logger.
			WithError(err).
			WithFields(logging.Fields{
				"channel_name": channelName,
				"payload":      payload,
			}).
			Error("notify failed")

		return err
	}

	return nil
}

// copy/paste from internal package https://github.com/jackc/pgx/blob/master/internal/sanitize/sanitize.go#L84
func QuoteString(str string) string {
	return "'" + strings.ReplaceAll(str, "'", "''") + "'"
}
