package router

import (
	systemHTTP "github.com/dinhtp/lee-goo/system/http"
)

// Register mounts the given HandlerRouter under /admin/modules on the platform HTTP server.
func Register(s *systemHTTP.Server, r HandlerRouter) {
	g := s.Echo.Group("/admin/modules")
	r.Register(g)
}
