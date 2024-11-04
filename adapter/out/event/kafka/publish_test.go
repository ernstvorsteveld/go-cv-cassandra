package kafka

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/ernstvorsteveld/go-cv-cassandra/domain/model"
	"github.com/ernstvorsteveld/go-cv-cassandra/pkg/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"

	log "github.com/sirupsen/logrus"
	"github.com/testcontainers/testcontainers-go/modules/kafka"
)

var kafkaContext *KafkaContext

func TestMain(m *testing.M) {
	log.Infof("About to start Kafka container in startContainer")
	ctx := context.Background()

	kafkaContainer, err := kafka.Run(ctx,
		"confluentinc/confluent-local:7.7.1",
		kafka.WithClusterID("testcontainers-cluster"),
	)
	defer func() {
		if err := testcontainers.TerminateContainer(kafkaContainer); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	}()
	if err != nil {
		log.Printf("failed to start container: %s", err)
		return
	}

	rawPort, _ := kafkaContainer.MappedPort(ctx, "9093")
	if rawPort == "" {
		log.Printf("rawPort is empty")
	}

	parts := strings.Split(rawPort.Port(), "/")
	log.Printf("Running on port %s", parts[0])

	config := &utils.Configuration{
		EH: utils.EventHandlerConfiguration{
			Kafka: utils.KafkaConfiguration{
				BootstrapServers: "localhost:" + parts[0],
			},
		},
	}
	producer := NewKafkaProducer(config)
	consumer := NewKafkaConsumer(config)
	defer producer.Close()
	defer consumer.Close()

	kafkaContext = NewKafkaContext(producer, consumer)

	m.Run()
}

func Test_should_publish_to_topic(t *testing.T) {
	log.Println("About to publish and consume")
	topic := "test.tag.created"

	for i := 0; i < 10; i++ {
		key := uuid.NewString()
		err := kafkaContext.Publish(context.Background(), topic, EventPayload{
			CorrelationId: uuid.NewString(),
			EventType:     "test",
			Key:           key,
			Payload: model.Tag{
				Id:   key,
				Name: uuid.NewString(),
			},
		})
		assert.Nil(t, err)
	}
	log.Printf("number of messages in producer queue %d\n", kafkaContext.p.Len())
	kafkaContext.p.Flush(15 * 1000)
	log.Printf("AFTER number of messages in producer queue %d\n", kafkaContext.p.Len())

	kafkaContext.c.SubscribeTopics([]string{topic}, nil)

	for {
		msg, err := kafkaContext.c.ReadMessage(100 * time.Millisecond)
		if err != nil {
			log.Printf("No message: %v\n", err)
			continue
		}
		fmt.Printf("Consumed event from topic %s: key = %-10s value = %s\n",
			*msg.TopicPartition.Topic, string(msg.Key), string(msg.Value))
	}
}
