package app

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/gkarman/demo/internal/config"
	"github.com/gkarman/demo/internal/infrastructure/mq"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WorkerNotify struct {
	log            *slog.Logger
	db             *pgxpool.Pool
	rabbitConsumer *mq.RabbitConsumer
}

func NewWorkerNotify(ctx context.Context) (*WorkerNotify, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	log := initLogger(cfg)

	log.Info("db connect...")
	db, err := initPostgres(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("connect to postgresDB: %w", err)
	}
	log.Info("db connected")

	log.Info("rabbitPusher rabbitConsumer connect...")
	consumer, err := initRabbitConsumer(cfg)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("rabbitPusher rabbitConsumer init: %w", err)
	}
	log.Info("rabbitPusher rabbitConsumer connected")

	return &WorkerNotify{
		log:            log,
		db:             db,
		rabbitConsumer: consumer,
	}, nil
}

func (w *WorkerNotify) Run(ctx context.Context) error {
	defer w.db.Close()
	defer w.rabbitConsumer.Close()

	w.log.Info("worker_notify started")
	go func() {
		err := w.rabbitConsumer.Consume(ctx, w.handleMessage)
		if err != nil {
			w.log.Error("consumer stopped", "error", err)
		}
	}()

	<-ctx.Done()
	w.log.Info("worker_notify shutting down", "reason", ctx.Err())

	return nil
}

func (a *WorkerNotify) handleMessage(body []byte) error {
	a.log.Info("message received", "body", string(body))
	return nil
}

func initRabbitConsumer(cfg *config.Config) (*mq.RabbitConsumer, error) {
	configRabbit := mq.Config{
		User:           cfg.RabbitMQ.User,
		Password:       cfg.RabbitMQ.Password,
		Host:           cfg.RabbitMQ.Host,
		Port:           cfg.RabbitMQ.Port,
		Exchange:       cfg.RabbitMQ.Exchange,
		ReconnectDelay: time.Duration(cfg.RabbitMQ.ReconnectDelay) * time.Second,
	}
	consumer, err := mq.NewRabbitConsumer(
		configRabbit,
		"notify_queue",
	)

	if err != nil {
		return nil, err
	}

	return consumer, nil
}
