package extension

import "go.uber.org/fx"

// FxOptions returns the fx.Option that provides the ExtensionRegistry singleton.
func FxOptions() fx.Option {
	return fx.Provide(NewExtensionRegistry)
}
