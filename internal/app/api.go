package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/gkarman/demo/internal/app/events"
	"github.com/gkarman/demo/internal/config"
	"github.com/gkarman/demo/internal/infrastructure/dispatcher"
	"github.com/gkarman/demo/internal/infrastructure/mq"
	grpcTransport "github.com/gkarman/demo/internal/transport/grpc"
	grpcinterceptor "github.com/gkarman/demo/internal/transport/grpc/interceptor"
	httpTransport "github.com/gkarman/demo/internal/transport/http"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

type Api struct {
	log        *slog.Logger
	db         *pgxpool.Pool
	serverHttp *httpTransport.Server
	grpcServer *grpcTransport.Server
	rabbit     *mq.RabbitPublisher
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

	log.Info("rabbit connect...")
	rabbitPublisher, err := initRabbitPublisher(cfg)
	if err != nil {
		return nil, fmt.Errorf("rabbit init: %w", err)
	}
	log.Info("rabbit connected")

	d := dispatcher.New()
	events.RegisterEventHandlers(d, log, rabbitPublisher)

	serverHttp := initHTTPServer(log, postgresDB, cfg, d)
	serverGrpc, err := initGRPCServer(log, postgresDB, cfg, d)
	if err != nil {
		rabbitPublisher.Close()
		return nil, fmt.Errorf("create gRPC server: %w", err)
	}

	return &Api{
		log:        log,
		db:         postgresDB,
		serverHttp: serverHttp,
		grpcServer: serverGrpc,
		rabbit:     rabbitPublisher,
	}, nil
}

func initHTTPServer(log *slog.Logger, db *pgxpool.Pool, cfg *config.Config, d *dispatcher.Dispatcher) *httpTransport.Server {
	router := httpTransport.NewRouter(log, db, d)
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

func initGRPCServer(log *slog.Logger, db *pgxpool.Pool, cfg *config.Config, d *dispatcher.Dispatcher) (*grpcTransport.Server, error) {
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
	grpcTransport.RegisterServices(grpcServer, log, db, d)

	return grpcServer, nil
}

func (a *Api) Run(ctx context.Context) error {
	defer a.db.Close()
	defer a.rabbit.Close()
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
