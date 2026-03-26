package platform

import (
	"fmt"
	"time"

	"github.com/gkarman/demo/internal/config"
	"github.com/gkarman/demo/internal/infrastructure/mq"
)

func NewRabbitConsumer(cfg *config.Config) (*mq.RabbitConsumer, error) {
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
		"worker.notify",
		[]string{
			"car.#",
			"user.#",
		},
	)

	if err != nil {
		return nil, fmt.Errorf("unable to create rabbit consumer: %w", err)
	}

	return consumer, nil
}
