package platformfx

import (
	"go.uber.org/fx"

	"github.com/dinhtp/lee-goo/system/config"
	"github.com/dinhtp/lee-goo/system/database/postgresql"
	"github.com/dinhtp/lee-goo/system/logger"
	"github.com/dinhtp/lee-goo/system/migrator"
)

// MigrateOptions returns a minimal fx.Option set for CLI migration commands.
// Excludes HTTP server, eventbus, extension, and security to keep CLI apps lean.
func MigrateOptions() fx.Option {
	return fx.Options(
		fx.Provide(config.Load),
		fx.Provide(logger.NewLogger),
		postgresql.FxOptions(),
		migrator.FxOptions(),
	)
}
