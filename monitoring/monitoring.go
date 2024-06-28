package libpack_monitoring

import (
	"flag"
	"fmt"
	"time"

	"github.com/VictoriaMetrics/metrics"
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/goutil/envutil"
	libpack_config "github.com/lukaszraczylo/graphql-monitoring-proxy/config"
	libpack_logger "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
)

type MetricsSetup struct {
	metrics_set        *metrics.Set
	metrics_set_custom *metrics.Set
	ic                 *InitConfig
	metrics_prefix     string
}

var log = libpack_logger.New().SetMinLogLevel(libpack_logger.LEVEL_INFO)

type InitConfig struct {
	PurgeOnCrawl bool
	PurgeEvery   int
}

func NewMonitoring(ic *InitConfig) *MetricsSetup {
	ms := &MetricsSetup{
		ic:                 ic,
		metrics_set:        metrics.NewSet(),
		metrics_set_custom: metrics.NewSet(),
	}

	if flag.Lookup("test.v") == nil {
		go ms.startPrometheusEndpoint()

		if ic.PurgeEvery > 0 {
			ticker := time.NewTicker(time.Duration(ic.PurgeEvery) * time.Second)
			go func() {
				for range ticker.C {
					ms.PurgeMetrics()
				}
			}()
		}
	}

	return ms
}

func (ms *MetricsSetup) startPrometheusEndpoint() {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		AppName:               fmt.Sprintf("GraphQL Monitoring Proxy - %s v%s", libpack_config.PKG_NAME, libpack_config.PKG_VERSION),
	})
	app.Get("/metrics", ms.metricsEndpoint)
	if err := app.Listen(fmt.Sprintf(":%d", envutil.GetInt("MONITORING_PORT", 9393))); err != nil {
		log.Critical(&libpack_logger.LogMessage{
			Message: "Can't start the service",
			Pairs:   map[string]interface{}{"error": err},
		})
	}
}

func (ms *MetricsSetup) metricsEndpoint(c *fiber.Ctx) error {
	ms.metrics_set.WritePrometheus(c.Response().BodyWriter())
	ms.metrics_set_custom.WritePrometheus(c.Response().BodyWriter())

	if ms.ic.PurgeOnCrawl && ms.ic.PurgeEvery == 0 {
		ms.PurgeMetrics()
	}
	return nil
}

func (ms *MetricsSetup) AddMetricsPrefix(prefix string) {
	ms.metrics_prefix = prefix
}

func (ms *MetricsSetup) ListActiveMetrics() []string {
	return ms.metrics_set.ListMetricNames()
}

func (ms *MetricsSetup) RegisterMetricsGauge(metric_name string, labels map[string]string, val float64) *metrics.Gauge {
	if err := validate_metrics_name(metric_name); err != nil {
		log.Critical(&libpack_logger.LogMessage{
			Message: "RegisterMetricsGauge() error",
			Pairs:   map[string]interface{}{"_error": "Invalid metric name", "_metric_name": metric_name},
		})
		return nil
	}
	return ms.metrics_set_custom.GetOrCreateGauge(ms.get_metrics_name(metric_name, labels), func() float64 {
		return val
	})
}

func (ms *MetricsSetup) RegisterMetricsCounter(metric_name string, labels map[string]string) *metrics.Counter {
	if err := validate_metrics_name(metric_name); err != nil {
		log.Critical(&libpack_logger.LogMessage{
			Message: "RegisterMetricsCounter() error",
			Pairs:   map[string]interface{}{"_error": "Invalid metric name", "_metric_name": metric_name},
		})
		return nil
	}
	if metric_name == MetricsSucceeded || metric_name == MetricsFailed || metric_name == MetricsSkipped {
		return ms.metrics_set.GetOrCreateCounter(ms.get_metrics_name(metric_name, labels))
	}
	return ms.metrics_set_custom.GetOrCreateCounter(ms.get_metrics_name(metric_name, labels))
}

func (ms *MetricsSetup) RegisterFloatCounter(metric_name string, labels map[string]string) *metrics.FloatCounter {
	if err := validate_metrics_name(metric_name); err != nil {
		log.Critical(&libpack_logger.LogMessage{
			Message: "RegisterFloatCounter() error",
			Pairs:   map[string]interface{}{"_error": "Invalid metric name", "_metric_name": metric_name},
		})
		return nil
	}
	return ms.metrics_set_custom.GetOrCreateFloatCounter(ms.get_metrics_name(metric_name, labels))
}

func (ms *MetricsSetup) RegisterMetricsSummary(metric_name string, labels map[string]string) *metrics.Summary {
	if err := validate_metrics_name(metric_name); err != nil {
		log.Critical(&libpack_logger.LogMessage{
			Message: "RegisterMetricsSummary() error",
			Pairs:   map[string]interface{}{"_error": "Invalid metric name", "_metric_name": metric_name},
		})
		return nil
	}
	return ms.metrics_set_custom.GetOrCreateSummary(ms.get_metrics_name(metric_name, labels))
}

func (ms *MetricsSetup) RegisterMetricsHistogram(metric_name string, labels map[string]string) *metrics.Histogram {
	if err := validate_metrics_name(metric_name); err != nil {
		log.Critical(&libpack_logger.LogMessage{
			Message: "RegisterMetricsHistogram() error",
			Pairs:   map[string]interface{}{"_error": "Invalid metric name", "_metric_name": metric_name},
		})
		return nil
	}
	return ms.metrics_set_custom.GetOrCreateHistogram(ms.get_metrics_name(metric_name, labels))
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
	ms.metrics_set_custom.UnregisterMetric(ms.get_metrics_name(metric_name, labels))
}

func (ms *MetricsSetup) PurgeMetrics() {
	ms.metrics_set_custom.UnregisterAllMetrics()
}
