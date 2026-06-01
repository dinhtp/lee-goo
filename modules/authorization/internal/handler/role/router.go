package role

import (
	"github.com/labstack/echo/v4"

	internalRouter "github.com/dinhtp/lee-goo/modules/authorization/internal/router"
)

type router struct {
	handler *Handler
}

// compile-time interface check
var _ internalRouter.HandlerRouter = (*router)(nil)

// NewRouter wraps the Handler and implements HandlerRouter for /roles.
func NewRouter(h *Handler) internalRouter.HandlerRouter {
	return &router{handler: h}
}

// Register mounts role CRUD routes onto the provided echo.Group.
func (r *router) Register(g *echo.Group) {
	g.POST("", r.handler.CreateRole)
	g.GET("", r.handler.ListRoles)
}
