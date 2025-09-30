package main

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProxyError(t *testing.T) {
	tests := []struct {
		name         string
		code         string
		message      string
		statusCode   int
		retryable    bool
		expectStatus int
	}{
		{
			name:         "connection refused error",
			code:         ErrCodeConnectionRefused,
			message:      "backend unavailable",
			statusCode:   http.StatusServiceUnavailable,
			retryable:    true,
			expectStatus: http.StatusServiceUnavailable,
		},
		{
			name:         "timeout error",
			code:         ErrCodeTimeout,
			message:      "request timeout",
			statusCode:   http.StatusGatewayTimeout,
			retryable:    true,
			expectStatus: http.StatusGatewayTimeout,
		},
		{
			name:         "circuit breaker open",
			code:         ErrCodeCircuitOpen,
			message:      "circuit breaker open",
			statusCode:   http.StatusServiceUnavailable,
			retryable:    false,
			expectStatus: http.StatusServiceUnavailable,
		},
		{
			name:         "rate limit exceeded",
			code:         ErrCodeRateLimited,
			message:      "too many requests",
			statusCode:   http.StatusTooManyRequests,
			retryable:    false,
			expectStatus: http.StatusTooManyRequests,
		},
		{
			name:         "service unavailable",
			code:         ErrCodeServiceUnavailable,
			message:      "no retry tokens available",
			statusCode:   http.StatusServiceUnavailable,
			retryable:    false,
			expectStatus: http.StatusServiceUnavailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewProxyError(tt.code, tt.message, tt.statusCode, tt.retryable)

			assert.NotNil(t, err)
			assert.Equal(t, tt.code, err.Code)
			assert.Equal(t, tt.message, err.Message)
			assert.Equal(t, tt.retryable, err.Retryable)
			assert.Equal(t, tt.expectStatus, err.StatusCode)
			assert.NotEmpty(t, err.Timestamp)
			assert.NotNil(t, err.Metadata)
		})
	}
}

func TestProxyError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *ProxyError
		expected string
	}{
		{
			name: "error with details",
			err: NewProxyError(ErrCodeConnectionRefused, "backend unavailable", http.StatusServiceUnavailable, true).
				WithDetails("connection refused"),
			expected: "CONNECTION_REFUSED: backend unavailable (connection refused)",
		},
		{
			name:     "error without details",
			err:      NewProxyError(ErrCodeCircuitOpen, "circuit breaker open", http.StatusServiceUnavailable, false),
			expected: "CIRCUIT_OPEN: circuit breaker open",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.err.Error())
		})
	}
}

func TestProxyError_Unwrap(t *testing.T) {
	cause := errors.New("original error")
	err := NewProxyError(ErrCodeTimeout, "timeout occurred", http.StatusGatewayTimeout, true).WithCause(cause)

	unwrapped := errors.Unwrap(err)
	assert.Equal(t, cause, unwrapped)
}

func TestProxyError_WithMethods(t *testing.T) {
	t.Run("with details", func(t *testing.T) {
		err := NewProxyError(ErrCodeTimeout, "timeout", http.StatusGatewayTimeout, true).
			WithDetails("operation timed out")

		assert.Equal(t, "operation timed out", err.Details)
	})

	t.Run("with cause", func(t *testing.T) {
		cause := errors.New("original error")
		err := NewProxyError(ErrCodeTimeout, "timeout", http.StatusGatewayTimeout, true).
			WithCause(cause)

		assert.Equal(t, cause, err.Cause)
	})

	t.Run("with trace ID", func(t *testing.T) {
		err := NewProxyError(ErrCodeTimeout, "timeout", http.StatusGatewayTimeout, true).
			WithTraceID("trace-123")

		assert.Equal(t, "trace-123", err.TraceID)
	})

	t.Run("with metadata", func(t *testing.T) {
		err := NewProxyError(ErrCodeTimeout, "timeout", http.StatusGatewayTimeout, true).
			WithMetadata("attempt", 3).
			WithMetadata("endpoint", "/graphql")

		assert.Equal(t, 3, err.Metadata["attempt"])
		assert.Equal(t, "/graphql", err.Metadata["endpoint"])
	})
}

func TestProxyError_MarshalJSON(t *testing.T) {
	cause := errors.New("connection refused")
	err := NewProxyError(ErrCodeConnectionRefused, "backend unavailable", http.StatusServiceUnavailable, true).
		WithDetails("network error").
		WithCause(cause).
		WithTraceID("trace-456")

	data, jsonErr := err.MarshalJSON()
	assert.NoError(t, jsonErr)
	assert.NotEmpty(t, data)
	assert.Contains(t, string(data), "CONNECTION_REFUSED")
	assert.Contains(t, string(data), "backend unavailable")
	assert.Contains(t, string(data), "connection refused")
}

func TestErrorCodes(t *testing.T) {
	// Verify all error codes are defined
	codes := []string{
		ErrCodeConnectionRefused,
		ErrCodeConnectionReset,
		ErrCodeTimeout,
		ErrCodeCircuitOpen,
		ErrCodeRateLimited,
		ErrCodeInvalidRequest,
		ErrCodeBackendError,
		ErrCodeInternalError,
		ErrCodeUnauthorized,
		ErrCodeForbidden,
		ErrCodeNotFound,
		ErrCodeServiceUnavailable,
		ErrCodeBadGateway,
		ErrCodeInvalidResponse,
		ErrCodeQueryTooComplex,
		ErrCodeCacheFailed,
		ErrCodeContextCanceled,
	}

	for _, code := range codes {
		assert.NotEmpty(t, code, "Error code should not be empty")
	}

	// Verify codes are unique
	codeMap := make(map[string]bool)
	for _, code := range codes {
		assert.False(t, codeMap[code], "Error code %s should be unique", code)
		codeMap[code] = true
	}
}

func TestProxyError_ChainableMethods(t *testing.T) {
	// Test that methods can be chained
	err := NewProxyError(ErrCodeTimeout, "timeout", http.StatusGatewayTimeout, true).
		WithDetails("operation timeout").
		WithCause(errors.New("deadline exceeded")).
		WithTraceID("trace-789").
		WithMetadata("attempt", 1).
		WithMetadata("duration_ms", 5000)

	assert.Equal(t, "operation timeout", err.Details)
	assert.NotNil(t, err.Cause)
	assert.Equal(t, "trace-789", err.TraceID)
	assert.Equal(t, 1, err.Metadata["attempt"])
	assert.Equal(t, 5000, err.Metadata["duration_ms"])
}

func TestProxyError_Retryable(t *testing.T) {
	tests := []struct {
		name      string
		code      string
		retryable bool
	}{
		{
			name:      "timeout is retryable",
			code:      ErrCodeTimeout,
			retryable: true,
		},
		{
			name:      "connection refused is retryable",
			code:      ErrCodeConnectionRefused,
			retryable: true,
		},
		{
			name:      "rate limited is not retryable",
			code:      ErrCodeRateLimited,
			retryable: false,
		},
		{
			name:      "circuit open is not retryable",
			code:      ErrCodeCircuitOpen,
			retryable: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewProxyError(tt.code, "test error", http.StatusInternalServerError, tt.retryable)
			assert.Equal(t, tt.retryable, err.Retryable)
		})
	}
}
