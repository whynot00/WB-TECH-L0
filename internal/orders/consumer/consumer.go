package consumer

import (
	"fmt"
	"kaf-interface/internal/orders/config"
	"kaf-interface/internal/orders/models"
	"log/slog"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

const (
	noTimeout = -1
)

type Handler interface {
	MessageHandler(message []byte, offset kafka.Offset) error
	MigrateOrdersFromDB()
}

type Consumer struct {
	consumer *kafka.Consumer
	handler  Handler
	logger   *slog.Logger
	stop     bool
}

func NewConsumer(handler Handler, logger *slog.Logger, cfg *config.Config) (*Consumer, error) {

	conf := &kafka.ConfigMap{
		"bootstrap.servers":        strings.Join(cfg.Kafka.BootstrapServers, ","),
		"group.id":                 cfg.Kafka.ConsumerGroup,
		"session.timeout.ms":       cfg.Kafka.SessionTimeout,
		"enable.auto.offset.store": false,
		"enable.auto.commit":       true,
		"auto.commit.interval.ms":  5000,
		"auto.offset.reset":        "earliest",
	}

	// создание консьюмера
	c, err := kafka.NewConsumer(conf)
	if err != nil {
		return nil, fmt.Errorf("error with new consumer: %w", err)
	}

	// подписка на топик "orders"
	if err := c.Subscribe(cfg.Kafka.Topic, nil); err != nil {
		return nil, fmt.Errorf("error with consumer subscribe: %w", err)
	}

	return &Consumer{
		consumer: c,
		handler:  handler,
		logger:   logger,
	}, nil
}

func (c *Consumer) Start() {

	// подгрузка данных в кеш из БД
	c.handler.MigrateOrdersFromDB()

	for {

		if c.stop {
			break
		}

		// чтение сообщений из брокера (offser reset - earliest)
		kafkaMsg, err := c.consumer.ReadMessage(noTimeout)
		if err != nil {
			c.logger.Error("error with read message", slog.Any("error", err))
		}

		if kafkaMsg == nil {
			continue
		}

		if err := c.handler.MessageHandler(kafkaMsg.Value, kafkaMsg.TopicPartition.Offset); err != nil {
			switch err {

			// валидация сообщения
			case models.JSONUnmarshalError: // в случае не соответствии модели orders
				c.logger.Warn("unexpected message content")
			default: // иные ошибки
				c.logger.Error("error with handle kaffka message", slog.Any("error", err))
			}
			continue
		}

		// подтверждение обратботки сообщения
		if _, err := c.consumer.StoreMessage(kafkaMsg); err != nil {
			c.logger.Error("error with offset kaffka message", slog.Any("error", err))
			continue
		}
	}
}

func (c *Consumer) Stop() error {
	c.stop = true

	if _, err := c.consumer.Commit(); err != nil {
		return fmt.Errorf("error with commit consumer: %w", err)
	}

	c.logger.Info("commitet offset")

	return c.consumer.Close()
}
