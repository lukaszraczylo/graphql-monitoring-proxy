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
	CacheSummary   map[string]interface{} `json:"cache_summary,omitempty"`
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

	// Create Redis client
	client := redis.NewClient(opt)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
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
func (ma *MetricsAggregator) publishMetrics() {
	ma.mu.RLock()
	defer ma.mu.RUnlock()

	// Gather all stats using the admin dashboard's method
	dashboard := NewAdminDashboard(ma.logger)
	allStats := dashboard.gatherAllStats()

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
		Stats:         allStats,
	}

	// Extract specific sections for easier aggregation
	if stats, ok := allStats["stats"].(map[string]interface{}); ok {
		metrics.Stats = stats
		if cacheSummary, ok := stats["cache_summary"].(map[string]interface{}); ok {
			metrics.CacheSummary = cacheSummary
		}
	}
	if health, ok := allStats["health"].(map[string]interface{}); ok {
		metrics.Health = health
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
	ctx, cancel := context.WithTimeout(ma.ctx, 2*time.Second)
	defer cancel()

	pipe := ma.redisClient.Pipeline()
	pipe.Set(ctx, key, data, ma.ttl)
	pipe.SAdd(ctx, ma.publishKey, ma.instanceID)
	pipe.Expire(ctx, ma.publishKey, ma.ttl*2) // Keep set alive

	if _, err := pipe.Exec(ctx); err != nil {
		if ma.logger != nil {
			ma.logger.Warning(&libpack_logger.LogMessage{
				Message: "Failed to publish metrics to Redis",
				Pairs:   map[string]interface{}{"error": err.Error()},
			})
		}
	}
}

// removeInstanceMetrics cleans up metrics from Redis on shutdown
func (ma *MetricsAggregator) removeInstanceMetrics() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	key := fmt.Sprintf("%s:%s", ma.publishKey, ma.instanceID)
	pipe := ma.redisClient.Pipeline()
	pipe.Del(ctx, key)
	pipe.SRem(ctx, ma.publishKey, ma.instanceID)
	pipe.Exec(ctx)

	if ma.logger != nil {
		ma.logger.Info(&libpack_logger.LogMessage{
			Message: "Removed instance metrics from Redis",
			Pairs:   map[string]interface{}{"instance_id": ma.instanceID},
		})
	}
}

// GetAggregatedMetrics retrieves and aggregates metrics from all instances
func (ma *MetricsAggregator) GetAggregatedMetrics() (*AggregatedMetrics, error) {
	ctx, cancel := context.WithTimeout(ma.ctx, 5*time.Second)
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

	for i, cmd := range cmds {
		data, err := cmd.Result()
		if err != nil {
			if ma.logger != nil {
				ma.logger.Debug(&libpack_logger.LogMessage{
					Message: "Failed to get instance metrics",
					Pairs: map[string]interface{}{
						"instance_id": instanceIDs[i],
						"error":       err.Error(),
					},
				})
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

		// Check if instance is stale (not updated in 2x TTL)
		if time.Since(metrics.LastUpdate) > ma.ttl*2 {
			continue // Skip stale instances
		}

		instances = append(instances, metrics)
		perInstance[metrics.InstanceID] = metrics

		// Count healthy instances
		if health, ok := metrics.Health["status"].(string); ok && health == "healthy" {
			healthyCount++
		}
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
		totalCurrentRPS        float64
		totalAvgRPS            float64
		totalActiveConnections int64
		totalWSConnections     int64
		totalCoalescedRequests int64
		totalPrimaryRequests   int64
		oldestUptime           float64
	)

	for _, instance := range instances {
		// Track oldest uptime for cluster uptime
		if oldestUptime == 0 || instance.UptimeSeconds < oldestUptime {
			oldestUptime = instance.UptimeSeconds
		}

		// Aggregate request stats
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
		}

		// Aggregate cache stats
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

	return map[string]interface{}{
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
		"connections": map[string]interface{}{
			"total_active": totalActiveConnections,
		},
		"websocket": map[string]interface{}{
			"total_connections": totalWSConnections,
		},
		"coalescing": map[string]interface{}{
			"total_coalesced_requests": totalCoalescedRequests,
			"total_primary_requests":   totalPrimaryRequests,
			"backend_savings_pct":      backendSavings,
		},
	}
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
	ctx, cancel := context.WithTimeout(ma.ctx, 2*time.Second)
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
