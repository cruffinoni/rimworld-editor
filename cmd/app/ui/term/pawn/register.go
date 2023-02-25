package pawn

import (
	"github.com/cruffinoni/rimworld-editor/generated"
	"github.com/cruffinoni/rimworld-editor/xml/types/iterator"
)

type PawnsRegisterer map[string]*generated.PawnsAlive

func RegisterPawns(sg *generated.Savegame) map[string]*generated.PawnsAlive {
	pawns := make(map[string]*generated.PawnsAlive, sg.Game.World.WorldPawns.PawnsAlive.Capacity())
	for i := iterator.NewSliceIterator[*generated.PawnsAlive](sg.Game.World.WorldPawns.PawnsAlive); i.HasNext(); i = i.Next() {
		v := i.Value()
		pawns["Thing_"+v.Id] = v
	}
	return pawns
}
