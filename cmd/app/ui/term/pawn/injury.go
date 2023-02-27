package pawn

import (
	"strconv"

	"github.com/cruffinoni/rimworld-editor/algorithm"
	"github.com/cruffinoni/rimworld-editor/cmd/app/ui/term/printer"
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
		printer.PrintErrorS("At least, one pawn name is required (e.g. Thing_Human1046")
		return nil
	}
	p, ok := i.RegisteredPawns[args[0]]
	if !ok {
		printer.PrintErrorS("Pawn not found.")
		return nil
	}
	p.HealthTracker.HediffSet.Hediffs.Reset()
	printer.PrintSf("All injury for %s has been removed", args[0])
	return nil
}

func (i *Injury) Remove(args []string) error {
	if len(args) == 0 {
		printer.PrintErrorS("One pawn name is required (e.g. Thing_Human1046) and one injury id is required")
		return nil
	}
	p, ok := i.RegisteredPawns[args[0]]
	if !ok {
		printer.PrintErrorS("Pawn not found.")
		return nil
	}
	loadID, err := strconv.Atoi(args[1])
	if err != nil {
		printer.PrintError(err)
		return nil
	}
	injury, okFind := algorithm.FindInSliceIf[*generated.Hediffs](p.HealthTracker.HediffSet.Hediffs, func(hediffs *generated.Hediffs) bool {
		return hediffs.LoadID == int64(loadID)
	})
	if !okFind {
		printer.PrintSf("Injury '%s' not found.", args[1])
		return nil
	}
	printer.PrintSf("Injury %+v found", injury)
	return nil
}

func (i *Injury) List(args []string) error {
	if len(args) == 0 {
		printer.PrintErrorS("At least, one pawn name is required (e.g. Thing_Human1046)")
		return nil
	}
	p, ok := i.RegisteredPawns[args[0]]
	if !ok {
		printer.PrintErrorS("Pawn not found.")
		return nil
	}
	printer.PrintSf("%s's injuries:", args[0])
	if p.HealthTracker.HediffSet.Hediffs.Capacity() == 0 {
		return nil
	}
	for i := iterator.NewSliceIterator[*generated.Hediffs](p.HealthTracker.HediffSet.Hediffs); i.HasNext(); i = i.Next() {
		v := i.Value()
		printer.PrintSf("%v (ID: %v) with a sevirity of %f", v.Def, v.LoadID, v.Severity*100.0)
		if v.CombatLogText != "null" {
			printer.PrintSf("Logs: %s", v.CombatLogText)
		}
	}
	return nil
}
