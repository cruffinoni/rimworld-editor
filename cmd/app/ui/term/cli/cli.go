package cli

import "github.com/cruffinoni/rimworld-editor/pkg/logging"

type Cli struct {
	command *Command
}

func NewCli(logger logging.Logger) *Cli {
	return &Cli{
		command: &Command{logger: logger},
	}
}

func (c *Cli) NewCommand(cmd *Command, f CmdSpan, params ...Params) {
	cmd.logger = c.command.logger
	c.command.NewCommand(cmd, f, params...)
}

func (c *Cli) Parse(args []string) error {
	return c.command.RunWithArgs(args)
}
