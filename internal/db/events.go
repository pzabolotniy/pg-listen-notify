package db

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pzabolotniy/logging/pkg/logging"
)

type Event struct {
	ReceivedAt time.Time              `db:"received_at"`
	Payload    map[string]interface{} `db:"payload"`
	ID         uuid.UUID              `db:"id"`
}

type Execer interface {
	Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
}

func CreateEvent(ctx context.Context, dbConn Execer, dbEvent *Event) error {
	logger := logging.FromContext(ctx)
	_, err := dbConn.Exec(ctx, `INSERT INTO events (id, payload, received_at) VALUES ($1, $2, $3)`,
		dbEvent.ID, dbEvent.Payload, dbEvent.ReceivedAt)
	if err != nil {
		logger.WithError(err).Error("insert event failed")

		return err
	}

	return nil
}
