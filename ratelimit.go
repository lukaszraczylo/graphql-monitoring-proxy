package main

import (
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/goccy/go-json"
	goratecounter "github.com/lukaszraczylo/go-ratecounter"
	libpack_logger "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
)

// RateLimitConfig holds the rate limit configuration for a role
type RateLimitConfig struct {
	RateCounterTicker *goratecounter.RateCounter
	// burstCounter enforces Burst over a short trailing window, independent
	// of RateCounterTicker's sustained-rate window. Nil when Burst <= 0
	// (burst disabled). See burstWindow and checkRateLimit for semantics.
	burstCounter *burstWindow
	Endpoints    []string      `json:"endpoints,omitempty"`
	Interval     time.Duration `json:"interval"`
	Req          int           `json:"req"`
	Burst        int           `json:"burst,omitempty"`
}

// UnmarshalJSON implements custom JSON unmarshaling for RateLimitConfig
func (r *RateLimitConfig) UnmarshalJSON(data []byte) error {
	// Use a temporary struct to unmarshal the JSON data.
	// It must carry every configurable field: Burst (used by the burst
	// allowance in checkRateLimit) and Endpoints would otherwise be
	// silently dropped from the loaded config and left at their zero values.
	type RateLimitConfigTemp struct {
		Interval  any      `json:"interval"`
		Burst     int      `json:"burst,omitempty"`
		Endpoints []string `json:"endpoints,omitempty"`
		Req       int      `json:"req"`
	}

	var temp RateLimitConfigTemp
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	// Set the Req, Burst and Endpoints fields directly
	r.Req = temp.Req
	r.Burst = temp.Burst
	r.Endpoints = temp.Endpoints

	// Handle the Interval field based on its type
	switch v := temp.Interval.(type) {
	case string:
		// Convert string to time.Duration
		switch v {
		case "second":
			r.Interval = time.Second
		case "minute":
			r.Interval = time.Minute
		case "hour":
			r.Interval = time.Hour
		case "day":
			r.Interval = 24 * time.Hour
		default:
			// Try to parse as a Go duration string (e.g. "1s", "5m")
			var err error
			r.Interval, err = time.ParseDuration(v)
			if err != nil {
				return fmt.Errorf("invalid duration format: %s", v)
			}
		}
	case float64:
		// Numeric value is assumed to be in seconds
		r.Interval = time.Duration(v * float64(time.Second))
	default:
		return fmt.Errorf("interval must be a string or number, got %T", v)
	}

	return nil
}

var (
	rateLimits  = make(map[string]RateLimitConfig)
	rateLimitMu sync.RWMutex
	// Use atomic.Value for safe concurrent config swapping
	rateLimitConfigAtomic atomic.Value
)

// burstWindowFraction sets the short burst window as a fraction of the
// sustained Interval: window = Interval / burstWindowFraction. Not
// specified by any upstream spec; chosen so the burst window is clearly
// shorter than the sustained window while remaining simple to reason about.
const burstWindowFraction = 10

// minBurstWindow floors the short burst window so a very small Interval
// (for example a few milliseconds) still leaves the burst window a usable,
// non-degenerate duration.
const minBurstWindow = 10 * time.Millisecond

// burstWindowDuration derives the short burst-window length from the
// sustained Interval. See burstWindowFraction and minBurstWindow.
func burstWindowDuration(interval time.Duration) time.Duration {
	w := interval / burstWindowFraction
	if w < minBurstWindow {
		w = minBurstWindow
	}
	return w
}

// burstWindow is a minimal, clock-injectable, trailing-window request
// counter used to enforce RateLimitConfig.Burst independently of the
// sustained-rate RateCounterTicker.
//
// go-ratecounter@v0.1.12 cannot provide this second, shorter window
// cleanly: the window length lives on *RateCounter itself, not per named
// Counter (helpers.go: GetRate -> c.parent.interval.Seconds()), so two
// different window lengths need two separate RateCounter instances. Each
// instance drives its own background goroutine off wall-clock time.Now with
// no seam to inject a fake clock (helpers.go: start(), cleanUpOldValues(),
// addValue() all call time.Now() directly). That would force any test of
// window rollover to sleep on real time. burstWindow is implemented here
// instead, with a swappable clock, so tests can advance time deterministically.
type burstWindow struct {
	mu     sync.Mutex
	now    func() time.Time
	window time.Duration
	hits   []time.Time
}

// newBurstWindow returns a burstWindow that tracks hits over the trailing
// duration window, using the real wall clock.
func newBurstWindow(window time.Duration) *burstWindow {
	return &burstWindow{now: time.Now, window: window}
}

// recordAndCount records a hit at the current time, drops hits older than
// the trailing window, and returns the resulting hit count within the
// window. Safe for concurrent use.
func (b *burstWindow) recordAndCount() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	now := b.now()
	cutoff := now.Add(-b.window)
	live := b.hits[:0]
	for _, t := range b.hits {
		if t.After(cutoff) {
			live = append(live, t)
		}
	}
	live = append(live, now)
	b.hits = live
	return len(b.hits)
}

// Variable to hold the current load config function - allows for testing
var loadConfigFunc = loadConfigFromPath

// loadRatelimitConfig loads the rate limit configurations from file
func loadRatelimitConfig() error {
	paths := []string{"/go/src/app/ratelimit.json", "./ratelimit.json", "./static/app/default-ratelimit.json"}
	configError := NewRateLimitConfigError(paths)

	// Try each path and collect detailed error information
	for _, path := range paths {
		if err := loadConfigFunc(path); err == nil {
			return nil
		} else {
			// Store the specific error for this path
			configError.PathErrors[path] = err.Error()
		}
	}

	// Log detailed error information
	cfg.Logger.Error(&libpack_logger.LogMessage{
		Message: "Failed to load rate limit configuration",
		Pairs: map[string]any{
			"paths":       paths,
			"path_errors": configError.PathErrors,
		},
	})

	return configError
}

func loadConfigFromPath(path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		// Provide more specific error message based on the error type
		errMsg := ""
		if os.IsNotExist(err) {
			errMsg = "File not found"
		} else if os.IsPermission(err) {
			errMsg = "Permission denied"
		} else {
			errMsg = "I/O error: " + err.Error()
		}

		cfg.Logger.Debug(&libpack_logger.LogMessage{
			Message: "Failed to load rate limit config",
			Pairs: map[string]any{
				"path":          path,
				"error":         errMsg,
				"error_details": err.Error(),
			},
		})
		return fmt.Errorf("%s", errMsg)
	}

	var config struct {
		RateLimit map[string]RateLimitConfig `json:"ratelimit"`
	}

	if err := json.Unmarshal(file, &config); err != nil {
		errMsg := fmt.Sprintf("Invalid JSON format: %s", err.Error())
		cfg.Logger.Debug(&libpack_logger.LogMessage{
			Message: "Failed to parse rate limit config",
			Pairs: map[string]any{
				"path":  path,
				"error": errMsg,
			},
		})
		return fmt.Errorf("%s", errMsg)
	}

	// Validate configuration
	if len(config.RateLimit) == 0 {
		errMsg := "Empty rate limit configuration"
		cfg.Logger.Debug(&libpack_logger.LogMessage{
			Message: "Invalid rate limit config",
			Pairs: map[string]any{
				"path":  path,
				"error": errMsg,
			},
		})
		return fmt.Errorf("%s", errMsg)
	}

	newRateLimits := make(map[string]RateLimitConfig, len(config.RateLimit))
	for key, value := range config.RateLimit {
		value.RateCounterTicker = goratecounter.NewRateCounter().WithConfig(goratecounter.RateCounterConfig{
			Interval: value.Interval,
		})

		// Burst disabled (<=0) leaves burstCounter nil, so checkRateLimit
		// performs the Req-only check exactly as before burst existed.
		if value.Burst > 0 {
			value.burstCounter = newBurstWindow(burstWindowDuration(value.Interval))
		}

		if cfg.LogLevel == "DEBUG" {
			cfg.Logger.Debug(&libpack_logger.LogMessage{
				Message: "Setting ratelimit config for role",
				Pairs: map[string]any{
					"role":          key,
					"interval_used": value.Interval,
					"ratelimit":     value.Req,
				},
			})
		}
		newRateLimits[key] = value
	}

	// Use atomic swap for thread-safe configuration updates
	rateLimitMu.Lock()
	rateLimits = newRateLimits
	// Store the new config atomically
	rateLimitConfigAtomic.Store(newRateLimits)
	rateLimitMu.Unlock()

	cfg.Logger.Debug(&libpack_logger.LogMessage{
		Message: "Rate limit config loaded",
		Pairs:   map[string]any{"ratelimit": rateLimits},
	})
	return nil
}

// rateLimitedRequest checks if a request should be rate-limited
func rateLimitedRequest(userID, userRole string) bool {
	// Try to get config from atomic value first for better performance
	if configInterface := rateLimitConfigAtomic.Load(); configInterface != nil {
		if config, ok := configInterface.(map[string]RateLimitConfig); ok {
			if roleConfig, exists := config[userRole]; exists && roleConfig.RateCounterTicker != nil {
				return checkRateLimit(userID, userRole, roleConfig, "")
			}
		}
	}

	// Fallback to mutex-protected access
	rateLimitMu.RLock()
	roleConfig, ok := rateLimits[userRole]
	rateLimitMu.RUnlock()

	if !ok || roleConfig.RateCounterTicker == nil {
		cfg.Logger.Warning(&libpack_logger.LogMessage{
			Message: "Rate limit role not found or ticker not initialized - defaulting to deny",
			Pairs:   map[string]any{"user_role": userRole},
		})
		// Default to deny when config not found (security fix)
		return false
	}

	return checkRateLimit(userID, userRole, roleConfig, "")
}

// checkRateLimit performs the actual rate limit check.
//
// Two independent gates apply, both must pass:
//
//  1. Sustained: tickerRate (requests measured over the full Interval) must
//     stay at or below Req. This is the only gate when Burst <= 0, and its
//     comparison is unchanged from before the burst fix.
//  2. Burst (only when roleConfig.burstCounter != nil, i.e. Burst > 0): the
//     raw request count within a short trailing window (see
//     burstWindowDuration), much shorter than Interval, must stay at or
//     below Burst.
//
// Because the burst window is short relative to Interval, a client can
// legally land up to Burst requests in a tight cluster - more than a
// steady Req-paced stream would deliver in that same short span - as long
// as the sustained gate still holds over the full Interval. That is what
// makes burst > req a real spike allowance instead of a no-op. Burst is a
// short-window ceiling layered ON TOP of the sustained Req rate, not a
// replacement for it. burst < req is an unusual config: the burst gate
// becomes the tighter, binding constraint whenever requests cluster in
// time, but the sustained cap for evenly-spread traffic is still Req
// alone. The pre-fix code compared the SAME sustained-window tickerRate
// against both Burst and Req, so when burst < req its Burst check tripped
// first and Burst was effectively the sustained cap too. This gate's
// Req-only sustained check can therefore allow MORE through than the
// pre-fix code did for evenly-spread traffic under a burst < req config.
// In both cases the burst gate never allows more than Burst requests
// within its own window (fail-safe).
func checkRateLimit(userID, userRole string, roleConfig RateLimitConfig, endpoint string) bool {
	roleConfig.RateCounterTicker.Incr(1)
	tickerRate := roleConfig.RateCounterTicker.GetRate()

	logDetails := map[string]any{
		"user_role":   userRole,
		"user_id":     userID,
		"rate":        tickerRate,
		"config_rate": roleConfig.Req,
		"interval":    roleConfig.Interval,
		"endpoint":    endpoint,
	}

	cfg.Logger.Debug(&libpack_logger.LogMessage{
		Message: "Rate limit ticker",
		Pairs:   map[string]any{"log_details": logDetails},
	})

	if tickerRate > float64(roleConfig.Req) {
		cfg.Logger.Debug(&libpack_logger.LogMessage{
			Message: "Rate limit exceeded",
			Pairs:   map[string]any{"log_details": logDetails},
		})
		return false
	}

	if roleConfig.burstCounter != nil {
		burstCount := roleConfig.burstCounter.recordAndCount()
		burstLogDetails := map[string]any{
			"user_role":    userRole,
			"user_id":      userID,
			"burst_count":  burstCount,
			"config_burst": roleConfig.Burst,
			"endpoint":     endpoint,
		}
		if burstCount > roleConfig.Burst {
			cfg.Logger.Debug(&libpack_logger.LogMessage{
				Message: "Burst limit exceeded",
				Pairs:   map[string]any{"log_details": burstLogDetails},
			})
			return false
		}
	}

	return true
}
