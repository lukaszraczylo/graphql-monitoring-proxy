package libpack_monitoring

func (ms *MetricsSetup) RegisterDefaultMetrics() {
	ms.RegisterMetricsCounter(MetricsSucceeded, nil)
	ms.RegisterMetricsCounter(MetricsFailed, nil)
	ms.RegisterMetricsCounter(MetricsSkipped, nil)
	ms.RegisterMetricsHistogram(MetricsDuration, nil)
}

func (ms *MetricsSetup) RegisterGoMetrics() {
	// TODO: metrics.WriteProcessMetrics(ms.metrics_set)
}
