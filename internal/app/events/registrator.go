package events

import (
	"context"
	"log/slog"

	carevents "github.com/gkarman/demo/internal/domain/car/events"
	"github.com/gkarman/demo/internal/infrastructure/eventbus"
	"github.com/gkarman/demo/internal/service/car/handlers"
)

func RegisterEventHandlers(bus *eventbus.EventBus, log *slog.Logger) {

	bus.Subscribe("car.created", func(ctx context.Context, e any) {
		event, ok := e.(*carevents.CarCreated)
		if !ok {
			log.Error("RegisterEventHandlers car.created no event")
			return
		}
		CarCreatedLogHandler := handlers.NewCarCreatedLogHandler()
		CarCreatedLogHandler.Handle(ctx, event)
	})
}

