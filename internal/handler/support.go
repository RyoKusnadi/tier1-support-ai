package handler

import (
	"net/http"

	"github.com/RyoKusnadi/tier1-support-ai/internal/knowledge"
	"github.com/RyoKusnadi/tier1-support-ai/internal/llm"
	"github.com/RyoKusnadi/tier1-support-ai/internal/logger"
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
}

// NewSupportHandler creates a new support handler
func NewSupportHandler(llmClient llm.Client) *SupportHandler {
	return &SupportHandler{
		llmClient:           llmClient,
		retriever:           knowledge.NewInMemoryRetriever(),
		confidenceThreshold: 0.7,
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

	// Apply confidence-based fallback (Phase 4 + API contract)
	isFallback := resp.Confidence < h.confidenceThreshold
	answer := resp.Content
	if isFallback {
		answer = "We are unable to confidently answer your question. Please contact customer support."
	}

	// Return response
	c.JSON(http.StatusOK, SupportQueryResponse{
		Answer:     answer,
		Confidence: resp.Confidence,
		TenantID:   req.TenantID,
		Language:   req.Language,
		Fallback:   isFallback,
	})
}
