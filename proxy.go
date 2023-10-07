package main

import (
	"crypto/tls"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	libpack_monitoring "github.com/telegram-bot-app/libpack/monitoring"
)

func proxyTheRequest(c *fiber.Ctx) error {
	c.Request().Header.Add("X-Real-IP", c.IP())
	c.Request().Header.Add("X-Forwarded-For", c.IP())

	proxy.WithTlsConfig(&tls.Config{
		InsecureSkipVerify: true,
	})

	err := proxy.DoRedirects(c, cfg.Server.HostGraphQL, 3)
	if err != nil {
		cfg.Logger.Error("Can't proxy the request", map[string]interface{}{"error": err.Error()})
		cfg.Monitoring.Increment(libpack_monitoring.MetricsFailed, nil)
		return err
	}

	c.Response().Header.Del(fiber.HeaderServer)
	return nil
}
