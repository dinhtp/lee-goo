package module

import (
	"context"
	"fmt"

	domainModule "github.com/dinhtp/lee-goo/modules/core/internal/domain/module"
)

// Graph returns a dependency map: module name → list of its required dependencies.
func (s *service) Graph(ctx context.Context) (map[string][]string, error) {
	manifests, err := s.discoverManifests()
	if err != nil {
		return nil, fmt.Errorf("service.Graph: %w", err)
	}
	return s.buildDependencyGraph(manifests), nil
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
	// Ensure every referenced dependency also has an entry in the graph,
	// so the node count is consistent.
	for _, deps := range graph {
		for _, dep := range deps {
			if _, exists := graph[dep]; !exists {
				graph[dep] = nil
			}
		}
	}

	// in-degree[node] = number of unresolved dependencies of node.
	// A node with in-degree 0 has no pending deps and can be scheduled.
	inDegree := make(map[string]int, len(graph))
	for node := range graph {
		if _, exists := inDegree[node]; !exists {
			inDegree[node] = 0
		}
		for _, dep := range graph[node] {
			inDegree[node]++ // node depends on dep → node's in-degree goes up
			_ = dep
		}
	}

	// Build reverse adjacency: dep → list of nodes that depend on dep.
	// When dep is resolved, we decrement in-degree of each of its dependents.
	reverse := make(map[string][]string, len(graph))
	for node, deps := range graph {
		for _, dep := range deps {
			reverse[dep] = append(reverse[dep], node)
		}
	}

	// Seed queue with nodes that have no dependencies.
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

		// Resolve node: decrement in-degree of everything that depended on it.
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
