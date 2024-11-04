package kafka

import (
	"context"
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type EventPayload struct {
	CorrelationId string      `json:"correlationId"`
	EventType     string      `json:"type"`
	Key           string      `json:"key"`
	Payload       interface{} `json:"payload"`
}

type KafkaContext struct {
	p *kafka.Producer
	c *kafka.Consumer
}

func NewKafkaContext(p *kafka.Producer, c *kafka.Consumer) *KafkaContext {
	return &KafkaContext{
		p: p,
		c: c,
	}
}

func (k *KafkaContext) Publish(ctx context.Context, eventType string, event EventPayload) error {
	value, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = k.p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &eventType, Partition: kafka.PartitionAny},
		Key:            []byte(event.Key),
		Value:          value,
	}, nil)

	return err
}
