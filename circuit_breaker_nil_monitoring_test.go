package main

import (
	"bytes"
	"testing"

	libpack_logger "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
	"github.com/stretchr/testify/require"
)

// TestInitCircuitBreaker_NilMonitoring_NoPanic is a regression test for a boot
// panic: enabling the circuit breaker (ENABLE_CIRCUIT_BREAKER=true) crashed at
// startup because parseConfig calls initCircuitBreaker before the monitoring
// server assigns cfg.Monitoring (StartMonitoringServer runs later). The metrics
// gauge registration then dereferenced a nil *MetricsSetup.
//
// The breaker (and the admin dashboard, which reads gobreaker state directly)
// must initialise without panicking even when monitoring is not yet available.
func TestInitCircuitBreaker_NilMonitoring_NoPanic(t *testing.T) {
	origCfg := cfg
	cbMutex.Lock()
	origCb, origMetrics := cb, cbMetrics
	cb, cbMetrics = nil, nil
	cbMutex.Unlock()
	t.Cleanup(func() {
		cbMutex.Lock()
		cb, cbMetrics = origCb, origMetrics
		cbMutex.Unlock()
		cfg = origCfg
	})

	cfg = &config{}
	cfg.Logger = libpack_logger.New().SetOutput(&bytes.Buffer{})
	cfg.Monitoring = nil // the production state when parseConfig runs
	cfg.CircuitBreaker.Enable = true
	cfg.CircuitBreaker.MaxFailures = 3
	cfg.CircuitBreaker.Timeout = 5
	cfg.CircuitBreaker.MaxRequestsInHalfOpen = 2

	require.NotPanics(t, func() { initCircuitBreaker(cfg) },
		"initCircuitBreaker must not panic when cfg.Monitoring is nil")
	require.NotNil(t, cb, "circuit breaker must initialise even without monitoring")
	require.NotNil(t, cbMetrics, "circuit breaker metrics manager must be created")
}
