package commandline

type ErrorHandling int

const (
	ErrorHandlingDefault ErrorHandling = iota
	ErrorHandlingContinue
)

type Command interface {
	Command(name, description string, configure func(Command))
	StringArg(name, value, desc string) *string
	StringsArg(name string, value []string, desc string) *[]string
	IntArg(name string, value int, desc string) *int
	Float64Arg(name string, value float64, desc string) *float64
	SetSpec(spec string)
	Action(action func())
}

type CommandRegistry interface {
	Command(name, description string, configure func(Command))
}

type App interface {
	CommandRegistry
	BoolOptPtr(dest *bool, name string, value bool, desc string)
	StringOptPtr(dest *string, name string, value string, desc string)
	IntOptPtr(dest *int, name string, value int, desc string)
	Version(command, version string)
	SetBefore(before func())
	SetAction(action func())
	SetHidden(hidden bool)
	SetErrorHandling(mode ErrorHandling)
	PrintHelp()
	Run(args []string) error
	Exit(code int)
}
