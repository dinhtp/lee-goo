package server

import (
	"context"
	"fmt"

	echolib "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/fx"

	"github.com/dinhtp/lee-goo/system/config"
	"github.com/dinhtp/lee-goo/system/logger"
)

// NewEchoEngine builds an Engine from config with default middleware applied.
// Returns both Engine (for lifecycle) and *echo.Echo (for route registration via fx).
func NewEchoEngine(cfg *config.Config) (Engine, *echolib.Echo, error) {
	port := cfg.Server.Port
	if port == "" {
		port = "8080"
	}

	e := NewEngine(fmt.Sprintf(":%s", port),
		func(app *echolib.Echo) {
			app.HideBanner = true
			app.Use(middleware.Recover())
			app.Use(middleware.Logger())
		},
	)

	app, err := e.Instance()
	if err != nil {
		return nil, nil, err
	}

	return e, app, nil
}

// RegisterLifecycle hooks the engine start/stop into the fx lifecycle.
func RegisterLifecycle(lc fx.Lifecycle, e Engine, log *logger.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := e.Startup(); err != nil {
					log.Sugar().Infow("http server stopped", "reason", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return e.Shutdown(ctx)
		},
	})
}

// FxOptions returns the fx module for the HTTP server.
func FxOptions() fx.Option {
	return fx.Options(
		fx.Provide(NewEchoEngine),
		fx.Invoke(RegisterLifecycle),
	)
}
