package eventing

import (
	"context"
	"fmt"
	"time"

	eventinfra "dift_backend_go/notification-service/internal/integration/event"
	"dift_backend_go/notification-service/internal/servicecore/queue"
)

type ReliablePublisher struct {
	raw    eventinfra.Publisher
	outbox *queue.Outbox
}

func NewReliablePublisher(raw eventinfra.Publisher) (*ReliablePublisher, error) {
	if raw == nil {
		return nil, fmt.Errorf("raw publisher is required")
	}

	outbox, err := queue.NewOutbox(queue.OutboxConfig{
		Relay: func(ctx context.Context, batch []*queue.Message) error {
			for _, msg := range batch {
				key := ""
				if msg.Metadata != nil {
					key = msg.Metadata["event_key"]
				}
				if err := raw.Publish(msg.Topic, key, msg.Body); err != nil {
					return err
				}
			}
			return nil
		},
		Log:                 newQueueLogger("notification-service.publisher"),
		ReloadPending:       false,
		WorkerCount:         2,
		MaxBatchSize:        32,
		BatchLingerDuration: 10 * time.Millisecond,
		MaxRelayAttempts:    5,
	})
	if err != nil {
		return nil, err
	}
	return &ReliablePublisher{raw: raw, outbox: outbox}, nil
}

func (p *ReliablePublisher) Publish(topic string, key string, value []byte) error {
	return p.outbox.Enqueue(context.Background(), &queue.Message{
		ID:        fmt.Sprintf("%d-%s-%s", time.Now().UnixNano(), topic, key),
		SenderID:  "notification-service",
		Topic:     topic,
		Priority:  queue.PriorityHigh,
		Body:      append([]byte(nil), value...),
		Timestamp: time.Now().UTC(),
		Metadata: map[string]string{
			"event_key": key,
		},
	})
}

func (p *ReliablePublisher) Close() error {
	if p.outbox != nil {
		_ = p.outbox.Close()
	}
	if p.raw != nil {
		return p.raw.Close()
	}
	return nil
}
