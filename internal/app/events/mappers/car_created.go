package mappers

import (
	contracts "github.com/gkarman/demo/internal/contracts/events"
	carevents "github.com/gkarman/demo/internal/domain/car/events"
)

func MapCarCreated(e *carevents.CarCreated) contracts.CarCreatedV1 {
	return contracts.CarCreatedV1{
		EventType:  "car.created",
		EventID:    e.EventID(),
		CarID:      e.ID,
		Name:       e.Name,
		OccurredAt: e.OccurredAt(),
	}
}
