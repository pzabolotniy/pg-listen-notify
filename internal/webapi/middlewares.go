package webapi

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/pzabolotniy/logging/pkg/logging"
	"go.opentelemetry.io/otel/trace"
)

const LogXTRaceID = "x_trace_id"

func (h *HandlerEnv) WithXTraceID(next http.Handler) http.Handler {
	handlerFn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		traceID := uuid.New().String()
		spCtx := trace.SpanContextFromContext(ctx)
		if spanTraceID := spCtx.TraceID(); spanTraceID.IsValid() {
			traceID = spanTraceID.String()
		}
		ctx = logging.ReplaceFieldsInContext(ctx, logging.Fields{LogXTRaceID: traceID})
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(handlerFn)
}
