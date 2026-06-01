package router

import (
	systemHTTP "github.com/dinhtp/lee-goo/system/http"
)

// Register mounts the HandlerRouter under /users on the platform HTTP server.
func Register(s *systemHTTP.Server, r HandlerRouter) {
	g := s.Echo.Group("/users")
	r.Register(g)
}
