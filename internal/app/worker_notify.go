package app

import (
	"context"
	"fmt"
	"time"

	"github.com/gkarman/demo/internal/config"
	"github.com/gkarman/demo/internal/infrastructure/mq"
	"github.com/gkarman/demo/internal/worker/notify"
)

func NewWorkerNotify(ctx context.Context) (*notify.Worker, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	log := initLogger(cfg)

	log.Info("db connect...")
	db, err := initPostgres(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("connect to postgres: %w", err)
	}
	log.Info("db connected")

	log.Info("rabbit consumer connect...")
	consumer, err := initRabbitConsumer(cfg)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("init rabbit consumer: %w", err)
	}
	log.Info("rabbit consumer connected")

	router := notify.NewRouterWithHandlers(log)
	worker := notify.New(log, consumer, router)

	return worker, nil
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

	return mq.NewRabbitConsumer(configRabbit, "notify_queue")
}
