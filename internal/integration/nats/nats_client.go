package natsintegration

import "context"

type NATSPublisherPort interface {
	Publish(ctx context.Context, subject string, data []byte) error
}
