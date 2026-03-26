package mappers

import (
	carevents "github.com/gkarman/demo/internal/domain/car/events"
	contracts "github.com/gkarman/demo/internal/infrastructure/contracts/events"
)

func MapCarCreated(e *carevents.CarCreated) contracts.CarCreatedV1 {
	return contracts.CarCreatedV1{
		EventType:  contracts.EventCarCreatedV1,
		EventID:    e.EventID(),
		CarID:      e.ID,
		Name:       e.Name,
		OccurredAt: e.OccurredAt(),
	}
}
