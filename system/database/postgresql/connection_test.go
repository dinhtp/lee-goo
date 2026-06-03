package postgresql

import (
	"errors"
	"testing"
)

func TestNewDefaultConfigSetsSQLDefaults(t *testing.T) {
	config := newDefaultConfig()
	if config == nil {
		t.Fatal("expected config")
	}
	if config.MaxIdleConnections != 2 {
		t.Fatalf("expected 2 max idle connections, got %d", config.MaxIdleConnections)
	}
	if config.MaxOpenConnections != 4 {
		t.Fatalf("expected 4 max open connections, got %d", config.MaxOpenConnections)
	}
}

func TestNewConnectionInitializesWithNilConfig(t *testing.T) {
	conn := NewConnection("dsn", nil)
	raw, ok := conn.(*connection)
	if !ok {
		t.Fatal("expected concrete connection")
	}
	if raw.DataSourceName() != "dsn" {
		t.Fatalf("expected dsn, got %s", raw.DataSourceName())
	}
	if raw.config == nil {
		t.Fatal("expected initialized config")
	}
}

func TestOpenReturnsErrMissingConfigWhenConfigIsNil(t *testing.T) {
	conn := &connection{dsn: "dsn"}

	if _, err := conn.Open(); !errors.Is(err, ErrMissingConfig) {
		t.Fatalf("expected %v, got %v", ErrMissingConfig, err)
	}
}

func TestCloseAndPingReturnErrUninitializedDatabaseWhenUnopened(t *testing.T) {
	conn := NewConnection("dsn", nil)

	if err := conn.Close(); !errors.Is(err, ErrUninitializedDatabase) {
		t.Fatalf("expected close error %v, got %v", ErrUninitializedDatabase, err)
	}
	if err := conn.Ping(); !errors.Is(err, ErrUninitializedDatabase) {
		t.Fatalf("expected ping error %v, got %v", ErrUninitializedDatabase, err)
	}
}

func TestInstanceReturnErrUninitializedDatabaseWhenUnopened(t *testing.T) {
	conn := NewConnection("dsn", nil)

	if _, err := conn.Instance(); !errors.Is(err, ErrUninitializedDatabase) {
		t.Fatalf("expected %v, got %v", ErrUninitializedDatabase, err)
	}
}

// TestOpenReturnsErrorWithoutStoringInstanceOnPingFailure verifies that a failed Open
// leaves the connection uninitialized — a bad host/port guarantees ping failure.
func TestOpenReturnsErrorWithoutStoringInstanceOnPingFailure(t *testing.T) {
	conn := NewConnection("host=127.0.0.1 port=1 user=postgres password=postgres dbname=lee_goo sslmode=disable connect_timeout=1", nil)

	if _, err := conn.Open(); err == nil {
		t.Fatal("expected open error with unreachable host")
	}

	if _, err := conn.Instance(); !errors.Is(err, ErrUninitializedDatabase) {
		t.Fatalf("expected %v after failed Open, got %v", ErrUninitializedDatabase, err)
	}
}
