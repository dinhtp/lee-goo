package migrator

import (
	"go.uber.org/fx"

	"github.com/dinhtp/lee-goo/system/config"
	"github.com/dinhtp/lee-goo/system/logger"
)

type params struct {
	fx.In
	Config  *config.Config
	Logger  *logger.Logger
	Sources []Source `group:"migration.sources"`
}

func newRunner(p params) *Runner {
	return NewRunner(p.Config.Database.DSN(), p.Logger, p.Sources)
}

// FxOptions returns the fx module for the migration runner.
// Collects all Sources registered under the "migration.sources" group and provides *Runner.
func FxOptions() fx.Option {
	return fx.Options(
		fx.Provide(newRunner),
	)
}
