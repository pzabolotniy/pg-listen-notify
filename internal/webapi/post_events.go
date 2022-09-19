package webapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pzabolotniy/logging/pkg/logging"

	"github.com/pzabolotniy/listen-notify/internal/db"
)

const MsgCreateEventFailed = "create event failed"

type CreateEventInput map[string]interface{}

type CreatedEvent struct {
	ReceivedAt time.Time        `json:"received_at"`
	Payload    CreateEventInput `json:"payload"`
	ID         uuid.UUID        `json:"id"`
}

//nolint:funlen // db communication should be moved to the separate func later
func (h *HandlerEnv) PostEvents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.FromContext(ctx)
	input := new(CreateEventInput)
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		logger.WithError(err).Error("decode input failed")
		InternalServerError(ctx, w, "decode input failed")

		return
	}

	eventID := uuid.New()
	eventReceivedAt := time.Now().UTC()
	payload := make(map[string]interface{})
	for k, v := range *input {
		payload[k] = v
	}
	dbEvent := &db.Event{
		ID:         eventID,
		Payload:    payload,
		ReceivedAt: eventReceivedAt,
	}
	dbConn := h.DbConn
	tx, err := dbConn.Begin(ctx)
	if err != nil {
		logger.WithError(err).Error("start tx failed")
		InternalServerError(ctx, w, MsgCreateEventFailed)

		return
	}
	defer func() {
		rollbackErr := tx.Rollback(ctx)
		if rollbackErr != nil && !errors.Is(rollbackErr, pgx.ErrTxClosed) {
			logger.WithError(rollbackErr).Error("rollback failed")
		}
	}()

	err = db.CreateEvent(ctx, tx, dbEvent)
	if err != nil {
		logger.WithError(err).Error("create event failed")
		InternalServerError(ctx, w, MsgCreateEventFailed)

		return
	}

	notifyPayload := &db.NotifyPayload{ID: eventID}
	err = db.NotifyEventCh(ctx, tx, h.EventsConf.ChannelName, notifyPayload)
	if err != nil {
		logger.WithError(err).Error("notify failed")
		InternalServerError(ctx, w, MsgCreateEventFailed)

		return
	}

	err = tx.Commit(ctx)
	if err != nil {
		logger.WithError(err).Error("commit failed")
		InternalServerError(ctx, w, MsgCreateEventFailed)

		return
	}
	logger.WithField("event_id", eventID).Trace("event created")
	resp := &CreatedEvent{
		ID:         eventID,
		Payload:    payload,
		ReceivedAt: eventReceivedAt,
	}
	OKResponse(ctx, w, resp)
}
