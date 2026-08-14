package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"
	libpack_cache "github.com/lukaszraczylo/graphql-monitoring-proxy/cache"
	libpack_cache_memory "github.com/lukaszraczylo/graphql-monitoring-proxy/cache/memory"
	libpack_logger "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
	libpack_monitoring "github.com/lukaszraczylo/graphql-monitoring-proxy/monitoring"
	"github.com/sony/gobreaker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/valyala/fasthttp"
)

// CircuitBreakerTestSuite is a test suite for circuit breaker functionality
type CircuitBreakerTestSuite struct {
	suite.Suite
	originalConfig *config
	outputBuffer   *bytes.Buffer // Used to capture logger output
}

func (suite *CircuitBreakerTestSuite) SetupTest() {

	// Store original config to restore later
	suite.originalConfig = cfg

	// Create a buffer to capture logger output
	suite.outputBuffer = &bytes.Buffer{}

	// Setup a new config with a real logger that writes to our buffer
	cfg = &config{}
	cfg.Logger = libpack_logger.New().SetOutput(suite.outputBuffer)

	// Initialize monitoring with a minimal configuration
	cfg.Monitoring = libpack_monitoring.NewMonitoring(&libpack_monitoring.InitConfig{
		PurgeOnCrawl: false,
		PurgeEvery:   0,
	})

	// Configure circuit breaker settings
	cfg.CircuitBreaker.Enable = true
	cfg.CircuitBreaker.MaxFailures = 3
	cfg.CircuitBreaker.Timeout = 5
	cfg.CircuitBreaker.MaxRequestsInHalfOpen = 2
	cfg.CircuitBreaker.ReturnCachedOnOpen = true
	cfg.CircuitBreaker.TripOn5xx = true

	// Initialize memory cache
	memCache := libpack_cache_memory.New(time.Minute)
	cacheConfig := &libpack_cache.CacheConfig{
		Logger: cfg.Logger,
		Client: memCache,
		TTL:    60,
	}
	libpack_cache.EnableCache(cacheConfig)
}

func (suite *CircuitBreakerTestSuite) TearDownTest() {
	// Restore original config
	cfg = suite.originalConfig

	// Reset circuit breaker and metrics
	cbMutex.Lock()
	defer cbMutex.Unlock()
	cb = nil
	// Circuit breaker metrics are now managed by cbMetrics
	cbMetrics = nil
	// Reset the C7 progressive-backoff tracker so a trip count/timestamp
	// from one test can't leak into the next.
	cbBackoff.reset()
}

// Helper function to check if a specific message appears in the logger output
func (suite *CircuitBreakerTestSuite) logContains(substring string) bool {
	return strings.Contains(suite.outputBuffer.String(), substring)
}

// TestCreateTripFunc tests the circuit breaker trip function logic
func (suite *CircuitBreakerTestSuite) TestCreateTripFunc() {
	// Create the trip function
	tripFunc := createTripFunc(cfg)

	// Test cases
	testCases := []struct {
		name           string
		counts         gobreaker.Counts
		expectedResult bool
	}{
		{
			name: "below threshold",
			counts: gobreaker.Counts{
				Requests:             10,
				TotalSuccesses:       8,
				TotalFailures:        2,
				ConsecutiveSuccesses: 0,
				ConsecutiveFailures:  2, // Below MaxFailures (3)
			},
			expectedResult: false,
		},
		{
			name: "at threshold",
			counts: gobreaker.Counts{
				Requests:             10,
				TotalSuccesses:       7,
				TotalFailures:        3,
				ConsecutiveSuccesses: 0,
				ConsecutiveFailures:  3, // Equal to MaxFailures (3)
			},
			expectedResult: true,
		},
		{
			name: "above threshold",
			counts: gobreaker.Counts{
				Requests:             10,
				TotalSuccesses:       5,
				TotalFailures:        5,
				ConsecutiveSuccesses: 0,
				ConsecutiveFailures:  5, // Above MaxFailures (3)
			},
			expectedResult: true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Reset the buffer before each test case
			suite.outputBuffer.Reset()

			// Test the trip function
			result := tripFunc(tc.counts)
			suite.Equal(tc.expectedResult, result, "Trip function result should match expected")

			// If it should trip, verify that a warning log was generated
			if tc.expectedResult {
				suite.True(suite.logContains("Circuit breaker tripped"),
					"Expected a warning log when circuit breaker trips")
				suite.True(suite.logContains(fmt.Sprintf(`"consecutive_failures":%d`, tc.counts.ConsecutiveFailures)),
					"Log should contain consecutive failures count")
			}
		})
	}
}

// TestCreateStateChangeFunc tests the state change function logic
func (suite *CircuitBreakerTestSuite) TestCreateStateChangeFunc() {
	// We'll skip this test as it's problematic with the gauge callback issue
	suite.T().Skip("Skipping due to gauge callback issues")
}

// TestCircuitBreakerInitialization tests the circuit breaker initialization
func (suite *CircuitBreakerTestSuite) TestCircuitBreakerInitialization() {
	// Reset the buffer before the test
	suite.outputBuffer.Reset()

	// Initialize circuit breaker
	initCircuitBreaker(cfg)

	// Verify circuit breaker was initialized
	suite.NotNil(cb, "Circuit breaker should be initialized")
	suite.NotNil(cbMetrics, "Circuit breaker metrics should be initialized")

	// Verify the log message
	suite.True(suite.logContains("Circuit breaker initialized"),
		"Log should contain initialization message")

	// Test with disabled circuit breaker
	suite.outputBuffer.Reset()
	cfg.CircuitBreaker.Enable = false

	// Reset circuit breaker
	cbMutex.Lock()
	cb = nil
	cbMetrics = nil
	cbMutex.Unlock()

	// Initialize again with disabled config
	initCircuitBreaker(cfg)

	// Verify circuit breaker was not initialized
	suite.Nil(cb, "Circuit breaker should not be initialized when disabled")

	// Verify the log message
	suite.True(suite.logContains("Circuit breaker is disabled"),
		"Log should contain disabled message")
}

// TestExecuteFunctionBehavior tests the basic behavior of Execute without circuit breaker
func (suite *CircuitBreakerTestSuite) TestExecuteFunctionBehavior() {
	// Reset for this test
	cfg.CircuitBreaker.Enable = true
	initCircuitBreaker(cfg)

	// Test with success
	result := "success"
	execResult, err := cb.Execute(func() (any, error) {
		return result, nil
	})

	suite.NoError(err, "Execute should not return error on success")
	suite.Equal(result, execResult, "Execute should return the correct result value")

	// Test with error
	testErr := errors.New("test error")
	_, err = cb.Execute(func() (any, error) {
		return nil, testErr
	})

	suite.Error(err, "Execute should return error when function returns error")
	suite.Equal(testErr.Error(), err.Error(), "Error message should match")
}

// TestClassifyCircuitFailure is a pure-function table test for the C5/C6
// gating logic: whether a proxy failure counts toward the circuit
// breaker's failure statistics is gated by the matching
// CIRCUIT_TRIP_ON_TIMEOUTS / CIRCUIT_TRIP_ON_4XX / CIRCUIT_TRIP_ON_5XX flag.
// Any other failure kind (e.g. a connection error) has no dedicated flag
// and always counts, regardless of flag values -- that's the pre-fix
// baseline (flags were parsed but never enforced) this preserves.
func TestClassifyCircuitFailure(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		statusCode     int
		tripOnTimeouts bool
		tripOn4xx      bool
		tripOn5xx      bool
		wantCounts     bool
	}{
		{
			name:           "timeout counts when TripOnTimeouts=true (library default)",
			err:            errors.New("timeout"),
			tripOnTimeouts: true,
			wantCounts:     true,
		},
		{
			name:           "timeout does not count when TripOnTimeouts=false",
			err:            errors.New("timeout"),
			tripOnTimeouts: false,
			wantCounts:     false,
		},
		{
			name:       "5xx counts when TripOn5xx=true (library default)",
			err:        fmt.Errorf("received non-200 response: %d", 500),
			statusCode: 500,
			tripOn5xx:  true,
			wantCounts: true,
		},
		{
			name:       "5xx does not count when TripOn5xx=false",
			err:        fmt.Errorf("received non-200 response: %d", 503),
			statusCode: 503,
			tripOn5xx:  false,
			wantCounts: false,
		},
		{
			name:       "4xx counts when TripOn4xx=true",
			err:        fmt.Errorf("received non-200 response: %d", 404),
			statusCode: 404,
			tripOn4xx:  true,
			wantCounts: true,
		},
		{
			name:       "4xx does not count when TripOn4xx=false (library default)",
			err:        fmt.Errorf("received non-200 response: %d", 400),
			statusCode: 400,
			tripOn4xx:  false,
			wantCounts: false,
		},
		{
			name:       "connection error always counts, regardless of any flag",
			err:        errors.New("connection refused"),
			statusCode: 0,
			wantCounts: true,
		},
		{
			name:       "nil fiber context error always counts, regardless of any flag",
			err:        errFiberCtxNilDuringRetry,
			statusCode: 0,
			wantCounts: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCfg := &config{}
			testCfg.CircuitBreaker.TripOnTimeouts = tt.tripOnTimeouts
			testCfg.CircuitBreaker.TripOn4xx = tt.tripOn4xx
			testCfg.CircuitBreaker.TripOn5xx = tt.tripOn5xx

			got := classifyCircuitFailure(testCfg, tt.err, tt.statusCode)
			assert.Equal(t, tt.wantCounts, got)
		})
	}
}

// TestCircuitBreakerIsSuccessful checks the gobreaker Settings.IsSuccessful
// hook directly: a nil error is success, a *circuitBreakerOutcome defers to
// its pre-computed verdict, and any other non-nil error (one that never
// went through performProxyRequestCore's classification) is a failure --
// matching gobreaker's own default.
func TestCircuitBreakerIsSuccessful(t *testing.T) {
	assert.True(t, circuitBreakerIsSuccessful(nil))
	assert.False(t, circuitBreakerIsSuccessful(errors.New("boom")))
	assert.False(t, circuitBreakerIsSuccessful(&circuitBreakerOutcome{err: errors.New("boom"), countsAsFailure: true}))
	assert.True(t, circuitBreakerIsSuccessful(&circuitBreakerOutcome{err: errors.New("boom"), countsAsFailure: false}))

	// Wrapping must not defeat classification.
	wrapped := fmt.Errorf("outer: %w", &circuitBreakerOutcome{err: errors.New("boom"), countsAsFailure: false})
	assert.True(t, circuitBreakerIsSuccessful(wrapped))
}

// TestCircuitBreakerFlagGatingIntegration exercises the full
// performProxyRequest path (not just the pure classifier) for the two
// fast, non-retried outcome kinds -- 4xx and timeout -- confirming both
// halves of the C5/C6 contract: (1) the client-visible result is an error
// either way (only the breaker bookkeeping differs), and (2) the circuit
// only trips when the matching CIRCUIT_TRIP_ON_* flag is true. The 5xx
// case is covered by TestClassifyCircuitFailure above and by
// TestCachingAndCircuitBreakerInteraction (integration_test.go); an
// end-to-end flag-gated 5xx test would additionally need to wait out the
// full 7-attempt retry backoff (~25s per request), which isn't worth the
// wall-clock cost here.
func (suite *CircuitBreakerTestSuite) TestCircuitBreakerFlagGatingIntegration() {
	runOnce := func(server *httptest.Server) error {
		app := fiber.New()
		reqCtx := &fasthttp.RequestCtx{}
		reqCtx.Request.SetRequestURI("/graphql")
		reqCtx.Request.Header.SetMethod("POST")
		reqCtx.Request.Header.SetContentType("application/json")
		reqCtx.Request.SetBody([]byte(`{"query":"query { test }"}`))
		ctx := app.AcquireCtx(reqCtx)
		defer app.ReleaseCtx(ctx)
		return performProxyRequest(ctx, server.URL)
	}

	suite.Run("4xx", func() {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"errors":[{"message":"bad request"}]}`))
		}))
		defer server.Close()

		for _, tripOn4xx := range []bool{true, false} {
			suite.Run(fmt.Sprintf("TripOn4xx=%v", tripOn4xx), func() {
				suite.outputBuffer.Reset()
				cfg.CircuitBreaker.MaxFailures = 1
				cfg.CircuitBreaker.Timeout = 5
				cfg.CircuitBreaker.TripOn4xx = tripOn4xx
				cfg.CircuitBreaker.ReturnCachedOnOpen = false
				cfg.Client.ClientTimeout = 5
				cfg.Client.FastProxyClient = createFasthttpClient(cfg)
				initCircuitBreaker(cfg)

				err := runOnce(server)
				assert.Error(suite.T(), err, "client must still see the 4xx failure regardless of the flag")

				wantState := gobreaker.StateClosed
				if tripOn4xx {
					wantState = gobreaker.StateOpen
				}
				assert.Equal(suite.T(), wantState.String(), cb.State().String(),
					"circuit state must match TripOn4xx=%v", tripOn4xx)
			})
		}
	})

	suite.Run("timeout", func() {
		// Handler sleeps well past the 1s client timeout below so every
		// attempt is a genuine client-side timeout, not a race.
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(2 * time.Second)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"data":{"ok":true}}`))
		}))
		defer server.Close()

		for _, tripOnTimeouts := range []bool{true, false} {
			suite.Run(fmt.Sprintf("TripOnTimeouts=%v", tripOnTimeouts), func() {
				suite.outputBuffer.Reset()
				cfg.CircuitBreaker.MaxFailures = 1
				cfg.CircuitBreaker.Timeout = 5
				cfg.CircuitBreaker.TripOnTimeouts = tripOnTimeouts
				cfg.CircuitBreaker.ReturnCachedOnOpen = false
				cfg.Client.ClientTimeout = 1
				cfg.Client.FastProxyClient = createFasthttpClient(cfg)
				initCircuitBreaker(cfg)

				err := runOnce(server)
				assert.Error(suite.T(), err, "client must still see the timeout failure regardless of the flag")

				wantState := gobreaker.StateClosed
				if tripOnTimeouts {
					wantState = gobreaker.StateOpen
				}
				assert.Equal(suite.T(), wantState.String(), cb.State().String(),
					"circuit state must match TripOnTimeouts=%v", tripOnTimeouts)
			})
		}
	})
}

// TestCircuitBreakerFailureRatioWindow covers the C12 fix: with
// Settings.Interval left at 0 (the pre-fix default), gobreaker never clears
// Counts while closed, so CIRCUIT_FAILURE_RATIO becomes lifetime-cumulative
// -- once enough historical successes have piled up, no realistic burst of
// new failures can move the ratio again. This test shrinks
// circuitBreakerCountWindow to a fraction of a second (restored via
// t.Cleanup) so the window boundary can be observed without sleeping for
// the real 60s default, then shows: (1) Counts actually reset for a new
// generation once the window elapses, instead of accumulating forever, and
// (2) a fresh in-window failure burst can trip the ratio check purely on
// its own -- which a lifetime-cumulative ratio, diluted by an earlier
// batch of successes, would not allow.
func (suite *CircuitBreakerTestSuite) TestCircuitBreakerFailureRatioWindow() {
	originalWindow := circuitBreakerCountWindow
	circuitBreakerCountWindow = 150 * time.Millisecond
	suite.T().Cleanup(func() { circuitBreakerCountWindow = originalWindow })

	// Disable the consecutive-failures trip path so only FailureRatio can
	// trip the breaker in this test; the circuit stays closed throughout.
	cfg.CircuitBreaker.MaxFailures = 1000
	cfg.CircuitBreaker.FailureRatio = 0.5
	cfg.CircuitBreaker.SampleSize = 4
	cfg.CircuitBreaker.Timeout = 5
	initCircuitBreaker(cfg)

	require := suite.Require()
	testErr := errors.New("test error")
	execute := func(fail bool) {
		_, _ = cb.Execute(func() (any, error) {
			if fail {
				return nil, testErr
			}
			return "ok", nil
		})
	}

	// Phase 1: a batch of mostly-successful requests with one failure --
	// well under the 50% ratio, so it must not trip, but it does
	// accumulate Requests/TotalFailures in Counts.
	for i := 0; i < 10; i++ {
		execute(false)
	}
	execute(true)
	require.Equal(gobreaker.StateClosed.String(), cb.State().String(), "low ratio must not trip")
	require.Equal(uint32(11), cb.Counts().Requests, "phase 1 requests should be counted")

	// Wait past the (shrunk) window. gobreaker clears Counts lazily, on
	// the next request, once Interval has elapsed while closed.
	time.Sleep(circuitBreakerCountWindow + 100*time.Millisecond)

	// This request both triggers the lazy clear and starts the new
	// generation's Counts.
	execute(false)
	assert.Equal(suite.T(), uint32(1), cb.Counts().Requests,
		"Counts must reset for the new window instead of accumulating across it (lifetime-cumulative regression)")

	// Phase 2: within the fresh window, a small failure burst alone (2
	// successes + 2 failures = 50% ratio, matching SampleSize) is enough to
	// trip. Under the pre-fix lifetime-cumulative behaviour, phase 1's 10
	// successes would still be diluting the ratio (13 requests, 3 failures
	// = 23%) and this burst alone would not trip the circuit.
	execute(false)
	execute(true)
	execute(true)
	assert.Equal(suite.T(), gobreaker.StateOpen.String(), cb.State().String(),
		"a fresh in-window failure burst must be able to trip the ratio check on its own")
}

// TestCircuitBackoffResetsOnSuccessfulRecovery is a suite method (needs
// cfg.Logger from SetupTest, which createStateChangeFunc's logging calls
// into) that exercises the C7 OnStateChange wiring end-to-end against a
// real *gobreaker.CircuitBreaker: consecutiveTrips increments on Closed ->
// Open, and resets to 0 exactly on the eventual HalfOpen -> Closed
// transition (a successful recovery).
//
// This breaker is constructed directly (bypassing initCircuitBreaker, whose
// Settings.Timeout only accepts whole seconds via cfg.CircuitBreaker.Timeout)
// with a deliberately tiny real Timeout so the Open -> HalfOpen -> Closed
// sequence completes after a bounded few-millisecond sleep -- gobreaker's
// own clock is not injectable, so some real wait is unavoidable to advance
// it. That sleep only advances gobreaker's OWN internal timeout to exercise
// the OnStateChange plumbing; it is unrelated to, and far shorter than, the
// C7 multiplier-scaled backoff window itself, which is verified
// deterministically (via injected trip counts and an injected `now`, no
// sleeping) by TestCircuitBackoffGateBlocks and
// TestCircuitBackoffLongerClosedDoorWithHigherMultiplier below.
func (suite *CircuitBreakerTestSuite) TestCircuitBackoffResetsOnSuccessfulRecovery() {
	cbMutex.Lock()
	cbBackoff.reset()
	cb = gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:          "recovery-test",
		MaxRequests:   1,
		Timeout:       5 * time.Millisecond,
		ReadyToTrip:   func(counts gobreaker.Counts) bool { return counts.ConsecutiveFailures >= 1 },
		OnStateChange: createStateChangeFunc(cfg),
		IsSuccessful:  circuitBreakerIsSuccessful,
	})
	cbMutex.Unlock()

	require := suite.Require()
	require.Equal(int64(0), cbBackoff.trips(), "fresh breaker must start with zero trips")

	// Closed -> Open: a single failure trips it immediately, no timing
	// dependency at all.
	_, _ = cb.Execute(func() (any, error) {
		return nil, &circuitBreakerOutcome{err: errors.New("boom"), countsAsFailure: true}
	})
	require.Equal(gobreaker.StateOpen.String(), cb.State().String(), "one failure must trip the breaker")
	require.Equal(int64(1), cbBackoff.trips(), "the trip must be recorded")
	require.False(cbBackoff.lastTrip().IsZero(), "lastTrip must be stamped")

	// Wait past gobreaker's own (tiny) Timeout so the next call observes the
	// lazy Open -> HalfOpen transition, then succeed to close it.
	time.Sleep(15 * time.Millisecond)
	_, err := cb.Execute(func() (any, error) { return "ok", nil })
	suite.NoError(err, "the half-open probe must succeed")
	suite.Equal(gobreaker.StateClosed.String(), cb.State().String(), "a successful probe must close the breaker")
	suite.Equal(int64(0), cbBackoff.trips(), "trip count must reset to 0 after a successful recovery to Closed")
}

// TestEffectiveCircuitBackoff is a pure-function table test of the C7
// backoff formula: effective = base * multiplier^(trips-1), capped at
// maxBackoff. It is the primary proof of default-preservation: with the
// default BackoffMultiplier=1.0, effective must equal base exactly, for
// every trip count, with no dependency on wall-clock or a live breaker.
func TestEffectiveCircuitBackoff(t *testing.T) {
	const base = 10 * time.Second

	t.Run("multiplier=1 (default) is constant=base for any trip count", func(t *testing.T) {
		for _, trips := range []int64{0, 1, 2, 5, 100} {
			got := effectiveCircuitBackoff(base, 1.0, 0, trips)
			assert.Equal(t, base, got, "trips=%d must equal base exactly under the default multiplier", trips)
		}
	})

	t.Run("multiplier=2 compounds from the second trip, capped at maxBackoff", func(t *testing.T) {
		const maxBackoff = 100 * time.Second
		tests := []struct {
			trips int64
			want  time.Duration
		}{
			{trips: 0, want: base},       // defensive: never tripped
			{trips: 1, want: base},       // first trip: exponent 0 -> multiplier^0=1
			{trips: 2, want: 2 * base},   // exponent 1
			{trips: 3, want: 4 * base},   // exponent 2
			{trips: 4, want: 8 * base},   // exponent 3, still under the 100s cap
			{trips: 5, want: maxBackoff}, // exponent 4 -> 16*base=160s, capped at 100s
		}
		for _, tt := range tests {
			got := effectiveCircuitBackoff(base, 2.0, maxBackoff, tt.trips)
			assert.Equal(t, tt.want, got, "trips=%d", tt.trips)
		}
	})

	t.Run("maxBackoff<=0 means uncapped", func(t *testing.T) {
		got := effectiveCircuitBackoff(base, 2.0, 0, 6) // exponent 5 -> 32*base
		assert.Equal(t, 32*base, got)
	})

	t.Run("maxBackoff<=0 with a high trip count is still bounded by the sanity ceiling", func(t *testing.T) {
		// finding C-overflow: CIRCUIT_MAX_BACKOFF_TIMEOUT=0 (uncapped) plus a
		// high enough consecutive-trip count makes
		// base*multiplier^(trips-1) overflow math.Pow's float64 arithmetic
		// (or the subsequent float64->time.Duration conversion) into a
		// huge-but-still-positive Duration -- on this platform the
		// out-of-range conversion saturates to math.MaxInt64 nanoseconds
		// (~292 years), which is NOT less than base, so the pre-existing
		// "effective < base" underflow guard alone does not catch it. The
		// hard ceiling must, regardless of trip count.
		got := effectiveCircuitBackoff(base, 2.0, 0, 1000)
		assert.Equal(t, maxSaneCircuitBackoff, got, "an uncapped config must never exceed the hard sanity ceiling")
		assert.LessOrEqual(t, got, maxSaneCircuitBackoff)
	})

	t.Run("misconfigured multiplier<1 is clamped to 1, never shrinks below base", func(t *testing.T) {
		got := effectiveCircuitBackoff(base, 0.5, 0, 3)
		assert.Equal(t, base, got)
	})

	t.Run("misconfigured maxBackoff below base never shrinks below base", func(t *testing.T) {
		got := effectiveCircuitBackoff(base, 2.0, 1*time.Second, 3)
		assert.Equal(t, base, got)
	})
}

// TestCircuitBackoffGateBlocks is a deterministic table test of
// circuitBackoffGateBlocks: trip count and elapsed time are both injected
// (via repeated cbBackoff.recordTrip calls and an explicit `now` parameter)
// rather than produced by sleeping, so the pass/fail boundary is exact and
// cannot flake. The breaker passed in is a freshly constructed, never-
// tripped *gobreaker.CircuitBreaker (State()==Closed) -- circuitBackoffGateBlocks
// only cares that breaker.State() != StateOpen, which a fresh breaker
// satisfies trivially without needing to wait out any real timeout.
func TestCircuitBackoffGateBlocks(t *testing.T) {
	t.Cleanup(func() { cbBackoff.reset() })

	const base = 10 * time.Second
	tripTime := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name       string
		trips      int64
		multiplier float64
		maxBackoff time.Duration
		elapsed    time.Duration
		wantBlock  bool
	}{
		{
			name:      "never tripped -> never blocks",
			trips:     0,
			elapsed:   0,
			wantBlock: false,
		},
		{
			name:       "first trip, multiplier=1, well within base -> blocks",
			trips:      1,
			multiplier: 1,
			elapsed:    1 * time.Second,
			wantBlock:  true,
		},
		{
			name:       "first trip, multiplier=1, base elapsed -> does not block",
			trips:      1,
			multiplier: 1,
			elapsed:    base,
			wantBlock:  false,
		},
		{
			name:       "second trip, multiplier=2, base elapsed but effective (2x base) not -> blocks",
			trips:      2,
			multiplier: 2,
			elapsed:    base + time.Second,
			wantBlock:  true,
		},
		{
			name:       "second trip, multiplier=2, effective (2x base) elapsed -> does not block",
			trips:      2,
			multiplier: 2,
			elapsed:    2 * base,
			wantBlock:  false,
		},
		{
			name:       "high trip count capped at maxBackoff, still within cap -> blocks",
			trips:      10,
			multiplier: 2,
			maxBackoff: 20 * time.Second,
			elapsed:    19 * time.Second,
			wantBlock:  true,
		},
		{
			name:       "high trip count capped at maxBackoff, cap elapsed -> does not block",
			trips:      10,
			multiplier: 2,
			maxBackoff: 20 * time.Second,
			elapsed:    20 * time.Second,
			wantBlock:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cbBackoff.reset()
			for i := int64(0); i < tt.trips; i++ {
				// Repeated recordTrip calls model consecutive trips; each
				// overwrites lastTripUnixNano, so the final timestamp is
				// exactly tripTime regardless of trip count.
				cbBackoff.recordTrip(tripTime)
			}

			testCfg := &config{}
			testCfg.CircuitBreaker.Timeout = int(base / time.Second)
			testCfg.CircuitBreaker.BackoffMultiplier = tt.multiplier
			testCfg.CircuitBreaker.MaxBackoffTimeout = int(tt.maxBackoff / time.Second)

			breaker := gobreaker.NewCircuitBreaker(gobreaker.Settings{Name: "gate-test"})
			now := tripTime.Add(tt.elapsed)

			got := circuitBackoffGateBlocks(testCfg, breaker, now)
			assert.Equal(t, tt.wantBlock, got)
		})
	}
}

// TestCircuitBackoffLongerClosedDoorWithHigherMultiplier directly proves
// the C7 requirement: for the SAME consecutive-trip count and the SAME
// elapsed time since the last trip (both injected, not slept), a higher
// BackoffMultiplier keeps the gate blocking while multiplier=1 (the
// default) has already let go -- i.e. progressive backoff strictly
// lengthens the closed-door window, and the default reproduces today's
// behaviour exactly.
func TestCircuitBackoffLongerClosedDoorWithHigherMultiplier(t *testing.T) {
	t.Cleanup(func() { cbBackoff.reset() })

	const base = 10 * time.Second
	const trips = 3 // third consecutive trip: multiplier exponent = trips-1 = 2
	tripTime := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	// Past base (so a multiplier=1 breaker would already let a probe
	// through) but well short of 2^2 * base = 40s.
	elapsed := base + time.Second
	now := tripTime.Add(elapsed)

	seed := func() {
		cbBackoff.reset()
		for i := 0; i < trips; i++ {
			cbBackoff.recordTrip(tripTime)
		}
	}

	cfgDefault := &config{}
	cfgDefault.CircuitBreaker.Timeout = int(base / time.Second)
	cfgDefault.CircuitBreaker.BackoffMultiplier = 1.0

	cfgProgressive := &config{}
	cfgProgressive.CircuitBreaker.Timeout = int(base / time.Second)
	cfgProgressive.CircuitBreaker.BackoffMultiplier = 2.0

	seed()
	blockedAtDefault := circuitBackoffGateBlocks(cfgDefault, gobreaker.NewCircuitBreaker(gobreaker.Settings{Name: "default"}), now)

	seed()
	blockedAtProgressive := circuitBackoffGateBlocks(cfgProgressive, gobreaker.NewCircuitBreaker(gobreaker.Settings{Name: "progressive"}), now)

	assert.False(t, blockedAtDefault, "multiplier=1.0 (default): base has already elapsed, gate must match today's behaviour and not block")
	assert.True(t, blockedAtProgressive, "multiplier=2.0: effective backoff has not yet elapsed at the same instant, gate must still block -- strictly longer closed-door time than the default")
}

// Start the test suite
func TestCircuitBreakerSuite(t *testing.T) {
	suite.Run(t, new(CircuitBreakerTestSuite))
}
