package http

import "go.uber.org/fx"

// FxOptions returns the fx module for the HTTP server.
func FxOptions() fx.Option {
	return fx.Options(
		fx.Provide(NewServer),
		fx.Invoke(RegisterLifecycle),
	)
}
