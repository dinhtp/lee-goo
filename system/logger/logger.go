package logger

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/dinhtp/lee-goo/system/config"
)

// NewLogger constructs a JSON slog.Logger at the level configured in cfg.Log.Level.
func NewLogger(cfg *config.Config) (*slog.Logger, error) {
	level, err := parseLevel(cfg.Log.Level)
	if err != nil {
		return nil, err
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	return slog.New(handler), nil
}

// parseLevel converts a string level name to slog.Level.
func parseLevel(s string) (slog.Level, error) {
	switch s {
	case "debug":
		return slog.LevelDebug, nil
	case "info", "":
		return slog.LevelInfo, nil
	case "warn", "warning":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return slog.LevelInfo, fmt.Errorf("unknown log level %q, defaulting to info", s)
	}
}
