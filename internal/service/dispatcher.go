package service

import (
	"context"

	"github.com/gkarman/demo/internal/domain/car/events"
)

type Dispatcher interface {
	Dispatch(ctx context.Context, events []events.Domain)
}
