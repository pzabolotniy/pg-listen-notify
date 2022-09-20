package webapi

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/pzabolotniy/logging/pkg/logging"
	"go.opentelemetry.io/otel/trace"
)

func WithXTraceID(next http.Handler) http.Handler {
	handlerFn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := logging.FromContext(ctx)
		traceID := uuid.New().String()
		spCtx := trace.SpanContextFromContext(ctx)
		if spanTraceID := spCtx.TraceID(); spanTraceID.IsValid() {
			traceID = spanTraceID.String()
		}
		logger = logger.WithField("x_trace_id", traceID)
		w.Header().Add("x-trace-id", traceID)
		ctx = logging.WithContext(ctx, logger)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(handlerFn)
}
