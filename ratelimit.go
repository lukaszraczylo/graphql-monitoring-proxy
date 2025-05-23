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
	Interval          time.Duration `json:"interval"`
	Req               int           `json:"req"`
}

// UnmarshalJSON implements custom JSON unmarshaling for RateLimitConfig
func (r *RateLimitConfig) UnmarshalJSON(data []byte) error {
	// Use a temporary struct to unmarshal the JSON data
	type RateLimitConfigTemp struct {
		Interval interface{} `json:"interval"`
		Req      int         `json:"req"`
	}

	var temp RateLimitConfigTemp
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	// Set the Req field directly
	r.Req = temp.Req

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
		Pairs: map[string]interface{}{
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
			Pairs: map[string]interface{}{
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
			Pairs: map[string]interface{}{
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
			Pairs: map[string]interface{}{
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

		if cfg.LogLevel == "DEBUG" {
			cfg.Logger.Debug(&libpack_logger.LogMessage{
				Message: "Setting ratelimit config for role",
				Pairs: map[string]interface{}{
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
		Pairs:   map[string]interface{}{"ratelimit": rateLimits},
	})
	return nil
}

// rateLimitedRequest checks if a request should be rate-limited
func rateLimitedRequest(userID, userRole string) bool {
	// Try to get config from atomic value first for better performance
	if configInterface := rateLimitConfigAtomic.Load(); configInterface != nil {
		if config, ok := configInterface.(map[string]RateLimitConfig); ok {
			if roleConfig, exists := config[userRole]; exists && roleConfig.RateCounterTicker != nil {
				return checkRateLimit(userID, userRole, roleConfig)
			}
		}
	}

	// Fallback to mutex-protected access
	rateLimitMu.RLock()
	roleConfig, ok := rateLimits[userRole]
	rateLimitMu.RUnlock()

	if !ok || roleConfig.RateCounterTicker == nil {
		cfg.Logger.Debug(&libpack_logger.LogMessage{
			Message: "Rate limit role not found or ticker not initialized - defaulting to deny",
			Pairs:   map[string]interface{}{"user_role": userRole},
		})
		// Default to deny when config not found (security fix)
		return false
	}

	return checkRateLimit(userID, userRole, roleConfig)
}

// checkRateLimit performs the actual rate limit check
func checkRateLimit(userID, userRole string, roleConfig RateLimitConfig) bool {
	roleConfig.RateCounterTicker.Incr(1)
	tickerRate := roleConfig.RateCounterTicker.GetRate()

	logDetails := map[string]interface{}{
		"user_role":   userRole,
		"user_id":     userID,
		"rate":        tickerRate,
		"config_rate": roleConfig.Req,
		"interval":    roleConfig.Interval,
	}

	cfg.Logger.Debug(&libpack_logger.LogMessage{
		Message: "Rate limit ticker",
		Pairs:   map[string]interface{}{"log_details": logDetails},
	})

	if tickerRate > float64(roleConfig.Req) {
		cfg.Logger.Debug(&libpack_logger.LogMessage{
			Message: "Rate limit exceeded",
			Pairs:   map[string]interface{}{"log_details": logDetails},
		})
		return false
	}

	return true
}
