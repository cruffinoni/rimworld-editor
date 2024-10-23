package pawn

import (
	"github.com/cruffinoni/printer"

	"github.com/cruffinoni/rimworld-editor/cmd/app/ui/term/faction"
	"github.com/cruffinoni/rimworld-editor/generated"
)

type List struct {
	sg *generated.Savegame
	rp PawnsRegisterer
	rf faction.Registerer
}

func NewList(sg *generated.Savegame, rp PawnsRegisterer, rf faction.Registerer) *List {
	return &List{
		sg: sg,
		rp: rp,
		rf: rf,
	}
}

func (l *List) ListAllPawns() {
	for k, v := range l.rp {
		printer.Debugf("Pawn {{{-BOLD}}}%s{{{-RESET}}} registered", k)
		if v.Name == nil {
			printer.Debugf("name is nil: %v", v.Name)
			continue
		}
		printer.Debugf("Full name: %s", getPawnFullNameColorFormatted(v))
		if fac, ok := l.rf[v.Faction]; ok {
			faction.PrintFactionInformation(l.rf, fac, false)
		}
		printer.Debugf("")
	}
}
