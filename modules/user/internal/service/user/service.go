package user

import (
	"context"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/dinhtp/lee-goo/modules/user/contracts"
	domainUser "github.com/dinhtp/lee-goo/modules/user/internal/domain/user"
	"github.com/dinhtp/lee-goo/system/eventbus"
	"github.com/dinhtp/lee-goo/system/extension"
)

type service struct {
	userPort domainUser.UserPort
	eventBus eventbus.EventBus
	registry *extension.ExtensionRegistry
}

// compile-time interface check — service must satisfy the internal UseCase.
var _ domainUser.UseCase = (*service)(nil)

// NewService constructs the user service with all required dependencies.
func NewService(
	port domainUser.UserPort,
	bus eventbus.EventBus,
	registry *extension.ExtensionRegistry,
) domainUser.UseCase {
	return &service{userPort: port, eventBus: bus, registry: registry}
}

// CreateUser validates uniqueness, hashes the password, persists the user,
// fires extension hooks, and publishes the user.created event.
func (s *service) CreateUser(ctx context.Context, email, name, password string) (*domainUser.User, error) {
	// 1. Guard: reject duplicate email.
	_, err := s.userPort.FindByEmail(ctx, email)
	if err == nil {
		return nil, contracts.ErrUserAlreadyExists
	}
	if err != contracts.ErrUserNotFound {
		return nil, err
	}

	// 2. Hash password.
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 3. Build and persist domain entity.
	now := time.Now()
	u := domainUser.User{
		ID:           uuid.New().String(),
		Email:        email,
		Name:         name,
		PasswordHash: string(hash),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := s.userPort.Create(ctx, u); err != nil {
		return nil, err
	}

	// 4. Call "user.after_created" extension hooks (non-fatal; errors are ignored
	//    intentionally — hooks are observers, not transactional participants).
	hooks := s.registry.Resolve("user.after_created")
	for _, h := range hooks {
		if fn, ok := h.(func(context.Context, string) error); ok {
			_ = fn(ctx, u.ID)
		}
	}

	// 5. Publish domain event.
	_ = s.eventBus.Publish(ctx, "user.created", contracts.UserCreatedEvent{
		UserID: u.ID,
		Email:  u.Email,
	})

	return &u, nil
}

// FindByID retrieves a user by primary key.
func (s *service) FindByID(ctx context.Context, id string) (*domainUser.User, error) {
	return s.userPort.FindByID(ctx, id)
}
