package module

import "context"

// UseCase defines all module lifecycle operations.
type UseCase interface {
	Discover(ctx context.Context) ([]Module, error)
	Install(ctx context.Context, name string) error
	Enable(ctx context.Context, name string) error
	Disable(ctx context.Context, name string) error
	Upgrade(ctx context.Context, name string) error
	Uninstall(ctx context.Context, name string, force bool) error
	Remove(ctx context.Context, name string) error
	Migrate(ctx context.Context, name string) error
	MigrateAll(ctx context.Context) error
	Status(ctx context.Context, name string) (*Module, error)
	Graph(ctx context.Context) (map[string][]string, error)
	Doctor(ctx context.Context) ([]string, error)
	Sync(ctx context.Context) error
}
