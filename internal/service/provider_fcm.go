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

type FCMProvider struct {
	cfg        *config.Config
	httpClient *http.Client
}

func NewFCMProvider(cfg *config.Config) *FCMProvider {
	timeout := time.Duration(cfg.FCMTimeoutMS) * time.Millisecond
	if timeout <= 0 {
		timeout = 4 * time.Second
	}
	return &FCMProvider{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (p *FCMProvider) Name() string { return "fcm" }

func (p *FCMProvider) Supports(channel string) bool {
	return channel == "fcm" || channel == "push" || channel == "firebase"
}

func (p *FCMProvider) Send(ctx context.Context, env *NotificationEnvelope) error {
	if p.cfg.FCMServerKey == "" {
		return errors.New("fcm server_key is empty")
	}
	if len(env.Recipients.DeviceTokens) == 0 && env.Recipients.Topic == "" {
		return errors.New("no fcm target token/topic")
	}

	reqBody := map[string]any{
		"priority": env.Priority,
		"notification": map[string]any{
			"title": env.Title,
			"body":  env.Message,
		},
		"data": env.Data,
	}
	if len(env.Recipients.DeviceTokens) > 0 {
		reqBody["registration_ids"] = env.Recipients.DeviceTokens
	}
	if env.Recipients.Topic != "" {
		reqBody["to"] = "/topics/" + env.Recipients.Topic
	}
	raw, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.cfg.FCMEndpoint, bytes.NewReader(raw))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "key="+p.cfg.FCMServerKey)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return errors.New("fcm responded with non-2xx")
	}
	return nil
}
