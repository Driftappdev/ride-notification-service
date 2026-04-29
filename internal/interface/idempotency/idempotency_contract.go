package idempotencyport

import "context"

type IdempotencyStorePort interface {
	Acquire(ctx context.Context, key string) (bool, error)
	Release(ctx context.Context, key string) error
}
