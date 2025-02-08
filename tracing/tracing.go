package tracing

import (
	"context"
	"encoding/json"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

type TracingSetup struct {
	tracerProvider *sdktrace.TracerProvider
	tracer         trace.Tracer
}

type TraceSpanInfo struct {
	TraceParent string `json:"traceparent"`
}

// NewTracing creates a new tracing setup with OTLP exporter
func NewTracing(ctx context.Context, endpoint string) (*TracingSetup, error) {
	if ctx.Err() != nil {
		return nil, fmt.Errorf("invalid context: %v", ctx.Err())
	}
	if endpoint == "" {
		return nil, fmt.Errorf("endpoint cannot be empty")
	}

	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("graphql-monitoring-proxy"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	tracer := tracerProvider.Tracer("graphql-monitoring-proxy")

	return &TracingSetup{
		tracerProvider: tracerProvider,
		tracer:         tracer,
	}, nil
}

// ExtractSpanContext extracts span context from TraceSpanInfo
func (ts *TracingSetup) ExtractSpanContext(spanInfo *TraceSpanInfo) (trace.SpanContext, error) {
	carrier := propagation.MapCarrier{
		"traceparent": spanInfo.TraceParent,
	}
	ctx := context.Background()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)
	spanCtx := trace.SpanContextFromContext(ctx)
	if !spanCtx.IsValid() {
		return trace.SpanContext{}, fmt.Errorf("invalid span context")
	}
	return spanCtx, nil
}

// ParseTraceHeader parses X-Trace-Span header content
func ParseTraceHeader(headerContent string) (*TraceSpanInfo, error) {
	var spanInfo TraceSpanInfo
	if err := json.Unmarshal([]byte(headerContent), &spanInfo); err != nil {
		return nil, fmt.Errorf("failed to parse trace header: %w", err)
	}
	return &spanInfo, nil
}

// Shutdown cleanly shuts down the tracer provider
func (ts *TracingSetup) Shutdown(ctx context.Context) error {
	if ts.tracerProvider == nil {
		return nil
	}
	return ts.tracerProvider.Shutdown(ctx)
}

// StartSpan starts a new span with the given name and parent context
func (ts *TracingSetup) StartSpan(ctx context.Context, name string) (trace.Span, context.Context) {
	if ts.tracer == nil {
		return trace.SpanFromContext(ctx), ctx
	}
	ctx, span := ts.tracer.Start(ctx, name)
	return span, ctx
}
