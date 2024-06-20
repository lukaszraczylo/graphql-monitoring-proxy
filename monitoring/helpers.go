package libpack_monitoring

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"unicode"

	libpack_config "github.com/lukaszraczylo/graphql-monitoring-proxy/config"
)

// Cache for sorted label keys to avoid repeated sorting
var sortedLabelKeysCache = struct {
	m map[string][]string
	sync.RWMutex
}{m: make(map[string][]string)}

func (ms *MetricsSetup) get_metrics_name(name string, labels map[string]string) string {
	const unknownPodName = "unknown"
	var buf bytes.Buffer

	// Prepare default labels without initializing a new map
	podName := getPodName()
	if labels == nil {
		labels = defaultLabels(podName)
	} else {
		ensureDefaultLabels(&labels, podName)
	}

	// Prefix handling
	if ms.metrics_prefix != "" {
		buf.WriteString(ms.metrics_prefix)
		buf.WriteString("_")
	}
	buf.WriteString(name)

	// Append labels if any
	if len(labels) > 0 {
		buf.WriteString("{")
		appendSortedLabels(&buf, labels)
		buf.WriteString("}")
	}

	return buf.String()
}

func getPodName() string {
	const unknownPodName = "unknown"
	if hn, err := os.Hostname(); err == nil {
		return hn
	}
	return unknownPodName
}

func defaultLabels(podName string) map[string]string {
	return map[string]string{
		"microservice": libpack_config.PKG_NAME,
		"pod":          podName,
	}
}

func ensureDefaultLabels(labels *map[string]string, podName string) {
	if *labels == nil {
		*labels = make(map[string]string)
	}
	if _, exists := (*labels)["microservice"]; !exists {
		(*labels)["microservice"] = libpack_config.PKG_NAME
	}
	if _, exists := (*labels)["pod"]; !exists {
		(*labels)["pod"] = podName
	}
}

func appendSortedLabels(buf *bytes.Buffer, labels map[string]string) {
	keys := getSortedKeys(labels)
	for i, k := range keys {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(k)
		buf.WriteString("=\"")
		buf.WriteString(labels[k])
		buf.WriteString("\"")
	}
}

func getSortedKeys(labels map[string]string) []string {
	labelsKey := labelsToString(labels)

	sortedLabelKeysCache.RLock()
	keys, exists := sortedLabelKeysCache.m[labelsKey]
	sortedLabelKeysCache.RUnlock()

	if !exists {
		keys = make([]string, 0, len(labels))
		for k := range labels {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		sortedLabelKeysCache.Lock()
		sortedLabelKeysCache.m[labelsKey] = keys
		sortedLabelKeysCache.Unlock()
	}

	return keys
}

func labelsToString(labels map[string]string) string {
	var sb strings.Builder
	for k, v := range labels {
		sb.WriteString(k)
		sb.WriteString("=")
		sb.WriteString(v)
		sb.WriteString(";")
	}
	return sb.String()
}

// validate_metrics_name validates the name of the metric to adhere to the Prometheus naming conventions
// https://prometheus.io/docs/practices/naming/
func validate_metrics_name(name string) error {
	cleanedName := clean_metric_name(name)

	// Trim leading and trailing underscores
	finalName := strings.Trim(cleanedName, "_")

	// Check if the processed name matches the original input
	if finalName != name {
		return fmt.Errorf("Invalid metric name: %s, expected %s", name, finalName)
	}
	return nil
}

// clean_metric_name processes the metric name according to Prometheus naming conventions
func clean_metric_name(name string) string {
	var buf bytes.Buffer
	lastWasUnderscore := false

	for _, r := range name {
		if is_allowed_rune(r) {
			if is_special_rune(r) {
				if lastWasUnderscore {
					continue // Skip if the previous character was also an underscore
				}
				r = '_' // Convert spaces and special characters to underscores
				lastWasUnderscore = true
			} else {
				lastWasUnderscore = false
			}
			buf.WriteRune(r)
		} else if !lastWasUnderscore {
			buf.WriteRune('_')
			lastWasUnderscore = true
		}
	}

	// Remove trailing underscore
	result := buf.String()
	return strings.Trim(result, "_")
}

// is_allowed_rune checks if the rune is allowed in the metric name
func is_allowed_rune(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == ' ' || r == '_'
}

// is_special_rune checks if the rune is a space or an underscore
func is_special_rune(r rune) bool {
	return r == ' ' || r == '_'
}

func compile_metrics_with_labels(name string, labels map[string]string) string {
	var buf bytes.Buffer

	buf.WriteString(name)

	// Collect keys and sort them
	keys := getSortedKeys(labels)

	// Append sorted key-value pairs to the buffer
	for _, k := range keys {
		buf.WriteString("_")
		buf.WriteString(k)
		buf.WriteString("_")
		buf.WriteString(labels[k])
	}

	return buf.String()
}
