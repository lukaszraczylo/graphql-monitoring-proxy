package tracing

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/trace"
)

func TestParseTraceHeader(t *testing.T) {
	tests := []struct {
		name    string
		header  string
		want    *TraceSpanInfo
		wantErr bool
	}{
		{
			name:   "valid trace header",
			header: `{"traceparent": "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01"}`,
			want: &TraceSpanInfo{
				TraceParent: "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01",
			},
			wantErr: false,
		},
		{
			name:    "invalid json",
			header:  `{"traceparent": invalid}`,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "empty header",
			header:  "",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseTraceHeader(tt.header)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTraceHeader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				gotJSON, _ := json.Marshal(got)
				wantJSON, _ := json.Marshal(tt.want)
				if string(gotJSON) != string(wantJSON) {
					t.Errorf("ParseTraceHeader() = %v, want %v", string(gotJSON), string(wantJSON))
				}
			}
		})
	}
}

func TestNewTracing(t *testing.T) {
	// Skip actual connection tests since they require a running collector
	t.Run("empty endpoint", func(t *testing.T) {
		ctx := context.Background()
		_, err := NewTracing(ctx, "")
		assert.Error(t, err, "Expected error for empty endpoint")
		assert.Contains(t, err.Error(), "endpoint cannot be empty")
	})

	t.Run("invalid endpoint", func(t *testing.T) {
		// We'll use a more severe syntax error in the endpoint to trigger a validation error
		ctx := context.Background()
		// Use a port that exceeds the maximum valid port number
		_, err := NewTracing(ctx, "localhost:999999")
		assert.Error(t, err, "Expected error for invalid endpoint format")
	})
}

func TestTracingSetup_ExtractSpanContext(t *testing.T) {
	ts := &TracingSetup{}
	spanInfo := &TraceSpanInfo{
		TraceParent: "invalid-traceparent",
	}

	_, err := ts.ExtractSpanContext(spanInfo)
	assert.Error(t, err, "Expected error for invalid traceparent")
	assert.Contains(t, err.Error(), "invalid span context")
}

func TestTracingSetup_StartSpan(t *testing.T) {
	ts := &TracingSetup{}
	ctx := context.Background()

	span, newCtx := ts.StartSpan(ctx, "test-span")
	assert.NotNil(t, span, "Expected non-nil span even when tracer is nil")
	assert.NotNil(t, newCtx, "Expected non-nil context")
	assert.Equal(t, trace.SpanFromContext(ctx), span, "Expected span from context when tracer is nil")
}

func TestTracingSetup_Shutdown(t *testing.T) {
	ts := &TracingSetup{}
	ctx := context.Background()

	err := ts.Shutdown(ctx)
	assert.NoError(t, err, "Expected no error when shutting down nil tracer provider")
}
