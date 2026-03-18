package car

import (
	"log/slog"

	grpcCarv1 "github.com/gkarman/demo/api/gen/go/car/v1"
	carservice "github.com/gkarman/demo/internal/service/car"
)

type Handler struct {
	grpcCarv1.UnimplementedCarServer
	log           *slog.Logger
	getService    *carservice.GetService
	listService   *carservice.List
	createService *carservice.CreateService
}

func New(
	log *slog.Logger,
	getService *carservice.GetService,
	listService *carservice.List,
	createService *carservice.CreateService,
) *Handler {
	return &Handler{
		log:           log,
		getService:    getService,
		listService:   listService,
		createService: createService,
	}
}
