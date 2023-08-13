package term

import (
	"github.com/cruffinoni/rimworld-editor/cmd/app/ui/term/printer"
	"github.com/cruffinoni/rimworld-editor/generated"
	"github.com/cruffinoni/rimworld-editor/xml/types/iterator"
)

func (c *Console) growAllPlants(percent float64) {
	growth := percent / 100.0
	if growth > 1.0 {
		growth = 1.0
	} else if growth < 0.0 {
		growth = 0.1
	}

	count := 0
	ite := iterator.NewSliceIterator[*generated.Thing](c.save.Game.Maps.GetFromIndex(0).Things)
	for i := ite; i.HasNext(); i = i.Next() {
		t := i.Value()
		if t.Attr.Get("Class") != "Plant" {
			continue
		}
		t.Growth = growth
		count++
	}
	printer.PrintSf("{-BOLD,GREEN}%d plants{-RESET} edited to %.2f%% growth percentage", count, growth*100.0)
	return
}
