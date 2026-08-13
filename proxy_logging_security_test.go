package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"unicode/utf8"

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
		input       map[string]any
		expected    map[string]any
		contentType string
		description string
	}{
		{
			name: "Password field redaction",
			input: map[string]any{
				"username": "user123",
				"password": "secret123",
				"email":    "user@example.com",
			},
			expected: map[string]any{
				"username": "user123",
				"password": "[REDACTED]",
				"email":    "[REDACTED]",
			},
			contentType: "application/json",
			description: "Should redact password and email fields",
		},
		{
			name: "API key and token redaction",
			input: map[string]any{
				"data":    "normal data",
				"api_key": "sk-123456789",
				"token":   "bearer-token-123",
				"auth":    "auth-value",
			},
			expected: map[string]any{
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
			input: map[string]any{
				"user": map[string]any{
					"name":     "John Doe",
					"password": "secret123",
					"profile": map[string]any{
						"api_key": "sk-nested-key",
						"bio":     "User bio",
					},
				},
				"public_data": "visible",
			},
			expected: map[string]any{
				"user": map[string]any{
					"name":     "John Doe",
					"password": "[REDACTED]",
					"profile": map[string]any{
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
			input: map[string]any{
				"users": []any{
					map[string]any{
						"name":     "User1",
						"password": "pass1",
					},
					map[string]any{
						"name":  "User2",
						"token": "token2",
					},
				},
			},
			expected: map[string]any{
				"users": []any{
					map[string]any{
						"name":     "User1",
						"password": "[REDACTED]",
					},
					map[string]any{
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
			input: map[string]any{
				"order_id":    "12345",
				"credit_card": "4111111111111111",
				"cvv":         "123",
				"amount":      100.50,
			},
			expected: map[string]any{
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
			input: map[string]any{
				"name":    "John Doe",
				"ssn":     "123-45-6789",
				"phone":   "+1-555-123-4567",
				"address": "123 Main St",
				"age":     30,
			},
			expected: map[string]any{
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
			input: map[string]any{
				"UserName": "john",
				"PASSWORD": "secret",
				"Api_Key":  "key123",
				"Bearer":   "token",
			},
			expected: map[string]any{
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
			input: map[string]any{
				"pwd":      "secret1",
				"passwd":   "secret2",
				"password": "secret3",
				"pass":     "secret4", // Now redacted for better security coverage
			},
			expected: map[string]any{
				"pwd":      "[REDACTED]",
				"passwd":   "[REDACTED]",
				"password": "[REDACTED]",
				"pass":     "[REDACTED]",
			},
			contentType: "application/json",
			description: "Should handle various password field patterns",
		},
		{
			name: "Various auth patterns",
			input: map[string]any{
				"authorization": "Bearer token123",
				"auth":          "basic auth",
				"bearer":        "token456",
				"session":       "sess123",
				"sessionid":     "session456",
				"session_id":    "session789",
				"cookie":        "cookie_value",
			},
			expected: map[string]any{
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
			var sanitized map[string]any
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
		data := map[string]any{
			"level1": map[string]any{
				"level2": map[string]any{
					"level3": map[string]any{
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
		level3 := data["level1"].(map[string]any)["level2"].(map[string]any)["level3"].(map[string]any)
		suite.Equal("[REDACTED]", level3["password"])
		suite.Equal("data", level3["public"])

		// Verify intermediate levels
		level2 := data["level1"].(map[string]any)["level2"].(map[string]any)
		suite.Equal("[REDACTED]", level2["token"])

		// Verify top level
		suite.Equal("[REDACTED]", data["secret"])
		level1 := data["level1"].(map[string]any)
		suite.Equal("value", level1["normal"])
	})

	suite.Run("Array of objects", func() {
		data := map[string]any{
			"users": []any{
				map[string]any{
					"name":     "User1",
					"password": "testpass1",
				},
				map[string]any{
					"name":  "User2",
					"token": "testtoken2",
				},
				"not-an-object", // Should be ignored
			},
		}

		redactSensitiveFields(data, sensitiveFields)

		users := data["users"].([]any)
		user1 := users[0].(map[string]any)
		user2 := users[1].(map[string]any)

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
	largeData := make(map[string]any)
	for i := 0; i < 1000; i++ {
		largeData[fmt.Sprintf("user_%d", i)] = map[string]any{
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
	var sanitized map[string]any
	err = json.Unmarshal([]byte(result), &sanitized)
	suite.NoError(err)

	// Verify sensitive data was redacted (spot check)
	user0 := sanitized["user_0"].(map[string]any)
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
			data := make(map[string]any)
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
	testData := map[string]any{
		"username": "testuser",
		"password": "secret123",
		"api_key":  "sk-123456789",
		"data":     "normal data",
		"nested": map[string]any{
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

func (suite *ProxyLoggingSecurityTestSuite) TestSanitizeTruncatedBodyRedactsAndKeepsUTF8() {
	// A large body whose JSON parse fails (malformed JSON) must still redact
	// sensitive fields in the retained prefix, and the truncation must not
	// cut a multi-byte UTF-8 codepoint.
	secret := "hunter2-S3cret"
	big := strings.Repeat("a", MaxLogBodySize+10)
	// Unterminated string -> JSON decode fails -> falls through to truncation.
	malformed := `{"password": "` + secret + `"` + big

	result := sanitizeForLogging([]byte(malformed), "application/json")
	if strings.Contains(result, secret) {
		suite.T().Errorf("truncated body leaked secret: %q", result)
	}
	if !strings.Contains(result, RedactedPlaceholder) {
		suite.T().Errorf("expected redaction placeholder in %q", result)
	}
	if !utf8.ValidString(result) {
		suite.T().Errorf("result is not valid UTF-8: %q", result)
	}
	if !strings.HasSuffix(result, TruncatedSuffix) {
		suite.T().Errorf("expected truncated suffix in %q", result)
	}
}

func (suite *ProxyLoggingSecurityTestSuite) TestTruncateUTF8RuneBoundary() {
	// A multi-byte rune straddling the byte cutoff must be backed off, not split.
	orig := strings.Repeat("a", MaxLogBodySize-1) + "界"
	out := truncateUTF8(orig, MaxLogBodySize)
	if !utf8.ValidString(out) {
		suite.T().Fatalf("truncateUTF8 produced invalid UTF-8: %x", out)
	}
	if len(out) != MaxLogBodySize-1 {
		suite.T().Fatalf("expected backoff to %d bytes, got %d", MaxLogBodySize-1, len(out))
	}
}

// TestTruncateUTF8 covers the truncateUTF8 backoff logic directly: it must
// never validate the whole string (that emptied any body with a single
// invalid byte anywhere), and it must stay O(1) backoff bounded by
// utf8.UTFMax-1, not O(n) whole-string re-validation per byte removed.
func (suite *ProxyLoggingSecurityTestSuite) TestTruncateUTF8() {
	const limit = 100

	tests := []struct {
		name      string
		input     string
		maxBytes  int
		checkFunc func(t *testing.T, out string)
	}{
		{
			name:     "multi-byte rune straddling the cut is backed off, not split",
			input:    strings.Repeat("a", limit-1) + "界", // 界 is 3 bytes, starts at limit-1
			maxBytes: limit,
			checkFunc: func(t *testing.T, out string) {
				if !utf8.ValidString(out) {
					t.Fatalf("output is not valid UTF-8: %x", out)
				}
				if out != strings.Repeat("a", limit-1) {
					t.Fatalf("expected the partial rune dropped entirely, got %q", out)
				}
			},
		},
		{
			name:     "invalid byte early in the string cuts near the limit, not to empty",
			input:    strings.Repeat("a", limit-1) + string([]byte{0x80}) + strings.Repeat("b", 500),
			maxBytes: limit,
			checkFunc: func(t *testing.T, out string) {
				if out == "" {
					t.Fatalf("truncateUTF8 emptied a string with a single stray invalid byte")
				}
				if len(out) < limit-utf8.UTFMax {
					t.Fatalf("expected cut near the %d-byte limit, got %d bytes: %q", limit, len(out), out)
				}
			},
		},
		{
			name:     "pure ASCII is cut exactly at maxBytes, unchanged behavior",
			input:    strings.Repeat("x", 500),
			maxBytes: limit,
			checkFunc: func(t *testing.T, out string) {
				if len(out) != limit {
					t.Fatalf("expected exact cut at %d bytes, got %d", limit, len(out))
				}
				if out != strings.Repeat("x", limit) {
					t.Fatalf("unexpected ASCII truncation result: %q", out)
				}
			},
		},
		{
			name:     "string shorter than maxBytes is returned unchanged",
			input:    "short string",
			maxBytes: limit,
			checkFunc: func(t *testing.T, out string) {
				if out != "short string" {
					t.Fatalf("expected input unchanged, got %q", out)
				}
			},
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			out := truncateUTF8(tt.input, tt.maxBytes)
			tt.checkFunc(suite.T(), out)
		})
	}
}

// TestTruncateUTF8NotQuadratic guards the O(n^2) regression: truncating a
// body that is invalid UTF-8 overall (one stray byte, otherwise huge) must
// stay fast and must not collapse the result to empty.
func (suite *ProxyLoggingSecurityTestSuite) TestTruncateUTF8NotQuadratic() {
	body := string([]byte{0x80}) + strings.Repeat("a", MaxLogBodySize*10)
	out := truncateUTF8(body, MaxLogBodySize)
	if len(out) == 0 {
		suite.T().Fatalf("truncateUTF8 collapsed an otherwise-valid large body to empty")
	}
	if len(out) < MaxLogBodySize-utf8.UTFMax {
		suite.T().Fatalf("expected cut near %d bytes, got %d", MaxLogBodySize, len(out))
	}
}

// TestSanitizeTruncatedBodySecretStraddlingCutNotExposed verifies that a
// sensitive field value whose closing quote falls past the MaxLogBodySize
// cut point (so the "field":"value" redaction regex cannot match it) is
// still not leaked in the truncated log line.
func (suite *ProxyLoggingSecurityTestSuite) TestSanitizeTruncatedBodySecretStraddlingCutNotExposed() {
	secret := "S3cretStraddle-" + strings.Repeat("Z", 100)
	head := strings.Repeat("x", MaxLogBodySize-50)
	body := head + `"password": "` + secret + `"}`
	suite.Greater(len(body), MaxLogBodySize, "test body must exceed MaxLogBodySize to hit the truncation path")

	// Sanity check: the opening quote of the secret value must land before
	// the cut, and the closing quote after it, so the value straddles it.
	openIdx := strings.Index(body, secret)
	suite.Less(openIdx, MaxLogBodySize, "secret must start before the cut")
	suite.Greater(openIdx+len(secret), MaxLogBodySize, "secret must extend past the cut")

	result := sanitizeForLogging([]byte(body), "text/plain")

	suite.NotContains(result, secret, "full secret leaked")
	suite.NotContains(result, "S3cretStraddle", "partial secret prefix leaked")
	suite.True(utf8.ValidString(result), "result must be valid UTF-8")
	suite.True(strings.HasSuffix(result, TruncatedSuffix), "result must carry the truncated suffix")
}

// TestSanitizeTruncatedBodyTrimsPartialRuneOnPlainBody verifies the
// rune-boundary trim in sanitizeTruncatedBody runs on the raw
// body[:MaxLogBodySize] cut, not only in the final truncateUTF8 call. The
// body here is plain non-JSON UTF-8 prose with no quotes and no sensitive
// field names, so redactPatternInString/stripTrailingUnterminatedQuote
// cannot match anything and cannot shrink the prefix: by the time the final
// truncateUTF8 safety-net call runs, len(prefix) is already <=
// MaxLogBodySize, so its len-guard makes it a no-op. Only trimming the raw
// cut up front (before redaction) can remove the multi-byte rune straddling
// byte MaxLogBodySize here.
func (suite *ProxyLoggingSecurityTestSuite) TestSanitizeTruncatedBodyTrimsPartialRuneOnPlainBody() {
	// "界" is 3 bytes and starts at byte MaxLogBodySize-1, so it straddles
	// the cut point exactly like TestTruncateUTF8RuneBoundary's input, but
	// driven through the full sanitizeForLogging -> sanitizeTruncatedBody
	// path instead of calling truncateUTF8 directly.
	body := strings.Repeat("a", MaxLogBodySize-1) + "界" + strings.Repeat("b", 200)
	suite.Greater(len(body), MaxLogBodySize, "test body must exceed MaxLogBodySize to hit the truncation path")

	result := sanitizeForLogging([]byte(body), "text/plain")

	suite.True(utf8.ValidString(result), "result must be valid UTF-8")
	suite.True(strings.HasSuffix(result, TruncatedSuffix), "result must carry the truncated suffix")
	suite.NotContains(result, "界", "a split multi-byte rune must not survive truncation")
}
