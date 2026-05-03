package controlplane

import (
	"context"
	"errors"
	"time"

	"github.com/driftappdev/libpackage/notification/domain"
	"github.com/driftappdev/libpackage/notification/port"
)

type ControlPlane struct {
	cfg         Config
	router      port.Router
	policy      port.PolicyEngine
	idem        port.IdempotencyStore
	limiter     port.RateLimiter
	store       port.DeliveryStore
	resolver    port.TokenResolver
}

type Option func(*ControlPlane)

func WithPolicy(p port.PolicyEngine) Option { return func(c *ControlPlane) { c.policy = p } }
func WithIdempotency(s port.IdempotencyStore) Option { return func(c *ControlPlane) { c.idem = s } }
func WithRateLimiter(l port.RateLimiter) Option { return func(c *ControlPlane) { c.limiter = l } }
func WithDeliveryStore(s port.DeliveryStore) Option { return func(c *ControlPlane) { c.store = s } }
func WithTokenResolver(r port.TokenResolver) Option { return func(c *ControlPlane) { c.resolver = r } }

func New(cfg Config, router port.Router, opts ...Option) *ControlPlane {
	cp := &ControlPlane{cfg: cfg.Normalize(), router: router}
	for _, opt := range opts { opt(cp) }
	return cp
}

func (c *ControlPlane) Dispatch(ctx context.Context, n domain.Notification) (domain.DeliveryResult, error) {
	if n.CreatedAt.IsZero() { n.CreatedAt = time.Now().UTC() }
	if err := n.Validate(); err != nil {
		return domain.DeliveryResult{NotificationID: n.ID, Status: domain.DeliveryFailed, ErrorCode: "VALIDATION_ERROR", ErrorMessage: err.Error()}, err
	}

	if n.IdempotencyKey != "" && c.idem != nil {
		seen, err := c.idem.Seen(ctx, n.TenantID, n.IdempotencyKey)
		if err != nil { return failedCP(n, "IDEMPOTENCY_CHECK_FAILED", err), err }
		if seen {
			res := domain.DeliveryResult{NotificationID: n.ID, Status: domain.DeliverySuppressed, ErrorCode: "DUPLICATE", ErrorMessage: domain.ErrDuplicate.Error()}
			c.save(ctx, res)
			return res, domain.ErrDuplicate
		}
	}

	if c.policy != nil {
		decision, err := c.policy.Evaluate(ctx, n)
		if err != nil { return failedCP(n, "POLICY_FAILED", err), err }
		if decision.Mutate != nil { decision.Mutate(&n) }
		if !decision.Allow {
			res := domain.DeliveryResult{NotificationID: n.ID, Status: domain.DeliverySuppressed, ErrorCode: "SUPPRESSED", ErrorMessage: decision.Reason}
			c.save(ctx, res)
			return res, domain.ErrSuppressed
		}
	}

	if c.limiter != nil {
		key := string(n.Channel) + ":" + n.TenantID
		allowed, retryAfter, err := c.limiter.Allow(ctx, n.TenantID, key, n.Priority.Weight())
		if err != nil { return failedCP(n, "RATELIMIT_CHECK_FAILED", err), err }
		if !allowed {
			res := domain.DeliveryResult{
				NotificationID: n.ID,
				Status: domain.DeliveryFailed,
				ErrorCode: "RATE_LIMITED",
				ErrorMessage: "retry_after=" + retryAfter.String(),
				Retryable: true,
			}
			c.save(ctx, res)
			return res, domain.ErrRateLimited
		}
	}

	resolved, err := c.resolveAudience(ctx, n)
	if err != nil { return failedCP(n, "AUDIENCE_RESOLVE_FAILED", err), err }

	sender, err := c.router.Route(ctx, resolved)
	if err != nil { return failedCP(n, "ROUTE_FAILED", err), err }

	var result domain.DeliveryResult
	var sendErr error
	plan := RetryPlan{MaxAttempts: c.cfg.MaxAttempts, BaseDelay: c.cfg.RetryBaseDelay, MaxDelay: c.cfg.RetryMaxDelay}

	for attempt := 1; attempt <= c.cfg.MaxAttempts; attempt++ {
		result, sendErr = sender.Send(ctx, resolved)
		result.Attempt = attempt

		if sendErr == nil || !result.Retryable {
			break
		}
		if attempt < c.cfg.MaxAttempts {
			select {
			case <-ctx.Done():
				return failedCP(n, "CONTEXT_CANCELLED", ctx.Err()), ctx.Err()
			case <-time.After(plan.Delay(attempt)):
			}
		}
	}

	if n.IdempotencyKey != "" && c.idem != nil && sendErr == nil {
		_ = c.idem.Mark(ctx, n.TenantID, n.IdempotencyKey, c.cfg.IdempotencyTTL)
	}

	c.save(ctx, result)
	return result, sendErr
}

func (c *ControlPlane) resolveAudience(ctx context.Context, n domain.Notification) (domain.Notification, error) {
	if c.resolver == nil {
		if n.Audience.Type == domain.AudienceUser || n.Audience.Type == domain.AudienceSegment {
			return n, errors.New("token resolver is required")
		}
		return n, nil
	}

	switch n.Audience.Type {
	case domain.AudienceUser:
		tokens, err := c.resolver.ResolveUserTokens(ctx, n.TenantID, n.Audience.UserID)
		if err != nil { return n, err }
		n.Audience = domain.Audience{Type: domain.AudienceTokens, Tokens: tokens}
	case domain.AudienceSegment:
		tokens, err := c.resolver.ResolveSegmentTokens(ctx, n.TenantID, n.Audience.SegmentID)
		if err != nil { return n, err }
		n.Audience = domain.Audience{Type: domain.AudienceTokens, Tokens: tokens}
	}
	return n, nil
}

func (c *ControlPlane) save(ctx context.Context, res domain.DeliveryResult) {
	if c.store != nil {
		_ = c.store.SaveResult(ctx, res)
	}
}

func failedCP(n domain.Notification, code string, err error) domain.DeliveryResult {
	msg := ""
	if err != nil { msg = err.Error() }
	return domain.DeliveryResult{NotificationID: n.ID, Status: domain.DeliveryFailed, ErrorCode: code, ErrorMessage: msg}
}
