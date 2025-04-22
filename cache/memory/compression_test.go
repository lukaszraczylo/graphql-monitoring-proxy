package libpack_cache_memory

import (
	"bytes"
	"compress/gzip"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestCompressionThreshold tests that values are only compressed when they exceed the threshold
func TestCompressionThreshold(t *testing.T) {
	cache := New(5 * time.Second)

	// Create test values
	smallValue := make([]byte, CompressionThreshold-100) // Below threshold
	largeValue := make([]byte, CompressionThreshold*2)   // Above threshold

	// Fill values with compressible data (repeating patterns compress well)
	for i := 0; i < len(smallValue); i++ {
		smallValue[i] = byte(i % 10)
	}
	for i := 0; i < len(largeValue); i++ {
		largeValue[i] = byte(i % 10)
	}

	// Test small value
	cache.Set("small-key", smallValue, 5*time.Second)

	// Extract the entry directly from the cache to check if it's compressed
	entryRaw, found := cache.entries.Load("small-key")
	assert.True(t, found, "Entry should exist")

	entry := entryRaw.(CacheEntry)
	assert.False(t, entry.Compressed, "Small value should not be compressed")
	assert.Equal(t, smallValue, entry.Value, "Small value should be stored as-is")

	// Test large value
	cache.Set("large-key", largeValue, 5*time.Second)

	entryRaw, found = cache.entries.Load("large-key")
	assert.True(t, found, "Entry should exist")

	entry = entryRaw.(CacheEntry)
	assert.True(t, entry.Compressed, "Large value should be compressed")

	// Ensure the stored value isn't the original
	assert.NotEqual(t, largeValue, entry.Value, "Large value should not be stored as-is")

	// Verify the value is actually compressed (should be smaller)
	assert.Less(t, len(entry.Value), len(largeValue), "Compressed value should be smaller than original")

	// Verify we can retrieve the uncompressed value correctly
	retrievedLarge, found := cache.Get("large-key")
	assert.True(t, found, "Large value should be retrievable")
	assert.Equal(t, largeValue, retrievedLarge, "Retrieved large value should match original")
}

// TestCompressionMemoryUsage tests that memory usage is calculated correctly for compressed entries
func TestCompressionMemoryUsage(t *testing.T) {
	cache := New(5 * time.Second)

	// Create a large, highly compressible value
	valueSize := CompressionThreshold * 4
	value := make([]byte, valueSize)
	for i := 0; i < valueSize; i++ {
		value[i] = byte(i % 2) // Highly compressible pattern (alternating 0s and 1s)
	}

	// Get initial memory usage
	initialMemUsage := cache.GetMemoryUsage()

	// Add the value
	key := "large-compressible-key"
	cache.Set(key, value, 5*time.Second)

	// Get memory usage after adding
	newMemUsage := cache.GetMemoryUsage()

	// The memory usage increase should be less than the full value size due to compression
	memUsageIncrease := newMemUsage - initialMemUsage

	// Extract the entry to check its compressed size
	entryRaw, found := cache.entries.Load(key)
	assert.True(t, found, "Entry should exist")

	entry := entryRaw.(CacheEntry)
	assert.True(t, entry.Compressed, "Value should be compressed")

	// Verify the reported memory usage matches the compressed size + overheads
	compressedSize := int64(len(entry.Value))
	keySize := int64(len(key))
	expectedUsage := compressedSize + keySize + approxEntryOverhead

	// The memory usage should reflect the compressed size, not the original size
	assert.InDelta(t, expectedUsage, memUsageIncrease, float64(approxEntryOverhead),
		"Memory usage should be based on compressed size")

	// Verify memory usage is correctly updated after deletion
	cache.Delete(key)
	finalMemUsage := cache.GetMemoryUsage()
	assert.Equal(t, initialMemUsage, finalMemUsage,
		"Memory usage should return to initial value after deletion")
}

// TestUncompressibleData tests the case where compression doesn't reduce size
func TestUncompressibleData(t *testing.T) {
	cache := New(5 * time.Second)

	// Create a large, random (less compressible) value
	valueSize := CompressionThreshold * 2

	// Create pseudo-random data that doesn't compress well
	// Using a custom PRNG for deterministic results across test runs
	value := make([]byte, valueSize)
	seed := uint32(42)
	for i := 0; i < valueSize; i++ {
		// Simple linear congruential generator
		seed = seed*1664525 + 1013904223
		value[i] = byte(seed)
	}

	// Try to compress it directly to see if it actually would reduce size
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	gw.Write(value)
	gw.Close()
	compressedDirectly := buf.Bytes()

	// Now use the cache's Set method
	key := "uncompressible-key"
	cache.Set(key, value, 5*time.Second)

	// Extract the entry to check if it's compressed
	entryRaw, found := cache.entries.Load(key)
	assert.True(t, found, "Entry should exist")

	entry := entryRaw.(CacheEntry)

	// If our test data actually compressed to a smaller size, we expect the cache to store it compressed
	if len(compressedDirectly) < len(value) {
		assert.True(t, entry.Compressed, "Value should be stored compressed if smaller")
		assert.Less(t, len(entry.Value), len(value), "Compressed value should be smaller")
	} else {
		// Uncommon case: our pseudo-random data actually expanded with gzip
		// In this case, the cache should store it uncompressed
		assert.False(t, entry.Compressed, "Value should not be compressed if it would expand")
		assert.Equal(t, value, entry.Value, "Value should be stored as-is")
	}

	// Regardless, we should be able to get the correct value back
	retrievedValue, found := cache.Get(key)
	assert.True(t, found, "Value should be retrievable")
	assert.Equal(t, value, retrievedValue, "Retrieved value should match original")
}

// TestCompressDecompressDirectly tests the compress and decompress methods directly
func TestCompressDecompressDirectly(t *testing.T) {
	cache := New(5 * time.Second)

	// Test with various sizes
	testSizes := []int{
		100,                      // Small
		CompressionThreshold - 1, // Just below threshold
		CompressionThreshold,     // At threshold
		CompressionThreshold + 1, // Just above threshold
		CompressionThreshold * 2, // Well above threshold
	}

	for _, size := range testSizes {
		t.Run("Size-"+string(rune('A'+len(testSizes)%26)), func(t *testing.T) {
			// Generate test data with a repeating pattern
			data := make([]byte, size)
			for i := 0; i < size; i++ {
				data[i] = byte(i % 256)
			}

			// Compress the data
			compressed, err := cache.compress(data)
			assert.NoError(t, err, "Compression should not error")

			// Small data may get larger when compressed, larger data should get smaller
			if size > CompressionThreshold {
				assert.Less(t, len(compressed), len(data),
					"Compression should reduce size for data above threshold")
			}

			// Decompress and verify it matches the original
			decompressed, err := cache.decompress(compressed)
			assert.NoError(t, err, "Decompression should not error")
			assert.Equal(t, data, decompressed, "Data should round-trip correctly through compression")
		})
	}
}

// TestDecompressInvalidData tests handling invalid data in decompress
func TestDecompressInvalidData(t *testing.T) {
	cache := New(5 * time.Second)

	// Try to decompress non-gzip data
	invalidData := []byte("This is not valid gzip data")
	_, err := cache.decompress(invalidData)
	assert.Error(t, err, "Decompressing invalid data should return error")

	// Set compressed flag but store invalid data
	key := "invalid-compressed-key"
	cache.entries.Store(key, CacheEntry{
		Value:      invalidData,
		ExpiresAt:  time.Now().Add(5 * time.Second),
		Compressed: true, // Flag as compressed even though it's not
		MemorySize: int64(len(invalidData) + len(key) + approxEntryOverhead),
	})

	// Try to get it - should fail gracefully
	_, found := cache.Get(key)
	assert.False(t, found, "Get should fail gracefully for invalid compressed data")
}
