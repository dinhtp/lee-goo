package router

import "github.com/labstack/echo/v4"

// HandlerRouter is implemented by any handler that can register routes on an Echo group.
type HandlerRouter interface {
	Register(g *echo.Group)
}
