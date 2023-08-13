package ui

import "github.com/cruffinoni/rimworld-editor/generated"

// Options is a struct that contains the Options for the application.
type Options struct {
	Verbose                 bool
	Generate                bool
	Save                    bool
	Output                  string
	Input                   string
	MaxSaveGameFileDiscover int
	Mode                    string
	OperatingSystem         string
}

type Mode interface {
	Execute(args []string) error
	Init(options *Options, save *generated.Savegame)
}
