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

// Helper functions

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
