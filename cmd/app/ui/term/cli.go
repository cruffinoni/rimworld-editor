package term

import (
	"fmt"
	"log"
	"strings"

	"github.com/cruffinoni/rimworld-editor/cmd/app/ui"
	"github.com/cruffinoni/rimworld-editor/cmd/app/ui/term/faction"
	"github.com/cruffinoni/rimworld-editor/generated"

	"github.com/c-bata/go-prompt"
)

type Console struct {
	ui.Mode
	opt        *ui.Options
	shouldExit bool
	commands   *TerminalCommand
	save       *generated.Savegame
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
		if c.shouldExit {
			log.Println("Execution ended.")
			break
		}
	}
	return lastError
}

func (c *Console) exit([]string) error {
	c.shouldExit = true
	return nil
}

func (c *Console) Init(options *ui.Options, save *generated.Savegame) {
	c.opt = options
	c.save = save
	c.commands = NewTerminalCommands()
	c.commands.RegisterCommand(&config{
		name:        "exit",
		description: "Exit the console",
		handler:     c.exit,
	})
	c.commands.RegisterCommand(&config{
		name:        "pawn",
		description: "Pawn commands",
	})
	fl := faction.List{SG: c.save}
	c.commands.RegisterCommand(&config{
		name:        "faction",
		description: "Faction commands",
	}).RegisterCommand(
		&config{
			name:        "list",
			description: "List all factions",
			handler:     fl.ListAllFactions,
		},
		//&config{
		//	name:        "create",
		//	description: "Create a faction",
		//	handler:     f.Create,
		//},
		//&config{
		//	name:        "delete",
		//	description: "Delete a faction",
		//	handler:     f.Delete,
		//},
	)
	c.commands.RegisterCommand(&config{
		name:        "world",
		description: "Commands to interact with the world of Rimeworld",
	}).RegisterCommand(
		&config{
			name:        "growth",
			description: "Make all plant to grow at a percentage",
			handler:     c.growAllPlants,
		})
}
