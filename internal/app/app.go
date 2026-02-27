package app

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/gkarman/demo/internal/config"
	"github.com/gkarman/demo/internal/db"
	"github.com/gkarman/demo/internal/logger"
	grpcTransport "github.com/gkarman/demo/internal/transport/grpc"
	httpTransport "github.com/gkarman/demo/internal/transport/http"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	log        *slog.Logger
	db         *pgxpool.Pool
	serverHttp *httpTransport.Server
	grpcServer *grpcTransport.Server
}

func New(ctx context.Context) (*App, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	log := initLogger(cfg)

	log.Info("connecting to database")
	postgresDB, err := initPostgres(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("connect to postgresDB: %w", err)
	}
	log.Info("database connected")

	serverHttp := initHTTPServer(log, postgresDB, cfg)
	serverGrpc, err := initGRPCServer(log, cfg)
	if err != nil {
		return nil, fmt.Errorf("create gRPC server: %w", err)
	}

	return &App{
		log:        log,
		db:         postgresDB,
		serverHttp: serverHttp,
		grpcServer: serverGrpc,
	}, nil
}

func initLogger(cfg *config.Config) *slog.Logger {
	log := logger.New(logger.Config{Level: cfg.Logger.Level})
	slog.SetDefault(log)
	return log
}

func initHTTPServer(log *slog.Logger, db *pgxpool.Pool, cfg *config.Config) *httpTransport.Server {
	router := httpTransport.NewRouter(log, db)
	return httpTransport.NewServer(
		log,
		router,
		httpTransport.Config{
			Addr:         cfg.ServerHttp.Addr,
			ReadTimeout:  time.Duration(cfg.ServerHttp.ReadTimeoutSeconds) * time.Second,
			WriteTimeout: time.Duration(cfg.ServerHttp.WriteTimeoutSeconds) * time.Second,
		},
	)
}

func initGRPCServer(log *slog.Logger, cfg *config.Config) (*grpcTransport.Server, error) {
	return grpcTransport.NewServer(
		log,
		grpcTransport.Config{
			Addr: cfg.ServerGRPC.Addr,
		},
	)
}

func (a *App) Run(ctx context.Context) error {
	defer a.db.Close()
	a.serverHttp.Start()
	a.grpcServer.Start()

	<-ctx.Done()

	a.log.Info("shutting down application", "reason", ctx.Err())

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := a.serverHttp.Stop(shutdownCtx); err != nil {
		a.log.Error("serverHttp shutdown failed", "error", err)
		return err
	}
	a.grpcServer.Stop()

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
