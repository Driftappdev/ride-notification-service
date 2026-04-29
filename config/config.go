package config

import (
	"log"

	"github.com/spf13/viper"
)

// Config holds application configuration for notification services.
type Config struct {
	// Server
	HTTPPort int
	GRPCPort int
	Timeout  int // ms

	// EventBus / EventBus
	EventBusBrokers []string
	EventBusGroup   string

	// Reserved topics (kept separate from ride flow)
	EventBusMatchingTopic     string
	EventBusDriverCancelTopic string
	EventBusOrderTopic        string
	EventBusDriverTopic       string

	// Generic notification topics
	EventBusGenericTopic  string
	EventBusDispatchTopic string
	EventBusStatusTopic   string
	EventBusEventTopics   []string

	// MQTT settings
	MQTTBroker         string
	MQTTTopicDriver    string
	MQTTTopicPassenger string
	MQTTTopicGeneric   string
	MQTTSecure         bool

	// Provider settings
	DefaultChannels       []string
	ProviderEventBusEnable   bool
	ProviderMQTTEnable    bool
	ProviderFCMEnable     bool
	ProviderWebhookEnable bool

	FCMEndpoint  string
	FCMServerKey string
	FCMTimeoutMS int

	WebhookURL       string
	WebhookAuthToken string
	WebhookTimeoutMS int

	// Redis / Cache
	RedisURL                string
	RedisPassword           string
	RedisDB                 int
	RedisTTL                string
	RedisNotificationPrefix string
}

// LoadConfig loads configuration from config.yaml and environment variables.
func LoadConfig() *Config {
	v := viper.New()
	v.SetConfigFile("./config/config.yaml")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		log.Printf("[Config] warning reading config file, using defaults/env: %v", err)
	}

	cfg := &Config{
		HTTPPort: v.GetInt("server.http_port"),
		GRPCPort: v.GetInt("server.grpc_port"),
		Timeout:  v.GetInt("server.timeout_ms"),

		EventBusBrokers:           v.GetStringSlice("event_bus.brokers"),
		EventBusGroup:             v.GetString("event_bus.group"),
		EventBusMatchingTopic:     v.GetString("event_bus.matching_topic"),
		EventBusDriverCancelTopic: v.GetString("event_bus.cancel_topic"),
		EventBusOrderTopic:        v.GetString("event_bus.order_topic"),
		EventBusDriverTopic:       v.GetString("event_bus.driver_topic"),
		EventBusGenericTopic:      v.GetString("event_bus.generic_topic"),
		EventBusDispatchTopic:     v.GetString("event_bus.dispatch_topic"),
		EventBusStatusTopic:       v.GetString("event_bus.status_topic"),
		EventBusEventTopics:       v.GetStringSlice("event_bus.event_topics"),

		MQTTBroker:         v.GetString("mqtt.broker"),
		MQTTTopicDriver:    v.GetString("mqtt.driver_topic"),
		MQTTTopicPassenger: v.GetString("mqtt.passenger_topic"),
		MQTTTopicGeneric:   v.GetString("mqtt.generic_topic"),
		MQTTSecure:         v.GetBool("mqtt.secure"),

		DefaultChannels:       v.GetStringSlice("providers.default_channels"),
		ProviderEventBusEnable:   v.GetBool("providers.event_bus.enabled"),
		ProviderMQTTEnable:    v.GetBool("providers.mqtt.enabled"),
		ProviderFCMEnable:     v.GetBool("providers.fcm.enabled"),
		ProviderWebhookEnable: v.GetBool("providers.webhook.enabled"),

		FCMEndpoint:  v.GetString("providers.fcm.endpoint"),
		FCMServerKey: v.GetString("providers.fcm.server_key"),
		FCMTimeoutMS: v.GetInt("providers.fcm.timeout_ms"),

		WebhookURL:       v.GetString("providers.webhook.url"),
		WebhookAuthToken: v.GetString("providers.webhook.auth_token"),
		WebhookTimeoutMS: v.GetInt("providers.webhook.timeout_ms"),

		RedisURL:                v.GetString("redis.url"),
		RedisPassword:           v.GetString("redis.password"),
		RedisDB:                 v.GetInt("redis.db"),
		RedisTTL:                v.GetString("redis.ttl"),
		RedisNotificationPrefix: v.GetString("redis.keys.notification_prefix"),
	}

	if cfg.HTTPPort == 0 {
		cfg.HTTPPort = 2222
	}
	if cfg.GRPCPort == 0 {
		cfg.GRPCPort = 5051
	}
	if cfg.EventBusGroup == "" {
		cfg.EventBusGroup = "notification-service-group"
	}
	if cfg.EventBusMatchingTopic == "" {
		cfg.EventBusMatchingTopic = "notification.generic.matching-events"
	}
	if cfg.EventBusDriverCancelTopic == "" {
		cfg.EventBusDriverCancelTopic = "notification.generic.cancel-events"
	}
	if cfg.EventBusOrderTopic == "" {
		cfg.EventBusOrderTopic = "notification.generic.order-events"
	}
	if cfg.EventBusDriverTopic == "" {
		cfg.EventBusDriverTopic = "notification.generic.driver-events"
	}
	if cfg.EventBusGenericTopic == "" {
		cfg.EventBusGenericTopic = "admin.notification.send"
	}
	if cfg.EventBusDispatchTopic == "" {
		cfg.EventBusDispatchTopic = "notification.generic.dispatch"
	}
	if cfg.EventBusStatusTopic == "" {
		cfg.EventBusStatusTopic = "notification.delivery.status"
	}
	if cfg.MQTTBroker == "" {
		cfg.MQTTBroker = "tcp://mqtt-broker:1883"
	}
	if cfg.MQTTTopicDriver == "" {
		cfg.MQTTTopicDriver = "ui/driver/general-notification"
	}
	if cfg.MQTTTopicPassenger == "" {
		cfg.MQTTTopicPassenger = "ui/passenger/general-notification"
	}
	if cfg.MQTTTopicGeneric == "" {
		cfg.MQTTTopicGeneric = "ui/general/notification"
	}
	if len(cfg.DefaultChannels) == 0 {
		cfg.DefaultChannels = []string{"event_bus", "mqtt"}
	}
	if !cfg.ProviderEventBusEnable && !cfg.ProviderMQTTEnable && !cfg.ProviderFCMEnable && !cfg.ProviderWebhookEnable {
		cfg.ProviderEventBusEnable = true
		cfg.ProviderMQTTEnable = true
	}
	if cfg.FCMEndpoint == "" {
		cfg.FCMEndpoint = "https://fcm.googleapis.com/fcm/send"
	}
	if cfg.FCMTimeoutMS <= 0 {
		cfg.FCMTimeoutMS = 4000
	}
	if cfg.WebhookTimeoutMS <= 0 {
		cfg.WebhookTimeoutMS = 4000
	}
	if cfg.RedisTTL == "" {
		cfg.RedisTTL = "30s"
	}
	if cfg.RedisNotificationPrefix == "" {
		cfg.RedisNotificationPrefix = "notification:"
	}

	log.Printf("[Config] loaded notification-service config on http_port=%d grpc_port=%d", cfg.HTTPPort, cfg.GRPCPort)
	return cfg
}
