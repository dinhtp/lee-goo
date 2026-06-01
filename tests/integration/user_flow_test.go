package integration_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	moduleFx "github.com/dinhtp/lee-goo/modules/core/fx"
	authzFx "github.com/dinhtp/lee-goo/modules/authorization/fx"
	userFx "github.com/dinhtp/lee-goo/modules/user/fx"
	"github.com/dinhtp/lee-goo/pkg/testapp"
)

// skipIfNoDB skips the test when no test database is configured.
func skipIfNoDB(t *testing.T) {
	t.Helper()
	if os.Getenv("TEST_DATABASE_URL") == "" && os.Getenv("DATABASE_HOST") == "" {
		t.Skip("no test DB: set DATABASE_HOST or TEST_DATABASE_URL to run integration tests")
	}
}

// TestUserCreateFlow verifies that the user module boots and integrates with
// the authorization module via the fx dependency graph.
func TestUserCreateFlow(t *testing.T) {
	skipIfNoDB(t)

	app := testapp.New(
		testapp.WithModules(
			moduleFx.Module(),
			userFx.Module(),
			authzFx.Module(),
		),
	)
	app.Start(t)

	require.NotNil(t, app)
	// DB-dependent assertions go here after migration setup.
}
