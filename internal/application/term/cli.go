package term

import (
	"flag"
	"fmt"
	"strings"

	cli "github.com/jawher/mow.cli"

	"github.com/cruffinoni/rimworld-editor/generated"
	"github.com/cruffinoni/rimworld-editor/internal/application/terminal/faction"
	"github.com/cruffinoni/rimworld-editor/internal/application/terminal/pawn"
	"github.com/cruffinoni/rimworld-editor/internal/application/userinterface"
	"github.com/cruffinoni/rimworld-editor/pkg/logging"

	"github.com/c-bata/go-prompt"
)

type Console struct {
	ui.Mode
	opt        *ui.Options
	shouldExit bool
	save       *generated.Savegame
	cli        *cli.Cli
	logger     logging.Logger
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
		f = append([]string{"rimworld"}, f...)
		_ = c.cli.Run(f)
		if c.shouldExit {
			c.logger.Debug("Execution ended")
			break
		}
	}
	return lastError
}

func (c *Console) exit() {
	c.shouldExit = true
	return
}

func (c *Console) Init(options *ui.Options, save *generated.Savegame) {
	c.opt = options
	c.save = save
	c.logger.Debug("Console init")
	rf := faction.RegisterFactions(c.logger, c.save)
	rp := pawn.RegisterPawns(c.save, rf)
	fl := faction.NewList(c.logger, c.save, rf)

	c.cli = cli.App("rimworld-editor", "Rimworld save game editor")
	c.cli.ErrorHandling = flag.ContinueOnError
	c.cli.Hidden = true
	c.cli.Command("exit", "Exit the console", cli.ActionCommand(c.exit))
	pawn.RegisterPawnCommands(c.logger, c.cli, rp, rf, c.save)

	c.cli.Command("faction", "Faction commands", func(cmd *cli.Cmd) {
		cmd.Command("list", "List all factions", cli.ActionCommand(fl.ListAllFactions))
	})

	c.cli.Command("world", "Interact with the world of Rimworld", func(cmd *cli.Cmd) {
		cmd.Command("growth-plants", "Make all plant to grow at a percentage", func(cmd *cli.Cmd) {
			percent := cmd.Float64Arg("PERCENTAGE", 100.0, "Percent of growth to attribute to plants")
			cmd.Action = func() {
				c.growAllPlants(*percent)
			}
		})
	})
}

func (c *Console) SetLogger(logger logging.Logger) {
	c.logger = logger
}
