package memory

import (
	"context"
	"sync"
	"time"
)

type idemItem struct {
	expireAt time.Time
}

type IdempotencyStore struct {
	mu sync.Mutex
	m  map[string]idemItem
}

func NewIdempotencyStore() *IdempotencyStore {
	return &IdempotencyStore{m: map[string]idemItem{}}
}

func (s *IdempotencyStore) Seen(ctx context.Context, tenantID, key string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	k := tenantID + ":" + key
	item, ok := s.m[k]
	if !ok { return false, nil }
	if time.Now().After(item.expireAt) {
		delete(s.m, k)
		return false, nil
	}
	return true, nil
}

func (s *IdempotencyStore) Mark(ctx context.Context, tenantID, key string, ttl time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.m[tenantID+":"+key] = idemItem{expireAt: time.Now().Add(ttl)}
	return nil
}
