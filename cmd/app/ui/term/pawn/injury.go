package pawn

import (
	"strconv"

	"github.com/cruffinoni/rimworld-editor/generated"
	"github.com/cruffinoni/rimworld-editor/internal/xml/types/iterator"
	"github.com/cruffinoni/rimworld-editor/pkg/logging"
)

type Injury struct {
	RegisteredPawns PawnsRegisterer
	logger          logging.Logger
}

func NewInjury(logger logging.Logger, rp PawnsRegisterer) *Injury {
	return &Injury{
		RegisteredPawns: rp,
		logger:          logger,
	}
}

func (i *Injury) RemoveInjuries(pawnID string) {
	p, ok := i.RegisteredPawns[pawnID]
	if !ok {
		i.logger.WithField("pawn_id", pawnID).Error("Pawn not found")
		return
	}
	p.HealthTracker.HediffSet.Hediffs.Reset()
	i.logger.WithField("pawn_id", pawnID).Info("All injuries removed")
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
		i.logger.Error("At least one injury id is required")
		return
	}
	p, ok := i.RegisteredPawns[pawnID]
	if !ok {
		i.logger.WithField("pawn_id", pawnID).Error("Pawn not found")
		return
	}
	totalInjuries := p.HealthTracker.HediffSet.Hediffs.Capacity()
	for idx := 0; idx < totalInjuries; idx++ {
		for _, j := range injuryIds {
			currentInjury := p.HealthTracker.HediffSet.Hediffs.At(idx)
			if j != strconv.FormatInt(currentInjury.LoadId, 10) {
				continue
			}
			i.logger.WithFields(logging.Fields{
				"pawn_id":  pawnID,
				"injuryID": currentInjury.LoadId,
			}).Info("Injury removed")
			p.HealthTracker.HediffSet.Hediffs.Remove(idx)
			injuryIds = append(injuryIds[:idx], injuryIds[idx+1:]...)
		}
	}
}

func (i *Injury) List(pawnID string) {
	p, ok := i.RegisteredPawns[pawnID]
	if !ok {
		i.logger.WithField("pawn_id", pawnID).Error("Pawn not found")
		return
	}
	i.logger.WithField("pawn_id", pawnID).Info("Listing injuries")
	if p.HealthTracker.HediffSet.Hediffs.Capacity() == 0 {
		return
	}
	for it := iterator.NewSliceIterator[*generated.Hediffs](p.HealthTracker.HediffSet.Hediffs); it.HasNext(); it = it.Next() {
		v := it.Value()
		i.logger.WithFields(logging.Fields{
			"def":      v.Def,
			"loadID":   v.LoadId,
			"severity": v.Severity * 100.0,
		}).Info("Injury")
		// if v.CombatLogText != "null" && v.CombatLogText != "" {
		// 	printer.Debugf("Logs: %s", v.CombatLogText)
		// }
	}
}
