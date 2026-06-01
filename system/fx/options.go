// Package platformfx wires all system providers into a single fx.Option
// that application entrypoints and modules can compose.
package platformfx

import (
	"go.uber.org/fx"

	"github.com/dinhtp/lee-goo/system/config"
	"github.com/dinhtp/lee-goo/system/database"
	"github.com/dinhtp/lee-goo/system/eventbus"
	"github.com/dinhtp/lee-goo/system/extension"
	systemhttp "github.com/dinhtp/lee-goo/system/http"
	"github.com/dinhtp/lee-goo/system/logger"
	"github.com/dinhtp/lee-goo/system/security"
)

// Options returns the complete set of system fx providers.
func Options() fx.Option {
	return fx.Options(
		fx.Provide(config.Load),
		fx.Provide(logger.NewLogger),
		database.FxOptions(),
		systemhttp.FxOptions(),
		eventbus.FxOptions(),
		extension.FxOptions(),
		// security.NewJWTSecurity returns *jwtSecurity; NewSigner/NewVerifier
		// extract the Signer and Verifier interfaces from it.
		fx.Provide(security.NewJWTSecurity),
		fx.Provide(security.NewSigner),
		fx.Provide(security.NewVerifier),
	)
}

// TestOptions returns system options suitable for integration tests.
// Tests supply DATABASE_URL and AUTH_JWT_SECRET via environment variables.
func TestOptions() fx.Option {
	return Options()
}
