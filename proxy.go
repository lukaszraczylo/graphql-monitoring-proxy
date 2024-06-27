package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/goccy/go-json"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	libpack_logger "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
	libpack_monitoring "github.com/lukaszraczylo/graphql-monitoring-proxy/monitoring"
	libpack_trace "github.com/lukaszraczylo/graphql-monitoring-proxy/tracing"
	"github.com/valyala/fasthttp"
)

func createFasthttpClient(timeout int) *fasthttp.Client {
	return &fasthttp.Client{
		Name:                     "graphql_proxy",
		NoDefaultUserAgentHeader: true,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		MaxConnsPerHost:               2048,
		ReadTimeout:                   time.Second * time.Duration(timeout),
		WriteTimeout:                  time.Second * time.Duration(timeout),
		MaxIdleConnDuration:           time.Second * time.Duration(timeout),
		MaxConnDuration:               time.Second * time.Duration(timeout),
		DisableHeaderNamesNormalizing: true,
	}
}

func proxyTheRequest(c *fiber.Ctx, currentEndpoint string, ctx context.Context) error {
	if !checkAllowedURLs(c) {
		cfg.Logger.Error(&libpack_logger.LogMessage{
			Message: "Request blocked",
			Pairs:   map[string]interface{}{"path": c.Path()},
		})
		if ifNotInTest() {
			cfg.Monitoring.Increment(libpack_monitoring.MetricsSkipped, nil)
		}
		c.Status(403).SendString("Request blocked - not allowed URL")
		return nil
	}
	c.Request().Header.DisableNormalizing()
	c.Request().Header.Add("X-Real-IP", c.IP())
	c.Request().Header.Add(fiber.HeaderXForwardedFor, string(c.Request().Header.Peek("X-Forwarded-For")))
	c.Request().Header.Del(fiber.HeaderAcceptEncoding)

	// added dummy check for the log level because it executes additional functions which could
	// potentially slow down the execution.
	if cfg.LogLevel == "debug" {
		cfg.Logger.Debug(&libpack_logger.LogMessage{
			Message: "Proxying the request",
			Pairs: map[string]interface{}{
				"path":         c.Path(),
				"body":         string(c.Request().Body()),
				"headers":      c.GetReqHeaders(),
				"request_uuid": c.Locals("request_uuid"),
			},
		})
	}

	err := retry.Do(
		func() error {
			errInt := proxy.DoRedirects(c, currentEndpoint+c.Path(), 3, cfg.Client.FastProxyClient)
			if errInt != nil {
				cfg.Logger.Error(&libpack_logger.LogMessage{
					Message: "Can't proxy the request",
					Pairs: map[string]interface{}{
						"error": errInt.Error(),
					},
				})
				if ifNotInTest() {
					cfg.Monitoring.Increment(libpack_monitoring.MetricsFailed, nil)
				}
				return errInt
			}
			return nil
		},
		retry.OnRetry(func(n uint, err error) {
			cfg.Logger.Warning(&libpack_logger.LogMessage{
				Message: "Retrying the request",
				Pairs: map[string]interface{}{
					"path":  c.Path(),
					"error": err.Error(),
				},
			})
		}),
		retry.Attempts(uint(3)),
		retry.DelayType(retry.BackOffDelay),
		retry.Delay(time.Duration(250*time.Millisecond)),
		retry.LastErrorOnly(true),
	)

	if err != nil {
		cfg.Logger.Warning(&libpack_logger.LogMessage{
			Message: "Can't proxy the request",
			Pairs: map[string]interface{}{
				"error": err.Error(),
			},
		})
		return err
	}

	if cfg.LogLevel == "debug" {
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

	if c.Response().StatusCode() != 200 {
		if ifNotInTest() {
			cfg.Monitoring.Increment(libpack_monitoring.MetricsFailed, nil)
		}
		cfg.Logger.Error(&libpack_logger.LogMessage{
			Message: "Received non-200 response from the GraphQL server",
			Pairs: map[string]interface{}{
				"status_code": c.Response().StatusCode(),
			},
		})
		return fmt.Errorf("Received non-200 response from the GraphQL server: %d", c.Response().StatusCode())
	}

	c.Response().Header.Del(fiber.HeaderServer)
	if cfg.Trace.Enable {
		tracingContext := libpack_trace.TraceContextInject(ctx)
		if tracingContext == nil {
			cfg.Logger.Error(&libpack_logger.LogMessage{
				Message: "Can't inject empty tracing context",
			})
			return nil
		}
		traceJsonEncoded, err := json.Marshal(tracingContext)
		if err != nil {
			cfg.Logger.Error(&libpack_logger.LogMessage{
				Message: "Can't convert tracing context to JSON",
				Pairs: map[string]interface{}{
					"error": err.Error(),
				},
			})
			return err
		}
		c.Response().Header.Set("X-Trace-Span", string(traceJsonEncoded))
	}
	return nil
}
