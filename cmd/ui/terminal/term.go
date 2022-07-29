package terminal

import "github.com/cruffinoni/rimworld-editor/cmd/ui"

type Console struct {
	opt ui.Options
}

func (c *Console) Execute() error {
	return nil
}

func (c *Console) NewMode(options *ui.Options) ui.Mode {
	return &Console{}
}
