package pkg

import (
	"context"
	gotimeout "github.com/driftappdev/libpackage/gotimeout"
	"time"
)

func DoWithTimeout(ctx context.Context, timeout time.Duration, fn func(context.Context) error) error {
	return gotimeout.Do(ctx, timeout, fn)
}
