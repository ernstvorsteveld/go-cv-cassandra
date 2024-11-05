package event_model

type EventPayload struct {
	CorrelationId string      `json:"correlationId"`
	EventType     string      `json:"type"`
	Key           string      `json:"key"`
	Payload       interface{} `json:"payload"`
}
