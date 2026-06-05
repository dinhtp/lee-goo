package health

import "github.com/labstack/echo/v4"

// Register mounts the liveness probe at GET /healthz on the root Echo instance.
// Registered directly on *echo.Echo (not a group) so no middleware applies.
func Register(app *echo.Echo, h *Handler) {
	app.GET("/healthz", h.Live)
}
