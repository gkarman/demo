package grpc

import (
	carv1 "github.com/gkarman/demo/api/gen/go/car/v1"
	carrepository "github.com/gkarman/demo/internal/repository/car"
	carservice "github.com/gkarman/demo/internal/service/car"
	carhandler "github.com/gkarman/demo/internal/transport/grpc/handler/car"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

func RegisterServices(server *Server, log *slog.Logger, db *pgxpool.Pool) {
	registerCarService(server, log, db)
}

func registerCarService(server *Server, log *slog.Logger, db *pgxpool.Pool) {
	repo := carrepository.New(db)
	getSvc := carservice.NewGet(repo)
	listSvc := carservice.NewList(repo)

	handler := carhandler.New(log, getSvc, listSvc)
	carv1.RegisterCarServer(server.Registrar(), handler)
}
