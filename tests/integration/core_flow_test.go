package integration_test

import (
	"testing"

	moduleFx "github.com/dinhtp/lee-goo/modules/core/fx"
	"github.com/dinhtp/lee-goo/pkg/testapp"
)

// TestModuleManagerFlowSkipsWithoutDB verifies that the module management
// module boots in isolation and wires correctly via the fx dependency graph.
func TestModuleManagerFlowSkipsWithoutDB(t *testing.T) {
	skipIfNoDB(t)

	app := testapp.New(
		testapp.WithModules(moduleFx.Module()),
	)
	app.Start(t)
}
