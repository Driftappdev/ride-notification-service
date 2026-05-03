package domain

import (
	"errors"
	"time"
)

type Notification struct {
	ID             string
	TenantID       string
	RequestID      string
	IdempotencyKey string

	Channel  Channel
	Priority Priority
	Audience Audience

	Title string
	Body  string
	Image string
	Data  map[string]string

	CollapseKey string
	TTL         time.Duration

	TemplateKey string
	Locale      string

	ScheduledAt *time.Time
	CreatedAt   time.Time
}

func (n Notification) Validate() error {
	if n.ID == "" { return errors.New("notification id is required") }
	if n.TenantID == "" { return errors.New("tenant id is required") }
	if !n.Channel.Valid() { return errors.New("invalid channel") }
	if err := n.Audience.Validate(); err != nil { return err }
	if n.Title == "" && n.Body == "" && len(n.Data) == 0 {
		return errors.New("notification content is required")
	}
	if n.Priority == 0 {
		return errors.New("priority is required")
	}
	if n.TTL < 0 {
		return errors.New("ttl cannot be negative")
	}
	return nil
}
