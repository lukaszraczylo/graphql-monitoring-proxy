package main

import (
	"os"
	"sync"
	"time"

	"github.com/goccy/go-json"
	goratecounter "github.com/lukaszraczylo/go-ratecounter"
	libpack_logger "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
)

type RateLimitConfig struct {
	RateCounterTicker *goratecounter.RateCounter
	Interval          string `json:"interval"`
	Req               int    `json:"req"`
}

var (
	rateLimits         map[string]RateLimitConfig
	ratelimitIntervals = map[string]time.Duration{
		"milli":  time.Millisecond,
		"micro":  time.Microsecond,
		"nano":   time.Nanosecond,
		"second": time.Second,
		"minute": time.Minute,
		"hour":   time.Hour,
		"day":    24 * time.Hour,
	}
	rateLimitMu sync.RWMutex
)

func loadRatelimitConfig() error {
	paths := []string{"/go/src/app/ratelimit.json", "./ratelimit.json", "./static/app/default-ratelimit.json"}
	for _, path := range paths {
		if err := loadConfigFromPath(path); err == nil {
			return nil
		}
	}
	cfg.Logger.Error(&libpack_logger.LogMessage{
		Message: "Rate limit config not found",
		Pairs:   map[string]interface{}{"paths": paths},
	})
	return os.ErrNotExist
}

func loadConfigFromPath(path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		cfg.Logger.Debug(&libpack_logger.LogMessage{
			Message: "Failed to load config",
			Pairs:   map[string]interface{}{"path": path, "error": err},
		})
		return err
	}

	var config struct {
		RateLimit map[string]RateLimitConfig `json:"ratelimit"`
	}

	if err := json.Unmarshal(file, &config); err != nil {
		return err
	}

	newRateLimits := make(map[string]RateLimitConfig, len(config.RateLimit))
	for key, value := range config.RateLimit {
		interval := ratelimitIntervals[value.Interval]
		value.RateCounterTicker = goratecounter.NewRateCounter().WithConfig(goratecounter.RateCounterConfig{
			Interval: time.Duration(value.Req) * interval,
		})

		if cfg.LogLevel == "debug" {
			cfg.Logger.Debug(&libpack_logger.LogMessage{
				Message: "Setting ratelimit config for role",
				Pairs: map[string]interface{}{
					"role":              key,
					"interval_provided": value.Interval,
					"interval_used":     interval,
					"ratelimit":         value.Req,
				},
			})
		}
		newRateLimits[key] = value
	}

	rateLimitMu.Lock()
	rateLimits = newRateLimits
	rateLimitMu.Unlock()

	cfg.Logger.Debug(&libpack_logger.LogMessage{
		Message: "Rate limit config loaded",
		Pairs:   map[string]interface{}{"ratelimit": rateLimits},
	})
	return nil
}

func rateLimitedRequest(userID, userRole string) bool {
	rateLimitMu.RLock()
	defer rateLimitMu.RUnlock()

	if rateLimits == nil {
		cfg.Logger.Debug(&libpack_logger.LogMessage{
			Message: "Rate limit config not found",
			Pairs:   map[string]interface{}{"user_role": userRole},
		})
		return true
	}

	roleConfig, ok := rateLimits[userRole]
	if !ok || roleConfig.RateCounterTicker == nil {
		cfg.Logger.Debug(&libpack_logger.LogMessage{
			Message: "Rate limit role not found or ticker not initialized",
			Pairs:   map[string]interface{}{"user_role": userRole},
		})
		return true
	}

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
