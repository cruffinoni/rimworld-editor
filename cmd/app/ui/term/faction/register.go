package faction

import (
	"strconv"

	"github.com/cruffinoni/rimworld-editor/generated"
	"github.com/cruffinoni/rimworld-editor/internal/xml/types/iterator"
	"github.com/cruffinoni/rimworld-editor/pkg/logging"
)

type Registerer map[string]*generated.AllFactions

func RegisterFactions(logger logging.Logger, sg *generated.Savegame) Registerer {
	allFac := Registerer{}
	ite := iterator.NewSliceIterator[*generated.AllFactions](sg.Game.World.FactionManager.AllFactions)
	for i := ite; i.HasNext(); i = i.Next() {
		v := i.Value()
		allFac[GetFactionID(v.LoadId)] = v
	}
	return allFac
}

func GetFactionID(loadID int64) string {
	return "Faction_" + strconv.FormatInt(loadID, 10)
}

func IsPlayerFaction(f *generated.AllFactions) bool {
	if f == nil {
		return false
	}
	return f.Def == "PlayerColony"
}

func PrintFactionInformation(logger logging.Logger, rf Registerer, f *generated.AllFactions, withRelations bool) {
	id := GetFactionID(f.LoadId)
	logger.WithFields(logging.Fields{
		"faction": f.Def,
		"name":    f.Name,
		"id":      id,
	}).Info("Faction")
	if f.Def == "PlayerColony" {
		logger.WithField("id", id).Debug("Player faction")
	} else {
		logger.WithField("id", id).Debug("AI-controlled faction")
	}
	if f.Leader == "null" {
		logger.WithField("id", id).Debug("Faction has no leader")
	} else {
		logger.WithFields(logging.Fields{
			"id":     id,
			"leader": f.Leader,
		}).Debug("Faction leader")
	}
	if withRelations {
		logger.WithField("id", id).Debug("Relations")
		for i := iterator.NewSliceIterator[*generated.Relations](f.Relations); i.HasNext(); i = i.Next() {
			r := i.Value()
			relationLogger := logger.WithFields(logging.Fields{
				"source":   f.Def,
				"sourceID": id,
				"target":   rf[r.Other].Def,
				"targetID": GetFactionID(rf[r.Other].LoadId),
				"goodwill": r.Goodwill,
			})
			if r.Goodwill > 75 {
				relationLogger.Info("Relation: allied")
			} else if r.Goodwill < -50 {
				relationLogger.Warn("Relation: hostile")
			} else {
				relationLogger.Info("Relation: neutral")
			}
		}
	}
}
