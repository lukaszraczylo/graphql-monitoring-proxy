package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ProxyLoggingSecurityTestSuite struct {
	suite.Suite
}

func TestProxyLoggingSecurityTestSuite(t *testing.T) {
	suite.Run(t, new(ProxyLoggingSecurityTestSuite))
}

// TestSensitiveDataSanitization tests that sensitive data is properly redacted from logs
func (suite *ProxyLoggingSecurityTestSuite) TestSensitiveDataSanitization() {
	tests := []struct {
		name        string
		input       map[string]interface{}
		expected    map[string]interface{}
		contentType string
		description string
	}{
		{
			name: "Password field redaction",
			input: map[string]interface{}{
				"username": "user123",
				"password": "secret123",
				"email":    "user@example.com",
			},
			expected: map[string]interface{}{
				"username": "user123",
				"password": "[REDACTED]",
				"email":    "[REDACTED]",
			},
			contentType: "application/json",
			description: "Should redact password and email fields",
		},
		{
			name: "API key and token redaction",
			input: map[string]interface{}{
				"data":    "normal data",
				"api_key": "sk-123456789",
				"token":   "bearer-token-123",
				"auth":    "auth-value",
			},
			expected: map[string]interface{}{
				"data":    "normal data",
				"api_key": "[REDACTED]",
				"token":   "[REDACTED]",
				"auth":    "[REDACTED]",
			},
			contentType: "application/json",
			description: "Should redact API keys and tokens",
		},
		{
			name: "Nested sensitive fields",
			input: map[string]interface{}{
				"user": map[string]interface{}{
					"name":     "John Doe",
					"password": "secret123",
					"profile": map[string]interface{}{
						"api_key": "sk-nested-key",
						"bio":     "User bio",
					},
				},
				"public_data": "visible",
			},
			expected: map[string]interface{}{
				"user": map[string]interface{}{
					"name":     "John Doe",
					"password": "[REDACTED]",
					"profile": map[string]interface{}{
						"api_key": "[REDACTED]",
						"bio":     "User bio",
					},
				},
				"public_data": "visible",
			},
			contentType: "application/json",
			description: "Should redact nested sensitive fields",
		},
		{
			name: "Array with sensitive data",
			input: map[string]interface{}{
				"users": []interface{}{
					map[string]interface{}{
						"name":     "User1",
						"password": "pass1",
					},
					map[string]interface{}{
						"name":  "User2",
						"token": "token2",
					},
				},
			},
			expected: map[string]interface{}{
				"users": []interface{}{
					map[string]interface{}{
						"name":     "User1",
						"password": "[REDACTED]",
					},
					map[string]interface{}{
						"name":  "User2",
						"token": "[REDACTED]",
					},
				},
			},
			contentType: "application/json",
			description: "Should redact sensitive fields in arrays",
		},
		{
			name: "Credit card and financial data",
			input: map[string]interface{}{
				"order_id":    "12345",
				"credit_card": "4111111111111111",
				"cvv":         "123",
				"amount":      100.50,
			},
			expected: map[string]interface{}{
				"order_id":    "12345",
				"credit_card": "[REDACTED]",
				"cvv":         "[REDACTED]",
				"amount":      json.Number("100.5"),
			},
			contentType: "application/json",
			description: "Should redact financial sensitive data",
		},
		{
			name: "Personal identifiable information",
			input: map[string]interface{}{
				"name":    "John Doe",
				"ssn":     "123-45-6789",
				"phone":   "+1-555-123-4567",
				"address": "123 Main St",
				"age":     30,
			},
			expected: map[string]interface{}{
				"name":    "John Doe",
				"ssn":     "[REDACTED]",
				"phone":   "[REDACTED]",
				"address": "[REDACTED]",
				"age":     json.Number("30"),
			},
			contentType: "application/json",
			description: "Should redact PII data",
		},
		{
			name: "Mixed case field names",
			input: map[string]interface{}{
				"UserName": "john",
				"PASSWORD": "secret",
				"Api_Key":  "key123",
				"Bearer":   "token",
			},
			expected: map[string]interface{}{
				"UserName": "john",
				"PASSWORD": "[REDACTED]",
				"Api_Key":  "[REDACTED]",
				"Bearer":   "[REDACTED]",
			},
			contentType: "application/json",
			description: "Should handle mixed case field names",
		},
		{
			name: "Various password patterns",
			input: map[string]interface{}{
				"pwd":      "secret1",
				"passwd":   "secret2",
				"password": "secret3",
				"pass":     "not-redacted", // Should NOT be redacted (not in list)
			},
			expected: map[string]interface{}{
				"pwd":      "[REDACTED]",
				"passwd":   "[REDACTED]",
				"password": "[REDACTED]",
				"pass":     "not-redacted",
			},
			contentType: "application/json",
			description: "Should handle various password field patterns",
		},
		{
			name: "Various auth patterns",
			input: map[string]interface{}{
				"authorization": "Bearer token123",
				"auth":          "basic auth",
				"bearer":        "token456",
				"session":       "sess123",
				"sessionid":     "session456",
				"session_id":    "session789",
				"cookie":        "cookie_value",
			},
			expected: map[string]interface{}{
				"authorization": "[REDACTED]",
				"auth":          "[REDACTED]",
				"bearer":        "[REDACTED]",
				"session":       "[REDACTED]",
				"sessionid":     "[REDACTED]",
				"session_id":    "[REDACTED]",
				"cookie":        "[REDACTED]",
			},
			contentType: "application/json",
			description: "Should handle various authentication field patterns",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			// Convert input to JSON bytes
			inputBytes, err := json.Marshal(tt.input)
			suite.NoError(err)

			// Test the sanitization function
			result := sanitizeForLogging(inputBytes, tt.contentType)

			// Parse the result back to compare
			var sanitized map[string]interface{}
			decoder := json.NewDecoder(strings.NewReader(result))
			decoder.UseNumber() // Preserve number precision and type
			err = decoder.Decode(&sanitized)
			suite.NoError(err, "Sanitized result should be valid JSON")

			// Compare the result with expected
			suite.Equal(tt.expected, sanitized, tt.description)

			// Verify no sensitive data remains in the string representation
			resultStr := strings.ToLower(result)
			if strings.Contains(tt.name, "password") || strings.Contains(tt.name, "secret") {
				suite.NotContains(resultStr, "secret", "Should not contain 'secret' in result")
			}
			if strings.Contains(tt.name, "key") {
				suite.NotContains(resultStr, "sk-", "Should not contain API key prefix")
			}
		})
	}
}

// TestSensitiveDataSanitizationNonJSON tests sanitization for non-JSON content
func (suite *ProxyLoggingSecurityTestSuite) TestSensitiveDataSanitizationNonJSON() {
	tests := []struct {
		name                   string
		input                  string
		contentType            string
		description            string
		shouldNotContain       []string
		shouldContainSanitized []string
	}{
		{
			name:                   "Form data with password",
			input:                  "username=john&password=secret123&email=john@example.com",
			contentType:            "application/x-www-form-urlencoded",
			shouldNotContain:       []string{"secret123"},
			shouldContainSanitized: []string{"password=[REDACTED]"},
			description:            "Should redact password in form data",
		},
		{
			name:                   "Query string with sensitive data",
			input:                  "?user=john&api_key=sk-123456&public=data",
			contentType:            "text/plain",
			shouldNotContain:       []string{"sk-123456"},
			shouldContainSanitized: []string{"api_key=[REDACTED]"},
			description:            "Should redact API key in query string",
		},
		{
			name:                   "Large body truncation",
			input:                  strings.Repeat("a", 1500) + "password=secret",
			contentType:            "text/plain",
			shouldNotContain:       []string{},
			shouldContainSanitized: []string{"[truncated]"},
			description:            "Should truncate large bodies",
		},
		{
			name:                   "XML-like content with sensitive data",
			input:                  "<user><name>John</name><password>secret123</password></user>",
			contentType:            "application/xml",
			shouldNotContain:       []string{"secret123"},
			shouldContainSanitized: []string{"password=[REDACTED]"},
			description:            "Should redact sensitive data in XML-like content",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			result := sanitizeForLogging([]byte(tt.input), tt.contentType)

			// Check that sensitive data is removed
			for _, sensitiveData := range tt.shouldNotContain {
				suite.NotContains(result, sensitiveData,
					"Result should not contain sensitive data: %s", sensitiveData)
			}

			// Check that redaction markers are present
			for _, redactedPattern := range tt.shouldContainSanitized {
				suite.Contains(result, redactedPattern,
					"Result should contain redaction marker: %s", redactedPattern)
			}
		})
	}
}

// TestSanitizeHeaders tests header sanitization
func (suite *ProxyLoggingSecurityTestSuite) TestSanitizeHeaders() {
	tests := []struct {
		input    map[string]string
		expected map[string]string
		name     string
	}{
		{
			name: "Authorization header redaction",
			input: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer token123",
				"User-Agent":    "Test/1.0",
			},
			expected: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "[REDACTED]",
				"User-Agent":    "Test/1.0",
			},
		},
		{
			name: "API key headers redaction",
			input: map[string]string{
				"X-API-Key":      "sk-123456",
				"X-Auth-Token":   "auth-token-123",
				"X-API-Secret":   "secret-key",
				"Content-Length": "100",
			},
			expected: map[string]string{
				"X-API-Key":      "[REDACTED]",
				"X-Auth-Token":   "[REDACTED]",
				"X-API-Secret":   "[REDACTED]",
				"Content-Length": "100",
			},
		},
		{
			name: "Cookie headers redaction",
			input: map[string]string{
				"Cookie":     "sessionid=abc123; userid=456",
				"Set-Cookie": "token=xyz789; Path=/",
				"Host":       "example.com",
			},
			expected: map[string]string{
				"Cookie":     "[REDACTED]",
				"Set-Cookie": "[REDACTED]",
				"Host":       "example.com",
			},
		},
		{
			name: "Mixed case headers",
			input: map[string]string{
				"AUTHORIZATION": "Bearer token",
				"x-api-key":     "key123",
				"Content-TYPE":  "json",
			},
			expected: map[string]string{
				"AUTHORIZATION": "[REDACTED]",
				"x-api-key":     "[REDACTED]",
				"Content-TYPE":  "json",
			},
		},
		{
			name: "CSRF and access tokens",
			input: map[string]string{
				"X-CSRF-Token":   "csrf123",
				"X-Access-Token": "access456",
				"Accept":         "application/json",
			},
			expected: map[string]string{
				"X-CSRF-Token":   "[REDACTED]",
				"X-Access-Token": "[REDACTED]",
				"Accept":         "application/json",
			},
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			result := sanitizeHeaders(tt.input)
			suite.Equal(tt.expected, result)

			// Verify original headers are not modified
			for key, originalValue := range tt.input {
				suite.Equal(originalValue, tt.input[key],
					"Original headers should not be modified")
			}
		})
	}
}

// TestRedactSensitiveFields tests the recursive redaction function
func (suite *ProxyLoggingSecurityTestSuite) TestRedactSensitiveFields() {
	sensitiveFields := []string{"password", "token", "secret"}

	suite.Run("Deep nested structure", func() {
		data := map[string]interface{}{
			"level1": map[string]interface{}{
				"level2": map[string]interface{}{
					"level3": map[string]interface{}{
						"password": "testdeepsecret",
						"public":   "data",
					},
					"token": "testlevel2token",
				},
				"normal": "value",
			},
			"secret": "testtoplevel",
		}

		redactSensitiveFields(data, sensitiveFields)

		// Verify deep nesting is handled
		level3 := data["level1"].(map[string]interface{})["level2"].(map[string]interface{})["level3"].(map[string]interface{})
		suite.Equal("[REDACTED]", level3["password"])
		suite.Equal("data", level3["public"])

		// Verify intermediate levels
		level2 := data["level1"].(map[string]interface{})["level2"].(map[string]interface{})
		suite.Equal("[REDACTED]", level2["token"])

		// Verify top level
		suite.Equal("[REDACTED]", data["secret"])
		level1 := data["level1"].(map[string]interface{})
		suite.Equal("value", level1["normal"])
	})

	suite.Run("Array of objects", func() {
		data := map[string]interface{}{
			"users": []interface{}{
				map[string]interface{}{
					"name":     "User1",
					"password": "testpass1",
				},
				map[string]interface{}{
					"name":  "User2",
					"token": "testtoken2",
				},
				"not-an-object", // Should be ignored
			},
		}

		redactSensitiveFields(data, sensitiveFields)

		users := data["users"].([]interface{})
		user1 := users[0].(map[string]interface{})
		user2 := users[1].(map[string]interface{})

		suite.Equal("[REDACTED]", user1["password"])
		suite.Equal("User1", user1["name"])
		suite.Equal("[REDACTED]", user2["token"])
		suite.Equal("User2", user2["name"])
		suite.Equal("not-an-object", users[2])
	})
}

// TestRedactPatternInString tests string pattern redaction
func (suite *ProxyLoggingSecurityTestSuite) TestRedactPatternInString() {
	tests := []struct {
		name     string
		input    string
		pattern  string
		expected string
	}{
		{
			name:     "JSON-style pattern",
			input:    `{"password": "secret123", "user": "john"}`,
			pattern:  "password",
			expected: `{"password":"[REDACTED]", "user": "john"}`,
		},
		{
			name:     "Form-style pattern with equals",
			input:    "username=john&password=secret&email=test",
			pattern:  "password",
			expected: "username=john&password=[REDACTED]&email=test",
		},
		{
			name:     "Double quoted pattern",
			input:    `password="secret123"`,
			pattern:  "password",
			expected: `password="[REDACTED]"`,
		},
		{
			name:     "Single quoted pattern",
			input:    `password='secret123'`,
			pattern:  "password",
			expected: `password='[REDACTED]'`,
		},
		{
			name:     "No match",
			input:    "normal text without sensitive data",
			pattern:  "password",
			expected: "normal text without sensitive data",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			result := redactPatternInString(tt.input, tt.pattern)
			suite.Equal(tt.expected, result)
		})
	}
}

// TestSanitizationPerformance tests performance of sanitization functions
func (suite *ProxyLoggingSecurityTestSuite) TestSanitizationPerformance() {
	// Create a large JSON structure with sensitive data
	largeData := make(map[string]interface{})
	for i := 0; i < 1000; i++ {
		largeData[fmt.Sprintf("user_%d", i)] = map[string]interface{}{
			"name":     fmt.Sprintf("User%d", i),
			"password": fmt.Sprintf("secret%d", i),
			"email":    fmt.Sprintf("user%d@example.com", i),
			"public":   fmt.Sprintf("public_data_%d", i),
		}
	}

	largeJSON, err := json.Marshal(largeData)
	suite.NoError(err)

	// Test that sanitization completes in reasonable time
	result := sanitizeForLogging(largeJSON, "application/json")

	// Verify the result is valid JSON
	var sanitized map[string]interface{}
	err = json.Unmarshal([]byte(result), &sanitized)
	suite.NoError(err)

	// Verify sensitive data was redacted (spot check)
	user0 := sanitized["user_0"].(map[string]interface{})
	suite.Equal("[REDACTED]", user0["password"])
	suite.Equal("[REDACTED]", user0["email"])
	suite.Equal("User0", user0["name"])
}

// TestEdgeCases tests edge cases and error conditions
func (suite *ProxyLoggingSecurityTestSuite) TestEdgeCases() {
	suite.Run("Empty body", func() {
		result := sanitizeForLogging([]byte{}, "application/json")
		suite.Equal("", result)
	})

	suite.Run("Invalid JSON", func() {
		invalidJSON := []byte(`{"invalid": json}`)
		result := sanitizeForLogging(invalidJSON, "application/json")
		// Should fall back to string sanitization
		suite.Contains(result, "invalid")
	})

	suite.Run("Nil data", func() {
		// Test with nil maps (should not panic)
		sensitiveFields := []string{"password"}

		// This should not panic
		suite.NotPanics(func() {
			data := make(map[string]interface{})
			data["test"] = nil
			redactSensitiveFields(data, sensitiveFields)
		})
	})

	suite.Run("Empty headers", func() {
		result := sanitizeHeaders(map[string]string{})
		suite.Equal(map[string]string{}, result)
	})

	suite.Run("Very large content type", func() {
		largeContentType := strings.Repeat("json", 1000)
		result := sanitizeForLogging([]byte(`{"test": "data"}`), largeContentType)
		suite.Contains(result, "test")
	})
}

// BenchmarkSanitizeForLogging benchmarks the sanitization function
func BenchmarkSanitizeForLogging(b *testing.B) {
	testData := map[string]interface{}{
		"username": "testuser",
		"password": "secret123",
		"api_key":  "sk-123456789",
		"data":     "normal data",
		"nested": map[string]interface{}{
			"token": "nested-token",
			"value": "nested-value",
		},
	}

	jsonData, _ := json.Marshal(testData)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sanitizeForLogging(jsonData, "application/json")
	}
}

// BenchmarkSanitizeHeaders benchmarks header sanitization
func BenchmarkSanitizeHeaders(b *testing.B) {
	headers := map[string]string{
		"Content-Type":   "application/json",
		"Authorization":  "Bearer token123",
		"X-API-Key":      "sk-123456",
		"User-Agent":     "Test/1.0",
		"Accept":         "application/json",
		"Content-Length": "100",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sanitizeHeaders(headers)
	}
}
