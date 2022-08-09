package ui

import "github.com/cruffinoni/rimworld-editor/generated"

// Options is a struct that contains the Options for the application.
type Options struct {
	Verbose  bool
	Generate bool
	Output   string
	Input    string
	Mode     string
}

type Mode interface {
	Execute(args []string) error
	Init(options *Options, save *generated.Save)
}
