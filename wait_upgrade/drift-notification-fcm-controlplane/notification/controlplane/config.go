package controlplane

import "time"

type Config struct {
	IdempotencyTTL time.Duration
	MaxAttempts    int
	RetryBaseDelay time.Duration
	RetryMaxDelay  time.Duration
}

func (c Config) Normalize() Config {
	if c.IdempotencyTTL <= 0 { c.IdempotencyTTL = 24 * time.Hour }
	if c.MaxAttempts <= 0 { c.MaxAttempts = 5 }
	if c.RetryBaseDelay <= 0 { c.RetryBaseDelay = 250 * time.Millisecond }
	if c.RetryMaxDelay <= 0 { c.RetryMaxDelay = 30 * time.Second }
	return c
}
