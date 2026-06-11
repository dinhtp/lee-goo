package migrator

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"sort"

	"github.com/dinhtp/lee-goo/system/logger"
	"github.com/golang-migrate/migrate/v4"
	pg "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Source describes one module's embedded SQL migration files.
type Source struct {
	Name string // module name; used for schema_migrations_<name> table
	FS   fs.FS  // embedded SQL files (//go:embed *.sql)
	Path string // directory within FS containing *.sql files (usually ".")
}

// Runner executes migrations for registered Sources.
type Runner struct {
	dsn     string
	sources []Source
	logger  *logger.Logger
}

// NewRunner constructs a Runner. sources are sorted by Name for determinism.
func NewRunner(dsn string, logger *logger.Logger, sources []Source) *Runner {
	sorted := make([]Source, len(sources))
	copy(sorted, sources)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Name < sorted[j].Name })
	return &Runner{dsn: dsn, logger: logger, sources: sorted}
}

// UpFor runs UP migrations for the named module.
func (r *Runner) UpFor(ctx context.Context, name string) error {
	s, err := r.findSource(name)
	if err != nil {
		return err
	}
	return r.runUp(ctx, s)
}

// UpAll runs UP migrations for all registered sources in alphabetical order.
func (r *Runner) UpAll(ctx context.Context) error {
	for _, s := range r.sources {
		if err := r.runUp(ctx, s); err != nil {
			return err
		}
	}
	return nil
}

// DownFor runs DOWN (rollback) migrations for the named module.
func (r *Runner) DownFor(ctx context.Context, name string) error {
	s, err := r.findSource(name)
	if err != nil {
		return err
	}
	return r.runDown(ctx, s)
}

func (r *Runner) findSource(name string) (Source, error) {
	for _, s := range r.sources {
		if s.Name == name {
			return s, nil
		}
	}
	return Source{}, fmt.Errorf("migration source %q not found", name)
}

func (r *Runner) runUp(ctx context.Context, s Source) error {
	m, db, err := r.newMigrate(s)
	if err != nil {
		return fmt.Errorf("migrator %s: %w", s.Name, err)
	}
	defer func() { _, _ = m.Close(); _ = db.Close() }()
	stopOnCancel(ctx, m)

	r.logger.WithFields(map[string]any{
		"module": s.Name,
	}).Info("running UP migrations")
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrator %s up: %w", s.Name, err)
	}
	return nil
}

func (r *Runner) runDown(ctx context.Context, s Source) error {
	m, db, err := r.newMigrate(s)
	if err != nil {
		return fmt.Errorf("migrator %s: %w", s.Name, err)
	}
	defer func() { _, _ = m.Close(); _ = db.Close() }()
	stopOnCancel(ctx, m)

	r.logger.WithFields(map[string]any{
		"module": s.Name,
	}).Info("running DOWN migrations")
	if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrator %s down: %w", s.Name, err)
	}
	return nil
}

// stopOnCancel signals the migrator to stop gracefully when ctx is cancelled.
func stopOnCancel(ctx context.Context, m *migrate.Migrate) {
	go func() {
		<-ctx.Done()
		m.GracefulStop <- true
	}()
}

// newMigrate opens a dedicated single-connection *sql.DB for this source and returns
// both the Migrate instance and the DB. The caller must defer m.Close() then db.Close().
// MaxOpenConns(1) ensures the advisory lock used by golang-migrate stays on one connection.
func (r *Runner) newMigrate(s Source) (*migrate.Migrate, *sql.DB, error) {
	db, err := sql.Open("pgx", r.dsn)
	if err != nil {
		return nil, nil, fmt.Errorf("open migration db: %w", err)
	}
	db.SetMaxOpenConns(1)

	driver, err := pg.WithInstance(db, &pg.Config{
		MigrationsTable: "schema_migrations_" + s.Name,
	})
	if err != nil {
		_ = db.Close()
		return nil, nil, fmt.Errorf("postgres driver: %w", err)
	}

	src, err := iofs.New(s.FS, s.Path)
	if err != nil {
		_ = db.Close()
		return nil, nil, fmt.Errorf("iofs source: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", src, "postgres", driver)
	if err != nil {
		_ = driver.Close() // closes both the advisory-lock conn and the pool reference
		return nil, nil, err
	}
	return m, db, nil
}
