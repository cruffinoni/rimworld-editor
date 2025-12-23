package logging

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/cruffinoni/rimworld-editor/internal/config"
)

type logger struct {
	entry *logrus.Entry
}

func (l *logger) Write(p []byte) (n int, err error) {
	return l.entry.WriterLevel(logrus.InfoLevel).Write(p)
}

// Implement Logger interface methods
func (l *logger) Trace(msg string) {
	l.entry.Trace(msg)
}

func (l *logger) Debug(msg string) {
	l.entry.Debug(msg)
}

func (l *logger) Info(msg string) {
	l.entry.Info(msg)
}

func (l *logger) Warn(msg string) {
	l.entry.Warn(msg)
}

func (l *logger) Error(msg string) {
	l.entry.Error(msg)
}

func (l *logger) Fatal(msg string) {
	l.entry.Fatal(msg)
}

func (l *logger) Panic(msg string) {
	l.entry.Panic(msg)
}

func (l *logger) WithField(key string, value any) Logger {
	return &logger{entry: l.entry.WithField(key, value)}
}

func (l *logger) WithFields(fields Fields) Logger {
	logrusFields := make(logrus.Fields)
	for k, v := range fields {
		logrusFields[k] = v
	}
	return &logger{entry: l.entry.WithFields(logrusFields)}
}

func (l *logger) WithError(err error) Logger {
	return &logger{entry: l.entry.WithError(err)}
}

func (l *logger) SetOutput(output io.Writer) {
	l.entry.Logger.SetOutput(output)
}

func getLogLevel(level string) logrus.Level {
	switch strings.ToLower(level) {
	case "trace":
		return logrus.TraceLevel
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warn":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	default:
		return logrus.InfoLevel
	}
}

func resolveOutput(output string) io.Writer {
	switch strings.ToLower(strings.TrimSpace(output)) {
	case "", "stdout":
		return os.Stdout
	case "stderr":
		return os.Stderr
	default:
		f, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
		if err != nil {
			return os.Stdout
		}
		return f
	}
}

// NewLogger creates a new logger instance with the specified component name and configuration
func NewLogger(component string, cfg *config.LoggingConfig) Logger {
	if cfg == nil {
		cfg = &config.LoggingConfig{
			Level:  "info",
			Format: "text",
			Output: "stdout",
		}
	}
	baseLogger := logrus.New()

	if cfg.Format == "json" {
		baseLogger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339Nano,
		})
	} else {
		baseLogger.SetFormatter(&logrus.TextFormatter{
			DisableColors:   true,
			FullTimestamp:   true,
			TimestampFormat: "15:04:05",
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				filename := path.Base(f.File)
				return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
			},
			DisableQuote: true,
			SortingFunc: func(keys []string) {
				sort.Strings(keys)
			},
			PadLevelText: true,
		})
	}

	baseLogger.SetOutput(resolveOutput(cfg.Output))
	baseLogger.SetLevel(getLogLevel(cfg.Level))
	entry := baseLogger.WithFields(logrus.Fields{
		"component": component,
	})

	return &logger{entry: entry}
}
