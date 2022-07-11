package editor

import (
	"github.com/cruffinoni/rimworld-editor/editor/game"
	"github.com/cruffinoni/rimworld-editor/editor/meta"
	"github.com/cruffinoni/rimworld-editor/xml"
)

type Savegame struct {
	Meta *meta.Meta `xml:"meta"`
	Game *game.Game `xml:"game" `
}

func (s *Savegame) Assign(_ *xml.Element) error {
	return nil
}

func (s *Savegame) GetPath() string {
	return ""
}
