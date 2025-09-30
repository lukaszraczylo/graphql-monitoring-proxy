package main

import (
	"sync"
	"sync/atomic"
	"time"
)

// RPSTracker tracks requests per second using periodic sampling
type RPSTracker struct {
	lastCount      atomic.Int64
	lastSampleTime atomic.Int64 // Unix nano
	currentRPS     uint64       // stored as uint64, accessed with atomic operations
	mu             sync.RWMutex // for currentRPS updates
}

// NewRPSTracker creates a new RPS tracker
func NewRPSTracker() *RPSTracker {
	tracker := &RPSTracker{}
	tracker.lastSampleTime.Store(time.Now().UnixNano())
	go tracker.updateLoop()
	return tracker
}

// RecordRequest increments the request counter
func (r *RPSTracker) RecordRequest() {
	// Just increment the counter, sampling happens in background
	r.lastCount.Add(1)
}

// updateLoop periodically calculates current RPS
func (r *RPSTracker) updateLoop() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		r.sample()
	}
}

// sample calculates RPS since last sample
func (r *RPSTracker) sample() {
	now := time.Now()
	nowNano := now.UnixNano()

	currentCount := r.lastCount.Load()
	lastSampleNano := r.lastSampleTime.Load()

	if lastSampleNano == 0 {
		r.lastSampleTime.Store(nowNano)
		return
	}

	elapsed := float64(nowNano-lastSampleNano) / float64(time.Second)
	if elapsed > 0 {
		rps := float64(currentCount) / elapsed
		// Store RPS as centirps for precision (multiply by 100)
		r.mu.Lock()
		atomic.StoreUint64(&r.currentRPS, uint64(rps*100))
		r.mu.Unlock()
	}

	// Reset for next sample
	r.lastCount.Store(0)
	r.lastSampleTime.Store(nowNano)
}

// GetCurrentRPS returns the current requests per second
func (r *RPSTracker) GetCurrentRPS() float64 {
	r.mu.RLock()
	centirps := atomic.LoadUint64(&r.currentRPS)
	r.mu.RUnlock()
	return float64(centirps) / 100.0
}

var globalRPSTracker *RPSTracker

// InitializeRPSTracker initializes the global RPS tracker
func InitializeRPSTracker() *RPSTracker {
	if globalRPSTracker == nil {
		globalRPSTracker = NewRPSTracker()
	}
	return globalRPSTracker
}

// GetRPSTracker returns the global RPS tracker
func GetRPSTracker() *RPSTracker {
	return globalRPSTracker
}
