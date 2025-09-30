package libpack_cache_memory

import (
	"bytes"

	"github.com/lukaszraczylo/graphql-monitoring-proxy/pkg/pools"
)

// GetBuffer gets a buffer from the pool (delegates to unified implementation)
func GetBuffer() *bytes.Buffer {
	return pools.GetBuffer()
}

// PutBuffer returns a buffer to the pool (delegates to unified implementation)
func PutBuffer(buf *bytes.Buffer) {
	pools.PutBuffer(buf)
}
