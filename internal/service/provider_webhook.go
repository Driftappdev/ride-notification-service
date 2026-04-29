package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"dift_backend_go/notification-service/config"
)

type WebhookProvider struct {
	cfg        *config.Config
	httpClient *http.Client
}

func NewWebhookProvider(cfg *config.Config) *WebhookProvider {
	timeout := time.Duration(cfg.WebhookTimeoutMS) * time.Millisecond
	if timeout <= 0 {
		timeout = 4 * time.Second
	}
	return &WebhookProvider{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (p *WebhookProvider) Name() string { return "webhook" }

func (p *WebhookProvider) Supports(channel string) bool {
	return channel == "webhook"
}

func (p *WebhookProvider) Send(ctx context.Context, env *NotificationEnvelope) error {
	if p.cfg.WebhookURL == "" {
		return errors.New("webhook url is empty")
	}
	body, err := marshalEnvelope(env)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.cfg.WebhookURL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if p.cfg.WebhookAuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+p.cfg.WebhookAuthToken)
	}
	resp, err := p.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return errors.New("webhook responded with non-2xx")
	}

	statusEvent := map[string]any{
		"notification_id": env.NotificationID,
		"provider":        "webhook",
		"status":          "delivered",
		"at":              nowRFC3339(),
	}
	_, _ = json.Marshal(statusEvent)
	return nil
}
