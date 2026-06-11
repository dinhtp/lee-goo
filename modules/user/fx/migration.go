package fx

import (
	"go.uber.org/fx"

	"github.com/dinhtp/lee-goo/modules/user/migrations"
	"github.com/dinhtp/lee-goo/system/migrator"
)

// Source returns the user module's migration source as a plain value.
func Source() migrator.Source {
	return migrator.Source{Name: "user", FS: migrations.FS, Path: "."}
}

// MigrationSource returns the fx.Option that registers the user module's
// embedded SQL migrations under the "migration.sources" group.
func MigrationSource() fx.Option {
	return fx.Provide(
		fx.Annotate(
			func() migrator.Source { return Source() },
			fx.ResultTags(`group:"migration.sources"`),
		),
	)
}
