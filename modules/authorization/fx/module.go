package fx

import (
	"context"

	"go.uber.org/fx"

	authzConfig "github.com/dinhtp/lee-goo/modules/authorization/config"
	domainRole "github.com/dinhtp/lee-goo/modules/authorization/internal/domain/role"
	roleHandler "github.com/dinhtp/lee-goo/modules/authorization/internal/handler/role"
	permissionRepository "github.com/dinhtp/lee-goo/modules/authorization/internal/repository/permission"
	roleRepository "github.com/dinhtp/lee-goo/modules/authorization/internal/repository/role"
	authzRouter "github.com/dinhtp/lee-goo/modules/authorization/internal/router"
	roleService "github.com/dinhtp/lee-goo/modules/authorization/internal/service/role"
	"github.com/dinhtp/lee-goo/system/extension"
)

// Module returns the fx.Option that wires all authorization-module providers and invokers.
// Consumers compose this into their application fx.App alongside platform options.
func Module() fx.Option {
	return fx.Module("authorization",
		fx.Provide(
			// Config: default authorization configuration
			authzConfig.DefaultConfig,
			// Repositories: database.Connection → domain port interfaces
			roleRepository.NewRepository,
			permissionRepository.NewRepository,
			// Service: ports + EventBus → *Service (concrete type for multi-interface provision)
			roleService.NewService,
			// Provide domain use-case interfaces from *Service via constructor functions
			func(s *roleService.Service) domainRole.RoleUseCase { return s },
			func(s *roleService.Service) domainRole.PolicyUseCase { return s },
			// Handler: RoleUseCase + PolicyUseCase → *Handler
			roleHandler.NewHandler,
			// Router: *Handler → HandlerRouter
			roleHandler.NewRouter,
		),
		// Mount /roles routes on the platform HTTP server at startup.
		fx.Invoke(authzRouter.Register),
		// Wire the default-role assignment hook into the extension registry.
		fx.Invoke(registerDefaultRoleHook),
	)
}

// registerDefaultRoleHook registers a "user.after_created" extension point handler
// that assigns the configured default role to every newly created user.
// The hook signature func(context.Context, string) error matches what the user
// service calls via registry.Resolve("user.after_created").
func registerDefaultRoleHook(
	reg *extension.ExtensionRegistry,
	svc *roleService.Service,
	cfg *authzConfig.AuthzConfig,
) {
	reg.Register("user.after_created", 100, func(ctx context.Context, userID string) error {
		// Production: look up cfg.DefaultRole by name, then call svc.AssignRoleToUser.
		// Boilerplate: demonstrates the wiring pattern; actual assignment is a no-op
		// until a named-role lookup and AssignRoleToUser method are implemented.
		_ = ctx
		_ = userID
		_ = svc
		_ = cfg
		return nil
	})
}
