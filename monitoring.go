package main

import (
	libpack_monitoring "github.com/telegram-bot-app/libpack/monitoring"
)

func StartMonitoringServer() {
	cfg.Monitoring = libpack_monitoring.NewMonitoring()
	cfg.Monitoring.AddMetricsPrefix("graphql_proxy")
	cfg.Monitoring.RegisterDefaultMetrics()
}
