package event

import (
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

// Publisher interface
type Publisher interface {
	Publish(topic string, key string, value []byte) error
	Close() error
}

// EventBusPublisher keeps legacy name but uses NATS transport.
type EventBusPublisher struct {
	conn *nats.Conn
}

func NewEventBusPublisher(brokers []string) *EventBusPublisher {
	url := "nats://localhost:4222"
	if len(brokers) > 0 {
		url = normalizeNATSURL(brokers[0])
	}
	nc, err := nats.Connect(url, nats.Name("notification-service-event-publisher"))
	if err != nil {
		panic(err)
	}
	return &EventBusPublisher{conn: nc}
}

func (p *EventBusPublisher) Publish(topic string, key string, value []byte) error {
	hdr := nats.Header{}
	if key != "" {
		hdr.Set("x-event-key", key)
	}
	err := p.conn.PublishMsg(&nats.Msg{Subject: topic, Header: hdr, Data: value})
	if err != nil {
		log.Printf("[NATSPublisher] Failed to publish message to subject %s, key=%s: %v", topic, key, err)
		return err
	}
	log.Printf("[NATSPublisher] Published message to subject %s, key=%s at=%s", topic, key, time.Now().UTC().Format(time.RFC3339))
	return nil
}

func (p *EventBusPublisher) Close() error {
	if p.conn != nil {
		p.conn.Flush()
		p.conn.Close()
	}
	return nil
}
