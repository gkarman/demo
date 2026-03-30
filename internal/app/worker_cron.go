package app

import (
	"context"

	"github.com/gkarman/demo/internal/config"
	"github.com/gkarman/demo/internal/platform"
	"github.com/gkarman/demo/internal/worker/cron"
)

func NewWorkerCron(_ context.Context) (*cron.Worker, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	log := platform.NewLogger(cfg)
	worker := cron.New(log)

	return worker, nil
}