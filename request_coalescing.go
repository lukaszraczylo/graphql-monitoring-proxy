package main

import (
	"sync"
	"sync/atomic"
	"time"

	libpack_logger "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
	libpack_monitoring "github.com/lukaszraczylo/graphql-monitoring-proxy/monitoring"
)

// CoalescedResponse represents the shared response
type CoalescedResponse struct {
	Body       []byte
	StatusCode int
	Headers    map[string]string
	Err        error
	CachedAt   time.Time
}

// RequestCoalescer implements the single-flight pattern to deduplicate identical concurrent requests
type RequestCoalescer struct {
	inflight   sync.Map // key: hash, value: *inflightRequest
	logger     *libpack_logger.Logger
	monitoring *libpack_monitoring.MetricsSetup
	enabled    bool

	// Statistics
	totalRequests     atomic.Int64
	coalescedRequests atomic.Int64
	inflightCount     atomic.Int64
}

// inflightRequest represents a request currently in flight
type inflightRequest struct {
	wg        sync.WaitGroup
	response  *CoalescedResponse
	waiters   atomic.Int32
	createdAt time.Time
	mu        sync.RWMutex
}

// NewRequestCoalescer creates a new request coalescer
func NewRequestCoalescer(enabled bool, logger *libpack_logger.Logger, monitoring *libpack_monitoring.MetricsSetup) *RequestCoalescer {
	rc := &RequestCoalescer{
		logger:     logger,
		monitoring: monitoring,
		enabled:    enabled,
	}

	if logger != nil && enabled {
		logger.Info(&libpack_logger.LogMessage{
			Message: "Request coalescing enabled",
		})
	}

	return rc
}

// Do executes a function, deduplicating concurrent calls with the same key
func (rc *RequestCoalescer) Do(key string, fn func() (*CoalescedResponse, error)) (*CoalescedResponse, error) {
	rc.totalRequests.Add(1)

	if !rc.enabled {
		return fn()
	}

	// Try to load existing inflight request
	if existing, loaded := rc.inflight.Load(key); loaded {
		inflight := existing.(*inflightRequest)

		// Increment waiter count
		waiters := inflight.waiters.Add(1)
		rc.coalescedRequests.Add(1)

		if rc.logger != nil {
			rc.logger.Debug(&libpack_logger.LogMessage{
				Message: "Request coalesced with in-flight request",
				Pairs: map[string]interface{}{
					"key":     key[:min(len(key), 32)] + "...",
					"waiters": waiters,
				},
			})
		}

		// Wait for the inflight request to complete
		inflight.wg.Wait()

		// Return the shared response
		inflight.mu.RLock()
		defer inflight.mu.RUnlock()

		if rc.monitoring != nil {
			rc.monitoring.Increment("graphql_proxy_coalesced_requests_total", nil)
		}

		return inflight.response, nil
	}

	// Create a new inflight request
	inflight := &inflightRequest{
		createdAt: time.Now(),
	}
	inflight.wg.Add(1)
	inflight.waiters.Store(1) // This request is the first waiter

	// Try to store it (another goroutine might have just done the same)
	actual, loaded := rc.inflight.LoadOrStore(key, inflight)
	if loaded {
		// Someone else beat us to it, wait for their result
		existingInflight := actual.(*inflightRequest)
		waiters := existingInflight.waiters.Add(1)
		rc.coalescedRequests.Add(1)

		if rc.logger != nil {
			rc.logger.Debug(&libpack_logger.LogMessage{
				Message: "Request coalesced (race condition)",
				Pairs: map[string]interface{}{
					"key":     key[:min(len(key), 32)] + "...",
					"waiters": waiters,
				},
			})
		}

		existingInflight.wg.Wait()

		existingInflight.mu.RLock()
		defer existingInflight.mu.RUnlock()

		if rc.monitoring != nil {
			rc.monitoring.Increment("graphql_proxy_coalesced_requests_total", nil)
		}

		return existingInflight.response, nil
	}

	// We're the primary request, execute the function
	rc.inflightCount.Add(1)
	defer rc.inflightCount.Add(-1)

	// Execute the request
	response, err := fn()

	// Store the result
	inflight.mu.Lock()
	if err != nil {
		inflight.response = &CoalescedResponse{
			Err: err,
		}
	} else {
		inflight.response = response
	}
	inflight.mu.Unlock()

	// Clean up and notify waiters
	rc.inflight.Delete(key)
	inflight.wg.Done()

	// Log statistics
	waiters := inflight.waiters.Load()
	duration := time.Since(inflight.createdAt)

	if rc.logger != nil && waiters > 1 {
		rc.logger.Info(&libpack_logger.LogMessage{
			Message: "Request completed, served coalesced waiters",
			Pairs: map[string]interface{}{
				"key":         key[:min(len(key), 32)] + "...",
				"waiters":     waiters,
				"duration_ms": duration.Milliseconds(),
				"saved_calls": waiters - 1,
			},
		})
	}

	if rc.monitoring != nil {
		rc.monitoring.Increment("graphql_proxy_primary_requests_total", nil)
		if waiters > 1 {
			rc.monitoring.Update("graphql_proxy_coalescing_wait_duration", nil, duration.Seconds())
		}
	}

	return inflight.response, nil
}

// GetStats returns coalescing statistics
func (rc *RequestCoalescer) GetStats() map[string]interface{} {
	totalRequests := rc.totalRequests.Load()
	coalescedRequests := rc.coalescedRequests.Load()

	var coalescingRate float64
	if totalRequests > 0 {
		coalescingRate = float64(coalescedRequests) / float64(totalRequests) * 100
	}

	primaryRequests := totalRequests - coalescedRequests

	var savings float64
	if primaryRequests > 0 {
		savings = float64(coalescedRequests) / float64(primaryRequests) * 100
	}

	return map[string]interface{}{
		"enabled":             rc.enabled,
		"total_requests":      totalRequests,
		"primary_requests":    primaryRequests,
		"coalesced_requests":  coalescedRequests,
		"inflight_count":      rc.inflightCount.Load(),
		"coalescing_rate_pct": coalescingRate,
		"backend_savings_pct": savings,
	}
}

// Reset resets coalescing statistics
func (rc *RequestCoalescer) Reset() {
	rc.totalRequests.Store(0)
	rc.coalescedRequests.Store(0)
}

// Global request coalescer
var (
	requestCoalescer     *RequestCoalescer
	requestCoalescerOnce sync.Once
)

// InitializeRequestCoalescer initializes the global request coalescer
func InitializeRequestCoalescer(enabled bool, logger *libpack_logger.Logger, monitoring *libpack_monitoring.MetricsSetup) *RequestCoalescer {
	requestCoalescerOnce.Do(func() {
		requestCoalescer = NewRequestCoalescer(enabled, logger, monitoring)
	})
	return requestCoalescer
}

// GetRequestCoalescer returns the global request coalescer
func GetRequestCoalescer() *RequestCoalescer {
	return requestCoalescer
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
