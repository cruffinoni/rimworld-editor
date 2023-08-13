package pawn

import (
	"reflect"

	"github.com/cruffinoni/rimworld-editor/cmd/app/ui/term/faction"
	"github.com/cruffinoni/rimworld-editor/generated"
	"github.com/cruffinoni/rimworld-editor/xml/types/iterator"
)

type PawnsRegisterer map[string]*generated.Thing

var commonFieldsPawnsAliveThing map[string]reflect.Type

func init() {
	comp := reflect.TypeOf(generated.Thing{})
	ref := reflect.TypeOf(generated.PawnsAlive{})
	commonFieldsPawnsAliveThing = make(map[string]reflect.Type, ref.NumField())
	for i := 0; i < ref.NumField(); i++ {
		field := ref.Field(i)
		if s, ok := comp.FieldByName(field.Name); ok {
			commonFieldsPawnsAliveThing[field.Name] = s.Type
		}
	}
}

func IsPlayerPawn(pawn *generated.Thing, rf faction.Registerer) bool {
	if pawn == nil {
		return false
	}
	if f, ok := rf[pawn.Faction]; !ok {
		return false
	} else {
		return faction.IsPlayerFaction(f)
	}
}

func RegisterPawns(sg *generated.Savegame, rf faction.Registerer) PawnsRegisterer {
	pawns := PawnsRegisterer{}
	for i := iterator.NewSliceIterator[*generated.PawnsAlive](sg.Game.World.WorldPawns.PawnsAlive); i.HasNext(); i = i.Next() {
		v := i.Value()
		thingName := "Thing_" + v.Id
		thing := &generated.Thing{}
		thingValue := reflect.ValueOf(thing).Elem()
		for f := range commonFieldsPawnsAliveThing {
			tf := thingValue.FieldByName(f)
			vf := reflect.ValueOf(v).Elem().FieldByName(f)
			if tf.Type() != vf.Type() {
				continue
			}
			thingValue.FieldByName(f).Set(reflect.ValueOf(reflect.ValueOf(v).Elem().FieldByName(f).Interface()))
		}
		pawns[thingName] = thing
	}
	for i := iterator.NewSliceIterator[*generated.Maps](sg.Game.Maps); i.HasNext(); i = i.Next() {
		for j := iterator.NewSliceIterator[*generated.Thing](i.Value().Things); j.HasNext(); j = j.Next() {
			v := j.Value()
			if v.Attr.Get("Class") == "Pawn" && v.Def == "Human" && IsPlayerPawn(v, rf) {
				pawns["Thing_"+v.Id] = v
			}
		}
	}
	return pawns
}
