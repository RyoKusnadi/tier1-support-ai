package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port string
	Env  string

	// LLM Configuration
	LLMProvider     string
	LLMAPIKey       string
	LLMBaseURL      string
	LLMDefaultModel string
	LLMMaxTokens    int
	LLMTemperature  float64
	LLMTimeout      int
	LLMMaxRetries   int
	LLMRetryDelay   int

	// Reliability & cost control (Phase 5)
	// Per-tenant rate limiting (token bucket)
	TenantRateLimitPerSec float64
	TenantRateLimitBurst  int

	// Response cache TTL in seconds
	ResponseCacheTTLSeconds int

	// Token usage tracking window (hours) and per-tenant token budget per window
	TokenUsageWindowHours int
	TenantTokenBudget     int
}

func Load() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	// LLM Configuration
	llmProvider := os.Getenv("LLM_PROVIDER")
	if llmProvider == "" {
		llmProvider = "openai"
	}

	llmAPIKey := os.Getenv("LLM_API_KEY")
	llmBaseURL := os.Getenv("LLM_BASE_URL")
	llmDefaultModel := os.Getenv("LLM_DEFAULT_MODEL")
	if llmDefaultModel == "" {
		llmDefaultModel = "gpt-3.5-turbo"
	}

	llmMaxTokens := getIntEnv("LLM_MAX_TOKENS", 500)
	llmTemperature := getFloatEnv("LLM_TEMPERATURE", 0.7)
	llmTimeout := getIntEnv("LLM_TIMEOUT", 30)
	llmMaxRetries := getIntEnv("LLM_MAX_RETRIES", 3)
	llmRetryDelay := getIntEnv("LLM_RETRY_DELAY", 100)

	// Reliability & cost control (Phase 5)
	tenantRateLimitPerSec := getFloatEnv("TENANT_RATE_LIMIT_PER_SEC", 5.0)
	tenantRateLimitBurst := getIntEnv("TENANT_RATE_LIMIT_BURST", 10)
	responseCacheTTLSeconds := getIntEnv("RESPONSE_CACHE_TTL_SECONDS", 300)
	tokenUsageWindowHours := getIntEnv("TOKEN_USAGE_WINDOW_HOURS", 24)
	tenantTokenBudget := getIntEnv("TENANT_TOKEN_BUDGET", 0) // 0 = disabled

	return Config{
		Port: port,
		Env:  env,

		LLMProvider:     llmProvider,
		LLMAPIKey:       llmAPIKey,
		LLMBaseURL:      llmBaseURL,
		LLMDefaultModel: llmDefaultModel,
		LLMMaxTokens:    llmMaxTokens,
		LLMTemperature:  llmTemperature,
		LLMTimeout:      llmTimeout,
		LLMMaxRetries:   llmMaxRetries,
		LLMRetryDelay:   llmRetryDelay,

		TenantRateLimitPerSec:   tenantRateLimitPerSec,
		TenantRateLimitBurst:    tenantRateLimitBurst,
		ResponseCacheTTLSeconds: responseCacheTTLSeconds,
		TokenUsageWindowHours:   tokenUsageWindowHours,
		TenantTokenBudget:       tenantTokenBudget,
	}
}

func getIntEnv(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}

func getFloatEnv(key string, defaultValue float64) float64 {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return defaultValue
	}
	return floatValue
}
