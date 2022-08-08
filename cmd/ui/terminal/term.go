package terminal

import (
	"errors"
	"flag"
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/cruffinoni/rimworld-editor/cmd/ui"
	cli "github.com/jawher/mow.cli"
	"strings"
)

type Console struct {
	ui.Mode
	opt *ui.Options
	*cli.Cli
}

func (c *Console) completer(d prompt.Document) []prompt.Suggest {
	return nil
}

var (
	errIncorrectUsage = errors.New("incorrect usage")
)

func (c *Console) Execute([]string) error {
	fmt.Println("Welcome...")
	for {
		input := prompt.Input("-> ", c.completer)
		if input == "exit" {
			return nil
		}
		f := append([]string{"cmd"}, strings.Fields(input)...)
		if err := c.Cli.Run(f); err != nil && err.Error() != "incorrect usage" {
			return err
		}
	}
}

func (c *Console) Init(options *ui.Options) {
	c.opt = options
	c.Cli = cli.App("console", "Rimworld save game editor console")
	/*
		Pawn
			- List
			- Edit
		Faction
			- List
			- Edit
		World
	*/
	c.ErrorHandling = flag.ContinueOnError
	c.Cli.Command("pawn", "Pawn commands", func(cmd *cli.Cmd) {
		cmd.Command("list", "List all pawns", func(cmd *cli.Cmd) {
			cmd.Action = func() {
				fmt.Println("list pawns")
			}
		})
	})
}
