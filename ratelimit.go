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
	"minute": time.Minute,
	"second": time.Second,
	"hour":   time.Hour,
	"day":    time.Hour * 24,
}

func loadRatelimitConfig() error {
	paths := [2]string{"/app/ratelimit.json", "./ratelimit.json"}
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
			config.RateLimit[key] = value
		}

		rateLimits = config.RateLimit
		return nil
	}
	return os.ErrNotExist
}

func rateLimitedRequest(userRole string, userId string) (shouldAllow bool) {
	if rateLimits == nil {
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
	if rateLimits[userRole].RateCounterTicker.GetRate() > float64(rateLimits[userRole].Req) {
		cfg.Logger.Warning("Rate limit exceeded", map[string]interface{}{"user_role": userRole, "user_id": userId})
		return false
	}
	return
}
