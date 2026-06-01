package auth

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	authConfig "github.com/dinhtp/lee-goo/modules/authentication/config"
	"github.com/dinhtp/lee-goo/modules/authentication/contracts"
	domainAuth "github.com/dinhtp/lee-goo/modules/authentication/internal/domain/auth"
	userContracts "github.com/dinhtp/lee-goo/modules/user/contracts"
	"github.com/dinhtp/lee-goo/system/eventbus"
	"github.com/dinhtp/lee-goo/system/security"
)

// service implements domainAuth.UseCase using stateless JWT tokens.
type service struct {
	userService userContracts.UserService
	signer      security.Signer
	verifier    security.Verifier
	eventBus    eventbus.EventBus
	config      *authConfig.AuthConfig
}

// compile-time interface check
var _ domainAuth.UseCase = (*service)(nil)

// NewService constructs the authentication service with all required dependencies.
func NewService(
	userSvc userContracts.UserService,
	signer security.Signer,
	verifier security.Verifier,
	bus eventbus.EventBus,
	cfg *authConfig.AuthConfig,
) domainAuth.UseCase {
	return &service{
		userService: userSvc,
		signer:      signer,
		verifier:    verifier,
		eventBus:    bus,
		config:      cfg,
	}
}

// Login verifies the provided credentials and issues an access+refresh token pair.
// It publishes login_failed on wrong password and login_succeeded on success.
func (s *service) Login(ctx context.Context, email, password string) (*domainAuth.TokenPair, error) {
	user, err := s.userService.FindByEmail(ctx, email)
	if err != nil {
		return nil, contracts.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		_ = s.eventBus.Publish(ctx, "authentication.login_failed", contracts.LoginFailedEvent{Email: email})
		return nil, contracts.ErrInvalidCredentials
	}

	accessToken, err := s.signer.Sign(jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"type":  "access",
		"exp":   time.Now().Add(s.config.AccessTokenTTL).Unix(),
	})
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.signer.Sign(jwt.MapClaims{
		"sub":  user.ID,
		"type": "refresh",
		"exp":  time.Now().Add(s.config.RefreshTokenTTL).Unix(),
	})
	if err != nil {
		return nil, err
	}

	_ = s.eventBus.Publish(ctx, "authentication.login_succeeded", contracts.LoginSucceededEvent{
		UserID: user.ID,
		Email:  user.Email,
	})

	return &domainAuth.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.config.AccessTokenTTL.Seconds()),
	}, nil
}

// Refresh validates a refresh token and issues a new access token.
// It rejects tokens of type "access" to prevent token-type confusion attacks.
func (s *service) Refresh(_ context.Context, refreshToken string) (*domainAuth.TokenPair, error) {
	claims, err := s.verifier.Verify(refreshToken)
	if err != nil {
		return nil, contracts.ErrTokenInvalid
	}

	mapClaims, ok := claims.(jwt.MapClaims)
	if !ok {
		return nil, contracts.ErrTokenInvalid
	}

	tokenType, _ := mapClaims["type"].(string)
	if tokenType != "refresh" {
		return nil, contracts.ErrTokenInvalid
	}

	userID, _ := mapClaims["sub"].(string)

	accessToken, err := s.signer.Sign(jwt.MapClaims{
		"sub":  userID,
		"type": "access",
		"exp":  time.Now().Add(s.config.AccessTokenTTL).Unix(),
	})
	if err != nil {
		return nil, err
	}

	return &domainAuth.TokenPair{
		AccessToken: accessToken,
		ExpiresIn:   int64(s.config.AccessTokenTTL.Seconds()),
	}, nil
}

// Logout is a no-op for stateless JWT — token invalidation is the client's responsibility.
func (s *service) Logout(_ context.Context) error {
	return nil
}

// ── AuthService adapter ──────────────────────────────────────────────────────

// authServiceAdapter bridges domainAuth.UseCase to the public contracts.AuthService interface.
type authServiceAdapter struct {
	uc domainAuth.UseCase
}

// compile-time interface check
var _ contracts.AuthService = (*authServiceAdapter)(nil)

// NewAuthServiceAdapter wraps a domainAuth.UseCase as a contracts.AuthService.
func NewAuthServiceAdapter(uc domainAuth.UseCase) contracts.AuthService {
	return &authServiceAdapter{uc: uc}
}

func (a *authServiceAdapter) Login(ctx context.Context, email, password string) (*contracts.TokenPair, error) {
	pair, err := a.uc.Login(ctx, email, password)
	if err != nil {
		return nil, err
	}
	return &contracts.TokenPair{
		AccessToken:  pair.AccessToken,
		RefreshToken: pair.RefreshToken,
		ExpiresIn:    pair.ExpiresIn,
	}, nil
}

func (a *authServiceAdapter) Refresh(ctx context.Context, refreshToken string) (*contracts.TokenPair, error) {
	pair, err := a.uc.Refresh(ctx, refreshToken)
	if err != nil {
		return nil, err
	}
	return &contracts.TokenPair{
		AccessToken: pair.AccessToken,
		ExpiresIn:   pair.ExpiresIn,
	}, nil
}
