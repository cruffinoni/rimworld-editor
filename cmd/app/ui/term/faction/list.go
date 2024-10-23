package faction

import (
	"github.com/cruffinoni/printer"

	"github.com/cruffinoni/rimworld-editor/generated"
	"github.com/cruffinoni/rimworld-editor/internal/xml/types/iterator"
)

type List struct {
	sg *generated.Savegame
	r  Registerer
}

func NewList(sg *generated.Savegame, r Registerer) *List {
	return &List{r: r, sg: sg}
}

func (l *List) ListAllFactions() {
	allFac := map[string]*generated.AllFactions{}
	printer.Print("Summary of all factions...")
	ite := iterator.NewSliceIterator[*generated.AllFactions](l.sg.Game.World.FactionManager.AllFactions)
	for i := ite; i.HasNext(); i = i.Next() {
		v := i.Value()
		allFac[GetFactionID(v.LoadId)] = v
	}
	for _, f := range allFac {
		PrintFactionInformation(l.r, f, true)
		printer.Print("")
	}
	return
}
