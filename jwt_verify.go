package main

import (
	"crypto"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
)

// keyFamily identifies which class of signing algorithms a configured JWT
// key source belongs to. Verify's keyfunc (see JWTVerifier.keyfn) enforces
// that a token's algorithm belongs to the SAME family as the configured
// key, so an HMAC-signed token can never be checked against an asymmetric
// key and vice versa. This is what stops the alg-confusion attack fixed by
// GHSA-9gqw-h2rw-44wv (e.g. an attacker forging an HS256 token whose HMAC
// secret is the server's own RSA public key bytes).
type keyFamily int

const (
	keyFamilyHMAC keyFamily = iota
	keyFamilyAsymmetric
)

// hmacMethods and asymmetricMethods are the JWT "alg" header values
// belonging to each keyFamily. "none" is deliberately absent from both
// lists: WithValidMethods (see Verify) only accepts algorithms in one of
// these lists, so a token with alg=none is rejected before JWTVerifier's
// keyfunc is even called.
var (
	hmacMethods       = []string{"HS256", "HS384", "HS512"}
	asymmetricMethods = []string{"RS256", "RS384", "RS512", "PS256", "PS384", "PS512", "ES256", "ES384", "ES512"}
)

// jwtVerifierConfig holds the resolved JWT_* environment settings
// newJWTVerifier needs to build a JWTVerifier. main.go's parseConfig reads
// the environment variables and populates this struct.
type jwtVerifierConfig struct {
	// Secret is the HMAC shared secret (JWT_SECRET). Selects keyFamilyHMAC.
	Secret string
	// PublicKeyPEM is a PEM-encoded RSA or ECDSA public key (JWT_PUBLIC_KEY),
	// either inline (contains "-----BEGIN") or a path to a file containing
	// it. Selects keyFamilyAsymmetric.
	PublicKeyPEM string
	// JWKSURL is a JWKS endpoint (JWT_JWKS_URL) fetched and kept up to date
	// by keyfunc. Selects keyFamilyAsymmetric.
	JWKSURL string
	// SigningMethods is an optional comma-separated allowlist
	// (JWT_SIGNING_METHODS), e.g. "RS256,ES256". Empty means the full
	// algorithm family for the configured key source is allowed.
	SigningMethods string
	// Issuer, when set, is required to match the token's "iss" claim.
	Issuer string
	// Audience, when set, is required to be present in the token's "aud" claim.
	Audience string
	// LeewaySeconds is the clock-skew tolerance applied to exp/nbf checks.
	LeewaySeconds int
}

// JWTVerifier verifies the signature of a JWT bearer token against exactly
// one configured key source (an HMAC secret, a PEM public key, or a JWKS
// endpoint), pinned to that source's algorithm family. Build one with
// newJWTVerifier; a nil *JWTVerifier on config.Client (the default) means
// signature verification is off and callers must use the legacy unverified
// decode path (see details.go extractClaimsFromJWTHeader).
type JWTVerifier struct {
	jwks           keyfunc.Keyfunc
	hmacSecret     []byte
	publicKey      crypto.PublicKey
	family         keyFamily
	allowedMethods []string
	issuer         string
	audience       string
	leeway         time.Duration
}

// newJWTVerifier builds a JWTVerifier from cfg. Exactly one of
// cfg.Secret, cfg.PublicKeyPEM, or cfg.JWKSURL must be set; any other count
// is a fatal configuration error, not a silent fallback, so the proxy never
// starts half-configured (GHSA-9gqw-h2rw-44wv requires verification to
// actually run once JWT_VERIFY_SIGNATURE is on).
func newJWTVerifier(cfg jwtVerifierConfig) (*JWTVerifier, error) {
	sourcesConfigured := 0
	if cfg.Secret != "" {
		sourcesConfigured++
	}
	if cfg.PublicKeyPEM != "" {
		sourcesConfigured++
	}
	if cfg.JWKSURL != "" {
		sourcesConfigured++
	}
	if sourcesConfigured != 1 {
		return nil, fmt.Errorf("JWT_VERIFY_SIGNATURE requires exactly one of JWT_SECRET, JWT_PUBLIC_KEY, JWT_JWKS_URL, got %d configured", sourcesConfigured)
	}

	v := &JWTVerifier{
		issuer:   cfg.Issuer,
		audience: cfg.Audience,
		leeway:   time.Duration(cfg.LeewaySeconds) * time.Second,
	}

	switch {
	case cfg.Secret != "":
		v.family = keyFamilyHMAC
		v.hmacSecret = []byte(cfg.Secret)
	case cfg.PublicKeyPEM != "":
		v.family = keyFamilyAsymmetric
		pemBytes, err := resolvePEMSource(cfg.PublicKeyPEM)
		if err != nil {
			return nil, fmt.Errorf("failed to read JWT_PUBLIC_KEY: %w", err)
		}
		pubKey, err := parsePublicKeyPEM(pemBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse JWT_PUBLIC_KEY: %w", err)
		}
		v.publicKey = pubKey
	case cfg.JWKSURL != "":
		v.family = keyFamilyAsymmetric
		kf, err := keyfunc.NewDefault([]string{cfg.JWKSURL})
		if err != nil {
			return nil, fmt.Errorf("failed to initialise JWKS from JWT_JWKS_URL: %w", err)
		}
		v.jwks = kf
	}

	allowed, err := resolveAllowedMethods(v.family, cfg.SigningMethods)
	if err != nil {
		return nil, err
	}
	v.allowedMethods = allowed

	return v, nil
}

// resolveAllowedMethods intersects the operator-configured
// JWT_SIGNING_METHODS allowlist with the algorithm family of the
// configured key. An empty allowlist means "allow the whole family". A
// configured method outside the key's family (e.g. JWT_SECRET set together
// with JWT_SIGNING_METHODS=RS256) is a fatal configuration error: it can
// never be satisfied, since JWTVerifier.keyfn refuses to return an HMAC
// secret for an asymmetric alg or vice versa.
func resolveAllowedMethods(family keyFamily, signingMethods string) ([]string, error) {
	familyMethods := hmacMethods
	if family == keyFamilyAsymmetric {
		familyMethods = asymmetricMethods
	}

	if strings.TrimSpace(signingMethods) == "" {
		return familyMethods, nil
	}

	allowedInFamily := make(map[string]bool, len(familyMethods))
	for _, m := range familyMethods {
		allowedInFamily[m] = true
	}

	configured := strings.Split(signingMethods, ",")
	allowed := make([]string, 0, len(configured))
	for _, m := range configured {
		m = strings.TrimSpace(m)
		if m == "" {
			continue
		}
		if !allowedInFamily[m] {
			return nil, fmt.Errorf("JWT_SIGNING_METHODS %q is not part of the configured key's algorithm family", m)
		}
		allowed = append(allowed, m)
	}
	if len(allowed) == 0 {
		return nil, errors.New("JWT_SIGNING_METHODS resolved to no valid algorithms for the configured key")
	}
	return allowed, nil
}

// resolvePEMSource returns the PEM bytes for JWT_PUBLIC_KEY. The value is
// used inline when it already contains a PEM header, otherwise it is
// treated as a path and read from disk.
func resolvePEMSource(value string) ([]byte, error) {
	if strings.Contains(value, "-----BEGIN") {
		return []byte(value), nil
	}
	data, err := os.ReadFile(value)
	if err != nil {
		return nil, fmt.Errorf("reading PEM file %q: %w", value, err)
	}
	return data, nil
}

// parsePublicKeyPEM parses pemBytes as either an RSA or an ECDSA public
// key, trying RSA first. It returns an error if neither parses.
func parsePublicKeyPEM(pemBytes []byte) (crypto.PublicKey, error) {
	if rsaKey, err := jwt.ParseRSAPublicKeyFromPEM(pemBytes); err == nil {
		return rsaKey, nil
	}
	if ecKey, err := jwt.ParseECPublicKeyFromPEM(pemBytes); err == nil {
		return ecKey, nil
	}
	return nil, errors.New("PEM does not contain a valid RSA or ECDSA public key")
}

// Verify checks the signature, algorithm family, and standard claims (exp,
// nbf, and the configured issuer/audience) of a bearer token taken from an
// Authorization header value. It returns the token's claims only when every
// check passes; on any failure it returns a non-nil error and the claims
// MUST NOT be used, so callers never trust a claim from a token that failed
// verification.
func (v *JWTVerifier) Verify(authorization string) (map[string]any, error) {
	tokenString := stripBearerPrefix(authorization)

	opts := []jwt.ParserOption{
		jwt.WithValidMethods(v.allowedMethods),
		jwt.WithLeeway(v.leeway),
	}
	if v.issuer != "" {
		opts = append(opts, jwt.WithIssuer(v.issuer))
	}
	if v.audience != "" {
		opts = append(opts, jwt.WithAudience(v.audience))
	}

	parser := jwt.NewParser(opts...)

	claims := jwt.MapClaims{}
	token, err := parser.ParseWithClaims(tokenString, claims, v.keyfn)
	if err != nil {
		return nil, fmt.Errorf("jwt verification failed: %w", err)
	}
	if !token.Valid {
		return nil, errors.New("jwt verification failed: token is not valid")
	}

	return map[string]any(claims), nil
}

// stripBearerPrefix removes a case-insensitive "Bearer " prefix from an
// Authorization header value, if present.
func stripBearerPrefix(authorization string) string {
	const prefix = "bearer "
	if len(authorization) > len(prefix) && strings.EqualFold(authorization[:len(prefix)], prefix) {
		return authorization[len(prefix):]
	}
	return authorization
}

// keyfn is the jwt.Keyfunc used by Verify to resolve the key for a parsed
// token. It is the enforcement point for algorithm-family pinning (the
// actual fix for GHSA-9gqw-h2rw-44wv): it identifies the concrete signing
// method type of the token and refuses to return a key at all unless that
// type belongs to v.family, so an HMAC secret is never handed back for an
// asymmetric token and a public key is never handed back for an HMAC
// token. Any other method, including SigningMethodNone, falls through to
// the default case and is rejected.
func (v *JWTVerifier) keyfn(token *jwt.Token) (any, error) {
	switch token.Method.(type) {
	case *jwt.SigningMethodHMAC:
		if v.family != keyFamilyHMAC {
			return nil, fmt.Errorf("unexpected alg %q for the configured key family", token.Method.Alg())
		}
		return v.hmacSecret, nil
	case *jwt.SigningMethodRSA, *jwt.SigningMethodRSAPSS, *jwt.SigningMethodECDSA:
		if v.family != keyFamilyAsymmetric {
			return nil, fmt.Errorf("unexpected alg %q for the configured key family", token.Method.Alg())
		}
		if v.jwks != nil {
			return v.jwks.Keyfunc(token)
		}
		return v.publicKey, nil
	default:
		return nil, fmt.Errorf("unexpected alg %q", token.Method.Alg())
	}
}
