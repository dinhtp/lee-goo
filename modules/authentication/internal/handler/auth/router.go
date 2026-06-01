package auth

import (
	"github.com/labstack/echo/v4"

	internalRouter "github.com/dinhtp/lee-goo/modules/authentication/internal/router"
)

// router wires the Handler into the HandlerRouter contract.
type router struct {
	handler *Handler
}

// compile-time interface check
var _ internalRouter.HandlerRouter = (*router)(nil)

// NewRouter constructs a HandlerRouter backed by the provided Handler.
func NewRouter(h *Handler) internalRouter.HandlerRouter {
	return &router{handler: h}
}

// Register mounts all authentication routes onto the provided Echo group.
func (r *router) Register(g *echo.Group) {
	g.POST("/login", r.handler.Login)
	g.POST("/refresh", r.handler.Refresh)
	g.POST("/logout", r.handler.Logout)
}
