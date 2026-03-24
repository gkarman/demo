package events

import (
	"time"

	"github.com/google/uuid"
)

type CarCreated struct {
	eventID    string
	ID         string
	Name       string
	occurredAt time.Time
}

func NewCarCreated(id string, name string) *CarCreated {
	return &CarCreated{
		eventID:    uuid.NewString(),
		ID:         id,
		Name:       name,
		occurredAt: time.Now(),
	}
}

func (e *CarCreated) EventID() string {
	return e.eventID
}

func (e *CarCreated) EventName() string {
	return "car.created"
}

func (e *CarCreated) OccurredAt() time.Time {
	return e.occurredAt
}
