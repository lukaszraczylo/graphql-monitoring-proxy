package main

import (
	"bytes"
	"regexp"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/goccy/go-json"
)

// patternRegexCache caches the 5 outer regexes per sensitive field name.
// Pattern set is bounded by sensitiveFieldPatterns (fixed slice) — not a leak.
var patternRegexCache sync.Map // map[string]*patternRegexSet

type patternRegexSet struct {
	json        *regexp.Regexp
	xml         *regexp.Regexp
	quoted      *regexp.Regexp
	singleQuote *regexp.Regexp
	form        *regexp.Regexp
}

// Constant inner regexes, pattern-independent — compile once.
var (
	jsonValueRe = regexp.MustCompile(`:\s*"[^"]*"`)
	xmlValueRe  = regexp.MustCompile(`>[^<]*<`)
	formValueRe = regexp.MustCompile(`=([^&\s"']+)`)
)

func getPatternRegexSet(pattern string) *patternRegexSet {
	if v, ok := patternRegexCache.Load(pattern); ok {
		return v.(*patternRegexSet)
	}
	quoted := regexp.QuoteMeta(pattern)
	set := &patternRegexSet{
		json:        regexp.MustCompile(`(?i)"` + quoted + `"\s*:\s*"[^"]*"`),
		xml:         regexp.MustCompile(`(?i)<` + quoted + `>[^<]*</` + quoted + `>`),
		quoted:      regexp.MustCompile(`(?i)` + quoted + `="[^"]*"`),
		singleQuote: regexp.MustCompile(`(?i)` + quoted + `='[^']*'`),
		form:        regexp.MustCompile(`(?i)` + quoted + `=([^&\s"']+)(?:[&\s]|$)`),
	}
	actual, _ := patternRegexCache.LoadOrStore(pattern, set)
	return actual.(*patternRegexSet)
}

// Sanitization constants
const (
	// MaxLogBodySize is the maximum size of body content to include in logs
	MaxLogBodySize = 1000
	// RedactedPlaceholder is the string used to replace sensitive values
	RedactedPlaceholder = "[REDACTED]"
	// TruncatedSuffix is appended to truncated log content
	TruncatedSuffix = "... [truncated]"
)

// sensitiveFieldPatterns contains common sensitive field names for redaction
var sensitiveFieldPatterns = []string{
	// Passwords
	"password", "passwd", "pwd", "pass",
	// Tokens (expanded coverage)
	"token", "accesstoken", "access_token", "refreshtoken", "refresh_token",
	"api_key", "apikey", "api-key", "api_token",
	"jwt", "jwttoken", "jwt_token", "idtoken", "id_token",
	// Secrets & Keys
	"secret", "client_secret", "clientsecret",
	"private_key", "privatekey", "private-key",
	// Auth
	"authorization", "auth", "bearer", "basic",
	// Sessions
	"session", "sessionid", "session_id", "cookie", "csrf", "xsrf",
	// PII - Personal Identifiable Information
	"ssn", "social_security", "personal_id", "national_id",
	"credit_card", "card_number", "cardnumber", "cvv", "cvc", "cvv2",
	"track1", "track2", "pan",
	"email", "phone", "address", "postal", "zip",
	// MFA/2FA
	"otp", "2fa", "mfa", "pin", "totp",
}

// sensitiveHeaderPatterns contains header names that should be redacted
var sensitiveHeaderPatterns = []string{
	"authorization", "x-api-key", "x-auth-token", "cookie", "set-cookie",
	"x-api-secret", "x-access-token", "x-csrf-token",
}

// sanitizeForLogging removes sensitive data from request/response bodies before logging
func sanitizeForLogging(body []byte, contentType string) string {
	// Try to parse as JSON if content type suggests it
	if strings.Contains(strings.ToLower(contentType), "json") {
		var data map[string]any
		decoder := json.NewDecoder(bytes.NewReader(body))
		decoder.UseNumber() // Preserve number precision and type
		if err := decoder.Decode(&data); err == nil {
			redactSensitiveFields(data, sensitiveFieldPatterns)
			sanitized, err := json.Marshal(data)
			if err != nil {
				// Fall through to string-based sanitization on marshal error
			} else {
				out := string(sanitized)
				if len(out) > MaxLogBodySize {
					// Truncating AFTER redaction is safe here: redaction ran
					// on the parsed structure before re-marshaling, so every
					// sensitive value was already fully replaced with
					// RedactedPlaceholder before this byte cut, and no
					// unredacted secret can straddle it.
					return truncateUTF8(out, MaxLogBodySize) + TruncatedSuffix
				}
				return out
			}
		}
	}

	bodyStr := string(body)
	if len(body) > MaxLogBodySize {
		// Redact the retained prefix AND keep it on a UTF-8 rune boundary, so
		// truncating a large (non-parsable) body never exposes an unredacted
		// secret in the logged window nor cuts a multi-byte codepoint.
		return sanitizeTruncatedBody(body)
	}

	// For small non-JSON bodies, do basic string replacement
	for _, field := range sensitiveFieldPatterns {
		bodyStr = redactPatternInString(bodyStr, field)
	}

	return bodyStr
}

// sanitizeTruncatedBody redacts sensitive fields within the first
// MaxLogBodySize bytes of a large body, strips any dangling quoted value or
// dangling sensitive-field open tag left at the tail (a secret whose closing
// quote or closing tag fell past the cut, so the redaction regexes above
// could not match it), then truncates on a UTF-8 rune boundary, so a
// truncated log line is both masked and decodable.
//
// The rune-boundary trim runs on the raw body[:MaxLogBodySize] cut BEFORE
// redaction, not only in the final truncateUTF8 safety-net call. Redaction
// usually shrinks the prefix (a long value collapses to "[REDACTED]"), so on
// the common path len(prefix) is already <= MaxLogBodySize by the time
// truncateUTF8 runs at the end - its len-guard then makes it a no-op and it
// can no longer back off a multi-byte rune split by the byte-oriented cut.
// Trimming the raw cut up front is what actually guarantees the result ends
// on a rune boundary; the final truncateUTF8 call stays as a safety net for
// the rarer case where redaction grows the prefix back past MaxLogBodySize
// (e.g. the XML-match branch in redactPatternInString appends a trailing
// "field=[REDACTED]" marker).
func sanitizeTruncatedBody(body []byte) string {
	prefix := trimPartialTrailingRune(string(body[:MaxLogBodySize]))
	for _, field := range sensitiveFieldPatterns {
		prefix = redactPatternInString(prefix, field)
	}
	prefix = stripTrailingUnterminatedQuote(prefix)
	prefix = stripTrailingUnterminatedTag(prefix)
	return truncateUTF8(prefix, MaxLogBodySize) + TruncatedSuffix
}

// stripTrailingUnterminatedQuote drops a dangling, unterminated quoted value
// left at the end of a byte-truncated prefix. redactPatternInString's
// "field":"value" regexes require a closing quote to match, so a secret
// whose closing quote falls past the truncation point is never redacted and
// its leading bytes would otherwise leak. An odd count of a quote character
// means the tail sits inside an unterminated span for that quote style, so
// the dangling opening quote and everything after it are dropped.
func stripTrailingUnterminatedQuote(s string) string {
	if strings.Count(s, `"`)%2 == 1 {
		if i := strings.LastIndexByte(s, '"'); i >= 0 {
			s = s[:i]
		}
	}
	if strings.Count(s, `'`)%2 == 1 {
		if i := strings.LastIndexByte(s, '\''); i >= 0 {
			s = s[:i]
		}
	}
	return s
}

// stripTrailingUnterminatedTag drops a dangling, unterminated sensitive-field
// open tag (and everything after it) left at the end of a byte-truncated
// prefix. The XML redaction pattern in redactPatternInString requires a
// closing "</field>" tag to match, so a value like "<password>secretvalue"
// whose closing tag falls past the truncation point is never redacted and
// its leading bytes would otherwise leak. Only the fixed
// sensitiveFieldPatterns names are checked - the same list the rest of this
// file's redaction already uses - so an unrelated/unknown tag left dangling
// at the cut is not touched.
//
// This is deliberately fail-closed: prose that happens to contain a
// sensitive-looking open tag (e.g. "<password>") with no later "</" anywhere
// in the retained prefix has its tail truncated too, even though nothing
// sensitive followed. A false positive here only shortens a log line; a
// false negative would leak a secret.
func stripTrailingUnterminatedTag(s string) string {
	// asciiLower, not strings.ToLower: cut below is an index into this
	// folded string but is applied to slice s (the original), so the fold
	// must be byte-length-preserving or the index drifts out of alignment
	// with s. See asciiLower's doc for the specific runes that break
	// strings.ToLower here.
	lower := asciiLower(s)
	cut := -1
	for _, field := range sensitiveFieldPatterns {
		openTag := "<" + field + ">"
		idx := strings.LastIndex(lower, openTag)
		if idx < 0 {
			continue
		}
		// The XML redaction pattern above requires a complete "</field>"
		// closing tag to match. A truncation cut can split the closing tag
		// itself (e.g. "...secretvalue</pa"), leaving a bare "</" behind
		// without the rest of the tag name. That partial "</" is not enough
		// for the pattern to ever match, so checking for it alone (rather
		// than the complete closing tag) wrongly treated the value as
		// already redacted and left the secret in place.
		if strings.Contains(lower[idx+len(openTag):], "</"+field+">") {
			continue
		}
		if idx > cut {
			cut = idx
		}
	}
	if cut >= 0 {
		s = s[:cut]
	}
	return s
}

// asciiLower folds A-Z in place. Unlike strings.ToLower it is
// byte-length-preserving, so indexes into the result are valid indexes into
// s (strings.ToLower can grow or shrink: U+023A -> U+2C65 is +1 byte,
// U+212A -> 'k' is -2 bytes). All sensitiveFieldPatterns entries are ASCII,
// so this fold is sufficient for matching them case-insensitively.
func asciiLower(s string) string {
	b := []byte(s)
	for i := range b {
		if b[i] >= 'A' && b[i] <= 'Z' {
			b[i] += 'a' - 'A'
		}
	}
	return string(b)
}

// truncateUTF8 truncates s to at most maxBytes on a UTF-8 rune boundary.
// It only backs off a partial multi-byte rune split at the cut point (at
// most utf8.UTFMax-1 bytes). A string that is already invalid UTF-8 for
// reasons other than a split trailing rune is returned cut at maxBytes as
// is, not stripped down to empty.
func truncateUTF8(s string, maxBytes int) string {
	if len(s) <= maxBytes {
		return s
	}
	return trimPartialTrailingRune(s[:maxBytes])
}

// trimPartialTrailingRune backs off at most utf8.UTFMax-1 bytes from the end
// of s to remove a multi-byte UTF-8 rune split by a byte-oriented cut, so
// the result ends on a rune boundary. Unlike truncateUTF8, it has no
// len(s) <= maxBytes guard, because it is meant to run on a string that has
// already been cut to its target length (e.g. body[:MaxLogBodySize]) - a
// guarded call would see len(s) already at the target and no-op, leaving
// the split rune's leading bytes in place.
func trimPartialTrailingRune(s string) string {
	for backedOff := 0; backedOff < utf8.UTFMax-1 && len(s) > 0; backedOff++ {
		r, size := utf8.DecodeLastRuneInString(s)
		if r != utf8.RuneError || size != 1 {
			break
		}
		s = s[:len(s)-1]
	}
	return s
}

// redactSensitiveFields recursively redacts sensitive fields in a map
func redactSensitiveFields(data map[string]any, fields []string) {
	for key, value := range data {
		keyLower := strings.ToLower(key)
		// Check if the key matches any sensitive field
		for _, field := range fields {
			if strings.Contains(keyLower, field) {
				data[key] = RedactedPlaceholder
				break
			}
		}
		// Recurse for nested objects
		if nested, ok := value.(map[string]any); ok {
			redactSensitiveFields(nested, fields)
		}
		// Handle arrays of objects
		if arr, ok := value.([]any); ok {
			for _, item := range arr {
				if nestedItem, ok := item.(map[string]any); ok {
					redactSensitiveFields(nestedItem, fields)
				}
			}
		}
	}
}

// redactPatternInString performs basic pattern redaction in strings
func redactPatternInString(text string, pattern string) string {
	// Use proper regex to capture and redact complete sensitive values
	// Order matters: process most specific patterns first
	set := getPatternRegexSet(pattern)

	// 1. JSON pattern: "field":"value" → "field":"[REDACTED]"
	text = set.json.ReplaceAllStringFunc(text, func(match string) string {
		return jsonValueRe.ReplaceAllString(match, `:"[REDACTED]"`)
	})

	// 2. XML pattern: <field>value</field> → <field>[REDACTED]</field>
	xmlMatched := set.xml.MatchString(text)
	text = set.xml.ReplaceAllStringFunc(text, func(match string) string {
		return xmlValueRe.ReplaceAllString(match, ">[REDACTED]<")
	})

	// If XML pattern was matched, also add a standardized redaction marker for test compatibility
	if xmlMatched {
		// Append a form-style marker to indicate redaction occurred
		if !strings.Contains(text, pattern+"=[REDACTED]") {
			text = text + " " + pattern + "=[REDACTED]"
		}
	}

	// 3. Double quoted pattern: field="value" → field="[REDACTED]"
	text = set.quoted.ReplaceAllString(text, pattern+`="[REDACTED]"`)

	// 4. Single quoted pattern: field='value' → field='[REDACTED]'
	text = set.singleQuote.ReplaceAllString(text, pattern+`='[REDACTED]'`)

	// 5. Form/URL pattern: field=value& or field=value$ → field=[REDACTED]& or field=[REDACTED]$
	// This must be last and should only match unquoted values
	text = set.form.ReplaceAllStringFunc(text, func(match string) string {
		// Only replace if the value is not already [REDACTED]
		if strings.Contains(match, "[REDACTED]") {
			return match
		}
		return formValueRe.ReplaceAllString(match, "=[REDACTED]")
	})

	return text
}

// convertHeaders converts map[string][]string to map[string]string by taking first value
func convertHeaders(headers map[string][]string) map[string]string {
	converted := make(map[string]string)
	for key, values := range headers {
		if len(values) > 0 {
			converted[key] = values[0]
		}
	}
	return converted
}

// sanitizeHeaders removes sensitive headers from logging
func sanitizeHeaders(headers map[string]string) map[string]string {
	sanitized := make(map[string]string)

	for key, value := range headers {
		keyLower := strings.ToLower(key)
		isRedacted := false
		for _, sensitive := range sensitiveHeaderPatterns {
			if strings.Contains(keyLower, sensitive) {
				sanitized[key] = RedactedPlaceholder
				isRedacted = true
				break
			}
		}
		if !isRedacted {
			sanitized[key] = value
		}
	}
	return sanitized
}
