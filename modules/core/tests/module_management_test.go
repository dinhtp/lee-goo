package tests

import (
	"testing"

	"github.com/stretchr/testify/require"

	serviceModule "github.com/dinhtp/lee-goo/modules/core/internal/service/module"
	domainModule "github.com/dinhtp/lee-goo/modules/core/internal/domain/module"
)

// TestTopologicalSort_Order verifies that dependencies are ordered before dependents.
func TestTopologicalSort_Order(t *testing.T) {
	// Graph: C depends on B, B depends on A → expected order: A, B, C
	graph := map[string][]string{
		"a": {},
		"b": {"a"},
		"c": {"b"},
	}

	sorted, err := serviceModule.TopologicalSort(graph)
	require.NoError(t, err)
	require.Equal(t, 3, len(sorted))

	pos := make(map[string]int, len(sorted))
	for i, name := range sorted {
		pos[name] = i
	}

	require.Less(t, pos["a"], pos["b"], "a must come before b")
	require.Less(t, pos["b"], pos["c"], "b must come before c")
}

// TestTopologicalSort_DetectsCycle verifies that a cycle returns ErrCircularDependency.
func TestTopologicalSort_DetectsCycle(t *testing.T) {
	// A → B → C → A  (cycle)
	graph := map[string][]string{
		"a": {"c"},
		"b": {"a"},
		"c": {"b"},
	}

	_, err := serviceModule.TopologicalSort(graph)
	require.ErrorIs(t, err, domainModule.ErrCircularDependency)
}

// TestTopologicalSort_NoDependencies verifies a flat graph with no edges.
func TestTopologicalSort_NoDependencies(t *testing.T) {
	graph := map[string][]string{
		"alpha": {},
		"beta":  {},
		"gamma": {},
	}

	sorted, err := serviceModule.TopologicalSort(graph)
	require.NoError(t, err)
	require.ElementsMatch(t, []string{"alpha", "beta", "gamma"}, sorted)
}

// TestTopologicalSort_DiamondDependency verifies correct ordering for a diamond shape.
// D depends on B and C; B and C both depend on A.
func TestTopologicalSort_DiamondDependency(t *testing.T) {
	graph := map[string][]string{
		"a": {},
		"b": {"a"},
		"c": {"a"},
		"d": {"b", "c"},
	}

	sorted, err := serviceModule.TopologicalSort(graph)
	require.NoError(t, err)
	require.Equal(t, 4, len(sorted))

	pos := make(map[string]int, len(sorted))
	for i, name := range sorted {
		pos[name] = i
	}

	require.Less(t, pos["a"], pos["b"], "a must come before b")
	require.Less(t, pos["a"], pos["c"], "a must come before c")
	require.Less(t, pos["b"], pos["d"], "b must come before d")
	require.Less(t, pos["c"], pos["d"], "c must come before d")
}

// TestTopologicalSort_SelfLoop verifies a self-referencing node is detected as a cycle.
func TestTopologicalSort_SelfLoop(t *testing.T) {
	graph := map[string][]string{
		"a": {"a"},
	}

	_, err := serviceModule.TopologicalSort(graph)
	require.ErrorIs(t, err, domainModule.ErrCircularDependency)
}
