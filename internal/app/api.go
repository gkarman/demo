package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/gkarman/demo/internal/config"
	grpcTransport "github.com/gkarman/demo/internal/transport/grpc"
	grpcinterceptor "github.com/gkarman/demo/internal/transport/grpc/interceptor"
	httpTransport "github.com/gkarman/demo/internal/transport/http"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
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
	serverGrpc, err := initGRPCServer(log, postgresDB, cfg)
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

func initGRPCServer(log *slog.Logger, db *pgxpool.Pool, cfg *config.Config) (*grpcTransport.Server, error) {
	grpcConf := grpcTransport.Config{
		Addr: cfg.ServerGRPC.Addr,
	}
	grpcServer, err := grpcTransport.NewServer(
		log,
		grpcConf,
		grpc.ChainUnaryInterceptor(
			grpcinterceptor.Recovery(),
			grpcinterceptor.Logger(log),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("create gRPC server with interceptors: %w", err)
	}
	grpcTransport.RegisterServices(grpcServer, log, db)

	return grpcServer, nil
}

func (a *App) Run(ctx context.Context) error {
	defer a.db.Close()
	a.serverHttp.Start()
	a.grpcServer.Start()

	<-ctx.Done()

	a.log.Info("shutting down application", "reason", ctx.Err())

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	return a.shutdownServers(shutdownCtx)
}

func (a *App) shutdownServers(ctx context.Context) error {

	var (
		wg       sync.WaitGroup
		errCh    = make(chan error, 2)
		joinedErr error
	)

	wg.Add(2)
	go func() {
		defer wg.Done()
		if err := a.serverHttp.Stop(ctx); err != nil {
			a.log.Error("serverHttp shutdown failed", "error", err)
			errCh <- fmt.Errorf("http shutdown: %w", err)
		}
	}()
	go func() {
		defer wg.Done()
		if err := a.grpcServer.Stop(ctx); err != nil {
			a.log.Error("gRPC server shutdown failed", "error", err)
			errCh <- fmt.Errorf("gRPC shutdown: %w", err)
		}
	}()

	wg.Wait()
	close(errCh)

	for err := range errCh {
		joinedErr = errors.Join(joinedErr, err)
	}

	return joinedErr
}
