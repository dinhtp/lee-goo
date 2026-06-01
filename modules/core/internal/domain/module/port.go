package module

import "context"

// ModulePort is the repository interface for persisting module state.
type ModulePort interface {
	FindAll(ctx context.Context) ([]Module, error)
	FindByName(ctx context.Context, name string) (*Module, error)
	Upsert(ctx context.Context, m Module) error
	UpdateStatus(ctx context.Context, name string, status Status) error
}
