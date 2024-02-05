package main

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/goccy/go-json"
	"github.com/lukaszraczylo/ask"
	libpack_monitoring "github.com/lukaszraczylo/graphql-monitoring-proxy/monitoring"
)

func extractClaimsFromJWTHeader(authorization string) (usr string, role string) {
	usr, role = "-", "-"

	handleError := func(msg string, details map[string]interface{}) {
		cfg.Monitoring.Increment(libpack_monitoring.MetricsFailed, nil)
		cfg.Logger.Error(msg, details)
	}

	tokenParts := strings.Split(authorization, ".")
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

	extractClaim := func(claimPath string, target *string, name string) {
		if len(claimPath) > 0 {
			var ok bool
			*target, ok = ask.For(claimMap, claimPath).String("-")
			if !ok {
				handleError(fmt.Sprintf("Can't find the %s", name), map[string]interface{}{"claim_map": claimMap, "path": claimPath})
			}
		}
	}

	extractClaim(cfg.Client.JWTUserClaimPath, &usr, "user id")
	extractClaim(cfg.Client.JWTRoleClaimPath, &role, "role")

	return
}
