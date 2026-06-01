package tests

import "testing"

// TestUserModuleCompiles is a compile-time smoke test.
// If this package builds, all internal packages resolve correctly.
func TestUserModuleCompiles(t *testing.T) {
	t.Log("user module compiles")
}
