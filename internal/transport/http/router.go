package http

import (
	"log/slog"

	"github.com/gkarman/demo/internal/repository/car"
	car_service "github.com/gkarman/demo/internal/service/car"
	car_handler "github.com/gkarman/demo/internal/transport/http/handler/car"
	"github.com/gkarman/demo/internal/transport/http/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/gkarman/demo/internal/transport/http/handler"
)

func NewRouter(log *slog.Logger, db *pgxpool.Pool) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger(log))
	registerHomeRoutes(r)
	registerCarRoutes(r, db)
	return r
}

func registerHomeRoutes(r *chi.Mux) {
	homeHandler := handler.NewHomeHandler()
	r.Get("/", homeHandler.Home)
}

func registerCarRoutes(r *chi.Mux, db *pgxpool.Pool) {
	carRepo := car.New(db)

	carListService := car_service.NewListService(carRepo)
	carListHandler := car_handler.NewCarListHandler(carListService)

	carGetByIdHandler := car_handler.NewGetByIdHandler(carRepo)
	r.Route("/v1", func(r chi.Router) {
		r.Get("/cars", carListHandler.Execute)
		r.Get("/car/{id}", carGetByIdHandler.Execute)
	})
}
