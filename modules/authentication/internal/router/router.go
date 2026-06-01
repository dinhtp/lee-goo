package router

import "github.com/labstack/echo/v4"

// HandlerRouter is the capability interface every handler must implement
// to register its routes onto an Echo group.
type HandlerRouter interface {
	Register(g *echo.Group)
}
