package pkg

import (
	"context"
	goretry "github.com/driftappdev/libpackage/goretry"
	"time"
)

type RetryConfig = goretry.Config

func DefaultRetryConfig(maxAttempts int, initialDelay, maxDelay time.Duration) RetryConfig {
	cfg := goretry.DefaultConfig()
	cfg.MaxAttempts = maxAttempts
	cfg.InitialDelay = initialDelay
	cfg.MaxDelay = maxDelay
	return cfg
}

func RetryDo(ctx context.Context, cfg RetryConfig, fn func(context.Context) error) error {
	return goretry.Do(ctx, cfg, fn).Err
}
