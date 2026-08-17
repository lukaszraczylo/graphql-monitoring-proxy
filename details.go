package main

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/goccy/go-json"
	"github.com/lukaszraczylo/ask"
	libpack_logger "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
	libpack_monitoring "github.com/lukaszraczylo/graphql-monitoring-proxy/monitoring"
)

const defaultValue = "-"

var emptyMetrics = map[string]string{}

// extractClaimsFromJWTHeader extracts the user ID and role claim from an
// Authorization header value.
//
// When cfg.Client.JWTVerifier is set (JWT_VERIFY_SIGNATURE=true), the token
// signature is checked first; a signature, expiry, or algorithm failure
// returns a non-nil err and the caller must not trust usr/role (this is the
// fix for GHSA-9gqw-h2rw-44wv, where a forged token with a guessed claim
// could otherwise read another user's cached response).
//
// When JWTVerifier is nil (the default), this falls back to the legacy
// path: the payload is decoded without checking its signature, exactly as
// before this field existed. err is always nil on that path, matching
// pre-verification behaviour byte for byte.
func extractClaimsFromJWTHeader(authorization string) (usr, role string, err error) {
	usr, role = defaultValue, defaultValue

	if cfg.Client.JWTVerifier != nil {
		claimMap, verr := cfg.Client.JWTVerifier.Verify(authorization)
		if verr != nil {
			handleError("JWT verification failed", map[string]any{"token": maskToken(authorization)})
			return defaultValue, defaultValue, verr
		}

		usr = extractClaim(claimMap, cfg.Client.JWTUserClaimPath, "user id")
		role = extractClaim(claimMap, cfg.Client.JWTRoleClaimPath, "role")
		return usr, role, nil
	}

	tokenParts := strings.SplitN(authorization, ".", 3)
	if len(tokenParts) != 3 {
		handleError("Can't split the token", map[string]any{"token": maskToken(authorization)})
		return usr, role, nil
	}

	claim, decErr := base64.RawURLEncoding.DecodeString(tokenParts[1])
	if decErr != nil {
		handleError("Can't decode the token", map[string]any{"token": maskToken(authorization)})
		return usr, role, nil
	}

	var claimMap map[string]any
	if jsonErr := json.Unmarshal(claim, &claimMap); jsonErr != nil {
		handleError("Can't unmarshal the claim", map[string]any{"token": maskToken(authorization)})
		return usr, role, nil
	}

	usr = extractClaim(claimMap, cfg.Client.JWTUserClaimPath, "user id")
	role = extractClaim(claimMap, cfg.Client.JWTRoleClaimPath, "role")

	return usr, role, nil
}

func extractClaim(claimMap map[string]any, claimPath, name string) string {
	if claimPath == "" {
		return defaultValue
	}

	// Validate claim path to prevent injection attacks
	if !isValidClaimPath(claimPath) {
		handleError(fmt.Sprintf("Invalid claim path for %s", name), map[string]any{"path": claimPath})
		return defaultValue
	}

	value, ok := ask.For(claimMap, claimPath).String(defaultValue)
	if !ok {
		handleError(fmt.Sprintf("Can't find the %s", name), map[string]any{"claim_map": sanitizeClaimMap(claimMap), "path": claimPath})
		return defaultValue
	}

	return value
}

// maskToken masks JWT tokens in logs to prevent exposure
func maskToken(token string) string {
	if len(token) <= 10 {
		return "***"
	}
	return token[:4] + "***" + token[len(token)-4:]
}

// isValidClaimPath validates JWT claim paths to prevent injection
func isValidClaimPath(path string) bool {
	if path == "" {
		return false
	}
	// Allow only alphanumeric characters, dots, underscores, and hyphens
	for _, char := range path {
		if (char < 'a' || char > 'z') &&
			(char < 'A' || char > 'Z') &&
			(char < '0' || char > '9') &&
			char != '.' && char != '_' && char != '-' {
			return false
		}
	}
	// Prevent path traversal attempts
	if strings.Contains(path, "..") || strings.Contains(path, "//") {
		return false
	}
	return true
}

// sanitizeClaimMap removes sensitive data from claim map for logging
func sanitizeClaimMap(claimMap map[string]any) map[string]any {
	sanitized := make(map[string]any)
	sensitiveKeys := map[string]bool{
		"password": true, "secret": true, "token": true, "key": true,
		"auth": true, "credential": true, "private": true,
	}

	for k, v := range claimMap {
		lowerKey := strings.ToLower(k)
		if sensitiveKeys[lowerKey] {
			sanitized[k] = "***"
		} else {
			sanitized[k] = v
		}
	}
	return sanitized
}

func handleError(msg string, details map[string]any) {
	cfg.Monitoring.Increment(libpack_monitoring.MetricsFailed, emptyMetrics)
	cfg.Logger.Error(&libpack_logger.LogMessage{
		Message: msg,
		Pairs:   details,
	})
}
