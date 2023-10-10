package libpack_logging

import (
	"os"
	"testing"
)

func BenchmarkNewLogger(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewLogger()
	}
}

func BenchmarkInfoLog(b *testing.B) {
	oldEnv := os.Getenv("LOG_LEVEL")
	os.Setenv("LOG_LEVEL", "info")
	oldStdout := os.Stdout
	oldStderr := os.Stderr
	os.Stdout, _ = os.Open(os.DevNull)
	os.Stderr, _ = os.Open(os.DevNull)
	defer func() {
		os.Stdout = oldStdout
		os.Stderr = oldStderr
		os.Setenv("LOG_LEVEL", oldEnv)
	}()

	testsLogger := NewLogger()
	for i := 0; i < b.N; i++ {
		testsLogger.Info("test", map[string]interface{}{"test": "test"})
	}
}
