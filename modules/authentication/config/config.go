package config

import "time"

// AuthConfig holds tunable parameters for the authentication module.
type AuthConfig struct {
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

// DefaultConfig returns sensible production defaults:
// access tokens expire in 15 minutes, refresh tokens in 7 days.
func DefaultConfig() *AuthConfig {
	return &AuthConfig{
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 7 * 24 * time.Hour,
	}
}
