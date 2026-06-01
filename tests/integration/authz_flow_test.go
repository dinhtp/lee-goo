package integration_test

import (
	"testing"

	authzFx "github.com/dinhtp/lee-goo/modules/authorization/fx"
	moduleFx "github.com/dinhtp/lee-goo/modules/core/fx"
	userFx "github.com/dinhtp/lee-goo/modules/user/fx"
	"github.com/dinhtp/lee-goo/pkg/testapp"
)

// TestAuthzFlowSkipsWithoutDB verifies that the authorization module boots
// and registers the default-role extension hook via the fx dependency graph.
func TestAuthzFlowSkipsWithoutDB(t *testing.T) {
	skipIfNoDB(t)

	app := testapp.New(
		testapp.WithModules(
			moduleFx.Module(),
			userFx.Module(),
			authzFx.Module(),
		),
	)
	app.Start(t)
}
