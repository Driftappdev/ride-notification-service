package eventconsumerport

import "context"

type EventConsumerPort interface {
	Consume(ctx context.Context, topic string, message []byte) error
}
