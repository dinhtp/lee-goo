package server

import (
	"errors"

	"github.com/labstack/echo/v4"
)

var (
	ErrMissingServerAddress = errors.New("server address is missing")
	ErrUninitializedEngine  = errors.New("server engine is not initialized")
)

// ConfigProvider configures the Echo instance during engine construction.
type ConfigProvider func(*echo.Echo)
