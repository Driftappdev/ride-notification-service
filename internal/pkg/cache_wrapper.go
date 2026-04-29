package pkg

import (
	gocache "github.com/driftappdev/libpackage/resilience/cache"
	"time"
)

type BoolCache = gocache.Cache[string, bool]

func NewBoolCache(maxSize int, ttl time.Duration) *BoolCache {
	return gocache.New[string, bool](gocache.Options{MaxSize: maxSize, DefaultTTL: ttl, Policy: gocache.PolicyLRU})
}

func CacheGetBool(c *BoolCache, key string) (bool, bool) {
	v, err := c.Get(key)
	if err != nil {
		return false, false
	}
	return v, true
}

func CacheSetBool(c *BoolCache, key string, value bool, ttl time.Duration) { c.Set(key, value, ttl) }
