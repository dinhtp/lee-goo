package fx

import (
	"go.uber.org/fx"

	"github.com/dinhtp/lee-goo/modules/core/migrations"
	"github.com/dinhtp/lee-goo/system/migrator"
)

// Source returns the core module's migration source as a plain value.
func Source() migrator.Source {
	return migrator.Source{Name: "core", FS: migrations.FS, Path: "."}
}

// MigrationSource returns the fx.Option that registers the core module's
// embedded SQL migrations under the "migration.sources" group.
func MigrationSource() fx.Option {
	return fx.Provide(
		fx.Annotate(
			func() migrator.Source { return Source() },
			fx.ResultTags(`group:"migration.sources"`),
		),
	)
}
