package webapi

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pzabolotniy/logging/pkg/logging"
	loggingMW "github.com/pzabolotniy/logging/pkg/middlewares"
	"github.com/riandyrn/otelchi"

	"github.com/pzabolotniy/listen-notify/internal/conf"
)

type HandlerEnv struct {
	DbConn     *pgxpool.Pool
	EventsConf *conf.Events
}

func PrepareRouter(h *HandlerEnv, logger logging.Logger) *chi.Mux {
	router := chi.NewRouter()
	router.Use(
		otelchi.Middleware("notifier-webapi", otelchi.WithChiRoutes(router)),
		loggingMW.WithLogger(logger),
		WithXRequestID,
	)
	router.Post("/events", h.PostEvents)

	return router
}
