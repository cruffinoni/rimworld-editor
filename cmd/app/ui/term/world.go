package term

import (
	"log"
	"strconv"

	"github.com/cruffinoni/rimworld-editor/generated"
	"github.com/cruffinoni/rimworld-editor/xml/types/iterator"
)

func (c *Console) growAllPlants(args []string) error {
	growth := 1.0
	if len(args) > 0 {
		var err error
		if growth, err = strconv.ParseFloat(args[0], 64); err != nil {
			return err
		}
	}
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
		t.Health = 1.0
		count++
	}
	log.Printf("%d plants edited to %.2f%% growth percentage", count, growth*100.0)
	return nil
}
