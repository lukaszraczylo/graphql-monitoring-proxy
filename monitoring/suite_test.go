package libpack_monitoring

import (
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// import (
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"testing"
// 	"time"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/suite"
// 	libpack_config "github.com/lukaszraczylo/graphql-monitoring-proxy/config"
// )

type MonitoringTestSuite struct {
	suite.Suite
	metrics_endpoint string
}

func (suite *MonitoringTestSuite) SetupTest() {
	suite.metrics_endpoint = "http://localhost:9393/metrics"
}

func (suite *MonitoringTestSuite) TearDownTest() {
}

func TestMonitoringSuite(t *testing.T) {
	suite.Run(t, new(MonitoringTestSuite))
}

func (suite *MonitoringTestSuite) testing_call_metrics_endpoint() *http.Response {
	resp, err := http.Get(suite.metrics_endpoint)
	if err != nil {
		suite.T().Error("Can't call metrics endpoint", err)
		suite.T().FailNow()
	}
	return resp
}

func (suite *MonitoringTestSuite) testing_get_body_as_text(resp *http.Response) string {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		suite.T().Error("Can't read response body", err)
		suite.T().FailNow()
	}
	return string(body)
}

func (suite *MonitoringTestSuite) TestNewMonitoring() {
	metrics_prefix := "within_test"

	suite.T().Run("TestWholeEndpoint", func(t *testing.T) {
		mon := NewMonitoring()
		mon.AddMetricsPrefix(metrics_prefix)
		mon.RegisterDefaultMetrics()

		resp := suite.testing_call_metrics_endpoint()
		assert.Equal(t, 200, resp.StatusCode)
		assert.Contains(t, suite.testing_get_body_as_text(resp), "within_test_requests_succesful", "Metrics endpoint should contain metrics with prefix %s", metrics_prefix)

		mon.RegisterMetricsGauge("test_gauge_metrics", nil, 14)
		resp = suite.testing_call_metrics_endpoint()
		assert.Equal(t, 200, resp.StatusCode)
		assert.Contains(t, suite.testing_get_body_as_text(resp), "within_test_test_gauge_metrics", "Metrics endpoint should contain metrics with prefix %s", metrics_prefix)

		// Triggering it again will increase the gauge by 1
		mon.RegisterMetricsGauge("test_gauge_metrics", nil, 7)
		time.Sleep(1 * time.Second)
		resp = suite.testing_call_metrics_endpoint()
		assert.Equal(t, 200, resp.StatusCode)
		assert.Contains(t, suite.testing_get_body_as_text(resp), "within_test_test_gauge_metrics 7", "Metrics endpoint should contain incremented () gauge metrics with prefix %s", metrics_prefix)

		mon.RegisterMetricsCounter("test_counter_metrics", nil)
		mon.Increment("test_counter_metrics", nil)
		time.Sleep(1 * time.Second)
		resp = suite.testing_call_metrics_endpoint()
		assert.Equal(t, 200, resp.StatusCode)
		assert.Contains(t, suite.testing_get_body_as_text(resp), "within_test_test_counter_metrics", "Metrics endpoint should contain metrics with prefix %s", metrics_prefix)
		mon.Increment("test_counter_metrics", nil)
		time.Sleep(1 * time.Second)
		resp = suite.testing_call_metrics_endpoint()
		assert.Equal(t, 200, resp.StatusCode)
		assert.Contains(t, suite.testing_get_body_as_text(resp), "within_test_test_counter_metrics 2", "Metrics endpoint should contain metrics with prefix %s", metrics_prefix)

		// 	mon.AddCustomMetrics(&CustomMetrics{
		// 		Name: "test_custom_metrics",
		// 		Help: "test custom metrics",
		// 		Type: TypeHistogram,
		// 	}, libpack_config.PKG_NAME)
		// 	resp = suite.testing_call_metrics_endpoint()
		// 	assert.Equal(t, 200, resp.StatusCode)
		// 	assert.Contains(t, suite.testing_get_body_as_text(resp), "within_test_test_custom_metrics", "Metrics endpoint should contain metrics with prefix %s", metrics_prefix)
		// 	fmt.Println(suite.testing_get_body_as_text(resp))

		// 	assert.Containsf(t, mon.ListActiveMetrics(), "test_custom_metrics", "ListActiveMetrics() should contain metrics with prefix %s", metrics_prefix)

		// 	mon.AddCustomMetrics(&CustomMetrics{
		// 		Name: "test_gauge_metrics",
		// 		Help: "test gauge metrics",
		// 		Type: TypeGauge,
		// 	}, libpack_config.PKG_NAME)
		// 	mon.Increment("test_gauge_metrics")
		// 	time.Sleep(2 * time.Second)
		// 	resp = suite.testing_call_metrics_endpoint()
		// 	assert.Equal(t, 200, resp.StatusCode)
		// 	assert.Contains(t, suite.testing_get_body_as_text(resp), "test_gauge_metrics{microservice=\""+libpack_config.PKG_NAME+"\"} 1", "Metrics endpoint should contain incremented (1) gauge metrics with prefix %s", metrics_prefix)

		// 	mon.Increment("test_gauge_metrics")
		// 	time.Sleep(2 * time.Second)
		// 	resp = suite.testing_call_metrics_endpoint()
		// 	assert.Equal(t, 200, resp.StatusCode)
		// 	assert.Contains(t, suite.testing_get_body_as_text(resp), "test_gauge_metrics{microservice=\""+libpack_config.PKG_NAME+"\"} 2", "Metrics endpoint should contain incremented (2) gauge metrics with prefix %s", metrics_prefix)

		// 	mon.AddCustomMetrics(&CustomMetrics{
		// 		Name: "test_counter_metrics",
		// 		Help: "test counter metrics",
		// 		Type: TypeCounter,
		// 	}, libpack_config.PKG_NAME)
		// 	mon.Set("test_counter_metrics", 7)
		// 	time.Sleep(2 * time.Second)
		// 	resp = suite.testing_call_metrics_endpoint()
		// 	assert.Equal(t, 200, resp.StatusCode)
		// 	assert.Contains(t, suite.testing_get_body_as_text(resp), "test_counter_metrics{microservice=\""+libpack_config.PKG_NAME+"\"} 7", "Metrics endpoint should contain metrics with prefix %s", metrics_prefix)

		resp = suite.testing_call_metrics_endpoint()
		assert.Equal(t, 200, resp.StatusCode)
		assert.Contains(t, suite.testing_get_body_as_text(resp), "go_goroutines", "Metrics endpoint should contain metrics with prefix %s", metrics_prefix)

	})
}
