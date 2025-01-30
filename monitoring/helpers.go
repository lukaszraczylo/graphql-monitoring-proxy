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

var sortedLabelKeysCache = struct {
	m sync.Map
}{}

func (ms *MetricsSetup) get_metrics_name(name string, labels map[string]string) string {
	var buf bytes.Buffer

	podName := getPodName()
	if labels == nil {
		labels = defaultLabels(podName)
	} else {
		ensureDefaultLabels(&labels, podName)
	}

	if ms.metrics_prefix != "" {
		buf.WriteString(ms.metrics_prefix)
		buf.WriteByte('_')
	}
	buf.WriteString(name)

	if len(labels) > 0 {
		buf.WriteByte('{')
		appendSortedLabels(&buf, labels)
		buf.WriteByte('}')
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
			buf.WriteByte(',')
		}
		buf.WriteString(k)
		buf.WriteString(`="`)
		buf.WriteString(labels[k])
		buf.WriteByte('"')
	}
}

func getSortedKeys(labels map[string]string) []string {
	labelsKey := labelsToString(labels)

	// Check if the sorted keys are already cached
	if keys, ok := sortedLabelKeysCache.m.Load(labelsKey); ok {
		return keys.([]string)
	}

	// Compute the sorted keys
	keys := make([]string, 0, len(labels))
	for k := range labels {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Store the sorted keys in the cache
	sortedLabelKeysCache.m.Store(labelsKey, keys)

	return keys
}

func labelsToString(labels map[string]string) string {
	keys := make([]string, 0, len(labels))
	for k := range labels {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		sb.WriteString(k)
		sb.WriteByte('=')
		sb.WriteString(labels[k])
		sb.WriteByte(';')
	}
	return sb.String()
}

func validate_metrics_name(name string) error {
	cleanedName := clean_metric_name(name)

	finalName := strings.Trim(cleanedName, "_")

	if finalName != name {
		return fmt.Errorf("invalid metric name: %s, expected %s", name, finalName)
	}
	return nil
}

func clean_metric_name(name string) string {
	var buf bytes.Buffer
	lastWasUnderscore := false

	for _, r := range name {
		if is_allowed_rune(r) {
			if is_special_rune(r) {
				if lastWasUnderscore {
					continue
				}
				r = '_'
				lastWasUnderscore = true
			} else {
				lastWasUnderscore = false
			}
			buf.WriteRune(r)
		} else if !lastWasUnderscore {
			buf.WriteByte('_')
			lastWasUnderscore = true
		}
	}

	return strings.Trim(buf.String(), "_")
}

func is_allowed_rune(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == ' ' || r == '_'
}

func is_special_rune(r rune) bool {
	return r == ' ' || r == '_'
}

func compile_metrics_with_labels(name string, labels map[string]string) string {
	var buf bytes.Buffer

	buf.WriteString(name)

	keys := getSortedKeys(labels)

	for _, k := range keys {
		buf.WriteByte('_')
		buf.WriteString(k)
		buf.WriteByte('_')
		buf.WriteString(labels[k])
	}

	return buf.String()
}
