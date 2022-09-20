package webapi

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/pzabolotniy/logging/pkg/logging"
	"go.opentelemetry.io/otel/trace"
)

func WithXRequestID(next http.Handler) http.Handler {
	handlerFn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := logging.FromContext(ctx)
		traceID := uuid.New().String()
		spCtx := trace.SpanContextFromContext(ctx)
		if spanTraceID := spCtx.TraceID(); spanTraceID.IsValid() {
			traceID = spanTraceID.String()
		}
		logger = logger.WithField("x_request_id", traceID)
		ctx = logging.WithContext(ctx, logger)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(handlerFn)
}
