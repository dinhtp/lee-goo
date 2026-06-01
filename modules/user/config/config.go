package config

// UserConfig holds user-module-specific configuration values.
type UserConfig struct {
	AllowRegistration bool `mapstructure:"allow_registration"`
}

// DefaultConfig returns sensible defaults for the user module.
func DefaultConfig() *UserConfig {
	return &UserConfig{AllowRegistration: true}
}
