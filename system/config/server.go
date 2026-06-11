package config

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	Port string `mapstructure:"port"`
	// LogLevel is the minimum log level. Valid values (case-insensitive): DEBUG, INFO, WARN, ERROR, PANIC, FATAL.
	// Env var: SERVER_LOG_LEVEL.
	LogLevel string `mapstructure:"log_level"`
}
