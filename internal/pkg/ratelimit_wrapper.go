package pkg

import (
	"context"
	goratelimit "github.com/driftappdev/libpackage/goratelimit"
	"net/http"
	"time"
)

type RateLimiter = goratelimit.Limiter
type KeyExtractor = goratelimit.KeyExtractor
type RateLimitResult = goratelimit.Result

func NewIPSlidingWindow(limit int, window time.Duration) (RateLimiter, KeyExtractor) {
	return goratelimit.NewSlidingWindow(limit, window), goratelimit.ByIP()
}

func Allow(limiter RateLimiter, ctx context.Context, key string) (goratelimit.Result, error) {
	return limiter.Allow(ctx, key)
}

func LimitKey(extractor KeyExtractor, r *http.Request) string { return extractor(r) }
