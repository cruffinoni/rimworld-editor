package pawn

import (
	"github.com/cruffinoni/rimworld-editor/generated"
	"github.com/cruffinoni/rimworld-editor/internal/application/term/faction"
	"github.com/cruffinoni/rimworld-editor/pkg/logging"
)

type List struct {
	sg     *generated.Savegame
	rp     PawnsRegisterer
	rf     faction.Registerer
	logger logging.Logger
}

func NewList(logger logging.Logger, sg *generated.Savegame, rp PawnsRegisterer, rf faction.Registerer) *List {
	return &List{
		sg:     sg,
		rp:     rp,
		rf:     rf,
		logger: logger,
	}
}

func (l *List) ListAllPawns() {
	for k, v := range l.rp {
		l.logger.WithField("pawn_id", k).Debug("Pawn registered")
		if v.Name == nil {
			l.logger.WithField("pawn_id", k).Debug("Pawn name is nil")
			continue
		}
		l.logger.WithFields(logging.Fields{
			"pawn_id":  k,
			"fullName": getPawnFullName(v),
		}).Info("Pawn name")
		if fac, ok := l.rf[v.Faction]; ok {
			faction.PrintFactionInformation(l.logger, l.rf, fac, false)
		}
		l.logger.Debug("Pawn summary end")
	}
}
