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

const unknownPodName = "unknown"

var (
	// podNameOnce guards the one-time resolution of the local hostname.
	// os.Hostname() is a syscall; resolving it on every metrics call (which
	// happens several times per request) is wasteful since the pod name
	// never changes for the lifetime of the process.
	podNameOnce  sync.Once
	podNameValue string
	// hostnameFunc resolves the local hostname. It is a package var (rather
	// than a direct os.Hostname() call) so tests can substitute a counting
	// stub and verify podNameOnce caches the result after the first call.
	hostnameFunc = os.Hostname
)

func getPodName() string {
	podNameOnce.Do(func() {
		if hn, err := hostnameFunc(); err == nil {
			podNameValue = hn
		} else {
			podNameValue = unknownPodName
		}
	})
	return podNameValue
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

// sanitizeLabelValue removes or replaces characters that are invalid in metric labels
// This includes null bytes, newlines, carriage returns, quotes, and backslashes
func sanitizeLabelValue(value string) string {
	if value == "" {
		return value
	}

	var buf strings.Builder
	buf.Grow(len(value))

	for _, r := range value {
		switch r {
		case '\x00': // null byte
			continue // Skip null bytes entirely
		case '\n', '\r', '\t': // newlines, carriage returns, tabs
			buf.WriteByte(' ') // Replace with space
		case '"', '\\': // quotes and backslashes need escaping
			buf.WriteByte('\\')
			buf.WriteRune(r)
		default:
			// Only allow printable ASCII and common unicode characters
			if unicode.IsPrint(r) {
				buf.WriteRune(r)
			}
		}
	}

	return buf.String()
}

func appendSortedLabels(buf *bytes.Buffer, labels map[string]string) {
	// Add defer/recover to prevent panics from crashing the application
	defer func() {
		if r := recover(); r != nil {
			// Log the panic but don't crash
			fmt.Fprintf(os.Stderr, "Recovered from panic in appendSortedLabels: %v\n", r)
		}
	}()

	if len(labels) == 0 || buf == nil {
		return
	}

	// Create a snapshot to avoid concurrent access issues
	labelsCopy := make(map[string]string, len(labels))
	for k, v := range labels {
		if k == "" {
			continue // Skip empty keys
		}
		// Sanitize the label value to remove null bytes and other invalid characters
		labelsCopy[k] = sanitizeLabelValue(v)
	}

	if len(labelsCopy) == 0 {
		return
	}

	keys := getSortedKeys(labelsCopy)
	for i, k := range keys {
		if v, ok := labelsCopy[k]; ok {
			if i > 0 {
				buf.WriteByte(',')
			}
			buf.WriteString(k)
			buf.WriteString(`="`)
			buf.WriteString(v)
			buf.WriteByte('"')
		}
	}
}

func getSortedKeys(labels map[string]string) []string {
	if labels == nil {
		return []string{}
	}

	// Compute the sorted keys - create a snapshot to avoid concurrent access issues
	keys := make([]string, 0, len(labels))
	for k := range labels {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// The cache is keyed by the sorted set of label NAMES only, not values.
	// The cached payload (the sorted key list) never depends on the label
	// values, so keying by value (as before) created a separate, never-
	// evicted cache entry per distinct value combination - an unbounded
	// cardinality bomb. Keying by names collapses all maps sharing the same
	// key set into a single entry, bounding growth to the number of
	// distinct label-name sets used by the application.
	namesKey := strings.Join(keys, ";")

	if cached, ok := sortedLabelKeysCache.m.Load(namesKey); ok {
		return cached.([]string)
	}

	sortedLabelKeysCache.m.Store(namesKey, keys)

	return keys
}

func labelsToString(labels map[string]string) string {
	// Add defer/recover to prevent panics from crashing the application
	defer func() {
		if r := recover(); r != nil {
			// Log the panic but don't crash
			fmt.Fprintf(os.Stderr, "Recovered from panic in labelsToString: %v\n", r)
		}
	}()

	if len(labels) == 0 {
		return ""
	}

	// Create a snapshot of the map to avoid concurrent access issues
	keys := make([]string, 0, len(labels))
	values := make(map[string]string, len(labels))

	for k, v := range labels {
		if k == "" {
			continue // Skip empty keys
		}
		keys = append(keys, k)
		values[k] = v
	}

	if len(keys) == 0 {
		return ""
	}

	sort.Strings(keys)

	// Pre-allocate the builder with estimated capacity to avoid reallocation
	var sb strings.Builder
	estimatedSize := 0
	for _, k := range keys {
		estimatedSize += len(k) + len(values[k]) + 2 // key + value + '=' + ';'
	}
	sb.Grow(estimatedSize)

	for _, k := range keys {
		if v, ok := values[k]; ok {
			sb.WriteString(k)
			sb.WriteByte('=')
			sb.WriteString(v)
			sb.WriteByte(';')
		}
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
	// Add defer/recover to prevent panics from crashing the application
	defer func() {
		if r := recover(); r != nil {
			// Log the panic but don't crash
			fmt.Fprintf(os.Stderr, "Recovered from panic in compile_metrics_with_labels: %v\n", r)
		}
	}()

	var buf bytes.Buffer

	buf.WriteString(name)

	if len(labels) == 0 {
		return buf.String()
	}

	// Create a snapshot to avoid concurrent access issues
	labelsCopy := make(map[string]string, len(labels))
	for k, v := range labels {
		if k == "" {
			continue // Skip empty keys
		}
		labelsCopy[k] = v
	}

	if len(labelsCopy) == 0 {
		return buf.String()
	}

	keys := getSortedKeys(labelsCopy)

	for _, k := range keys {
		if v, ok := labelsCopy[k]; ok {
			buf.WriteByte('_')
			buf.WriteString(k)
			buf.WriteByte('_')
			buf.WriteString(v)
		}
	}

	return buf.String()
}
