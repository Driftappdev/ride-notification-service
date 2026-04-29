package event_busintegration

import "context"

type EventBusPublisherPort interface {
	Publish(ctx context.Context, topic string, key, value []byte) error
}
