package libpack_monitoring

import (
	"os"
	"testing"

	libpack_config "github.com/lukaszraczylo/graphql-monitoring-proxy/config"
	"github.com/stretchr/testify/assert"
)

func TestGetMetricsName(t *testing.T) {
	ms := &MetricsSetup{metrics_prefix: "prefix"}
	libpack_config.PKG_NAME = "example_microservice"

	tests := []struct {
		name           string
		metricName     string
		labels         map[string]string
		expectedOutput string
	}{
		{
			name:           "No labels",
			metricName:     "test_metric",
			labels:         nil,
			expectedOutput: "prefix_test_metric{microservice=\"example_microservice\",pod=\"" + getPodName() + "\"}",
		},
		{
			name:       "With labels",
			metricName: "test_metric",
			labels: map[string]string{
				"label1": "value1",
				"label2": "value2",
			},
			expectedOutput: "prefix_test_metric{label1=\"value1\",label2=\"value2\",microservice=\"example_microservice\",pod=\"" + getPodName() + "\"}",
		},
		{
			name:       "Alphabetical order labels",
			metricName: "test_metric",
			labels: map[string]string{
				"label2": "value2",
				"label1": "value1",
			},
			expectedOutput: "prefix_test_metric{label1=\"value1\",label2=\"value2\",microservice=\"example_microservice\",pod=\"" + getPodName() + "\"}",
		},
		{
			name:           "Empty metric name",
			metricName:     "",
			labels:         nil,
			expectedOutput: "prefix_{microservice=\"example_microservice\",pod=\"" + getPodName() + "\"}",
		},
		{
			name:           "Empty labels map",
			metricName:     "test_metric",
			labels:         map[string]string{},
			expectedOutput: "prefix_test_metric{microservice=\"example_microservice\",pod=\"" + getPodName() + "\"}",
		},
		{
			name:       "Single label",
			metricName: "test_metric",
			labels: map[string]string{
				"label1": "value1",
			},
			expectedOutput: "prefix_test_metric{label1=\"value1\",microservice=\"example_microservice\",pod=\"" + getPodName() + "\"}",
		},
		{
			name:       "Multiple labels with special characters",
			metricName: "test_metric",
			labels: map[string]string{
				"label-2": "value-2",
				"label_1": "value_1",
			},
			expectedOutput: "prefix_test_metric{label-2=\"value-2\",label_1=\"value_1\",microservice=\"example_microservice\",pod=\"" + getPodName() + "\"}",
		},
		{
			name:       "Prefix only",
			metricName: "",
			labels: map[string]string{
				"label1": "value1",
			},
			expectedOutput: "prefix_{label1=\"value1\",microservice=\"example_microservice\",pod=\"" + getPodName() + "\"}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ms.get_metrics_name(tt.metricName, tt.labels)
			assert.Equal(t, tt.expectedOutput, result)
		})
	}
}

func TestCompileMetricsWithLabels(t *testing.T) {
	tests := []struct {
		name   string
		labels map[string]string
		want   string
	}{
		{"request_count", map[string]string{"env": "production", "region": "us-west-2"}, "request_count_env_production_region_us-west-2"},
		{"metric_name", map[string]string{}, "metric_name"},
		{"metric_name", nil, "metric_name"},
		{"metric_name", map[string]string{"key1": "value1"}, "metric_name_key1_value1"},
		{"metric_name", map[string]string{"k": "v", "key2": "value2"}, "metric_name_k_v_key2_value2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := compile_metrics_with_labels(tt.name, tt.labels); got != tt.want {
				t.Errorf("compile_metrics_with_labels() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateMetricsName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"Valid name", "valid_metric_name", false},
		{"Name with spaces", "valid metric name", true},
		{"Name with special chars", "valid@metric#name!", true},
		{"Name with leading underscore", "_valid_metric_name", true},
		{"Name with trailing underscore", "valid_metric_name_", true},
		{"Name with consecutive underscores", "valid__metric__name", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validate_metrics_name(tt.input); (err != nil) != tt.wantErr {
				t.Errorf("validate_metrics_name() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func getPodName() string {
	podName, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return podName
}
