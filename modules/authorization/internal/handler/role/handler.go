package role

import (
	"net/http"

	"github.com/labstack/echo/v4"

	domainRole "github.com/dinhtp/lee-goo/modules/authorization/internal/domain/role"
	"github.com/dinhtp/lee-goo/pkg/validate"
)

// Handler handles HTTP requests for role management and policy evaluation.
type Handler struct {
	roleUseCase   domainRole.RoleUseCase
	policyUseCase domainRole.PolicyUseCase
}

// NewHandler constructs a Handler with the required use cases.
func NewHandler(roleUC domainRole.RoleUseCase, policyUC domainRole.PolicyUseCase) *Handler {
	return &Handler{roleUseCase: roleUC, policyUseCase: policyUC}
}

// CreateRole handles POST /roles — creates a new role.
func (h *Handler) CreateRole(c echo.Context) error {
	var req CreateRoleRequest
	if err := validate.BindAndValidate(c, &req); err != nil {
		return err
	}
	role, err := h.roleUseCase.CreateRole(c.Request().Context(), req.Name, req.Description)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, toRoleResponse(role))
}

// AssignRole handles POST /users/:id/roles — assigns a role to a user.
func (h *Handler) AssignRole(c echo.Context) error {
	userID := c.Param("id")
	var req AssignRoleRequest
	if err := validate.BindAndValidate(c, &req); err != nil {
		return err
	}
	if err := h.roleUseCase.AssignRoleToUser(c.Request().Context(), userID, req.RoleID); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "role assigned"})
}

// ListRoles handles GET /roles — returns all defined roles.
func (h *Handler) ListRoles(c echo.Context) error {
	roles, err := h.roleUseCase.ListRoles(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	resp := make([]RoleResponse, 0, len(roles))
	for _, r := range roles {
		resp = append(resp, toRoleResponse(&r))
	}
	return c.JSON(http.StatusOK, resp)
}

func toRoleResponse(r *domainRole.Role) RoleResponse {
	return RoleResponse{ID: r.ID, Name: r.Name, Description: r.Description}
}
