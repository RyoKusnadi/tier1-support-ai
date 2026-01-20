package reliability

import (
	"sync"
	"time"
)

type CacheEntry[T any] struct {
	Value     T
	ExpiresAt time.Time
}

// ResponseCache is a small in-memory TTL cache for support responses.
// It is safe for concurrent use.
type ResponseCache[T any] struct {
	ttl time.Duration

	mu    sync.RWMutex
	items map[string]CacheEntry[T]
	now   func() time.Time
}

func NewResponseCache[T any](ttl time.Duration) *ResponseCache[T] {
	if ttl <= 0 {
		ttl = 5 * time.Minute
	}
	return &ResponseCache[T]{
		ttl:   ttl,
		items: map[string]CacheEntry[T]{},
		now:   time.Now,
	}
}

func (c *ResponseCache[T]) Get(key string) (T, bool) {
	var zero T
	if key == "" {
		return zero, false
	}

	c.mu.RLock()
	entry, ok := c.items[key]
	c.mu.RUnlock()
	if !ok {
		return zero, false
	}

	if !entry.ExpiresAt.IsZero() && c.now().After(entry.ExpiresAt) {
		c.mu.Lock()
		delete(c.items, key)
		c.mu.Unlock()
		return zero, false
	}

	return entry.Value, true
}

func (c *ResponseCache[T]) Set(key string, value T) {
	if key == "" {
		return
	}
	c.mu.Lock()
	c.items[key] = CacheEntry[T]{
		Value:     value,
		ExpiresAt: c.now().Add(c.ttl),
	}
	c.mu.Unlock()
}



