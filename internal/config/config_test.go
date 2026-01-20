package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Save original values
	originalPort := os.Getenv("PORT")
	originalEnv := os.Getenv("APP_ENV")
	originalLLMProvider := os.Getenv("LLM_PROVIDER")

	// Clean up
	defer func() {
		if originalPort != "" {
			os.Setenv("PORT", originalPort)
		} else {
			os.Unsetenv("PORT")
		}
		if originalEnv != "" {
			os.Setenv("APP_ENV", originalEnv)
		} else {
			os.Unsetenv("APP_ENV")
		}
		if originalLLMProvider != "" {
			os.Setenv("LLM_PROVIDER", originalLLMProvider)
		} else {
			os.Unsetenv("LLM_PROVIDER")
		}
	}()

	// Test defaults
	os.Unsetenv("PORT")
	os.Unsetenv("APP_ENV")
	os.Unsetenv("LLM_PROVIDER")

	cfg := Load()
	if cfg.Port != "8080" {
		t.Errorf("Load() Port = %s, want 8080", cfg.Port)
	}
	if cfg.Env != "development" {
		t.Errorf("Load() Env = %s, want development", cfg.Env)
	}
	if cfg.LLMProvider != "openai" {
		t.Errorf("Load() LLMProvider = %s, want openai", cfg.LLMProvider)
	}

	// Test custom values
	os.Setenv("PORT", "9000")
	os.Setenv("APP_ENV", "production")
	os.Setenv("LLM_PROVIDER", "anthropic")

	cfg = Load()
	if cfg.Port != "9000" {
		t.Errorf("Load() Port = %s, want 9000", cfg.Port)
	}
	if cfg.Env != "production" {
		t.Errorf("Load() Env = %s, want production", cfg.Env)
	}
	if cfg.LLMProvider != "anthropic" {
		t.Errorf("Load() LLMProvider = %s, want anthropic", cfg.LLMProvider)
	}
}


