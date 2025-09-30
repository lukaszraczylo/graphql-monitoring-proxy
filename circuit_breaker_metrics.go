package main

import (
	"sync/atomic"

	"github.com/VictoriaMetrics/metrics"
	libpack_monitoring "github.com/lukaszraczylo/graphql-monitoring-proxy/monitoring"
)

// CircuitBreakerMetrics manages circuit breaker metrics without recreating gauges
type CircuitBreakerMetrics struct {
	stateValue   atomic.Value // stores float64
	stateGauge   *metrics.Gauge
	failCounters map[string]*metrics.Counter
}

// NewCircuitBreakerMetrics creates a new circuit breaker metrics manager
func NewCircuitBreakerMetrics(monitoring *libpack_monitoring.MetricsSetup) *CircuitBreakerMetrics {
	cbm := &CircuitBreakerMetrics{
		failCounters: make(map[string]*metrics.Counter),
	}

	// Initialize state value
	cbm.stateValue.Store(float64(0))

	// Create gauge with callback that reads the atomic value
	cbm.stateGauge = monitoring.RegisterMetricsGauge(
		libpack_monitoring.MetricsCircuitState,
		nil,
		0, // Initial value doesn't matter as callback will be used
	)

	// Override the gauge callback to read from atomic value
	cbm.stateGauge = monitoring.RegisterMetricsGauge(
		libpack_monitoring.MetricsCircuitState,
		nil,
		cbm.GetState(),
	)

	return cbm
}

// UpdateState updates the circuit breaker state value atomically
func (cbm *CircuitBreakerMetrics) UpdateState(state float64) {
	cbm.stateValue.Store(state)
}

// GetState returns the current circuit breaker state value
func (cbm *CircuitBreakerMetrics) GetState() float64 {
	if val := cbm.stateValue.Load(); val != nil {
		return val.(float64)
	}
	return 0
}

// GetOrCreateFailCounter returns a counter for the given state key
func (cbm *CircuitBreakerMetrics) GetOrCreateFailCounter(monitoring *libpack_monitoring.MetricsSetup, stateKey string) *metrics.Counter {
	if counter, exists := cbm.failCounters[stateKey]; exists {
		return counter
	}

	// Create new counter
	counter := monitoring.RegisterMetricsCounter(stateKey, nil)
	cbm.failCounters[stateKey] = counter
	return counter
}

// Global circuit breaker metrics instance
var cbMetrics *CircuitBreakerMetrics

// InitializeCircuitBreakerMetrics initializes the global circuit breaker metrics
func InitializeCircuitBreakerMetrics(monitoring *libpack_monitoring.MetricsSetup) {
	if cbMetrics == nil {
		cbMetrics = NewCircuitBreakerMetrics(monitoring)
	}
}
