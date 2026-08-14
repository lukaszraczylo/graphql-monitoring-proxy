package libpack_logger

import (
	"bytes"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/goccy/go-json"
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

// TestConcurrentSettersDuringLog verifies that SetFieldName, SetMinLogLevel,
// SetShowCaller, and SetTimeFormat can run concurrently with active log()
// calls without a data race.
//
// Before the fix, SetFieldName wrote the package-global fieldNames map with
// no lock while log() read it from many goroutines, and
// SetMinLogLevel/SetShowCaller/SetTimeFormat mutated Logger fields with no
// lock while shouldLog()/log() read them. This test fails under
// `go test -race` on that code and passes once the setters and readers are
// synchronized.
func TestConcurrentSettersDuringLog(t *testing.T) {
	// fieldNames is package-global, so snapshot and restore it: this test
	// mutates it concurrently and must not leak state into other tests.
	originalFieldNames := make(map[string]string)
	for k, v := range fieldNames {
		originalFieldNames[k] = v
	}
	defer func() {
		for k, v := range originalFieldNames {
			fieldNames[k] = v
		}
	}()

	output := &bytes.Buffer{}
	logger := New().SetOutput(output).SetMinLogLevel(LEVEL_DEBUG)

	const (
		numLoggers = 50
		numSetters = 50
		iterations = 200
	)

	var wg sync.WaitGroup
	wg.Add(numLoggers + numSetters)

	// Goroutines that continuously log, exercising log()'s reads of the
	// global fieldNames map and the Logger's minLogLevel/showCaller/
	// timeFormat fields.
	for i := 0; i < numLoggers; i++ {
		go func(id int) {
			defer wg.Done()
			msg := &LogMessage{Message: "concurrent setter test"}
			for j := 0; j < iterations; j++ {
				logger.Info(msg)
				logger.Debug(msg)
			}
		}(i)
	}

	// Goroutines that continuously mutate logger-wide configuration while
	// logging is in flight from the goroutines above.
	for i := 0; i < numSetters; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				logger.SetFieldName("timestamp", "ts")
				logger.SetMinLogLevel(LEVEL_DEBUG)
				logger.SetShowCaller(id%2 == 0)
				logger.SetTimeFormat(time.RFC3339Nano)
			}
		}(i)
	}

	wg.Wait()

	if output.Len() == 0 {
		t.Fatal("Expected log output, but got none")
	}

	// Each emitted line must still be well-formed JSON: the concurrent
	// setter calls must not corrupt log() output, only race on which
	// setting "wins" at each point in time.
	lines := strings.Split(strings.TrimSpace(output.String()), "\n")
	for i, line := range lines {
		var decoded map[string]any
		if err := json.Unmarshal([]byte(line), &decoded); err != nil {
			t.Fatalf("line %d is not valid JSON: %v\nline: %s", i, err, line)
		}
	}
}
