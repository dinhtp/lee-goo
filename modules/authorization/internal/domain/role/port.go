package role

import "context"

// RolePort defines persistence operations for Role entities.
type RolePort interface {
	Create(ctx context.Context, r Role) error
	FindByID(ctx context.Context, id string) (*Role, error)
	List(ctx context.Context) ([]Role, error)
	AssignToUser(ctx context.Context, userID, roleID string) error
	RemoveFromUser(ctx context.Context, userID, roleID string) error
	FindByUserID(ctx context.Context, userID string) ([]Role, error)
}

// PermissionPort defines persistence operations for Permission entities.
type PermissionPort interface {
	Create(ctx context.Context, p Permission) error
	AssignToRole(ctx context.Context, roleID, permissionID string) error
	FindByRoleIDs(ctx context.Context, roleIDs []string) ([]Permission, error)
}
