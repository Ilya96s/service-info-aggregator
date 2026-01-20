package kafka

import (
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/service-info-aggregator/internal/models/events"
)

type KafkaConsumer struct {
	consumer *kafka.Consumer
	router   *EventRouter
}

func NewKafkaConsumer(brokers []string, groupID string, router *EventRouter) (*KafkaConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  strings.Join(brokers, ","),
		"group.id":           groupID,
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": false,
	})
	if err != nil {
		return nil, err
	}

	return &KafkaConsumer{c, router}, nil
}

func (c *KafkaConsumer) Run(ctx context.Context, topics []string) error {
	if err := c.consumer.SubscribeTopics(topics, nil); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			msg, err := c.consumer.ReadMessage(1 * time.Second)
			if err != nil {
				continue
			}

			if err := c.processMessage(ctx, msg); err != nil {
				slog.Error("message processing failed", "error", err)
				continue
			}

			c.consumer.CommitMessage(msg)
		}
	}
}

func (c *KafkaConsumer) processMessage(ctx context.Context, msg *kafka.Message) error {
	var event events.GenericUpdatedEvent

	if err := json.Unmarshal(msg.Value, &event); err != nil {
		return err
	}

	return c.router.Route(
		ctx,
		event.Type,
		event.Key,
		event.Payload,
	)
}

func (c *KafkaConsumer) Close() {
	c.consumer.Close()
}
