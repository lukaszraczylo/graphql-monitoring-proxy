package libpack_trace

import (
	"context"
	"fmt"
	"time"

	libpack_config "github.com/lukaszraczylo/graphql-monitoring-proxy/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"

	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func NewClient(otelGRPCCollector string, attr ...attribute.KeyValue) (func(), error) {
	attr = append(attr, semconv.ServiceNameKey.String(libpack_config.PKG_NAME))
	fmt.Printf("Starting OpenTelemetry tracer: otlp, configured with endpoint: %s\n", otelGRPCCollector)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(otelGRPCCollector),
	)

	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		fmt.Printf("ERROR: failed to create exporter: %v\n", err)
		return nil, err
	}

	tp := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithBatcher(exporter, trace.WithMaxExportBatchSize(1), trace.WithBatchTimeout(30*time.Second)),
		trace.WithResource(resource.NewWithAttributes(semconv.SchemaURL, attr...)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	shutdownFunc := func() {
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()
		fmt.Println("Shutting down otlp tracer")
		if err := tp.Shutdown(shutdownCtx); err != nil {
			fmt.Printf("ERROR: failed to shutdown tracer provider: %v\n", err)
		}
	}

	return shutdownFunc, nil
}

func TraceContextInject(ctx context.Context) map[string]string {
	carrier := propagation.MapCarrier{}
	propagator := otel.GetTextMapPropagator()
	propagator.Inject(ctx, carrier)
	return map[string]string(carrier)
}

func TraceContextExtract(ctx context.Context, traceContext map[string]string) context.Context {
	carrier := propagation.MapCarrier(traceContext)
	propagator := otel.GetTextMapPropagator()
	return propagator.Extract(ctx, carrier)
}

func StartSpanFromContext(ctx context.Context, operationName string) (context.Context, oteltrace.Span) {
	tr := otel.GetTracerProvider().Tracer("")
	return tr.Start(ctx, operationName, oteltrace.WithSpanKind(oteltrace.SpanKindServer))
}

func ContinueSpanFromContext(ctx context.Context, operationName string) (context.Context, oteltrace.Span) {
	tr := otel.GetTracerProvider().Tracer("")
	options := []oteltrace.SpanStartOption{
		oteltrace.WithSpanKind(oteltrace.SpanKindInternal),
		oteltrace.WithAttributes(attribute.String("cont", "true")),
	}
	return tr.Start(ctx, operationName, options...)
}

func AddAttributesToSpan(span oteltrace.Span, attributes ...attribute.KeyValue) {
	span.SetAttributes(attributes...)
}
