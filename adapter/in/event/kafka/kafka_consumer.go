package kafka_consumer

import (
	"fmt"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/ernstvorsteveld/go-cv-cassandra/pkg/utils"
)

type KafkaConsumerContext struct {
	c *kafka.Consumer
}

func NewKafkaConsumer(config *utils.Configuration) *KafkaConsumerContext {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": config.EH.Kafka.BootstrapServers,
		"group.id":          "test",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		fmt.Printf("Failed to create consumer: %s", err)
		os.Exit(1)
	}
	return &KafkaConsumerContext{
		c: c,
	}
}

func (c *KafkaConsumerContext) ReadMessage(timeout time.Duration) (*kafka.Message, error) {
	return c.c.ReadMessage(timeout)
}

func (c *KafkaConsumerContext) SubscribeTopics(topics []string) error {
	return c.c.SubscribeTopics(topics, nil)
}

func (c *KafkaConsumerContext) Close() {
	c.c.Close()
}
