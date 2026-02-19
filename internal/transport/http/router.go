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

	listService := car_service.NewListService(carRepo)
	listHandler := car_handler.NewCarListHandler(listService)

	getByIDService := car_service.NewGetByIDService(carRepo)
	getByIDHandler := car_handler.NewGetByIDHandler(getByIDService)

	r.Route("/v1", func(r chi.Router) {
		r.Get("/cars", listHandler.Execute)
		r.Get("/car/{id}", getByIDHandler.Execute)
	})
}
