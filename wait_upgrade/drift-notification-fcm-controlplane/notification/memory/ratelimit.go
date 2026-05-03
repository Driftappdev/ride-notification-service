package memory

import (
	"context"
	"sync"
	"time"
)

type BucketLimiter struct {
	mu       sync.Mutex
	capacity int
	refill   int
	window   time.Duration
	buckets  map[string]*bucket
}

type bucket struct {
	tokens    int
	updatedAt time.Time
}

func NewBucketLimiter(capacity int, refill int, window time.Duration) *BucketLimiter {
	if capacity <= 0 { capacity = 1000 }
	if refill <= 0 { refill = capacity }
	if window <= 0 { window = time.Minute }
	return &BucketLimiter{
		capacity: capacity,
		refill: refill,
		window: window,
		buckets: map[string]*bucket{},
	}
}

func (l *BucketLimiter) Allow(ctx context.Context, tenantID string, key string, cost int) (bool, time.Duration, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if cost <= 0 { cost = 1 }
	now := time.Now()
	b := l.buckets[key]
	if b == nil {
		b = &bucket{tokens: l.capacity, updatedAt: now}
		l.buckets[key] = b
	}

	elapsed := now.Sub(b.updatedAt)
	if elapsed >= l.window {
		windows := int(elapsed / l.window)
		b.tokens += windows * l.refill
		if b.tokens > l.capacity { b.tokens = l.capacity }
		b.updatedAt = now
	}

	if b.tokens < cost {
		return false, l.window - now.Sub(b.updatedAt), nil
	}
	b.tokens -= cost
	return true, 0, nil
}
