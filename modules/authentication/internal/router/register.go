package router

import (
	systemHTTP "github.com/dinhtp/lee-goo/system/http"
)

// Register mounts the HandlerRouter under /auth on the platform HTTP server.
func Register(s *systemHTTP.Server, r HandlerRouter) {
	g := s.Echo.Group("/auth")
	r.Register(g)
}
