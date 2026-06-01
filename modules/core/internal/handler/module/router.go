package module

import "github.com/labstack/echo/v4"

// router wires Handler methods to an Echo group.
type router struct {
	handler *Handler
}

// NewRouter constructs a router for the module handler.
func NewRouter(h *Handler) *router {
	return &router{handler: h}
}

// Register attaches module routes to the given Echo group.
func (r *router) Register(g *echo.Group) {
	g.GET("", r.handler.List)
	g.GET("/:name", r.handler.Status)
}
