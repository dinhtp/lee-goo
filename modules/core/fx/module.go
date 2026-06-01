// Package fx wires the module management system into the fx dependency graph.
package fx

import (
	"go.uber.org/fx"

	moduleHandler "github.com/dinhtp/lee-goo/modules/core/internal/handler/module"
	moduleRepository "github.com/dinhtp/lee-goo/modules/core/internal/repository/module"
	moduleRouter "github.com/dinhtp/lee-goo/modules/core/internal/router"
	moduleService "github.com/dinhtp/lee-goo/modules/core/internal/service/module"
)

// Module returns the fx.Option that wires the complete module management package.
func Module() fx.Option {
	return fx.Module("module",
		fx.Provide(
			moduleRepository.NewRepository,
			moduleService.NewService,
			moduleHandler.NewHandler,
			// NewRouter returns *router (unexported type); annotate as HandlerRouter interface
			// so moduleRouter.Register receives the correct dependency type.
			fx.Annotate(
				moduleHandler.NewRouter,
				fx.As(new(moduleRouter.HandlerRouter)),
			),
		),
		fx.Invoke(moduleRouter.Register),
	)
}
