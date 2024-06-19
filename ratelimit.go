package main

import (
	"os"
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

var rateLimits map[string]RateLimitConfig
var ratelimit_intervals = map[string]time.Duration{
	"milli":  time.Millisecond,
	"micro":  time.Microsecond,
	"nano":   time.Nanosecond,
	"second": time.Second,
	"minute": time.Minute,
	"hour":   time.Hour,
	"day":    time.Hour * 24,
}

func loadRatelimitConfig() error {
	paths := []string{"/go/src/app/ratelimit.json", "./ratelimit.json", "./static/app/default-ratelimit.json"}

	for _, path := range paths {
		err := loadConfigFromPath(path)
		if err == nil {
			return nil
		}
		cfg.Logger.Debug(&libpack_logger.LogMessage{
			Message: "Failed to load config",
			Pairs:   map[string]interface{}{"path": path, "error": err},
		})
	}

	cfg.Logger.Error(&libpack_logger.LogMessage{
		Message: "Rate limit config not found",
		Pairs:   map[string]interface{}{"paths": paths},
	})

	return os.ErrNotExist
}

func loadConfigFromPath(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	config := struct {
		RateLimit map[string]RateLimitConfig `json:"ratelimit"`
	}{}

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return err
	}

	for key, value := range config.RateLimit {
		value.RateCounterTicker = goratecounter.NewRateCounter().WithConfig(goratecounter.RateCounterConfig{
			Interval: time.Duration(value.Req) * ratelimit_intervals[value.Interval],
		})

		if cfg.LogLevel == "debug" {
			cfg.Logger.Debug(&libpack_logger.LogMessage{
				Message: "Setting ratelimit config for role",
				Pairs: map[string]interface{}{
					"role":              key,
					"interval_provided": value.Interval,
					"interval_used":     ratelimit_intervals[value.Interval],
					"ratelimit":         value.Req,
				},
			})
		}
		config.RateLimit[key] = value
	}

	rateLimits = config.RateLimit
	cfg.Logger.Debug(&libpack_logger.LogMessage{
		Message: "Rate limit config loaded",
		Pairs:   map[string]interface{}{"ratelimit": rateLimits},
	})
	return nil
}

func rateLimitedRequest(userID string, userRole string) (shouldAllow bool) {
	if rateLimits == nil {
		cfg.Logger.Debug(&libpack_logger.LogMessage{
			Message: "Rate limit config not found",
			Pairs:   map[string]interface{}{"user_role": userRole},
		})
		return true
	}

	// Fetch role config once to avoid multiple map lookups
	roleConfig, ok := rateLimits[userRole]
	if !ok {
		cfg.Logger.Debug(&libpack_logger.LogMessage{
			Message: "Rate limit role not found",
			Pairs:   map[string]interface{}{"user_role": userRole},
		})
		return true
	}

	if roleConfig.RateCounterTicker == nil {
		cfg.Logger.Debug(&libpack_logger.LogMessage{
			Message: "Rate limit ticker not found",
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
