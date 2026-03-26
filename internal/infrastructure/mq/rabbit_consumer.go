package mq

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitConsumer struct {
	conn *amqp.Connection
	ch   *amqp.Channel

	queue string
}

func NewRabbitConsumer(cfg Config, queue string, bindings []string) (*RabbitConsumer, error) {
	dsn := fmt.Sprintf("amqp://%s:%s@%s:%s/", cfg.User, cfg.Password, cfg.Host, cfg.Port)

	conn, err := amqp.Dial(dsn)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	_, err = ch.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		conn.Close()
		return nil, err
	}

	for _, key := range bindings {
		err = ch.QueueBind(
			queue,
			key,
			cfg.Exchange,
			false,
			nil,
		)
		if err != nil {
			conn.Close()
			return nil, err
		}
	}

	return &RabbitConsumer{
		conn:  conn,
		ch:    ch,
		queue: queue,
	}, nil
}

func (c *RabbitConsumer) Consume(ctx context.Context, handler func([]byte) error) error {
	msgs, err := c.ch.Consume(
		c.queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	for {
		select {
		case msg := <-msgs:
			err := handler(msg.Body)
			if err != nil {
				_ = msg.Nack(false, true)
				continue
			}
			_ = msg.Ack(false)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (c *RabbitConsumer) Close() error {
	if c.ch != nil {
		_ = c.ch.Close()
	}
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
