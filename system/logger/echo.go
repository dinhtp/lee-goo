package logger

import (
	"fmt"
	"io"

	echoLog "github.com/labstack/gommon/log"
	"github.com/mattn/go-colorable"
)

func (l *Logger) Output() io.Writer {
	return colorable.NewColorableStdout()
}

func (l *Logger) SetOutput(w io.Writer) {}

func (l *Logger) Prefix() string {
	return l.prefix
}

func (l *Logger) SetPrefix(p string) {
	l.prefix = p
}

func (l *Logger) Level() echoLog.Lvl {
	switch l.level {
	case LevelDebug:
		return echoLog.DEBUG
	case LevelInfo:
		return echoLog.INFO
	case LevelWarn:
		return echoLog.WARN
	case LevelError:
		return echoLog.ERROR
	case LevelPanic:
		return echoLog.Lvl(6)
	case LevelFatal:
		return echoLog.Lvl(7)
	default:
		return echoLog.OFF
	}
}

func (l *Logger) SetLevel(v echoLog.Lvl) {
	switch v {
	case echoLog.DEBUG:
		l.level = LevelDebug
	case echoLog.INFO:
		l.level = LevelInfo
	case echoLog.WARN:
		l.level = LevelWarn
	case echoLog.ERROR:
		l.level = LevelError
	case echoLog.Lvl(6):
		l.level = LevelPanic
	case echoLog.Lvl(7):
		l.level = LevelFatal
	}
}

func (l *Logger) SetHeader(h string) {
	l.WithErr(fmt.Errorf("%s", h))
}

func (l *Logger) Print(i ...interface{}) {
	fmt.Println(i...)
}

func (l *Logger) Printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (l *Logger) Printj(j echoLog.JSON) {
	fmt.Println(j)
}

func (l *Logger) Debug(i ...interface{}) {
	l.Sugar().Debug(i...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.Sugar().Debugf(format, args...)
}

func (l *Logger) Debugj(j echoLog.JSON) {
	l.Sugar().Debug(j)
}

func (l *Logger) Info(i ...interface{}) {
	l.Sugar().Info(i...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.Sugar().Infof(format, args...)
}

func (l *Logger) Infoj(j echoLog.JSON) {
	l.Sugar().Info(j)
}

func (l *Logger) Warn(i ...interface{}) {
	l.Sugar().Warn(i...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.Sugar().Warnf(format, args...)
}

func (l *Logger) Warnj(j echoLog.JSON) {
	l.Sugar().Warn(j)
}

func (l *Logger) Error(i ...interface{}) {
	l.Sugar().Error(i...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Sugar().Errorf(format, args...)
}

func (l *Logger) Errorj(j echoLog.JSON) {
	l.Sugar().Error(j)
}

func (l *Logger) Fatal(i ...interface{}) {
	l.Sugar().Fatal(i...)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.Sugar().Fatalf(format, args...)
}

func (l *Logger) Fatalj(j echoLog.JSON) {
	l.Sugar().Fatal(j)
}

func (l *Logger) Panic(i ...interface{}) {
	l.Sugar().Panic(i...)
}

func (l *Logger) Panicf(format string, args ...interface{}) {
	l.Sugar().Panicf(format, args...)
}

func (l *Logger) Panicj(j echoLog.JSON) {
	l.Sugar().Panic(j)
}
