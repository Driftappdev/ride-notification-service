package event

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

// EventBusConsumer implements Client with per-topic NATS queue subscriptions.
type EventBusConsumer struct {
	url     string
	groupID string
	conn    *nats.Conn

	mu   sync.Mutex
	subs []*nats.Subscription
}

func NewEventBusConsumer(brokers []string, groupID string) *EventBusConsumer {
	url := "nats://localhost:4222"
	if len(brokers) > 0 {
		url = normalizeNATSURL(brokers[0])
	}
	nc, err := nats.Connect(url, nats.Name("notification-service-event-consumer"))
	if err != nil {
		panic(fmt.Sprintf("failed to connect nats: %v", err))
	}
	return &EventBusConsumer{url: url, groupID: groupID, conn: nc}
}

func (c *EventBusConsumer) Subscribe(ctx context.Context, topic string, handler func([]byte)) {
	queue := c.groupID
	if queue == "" {
		queue = "notification-service"
	}
	sub, err := c.conn.QueueSubscribe(topic, queue, func(m *nats.Msg) {
		handler(m.Data)
	})
	if err != nil {
		log.Printf("[NATSConsumer] subscribe failed topic=%s err=%v", topic, err)
		return
	}

	c.mu.Lock()
	c.subs = append(c.subs, sub)
	c.mu.Unlock()

	log.Printf("[NATSConsumer] Subscribing to subject: %s", topic)
	go func() {
		<-ctx.Done()
		_ = sub.Unsubscribe()
	}()
}

func (c *EventBusConsumer) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, sub := range c.subs {
		_ = sub.Unsubscribe()
	}
	c.subs = nil
	if c.conn != nil {
		c.conn.FlushTimeout(2 * time.Second)
		c.conn.Close()
	}
	return nil
}

func normalizeNATSURL(addr string) string {
	if addr == "" {
		return "nats://localhost:4222"
	}
	if len(addr) >= 7 && (addr[:7] == "nats://" || addr[:6] == "tls://") {
		return addr
	}
	return "nats://" + addr
}
