package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/RyoKusnadi/tier1-support-ai/internal/logger"
	"github.com/gin-gonic/gin"
)

const (
	HeaderRequestID = "X-Request-Id"
	CtxRequestID    = "request_id"
)

// RequestLogger adds a request_id (if missing), logs request completion, and records latency.
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		reqID := c.GetHeader(HeaderRequestID)
		if reqID == "" {
			reqID = newRequestID()
		}
		c.Set(CtxRequestID, reqID)
		c.Writer.Header().Set(HeaderRequestID, reqID)

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		fields := map[string]interface{}{
			"request_id":  reqID,
			"method":      c.Request.Method,
			"path":        c.FullPath(),
			"status":      status,
			"latency_ms":  latency.Milliseconds(),
			"client_ip":   c.ClientIP(),
			"user_agent":  c.Request.UserAgent(),
		}
		// Tenant ID is best-effort; handlers can set it.
		if tenantID, ok := c.Get("tenant_id"); ok {
			fields["tenant_id"] = tenantID
		}

		logger.Info("request completed", fields)
	}
}

func newRequestID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}


