package db

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pzabolotniy/logging/pkg/logging"
)

type EventRepository struct {
	Logger logging.Logger
}

func (s *DBService) NewEventRepository() *EventRepository {
	return &EventRepository{Logger: s.Logger}
}

type Event struct {
	ReceivedAt time.Time              `db:"received_at"`
	Payload    map[string]interface{} `db:"payload"`
	ID         uuid.UUID              `db:"id"`
}

type Execer interface {
	Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
}

func (er *EventRepository) CreateEvent(ctx context.Context, dbConn Execer, dbEvent *Event) error {
	logger := logging.FromContext(ctx, er.Logger)
	_, err := dbConn.Exec(ctx, `INSERT INTO events (id, payload, received_at) VALUES ($1, $2, $3)`,
		dbEvent.ID, dbEvent.Payload, dbEvent.ReceivedAt)
	if err != nil {
		logger.WithError(err).Error("insert event failed")

		return err
	}

	return nil
}

type RowContextQueryer interface {
	QueryRow(ctx context.Context, query string, args ...any) pgx.Row
}

func (er *EventRepository) FetchAndLockEvent(ctx context.Context,
	dbConn RowContextQueryer,
	eventID uuid.UUID,
) (*Event, error) {
	logger := logging.FromContext(ctx, er.Logger)
	query := `UPDATE events
SET locked = TRUE
WHERE id IN (
	SELECT id
	FROM events
	WHERE id = $1
	  AND locked = FALSE
	FOR UPDATE SKIP LOCKED)
RETURNING id, payload, received_at`

	dbEvent := new(Event)
	err := dbConn.QueryRow(ctx, query, eventID).Scan(&dbEvent.ID, &dbEvent.Payload, &dbEvent.ReceivedAt)
	if err != nil {
		logger.WithError(err).WithField("event_id", eventID).Error("select event failed")

		return nil, fmt.Errorf("select event failed: %w", err)
	}

	return dbEvent, nil
}
