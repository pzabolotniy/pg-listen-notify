package webapi

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

type HandlerEnv struct {
	DbConn *pgx.Conn
}

func PrepareRouter(h *HandlerEnv) *chi.Mux {
	router := chi.NewRouter()
	router.Post("/events", h.PostEvents)
	return router
}
