package faction

import (
	"strings"
	"sync"

	"github.com/cruffinoni/rimworld-editor/generated"
	"github.com/cruffinoni/rimworld-editor/pkg/logging"
)

type Delete struct {
	sg     *generated.Savegame
	reg    Registerer
	logger logging.Logger
}

func NewDelete(logger logging.Logger, sg *generated.Savegame, reg Registerer) *Delete {
	return &Delete{
		sg:     sg,
		reg:    reg,
		logger: logger,
	}
}

func (d *Delete) Handle(args []string) error {
	if len(args) == 0 {
		d.logger.Error("At least one argument is required: faction id (e.g. Faction_12)")
		return nil
	}
	// TODO: Do it automatically from the command handler
	args[0] = strings.Trim(args[0], " ")
	if _, ok := d.reg[args[0]]; !ok {
		d.logger.WithField("id", args[0]).Error("Faction not found")
		return nil
	}
	wg := &sync.WaitGroup{}
	wg.Add(3)

	go func() {
		d.logger.WithField("id", args[0]).Info("Deleting faction from registry")
		nbOfFac := d.sg.Game.World.FactionManager.AllFactions.Capacity()
		for i := 0; i < nbOfFac; i++ {
			if GetFactionID(d.sg.Game.World.FactionManager.AllFactions.At(i).LoadId) == args[0] {
				d.sg.Game.World.FactionManager.AllFactions.Remove(i)
				break
			}
		}
		delete(d.reg, args[0])
		d.logger.WithField("id", args[0]).Info("Faction deleted")
		wg.Done()
	}()

	go func() {
		count := 0
		d.logger.WithField("id", args[0]).Info("Deleting faction-owned objects")
		for i, j := 0, d.sg.Game.Maps.Capacity(); i < j; i++ {
			for k, l := 0, d.sg.Game.Maps.At(i).Things.Capacity(); k < l; k++ {
				t := d.sg.Game.Maps.At(i).Things.At(k)
				if t.Faction == args[0] {
					d.sg.Game.Maps.At(i).Things.Remove(k)
					count++
				}
			}
		}
		d.logger.WithField("removed", count).Info("Faction-owned objects removed")
		wg.Done()
	}()

	go func() {
		count := 0
		d.logger.WithField("id", args[0]).Info("Removing faction archives")
		for i, j := 0, d.sg.Game.History.Archive.Archivables.Capacity(); i < j; i++ {
			a := d.sg.Game.History.Archive.Archivables.At(i)
			if a.RelatedFaction == args[0] {
				d.sg.Game.History.Archive.Archivables.Remove(i)
				count++
			}
		}
		d.logger.WithField("removed", count).Info("Faction archives removed")
		wg.Done()
	}()
	wg.Wait()
	return nil
}
