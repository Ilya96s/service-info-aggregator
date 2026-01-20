package kafka

import "context"

type EventHandler interface {
	Type() string
	Handle(ctx context.Context, key string, payload any) error
}
