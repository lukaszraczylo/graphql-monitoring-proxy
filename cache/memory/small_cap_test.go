package libpack_cache_memory

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Pins theorems/Theorems/MemCacheCap.lean: every cap >= 1 bounds the entry count.
func TestSet_SmallEntryCapIsEnforced(t *testing.T) {
	for _, maxEntries := range []int64{1, 2, 5, 9, 10, 25} {
		t.Run(fmt.Sprintf("cap_%d", maxEntries), func(t *testing.T) {
			c := NewWithSize(time.Minute, DefaultMaxMemorySize, maxEntries)
			defer c.Shutdown()
			for i := 0; i < 100; i++ {
				c.Set(fmt.Sprintf("k%d", i), []byte("v"), time.Minute)
				assert.LessOrEqual(t, c.CountQueries(), maxEntries)
			}
		})
	}
}
