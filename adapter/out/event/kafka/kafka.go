package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/ernstvorsteveld/go-cv-cassandra/pkg/utils"
)

func NewKafkaProducer(config *utils.Configuration) *kafka.Producer {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": config.EH.Kafka.BootstrapServers,
		"acks":              "all",
	})
	if err != nil {
		fmt.Printf("Failed to create producer: %s", err)
		os.Exit(1)
	}

	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Failed to deliver message: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Produced event to topic %s: key = %-10s value = %s\n",
						*ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))
				}
			}
		}
	}()
	return p
}

func NewKafkaConsumer(config *utils.Configuration) *kafka.Consumer {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": config.EH.Kafka.BootstrapServers,
		"group.id":          "test",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		fmt.Printf("Failed to create consumer: %s", err)
		os.Exit(1)
	}
	return c
}

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
