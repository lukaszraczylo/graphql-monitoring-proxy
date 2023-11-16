package main

import (
	libpack_monitoring "github.com/lukaszraczylo/graphql-monitoring-proxy/monitoring"
)

func StartMonitoringServer() {
	cfg.Monitoring = libpack_monitoring.NewMonitoring(cfg.Server.PurgeOnCrawl)
	cfg.Monitoring.AddMetricsPrefix("graphql_proxy")
	cfg.Monitoring.RegisterDefaultMetrics()
}
