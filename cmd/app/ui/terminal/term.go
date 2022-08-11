package terminal

import (
	"fmt"
	"github.com/cruffinoni/rimworld-editor/xml/saver"
	"github.com/cruffinoni/rimworld-editor/xml/saver/file"
	"log"
	"strings"

	"github.com/cruffinoni/rimworld-editor/cmd/app/ui"
	"github.com/cruffinoni/rimworld-editor/generated"

	"github.com/c-bata/go-prompt"
)

type commandHandler func([]string) error

type Console struct {
	ui.Mode
	opt        *ui.Options
	shouldExit bool
	commands   map[string]commandHandler
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
		if h, ok := c.commands[f[0]]; ok {
			lastError = h(f[1:])
			if lastError != nil {
				c.shouldExit = true
				log.Println(lastError)
			}
		} else {
			log.Printf("Unknown command: %s", f[0])
		}
		if c.shouldExit {
			break
		}
	}
	return lastError
}

func (c *Console) exit(_ []string) error {
	c.shouldExit = true
	return nil
}

func (c *Console) pawn(args []string) error {
	switch args[0] {
	case "list":
		log.Println("Listing pawns...")
	}
	return nil
}

func (c *Console) faction(args []string) error {
	switch args[0] {
	case "list":
		log.Println("Listing factions...")
		for i := c.save.Game.World.Info.FactionCounts.Iterator(); i != nil; i = i.Next() {
			log.Printf("%s: %d", i.Key(), i.Value())
		}
	}
	return nil
}

func (c *Console) savegame(args []string) error {
	path := "test/output_savegame.rws"
	if len(args) > 0 {
		path = args[0]
	}
	b := saver.NewBuffer()
	if err := file.Save(c.save, b, "savegame"); err != nil {
		return err
	}
	return b.ToFile(path)
}

func (c *Console) Init(options *ui.Options, save *generated.Save) {
	c.opt = options
	c.save = save
	c.commands = map[string]commandHandler{
		"exit":    c.exit,
		"pawn":    c.pawn,
		"faction": c.faction,
		"save":    c.savegame,
	}
}
