package libpack_monitoring

import (
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	libpack_config "github.com/lukaszraczylo/graphql-monitoring-proxy/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestCleanMetricName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"valid metric name", "valid_metric_name"},
		{"valid@metric#name!", "valid_metric_name"},
		{"__valid__metric__name__", "valid_metric_name"},
		{" valid metric name ", "valid_metric_name"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expected, clean_metric_name(tt.input))
		})
	}
}

func TestDefaultLabels(t *testing.T) {
	podName := "test-pod"
	libpack_config.PKG_NAME = "example_microservice"
	expected := map[string]string{
		"microservice": "example_microservice",
		"pod":          podName,
	}

	assert.Equal(t, expected, defaultLabels(podName))
}

func TestEnsureDefaultLabels(t *testing.T) {
	podName := "test-pod"
	libpack_config.PKG_NAME = "example_microservice"

	tests := []struct {
		inputLabels    map[string]string
		expectedLabels map[string]string
		name           string
	}{
		{
			name:           "Nil labels",
			inputLabels:    nil,
			expectedLabels: map[string]string{"microservice": "example_microservice", "pod": podName},
		},
		{
			name:           "Empty labels",
			inputLabels:    map[string]string{},
			expectedLabels: map[string]string{"microservice": "example_microservice", "pod": podName},
		},
		{
			name:           "Partial labels",
			inputLabels:    map[string]string{"microservice": "test_service"},
			expectedLabels: map[string]string{"microservice": "test_service", "pod": podName},
		},
		{
			name:           "Complete labels",
			inputLabels:    map[string]string{"microservice": "test_service", "pod": "custom_pod"},
			expectedLabels: map[string]string{"microservice": "test_service", "pod": "custom_pod"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ensureDefaultLabels(&tt.inputLabels, podName)
			assert.Equal(t, tt.expectedLabels, tt.inputLabels)
		})
	}
}

// TestGetPodName_HostnameResolvedOnce verifies getPodName() resolves the
// hostname exactly once (via podNameOnce) no matter how many times it is
// called, including concurrently. Repeated os.Hostname() syscalls on every
// metrics call were the dominant per-request cost this caching removes.
func TestGetPodName_HostnameResolvedOnce(t *testing.T) {
	origFn := hostnameFunc
	t.Cleanup(func() {
		// Restore the real hostname resolver and reset to a fresh, unfired
		// Once (never copy an existing sync.Once - it must not be copied by
		// value once it may have been used). The next getPodName() call will
		// re-resolve the real hostname and re-cache it for later tests.
		hostnameFunc = origFn
		podNameOnce = sync.Once{}
	})

	podNameOnce = sync.Once{}
	var callCount int32
	hostnameFunc = func() (string, error) {
		atomic.AddInt32(&callCount, 1)
		return "cached-test-host", nil
	}

	// Sequential calls must all observe the cached value.
	first := getPodName()
	second := getPodName()
	assert.Equal(t, "cached-test-host", first)
	assert.Equal(t, first, second)
	assert.Equal(t, int32(1), atomic.LoadInt32(&callCount), "hostnameFunc must be invoked exactly once")

	// Concurrent calls must not re-invoke hostnameFunc either.
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			assert.Equal(t, "cached-test-host", getPodName())
		}()
	}
	wg.Wait()

	assert.Equal(t, int32(1), atomic.LoadInt32(&callCount), "hostnameFunc must still be invoked exactly once after concurrent calls")
}

// TestGetSortedKeysCache_SharedAcrossValues verifies the sortedLabelKeysCache
// is keyed by the sorted set of label NAMES only. Two maps sharing the same
// keys but different values must collapse into a single cache entry rather
// than growing the cache per distinct value combination (unbounded
// cardinality bomb).
func TestGetSortedKeysCache_SharedAcrossValues(t *testing.T) {
	labelsA := map[string]string{"cache_key_test_a": "value1", "cache_key_test_b": "value2"}
	labelsB := map[string]string{"cache_key_test_a": "totally-different", "cache_key_test_b": "also-different"}

	keysA := getSortedKeys(labelsA)
	keysB := getSortedKeys(labelsB)

	require.Equal(t, []string{"cache_key_test_a", "cache_key_test_b"}, keysA)
	require.Equal(t, keysA, keysB)

	// Same underlying slice (same cache entry) must be returned for both
	// maps, proving the cache did not create a second, value-keyed entry.
	assert.Same(t, &keysA[0], &keysB[0])

	// The cache entry itself must be reachable only via a names-only key.
	namesKey := strings.Join(keysA, ";")
	cached, ok := sortedLabelKeysCache.m.Load(namesKey)
	require.True(t, ok)
	assert.Equal(t, keysA, cached)
}

func TestLabelsToString(t *testing.T) {
	tests := []struct {
		labels   map[string]string
		expected string
	}{
		{
			labels:   map[string]string{"key1": "value1", "key2": "value2"},
			expected: "key1=value1;key2=value2;",
		},
		{
			labels:   map[string]string{"a": "1", "b": "2"},
			expected: "a=1;b=2;",
		},
		{
			labels:   map[string]string{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			assert.Equal(t, tt.expected, labelsToString(tt.labels))
		})
	}
}
