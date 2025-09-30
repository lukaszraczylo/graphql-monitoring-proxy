package main

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	libpack_logger "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
	libpack_monitoring "github.com/lukaszraczylo/graphql-monitoring-proxy/monitoring"
	"github.com/stretchr/testify/assert"
)

func TestNewRequestCoalescer(t *testing.T) {
	tests := []struct {
		name    string
		enabled bool
	}{
		{
			name:    "enabled coalescer",
			enabled: true,
		},
		{
			name:    "disabled coalescer",
			enabled: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := libpack_logger.New()
			monitoring := libpack_monitoring.NewMonitoring(&libpack_monitoring.InitConfig{})

			rc := NewRequestCoalescer(tt.enabled, logger, monitoring)

			assert.NotNil(t, rc)
			assert.Equal(t, tt.enabled, rc.enabled)
		})
	}
}

func TestRequestCoalescer_Do_SingleRequest(t *testing.T) {
	rc := NewRequestCoalescer(true, libpack_logger.New(), nil)

	executed := false
	response := &CoalescedResponse{
		Body:       []byte("test response"),
		StatusCode: 200,
	}

	fn := func() (*CoalescedResponse, error) {
		executed = true
		return response, nil
	}

	result, err := rc.Do("test-key", fn)

	assert.NoError(t, err)
	assert.True(t, executed)
	assert.Equal(t, response, result)

	stats := rc.GetStats()
	assert.Equal(t, int64(1), stats["total_requests"])
	assert.Equal(t, int64(1), stats["primary_requests"])
	assert.Equal(t, int64(0), stats["coalesced_requests"])
}

func TestRequestCoalescer_Do_ConcurrentRequests(t *testing.T) {
	rc := NewRequestCoalescer(true, libpack_logger.New(), nil)

	var executionCount atomic.Int32
	response := &CoalescedResponse{
		Body:       []byte("test response"),
		StatusCode: 200,
	}

	fn := func() (*CoalescedResponse, error) {
		executionCount.Add(1)
		time.Sleep(50 * time.Millisecond) // Simulate work
		return response, nil
	}

	// Launch concurrent requests with the same key
	concurrentRequests := 10
	var wg sync.WaitGroup
	wg.Add(concurrentRequests)

	results := make([]*CoalescedResponse, concurrentRequests)
	errs := make([]error, concurrentRequests)

	for i := 0; i < concurrentRequests; i++ {
		go func(index int) {
			defer wg.Done()
			results[index], errs[index] = rc.Do("same-key", fn)
		}(i)
	}

	wg.Wait()

	// Function should only execute once
	assert.Equal(t, int32(1), executionCount.Load())

	// All requests should get the same response
	for i := 0; i < concurrentRequests; i++ {
		assert.NoError(t, errs[i])
		assert.Equal(t, response, results[i])
	}

	stats := rc.GetStats()
	assert.Equal(t, int64(concurrentRequests), stats["total_requests"])
	assert.Equal(t, int64(1), stats["primary_requests"])
	assert.Equal(t, int64(concurrentRequests-1), stats["coalesced_requests"])

	// Check backend savings
	backendSavings := stats["backend_savings_pct"].(float64)
	assert.Greater(t, backendSavings, 0.0)
}

func TestRequestCoalescer_Do_DifferentKeys(t *testing.T) {
	rc := NewRequestCoalescer(true, libpack_logger.New(), nil)

	var executionCount atomic.Int32

	fn := func() (*CoalescedResponse, error) {
		executionCount.Add(1)
		return &CoalescedResponse{Body: []byte("response")}, nil
	}

	// Concurrent requests with different keys
	var wg sync.WaitGroup
	keys := []string{"key1", "key2", "key3"}

	for _, key := range keys {
		wg.Add(1)
		go func(k string) {
			defer wg.Done()
			rc.Do(k, fn)
		}(key)
	}

	wg.Wait()

	// Function should execute for each unique key
	assert.Equal(t, int32(len(keys)), executionCount.Load())

	stats := rc.GetStats()
	assert.Equal(t, int64(3), stats["primary_requests"])
	assert.Equal(t, int64(0), stats["coalesced_requests"])
}

func TestRequestCoalescer_Do_Error(t *testing.T) {
	rc := NewRequestCoalescer(true, libpack_logger.New(), nil)

	expectedErr := errors.New("test error")

	fn := func() (*CoalescedResponse, error) {
		return nil, expectedErr
	}

	result, err := rc.Do("error-key", fn)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Error(t, result.Err)
	assert.Equal(t, expectedErr, result.Err)
}

func TestRequestCoalescer_Do_ConcurrentWithError(t *testing.T) {
	rc := NewRequestCoalescer(true, libpack_logger.New(), nil)

	expectedErr := errors.New("test error")
	var executionCount atomic.Int32

	fn := func() (*CoalescedResponse, error) {
		executionCount.Add(1)
		time.Sleep(50 * time.Millisecond)
		return nil, expectedErr
	}

	// Launch concurrent requests
	concurrentRequests := 5
	var wg sync.WaitGroup
	wg.Add(concurrentRequests)

	results := make([]*CoalescedResponse, concurrentRequests)

	for i := 0; i < concurrentRequests; i++ {
		go func(index int) {
			defer wg.Done()
			results[index], _ = rc.Do("error-key", fn)
		}(i)
	}

	wg.Wait()

	// Function should only execute once
	assert.Equal(t, int32(1), executionCount.Load())

	// All requests should get the same error in response
	for i := 0; i < concurrentRequests; i++ {
		assert.NotNil(t, results[i])
		assert.Error(t, results[i].Err)
		assert.Equal(t, expectedErr, results[i].Err)
	}
}

func TestRequestCoalescer_Do_Disabled(t *testing.T) {
	rc := NewRequestCoalescer(false, libpack_logger.New(), nil)

	var executionCount atomic.Int32

	fn := func() (*CoalescedResponse, error) {
		executionCount.Add(1)
		return &CoalescedResponse{Body: []byte("response")}, nil
	}

	// Launch concurrent requests with the same key
	concurrentRequests := 5
	var wg sync.WaitGroup
	wg.Add(concurrentRequests)

	for i := 0; i < concurrentRequests; i++ {
		go func() {
			defer wg.Done()
			rc.Do("same-key", fn)
		}()
	}

	wg.Wait()

	// When disabled, function should execute for each request
	assert.Equal(t, int32(concurrentRequests), executionCount.Load())

	stats := rc.GetStats()
	assert.Equal(t, false, stats["enabled"])
}

func TestRequestCoalescer_GetStats(t *testing.T) {
	rc := NewRequestCoalescer(true, libpack_logger.New(), nil)

	fn := func() (*CoalescedResponse, error) {
		time.Sleep(10 * time.Millisecond)
		return &CoalescedResponse{Body: []byte("response")}, nil
	}

	// Simulate some coalesced requests
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			rc.Do("key1", fn)
		}()
	}
	wg.Wait()

	// Add some non-coalesced requests
	rc.Do("key2", fn)
	rc.Do("key3", fn)

	stats := rc.GetStats()

	assert.Equal(t, true, stats["enabled"])
	assert.Equal(t, int64(12), stats["total_requests"])
	assert.Equal(t, int64(3), stats["primary_requests"])
	assert.Equal(t, int64(9), stats["coalesced_requests"])
	assert.Equal(t, int64(0), stats["inflight_count"])

	coalescingRate := stats["coalescing_rate_pct"].(float64)
	assert.Greater(t, coalescingRate, 0.0)
	assert.LessOrEqual(t, coalescingRate, 100.0)

	backendSavings := stats["backend_savings_pct"].(float64)
	assert.Greater(t, backendSavings, 0.0)
}

func TestRequestCoalescer_Reset(t *testing.T) {
	rc := NewRequestCoalescer(true, libpack_logger.New(), nil)

	fn := func() (*CoalescedResponse, error) {
		return &CoalescedResponse{Body: []byte("response")}, nil
	}

	// Generate some activity
	rc.Do("key1", fn)
	rc.Do("key2", fn)

	statsBefore := rc.GetStats()
	assert.Greater(t, statsBefore["total_requests"].(int64), int64(0))

	// Reset
	rc.Reset()

	statsAfter := rc.GetStats()
	assert.Equal(t, int64(0), statsAfter["total_requests"])
	assert.Equal(t, int64(0), statsAfter["primary_requests"])
	assert.Equal(t, int64(0), statsAfter["coalesced_requests"])
}

func TestRequestCoalescer_RaceCondition(t *testing.T) {
	rc := NewRequestCoalescer(true, libpack_logger.New(), nil)

	var executionCount atomic.Int32

	fn := func() (*CoalescedResponse, error) {
		executionCount.Add(1)
		time.Sleep(5 * time.Millisecond)
		return &CoalescedResponse{Body: []byte("response")}, nil
	}

	// Launch many concurrent requests in waves
	waves := 5
	requestsPerWave := 20

	for wave := 0; wave < waves; wave++ {
		var wg sync.WaitGroup
		wg.Add(requestsPerWave)

		for i := 0; i < requestsPerWave; i++ {
			go func() {
				defer wg.Done()
				rc.Do("race-key", fn)
			}()
		}

		wg.Wait()
		time.Sleep(10 * time.Millisecond) // Small delay between waves
	}

	// Execution count should be much less than total requests
	totalRequests := waves * requestsPerWave
	assert.Less(t, int(executionCount.Load()), totalRequests)

	stats := rc.GetStats()
	assert.Equal(t, int64(totalRequests), stats["total_requests"])
}

func TestRequestCoalescer_BackendSavingsCalculation(t *testing.T) {
	tests := []struct {
		name              string
		totalRequests     int64
		coalescedRequests int64
		expectedSavings   float64
	}{
		{
			name:              "50% savings",
			totalRequests:     100,
			coalescedRequests: 50,
			expectedSavings:   100.0, // 50 coalesced / 50 primary = 100%
		},
		{
			name:              "90% savings",
			totalRequests:     100,
			coalescedRequests: 90,
			expectedSavings:   900.0, // 90 coalesced / 10 primary = 900%
		},
		{
			name:              "no savings",
			totalRequests:     100,
			coalescedRequests: 0,
			expectedSavings:   0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc := NewRequestCoalescer(true, libpack_logger.New(), nil)

			rc.totalRequests.Store(tt.totalRequests)
			rc.coalescedRequests.Store(tt.coalescedRequests)

			stats := rc.GetStats()
			savings := stats["backend_savings_pct"].(float64)

			assert.InDelta(t, tt.expectedSavings, savings, 0.1)
		})
	}
}

func TestRequestCoalescer_GlobalInstance(t *testing.T) {
	rc := InitializeRequestCoalescer(true, libpack_logger.New(), nil)
	assert.NotNil(t, rc)

	// Should return the same instance
	rc2 := GetRequestCoalescer()
	assert.Equal(t, rc, rc2)
}

func TestMin(t *testing.T) {
	tests := []struct {
		a        int
		b        int
		expected int
	}{
		{a: 5, b: 10, expected: 5},
		{a: 10, b: 5, expected: 5},
		{a: 5, b: 5, expected: 5},
		{a: 0, b: 10, expected: 0},
		{a: -5, b: 5, expected: -5},
	}

	for _, tt := range tests {
		result := min(tt.a, tt.b)
		assert.Equal(t, tt.expected, result)
	}
}
