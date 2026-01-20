package llm

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRetry_Success(t *testing.T) {
	attempts := 0
	fn := func() error {
		attempts++
		return nil
	}

	err := Retry(context.Background(), fn, DefaultRetryConfig())
	if err != nil {
		t.Errorf("Retry() error = %v, want nil", err)
	}
	if attempts != 1 {
		t.Errorf("Retry() attempts = %d, want 1", attempts)
	}
}

func TestRetry_RetryableError(t *testing.T) {
	attempts := 0
	fn := func() error {
		attempts++
		if attempts < 2 {
			return &RetryableError{Err: errors.New("temporary error"), Retryable: true}
		}
		return nil
	}

	err := Retry(context.Background(), fn, DefaultRetryConfig())
	if err != nil {
		t.Errorf("Retry() error = %v, want nil", err)
	}
	if attempts != 2 {
		t.Errorf("Retry() attempts = %d, want 2", attempts)
	}
}

func TestRetry_NonRetryableError(t *testing.T) {
	attempts := 0
	fn := func() error {
		attempts++
		return &RetryableError{Err: errors.New("permanent error"), Retryable: false}
	}

	err := Retry(context.Background(), fn, DefaultRetryConfig())
	if err == nil {
		t.Error("Retry() error = nil, want error")
	}
	if attempts != 1 {
		t.Errorf("Retry() attempts = %d, want 1", attempts)
	}
}

func TestRetry_MaxRetries(t *testing.T) {
	attempts := 0
	fn := func() error {
		attempts++
		return &RetryableError{Err: errors.New("temporary error"), Retryable: true}
	}

	config := DefaultRetryConfig()
	config.MaxRetries = 2
	config.InitialDelay = 10 * time.Millisecond

	err := Retry(context.Background(), fn, config)
	if err == nil {
		t.Error("Retry() error = nil, want error")
	}
	if attempts != 3 { // Initial attempt + 2 retries
		t.Errorf("Retry() attempts = %d, want 3", attempts)
	}
}

func TestRetry_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	fn := func() error {
		return &RetryableError{Err: errors.New("temporary error"), Retryable: true}
	}

	err := Retry(ctx, fn, DefaultRetryConfig())
	if err != context.Canceled {
		t.Errorf("Retry() error = %v, want context.Canceled", err)
	}
}


