package config

import (
	"fmt"
	"time"
)

// DatabaseConfig holds PostgreSQL connection and pool settings.
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`

	// Pool settings — zero value means "use driver default".
	MaxOpenConnections int           `mapstructure:"max_open_connections"`
	MaxIdleConnections int           `mapstructure:"max_idle_connections"`
	ConnectionMaxTime  time.Duration `mapstructure:"connection_max_time"`
	ConnectionIdleTime time.Duration `mapstructure:"connection_idle_time"`
}

// DSN builds the PostgreSQL connection string from the database config fields.
func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.DBName, d.SSLMode,
	)
}
