package faction

import "github.com/cruffinoni/rimworld-editor/generated"

type Create struct {
	SG *generated.Savegame
}

func (a *Create) Handle(args []string) error {
	return nil
}
