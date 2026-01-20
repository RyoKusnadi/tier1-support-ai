package handler

import (
	"net/http"
	"time"

	"github.com/RyoKusnadi/tier1-support-ai/internal/knowledge"
	"github.com/RyoKusnadi/tier1-support-ai/internal/llm"
	"github.com/RyoKusnadi/tier1-support-ai/internal/logger"
	"github.com/RyoKusnadi/tier1-support-ai/internal/reliability"
	"github.com/gin-gonic/gin"
)

// SupportQueryRequest represents the request body for support queries
type SupportQueryRequest struct {
	Question      string   `json:"question" binding:"required"`
	TenantID      string   `json:"tenant_id" binding:"required"`
	Language      string   `json:"language" binding:"required"`
	KnowledgeBase []string `json:"knowledge_base,omitempty"` // Optional knowledge base for RAG
}

// SupportQueryResponse represents the response for support queries
type SupportQueryResponse struct {
	Answer     string  `json:"answer"`
	Confidence float64 `json:"confidence"`
	TenantID   string  `json:"tenant_id"`
	Language   string  `json:"language"`
	Fallback   bool    `json:"fallback,omitempty"`
}

// SupportHandler handles support-related requests
type SupportHandler struct {
	llmClient           llm.Client
	retriever           knowledge.Retriever
	confidenceThreshold float64

	// Phase 5 — Reliability & Cost Control
	rateLimiter   *reliability.TenantRateLimiter
	responseCache *reliability.ResponseCache[SupportQueryResponse]
	tokenUsage    *reliability.TokenUsageTracker
	budgetGuard   *reliability.BudgetGuard
}

// NewSupportHandler creates a new support handler
func NewSupportHandler(
	llmClient llm.Client,
	rateLimiter *reliability.TenantRateLimiter,
	responseCache *reliability.ResponseCache[SupportQueryResponse],
	tokenUsage *reliability.TokenUsageTracker,
	budgetGuard *reliability.BudgetGuard,
) *SupportHandler {
	return &SupportHandler{
		llmClient:           llmClient,
		retriever:           knowledge.NewInMemoryRetriever(),
		confidenceThreshold: 0.7,
		rateLimiter:         rateLimiter,
		responseCache:       responseCache,
		tokenUsage:          tokenUsage,
		budgetGuard:         budgetGuard,
	}
}

// SupportQuery handles POST /v1/support/query requests
func (h *SupportHandler) SupportQuery(c *gin.Context) {
	var req SupportQueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("invalid request", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request: " + err.Error(),
		})
		return
	}

	// Phase 5: per-tenant rate limiting
	if h.rateLimiter != nil && !h.rateLimiter.Allow(req.TenantID) {
		logger.Error("rate limit exceeded", map[string]interface{}{
			"tenant_id": req.TenantID,
		})
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error": gin.H{
				"code":    "RATE_LIMIT_EXCEEDED",
				"message": "Rate limit exceeded for tenant",
			},
		})
		return
	}

	// Phase 5: response caching (keyed by tenant, language, question)
	if h.responseCache != nil {
		cacheKey := buildCacheKey(req.TenantID, req.Language, req.Question)
		if cached, ok := h.responseCache.Get(cacheKey); ok {
			c.JSON(http.StatusOK, cached)
			return
		}
	}

	// Retrieve relevant knowledge (Phase 4 - Knowledge Retrieval)
	var retrievedKB []string
	if h.retriever != nil {
		kb, err := h.retriever.Retrieve(c.Request.Context(), req.TenantID, req.Language, req.Question)
		if err != nil {
			logger.Error("knowledge retrieval failed", map[string]interface{}{
				"error":     err.Error(),
				"tenant_id": req.TenantID,
				"language":  req.Language,
			})
		} else {
			retrievedKB = kb
		}
	}

	// Fallback when no relevant knowledge is found (Phase 4 requirement)
	if len(retrievedKB) == 0 {
		c.JSON(http.StatusOK, SupportQueryResponse{
			Answer:     "We are unable to confidently answer your question. Please contact customer support.",
			Confidence: 0.0,
			TenantID:   req.TenantID,
			Language:   req.Language,
			Fallback:   true,
		})
		return
	}

	// Merge retrieved knowledge with any explicit knowledge from the request
	mergedKB := append(retrievedKB, req.KnowledgeBase...)

	// Create LLM request (RAG-style: question + retrieved knowledge)
	llmReq := &llm.Request{
		Messages: []llm.Message{
			{
				Role:    "user",
				Content: req.Question,
			},
		},
		KnowledgeBase: mergedKB,
		Language:      req.Language,
		TenantID:      req.TenantID,
	}

	// Generate answer using LLM
	// Phase 5: budget guardrails (pre-call check)
	if h.budgetGuard != nil && h.budgetGuard.Enabled() && !h.budgetGuard.Allow(req.TenantID) {
		remaining, enabled, resetAt := h.budgetGuard.Remaining(req.TenantID)
		logger.Error("token budget exceeded", map[string]interface{}{
			"tenant_id": req.TenantID,
			"remaining": remaining,
			"enabled":   enabled,
			"reset_at":  resetAt.Format(time.RFC3339),
		})
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error": gin.H{
				"code":    "BUDGET_EXCEEDED",
				"message": "Token budget exceeded for tenant",
			},
		})
		return
	}

	resp, err := h.llmClient.GenerateAnswer(c.Request.Context(), llmReq)
	if err != nil {
		logger.Error("failed to generate answer", map[string]interface{}{
			"error":     err.Error(),
			"tenant_id": req.TenantID,
			"language":  req.Language,
		})
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate answer",
		})
		return
	}

	// Phase 5: token usage tracking (post-call)
	if h.tokenUsage != nil {
		usage := h.tokenUsage.Add(req.TenantID, resp.TokensUsed)
		logger.Info("token usage updated", map[string]interface{}{
			"tenant_id":   usage.TenantID,
			"tokens_used": usage.TokensUsed,
			"requests":    usage.Requests,
			"window":      usage.Window.String(),
		})
	}

	// Apply confidence-based fallback (Phase 4 + API contract)
	isFallback := resp.Confidence < h.confidenceThreshold
	answer := resp.Content
	if isFallback {
		answer = "We are unable to confidently answer your question. Please contact customer support."
	}

	// Return response
	finalResp := SupportQueryResponse{
		Answer:     answer,
		Confidence: resp.Confidence,
		TenantID:   req.TenantID,
		Language:   req.Language,
		Fallback:   isFallback,
	}

	// Store in cache for subsequent identical questions
	if h.responseCache != nil {
		cacheKey := buildCacheKey(req.TenantID, req.Language, req.Question)
		h.responseCache.Set(cacheKey, finalResp)
	}

	c.JSON(http.StatusOK, finalResp)
}

func buildCacheKey(tenantID, language, question string) string {
	return tenantID + "|" + language + "|" + question
}
