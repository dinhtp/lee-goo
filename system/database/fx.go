package database

import "go.uber.org/fx"

// FxOptions returns the fx module for database connection.
func FxOptions() fx.Option {
	return fx.Provide(NewConnection)
}
