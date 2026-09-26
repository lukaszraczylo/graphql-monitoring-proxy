package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	libpack_logger "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
)

func TestBackendReadinessProtectedGraphQL(t *testing.T) {
	var readyStatus atomic.Int32
	readyStatus.Store(http.StatusOK)
	var graphqlCalls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/graphql":
			graphqlCalls.Add(1)
			w.WriteHeader(http.StatusUnauthorized)
		case "/ready":
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Empty(t, r.Header.Get("Authorization"))
			assert.Empty(t, r.Header.Get("X-Husaria-Admin-Secret"))
			assert.Empty(t, r.Header.Get("X-Hasura-Admin-Secret"))
			body, err := io.ReadAll(r.Body)
			assert.NoError(t, err)
			assert.Empty(t, body)
			w.WriteHeader(int(readyStatus.Load()))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	t.Cleanup(server.Close)
	client := &fasthttp.Client{}
	t.Cleanup(client.CloseIdleConnections)
	logger := libpack_logger.New().SetOutput(io.Discard)
	legacy := NewBackendHealthManager(client, server.URL, "", logger)
	t.Cleanup(legacy.Shutdown)
	require.False(t, legacy.checkBackendHealth(), "GraphQL must reject the unauthenticated probe")

	manager := NewBackendHealthManager(client, server.URL, server.URL+"/ready", logger)
	t.Cleanup(manager.Shutdown)
	manager.checkInterval = 10 * time.Millisecond
	manager.StartHealthChecking()
	require.NoError(t, manager.WaitForBackendReady(time.Second))
	require.True(t, manager.IsHealthy())

	readyStatus.Store(http.StatusServiceUnavailable)
	require.Eventually(t, func() bool { return !manager.IsHealthy() }, time.Second, 5*time.Millisecond)
	readyStatus.Store(http.StatusOK)
	require.Eventually(t, manager.IsHealthy, time.Second, 5*time.Millisecond)
	assert.Equal(t, int32(1), graphqlCalls.Load(), "readiness checks must not query protected GraphQL")
	require.False(t, legacy.checkBackendHealth(), "readiness must not bypass GraphQL authorization")
}

func TestBackendReadinessExactTarget(t *testing.T) {
	const target = "/checks//ready%2Fdatabase?probe=one%2Ftwo&probe=three"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, target, r.RequestURI)
		w.WriteHeader(http.StatusNoContent)
	}))
	t.Cleanup(server.Close)
	client := &fasthttp.Client{}
	t.Cleanup(client.CloseIdleConnections)
	manager := NewBackendHealthManager(client, "http://unused.invalid", server.URL+target, libpack_logger.New().SetOutput(io.Discard))
	t.Cleanup(manager.Shutdown)
	assert.True(t, manager.checkBackendHealth())
}

func TestBackendReadinessStatus(t *testing.T) {
	for _, tc := range []struct {
		name    string
		status  int
		healthy bool
	}{
		{name: "good OK", status: http.StatusOK, healthy: true},
		{name: "good no content", status: http.StatusNoContent, healthy: true},
		{name: "edge final success", status: 299, healthy: true},
		{name: "bad redirect", status: http.StatusFound},
		{name: "bad unauthorized", status: http.StatusUnauthorized},
		{name: "bad unavailable", status: http.StatusServiceUnavailable},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/ready" {
					t.Error("readiness probe followed a redirect")
					w.WriteHeader(http.StatusOK)
					return
				}
				w.Header().Set("Location", "/redirected")
				w.WriteHeader(tc.status)
			}))
			t.Cleanup(server.Close)
			client := &fasthttp.Client{}
			t.Cleanup(client.CloseIdleConnections)
			manager := NewBackendHealthManager(client, server.URL, server.URL+"/ready", libpack_logger.New().SetOutput(io.Discard))
			t.Cleanup(manager.Shutdown)
			assert.Equal(t, tc.healthy, manager.checkBackendHealth())
		})
	}
}

func TestBackendReadinessNetworkFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	server.Close()
	client := &fasthttp.Client{}
	t.Cleanup(client.CloseIdleConnections)
	manager := NewBackendHealthManager(client, server.URL, server.URL+"/ready", libpack_logger.New().SetOutput(io.Discard))
	t.Cleanup(manager.Shutdown)
	assert.False(t, manager.checkBackendHealth())
}

func TestBackendReadinessTimeout(t *testing.T) {
	release := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		<-release
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)
	t.Cleanup(func() { close(release) })
	client := &fasthttp.Client{ReadTimeout: 20 * time.Millisecond}
	t.Cleanup(client.CloseIdleConnections)
	manager := NewBackendHealthManager(client, server.URL, server.URL+"/ready", libpack_logger.New().SetOutput(io.Discard))
	t.Cleanup(manager.Shutdown)
	assert.False(t, manager.checkBackendHealth())
}

func TestBackendReadinessLegacyGraphQL(t *testing.T) {
	for _, tc := range []struct {
		name           string
		configuredPath string
		requestURI     string
	}{
		{name: "good bare host", requestURI: "/v1/graphql"},
		{name: "edge root path", configuredPath: "/", requestURI: "/v1/graphql"},
		{name: "good custom path and query", configuredPath: "/graphql?health=1", requestURI: "/graphql?health=1"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				assert.Equal(t, tc.requestURI, r.RequestURI)
				assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
				assert.Empty(t, r.Header.Get("Authorization"))
				body, err := io.ReadAll(r.Body)
				assert.NoError(t, err)
				assert.JSONEq(t, `{"query":"{__typename}"}`, string(body))
				w.WriteHeader(http.StatusOK)
			}))
			t.Cleanup(server.Close)
			client := &fasthttp.Client{}
			t.Cleanup(client.CloseIdleConnections)
			manager := NewBackendHealthManager(client, server.URL+tc.configuredPath, "", libpack_logger.New().SetOutput(io.Discard))
			t.Cleanup(manager.Shutdown)
			assert.True(t, manager.checkBackendHealth())
		})
	}
}

func readinessURLWithUserinfo(user *url.Userinfo) string {
	u := url.URL{Scheme: "http", Host: "example.test", Path: "/ready", User: user}
	return u.String()
}

func TestBackendReadinessURLValidation(t *testing.T) {
	for _, tc := range []struct {
		name  string
		url   string
		valid bool
	}{
		{name: "edge unset", valid: true},
		{name: "good HTTP", url: "http://127.0.0.1:8080/ready", valid: true},
		{name: "good HTTPS query", url: "https://example.test/check?database=primary", valid: true},
		{name: "good IPv6", url: "http://[::1]:8080/ready", valid: true},
		{name: "bad relative", url: "/ready"},
		{name: "bad scheme", url: "ftp://example.test/ready"},
		{name: "bad missing host", url: "http:///ready"},
		{name: "bad port only host", url: "http://:8080/ready"},
		{name: "bad opaque", url: "http:ready"},
		{name: "bad userinfo", url: readinessURLWithUserinfo(url.UserPassword("operator", "secret"))},
		{name: "bad user only", url: readinessURLWithUserinfo(url.User("operator"))},
		{name: "bad escaped userinfo", url: strings.Replace(readinessURLWithUserinfo(url.UserPassword("operator", "%zz")), "%25zz", "%zz", 1)},
		{name: "bad path escape", url: "http://example.test/%zz"},
		{name: "bad port", url: "http://example.test:invalid/ready"},
		{name: "bad port range", url: "http://example.test:65536/ready"},
		{name: "bad fragment", url: "http://example.test/ready#ignored"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := validateBackendHealthcheckURL(tc.url)
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				assert.NotContains(t, err.Error(), tc.url)
				assert.NotContains(t, err.Error(), "operator")
				assert.NotContains(t, err.Error(), "secret")
			}
		})
	}
}

func TestBackendReadinessInvalidConfigStopsStartup(t *testing.T) {
	const helper = "GMP_TEST_INVALID_READINESS_CONFIG"
	if os.Getenv(helper) == "1" {
		parseConfig()
		t.Fatal("invalid readiness configuration did not stop startup")
	}
	command := exec.Command(os.Args[0], "-test.run=^TestBackendReadinessInvalidConfigStopsStartup$")
	for _, variable := range os.Environ() {
		if !strings.HasPrefix(variable, "GMP_BACKEND_HEALTHCHECK_URL=") && !strings.HasPrefix(variable, helper+"=") {
			command.Env = append(command.Env, variable)
		}
	}
	command.Env = append(command.Env, helper+"=1", "GMP_BACKEND_HEALTHCHECK_URL="+readinessURLWithUserinfo(url.UserPassword("operator", "secret")))
	output, err := command.CombinedOutput()
	var exitError *exec.ExitError
	require.ErrorAs(t, err, &exitError)
	assert.Equal(t, 1, exitError.ExitCode())
	assert.Contains(t, string(output), "Invalid BACKEND_HEALTHCHECK_URL")
	assert.NotContains(t, string(output), "operator")
	assert.NotContains(t, string(output), "secret")
}
