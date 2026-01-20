package reliability

import (
	"sync"
	"time"
)

// TenantRateLimiter is a simple per-tenant token bucket rate limiter.
// It is intentionally in-memory (Phase 5); can be swapped for Redis later.
type TenantRateLimiter struct {
	ratePerSec float64
	burst      float64

	mu      sync.Mutex
	buckets map[string]*tokenBucket
	now     func() time.Time
}

type tokenBucket struct {
	tokens     float64
	lastRefill time.Time
}

func NewTenantRateLimiter(ratePerSec float64, burst int) *TenantRateLimiter {
	if ratePerSec <= 0 {
		ratePerSec = 5
	}
	if burst <= 0 {
		burst = 10
	}
	return &TenantRateLimiter{
		ratePerSec: ratePerSec,
		burst:      float64(burst),
		buckets:    map[string]*tokenBucket{},
		now:        time.Now,
	}
}

// Allow returns true if the tenant is allowed to proceed right now.
func (l *TenantRateLimiter) Allow(tenantID string) bool {
	if tenantID == "" {
		tenantID = "_unknown"
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	b := l.buckets[tenantID]
	now := l.now()
	if b == nil {
		l.buckets[tenantID] = &tokenBucket{tokens: l.burst - 1, lastRefill: now}
		return true
	}

	// Refill based on elapsed time.
	elapsed := now.Sub(b.lastRefill).Seconds()
	if elapsed > 0 {
		b.tokens = minFloat(l.burst, b.tokens+(elapsed*l.ratePerSec))
		b.lastRefill = now
	}

	if b.tokens < 1 {
		return false
	}
	b.tokens -= 1
	return true
}

func minFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

