package faction

import (
	"github.com/cruffinoni/rimworld-editor/cmd/app/ui/term/printer"
	"github.com/cruffinoni/rimworld-editor/generated"
	"github.com/cruffinoni/rimworld-editor/xml/types/iterator"
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
	printer.PrintS("Summary of all factions...")
	ite := iterator.NewSliceIterator[*generated.AllFactions](l.sg.Game.World.FactionManager.AllFactions)
	for i := ite; i.HasNext(); i = i.Next() {
		v := i.Value()
		allFac[GetFactionID(v.LoadID)] = v
	}
	for _, f := range allFac {
		PrintFactionInformation(l.r, f, true)
		printer.PrintS("")
	}
	return
}
