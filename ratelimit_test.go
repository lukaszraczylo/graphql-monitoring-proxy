package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/goccy/go-json"
	goratecounter "github.com/lukaszraczylo/go-ratecounter"
	libpack_logger "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
)

func (suite *Tests) Test_loadRatelimitConfig() {
	// Setup
	cfg = &config{}
	parseConfig()
	cfg.Logger = libpack_logger.New()

	// Create a temporary test ratelimit.json file
	tempDir := os.TempDir()
	testConfigPath := filepath.Join(tempDir, "test_ratelimit.json")

	testConfig := struct {
		RateLimit map[string]RateLimitConfig `json:"ratelimit"`
	}{
		RateLimit: map[string]RateLimitConfig{
			"admin": {
				Interval: 1 * time.Second,
				Req:      100,
			},
			"user": {
				Interval: 1 * time.Second,
				Req:      10,
			},
		},
	}

	configData, err := json.Marshal(testConfig)
	suite.NoError(err)

	err = os.WriteFile(testConfigPath, configData, 0o644)
	suite.NoError(err)
	defer func() { _ = os.Remove(testConfigPath) }()

	// Test loading config from custom path
	suite.Run("load from custom path", func() {
		// Clear existing rate limits
		rateLimitMu.Lock()
		rateLimits = make(map[string]RateLimitConfig)
		rateLimitMu.Unlock()

		err := loadConfigFromPath(testConfigPath)
		suite.NoError(err)

		// Verify rate limits were loaded
		rateLimitMu.RLock()
		defer rateLimitMu.RUnlock()

		suite.Equal(2, len(rateLimits))
		suite.Contains(rateLimits, "admin")
		suite.Contains(rateLimits, "user")
		suite.Equal(100, rateLimits["admin"].Req)
		suite.Equal(10, rateLimits["user"].Req)
		suite.NotNil(rateLimits["admin"].RateCounterTicker)
		suite.NotNil(rateLimits["user"].RateCounterTicker)
	})

	// Test loading config from non-existent path
	suite.Run("load from non-existent path", func() {
		err := loadConfigFromPath("/non/existent/path.json")
		suite.Error(err)
	})

	// Test loading config with invalid JSON
	suite.Run("load invalid JSON", func() {
		invalidPath := filepath.Join(tempDir, "invalid_ratelimit.json")
		err := os.WriteFile(invalidPath, []byte("{invalid json}"), 0o644)
		suite.NoError(err)
		defer func() { _ = os.Remove(invalidPath) }()

		err = loadConfigFromPath(invalidPath)
		suite.Error(err)
	})

	// Test with a temporary ratelimit.json file in the current directory
	suite.Run("load from current directory", func() {
		// Create a temporary ratelimit.json in current directory
		currentDirPath := "./ratelimit.json"
		err := os.WriteFile(currentDirPath, configData, 0o644)
		suite.NoError(err)
		defer func() { _ = os.Remove(currentDirPath) }()

		// Clear existing rate limits
		rateLimitMu.Lock()
		rateLimits = make(map[string]RateLimitConfig)
		rateLimitMu.Unlock()

		// This should find the file in the current directory
		err = loadRatelimitConfig()
		suite.NoError(err)

		// Verify rate limits were loaded
		rateLimitMu.RLock()
		defer rateLimitMu.RUnlock()

		suite.Equal(2, len(rateLimits))
	})

	// Test with all files missing
	suite.Run("all files missing", func() {
		// Save the original load function and restore it when done
		originalLoadFunc := loadConfigFunc
		defer func() {
			loadConfigFunc = originalLoadFunc
		}()

		// Replace with a mock function that always returns "file does not exist" error
		loadConfigFunc = func(string) error {
			return fmt.Errorf("file does not exist")
		}

		// Clear existing rate limits
		rateLimitMu.Lock()
		rateLimits = make(map[string]RateLimitConfig)
		rateLimitMu.Unlock()

		// This should fail as our mock returns errors for all paths
		err = loadRatelimitConfig()
		suite.Error(err)

		// The error should be a RateLimitConfigError
		configErr, ok := err.(*RateLimitConfigError)
		suite.True(ok, "Expected *RateLimitConfigError but got %T", err)

		// All path errors should contain our mock error message
		for _, errMsg := range configErr.PathErrors {
			suite.Equal("file does not exist", errMsg)
		}
	})
}

func (suite *Tests) Test_rateLimitedRequest() {
	// Setup
	cfg = &config{}
	parseConfig()
	cfg.Logger = libpack_logger.New()

	// Create test rate limits
	rateLimitMu.Lock()
	rateLimits = make(map[string]RateLimitConfig)

	// Admin role with high limit
	adminCounter := goratecounter.NewRateCounter().WithConfig(goratecounter.RateCounterConfig{
		Interval: 1 * time.Second,
	})
	rateLimits["admin"] = RateLimitConfig{
		RateCounterTicker: adminCounter,
		Interval:          1 * time.Second,
		Req:               100,
	}

	// User role with low limit
	userCounter := goratecounter.NewRateCounter().WithConfig(goratecounter.RateCounterConfig{
		Interval: 1 * time.Second,
	})
	rateLimits["user"] = RateLimitConfig{
		RateCounterTicker: userCounter,
		Interval:          1 * time.Second,
		Req:               2, // Set very low for testing
	}
	rateLimitMu.Unlock()

	// Test non-existent role - should be denied for security
	suite.Run("non-existent role", func() {
		allowed := rateLimitedRequest("test-user-1", "non-existent-role")
		suite.False(allowed, "Unknown roles should be denied for security")
	})

	// Test admin role (high limit)
	suite.Run("admin role within limit", func() {
		allowed := rateLimitedRequest("admin-user", "admin")
		suite.True(allowed, "Admin should be within rate limit")
	})

	// Test user role (low limit)
	suite.Run("user role within limit", func() {
		// First request should be allowed
		allowed := rateLimitedRequest("regular-user", "user")
		suite.True(allowed, "First request should be within rate limit")

		// Second request should be allowed
		allowed = rateLimitedRequest("regular-user", "user")
		suite.True(allowed, "Second request should be within rate limit")

		// Third request should exceed limit
		allowed = rateLimitedRequest("regular-user", "user")
		suite.False(allowed, "Third request should exceed rate limit")
	})
}

func (suite *Tests) Test_RateLimitConfig_UnmarshalJSON() {
	// Test unmarshaling of string-based intervals
	suite.Run("unmarshal string intervals", func() {
		// Test JSON with string-based intervals
		jsonString := `{
			"ratelimit": {
				"admin": {
					"req": 100,
					"interval": "second"
				},
				"guest": {
					"req": 5,
					"interval": "minute"
				},
				"user": {
					"req": 1000,
					"interval": "hour"
				},
				"service": {
					"req": 10000,
					"interval": "day"
				},
				"custom": {
					"req": 50,
					"interval": "5s"
				}
			}
		}`

		var config struct {
			RateLimit map[string]RateLimitConfig `json:"ratelimit"`
		}

		err := json.Unmarshal([]byte(jsonString), &config)
		suite.NoError(err)

		// Verify correct parsing of intervals
		suite.Equal(time.Second, config.RateLimit["admin"].Interval)
		suite.Equal(time.Minute, config.RateLimit["guest"].Interval)
		suite.Equal(time.Hour, config.RateLimit["user"].Interval)
		suite.Equal(24*time.Hour, config.RateLimit["service"].Interval)
		suite.Equal(5*time.Second, config.RateLimit["custom"].Interval)

		// Verify req values
		suite.Equal(100, config.RateLimit["admin"].Req)
		suite.Equal(5, config.RateLimit["guest"].Req)
	})

	// Test unmarshaling of invalid interval formats
	suite.Run("unmarshal invalid intervals", func() {
		// Test with an invalid interval format
		jsonString := `{
			"req": 100,
			"interval": "invalid_format"
		}`

		var config RateLimitConfig
		err := json.Unmarshal([]byte(jsonString), &config)
		suite.Error(err)
		suite.Contains(err.Error(), "invalid duration format")
	})

	// Test unmarshaling of numeric intervals
	suite.Run("unmarshal numeric intervals", func() {
		// Test with a numeric interval (seconds)
		jsonString := `{
			"req": 100,
			"interval": 60
		}`

		var config RateLimitConfig
		err := json.Unmarshal([]byte(jsonString), &config)
		suite.NoError(err)
		suite.Equal(60*time.Second, config.Interval)
	})
}

func TestRateLimitConfigUnmarshalKeepsBurstAndEndpoints(t *testing.T) {
	// Regression: the custom UnmarshalJSON temp struct used to carry only
	// interval/req, silently dropping the configured burst and endpoints
	// fields. A config with "burst" or "endpoints" must keep them.
	raw := `{"interval":"minute","req":100,"burst":250,"endpoints":["/v1/graphql"]}`
	var rc RateLimitConfig
	if err := json.Unmarshal([]byte(raw), &rc); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if rc.Interval != time.Minute {
		t.Errorf("interval = %v, want %v", rc.Interval, time.Minute)
	}
	if rc.Req != 100 {
		t.Errorf("req = %d, want 100", rc.Req)
	}
	if rc.Burst != 250 {
		t.Errorf("burst = %d, want 250", rc.Burst)
	}
	if len(rc.Endpoints) != 1 || rc.Endpoints[0] != "/v1/graphql" {
		t.Errorf("endpoints = %v, want [\"/v1/graphql\"]", rc.Endpoints)
	}
}

// fakeClock is a manually-advanced clock injected into a burstWindow so
// tests can move the burst window forward deterministically, without
// sleeping on real time.
type fakeClock struct {
	mu  sync.Mutex
	now time.Time
}

func newFakeClock(start time.Time) *fakeClock {
	return &fakeClock{now: start}
}

func (f *fakeClock) Now() time.Time {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.now
}

func (f *fakeClock) Advance(d time.Duration) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.now = f.now.Add(d)
}

// newBurstTestRoleConfig mirrors the RateLimitConfig construction performed
// by loadConfigFromPath: a real RateCounterTicker for the sustained check,
// and (only when burst > 0) a burstCounter sized by burstWindowDuration.
func newBurstTestRoleConfig(interval time.Duration, req, burst int) RateLimitConfig {
	rc := RateLimitConfig{
		RateCounterTicker: goratecounter.NewRateCounter().WithConfig(goratecounter.RateCounterConfig{
			Interval: interval,
		}),
		Interval: interval,
		Req:      req,
		Burst:    burst,
	}
	if burst > 0 {
		rc.burstCounter = newBurstWindow(burstWindowDuration(interval))
	}
	return rc
}

// Test_checkRateLimit_BurstSemantics is the table test for the C10 burst
// fix. Each case is self-contained and deterministic: sustained-rate
// timing relies only on test execution being far faster than the
// configured Interval (the same assumption Test_rateLimitedRequest already
// makes), and burst-window timing is controlled directly via fakeClock, so
// nothing here depends on a wall-clock sleep as the pass/fail gate.
func Test_checkRateLimit_BurstSemantics(t *testing.T) {
	cfg = &config{}
	parseConfig()
	cfg.Logger = libpack_logger.New()

	t.Run("a_burst_unset_only_req_enforced", func(t *testing.T) {
		rc := newBurstTestRoleConfig(1*time.Second, 2, 0)
		if rc.burstCounter != nil {
			t.Fatalf("burstCounter = %v, want nil when Burst <= 0", rc.burstCounter)
		}

		if !checkRateLimit("u", "role-a", rc, "") {
			t.Errorf("request 1: got denied, want allowed (within Req)")
		}
		if !checkRateLimit("u", "role-a", rc, "") {
			t.Errorf("request 2: got denied, want allowed (at Req)")
		}
		if checkRateLimit("u", "role-a", rc, "") {
			t.Errorf("request 3: got allowed, want denied (exceeds Req, burst disabled)")
		}
	})

	t.Run("b_burst_above_req_allows_short_spike_sustained_still_capped", func(t *testing.T) {
		const (
			interval = 10 * time.Second
			req      = 3
			burst    = 6 // burst > req
		)
		rc := newBurstTestRoleConfig(interval, req, burst)
		fc := newFakeClock(time.Now())
		rc.burstCounter.now = fc.Now

		// A rapid cluster of `burst` requests must all be allowed, even
		// though burst (6) > req (3): the sustained gate only bounds the
		// total over the full 10s Interval (tickerRate <= 3 needs count <=
		// 30), so a small cluster like this never trips it.
		for i := 1; i <= burst; i++ {
			if !checkRateLimit("u", "role-b", rc, "") {
				t.Fatalf("request %d in first burst: got denied, want allowed (<= Burst=%d)", i, burst)
			}
		}

		// The (burst+1)th request in the SAME short window must be denied
		// by the burst gate specifically - not by the sustained gate,
		// which alone would still allow it (tickerRate = 7/10 = 0.7 <= 3).
		// This is the proof burst > req is not a no-op: it is the gate
		// that actually fires here.
		if checkRateLimit("u", "role-b", rc, "") {
			t.Fatalf("request %d: got allowed, want denied by burst gate (> Burst=%d)", burst+1, burst)
		}

		// Roll the burst window forward (fake clock only - the real
		// sustained ticker has not aged, since no real time passed) and
		// repeat clusters of `burst` requests. Each cluster alone passes
		// burst, but the cumulative sustained count keeps climbing until
		// it exceeds req*interval.Seconds() = 30, at which point the
		// sustained gate must deny regardless of how the burst window
		// resets. This proves "sustained still capped at Req".
		sustainedTripped := false
		for batch := 0; batch < 10 && !sustainedTripped; batch++ {
			fc.Advance(2 * burstWindowDuration(interval)) // clears the burst window
			for i := 1; i <= burst; i++ {
				allowed := checkRateLimit("u", "role-b", rc, "")
				rate := rc.RateCounterTicker.GetRate()
				if rate > float64(req) {
					if allowed {
						t.Fatalf("batch %d req %d: allowed with sustained rate %.2f > Req=%d", batch, i, rate, req)
					}
					sustainedTripped = true
					break
				}
			}
		}
		if !sustainedTripped {
			t.Fatalf("sustained gate never tripped after repeated bursts; expected it to cap total requests at Req over the Interval")
		}
	})

	t.Run("c_burst_below_req_documented_behaviour", func(t *testing.T) {
		t.Run("clustered_traffic_capped_below_req_by_burst", func(t *testing.T) {
			const (
				interval = 10 * time.Second
				req      = 6
				burst    = 3 // burst < req
			)
			rc := newBurstTestRoleConfig(interval, req, burst)
			fc := newFakeClock(time.Now())
			rc.burstCounter.now = fc.Now

			for i := 1; i <= burst; i++ {
				if !checkRateLimit("u", "role-c1", rc, "") {
					t.Fatalf("request %d: got denied, want allowed (<= Burst=%d)", i, burst)
				}
			}
			// The 4th rapid request is within Req (tickerRate = 4/10 = 0.4
			// <= 6) but must still be denied: with burst < req, the burst
			// gate is the tighter, binding constraint for clustered
			// traffic - it is not a no-op the way the pre-fix code made
			// burst > req.
			if checkRateLimit("u", "role-c1", rc, "") {
				t.Fatalf("request %d: got allowed, want denied by burst gate (> Burst=%d) though within Req", burst+1, burst)
			}
		})

		t.Run("evenly_spread_traffic_still_governed_by_req_alone", func(t *testing.T) {
			const (
				interval = 1 * time.Second
				req      = 3
				burst    = 1 // burst < req
			)
			rc := newBurstTestRoleConfig(interval, req, burst)
			fc := newFakeClock(time.Now())
			rc.burstCounter.now = fc.Now
			window := burstWindowDuration(interval)

			// One request per (window*2): each lands in its own fresh
			// burst window (count resets to 1 <= Burst=1 every time), so
			// the burst gate never fires. Only the sustained gate, over
			// req=3, eventually denies.
			for i := 1; i <= req; i++ {
				if !checkRateLimit("u", "role-c2", rc, "") {
					t.Fatalf("spread request %d: got denied, want allowed (within Req=%d, burst resets every time)", i, req)
				}
				fc.Advance(2 * window)
			}
			if checkRateLimit("u", "role-c2", rc, "") {
				t.Fatalf("spread request %d: got allowed, want denied by sustained gate (> Req=%d)", req+1, req)
			}
		})
	})

	t.Run("d_fail_safe_never_allows_more_than_burst_in_window", func(t *testing.T) {
		const (
			interval   = 10 * time.Second
			req        = 1000 // high enough it never binds; isolates the burst gate
			burst      = 5
			goroutines = 50
		)
		rc := newBurstTestRoleConfig(interval, req, burst)

		var allowed int64
		var wg sync.WaitGroup
		wg.Add(goroutines)
		for i := 0; i < goroutines; i++ {
			go func() {
				defer wg.Done()
				if checkRateLimit("u", "role-d", rc, "") {
					atomic.AddInt64(&allowed, 1)
				}
			}()
		}
		wg.Wait()

		if allowed != burst {
			t.Fatalf("allowed = %d concurrent requests, want exactly Burst=%d (fail-safe: never more than Burst allowed in the burst window)", allowed, burst)
		}
	})
}
