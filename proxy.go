package main

import (
	"crypto/tls"
	"fmt"
	"time"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	libpack_monitoring "github.com/lukaszraczylo/graphql-monitoring-proxy/monitoring"
	"github.com/valyala/fasthttp"
)

func createFasthttpClient() *fasthttp.Client {
	return &fasthttp.Client{
		Name:                     "graphql_proxy",
		NoDefaultUserAgentHeader: true,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		MaxConnsPerHost:               100,
		MaxIdleConnDuration:           2 * time.Minute,
		ReadTimeout:                   time.Second * 10,
		WriteTimeout:                  time.Second * 10,
		DisableHeaderNamesNormalizing: true,
	}
}

func proxyTheRequest(c *fiber.Ctx) error {
	if !checkAllowedURLs(c) {
		cfg.Logger.Error("Request blocked", map[string]interface{}{"path": c.Path()})
		cfg.Monitoring.Increment(libpack_monitoring.MetricsSkipped, nil)
		c.Status(403).SendString("Request blocked - not allowed URL")
		return nil
	}

	c.Request().Header.Add("X-Real-IP", c.IP())
	c.Request().Header.Add(fiber.HeaderXForwardedFor, string(c.Request().Header.Peek("X-Forwarded-For")))

	proxy.WithClient(cfg.Client.FastProxyClient)

	cfg.Logger.Debug("Proxying the request", map[string]interface{}{"path": c.Path(), "body": string(c.Request().Body())})
	err := proxy.DoRedirects(c, cfg.Server.HostGraphQL+c.Path(), 3)
	if err != nil {
		cfg.Logger.Error("Can't proxy the request", map[string]interface{}{"error": err.Error()})
		cfg.Monitoring.Increment(libpack_monitoring.MetricsFailed, nil)
		return err
	}
	cfg.Logger.Debug("Received proxied response", map[string]interface{}{"path": c.Path(), "response_body": string(c.Response().Body()), "response_code": c.Response().StatusCode()})

	if c.Response().StatusCode() != 200 {
		cfg.Monitoring.Increment(libpack_monitoring.MetricsFailed, nil)
		return fmt.Errorf("Received non-200 response from the GraphQL server: %d", c.Response().StatusCode())
	}

	c.Response().Header.Del(fiber.HeaderServer)
	return nil
}
