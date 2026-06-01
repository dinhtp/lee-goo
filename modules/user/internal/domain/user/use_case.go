package user

import "context"

// UseCase defines all application-level operations on the user domain.
// The service layer implements this; the handler layer depends on it.
type UseCase interface {
	CreateUser(ctx context.Context, email, name, password string) (*User, error)
	FindByID(ctx context.Context, id string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	ListUsers(ctx context.Context) ([]User, error)
	UpdateUser(ctx context.Context, id, name string) (*User, error)
}
