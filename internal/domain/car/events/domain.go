package events

import "time"

type Domain interface {
	EventID() string
	EventName() string
	OccurredAt() time.Time
}
