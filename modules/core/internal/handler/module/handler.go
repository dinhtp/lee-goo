package module

import (
	"net/http"

	"github.com/labstack/echo/v4"

	domainModule "github.com/dinhtp/lee-goo/modules/core/internal/domain/module"
)

// Handler exposes module management endpoints for the admin API.
type Handler struct {
	useCase domainModule.UseCase
}

// NewHandler constructs a Handler wired to the given use case.
func NewHandler(uc domainModule.UseCase) *Handler {
	return &Handler{useCase: uc}
}

// List returns all modules tracked in the database.
func (h *Handler) List(c echo.Context) error {
	modules, err := h.useCase.Discover(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	resp := make([]ModuleResponse, 0, len(modules))
	for _, m := range modules {
		resp = append(resp, ModuleResponse{
			Name:    m.Name,
			Version: m.Version,
			Status:  string(m.Status),
			Path:    m.Path,
		})
	}
	return c.JSON(http.StatusOK, resp)
}

// Status returns the current state of a single module by name.
func (h *Handler) Status(c echo.Context) error {
	name := c.Param("name")
	m, err := h.useCase.Status(c.Request().Context(), name)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if m == nil {
		return echo.NewHTTPError(http.StatusNotFound, "module not found")
	}
	return c.JSON(http.StatusOK, ModuleResponse{
		Name:    m.Name,
		Version: m.Version,
		Status:  string(m.Status),
		Path:    m.Path,
	})
}
