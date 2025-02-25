package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/url"
	"time"

	"go.opentelemetry.io/otel/trace"

	"github.com/avast/retry-go/v4"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	libpack_logger "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
	libpack_monitoring "github.com/lukaszraczylo/graphql-monitoring-proxy/monitoring"
	libpack_tracing "github.com/lukaszraczylo/graphql-monitoring-proxy/tracing"
	"github.com/valyala/fasthttp"
)

// createFasthttpClient creates and configures a fasthttp client.
func createFasthttpClient(timeout int) *fasthttp.Client {
	return &fasthttp.Client{
		Name:                     "graphql_proxy",
		NoDefaultUserAgentHeader: true,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		MaxConnsPerHost:               2048,
		ReadTimeout:                   time.Duration(timeout) * time.Second,
		WriteTimeout:                  time.Duration(timeout) * time.Second,
		MaxIdleConnDuration:           time.Duration(timeout) * time.Second,
		MaxConnDuration:               time.Duration(timeout) * time.Second,
		DisableHeaderNamesNormalizing: false,
	}
}

// proxyTheRequest handles the request proxying logic.
func proxyTheRequest(c *fiber.Ctx, currentEndpoint string) error {
	// Setup tracing if enabled
	var span trace.Span
	ctx := setupTracing(c)
	
	if cfg.Tracing.Enable && tracer != nil {
		span, ctx = tracer.StartSpan(ctx, "proxy_request")
		defer span.End()
	}

	// Check if URL is allowed
	if !checkAllowedURLs(c) {
		if ifNotInTest() {
			cfg.Monitoring.Increment(libpack_monitoring.MetricsSkipped, nil)
		}
		return fmt.Errorf("request blocked - not allowed URL: %s", c.Path())
	}

	// Construct and validate proxy URL
	proxyURL := currentEndpoint + c.Path()
	if _, err := url.Parse(proxyURL); err != nil {
		return fmt.Errorf("invalid URL: %v", err)
	}

	// Log request details in debug mode
	if cfg.LogLevel == "DEBUG" {
		logDebugRequest(c)
	}

	// Perform the proxy request with retries
	if err := performProxyRequest(c, proxyURL); err != nil {
		if ifNotInTest() {
			cfg.Monitoring.Increment(libpack_monitoring.MetricsFailed, nil)
		}
		return err
	}

	// Log response details in debug mode
	if cfg.LogLevel == "DEBUG" {
		logDebugResponse(c)
	}

	// Handle gzipped responses
	if err := handleGzippedResponse(c); err != nil {
		return err
	}

	// Final status check
	if c.Response().StatusCode() != fiber.StatusOK {
		if ifNotInTest() {
			cfg.Monitoring.Increment(libpack_monitoring.MetricsFailed, nil)
		}
		return fmt.Errorf("received non-200 response from the GraphQL server: %d", c.Response().StatusCode())
	}

	// Remove server header for security
	c.Response().Header.Del(fiber.HeaderServer)
	return nil
}

// setupTracing extracts and sets up tracing context from request headers
func setupTracing(c *fiber.Ctx) context.Context {
	ctx := context.Background()
	
	if !cfg.Tracing.Enable || tracer == nil {
		return ctx
	}
	
	// Extract trace information from header
	if traceHeader := c.Get("X-Trace-Span"); traceHeader != "" {
		spanInfo, err := libpack_tracing.ParseTraceHeader(traceHeader)
		if err != nil {
			cfg.Logger.Warning(&libpack_logger.LogMessage{
				Message: "Failed to parse trace header",
				Pairs:   map[string]interface{}{"error": err.Error()},
			})
		} else if spanCtx, err := tracer.ExtractSpanContext(spanInfo); err == nil {
			ctx = trace.ContextWithSpanContext(ctx, spanCtx)
		}
	}
	
	return ctx
}

// performProxyRequest executes the proxy request with retries
func performProxyRequest(c *fiber.Ctx, proxyURL string) error {
	return retry.Do(
		func() error {
			if err := proxy.DoRedirects(c, proxyURL, 3, cfg.Client.FastProxyClient); err != nil {
				return err
			}
			if c.Response().StatusCode() != fiber.StatusOK {
				return fmt.Errorf("received non-200 response: %d", c.Response().StatusCode())
			}
			return nil
		},
		retry.Attempts(5),
		retry.DelayType(retry.BackOffDelay),
		retry.Delay(250*time.Millisecond),
		retry.MaxDelay(5*time.Second),
		retry.OnRetry(func(n uint, err error) {
			cfg.Logger.Warning(&libpack_logger.LogMessage{
				Message: "Retrying the request",
				Pairs: map[string]interface{}{
					"path":    c.Path(),
					"attempt": n + 1,
					"error":   err.Error(),
				},
			})
		}),
		retry.LastErrorOnly(true),
	)
}

// handleGzippedResponse decompresses gzipped responses
func handleGzippedResponse(c *fiber.Ctx) error {
	if !bytes.EqualFold(c.Response().Header.Peek("Content-Encoding"), []byte("gzip")) {
		return nil
	}
	
	// Create a pooled gzip reader
	reader, err := gzip.NewReader(bytes.NewReader(c.Response().Body()))
	if err != nil {
		cfg.Logger.Error(&libpack_logger.LogMessage{
			Message: "Failed to create gzip reader",
			Pairs:   map[string]interface{}{"error": err.Error()},
		})
		return err
	}
	defer reader.Close()

	// Read decompressed data
	decompressed, err := io.ReadAll(reader)
	if err != nil {
		cfg.Logger.Error(&libpack_logger.LogMessage{
			Message: "Failed to decompress response",
			Pairs:   map[string]interface{}{"error": err.Error()},
		})
		return err
	}

	// Update response
	c.Response().SetBody(decompressed)
	c.Response().Header.Del("Content-Encoding")
	return nil
}

// logDebugRequest logs the request details when in debug mode.
func logDebugRequest(c *fiber.Ctx) {
	cfg.Logger.Debug(&libpack_logger.LogMessage{
		Message: "Proxying the request",
		Pairs: map[string]interface{}{
			"path":         c.Path(),
			"body":         string(c.Body()),
			"headers":      c.GetReqHeaders(),
			"request_uuid": c.Locals("request_uuid"),
		},
	})
}

// logDebugResponse logs the response details when in debug mode.
func logDebugResponse(c *fiber.Ctx) {
	cfg.Logger.Debug(&libpack_logger.LogMessage{
		Message: "Received proxied response",
		Pairs: map[string]interface{}{
			"path":          c.Path(),
			"response_body": string(c.Response().Body()),
			"response_code": c.Response().StatusCode(),
			"headers":       c.GetRespHeaders(),
			"request_uuid":  c.Locals("request_uuid"),
		},
	})
}
