package logger

import (
	"fmt"
	"strings"

	"go.uber.org/zap/zapcore"
)

type Level string

const (
	LevelDebug Level = "DEBUG"
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelError Level = "ERROR"
	LevelPanic Level = "PANIC"
	LevelFatal Level = "FATAL"
)

// ParseLevel normalises s to uppercase and validates it against the known Level constants.
// Returns an error if s does not match any known level.
func ParseLevel(s string) (Level, error) {
	lvl := Level(strings.ToUpper(s))
	switch lvl {
	case LevelDebug, LevelInfo, LevelWarn, LevelError, LevelPanic, LevelFatal:
		return lvl, nil
	default:
		return LevelInfo, fmt.Errorf("unknown log level %q", s)
	}
}

func (l Level) ToZapLevel() zapcore.Level {
	switch l {
	case LevelDebug:
		return zapcore.DebugLevel
	case LevelInfo:
		return zapcore.InfoLevel
	case LevelWarn:
		return zapcore.WarnLevel
	case LevelError:
		return zapcore.ErrorLevel
	case LevelPanic:
		return zapcore.PanicLevel
	case LevelFatal:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}
