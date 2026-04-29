package cacheintegration

import "context"

type Store interface {
	Get(ctx context.Context, key string) ([]byte, bool, error)
	Set(ctx context.Context, key string, value []byte, ttlSeconds int64) error
	Delete(ctx context.Context, key string) error
}
