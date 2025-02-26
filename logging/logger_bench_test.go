package libpack_logger

import (
	"bytes"
	"testing"
	"time"
)

func Benchmark_NewLogger(b *testing.B) {
	type triggers struct {
		ModFormat struct {
			Format string
		}
		ModLevel struct {
			Level int
		}
	}

	tests := []struct {
		name     string
		triggers triggers
	}{
		{
			name: "BenchmarkNew",
		},
		{
			name: "BenchmarkNewChangeTimeFormat",
			triggers: triggers{
				ModFormat: struct{ Format string }{
					Format: time.RFC3339Nano,
				},
			},
		},
		{
			name: "BenchmarkNewChangeLogLevel",
			triggers: triggers{
				ModLevel: struct{ Level int }{
					Level: LEVEL_DEBUG,
				},
			},
		},
		{
			name: "BenchmarkNewChangeTimeFormatAndLogLevel",
			triggers: triggers{
				ModFormat: struct{ Format string }{
					Format: time.RFC3339Nano,
				},
				ModLevel: struct{ Level int }{
					Level: LEVEL_DEBUG,
				},
			},
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = New()
			}
		})
	}
}

func Benchmark_Log_Debug(b *testing.B) {
	output := &bytes.Buffer{}
	logger := New().SetMinLogLevel(LEVEL_DEBUG).SetOutput(output)
	msg := &LogMessage{
		Message: "debug message",
		Pairs:   make(map[string]any),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Debug(msg)
	}
}

func Benchmark_Log_Info(b *testing.B) {
	output := &bytes.Buffer{}
	logger := New().SetMinLogLevel(LEVEL_INFO).SetOutput(output)
	msg := &LogMessage{
		Message: "info message",
		Pairs:   make(map[string]any),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info(msg)
	}
}

func Benchmark_Log_Warn(b *testing.B) {
	output := &bytes.Buffer{}
	logger := New().SetMinLogLevel(LEVEL_WARN).SetOutput(output)
	msg := &LogMessage{
		Message: "warn message",
		Pairs:   make(map[string]any),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Warn(msg)
	}
}

func Benchmark_Log_Error(b *testing.B) {
	output := &bytes.Buffer{}
	logger := New().SetMinLogLevel(LEVEL_ERROR).SetOutput(output)
	msg := &LogMessage{
		Message: "error message",
		Pairs:   map[string]any{"key": "value"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Error(msg)
	}
}

func Benchmark_Log_Fatal(b *testing.B) {
	output := &bytes.Buffer{}
	logger := New().SetMinLogLevel(LEVEL_FATAL).SetOutput(output)
	msg := &LogMessage{
		Message: "fatal message",
		Pairs:   make(map[string]any),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Fatal(msg)
	}
}
