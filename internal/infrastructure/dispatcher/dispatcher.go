package dispatcher

import (
	"context"

	carevents "github.com/gkarman/demo/internal/domain/car/events"
)

type Handler func(ctx context.Context, e any)

type Dispatcher struct {
	handlers map[string][]Handler
}

func New() *Dispatcher {
	return &Dispatcher{
		handlers: make(map[string][]Handler),
	}
}

func (d *Dispatcher) Register(eventName string, h Handler) {
	d.handlers[eventName] = append(d.handlers[eventName], h)
}

func (d *Dispatcher) Dispatch(ctx context.Context, events []carevents.Domain) {
	for _, event := range events {
		if hs, ok := d.handlers[event.EventName()]; ok {
			for _, h := range hs {
				h(ctx, event)
			}
		}
	}
}
