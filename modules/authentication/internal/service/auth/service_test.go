package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	authConfig "github.com/dinhtp/lee-goo/modules/authentication/config"
	"github.com/dinhtp/lee-goo/modules/authentication/contracts"
	userContracts "github.com/dinhtp/lee-goo/modules/user/contracts"
	"github.com/dinhtp/lee-goo/system/eventbus"
)

// ── manual mocks ─────────────────────────────────────────────────────────────

type mockSigner struct{}

func (m *mockSigner) Sign(_ jwt.MapClaims) (string, error) { return "signed-token", nil }

type mockVerifier struct {
	shouldFail bool
	tokenType  string
}

func (m *mockVerifier) Verify(_ string) (jwt.Claims, error) {
	if m.shouldFail {
		return nil, errors.New("invalid token")
	}
	return jwt.MapClaims{
		"sub":  "user-1",
		"type": m.tokenType,
		"exp":  float64(time.Now().Add(time.Hour).Unix()),
	}, nil
}

type mockUserService struct {
	user *userContracts.User
	err  error
}

func (m *mockUserService) FindByEmail(_ context.Context, _ string) (*userContracts.User, error) {
	return m.user, m.err
}

func (m *mockUserService) FindByID(_ context.Context, _ string) (*userContracts.User, error) {
	return m.user, m.err
}

func (m *mockUserService) CreateUser(_ context.Context, _, _, _ string) (*userContracts.User, error) {
	return m.user, m.err
}

type mockEventBus struct{}

func (m *mockEventBus) Publish(_ context.Context, _ string, _ any) error { return nil }
func (m *mockEventBus) Subscribe(_ string, _ eventbus.Handler)           {}

// ── helpers ───────────────────────────────────────────────────────────────────

// hashedUser creates a User whose PasswordHash is the bcrypt hash of rawPassword.
// Uses bcrypt.MinCost to keep tests fast.
func hashedUser(t *testing.T, rawPassword string) *userContracts.User {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.MinCost)
	require.NoError(t, err)
	return &userContracts.User{
		ID:           "user-1",
		Email:        "user@example.com",
		Name:         "Test User",
		PasswordHash: string(hash),
	}
}

func defaultConfig() *authConfig.AuthConfig {
	return &authConfig.AuthConfig{
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 7 * 24 * time.Hour,
	}
}

// ── tests ─────────────────────────────────────────────────────────────────────

func TestLogin_Success(t *testing.T) {
	const rawPassword = "correct-password"
	user := hashedUser(t, rawPassword)

	svc := NewService(
		&mockUserService{user: user},
		&mockSigner{},
		&mockVerifier{tokenType: "access"},
		&mockEventBus{},
		defaultConfig(),
	)

	pair, err := svc.Login(context.Background(), user.Email, rawPassword)

	require.NoError(t, err)
	assert.NotEmpty(t, pair.AccessToken)
	assert.NotEmpty(t, pair.RefreshToken)
	assert.Greater(t, pair.ExpiresIn, int64(0))
}

func TestLogin_WrongPassword_ReturnsErrInvalidCredentials(t *testing.T) {
	user := hashedUser(t, "correct-password")

	svc := NewService(
		&mockUserService{user: user},
		&mockSigner{},
		&mockVerifier{tokenType: "access"},
		&mockEventBus{},
		defaultConfig(),
	)

	_, err := svc.Login(context.Background(), user.Email, "wrong-password")

	assert.ErrorIs(t, err, contracts.ErrInvalidCredentials)
}

func TestLogin_UserNotFound_ReturnsErrInvalidCredentials(t *testing.T) {
	svc := NewService(
		&mockUserService{err: errors.New("not found")},
		&mockSigner{},
		&mockVerifier{tokenType: "access"},
		&mockEventBus{},
		defaultConfig(),
	)

	_, err := svc.Login(context.Background(), "ghost@example.com", "any-password")

	assert.ErrorIs(t, err, contracts.ErrInvalidCredentials)
}

func TestRefresh_ValidRefreshToken_ReturnsNewAccessToken(t *testing.T) {
	svc := NewService(
		&mockUserService{},
		&mockSigner{},
		&mockVerifier{tokenType: "refresh"},
		&mockEventBus{},
		defaultConfig(),
	)

	pair, err := svc.Refresh(context.Background(), "valid-refresh-token")

	require.NoError(t, err)
	assert.NotEmpty(t, pair.AccessToken)
	assert.Empty(t, pair.RefreshToken, "refresh should not be re-issued")
	assert.Greater(t, pair.ExpiresIn, int64(0))
}

func TestRefresh_WithAccessToken_ReturnsErrTokenInvalid(t *testing.T) {
	// Passing an access token (type:"access") must be rejected.
	svc := NewService(
		&mockUserService{},
		&mockSigner{},
		&mockVerifier{tokenType: "access"}, // verifier accepts but type is wrong
		&mockEventBus{},
		defaultConfig(),
	)

	_, err := svc.Refresh(context.Background(), "some-access-token")

	assert.ErrorIs(t, err, contracts.ErrTokenInvalid)
}

func TestRefresh_InvalidToken_ReturnsErrTokenInvalid(t *testing.T) {
	svc := NewService(
		&mockUserService{},
		&mockSigner{},
		&mockVerifier{shouldFail: true},
		&mockEventBus{},
		defaultConfig(),
	)

	_, err := svc.Refresh(context.Background(), "malformed-token")

	assert.ErrorIs(t, err, contracts.ErrTokenInvalid)
}

func TestLogout_AlwaysSucceeds(t *testing.T) {
	svc := NewService(
		&mockUserService{},
		&mockSigner{},
		&mockVerifier{},
		&mockEventBus{},
		defaultConfig(),
	)

	err := svc.Logout(context.Background())
	assert.NoError(t, err)
}
