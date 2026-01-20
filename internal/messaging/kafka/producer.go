package kafka

import (
	"context"
	"strings"

	ckafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type KafkaProducer struct {
	producer *ckafka.Producer
}

func NewKafkaProducer(brokers []string, clientID string) (*KafkaProducer, error) {
	p, err := ckafka.NewProducer(&ckafka.ConfigMap{
		"bootstrap.servers": strings.Join(brokers, ","),
		"client.id":         clientID,
		"acks":              "all",
	})
	if err != nil {
		return nil, err
	}
	return &KafkaProducer{producer: p}, nil
}

func (p *KafkaProducer) Publish(ctx context.Context, topic, key string, payload []byte) error {
	return p.producer.Produce(&ckafka.Message{
		TopicPartition: ckafka.TopicPartition{
			Topic:     &topic,
			Partition: ckafka.PartitionAny,
		},
		Key:   []byte(key),
		Value: payload,
	}, nil)
}

func (p *KafkaProducer) Close() {
	p.producer.Flush(500)
	p.producer.Close()
}
