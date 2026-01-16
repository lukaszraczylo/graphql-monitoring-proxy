package main

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	libpack_logger "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
)

// RetryBudget implements a token bucket algorithm to limit the rate of retries
// This prevents retry storms and cascading failures
type RetryBudget struct {
	tokensPerSecond float64
	maxTokens       int64
	currentTokens   atomic.Int64
	lastRefill      atomic.Int64 // Unix timestamp in nanoseconds
	mu              sync.RWMutex
	enabled         bool
	logger          *libpack_logger.Logger
	ctx             context.Context
	cancel          context.CancelFunc

	// Statistics
	totalAttempts  atomic.Int64
	allowedRetries atomic.Int64
	deniedRetries  atomic.Int64
}

// RetryBudgetConfig holds configuration for retry budget
type RetryBudgetConfig struct {
	TokensPerSecond float64 // Rate at which tokens are refilled
	MaxTokens       int     // Maximum number of tokens (burst capacity)
	Enabled         bool    // Whether retry budget is enabled
}

// NewRetryBudget creates a new retry budget (deprecated, use NewRetryBudgetWithContext)
func NewRetryBudget(config RetryBudgetConfig, logger *libpack_logger.Logger) *RetryBudget {
	return NewRetryBudgetWithContext(context.Background(), config, logger)
}

// NewRetryBudgetWithContext creates a new retry budget with context for graceful shutdown
func NewRetryBudgetWithContext(ctx context.Context, config RetryBudgetConfig, logger *libpack_logger.Logger) *RetryBudget {
	budgetCtx, cancel := context.WithCancel(ctx)
	rb := &RetryBudget{
		tokensPerSecond: config.TokensPerSecond,
		maxTokens:       int64(config.MaxTokens),
		enabled:         config.Enabled,
		logger:          logger,
		ctx:             budgetCtx,
		cancel:          cancel,
	}

	// Initialize with full bucket
	rb.currentTokens.Store(rb.maxTokens)
	rb.lastRefill.Store(time.Now().UnixNano())

	// Start refill goroutine
	if rb.enabled {
		go rb.refillLoop()
	}

	return rb
}

// AllowRetry checks if a retry is allowed based on the current budget
func (rb *RetryBudget) AllowRetry() bool {
	rb.totalAttempts.Add(1)

	if !rb.enabled {
		rb.allowedRetries.Add(1)
		return true
	}

	// Try to consume a token
	for {
		current := rb.currentTokens.Load()
		if current <= 0 {
			rb.deniedRetries.Add(1)
			if rb.logger != nil {
				rb.logger.Debug(&libpack_logger.LogMessage{
					Message: "Retry denied: budget exhausted",
					Pairs: map[string]any{
						"current_tokens": current,
						"denied_count":   rb.deniedRetries.Load(),
					},
				})
			}
			return false
		}

		if rb.currentTokens.CompareAndSwap(current, current-1) {
			rb.allowedRetries.Add(1)
			return true
		}
	}
}

// refillLoop periodically refills tokens
func (rb *RetryBudget) refillLoop() {
	ticker := time.NewTicker(100 * time.Millisecond) // Refill every 100ms
	defer ticker.Stop()

	for {
		select {
		case <-rb.ctx.Done():
			return
		case <-ticker.C:
			rb.refill()
		}
	}
}

// Shutdown stops the retry budget goroutine
func (rb *RetryBudget) Shutdown() {
	if rb.cancel != nil {
		rb.cancel()
	}
}

// refill adds tokens to the bucket based on elapsed time
func (rb *RetryBudget) refill() {
	now := time.Now().UnixNano()
	last := rb.lastRefill.Load()

	// Calculate elapsed time in seconds
	elapsed := float64(now-last) / float64(time.Second)

	// Calculate tokens to add
	tokensToAdd := int64(elapsed * rb.tokensPerSecond)

	if tokensToAdd > 0 {
		// Update last refill time
		if rb.lastRefill.CompareAndSwap(last, now) {
			// Add tokens, capped at maxTokens
			for {
				current := rb.currentTokens.Load()
				newValue := current + tokensToAdd
				if newValue > rb.maxTokens {
					newValue = rb.maxTokens
				}

				if rb.currentTokens.CompareAndSwap(current, newValue) {
					break
				}
			}
		}
	}
}

// GetStats returns current statistics
func (rb *RetryBudget) GetStats() map[string]any {
	totalAttempts := rb.totalAttempts.Load()
	allowedRetries := rb.allowedRetries.Load()
	deniedRetries := rb.deniedRetries.Load()

	var denialRate float64
	if totalAttempts > 0 {
		denialRate = float64(deniedRetries) / float64(totalAttempts) * 100
	}

	return map[string]any{
		"enabled":         rb.enabled,
		"current_tokens":  rb.currentTokens.Load(),
		"max_tokens":      rb.maxTokens,
		"tokens_per_sec":  rb.tokensPerSecond,
		"total_attempts":  totalAttempts,
		"allowed_retries": allowedRetries,
		"denied_retries":  deniedRetries,
		"denial_rate_pct": denialRate,
	}
}

// Reset resets the retry budget statistics
func (rb *RetryBudget) Reset() {
	rb.totalAttempts.Store(0)
	rb.allowedRetries.Store(0)
	rb.deniedRetries.Store(0)
	rb.currentTokens.Store(rb.maxTokens)
}

// UpdateConfig updates the retry budget configuration
func (rb *RetryBudget) UpdateConfig(config RetryBudgetConfig) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	rb.tokensPerSecond = config.TokensPerSecond
	rb.maxTokens = int64(config.MaxTokens)
	rb.enabled = config.Enabled

	// Reset to full capacity
	rb.currentTokens.Store(rb.maxTokens)

	if rb.logger != nil {
		rb.logger.Info(&libpack_logger.LogMessage{
			Message: "Retry budget configuration updated",
			Pairs: map[string]any{
				"tokens_per_sec": config.TokensPerSecond,
				"max_tokens":     config.MaxTokens,
				"enabled":        config.Enabled,
			},
		})
	}
}

// Global retry budget instance
var (
	retryBudget     *RetryBudget
	retryBudgetOnce sync.Once
)

// InitializeRetryBudget initializes the global retry budget (deprecated, use InitializeRetryBudgetWithContext)
func InitializeRetryBudget(config RetryBudgetConfig, logger *libpack_logger.Logger) *RetryBudget {
	return InitializeRetryBudgetWithContext(context.Background(), config, logger)
}

// InitializeRetryBudgetWithContext initializes the global retry budget with context for graceful shutdown
func InitializeRetryBudgetWithContext(ctx context.Context, config RetryBudgetConfig, logger *libpack_logger.Logger) *RetryBudget {
	retryBudgetOnce.Do(func() {
		retryBudget = NewRetryBudgetWithContext(ctx, config, logger)
		if logger != nil && config.Enabled {
			logger.Info(&libpack_logger.LogMessage{
				Message: "Retry budget initialized",
				Pairs: map[string]any{
					"tokens_per_sec": config.TokensPerSecond,
					"max_tokens":     config.MaxTokens,
				},
			})
		}
	})
	return retryBudget
}

// GetRetryBudget returns the global retry budget instance
func GetRetryBudget() *RetryBudget {
	return retryBudget
}
