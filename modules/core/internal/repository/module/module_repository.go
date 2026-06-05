package module

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	domainModule "github.com/dinhtp/lee-goo/modules/core/internal/domain/module"
)

// moduleRepository implements domainModule.ModulePort using sqlx.
type moduleRepository struct {
	db *sqlx.DB
}

// compile-time interface check
var _ domainModule.ModulePort = (*moduleRepository)(nil)

// NewRepository constructs a moduleRepository from a *sqlx.DB.
func NewRepository(db *sqlx.DB) domainModule.ModulePort {
	return &moduleRepository{db: db}
}

const findAllQuery = `
SELECT name, version, status, path, checksum,
       installed_at, enabled_at, disabled_at,
       uninstalled_at,
       created_at, updated_at
FROM modules
ORDER BY name`

// FindAll returns all modules ordered by name.
func (r *moduleRepository) FindAll(ctx context.Context) ([]domainModule.Module, error) {
	rows, err := r.db.QueryContext(ctx, findAllQuery)
	if err != nil {
		return nil, fmt.Errorf("moduleRepository.FindAll: %w", err)
	}
	defer rows.Close()

	var modules []domainModule.Module
	for rows.Next() {
		m, err := scanModule(rows)
		if err != nil {
			return nil, fmt.Errorf("moduleRepository.FindAll scan: %w", err)
		}
		modules = append(modules, m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("moduleRepository.FindAll rows: %w", err)
	}
	return modules, nil
}

const findByNameQuery = `
SELECT name, version, status, path, checksum,
       installed_at, enabled_at, disabled_at,
       uninstalled_at,
       created_at, updated_at
FROM modules
WHERE name = $1`

// FindByName retrieves a single module by its unique name.
func (r *moduleRepository) FindByName(ctx context.Context, name string) (*domainModule.Module, error) {
	row := r.db.QueryRowContext(ctx, findByNameQuery, name)
	m, err := scanModule(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("moduleRepository.FindByName: %w", err)
	}
	return &m, nil
}

const upsertQuery = `
INSERT INTO modules (name, version, status, path, checksum, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
ON CONFLICT (name) DO UPDATE
SET version    = EXCLUDED.version,
    status     = EXCLUDED.status,
    path       = EXCLUDED.path,
    checksum   = EXCLUDED.checksum,
    updated_at = NOW()`

// Upsert inserts or updates a module record.
func (r *moduleRepository) Upsert(ctx context.Context, m domainModule.Module) error {
	_, err := r.db.ExecContext(ctx, upsertQuery,
		m.Name, m.Version, string(m.Status), m.Path, m.Checksum,
	)
	if err != nil {
		return fmt.Errorf("moduleRepository.Upsert: %w", err)
	}
	return nil
}

const updateStatusQuery = `UPDATE modules SET status = $2, updated_at = NOW() WHERE name = $1`

// UpdateStatus sets a module's status field without touching other columns.
func (r *moduleRepository) UpdateStatus(ctx context.Context, name string, status domainModule.Status) error {
	_, err := r.db.ExecContext(ctx, updateStatusQuery, name, string(status))
	if err != nil {
		return fmt.Errorf("moduleRepository.UpdateStatus: %w", err)
	}
	return nil
}

// scanner is satisfied by both *sql.Row and *sql.Rows so scanModule can serve both.
type scanner interface {
	Scan(dest ...any) error
}

func scanModule(s scanner) (domainModule.Module, error) {
	var m domainModule.Module
	var status string
	var installedAt, enabledAt, disabledAt, uninstalledAt *time.Time
	err := s.Scan(
		&m.Name, &m.Version, &status, &m.Path, &m.Checksum,
		&installedAt, &enabledAt, &disabledAt,
		&uninstalledAt,
		&m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		return domainModule.Module{}, err
	}
	m.Status = domainModule.Status(status)
	m.InstalledAt = installedAt
	m.EnabledAt = enabledAt
	m.DisabledAt = disabledAt
	m.UninstalledAt = uninstalledAt
	return m, nil
}
