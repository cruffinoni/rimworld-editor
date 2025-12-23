package _interface

import "github.com/cruffinoni/rimworld-editor/pkg/logging"

// LoggerSetter allows callers to inject a logger into XML-related types.
type LoggerSetter interface {
	SetLogger(logger logging.Logger)
}
