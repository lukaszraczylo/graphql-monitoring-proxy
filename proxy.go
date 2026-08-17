package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"go.opentelemetry.io/otel/trace"

	"github.com/avast/retry-go/v4"
	"github.com/gofiber/fiber/v3"
	libpack_cache "github.com/lukaszraczylo/graphql-monitoring-proxy/cache"
	libpack_logger "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
	libpack_monitoring "github.com/lukaszraczylo/graphql-monitoring-proxy/monitoring"
	libpack_tracing "github.com/lukaszraczylo/graphql-monitoring-proxy/tracing"
	"github.com/sony/gobreaker"
	"github.com/valyala/fasthttp"
)

// Errors related to circuit breaker
var (
	ErrCircuitOpen = errors.New("circuit breaker is open")
)

// Sentinel errors for the proxy request retry path. Grouped here so callers
// can use errors.Is for comparison instead of brittle string matching.
// Message text MUST match the historical fmt.Errorf strings — tests and
// callers may assert on .Error().
var (
	// errFiberCtxNilDuringRetry — fiber context dropped while retrying.
	errFiberCtxNilDuringRetry = errors.New("fiber context became nil during retry")
	// errFiberRespNil — fiber response object became nil mid-request.
	errFiberRespNil = errors.New("fiber response became nil")
	// errFiberCtxNil — fiber context was nil before the request started.
	errFiberCtxNil = errors.New("fiber context is nil")
)

// Default values for circuit breaker
const (
	defaultMaxRequestsInHalfOpen = 10 // Default maximum requests in half-open state
)

// circuitBreakerCountWindow bounds how long the circuit breaker accumulates
// Counts while CLOSED before clearing them (gobreaker Settings.Interval).
// gobreaker treats Interval<=0 as "never clear Counts during the closed
// state" (see sony/gobreaker Settings docs), which makes
// CIRCUIT_FAILURE_RATIO lifetime-cumulative instead of a recent-window
// ratio: after long uptime a handful of new failures can no longer move a
// ratio diluted by millions of past successes. A positive Interval makes
// gobreaker periodically clear Counts while closed, so the ratio measures a
// rolling window instead.
//
// The window must comfortably exceed the worst-case time a single logical
// request can take, not just typical request latency. gobreaker checks the
// interval lazily, both when a request starts (beforeRequest) and when it
// finishes (afterRequest); if the interval elapses *while* a request is
// in flight, afterRequest rolls Counts over to a new generation before
// recording that request's outcome, and then silently drops it (its
// generation no longer matches the one captured at the start). Each
// logical proxied request here can itself retry up to 7x with exponential
// backoff capped at 10s (performProxyRequestWithEnhancedRetries), i.e.
// ~25.5s worst case (0.5+1+2+4+8+10s) before it finally fails -- so with
// the library default MaxFailures=10, a run of consecutive failing
// requests can legitimately take ~4-5 minutes. A short window (e.g. the
// more typical 60s) reproduces exactly that: it was caught by
// TestCachingAndCircuitBreakerInteraction, where 3 sequential
// retry-exhausting failures (~76s total) spuriously left the circuit
// closed because the last failure's generation had already rolled over
// out from under it. 10 minutes keeps the ratio meaningfully "recent"
// (instead of lifetime-cumulative) while leaving multiples of headroom
// over that worst case. With that much headroom, a single request is very
// unlikely to straddle the window boundary; it is not impossible, since
// gobreaker checks the interval lazily against wall-clock time rather than
// reserving a slot up front. If one still does straddle it, gobreaker drops
// only that one outcome, which delays a trip by at most one request rather
// than breaking correctness.
//
// It is a package var (not const) purely so tests can shrink it to observe
// window-clearing behaviour without sleeping for the real duration.
var circuitBreakerCountWindow = 10 * time.Minute

// Global circuit breaker
var (
	cb      *gobreaker.CircuitBreaker
	cbMutex sync.RWMutex
)

// circuitBackoffState tracks the progressive open-state backoff (finding
// C7) for the global circuit breaker cb. gobreaker v1.0.0's own open-state
// Timeout (cb.timeout) is fixed at construction and has no setter, so it
// cannot be varied per-trip from outside the library. This state instead
// backs a decorator layered IN FRONT of gobreaker -- see
// circuitBackoffGateBlocks and its use in performProxyRequestCore -- that
// never reaches into gobreaker's internals; gobreaker's own Timeout keeps
// governing its internal half-open transition unchanged.
//
// consecutiveTrips counts transitions to StateOpen since the last
// successful recovery (a transition to StateClosed); lastTripUnixNano is
// the UnixNano() of the most recent such transition (0 = never tripped).
// Both are written from createStateChangeFunc's callback, which gobreaker
// invokes synchronously while still holding its own internal mutex (see
// sony/gobreaker CircuitBreaker.setState) -- so that callback must never
// call back into the breaker (e.g. cb.State()), which would deadlock on
// gobreaker's own non-reentrant lock. circuitBackoffGateBlocks reads this
// state on every proxied request. Atomics (not a mutex) keep both
// directions lock-free and sidestep any lock-ordering dependency on
// gobreaker's internal lock.
type circuitBackoffState struct {
	consecutiveTrips atomic.Int64
	lastTripUnixNano atomic.Int64
}

// recordTrip increments the consecutive-trip counter and stamps the current
// time. Called from createStateChangeFunc on transition to StateOpen.
func (s *circuitBackoffState) recordTrip(now time.Time) {
	s.consecutiveTrips.Add(1)
	s.lastTripUnixNano.Store(now.UnixNano())
}

// recordRecovery resets the consecutive-trip counter. Called from
// createStateChangeFunc on transition to StateClosed (a successful
// recovery), so the next trip after a clean close starts progressive
// backoff over again at the base timeout.
func (s *circuitBackoffState) recordRecovery() {
	s.consecutiveTrips.Store(0)
}

// reset clears all state. Called when a new circuit breaker is constructed
// (initCircuitBreaker) so a fresh breaker never inherits a stale trip count
// or timestamp left over from a previous one (e.g. across tests, or a
// config reload).
func (s *circuitBackoffState) reset() {
	s.consecutiveTrips.Store(0)
	s.lastTripUnixNano.Store(0)
}

// trips returns the current consecutive-trip count (0 = never tripped, or
// recovered since the last trip).
func (s *circuitBackoffState) trips() int64 {
	return s.consecutiveTrips.Load()
}

// lastTrip returns the time of the most recent transition to StateOpen, or
// the zero time.Time if the breaker has never tripped.
func (s *circuitBackoffState) lastTrip() time.Time {
	ns := s.lastTripUnixNano.Load()
	if ns == 0 {
		return time.Time{}
	}
	return time.Unix(0, ns)
}

// cbBackoff is the progressive-backoff tracker for the global circuit
// breaker cb -- one breaker, one tracker, mirroring the existing cb/cbMutex
// global pattern.
var cbBackoff circuitBackoffState

// Package-level substring tables used by isConnectionError / isTimeoutError.
// Hoisted to avoid per-call slice allocations on the hot path. All entries
// must be lower-case; callers lower-case the error string once before matching.
var (
	connectionErrorSubstrings = []string{
		"connection refused",
		"connection reset",
		"no route to host",
		"network is unreachable",
		"broken pipe",
		"connection closed",
		"eof",
		"no such host",
		"dial tcp",
		"dial udp",
	}

	timeoutErrorSubstrings = []string{
		"timeout",
		"deadline exceeded",
		"context deadline exceeded",
	}
)

// safeUint32 converts an int to uint32 safely, handling negative values and values exceeding uint32 max
func safeUint32(value int) uint32 {
	// Handle negative values
	if value < 0 {
		return 0
	}

	// Handle values exceeding uint32 max
	if value > math.MaxUint32 {
		return math.MaxUint32
	}

	return uint32(value)
}

// initCircuitBreaker initializes the circuit breaker with configured settings
func initCircuitBreaker(config *config) {
	// Only initialize if enabled
	if !config.CircuitBreaker.Enable {
		config.Logger.Info(&libpack_logger.LogMessage{
			Message: "Circuit breaker is disabled",
		})
		return
	}

	cbMutex.Lock()
	defer cbMutex.Unlock()

	// A freshly constructed breaker must never inherit a stale
	// consecutive-trip count or timestamp from a previous one (see
	// circuitBackoffState.reset doc) -- e.g. across tests, or a config
	// reload that reinitializes cb.
	cbBackoff.reset()

	// Initialize circuit breaker metrics
	InitializeCircuitBreakerMetrics(config.Monitoring)

	// Create circuit breaker settings
	cbSettings := gobreaker.Settings{
		Name:        "graphql-proxy-circuit",
		MaxRequests: safeMaxRequests(config.CircuitBreaker.MaxRequestsInHalfOpen),
		// Periodic reset while closed keeps CIRCUIT_FAILURE_RATIO windowed
		// instead of lifetime-cumulative; see circuitBreakerCountWindow doc.
		Interval:      circuitBreakerCountWindow,
		Timeout:       time.Duration(config.CircuitBreaker.Timeout) * time.Second,
		ReadyToTrip:   createTripFunc(config),
		OnStateChange: createStateChangeFunc(config),
		IsSuccessful:  circuitBreakerIsSuccessful,
	}

	// Initialize the circuit breaker
	cb = gobreaker.NewCircuitBreaker(cbSettings)

	config.Logger.Info(&libpack_logger.LogMessage{
		Message: "Circuit breaker initialized",
		Pairs: map[string]any{
			"max_failures":       config.CircuitBreaker.MaxFailures,
			"timeout_seconds":    config.CircuitBreaker.Timeout,
			"max_half_open_reqs": config.CircuitBreaker.MaxRequestsInHalfOpen,
		},
	})
}

// createTripFunc returns a function that determines when to trip the circuit
func createTripFunc(config *config) func(counts gobreaker.Counts) bool {
	return func(counts gobreaker.Counts) bool {
		// Check consecutive failures first
		if counts.ConsecutiveFailures >= safeUint32(config.CircuitBreaker.MaxFailures) {
			config.Logger.Warning(&libpack_logger.LogMessage{
				Message: "Circuit breaker tripped due to consecutive failures",
				Pairs: map[string]any{
					"consecutive_failures": counts.ConsecutiveFailures,
					"max_failures":         config.CircuitBreaker.MaxFailures,
					"total_requests":       counts.Requests,
				},
			})
			return true
		}

		// Check failure ratio if configured and enough samples
		if config.CircuitBreaker.FailureRatio > 0 &&
			config.CircuitBreaker.SampleSize > 0 &&
			counts.Requests >= safeUint32(config.CircuitBreaker.SampleSize) {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			if failureRatio >= config.CircuitBreaker.FailureRatio {
				config.Logger.Warning(&libpack_logger.LogMessage{
					Message: "Circuit breaker tripped due to failure ratio",
					Pairs: map[string]any{
						"failure_ratio":  failureRatio,
						"threshold":      config.CircuitBreaker.FailureRatio,
						"total_failures": counts.TotalFailures,
						"total_requests": counts.Requests,
					},
				})
				return true
			}
		}

		return false
	}
}

// circuitBreakerOutcome wraps a proxy failure with a pre-computed verdict on
// whether it should count against the circuit breaker's failure statistics.
// gobreaker.Settings.IsSuccessful only receives the error returned from the
// Execute closure and cannot see the per-request fiber context, so the
// classification (which needs the response status code) is made once, at
// the point where both are available -- see performProxyRequestCore -- and
// travels with the error. Error() and Unwrap() forward to the wrapped
// error, so retry policy, StatusCodeForError, and logging all see the
// original, unwrapped error.
type circuitBreakerOutcome struct {
	err             error
	countsAsFailure bool
}

func (e *circuitBreakerOutcome) Error() string { return e.err.Error() }
func (e *circuitBreakerOutcome) Unwrap() error { return e.err }

// circuitBreakerIsSuccessful is gobreaker's IsSuccessful hook (see
// Settings.IsSuccessful). It treats nil errors as success and defers to the
// pre-computed circuitBreakerOutcome.countsAsFailure verdict for classified
// proxy failures. Any other non-nil error -- one that never went through
// performProxyRequestCore's classification -- is treated as a failure,
// matching gobreaker's own default (defaultIsSuccessful) and pre-fix
// behaviour.
func circuitBreakerIsSuccessful(err error) bool {
	if err == nil {
		return true
	}
	var outcome *circuitBreakerOutcome
	if errors.As(err, &outcome) {
		return !outcome.countsAsFailure
	}
	return false
}

// classifyCircuitFailure decides whether a failed proxy attempt should count
// toward the circuit breaker's failure statistics, gated by
// CIRCUIT_TRIP_ON_TIMEOUTS / CIRCUIT_TRIP_ON_4XX / CIRCUIT_TRIP_ON_5XX.
// statusCode is the last response status observed for this request (0 when
// no HTTP response was ever received, e.g. a connection error).
//
// Every failure that isn't a classified timeout/4xx/5xx (connection errors,
// nil fiber context, etc.) always counts -- those outcomes have no
// dedicated flag and always tripped the breaker before this fix, so that
// baseline is preserved regardless of flag values.
func classifyCircuitFailure(config *config, err error, statusCode int) bool {
	if isTimeoutError(err) {
		return config.CircuitBreaker.TripOnTimeouts
	}

	// doProxyRequestWithTimeout's non-200 error always carries the real,
	// freshly-copied status in c.Response() at the moment it was produced
	// (see the comment in executeProxyAttempt). Guarding on the message
	// first means a stale statusCode left over from an earlier attempt
	// (e.g. this attempt was a connection error, which never copies a
	// response) can't be misread as an HTTP status classification.
	if strings.Contains(err.Error(), "non-200 response") {
		switch {
		case statusCode >= 500 && statusCode < 600:
			return config.CircuitBreaker.TripOn5xx
		case statusCode >= 400 && statusCode < 500:
			return config.CircuitBreaker.TripOn4xx
		}
	}

	return true
}

// createStateChangeFunc returns a function that handles circuit state changes
func createStateChangeFunc(config *config) func(name string, from gobreaker.State, to gobreaker.State) {
	return func(name string, from gobreaker.State, to gobreaker.State) {
		// Progressive backoff bookkeeping (C7). gobreaker invokes
		// OnStateChange synchronously while still holding its own internal
		// mutex, so this must only touch cbBackoff's atomics -- never call
		// back into the breaker (e.g. cb.State()), which would deadlock.
		switch to {
		case gobreaker.StateOpen:
			cbBackoff.recordTrip(time.Now())
		case gobreaker.StateClosed:
			cbBackoff.recordRecovery()
		}

		var stateValue float64
		var stateName string

		switch to {
		case gobreaker.StateOpen:
			stateValue = float64(libpack_monitoring.CircuitOpen)
			stateName = "open"
		case gobreaker.StateHalfOpen:
			stateValue = float64(libpack_monitoring.CircuitHalfOpen)
			stateName = "half-open"
		case gobreaker.StateClosed:
			stateValue = float64(libpack_monitoring.CircuitClosed)
			stateName = "closed"
		}

		// Update metrics using atomic operations to prevent race conditions
		// Use a separate atomic variable to track state instead of recreating gauges
		updateCircuitBreakerState(config, stateValue)

		// Log state change
		config.Logger.Info(&libpack_logger.LogMessage{
			Message: "Circuit breaker state changed",
			Pairs: map[string]any{
				"from": from.String(),
				"to":   to.String(),
				"name": name,
			},
		})

		// Use the new metrics system
		if cbMetrics != nil {
			// Replace hyphens with underscores to avoid validation errors
			safeStateName := strings.ReplaceAll(stateName, "-", "_")
			stateKey := fmt.Sprintf("circuit_state_%s", safeStateName)
			counter := cbMetrics.GetOrCreateFailCounter(config.Monitoring, stateKey)
			counter.Inc()
		}
	}
}

// gobreakerDefaultTimeout mirrors sony/gobreaker's own fallback for
// Settings.Timeout<=0 ("If Timeout is less than or equal to 0, the timeout
// value of the CircuitBreaker is set to 60 seconds" -- gobreaker.go
// NewCircuitBreaker). Duplicated here, since *gobreaker.CircuitBreaker
// exposes no getter for its internal timeout, so circuitBreakerBaseTimeout's
// result always matches what gobreaker itself actually uses as its
// open-state Timeout.
const gobreakerDefaultTimeout = 60 * time.Second

// circuitBreakerBaseTimeout returns the open-state Timeout gobreaker was
// constructed with (see initCircuitBreaker), applying the same <=0
// fallback gobreaker applies internally.
func circuitBreakerBaseTimeout(config *config) time.Duration {
	base := time.Duration(config.CircuitBreaker.Timeout) * time.Second
	if base <= 0 {
		return gobreakerDefaultTimeout
	}
	return base
}

// maxSaneCircuitBackoff hard-caps effectiveCircuitBackoff's result
// independently of the configured maxBackoff. When an operator sets
// CIRCUIT_MAX_BACKOFF_TIMEOUT=0 (uncapped) together with a
// BackoffMultiplier>1, a high enough consecutive-trip count makes
// base*multiplier^(trips-1) overflow math.Pow's float64 arithmetic (or the
// subsequent float64->time.Duration conversion) into a garbage or absurdly
// large Duration. This ceiling makes sure the breaker can never block
// requests for longer than a sane maximum, no matter what maxBackoff is
// configured to.
const maxSaneCircuitBackoff = 24 * time.Hour

// effectiveCircuitBackoff computes the progressive open-state backoff
// duration for the given consecutive-trip count (finding C7):
//
//	effective = base * multiplier^(trips-1), capped at maxBackoff.
//
// trips-1, not trips, so the FIRST trip (trips<=1) always backs off by
// exactly base -- identical to gobreaker's own fixed Timeout for that first
// open period. The multiplier only compounds starting from the SECOND
// consecutive trip. trips<=0 is defensive (a breaker that never tripped
// never reaches this path) and also returns base.
//
// With the default BackoffMultiplier=1.0, effective==base unconditionally
// for every trip count -- an exact fast path below skips the float
// round-trip entirely, so this is byte-identical to pre-C7 behaviour, not
// merely float-equal. A misconfigured multiplier<1 is clamped to 1, and the
// result is never allowed to fall below base or above a maxBackoff that is
// itself (misconfigured) smaller than base -- effective is always >= base,
// so this can only ever ADD closed-door time on top of gobreaker's own
// Timeout, never remove it. See maxSaneCircuitBackoff for the absolute
// ceiling applied regardless of maxBackoff.
func effectiveCircuitBackoff(base time.Duration, multiplier float64, maxBackoff time.Duration, trips int64) time.Duration {
	if base <= 0 || trips <= 1 {
		return base
	}

	if multiplier < 1 {
		multiplier = 1
	}
	if multiplier == 1 {
		return base
	}

	exponent := float64(trips - 1)
	raw := float64(base) * math.Pow(multiplier, exponent)

	// On amd64 the float64->int64 conversion for +Inf saturates to MinInt64
	// rather than MaxInt64, so the subsequent "< base" guard would return base
	// instead of the hard ceiling. Detect overflow before the cast.
	if math.IsInf(raw, 1) || math.IsNaN(raw) {
		return maxSaneCircuitBackoff
	}

	effective := time.Duration(raw)
	if effective < base {
		// Overflow or float rounding pushed us below base; never shorten.
		return base
	}

	// Hard sanity ceiling, applied before the configured cap below and
	// independent of it -- see maxSaneCircuitBackoff.
	if effective > maxSaneCircuitBackoff {
		effective = maxSaneCircuitBackoff
	}

	if maxBackoff > 0 && effective > maxBackoff {
		if maxBackoff < base {
			return base // misconfigured cap below base; never shorten
		}
		return maxBackoff
	}

	return effective
}

// circuitBackoffGateBlocks reports whether the C7 progressive-backoff gate
// should reject a request IN FRONT of gobreaker, without calling
// breaker.Execute -- and therefore without calling the backend or touching
// gobreaker's own Counts/generation.
//
// gobreaker's own Timeout still fully governs its internal half-open
// transition (see sony/gobreaker CircuitBreaker.currentState); this gate
// only adds an extra rejection window on top of it. It blocks exactly when
// both hold:
//   - the progressively-extended effective backoff (computed from the
//     consecutive-trip count) has not yet elapsed since the last trip, and
//   - gobreaker's own (shorter-or-equal) Timeout HAS already elapsed, i.e.
//     breaker.State() has lazily moved off StateOpen and would otherwise
//     let a half-open probe through.
//
// When gobreaker itself is still StateOpen, calling Execute would reject
// with ErrOpenState anyway, so no gate is needed there. When there have
// been no trips yet (trips<=0), this always returns false without touching
// breaker.State() at all, keeping the common (never-tripped) case
// allocation- and lock-free on the hot path. now is a parameter (not
// time.Now() read internally) so tests can drive this deterministically
// without wall-clock sleeps.
func circuitBackoffGateBlocks(config *config, breaker *gobreaker.CircuitBreaker, now time.Time) bool {
	trips := cbBackoff.trips()
	if trips <= 0 {
		return false
	}

	lastTrip := cbBackoff.lastTrip()
	if lastTrip.IsZero() {
		return false
	}

	base := circuitBreakerBaseTimeout(config)
	maxBackoff := time.Duration(config.CircuitBreaker.MaxBackoffTimeout) * time.Second
	effective := effectiveCircuitBackoff(base, config.CircuitBreaker.BackoffMultiplier, maxBackoff, trips)

	if now.Sub(lastTrip) >= effective {
		return false // extended window elapsed; defer to gobreaker as normal
	}

	return breaker.State() != gobreaker.StateOpen
}

// createFasthttpClient creates and configures a fasthttp client with optimized settings.
// The client is configured based on the provided configuration settings, with careful
// attention to performance and security considerations.
func createFasthttpClient(clientConfig *config) *fasthttp.Client {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: clientConfig.Client.DisableTLSVerify,
	}

	// Calculate timeout values, ensuring they're always positive
	clientTimeout := time.Duration(clientConfig.Client.ClientTimeout) * time.Second
	if clientTimeout <= 0 {
		clientTimeout = 30 * time.Second // Default timeout of 30 seconds
	}

	// Honour the per-side read/write timeouts (CLIENT_READ_TIMEOUT /
	// CLIENT_WRITE_TIMEOUT); they fall back to the client timeout when not
	// explicitly configured, so behaviour is unchanged for operators who only
	// set CLIENT_TIMEOUT.
	readTimeout := time.Duration(clientConfig.Client.ReadTimeout) * time.Second
	if readTimeout <= 0 {
		readTimeout = clientTimeout
	}
	writeTimeout := time.Duration(clientConfig.Client.WriteTimeout) * time.Second
	if writeTimeout <= 0 {
		writeTimeout = clientTimeout
	}

	// Create a custom dialer with timeout
	dialer := &fasthttp.TCPDialer{
		Concurrency:      1000,
		DNSCacheDuration: time.Hour,
	}

	client := &fasthttp.Client{
		Name:                     "graphql_proxy",
		NoDefaultUserAgentHeader: true,
		TLSConfig:                tlsConfig,
		// Control connection pool size to prevent overwhelming backend services
		MaxConnsPerHost: clientConfig.Client.MaxConnsPerHost,
		// Configure timeouts to handle different network scenarios
		// Setting all timeout-related parameters to ensure proper timeout behavior
		Dial: func(addr string) (net.Conn, error) {
			return dialer.DialTimeout(addr, clientTimeout)
		},
		ReadTimeout:                   readTimeout,
		WriteTimeout:                  writeTimeout,
		MaxIdleConnDuration:           time.Duration(clientConfig.Client.MaxIdleConnDuration) * time.Second,
		MaxConnDuration:               clientTimeout,
		DisableHeaderNamesNormalizing: false,
		// Performance tuning
		ReadBufferSize:         4096,
		WriteBufferSize:        4096,
		MaxResponseBodySize:    1024 * 1024 * 10, // 10MB max response size
		DisablePathNormalizing: false,
	}

	// Initialize connection pool manager
	InitializeConnectionPool(client)

	return client
}

// proxyTheRequest handles the request proxying logic.
func proxyTheRequest(c fiber.Ctx, currentEndpoint string) error {
	// Record request for RPS tracking
	if rpsTracker := GetRPSTracker(); rpsTracker != nil {
		rpsTracker.RecordRequest()
	}

	// Setup tracing if enabled
	var span trace.Span
	var ctx context.Context

	if cfg.Tracing.Enable && tracer != nil {
		ctx = setupTracing(c)
		span, _ = tracer.StartSpan(ctx, "proxy_request")
		defer span.End()
	}

	// Check if URL is allowed
	if !checkAllowedURLs(c) {
		if ifNotInTest() {
			cfg.Monitoring.Increment(libpack_monitoring.MetricsSkipped, nil)
		}
		return fmt.Errorf("request blocked - not allowed URL: %s", c.Path())
	}

	// Construct and validate proxy URL
	proxyURL := currentEndpoint + c.OriginalURL()
	if _, err := url.Parse(proxyURL); err != nil {
		return fmt.Errorf("invalid URL: %v", err)
	}

	// Log request details in debug mode
	if cfg.LogLevel == "DEBUG" {
		logDebugRequest(c)
	}

	// Perform the proxy request with retries
	if err := performProxyRequest(c, proxyURL); err != nil {
		if ifNotInTest() {
			cfg.Monitoring.Increment(libpack_monitoring.MetricsFailed, nil)
		}
		return err
	}

	// Log response details in debug mode
	if cfg.LogLevel == "DEBUG" {
		logDebugResponse(c)
	}

	// Handle gzipped responses
	if err := handleGzippedResponse(c); err != nil {
		return err
	}

	// Final status check
	if c.Response().StatusCode() != fiber.StatusOK {
		if ifNotInTest() {
			cfg.Monitoring.Increment(libpack_monitoring.MetricsFailed, nil)
		}
		return fmt.Errorf("received non-200 response from the GraphQL server: %d", c.Response().StatusCode())
	}

	// Remove server header for security
	c.Response().Header.Del(fiber.HeaderServer)
	return nil
}

// setupTracing extracts and sets up tracing context from request headers
func setupTracing(c fiber.Ctx) context.Context {
	ctx := context.Background()

	if !cfg.Tracing.Enable || tracer == nil {
		return ctx
	}

	// Extract trace information from header
	if traceHeader := c.Get("X-Trace-Span"); traceHeader != "" {
		spanInfo, err := libpack_tracing.ParseTraceHeader(traceHeader)
		if err != nil {
			cfg.Logger.Warning(&libpack_logger.LogMessage{
				Message: "Failed to parse trace header",
				Pairs:   map[string]any{"error": err.Error()},
			})
		} else if spanCtx, err := tracer.ExtractSpanContext(spanInfo); err == nil {
			ctx = trace.ContextWithSpanContext(ctx, spanCtx)
		}
	}

	return ctx
}

// performProxyRequest executes the proxy request with retries, circuit breaker, and request coalescing
func performProxyRequest(c fiber.Ctx, proxyURL string) error {
	// Reuse the cache key already computed by handleCaching (stored on the
	// request context) to avoid a second MD5 of the full request body per
	// proxied request. Fall back to computing it when this function is
	// reached without going through handleCaching (direct callers).
	cacheKey, _ := c.Locals("query_cache_hash").(string)
	if cacheKey == "" {
		// Error intentionally ignored: the primary rejection for an
		// invalid/forged token already happened in processGraphQLRequest
		// (it 401s before the request ever reaches here). On a verify
		// failure extractUserInfo returns ("-", "-", err), so a forged
		// token still only maps into the shared anonymous cache bucket
		// below, never into a genuine, verified user's own bucket
		// (GHSA-9gqw-h2rw-44wv).
		userID, userRole, _ := extractUserInfo(c)
		cacheKey = libpack_cache.CalculateHash(c, userID, userRole)
	}

	// Check if request coalescing is enabled
	rc := GetRequestCoalescer()
	if rc != nil && cfg.RequestCoalescing.Enable {
		// Use request coalescing to deduplicate identical concurrent requests
		response, err := rc.Do(cacheKey, func() (*CoalescedResponse, error) {
			// Execute the actual proxy request
			proxyErr := performProxyRequestCore(c, proxyURL, cacheKey)

			// Capture the response for coalescing
			if proxyErr != nil {
				return &CoalescedResponse{
					Err:        proxyErr,
					StatusCode: c.Response().StatusCode(),
				}, proxyErr
			}

			return &CoalescedResponse{
				Body:       c.Response().Body(),
				StatusCode: c.Response().StatusCode(),
				// Headers intentionally left nil; not populated or read anywhere.
			}, nil
		})

		// Check for error from rc.Do (though it typically returns nil)
		if err != nil {
			return err
		}

		// Check for error stored in the response (for coalesced requests)
		if response != nil && response.Err != nil {
			return response.Err
		}

		// For coalesced requests (not the primary), we need to copy the response
		if response != nil && response.Body != nil && len(response.Body) > 0 {
			// Only set response if this is a coalesced request (body would be empty otherwise)
			if len(c.Response().Body()) == 0 {
				c.Response().SetStatusCode(response.StatusCode)
				c.Response().SetBody(response.Body)
			}
		}

		return nil
	}

	// No coalescing - execute directly
	return performProxyRequestCore(c, proxyURL, cacheKey)
}

// performProxyRequestCore executes the proxy request with retries and circuit breaker
// This is the core implementation used by both direct calls and coalesced requests
func performProxyRequestCore(c fiber.Ctx, proxyURL string, cacheKey string) error {
	// If circuit breaker is not enabled, use the original method
	if !cfg.CircuitBreaker.Enable || cb == nil {
		return performProxyRequestWithRetries(c, proxyURL)
	}

	var err error
	if circuitBackoffGateBlocks(cfg, cb, time.Now()) {
		// Progressive backoff (C7): gobreaker's own Timeout has already
		// elapsed -- it would otherwise let a half-open probe through -- but
		// the extended per-trip backoff window has not. Reject exactly like
		// an open breaker, without calling cb.Execute: no backend call, and
		// gobreaker's own Counts/generation are left untouched. gobreaker
		// stays idle in half-open until the extended window also elapses,
		// at which point the Execute call below performs the real probe.
		err = gobreaker.ErrOpenState
	} else {
		// Execute request through circuit breaker
		_, err = cb.Execute(func() (any, error) {
			// Execute the request with retries
			err := performProxyRequestWithRetries(c, proxyURL)
			if err != nil {
				// Log error that could potentially trip the circuit
				cfg.Logger.Warning(&libpack_logger.LogMessage{
					Message: "Error in circuit-protected request",
					Pairs: map[string]any{
						"path":  c.Path(),
						"error": err.Error(),
					},
				})
				// Classify the failure so CIRCUIT_TRIP_ON_TIMEOUTS/_4XX/_5XX gate
				// whether it counts toward the breaker's failure statistics
				// (see circuitBreakerIsSuccessful). The error returned to the
				// caller is unchanged either way -- only the breaker bookkeeping
				// is affected.
				countsAsFailure := classifyCircuitFailure(cfg, err, c.Response().StatusCode())
				if countsAsFailure {
					// Mirrors MetricsCircuitSuccessful below: incremented once per
					// Execute call, only for outcomes that actually count against
					// the breaker's failure statistics, so a single request is
					// never double-counted.
					cfg.Monitoring.Increment(libpack_monitoring.MetricsCircuitFailed, nil)
				}
				return nil, &circuitBreakerOutcome{err: err, countsAsFailure: countsAsFailure}
			}

			// Request was successful
			cfg.Monitoring.Increment(libpack_monitoring.MetricsCircuitSuccessful, nil)
			return nil, nil
		})
	}

	// If the circuit is open, implement graceful degradation
	if err == gobreaker.ErrOpenState {
		cfg.Monitoring.Increment(libpack_monitoring.MetricsCircuitRejected, nil)
		// If cache fallback is disabled, return the original circuit breaker error
		if !cfg.CircuitBreaker.ReturnCachedOnOpen {
			return gobreaker.ErrOpenState
		}
		return handleCircuitOpenGracefulDegradation(c, cacheKey)
	}

	return err
}

// performProxyRequestWithRetries executes the proxy request with retries
// This is the original implementation extracted for reuse
func performProxyRequestWithRetries(c fiber.Ctx, proxyURL string) error {
	// Check backend health first if available
	healthMgr := GetBackendHealthManager()
	if healthMgr != nil && !healthMgr.IsHealthy() {
		// If backend is unhealthy, use more aggressive retry strategy
		return performProxyRequestWithEnhancedRetries(c, proxyURL, true)
	}

	return performProxyRequestWithEnhancedRetries(c, proxyURL, false)
}

// executeProxyAttempt performs a single proxy attempt with error handling
func executeProxyAttempt(c fiber.Ctx, req *fasthttp.Request) error {
	// Additional safety check inside retry loop
	if c == nil {
		return retry.Unrecoverable(errFiberCtxNilDuringRetry)
	}

	// Get connection pool manager for stats tracking
	poolMgr := GetConnectionPoolManager()

	// Execute the proxy request. req was built once by
	// performProxyRequestWithEnhancedRetries and is reused unmodified across
	// every attempt (see the comment there for why that's safe).
	proxyErr := doProxyRequestWithTimeout(c, req, cfg.Client.FastProxyClient)
	if proxyErr != nil {
		// Check if this is a connection error
		if isConnectionError(proxyErr) {
			notifyHealthManager(false)
			// Track connection failure
			if poolMgr != nil {
				poolMgr.RecordConnectionFailure()
			}
			return proxyErr // Connection errors are retryable
		}

		// Check if this is a timeout error - don't retry timeouts
		if isTimeoutError(proxyErr) {
			return retry.Unrecoverable(proxyErr)
		}

		// The server responded with an HTTP status code. Decide retry policy
		// from the status: 5xx and 429 are retryable; other 4xx client
		// errors fail fast instead of burning backoff retries.
		// doProxyRequestWithTimeout already copied the status into c.Response().
		if strings.Contains(proxyErr.Error(), "non-200 response") {
			statusCode := c.Response().StatusCode()
			shouldRetry, _ := isRetryableStatusCode(statusCode)
			if !shouldRetry {
				return retry.Unrecoverable(proxyErr)
			}
			// Track as a failure for retryable HTTP errors (5xx, 429)
			if poolMgr != nil {
				poolMgr.RecordConnectionFailure()
			}
		}
		return proxyErr
	}

	// Safety check before accessing response (c is already validated at function entry)
	if c.Response() == nil {
		return retry.Unrecoverable(errFiberRespNil)
	}

	// Check status code and determine retry strategy
	statusCode := c.Response().StatusCode()
	shouldRetry, err := isRetryableStatusCode(statusCode)

	if err == nil {
		// Success case
		notifyHealthManager(true)
		// Track successful connection
		if poolMgr != nil {
			poolMgr.RecordConnectionSuccess()
		}
		return nil
	}

	if shouldRetry {
		// Track connection failure for retryable errors (5xx, etc)
		if poolMgr != nil {
			poolMgr.RecordConnectionFailure()
		}
		return err // Retryable error
	}

	return err // Non-retryable error (already wrapped with retry.Unrecoverable)
}

// performProxyRequestWithEnhancedRetries executes the proxy request with intelligent retry strategy
func performProxyRequestWithEnhancedRetries(c fiber.Ctx, proxyURL string, backendUnhealthy bool) error {
	// Safety check for nil context
	if c == nil {
		return errFiberCtxNil
	}

	// Build the outbound request once, before the retry loop, instead of
	// re-copying c.Request() (headers + body) on every attempt. A retry
	// replays the exact same request, so the copy-in is invariant across
	// attempts; only the response is fresh per attempt (doProxyRequestWithTimeout
	// acquires/releases its own *fasthttp.Response each call). This previously
	// ran c.Request().CopyTo(req) up to 7x per proxied request.
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	c.Request().CopyTo(req)
	req.SetRequestURI(proxyURL)

	var attempts uint
	var initialDelay time.Duration
	var maxDelayTime time.Duration

	// Read per-request base and max delay from config; fall back to coded
	// defaults when zero (e.g. during very early startup before parseConfig).
	baseDelayMs := cfg.Client.RetryBaseDelayMs
	if baseDelayMs <= 0 {
		baseDelayMs = 500
	}
	maxDelayMs := cfg.Client.RetryMaxDelayMs
	if maxDelayMs <= 0 {
		maxDelayMs = 10000
	}

	if backendUnhealthy {
		// Backend is known to be unhealthy; fail fast with fewer attempts.
		// The circuit breaker should handle it, so the max delay matches
		// the normal path (avoids invisible behaviour differences).
		attempts = 3
		initialDelay = time.Duration(baseDelayMs) * time.Millisecond
		maxDelayTime = time.Duration(maxDelayMs) * time.Millisecond
	} else {
		// Normal retry strategy
		attempts = 7
		initialDelay = time.Duration(baseDelayMs) * time.Millisecond
		maxDelayTime = time.Duration(maxDelayMs) * time.Millisecond
	}

	return retry.Do(
		func() error {
			return executeProxyAttempt(c, req)
		},
		retry.Attempts(attempts),
		retry.DelayType(retry.BackOffDelay),
		retry.Delay(initialDelay),
		retry.MaxDelay(maxDelayTime),
		retry.OnRetry(func(n uint, err error) {
			cfg.Logger.Warning(&libpack_logger.LogMessage{
				Message: "Retrying the request",
				Pairs: map[string]any{
					"path":              c.Path(),
					"attempt":           n + 1,
					"max_attempts":      attempts,
					"error":             err.Error(),
					"error_type":        fmt.Sprintf("%T", err),
					"is_timeout":        strings.Contains(strings.ToLower(err.Error()), "timeout"),
					"is_connection":     isConnectionError(err),
					"backend_unhealthy": backendUnhealthy,
				},
			})
		}),
		retry.LastErrorOnly(true),
		retry.RetryIf(func(err error) bool {
			// Unrecoverable errors (non-retryable 4xx, timeouts) must not be
			// retried. retry-go's finite-attempt loop only consults RetryIf,
			// so without this the default IsRecoverable gate is lost.
			if !retry.IsRecoverable(err) {
				return false
			}
			// Don't retry if context is cancelled or context is nil
			if c == nil {
				return false
			}

			// Safely check if context is done/cancelled
			// Note: fasthttp.RequestCtx.Done() can panic if not properly initialized
			// If we panic, don't retry (maintains backward compatibility with test behavior)
			shouldRetry := true
			func() {
				defer func() {
					if r := recover(); r != nil {
						// If we panic accessing context, don't retry
						// This typically happens in test scenarios with mock contexts
						shouldRetry = false
					}
				}()
				ctx := c.Context()
				if ctx == nil {
					return
				}
				select {
				case <-ctx.Done():
					shouldRetry = false
				default:
				}
			}()

			if !shouldRetry {
				return false
			}

			// Check retry budget before allowing retry
			if rb := GetRetryBudget(); rb != nil {
				if !rb.AllowRetry() {
					cfg.Logger.Warning(&libpack_logger.LogMessage{
						Message: "Retry denied by budget",
						Pairs: map[string]any{
							"path":  c.Path(),
							"error": err.Error(),
						},
					})
					return false
				}
			}
			return true
		}),
	)
}

// isConnectionError checks if the error is a connection-related error
func isConnectionError(err error) bool {
	if err == nil {
		return false
	}

	errStr := strings.ToLower(err.Error())
	for _, connErr := range connectionErrorSubstrings {
		if strings.Contains(errStr, connErr) {
			return true
		}
	}

	return false
}

// isTimeoutError checks if the error is a timeout-related error
func isTimeoutError(err error) bool {
	if err == nil {
		return false
	}
	errStr := strings.ToLower(err.Error())
	for _, tErr := range timeoutErrorSubstrings {
		if strings.Contains(errStr, tErr) {
			return true
		}
	}
	return false
}

// isRetryableStatusCode determines if an HTTP status code should trigger a retry
func isRetryableStatusCode(statusCode int) (bool, error) {
	// Don't retry client errors (4xx) except for specific cases
	if statusCode >= 400 && statusCode < 500 {
		// Retry on 429 (rate limit) and 503 (service unavailable - misclassified as 4xx)
		if statusCode == 429 || statusCode == 503 {
			return true, fmt.Errorf("retryable status code: %d", statusCode)
		}
		// Other 4xx errors are not retryable
		return false, retry.Unrecoverable(fmt.Errorf("client error: %d", statusCode))
	}

	// Retry on 5xx errors
	if statusCode >= 500 {
		return true, fmt.Errorf("server error: %d", statusCode)
	}

	// Success for 2xx and 3xx
	if statusCode >= 200 && statusCode < 400 {
		return false, nil // No error, no retry needed
	}

	return true, fmt.Errorf("unexpected status code: %d", statusCode)
}

// notifyHealthManager notifies the backend health manager of request success or failure
func notifyHealthManager(success bool) {
	if healthMgr := GetBackendHealthManager(); healthMgr != nil {
		healthMgr.updateHealthStatus(success)
	}
}

// handleCircuitOpenGracefulDegradation handles requests when the circuit breaker is open
func handleCircuitOpenGracefulDegradation(c fiber.Ctx, cacheKey string) error {
	// Try to serve from cache if configured and available
	if cfg.CircuitBreaker.ReturnCachedOnOpen {
		if cachedResponse := libpack_cache.CacheLookup(cacheKey); cachedResponse != nil {
			cfg.Logger.Info(&libpack_logger.LogMessage{
				Message: "Circuit open - serving from cache",
				Pairs: map[string]any{
					"path": c.Path(),
				},
			})

			// Set response from cache
			c.Response().SetBody(cachedResponse)
			c.Response().SetStatusCode(fiber.StatusOK)

			// Mark as cache hit since we're serving from cache
			cfg.Monitoring.Increment(libpack_monitoring.MetricsCacheHit, nil)
			cfg.Monitoring.Increment(libpack_monitoring.MetricsCircuitFallbackSuccess, nil)

			return nil
		}
	}

	// No cached response available - provide helpful error response
	cfg.Logger.Warning(&libpack_logger.LogMessage{
		Message: "Circuit open - no cached response available",
		Pairs: map[string]any{
			"path": c.Path(),
		},
	})

	cfg.Monitoring.Increment(libpack_monitoring.MetricsCircuitFallbackFailed, nil)

	return ErrCircuitOpen
}

// doProxyRequestWithTimeout performs a proxy request with proper timeout
// handling. req is the already-built outbound request (see
// performProxyRequestWithEnhancedRetries); only the response is
// acquired/released per call, so each retry attempt gets fresh
// per-attempt response state while reusing the same request.
func doProxyRequestWithTimeout(c fiber.Ctx, req *fasthttp.Request, client *fasthttp.Client) error {
	// Calculate timeout from client configuration
	clientTimeout := time.Duration(cfg.Client.ClientTimeout) * time.Second
	if clientTimeout <= 0 {
		clientTimeout = 30 * time.Second
	}

	// Acquire a fresh response object for this attempt
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	// Perform the request with timeout
	err := client.DoTimeout(req, resp, clientTimeout)
	if err != nil {
		return err
	}

	// Copy response back to fiber context
	resp.CopyTo(c.Response())

	// Check for non-200 responses and return error for tests
	if c.Response().StatusCode() != fiber.StatusOK {
		return fmt.Errorf("received non-200 response: %d", c.Response().StatusCode())
	}

	return nil
}

// handleGzippedResponse decompresses gzipped responses
func handleGzippedResponse(c fiber.Ctx) error {
	if !bytes.EqualFold(c.Response().Header.Peek("Content-Encoding"), []byte("gzip")) {
		return nil
	}

	// Use pooled gzip reader
	reader, err := GetGzipReader(bytes.NewReader(c.Response().Body()))
	if err != nil {
		cfg.Logger.Error(&libpack_logger.LogMessage{
			Message: "Failed to create gzip reader",
			Pairs:   map[string]any{"error": err.Error()},
		})
		return err
	}
	defer func() {
		// Return reader to pool
		PutGzipReader(reader)
	}()

	// Use pooled buffer for reading
	buf := GetHTTPBuffer()
	defer PutHTTPBuffer(buf)

	// Read decompressed data into pooled buffer
	_, err = io.Copy(buf, reader)
	if err != nil {
		cfg.Logger.Error(&libpack_logger.LogMessage{
			Message: "Failed to decompress response",
			Pairs:   map[string]any{"error": err.Error()},
		})
		return err
	}

	// Get decompressed data
	decompressed := buf.Bytes()

	// Update response
	c.Response().SetBody(decompressed)
	c.Response().Header.Del("Content-Encoding")
	return nil
}

// logDebugRequest logs the request details when in debug mode with sanitization.
func logDebugRequest(c fiber.Ctx) {
	contentType := string(c.Request().Header.ContentType())
	sanitizedBody := sanitizeForLogging(c.Body(), contentType)
	sanitizedHeaders := sanitizeHeaders(convertHeaders(c.GetReqHeaders()))

	cfg.Logger.Debug(&libpack_logger.LogMessage{
		Message: "Proxying the request",
		Pairs: map[string]any{
			"path":         c.Path(),
			"body":         sanitizedBody,
			"headers":      sanitizedHeaders,
			"request_uuid": c.Locals("request_uuid"),
		},
	})
}

// logDebugResponse logs the response details when in debug mode with sanitization.
func logDebugResponse(c fiber.Ctx) {
	contentType := string(c.Response().Header.ContentType())
	sanitizedBody := sanitizeForLogging(c.Response().Body(), contentType)
	sanitizedHeaders := sanitizeHeaders(convertHeaders(c.GetRespHeaders()))

	cfg.Logger.Debug(&libpack_logger.LogMessage{
		Message: "Received proxied response",
		Pairs: map[string]any{
			"path":          c.Path(),
			"response_body": sanitizedBody,
			"response_code": c.Response().StatusCode(),
			"headers":       sanitizedHeaders,
			"request_uuid":  c.Locals("request_uuid"),
		},
	})
}

// safeMaxRequests converts MaxRequestsInHalfOpen safely to uint32, providing a fallback value if out of bounds
func safeMaxRequests(maxRequestsInHalfOpen int) uint32 {
	// Check if value is invalid (negative or too large)
	if maxRequestsInHalfOpen < 0 || maxRequestsInHalfOpen > math.MaxUint32 {
		// Log warning and return a default value
		if cfg != nil && cfg.Logger != nil {
			cfg.Logger.Warning(&libpack_logger.LogMessage{
				Message: "Invalid MaxRequestsInHalfOpen value, using default",
				Pairs: map[string]any{
					"requested_value": maxRequestsInHalfOpen,
					"default_value":   defaultMaxRequestsInHalfOpen,
				},
			})
		}
		return uint32(defaultMaxRequestsInHalfOpen)
	}

	return uint32(maxRequestsInHalfOpen)
}

// updateCircuitBreakerState safely updates the circuit breaker state using atomic operations
func updateCircuitBreakerState(config *config, stateValue float64) {
	// Update the state atomically using the new metrics system
	if cbMetrics != nil {
		cbMetrics.UpdateState(stateValue)
	}
}
