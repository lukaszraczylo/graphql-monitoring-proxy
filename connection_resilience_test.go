package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	libpack_logger "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
	"github.com/stretchr/testify/suite"
)

// ConnectionResilienceTestSuite tests connection resilience features
type ConnectionResilienceTestSuite struct {
	suite.Suite
	originalConfig   *config
	outputBuffer     *bytes.Buffer
	mockServer       *httptest.Server
	mockServerCalled atomic.Int32
}

func (suite *ConnectionResilienceTestSuite) SetupTest() {
	// Store original config
	suite.originalConfig = cfg

	// Create a buffer to capture logger output
	suite.outputBuffer = &bytes.Buffer{}

	// Setup a new config with a real logger that writes to our buffer
	cfg = &config{}
	cfg.Logger = libpack_logger.New().SetOutput(suite.outputBuffer)

	// Reset call counter
	suite.mockServerCalled.Store(0)

	// Create a mock GraphQL server
	suite.mockServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		suite.mockServerCalled.Add(1)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":{"__typename":"Query"}}`))
	}))

	// Configure the test with mock server URL
	cfg.Server.HostGraphQL = suite.mockServer.URL
	cfg.Client.ClientTimeout = 5
	cfg.Client.MaxConnsPerHost = 10
	cfg.Client.MaxIdleConnDuration = 30
	cfg.Client.DisableTLSVerify = true

	// Create fasthttp client
	cfg.Client.FastProxyClient = createFasthttpClient(cfg)
}

func (suite *ConnectionResilienceTestSuite) TearDownTest() {
	// Close mock server
	if suite.mockServer != nil {
		suite.mockServer.Close()
	}

	// Clean up global instances with proper shutdown
	if backendHealthManager != nil {
		backendHealthManager.Shutdown()
		backendHealthManager = nil
	}

	if connectionPoolManager != nil {
		connectionPoolManager.Shutdown()
		connectionPoolManager = nil
	}

	// Restore original config
	cfg = suite.originalConfig
}

// TestBackendHealthManager tests the backend health monitoring
func (suite *ConnectionResilienceTestSuite) TestBackendHealthManager() {
	suite.Run("initialization", func() {
		healthMgr := NewBackendHealthManager(cfg.Client.FastProxyClient, cfg.Server.HostGraphQL, cfg.Logger)
		suite.NotNil(healthMgr)
		suite.Equal(cfg.Server.HostGraphQL, healthMgr.backendURL)
		suite.Equal(5*time.Second, healthMgr.checkInterval)
		suite.Equal(30, healthMgr.maxRetries)
	})

	suite.Run("health check success", func() {
		healthMgr := NewBackendHealthManager(cfg.Client.FastProxyClient, cfg.Server.HostGraphQL, cfg.Logger)
		isHealthy := healthMgr.checkBackendHealth()
		suite.True(isHealthy)
		suite.GreaterOrEqual(suite.mockServerCalled.Load(), int32(1))
	})

	suite.Run("health check failure", func() {
		// Use invalid URL to simulate failure
		healthMgr := NewBackendHealthManager(cfg.Client.FastProxyClient, "http://invalid-url:99999", cfg.Logger)
		isHealthy := healthMgr.checkBackendHealth()
		suite.False(isHealthy)
	})

	suite.Run("startup readiness with healthy backend", func() {
		healthMgr := NewBackendHealthManager(cfg.Client.FastProxyClient, cfg.Server.HostGraphQL, cfg.Logger)
		err := healthMgr.WaitForBackendReady(10 * time.Second)
		suite.NoError(err)
		suite.True(healthMgr.IsHealthy())
	})

	suite.Run("startup readiness timeout", func() {
		// Use invalid URL to simulate backend not ready
		healthMgr := NewBackendHealthManager(cfg.Client.FastProxyClient, "http://invalid-url:99999", cfg.Logger)
		err := healthMgr.WaitForBackendReady(2 * time.Second)
		suite.Error(err)
		suite.Contains(err.Error(), "did not become ready")
	})
}

// TestConnectionPoolManager tests the connection pool management
func (suite *ConnectionResilienceTestSuite) TestConnectionPoolManager() {
	suite.Run("initialization", func() {
		poolMgr := NewConnectionPoolManager(cfg.Client.FastProxyClient)
		suite.NotNil(poolMgr)
		suite.NotNil(poolMgr.client)
		suite.Equal(15*time.Second, poolMgr.keepAliveInterval)
		suite.Equal(30*time.Second, poolMgr.cleanupInterval)
		suite.Equal(60*time.Second, poolMgr.recoveryCheckInterval)
	})

	suite.Run("connection statistics", func() {
		poolMgr := NewConnectionPoolManager(cfg.Client.FastProxyClient)

		// Record some connections
		poolMgr.RecordConnectionSuccess()
		poolMgr.RecordConnectionSuccess()
		poolMgr.RecordConnectionFailure()

		stats := poolMgr.GetConnectionStats()
		suite.Equal(int64(2), stats["active_connections"])
		suite.Equal(int64(2), stats["total_connections"])
		suite.Equal(int64(1), stats["connection_failures"])
	})

	suite.Run("keep alive functionality", func() {
		poolMgr := NewConnectionPoolManager(cfg.Client.FastProxyClient)
		poolMgr.logger = cfg.Logger

		// Test keep-alive with valid backend
		poolMgr.performKeepAlive()

		// Should have made a request to the mock server
		suite.GreaterOrEqual(suite.mockServerCalled.Load(), int32(1))
	})

	suite.Run("recovery mechanism", func() {
		poolMgr := NewConnectionPoolManager(cfg.Client.FastProxyClient)
		poolMgr.logger = cfg.Logger

		// Simulate many failures to trigger recovery
		for i := 0; i < 10; i++ {
			poolMgr.RecordConnectionFailure()
		}

		// Check recovery triggers
		poolMgr.checkAndRecover()

		// Verify failure count was reset
		stats := poolMgr.GetConnectionStats()
		suite.Equal(int64(0), stats["connection_failures"])
	})
}

// TestIntegratedHealthManagement tests integration between health manager and connection pool
func (suite *ConnectionResilienceTestSuite) TestIntegratedHealthManagement() {
	suite.Run("global initialization", func() {
		// Initialize global instances
		healthMgr := InitializeBackendHealth(cfg.Client.FastProxyClient, cfg.Server.HostGraphQL, cfg.Logger)
		poolMgr := NewConnectionPoolManager(cfg.Client.FastProxyClient)

		// Set global instances
		backendHealthManager = healthMgr
		connectionPoolManager = poolMgr

		// Test global access
		suite.Equal(healthMgr, GetBackendHealthManager())
		suite.Equal(poolMgr, GetConnectionPoolManager())
	})

	suite.Run("health manager startup", func() {
		healthMgr := InitializeBackendHealth(cfg.Client.FastProxyClient, cfg.Server.HostGraphQL, cfg.Logger)
		backendHealthManager = healthMgr

		// Start health checking
		healthMgr.StartHealthChecking()

		// Wait for backend to be ready
		err := healthMgr.WaitForBackendReady(10 * time.Second)
		suite.NoError(err)

		// Give some time for health checks to run
		time.Sleep(100 * time.Millisecond)

		// Verify health status
		suite.True(healthMgr.IsHealthy())
		suite.Equal(int32(0), healthMgr.GetConsecutiveFailures())
	})
}

// TestConnectionErrorDetection tests connection error detection
func (suite *ConnectionResilienceTestSuite) TestConnectionErrorDetection() {
	testCases := []struct {
		name     string
		errorMsg string
		expected bool
	}{
		{"connection refused", "connection refused", true},
		{"connection reset", "connection reset by peer", true},
		{"no route to host", "no route to host", true},
		{"network unreachable", "network is unreachable", true},
		{"broken pipe", "broken pipe", true},
		{"EOF", "EOF", true},
		{"dial tcp", "dial tcp 127.0.0.1:99999: connect: connection refused", true},
		{"regular error", "some other error", false},
		{"timeout error", "timeout exceeded", false},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			fakeErr := &mockError{msg: tc.errorMsg}
			isConn := isConnectionError(fakeErr)
			suite.Equal(tc.expected, isConn)
		})
	}
}

// mockError is a simple error implementation for testing
type mockError struct {
	msg string
}

func (e *mockError) Error() string {
	return e.msg
}

// TestRetryLogic tests the enhanced retry mechanism
func (suite *ConnectionResilienceTestSuite) TestRetryLogic() {
	suite.Run("connection error classification", func() {
		// Test that connection errors are properly identified
		connErr := &mockError{msg: "connection refused"}
		suite.True(isConnectionError(connErr))

		timeoutErr := &mockError{msg: "timeout exceeded"}
		suite.False(isConnectionError(timeoutErr))
	})
}

// Start the test suite
func TestConnectionResilienceSuite(t *testing.T) {
	suite.Run(t, new(ConnectionResilienceTestSuite))
}
