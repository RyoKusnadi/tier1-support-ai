package router

import (
	"github.com/RyoKusnadi/tier1-support-ai/internal/handler"
	"github.com/gin-gonic/gin"
)

func RegisterV1Routes(r *gin.Engine, supportHandler *handler.SupportHandler) {
	v1 := r.Group("/v1")
	{
		v1.POST("/support/query", supportHandler.SupportQuery)
	}
}
