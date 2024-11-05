package kafka_producer

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	model "github.com/ernstvorsteveld/go-cv-cassandra/adapter/domain/event"
	"github.com/ernstvorsteveld/go-cv-cassandra/pkg/utils"
)

type KafkaProducerContext struct {
	p *kafka.Producer
}

func NewKafkaProducer(config *utils.Configuration) *KafkaProducerContext {
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
	return &KafkaProducerContext{
		p: p,
	}
}

func (k *KafkaProducerContext) Publish(ctx context.Context, eventType string, event model.EventPayload) error {
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

func (k *KafkaProducerContext) Close() error {
	k.p.Close()
	return nil
}

func (k *KafkaProducerContext) Flush(duration time.Duration) int {
	return k.p.Flush(int(duration.Microseconds()))
}

func (k *KafkaProducerContext) Len() int {
	return k.p.Len()
}
