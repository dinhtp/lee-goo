package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"

	domainAuth "github.com/dinhtp/lee-goo/modules/authentication/internal/domain/auth"
	"github.com/dinhtp/lee-goo/pkg/validate"
)

// Handler exposes authentication endpoints via Echo.
type Handler struct {
	useCase domainAuth.UseCase
}

// NewHandler constructs an auth Handler with the provided domain use case.
func NewHandler(uc domainAuth.UseCase) *Handler {
	return &Handler{useCase: uc}
}

// Login handles POST /auth/login.
// Returns 401 on invalid credentials, 200 with token pair on success.
func (h *Handler) Login(c echo.Context) error {
	var req LoginRequest
	if err := validate.BindAndValidate(c, &req); err != nil {
		return err
	}

	pair, err := h.useCase.Login(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	return c.JSON(http.StatusOK, TokenResponse{
		AccessToken:  pair.AccessToken,
		RefreshToken: pair.RefreshToken,
		ExpiresIn:    pair.ExpiresIn,
	})
}

// Refresh handles POST /auth/refresh.
// Returns 401 on invalid or wrong-type token, 200 with new access token on success.
func (h *Handler) Refresh(c echo.Context) error {
	var req RefreshRequest
	if err := validate.BindAndValidate(c, &req); err != nil {
		return err
	}

	pair, err := h.useCase.Refresh(c.Request().Context(), req.RefreshToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	return c.JSON(http.StatusOK, TokenResponse{
		AccessToken: pair.AccessToken,
		ExpiresIn:   pair.ExpiresIn,
	})
}

// Logout handles POST /auth/logout.
// Always succeeds — stateless JWT logout is client-side token discard.
func (h *Handler) Logout(c echo.Context) error {
	_ = h.useCase.Logout(c.Request().Context())
	return c.JSON(http.StatusOK, map[string]string{"message": "logged out"})
}
