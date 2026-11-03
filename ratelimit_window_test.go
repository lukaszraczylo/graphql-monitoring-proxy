package main

import (
	"math/rand"
	"testing"
	"time"
)

// naiveExceeds is the exact reference: count every attempt in the trailing window.
func naiveExceeds(attempts []time.Time, now time.Time, window time.Duration, limit int) bool {
	cutoff := now.Add(-window)
	n := 0
	for _, t := range attempts {
		if t.After(cutoff) {
			n++
		}
	}
	return n > limit
}

// Pins theorems/Theorems/SlidingWindow.lean: storing the newest limit+1 attempts decides like the full count.
func TestAttemptWindow_MatchesExactCount(t *testing.T) {
	for _, limit := range []int{0, 1, 3, 10} {
		rng := rand.New(rand.NewSource(int64(limit) + 1))
		window := 10 * time.Second
		fc := newFakeClock(time.Unix(1_000_000, 0))
		w := newAttemptWindow(window, limit)
		w.now = fc.Now
		w.base = fc.Now()

		var attempts []time.Time
		for i := 0; i < 5000; i++ {
			fc.Advance(time.Duration(rng.Intn(3000)) * time.Millisecond)
			attempts = append(attempts, fc.Now())
			want := naiveExceeds(attempts, fc.Now(), window, limit)
			if got := w.recordAndExceeds(); got != want {
				t.Fatalf("limit=%d step=%d: got %v, want %v", limit, i, got, want)
			}
		}
	}
}

func TestAttemptWindow_StorageIsBounded(t *testing.T) {
	const limit = 10
	fc := newFakeClock(time.Unix(1_000_000, 0))
	w := newAttemptWindow(time.Hour, limit)
	w.now = fc.Now
	w.base = fc.Now()
	for i := 0; i < 100_000; i++ {
		w.recordAndExceeds()
		fc.Advance(time.Millisecond)
	}
	if got := len(w.ring); got > limit+1 {
		t.Fatalf("stored %d timestamps, want <= %d", got, limit+1)
	}
}

func BenchmarkAttemptWindow(b *testing.B) {
	w := newAttemptWindow(time.Hour, 1000)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		w.recordAndExceeds()
	}
}
