package handlers

import (
	"context"
	"encoding/json"
	"log/slog"

	carevents "github.com/gkarman/demo/internal/domain/car/events"
	"github.com/gkarman/demo/internal/infrastructure/events/mappers"
	"github.com/gkarman/demo/internal/infrastructure/mq"
)

func CarCreatedToRabbitHandler(publisher mq.Publisher, log *slog.Logger) func(ctx context.Context, e any) {

	return func(ctx context.Context, e any) {
		event, ok := e.(*carevents.CarCreated)
		if !ok {
			log.Error("invalid event type for car.created (rabbit)")
			return
		}

		msg := mappers.MapCarCreated(event)
		body, err := json.Marshal(msg)
		if err != nil {
			log.Error("marshal failed in CarCreatedToRabbitHandler", "err", err)
			return
		}

		err = publisher.Publish(ctx, "car.created.v1", body)
		if err != nil {
			log.Error("failed to publish to rabbitmq", "err", err)
			return
		}

		log.Info("event published to rabbitmq", "event", event.EventID())
	}
}
