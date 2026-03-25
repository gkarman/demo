package handlers

import (
	"context"
	"log/slog"

	carevents "github.com/gkarman/demo/internal/domain/car/events"
)

func CarCreatedLogHandler(log *slog.Logger) func(ctx context.Context, e any) {
	return func(ctx context.Context, e any) {

		event, ok := e.(*carevents.CarCreated)
		if !ok {
			log.Error("invalid event type for car.created")
			return
		}

		log.Info("car created",
			"id", event.ID,
			"name", event.Name,
			"event_id", event.EventID(),
		)
	}
}