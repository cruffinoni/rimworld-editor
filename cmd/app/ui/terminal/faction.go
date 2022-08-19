package terminal

import (
	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/types/algorithm"
	"log"
)

func (c *Console) factionList(args []string) error {
	if args[0] == "all" || len(args) == 0 {
		log.Println("Summary of all factions...")
		for i := c.save.Game.World.Info.FactionCounts.Iterator(); i != nil; i = i.Next() {
			log.Printf("%s: %d", i.Key(), i.Value())
		}
	} else {
		f, ok := algorithm.FindIf(&c.save.Game.World.FactionManager.AllFactions, func(e *xml.Element) bool {
			return true
		})
	}
	return nil
}

func (c *Console) factionCreate(args []string) error {
	return nil
}

func (c *Console) factionDelete(args []string) error {
	return nil
}
