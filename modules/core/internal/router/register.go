package router

import "github.com/labstack/echo/v4"

// Register mounts the given HandlerRouter under /admin/modules on the platform HTTP server.
func Register(app *echo.Echo, r HandlerRouter) {
	g := app.Group("/admin/modules")
	r.Register(g)
}
