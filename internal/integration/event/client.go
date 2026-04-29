package event

import "context"

// Client defines the interface that NotificationService uses to subscribe to events
type Client interface {
	Subscribe(ctx context.Context, topic string, handler func([]byte))
}
