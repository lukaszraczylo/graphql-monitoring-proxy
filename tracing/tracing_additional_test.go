package tracing

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

func TestStartSpanWithAttributes(t *testing.T) {
	// Create a minimal tracing setup without actual connection
	ts := &TracingSetup{
		tracer: noop.NewTracerProvider().Tracer("test"),
	}

	// Test with attributes
	t.Run("with attributes", func(t *testing.T) {
		ctx := context.Background()
		attrs := map[string]string{
			"key1": "value1",
			"key2": "value2",
		}

		span, newCtx := ts.StartSpanWithAttributes(ctx, "test-span", attrs)
		assert.NotNil(t, span)
		assert.NotNil(t, newCtx)

		// We can't easily test the attributes were set since it's a noop tracer,
		// but we can verify the function doesn't panic
		span.End()
	})

	// Test with nil attributes
	t.Run("with nil attributes", func(t *testing.T) {
		ctx := context.Background()

		span, newCtx := ts.StartSpanWithAttributes(ctx, "test-span", nil)
		assert.NotNil(t, span)
		assert.NotNil(t, newCtx)
		span.End()
	})

	// Test with nil tracer
	t.Run("with nil tracer", func(t *testing.T) {
		ctx := context.Background()
		nilTS := &TracingSetup{tracer: nil}

		span, newCtx := nilTS.StartSpanWithAttributes(ctx, "test-span", map[string]string{"key": "value"})
		assert.NotNil(t, span)
		assert.NotNil(t, newCtx)
		// Should not panic when ending the span
		span.End()
	})
}

func TestNewTracingWithInvalidEndpoint(t *testing.T) {
	// Skip endpoint tests that are already covered in the main test file
	t.Run("invalid endpoint format", func(t *testing.T) {
		t.Skip("This test is now handled in the main test file")
	})

	// Skip the unreachable endpoint test as it's flaky and already tested
	t.Run("unreachable endpoint", func(t *testing.T) {
		t.Skip("This test is now handled in the main test file")
	})
}

func TestTracingSetupWithMockTracer(t *testing.T) {
	// Create a mock tracer provider
	mockTracerProvider := noop.NewTracerProvider()
	mockTracer := mockTracerProvider.Tracer("mock-tracer")

	ts := &TracingSetup{
		tracerProvider: nil, // We don't need the provider for these tests
		tracer:         mockTracer,
	}

	// Test StartSpan
	t.Run("start span", func(t *testing.T) {
		ctx := context.Background()
		span, newCtx := ts.StartSpan(ctx, "test-span")

		assert.NotNil(t, span)
		assert.NotNil(t, newCtx)

		// Add some attributes and events to ensure no panics
		span.SetAttributes(attribute.String("test", "value"))
		span.AddEvent("test-event")

		// End the span
		span.End()
	})

	// Test StartSpanWithAttributes
	t.Run("start span with attributes", func(t *testing.T) {
		ctx := context.Background()
		attrs := map[string]string{
			"service": "test-service",
			"version": "1.0.0",
		}

		span, newCtx := ts.StartSpanWithAttributes(ctx, "test-span-with-attrs", attrs)

		assert.NotNil(t, span)
		assert.NotNil(t, newCtx)

		// End the span
		span.End()
	})
}

func TestShutdownWithNilProvider(t *testing.T) {
	ts := &TracingSetup{
		tracerProvider: nil,
		tracer:         noop.NewTracerProvider().Tracer("test"),
	}

	ctx := context.Background()
	err := ts.Shutdown(ctx)

	assert.NoError(t, err)
}

func TestExtractSpanContextWithInvalidTraceParent(t *testing.T) {
	ts := &TracingSetup{
		tracer: noop.NewTracerProvider().Tracer("test"),
	}

	// Test with invalid traceparent format
	t.Run("invalid traceparent format", func(t *testing.T) {
		spanInfo := &TraceSpanInfo{
			TraceParent: "invalid-format",
		}

		// Explicitly type the result to use trace package
		var spanCtx trace.SpanContext
		var err error
		spanCtx, err = ts.ExtractSpanContext(spanInfo)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid span context")
		assert.False(t, spanCtx.IsValid())
	})
}

func TestParseTraceHeaderWithEmptyHeader(t *testing.T) {
	// Test with empty header
	t.Run("empty header", func(t *testing.T) {
		_, err := ParseTraceHeader("")
		assert.Error(t, err)
	})

	// Test with invalid JSON
	t.Run("invalid JSON", func(t *testing.T) {
		_, err := ParseTraceHeader("{invalid json}")
		assert.Error(t, err)
	})

	// Test with valid JSON but missing traceparent
	t.Run("missing traceparent", func(t *testing.T) {
		_, err := ParseTraceHeader(`{"other": "value"}`)
		assert.NoError(t, err) // This should parse but the traceparent will be empty
	})
}
