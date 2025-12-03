package libpack_monitoring

func (ms *MetricsSetup) RegisterDefaultMetrics() {
	ms.RegisterMetricsCounter(MetricsSucceeded, nil)
	ms.RegisterMetricsCounter(MetricsFailed, nil)
	ms.RegisterMetricsCounter(MetricsSkipped, nil)
	ms.RegisterMetricsHistogram(MetricsDuration, nil)
	ms.RegisterMetricsCounter(MetricsCacheHit, nil)
	ms.RegisterMetricsCounter(MetricsCacheMiss, nil)
	ms.RegisterMetricsCounter(MetricsQueriesCached, nil)
}
