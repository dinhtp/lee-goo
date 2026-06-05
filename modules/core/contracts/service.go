package contracts

import "context"

// ModuleService is the public-facing contract exposed to other modules.
type ModuleService interface {
	Install(ctx context.Context, name string) error
	Enable(ctx context.Context, name string) error
	Disable(ctx context.Context, name string) error
	Uninstall(ctx context.Context, name string, force bool) error
	Status(ctx context.Context, name string) (*ModuleInfo, error)
}

// ModuleInfo is the public read-only view of a module.
type ModuleInfo struct {
	Name    string
	Version string
	Status  string
	Path    string
}
