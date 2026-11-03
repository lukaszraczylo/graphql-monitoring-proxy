package main

import (
	"encoding/base64"
	"net/http"
	"testing"

	libpack_logger "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
	"github.com/stretchr/testify/assert"
)

func unsignedJWT(sub, role string) string {
	enc := base64.RawURLEncoding.EncodeToString
	return "Bearer " + enc([]byte(`{"alg":"none"}`)) + "." +
		enc([]byte(`{"sub":"`+sub+`","role":"`+role+`"}`)) + "."
}

// Pins theorems/Theorems/WsIdentityGate.lean: the identity the backend acts as is the one that is gated.
func TestAuthorizePayloadIdentity(t *testing.T) {
	cfg = &config{Logger: libpack_logger.New()}
	cfg.Client.JWTUserClaimPath = "sub"
	cfg.Client.JWTRoleClaimPath = "role"
	bannedUsersIDs.Store("banned-user", "test")
	defer bannedUsersIDs.Delete("banned-user")

	tests := []struct {
		name     string
		upgrade  http.Header
		enriched http.Header
		wantErr  error
	}{
		{"payload adds banned identity", http.Header{}, http.Header{"Authorization": {unsignedJWT("banned-user", "user")}}, errWSBanned},
		{"payload replaces clean identity with banned", http.Header{"Authorization": {unsignedJWT("ok-user", "user")}}, http.Header{"Authorization": {unsignedJWT("banned-user", "user")}}, errWSBanned},
		{"payload identity clean", http.Header{}, http.Header{"Authorization": {unsignedJWT("ok-user", "user")}}, nil},
		{"payload does not change identity", http.Header{"Authorization": {unsignedJWT("ok-user", "user")}}, http.Header{"Authorization": {unsignedJWT("ok-user", "user")}}, nil},
		{"no authorization anywhere", http.Header{}, http.Header{}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.ErrorIs(t, authorizePayloadIdentity(tt.upgrade, tt.enriched), tt.wantErr)
		})
	}
}
