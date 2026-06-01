package cli_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

// projectRoot resolves the repository root from this file's location.
func projectRoot() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..")
}

// TestModuleCLI_List verifies that `go run . module list` exits 0.
func TestModuleCLI_List(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "module", "list")
	cmd.Dir = projectRoot()
	cmd.Env = os.Environ()

	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "module list failed: %s", string(out))
	require.NotEmpty(t, string(out))
}

// TestModuleCLI_Doctor verifies that `go run . module doctor` exits 0.
func TestModuleCLI_Doctor(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "module", "doctor")
	cmd.Dir = projectRoot()
	cmd.Env = os.Environ()

	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "module doctor failed: %s", string(out))
}

// TestModuleCLI_Graph verifies that `go run . module graph` exits 0.
func TestModuleCLI_Graph(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "module", "graph")
	cmd.Dir = projectRoot()
	cmd.Env = os.Environ()

	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "module graph failed: %s", string(out))
}
