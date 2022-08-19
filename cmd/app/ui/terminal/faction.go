package terminal

import (
	"github.com/cruffinoni/rimworld-editor/generated"
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
		f, ok := algorithm.FindIf(c.save.Game.World.FactionManager.AllFactions, func(a *generated.AllFactions) bool {
			return a.Name == args[0]
		})
		if !ok {
			log.Printf("faction '%v' not found", args[0])
			return nil
		}
		log.Printf("Faction %v => %v", f.Name, f.Def)
	}
	return nil
}

func (c *Console) factionCreate(args []string) error {
	return nil
}

func (c *Console) factionDelete(args []string) error {
	return nil
}
