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

func extractClaimsFromJWTHeader(authorization string) (usr, role string) {
	usr, role = defaultValue, defaultValue

	tokenParts := strings.SplitN(authorization, ".", 3)
	if len(tokenParts) != 3 {
		handleError("Can't split the token", map[string]any{"token": maskToken(authorization)})
		return
	}

	claim, err := base64.RawURLEncoding.DecodeString(tokenParts[1])
	if err != nil {
		handleError("Can't decode the token", map[string]any{"token": maskToken(authorization)})
		return
	}

	var claimMap map[string]any
	if err = json.Unmarshal(claim, &claimMap); err != nil {
		handleError("Can't unmarshal the claim", map[string]any{"token": maskToken(authorization)})
		return
	}

	usr = extractClaim(claimMap, cfg.Client.JWTUserClaimPath, "user id")
	role = extractClaim(claimMap, cfg.Client.JWTRoleClaimPath, "role")

	return
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
