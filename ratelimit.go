package main

import (
	"os"
	"time"

	goratecounter "github.com/lukaszraczylo/go-ratecounter"
)

type RateLimitConfig struct {
	Req               int    `json:"req"`
	Interval          string `json:"interval"`
	RateCounterTicker *goratecounter.RateCounter
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
	paths := [3]string{"/app/ratelimit.json", "./ratelimit.json", "./static/default-ratelimit.json"}
	for _, path := range paths {
		file, err := os.Open(path)
		if err != nil {
			continue
		}
		defer file.Close()
		decoder := json.NewDecoder(file)
		config := struct {
			RateLimit map[string]RateLimitConfig `json:"ratelimit"`
		}{}
		err = decoder.Decode(&config)
		if err != nil {
			return err
		}

		for key, value := range config.RateLimit {
			value.RateCounterTicker = goratecounter.NewRateCounter().WithConfig(goratecounter.RateCounterConfig{
				Interval: time.Duration(value.Req) * ratelimit_intervals[value.Interval],
			})
			cfg.Logger.Debug("Setting ratelimit config for role", map[string]interface{}{"role": key, "interval_provided": value.Interval, "interval_used": ratelimit_intervals[value.Interval], "ratelimit": value.Req})
			config.RateLimit[key] = value
		}

		rateLimits = config.RateLimit
		cfg.Logger.Debug("Rate limit config loaded", map[string]interface{}{"ratelimit": rateLimits})
		return nil
	}
	cfg.Logger.Debug("Rate limit config not found")
	return os.ErrNotExist
}

func rateLimitedRequest(userId string, userRole string) (shouldAllow bool) {
	if rateLimits == nil {
		cfg.Logger.Debug("Rate limit config not found", map[string]interface{}{"user_role": userRole})
		return true
	}
	// check if userRole is in rateLimits
	if _, ok := rateLimits[userRole]; !ok {
		cfg.Logger.Warning("Rate limit role not found", map[string]interface{}{"user_role": userRole})
		return true
	}

	if rateLimits[userRole].RateCounterTicker == nil {
		cfg.Logger.Warning("Rate limit ticker not found", map[string]interface{}{"user_role": userRole})
		return true
	}

	rateLimits[userRole].RateCounterTicker.Incr(1)
	ticker_rate := rateLimits[userRole].RateCounterTicker.GetRate()
	cfg.Logger.Debug("Rate limit ticker", map[string]interface{}{"user_role": userRole, "user_id": userId, "rate": ticker_rate, "config_rate": rateLimits[userRole].Req, "interval": rateLimits[userRole].Interval, "interval_duration": rateLimits[userRole].Interval})
	if ticker_rate > float64(rateLimits[userRole].Req) {
		cfg.Logger.Debug("Rate limit exceeded", map[string]interface{}{"user_role": userRole, "user_id": userId, "rate": ticker_rate, "config_rate": rateLimits[userRole].Req, "interval": rateLimits[userRole].Interval, "interval_duration": rateLimits[userRole].Interval})
		return false
	}
	return true
}
