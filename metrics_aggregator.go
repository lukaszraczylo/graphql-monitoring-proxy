package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	libpack_logger "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
	"github.com/redis/go-redis/v9"
)

// MetricsAggregator handles distributed metrics collection via Redis
type MetricsAggregator struct {
	redisClient  *redis.Client
	logger       *libpack_logger.Logger
	instanceID   string
	publishKey   string
	ttl          time.Duration
	publishTimer *time.Ticker
	ctx          context.Context
	cancel       context.CancelFunc
	mu           sync.RWMutex
}

// InstanceMetrics represents metrics for a single proxy instance
type InstanceMetrics struct {
	InstanceID     string                 `json:"instance_id"`
	Hostname       string                 `json:"hostname"`
	LastUpdate     time.Time              `json:"last_update"`
	UptimeSeconds  float64                `json:"uptime_seconds"`
	Stats          map[string]interface{} `json:"stats"`
	Cache          map[string]interface{} `json:"cache,omitempty"`         // Full cache details including memory
	CacheSummary   map[string]interface{} `json:"cache_summary,omitempty"` // Deprecated: kept for compatibility
	Health         map[string]interface{} `json:"health"`
	CircuitBreaker map[string]interface{} `json:"circuit_breaker,omitempty"`
	RetryBudget    map[string]interface{} `json:"retry_budget,omitempty"`
	Coalescing     map[string]interface{} `json:"coalescing,omitempty"`
	WebSocketStats map[string]interface{} `json:"websocket,omitempty"`
	Connections    map[string]interface{} `json:"connections,omitempty"`
}

// AggregatedMetrics represents combined metrics from all instances
type AggregatedMetrics struct {
	TotalInstances   int                        `json:"total_instances"`
	HealthyInstances int                        `json:"healthy_instances"`
	LastUpdate       time.Time                  `json:"last_update"`
	CombinedStats    map[string]interface{}     `json:"combined_stats"`
	Instances        []InstanceMetrics          `json:"instances"`
	PerInstanceStats map[string]InstanceMetrics `json:"per_instance_stats"`
}

var (
	metricsAggregator *MetricsAggregator
	aggregatorMutex   sync.RWMutex
)

// InitializeMetricsAggregator creates and starts the metrics aggregator
func InitializeMetricsAggregator(redisURL, redisPassword string, redisDB int, logger *libpack_logger.Logger) error {
	aggregatorMutex.Lock()
	defer aggregatorMutex.Unlock()

	if metricsAggregator != nil {
		return nil // Already initialized
	}

	// Parse Redis URL
	opt, err := redis.ParseURL(fmt.Sprintf("redis://%s/%d", redisURL, redisDB))
	if err != nil {
		return fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	if redisPassword != "" {
		opt.Password = redisPassword
	}

	// Create Redis client with connection timeouts
	opt.DialTimeout = 2 * time.Second
	opt.ReadTimeout = 2 * time.Second
	opt.WriteTimeout = 2 * time.Second
	opt.PoolTimeout = 3 * time.Second
	opt.MaxRetries = 2

	client := redis.NewClient(opt)

	// Test connection with detailed error reporting
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		// Log detailed connection error
		if logger != nil {
			logger.Error(&libpack_logger.LogMessage{
				Message: "❌ CRITICAL: Redis connection test FAILED during initialization",
				Pairs: map[string]interface{}{
					"error":        err.Error(),
					"redis_url":    redisURL,
					"redis_db":     redisDB,
					"has_password": redisPassword != "",
				},
			})
		}
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	// Log successful connection
	if logger != nil {
		logger.Info(&libpack_logger.LogMessage{
			Message: "✓ Redis connection test PASSED",
			Pairs: map[string]interface{}{
				"redis_url": redisURL,
				"redis_db":  redisDB,
			},
		})
	}

	// Generate unique instance ID (hostname + UUID for uniqueness)
	hostname, _ := os.Hostname()
	if hostname == "" {
		hostname = "unknown"
	}
	instanceID := fmt.Sprintf("%s-%s", hostname, uuid.New().String()[:8])

	ctx, cancel = context.WithCancel(context.Background())

	aggregator := &MetricsAggregator{
		redisClient:  client,
		logger:       logger,
		instanceID:   instanceID,
		publishKey:   "graphql-proxy:metrics:instances",
		ttl:          30 * time.Second, // Metrics expire after 30s if not updated
		publishTimer: time.NewTicker(5 * time.Second),
		ctx:          ctx,
		cancel:       cancel,
	}

	metricsAggregator = aggregator

	// Start publishing metrics
	go aggregator.startPublishing()

	if logger != nil {
		logger.Info(&libpack_logger.LogMessage{
			Message: "Metrics aggregator initialized",
			Pairs: map[string]interface{}{
				"instance_id": instanceID,
				"redis_url":   redisURL,
				"publish_key": aggregator.publishKey,
			},
		})
	}

	return nil
}

// GetMetricsAggregator returns the singleton instance
func GetMetricsAggregator() *MetricsAggregator {
	aggregatorMutex.RLock()
	defer aggregatorMutex.RUnlock()
	return metricsAggregator
}

// startPublishing periodically publishes metrics to Redis
func (ma *MetricsAggregator) startPublishing() {
	defer ma.publishTimer.Stop()

	// Publish immediately on start
	ma.publishMetrics()

	for {
		select {
		case <-ma.ctx.Done():
			// Clean up our metrics on shutdown
			ma.removeInstanceMetrics()
			return
		case <-ma.publishTimer.C:
			ma.publishMetrics()
		}
	}
}

// publishMetrics collects current metrics and stores them in Redis
// Note: This is exported for testing/debugging via admin API
func (ma *MetricsAggregator) publishMetrics() {
	// Defensive: check if aggregator is still valid
	if ma == nil {
		return
	}

	ma.mu.RLock()
	defer ma.mu.RUnlock()

	// Safety check: ensure global config is initialized
	if cfg == nil {
		if ma.logger != nil {
			ma.logger.Warning(&libpack_logger.LogMessage{
				Message: "Cannot publish metrics - global config not initialized yet",
				Pairs: map[string]interface{}{
					"instance_id": ma.instanceID,
				},
			})
		}
		return
	}

	// Gather all stats using the admin dashboard's method
	dashboard := NewAdminDashboard(ma.logger)
	allStats := dashboard.gatherAllStats()

	if len(allStats) == 0 {
		if ma.logger != nil {
			ma.logger.Warning(&libpack_logger.LogMessage{
				Message: "gatherAllStats returned empty/nil result",
				Pairs: map[string]interface{}{
					"instance_id": ma.instanceID,
				},
			})
		}
		return
	}

	// Create instance metrics
	hostname, _ := os.Hostname()
	if hostname == "" {
		hostname = "unknown"
	}

	metrics := InstanceMetrics{
		InstanceID:    ma.instanceID,
		Hostname:      hostname,
		LastUpdate:    time.Now(),
		UptimeSeconds: time.Since(startTime).Seconds(),
	}

	// Extract specific sections - CRITICAL: we must set the correct structure
	// Stats should contain the inner stats object with requests, cache_summary, etc.
	if stats, ok := allStats["stats"].(map[string]interface{}); ok {
		metrics.Stats = stats

		// Also extract cache summary separately for easier access (deprecated but kept for compatibility)
		if cacheSummary, ok := stats["cache_summary"].(map[string]interface{}); ok {
			metrics.CacheSummary = cacheSummary
		}

	} else {
		// Fallback: if stats extraction fails, use empty map
		if ma.logger != nil {
			ma.logger.Error(&libpack_logger.LogMessage{
				Message: "Failed to extract stats from allStats - using empty stats",
				Pairs: map[string]interface{}{
					"instance_id": ma.instanceID,
					"allStats_keys": func() []string {
						keys := make([]string, 0, len(allStats))
						for k := range allStats {
							keys = append(keys, k)
						}
						return keys
					}(),
				},
			})
		}
		metrics.Stats = make(map[string]interface{})
	}

	// Extract full cache details (includes memory usage)
	if cache, ok := allStats["cache"].(map[string]interface{}); ok {
		metrics.Cache = cache
	}

	if health, ok := allStats["health"].(map[string]interface{}); ok {
		metrics.Health = health
	} else {
		metrics.Health = make(map[string]interface{})
	}
	if cb, ok := allStats["circuit_breaker"].(map[string]interface{}); ok {
		metrics.CircuitBreaker = cb
	}
	if rb, ok := allStats["retry_budget"].(map[string]interface{}); ok {
		metrics.RetryBudget = rb
	}
	if coal, ok := allStats["coalescing"].(map[string]interface{}); ok {
		metrics.Coalescing = coal
	}
	if ws, ok := allStats["websocket"].(map[string]interface{}); ok {
		metrics.WebSocketStats = ws
	}
	if conn, ok := allStats["connections"].(map[string]interface{}); ok {
		metrics.Connections = conn
	}

	// Marshal to JSON
	data, err := json.Marshal(metrics)
	if err != nil {
		if ma.logger != nil {
			ma.logger.Error(&libpack_logger.LogMessage{
				Message: "Failed to marshal metrics for Redis",
				Pairs:   map[string]interface{}{"error": err.Error()},
			})
		}
		return
	}

	// Store in Redis hash with TTL
	key := fmt.Sprintf("%s:%s", ma.publishKey, ma.instanceID)

	// Create a fresh context with timeout to avoid inheriting cancelled parent context
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	pipe := ma.redisClient.Pipeline()
	pipe.Set(ctx, key, data, ma.ttl)
	pipe.SAdd(ctx, ma.publishKey, ma.instanceID)
	pipe.Expire(ctx, ma.publishKey, ma.ttl*2) // Keep set alive

	_, err = pipe.Exec(ctx)
	if err != nil {
		if ma.logger != nil {
			ma.logger.Error(&libpack_logger.LogMessage{
				Message: "❌ CRITICAL: Failed to publish metrics to Redis - cluster mode will not work!",
				Pairs: map[string]interface{}{
					"error":       err.Error(),
					"instance_id": ma.instanceID,
					"key":         key,
					"redis_key":   ma.publishKey,
				},
			})
		}
		return
	}
}

// removeInstanceMetrics cleans up metrics from Redis on shutdown
func (ma *MetricsAggregator) removeInstanceMetrics() {
	// Create a fresh context with timeout for cleanup
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	key := fmt.Sprintf("%s:%s", ma.publishKey, ma.instanceID)
	pipe := ma.redisClient.Pipeline()
	pipe.Del(ctx, key)
	pipe.SRem(ctx, ma.publishKey, ma.instanceID)
	_, err := pipe.Exec(ctx)

	if err != nil && ma.logger != nil {
		ma.logger.Warning(&libpack_logger.LogMessage{
			Message: "Failed to remove instance metrics from Redis during shutdown",
			Pairs:   map[string]interface{}{"instance_id": ma.instanceID, "error": err.Error()},
		})
		return
	}

	if ma.logger != nil {
		ma.logger.Info(&libpack_logger.LogMessage{
			Message: "Removed instance metrics from Redis",
			Pairs:   map[string]interface{}{"instance_id": ma.instanceID},
		})
	}
}

// GetAggregatedMetrics retrieves and aggregates metrics from all instances
func (ma *MetricsAggregator) GetAggregatedMetrics() (*AggregatedMetrics, error) {
	// Create a fresh context with timeout to avoid inheriting cancelled parent context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Get all instance IDs
	instanceIDs, err := ma.redisClient.SMembers(ctx, ma.publishKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get instance list: %w", err)
	}

	if len(instanceIDs) == 0 {
		return &AggregatedMetrics{
			TotalInstances:   0,
			HealthyInstances: 0,
			LastUpdate:       time.Now(),
			CombinedStats:    make(map[string]interface{}),
			Instances:        []InstanceMetrics{},
			PerInstanceStats: make(map[string]InstanceMetrics),
		}, nil
	}

	// Fetch metrics for all instances
	pipe := ma.redisClient.Pipeline()
	cmds := make([]*redis.StringCmd, len(instanceIDs))
	for i, instanceID := range instanceIDs {
		key := fmt.Sprintf("%s:%s", ma.publishKey, instanceID)
		cmds[i] = pipe.Get(ctx, key)
	}
	pipe.Exec(ctx)

	// Parse metrics
	instances := make([]InstanceMetrics, 0, len(instanceIDs))
	perInstance := make(map[string]InstanceMetrics)
	healthyCount := 0
	staleCount := 0
	errorCount := 0

	for i, cmd := range cmds {
		data, err := cmd.Result()
		if err != nil {
			errorCount++
			// Clean up stale instance ID from the set
			if err == redis.Nil {
				staleCount++
				// Remove stale instance from set in background
				go func(instID string) {
					cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 2*time.Second)
					defer cleanupCancel()
					ma.redisClient.SRem(cleanupCtx, ma.publishKey, instID)
				}(instanceIDs[i])
			}
			continue
		}

		var metrics InstanceMetrics
		if err := json.Unmarshal([]byte(data), &metrics); err != nil {
			if ma.logger != nil {
				ma.logger.Warning(&libpack_logger.LogMessage{
					Message: "Failed to unmarshal instance metrics",
					Pairs:   map[string]interface{}{"error": err.Error()},
				})
			}
			continue
		}

		// Check if instance is stale (not updated in 1 minute)
		instanceAge := time.Since(metrics.LastUpdate)
		if instanceAge > 1*time.Minute {
			staleCount++
			// Clean up stale instance from set in background
			go func(instID string, age time.Duration) {
				cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 2*time.Second)
				defer cleanupCancel()
				ma.redisClient.SRem(cleanupCtx, ma.publishKey, instID)
				if ma.logger != nil {
					ma.logger.Info(&libpack_logger.LogMessage{
						Message: "Removed inactive instance",
						Pairs: map[string]interface{}{
							"instance_id":      instID,
							"inactive_seconds": age.Seconds(),
						},
					})
				}
			}(instanceIDs[i], instanceAge)
			continue // Skip stale instances
		}

		instances = append(instances, metrics)
		perInstance[metrics.InstanceID] = metrics

		// Count healthy instances
		if health, ok := metrics.Health["status"].(string); ok && health == "healthy" {
			healthyCount++
		}
	}

	// Log cleanup stats if we found stale instances
	if ma.logger != nil && (staleCount > 0 || errorCount > 0) {
		ma.logger.Info(&libpack_logger.LogMessage{
			Message: "Cleaned up stale instance IDs from Redis",
			Pairs: map[string]interface{}{
				"total_in_set":    len(instanceIDs),
				"valid_instances": len(instances),
				"stale_cleaned":   staleCount,
				"errors":          errorCount,
			},
		})
	}

	// Aggregate statistics
	aggregated := &AggregatedMetrics{
		TotalInstances:   len(instances),
		HealthyInstances: healthyCount,
		LastUpdate:       time.Now(),
		CombinedStats:    ma.aggregateStats(instances),
		Instances:        instances,
		PerInstanceStats: perInstance,
	}

	return aggregated, nil
}

// aggregateStats combines statistics from multiple instances
func (ma *MetricsAggregator) aggregateStats(instances []InstanceMetrics) map[string]interface{} {
	if len(instances) == 0 {
		if ma.logger != nil {
			ma.logger.Warning(&libpack_logger.LogMessage{
				Message: "No instances to aggregate",
			})
		}
		return make(map[string]interface{})
	}

	// Initialize aggregated values
	var (
		totalRequests          int64
		totalSucceeded         int64
		totalFailed            int64
		totalSkipped           int64
		totalCacheHits         int64
		totalCacheMisses       int64
		totalCachedQueries     int64
		totalMemoryUsageMB     float64
		hasValidMemoryStats    bool // Track if any instance has valid memory stats
		totalCurrentRPS        float64
		totalAvgRPS            float64
		totalActiveConnections int64
		totalWSConnections     int64
		totalCoalescedRequests int64
		totalPrimaryRequests   int64
		oldestUptime           float64

		// Retry budget stats
		totalRetryAllowed  int64
		totalRetryDenied   int64
		totalRetryAttempts int64
		retryBudgetEnabled = false

		// Circuit breaker stats
		cbOpenCount           int
		cbHalfOpenCount       int
		cbClosedCount         int
		circuitBreakerEnabled = false
	)

	for idx, instance := range instances {
		// Track oldest uptime for cluster uptime
		if oldestUptime == 0 || instance.UptimeSeconds < oldestUptime {
			oldestUptime = instance.UptimeSeconds
		}

		// Aggregate request stats
		if instance.Stats == nil {
			if ma.logger != nil {
				ma.logger.Warning(&libpack_logger.LogMessage{
					Message: "Instance has nil Stats",
					Pairs: map[string]interface{}{
						"instance_id": instance.InstanceID,
						"index":       idx,
					},
				})
			}
			continue
		}

		if stats, ok := instance.Stats["requests"].(map[string]interface{}); ok {
			if total, ok := stats["total"].(float64); ok {
				totalRequests += int64(total)
			}
			if succeeded, ok := stats["succeeded"].(float64); ok {
				totalSucceeded += int64(succeeded)
			}
			if failed, ok := stats["failed"].(float64); ok {
				totalFailed += int64(failed)
			}
			if skipped, ok := stats["skipped"].(float64); ok {
				totalSkipped += int64(skipped)
			}
			if currentRPS, ok := stats["current_requests_per_second"].(float64); ok {
				totalCurrentRPS += currentRPS
			}
			if avgRPS, ok := stats["avg_requests_per_second"].(float64); ok {
				totalAvgRPS += avgRPS
			}
		} else {
			if ma.logger != nil {
				// Log what keys are actually in Stats for debugging
				keys := make([]string, 0, len(instance.Stats))
				for k := range instance.Stats {
					keys = append(keys, k)
				}
				ma.logger.Warning(&libpack_logger.LogMessage{
					Message: "Instance Stats missing 'requests' key",
					Pairs: map[string]interface{}{
						"instance_id": instance.InstanceID,
						"stats_keys":  keys,
						"index":       idx,
					},
				})
			}
		}

		// Aggregate cache stats from CacheSummary (backward compatibility)
		if len(instance.CacheSummary) > 0 {
			if hits, ok := instance.CacheSummary["hits"].(float64); ok {
				totalCacheHits += int64(hits)
			}
			if misses, ok := instance.CacheSummary["misses"].(float64); ok {
				totalCacheMisses += int64(misses)
			}
			if cached, ok := instance.CacheSummary["total_cached"].(float64); ok {
				totalCachedQueries += int64(cached)
			}
		}

		// Aggregate memory usage from full cache details
		// Skip -1 values which indicate Redis cache (memory tracking not available)
		if len(instance.Cache) > 0 {
			if memMB, ok := instance.Cache["memory_usage_mb"].(float64); ok && memMB >= 0 {
				totalMemoryUsageMB += memMB
				hasValidMemoryStats = true
			}
		}

		// Aggregate connection stats
		if len(instance.Connections) > 0 {
			if active, ok := instance.Connections["active_connections"].(float64); ok {
				totalActiveConnections += int64(active)
			}
		}

		// Aggregate WebSocket connections
		if len(instance.WebSocketStats) > 0 {
			if active, ok := instance.WebSocketStats["active_connections"].(float64); ok {
				totalWSConnections += int64(active)
			}
		}

		// Aggregate coalescing stats
		if len(instance.Coalescing) > 0 {
			if coalesced, ok := instance.Coalescing["coalesced_requests"].(float64); ok {
				totalCoalescedRequests += int64(coalesced)
			}
			if primary, ok := instance.Coalescing["primary_requests"].(float64); ok {
				totalPrimaryRequests += int64(primary)
			}
		}

		// Aggregate retry budget stats
		if len(instance.RetryBudget) > 0 {
			if enabled, ok := instance.RetryBudget["enabled"].(bool); ok && enabled {
				retryBudgetEnabled = true
			}
			if allowed, ok := instance.RetryBudget["allowed_retries"].(float64); ok {
				totalRetryAllowed += int64(allowed)
			}
			if denied, ok := instance.RetryBudget["denied_retries"].(float64); ok {
				totalRetryDenied += int64(denied)
			}
			if attempts, ok := instance.RetryBudget["total_attempts"].(float64); ok {
				totalRetryAttempts += int64(attempts)
			}
		}

		// Aggregate circuit breaker stats
		if len(instance.CircuitBreaker) > 0 {
			if enabled, ok := instance.CircuitBreaker["enabled"].(bool); ok && enabled {
				circuitBreakerEnabled = true
			}
			if state, ok := instance.CircuitBreaker["state"].(string); ok {
				switch state {
				case "open":
					cbOpenCount++
				case "half-open":
					cbHalfOpenCount++
				case "closed":
					cbClosedCount++
				}
			}
		}
	}

	// Calculate derived metrics
	successRate := 0.0
	if totalRequests > 0 {
		successRate = float64(totalSucceeded) / float64(totalRequests) * 100
	}

	cacheHitRate := 0.0
	totalCacheRequests := totalCacheHits + totalCacheMisses
	if totalCacheRequests > 0 {
		cacheHitRate = float64(totalCacheHits) / float64(totalCacheRequests) * 100
	}

	backendSavings := 0.0
	totalCoalRequests := totalCoalescedRequests + totalPrimaryRequests
	if totalCoalRequests > 0 {
		backendSavings = float64(totalCoalescedRequests) / float64(totalCoalRequests) * 100
	}

	// Calculate retry budget denial rate
	retryDenialRate := 0.0
	if totalRetryAttempts > 0 {
		retryDenialRate = float64(totalRetryDenied) / float64(totalRetryAttempts) * 100
	}

	// Determine overall circuit breaker state
	cbState := "unknown"
	if circuitBreakerEnabled {
		if cbOpenCount > 0 {
			cbState = "open" // If any instance is open, cluster is in degraded state
		} else if cbHalfOpenCount > 0 {
			cbState = "half-open"
		} else if cbClosedCount > 0 {
			cbState = "closed"
		}
	}

	result := map[string]interface{}{
		"cluster_mode":    true,
		"total_instances": len(instances),
		"cluster_uptime":  oldestUptime,
		"requests": map[string]interface{}{
			"total":                       totalRequests,
			"succeeded":                   totalSucceeded,
			"failed":                      totalFailed,
			"skipped":                     totalSkipped,
			"success_rate_pct":            successRate,
			"current_requests_per_second": totalCurrentRPS,
			"avg_requests_per_second":     totalAvgRPS,
		},
		"cache_summary": map[string]interface{}{
			"hits":         totalCacheHits,
			"misses":       totalCacheMisses,
			"hit_rate_pct": cacheHitRate,
			"total_cached": totalCachedQueries,
		},
		"memory": map[string]interface{}{
			"total_usage_mb": func() float64 {
				if hasValidMemoryStats {
					return totalMemoryUsageMB
				}
				return -1
			}(),
			"available": hasValidMemoryStats,
		},
		"connections": map[string]interface{}{
			"total_active": totalActiveConnections,
		},
		"websocket": map[string]interface{}{
			"total_connections": totalWSConnections,
		},
		"coalescing": map[string]interface{}{
			"enabled":                  len(instances) > 0, // enabled if we have instances with data
			"total_coalesced_requests": totalCoalescedRequests,
			"total_primary_requests":   totalPrimaryRequests,
			"backend_savings_pct":      backendSavings,
			"coalescing_rate_pct":      backendSavings,
		},
		"retry_budget": map[string]interface{}{
			"enabled":         retryBudgetEnabled,
			"allowed_retries": totalRetryAllowed,
			"denied_retries":  totalRetryDenied,
			"total_attempts":  totalRetryAttempts,
			"denial_rate_pct": retryDenialRate,
		},
		"circuit_breaker": map[string]interface{}{
			"enabled":            circuitBreakerEnabled,
			"state":              cbState,
			"instances_open":     cbOpenCount,
			"instances_closed":   cbClosedCount,
			"instances_halfopen": cbHalfOpenCount,
		},
	}

	return result
}

// Shutdown stops the metrics aggregator
func (ma *MetricsAggregator) Shutdown() {
	ma.mu.Lock()
	defer ma.mu.Unlock()

	if ma.cancel != nil {
		ma.cancel()
	}

	if ma.redisClient != nil {
		ma.redisClient.Close()
	}

	if ma.logger != nil {
		ma.logger.Info(&libpack_logger.LogMessage{
			Message: "Metrics aggregator shut down",
		})
	}
}

// GetInstanceID returns the current instance ID
func (ma *MetricsAggregator) GetInstanceID() string {
	return ma.instanceID
}

// IsClusterMode returns true if multiple instances are detected
func (ma *MetricsAggregator) IsClusterMode() bool {
	// Create a fresh context with timeout to avoid inheriting cancelled parent context
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	count, err := ma.redisClient.SCard(ctx, ma.publishKey).Result()
	if err != nil {
		return false
	}

	return count > 1
}

// GetInstanceHostname returns a human-readable instance identifier
func GetInstanceHostname() string {
	hostname, _ := os.Hostname()
	if hostname == "" {
		hostname = "unknown"
	}
	// Remove domain suffix for cleaner display
	if idx := strings.Index(hostname, "."); idx > 0 {
		hostname = hostname[:idx]
	}
	return hostname
}
