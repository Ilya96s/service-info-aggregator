package kafka

import (
	"context"
	"fmt"
)

type EventRouter struct {
	handlers map[string]EventHandler
}

func NewEventRouter(handlers ...EventHandler) *EventRouter {
	m := make(map[string]EventHandler)
	for _, h := range handlers {
		m[h.Type()] = h
	}
	return &EventRouter{handlers: m}
}

func (r *EventRouter) Route(cxt context.Context, eventType string, key string, payload any) error {
	h, ok := r.handlers[eventType]
	if !ok {
		return fmt.Errorf("no handler for event type: %s", eventType)
	}

	return h.Handle(cxt, key, payload)
}
