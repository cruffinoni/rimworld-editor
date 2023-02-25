package faction

import (
	"log"
	"strconv"

	"github.com/cruffinoni/rimworld-editor/generated"
	"github.com/cruffinoni/rimworld-editor/xml/types/iterator"
)

type List struct {
	SG *generated.Savegame
}

func GetFactionID(loadID int64) string {
	return "Faction_" + strconv.FormatInt(loadID, 10)
}

func (l *List) ListAllFactions(args []string) error {
	allFac := map[string]*generated.AllFactions{}
	log.Println("Summary of all factions...")
	ite := iterator.NewSliceIterator[*generated.AllFactions](l.SG.Game.World.FactionManager.AllFactions)
	for i := ite; i.HasNext(); i = i.Next() {
		v := i.Value()
		allFac[GetFactionID(v.LoadID)] = v
	}
	for id, f := range allFac {
		log.Printf("Faction %s (named '%s') - %s", f.Def, f.Name, id)
		if f.Def == "PlayerColony" {
			log.Println("This is the player's faction")
		} else {
			log.Println("This is a faction controlled by the IA")
		}
		if f.Leader == "null" {
			log.Println("The faction doesn't have any leader")
		} else {
			log.Printf("Faction's leader: %s", f.Leader)
		}
		log.Println("Relations:")
		for i := iterator.NewSliceIterator[*generated.Relations](f.Relations); i.HasNext(); i = i.Next() {
			r := i.Value()
			log.Printf("\t- %s (%s) => %s (%s) : %d",
				f.Def, id, allFac[r.Other].Def, GetFactionID(allFac[r.Other].LoadID), r.Goodwill)
		}
		log.Println("")
	}
	return nil
}
