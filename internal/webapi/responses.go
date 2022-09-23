package webapi

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/pzabolotniy/logging/pkg/logging"
)

type ResponseBody struct {
	Data  any    `json:"data,omitempty"`
	Error string `json:"error,omitempty"`
}

type Response struct {
	HTTPBody   *ResponseBody
	HTTPStatus int
}

func (h *HandlerEnv) InternalServerError(ctx context.Context, w http.ResponseWriter, msg string) {
	respBody := &ResponseBody{
		Error: msg,
	}
	resp := &Response{
		HTTPStatus: http.StatusInternalServerError,
		HTTPBody:   respBody,
	}
	h.makeJSONResponse(ctx, w, resp)
}

func (h *HandlerEnv) OKResponse(ctx context.Context, w http.ResponseWriter, data any) {
	respBody := &ResponseBody{
		Data: data,
	}
	resp := &Response{
		HTTPStatus: http.StatusCreated,
		HTTPBody:   respBody,
	}
	h.makeJSONResponse(ctx, w, resp)
}

func (h *HandlerEnv) makeJSONResponse(ctx context.Context, w http.ResponseWriter, resp *Response) {
	logger := logging.FromContext(ctx, h.Logger)
	fields := logging.FieldsFromContext(ctx)
	if v, ok := fields[LogXTRaceID]; ok {
		if httpHeaderValue, headerOK := v.(string); headerOK {
			w.Header().Add("x-trace-id", httpHeaderValue)
		}
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(resp.HTTPStatus)
	if encodeErr := json.NewEncoder(w).Encode(resp.HTTPBody); encodeErr != nil {
		logger.WithError(encodeErr).Error("encode response failed")
		w.WriteHeader(http.StatusInternalServerError)
	}
}
