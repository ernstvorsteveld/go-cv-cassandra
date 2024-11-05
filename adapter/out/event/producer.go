package adapter_out_event

import (
	"context"
	"time"

	model "github.com/ernstvorsteveld/go-cv-cassandra/adapter/domain/event"
)

type Producer interface {
	Publish(ctx context.Context, eventType string, event model.EventPayload) error
	Close() error
	Flush(duration time.Duration) int
	Len() int
}
