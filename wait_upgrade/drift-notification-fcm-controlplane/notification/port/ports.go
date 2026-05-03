package port

import (
	"context"
	"time"

	"github.com/driftappdev/libpackage/notification/domain"
)

type Sender interface {
	Send(ctx context.Context, n domain.Notification) (domain.DeliveryResult, error)
}

type TokenResolver interface {
	ResolveUserTokens(ctx context.Context, tenantID, userID string) ([]string, error)
	ResolveSegmentTokens(ctx context.Context, tenantID, segmentID string) ([]string, error)
}

type IdempotencyStore interface {
	Seen(ctx context.Context, tenantID, key string) (bool, error)
	Mark(ctx context.Context, tenantID, key string, ttl time.Duration) error
}

type RateLimiter interface {
	Allow(ctx context.Context, tenantID string, key string, cost int) (bool, time.Duration, error)
}

type DeliveryStore interface {
	SaveResult(ctx context.Context, result domain.DeliveryResult) error
}

type PolicyEngine interface {
	Evaluate(ctx context.Context, n domain.Notification) (PolicyDecision, error)
}

type PolicyDecision struct {
	Allow  bool
	Reason string
	Mutate func(*domain.Notification)
}

type Router interface {
	Route(ctx context.Context, n domain.Notification) (Sender, error)
}
