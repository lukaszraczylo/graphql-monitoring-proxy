package libpack_monitoring

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestRegistrationMethods_NilReceiver_NoPanic verifies that the metric
// registration methods tolerate a nil *MetricsSetup receiver. Callers may
// register metrics before the monitoring server has constructed the global
// MetricsSetup (e.g. circuit breaker init during config parsing); these
// methods must return a non-nil dummy instead of panicking.
func TestRegistrationMethods_NilReceiver_NoPanic(t *testing.T) {
	var ms *MetricsSetup // nil

	require.NotPanics(t, func() {
		require.NotNil(t, ms.RegisterMetricsGauge("g", nil, 1))
		require.NotNil(t, ms.RegisterMetricsGaugeFunc("gf", nil, func() float64 { return 1 }))
		require.NotNil(t, ms.RegisterMetricsCounter("c", nil))
		require.NotNil(t, ms.RegisterFloatCounter("fc", nil))
		require.NotNil(t, ms.RegisterMetricsSummary("s", nil))
		require.NotNil(t, ms.RegisterMetricsHistogram("h", nil))
	})
}
