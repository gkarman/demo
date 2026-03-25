package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/gkarman/demo/internal/config"
	"github.com/gkarman/demo/internal/infrastructure/dispatcher"
	"github.com/gkarman/demo/internal/infrastructure/events"
	"github.com/gkarman/demo/internal/infrastructure/mq"
	grpc2 "github.com/gkarman/demo/internal/infrastructure/transport/grpc"
	"github.com/gkarman/demo/internal/infrastructure/transport/grpc/interceptor"
	"github.com/gkarman/demo/internal/infrastructure/transport/http"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

type Api struct {
	log        *slog.Logger
	db         *pgxpool.Pool
	serverHttp *http.Server
	grpcServer   *grpc2.Server
	rabbitPusher *mq.RabbitPublisher
}

func NewApi(ctx context.Context) (*Api, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	log := initLogger(cfg)

	log.Info("db connect...")
	postgresDB, err := initPostgres(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("connect to postgresDB: %w", err)
	}
	log.Info("db connected")

	log.Info("rabbitPusher connect...")
	rabbitPublisher, err := initRabbitPublisher(cfg)
	if err != nil {
		return nil, fmt.Errorf("rabbitPusher init: %w", err)
	}
	log.Info("rabbitPusher connected")

	d := dispatcher.New()
	events.RegisterEventHandlers(d, log, rabbitPublisher)

	serverHttp := initHTTPServer(log, postgresDB, cfg, d)
	serverGrpc, err := initGRPCServer(log, postgresDB, cfg, d)
	if err != nil {
		rabbitPublisher.Close()
		return nil, fmt.Errorf("create gRPC server: %w", err)
	}

	return &Api{
		log:          log,
		db:           postgresDB,
		serverHttp:   serverHttp,
		grpcServer:   serverGrpc,
		rabbitPusher: rabbitPublisher,
	}, nil
}

func initHTTPServer(log *slog.Logger, db *pgxpool.Pool, cfg *config.Config, d *dispatcher.Dispatcher) *http.Server {
	router := http.NewRouter(log, db, d)
	return http.NewServer(
		log,
		router,
		http.Config{
			Addr:         cfg.ServerHttp.Addr,
			ReadTimeout:  time.Duration(cfg.ServerHttp.ReadTimeoutSeconds) * time.Second,
			WriteTimeout: time.Duration(cfg.ServerHttp.WriteTimeoutSeconds) * time.Second,
		},
	)
}

func initGRPCServer(log *slog.Logger, db *pgxpool.Pool, cfg *config.Config, d *dispatcher.Dispatcher) (*grpc2.Server, error) {
	grpcConf := grpc2.Config{
		Addr: cfg.ServerGRPC.Addr,
	}
	grpcServer, err := grpc2.NewServer(
		log,
		grpcConf,
		grpc.ChainUnaryInterceptor(
			interceptor.Recovery(),
			interceptor.Logger(log),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("create gRPC server with interceptors: %w", err)
	}
	grpc2.RegisterServices(grpcServer, log, db, d)

	return grpcServer, nil
}

func (a *Api) Run(ctx context.Context) error {
	defer a.db.Close()
	defer a.rabbitPusher.Close()
	a.serverHttp.Start()
	a.grpcServer.Start()

	<-ctx.Done()

	a.log.Info("shutting down application", "reason", ctx.Err())

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	return a.shutdownServers(shutdownCtx)
}

func (a *Api) shutdownServers(ctx context.Context) error {

	var (
		wg        sync.WaitGroup
		errCh     = make(chan error, 2)
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

func initRabbitPublisher(cfg *config.Config) (*mq.RabbitPublisher, error) {
	configRabbit := mq.Config{
		User:           cfg.RabbitMQ.User,
		Password:       cfg.RabbitMQ.Password,
		Host:           cfg.RabbitMQ.Host,
		Port:           cfg.RabbitMQ.Port,
		Exchange:       cfg.RabbitMQ.Exchange,
		ReconnectDelay: time.Duration(cfg.RabbitMQ.ReconnectDelay) * time.Second,
	}

	publisher, err := mq.NewRabbitPublisher(configRabbit)
	if err != nil {
		return nil, err
	}

	return publisher, nil
}
