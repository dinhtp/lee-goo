package module

import (
	"context"
	"fmt"
	"time"

	"github.com/dinhtp/lee-goo/modules/core/contracts"
	domainModule "github.com/dinhtp/lee-goo/modules/core/internal/domain/module"
)

// Discover scans the workspace for module.yaml manifests and returns them as
// Module values (not persisted — use Sync to persist).
func (s *service) Discover(ctx context.Context) ([]domainModule.Module, error) {
	manifests, err := s.discoverManifests()
	if err != nil {
		return nil, fmt.Errorf("service.Discover: %w", err)
	}
	modules := make([]domainModule.Module, 0, len(manifests))
	for _, mf := range manifests {
		modules = append(modules, domainModule.Module{
			Name:    mf.Name,
			Version: mf.Version,
			Status:  domainModule.StatusDiscovered,
		})
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

// Enable transitions an installed module to enabled status.
func (s *service) Enable(ctx context.Context, name string) error {
	if err := s.assertNotProtected(name); err != nil {
		return err
	}
	m, err := s.repo.FindByName(ctx, name)
	if err != nil {
		return fmt.Errorf("service.Enable find: %w", err)
	}
	if m == nil {
		return domainModule.ErrModuleNotInstalled
	}
	if m.Status != domainModule.StatusInstalled && m.Status != domainModule.StatusDisabled {
		return domainModule.ErrModuleNotInstalled
	}
	now := time.Now()
	m.Status = domainModule.StatusEnabled
	m.EnabledAt = &now
	if err := s.repo.Upsert(ctx, *m); err != nil {
		return fmt.Errorf("service.Enable upsert: %w", err)
	}
	_ = s.eventBus.Publish(ctx, "module.enabled", contracts.ModuleEnabledEvent{Name: name})
	return nil
}

// Disable transitions an enabled module to disabled status.
func (s *service) Disable(ctx context.Context, name string) error {
	if err := s.assertNotProtected(name); err != nil {
		return err
	}
	m, err := s.repo.FindByName(ctx, name)
	if err != nil {
		return fmt.Errorf("service.Disable find: %w", err)
	}
	if m == nil || m.Status != domainModule.StatusEnabled {
		return domainModule.ErrModuleNotEnabled
	}
	now := time.Now()
	m.Status = domainModule.StatusDisabled
	m.DisabledAt = &now
	if err := s.repo.Upsert(ctx, *m); err != nil {
		return fmt.Errorf("service.Disable upsert: %w", err)
	}
	_ = s.eventBus.Publish(ctx, "module.disabled", contracts.ModuleDisabledEvent{Name: name})
	return nil
}

// Upgrade marks a module as upgrading, then re-installs it at the new version.
func (s *service) Upgrade(ctx context.Context, name string) error {
	if err := s.assertNotProtected(name); err != nil {
		return err
	}
	m, err := s.repo.FindByName(ctx, name)
	if err != nil {
		return fmt.Errorf("service.Upgrade find: %w", err)
	}
	if m == nil {
		return domainModule.ErrModuleNotInstalled
	}
	// Mark upgrading
	if err := s.repo.UpdateStatus(ctx, name, domainModule.StatusUpgrading); err != nil {
		return fmt.Errorf("service.Upgrade status: %w", err)
	}
	// Re-discover manifest to pick up new version from disk
	manifests, err := s.discoverManifests()
	if err != nil {
		_ = s.repo.UpdateStatus(ctx, name, domainModule.StatusFailed)
		return fmt.Errorf("service.Upgrade discover: %w", err)
	}
	var newVersion string
	for _, mf := range manifests {
		if mf.Name == name {
			newVersion = mf.Version
			break
		}
	}
	now := time.Now()
	m.Version = newVersion
	m.Status = domainModule.StatusInstalled
	m.UpgradedAt = &now
	if err := s.repo.Upsert(ctx, *m); err != nil {
		return fmt.Errorf("service.Upgrade upsert: %w", err)
	}
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

// Remove marks a module as removed from codebase (source no longer on disk).
func (s *service) Remove(ctx context.Context, name string) error {
	if err := s.assertNotProtected(name); err != nil {
		return err
	}
	m, err := s.repo.FindByName(ctx, name)
	if err != nil {
		return fmt.Errorf("service.Remove find: %w", err)
	}
	if m == nil {
		return domainModule.ErrModuleNotFound
	}
	if m.Status != domainModule.StatusUninstalled {
		return fmt.Errorf("service.Remove: module must be uninstalled before removal")
	}
	now := time.Now()
	m.Status = domainModule.StatusRemoved
	m.RemovedFromCodebaseAt = &now
	return s.repo.Upsert(ctx, *m)
}

// assertNotProtected returns ErrProtectedModule for built-in modules.
func (s *service) assertNotProtected(name string) error {
	if _, ok := protectedModules[name]; ok {
		return domainModule.ErrProtectedModule
	}
	return nil
}

// assertNoDependents returns ErrModuleHasDependents if any enabled module
// declares name as a required dependency.
func (s *service) assertNoDependents(ctx context.Context, name string) error {
	graph, err := s.Graph(ctx)
	if err != nil {
		return err
	}
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
