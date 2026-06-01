package eventbus

import "go.uber.org/fx"

// FxOptions returns the fx.Option that provides the local in-process EventBus.
func FxOptions() fx.Option {
	return fx.Provide(NewLocalEventBus)
}
