package migrator

import (
	"fmt"

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
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		p.Config.Database.Host,
		p.Config.Database.Port,
		p.Config.Database.User,
		p.Config.Database.Password,
		p.Config.Database.DBName,
		p.Config.Database.SSLMode,
	)
	return NewRunner(dsn, p.Logger, p.Sources)
}

// FxOptions returns the fx module for the migration runner.
// Collects all Sources registered under the "migration.sources" group and provides *Runner.
func FxOptions() fx.Option {
	return fx.Options(
		fx.Provide(newRunner),
	)
}
