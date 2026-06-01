package module

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	domainModule "github.com/dinhtp/lee-goo/modules/core/internal/domain/module"
)

// Migrate runs pending database migrations for the named module.
// The migration files are expected at <workspace>/modules/<name>/migrations/.
func (s *service) Migrate(ctx context.Context, name string) error {
	m, err := s.repo.FindByName(ctx, name)
	if err != nil {
		return fmt.Errorf("service.Migrate find: %w", err)
	}
	if m == nil {
		return domainModule.ErrModuleNotFound
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return fmt.Errorf("service.Migrate: DATABASE_URL not set")
	}

	migPath := filepath.Join(s.workspace, "modules", name, "migrations")
	sourceURL := "file://" + migPath

	mig, err := migrate.New(sourceURL, dsn)
	if err != nil {
		return fmt.Errorf("service.Migrate init for %s: %w", name, err)
	}
	defer mig.Close()

	if err := mig.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("service.Migrate up for %s: %w", name, err)
	}
	return nil
}

// MigrateAll runs migrations for every module in topological dependency order.
func (s *service) MigrateAll(ctx context.Context) error {
	manifests, err := s.discoverManifests()
	if err != nil {
		return fmt.Errorf("service.MigrateAll discover: %w", err)
	}

	graph := s.buildDependencyGraph(manifests)
	order, err := TopologicalSort(graph)
	if err != nil {
		return fmt.Errorf("service.MigrateAll sort: %w", err)
	}

	for _, name := range order {
		migPath := filepath.Join(s.workspace, "modules", name, "migrations")
		if _, statErr := os.Stat(migPath); os.IsNotExist(statErr) {
			continue // module has no migrations directory
		}
		if err := s.Migrate(ctx, name); err != nil {
			return fmt.Errorf("service.MigrateAll %s: %w", name, err)
		}
	}
	return nil
}
