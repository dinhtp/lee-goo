package module

import (
	"context"
	"fmt"
	"time"

	"github.com/dinhtp/lee-goo/modules/core/contracts"
	domainModule "github.com/dinhtp/lee-goo/modules/core/internal/domain/module"
)

// Discover syncs manifests to the DB then returns all tracked modules.
func (s *service) Discover(ctx context.Context) ([]domainModule.Module, error) {
	_ = s.sync(ctx) // best-effort: register new manifests before listing
	modules, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("service.Discover: %w", err)
	}
	return modules, nil
}

// Install transitions a discovered module to installed status.
func (s *service) Install(ctx context.Context, name string) error {
	if err := s.assertNotProtected(name); err != nil {
		return err
	}
	m, err := s.findOrDiscover(ctx, name)
	if err != nil {
		return fmt.Errorf("service.Install find: %w", err)
	}
	now := time.Now()
	m.Status = domainModule.StatusInstalled
	m.InstalledAt = &now
	if err := s.repo.Upsert(ctx, *m); err != nil {
		return fmt.Errorf("service.Install upsert: %w", err)
	}
	_ = s.eventBus.Publish(ctx, "module.installed", contracts.ModuleInstalledEvent{
		Name: name, Version: m.Version,
	})
	return nil
}

// Uninstall transitions a module to uninstalled; force bypasses dependent checks.
func (s *service) Uninstall(ctx context.Context, name string, force bool) error {
	if err := s.assertNotProtected(name); err != nil {
		return err
	}
	if !force {
		if err := s.assertNoDependents(ctx, name); err != nil {
			return err
		}
	}
	m, err := s.repo.FindByName(ctx, name)
	if err != nil {
		return fmt.Errorf("service.Uninstall find: %w", err)
	}
	if m == nil {
		return domainModule.ErrModuleNotInstalled
	}
	now := time.Now()
	m.Status = domainModule.StatusUninstalled
	m.UninstalledAt = &now
	if err := s.repo.Upsert(ctx, *m); err != nil {
		return fmt.Errorf("service.Uninstall upsert: %w", err)
	}
	_ = s.eventBus.Publish(ctx, "module.uninstalled", contracts.ModuleUninstalledEvent{Name: name})
	return nil
}

// assertNotProtected returns ErrProtectedModule for built-in modules.
func (s *service) assertNotProtected(name string) error {
	if _, ok := protectedModules[name]; ok {
		return domainModule.ErrProtectedModule
	}
	return nil
}

// assertNoDependents returns ErrModuleHasDependents if any module declares name as a required dependency.
func (s *service) assertNoDependents(_ context.Context, name string) error {
	manifests, err := s.discoverManifests()
	if err != nil {
		return err
	}
	graph := s.buildDependencyGraph(manifests)
	for mod, deps := range graph {
		if mod == name {
			continue
		}
		for _, dep := range deps {
			if dep == name {
				return domainModule.ErrModuleHasDependents
			}
		}
	}
	return nil
}

// findOrDiscover looks up a module in the DB; if absent, discovers it from disk.
func (s *service) findOrDiscover(ctx context.Context, name string) (*domainModule.Module, error) {
	m, err := s.repo.FindByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if m != nil {
		return m, nil
	}
	// Not in DB yet — try disk
	manifests, err := s.discoverManifests()
	if err != nil {
		return nil, err
	}
	for _, mf := range manifests {
		if mf.Name == name {
			mod := domainModule.Module{
				Name:    mf.Name,
				Version: mf.Version,
				Status:  domainModule.StatusDiscovered,
			}
			return &mod, nil
		}
	}
	return nil, domainModule.ErrModuleNotFound
}
