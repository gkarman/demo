package http

import (
	"log/slog"

	"github.com/gkarman/demo/internal/repository/car"
	"github.com/gkarman/demo/internal/transport/http/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/gkarman/demo/internal/transport/http/handler"
)

func NewRouter(log *slog.Logger, db *pgxpool.Pool) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger(log))

	homeHandler := handler.NewHomeHandler()

	carRepo := car.New(db)
	carHandler := handler.NewCarHandler(carRepo)

	r.Get("/", homeHandler.Home)
	r.Route("/v1", func(r chi.Router) {
		r.Get("/cars", carHandler.GetCars)
	})

	return r
}
