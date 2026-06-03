package postgresql

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"go.uber.org/fx"

	"github.com/dinhtp/lee-goo/system/config"
)

// newConnection builds a DSN from cfg, applies pool settings, and opens the connection.
func newConnection(cfg *config.Config) (Connection, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	poolCfg := &Config{
		MaxOpenConnections: cfg.Database.MaxOpenConnections,
		MaxIdleConnections: cfg.Database.MaxIdleConnections,
		ConnectionMaxTime:  cfg.Database.ConnectionMaxTime,
		ConnectionIdleTime: cfg.Database.ConnectionIdleTime,
	}

	conn := NewConnection(dsn, poolCfg)
	if _, err := conn.Open(); err != nil {
		return nil, fmt.Errorf("database open: %w", err)
	}

	return conn, nil
}

// newDB extracts the *sqlx.DB from an open Connection for injection into repositories.
func newDB(conn Connection) (*sqlx.DB, error) {
	return conn.Instance()
}

// FxOptions returns the fx module for the PostgreSQL database connection.
// It provides both Connection (for lifecycle management) and *sqlx.DB (for repositories).
func FxOptions() fx.Option {
	return fx.Options(
		fx.Provide(newConnection),
		fx.Provide(newDB),
	)
}
