package user

import "context"

// UserPort is the persistence abstraction for the user domain.
// The repository layer implements this; the service layer depends on it.
type UserPort interface {
	FindByID(ctx context.Context, id string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	List(ctx context.Context) ([]User, error)
	Create(ctx context.Context, u User) error
	Update(ctx context.Context, u User) error
}
