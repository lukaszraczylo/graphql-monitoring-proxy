package libpack_cache_memory

import (
	"fmt"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Default constants for testing
const (
	DefaultTestExpiration = 5 * time.Second
)

func TestMemoryCacheClear(t *testing.T) {
	cache := New(DefaultTestExpiration)

	// Add some entries
	cache.Set("key1", []byte("value1"), DefaultTestExpiration)
	cache.Set("key2", []byte("value2"), DefaultTestExpiration)

	// Verify entries exist
	_, found := cache.Get("key1")
	assert.True(t, found, "Expected key1 to exist before clearing cache")

	// Clear the cache
	cache.Clear()

	// Verify cache is empty
	_, found = cache.Get("key1")
	assert.False(t, found, "Expected key1 to be removed after clearing cache")
	_, found = cache.Get("key2")
	assert.False(t, found, "Expected key2 to be removed after clearing cache")

	// Check that counter was reset
	assert.Equal(t, int64(0), cache.CountQueries(), "Expected count to be 0 after clearing cache")
}

func TestMemoryCacheCountQueries(t *testing.T) {
	cache := New(DefaultTestExpiration)

	// Check initial count
	assert.Equal(t, int64(0), cache.CountQueries(), "Expected initial count to be 0")

	// Add some entries
	cache.Set("key1", []byte("value1"), DefaultTestExpiration)
	cache.Set("key2", []byte("value2"), DefaultTestExpiration)
	cache.Set("key3", []byte("value3"), DefaultTestExpiration)

	// Check count
	assert.Equal(t, int64(3), cache.CountQueries(), "Expected count to be 3 after adding 3 entries")

	// Delete an entry
	cache.Delete("key1")

	// Check count after deletion
	assert.Equal(t, int64(2), cache.CountQueries(), "Expected count to be 2 after deleting 1 entry")
}

func TestMemoryCacheCleanExpiredEntries(t *testing.T) {
	// Create a cache with default expiration
	cache := New(10 * time.Second)

	// Add an entry that will expire quickly
	cache.Set("expire-soon", []byte("value1"), 10*time.Millisecond)

	// Add an entry that will not expire during the test
	cache.Set("expire-later", []byte("value3"), 10*time.Minute)

	// Initial count should be 2
	assert.Equal(t, int64(2), cache.CountQueries(), "Expected count to be 2 after adding entries")

	// Wait for short expiration
	time.Sleep(20 * time.Millisecond)

	// Get the expired key directly to verify it's expired
	_, expiredFound := cache.Get("expire-soon")
	assert.False(t, expiredFound, "Key 'expire-soon' should be expired now")

	// Verify the not-expired key is still there
	val, nonExpiredFound := cache.Get("expire-later")
	assert.True(t, nonExpiredFound, "Key 'expire-later' should not be expired")
	assert.Equal(t, []byte("value3"), val, "Expected correct value for 'expire-later'")

	// Manually clean expired entries
	cache.CleanExpiredEntries()

	// Count should be 1 now (only the non-expired entry)
	assert.Equal(t, int64(1), cache.CountQueries(), "Expected count to be 1 after cleaning expired entries")
}

func TestMemoryCacheConcurrentGetOnExpiredEntry_CountersConsistent(t *testing.T) {
	// Regression: concurrent Gets on the same expired entry must decrement
	// entryCount exactly once. Previously Get used an unconditional Delete +
	// decrement, so N concurrent Gets drove entryCount negative.
	cache := NewWithSize(DefaultTestExpiration, DefaultMaxMemorySize, DefaultMaxCacheSize)

	cache.Set("k", []byte("v"), 10*time.Millisecond)
	time.Sleep(50 * time.Millisecond) // let it expire

	var wg sync.WaitGroup
	for range 32 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, ok := cache.Get("k"); ok {
				t.Error("expired entry should not be returned")
			}
		}()
	}
	wg.Wait()

	if got := cache.CountQueries(); got != 0 {
		t.Fatalf("entryCount = %d after concurrent expired gets, want 0", got)
	}
	if got := cache.GetMemoryUsage(); got != 0 {
		t.Fatalf("memoryUsage = %d after concurrent expired gets, want 0", got)
	}
}

func TestMemoryCacheConcurrentSetSameNewKey_CounterNotDoubled(t *testing.T) {
	// Regression: two (or more) concurrent Sets of the same brand-new key must
	// not both take the "new entry" branch. Previously the exists-check and
	// counter increment were separate atomic ops, so concurrent Sets of a new
	// key each incremented entryCount, inflating it for a single stored entry.
	cache := NewWithSize(DefaultTestExpiration, DefaultMaxMemorySize, DefaultMaxCacheSize)

	const workers = 32
	var wg sync.WaitGroup
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cache.Set("same-key", []byte("value"), DefaultTestExpiration)
		}()
	}
	wg.Wait()

	if got := cache.CountQueries(); got != 1 {
		t.Fatalf("entryCount = %d after %d concurrent Sets of one new key, want 1", got, workers)
	}

	// Sanity: the entry is actually stored and retrievable.
	if _, ok := cache.Get("same-key"); !ok {
		t.Fatal("key should be retrievable after Set")
	}
}

func TestMemoryCacheConcurrentSetDeleteEvict_CountersMatchActualState(t *testing.T) {
	// Regression for M4 (evictToFreeMemory pre-allocation sized from live
	// entryCount, not configured maxCacheSize) and S2 (lock-free Set via
	// sync.Map.Swap, so Set no longer double-accounts when it races with
	// Delete/evictOldest/evictToFreeMemory/CleanExpiredEntries on the same
	// key).
	//
	// The specific key each goroutine touches, and which operation wins a
	// given race, is inherently nondeterministic. What must always hold,
	// deterministically, once every goroutine has joined, is the invariant
	// that entryCount/memoryUsage are derived state: they must exactly match
	// the cache's real final contents. Recomputing "expected" from the
	// actual final map (the known final key set) and diffing against the
	// atomic counters catches any lost or double-counted update.
	const (
		numKeys    = 50
		numWorkers = 16
		numOpsEach = 300
	)

	// Small maxCacheSize so Set's own eviction path (evictOldest) also fires
	// under the concurrent load, in addition to the evictors driven directly
	// below.
	cache := NewWithSize(DefaultTestExpiration, DefaultMaxMemorySize, 40)
	defer cache.Shutdown()

	keyFor := func(i int) string { return fmt.Sprintf("hammer-key-%d", i) }

	// Pre-seed so evictors have entries to act on from the start.
	for i := range numKeys {
		cache.Set(keyFor(i), []byte("seed-value"), DefaultTestExpiration)
	}

	var wg sync.WaitGroup
	for w := range numWorkers {
		wg.Add(1)
		go func(worker int) {
			defer wg.Done()
			for op := range numOpsEach {
				k := keyFor((worker + op) % numKeys)
				switch op % 3 {
				case 0:
					cache.Set(k, []byte("value"), DefaultTestExpiration)
				case 1:
					cache.Delete(k)
				case 2:
					cache.evictOldest(2)
				}
			}
		}(w)
	}

	wg.Add(2)
	go func() {
		defer wg.Done()
		for range 100 {
			cache.evictToFreeMemory(128)
		}
	}()
	go func() {
		defer wg.Done()
		for range 100 {
			cache.CleanExpiredEntries()
		}
	}()

	wg.Wait()

	// Stop the background cleanup goroutine before recomputing expected
	// state, so it cannot race with the Range below (nothing has expired
	// yet - DefaultTestExpiration is 5s - but this removes any doubt).
	cache.Shutdown()

	var wantCount, wantMemory int64
	cache.entries.Range(func(_, v any) bool {
		entry := v.(CacheEntry)
		wantCount++
		wantMemory += entry.MemorySize
		return true
	})

	if got := cache.CountQueries(); got != wantCount {
		t.Fatalf("entryCount = %d, want %d (actual live entries after join)", got, wantCount)
	}
	if got := cache.GetMemoryUsage(); got != wantMemory {
		t.Fatalf("memoryUsage = %d, want %d (sum of live entries' MemorySize after join)", got, wantMemory)
	}
}

func TestMemoryCacheConcurrentSetDeleteClear_CountersMatchActualState(t *testing.T) {
	// Regression for Clear(): the previous implementation Ranged the
	// sync.Map calling Delete per key, then unconditionally reset both
	// counters with StoreInt64(0). A Set racing Clear could land after the
	// Range pass but before or after the StoreInt64 calls, leaving a live
	// entry in the map while the counters read zero - the next Delete of
	// that key then drove the counters permanently negative.
	//
	// Clear now follows the same linearized LoadAndDelete-then-subtract
	// pattern as Delete/Get/CleanExpiredEntries/the evictors, so it no
	// longer promises a consistent point-in-time snapshot (a concurrent Set
	// can survive a Clear) but does guarantee entryCount/memoryUsage always
	// match the entries actually removed. As in
	// TestMemoryCacheConcurrentSetDeleteEvict_CountersMatchActualState,
	// which key wins a given race is nondeterministic, so the assertion
	// recomputes "expected" from the final map contents rather than
	// asserting a fixed number.
	const (
		numKeys    = 50
		numWorkers = 16
		numOpsEach = 300
	)

	cache := NewWithSize(DefaultTestExpiration, DefaultMaxMemorySize, DefaultMaxCacheSize)
	defer cache.Shutdown()

	keyFor := func(i int) string { return fmt.Sprintf("clear-key-%d", i) }

	// Pre-seed so Clear and the Delete workers have entries to act on from
	// the start.
	for i := range numKeys {
		cache.Set(keyFor(i), []byte("seed-value"), DefaultTestExpiration)
	}

	var wg sync.WaitGroup
	for w := range numWorkers {
		wg.Add(1)
		go func(worker int) {
			defer wg.Done()
			for op := range numOpsEach {
				k := keyFor((worker + op) % numKeys)
				switch op % 2 {
				case 0:
					cache.Set(k, []byte("value"), DefaultTestExpiration)
				case 1:
					cache.Delete(k)
				}
			}
		}(w)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for range 100 {
			cache.Clear()
		}
	}()

	wg.Wait()

	// Stop the background cleanup goroutine before recomputing expected
	// state, so it cannot race with the Range below (nothing has expired
	// yet - DefaultTestExpiration is 5s - but this removes any doubt).
	cache.Shutdown()

	var wantCount, wantMemory int64
	cache.entries.Range(func(_, v any) bool {
		entry := v.(CacheEntry)
		wantCount++
		wantMemory += entry.MemorySize
		return true
	})

	if got := cache.CountQueries(); got != wantCount {
		t.Fatalf("entryCount = %d, want %d (actual live entries after join)", got, wantCount)
	}
	if got := cache.GetMemoryUsage(); got != wantMemory {
		t.Fatalf("memoryUsage = %d, want %d (sum of live entries' MemorySize after join)", got, wantMemory)
	}
}

func TestMemoryCacheClose_IdempotentAndStillUsable(t *testing.T) {
	// Compile-time guard: the memory backend must satisfy io.Closer so
	// cache.Shutdown() can stop its cleanup goroutine (it otherwise leaks
	// the background timer until process exit), matching the Redis backend.
	var _ io.Closer = (*Cache)(nil)

	cache := New(DefaultTestExpiration)
	cache.Set("key", []byte("value"), DefaultTestExpiration)

	// Close is idempotent: it only cancels the background cleaner.
	if err := cache.Close(); err != nil {
		t.Fatalf("first Close: %v", err)
	}
	if err := cache.Close(); err != nil {
		t.Fatalf("second Close: %v", err)
	}

	// Stopping the cleaner must not render the cache unusable.
	if _, ok := cache.Get("key"); !ok {
		t.Fatal("entry should survive Close; Close only stops the background cleaner")
	}
	cache.Set("key2", []byte("v2"), DefaultTestExpiration)
	if _, ok := cache.Get("key2"); !ok {
		t.Fatal("cache should remain writable after Close")
	}
}

// TestNewWithSize_NonPositiveTTL_DoesNotPanic is a regression test for C1: a
// non-positive globalTTL (e.g. a misconfigured CACHE_TTL<=0) used to reach
// time.NewTicker(globalTTL/4) unguarded inside cleanupRoutine, and
// time.NewTicker panics for d<=0. That panic happens in the background
// goroutine spawned by NewWithSize, so a regression here crashes the whole
// process, not just this test.
func TestNewWithSize_NonPositiveTTL_DoesNotPanic(t *testing.T) {
	tests := []struct {
		name string
		ttl  time.Duration
	}{
		{"zero TTL", 0},
		{"negative TTL", -5 * time.Second},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := NewWithSize(tt.ttl, DefaultMaxMemorySize, DefaultMaxCacheSize)
			defer cache.Shutdown()

			// Give the background cleanup goroutine a moment to construct its
			// ticker.
			time.Sleep(20 * time.Millisecond)

			// Cache must remain fully usable.
			cache.Set("key", []byte("value"), DefaultTestExpiration)
			retrieved, found := cache.Get("key")
			assert.True(t, found, "cache should be usable after construction with a non-positive TTL")
			assert.Equal(t, []byte("value"), retrieved)
		})
	}
}

// TestDefaultTTLValue documents the sane fallback used both to size the
// cleanup ticker (C1) and to clamp a non-positive per-entry TTL (C11).
func TestDefaultTTLValue(t *testing.T) {
	assert.Equal(t, 60*time.Second, defaultTTL)
}
