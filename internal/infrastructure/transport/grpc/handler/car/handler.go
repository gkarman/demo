package car

import (
	"log/slog"

	grpcCarv1 "github.com/gkarman/demo/api/gen/go/car/v1"
	car2 "github.com/gkarman/demo/internal/application/service/car"
)

type Handler struct {
	grpcCarv1.UnimplementedCarServer
	log           *slog.Logger
	getService    *car2.GetService
	listService   *car2.List
	createService *car2.CreateService
	updateService *car2.UpdateService
	deleteService *car2.DeleteService
}

func New(
	log *slog.Logger,
	getService *car2.GetService,
	listService *car2.List,
	createService *car2.CreateService,
	updateService *car2.UpdateService,
	deleteService *car2.DeleteService,
) *Handler {
	return &Handler{
		log:           log,
		getService:    getService,
		listService:   listService,
		createService: createService,
		updateService: updateService,
		deleteService: deleteService,
	}
}
