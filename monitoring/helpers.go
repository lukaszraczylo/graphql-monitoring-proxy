package libpack_monitoring

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"unicode"

	libpack_config "github.com/lukaszraczylo/graphql-monitoring-proxy/config"
)

func (ms *MetricsSetup) get_metrics_name(name string, labels map[string]string) (complete_name string) {
	const unknownPodName = "unknown"
	var sb strings.Builder

	// Prepare default labels without initializing a new map
	podName := unknownPodName
	if hn, err := os.Hostname(); err == nil {
		podName = hn
	}
	if labels == nil {
		labels = map[string]string{
			"microservice": libpack_config.PKG_NAME,
			"pod":          podName,
		}
	} else {
		if _, exists := labels["microservice"]; !exists {
			labels["microservice"] = libpack_config.PKG_NAME
		}
		if _, exists := labels["pod"]; !exists {
			labels["pod"] = podName
		}
	}

	// Prefix handling
	if ms.metrics_prefix != "" {
		sb.WriteString(ms.metrics_prefix)
		sb.WriteString("_")
	}
	sb.WriteString(name)

	// Append labels if any
	if len(labels) > 0 {
		sb.WriteString("{")

		keys := make([]string, 0, len(labels))
		for k := range labels {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for i, k := range keys {
			if i > 0 {
				sb.WriteString(",")
			}
			sb.WriteString(k)
			sb.WriteString("=\"")
			sb.WriteString(labels[k])
			sb.WriteString("\"")
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
	var totalLength int
	totalLength += len(name)
	for k, v := range labels {
		totalLength += len(k) + len(v) + 2
	}

	var sb strings.Builder
	sb.Grow(totalLength + 1)

	sb.WriteString(name)

	// Collect keys and sort them
	keys := make([]string, 0, len(labels))
	for k := range labels {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Append sorted key-value pairs to the builder
	for _, k := range keys {
		sb.WriteString("_")
		sb.WriteString(k)
		sb.WriteString("_")
		sb.WriteString(labels[k])
	}

	return sb.String()
}
