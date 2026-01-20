package reliability

import "time"

// BudgetGuard enforces a per-tenant token budget per window.
// This is a "pre-call" guardrail: if a tenant has already exceeded budget,
// we reject further LLM calls until the window resets.
type BudgetGuard struct {
	tracker *TokenUsageTracker
	budget  int
}

func NewBudgetGuard(tracker *TokenUsageTracker, budget int) *BudgetGuard {
	if budget < 0 {
		budget = 0
	}
	return &BudgetGuard{tracker: tracker, budget: budget}
}

func (b *BudgetGuard) Enabled() bool {
	return b != nil && b.tracker != nil && b.budget > 0
}

func (b *BudgetGuard) BudgetTokens() int {
	if b == nil {
		return 0
	}
	return b.budget
}

// Allow returns false if tenant is already over budget for the current window.
func (b *BudgetGuard) Allow(tenantID string) bool {
	if !b.Enabled() {
		return true
	}
	u, ok := b.tracker.Get(tenantID)
	if !ok {
		return true
	}
	return u.TokensUsed < b.budget
}

// Remaining returns remaining budget and a boolean indicating if budget is enabled.
func (b *BudgetGuard) Remaining(tenantID string) (int, bool, time.Time) {
	if !b.Enabled() {
		return 0, false, time.Time{}
	}
	u, ok := b.tracker.Get(tenantID)
	if !ok {
		return b.budget, true, time.Time{}
	}
	remaining := b.budget - u.TokensUsed
	if remaining < 0 {
		remaining = 0
	}
	return remaining, true, u.WindowStart.Add(u.Window)
}
