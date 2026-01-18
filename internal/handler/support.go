package handler

import (
	"net/http"

	"github.com/RyoKusnadi/tier1-support-ai/internal/llm"
	"github.com/RyoKusnadi/tier1-support-ai/internal/logger"
	"github.com/gin-gonic/gin"
)

// SupportQueryRequest represents the request body for support queries
type SupportQueryRequest struct {
	Question    string   `json:"question" binding:"required"`
	TenantID    string   `json:"tenant_id" binding:"required"`
	Language    string   `json:"language" binding:"required"`
	KnowledgeBase []string `json:"knowledge_base,omitempty"` // Optional knowledge base for RAG
}

// SupportQueryResponse represents the response for support queries
type SupportQueryResponse struct {
	Answer      string  `json:"answer"`
	Confidence  float64 `json:"confidence"`
	TenantID    string  `json:"tenant_id"`
	Language    string  `json:"language"`
}

// SupportQuery handles POST /v1/support/query requests
func SupportQuery(llmClient llm.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
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

		// Create LLM request
		llmReq := &llm.Request{
			Messages: []llm.Message{
				{
					Role:    "user",
					Content: req.Question,
				},
			},
			KnowledgeBase: req.KnowledgeBase,
			Language:      req.Language,
			TenantID:      req.TenantID,
		}

		// Generate answer using LLM
		resp, err := llmClient.GenerateAnswer(c.Request.Context(), llmReq)
		if err != nil {
			logger.Error("failed to generate answer", map[string]interface{}{
				"error":    err.Error(),
				"tenant_id": req.TenantID,
				"language": req.Language,
			})
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to generate answer",
			})
			return
		}

		// Return response
		c.JSON(http.StatusOK, SupportQueryResponse{
			Answer:     resp.Content,
			Confidence: resp.Confidence,
			TenantID:   req.TenantID,
			Language:   req.Language,
		})
	}
}

