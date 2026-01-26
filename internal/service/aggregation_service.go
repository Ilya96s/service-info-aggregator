package service

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/service-info-aggregator/internal/messaging/kafka"
	"github.com/service-info-aggregator/internal/model/events"
)

type AggregationService struct {
	producer *kafka.KafkaProducer
	topic    string
}

func NewAggregationService(p *kafka.KafkaProducer, topic string) *AggregationService {
	return &AggregationService{
		producer: p,
		topic:    topic,
	}
}

func (s *AggregationService) Execute(ctx context.Context, provider Provider, param string) (any, error) {
	result, err := provider.Fetch(ctx, param)
	if err != nil {
		return nil, err
	}

	event := events.GenericUpdatedEvent{
		Type:      provider.Name(),
		Key:       param,
		Payload:   result,
		Timestamp: time.Now().UTC(),
	}

	bytes, err := json.Marshal(event)
	if err != nil {
		slog.Error("failed to marshal event", "error", err)
	} else {
		if err := s.producer.Publish(ctx, s.topic, param, bytes); err != nil {
			slog.Error("failed to publish event to kafka", "error", err)
		}
	}

	return result, nil
}
