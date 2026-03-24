package handlers

import (
	"context"

	carevents "github.com/gkarman/demo/internal/domain/car/events"
	"github.com/gkarman/demo/internal/logger"
)

type СarCreatedLogHandler struct {
}

func NewCarCreatedLogHandler() *СarCreatedLogHandler {
	return &СarCreatedLogHandler{}
}

func (h *СarCreatedLogHandler) Handle(ctx context.Context, e *carevents.CarCreated) {
	log := logger.FromContext(ctx)
	log.Info("car created",
		"id", e.ID,
		"name", e.Name,
	)
}