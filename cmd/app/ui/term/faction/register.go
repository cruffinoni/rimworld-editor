package faction

import (
	"strconv"

	"github.com/cruffinoni/printer"

	"github.com/cruffinoni/rimworld-editor/generated"
	"github.com/cruffinoni/rimworld-editor/internal/xml/types/iterator"
)

type Registerer map[string]*generated.AllFactions

func RegisterFactions(sg *generated.Savegame) Registerer {
	allFac := Registerer{}
	ite := iterator.NewSliceIterator[*generated.AllFactions](sg.Game.World.FactionManager.AllFactions)
	for i := ite; i.HasNext(); i = i.Next() {
		v := i.Value()
		allFac[GetFactionID(v.LoadId)] = v
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
	id := GetFactionID(f.LoadId)
	printer.Debugf("Faction %s (named '%s') - %s", f.Def, f.Name, id)
	if f.Def == "PlayerColony" {
		printer.Debugf("{{{-BOLD,F_GREEN}}}This is the player's faction")
	} else {
		printer.Debugf("{{{-F_MAGENTA}}}This is a faction controlled by the IA")
	}
	if f.Leader == "null" {
		printer.Debugf("The faction doesn't have any leader")
	} else {
		printer.Debugf("Faction's leader: {{{-BOLD}}}%s", f.Leader)
	}
	if withRelations {
		printer.Debugf("Relations:")
		for i := iterator.NewSliceIterator[*generated.Relations](f.Relations); i.HasNext(); i = i.Next() {
			r := i.Value()
			if r.Goodwill > 75 {
				printer.Debugf("\t- %s (%s) => %s (%s) : {{{-BOLD,F_GREEN}}}%d",
					f.Def, id, rf[r.Other].Def, GetFactionID(rf[r.Other].LoadId), r.Goodwill)
			} else if r.Goodwill < -50 {
				printer.Debugf("\t- %s (%s) => %s (%s) : {{{-BOLD,F_RED}}}%d",
					f.Def, id, rf[r.Other].Def, GetFactionID(rf[r.Other].LoadId), r.Goodwill)
			} else {
				printer.Debugf("\t- %s (%s) => %s (%s) : {{{-BOLD,F_YELLOW}}}%d",
					f.Def, id, rf[r.Other].Def, GetFactionID(rf[r.Other].LoadId), r.Goodwill)
			}
		}
	}
}
