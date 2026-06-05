package logger

import (
	"fmt"

	"github.com/dinhtp/lee-goo/system/config"
)

// NewLogger is the uber/fx constructor for *Logger.
// It reads the log level from cfg.Server.LogLevel (env: SERVER_LOG_LEVEL) and validates it against the Level constants.
func NewLogger(cfg *config.Config) (*Logger, error) {
	lvl, err := ParseLevel(cfg.Server.LogLevel)
	if err != nil {
		return nil, fmt.Errorf("config log level: %w", err)
	}
	return NewLog("lee-goo", WithLevel(lvl)), nil
}
