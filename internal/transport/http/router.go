package http

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/gkarman/demo/internal/transport/http/handler"
)

func NewRouter(log *slog.Logger, db *pgxpool.Pool, ) *chi.Mux {
	r := chi.NewRouter()

	homeHandler := handler.NewHomeHandler(log)
	r.Get("/", homeHandler.Home)

	carHandler := handler.NewCarHandler(log, db)
	r.Route("/v1", func(r chi.Router) {
		r.Get("/cars", carHandler.GetCars)
	})

	return r
}
