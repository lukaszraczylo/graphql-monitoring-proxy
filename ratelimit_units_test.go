package main

import (
	"testing"
	"time"

	libpack_logger "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
)

func setupRateLimitCfg() {
	cfg = &config{}
	parseConfig()
	cfg.Logger = libpack_logger.New()
}

// Pins theorems/Theorems/RateLimitUnits.lean: Req counts requests per Interval, not per second.
func TestCheckRateLimit_ReqIsPerInterval(t *testing.T) {
	const req = 3
	setupRateLimitCfg()
	rc := newBurstTestRoleConfig(time.Minute, req, 0)
	for i := 1; i <= req; i++ {
		if !checkRateLimit("u", "r", rc, "") {
			t.Fatalf("request %d denied, want allowed (Req=%d per minute)", i, req)
		}
	}
	if checkRateLimit("u", "r", rc, "") {
		t.Fatalf("request %d allowed, want denied: Req=%d per minute", req+1, req)
	}
}

// With Burst >= Req the burst gate never rejects what the sustained gate admits.
func TestCheckRateLimit_BurstAboveReqIsNoOp(t *testing.T) {
	const req, burst = 3, 6
	setupRateLimitCfg()
	rc := newBurstTestRoleConfig(time.Minute, req, burst)
	for i := 1; i <= req; i++ {
		if !checkRateLimit("u", "r", rc, "") {
			t.Fatalf("request %d denied", i)
		}
	}
	if checkRateLimit("u", "r", rc, "") {
		t.Fatalf("request %d allowed, want denied by the sustained gate", req+1)
	}
}
