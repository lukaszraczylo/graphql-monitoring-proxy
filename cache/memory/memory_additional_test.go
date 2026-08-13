package libpack_cache_memory

import (
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
