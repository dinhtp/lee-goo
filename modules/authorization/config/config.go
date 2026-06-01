package config

// AuthzConfig holds configuration values for the authorization module.
type AuthzConfig struct {
	// DefaultRole is assigned to new users when no explicit role is provided.
	DefaultRole string `mapstructure:"default_role"`
}

// DefaultConfig returns sensible defaults for the authorization module.
func DefaultConfig() *AuthzConfig {
	return &AuthzConfig{DefaultRole: "user"}
}
