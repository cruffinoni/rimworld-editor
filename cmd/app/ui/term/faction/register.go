package faction

import (
	"strconv"

	"github.com/cruffinoni/rimworld-editor/cmd/app/ui/term/printer"
	"github.com/cruffinoni/rimworld-editor/generated"
	"github.com/cruffinoni/rimworld-editor/xml/types/iterator"
)

type Registerer map[string]*generated.AllFactions

func RegisterFactions(sg *generated.Savegame) Registerer {
	allFac := Registerer{}
	ite := iterator.NewSliceIterator[*generated.AllFactions](sg.Game.World.FactionManager.AllFactions)
	for i := ite; i.HasNext(); i = i.Next() {
		v := i.Value()
		allFac[GetFactionID(v.LoadID)] = v
	}
	return allFac
}

func GetFactionID(loadID int64) string {
	return "Faction_" + strconv.FormatInt(loadID, 10)
}

func IsPlayerFaction(f *generated.AllFactions) bool {
	if f == nil {
		return false
	}
	return f.Def == "PlayerColony"
}

func PrintFactionInformation(rf Registerer, f *generated.AllFactions, withRelations bool) {
	id := GetFactionID(f.LoadID)
	printer.PrintSf("Faction %s (named '%s') - %s", f.Def, f.Name, id)
	if f.Def == "PlayerColony" {
		printer.PrintS("This is the player's faction")
	} else {
		printer.PrintS("This is a faction controlled by the IA")
	}
	if f.Leader == "null" {
		printer.PrintS("The faction doesn't have any leader")
	} else {
		printer.PrintSf("Faction's leader: %s", f.Leader)
	}
	if withRelations {
		printer.PrintS("Relations:")
		for i := iterator.NewSliceIterator[*generated.Relations](f.Relations); i.HasNext(); i = i.Next() {
			r := i.Value()
			if r.Goodwill > 75 {
				printer.PrintSf("\t- %s (%s) => %s (%s) : {-F_GREEN}%d",
					f.Def, id, rf[r.Other].Def, GetFactionID(rf[r.Other].LoadID), r.Goodwill)
			} else if r.Goodwill < -50 {
				printer.PrintSf("\t- %s (%s) => %s (%s) : {-F_RED}%d",
					f.Def, id, rf[r.Other].Def, GetFactionID(rf[r.Other].LoadID), r.Goodwill)
			} else {
				printer.PrintSf("\t- %s (%s) => %s (%s) : {-F_YELLOW}%d",
					f.Def, id, rf[r.Other].Def, GetFactionID(rf[r.Other].LoadID), r.Goodwill)
			}
		}
	}
}
