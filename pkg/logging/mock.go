package logging

import "io"

// MockLogger provides a no-op implementation of the Logger interface for testing.
// It can be used in unit tests to avoid actual logging output while maintaining
// the same interface contract.
type MockLogger struct{}

// NewMockLogger creates a new instance of MockLogger.
func NewMockLogger() Logger {
	return &MockLogger{}
}

// WithField returns the same MockLogger instance (no-op).
func (m *MockLogger) WithField(_ string, _ any) Logger {
	return m
}

// WithFields returns the same MockLogger instance (no-op).
func (m *MockLogger) WithFields(_ Fields) Logger {
	return m
}

// WithError returns the same MockLogger instance (no-op).
func (m *MockLogger) WithError(_ error) Logger {
	return m
}

// Trace performs no operation.
func (m *MockLogger) Trace(_ string) {}

// Debug performs no operation.
func (m *MockLogger) Debug(_ string) {}

// Info performs no operation.
func (m *MockLogger) Info(_ string) {}

// Warn performs no operation.
func (m *MockLogger) Warn(_ string) {}

// Error performs no operation.
func (m *MockLogger) Error(_ string) {}

// Fatal performs no operation.
func (m *MockLogger) Fatal(_ string) {}

// Panic performs no operation.
func (m *MockLogger) Panic(_ string) {}

// Tracef performs no operation.
func (m *MockLogger) Tracef(_ string, _ ...any) {}

// Debugf performs no operation.
func (m *MockLogger) Debugf(_ string, _ ...any) {}

// Infof performs no operation.
func (m *MockLogger) Infof(_ string, _ ...any) {}

// Warnf performs no operation.
func (m *MockLogger) Warnf(_ string, _ ...any) {}

// Errorf performs no operation.
func (m *MockLogger) Errorf(_ string, _ ...any) {}

// Fatalf performs no operation.
func (m *MockLogger) Fatalf(_ string, _ ...any) {}

// Panicf performs no operation.
func (m *MockLogger) Panicf(_ string, _ ...any) {}

// SetOutput performs no operation.
func (m *MockLogger) SetOutput(_ io.Writer) {}
