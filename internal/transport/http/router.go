package http

import (
	"log/slog"

	"github.com/gkarman/demo/internal/repository/car"
	carservice "github.com/gkarman/demo/internal/service/car"
	carhandler "github.com/gkarman/demo/internal/transport/http/handler/car"
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
	repo := car.New(db)

	listSvc := carservice.NewList(repo)
	listHandler := carhandler.NewList(listSvc)

	getCarSvc := carservice.NewGet(repo)
	getCarHandler := carhandler.NewGetCarHandler(getCarSvc)

	r.Route("/v1", func(r chi.Router) {
		r.Get("/cars", listHandler.Handle)
		r.Get("/cars/{id}", getCarHandler.Handle)
	})
}
