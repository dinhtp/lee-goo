package fx

import (
	"go.uber.org/fx"

	authConfig "github.com/dinhtp/lee-goo/modules/authentication/config"
	authHandler "github.com/dinhtp/lee-goo/modules/authentication/internal/handler/auth"
	authRouter "github.com/dinhtp/lee-goo/modules/authentication/internal/router"
	authService "github.com/dinhtp/lee-goo/modules/authentication/internal/service/auth"
)

// Module returns the fx.Option that wires all authentication-module providers and invokers.
// Consumers compose this into their application fx.App alongside platform options.
func Module() fx.Option {
	return fx.Module("authentication",
		fx.Provide(
			// Config: static defaults for token TTLs
			authConfig.DefaultConfig,
			// Service: UserService + Signer + Verifier + EventBus + AuthConfig → domainAuth.UseCase
			authService.NewService,
			// Adapter: domainAuth.UseCase → contracts.AuthService (public interface)
			authService.NewAuthServiceAdapter,
			// Handler: domainAuth.UseCase → *Handler
			authHandler.NewHandler,
			// Router: *Handler → internalRouter.HandlerRouter
			authHandler.NewRouter,
		),
		// Register routes on the platform HTTP server at startup.
		fx.Invoke(authRouter.Register),
	)
}
