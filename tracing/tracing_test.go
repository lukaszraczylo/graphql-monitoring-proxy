package libpack_trace

import (
	"testing"

	libpack_logging "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TraceTestSuite struct {
	suite.Suite
	logger *libpack_logging.Logger
}

func (suite *TraceTestSuite) SetupTest() {
	suite.logger = libpack_logging.New()
}

func (suite *TraceTestSuite) TearDownTest() {
	// Any cleanup logic can be added here
}

func TestTraceTestSuite(t *testing.T) {
	suite.Run(t, new(TraceTestSuite))
}

func (suite *TraceTestSuite) Test_NewClient() {
	shutdownFunc, err := NewClient(suite.logger, "localhost:4317")
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), shutdownFunc)

	shutdownFunc()
}

// func (suite *TraceTestSuite) Test_TraceContextInjectExtract() {
// 	ctx := context.Background()
// 	traceContext := TraceContextInject(ctx)
// 	assert.NotEmpty(suite.T(), traceContext)

// 	extractedCtx := TraceContextExtract(ctx, traceContext)
// 	assert.NotNil(suite.T(), extractedCtx)
// }

// func (suite *TraceTestSuite) Test_StartSpanFromContext() {
// 	ctx := context.Background()
// 	ctx, span := StartSpanFromContext(ctx, "operation")
// 	assert.NotNil(suite.T(), ctx)
// 	assert.NotNil(suite.T(), span)
// 	span.End()
// }

// func (suite *TraceTestSuite) Test_ContinueSpanFromContext() {
// 	ctx := context.Background()
// 	ctx, span := ContinueSpanFromContext(ctx, "operation")
// 	assert.NotNil(suite.T(), ctx)
// 	assert.NotNil(suite.T(), span)
// 	span.End()
// }

// func (suite *TraceTestSuite) Test_AddAttributesToSpan() {
// 	ctx := context.Background()
// 	_, span := StartSpanFromContext(ctx, "operation")

// 	attributes := []attribute.KeyValue{
// 		attribute.String("key1", "value1"),
// 		attribute.Int("key2", 2),
// 	}
// 	AddAttributesToSpan(span, attributes...)
// 	span.End()

// 	// Create an in-memory span exporter
// 	exporter := tracetest.NewSpanRecorder()
// 	tracerProvider := trace.NewTracerProvider(trace.WithSpanProcessor(exporter))
// 	otel.SetTracerProvider(tracerProvider)

// 	// Verify the span attributes
// 	spans := exporter.Ended()
// 	assert.Len(suite.T(), spans, 1)
// 	exportedSpan := spans[0]

// 	for _, attr := range attributes {
// 		assert.Contains(suite.T(), exportedSpan.Attributes(), attr)
// 	}
// }

// func (suite *TraceTestSuite) Test_Shutdown() {
// 	shutdownFunc, err := NewClient(suite.logger, "localhost:4317")
// 	assert.NoError(suite.T(), err)
// 	assert.NotNil(suite.T(), shutdownFunc)

// 	shutdownFunc()
// 	logOutput := captureStdOut(func() { suite.logger.Info(&libpack_logging.LogMessage{Message: "Shutting down tracer"}) })
// 	assert.Contains(suite.T(), logOutput, "Shutting down tracer")
// }

// // Helper function to capture standard output for testing logs
// func captureStdOut(f func()) string {
// 	originalStdout := os.Stdout
// 	r, w, _ := os.Pipe()
// 	os.Stdout = w
// 	f()
// 	w.Close()
// 	var buf bytes.Buffer
// 	buf.ReadFrom(r)
// 	os.Stdout = originalStdout
// 	return buf.String()
// }
