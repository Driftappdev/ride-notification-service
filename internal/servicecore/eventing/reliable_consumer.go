package eventing

import (
	"context"
	"fmt"
	"hash/fnv"
	"sync"
	"time"

	eventinfra "dift_backend_go/notification-service/internal/integration/event"
	"dift_backend_go/notification-service/internal/servicecore/queue"
)

type ReliableConsumer struct {
	raw   eventinfra.Client
	inbox *queue.Inbox
	dlq   *queue.DLQ
	mu    sync.Mutex
}

func NewReliableConsumer(raw eventinfra.Client) (*ReliableConsumer, error) {
	inbox := queue.NewInbox(queue.InboxConfig{
		DedupeEnabled:      true,
		DefaultMaxAttempts: 3,
		Log:                newQueueLogger("notification-service.consumer"),
	})
	dlq, err := queue.NewDLQ(queue.DLQConfig{
		Storage:         queue.NewInMemoryDLQStorage(),
		Log:             newQueueLogger("notification-service.consumer.dlq"),
		DefaultTTL:      24 * time.Hour,
		PoisonThreshold: 3,
		ReplayInbox:     inbox,
	})
	if err != nil {
		return nil, err
	}
	return &ReliableConsumer{raw: raw, inbox: inbox, dlq: dlq}, nil
}

func (c *ReliableConsumer) Subscribe(ctx context.Context, topic string, handler func([]byte)) {
	c.raw.Subscribe(ctx, topic, func(payload []byte) {
		msg := &queue.Message{
			ID:          buildMessageID(topic, payload),
			SenderID:    "event-bus",
			Topic:       topic,
			Priority:    queue.PriorityHigh,
			Body:        append([]byte(nil), payload...),
			MaxAttempts: 3,
			Timestamp:   time.Now().UTC(),
		}
		if err := c.inbox.Submit(msg); err != nil {
			return
		}
		_ = c.process(ctx, handler)
	})
}

func (c *ReliableConsumer) process(ctx context.Context, handler func([]byte)) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	msg, err := c.inbox.Consume(ctx)
	if err != nil {
		return err
	}
	msg.Attempts++
	done := make(chan struct{})
	go func() {
		defer close(done)
		handler(msg.Body)
	}()
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		_ = c.dlq.Send(context.Background(), msg, ctx.Err())
		return ctx.Err()
	}
}

func (c *ReliableConsumer) Close() error {
	if c.inbox != nil {
		_ = c.inbox.Close()
	}
	if c.dlq != nil {
		c.dlq.Close()
	}
	if closer, ok := c.raw.(interface{ Close() error }); ok {
		return closer.Close()
	}
	return nil
}

func buildMessageID(topic string, payload []byte) string {
	h := fnv.New64a()
	_, _ = h.Write([]byte(topic))
	_, _ = h.Write(payload)
	return fmt.Sprintf("%x", h.Sum64())
}
