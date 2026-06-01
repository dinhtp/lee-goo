package role

import "context"

// RoleUseCase exposes role management operations.
type RoleUseCase interface {
	CreateRole(ctx context.Context, name, description string) (*Role, error)
	AssignRoleToUser(ctx context.Context, userID, roleID string) error
	RemoveRoleFromUser(ctx context.Context, userID, roleID string) error
	ListRoles(ctx context.Context) ([]Role, error)
}

// PolicyUseCase exposes authorization policy evaluation.
type PolicyUseCase interface {
	CanDo(ctx context.Context, userID, action, resource string) (bool, error)
	GetUserPermissions(ctx context.Context, userID string) ([]Permission, error)
}
