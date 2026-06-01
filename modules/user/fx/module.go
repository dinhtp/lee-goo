package fx

import (
	"go.uber.org/fx"

	userHandler "github.com/dinhtp/lee-goo/modules/user/internal/handler/user"
	userRepository "github.com/dinhtp/lee-goo/modules/user/internal/repository/user"
	internalRouter "github.com/dinhtp/lee-goo/modules/user/internal/router"
	userService "github.com/dinhtp/lee-goo/modules/user/internal/service/user"
)

// Module returns the fx.Option that wires all user-module providers and invokers.
// Consumers compose this into their application fx.App alongside platform options.
func Module() fx.Option {
	return fx.Module("user",
		fx.Provide(
			// Repository: database.Connection → domainUser.UserPort
			userRepository.NewRepository,
			// Service: UserPort + EventBus + ExtensionRegistry → domainUser.UseCase
			userService.NewService,
			// Adapter: domainUser.UseCase → contracts.UserService (public interface)
			userService.NewUserServiceAdapter,
			// Handler: domainUser.UseCase → *Handler
			userHandler.NewHandler,
			// Router: *Handler → internalRouter.HandlerRouter
			userHandler.NewRouter,
		),
		// Register routes on the platform HTTP server at startup.
		fx.Invoke(internalRouter.Register),
	)
}
