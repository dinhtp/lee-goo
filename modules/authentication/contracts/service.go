package contracts

import "context"

// TokenPair holds the access and optional refresh tokens returned after authentication.
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn    int64  `json:"expires_in"`
}

// AuthService is the public cross-module interface for authentication operations.
type AuthService interface {
	Login(ctx context.Context, email, password string) (*TokenPair, error)
	Refresh(ctx context.Context, refreshToken string) (*TokenPair, error)
}
