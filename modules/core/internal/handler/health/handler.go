package health

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Handler exposes the liveness probe endpoint.
type Handler struct{}

// NewHandler constructs a Handler.
func NewHandler() *Handler { return &Handler{} }

// Live responds with {"status":"ok"} to signal the process is alive.
func (h *Handler) Live(c echo.Context) error {
	return c.JSON(http.StatusOK, HealthResponse{Status: "ok"})
}
