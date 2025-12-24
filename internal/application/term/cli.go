package term

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cruffinoni/rimworld-editor/generated"
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
	rootCmd    *cobra.Command
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
		c.rootCmd.SetArgs(f)
		if err := c.rootCmd.Execute(); err != nil {
			lastError = err
			c.logger.WithError(err).Error("Command failed")
		}
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
	c.rootCmd = c.newRootCommand()
}

func (c *Console) SetLogger(logger logging.Logger) {
	c.logger = logger
}

func (c *Console) newRootCommand() *cobra.Command {
	rf := faction.RegisterFactions(c.logger, c.save)
	rp := pawn.RegisterPawns(c.save, rf)

	rootCmd := &cobra.Command{
		Use:           "rimworld-editor",
		Short:         "Rimworld save game editor",
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	rootCmd.AddCommand(&cobra.Command{
		Use:   "exit",
		Short: "Exit the console",
		Run: func(cmd *cobra.Command, args []string) {
			c.exit()
		},
	})

	rootCmd.AddCommand(faction.NewCommand(c.logger, c.save, rf))
	rootCmd.AddCommand(pawn.NewCommand(c.logger, c.save, rp, rf))

	worldCmd := &cobra.Command{
		Use:   "world",
		Short: "Interact with the world of Rimworld",
	}
	worldCmd.AddCommand(&cobra.Command{
		Use:   "growth-plants <PERCENTAGE>",
		Short: "Make all plant to grow at a percentage",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			percent, err := strconv.ParseFloat(args[0], 64)
			if err != nil {
				return fmt.Errorf("invalid percentage: %w", err)
			}
			c.growAllPlants(percent)
			return nil
		},
	})
	rootCmd.AddCommand(worldCmd)

	return rootCmd
}
