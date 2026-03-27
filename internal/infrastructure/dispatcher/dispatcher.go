package dispatcher

import (
	"context"
	"reflect"
)

type Handler func(ctx context.Context, e any)

type Dispatcher struct {
	handlers map[reflect.Type][]Handler
}

func New() *Dispatcher {
	return &Dispatcher{
		handlers: make(map[reflect.Type][]Handler),
	}
}

func (d *Dispatcher) Register(eventType any, h Handler) {
	t := reflect.TypeOf(eventType)
	d.handlers[t] = append(d.handlers[t], h)
}

func (d *Dispatcher) Dispatch(ctx context.Context, events []any) {
	for _, event := range events {
		t := reflect.TypeOf(event)
		if hs, ok := d.handlers[t]; ok {
			for _, h := range hs {
				h(ctx, event)
			}
		}
	}
}
