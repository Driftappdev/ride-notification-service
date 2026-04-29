package service

import (
	"context"

	"dift_backend_go/notification-service/config"
	mqtt "dift_backend_go/notification-service/internal/integration/mqtt_client"
)

type MQTTProvider struct {
	cfg    *config.Config
	client mqtt.Client
}

func NewMQTTProvider(cfg *config.Config, client mqtt.Client) *MQTTProvider {
	return &MQTTProvider{cfg: cfg, client: client}
}

func (p *MQTTProvider) Name() string { return "mqtt" }

func (p *MQTTProvider) Supports(channel string) bool {
	return channel == "mqtt" || channel == "push_mqtt"
}

func (p *MQTTProvider) Send(ctx context.Context, env *NotificationEnvelope) error {
	data, err := marshalEnvelope(env)
	if err != nil {
		return err
	}
	topic := p.cfg.MQTTTopicGeneric
	if len(env.Recipients.DriverIDs) > 0 {
		topic = p.cfg.MQTTTopicDriver
	} else if len(env.Recipients.UserIDs) > 0 {
		topic = p.cfg.MQTTTopicPassenger
	}
	return p.client.SendRaw(ctx, topic, data)
}
