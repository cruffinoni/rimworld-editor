package commandline

import (
	"flag"

	cli "github.com/jawher/mow.cli"
)

type MowApp struct {
	app *cli.Cli
}

type MowCommand struct {
	cmd *cli.Cmd
}

func NewMowApp(name, description string) *MowApp {
	return &MowApp{
		app: cli.App(name, description),
	}
}

func (m *MowApp) Command(name, description string, configure func(Command)) {
	m.app.Command(name, description, func(cmd *cli.Cmd) {
		if configure != nil {
			configure(&MowCommand{cmd: cmd})
		}
	})
}

func (m *MowApp) BoolOptPtr(dest *bool, name string, value bool, desc string) {
	m.app.BoolOptPtr(dest, name, value, desc)
}

func (m *MowApp) StringOptPtr(dest *string, name string, value string, desc string) {
	m.app.StringOptPtr(dest, name, value, desc)
}

func (m *MowApp) IntOptPtr(dest *int, name string, value int, desc string) {
	m.app.IntOptPtr(dest, name, value, desc)
}

func (m *MowApp) Version(command, version string) {
	m.app.Version(command, version)
}

func (m *MowApp) SetBefore(before func()) {
	m.app.Before = before
}

func (m *MowApp) SetAction(action func()) {
	m.app.Action = action
}

func (m *MowApp) SetHidden(hidden bool) {
	m.app.Hidden = hidden
}

func (m *MowApp) SetErrorHandling(mode ErrorHandling) {
	switch mode {
	case ErrorHandlingContinue:
		m.app.ErrorHandling = flag.ContinueOnError
	default:
		m.app.ErrorHandling = flag.ExitOnError
	}
}

func (m *MowApp) PrintHelp() {
	m.app.PrintHelp()
}

func (m *MowApp) Run(args []string) error {
	return m.app.Run(args)
}

func (m *MowApp) Exit(code int) {
	cli.Exit(code)
}

func (m *MowCommand) Command(name, description string, configure func(Command)) {
	m.cmd.Command(name, description, func(cmd *cli.Cmd) {
		if configure != nil {
			configure(&MowCommand{cmd: cmd})
		}
	})
}

func (m *MowCommand) StringArg(name, value, desc string) *string {
	return m.cmd.StringArg(name, value, desc)
}

func (m *MowCommand) StringsArg(name string, value []string, desc string) *[]string {
	return m.cmd.StringsArg(name, value, desc)
}

func (m *MowCommand) IntArg(name string, value int, desc string) *int {
	return m.cmd.IntArg(name, value, desc)
}

func (m *MowCommand) Float64Arg(name string, value float64, desc string) *float64 {
	return m.cmd.Float64Arg(name, value, desc)
}

func (m *MowCommand) SetSpec(spec string) {
	m.cmd.Spec = spec
}

func (m *MowCommand) Action(action func()) {
	m.cmd.Action = action
}
