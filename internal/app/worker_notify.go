package app

import (
	"context"
	"log/slog"

	"github.com/gkarman/demo/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WorkerNotify struct {
	log *slog.Logger
	db  *pgxpool.Pool
}

func NewWorkerNotify(ctx context.Context) (*WorkerNotify, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	log := initLogger(cfg)

	log.Info("db connect...")
	db, err := initPostgres(ctx, cfg)
	log.Info("db connected")

	if err != nil {
		return nil, err
	}

	return &WorkerNotify{
		log: log,
		db:  db,
	}, nil
}

func (a *WorkerNotify) Run(ctx context.Context) error {
	defer a.db.Close()

	a.log.Info("worker_notify started")

	<-ctx.Done()

	a.log.Info("worker_notify shutting down", "reason", ctx.Err())

	return nil
}
