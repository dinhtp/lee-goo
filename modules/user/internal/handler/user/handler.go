package user

import (
	"net/http"

	"github.com/labstack/echo/v4"

	domainUser "github.com/dinhtp/lee-goo/modules/user/internal/domain/user"
	"github.com/dinhtp/lee-goo/pkg/validate"
)

// Handler holds HTTP handler methods for the user resource.
type Handler struct {
	useCase domainUser.UseCase
}

// NewHandler constructs a Handler with the given UseCase.
func NewHandler(useCase domainUser.UseCase) *Handler {
	return &Handler{useCase: useCase}
}

// CreateUser handles POST /users.
func (h *Handler) CreateUser(c echo.Context) error {
	var req CreateUserRequest
	if err := validate.BindAndValidate(c, &req); err != nil {
		return err
	}
	user, err := h.useCase.CreateUser(c.Request().Context(), req.Email, req.Name, req.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, toUserResponse(user))
}

// GetUser handles GET /users/:id.
func (h *Handler) GetUser(c echo.Context) error {
	id := c.Param("id")
	user, err := h.useCase.FindByID(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, toUserResponse(user))
}

// ListUsers handles GET /users.
func (h *Handler) ListUsers(c echo.Context) error {
	users, err := h.useCase.ListUsers(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	resp := make([]UserResponse, 0, len(users))
	for i := range users {
		resp = append(resp, toUserResponse(&users[i]))
	}
	return c.JSON(http.StatusOK, resp)
}

// UpdateUser handles PUT /users/:id.
func (h *Handler) UpdateUser(c echo.Context) error {
	id := c.Param("id")
	var req UpdateUserRequest
	if err := validate.BindAndValidate(c, &req); err != nil {
		return err
	}
	user, err := h.useCase.UpdateUser(c.Request().Context(), id, req.Name)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, toUserResponse(user))
}

func toUserResponse(u *domainUser.User) UserResponse {
	return UserResponse{ID: u.ID, Email: u.Email, Name: u.Name}
}
