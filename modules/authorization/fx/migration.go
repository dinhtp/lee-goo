package fx

import (
	"go.uber.org/fx"

	"github.com/dinhtp/lee-goo/modules/authorization/migrations"
	"github.com/dinhtp/lee-goo/system/migrator"
)

// MigrationSource returns the fx.Option that registers the authorization module's
// embedded SQL migrations under the "migration.sources" group.
func MigrationSource() fx.Option {
	return fx.Provide(
		fx.Annotate(
			func() migrator.Source {
				return migrator.Source{Name: "authorization", FS: migrations.FS, Path: "."}
			},
			fx.ResultTags(`group:"migration.sources"`),
		),
	)
}
