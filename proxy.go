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
	if cfg.Tracing.Enable && tracer != nil {
		var span trace.Span
		spanCtx := context.Background()
		// Extract trace information from header
		if traceHeader := c.Get("X-Trace-Span"); traceHeader != "" {
			spanInfo, err := libpack_tracing.ParseTraceHeader(traceHeader)
			if err != nil {
				cfg.Logger.Warning(&libpack_logger.LogMessage{
					Message: "Failed to parse trace header",
					Pairs:   map[string]interface{}{"error": err.Error()},
				})
			} else {
				if extractedSpanCtx, err := tracer.ExtractSpanContext(spanInfo); err == nil {
					spanCtx = trace.ContextWithSpanContext(spanCtx, extractedSpanCtx)
				}
			}
		}

		// Start a new span
		span, _ = tracer.StartSpan(spanCtx, "proxy_request")
		defer span.End()
	}

	if !checkAllowedURLs(c) {
		cfg.Logger.Error(&libpack_logger.LogMessage{
			Message: "Request blocked",
			Pairs:   map[string]interface{}{"path": c.Path()},
		})
		if ifNotInTest() {
			cfg.Monitoring.Increment(libpack_monitoring.MetricsSkipped, nil)
		}
		return fmt.Errorf("request blocked - not allowed URL: %s", c.Path())
	}

	proxyURL := currentEndpoint + c.Path()
	_, err := url.Parse(proxyURL)
	if err != nil {
		return fmt.Errorf("invalid URL: %v", err)
	}

	if cfg.LogLevel == "DEBUG" {
		logDebugRequest(c)
	}

	err = retry.Do(
		func() error {
			proxyErr := proxy.DoRedirects(c, proxyURL, 3, cfg.Client.FastProxyClient)
			if proxyErr != nil {
				return proxyErr
			}
			if c.Response().StatusCode() != fiber.StatusOK {
				return fmt.Errorf("received non-200 response from the GraphQL server: %d", c.Response().StatusCode())
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

	if err != nil {
		cfg.Logger.Warning(&libpack_logger.LogMessage{
			Message: "Can't proxy the request",
			Pairs:   map[string]interface{}{"error": err.Error()},
		})
		if ifNotInTest() {
			cfg.Monitoring.Increment(libpack_monitoring.MetricsFailed, nil)
		}
		return fmt.Errorf("failed to proxy request: %v", err)
	}

	if cfg.LogLevel == "DEBUG" {
		logDebugResponse(c)
	}

	if bytes.EqualFold(c.Response().Header.Peek("Content-Encoding"), []byte("gzip")) {
		// Decompress gzip response
		reader, err := gzip.NewReader(bytes.NewReader(c.Response().Body()))
		if err != nil {
			cfg.Logger.Error(&libpack_logger.LogMessage{
				Message: "Failed to create gzip reader",
				Pairs:   map[string]interface{}{"error": err.Error()},
			})
			return err
		}
		defer reader.Close()

		decompressed, err := io.ReadAll(reader)
		if err != nil {
			cfg.Logger.Error(&libpack_logger.LogMessage{
				Message: "Failed to decompress response",
				Pairs:   map[string]interface{}{"error": err.Error()},
			})
			return err
		}

		c.Response().SetBody(decompressed)
		c.Response().Header.Del("Content-Encoding")
	}

	if c.Response().StatusCode() != fiber.StatusOK {
		if ifNotInTest() {
			cfg.Monitoring.Increment(libpack_monitoring.MetricsFailed, nil)
		}
		return fmt.Errorf("received non-200 response from the GraphQL server: %d", c.Response().StatusCode())
	}

	c.Response().Header.Del(fiber.HeaderServer)
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
