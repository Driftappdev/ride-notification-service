package service

import (
	"context"
	"encoding/json"
	"strings"
	"time"
)

type NotificationProvider interface {
	Name() string
	Supports(channel string) bool
	Send(ctx context.Context, env *NotificationEnvelope) error
}

type DeliveryResult struct {
	NotificationID string `json:"notification_id"`
	Provider       string `json:"provider"`
	Channel        string `json:"channel"`
	Success        bool   `json:"success"`
	Error          string `json:"error,omitempty"`
	At             string `json:"at"`
}

func uniqueChannels(channels []string) []string {
	seen := map[string]bool{}
	out := make([]string, 0, len(channels))
	for _, ch := range channels {
		ch = strings.ToLower(strings.TrimSpace(ch))
		if ch == "" || seen[ch] {
			continue
		}
		seen[ch] = true
		out = append(out, ch)
	}
	return out
}

func marshalEnvelope(env *NotificationEnvelope) ([]byte, error) {
	return json.Marshal(env)
}

func nowRFC3339() string {
	return time.Now().UTC().Format(time.RFC3339)
}
