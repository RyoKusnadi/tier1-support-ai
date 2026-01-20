package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/RyoKusnadi/tier1-support-ai/internal/config"
	"github.com/RyoKusnadi/tier1-support-ai/internal/handler"
	"github.com/RyoKusnadi/tier1-support-ai/internal/llm"
	"github.com/RyoKusnadi/tier1-support-ai/internal/logger"
	"github.com/RyoKusnadi/tier1-support-ai/internal/middleware"
	"github.com/RyoKusnadi/tier1-support-ai/internal/observability"
	"github.com/RyoKusnadi/tier1-support-ai/internal/reliability"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	// Initialize LLM client
	llmConfig := llm.Config{
		Provider:     cfg.LLMProvider,
		APIKey:       cfg.LLMAPIKey,
		BaseURL:      cfg.LLMBaseURL,
		DefaultModel: cfg.LLMDefaultModel,
		MaxTokens:    cfg.LLMMaxTokens,
		Temperature:  cfg.LLMTemperature,
		Timeout:      cfg.LLMTimeout,
		MaxRetries:   cfg.LLMMaxRetries,
		RetryDelay:   cfg.LLMRetryDelay,
	}

	llmClient, err := llm.NewClient(llmConfig)
	if err != nil {
		log.Fatalf("failed to initialize LLM client: %v", err)
	}

	// Initialize reliability & cost-control primitives (Phase 5)
	rateLimiter := reliability.NewTenantRateLimiter(cfg.TenantRateLimitPerSec, cfg.TenantRateLimitBurst)
	responseCache := reliability.NewResponseCache[handler.SupportQueryResponse](time.Duration(cfg.ResponseCacheTTLSeconds) * time.Second)
	tokenUsageTracker := reliability.NewTokenUsageTracker(time.Duration(cfg.TokenUsageWindowHours) * time.Hour)
	var budgetGuard *reliability.BudgetGuard
	if cfg.TenantTokenBudget > 0 {
		budgetGuard = reliability.NewBudgetGuard(tokenUsageTracker, cfg.TenantTokenBudget)
	}

	// Initialize handlers
	metrics := observability.New()
	supportHandler := handler.NewSupportHandler(llmClient, rateLimiter, responseCache, tokenUsageTracker, budgetGuard, metrics)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.RequestLogger())
	router.Use(observability.Middleware(metrics))

	router.GET("/health", handler.Health)
	router.GET("/metrics", metrics.Handler)

	// Register support query endpoint
	v1 := router.Group("/v1")
	{
		support := v1.Group("/support")
		{
			support.POST("/query", supportHandler.SupportQuery)
		}
	}

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	go func() {
		logger.Info("server starting on :8080", nil)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("server shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	log.Println("server exited properly")
}
