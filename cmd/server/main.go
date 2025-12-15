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
	"github.com/RyoKusnadi/tier1-support-ai/internal/logger"
	router "github.com/RyoKusnadi/tier1-support-ai/internal/router"
	"github.com/gin-gonic/gin"
)

func main() {
	service := gin.New()
	service.Use(gin.Recovery())

	service.GET("/health", handler.Health)
	router.RegisterV1Routes(service)

	cfg := config.Load()
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: service,
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
