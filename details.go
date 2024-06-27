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

const (
	defaultValue = "-"
)

var (
	emptyMetrics = map[string]string{}
)

func extractClaimsFromJWTHeader(authorization string) (usr, role string) {
	usr, role = defaultValue, defaultValue

	tokenParts := strings.SplitN(authorization, ".", 3)
	if len(tokenParts) != 3 {
		handleError("Can't split the token", map[string]interface{}{"token": authorization})
		return
	}

	claim, err := base64.RawURLEncoding.DecodeString(tokenParts[1])
	if err != nil {
		handleError("Can't decode the token", map[string]interface{}{"token": authorization})
		return
	}

	var claimMap map[string]interface{}
	if err = json.Unmarshal(claim, &claimMap); err != nil {
		handleError("Can't unmarshal the claim", map[string]interface{}{"token": authorization})
		return
	}

	usr = extractClaim(claimMap, cfg.Client.JWTUserClaimPath, "user id")
	role = extractClaim(claimMap, cfg.Client.JWTRoleClaimPath, "role")

	return
}

func extractClaim(claimMap map[string]interface{}, claimPath, name string) string {
	if claimPath == "" {
		return defaultValue
	}

	value, ok := ask.For(claimMap, claimPath).String(defaultValue)
	if !ok {
		handleError(fmt.Sprintf("Can't find the %s", name), map[string]interface{}{"claim_map": claimMap, "path": claimPath})
		return defaultValue
	}

	return value
}

func handleError(msg string, details map[string]interface{}) {
	cfg.Monitoring.Increment(libpack_monitoring.MetricsFailed, emptyMetrics)
	cfg.Logger.Error(&libpack_logger.LogMessage{
		Message: msg,
		Pairs:   details,
	})
}
