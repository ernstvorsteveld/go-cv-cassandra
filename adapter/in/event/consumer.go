package adapter_in_event

import (
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Consumer interface {
	ReadMessage(timeout time.Duration) (*kafka.Message, error)
	SubscribeTopics(topics []string) error
	Close()
}
