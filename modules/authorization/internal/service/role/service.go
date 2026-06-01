package role

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/dinhtp/lee-goo/modules/authorization/contracts"
	domainRole "github.com/dinhtp/lee-goo/modules/authorization/internal/domain/role"
	"github.com/dinhtp/lee-goo/system/eventbus"
)

// Service implements RoleUseCase, PolicyUseCase, contracts.RoleService, and contracts.PolicyService.
type Service struct {
	rolePort domainRole.RolePort
	permPort domainRole.PermissionPort
	eventBus eventbus.EventBus
	// cache holds []domainRole.Permission keyed by userID to avoid repeated DB lookups.
	cache sync.Map
}

// compile-time interface checks
var _ domainRole.RoleUseCase = (*Service)(nil)
var _ domainRole.PolicyUseCase = (*Service)(nil)
var _ contracts.RoleService = (*Service)(nil)
var _ contracts.PolicyService = (*Service)(nil)

// NewService constructs the authorization Service with all required dependencies.
func NewService(
	rolePort domainRole.RolePort,
	permPort domainRole.PermissionPort,
	bus eventbus.EventBus,
) *Service {
	return &Service{rolePort: rolePort, permPort: permPort, eventBus: bus}
}

// CreateRole builds a new Role with a generated UUID, persists it, and returns it.
func (s *Service) CreateRole(ctx context.Context, name, description string) (*domainRole.Role, error) {
	r := domainRole.Role{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
	}
	if err := s.rolePort.Create(ctx, r); err != nil {
		return nil, err
	}
	return &r, nil
}

// AssignRoleToUser assigns a role to a user, invalidates the permission cache,
// and publishes the authorization.role_assigned event.
func (s *Service) AssignRoleToUser(ctx context.Context, userID, roleID string) error {
	if err := s.rolePort.AssignToUser(ctx, userID, roleID); err != nil {
		return err
	}
	s.cache.Delete(userID)
	_ = s.eventBus.Publish(ctx, "authorization.role_assigned", contracts.RoleAssignedEvent{
		UserID: userID,
		RoleID: roleID,
	})
	return nil
}

// RemoveRoleFromUser removes a role from a user, invalidates the permission cache,
// and publishes the authorization.role_removed event.
func (s *Service) RemoveRoleFromUser(ctx context.Context, userID, roleID string) error {
	if err := s.rolePort.RemoveFromUser(ctx, userID, roleID); err != nil {
		return err
	}
	s.cache.Delete(userID)
	_ = s.eventBus.Publish(ctx, "authorization.role_removed", contracts.RoleRemovedEvent{
		UserID: userID,
		RoleID: roleID,
	})
	return nil
}

// ListRoles returns all defined roles ordered by name.
func (s *Service) ListRoles(ctx context.Context) ([]domainRole.Role, error) {
	return s.rolePort.List(ctx)
}

// CanDo returns true when the user holds a permission matching action+resource.
func (s *Service) CanDo(ctx context.Context, userID, action, resource string) (bool, error) {
	perms, err := s.getUserPermissions(ctx, userID)
	if err != nil {
		return false, err
	}
	for _, p := range perms {
		if p.Action == action && p.Resource == resource {
			return true, nil
		}
	}
	return false, nil
}

// GetUserPermissions returns all permissions the user inherits through their roles.
func (s *Service) GetUserPermissions(ctx context.Context, userID string) ([]domainRole.Permission, error) {
	return s.getUserPermissions(ctx, userID)
}

// getUserPermissions is the shared implementation with in-memory caching.
func (s *Service) getUserPermissions(ctx context.Context, userID string) ([]domainRole.Permission, error) {
	if cached, ok := s.cache.Load(userID); ok {
		return cached.([]domainRole.Permission), nil
	}

	roles, err := s.rolePort.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	roleIDs := make([]string, len(roles))
	for i, r := range roles {
		roleIDs[i] = r.ID
	}

	if len(roleIDs) == 0 {
		return nil, nil
	}

	perms, err := s.permPort.FindByRoleIDs(ctx, roleIDs)
	if err != nil {
		return nil, err
	}

	s.cache.Store(userID, perms)
	return perms, nil
}
