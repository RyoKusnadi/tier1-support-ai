package llm

import (
	"context"
	"fmt"
	"math"
	"time"
)

// RetryConfig holds retry configuration
type RetryConfig struct {
	MaxRetries int           // Maximum number of retries
	InitialDelay time.Duration // Initial delay before first retry
	MaxDelay     time.Duration // Maximum delay between retries
	Multiplier   float64       // Exponential backoff multiplier
}

// DefaultRetryConfig returns a default retry configuration
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries:   3,
		InitialDelay: 100 * time.Millisecond,
		MaxDelay:     5 * time.Second,
		Multiplier:   2.0,
	}
}

// RetryableError represents an error that can be retried
type RetryableError struct {
	Err      error
	Retryable bool
}

func (e *RetryableError) Error() string {
	return e.Err.Error()
}

// IsRetryable checks if an error is retryable
func IsRetryable(err error) bool {
	if retryErr, ok := err.(*RetryableError); ok {
		return retryErr.Retryable
	}
	return false
}

// Retry executes a function with retry logic
func Retry(ctx context.Context, fn func() error, config RetryConfig) error {
	var lastErr error
	
	for attempt := 0; attempt <= config.MaxRetries; attempt++ {
		// Check if context is cancelled
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Execute the function
		err := fn()
		if err == nil {
			return nil
		}

		lastErr = err

		// Check if error is retryable
		if !IsRetryable(err) {
			return err
		}

		// Don't retry on the last attempt
		if attempt == config.MaxRetries {
			break
		}

		// Calculate delay with exponential backoff
		delay := time.Duration(float64(config.InitialDelay) * math.Pow(config.Multiplier, float64(attempt)))
		if delay > config.MaxDelay {
			delay = config.MaxDelay
		}

		// Wait before retrying
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}
	}

	return fmt.Errorf("max retries exceeded: %w", lastErr)
}


