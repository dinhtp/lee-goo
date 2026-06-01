// Package module implements the module management use case.
package module

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/dinhtp/lee-goo/system/eventbus"
	domainModule "github.com/dinhtp/lee-goo/modules/core/internal/domain/module"
)

// protectedModules cannot be modified by end-users.
var protectedModules = map[string]struct{}{
	"module":         {},
	"user":           {},
	"authentication": {},
	"authorization":  {},
}

// service is the concrete UseCase implementation.
type service struct {
	repo      domainModule.ModulePort
	eventBus  eventbus.EventBus
	workspace string // absolute path to workspace root
}

// compile-time interface check
var _ domainModule.UseCase = (*service)(nil)

// NewService constructs a service wired to the given repository and event bus.
// The workspace root is resolved via os.Getwd() at construction time.
func NewService(repo domainModule.ModulePort, bus eventbus.EventBus) domainModule.UseCase {
	wd, _ := os.Getwd()
	return &service{
		repo:      repo,
		eventBus:  bus,
		workspace: wd,
	}
}

// Status returns the persisted state of a single module.
func (s *service) Status(ctx context.Context, name string) (*domainModule.Module, error) {
	m, err := s.repo.FindByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("service.Status: %w", err)
	}
	return m, nil
}

// Sync reconciles the DB with discovered manifests: modules found on disk
// but not in the DB are upserted with status=discovered.
func (s *service) Sync(ctx context.Context) error {
	manifests, err := s.discoverManifests()
	if err != nil {
		return fmt.Errorf("service.Sync discover: %w", err)
	}

	for _, mf := range manifests {
		existing, err := s.repo.FindByName(ctx, mf.Name)
		if err != nil {
			return fmt.Errorf("service.Sync find %s: %w", mf.Name, err)
		}
		if existing != nil {
			continue // already tracked
		}
		modPath := filepath.Join(s.workspace, "modules", mf.Name)
		m := domainModule.Module{
			Name:    mf.Name,
			Version: mf.Version,
			Status:  domainModule.StatusDiscovered,
			Path:    modPath,
		}
		if err := s.repo.Upsert(ctx, m); err != nil {
			return fmt.Errorf("service.Sync upsert %s: %w", mf.Name, err)
		}
	}
	return nil
}

// Doctor checks for known issues: missing manifests, broken dependency refs.
func (s *service) Doctor(ctx context.Context) ([]string, error) {
	var issues []string

	manifests, err := s.discoverManifests()
	if err != nil {
		return nil, fmt.Errorf("service.Doctor: %w", err)
	}

	known := make(map[string]struct{}, len(manifests))
	for _, mf := range manifests {
		known[mf.Name] = struct{}{}
	}

	for _, mf := range manifests {
		for _, dep := range mf.Dependencies.Required {
			if _, ok := known[dep]; !ok {
				issues = append(issues, fmt.Sprintf(
					"module %q requires %q but it was not found on disk", mf.Name, dep,
				))
			}
		}
	}
	return issues, nil
}

// discoverManifests reads all modules/*/module.yaml files under the workspace.
func (s *service) discoverManifests() ([]domainModule.Manifest, error) {
	modulesDir := filepath.Join(s.workspace, "modules")
	entries, err := os.ReadDir(modulesDir)
	if err != nil {
		return nil, fmt.Errorf("discoverManifests ReadDir: %w", err)
	}

	var manifests []domainModule.Manifest
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		yamlPath := filepath.Join(modulesDir, entry.Name(), "module.yaml")
		if err := s.validatePathSafety(yamlPath); err != nil {
			continue // skip unsafe paths silently
		}
		data, err := os.ReadFile(yamlPath)
		if err != nil {
			continue // manifest absent — not an error
		}
		var mf domainModule.Manifest
		if err := yaml.Unmarshal(data, &mf); err != nil {
			continue // malformed manifest — skip
		}
		manifests = append(manifests, mf)
	}
	return manifests, nil
}

// validatePathSafety ensures a path is beneath the workspace modules/ directory
// and contains no path-traversal sequences.
func (s *service) validatePathSafety(path string) error {
	clean := filepath.Clean(path)
	modulesRoot := filepath.Clean(filepath.Join(s.workspace, "modules"))

	rel, err := filepath.Rel(modulesRoot, clean)
	if err != nil {
		return domainModule.ErrInvalidModulePath
	}
	// rel must not start with ".." (escaping modules root)
	if len(rel) >= 2 && rel[:2] == ".." {
		return domainModule.ErrInvalidModulePath
	}
	return nil
}
