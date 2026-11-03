package main

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Pins theorems/Theorems/RpsSample.lean: no recorded request is dropped by sample().
func TestRPSTracker_SampleKeepsConcurrentRequests(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	tr := NewRPSTracker(ctx)
	tr.Shutdown() // stop the background sampler; sample() is driven by the test

	const writers, perWriter = 8, 20000
	var wg sync.WaitGroup
	for i := 0; i < writers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < perWriter; j++ {
				tr.RecordRequest()
			}
		}()
	}
	var sampled int64
	done := make(chan struct{})
	go func() { wg.Wait(); close(done) }()
loop:
	for {
		select {
		case <-done:
			break loop
		default:
			sampled += tr.sampleCount()
		}
	}
	sampled += tr.sampleCount()
	assert.Equal(t, int64(writers*perWriter), sampled)
}
