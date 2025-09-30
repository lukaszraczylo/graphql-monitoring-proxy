package main

import (
	"fmt"
	"strings"
)

// RateLimitConfigError represents a detailed error when loading rate limit configuration
type RateLimitConfigError struct {
	PathErrors map[string]string
	Paths      []string
}

// Error implements the error interface
func (e *RateLimitConfigError) Error() string {
	sb := strings.Builder{}
	sb.WriteString("Failed to load rate limit configuration. Please ensure a valid configuration file exists at one of these locations:\n")

	for _, path := range e.Paths {
		errMsg := e.PathErrors[path]
		sb.WriteString(fmt.Sprintf("  - %s: %s\n", path, errMsg))
	}

	sb.WriteString("\nTo resolve this issue:\n")
	sb.WriteString("1. Create a valid JSON file using the following template:\n")
	sb.WriteString(`   {
	    "ratelimit": {
	      "admin": {
	        "req": 100,
	        "interval": "second"
	      },
	      "guest": {
	        "req": 3,
	        "interval": "second"
	      },
	      "-": {
	        "req": 10,
	        "interval": "minute"
	      }
	    }
	  }`)
	sb.WriteString("\n\nThe 'interval' field supports the following formats:\n")
	sb.WriteString("  - String values: \"second\", \"minute\", \"hour\", \"day\"\n")
	sb.WriteString("  - Go duration strings: \"5s\", \"10m\", \"1h\"\n")
	sb.WriteString("  - Numeric values (in seconds): 60, 3600\n")
	sb.WriteString("\n2. Save it as 'ratelimit.json' in the current directory or in '/go/src/app/' (in Docker)\n")
	sb.WriteString("3. Ensure the file has correct permissions and is accessible by the service\n")

	return sb.String()
}

// NewRateLimitConfigError creates a new rate limit configuration error
func NewRateLimitConfigError(paths []string) *RateLimitConfigError {
	return &RateLimitConfigError{
		Paths:      paths,
		PathErrors: make(map[string]string),
	}
}
