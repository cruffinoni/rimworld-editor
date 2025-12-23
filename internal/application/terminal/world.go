package terminal

import (
	"github.com/cruffinoni/rimworld-editor/generated"
	"github.com/cruffinoni/rimworld-editor/internal/xml/iterator"
	"github.com/cruffinoni/rimworld-editor/pkg/logging"
)

func (c *Console) growAllPlants(percent float64) {
	growth := percent / 100.0
	if growth > 1.0 {
		growth = 1.0
	} else if growth < 0.0 {
		growth = 0.1
	}

	count := 0
	ite := iterator.NewSliceIterator[*generated.Thing](c.save.Game.Maps.At(0).Things)
	for i := ite; i.HasNext(); i = i.Next() {
		t := i.Value()
		if t.Attr.Get("Class") != "Plant" {
			continue
		}
		t.Growth = growth
		count++
	}
	c.logger.WithFields(logging.Fields{
		"plants": count,
		"growth": growth * 100.0,
	}).Info("Updated plant growth percentage")
	return
}
