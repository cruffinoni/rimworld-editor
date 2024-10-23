package pawn

import (
	"strconv"

	"github.com/cruffinoni/printer"

	"github.com/cruffinoni/rimworld-editor/generated"
	"github.com/cruffinoni/rimworld-editor/internal/xml/types/iterator"
)

type Injury struct {
	RegisteredPawns PawnsRegisterer
}

func NewInjury(rp PawnsRegisterer) *Injury {
	return &Injury{
		RegisteredPawns: rp,
	}
}

func (i *Injury) RemoveInjuries(pawnID string) {
	p, ok := i.RegisteredPawns[pawnID]
	if !ok {
		printer.PrintErrorSf("Pawn %s not found", pawnID)
		return
	}
	p.HealthTracker.HediffSet.Hediffs.Reset()
	printer.Printf("All injuries has been removed for {-BOLD,GREEN}%s", pawnID)
}

func (i *Injury) Heal(pawnID string) {
	//p, ok := i.RegisteredPawns[pawnID]
	//if !ok {
	//	printer.PrintErrorSf("Pawn %s not found", pawnID)
	//	return
	//}
	//algorithm.SliceForeach[*generated.Hediffs](p.HealthTracker.HediffSet.Hediffs, func(h *generated.Hediffs) {
	//
	//})
}

func (i *Injury) Remove(pawnID string, injuryIds []string) {
	if len(injuryIds) == 0 {
		printer.PrintErrorS("At least one injury id is required")
		return
	}
	p, ok := i.RegisteredPawns[pawnID]
	if !ok {
		printer.PrintErrorSf("Pawn %s not found", pawnID)
		return
	}
	totalInjuries := p.HealthTracker.HediffSet.Hediffs.Capacity()
	for i := 0; i < totalInjuries; i++ {
		for _, j := range injuryIds {
			currentInjury := p.HealthTracker.HediffSet.Hediffs.At(i)
			if j != strconv.FormatInt(currentInjury.LoadId, 10) {
				continue
			}
			printer.Printf("Injury {-RED,BOLD}%s {-RESET}removed", currentInjury.LoadId)
			p.HealthTracker.HediffSet.Hediffs.Remove(i)
			injuryIds = append(injuryIds[:i], injuryIds[i+1:]...)
		}
	}
}

func (i *Injury) List(pawnID string) {
	p, ok := i.RegisteredPawns[pawnID]
	if !ok {
		printer.PrintErrorSf("Pawn %s not found", pawnID)
		return
	}
	printer.Printf("{-BOLD,MAGENTA}%s{-RESET}'s injuries:", pawnID)
	if p.HealthTracker.HediffSet.Hediffs.Capacity() == 0 {
		return
	}
	for i := iterator.NewSliceIterator[*generated.Hediffs](p.HealthTracker.HediffSet.Hediffs); i.HasNext(); i = i.Next() {
		v := i.Value()
		printer.Printf("%v ({-BOLD}load ID: %v{-RESET}) with a severity of {-RED}%f", v.Def, v.LoadId, v.Severity*100.0)
		// if v.CombatLogText != "null" && v.CombatLogText != "" {
		// 	printer.Printf("Logs: %s", v.CombatLogText)
		// }
	}
}
