package main

import (
	"crypto/tls"
	"fmt"
	"net/url"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	libpack_logger "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
	libpack_monitoring "github.com/lukaszraczylo/graphql-monitoring-proxy/monitoring"
	"github.com/valyala/fasthttp"
)

var (
	httpClient *fasthttp.Client
)

func init() {
	httpClient = createFasthttpClient(30) // Assuming a default timeout of 30 seconds
}

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
		DisableHeaderNamesNormalizing: true,
	}
}
func proxyTheRequest(c *fiber.Ctx, currentEndpoint string) error {
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

	if cfg.LogLevel == "debug" {
		logDebugRequest(c)
	}

	err = retry.Do(
		func() error {
			return proxy.DoRedirects(c, proxyURL, 3, httpClient)
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
		retry.Attempts(3),
		retry.DelayType(retry.BackOffDelay),
		retry.Delay(250*time.Millisecond),
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

	if cfg.LogLevel == "debug" {
		logDebugResponse(c)
	}

	if c.Response().StatusCode() != 200 {
		if ifNotInTest() {
			cfg.Monitoring.Increment(libpack_monitoring.MetricsFailed, nil)
		}
		return fmt.Errorf("received non-200 response from the GraphQL server: %d", c.Response().StatusCode())
	}

	c.Response().Header.Del(fiber.HeaderServer)
	return nil
}

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
