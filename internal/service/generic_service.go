package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"dift_backend_go/notification-service/config"
	"dift_backend_go/notification-service/internal/integration/event"
	mqtt "dift_backend_go/notification-service/internal/integration/mqtt_client"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

type GenericNotificationService struct {
	cfg           *config.Config
	eventClient   event.Client
	mqttClient    mqtt.Client
	eventProducer event.Publisher
	providers     []NotificationProvider
}

func NewGenericNotificationService(
	cfg *config.Config,
	eventClient event.Client,
	mqttClient mqtt.Client,
	eventProducer event.Publisher,
) *GenericNotificationService {
	svc := &GenericNotificationService{
		cfg:           cfg,
		eventClient:   eventClient,
		mqttClient:    mqttClient,
		eventProducer: eventProducer,
	}
	svc.providers = svc.buildProviders()
	return svc
}

func (s *GenericNotificationService) StartListening(ctx context.Context) {
	topics := make([]string, 0, 1+len(s.cfg.EventBusEventTopics))
	if s.cfg.EventBusGenericTopic != "" {
		topics = append(topics, s.cfg.EventBusGenericTopic)
	}
	topics = append(topics, s.cfg.EventBusEventTopics...)
	topics = uniqueChannels(topics)

	if len(topics) == 0 {
		log.Println("[GenericNotificationService] no topics configured, skip consumer")
		return
	}
	for _, topic := range topics {
		log.Printf("[GenericNotificationService] listening topic=%s", topic)
		s.eventClient.Subscribe(ctx, topic, s.handleInboundNotification)
	}
}

func (s *GenericNotificationService) Dispatch(ctx context.Context, payload map[string]any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return s.dispatchBytes(ctx, data)
}

func (s *GenericNotificationService) handleInboundNotification(raw []byte) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.dispatchBytes(ctx, raw); err != nil {
		log.Printf("[GenericNotificationService] dispatch failed: %v", err)
	}
}

func (s *GenericNotificationService) dispatchBytes(ctx context.Context, raw []byte) error {
	payload, err := normalizePayload(raw)
	if err != nil {
		return err
	}

	env, err := toEnvelope(payload)
	if err != nil {
		return err
	}

	channels := uniqueChannels(env.Channels)
	if len(channels) == 0 {
		channels = uniqueChannels(s.cfg.DefaultChannels)
	}

	successCount := 0
	results := make([]DeliveryResult, 0, len(channels))
	for _, channel := range channels {
		provider := s.findProvider(channel)
		if provider == nil {
			results = append(results, DeliveryResult{
				NotificationID: env.NotificationID,
				Provider:       "none",
				Channel:        channel,
				Success:        false,
				Error:          "unsupported_channel",
				At:             nowRFC3339(),
			})
			continue
		}

		sendErr := provider.Send(ctx, env)
		result := DeliveryResult{
			NotificationID: env.NotificationID,
			Provider:       provider.Name(),
			Channel:        channel,
			Success:        sendErr == nil,
			At:             nowRFC3339(),
		}
		if sendErr != nil {
			result.Error = sendErr.Error()
			log.Printf("[GenericNotificationService] provider=%s channel=%s failed: %v", provider.Name(), channel, sendErr)
		} else {
			successCount++
		}
		results = append(results, result)
	}

	s.publishDeliveryStatus(env.NotificationID, results)
	if successCount == 0 {
		return errors.New("notification delivery failed for all channels")
	}
	return nil
}

func normalizePayload(raw []byte) (map[string]any, error) {
	out := map[string]any{}

	pbStruct := &structpb.Struct{}
	if err := proto.Unmarshal(raw, pbStruct); err == nil {
		return pbStruct.AsMap(), nil
	}
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *GenericNotificationService) findProvider(channel string) NotificationProvider {
	for _, p := range s.providers {
		if p.Supports(channel) {
			return p
		}
	}
	return nil
}

func (s *GenericNotificationService) buildProviders() []NotificationProvider {
	out := make([]NotificationProvider, 0, 4)
	if s.cfg.ProviderEventBusEnable {
		out = append(out, NewEventBusDispatchProvider(s.cfg, s.eventProducer))
	}
	if s.cfg.ProviderMQTTEnable {
		out = append(out, NewMQTTProvider(s.cfg, s.mqttClient))
	}
	if s.cfg.ProviderFCMEnable {
		out = append(out, NewFCMProvider(s.cfg))
	}
	if s.cfg.ProviderWebhookEnable {
		out = append(out, NewWebhookProvider(s.cfg))
	}
	log.Printf("[GenericNotificationService] providers initialized: %d", len(out))
	return out
}

func (s *GenericNotificationService) publishDeliveryStatus(notificationID string, results []DeliveryResult) {
	if s.cfg.EventBusStatusTopic == "" || len(results) == 0 {
		return
	}
	raw, err := json.Marshal(map[string]any{
		"notification_id": notificationID,
		"results":         results,
		"at":              nowRFC3339(),
	})
	if err != nil {
		return
	}
	if err := s.eventProducer.Publish(s.cfg.EventBusStatusTopic, buildStatusKey(notificationID), raw); err != nil {
		log.Printf("[GenericNotificationService] publish delivery status failed: %v", err)
	}
}

func buildStatusKey(notificationID string) string {
	if strings.TrimSpace(notificationID) != "" {
		return notificationID
	}
	return fmt.Sprintf("status-%d", time.Now().UTC().UnixNano())
}
