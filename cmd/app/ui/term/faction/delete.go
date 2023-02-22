package faction

import "github.com/cruffinoni/rimworld-editor/generated"

type Delete struct {
	SG *generated.Savegame
}

func (d *Delete) Handle(args []string) error {
	return nil
}
