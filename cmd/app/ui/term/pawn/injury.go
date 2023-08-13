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

func (i *Injury) RemoveInjuries(pawnID string) {
	p, ok := i.RegisteredPawns[pawnID]
	if !ok {
		printer.PrintErrorSf("Pawn %s not found", pawnID)
		return
	}
	p.HealthTracker.HediffSet.Hediffs.Reset()
	printer.PrintSf("All injuries has been removed for {-BOLD,GREEN}%s", pawnID)
	return
}

func (i *Injury) Heal(pawnID string) {
	p, ok := i.RegisteredPawns[pawnID]
	if !ok {
		printer.PrintErrorSf("Pawn %s not found", pawnID)
		return
	}
	algorithm.SliceForeach[*generated.Hediffs](p.HealthTracker.HediffSet.Hediffs, func(h *generated.Hediffs) {

	})
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
			currentInjury := p.HealthTracker.HediffSet.Hediffs.GetFromIndex(i)
			if j != strconv.FormatInt(currentInjury.LoadID, 10) {
				continue
			}
			printer.PrintSf("Injury {-RED,BOLD}%s {-RESET}removed", currentInjury.LoadID)
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
	printer.PrintSf("{-BOLD,MAGENTA}%s{-RESET}'s injuries:", pawnID)
	if p.HealthTracker.HediffSet.Hediffs.Capacity() == 0 {
		return
	}
	for i := iterator.NewSliceIterator[*generated.Hediffs](p.HealthTracker.HediffSet.Hediffs); i.HasNext(); i = i.Next() {
		v := i.Value()
		printer.PrintSf("%v ({-BOLD}load ID: %v{-RESET}) with a severity of {-RED}%f", v.Def, v.LoadID, v.Severity*100.0)
		// if v.CombatLogText != "null" && v.CombatLogText != "" {
		// 	printer.PrintSf("Logs: %s", v.CombatLogText)
		// }
	}
}
