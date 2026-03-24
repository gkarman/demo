package eventbus

import (
	"context"

	"github.com/gkarman/demo/internal/domain/car/events"
)

type EventHandler func(ctx context.Context, event any)

type EventBus struct {
	handlers map[string][]EventHandler
}

func New() *EventBus {
	return &EventBus{
		handlers: make(map[string][]EventHandler),
	}
}

func (b *EventBus) Subscribe(eventName string, handler EventHandler) {
	b.handlers[eventName] = append(b.handlers[eventName], handler)
}

func (b *EventBus) Publish(ctx context.Context, event events.Domain) {
	name := event.EventName()

	if handlers, ok := b.handlers[name]; ok {
		for _, h := range handlers {
			h(ctx, event)
		}
	}
}
