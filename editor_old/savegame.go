package editor_old

import (
	"github.com/cruffinoni/rimworld-editor/editor_old/game"
	"github.com/cruffinoni/rimworld-editor/editor_old/meta"
	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/attributes"
	"github.com/cruffinoni/rimworld-editor/xml/saver"
	"github.com/cruffinoni/rimworld-editor/xml/saver/file"
)

type Savegame struct {
	Meta *meta.Meta `xml:"meta"`
	Game *game.Game `xml:"game"`
}

func (s *Savegame) SetAttributes(_ attributes.Attributes) {
	// No attributes need to be set.
}

func (s *Savegame) GetAttributes() attributes.Attributes {
	return nil
}

func (s *Savegame) Assign(_ *xml.Element) error {
	return nil
}

func (s *Savegame) GetPath() string {
	return ""
}

func (s *Savegame) GetMeta() *meta.Meta {
	return s.Meta
}

func (s *Savegame) SaveXML(path string) error {
	b := saver.NewBuffer()
	if err := file.Save(s, b, "savegame"); err != nil {
		return err
	}
	return b.ToFile(path)
}
