package main

import (
	"testing"
	"time"

	libpack_logger "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
	"github.com/stretchr/testify/assert"
)

func TestNewRetryBudget(t *testing.T) {
	tests := []struct {
		name   string
		config RetryBudgetConfig
	}{
		{
			name: "default config",
			config: RetryBudgetConfig{
				TokensPerSecond: 10.0,
				MaxTokens:       100,
				Enabled:         true,
			},
		},
		{
			name: "custom config",
			config: RetryBudgetConfig{
				TokensPerSecond: 50.0,
				MaxTokens:       500,
				Enabled:         true,
			},
		},
		{
			name: "disabled config",
			config: RetryBudgetConfig{
				TokensPerSecond: 10.0,
				MaxTokens:       100,
				Enabled:         false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := libpack_logger.New()

			rb := NewRetryBudget(tt.config, logger)

			assert.NotNil(t, rb)
			assert.Equal(t, tt.config.Enabled, rb.enabled)
			assert.Equal(t, tt.config.TokensPerSecond, rb.tokensPerSecond)
			assert.Equal(t, int64(tt.config.MaxTokens), rb.maxTokens)

			if tt.config.Enabled {
				// Should start with max tokens
				assert.Equal(t, int64(tt.config.MaxTokens), rb.currentTokens.Load())
			}
		})
	}
}

func TestRetryBudget_Allow(t *testing.T) {
	t.Run("allow when enabled and tokens available", func(t *testing.T) {
		config := RetryBudgetConfig{
			TokensPerSecond: 10.0,
			MaxTokens:       100,
			Enabled:         true,
		}

		rb := NewRetryBudget(config, libpack_logger.New())

		// Should allow first request
		allowed := rb.AllowRetry()
		assert.True(t, allowed)

		// Tokens should be decremented
		assert.Less(t, rb.currentTokens.Load(), int64(100))
	})

	t.Run("deny when tokens exhausted", func(t *testing.T) {
		config := RetryBudgetConfig{
			TokensPerSecond: 10.0,
			MaxTokens:       2,
			Enabled:         true,
		}

		rb := NewRetryBudget(config, libpack_logger.New())

		// Consume all tokens
		assert.True(t, rb.AllowRetry())
		assert.True(t, rb.AllowRetry())

		// Should deny when exhausted
		assert.False(t, rb.AllowRetry())

		stats := rb.GetStats()
		assert.Greater(t, stats["denied_retries"].(int64), int64(0))
	})

	t.Run("always allow when disabled", func(t *testing.T) {
		config := RetryBudgetConfig{
			TokensPerSecond: 10.0,
			MaxTokens:       0,
			Enabled:         false,
		}

		rb := NewRetryBudget(config, libpack_logger.New())

		// Should always allow when disabled
		for i := 0; i < 100; i++ {
			assert.True(t, rb.AllowRetry())
		}
	})
}

func TestRetryBudget_Refill(t *testing.T) {
	t.Run("tokens refill over time", func(t *testing.T) {
		config := RetryBudgetConfig{
			TokensPerSecond: 100.0, // Fast refill for testing
			MaxTokens:       100,
			Enabled:         true,
		}

		rb := NewRetryBudget(config, libpack_logger.New())

		// Consume some tokens
		for i := 0; i < 50; i++ {
			rb.AllowRetry()
		}

		tokensBefore := rb.currentTokens.Load()

		// Wait for refill (multiple refill cycles at 100ms each)
		time.Sleep(300 * time.Millisecond)

		tokensAfter := rb.currentTokens.Load()

		// Tokens should have increased
		assert.Greater(t, tokensAfter, tokensBefore)
	})

	t.Run("tokens don't exceed max", func(t *testing.T) {
		config := RetryBudgetConfig{
			TokensPerSecond: 100.0,
			MaxTokens:       50,
			Enabled:         true,
		}

		rb := NewRetryBudget(config, libpack_logger.New())

		// Wait for potential overflow
		time.Sleep(200 * time.Millisecond)

		tokens := rb.currentTokens.Load()
		assert.LessOrEqual(t, tokens, int64(50))
	})
}

func TestRetryBudget_GetStats(t *testing.T) {
	t.Run("tracks statistics correctly", func(t *testing.T) {
		config := RetryBudgetConfig{
			TokensPerSecond: 10.0,
			MaxTokens:       5,
			Enabled:         true,
		}

		rb := NewRetryBudget(config, libpack_logger.New())

		// Allow some requests
		rb.AllowRetry()
		rb.AllowRetry()
		rb.AllowRetry()

		// Consume all tokens to trigger denials
		rb.AllowRetry()
		rb.AllowRetry()
		rb.AllowRetry() // Should be denied
		rb.AllowRetry() // Should be denied

		stats := rb.GetStats()

		assert.Equal(t, true, stats["enabled"])
		assert.Equal(t, 10.0, stats["tokens_per_sec"])
		assert.Equal(t, int64(5), stats["max_tokens"])
		assert.GreaterOrEqual(t, stats["current_tokens"].(int64), int64(0))
		assert.Equal(t, int64(7), stats["total_attempts"])
		assert.GreaterOrEqual(t, stats["denied_retries"].(int64), int64(2))
		assert.Greater(t, stats["denial_rate_pct"].(float64), 0.0)
	})

	t.Run("stats when disabled", func(t *testing.T) {
		config := RetryBudgetConfig{
			TokensPerSecond: 10.0,
			MaxTokens:       100,
			Enabled:         false,
		}

		rb := NewRetryBudget(config, libpack_logger.New())

		stats := rb.GetStats()

		assert.Equal(t, false, stats["enabled"])
		assert.Equal(t, int64(0), stats["total_attempts"])
		assert.Equal(t, int64(0), stats["denied_retries"])
	})
}

func TestRetryBudget_Reset(t *testing.T) {
	config := RetryBudgetConfig{
		TokensPerSecond: 10.0,
		MaxTokens:       10,
		Enabled:         true,
	}

	rb := NewRetryBudget(config, libpack_logger.New())

	// Generate some activity
	for i := 0; i < 15; i++ {
		rb.AllowRetry()
	}

	statsBefore := rb.GetStats()
	assert.Greater(t, statsBefore["total_attempts"].(int64), int64(0))

	// Reset
	rb.Reset()

	statsAfter := rb.GetStats()
	assert.Equal(t, int64(0), statsAfter["total_attempts"])
	assert.Equal(t, int64(0), statsAfter["denied_retries"])
	assert.Equal(t, int64(10), statsAfter["current_tokens"]) // Should reset to max
}

func TestRetryBudget_ConcurrentAccess(t *testing.T) {
	config := RetryBudgetConfig{
		TokensPerSecond: 100.0,
		MaxTokens:       1000,
		Enabled:         true,
	}

	rb := NewRetryBudget(config, libpack_logger.New())

	// Concurrent access test
	done := make(chan bool)
	goroutines := 100
	requestsPerGoroutine := 10

	for i := 0; i < goroutines; i++ {
		go func() {
			for j := 0; j < requestsPerGoroutine; j++ {
				rb.AllowRetry()
			}
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < goroutines; i++ {
		<-done
	}

	stats := rb.GetStats()
	totalAttempts := stats["total_attempts"].(int64)

	// Should have processed all requests
	assert.Equal(t, int64(goroutines*requestsPerGoroutine), totalAttempts)
}

func TestRetryBudget_DenialRate(t *testing.T) {
	config := RetryBudgetConfig{
		TokensPerSecond: 1.0,
		MaxTokens:       10,
		Enabled:         true,
	}

	rb := NewRetryBudget(config, libpack_logger.New())

	// Consume all tokens
	for i := 0; i < 10; i++ {
		rb.AllowRetry()
	}

	// These should be denied
	deniedCount := 0
	for i := 0; i < 10; i++ {
		if !rb.AllowRetry() {
			deniedCount++
		}
	}

	assert.Greater(t, deniedCount, 0)

	stats := rb.GetStats()
	denialRate := stats["denial_rate_pct"].(float64)

	assert.Greater(t, denialRate, 0.0)
	assert.LessOrEqual(t, denialRate, 100.0)
}

func TestRetryBudget_GlobalInstance(t *testing.T) {
	config := RetryBudgetConfig{
		TokensPerSecond: 10.0,
		MaxTokens:       100,
		Enabled:         true,
	}

	rb := InitializeRetryBudget(config, libpack_logger.New())
	assert.NotNil(t, rb)

	// Should return the same instance
	rb2 := GetRetryBudget()
	assert.Equal(t, rb, rb2)
}
