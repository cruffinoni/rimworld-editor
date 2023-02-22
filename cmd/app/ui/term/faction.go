package term

import (
	"github.com/cruffinoni/rimworld-editor/algorithm"
	"github.com/cruffinoni/rimworld-editor/generated"
	"github.com/cruffinoni/rimworld-editor/xml/types/iterator"

	"log"
)

func (c *Console) factionList(args []string) error {
	if len(args) == 0 || args[0] == "all" {
		log.Println("Summary of all factions...")
		log.Printf("Factions: %+v", c.save.Game.World.Info.Factions)
		ite := iterator.NewSliceIterator[string](c.save.Game.World.Info.Factions)
		for i := ite; i.HasNext(); i = i.Next() {
			log.Printf("%s", i.Value())
		}
	} else {
		f, ok := algorithm.FindInSlice[*generated.AllFactions](c.save.Game.World.FactionManager.AllFactions, &generated.AllFactions{Name: args[0]})
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
