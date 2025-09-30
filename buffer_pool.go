package main

import (
	"bytes"
	"compress/gzip"
	"io"

	"github.com/lukaszraczylo/graphql-monitoring-proxy/pkg/pools"
)

// Legacy compatibility layer - delegates to unified pool implementation

// GetHTTPBuffer gets a buffer from the global pool
func GetHTTPBuffer() *bytes.Buffer {
	return pools.GetBuffer()
}

// PutHTTPBuffer returns a buffer to the global pool
func PutHTTPBuffer(buf *bytes.Buffer) {
	pools.PutBuffer(buf)
}

// GetGzipWriter gets a gzip writer from the global pool
func GetGzipWriter(w io.Writer) *gzip.Writer {
	return pools.GetGzipWriter(w)
}

// PutGzipWriter returns a gzip writer to the global pool
func PutGzipWriter(gz *gzip.Writer) {
	pools.PutGzipWriter(gz)
}

// GetGzipReader gets a gzip reader from the global pool
func GetGzipReader(r io.Reader) (*gzip.Reader, error) {
	return pools.GetGzipReader(r)
}

// PutGzipReader returns a gzip reader to the global pool
func PutGzipReader(gr *gzip.Reader) {
	pools.PutGzipReader(gr)
}
