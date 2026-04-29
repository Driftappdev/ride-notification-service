package service

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type NotificationEnvelope struct {
	NotificationID string                 `json:"notification_id"`
	EventType      string                 `json:"event_type"`
	Status         string                 `json:"status"`
	Title          string                 `json:"title"`
	Message        string                 `json:"message"`
	Priority       string                 `json:"priority"`
	Channels       []string               `json:"channels"`
	Recipients     NotificationRecipients `json:"recipients"`
	Data           map[string]any         `json:"data"`
	Metadata       map[string]any         `json:"metadata"`
	CreatedAt      string                 `json:"created_at"`
}

type NotificationRecipients struct {
	UserIDs      []string `json:"user_ids"`
	DriverIDs    []string `json:"driver_ids"`
	DeviceTokens []string `json:"device_tokens"`
	Topic        string   `json:"topic"`
}

func toEnvelope(payload map[string]any) (*NotificationEnvelope, error) {
	raw, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	var env NotificationEnvelope
	if err := json.Unmarshal(raw, &env); err != nil {
		return nil, err
	}

	if env.NotificationID == "" {
		env.NotificationID = time.Now().UTC().Format("20060102T150405.000000000")
	}
	if env.Priority == "" {
		env.Priority = "normal"
	}
	if env.CreatedAt == "" {
		env.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	}
	if env.Data == nil {
		env.Data = map[string]any{}
	}
	if env.Metadata == nil {
		env.Metadata = map[string]any{}
	}
	applyDefaultTemplate(&env)
	return &env, nil
}

func applyDefaultTemplate(env *NotificationEnvelope) {
	if strings.TrimSpace(env.Title) != "" && strings.TrimSpace(env.Message) != "" {
		return
	}
	key := strings.ToLower(strings.TrimSpace(env.EventType + ":" + env.Status))
	switch key {
	case "order:accepted", "order_status:accepted":
		env.Title = defaultTitle(env.Title, "Order Accepted")
		env.Message = defaultMsg(env.Message, "Your order has been accepted.")
	case "order:completed", "order_status:completed":
		env.Title = defaultTitle(env.Title, "Order Completed")
		env.Message = defaultMsg(env.Message, "Your trip/order is completed successfully.")
	case "order:cancelled", "order_status:cancelled":
		env.Title = defaultTitle(env.Title, "Order Cancelled")
		env.Message = defaultMsg(env.Message, "Your order was cancelled.")
	case "driver:online", "driver_status:online":
		env.Title = defaultTitle(env.Title, "Driver Online")
		env.Message = defaultMsg(env.Message, "Driver is now available.")
	case "driver:offline", "driver_status:offline":
		env.Title = defaultTitle(env.Title, "Driver Offline")
		env.Message = defaultMsg(env.Message, "Driver is currently unavailable.")
	case "payment:success", "payment_status:success":
		env.Title = defaultTitle(env.Title, "Payment Success")
		env.Message = defaultMsg(env.Message, "Payment has been completed.")
	case "payment:failed", "payment_status:failed":
		env.Title = defaultTitle(env.Title, "Payment Failed")
		env.Message = defaultMsg(env.Message, "Payment failed, please try again.")
	default:
		env.Title = defaultTitle(env.Title, "Notification")
		if strings.TrimSpace(env.Message) == "" {
			env.Message = fmt.Sprintf("Event '%s' status '%s'", env.EventType, env.Status)
		}
	}
}

func defaultTitle(cur, fallback string) string {
	if strings.TrimSpace(cur) != "" {
		return cur
	}
	return fallback
}

func defaultMsg(cur, fallback string) string {
	if strings.TrimSpace(cur) != "" {
		return cur
	}
	return fallback
}
