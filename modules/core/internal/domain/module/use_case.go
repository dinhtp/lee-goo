package module

import "context"

// UseCase defines module lifecycle operations backed by CLI commands or HTTP handlers.
type UseCase interface {
	Discover(ctx context.Context) ([]Module, error)
	Install(ctx context.Context, name string) error
	Uninstall(ctx context.Context, name string, force bool) error
	Status(ctx context.Context, name string) (*Module, error)
}
