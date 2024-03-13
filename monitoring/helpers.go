package libpack_monitoring

import (
	"fmt"
	"os"
	"strings"
	"unicode"

	libpack_config "github.com/lukaszraczylo/graphql-monitoring-proxy/config"
)

func (ms *MetricsSetup) get_metrics_name(name string, labels map[string]string) (complete_name string) {
	if labels == nil {
		labels = make(map[string]string)
	}

	// Adding default labels
	labels["microservice"] = libpack_config.PKG_NAME
	if podName, err := os.Hostname(); err == nil {
		labels["pod"] = podName
	} else {
		labels["pod"] = "unknown"
	}

	var sb strings.Builder
	if ms.metrics_prefix != "" {
		sb.WriteString(ms.metrics_prefix)
		sb.WriteString("_")
	}
	sb.WriteString(name)

	if len(labels) > 0 {
		sb.WriteString("{")
		first := true
		for k, v := range labels {
			if !first {
				sb.WriteString(",")
			}
			sb.WriteString(k)
			sb.WriteString("=\"")
			sb.WriteString(v)
			sb.WriteString("\"")
			first = false
		}
		sb.WriteString("}")
	}

	return sb.String()
}

// validate_metrics_name validates the name of the metric to adhere to the Prometheus naming conventions
// https://prometheus.io/docs/practices/naming/
func validate_metrics_name(name string) error {
	var sb strings.Builder // Use strings.Builder for efficient string concatenation

	// Track if the last character was an underscore to avoid duplicate underscores
	lastWasUnderscore := false

	for _, r := range name {
		// Convert spaces to underscores and skip non-alphanumeric characters except underscores
		if r == ' ' || (unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_') {
			if r == ' ' || r == '_' {
				if lastWasUnderscore {
					continue // Skip if the previous character was also an underscore
				}
				r = '_' // Convert spaces to underscores
				lastWasUnderscore = true
			} else {
				lastWasUnderscore = false
			}
			sb.WriteRune(r) // Add valid characters to the builder
		}
	}
	// Trim leading and trailing underscores
	name_new := strings.Trim(sb.String(), "_")

	// Check if the processed name matches the original input
	if name_new != name {
		return fmt.Errorf("Invalid metric name: %s, expected %s", name, name_new)
	}
	return nil
}

func compile_metrics_with_labels(name string, labels map[string]string) string {
	metric_name := name
	for k, v := range labels {
		metric_name += "_" + k + "_" + v
	}
	return metric_name
}
