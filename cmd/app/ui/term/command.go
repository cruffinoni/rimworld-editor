package term

import (
	"errors"
	"fmt"
	"strings"
)

type config struct {
	name        string
	description string
	handler     func([]string) error
}
type TerminalCommand struct {
	config   *config
	commands map[string]*TerminalCommand
}

func NewTerminalCommands() *TerminalCommand {
	return &TerminalCommand{
		commands: make(map[string]*TerminalCommand),
	}
}

func (t *TerminalCommand) RegisterCommand(cmd ...*config) *TerminalCommand {
	lastCreatedCmd := t
	for _, c := range cmd {
		if _, ok := t.commands[c.name]; ok {
			panic(fmt.Sprintf("Command %s already registered", c.name))
		}
		newT := &TerminalCommand{
			config:   c,
			commands: make(map[string]*TerminalCommand),
		}
		t.commands[c.name] = newT
		lastCreatedCmd = newT
	}
	if len(cmd) > 1 {
		return t
	}
	return lastCreatedCmd
}

var errUnknownCommand = errors.New("unknown command")

func (t *TerminalCommand) showHelp() {
	fmt.Printf("%s: %s\n", t.config.name, t.config.description)
	for n, c := range t.commands {
		fmt.Printf("%s: %s\n", n, c.config.description)
	}
}

func (t *TerminalCommand) Parse(input []string) error {
	//log.Printf("cmd: %+#v", t)
	if len(t.commands) == 0 || len(input) == 0 {
		if t.config.handler == nil {
			t.showHelp()
			return nil
		}
		return t.config.handler(input)
	}
	for n, c := range t.commands {
		if n == strings.ToLower(input[0]) {
			//log.Printf("Found command %s", n)
			return c.Parse(input[1:])
		}
	}
	return errUnknownCommand
}
