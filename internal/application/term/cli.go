package term

import (
	"fmt"
	"strings"

	"github.com/cruffinoni/rimworld-editor/generated"
	"github.com/cruffinoni/rimworld-editor/internal/application/commandline"
	"github.com/cruffinoni/rimworld-editor/internal/application/term/faction"
	"github.com/cruffinoni/rimworld-editor/internal/application/term/pawn"
	"github.com/cruffinoni/rimworld-editor/internal/application/ui"
	"github.com/cruffinoni/rimworld-editor/pkg/logging"

	"github.com/c-bata/go-prompt"
)

type Console struct {
	ui.Mode
	opt        *ui.Options
	shouldExit bool
	save       *generated.Savegame
	cliApp     commandline.App
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
		_ = c.cliApp.Run(f)
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

	c.cliApp = commandline.NewMowApp("rimworld-editor", "Rimworld save game editor")
	c.cliApp.SetErrorHandling(commandline.ErrorHandlingContinue)
	c.cliApp.SetHidden(true)
	c.cliApp.Command("exit", "Exit the console", func(cmd commandline.Command) {
		cmd.Action(c.exit)
	})
	pawn.RegisterPawnCommands(c.logger, c.cliApp, rp, rf, c.save)

	c.cliApp.Command("faction", "Faction commands", func(cmd commandline.Command) {
		cmd.Command("list", "List all factions", func(cmd commandline.Command) {
			cmd.Action(fl.ListAllFactions)
		})
	})

	c.cliApp.Command("world", "Interact with the world of Rimworld", func(cmd commandline.Command) {
		cmd.Command("growth-plants", "Make all plant to grow at a percentage", func(cmd commandline.Command) {
			percent := cmd.Float64Arg("PERCENTAGE", 100.0, "Percent of growth to attribute to plants")
			cmd.Action(func() {
				c.growAllPlants(*percent)
			})
		})
	})
}

func (c *Console) SetLogger(logger logging.Logger) {
	c.logger = logger
}
