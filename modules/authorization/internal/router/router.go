package router

import "github.com/labstack/echo/v4"

// HandlerRouter is implemented by each handler package to register its routes
// onto a given echo.Group.
type HandlerRouter interface {
	Register(g *echo.Group)
}
