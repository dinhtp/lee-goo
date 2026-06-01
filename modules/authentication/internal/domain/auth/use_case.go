package auth

import "context"

// UseCase defines the authentication domain operations.
type UseCase interface {
	Login(ctx context.Context, email, password string) (*TokenPair, error)
	Refresh(ctx context.Context, refreshToken string) (*TokenPair, error)
	Logout(ctx context.Context) error
}
