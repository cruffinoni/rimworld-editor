package terminal

import (
	"fmt"
	"log"
	"strings"

	"github.com/cruffinoni/rimworld-editor/cmd/app/ui"
	"github.com/cruffinoni/rimworld-editor/generated"

	"github.com/c-bata/go-prompt"
)

type Console struct {
	ui.Mode
	opt        *ui.Options
	shouldExit bool
	commands   *terminalCommand
	save       *generated.Save
}

func (c *Console) completer(d prompt.Document) []prompt.Suggest {
	return nil
}

func (c *Console) Execute([]string) error {
	fmt.Println("Welcome...")
	var lastError error
	for {
		input := prompt.Input("-> ", c.completer)
		f := strings.Fields(input)
		if len(f) == 0 {
			continue
		}
		if err := c.commands.Parse(f); err != nil {
			log.Println(err)
			if err != errUnknownCommand {
				c.shouldExit = true
			}
		}
		log.Println("Execution ended.")
		if c.shouldExit {
			break
		}
	}
	return lastError
}

func (c *Console) exit([]string) error {
	c.shouldExit = true
	return nil
}

func (c *Console) Init(options *ui.Options, save *generated.Save) {
	c.opt = options
	c.save = save
	c.commands = NewTerminalCommands()
	c.commands.RegisterCommand(&details{
		name:        "exit",
		description: "Exit the console",
		handler:     c.exit,
	})
	c.commands.RegisterCommand(&details{
		name:        "pawn",
		description: "Pawn commands",
	})
	c.commands.RegisterCommand(&details{
		name:        "faction",
		description: "Faction commands",
	}).RegisterCommand(
		&details{
			name:        "list",
			description: "List factions",
			handler:     c.factionList,
		},
		&details{
			name:        "create",
			description: "Create a faction",
			handler:     c.factionCreate,
		},
		&details{
			name:        "delete",
			description: "Delete a faction",
			handler:     c.factionDelete,
		},
	)
}
