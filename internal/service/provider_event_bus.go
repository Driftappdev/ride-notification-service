package service

import (
	"context"

	"dift_backend_go/notification-service/config"
	"dift_backend_go/notification-service/internal/integration/event"
)

type EventBusDispatchProvider struct {
	cfg      *config.Config
	producer event.Publisher
}

func NewEventBusDispatchProvider(cfg *config.Config, producer event.Publisher) *EventBusDispatchProvider {
	return &EventBusDispatchProvider{cfg: cfg, producer: producer}
}

func (p *EventBusDispatchProvider) Name() string { return "event_bus" }

func (p *EventBusDispatchProvider) Supports(channel string) bool {
	return channel == "event_bus" || channel == "event_bus_dispatch"
}

func (p *EventBusDispatchProvider) Send(_ context.Context, env *NotificationEnvelope) error {
	data, err := marshalEnvelope(env)
	if err != nil {
		return err
	}
	return p.producer.Publish(p.cfg.EventBusDispatchTopic, env.NotificationID, data)
}
