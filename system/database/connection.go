package database

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib" // register pgx as database/sql driver

	"github.com/dinhtp/lee-goo/system/config"
)

// Connection abstracts the database handle lifecycle.
type Connection interface {
	DB() *sql.DB
	Close() error
}

type connection struct {
	db *sql.DB
}

// compile-time interface check
var _ Connection = (*connection)(nil)

// NewConnection opens a PostgreSQL connection using pgx stdlib driver.
func NewConnection(cfg *config.Config) (Connection, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("database open: %w", err)
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("database ping: %w", err)
	}

	return &connection{db: db}, nil
}

func (c *connection) DB() *sql.DB { return c.db }

func (c *connection) Close() error { return c.db.Close() }
