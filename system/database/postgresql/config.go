// Package postgresql provides a PostgreSQL database connection backed by sqlx (github.com/jmoiron/sqlx).
// It exposes a Connection interface with a standard lifecycle: NewConnection → Open → Instance → Close.
// Use the Transact helper for transactional operations.
package postgresql

import (
	"errors"
	"time"
)

var (
	// ErrMissingConfig is returned by Open when Config is nil.
	ErrMissingConfig = errors.New("database config is missing")
	// ErrUninitializedDatabase is returned by Instance, Close, and Ping when Open has not
	// been called successfully yet.
	ErrUninitializedDatabase = errors.New("database instance is not initialized")
)

// Config holds the PostgreSQL connection pool configuration.
// Pool settings are applied only when their value is greater than zero; zero means
// "not applied" (the underlying database/sql pool uses its own defaults).
type Config struct {
	// ConnectionMaxTime sets the maximum lifetime of a connection. Zero means no limit.
	ConnectionMaxTime time.Duration
	// ConnectionIdleTime sets the maximum time a connection may remain idle. Zero means no limit.
	ConnectionIdleTime time.Duration
	// MaxIdleConnections sets the maximum number of idle connections. Default: 2.
	MaxIdleConnections int
	// MaxOpenConnections sets the maximum number of open connections. Default: 4.
	MaxOpenConnections int
}

// newDefaultConfig returns a default pool config.
func newDefaultConfig() *Config {
	return &Config{
		MaxIdleConnections: 2,
		MaxOpenConnections: 4,
	}
}
