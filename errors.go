package main

import (
	"encoding/json"
	"fmt"
	"time"
)

// Error codes for structured error responses
const (
	ErrCodeConnectionRefused  = "CONNECTION_REFUSED"
	ErrCodeConnectionReset    = "CONNECTION_RESET"
	ErrCodeTimeout            = "TIMEOUT"
	ErrCodeCircuitOpen        = "CIRCUIT_OPEN"
	ErrCodeRateLimited        = "RATE_LIMITED"
	ErrCodeInvalidRequest     = "INVALID_REQUEST"
	ErrCodeBackendError       = "BACKEND_ERROR"
	ErrCodeInternalError      = "INTERNAL_ERROR"
	ErrCodeUnauthorized       = "UNAUTHORIZED"
	ErrCodeForbidden          = "FORBIDDEN"
	ErrCodeNotFound           = "NOT_FOUND"
	ErrCodeServiceUnavailable = "SERVICE_UNAVAILABLE"
	ErrCodeBadGateway         = "BAD_GATEWAY"
	ErrCodeInvalidResponse    = "INVALID_RESPONSE"
	ErrCodeQueryTooComplex    = "QUERY_TOO_COMPLEX"
	ErrCodeCacheFailed        = "CACHE_FAILED"
	ErrCodeContextCanceled    = "CONTEXT_CANCELED"
)

// ProxyError represents a structured error response
type ProxyError struct {
	Code       string                 `json:"code"`               // Machine-readable error code
	Message    string                 `json:"message"`            // Human-readable error message
	Details    string                 `json:"details,omitempty"`  // Additional error details
	Retryable  bool                   `json:"retryable"`          // Whether the request can be retried
	StatusCode int                    `json:"status_code"`        // HTTP status code
	Timestamp  time.Time              `json:"timestamp"`          // When the error occurred
	TraceID    string                 `json:"trace_id,omitempty"` // Trace ID for correlation
	Metadata   map[string]interface{} `json:"metadata,omitempty"` // Additional context
	Cause      error                  `json:"-"`                  // Original error (not serialized)
}

// Error implements the error interface
func (e *ProxyError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s (%s)", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the underlying error
func (e *ProxyError) Unwrap() error {
	return e.Cause
}

// MarshalJSON implements custom JSON marshaling
func (e *ProxyError) MarshalJSON() ([]byte, error) {
	type Alias ProxyError
	return json.Marshal(&struct {
		*Alias
		CauseMessage string `json:"cause,omitempty"`
	}{
		Alias: (*Alias)(e),
		CauseMessage: func() string {
			if e.Cause != nil {
				return e.Cause.Error()
			}
			return ""
		}(),
	})
}

// NewProxyError creates a new structured error
func NewProxyError(code, message string, statusCode int, retryable bool) *ProxyError {
	return &ProxyError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
		Retryable:  retryable,
		Timestamp:  time.Now(),
		Metadata:   make(map[string]interface{}),
	}
}

// WithDetails adds details to the error
func (e *ProxyError) WithDetails(details string) *ProxyError {
	e.Details = details
	return e
}

// WithCause adds the underlying cause
func (e *ProxyError) WithCause(cause error) *ProxyError {
	e.Cause = cause
	return e
}

// WithTraceID adds a trace ID
func (e *ProxyError) WithTraceID(traceID string) *ProxyError {
	e.TraceID = traceID
	return e
}

// WithMetadata adds metadata
func (e *ProxyError) WithMetadata(key string, value interface{}) *ProxyError {
	e.Metadata[key] = value
	return e
}

// Common error constructors

// NewConnectionError creates a connection-related error
func NewConnectionError(err error) *ProxyError {
	code := ErrCodeConnectionRefused
	if err != nil {
		errStr := err.Error()
		if contains(errStr, "reset") {
			code = ErrCodeConnectionReset
		}
	}

	return NewProxyError(code, "Failed to connect to backend", 502, true).
		WithCause(err)
}

// NewTimeoutError creates a timeout error
func NewTimeoutError(err error) *ProxyError {
	return NewProxyError(ErrCodeTimeout, "Request timed out", 504, false).
		WithCause(err)
}

// NewCircuitOpenError creates a circuit breaker open error
func NewCircuitOpenError() *ProxyError {
	return NewProxyError(ErrCodeCircuitOpen, "Service temporarily unavailable due to circuit breaker", 503, false).
		WithDetails("The backend service is currently experiencing issues. Please try again later.")
}

// NewRateLimitError creates a rate limit error
func NewRateLimitError(userID, role string) *ProxyError {
	return NewProxyError(ErrCodeRateLimited, "Rate limit exceeded", 429, false).
		WithDetails("You have exceeded the rate limit for your role").
		WithMetadata("user_id", userID).
		WithMetadata("role", role)
}

// NewBackendError creates a backend error from status code
func NewBackendError(statusCode int, body string) *ProxyError {
	code := ErrCodeBackendError
	message := "Backend returned an error"
	retryable := false

	switch {
	case statusCode == 429:
		code = ErrCodeRateLimited
		message = "Backend rate limit exceeded"
		retryable = true
	case statusCode == 503:
		code = ErrCodeServiceUnavailable
		message = "Backend service unavailable"
		retryable = true
	case statusCode == 502 || statusCode == 504:
		code = ErrCodeBadGateway
		message = "Bad gateway"
		retryable = true
	case statusCode >= 500:
		code = ErrCodeBackendError
		message = "Backend server error"
		retryable = true
	case statusCode == 404:
		code = ErrCodeNotFound
		message = "Resource not found"
	case statusCode == 403:
		code = ErrCodeForbidden
		message = "Access forbidden"
	case statusCode == 401:
		code = ErrCodeUnauthorized
		message = "Unauthorized"
	case statusCode >= 400:
		code = ErrCodeInvalidRequest
		message = "Invalid request"
	}

	return NewProxyError(code, message, statusCode, retryable).
		WithMetadata("backend_status", statusCode).
		WithMetadata("backend_body", truncateString(body, 500))
}

// NewInvalidResponseError creates an invalid response error
func NewInvalidResponseError(details string) *ProxyError {
	return NewProxyError(ErrCodeInvalidResponse, "Backend returned invalid response", 502, false).
		WithDetails(details)
}

// NewInternalError creates an internal error
func NewInternalError(err error) *ProxyError {
	return NewProxyError(ErrCodeInternalError, "Internal proxy error", 500, false).
		WithCause(err)
}

// NewContextCanceledError creates a context canceled error
func NewContextCanceledError() *ProxyError {
	return NewProxyError(ErrCodeContextCanceled, "Request canceled", 499, false).
		WithDetails("The request was canceled by the client")
}

// Helper functions

func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// IsRetryable checks if an error is retryable
func IsRetryable(err error) bool {
	if err == nil {
		return false
	}

	if proxyErr, ok := err.(*ProxyError); ok {
		return proxyErr.Retryable
	}

	return false
}

// GetStatusCode extracts the status code from an error
func GetStatusCode(err error) int {
	if err == nil {
		return 200
	}

	if proxyErr, ok := err.(*ProxyError); ok {
		return proxyErr.StatusCode
	}

	return 500
}
