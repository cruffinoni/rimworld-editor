package faction

import (
	"github.com/cruffinoni/rimworld-editor/generated"
	"github.com/cruffinoni/rimworld-editor/internal/xml/iterator"
	"github.com/cruffinoni/rimworld-editor/pkg/logging"
)

type List struct {
	sg     *generated.Savegame
	r      Registerer
	logger logging.Logger
}

func NewList(logger logging.Logger, sg *generated.Savegame, r Registerer) *List {
	return &List{r: r, sg: sg, logger: logger}
}

func (l *List) ListAllFactions() {
	allFac := map[string]*generated.AllFactions{}
	l.logger.Info("Summary of all factions")
	ite := iterator.NewSliceIterator[*generated.AllFactions](l.sg.Game.World.FactionManager.AllFactions)
	for i := ite; i.HasNext(); i = i.Next() {
		v := i.Value()
		allFac[GetFactionID(v.LoadId)] = v
	}
	for _, f := range allFac {
		PrintFactionInformation(l.logger, l.r, f, true)
		l.logger.Debug("Faction summary end")
	}
	return
}
