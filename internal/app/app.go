package app

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/gkarman/demo/internal/config"
	"github.com/gkarman/demo/internal/db"
	"github.com/gkarman/demo/internal/logger"
	httpTransport "github.com/gkarman/demo/internal/transport/http"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	log    *slog.Logger
	pool   *pgxpool.Pool
	server *httpTransport.Server
}

func New(ctx context.Context) (*App, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	log := logger.New(logger.Config{Level: cfg.Logger.Level})
	slog.SetDefault(log)

	log.Info("connecting to database")
	pool, err := initPostgres(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("connect to pool: %w", err)
	}
	log.Info("database connected")

	router := httpTransport.NewRouter(log, pool)
	server := httpTransport.NewServer(
		log,
		router,
		httpTransport.Config{
			Addr:         cfg.ServerHttp.Addr,
			ReadTimeout:  time.Duration(cfg.ServerHttp.ReadTimeoutSeconds) * time.Second,
			WriteTimeout: time.Duration(cfg.ServerHttp.WriteTimeoutSeconds) * time.Second,
		},
	)

	return &App{
		log:    log,
		pool:   pool,
		server: server,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	defer a.pool.Close()
	a.server.Start()

	<-ctx.Done()

	a.log.Info("shutting down application", "reason", ctx.Err())

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := a.server.Stop(shutdownCtx); err != nil {
		a.log.Error("server shutdown failed", "error", err)
		return err
	}

	return nil
}

func initPostgres(parent context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(parent, 10*time.Second)
	defer cancel()

	return db.NewPool(ctx, db.Config{
		DSN:             cfg.DB.DSN(),
		MaxConns:        cfg.DB.MaxConnections,
		MinConns:        cfg.DB.MinConnections,
		MaxConnLifetime: time.Duration(cfg.DB.MaxConnectionLifeTimeMinutes) * time.Minute,
		MaxConnIdleTime: time.Duration(cfg.DB.MaxConnectionIdleTimeMinutes) * time.Minute,
	})
}
