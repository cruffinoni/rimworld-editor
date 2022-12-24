package term

import (
	"errors"
	"fmt"
	"log"
)

type commandHandler func([]string) error

type details struct {
	name        string
	description string
	handler     commandHandler
}
type TerminalCommand struct {
	*details
	commands map[string]*TerminalCommand
}

func NewTerminalCommands() *TerminalCommand {
	return &TerminalCommand{
		commands: make(map[string]*TerminalCommand),
	}
}

func (t *TerminalCommand) RegisterCommand(cmd ...*details) *TerminalCommand {
	lastCreatedCmd := t
	for _, c := range cmd {
		if _, ok := t.commands[c.name]; ok {
			panic(fmt.Sprintf("Command %s already registered", c.name))
		}
		newT := &TerminalCommand{
			details:  c,
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
	fmt.Printf("%s: %s\n", t.details.name, t.details.description)
	for n, c := range t.commands {
		fmt.Printf("%s: %s\n", n, c.details.description)
	}
}

func (t *TerminalCommand) Parse(input []string) error {
	log.Printf("cmd: %+#v", t)
	if len(t.commands) == 0 || len(input) == 0 {
		if t.details.handler == nil {
			t.showHelp()
			return nil
		}
		return t.details.handler(input)
	}
	for n, c := range t.commands {
		if n == input[0] {
			log.Printf("Found command %s", n)
			return c.Parse(input[1:])
		}
	}
	return errUnknownCommand
}