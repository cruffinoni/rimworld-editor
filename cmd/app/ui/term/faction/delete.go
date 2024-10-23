package faction

import (
	"strings"
	"sync"

	"github.com/cruffinoni/printer"

	"github.com/cruffinoni/rimworld-editor/generated"
)

type Delete struct {
	sg  *generated.Savegame
	reg Registerer
}

func NewDelete(sg *generated.Savegame, reg Registerer) *Delete {
	return &Delete{
		sg:  sg,
		reg: reg,
	}
}

func (d *Delete) Handle(args []string) error {
	if len(args) == 0 {
		printer.PrintErrorS("at least one argument is required: faction id (e.g. Faction_12)")
		return nil
	}
	// TODO: Do it automatically from the command handler
	args[0] = strings.Trim(args[0], " ")
	if _, ok := d.reg[args[0]]; !ok {
		printer.PrintErrorSf("Faction '%s' not found", args[0])
		return nil
	}
	wg := &sync.WaitGroup{}
	wg.Add(3)

	go func() {
		printer.Debugf("Deleting faction {{{-BOLD,F_RED}}}'%s'{{{-RESET}}} from the available factions", args[0])
		nbOfFac := d.sg.Game.World.FactionManager.AllFactions.Capacity()
		for i := 0; i < nbOfFac; i++ {
			if GetFactionID(d.sg.Game.World.FactionManager.AllFactions.At(i).LoadId) == args[0] {
				d.sg.Game.World.FactionManager.AllFactions.Remove(i)
				break
			}
		}
		delete(d.reg, args[0])
		printer.Debugf("Faction {{{-BOLD,F_RED}}}'%s'{{{-RESET}}} deleted", args[0])
		wg.Done()
	}()

	go func() {
		count := 0
		printer.Debugf("Deleting things (animals, objects, etc.) that belong to {{{-BOLD,F_RED}}}'%s'{{{-RESET}}}", args[0])
		for i, j := 0, d.sg.Game.Maps.Capacity(); i < j; i++ {
			for k, l := 0, d.sg.Game.Maps.At(i).Things.Capacity(); k < l; k++ {
				t := d.sg.Game.Maps.At(i).Things.At(k)
				if t.Faction == args[0] {
					d.sg.Game.Maps.At(i).Things.Remove(k)
					count++
				}
			}
		}
		printer.Debugf("{{{-BOLD,F_BLUE}}}%d{{{-RESET}}} objects removed", count)
		wg.Done()
	}()

	go func() {
		count := 0
		printer.Debugf("Removing archives that are related to the faction {{{-BOLD,F_RED}}}'%s'{{{-RESET}}}", args[0])
		for i, j := 0, d.sg.Game.History.Archive.Archivables.Capacity(); i < j; i++ {
			a := d.sg.Game.History.Archive.Archivables.At(i)
			if a.RelatedFaction == args[0] {
				d.sg.Game.History.Archive.Archivables.Remove(i)
				count++
			}
		}
		printer.Debugf("{{{-BOLD,F_BLUE}}}%d{{{-RESET}}} archives removed", count)
		wg.Done()
	}()
	wg.Wait()
	return nil
}
