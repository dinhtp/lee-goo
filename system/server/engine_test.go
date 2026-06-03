package server

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	echohttp "github.com/labstack/echo/v4"
)

func TestNewEngineAppliesConfigProvidersInOrder(t *testing.T) {
	order := make([]string, 0, 2)

	eng := NewEngine(":8080",
		func(app *echohttp.Echo) {
			order = append(order, "first")
			app.GET("/health", func(c echohttp.Context) error {
				return c.String(http.StatusOK, "ok")
			})
		},
		func(app *echohttp.Echo) {
			order = append(order, "second")
		},
	)

	if eng.Address() != ":8080" {
		t.Fatalf("Address() = %q, want %q", eng.Address(), ":8080")
	}

	if !reflect.DeepEqual(order, []string{"first", "second"}) {
		t.Fatalf("config providers ran in order %v", order)
	}

	app, err := eng.Instance()
	if err != nil {
		t.Fatalf("Instance() error = %v", err)
	}

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	app.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("ServeHTTP() status = %d, want %d", recorder.Code, http.StatusOK)
	}

	if recorder.Body.String() != "ok" {
		t.Fatalf("ServeHTTP() body = %q, want %q", recorder.Body.String(), "ok")
	}
}

func TestInstanceReturnsErrorWhenEngineIsUninitialized(t *testing.T) {
	eng := &engine{}

	app, err := eng.Instance()
	if !errors.Is(err, ErrUninitializedEngine) {
		t.Fatalf("Instance() error = %v, want %v", err, ErrUninitializedEngine)
	}

	if app != nil {
		t.Fatalf("Instance() app = %v, want nil", app)
	}
}

func TestStartupReturnsErrorWhenAddressIsMissing(t *testing.T) {
	eng := NewEngine("")

	err := eng.Startup()
	if !errors.Is(err, ErrMissingServerAddress) {
		t.Fatalf("Startup() error = %v, want %v", err, ErrMissingServerAddress)
	}
}

func TestShutdownReturnsErrorWhenEngineIsUninitialized(t *testing.T) {
	eng := &engine{}

	err := eng.Shutdown(context.Background())
	if !errors.Is(err, ErrUninitializedEngine) {
		t.Fatalf("Shutdown() error = %v, want %v", err, ErrUninitializedEngine)
	}
}

func TestStartupAndShutdownServeRequests(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Listen() error = %v", err)
	}

	eng := NewEngine(listener.Addr().String(), func(app *echohttp.Echo) {
		app.GET("/health", func(c echohttp.Context) error {
			return c.String(http.StatusOK, "ok")
		})
	})

	impl := eng.(*engine)
	impl.instance.Listener = listener
	impl.instance.HideBanner = true
	impl.instance.HidePort = true

	serveErr := make(chan error, 1)
	go func() {
		serveErr <- impl.Startup()
	}()

	response, err := getEventually("http://" + listener.Addr().String() + "/health")
	if err != nil {
		t.Fatalf("GET /health error = %v", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("ReadAll() error = %v", err)
	}

	if response.StatusCode != http.StatusOK {
		t.Fatalf("GET /health status = %d, want %d", response.StatusCode, http.StatusOK)
	}

	if string(body) != "ok" {
		t.Fatalf("GET /health body = %q, want %q", string(body), "ok")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := impl.Shutdown(ctx); err != nil {
		t.Fatalf("Shutdown() error = %v", err)
	}

	if err := <-serveErr; !errors.Is(err, http.ErrServerClosed) {
		t.Fatalf("Startup() error = %v, want %v", err, http.ErrServerClosed)
	}
}

func getEventually(url string) (*http.Response, error) {
	client := &http.Client{Timeout: time.Second}
	deadline := time.Now().Add(2 * time.Second)

	for time.Now().Before(deadline) {
		response, err := client.Get(url)
		if err == nil {
			return response, nil
		}

		time.Sleep(10 * time.Millisecond)
	}

	return nil, errors.New("server did not start before timeout")
}
