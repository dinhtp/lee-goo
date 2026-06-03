package router

import "github.com/labstack/echo/v4"

// Register mounts the HandlerRouter under /roles on the platform HTTP server.
func Register(app *echo.Echo, roleRouter HandlerRouter) {
	g := app.Group("/roles")
	roleRouter.Register(g)
}
