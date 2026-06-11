package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"

	"github.com/dinhtp/lee-goo/static"
)

// Config holds all application configuration loaded from environment variables.
type Config struct {
	Env      string         `mapstructure:"env"`
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Auth     AuthConfig     `mapstructure:"auth"`
}

// IsDevelopment returns true when running in a local or dev environment.
func (c *Config) IsDevelopment() bool {
	return c.Env == static.EnvLocal || c.Env == static.EnvDev
}

// Load reads configuration from environment variables (and optionally a .env file).
// The path parameter is optional — pass an empty string to skip file loading.
func Load() (*Config, error) {
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Attempt to read .env file if present; errors are ignored (file is optional).
	v.SetConfigFile(".env")
	v.SetConfigType("env")
	_ = v.ReadInConfig()

	// Sensible defaults so the app starts without full configuration.
	v.SetDefault("env", "local")
	v.SetDefault("server.port", "8080")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.sslmode", "disable")
	v.SetDefault("database.max_open_connections", 4)
	v.SetDefault("database.max_idle_connections", 2)
	v.SetDefault("auth.access_token_ttl", "15m")
	v.SetDefault("auth.refresh_token_ttl", "168h")
	v.SetDefault("server.log_level", "info")

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("config unmarshal: %w", err)
	}
	return &cfg, nil
}
