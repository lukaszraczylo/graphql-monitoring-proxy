package main

import (
	"sync"
	"testing"
	"time"

	libpack_logging "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/valyala/fasthttp"
)

type ConnectionPoolTestSuite struct {
	suite.Suite
	origCfg               *config
	origConnectionManager *ConnectionPoolManager
}

func TestConnectionPoolTestSuite(t *testing.T) {
	suite.Run(t, new(ConnectionPoolTestSuite))
}

func (suite *ConnectionPoolTestSuite) SetupTest() {
	suite.origCfg = cfg
	cfg = &config{
		Logger: libpack_logging.New(),
	}
	suite.origConnectionManager = connectionPoolManager
	connectionPoolManager = nil
}

func (suite *ConnectionPoolTestSuite) TearDownTest() {
	if connectionPoolManager != nil {
		connectionPoolManager.Shutdown()
		connectionPoolManager = nil
	}
	cfg = suite.origCfg
	connectionPoolManager = suite.origConnectionManager
}

func (suite *ConnectionPoolTestSuite) TestNewConnectionPoolManager() {
	client := &fasthttp.Client{
		MaxConnsPerHost: 100,
	}

	cpm := NewConnectionPoolManager(client)
	assert.NotNil(suite.T(), cpm)
	assert.NotNil(suite.T(), cpm.client)
	assert.NotNil(suite.T(), cpm.ctx)
	assert.NotNil(suite.T(), cpm.cancel)

	// Cleanup
	cpm.Shutdown()
}

func (suite *ConnectionPoolTestSuite) TestGetClient() {
	client := &fasthttp.Client{
		MaxConnsPerHost: 100,
	}

	cpm := NewConnectionPoolManager(client)
	defer cpm.Shutdown()

	retrievedClient := cpm.GetClient()
	assert.Equal(suite.T(), client, retrievedClient)
}

func (suite *ConnectionPoolTestSuite) TestShutdown() {
	client := &fasthttp.Client{
		MaxConnsPerHost: 100,
	}

	cpm := NewConnectionPoolManager(client)

	// Shutdown should be safe
	err := cpm.Shutdown()
	assert.NoError(suite.T(), err)

	// Multiple shutdowns should be safe
	err = cpm.Shutdown()
	assert.NoError(suite.T(), err)
}

func (suite *ConnectionPoolTestSuite) TestShutdownNil() {
	var cpm *ConnectionPoolManager
	err := cpm.Shutdown()
	assert.NoError(suite.T(), err)
}

func (suite *ConnectionPoolTestSuite) TestPeriodicCleanup() {
	client := &fasthttp.Client{
		MaxConnsPerHost: 100,
	}

	cpm := NewConnectionPoolManager(client)

	// Let the cleanup goroutine run
	time.Sleep(50 * time.Millisecond)

	// Shutdown should stop the cleanup goroutine
	err := cpm.Shutdown()
	assert.NoError(suite.T(), err)
}

func (suite *ConnectionPoolTestSuite) TestCleanIdleConnections() {
	client := &fasthttp.Client{
		MaxConnsPerHost: 100,
	}

	cpm := NewConnectionPoolManager(client)
	defer cpm.Shutdown()

	// Manually trigger cleanup
	cpm.cleanIdleConnections()

	// Should not panic or error
	assert.NotNil(suite.T(), cpm.client)
}

func (suite *ConnectionPoolTestSuite) TestConcurrentAccess() {
	client := &fasthttp.Client{
		MaxConnsPerHost: 100,
	}

	cpm := NewConnectionPoolManager(client)
	defer cpm.Shutdown()

	var wg sync.WaitGroup

	// Concurrent reads
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				c := cpm.GetClient()
				assert.NotNil(suite.T(), c)
				time.Sleep(time.Microsecond)
			}
		}()
	}

	// Concurrent cleanups
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				cpm.cleanIdleConnections()
				time.Sleep(time.Millisecond)
			}
		}()
	}

	wg.Wait()
}

func (suite *ConnectionPoolTestSuite) TestInitializeConnectionPool() {
	client := &fasthttp.Client{
		MaxConnsPerHost: 200,
	}

	InitializeConnectionPool(client)
	assert.NotNil(suite.T(), connectionPoolManager)
	assert.Equal(suite.T(), client, connectionPoolManager.GetClient())

	// Initialize again should replace the old one
	newClient := &fasthttp.Client{
		MaxConnsPerHost: 300,
	}
	InitializeConnectionPool(newClient)
	assert.Equal(suite.T(), newClient, connectionPoolManager.GetClient())
}

func (suite *ConnectionPoolTestSuite) TestShutdownConnectionPool() {
	client := &fasthttp.Client{
		MaxConnsPerHost: 100,
	}

	InitializeConnectionPool(client)
	assert.NotNil(suite.T(), connectionPoolManager)

	ShutdownConnectionPool()
	assert.Nil(suite.T(), connectionPoolManager)

	// Shutdown again should be safe
	ShutdownConnectionPool()
	assert.Nil(suite.T(), connectionPoolManager)
}

func (suite *ConnectionPoolTestSuite) TestGetConnectionPoolManager() {
	assert.Nil(suite.T(), GetConnectionPoolManager())

	client := &fasthttp.Client{
		MaxConnsPerHost: 100,
	}
	InitializeConnectionPool(client)

	manager := GetConnectionPoolManager()
	assert.NotNil(suite.T(), manager)
	assert.Equal(suite.T(), connectionPoolManager, manager)

	ShutdownConnectionPool()
	assert.Nil(suite.T(), GetConnectionPoolManager())
}

func (suite *ConnectionPoolTestSuite) TestContextCancellation() {
	client := &fasthttp.Client{
		MaxConnsPerHost: 100,
	}

	cpm := NewConnectionPoolManager(client)

	// Cancel the context
	cpm.cancel()

	// Give the cleanup goroutine time to exit
	time.Sleep(50 * time.Millisecond)

	// Shutdown should still work
	err := cpm.Shutdown()
	assert.NoError(suite.T(), err)
}

func (suite *ConnectionPoolTestSuite) TestRaceConditions() {
	client := &fasthttp.Client{
		MaxConnsPerHost: 100,
	}

	var wg sync.WaitGroup

	// Concurrent initialization and shutdown
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			InitializeConnectionPool(client)
		}()
	}

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(time.Microsecond)
			ShutdownConnectionPool()
		}()
	}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			manager := GetConnectionPoolManager()
			if manager != nil {
				_ = manager.GetClient()
			}
		}()
	}

	wg.Wait()
}

func (suite *ConnectionPoolTestSuite) TestCleanupWithNilLogger() {
	// Test cleanup when cfg or logger is nil
	origCfg := cfg
	cfg = nil

	client := &fasthttp.Client{
		MaxConnsPerHost: 100,
	}

	cpm := NewConnectionPoolManager(client)

	// Should not panic
	cpm.cleanIdleConnections()
	err := cpm.Shutdown()
	assert.NoError(suite.T(), err)

	cfg = origCfg
}

func (suite *ConnectionPoolTestSuite) TestMemoryManagement() {
	// Test that connection pool properly manages memory
	client := &fasthttp.Client{
		MaxConnsPerHost:     10,
		MaxIdleConnDuration: 100 * time.Millisecond,
	}

	cpm := NewConnectionPoolManager(client)
	defer cpm.Shutdown()

	// Simulate connections being created and becoming idle
	// The periodic cleanup should handle them
	time.Sleep(150 * time.Millisecond)

	// Manual cleanup to ensure connections are released
	cpm.cleanIdleConnections()

	// Verify client is still accessible
	assert.NotNil(suite.T(), cpm.GetClient())
}

// Benchmark tests
func BenchmarkConnectionPoolGetClient(b *testing.B) {
	client := &fasthttp.Client{
		MaxConnsPerHost: 100,
	}

	cpm := NewConnectionPoolManager(client)
	defer cpm.Shutdown()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = cpm.GetClient()
		}
	})
}

func BenchmarkConnectionPoolCleanup(b *testing.B) {
	client := &fasthttp.Client{
		MaxConnsPerHost: 100,
	}

	cpm := NewConnectionPoolManager(client)
	defer cpm.Shutdown()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cpm.cleanIdleConnections()
	}
}
