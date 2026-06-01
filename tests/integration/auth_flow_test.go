package integration_test

import (
	"testing"

	authFx "github.com/dinhtp/lee-goo/modules/authentication/fx"
	moduleFx "github.com/dinhtp/lee-goo/modules/core/fx"
	userFx "github.com/dinhtp/lee-goo/modules/user/fx"
	"github.com/dinhtp/lee-goo/pkg/testapp"
)

// TestAuthFlowSkipsWithoutDB verifies that the authentication module boots
// alongside the user module and wires correctly via the fx dependency graph.
func TestAuthFlowSkipsWithoutDB(t *testing.T) {
	skipIfNoDB(t)

	app := testapp.New(
		testapp.WithModules(
			moduleFx.Module(),
			userFx.Module(),
			authFx.Module(),
		),
	)
	app.Start(t)
}
