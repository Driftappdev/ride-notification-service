package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
)

type AdminControlRequest struct {
	Action         string         `json:"action"`
	NotificationID string         `json:"notification_id"`
	TemplateID     string         `json:"template_id"`
	Channels       []string       `json:"channels"`
	Payload        map[string]any `json:"payload"`
	Recipients     map[string]any `json:"recipients"`
}

type AdminControlService struct {
	dispatch     func(context.Context, map[string]any) error
	timeout      time.Duration
	mu           sync.Mutex
	idempotentOK map[string]time.Time
}

func NewAdminControlService(dispatch func(context.Context, map[string]any) error, timeoutMS int) *AdminControlService {
	timeout := 5 * time.Second
	if timeoutMS > 0 {
		timeout = time.Duration(timeoutMS) * time.Millisecond
	}
	return &AdminControlService{
		dispatch:     dispatch,
		timeout:      timeout,
		idempotentOK: make(map[string]time.Time),
	}
}

func (s *AdminControlService) Execute(ctx context.Context, idempotencyKey string, req AdminControlRequest) (map[string]any, error) {
	if strings.TrimSpace(req.Action) == "" {
		return nil, errors.New("action is required")
	}
	if idempotencyKey != "" && s.seen(idempotencyKey) {
		return map[string]any{
			"accepted":       true,
			"idempotent_hit": true,
			"action":         req.Action,
		}, nil
	}

	switch req.Action {
	case "health.check":
		return map[string]any{"service": "notification-service", "status": "ok", "accepted": true}, nil
	case "notification.dispatch", "notification.send", "notification.preview":
		payload := s.toDispatchPayload(req)
		runCtx, cancel := context.WithTimeout(ctx, s.timeout)
		defer cancel()
		if err := s.dispatch(runCtx, payload); err != nil {
			return nil, err
		}
		if idempotencyKey != "" {
			s.mark(idempotencyKey)
		}
		return map[string]any{
			"accepted":        true,
			"action":          req.Action,
			"notification_id": payload["notification_id"],
		}, nil
	default:
		return nil, fmt.Errorf("unsupported action: %s", req.Action)
	}
}

func (s *AdminControlService) toDispatchPayload(req AdminControlRequest) map[string]any {
	out := map[string]any{
		"notification_id": req.NotificationID,
		"template_id":     req.TemplateID,
		"channels":        req.Channels,
		"payload":         req.Payload,
		"recipients":      req.Recipients,
	}
	if strings.TrimSpace(req.NotificationID) == "" {
		out["notification_id"] = fmt.Sprintf("admin-%d", time.Now().UTC().UnixNano())
	}
	return out
}

func (s *AdminControlService) seen(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	exp, ok := s.idempotentOK[key]
	if !ok {
		return false
	}
	if time.Now().After(exp) {
		delete(s.idempotentOK, key)
		return false
	}
	return true
}

func (s *AdminControlService) mark(key string) {
	s.mu.Lock()
	s.idempotentOK[key] = time.Now().Add(10 * time.Minute)
	s.mu.Unlock()
}
