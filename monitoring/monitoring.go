// Package `libpack_monitoring` provides and easy way to add prometheus metrics to your application.
// It also provides a way to add custom metrics to the already started prometheus registry.

package libpack_monitoring

import (
	"fmt"
	"time"

	"github.com/VictoriaMetrics/metrics"
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/goutil/envutil"
	libpack_config "github.com/lukaszraczylo/graphql-monitoring-proxy/config"
	logging "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
)

type MetricsSetup struct {
	metrics_prefix string
	metrics_set    *metrics.Set
}

var (
	log *logging.LogConfig
)

func NewMonitoring() *MetricsSetup {
	log = logging.NewLogger()
	ms := &MetricsSetup{}
	ms.metrics_set = metrics.NewSet()
	go ms.startPrometheusEndpoint()
	return ms
}

func (ms *MetricsSetup) startPrometheusEndpoint() {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		AppName:               fmt.Sprintf("GraphQL Monitoring Proxy - %s v%s", libpack_config.PKG_NAME, libpack_config.PKG_VERSION),
	})
	app.Get("/metrics", ms.metricsEndpoint)
	err := app.Listen(fmt.Sprintf(":%d", envutil.GetInt("MONITORING_PORT", 9393)))
	if err != nil {
		fmt.Println("Can't start the service: ", err)
	}
}

func (ms *MetricsSetup) metricsEndpoint(c *fiber.Ctx) error {
	ms.metrics_set.WritePrometheus(c.Response().BodyWriter())
	return nil
}

func (ms *MetricsSetup) AddMetricsPrefix(prefix string) {
	ms.metrics_prefix = prefix
}

func (ms *MetricsSetup) ListActiveMetrics() []string {
	return ms.metrics_set.ListMetricNames()
}

func (ms *MetricsSetup) RegisterMetricsGauge(metric_name string, labels map[string]string, val float64) *metrics.Gauge {
	if validate_metrics_name(metric_name) != nil {
		log.Critical("RegisterMetricsGauge() error", map[string]interface{}{"_error": "Invalid metric name", "_metric_name": metric_name})
		return nil
	}
	return ms.metrics_set.GetOrCreateGauge(ms.get_metrics_name(metric_name, labels), func() float64 {
		// get current value of the gauge and add val to it
		return val
	})
}

func (ms *MetricsSetup) RegisterMetricsCounter(metric_name string, labels map[string]string) *metrics.Counter {
	if validate_metrics_name(metric_name) != nil {
		log.Critical("RegisterMetricsCounter() error", map[string]interface{}{"_error": "Invalid metric name", "_metric_name": metric_name})
		return nil
	}
	return ms.metrics_set.GetOrCreateCounter(ms.get_metrics_name(metric_name, labels))
}

func (ms *MetricsSetup) RegisterFloatCounter(metric_name string, labels map[string]string) *metrics.FloatCounter {
	if validate_metrics_name(metric_name) != nil {
		log.Critical("RegisterFloatCounter() error", map[string]interface{}{"_error": "Invalid metric name", "_metric_name": metric_name})
		return nil
	}
	return ms.metrics_set.GetOrCreateFloatCounter(ms.get_metrics_name(metric_name, labels))
}

func (ms *MetricsSetup) RegisterMetricsSummary(metric_name string, labels map[string]string) *metrics.Summary {
	if validate_metrics_name(metric_name) != nil {
		log.Critical("RegisterMetricsSummary() error", map[string]interface{}{"_error": "Invalid metric name", "_metric_name": metric_name})
		return nil
	}
	return ms.metrics_set.GetOrCreateSummary(ms.get_metrics_name(metric_name, labels))
}

func (ms *MetricsSetup) RegisterMetricsHistogram(metric_name string, labels map[string]string) *metrics.Histogram {
	if validate_metrics_name(metric_name) != nil {
		log.Critical("RegisterMetricsHistogram() error", map[string]interface{}{"_error": "Invalid metric name", "_metric_name": metric_name})
		return nil
	}
	return ms.metrics_set.GetOrCreateHistogram(ms.get_metrics_name(metric_name, labels))
}

func (ms *MetricsSetup) Increment(metric_name string, labels map[string]string) {
	ms.RegisterMetricsCounter(metric_name, labels).Inc()
}

func (ms *MetricsSetup) IncrementFloat(metric_name string, labels map[string]string, value float64) {
	ms.RegisterFloatCounter(metric_name, labels).Add(value)
}

func (ms *MetricsSetup) Set(metric_name string, labels map[string]string, value uint64) {
	ms.RegisterMetricsCounter(metric_name, labels).Set(value)
}

func (ms *MetricsSetup) Update(metric_name string, labels map[string]string, value float64) {
	ms.RegisterMetricsHistogram(metric_name, labels).Update(value)
}

func (ms *MetricsSetup) UpdateDuration(metric_name string, labels map[string]string, value time.Time) {
	ms.RegisterMetricsHistogram(metric_name, labels).UpdateDuration(value)
}

func (ms *MetricsSetup) UpdateSummary(metric_name string, labels map[string]string, value float64) {
	ms.RegisterMetricsSummary(metric_name, labels).Update(value)
}

func (ms *MetricsSetup) RemoveMetrics(metric_name string, labels map[string]string) {
	ms.metrics_set.UnregisterMetric(ms.get_metrics_name(metric_name, labels))
}

func (ms *MetricsSetup) PurgeMetrics() {
	ms.metrics_set.UnregisterAllMetrics()
}
