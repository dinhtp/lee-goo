package logger

import (
	"strings"

	"github.com/dinhtp/lee-goo/system/config"
)

// NewLogger is the uber/fx constructor for *Logger.
// It normalises cfg.Log.Level to uppercase to match the Level constants (e.g. "info" → "INFO").
func NewLogger(cfg *config.Config) (*Logger, error) {
	lvl := Level(strings.ToUpper(cfg.Log.Level))
	return NewLog("lee-goo", WithLevel(lvl)), nil
}
