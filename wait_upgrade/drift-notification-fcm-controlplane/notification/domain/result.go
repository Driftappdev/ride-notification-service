package domain

import "time"

type DeliveryStatus string

const (
	DeliveryAccepted   DeliveryStatus = "accepted"
	DeliverySent       DeliveryStatus = "sent"
	DeliveryPartial    DeliveryStatus = "partial"
	DeliveryFailed     DeliveryStatus = "failed"
	DeliverySuppressed DeliveryStatus = "suppressed"
	DeliveryRetried    DeliveryStatus = "retried"
)

type TokenResult struct {
	Token        string
	MessageID    string
	OK           bool
	ErrorCode    string
	ErrorMessage string
	Retryable    bool
	RemoveToken  bool
}

type DeliveryResult struct {
	NotificationID string
	Provider       string
	Status         DeliveryStatus
	MessageID      string
	AcceptedAt    time.Time
	SentAt         time.Time
	TokenResults   []TokenResult
	ErrorCode      string
	ErrorMessage   string
	Retryable      bool
	Attempt        int
}
