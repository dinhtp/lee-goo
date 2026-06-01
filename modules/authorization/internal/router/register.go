package router

import (
	systemHTTP "github.com/dinhtp/lee-goo/system/http"
)

// Register mounts the HandlerRouter under /roles on the platform HTTP server.
func Register(s *systemHTTP.Server, roleRouter HandlerRouter) {
	g := s.Echo.Group("/roles")
	roleRouter.Register(g)
}
