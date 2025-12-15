package handler

import (
	"net/http"

	"github.com/RyoKusnadi/tier1-support-ai/internal/config"
	"github.com/gin-gonic/gin"
)

type SupportQueryRequest struct {
	TenantID string `json:"tenant_id" binding:"required"`
	Language string `json:"language" binding:"required"`
	Question string `json:"question" binding:"required"`
}

type SupportQueryResponse struct {
	Answer     string  `json:"answer"`
	Confidence float64 `json:"confidence"`
	Fallback   bool    `json:"fallback,omitempty"`
}

func SupportQuery(c *gin.Context) {
	var req SupportQueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": err.Error(),
			},
		})
		return
	}

	if !config.Tenants[req.TenantID] {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code":    "TENANT_NOT_FOUND",
				"message": "tenant not found",
			},
		})
		return
	}

	if !config.SupportedLanguages[req.Language] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "UNSUPPORTED_LANGUAGE",
				"message": "language not supported",
			},
		})
		return
	}

	resp := SupportQueryResponse{
		Answer:     "Thank you for your question. Our support team will assist you shortly.",
		Confidence: 0.5,
		Fallback:   true,
	}

	c.JSON(http.StatusOK, resp)
}
