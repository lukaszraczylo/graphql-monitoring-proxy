package main

import (
	"sync"
	"time"
)

// attemptWindow decides "more than limit attempts in the trailing window" while
// storing only the newest limit+1 attempt times. The attempt number limit+1
// counted back from now lies inside the window exactly when the full count
// exceeds limit (see theorems/Theorems/SlidingWindow.lean), so memory and time
// per attempt stay constant however hard a client floods.
type attemptWindow struct {
	now    func() time.Time
	base   time.Time // monotonic origin for stored offsets
	ring   []int64   // offsets from base, grows to capacity then wraps
	mu     sync.Mutex
	window time.Duration
	head   int // index of the oldest stored attempt once the ring is full
	size   int // limit+1: attempts kept; the ring grows lazily up to this
}

func newAttemptWindow(window time.Duration, limit int) *attemptWindow {
	return &attemptWindow{
		now:    time.Now,
		base:   time.Now(),
		window: window,
		size:   max(limit, 0) + 1,
	}
}

// recordAndExceeds records an attempt now and reports whether the attempts in
// the trailing window, this one included, outnumber the limit.
func (w *attemptWindow) recordAndExceeds() bool {
	w.mu.Lock()
	defer w.mu.Unlock()

	t := int64(w.now().Sub(w.base))
	if len(w.ring) < w.size {
		w.ring = append(w.ring, t)
	} else {
		w.ring[w.head] = t
		w.head = (w.head + 1) % len(w.ring)
	}
	if len(w.ring) < w.size {
		return false
	}
	return w.ring[w.head] > t-int64(w.window)
}
