package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"math/rand"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/lukaszraczylo/graphql-monitoring-proxy/pkg/pools"
	"github.com/stretchr/testify/suite"
)

type PoolsSecurityTestSuite struct {
	suite.Suite
}

func TestPoolsSecurityTestSuite(t *testing.T) {
	suite.Run(t, new(PoolsSecurityTestSuite))
}

// TestBufferPoolConcurrency tests concurrent Get/Put operations for thread safety
func (suite *PoolsSecurityTestSuite) TestBufferPoolConcurrency() {
	const numGoroutines = 100
	const numOperationsPerGoroutine = 100

	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines*numOperationsPerGoroutine)

	suite.Run("Concurrent buffer pool operations", func() {
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(goroutineID int) {
				defer wg.Done()

				for j := 0; j < numOperationsPerGoroutine; j++ {
					// Get buffer from pool
					buf := pools.GetBuffer()
					if buf == nil {
						errors <- fmt.Errorf("goroutine %d, iteration %d: got nil buffer", goroutineID, j)
						continue
					}

					// Verify buffer is reset/clean
					if buf.Len() != 0 {
						errors <- fmt.Errorf("goroutine %d, iteration %d: buffer not reset, length: %d", goroutineID, j, buf.Len())
						continue
					}

					// Use the buffer
					testData := fmt.Sprintf("test data from goroutine %d iteration %d", goroutineID, j)
					buf.WriteString(testData)

					// Verify data was written correctly
					if buf.String() != testData {
						errors <- fmt.Errorf("goroutine %d, iteration %d: data corruption", goroutineID, j)
						continue
					}

					// Return buffer to pool
					pools.PutBuffer(buf)

					// Small random delay to increase chance of race conditions
					if rand.Intn(10) == 0 {
						time.Sleep(time.Microsecond)
					}
				}
			}(i)
		}

		wg.Wait()
		close(errors)

		// Check for any errors
		errorCount := 0
		for err := range errors {
			suite.T().Errorf("Concurrent operation failed: %v", err)
			errorCount++
		}

		suite.Equal(0, errorCount, "Should have no errors in concurrent operations")
	})
}

// TestBufferPoolMemoryLeak tests for memory leaks in buffer pooling
func (suite *PoolsSecurityTestSuite) TestBufferPoolMemoryLeak() {
	suite.Run("Memory leak prevention", func() {
		var memBefore runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&memBefore)

		// Create many buffers and return them to pool
		const numBuffers = 1000
		buffers := make([]*bytes.Buffer, numBuffers)

		for i := 0; i < numBuffers; i++ {
			buffers[i] = pools.GetBuffer()
			// Write some data
			buffers[i].WriteString(strings.Repeat("a", 1024))
		}

		// Return all buffers to pool
		for i := 0; i < numBuffers; i++ {
			pools.PutBuffer(buffers[i])
		}

		// Clear references
		for i := range buffers {
			buffers[i] = nil
		}
		buffers = nil

		// Force garbage collection
		runtime.GC()
		runtime.GC() // Second GC to ensure cleanup

		var memAfter runtime.MemStats
		runtime.ReadMemStats(&memAfter)

		// Memory usage shouldn't increase dramatically
		memDiff := int64(memAfter.Alloc) - int64(memBefore.Alloc)
		maxAcceptableIncrease := int64(1024 * 1024) // 1MB

		suite.LessOrEqual(memDiff, maxAcceptableIncrease,
			"Memory usage increased by %d bytes, should be less than %d bytes",
			memDiff, maxAcceptableIncrease)
	})
}

// TestBufferSizeLimit tests that oversized buffers are not pooled
func (suite *PoolsSecurityTestSuite) TestBufferSizeLimit() {
	suite.Run("Oversized buffer rejection", func() {
		buf := pools.GetBuffer()

		// Write data larger than MaxBufferSize
		largeData := make([]byte, pools.MaxBufferSize+1)
		for i := range largeData {
			largeData[i] = 'a'
		}
		buf.Write(largeData)

		// Verify buffer is oversized
		suite.Greater(buf.Cap(), pools.MaxBufferSize,
			"Buffer capacity should exceed MaxBufferSize")

		// Return oversized buffer to pool
		pools.PutBuffer(buf)

		// Get a new buffer - should be a fresh one, not the oversized one
		newBuf := pools.GetBuffer()
		suite.Equal(0, newBuf.Len(), "New buffer should be empty")
		suite.LessOrEqual(newBuf.Cap(), pools.MaxBufferSize,
			"New buffer capacity should be within limits")

		pools.PutBuffer(newBuf)
	})
}

// TestBufferPoolRaceConditions tests for race conditions in buffer pooling
func (suite *PoolsSecurityTestSuite) TestBufferPoolRaceConditions() {
	suite.Run("Race condition detection", func() {
		const numGoroutines = 50
		var wg sync.WaitGroup
		bufferMap := sync.Map{} // Track buffers to detect sharing

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(goroutineID int) {
				defer wg.Done()

				for j := 0; j < 50; j++ {
					buf := pools.GetBuffer()
					bufferAddr := fmt.Sprintf("%p", buf)

					// Check if this buffer is already in use
					if _, exists := bufferMap.LoadOrStore(bufferAddr, goroutineID); exists {
						suite.T().Errorf("Buffer %s is being used by multiple goroutines", bufferAddr)
						return
					}

					// Use buffer
					buf.WriteString(fmt.Sprintf("goroutine-%d-op-%d", goroutineID, j))

					// Simulate some work
					time.Sleep(time.Microsecond * time.Duration(rand.Intn(10)))

					// Remove from tracking and return to pool
					bufferMap.Delete(bufferAddr)
					pools.PutBuffer(buf)
				}
			}(i)
		}

		wg.Wait()
	})
}

// TestGzipWriterPoolConcurrency tests concurrent operations on gzip writer pool
func (suite *PoolsSecurityTestSuite) TestGzipWriterPoolConcurrency() {
	const numGoroutines = 50
	const numOperationsPerGoroutine = 20

	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines*numOperationsPerGoroutine)

	suite.Run("Concurrent gzip writer pool operations", func() {
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(goroutineID int) {
				defer wg.Done()

				for j := 0; j < numOperationsPerGoroutine; j++ {
					// Create a buffer for compressed data
					buf := &bytes.Buffer{}

					// Get gzip writer from pool
					gz := pools.GetGzipWriter(buf)
					if gz == nil {
						errors <- fmt.Errorf("goroutine %d, iteration %d: got nil gzip writer", goroutineID, j)
						continue
					}

					// Write test data
					testData := fmt.Sprintf("test data from goroutine %d iteration %d", goroutineID, j)
					if _, err := gz.Write([]byte(testData)); err != nil {
						errors <- fmt.Errorf("goroutine %d, iteration %d: write error: %v", goroutineID, j, err)
						continue
					}

					if err := gz.Close(); err != nil {
						errors <- fmt.Errorf("goroutine %d, iteration %d: close error: %v", goroutineID, j, err)
						continue
					}

					// Verify compression worked
					if buf.Len() == 0 {
						errors <- fmt.Errorf("goroutine %d, iteration %d: no compressed data", goroutineID, j)
						continue
					}

					// Return writer to pool
					pools.PutGzipWriter(gz)
				}
			}(i)
		}

		wg.Wait()
		close(errors)

		// Check for any errors
		errorCount := 0
		for err := range errors {
			suite.T().Errorf("Concurrent gzip writer operation failed: %v", err)
			errorCount++
		}

		suite.Equal(0, errorCount, "Should have no errors in concurrent gzip writer operations")
	})
}

// TestGzipReaderPoolConcurrency tests concurrent operations on gzip reader pool
func (suite *PoolsSecurityTestSuite) TestGzipReaderPoolConcurrency() {
	// First, prepare some compressed data
	testData := "Hello, World! This is test data for gzip reader pool testing."
	var compressedBuf bytes.Buffer
	gz := gzip.NewWriter(&compressedBuf)
	gz.Write([]byte(testData))
	gz.Close()
	compressedData := compressedBuf.Bytes()

	const numGoroutines = 30
	const numOperationsPerGoroutine = 10

	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines*numOperationsPerGoroutine)

	suite.Run("Concurrent gzip reader pool operations", func() {
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(goroutineID int) {
				defer wg.Done()

				for j := 0; j < numOperationsPerGoroutine; j++ {
					// Create reader from compressed data
					reader := bytes.NewReader(compressedData)

					// Get gzip reader from pool
					gr, err := pools.GetGzipReader(reader)
					if err != nil {
						errors <- fmt.Errorf("goroutine %d, iteration %d: error getting gzip reader: %v", goroutineID, j, err)
						continue
					}

					// Read decompressed data
					decompressed, err := io.ReadAll(gr)
					if err != nil {
						errors <- fmt.Errorf("goroutine %d, iteration %d: read error: %v", goroutineID, j, err)
						continue
					}

					// Verify data integrity
					if string(decompressed) != testData {
						errors <- fmt.Errorf("goroutine %d, iteration %d: data mismatch", goroutineID, j)
						continue
					}

					// Return reader to pool
					pools.PutGzipReader(gr)
				}
			}(i)
		}

		wg.Wait()
		close(errors)

		// Check for any errors
		errorCount := 0
		for err := range errors {
			suite.T().Errorf("Concurrent gzip reader operation failed: %v", err)
			errorCount++
		}

		suite.Equal(0, errorCount, "Should have no errors in concurrent gzip reader operations")
	})
}

// TestPoolNilHandling tests proper handling of nil parameters
func (suite *PoolsSecurityTestSuite) TestPoolNilHandling() {
	suite.Run("Nil buffer handling", func() {
		// Should not panic when putting nil buffer
		suite.NotPanics(func() {
			pools.PutBuffer(nil)
		})
	})

	suite.Run("Nil gzip writer handling", func() {
		// Should not panic when putting nil gzip writer
		suite.NotPanics(func() {
			pools.PutGzipWriter(nil)
		})
	})

	suite.Run("Nil gzip reader handling", func() {
		// Should not panic when putting nil gzip reader
		suite.NotPanics(func() {
			pools.PutGzipReader(nil)
		})
	})
}

// TestPoolResourceExhaustion tests behavior under resource exhaustion
func (suite *PoolsSecurityTestSuite) TestPoolResourceExhaustion() {
	suite.Run("Buffer pool under pressure", func() {
		// Get many buffers without returning them
		const numBuffers = 10000
		buffers := make([]*bytes.Buffer, numBuffers)

		for i := 0; i < numBuffers; i++ {
			buffers[i] = pools.GetBuffer()
			suite.NotNil(buffers[i], "Should always get a buffer (pool should create new ones)")
		}

		// Each buffer should be functional
		for i := 0; i < numBuffers; i++ {
			buffers[i].WriteString("test")
			suite.Equal("test", buffers[i].String())
		}

		// Return all buffers
		for i := 0; i < numBuffers; i++ {
			pools.PutBuffer(buffers[i])
		}
	})
}

// TestPoolBufferReset tests that buffers are properly reset
func (suite *PoolsSecurityTestSuite) TestPoolBufferReset() {
	suite.Run("Buffer reset verification", func() {
		// Get a buffer and write data
		buf1 := pools.GetBuffer()
		buf1.WriteString("sensitive data")
		suite.Equal("sensitive data", buf1.String())

		// Return to pool
		pools.PutBuffer(buf1)

		// Get another buffer (might be the same one)
		buf2 := pools.GetBuffer()

		// Should be empty (reset)
		suite.Equal(0, buf2.Len(), "Buffer should be reset to empty")
		suite.Equal("", buf2.String(), "Buffer content should be empty")

		pools.PutBuffer(buf2)
	})
}

// TestPoolGzipWriterReset tests that gzip writers are properly reset
func (suite *PoolsSecurityTestSuite) TestPoolGzipWriterReset() {
	suite.Run("Gzip writer reset verification", func() {
		// First usage
		buf1 := &bytes.Buffer{}
		gz1 := pools.GetGzipWriter(buf1)
		gz1.Write([]byte("data1"))
		gz1.Close()

		pools.PutGzipWriter(gz1)

		// Second usage
		buf2 := &bytes.Buffer{}
		gz2 := pools.GetGzipWriter(buf2)
		gz2.Write([]byte("data2"))
		gz2.Close()

		// Decompress to verify only "data2" is present
		reader, err := gzip.NewReader(buf2)
		suite.NoError(err)

		decompressed, err := io.ReadAll(reader)
		suite.NoError(err)
		reader.Close()

		suite.Equal("data2", string(decompressed),
			"Gzip writer should be reset and not contain previous data")

		pools.PutGzipWriter(gz2)
	})
}

// TestPoolDataIsolation tests that data doesn't leak between pool uses
func (suite *PoolsSecurityTestSuite) TestPoolDataIsolation() {
	suite.Run("Buffer data isolation", func() {
		// Create sensitive data pattern
		sensitiveData := "password=secret123&api_key=sk-sensitive"

		// Use buffer with sensitive data
		buf1 := pools.GetBuffer()
		buf1.WriteString(sensitiveData)
		suite.Contains(buf1.String(), "secret123")

		// Return to pool
		pools.PutBuffer(buf1)

		// Get new buffer and use it
		buf2 := pools.GetBuffer()
		buf2.WriteString("public data")

		// Verify no sensitive data leaks
		bufContent := buf2.String()
		suite.NotContains(bufContent, "secret123", "Sensitive data should not leak")
		suite.NotContains(bufContent, "sk-sensitive", "API key should not leak")
		suite.Equal("public data", bufContent)

		pools.PutBuffer(buf2)
	})
}

// TestPoolIntegration tests integration between different pool types
func (suite *PoolsSecurityTestSuite) TestPoolIntegration() {
	suite.Run("Combined buffer and gzip operations", func() {
		const numOperations = 100
		var wg sync.WaitGroup
		errors := make(chan error, numOperations)

		for i := 0; i < numOperations; i++ {
			wg.Add(1)
			go func(opID int) {
				defer wg.Done()

				// Get buffer and gzip writer
				buf := pools.GetBuffer()
				gz := pools.GetGzipWriter(buf)

				// Write test data
				testData := fmt.Sprintf("operation %d test data", opID)
				if _, err := gz.Write([]byte(testData)); err != nil {
					errors <- fmt.Errorf("operation %d: write error: %v", opID, err)
					return
				}

				if err := gz.Close(); err != nil {
					errors <- fmt.Errorf("operation %d: close error: %v", opID, err)
					return
				}

				// Verify compression worked
				if buf.Len() == 0 {
					errors <- fmt.Errorf("operation %d: no compressed data", opID)
					return
				}

				// Test decompression with pool reader
				gr, err := pools.GetGzipReader(bytes.NewReader(buf.Bytes()))
				if err != nil {
					errors <- fmt.Errorf("operation %d: reader error: %v", opID, err)
					return
				}

				decompressed, err := io.ReadAll(gr)
				if err != nil {
					errors <- fmt.Errorf("operation %d: decompress error: %v", opID, err)
					return
				}

				if string(decompressed) != testData {
					errors <- fmt.Errorf("operation %d: data mismatch", opID)
					return
				}

				// Return everything to pools
				pools.PutGzipWriter(gz)
				pools.PutBuffer(buf)
				pools.PutGzipReader(gr)
			}(i)
		}

		wg.Wait()
		close(errors)

		// Check for errors
		errorCount := 0
		for err := range errors {
			suite.T().Errorf("Integration test failed: %v", err)
			errorCount++
		}

		suite.Equal(0, errorCount, "Should have no errors in integration tests")
	})
}

// BenchmarkBufferPoolOperations benchmarks buffer pool performance
func BenchmarkBufferPoolOperations(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			buf := pools.GetBuffer()
			buf.WriteString("benchmark test data")
			pools.PutBuffer(buf)
		}
	})
}

// BenchmarkGzipWriterPoolOperations benchmarks gzip writer pool performance
func BenchmarkGzipWriterPoolOperations(b *testing.B) {
	testData := []byte("benchmark test data for gzip compression")

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			buf := &bytes.Buffer{}
			gz := pools.GetGzipWriter(buf)
			gz.Write(testData)
			gz.Close()
			pools.PutGzipWriter(gz)
		}
	})
}
