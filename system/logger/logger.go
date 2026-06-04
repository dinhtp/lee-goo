package logger

import (
	"context"
	"fmt"
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	nanoId "github.com/matoous/go-nanoid/v2"
)

// Logger wraps zap.Logger with level, log ID, and development-mode support.
type Logger struct {
	*zap.Logger
	level  Level
	prefix string
	logID  string
}

// NewLog creates a named zap logger with the given options.
func NewLog(name string, options ...Option) *Logger {
	result := &Logger{level: LevelInfo, logID: nanoId.Must()}

	for _, opt := range options {
		opt(result)
	}

	result.Logger = result.newZapLogger(name)

	return result
}

func (l *Logger) newZapLogger(name string) *zap.Logger {
	zapConfig := zap.NewProductionConfig()
	zapConfig.Level = zap.NewAtomicLevelAt(l.level.ToZapLevel())
	zapConfig.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	zapConfig.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder

	if l.level == LevelDebug {
		zapConfig.Development = true
	}

	zapLogger, err := zapConfig.Build()
	if err != nil {
		log.Println(err)
		return nil
	}

	defer func() {
		_ = zapLogger.Sync()
	}()

	// Skip 1 caller level so log lines point to the application code, not this wrapper.
	zapLogger = zapLogger.WithOptions(zap.AddCallerSkip(1))
	zapLogger = zapLogger.Named(name)

	if l.logID != "" {
		zapLogger = zapLogger.With(zap.String(FieldLogID, l.logID))
	}

	return zapLogger
}

func (l *Logger) clone() *Logger {
	cloned := *l
	return &cloned
}

// WithCtx returns a logger with the log ID extracted from ctx (if present).
func (l *Logger) WithCtx(ctx context.Context) *Logger {
	logId := ctx.Value(FieldLogID)
	if logId == nil {
		return l
	}
	return l.WithLogID(fmt.Sprintf("%s", logId))
}

// WithLogID returns a cloned logger with the given log ID attached.
func (l *Logger) WithLogID(id string) *Logger {
	if id == "" {
		id = nanoId.Must()
	}
	cloned := l.clone()
	cloned.logID = id
	cloned.Logger = cloned.With(zap.String(FieldLogID, id))
	return cloned
}

// WithErr returns a cloned logger with err attached as a field.
func (l *Logger) WithErr(err error) *Logger {
	if err == nil {
		return l
	}
	cloned := l.clone()
	cloned.Logger = cloned.With(zap.String(FieldError, err.Error()))
	return cloned
}

func (l *Logger) WithFields(fields map[string]any) *Logger {
	cloned := l.clone()
	logFields := make([]zap.Field, 0)

	for key, value := range fields {
		logFields = append(logFields, zap.Any(key, value))
	}

	cloned.Logger = cloned.With(logFields...)

	return cloned
}
