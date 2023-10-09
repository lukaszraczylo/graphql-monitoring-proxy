package main

import (
	"encoding/base64"
	"strings"

	"github.com/lukaszraczylo/ask"
	libpack_monitoring "github.com/telegram-bot-app/libpack/monitoring"
)

func extractClaimsFromJWTHeader(authorization string) (usr string, role string) {
	tokenParts := strings.Split(authorization, ".")
	if len(tokenParts) != 3 {
		cfg.Monitoring.Increment(libpack_monitoring.MetricsFailed, nil)
		cfg.Logger.Error("Can't split the token", map[string]interface{}{"token": authorization})
		return
	}
	claim, err := base64.RawURLEncoding.DecodeString(tokenParts[1])
	if err != nil {
		cfg.Monitoring.Increment(libpack_monitoring.MetricsFailed, nil)
		cfg.Logger.Error("Can't decode the token", map[string]interface{}{"token": authorization})
		return
	}
	var claimMap map[string]interface{}
	err = json.Unmarshal(claim, &claimMap)
	if err != nil {
		cfg.Monitoring.Increment(libpack_monitoring.MetricsFailed, nil)
		cfg.Logger.Error("Can't unmarshal the claim", map[string]interface{}{"token": authorization})
		return
	}

	if len(cfg.Client.JWTUserClaimPath) > 0 {
		var ok bool
		usr, ok = ask.For(claimMap, cfg.Client.JWTUserClaimPath).String("-")
		if !ok {
			cfg.Monitoring.Increment(libpack_monitoring.MetricsFailed, nil)
			cfg.Logger.Error("Can't find the user id", map[string]interface{}{"claim_map": claimMap, "path": cfg.Client.JWTUserClaimPath})
			return
		}
	}

	if len(cfg.Client.JWTRoleClaimPath) > 0 {
		var ok bool
		role, ok = ask.For(claimMap, cfg.Client.JWTRoleClaimPath).String("-")
		if !ok {
			cfg.Monitoring.Increment(libpack_monitoring.MetricsFailed, nil)
			cfg.Logger.Error("Can't find the role", map[string]interface{}{"claim_map": claimMap, "path": cfg.Client.JWTRoleClaimPath})
			return
		}
	}

	return
}
