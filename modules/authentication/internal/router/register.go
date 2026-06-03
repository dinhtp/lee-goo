package router

import "github.com/labstack/echo/v4"

// Register mounts the HandlerRouter under /auth on the platform HTTP server.
func Register(app *echo.Echo, r HandlerRouter) {
	g := app.Group("/auth")
	r.Register(g)
}
