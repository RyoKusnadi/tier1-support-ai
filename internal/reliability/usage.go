package reliability

import (
	"sync"
	"time"
)

// TokenUsage is a snapshot of token usage for a tenant in the current window.
type TokenUsage struct {
	TenantID    string
	WindowStart time.Time
	Window      time.Duration

	Requests   int
	TokensUsed int
}

// TokenUsageTracker tracks per-tenant token usage in-memory.
// Budget enforcement is handled by BudgetGuard, which consumes this tracker.
type TokenUsageTracker struct {
	window time.Duration

	mu     sync.Mutex
	tenants map[string]*TokenUsage
	now    func() time.Time
}

func NewTokenUsageTracker(window time.Duration) *TokenUsageTracker {
	if window <= 0 {
		window = 24 * time.Hour
	}
	return &TokenUsageTracker{
		window:  window,
		tenants: map[string]*TokenUsage{},
		now:     time.Now,
	}
}

func (t *TokenUsageTracker) Add(tenantID string, tokensUsed int) TokenUsage {
	if tenantID == "" {
		tenantID = "_unknown"
	}
	if tokensUsed < 0 {
		tokensUsed = 0
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.now()
	u := t.tenants[tenantID]
	if u == nil || now.Sub(u.WindowStart) >= t.window {
		u = &TokenUsage{
			TenantID:    tenantID,
			WindowStart: now,
			Window:      t.window,
		}
		t.tenants[tenantID] = u
	}

	u.Requests += 1
	u.TokensUsed += tokensUsed
	return *u
}

func (t *TokenUsageTracker) Get(tenantID string) (TokenUsage, bool) {
	if tenantID == "" {
		tenantID = "_unknown"
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	u := t.tenants[tenantID]
	if u == nil {
		return TokenUsage{}, false
	}

	now := t.now()
	if now.Sub(u.WindowStart) >= t.window {
		delete(t.tenants, tenantID)
		return TokenUsage{}, false
	}

	return *u, true
}


