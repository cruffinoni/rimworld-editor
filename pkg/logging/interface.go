package logging

import (
	"io"
)

// Fields represents a map of fields for structured logging
type Fields map[string]any

// Logger defines the interface for structured logging
type Logger interface {
	Trace(msg string)
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Fatal(msg string)
	Panic(msg string)

	// WithField returns a new logger with an additional field for structured logging
	WithField(key string, value any) Logger
	// WithFields returns a new logger with multiple additional fields for structured logging
	WithFields(fields Fields) Logger
	// WithError returns a new logger with an error field for structured logging
	WithError(err error) Logger

	SetOutput(output io.Writer)
}
