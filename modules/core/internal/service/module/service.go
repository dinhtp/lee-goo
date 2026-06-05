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

// sync reconciles the DB with discovered manifests: modules found on disk
// but not in the DB are upserted with status=discovered.
func (s *service) sync(ctx context.Context) error {
	manifests, err := s.discoverManifests()
	if err != nil {
		return fmt.Errorf("service.sync discover: %w", err)
	}

	for _, mf := range manifests {
		existing, err := s.repo.FindByName(ctx, mf.Name)
		if err != nil {
			return fmt.Errorf("service.sync find %s: %w", mf.Name, err)
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
			return fmt.Errorf("service.sync upsert %s: %w", mf.Name, err)
		}
	}
	return nil
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

// buildDependencyGraph converts manifests into a name→deps adjacency map.
func (s *service) buildDependencyGraph(manifests []domainModule.Manifest) map[string][]string {
	graph := make(map[string][]string, len(manifests))
	for _, mf := range manifests {
		deps := make([]string, 0, len(mf.Dependencies.Required))
		deps = append(deps, mf.Dependencies.Required...)
		graph[mf.Name] = deps
	}
	return graph
}

// TopologicalSort performs Kahn's algorithm on the dependency graph and returns
// modules in a safe installation order (dependencies before dependents).
// graph[node] = list of modules that node depends on.
// Returns ErrCircularDependency when a cycle is detected.
// Exported so tests can exercise it directly.
func TopologicalSort(graph map[string][]string) ([]string, error) {
	// Ensure every referenced dependency also has an entry in the graph.
	for _, deps := range graph {
		for _, dep := range deps {
			if _, exists := graph[dep]; !exists {
				graph[dep] = nil
			}
		}
	}

	// in-degree[node] = number of unresolved dependencies of node.
	inDegree := make(map[string]int, len(graph))
	for node := range graph {
		if _, exists := inDegree[node]; !exists {
			inDegree[node] = 0
		}
		for range graph[node] {
			inDegree[node]++
		}
	}

	// reverse adjacency: dep → nodes that depend on dep.
	reverse := make(map[string][]string, len(graph))
	for node, deps := range graph {
		for _, dep := range deps {
			reverse[dep] = append(reverse[dep], node)
		}
	}

	queue := make([]string, 0, len(inDegree))
	for node, deg := range inDegree {
		if deg == 0 {
			queue = append(queue, node)
		}
	}
	sortStrings(queue)

	var sorted []string
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		sorted = append(sorted, node)

		for _, dependent := range reverse[node] {
			inDegree[dependent]--
			if inDegree[dependent] == 0 {
				queue = append(queue, dependent)
				sortStrings(queue)
			}
		}
	}

	if len(sorted) != len(graph) {
		return nil, domainModule.ErrCircularDependency
	}
	return sorted, nil
}

// sortStrings sorts a string slice in place (insertion sort — small N).
func sortStrings(s []string) {
	for i := 1; i < len(s); i++ {
		key := s[i]
		j := i - 1
		for j >= 0 && s[j] > key {
			s[j+1] = s[j]
			j--
		}
		s[j+1] = key
	}
}
