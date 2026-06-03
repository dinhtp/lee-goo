package server

import (
	"context"

	"github.com/labstack/echo/v4"
)

// Engine is the HTTP server abstraction.
type Engine interface {
	Address() string
	Instance() (*echo.Echo, error)
	Startup() error
	Shutdown(ctx context.Context) error
}

type engine struct {
	address  string
	instance *echo.Echo
}

var _ Engine = &engine{}

// NewEngine creates an Echo server engine, applying each ConfigProvider in order.
func NewEngine(address string, configs ...ConfigProvider) Engine {
	echoServer := echo.New()
	for _, provide := range configs {
		provide(echoServer)
	}
	return &engine{address: address, instance: echoServer}
}

func (e *engine) Address() string { return e.address }

func (e *engine) Instance() (*echo.Echo, error) {
	if e.instance == nil {
		return nil, ErrUninitializedEngine
	}
	return e.instance, nil
}

func (e *engine) Startup() error {
	if e.Address() == "" {
		return ErrMissingServerAddress
	}
	return e.instance.Start(e.address)
}

func (e *engine) Shutdown(ctx context.Context) error {
	if e.instance == nil {
		return ErrUninitializedEngine
	}
	return e.instance.Shutdown(ctx)
}
