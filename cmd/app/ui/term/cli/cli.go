package cli

type Cli struct {
	command *Command
}

func NewCli() *Cli {
	return &Cli{
		command: &Command{},
	}
}

func (c *Cli) NewCommand(cmd *Command, f CmdSpan, params ...Params) {
	c.command.NewCommand(cmd, f, params...)
}

func (c *Cli) Parse(args []string) error {
	return c.command.RunWithArgs(args)
}
