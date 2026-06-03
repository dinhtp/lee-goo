package postgresql

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

// Connection is the lifecycle interface for a PostgreSQL database connection.
type Connection interface {
	// DataSourceName returns the DSN used to open this connection.
	DataSourceName() string
	// Open creates the underlying sqlx.DB, applies pool settings, and pings the server.
	// Returns ErrMissingConfig when Config is nil. On ping failure the connection is closed
	// before the error is returned.
	Open() (*sqlx.DB, error)
	// Close releases the underlying database/sql connection pool.
	// Returns ErrUninitializedDatabase when Open has not been called.
	Close() error
	// Instance returns the active sqlx.DB opened by Open.
	// Returns ErrUninitializedDatabase when Open has not been called.
	Instance() (*sqlx.DB, error)
	// Ping verifies the connection is alive by sending a ping to the PostgreSQL server.
	Ping() error
}

// connection is an implementation of the database Connection.
type connection struct {
	dsn      string
	config   *Config
	instance *sqlx.DB
}

// NewConnection constructs a Connection for the given DSN.
// When config is nil, a default config with MaxIdleConnections=2 and MaxOpenConnections=4 is used.
func NewConnection(dsn string, config *Config) Connection {
	if config == nil {
		config = newDefaultConfig()
	}

	return &connection{dsn: dsn, config: config}
}

// DataSourceName returns the data source connection string.
func (c *connection) DataSourceName() string {
	return c.dsn
}

// Open initializes and verifies a new database client.
func (c *connection) Open() (*sqlx.DB, error) {
	if c.config == nil {
		return nil, ErrMissingConfig
	}

	instance, err := sqlx.Open("pgx", c.dsn)
	if err != nil {
		return nil, err
	}

	if c.config.MaxOpenConnections > 0 {
		instance.SetMaxOpenConns(c.config.MaxOpenConnections)
	}

	if c.config.MaxIdleConnections > 0 {
		instance.SetMaxIdleConns(c.config.MaxIdleConnections)
	}

	if c.config.ConnectionMaxTime > 0 {
		instance.SetConnMaxLifetime(c.config.ConnectionMaxTime)
	}

	if c.config.ConnectionIdleTime > 0 {
		instance.SetConnMaxIdleTime(c.config.ConnectionIdleTime)
	}

	if err = instance.Ping(); err != nil {
		_ = instance.Close()
		return nil, err
	}

	c.instance = instance
	return c.instance, nil
}

// Close closes the current database client.
func (c *connection) Close() error {
	if c.instance == nil {
		return ErrUninitializedDatabase
	}

	return c.instance.Close()
}

// Instance returns the active PostgreSQL database client opened by Open.
func (c *connection) Instance() (*sqlx.DB, error) {
	if c.instance == nil {
		return nil, ErrUninitializedDatabase
	}

	return c.instance, nil
}

// Ping verifies if the current database client is active and healthy.
func (c *connection) Ping() error {
	instance, err := c.Instance()
	if err != nil {
		return err
	}

	return instance.Ping()
}
