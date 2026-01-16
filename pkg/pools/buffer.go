// Package pools provides memory-efficient buffer and gzip reader pools
// for reducing allocations in high-throughput request processing.
// Buffers are automatically sized and recycled to minimize GC pressure.
package pools

import (
	"bytes"
	"compress/gzip"
	"io"
	"sync"
)

const (
	// MaxBufferSize is the maximum size of a buffer that will be returned to the pool
	MaxBufferSize = 1024 * 1024 // 1MB
	// InitialBufferSize is the initial capacity of buffers in the pool
	InitialBufferSize = 4096 // 4KB
)

// bufferPool is the global pool for reusable buffers
var bufferPool = &sync.Pool{
	New: func() any {
		return bytes.NewBuffer(make([]byte, 0, InitialBufferSize))
	},
}

// gzipWriterPool is the global pool for reusable gzip writers
var gzipWriterPool = &sync.Pool{
	New: func() any {
		return gzip.NewWriter(nil)
	},
}

// gzipReaderPool is the global pool for reusable gzip readers
var gzipReaderPool = &sync.Pool{
	New: func() any {
		return new(gzip.Reader)
	},
}

// GetBuffer retrieves a buffer from the pool
func GetBuffer() *bytes.Buffer {
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	return buf
}

// PutBuffer returns a buffer to the pool
func PutBuffer(buf *bytes.Buffer) {
	if buf == nil {
		return
	}
	// Don't pool large buffers to avoid memory bloat
	if buf.Cap() > MaxBufferSize {
		return
	}
	buf.Reset()
	bufferPool.Put(buf)
}

// GetGzipWriter retrieves a gzip writer from the pool
func GetGzipWriter(w io.Writer) *gzip.Writer {
	gz := gzipWriterPool.Get().(*gzip.Writer)
	gz.Reset(w)
	return gz
}

// PutGzipWriter returns a gzip writer to the pool
func PutGzipWriter(gz *gzip.Writer) {
	if gz == nil {
		return
	}
	gz.Reset(nil)
	gzipWriterPool.Put(gz)
}

// GetGzipReader retrieves a gzip reader from the pool
func GetGzipReader(r io.Reader) (*gzip.Reader, error) {
	gr := gzipReaderPool.Get().(*gzip.Reader)
	if err := gr.Reset(r); err != nil {
		// If reset fails, create a new reader
		return gzip.NewReader(r)
	}
	return gr, nil
}

// PutGzipReader returns a gzip reader to the pool
func PutGzipReader(gr *gzip.Reader) {
	if gr == nil {
		return
	}
	gr.Close()
	gzipReaderPool.Put(gr)
}

// Stats provides statistics about the buffer pool usage
type Stats struct {
	BuffersInUse  int
	MaxBufferSize int
}

// GetStats returns current pool statistics (placeholder for future monitoring)
func GetStats() Stats {
	// This is a placeholder for future implementation
	// sync.Pool doesn't provide direct statistics access
	return Stats{
		BuffersInUse:  0,
		MaxBufferSize: MaxBufferSize,
	}
}
