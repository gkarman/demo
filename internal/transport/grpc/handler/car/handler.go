package car

import (
	"log/slog"

	grpcCarv1 "github.com/gkarman/demo/api/gen/go/car/v1"
	carservice "github.com/gkarman/demo/internal/service/car"
)

type Handler struct {
	grpcCarv1.UnimplementedCarServer
	log         *slog.Logger
	getService  *carservice.GetService
	listService *carservice.List
}

func New(log *slog.Logger, getService *carservice.GetService, listService *carservice.List) *Handler {
	return &Handler{
		log:         log,
		getService:  getService,
		listService: listService,
	}
}
