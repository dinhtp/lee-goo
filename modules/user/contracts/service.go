package contracts

import "context"

// UserService is the public cross-module interface for user operations.
// Authentication and authorization modules import only this package.
type UserService interface {
	FindByID(ctx context.Context, id string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	CreateUser(ctx context.Context, email, name, password string) (*User, error)
}

// User is the cross-module value object.
// PasswordHash is included so the authentication module can verify credentials
// without importing user internals. Only the auth module should read this field.
type User struct {
	ID           string
	Email        string
	Name         string
	PasswordHash string
}
