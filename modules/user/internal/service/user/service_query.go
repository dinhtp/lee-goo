package user

import (
	"context"
	"time"

	"github.com/dinhtp/lee-goo/modules/user/contracts"
	domainUser "github.com/dinhtp/lee-goo/modules/user/internal/domain/user"
)

// FindByEmail retrieves a user by email address.
func (s *service) FindByEmail(ctx context.Context, email string) (*domainUser.User, error) {
	return s.userPort.FindByEmail(ctx, email)
}

// ListUsers returns all users ordered by creation date descending.
func (s *service) ListUsers(ctx context.Context) ([]domainUser.User, error) {
	return s.userPort.List(ctx)
}

// UpdateUser changes a user's display name and publishes user.updated event.
func (s *service) UpdateUser(ctx context.Context, id, name string) (*domainUser.User, error) {
	u, err := s.userPort.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	u.Name = name
	u.UpdatedAt = time.Now()

	if err := s.userPort.Update(ctx, *u); err != nil {
		return nil, err
	}

	_ = s.eventBus.Publish(ctx, "user.updated", contracts.UserUpdatedEvent{
		UserID: u.ID,
		Email:  u.Email,
		Name:   u.Name,
	})

	return u, nil
}
