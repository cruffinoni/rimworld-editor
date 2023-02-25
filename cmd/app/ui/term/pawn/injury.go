package pawn

import (
	"log"
	"strconv"

	"github.com/cruffinoni/rimworld-editor/algorithm"
	"github.com/cruffinoni/rimworld-editor/generated"
	"github.com/cruffinoni/rimworld-editor/xml/types/iterator"
)

type Injury struct {
	RegisteredPawns PawnsRegisterer
}

func NewInjury(rp PawnsRegisterer) *Injury {
	return &Injury{
		RegisteredPawns: rp,
	}
}

func (i *Injury) RemoveInjuries(args []string) error {
	if len(args) == 0 {
		log.Printf("At least, one pawn name is required (e.g. Thing_Human1046")
		return nil
	}
	p, ok := i.RegisteredPawns[args[0]]
	if !ok {
		log.Printf("Pawn not found.")
		return nil
	}
	p.HealthTracker.HediffSet.Hediffs.Reset()
	log.Printf("All injury for %s has been removed", args[0])
	return nil
}

func (i *Injury) Remove(args []string) error {
	if len(args) == 0 {
		log.Printf("One pawn name is required (e.g. Thing_Human1046) and one injury id is required")
		return nil
	}
	p, ok := i.RegisteredPawns[args[0]]
	if !ok {
		log.Printf("Pawn not found.")
		return nil
	}
	injury, okFind := algorithm.FindInSliceIf[*generated.Hediffs](p.HealthTracker.HediffSet.Hediffs, func(hediffs *generated.Hediffs) bool {
		return strconv.FormatInt(hediffs.LoadID, 10) == args[1]
	})
	if !okFind {
		log.Printf("Injury '%s' not found.", args[1])
		return nil
	}
	log.Printf("Injury %+v found", injury)
	return nil
}

func (i *Injury) List(args []string) error {
	if len(args) == 0 {
		log.Printf("At least, one pawn name is required (e.g. Thing_Human1046")
		return nil
	}
	p, ok := i.RegisteredPawns[args[0]]
	if !ok {
		log.Printf("Pawn not found.")
		return nil
	}
	log.Printf("%s's injuries:", args[0])
	if p.HealthTracker.HediffSet.Hediffs.Capacity() == 0 {
		return nil
	}
	for i := iterator.NewSliceIterator[*generated.Hediffs](p.HealthTracker.HediffSet.Hediffs); i.HasNext(); i = i.Next() {
		v := i.Value()
		log.Printf("%v (ID: %v) with a sevirity of %f", v.Def, v.LoadID, v.Severity*100.0)
		if v.CombatLogText != "null" {
			log.Printf("Logs: %s", v.CombatLogText)
		}
	}
	return nil
}
