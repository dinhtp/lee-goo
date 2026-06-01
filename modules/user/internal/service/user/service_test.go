package user_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dinhtp/lee-goo/modules/user/contracts"
	domainUser "github.com/dinhtp/lee-goo/modules/user/internal/domain/user"
	userService "github.com/dinhtp/lee-goo/modules/user/internal/service/user"
	"github.com/dinhtp/lee-goo/system/eventbus"
	"github.com/dinhtp/lee-goo/system/extension"
)

// --- manual mocks ---

type mockUserPort struct {
	users map[string]*domainUser.User
}

func newMockUserPort() *mockUserPort {
	return &mockUserPort{users: make(map[string]*domainUser.User)}
}

func (m *mockUserPort) FindByID(_ context.Context, id string) (*domainUser.User, error) {
	u, ok := m.users[id]
	if !ok {
		return nil, contracts.ErrUserNotFound
	}
	return u, nil
}

func (m *mockUserPort) FindByEmail(_ context.Context, email string) (*domainUser.User, error) {
	for _, u := range m.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, contracts.ErrUserNotFound
}

func (m *mockUserPort) List(_ context.Context) ([]domainUser.User, error) {
	out := make([]domainUser.User, 0, len(m.users))
	for _, u := range m.users {
		out = append(out, *u)
	}
	return out, nil
}

func (m *mockUserPort) Create(_ context.Context, u domainUser.User) error {
	for _, existing := range m.users {
		if existing.Email == u.Email {
			return contracts.ErrUserAlreadyExists
		}
	}
	m.users[u.ID] = &u
	return nil
}

func (m *mockUserPort) Update(_ context.Context, u domainUser.User) error {
	if _, ok := m.users[u.ID]; !ok {
		return contracts.ErrUserNotFound
	}
	m.users[u.ID] = &u
	return nil
}

// --- helpers ---

func newTestService(port domainUser.UserPort) domainUser.UseCase {
	return userService.NewService(port, eventbus.NewNoopEventBus(), extension.NewExtensionRegistry())
}

// --- tests ---

func TestCreateUser_Success(t *testing.T) {
	svc := newTestService(newMockUserPort())

	u, err := svc.CreateUser(context.Background(), "alice@example.com", "Alice", "secret123")
	require.NoError(t, err)
	assert.NotEmpty(t, u.ID)
	assert.Equal(t, "alice@example.com", u.Email)
	assert.Equal(t, "Alice", u.Name)
	assert.NotEmpty(t, u.PasswordHash)
	// Hash must never equal the plaintext password.
	assert.NotEqual(t, "secret123", u.PasswordHash)
}

func TestCreateUser_DuplicateEmail_ReturnsErrUserAlreadyExists(t *testing.T) {
	svc := newTestService(newMockUserPort())

	_, err := svc.CreateUser(context.Background(), "bob@example.com", "Bob", "pass1")
	require.NoError(t, err)

	_, err = svc.CreateUser(context.Background(), "bob@example.com", "Bob2", "pass2")
	assert.ErrorIs(t, err, contracts.ErrUserAlreadyExists)
}

func TestFindByEmail_NotFound_ReturnsErrUserNotFound(t *testing.T) {
	svc := newTestService(newMockUserPort())

	_, err := svc.FindByEmail(context.Background(), "nobody@example.com")
	assert.ErrorIs(t, err, contracts.ErrUserNotFound)
}

func TestFindByID_Success(t *testing.T) {
	svc := newTestService(newMockUserPort())

	created, err := svc.CreateUser(context.Background(), "carol@example.com", "Carol", "pw")
	require.NoError(t, err)

	found, err := svc.FindByID(context.Background(), created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, found.ID)
}

func TestUpdateUser_Success(t *testing.T) {
	svc := newTestService(newMockUserPort())

	created, err := svc.CreateUser(context.Background(), "dave@example.com", "Dave", "pw")
	require.NoError(t, err)

	updated, err := svc.UpdateUser(context.Background(), created.ID, "David")
	require.NoError(t, err)
	assert.Equal(t, "David", updated.Name)
}

func TestListUsers_ReturnsAllUsers(t *testing.T) {
	svc := newTestService(newMockUserPort())

	_, _ = svc.CreateUser(context.Background(), "u1@example.com", "U1", "pw")
	_, _ = svc.CreateUser(context.Background(), "u2@example.com", "U2", "pw")

	users, err := svc.ListUsers(context.Background())
	require.NoError(t, err)
	assert.Len(t, users, 2)
}
