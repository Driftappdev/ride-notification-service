package controlplane

import (
	"math"
	"math/rand"
	"time"
)

type RetryPlan struct {
	MaxAttempts int
	BaseDelay   time.Duration
	MaxDelay    time.Duration
	JitterRatio float64
}

func (p RetryPlan) Delay(attempt int) time.Duration {
	if attempt <= 0 { attempt = 1 }
	if p.BaseDelay <= 0 { p.BaseDelay = 250 * time.Millisecond }
	if p.MaxDelay <= 0 { p.MaxDelay = 30 * time.Second }
	if p.JitterRatio <= 0 { p.JitterRatio = 0.25 }

	exp := float64(p.BaseDelay) * math.Pow(2, float64(attempt-1))
	delay := time.Duration(exp)
	if delay > p.MaxDelay { delay = p.MaxDelay }

	jitter := int64(float64(delay) * p.JitterRatio)
	if jitter <= 0 { return delay }

	return delay - time.Duration(jitter) + time.Duration(rand.Int63n(jitter*2+1))
}
