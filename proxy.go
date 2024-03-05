package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"time"

	"github.com/avast/retry-go/v4"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	libpack_monitoring "github.com/lukaszraczylo/graphql-monitoring-proxy/monitoring"
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

func proxyTheRequest(c *fiber.Ctx) error {
	if !checkAllowedURLs(c) {
		cfg.Logger.Error("Request blocked", map[string]interface{}{"path": c.Path()})
		if flag.Lookup("test.v") == nil {
			cfg.Monitoring.Increment(libpack_monitoring.MetricsSkipped, nil)
		}
		c.Status(403).SendString("Request blocked - not allowed URL")
		return nil
	}
	c.Request().Header.DisableNormalizing()
	c.Request().Header.Add("X-Real-IP", c.IP())
	c.Request().Header.Add(fiber.HeaderXForwardedFor, string(c.Request().Header.Peek("X-Forwarded-For")))
	c.Request().Header.Del(fiber.HeaderAcceptEncoding)

	cfg.Logger.Debug("Proxying the request", map[string]interface{}{"path": c.Path(), "body": string(c.Request().Body()), "headers": c.GetReqHeaders(), "request_uuid": c.Locals("request_uuid")})

	err := retry.Do(
		func() error {
			errInt := proxy.DoRedirects(c, cfg.Server.HostGraphQL+c.Path(), 3, cfg.Client.FastProxyClient)
			if errInt != nil {
				cfg.Logger.Error("Can't proxy the request", map[string]interface{}{"error": errInt.Error()})
				if flag.Lookup("test.v") == nil {
					cfg.Monitoring.Increment(libpack_monitoring.MetricsFailed, nil)
				}
				return errInt
			}
			return nil
		},
		retry.OnRetry(func(n uint, err error) {
			cfg.Logger.Warning("Retrying the request", map[string]interface{}{"path": c.Path(), "error": err.Error()})
		}),
		retry.Attempts(uint(3)),
		retry.DelayType(retry.BackOffDelay),
		retry.Delay(time.Duration(250*time.Millisecond)),
		retry.LastErrorOnly(true),
	)

	if err != nil {
		cfg.Logger.Warning("Can't proxy the request", map[string]interface{}{"error": err.Error()})
		return err
	}

	cfg.Logger.Debug("Received proxied response", map[string]interface{}{"path": c.Path(), "response_body": string(c.Response().Body()), "response_code": c.Response().StatusCode(), "headers": c.GetRespHeaders(), "request_uuid": c.Locals("request_uuid")})

	if c.Response().StatusCode() != 200 {
		if flag.Lookup("test.v") == nil {
			cfg.Monitoring.Increment(libpack_monitoring.MetricsFailed, nil)
		}
		return fmt.Errorf("Received non-200 response from the GraphQL server: %d", c.Response().StatusCode())
	}

	c.Response().Header.Del(fiber.HeaderServer)
	return nil
}
