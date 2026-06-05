package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"

	"github.com/dinhtp/lee-goo/modules/core/internal/handler/health"
)

// TestHealthzEndpoint verifies GET /healthz returns 200 {"status":"ok"} with no auth.
func TestHealthzEndpoint(t *testing.T) {
	e := echo.New()
	h := health.NewHandler()
	health.Register(e, h)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), `"status":"ok"`)
}
