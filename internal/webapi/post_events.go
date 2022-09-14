package webapi

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/pzabolotniy/logging/pkg/logging"

	"github.com/pzabolotniy/listen-notify/internal/db"
)

type CreateEventInput map[string]interface{}

type CreatedEvent struct {
	ReceivedAt time.Time        `json:"received_at"`
	Payload    CreateEventInput `json:"payload"`
	ID         uuid.UUID        `json:"id"`
}

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
	err = db.CreateEvent(ctx, dbConn, dbEvent)
	if err != nil {
		logger.WithError(err).Error("create event failed")
		InternalServerError(ctx, w, "create event failed")

		return
	}

	resp := &CreatedEvent{
		ID:         eventID,
		Payload:    payload,
		ReceivedAt: eventReceivedAt,
	}
	OKResponse(ctx, w, resp)
}
