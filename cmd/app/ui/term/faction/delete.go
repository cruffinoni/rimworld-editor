package faction

import (
	"strings"
	"sync"

	"github.com/cruffinoni/rimworld-editor/cmd/app/ui/term/printer"
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
		printer.PrintSf("Deleting faction {-BOLD,F_RED}'%s'{-RESET} from the available factions", args[0])
		nbOfFac := d.sg.Game.World.FactionManager.AllFactions.Capacity()
		for i := 0; i < nbOfFac; i++ {
			if GetFactionID(d.sg.Game.World.FactionManager.AllFactions.GetFromIndex(i).LoadID) == args[0] {
				d.sg.Game.World.FactionManager.AllFactions.Remove(i)
				break
			}
		}
		delete(d.reg, args[0])
		printer.PrintSf("Faction {-BOLD,F_RED}'%s'{-RESET} deleted", args[0])
		wg.Done()
	}()

	go func() {
		count := 0
		printer.PrintSf("Deleting things (animals, objects, etc.) that belong to {-BOLD,F_RED}'%s'{-RESET}", args[0])
		for i, j := 0, d.sg.Game.Maps.Capacity(); i < j; i++ {
			for k, l := 0, d.sg.Game.Maps.GetFromIndex(i).Things.Capacity(); k < l; k++ {
				t := d.sg.Game.Maps.GetFromIndex(i).Things.GetFromIndex(k)
				if t.Faction == args[0] {
					d.sg.Game.Maps.GetFromIndex(i).Things.Remove(k)
					count++
				}
			}
		}
		printer.PrintSf("{-BOLD,F_BLUE}%d{-RESET} objects removed", count)
		wg.Done()
	}()

	go func() {
		count := 0
		printer.PrintSf("Removing archives that are related to the faction {-BOLD,F_RED}'%s'{-RESET}", args[0])
		for i, j := 0, d.sg.Game.History.Archive.Archivables.Capacity(); i < j; i++ {
			a := d.sg.Game.History.Archive.Archivables.GetFromIndex(i)
			if a.RelatedFaction == args[0] {
				d.sg.Game.History.Archive.Archivables.Remove(i)
				count++
			}
		}
		printer.PrintSf("{-BOLD,F_BLUE}%d{-RESET} archives removed", count)
		wg.Done()
	}()
	wg.Wait()
	return nil
}
