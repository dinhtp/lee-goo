package contracts

import "context"

// RoleService is the public interface for role assignment operations.
// Consumed by other modules (e.g. authentication middleware).
type RoleService interface {
	AssignRoleToUser(ctx context.Context, userID, roleID string) error
	RemoveRoleFromUser(ctx context.Context, userID, roleID string) error
}

// PolicyService is the public interface for authorization checks.
type PolicyService interface {
	CanDo(ctx context.Context, userID, action, resource string) (bool, error)
}
