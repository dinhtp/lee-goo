package config

// AuthConfig holds JWT authentication settings.
type AuthConfig struct {
	JWTSecret       string `mapstructure:"jwt_secret"`
	AccessTokenTTL  string `mapstructure:"access_token_ttl"`  // e.g. "15m"
	RefreshTokenTTL string `mapstructure:"refresh_token_ttl"` // e.g. "168h"
}
