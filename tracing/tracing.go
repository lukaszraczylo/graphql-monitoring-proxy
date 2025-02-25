package tracing

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	if ctx == nil {
		return nil, fmt.Errorf("context cannot be nil")
	}
	if endpoint == "" {
		return nil, fmt.Errorf("endpoint cannot be empty")
	}

	// Create a timeout context for connection establishment
	dialCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Connect to the collector with improved options
	conn, err := grpc.DialContext(dialCtx, endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithReturnConnectionError(),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(16*1024*1024)), // 16MB max message size
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	// Create the exporter
	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithGRPCConn(conn),
		otlptracegrpc.WithTimeout(5*time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	// Create a resource with more detailed attributes
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("graphql-monitoring-proxy"),
			semconv.ServiceVersion("1.0"),
			semconv.DeploymentEnvironment("production"),
			attribute.String("application.type", "proxy"),
		),
		resource.WithHost(),      // Add host information
		resource.WithOSType(),    // Add OS information
		resource.WithProcessPID(), // Add process information
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Create the tracer provider with improved configuration
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter,
			// Configure batch processing
			sdktrace.WithMaxExportBatchSize(512),
			sdktrace.WithBatchTimeout(3*time.Second),
			sdktrace.WithMaxQueueSize(2048),
		),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(0.1)), // Sample 10% of traces
	)
	
	// Set the global tracer provider and propagator
	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	// Create a tracer
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
	if ts == nil || ts.tracer == nil {
		// Return a no-op span if tracing is not configured
		return trace.SpanFromContext(ctx), ctx
	}
	
	// Add common attributes to all spans
	opts := []trace.SpanStartOption{
		trace.WithAttributes(
			semconv.ServiceName("graphql-monitoring-proxy"),
			semconv.ServiceVersion("1.0"),
		),
	}
	
	ctx, span := ts.tracer.Start(ctx, name, opts...)
	return span, ctx
}

// StartSpanWithAttributes starts a new span with custom attributes
func (ts *TracingSetup) StartSpanWithAttributes(ctx context.Context, name string, attrs map[string]string) (trace.Span, context.Context) {
	if ts == nil || ts.tracer == nil {
		return trace.SpanFromContext(ctx), ctx
	}
	
	// Convert string attributes to KeyValue pairs
	attributes := make([]attribute.KeyValue, 0, len(attrs)+2)
	attributes = append(attributes,
		semconv.ServiceName("graphql-monitoring-proxy"),
		semconv.ServiceVersion("1.0"),
	)
	
	for k, v := range attrs {
		attributes = append(attributes, attribute.String(k, v))
	}
	
	ctx, span := ts.tracer.Start(ctx, name, trace.WithAttributes(attributes...))
	return span, ctx
}
