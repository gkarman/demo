package grpc

import (
	"log/slog"

	carv1 "github.com/gkarman/demo/api/gen/go/car/v1"
	"github.com/gkarman/demo/internal/application/service/car"
	"github.com/gkarman/demo/internal/infrastructure/dispatcher"
	carrepository "github.com/gkarman/demo/internal/infrastructure/repository/car"
	carhandler "github.com/gkarman/demo/internal/infrastructure/transport/grpc/handler/car"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterServices(server *Server, log *slog.Logger, db *pgxpool.Pool, d *dispatcher.Dispatcher) {
	registerCarService(server, log, db, d)
}

func registerCarService(server *Server, log *slog.Logger, db *pgxpool.Pool, d *dispatcher.Dispatcher) {
	repo := carrepository.New(db)
	getSvc := car.NewGet(repo)
	listSvc := car.NewList(repo)
	createSvc := car.NewCreate(repo, d)
	updateSvc := car.NewUpdate(repo)
	deleteSvc := car.NewDelete(repo)

	handler := carhandler.New(log, getSvc, listSvc, createSvc, updateSvc, deleteSvc)
	carv1.RegisterCarServer(server.Registrar(), handler)
}
