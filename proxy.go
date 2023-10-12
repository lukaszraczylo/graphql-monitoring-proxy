package main

import (
	"crypto/tls"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	libpack_monitoring "github.com/lukaszraczylo/graphql-monitoring-proxy/monitoring"
)

func proxyTheRequest(c *fiber.Ctx) error {
	if !checkAllowedURLs(c) {
		cfg.Logger.Error("Request blocked", map[string]interface{}{"path": c.Path()})
		cfg.Monitoring.Increment(libpack_monitoring.MetricsSkipped, nil)
		c.Status(403).SendString("Request blocked - not allowed URL")
		return nil
	}

	c.Request().Header.Add("X-Real-IP", c.IP())
	c.Request().Header.Add(fiber.HeaderXForwardedFor, string(c.Request().Header.Peek("X-Forwarded-For")))

	proxy.WithTlsConfig(&tls.Config{
		InsecureSkipVerify: true,
	})

	err := proxy.DoRedirects(c, cfg.Server.HostGraphQL+c.Path(), 3)
	if err != nil {
		cfg.Logger.Error("Can't proxy the request", map[string]interface{}{"error": err.Error()})
		cfg.Monitoring.Increment(libpack_monitoring.MetricsFailed, nil)
		return err
	}

	c.Response().Header.Del(fiber.HeaderServer)
	return nil
}
