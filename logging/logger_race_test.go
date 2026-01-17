package libpack_logger

import (
	"bytes"
	"sync"
	"testing"
)

// Test_LogConcurrentAccess verifies that the logger correctly handles concurrent access
// without race conditions
func TestLogConcurrentAccess(t *testing.T) {
	output := &bytes.Buffer{}
	logger := New().SetOutput(output).SetMinLogLevel(LEVEL_DEBUG)

	// Number of concurrent goroutines
	numGoroutines := 100
	// Wait group to synchronize goroutines
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Launch multiple goroutines to log concurrently
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			msg := &LogMessage{
				Message: "concurrent log test",
				Pairs: map[string]any{
					"goroutine_id": id,
				},
			}
			// Use different log levels to test all paths
			switch id % 5 {
			case 0:
				logger.Debug(msg)
			case 1:
				logger.Info(msg)
			case 2:
				logger.Warn(msg)
			case 3:
				logger.Error(msg)
			case 4:
				logger.Fatal(msg)
			}
		}(i)
	}

	// Wait for all goroutines to complete
	wg.Wait()

	// If we make it here without a race detector failure, the test passes
	if output.Len() == 0 {
		t.Error("Expected log output, but got none")
	}
}
