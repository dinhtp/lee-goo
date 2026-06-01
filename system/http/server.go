package http

import (
	"context"
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/fx"

	"github.com/dinhtp/lee-goo/system/config"
)

// Server wraps the Echo instance.
type Server struct {
	Echo *echo.Echo
}

// NewServer creates a pre-configured Echo server with recover and logger middleware.
func NewServer(logger *slog.Logger) *Server {
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	return &Server{Echo: e}
}

// RegisterLifecycle hooks the server start/stop into the fx lifecycle.
func RegisterLifecycle(lc fx.Lifecycle, s *Server, cfg *config.Config) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			port := cfg.Server.Port
			if port == "" {
				port = "8080"
			}
			// Non-blocking start; Echo logs any bind errors internally.
			go func() {
				if err := s.Echo.Start(":" + port); err != nil {
					logger := slog.Default()
					logger.Info("http server stopped", "reason", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return s.Echo.Shutdown(ctx)
		},
	})
}
