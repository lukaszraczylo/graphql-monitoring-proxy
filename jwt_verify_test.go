package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	libpack_logger "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
)

// ---------------------------------------------------------------------------
// Test helpers: in-test key generation and token signing.
// ---------------------------------------------------------------------------

// generateRSAKeyPair generates an RSA key pair and PEM-encodes the public
// key (PKIX, matching what jwt.ParseRSAPublicKeyFromPEM expects).
func generateRSAKeyPair(t *testing.T) (*rsa.PrivateKey, []byte) {
	t.Helper()
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate RSA key: %v", err)
	}
	pubBytes, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	if err != nil {
		t.Fatalf("failed to marshal RSA public key: %v", err)
	}
	return priv, pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubBytes})
}

// generateECKeyPair generates an ECDSA (P-256) key pair and PEM-encodes the
// public key (PKIX, matching what jwt.ParseECPublicKeyFromPEM expects).
func generateECKeyPair(t *testing.T) (*ecdsa.PrivateKey, []byte) {
	t.Helper()
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate EC key: %v", err)
	}
	pubBytes, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	if err != nil {
		t.Fatalf("failed to marshal EC public key: %v", err)
	}
	return priv, pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubBytes})
}

func signToken(t *testing.T, method jwt.SigningMethod, key any, claims jwt.MapClaims) string {
	t.Helper()
	signed, err := jwt.NewWithClaims(method, claims).SignedString(key)
	if err != nil {
		t.Fatalf("failed to sign token with %s: %v", method.Alg(), err)
	}
	return signed
}

// ---------------------------------------------------------------------------
// Case 1 & 2: HS256 valid token verifies; wrong secret is rejected.
// ---------------------------------------------------------------------------

func TestJWTVerifier_HMAC(t *testing.T) {
	const secret = "correct-horse-battery-staple"
	token := signToken(t, jwt.SigningMethodHS256, []byte(secret), jwt.MapClaims{
		"sub": "alice",
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	tests := []struct {
		name    string
		secret  string
		wantErr bool
	}{
		{name: "valid HS256 token verifies and returns claims", secret: secret, wantErr: false},
		{name: "wrong HMAC secret is rejected", secret: "totally-different-secret", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := newJWTVerifier(jwtVerifierConfig{Secret: tt.secret})
			if err != nil {
				t.Fatalf("newJWTVerifier() error = %v", err)
			}

			claims, err := v.Verify(token)
			if tt.wantErr {
				if err == nil {
					t.Fatal("Verify() expected an error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("Verify() unexpected error: %v", err)
			}
			if claims["sub"] != "alice" {
				t.Fatalf("Verify() sub claim = %v, want %q", claims["sub"], "alice")
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Case 3: alg=none is always rejected.
// ---------------------------------------------------------------------------

func TestJWTVerifier_AlgNoneRejected(t *testing.T) {
	v, err := newJWTVerifier(jwtVerifierConfig{Secret: "some-secret"})
	if err != nil {
		t.Fatalf("newJWTVerifier() error = %v", err)
	}

	// jwt.SigningMethodNone.Sign refuses to sign unless the caller passes the
	// library's own unsafe sentinel key, which is exactly the point: this is
	// the only way to construct a syntactically valid alg=none token at all.
	noneToken := signToken(t, jwt.SigningMethodNone, jwt.UnsafeAllowNoneSignatureType, jwt.MapClaims{
		"sub": "alice",
	})

	if _, err := v.Verify(noneToken); err == nil {
		t.Fatal("Verify() accepted an alg=none token")
	}
}

// ---------------------------------------------------------------------------
// Case 4: alg confusion. An RSA-configured verifier must reject an HS256
// token forged with the RSA public key's own PEM bytes as the HMAC secret.
// ---------------------------------------------------------------------------

func TestJWTVerifier_AlgConfusionRejected(t *testing.T) {
	_, pubPEM := generateRSAKeyPair(t)

	v, err := newJWTVerifier(jwtVerifierConfig{PublicKeyPEM: string(pubPEM)})
	if err != nil {
		t.Fatalf("newJWTVerifier() error = %v", err)
	}

	forged := signToken(t, jwt.SigningMethodHS256, pubPEM, jwt.MapClaims{
		"sub": "victim",
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	if _, err := v.Verify(forged); err == nil {
		t.Fatal("Verify() accepted an HS256 token signed with the RSA public key bytes as the HMAC secret (alg confusion)")
	}
}

// ---------------------------------------------------------------------------
// Case 5: expired token is rejected; sufficient leeway accepts the same
// token.
// ---------------------------------------------------------------------------

func TestJWTVerifier_ExpiryAndLeeway(t *testing.T) {
	const secret = "leeway-secret"
	token := signToken(t, jwt.SigningMethodHS256, []byte(secret), jwt.MapClaims{
		"sub": "alice",
		"exp": time.Now().Add(-30 * time.Second).Unix(),
	})

	t.Run("expired token without leeway is rejected", func(t *testing.T) {
		v, err := newJWTVerifier(jwtVerifierConfig{Secret: secret})
		if err != nil {
			t.Fatalf("newJWTVerifier() error = %v", err)
		}
		if _, err := v.Verify(token); err == nil {
			t.Fatal("Verify() expected an error for an expired token, got nil")
		}
	})

	t.Run("sufficient leeway accepts the same token", func(t *testing.T) {
		v, err := newJWTVerifier(jwtVerifierConfig{Secret: secret, LeewaySeconds: 60})
		if err != nil {
			t.Fatalf("newJWTVerifier() error = %v", err)
		}
		if _, err := v.Verify(token); err != nil {
			t.Fatalf("Verify() unexpected error with sufficient leeway: %v", err)
		}
	})
}

// ---------------------------------------------------------------------------
// Case 6: issuer mismatch and audience mismatch are both rejected.
// ---------------------------------------------------------------------------

func TestJWTVerifier_IssuerAndAudienceMismatch(t *testing.T) {
	const secret = "iss-aud-secret"

	t.Run("issuer mismatch is rejected", func(t *testing.T) {
		token := signToken(t, jwt.SigningMethodHS256, []byte(secret), jwt.MapClaims{
			"sub": "alice",
			"iss": "wrong-issuer",
			"exp": time.Now().Add(time.Hour).Unix(),
		})

		v, err := newJWTVerifier(jwtVerifierConfig{Secret: secret, Issuer: "expected-issuer"})
		if err != nil {
			t.Fatalf("newJWTVerifier() error = %v", err)
		}
		if _, err := v.Verify(token); err == nil {
			t.Fatal("Verify() expected an issuer-mismatch error, got nil")
		}
	})

	t.Run("audience mismatch is rejected", func(t *testing.T) {
		token := signToken(t, jwt.SigningMethodHS256, []byte(secret), jwt.MapClaims{
			"sub": "alice",
			"aud": "wrong-audience",
			"exp": time.Now().Add(time.Hour).Unix(),
		})

		v, err := newJWTVerifier(jwtVerifierConfig{Secret: secret, Audience: "expected-audience"})
		if err != nil {
			t.Fatalf("newJWTVerifier() error = %v", err)
		}
		if _, err := v.Verify(token); err == nil {
			t.Fatal("Verify() expected an audience-mismatch error, got nil")
		}
	})
}

// ---------------------------------------------------------------------------
// Supplementary: both asymmetric families (RSA and EC) round-trip
// correctly, and a token from one asymmetric algorithm is rejected by a
// verifier pinned to the other via JWT_SIGNING_METHODS.
// ---------------------------------------------------------------------------

func TestJWTVerifier_AsymmetricFamilies(t *testing.T) {
	t.Run("RS256 valid token verifies against an RSA public key", func(t *testing.T) {
		priv, pubPEM := generateRSAKeyPair(t)
		token := signToken(t, jwt.SigningMethodRS256, priv, jwt.MapClaims{
			"sub": "alice",
			"exp": time.Now().Add(time.Hour).Unix(),
		})

		v, err := newJWTVerifier(jwtVerifierConfig{PublicKeyPEM: string(pubPEM)})
		if err != nil {
			t.Fatalf("newJWTVerifier() error = %v", err)
		}
		claims, err := v.Verify(token)
		if err != nil {
			t.Fatalf("Verify() unexpected error: %v", err)
		}
		if claims["sub"] != "alice" {
			t.Fatalf("Verify() sub claim = %v, want %q", claims["sub"], "alice")
		}
	})

	t.Run("ES256 valid token verifies against an EC public key", func(t *testing.T) {
		priv, pubPEM := generateECKeyPair(t)
		token := signToken(t, jwt.SigningMethodES256, priv, jwt.MapClaims{
			"sub": "alice",
			"exp": time.Now().Add(time.Hour).Unix(),
		})

		v, err := newJWTVerifier(jwtVerifierConfig{PublicKeyPEM: string(pubPEM)})
		if err != nil {
			t.Fatalf("newJWTVerifier() error = %v", err)
		}
		claims, err := v.Verify(token)
		if err != nil {
			t.Fatalf("Verify() unexpected error: %v", err)
		}
		if claims["sub"] != "alice" {
			t.Fatalf("Verify() sub claim = %v, want %q", claims["sub"], "alice")
		}
	})

	t.Run("ES256 token rejected by a verifier pinned to RS256 only", func(t *testing.T) {
		ecPriv, _ := generateECKeyPair(t)
		_, rsaPubPEM := generateRSAKeyPair(t)
		token := signToken(t, jwt.SigningMethodES256, ecPriv, jwt.MapClaims{
			"sub": "alice",
			"exp": time.Now().Add(time.Hour).Unix(),
		})

		v, err := newJWTVerifier(jwtVerifierConfig{PublicKeyPEM: string(rsaPubPEM), SigningMethods: "RS256"})
		if err != nil {
			t.Fatalf("newJWTVerifier() error = %v", err)
		}
		if _, err := v.Verify(token); err == nil {
			t.Fatal("Verify() accepted an ES256 token against a verifier pinned to RS256 only")
		}
	})
}

// ---------------------------------------------------------------------------
// Startup validation: newJWTVerifier must fail loudly on a bad key-source
// count or an unparsable key, never boot half-configured.
// ---------------------------------------------------------------------------

func TestNewJWTVerifier_StartupValidation(t *testing.T) {
	_, rsaPubPEM := generateRSAKeyPair(t)

	tests := []struct {
		name    string
		cfg     jwtVerifierConfig
		wantErr bool
	}{
		{name: "zero key sources is a fatal error", cfg: jwtVerifierConfig{}, wantErr: true},
		{
			name:    "two key sources is a fatal error",
			cfg:     jwtVerifierConfig{Secret: "s", PublicKeyPEM: string(rsaPubPEM)},
			wantErr: true,
		},
		{name: "exactly one key source is valid", cfg: jwtVerifierConfig{Secret: "s"}, wantErr: false},
		{
			name:    "signing method outside the key's family is a fatal error",
			cfg:     jwtVerifierConfig{Secret: "s", SigningMethods: "RS256"},
			wantErr: true,
		},
		{
			name:    "unparsable PEM is a fatal error",
			cfg:     jwtVerifierConfig{PublicKeyPEM: "-----BEGIN PUBLIC KEY-----\nnot-a-real-key\n-----END PUBLIC KEY-----"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := newJWTVerifier(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Fatalf("newJWTVerifier() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Case 7: verification OFF (JWTVerifier nil) preserves the legacy,
// unverified decode path byte for byte -- required for a no-config-change
// upgrade to be default-preserving.
// ---------------------------------------------------------------------------

func TestExtractClaimsFromJWTHeader_VerifyOffBackwardCompat(t *testing.T) {
	cfgMutex.Lock()
	if cfg == nil {
		cfg = &config{}
	}
	if cfg.Logger == nil {
		cfg.Logger = libpack_logger.New()
	}
	origClient := cfg.Client
	cfg.Client.JWTVerifier = nil
	cfg.Client.JWTUserClaimPath = "sub"
	cfg.Client.JWTRoleClaimPath = ""
	cfgMutex.Unlock()
	t.Cleanup(func() {
		cfgMutex.Lock()
		cfg.Client = origClient
		cfgMutex.Unlock()
	})

	// An unsigned, forged token. With verification off, extractClaimsFromJWTHeader
	// must still decode it (today's behaviour) instead of erroring.
	forged := signToken(t, jwt.SigningMethodNone, jwt.UnsafeAllowNoneSignatureType, jwt.MapClaims{
		"sub": "mallory",
	})

	usr, _, err := extractClaimsFromJWTHeader(forged)
	if err != nil {
		t.Fatalf("extractClaimsFromJWTHeader() unexpected error with verification off: %v", err)
	}
	if usr != "mallory" {
		t.Fatalf("extractClaimsFromJWTHeader() usr = %q, want %q (legacy decode must still work)", usr, "mallory")
	}
}

// ---------------------------------------------------------------------------
// Case 8: verification ON, a forged token sharing a victim's "sub" claim but
// signed with the wrong secret must make extractUserInfo return a non-nil
// error -- never the victim's identity -- so the caller 401s the request
// instead of letting it read (or write) the victim's cache entry.
// ---------------------------------------------------------------------------

func TestExtractUserInfo_ForgedTokenNeverReachesVictimCacheKey(t *testing.T) {
	const secret = "cache-isolation-secret"

	cfgMutex.Lock()
	if cfg == nil {
		cfg = &config{}
	}
	if cfg.Logger == nil {
		cfg.Logger = libpack_logger.New()
	}
	origClient := cfg.Client
	verifier, err := newJWTVerifier(jwtVerifierConfig{Secret: secret})
	if err != nil {
		cfgMutex.Unlock()
		t.Fatalf("newJWTVerifier() error = %v", err)
	}
	cfg.Client.JWTVerifier = verifier
	cfg.Client.JWTUserClaimPath = "sub"
	cfgMutex.Unlock()
	t.Cleanup(func() {
		cfgMutex.Lock()
		cfg.Client = origClient
		cfgMutex.Unlock()
	})

	// Same "sub" as a legitimate victim, but signed with the wrong secret.
	// If extractUserInfo ever returned "victim" here, the forged request
	// would land in the victim's own per-user cache bucket.
	forged := signToken(t, jwt.SigningMethodHS256, []byte("wrong-secret"), jwt.MapClaims{
		"sub": "victim",
	})

	app := fiber.New()
	var capturedErr error
	var capturedUserID, capturedRole string
	app.Get("/", func(c fiber.Ctx) error {
		capturedUserID, capturedRole, capturedErr = extractUserInfo(c)
		return c.SendStatus(fiber.StatusOK)
	})

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+forged)
	resp, err := app.Test(req, fiber.TestConfig{Timeout: 0})
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	_ = resp.Body.Close()

	if capturedErr == nil {
		t.Fatal("extractUserInfo() expected a verification error for a forged token, got nil")
	}
	if capturedUserID != defaultValue || capturedRole != defaultValue {
		t.Fatalf("extractUserInfo() on verify failure = (%q, %q), want (%q, %q) so a forged token never maps to the victim's cache bucket",
			capturedUserID, capturedRole, defaultValue, defaultValue)
	}
}
