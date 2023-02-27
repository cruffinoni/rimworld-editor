package pawn

import (
	"github.com/cruffinoni/rimworld-editor/cmd/app/ui/term/faction"
	"github.com/cruffinoni/rimworld-editor/cmd/app/ui/term/printer"
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

func (l *List) ListAllPawns(_ []string) error {
	for k, v := range l.rp {
		printer.PrintSf("Pawns {-BOLD}%s{-RESET} registered", k)
		if v.Name == nil {
			printer.PrintSf("name is nil: %v", v.Name)
			continue
		}
		if v.Name.Name != "" {
			printer.PrintSf("Full name: {-F_GREEN}%s {-F_MAGENTA}%s {-F_CYAN}%s", v.Name.First, v.Name.Name, v.Name.Last)
		} else {
			printer.PrintSf("Full name: {-F_GREEN}%s {-F_CYAN}%s", v.Name.First, v.Name.Last)
		}
		if fac, ok := l.rf[v.Faction]; ok {
			faction.PrintFactionInformation(l.rf, fac, false)
		}
		printer.Print([]byte{})
	}
	return nil
}
