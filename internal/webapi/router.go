package webapi

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/pzabolotniy/logging/pkg/logging"
	loggingMW "github.com/pzabolotniy/logging/pkg/middlewares"
)

type HandlerEnv struct {
	DbConn *pgx.Conn
}

func PrepareRouter(h *HandlerEnv, logger logging.Logger) *chi.Mux {
	router := chi.NewRouter()
	router.Use(loggingMW.WithLogger(logger))
	router.Post("/events", h.PostEvents)

	return router
}
