package fx

import (
	"go.uber.org/fx"

	"github.com/dinhtp/lee-goo/modules/user/migrations"
	"github.com/dinhtp/lee-goo/system/migrator"
)

// MigrationSource returns the fx.Option that registers the user module's
// embedded SQL migrations under the "migration.sources" group.
func MigrationSource() fx.Option {
	return fx.Provide(
		fx.Annotate(
			func() migrator.Source {
				return migrator.Source{Name: "user", FS: migrations.FS, Path: "."}
			},
			fx.ResultTags(`group:"migration.sources"`),
		),
	)
}
