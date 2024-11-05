package kafka_producer

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	event_model "github.com/ernstvorsteveld/go-cv-cassandra/adapter/domain/event"
	adapter_in_event "github.com/ernstvorsteveld/go-cv-cassandra/adapter/in/event"
	kafka_consumer "github.com/ernstvorsteveld/go-cv-cassandra/adapter/in/event/kafka"
	adapter_out_event "github.com/ernstvorsteveld/go-cv-cassandra/adapter/out/event"
	"github.com/ernstvorsteveld/go-cv-cassandra/domain/model"
	"github.com/ernstvorsteveld/go-cv-cassandra/pkg/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"

	log "github.com/sirupsen/logrus"
	"github.com/testcontainers/testcontainers-go/modules/kafka"
)

var producer adapter_out_event.Producer
var consumer adapter_in_event.Consumer

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
	producer = NewKafkaProducer(config)
	consumer = kafka_consumer.NewKafkaConsumer(config)
	defer producer.Close()
	defer consumer.Close()

	m.Run()
}

func Test_should_publish_to_topic(t *testing.T) {
	log.Println("About to publish and consume")
	topic := "test.tag.created"

	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("key_%d", i)
		err := producer.Publish(context.Background(), topic, event_model.EventPayload{
			CorrelationId: uuid.NewString(),
			EventType:     "test",
			Key:           key,
			Payload: model.Tag{
				Id:   key,
				Name: fmt.Sprintf("tagname_%d", i),
			},
		})
		assert.Nil(t, err)
	}

	producer.Flush(1500 * 1000)
	assert.Equal(t, 0, producer.Len(), "Not all messages flushed")

	consumer.SubscribeTopics([]string{topic})

	nr := 0
	others := true
	for others {
		msg, err := consumer.ReadMessage(100 * time.Millisecond)
		if err != nil {
			log.Printf("No message: %v\n", err)
			continue
		}

		var event event_model.EventPayload
		json.Unmarshal(msg.Value, &event)
		assert.Equal(t, fmt.Sprintf("key_%d", nr), event.Key)
		fmt.Printf("Consumed event from topic %s: key = %-10s value = %s\n",
			*msg.TopicPartition.Topic, string(msg.Key), string(msg.Value))
		nr = nr + 1
		others = nr != 10
	}
}
