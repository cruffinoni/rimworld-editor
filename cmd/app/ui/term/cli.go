package term

import (
	"fmt"
	"strings"

	"github.com/cruffinoni/rimworld-editor/cmd/app/ui"
	"github.com/cruffinoni/rimworld-editor/cmd/app/ui/term/faction"
	"github.com/cruffinoni/rimworld-editor/cmd/app/ui/term/pawn"
	"github.com/cruffinoni/rimworld-editor/cmd/app/ui/term/printer"
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
			printer.PrintError(err)
			if err != errUnknownCommand {
				c.shouldExit = true
			}
		}
		if c.shouldExit {
			printer.PrintS("Execution ended.")
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
	rf := faction.RegisterFactions(c.save)
	rp := pawn.RegisterPawns(c.save, rf)
	pl := pawn.NewList(c.save, rp, rf)
	pawnCmd := c.commands.RegisterCommand(&config{
		name:        "pawn",
		description: "Pawn commands",
	})
	pawnCmd.RegisterCommand(&config{
		name:        "world",
		description: "List all pawns that alive in the game including your pawns and the faction leaders",
		handler:     pl.ListAllPawns,
	})

	pi := pawn.NewInjury(rp)
	pawnCmd.RegisterCommand(&config{
		name:        "injury",
		description: "Commands to manipulate pawn's injury",
	}).RegisterCommand(&config{
		name:        "remove-all",
		description: "Remove all injuries from a pawn",
		handler:     pi.RemoveInjuries,
	}, &config{
		name:        "list",
		description: "List all injuries from a pawn",
		handler:     pi.List,
	})

	fl := faction.NewList(c.save, rf)
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
