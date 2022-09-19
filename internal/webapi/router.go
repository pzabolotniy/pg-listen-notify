package webapi

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/pzabolotniy/logging/pkg/logging"
	loggingMW "github.com/pzabolotniy/logging/pkg/middlewares"

	"github.com/pzabolotniy/listen-notify/internal/conf"
)

type HandlerEnv struct {
	DbConn     *pgx.Conn
	EventsConf *conf.Events
}

func PrepareRouter(h *HandlerEnv, logger logging.Logger) *chi.Mux {
	router := chi.NewRouter()
	router.Use(loggingMW.WithLogger(logger))
	router.Post("/events", h.PostEvents)

	return router
}
