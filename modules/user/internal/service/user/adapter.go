package user

import (
	"context"

	"github.com/dinhtp/lee-goo/modules/user/contracts"
	domainUser "github.com/dinhtp/lee-goo/modules/user/internal/domain/user"
)

// UserServiceAdapter wraps the internal UseCase and exposes contracts.UserService
// so other modules (auth, authz) can depend on the public interface only.
type UserServiceAdapter struct {
	uc domainUser.UseCase
}

// compile-time interface check
var _ contracts.UserService = (*UserServiceAdapter)(nil)

// NewUserServiceAdapter returns a contracts.UserService backed by the internal UseCase.
func NewUserServiceAdapter(uc domainUser.UseCase) contracts.UserService {
	return &UserServiceAdapter{uc: uc}
}

func (a *UserServiceAdapter) FindByID(ctx context.Context, id string) (*contracts.User, error) {
	u, err := a.uc.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toContractsUser(u), nil
}

func (a *UserServiceAdapter) FindByEmail(ctx context.Context, email string) (*contracts.User, error) {
	u, err := a.uc.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return toContractsUser(u), nil
}

func (a *UserServiceAdapter) CreateUser(ctx context.Context, email, name, password string) (*contracts.User, error) {
	u, err := a.uc.CreateUser(ctx, email, name, password)
	if err != nil {
		return nil, err
	}
	return toContractsUser(u), nil
}

// toContractsUser maps the internal domain User to the public cross-module value object.
func toContractsUser(u *domainUser.User) *contracts.User {
	return &contracts.User{
		ID:           u.ID,
		Email:        u.Email,
		Name:         u.Name,
		PasswordHash: u.PasswordHash,
	}
}
