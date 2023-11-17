package libpack_monitoring

import (
	"fmt"
	"os"
	"sort"
	"strings"

	libpack_config "github.com/lukaszraczylo/graphql-monitoring-proxy/config"
)

func (ms *MetricsSetup) get_metrics_name(name string, labels map[string]string) (complete_name string) {
	var err error
	if labels == nil {
		labels = make(map[string]string)
	}
	labels["microservice"] = libpack_config.PKG_NAME
	labels["pod"], err = os.Hostname()
	if err != nil {
		labels["pod"] = "unknown"
	}

	if ms.metrics_prefix != "" {
		complete_name = ms.metrics_prefix + "_" + name
	} else {
		complete_name = name
	}
	if labels != nil {
		keys := make([]string, 0, len(labels))
		for k := range labels {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		complete_name += "{"
		for _, k := range keys {
			complete_name += k + "=\"" + labels[k] + "\","
		}
		complete_name = strings.TrimSuffix(complete_name, ",")
		complete_name += "}"
	}
	return
}

// validate_metrics_name validates the name of the metric to adhere to the Prometheus naming conventions
// https://prometheus.io/docs/practices/naming/
func validate_metrics_name(name string) error {
	// replace all spaces with underscores and remove all other non-alphanumeric characters
	name_new := strings.ReplaceAll(name, " ", "_")
	name_new = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' {
			return r
		}
		return -1
	}, name_new)
	name_new = strings.ReplaceAll(name_new, "__", "_")
	name_new = strings.Trim(name_new, "_")
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
