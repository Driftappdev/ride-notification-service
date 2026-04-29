package pkg

import (
	gocircuit "github.com/driftappdev/libpackage/gocircuit"
	"time"
)

type CircuitBreaker = gocircuit.Breaker

func NewCircuitBreaker(name string, openTimeout time.Duration, minVolume int64) *CircuitBreaker {
	cfg := gocircuit.DefaultConfig(name)
	cfg.OpenTimeout = openTimeout
	cfg.MinimumRequestVolume = minVolume
	return gocircuit.New(cfg)
}
