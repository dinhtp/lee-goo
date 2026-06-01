// Package testapp provides a lightweight fx application wrapper for integration tests.
// Tests supply database and auth config via environment variables:
//   - DATABASE_HOST, DATABASE_USER, DATABASE_PASSWORD, DATABASE_DBNAME
//   - AUTH_JWT_SECRET
package testapp

import (
	"context"
	"testing"

	"go.uber.org/fx"

	systemfx "github.com/dinhtp/lee-goo/system/fx"
)

// App wraps an fx application for integration testing.
type App struct {
	fxApp *fx.App
}

// Option configures a test App.
type Option func(*options)

type options struct {
	modules []fx.Option
}

// WithModules adds fx module options to the test app.
func WithModules(mods ...fx.Option) Option {
	return func(o *options) {
		o.modules = append(o.modules, mods...)
	}
}

// New constructs a test App with the given module options.
// Requires DATABASE_HOST (or TEST_DATABASE_URL) and AUTH_JWT_SECRET env vars.
func New(opts ...Option) *App {
	o := &options{}
	for _, opt := range opts {
		opt(o)
	}

	fxOpts := []fx.Option{
		fx.NopLogger,
		systemfx.TestOptions(),
	}
	fxOpts = append(fxOpts, o.modules...)

	return &App{
		fxApp: fx.New(fxOpts...),
	}
}

// Start starts the fx app and registers cleanup via t.Cleanup.
func (a *App) Start(t *testing.T) {
	t.Helper()
	err := a.fxApp.Start(context.Background())
	if err != nil {
		t.Fatalf("testapp Start: %v", err)
	}
	t.Cleanup(func() {
		_ = a.fxApp.Stop(context.Background())
	})
}
