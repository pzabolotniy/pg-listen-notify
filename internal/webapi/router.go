package webapi

import (
	"github.com/go-chi/chi/v5"
	"github.com/pzabolotniy/logging/pkg/logging"
	"github.com/riandyrn/otelchi"

	"github.com/pzabolotniy/listen-notify/internal/conf"
	"github.com/pzabolotniy/listen-notify/internal/db"
)

type HandlerEnv struct {
	DBService  *db.DBService
	EventsConf *conf.Events
	Logger     logging.Logger
}

func PrepareRouter(h *HandlerEnv) *chi.Mux {
	router := chi.NewRouter()
	router.Use(
		otelchi.Middleware("notifier-webapi", otelchi.WithChiRoutes(router)),
		h.WithXTraceID,
	)
	router.Post("/events", h.PostEvents)

	return router
}
