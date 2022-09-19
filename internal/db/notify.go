package db

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/pzabolotniy/logging/pkg/logging"
)

func NotifyEventCh(ctx context.Context, dbConn Execer, channelName string, eventID uuid.UUID) error {
	logger := logging.FromContext(ctx)
	query := fmt.Sprintf(`NOTIFY %s, %s`, channelName, QuoteString(eventID.String()))
	_, err := dbConn.Exec(ctx, query)
	if err != nil {
		logger.
			WithError(err).
			WithFields(logging.Fields{
				"channel_name": channelName,
				"payload":      eventID,
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
