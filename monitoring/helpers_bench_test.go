package libpack_monitoring

import (
	"testing"

	libpack_config "github.com/lukaszraczylo/graphql-monitoring-proxy/config"
)

func BenchmarkGetMetricsName(b *testing.B) {
	// Setup environment
	libpack_config.PKG_NAME = "test_service"

	ms := &MetricsSetup{metrics_prefix: "test_prefix"}

	labels := map[string]string{
		"env":    "production",
		"region": "us-west-2",
	}

	// Run the benchmark
	for n := 0; n < b.N; n++ {
		ms.get_metrics_name("request_count", labels)
	}
}

func BenchmarkCompileMetricsWithLabels(b *testing.B) {
	labels := map[string]string{
		"env":    "production",
		"region": "us-west-2",
		"app":    "api-server",
	}

	for n := 0; n < b.N; n++ {
		compile_metrics_with_labels("request_count", labels)
	}
}

func BenchmarkValidateMetricsName(b *testing.B) {
	input := "valid metric name with special chars @#! and underscores__"

	for n := 0; n < b.N; n++ {
		_ = validate_metrics_name(input)
	}
}
