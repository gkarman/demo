package cron

import (
	"context"
	"log/slog"
)

type Worker struct {
	log *slog.Logger
}

func New(log *slog.Logger) *Worker {
	return &Worker{
		log: log,
	}
}

func (w *Worker) Run(ctx context.Context) error {
	w.log.Info("worker_cron started")

	go func() {

	}()

	<-ctx.Done()
	w.log.Info("worker_cron shutting down", "reason", ctx.Err())

	return nil
}