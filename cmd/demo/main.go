package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gkarman/demo/internal/config"
	"github.com/gkarman/demo/internal/db"
	"github.com/gkarman/demo/internal/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	if err := run(ctx); err != nil {
		slog.Error("application failed", "error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	log := logger.New(logger.Config{Level: cfg.Logger.Level})
	slog.SetDefault(log)

	log.Info("connecting to database")
	pool, err := initPostgres(ctx, cfg)
	if err != nil {
		return fmt.Errorf("connect to db: %w", err)
	}
	log.Info("database connected")

	<-ctx.Done()
	pool.Close()
	log.Info("shutting down application", "reason", ctx.Err())

	return nil
}

func initPostgres(parent context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(parent, 10*time.Second)
	defer cancel()

	pool, err := db.NewPool(ctx, db.Config{
		DSN:             cfg.DB.DSN(),
		MaxConns:        cfg.DB.MaxConnections,
		MinConns:        cfg.DB.MinConnections,
		MaxConnLifetime: time.Duration(cfg.DB.MaxConnectionLifeTimeMinutes) * time.Minute,
		MaxConnIdleTime: time.Duration(cfg.DB.MaxConnectionIdleTimeMinutes) * time.Minute,
	})

	if err != nil {
		return nil, err
	}

	return pool, nil
}
