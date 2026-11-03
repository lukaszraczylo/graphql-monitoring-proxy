package main

import (
	"net"
	"testing"
	"time"

	libpack_logger "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
	"github.com/valyala/fasthttp"
)

// Pins theorems/Theorems/ReadyBackoff.lean: the startup wait never sleeps past its timeout.
func TestWaitForBackendReady_HonoursTimeout(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := ln.Addr().String()
	_ = ln.Close() // closed port: every probe is refused immediately

	bhm := NewBackendHealthManager(&fasthttp.Client{}, "http://"+addr, "", libpack_logger.New())
	defer bhm.cancel()

	const timeout = 3 * time.Second
	start := time.Now()
	if err := bhm.WaitForBackendReady(timeout); err == nil {
		t.Fatal("want error for unreachable backend")
	}
	if elapsed := time.Since(start); elapsed > timeout+500*time.Millisecond {
		t.Fatalf("waited %v, want <= %v", elapsed, timeout+500*time.Millisecond)
	}
}
