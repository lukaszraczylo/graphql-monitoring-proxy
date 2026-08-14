package main

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	libpack_logging "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ShutdownTestSuite struct {
	suite.Suite
	origCfg *config
}

func TestShutdownTestSuite(t *testing.T) {
	suite.Run(t, new(ShutdownTestSuite))
}

func (suite *ShutdownTestSuite) SetupTest() {
	cfgMutex.RLock()
	suite.origCfg = cfg
	cfgMutex.RUnlock()
	cfgMutex.Lock()
	cfg = &config{
		Logger: libpack_logging.New(),
	}
	cfgMutex.Unlock()
}

func (suite *ShutdownTestSuite) TearDownTest() {
	cfgMutex.Lock()
	cfg = suite.origCfg
	cfgMutex.Unlock()
}

func (suite *ShutdownTestSuite) TestNewShutdownManager() {
	ctx := context.Background()
	sm := NewShutdownManager(ctx)

	assert.NotNil(suite.T(), sm)
	assert.NotNil(suite.T(), sm.ctx)
	assert.NotNil(suite.T(), sm.cancel)
	assert.Empty(suite.T(), sm.components)
}

func (suite *ShutdownTestSuite) TestRegisterComponent() {
	sm := NewShutdownManager(context.Background())

	// Register multiple components
	sm.RegisterComponent("component1", func(ctx context.Context) error {
		return nil
	})

	sm.RegisterComponent("component2", func(ctx context.Context) error {
		return nil
	})

	assert.Len(suite.T(), sm.components, 2)
	assert.Equal(suite.T(), "component1", sm.components[0].Name)
	assert.Equal(suite.T(), "component2", sm.components[1].Name)
}

func (suite *ShutdownTestSuite) TestRegisterComponentConcurrent() {
	sm := NewShutdownManager(context.Background())

	var wg sync.WaitGroup
	numComponents := 100

	for i := 0; i < numComponents; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			sm.RegisterComponent(
				"component"+string(rune(idx)),
				func(ctx context.Context) error {
					return nil
				},
			)
		}(i)
	}

	wg.Wait()
	assert.Len(suite.T(), sm.components, numComponents)
}

func (suite *ShutdownTestSuite) TestRunGoroutine() {
	sm := NewShutdownManager(context.Background())

	goroutineStarted := make(chan bool, 1)
	goroutineFinished := make(chan bool, 1)

	sm.RunGoroutine("test-goroutine", func(ctx context.Context) {
		goroutineStarted <- true
		<-ctx.Done()
		goroutineFinished <- true
	})

	// Wait for goroutine to start
	select {
	case <-goroutineStarted:
		// Good, goroutine started
	case <-time.After(100 * time.Millisecond):
		suite.T().Fatal("Goroutine did not start")
	}

	// Cancel context to trigger shutdown
	sm.cancel()

	// Wait for goroutine to finish
	select {
	case <-goroutineFinished:
		// Good, goroutine finished
	case <-time.After(100 * time.Millisecond):
		suite.T().Fatal("Goroutine did not finish")
	}
}

func (suite *ShutdownTestSuite) TestRunGoroutineMultiple() {
	sm := NewShutdownManager(context.Background())

	var counter int32
	numGoroutines := 10

	for i := 0; i < numGoroutines; i++ {
		sm.RunGoroutine("goroutine"+string(rune(i)), func(ctx context.Context) {
			atomic.AddInt32(&counter, 1)
			<-ctx.Done()
			atomic.AddInt32(&counter, -1)
		})
	}

	// Give goroutines time to start
	time.Sleep(50 * time.Millisecond)
	assert.Equal(suite.T(), int32(numGoroutines), atomic.LoadInt32(&counter))

	// Cancel and wait for shutdown
	sm.cancel()
	sm.wg.Wait()

	assert.Equal(suite.T(), int32(0), atomic.LoadInt32(&counter))
}

func (suite *ShutdownTestSuite) TestShutdownSuccess() {
	sm := NewShutdownManager(context.Background())

	component1Shutdown := false
	sm.RegisterComponent("component1", func(ctx context.Context) error {
		component1Shutdown = true
		return nil
	})

	component2Shutdown := false
	sm.RegisterComponent("component2", func(ctx context.Context) error {
		component2Shutdown = true
		return nil
	})

	goroutineShutdown := make(chan bool, 1)
	sm.RunGoroutine("test-goroutine", func(ctx context.Context) {
		<-ctx.Done()
		goroutineShutdown <- true
	})

	// Perform shutdown
	err := sm.Shutdown(1 * time.Second)
	assert.NoError(suite.T(), err)

	// Verify all components were shut down
	assert.True(suite.T(), component1Shutdown)
	assert.True(suite.T(), component2Shutdown)

	// Verify goroutine was shut down
	select {
	case <-goroutineShutdown:
		// Good
	case <-time.After(100 * time.Millisecond):
		suite.T().Fatal("Goroutine did not shut down")
	}
}

func (suite *ShutdownTestSuite) TestShutdownWithError() {
	sm := NewShutdownManager(context.Background())

	componentShutdown := false
	sm.RegisterComponent("failing-component", func(ctx context.Context) error {
		componentShutdown = true
		return errors.New("shutdown failed")
	})

	// Shutdown should continue even if a component fails
	err := sm.Shutdown(1 * time.Second)
	assert.NoError(suite.T(), err) // Shutdown manager doesn't return component errors
	assert.True(suite.T(), componentShutdown)
}

func (suite *ShutdownTestSuite) TestShutdownTimeout() {
	sm := NewShutdownManager(context.Background())

	// Register a component that takes too long to shutdown
	sm.RegisterComponent("slow-component", func(ctx context.Context) error {
		select {
		case <-time.After(2 * time.Second):
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	})

	// Shutdown with short timeout
	start := time.Now()
	err := sm.Shutdown(100 * time.Millisecond)
	elapsed := time.Since(start)

	// Should timeout quickly
	assert.NoError(suite.T(), err)
	assert.Less(suite.T(), elapsed, 500*time.Millisecond)
}

func (suite *ShutdownTestSuite) TestShutdownConcurrentComponents() {
	sm := NewShutdownManager(context.Background())

	var shutdownOrder []int
	var mu sync.Mutex

	// Register multiple components that shutdown concurrently
	for i := 0; i < 5; i++ {
		idx := i
		sm.RegisterComponent("component"+string(rune(idx)), func(ctx context.Context) error {
			time.Sleep(time.Duration(idx*10) * time.Millisecond)
			mu.Lock()
			shutdownOrder = append(shutdownOrder, idx)
			mu.Unlock()
			return nil
		})
	}

	err := sm.Shutdown(1 * time.Second)
	assert.NoError(suite.T(), err)

	// All components should have shut down
	assert.Len(suite.T(), shutdownOrder, 5)
}

func (suite *ShutdownTestSuite) TestShutdownIdempotent() {
	sm := NewShutdownManager(context.Background())

	shutdownCount := int32(0)
	sm.RegisterComponent("component", func(ctx context.Context) error {
		atomic.AddInt32(&shutdownCount, 1)
		return nil
	})

	// First shutdown
	err := sm.Shutdown(100 * time.Millisecond)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int32(1), atomic.LoadInt32(&shutdownCount))

	// Second shutdown should be safe but not call components again
	err = sm.Shutdown(100 * time.Millisecond)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int32(1), atomic.LoadInt32(&shutdownCount))
}

func (suite *ShutdownTestSuite) TestShutdownEmptyManager() {
	sm := NewShutdownManager(context.Background())

	// Shutdown with no components should be safe
	err := sm.Shutdown(100 * time.Millisecond)
	assert.NoError(suite.T(), err)
}

func (suite *ShutdownTestSuite) TestContextCancellation() {
	ctx, cancel := context.WithCancel(context.Background())
	sm := NewShutdownManager(ctx)

	goroutineExited := make(chan bool, 1)
	sm.RunGoroutine("test-goroutine", func(ctx context.Context) {
		<-ctx.Done()
		goroutineExited <- true
	})

	// Cancel the parent context
	cancel()

	// Goroutine should still exit properly
	select {
	case <-goroutineExited:
		// Good
	case <-time.After(100 * time.Millisecond):
		suite.T().Fatal("Goroutine did not exit after context cancellation")
	}
}

func (suite *ShutdownTestSuite) TestRegisterComponentDefaultPhase() {
	sm := NewShutdownManager(context.Background())

	sm.RegisterComponent("component1", func(ctx context.Context) error {
		return nil
	})

	assert.Len(suite.T(), sm.components, 1)
	assert.Equal(suite.T(), DefaultShutdownPhase, sm.components[0].Phase)
}

func (suite *ShutdownTestSuite) TestRegisterComponentWithPhase() {
	sm := NewShutdownManager(context.Background())

	sm.RegisterComponentWithPhase("early", 1, func(ctx context.Context) error {
		return nil
	})
	sm.RegisterComponentWithPhase("late", 5, func(ctx context.Context) error {
		return nil
	})

	assert.Len(suite.T(), sm.components, 2)
	assert.Equal(suite.T(), 1, sm.components[0].Phase)
	assert.Equal(suite.T(), 5, sm.components[1].Phase)
}

// TestShutdownPhasesRunSequentially registers components in two phases and
// records each component's start/end wall-clock time. It asserts every
// phase-2 component starts only after every phase-1 component has finished,
// proving phases execute sequentially rather than all at once.
func (suite *ShutdownTestSuite) TestShutdownPhasesRunSequentially() {
	sm := NewShutdownManager(context.Background())

	type span struct {
		start time.Time
		end   time.Time
	}
	var mu sync.Mutex
	spans := make(map[string]span)

	recordSpan := func(name string, work func()) {
		mu.Lock()
		spans[name] = span{start: time.Now()}
		mu.Unlock()

		work()

		mu.Lock()
		s := spans[name]
		s.end = time.Now()
		spans[name] = s
		mu.Unlock()
	}

	sm.RegisterComponentWithPhase("phase1-a", 1, func(ctx context.Context) error {
		recordSpan("phase1-a", func() { time.Sleep(30 * time.Millisecond) })
		return nil
	})
	sm.RegisterComponentWithPhase("phase1-b", 1, func(ctx context.Context) error {
		recordSpan("phase1-b", func() { time.Sleep(30 * time.Millisecond) })
		return nil
	})
	sm.RegisterComponentWithPhase("phase2-a", 2, func(ctx context.Context) error {
		recordSpan("phase2-a", func() {})
		return nil
	})
	sm.RegisterComponentWithPhase("phase2-b", 2, func(ctx context.Context) error {
		recordSpan("phase2-b", func() {})
		return nil
	})

	err := sm.Shutdown(2 * time.Second)
	assert.NoError(suite.T(), err)

	mu.Lock()
	defer mu.Unlock()

	assert.Len(suite.T(), spans, 4, "all four components must have run")

	phase1End := spans["phase1-a"].end
	if spans["phase1-b"].end.After(phase1End) {
		phase1End = spans["phase1-b"].end
	}

	assert.False(suite.T(), spans["phase2-a"].start.Before(phase1End),
		"phase2-a started before every phase-1 component finished")
	assert.False(suite.T(), spans["phase2-b"].start.Before(phase1End),
		"phase2-b started before every phase-1 component finished")
}

// TestShutdownSamePhaseComponentsRunConcurrently proves components sharing a
// phase still run concurrently. Every registered component blocks on a
// shared release gate that only opens once all of them have signaled they
// started; if the components ran one at a time, the first would block
// forever on the gate and the test would time out.
func (suite *ShutdownTestSuite) TestShutdownSamePhaseComponentsRunConcurrently() {
	sm := NewShutdownManager(context.Background())

	const numComponents = 3
	started := make(chan struct{}, numComponents)
	release := make(chan struct{})

	for i := 0; i < numComponents; i++ {
		sm.RegisterComponentWithPhase("component", 7, func(ctx context.Context) error {
			started <- struct{}{}
			<-release
			return nil
		})
	}

	shutdownDone := make(chan struct{})
	go func() {
		err := sm.Shutdown(2 * time.Second)
		assert.NoError(suite.T(), err)
		close(shutdownDone)
	}()

	for i := 0; i < numComponents; i++ {
		select {
		case <-started:
		case <-time.After(1 * time.Second):
			suite.T().Fatal("same-phase components did not all start concurrently")
		}
	}
	close(release)

	select {
	case <-shutdownDone:
	case <-time.After(1 * time.Second):
		suite.T().Fatal("shutdown did not complete after releasing same-phase components")
	}
}

// TestShutdownPhaseTimeoutStillAttemptsNextPhase proves a phase-1 component
// that ignores the ctx deadline and outlives the shutdown budget does not
// prevent phase 2 from being attempted, and that Shutdown still returns
// promptly instead of blocking on the hung component.
func (suite *ShutdownTestSuite) TestShutdownPhaseTimeoutStillAttemptsNextPhase() {
	sm := NewShutdownManager(context.Background())

	phase1Started := make(chan struct{})
	sm.RegisterComponentWithPhase("phase1-hung", 1, func(ctx context.Context) error {
		close(phase1Started)
		// Ignore ctx entirely to simulate a component that does not
		// respect the shutdown deadline.
		time.Sleep(5 * time.Second)
		return nil
	})

	phase2Attempted := make(chan struct{})
	sm.RegisterComponentWithPhase("phase2-after-hang", 2, func(ctx context.Context) error {
		close(phase2Attempted)
		return nil
	})

	start := time.Now()
	err := sm.Shutdown(100 * time.Millisecond)
	elapsed := time.Since(start)

	assert.NoError(suite.T(), err)
	assert.Less(suite.T(), elapsed, 1*time.Second,
		"shutdown must not block waiting on a hung phase-1 component")

	select {
	case <-phase1Started:
	default:
		suite.T().Fatal("phase-1 component was never started")
	}

	select {
	case <-phase2Attempted:
		// Good: phase 2 was attempted despite phase 1 exceeding the deadline.
	case <-time.After(500 * time.Millisecond):
		suite.T().Fatal("phase-2 component was never attempted after phase-1 exceeded the shutdown deadline")
	}
}

// Benchmark tests
func BenchmarkRegisterComponent(b *testing.B) {
	sm := NewShutdownManager(context.Background())
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sm.RegisterComponent("component", func(ctx context.Context) error {
			return nil
		})
	}
}

func BenchmarkShutdown(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		sm := NewShutdownManager(context.Background())
		for j := 0; j < 10; j++ {
			sm.RegisterComponent("component"+string(rune(j)), func(ctx context.Context) error {
				return nil
			})
		}
		b.StartTimer()

		sm.Shutdown(100 * time.Millisecond)
	}
}
