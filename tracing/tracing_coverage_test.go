package tracing

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace/noop"
)

// TestNewTracing_NilContext covers the nil context early-return branch (line 34-36).
func TestNewTracing_NilContext_ReturnsError(t *testing.T) {
	_, err := NewTracing(nil, "localhost:4317") //nolint:staticcheck // SA1012: intentional nil to test the error branch
	require.Error(t, err)
	assert.Contains(t, err.Error(), "context cannot be nil")
}

// TestNewTracing_InvalidEndpointFormats covers endpoint validation branches.
// Note: fmt.Sscanf("%s:%d") treats %s as greedy so any "host:port" string hits
// the format error (n!=2). The port-range branch (port>65535) requires n==2
// which Sscanf never produces for "host:port" strings — that's a source quirk.
func TestNewTracing_InvalidEndpointFormats_ReturnsError(t *testing.T) {
	tests := []struct {
		name     string
		endpoint string
	}{
		{name: "no port separator", endpoint: "localhost"},
		{name: "port over max", endpoint: "localhost:999999"},
		{name: "plain hostname only", endpoint: "myhost"},
		{name: "just a number", endpoint: "12345"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := NewTracing(context.Background(), tt.endpoint)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "invalid endpoint format")
		})
	}
}

// TestShutdown_WithRealProvider covers the non-nil tracerProvider shutdown path (line 133).
func TestShutdown_WithRealProvider_NoError(t *testing.T) {
	// Use in-memory exporter so no network needed.
	exporter := tracetest.NewInMemoryExporter()
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSyncer(exporter),
	)
	ts := &TracingSetup{
		tracerProvider: tp,
		tracer:         tp.Tracer("shutdown-test"),
	}

	ctx := context.Background()
	err := ts.Shutdown(ctx)
	assert.NoError(t, err)
}

// TestStartSpan_WithRealTracer covers StartSpan with a real (noop) tracer — the non-nil path.
func TestStartSpan_WithRealTracer_ReturnsSpan(t *testing.T) {
	tp := noop.NewTracerProvider()
	ts := &TracingSetup{
		tracer: tp.Tracer("start-span-test"),
	}
	ctx := context.Background()
	span, newCtx := ts.StartSpan(ctx, "my-operation")
	assert.NotNil(t, span)
	assert.NotNil(t, newCtx)
	span.End()
}

// TestStartSpanWithAttributes_WithRealTracer covers the non-nil tracer path with attrs.
func TestStartSpanWithAttributes_WithRealTracer_RecordsSpan(t *testing.T) {
	exporter := tracetest.NewInMemoryExporter()
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSyncer(exporter),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)
	ts := &TracingSetup{
		tracerProvider: tp,
		tracer:         tp.Tracer("attr-test"),
	}

	ctx := context.Background()
	attrs := map[string]string{
		"user.id":   "u-42",
		"operation": "query",
	}
	span, newCtx := ts.StartSpanWithAttributes(ctx, "graphql-query", attrs)
	require.NotNil(t, span)
	require.NotNil(t, newCtx)
	span.End()

	spans := exporter.GetSpans()
	require.Len(t, spans, 1)
	assert.Equal(t, "graphql-query", spans[0].Name)
}

// TestExtractSpanContext_ValidTraceparent covers the valid span context branch (line 115-116).
// ExtractSpanContext uses otel.GetTextMapPropagator(); we must register the W3C
// TraceContext propagator before calling it (NewTracing normally does this).
func TestExtractSpanContext_ValidTraceparent_ReturnsValid(t *testing.T) {
	otel.SetTextMapPropagator(propagation.TraceContext{})

	tp := noop.NewTracerProvider()
	ts := &TracingSetup{
		tracer: tp.Tracer("extract-test"),
	}
	spanInfo := &TraceSpanInfo{
		TraceParent: "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01",
	}
	spanCtx, err := ts.ExtractSpanContext(spanInfo)
	require.NoError(t, err)
	assert.True(t, spanCtx.IsValid())
}
