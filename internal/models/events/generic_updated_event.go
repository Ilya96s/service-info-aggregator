package events

import "time"

type GenericUpdatedEvent struct {
	Type      string    `json:"type"`
	Key       string    `json:"key"`
	Payload   any       `json:"payload"`
	Timestamp time.Time `json:"timestamp"`
}
