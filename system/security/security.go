package security

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"

	"github.com/dinhtp/lee-goo/system/config"
)

// Signer creates signed JWT tokens from claims.
type Signer interface {
	Sign(claims jwt.MapClaims) (string, error)
}

// Verifier parses and validates JWT tokens.
type Verifier interface {
	Verify(token string) (jwt.Claims, error)
}

// jwtSecurity implements both Signer and Verifier using HMAC-SHA256.
type jwtSecurity struct {
	secret []byte
}

// compile-time interface checks
var _ Signer = (*jwtSecurity)(nil)
var _ Verifier = (*jwtSecurity)(nil)

// NewJWTSecurity constructs a jwtSecurity instance from config.
// Returns an error if AUTH_JWT_SECRET is not set.
func NewJWTSecurity(cfg *config.Config) (*jwtSecurity, error) {
	secret := cfg.Auth.JWTSecret
	if secret == "" {
		return nil, fmt.Errorf("AUTH_JWT_SECRET must be set")
	}
	return &jwtSecurity{secret: []byte(secret)}, nil
}

// NewSigner extracts the Signer interface from a jwtSecurity instance.
func NewSigner(js *jwtSecurity) Signer { return js }

// NewVerifier extracts the Verifier interface from a jwtSecurity instance.
func NewVerifier(js *jwtSecurity) Verifier { return js }

// Sign creates a signed JWT string from the provided claims.
func (j *jwtSecurity) Sign(claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(j.secret)
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}
	return signed, nil
}

// Verify parses and validates a JWT string, returning its claims.
func (j *jwtSecurity) Verify(tokenStr string) (jwt.Claims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return j.secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("verify token: %w", err)
	}
	return token.Claims, nil
}
