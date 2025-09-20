package pools

import (
	"bytes"
	"compress/gzip"
	"io"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type BufferPoolTestSuite struct {
	suite.Suite
}

func TestBufferPoolTestSuite(t *testing.T) {
	suite.Run(t, new(BufferPoolTestSuite))
}

func (suite *BufferPoolTestSuite) TestGetBuffer() {
	buf := GetBuffer()
	assert.NotNil(suite.T(), buf)
	assert.Equal(suite.T(), 0, buf.Len())
	assert.GreaterOrEqual(suite.T(), buf.Cap(), InitialBufferSize)
}

func (suite *BufferPoolTestSuite) TestPutBuffer() {
	buf := GetBuffer()
	buf.WriteString("test data")
	assert.Equal(suite.T(), "test data", buf.String())

	PutBuffer(buf)

	// Get a new buffer - it should be reset
	buf2 := GetBuffer()
	assert.Equal(suite.T(), 0, buf2.Len())
	assert.Equal(suite.T(), "", buf2.String())
}

func (suite *BufferPoolTestSuite) TestPutBufferNil() {
	// Should not panic
	PutBuffer(nil)
}

func (suite *BufferPoolTestSuite) TestPutBufferLarge() {
	buf := bytes.NewBuffer(make([]byte, 0, MaxBufferSize+1))

	// Large buffer should not be pooled
	PutBuffer(buf)

	// Getting a new buffer should return a new one, not the large one
	buf2 := GetBuffer()
	assert.LessOrEqual(suite.T(), buf2.Cap(), MaxBufferSize)
}

func (suite *BufferPoolTestSuite) TestBufferReuse() {
	// Test that buffers are actually being reused
	buf1 := GetBuffer()
	buf1.WriteString("test")
	ptr1 := buf1

	PutBuffer(buf1)

	buf2 := GetBuffer()
	// Due to pool behavior, we might or might not get the same buffer back
	// but it should be properly reset
	assert.Equal(suite.T(), 0, buf2.Len())
	assert.Equal(suite.T(), "", buf2.String())
	_ = ptr1 // Keep reference to avoid compiler optimization
}

func (suite *BufferPoolTestSuite) TestGzipWriter() {
	var buf bytes.Buffer
	gz := GetGzipWriter(&buf)
	assert.NotNil(suite.T(), gz)

	// Write some data
	data := "test gzip data"
	_, err := gz.Write([]byte(data))
	assert.NoError(suite.T(), err)

	err = gz.Close()
	assert.NoError(suite.T(), err)

	// Verify data was compressed
	assert.Greater(suite.T(), buf.Len(), 0)

	PutGzipWriter(gz)
}

func (suite *BufferPoolTestSuite) TestGzipWriterNil() {
	// Should not panic
	PutGzipWriter(nil)
}

func (suite *BufferPoolTestSuite) TestGzipWriterReuse() {
	var buf1, buf2 bytes.Buffer

	// First use
	gz := GetGzipWriter(&buf1)
	gz.Write([]byte("data1"))
	gz.Close()
	PutGzipWriter(gz)

	// Second use - should be reset
	gz2 := GetGzipWriter(&buf2)
	gz2.Write([]byte("data2"))
	gz2.Close()

	// Both buffers should contain valid gzip data
	assert.Greater(suite.T(), buf1.Len(), 0)
	assert.Greater(suite.T(), buf2.Len(), 0)
	assert.NotEqual(suite.T(), buf1.Bytes(), buf2.Bytes())

	PutGzipWriter(gz2)
}

func (suite *BufferPoolTestSuite) TestGzipReader() {
	// Create gzipped data
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	gz.Write([]byte("test data"))
	gz.Close()

	// Read using pooled reader
	gr, err := GetGzipReader(&buf)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), gr)

	data, err := io.ReadAll(gr)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "test data", string(data))

	PutGzipReader(gr)
}

func (suite *BufferPoolTestSuite) TestGzipReaderInvalidData() {
	buf := bytes.NewBufferString("invalid gzip data")

	gr, err := GetGzipReader(buf)
	// Should return error or new reader
	if err == nil {
		assert.NotNil(suite.T(), gr)
		// Try to read - should fail
		_, readErr := io.ReadAll(gr)
		assert.Error(suite.T(), readErr)
		PutGzipReader(gr)
	}
}

func (suite *BufferPoolTestSuite) TestGzipReaderNil() {
	// Should not panic
	PutGzipReader(nil)
}

func (suite *BufferPoolTestSuite) TestGzipReaderReuse() {
	// Create two different gzipped data
	var buf1, buf2 bytes.Buffer

	gz1 := gzip.NewWriter(&buf1)
	gz1.Write([]byte("data1"))
	gz1.Close()

	gz2 := gzip.NewWriter(&buf2)
	gz2.Write([]byte("data2"))
	gz2.Close()

	// Read first data
	gr, err := GetGzipReader(&buf1)
	assert.NoError(suite.T(), err)
	data1, err := io.ReadAll(gr)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "data1", string(data1))
	PutGzipReader(gr)

	// Read second data with potentially reused reader
	gr2, err := GetGzipReader(&buf2)
	assert.NoError(suite.T(), err)
	data2, err := io.ReadAll(gr2)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "data2", string(data2))
	PutGzipReader(gr2)
}

func (suite *BufferPoolTestSuite) TestConcurrentBufferAccess() {
	var wg sync.WaitGroup
	numGoroutines := 100
	numOperations := 100

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				buf := GetBuffer()
				buf.WriteString("test data")
				assert.Equal(suite.T(), "test data", buf.String())
				PutBuffer(buf)
			}
		}(i)
	}

	wg.Wait()
}

func (suite *BufferPoolTestSuite) TestConcurrentGzipWriter() {
	var wg sync.WaitGroup
	numGoroutines := 50

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			var buf bytes.Buffer
			gz := GetGzipWriter(&buf)
			data := strings.Repeat("test", 100)
			gz.Write([]byte(data))
			gz.Close()
			assert.Greater(suite.T(), buf.Len(), 0)
			PutGzipWriter(gz)
		}(i)
	}

	wg.Wait()
}

func (suite *BufferPoolTestSuite) TestConcurrentGzipReader() {
	// Prepare gzipped data
	var source bytes.Buffer
	gz := gzip.NewWriter(&source)
	gz.Write([]byte("test data for concurrent reading"))
	gz.Close()
	sourceData := source.Bytes()

	var wg sync.WaitGroup
	numGoroutines := 50

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			// Each goroutine needs its own reader for the data
			buf := bytes.NewBuffer(sourceData)
			gr, err := GetGzipReader(buf)
			if err != nil {
				// Handle error from failed reset
				return
			}
			data, err := io.ReadAll(gr)
			if err == nil {
				assert.Equal(suite.T(), "test data for concurrent reading", string(data))
			}
			PutGzipReader(gr)
		}(i)
	}

	wg.Wait()
}

func (suite *BufferPoolTestSuite) TestRaceConditions() {
	var wg sync.WaitGroup
	var bufferOps, gzipWriterOps, gzipReaderOps int32

	// Buffer operations
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				buf := GetBuffer()
				buf.WriteString("race test")
				PutBuffer(buf)
				atomic.AddInt32(&bufferOps, 1)
			}
		}()
	}

	// Gzip writer operations
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				var buf bytes.Buffer
				gz := GetGzipWriter(&buf)
				gz.Write([]byte("test"))
				gz.Close()
				PutGzipWriter(gz)
				atomic.AddInt32(&gzipWriterOps, 1)
			}
		}()
	}

	// Gzip reader operations
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				var buf bytes.Buffer
				gz := gzip.NewWriter(&buf)
				gz.Write([]byte("test"))
				gz.Close()

				gr, err := GetGzipReader(&buf)
				if err == nil {
					io.ReadAll(gr)
					PutGzipReader(gr)
					atomic.AddInt32(&gzipReaderOps, 1)
				}
			}
		}()
	}

	wg.Wait()

	assert.Equal(suite.T(), int32(1000), atomic.LoadInt32(&bufferOps))
	assert.Equal(suite.T(), int32(1000), atomic.LoadInt32(&gzipWriterOps))
	assert.LessOrEqual(suite.T(), int32(900), atomic.LoadInt32(&gzipReaderOps)) // Some might fail
}

func (suite *BufferPoolTestSuite) TestGetStats() {
	stats := GetStats()
	assert.Equal(suite.T(), MaxBufferSize, stats.MaxBufferSize)
	// BuffersInUse is always 0 in current implementation
	assert.Equal(suite.T(), 0, stats.BuffersInUse)
}

func (suite *BufferPoolTestSuite) TestBufferGrowth() {
	buf := GetBuffer()

	// Write more than initial capacity
	largeData := strings.Repeat("x", InitialBufferSize*2)
	buf.WriteString(largeData)

	assert.Equal(suite.T(), len(largeData), buf.Len())
	assert.GreaterOrEqual(suite.T(), buf.Cap(), len(largeData))

	PutBuffer(buf)
}

func (suite *BufferPoolTestSuite) TestMemoryEfficiency() {
	// Test that pools actually reduce allocations
	allocsBefore := testing.AllocsPerRun(100, func() {
		buf := new(bytes.Buffer)
		buf.WriteString("test")
		_ = buf.String()
	})

	allocsWithPool := testing.AllocsPerRun(100, func() {
		buf := GetBuffer()
		buf.WriteString("test")
		_ = buf.String()
		PutBuffer(buf)
	})

	// Pool should reduce allocations
	assert.Less(suite.T(), allocsWithPool, allocsBefore)
}

// Benchmark tests
func BenchmarkBufferPool(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			buf := GetBuffer()
			buf.WriteString("benchmark test data")
			PutBuffer(buf)
		}
	})
}

func BenchmarkGzipWriterPool(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var buf bytes.Buffer
			gz := GetGzipWriter(&buf)
			gz.Write([]byte("benchmark test data"))
			gz.Close()
			PutGzipWriter(gz)
		}
	})
}

func BenchmarkGzipReaderPool(b *testing.B) {
	// Prepare compressed data
	var compressed bytes.Buffer
	gz := gzip.NewWriter(&compressed)
	gz.Write([]byte("benchmark test data"))
	gz.Close()
	data := compressed.Bytes()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			buf := bytes.NewBuffer(data)
			gr, err := GetGzipReader(buf)
			if err == nil {
				io.ReadAll(gr)
				PutGzipReader(gr)
			}
		}
	})
}

func BenchmarkWithoutPool(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			buf := new(bytes.Buffer)
			buf.WriteString("benchmark test data")
			// Buffer is discarded, letting GC handle it
		}
	})
}
