package ui

// Options is a struct that contains the Options for the application.
type Options struct {
	Verbose  bool
	Generate bool
	Output   string
	Input    string
	Mode     string
}

type Mode interface {
	Execute() error
	NewMode(options *Options) Mode
}
