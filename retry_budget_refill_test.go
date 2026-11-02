package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Pins theorems/Theorems/RetryRefill.lean: total_added(15 tps, 100ms ticks, 10s) = 150.
func TestRefillStep_KeepsFractionalTokens(t *testing.T) {
	tests := []struct {
		name string
		tps  float64
		want int64
	}{
		{"15 tps", 15, 150},
		{"25 tps", 25, 250},
		{"5 tps (sub-token per tick)", 5, 50},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tick := 100 * time.Millisecond
			var last, total int64
			now := int64(tick)
			for i := 0; i < 100; i++ {
				added, newLast := refillStep(tt.tps, last, now)
				total += added
				last = newLast
				now += int64(tick)
			}
			assert.Equal(t, tt.want, total)
		})
	}
}

func TestRefillStep_NoElapsedNoChange(t *testing.T) {
	added, last := refillStep(10, 1000, 1000)
	assert.Zero(t, added)
	assert.Equal(t, int64(1000), last)
}

func TestRefillStep_ZeroRateNeverAdds(t *testing.T) {
	added, last := refillStep(0, 0, int64(time.Hour))
	assert.Zero(t, added)
	assert.Zero(t, last)
}
