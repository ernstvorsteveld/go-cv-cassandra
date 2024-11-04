package out

import "context"

type EventType string

const (
	EventTypeTagCreated EventType = "tag.created"
	EventTypeTagUpdated EventType = "tag.updated"
	EventTypeTagDeleted EventType = "tag.deleted"

	EventTypeExperienceCreated EventType = "experience.created"
	EventTypeExperienceUpdated EventType = "experience.updated"
	EventTypeExperienceDeleted EventType = "experience.deleted"
)

type EventPort interface {
	Publish(ctx context.Context, eventType string, event interface{}) error
}
