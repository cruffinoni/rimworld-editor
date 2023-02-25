package pawn

import (
	"log"

	"github.com/cruffinoni/rimworld-editor/generated"
)

type List struct {
	sg *generated.Savegame
	rp PawnsRegisterer
}

func NewList(sg *generated.Savegame, rp PawnsRegisterer) *List {
	return &List{
		sg: sg,
		rp: rp,
	}
}

func (l *List) ListAllPawns(_ []string) error {
	for k, v := range l.rp {
		log.Printf("Pawns %s registered", k)
		log.Printf("Full name: %s %s %s", v.Name.First, v.Name.Name, v.Name.Last)
	}
	return nil
}
