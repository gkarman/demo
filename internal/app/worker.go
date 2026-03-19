package app

import (
	"context"
	"log/slog"

	"github.com/gkarman/demo/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WorkerApp struct {
	log *slog.Logger
	db  *pgxpool.Pool
}

func NewWorker(ctx context.Context) (*WorkerApp, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	log := initLogger(cfg)
	db, err := initPostgres(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return &WorkerApp{
		log: log,
		db:  db,
	}, nil
}

func (a *WorkerApp) Run(ctx context.Context) error {
	defer a.db.Close()

	a.log.Info("worker started")

	<-ctx.Done()

	a.log.Info("worker shutting down", "reason", ctx.Err())

	return nil
}
