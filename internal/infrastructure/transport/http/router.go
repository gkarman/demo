package http

import (
	"log/slog"

	car3 "github.com/gkarman/demo/internal/application/service/car"
	"github.com/gkarman/demo/internal/infrastructure/dispatcher"
	"github.com/gkarman/demo/internal/infrastructure/repository/car"
	"github.com/gkarman/demo/internal/infrastructure/transport/http/handler"
	car2 "github.com/gkarman/demo/internal/infrastructure/transport/http/handler/car"
	middleware2 "github.com/gkarman/demo/internal/infrastructure/transport/http/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewRouter(log *slog.Logger, db *pgxpool.Pool, d *dispatcher.Dispatcher) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware2.Logger(log))
	r.Use(middleware2.Recovery())
	registerHomeRoutes(r)
	registerCarRoutes(r, db, d)
	return r
}

func registerHomeRoutes(r *chi.Mux) {
	homeHandler := handler.NewHomeHandler()
	r.Get("/", homeHandler.Home)
}

func registerCarRoutes(r *chi.Mux, db *pgxpool.Pool, d *dispatcher.Dispatcher) {
	repo := car.New(db)

	listSvc := car3.NewList(repo)
	listHandler := car2.NewList(listSvc)

	getCarSvc := car3.NewGet(repo)
	getCarHandler := car2.NewGetCarHandler(getCarSvc)

	createCarSvc := car3.NewCreate(repo, d)
	createHandler := car2.NewCreate(createCarSvc)

	updateCarSvc := car3.NewUpdate(repo)
	updateHandler := car2.NewUpdate(updateCarSvc)

	deleteCarSvc := car3.NewDelete(repo)
	deleteHandler := car2.NewDelete(deleteCarSvc)

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/cars", createHandler.Handle)
		r.Get("/cars", listHandler.Handle)
		r.Get("/cars/{id}", getCarHandler.Handle)
		r.Put("/cars/{id}", updateHandler.Handle)
		r.Delete("/cars/{id}", deleteHandler.Handle)
	})
}
