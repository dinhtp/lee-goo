package user

import (
	"github.com/labstack/echo/v4"

	internalRouter "github.com/dinhtp/lee-goo/modules/user/internal/router"
)

type router struct {
	handler *Handler
}

// compile-time interface check
var _ internalRouter.HandlerRouter = (*router)(nil)

// NewRouter wraps the Handler and returns a HandlerRouter ready for registration.
func NewRouter(h *Handler) internalRouter.HandlerRouter {
	return &router{handler: h}
}

// Register mounts all user routes on the given echo Group.
func (r *router) Register(g *echo.Group) {
	g.POST("", r.handler.CreateUser)
	g.GET("", r.handler.ListUsers)
	g.GET("/:id", r.handler.GetUser)
	g.PUT("/:id", r.handler.UpdateUser)
}
