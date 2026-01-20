package observability

import (
	"encoding/json"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
)

// Metrics is a small in-process metrics registry (Phase 6).
// It is intentionally dependency-free; can be swapped for Prometheus later.
type Metrics struct {
	RequestsTotal      atomic.Int64
	ErrorsTotal        atomic.Int64
	RateLimitedTotal   atomic.Int64
	BudgetBlockedTotal atomic.Int64
	CacheHitsTotal     atomic.Int64
	CacheMissesTotal   atomic.Int64

	LatencyCount atomic.Int64
	LatencySumMs atomic.Int64
}

func New() *Metrics { return &Metrics{} }

func (m *Metrics) ObserveRequest(latency time.Duration, status int) {
	m.RequestsTotal.Add(1)
	if status >= 400 {
		m.ErrorsTotal.Add(1)
	}
	m.LatencyCount.Add(1)
	m.LatencySumMs.Add(latency.Milliseconds())
}

func (m *Metrics) Snapshot() map[string]interface{} {
	count := m.LatencyCount.Load()
	sum := m.LatencySumMs.Load()
	var avg float64
	if count > 0 {
		avg = float64(sum) / float64(count)
	}

	return map[string]interface{}{
		"requests_total":       m.RequestsTotal.Load(),
		"errors_total":         m.ErrorsTotal.Load(),
		"rate_limited_total":   m.RateLimitedTotal.Load(),
		"budget_blocked_total": m.BudgetBlockedTotal.Load(),
		"cache_hits_total":     m.CacheHitsTotal.Load(),
		"cache_misses_total":   m.CacheMissesTotal.Load(),
		"latency_count":        count,
		"latency_sum_ms":       sum,
		"latency_avg_ms":       avg,
	}
}

// Handler exposes metrics as JSON at GET /metrics.
func (m *Metrics) Handler(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	enc := json.NewEncoder(c.Writer)
	_ = enc.Encode(m.Snapshot())
}

// Middleware records per-request latency + status.
func Middleware(m *Metrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		m.ObserveRequest(time.Since(start), c.Writer.Status())
	}
}

// Ensure we don't accidentally import net/http without using it in some build tags.
var _ = http.StatusOK

