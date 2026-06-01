package role_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	domainRole "github.com/dinhtp/lee-goo/modules/authorization/internal/domain/role"
	roleService "github.com/dinhtp/lee-goo/modules/authorization/internal/service/role"
	"github.com/dinhtp/lee-goo/system/eventbus"
)

// --- mock RolePort ---

type mockRolePort struct {
	roles      map[string]domainRole.Role
	userRoles  map[string][]string // userID -> []roleID
	assignErr  error
	removeErr  error
	findByUser func(userID string) ([]domainRole.Role, error)
}

func newMockRolePort() *mockRolePort {
	return &mockRolePort{
		roles:     make(map[string]domainRole.Role),
		userRoles: make(map[string][]string),
	}
}

func (m *mockRolePort) Create(_ context.Context, r domainRole.Role) error {
	m.roles[r.ID] = r
	return nil
}

func (m *mockRolePort) FindByID(_ context.Context, id string) (*domainRole.Role, error) {
	r, ok := m.roles[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return &r, nil
}

func (m *mockRolePort) List(_ context.Context) ([]domainRole.Role, error) {
	roles := make([]domainRole.Role, 0, len(m.roles))
	for _, r := range m.roles {
		roles = append(roles, r)
	}
	return roles, nil
}

func (m *mockRolePort) AssignToUser(_ context.Context, userID, roleID string) error {
	if m.assignErr != nil {
		return m.assignErr
	}
	m.userRoles[userID] = append(m.userRoles[userID], roleID)
	return nil
}

func (m *mockRolePort) RemoveFromUser(_ context.Context, userID, roleID string) error {
	if m.removeErr != nil {
		return m.removeErr
	}
	existing := m.userRoles[userID]
	updated := existing[:0]
	for _, id := range existing {
		if id != roleID {
			updated = append(updated, id)
		}
	}
	m.userRoles[userID] = updated
	return nil
}

func (m *mockRolePort) FindByUserID(_ context.Context, userID string) ([]domainRole.Role, error) {
	if m.findByUser != nil {
		return m.findByUser(userID)
	}
	roleIDs := m.userRoles[userID]
	roles := make([]domainRole.Role, 0, len(roleIDs))
	for _, id := range roleIDs {
		if r, ok := m.roles[id]; ok {
			roles = append(roles, r)
		}
	}
	return roles, nil
}

// --- mock PermissionPort ---

type mockPermissionPort struct {
	perms map[string][]domainRole.Permission // roleID -> []Permission
}

func newMockPermissionPort() *mockPermissionPort {
	return &mockPermissionPort{perms: make(map[string][]domainRole.Permission)}
}

func (m *mockPermissionPort) Create(_ context.Context, _ domainRole.Permission) error {
	return nil
}

func (m *mockPermissionPort) AssignToRole(_ context.Context, roleID, _ string) error {
	return nil
}

func (m *mockPermissionPort) FindByRoleIDs(_ context.Context, roleIDs []string) ([]domainRole.Permission, error) {
	seen := make(map[string]bool)
	var result []domainRole.Permission
	for _, id := range roleIDs {
		for _, p := range m.perms[id] {
			key := p.Action + ":" + p.Resource
			if !seen[key] {
				seen[key] = true
				result = append(result, p)
			}
		}
	}
	return result, nil
}

// --- helpers ---

func buildService(rolePort *mockRolePort, permPort *mockPermissionPort) *roleService.Service {
	return roleService.NewService(rolePort, permPort, eventbus.NewNoopEventBus())
}

// --- tests ---

func TestCanDo_ReturnsTrueWhenPermissionMatches(t *testing.T) {
	rolePort := newMockRolePort()
	permPort := newMockPermissionPort()

	// Seed: role "r1" with permission "create:orders"
	rolePort.roles["r1"] = domainRole.Role{ID: "r1", Name: "editor"}
	rolePort.userRoles["user1"] = []string{"r1"}
	permPort.perms["r1"] = []domainRole.Permission{
		{ID: "p1", Action: "create", Resource: "orders"},
	}

	svc := buildService(rolePort, permPort)
	ok, err := svc.CanDo(context.Background(), "user1", "create", "orders")
	require.NoError(t, err)
	assert.True(t, ok)
}

func TestCanDo_ReturnsFalseWhenNoMatchingPermission(t *testing.T) {
	rolePort := newMockRolePort()
	permPort := newMockPermissionPort()

	rolePort.roles["r1"] = domainRole.Role{ID: "r1", Name: "viewer"}
	rolePort.userRoles["user1"] = []string{"r1"}
	permPort.perms["r1"] = []domainRole.Permission{
		{ID: "p1", Action: "read", Resource: "orders"},
	}

	svc := buildService(rolePort, permPort)
	ok, err := svc.CanDo(context.Background(), "user1", "delete", "orders")
	require.NoError(t, err)
	assert.False(t, ok)
}

func TestCanDo_ReturnsFalseWhenUserHasNoRoles(t *testing.T) {
	rolePort := newMockRolePort()
	permPort := newMockPermissionPort()

	svc := buildService(rolePort, permPort)
	ok, err := svc.CanDo(context.Background(), "unknown-user", "read", "users")
	require.NoError(t, err)
	assert.False(t, ok)
}

func TestAssignRoleToUser_InvalidatesCache(t *testing.T) {
	rolePort := newMockRolePort()
	permPort := newMockPermissionPort()

	rolePort.roles["r1"] = domainRole.Role{ID: "r1", Name: "admin"}
	permPort.perms["r1"] = []domainRole.Permission{
		{ID: "p1", Action: "delete", Resource: "users"},
	}

	svc := buildService(rolePort, permPort)
	ctx := context.Background()

	// First call: user has no roles — cache stores empty result.
	ok, err := svc.CanDo(ctx, "user2", "delete", "users")
	require.NoError(t, err)
	assert.False(t, ok, "user2 has no roles yet")

	// Assign role — must invalidate cache.
	rolePort.userRoles["user2"] = []string{"r1"}
	err = svc.AssignRoleToUser(ctx, "user2", "r1")
	require.NoError(t, err)

	// Second call: cache invalidated, fresh lookup should return true.
	ok, err = svc.CanDo(ctx, "user2", "delete", "users")
	require.NoError(t, err)
	assert.True(t, ok, "user2 should now have delete:users after cache invalidation")
}

func TestRemoveRoleFromUser_InvalidatesCache(t *testing.T) {
	rolePort := newMockRolePort()
	permPort := newMockPermissionPort()

	rolePort.roles["r1"] = domainRole.Role{ID: "r1", Name: "admin"}
	rolePort.userRoles["user3"] = []string{"r1"}
	permPort.perms["r1"] = []domainRole.Permission{
		{ID: "p1", Action: "delete", Resource: "users"},
	}

	svc := buildService(rolePort, permPort)
	ctx := context.Background()

	// Warm cache with permission.
	ok, err := svc.CanDo(ctx, "user3", "delete", "users")
	require.NoError(t, err)
	assert.True(t, ok)

	// Remove role — must invalidate cache.
	err = svc.RemoveRoleFromUser(ctx, "user3", "r1")
	require.NoError(t, err)

	// Now CanDo must re-query and return false.
	ok, err = svc.CanDo(ctx, "user3", "delete", "users")
	require.NoError(t, err)
	assert.False(t, ok, "cache must be invalidated after role removal")
}
