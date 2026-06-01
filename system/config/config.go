package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config holds all application configuration loaded from environment variables.
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Auth     AuthConfig     `mapstructure:"auth"`
	Log      LogConfig      `mapstructure:"log"`
}

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	Port string `mapstructure:"port"`
}

// DatabaseConfig holds PostgreSQL connection settings.
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

// AuthConfig holds JWT authentication settings.
type AuthConfig struct {
	JWTSecret       string `mapstructure:"jwt_secret"`
	AccessTokenTTL  string `mapstructure:"access_token_ttl"`  // e.g. "15m"
	RefreshTokenTTL string `mapstructure:"refresh_token_ttl"` // e.g. "168h"
}

// LogConfig holds logging settings.
type LogConfig struct {
	Level string `mapstructure:"level"` // debug, info, warn, error
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
	v.SetDefault("server.port", "8080")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.sslmode", "disable")
	v.SetDefault("auth.access_token_ttl", "15m")
	v.SetDefault("auth.refresh_token_ttl", "168h")
	v.SetDefault("log.level", "info")

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("config unmarshal: %w", err)
	}
	return &cfg, nil
}
