package logger

type FieldCtxKey string

const (
	FieldLogID  = "log_id"
	FieldError  = "error"
	FieldHeader = "header"
)

const (
	LogFieldLogIDCtxKey FieldCtxKey = "log_id"
)

type Option func(*Logger)

func WithLevel(lvl Level) Option {
	return func(l *Logger) {
		l.level = lvl
	}
}

func WithLogID(logID string) Option {
	return func(l *Logger) {
		l.logID = logID
	}
}
