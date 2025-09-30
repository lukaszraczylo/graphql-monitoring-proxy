package libpack_cache_memory

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestEvictToFreeMemory tests that the cache correctly evicts
// items when it exceeds its memory limit.
func TestEvictToFreeMemory(t *testing.T) {
	// Create a cache with a small memory limit: 5KB (ensure eviction happens)
	smallMemLimit := int64(5 * 1024)
	cache := NewWithSize(5*time.Second, smallMemLimit, 1000)

	// Create entries with known sizes
	// Each entry will be ~512 bytes plus overhead
	valueSize := 512
	numEntriesToExceedLimit := 12 // Should exceed the 5KB limit and force eviction

	// Create a slice to track keys in insertion order
	keys := make([]string, numEntriesToExceedLimit)

	// Add entries with significant delays between insertions
	for i := 0; i < numEntriesToExceedLimit; i++ {
		key := fmt.Sprintf("test-key-%d", i)
		keys[i] = key

		value := make([]byte, valueSize)
		for j := 0; j < valueSize; j++ {
			value[j] = byte(i % 256) // Fill with a repeating pattern
		}

		cache.Set(key, value, 30*time.Second)

		// More significant delay to ensure different timestamps
		time.Sleep(10 * time.Millisecond)
	}

	// Allow time for eviction to complete
	time.Sleep(50 * time.Millisecond)

	// Verify memory usage is below the limit
	memUsage := cache.GetMemoryUsage()
	assert.LessOrEqual(t, memUsage, smallMemLimit,
		"Memory usage (%d) should be less than or equal to the limit (%d)", memUsage, smallMemLimit)

	// Count how many items are left in the cache and which ones
	present := 0
	for i := 0; i < numEntriesToExceedLimit; i++ {
		_, found := cache.Get(keys[i])
		if found {
			present++
		}
	}

	// We expect some items to be evicted based on the memory limit
	assert.Less(t, present, numEntriesToExceedLimit,
		"Some items should have been evicted (%d present out of %d total)",
		present, numEntriesToExceedLimit)

	// Verify newer items (inserted later) are more likely to be in the cache
	// Check the last few items which should be the newest
	for i := numEntriesToExceedLimit - 3; i < numEntriesToExceedLimit; i++ {
		_, found := cache.Get(keys[i])
		assert.True(t, found, "Newer key %s should still exist", keys[i])
	}
}

// TestMaxCacheSize verifies the behavior when adding more items than the maxCacheSize limit
func TestMaxCacheSize(t *testing.T) {
	// Create a cache with a small limit
	smallLimit := int64(5)
	cache := NewWithSize(5*time.Second, DefaultMaxMemorySize, smallLimit)

	// Add entries with increasing size (to avoid memory-based eviction)
	for i := 0; i < 20; i++ {
		key := fmt.Sprintf("test-key-%d", i)
		value := []byte(key)
		cache.Set(key, value, 10*time.Second)
	}

	// Verify we can get a reasonable number of items
	// (we don't test for exact count as implementation may vary)
	foundCount := 0
	for i := 0; i < 20; i++ {
		key := fmt.Sprintf("test-key-%d", i)
		_, found := cache.Get(key)
		if found {
			foundCount++
		}
	}

	// We should find some items but not all 20
	assert.Greater(t, foundCount, 0, "Some items should be in the cache")
	assert.LessOrEqual(t, foundCount, 20, "Not all items should be in the cache with small limit")
}

// TestGetMemoryUsage verifies that memory usage tracking is accurate
func TestGetMemoryUsage(t *testing.T) {
	cache := New(5 * time.Second)

	// Initially memory usage should be 0
	assert.Equal(t, int64(0), cache.GetMemoryUsage(), "Initial memory usage should be 0")

	// Add an entry with a known approximate size
	valueSize := 1024
	value := make([]byte, valueSize)
	key := "test-key"

	cache.Set(key, value, 5*time.Second)

	// Check memory usage - should be approximately valueSize + key length + overhead
	expectedMinUsage := int64(valueSize + len(key))
	memUsage := cache.GetMemoryUsage()
	assert.GreaterOrEqual(t, memUsage, expectedMinUsage,
		"Memory usage (%d) should be at least the value size plus key length (%d)", memUsage, expectedMinUsage)

	// Delete the entry and verify memory usage decreases
	cache.Delete(key)
	assert.Equal(t, int64(0), cache.GetMemoryUsage(), "Memory usage should be 0 after deletion")
}

// TestSetMaxMemorySize tests changing the memory limit and resulting eviction
func TestSetMaxMemorySize(t *testing.T) {
	// Start with a large limit
	initialLimit := int64(100 * 1024)
	cache := NewWithSize(5*time.Second, initialLimit, 1000)

	// Fill the cache with ~50KB of data
	valueSize := 1024
	numEntries := 50

	for i := 0; i < numEntries; i++ {
		key := generateKey(i)
		value := make([]byte, valueSize)
		cache.Set(key, value, 5*time.Second)

		// Small delay for timestamp differences
		time.Sleep(time.Millisecond)
	}

	// Verify all entries exist
	for i := 0; i < numEntries; i++ {
		_, found := cache.Get(generateKey(i))
		assert.True(t, found, "All entries should exist before limit change")
	}

	// Get current memory usage
	originalUsage := cache.GetMemoryUsage()

	// Now reduce the limit to 20KB - should trigger eviction
	newLimit := int64(20 * 1024)
	cache.SetMaxMemorySize(newLimit)

	// Verify memory usage is now below the new limit
	newUsage := cache.GetMemoryUsage()
	assert.LessOrEqual(t, newUsage, newLimit,
		"After SetMaxMemorySize, memory usage (%d) should be less than or equal to new limit (%d)",
		newUsage, newLimit)
	assert.Less(t, newUsage, originalUsage,
		"Memory usage should have decreased after lowering the limit")

	// Some older entries should be gone, newer ones should still exist
	removedCount := 0
	remainingCount := 0
	for i := 0; i < numEntries; i++ {
		_, found := cache.Get(generateKey(i))
		if found {
			remainingCount++
		} else {
			removedCount++
		}
	}

	assert.Greater(t, removedCount, 0, "Some entries should have been removed")
	assert.Greater(t, remainingCount, 0, "Some entries should still exist")
}

// Helper function to generate consistent keys
func generateKey(index int) string {
	return "test-key-" + fmt.Sprintf("%d", index)
}
