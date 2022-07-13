package editor

import (
	"github.com/cruffinoni/rimworld-editor/editor/game"
	"github.com/cruffinoni/rimworld-editor/editor/meta"
	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/saver"
	"log"
	"reflect"
)

type Savegame struct {
	Meta *meta.Meta `xml:"meta"`
	Game *game.Game `xml:"game"`
}

func (s *Savegame) SetAttributes(_ xml.Attributes) {
	// No attributes need to be set.
}

func (s *Savegame) GetAttributes() xml.Attributes {
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

func saveValue(val any, b *saver.Buffer) error {
	t := reflect.TypeOf(val).Elem()
	v := reflect.ValueOf(val)
	log.Printf("t name: '%s'", t.Name())

	var attr xml.Attributes
	if attributeAssigner, ok := val.(xml.AttributeAssigner); ok {
		attr = attributeAssigner.GetAttributes()
	}
	b.OpenTag(t.Name(), attr)
	switch t.Kind() {
	case reflect.Struct:
		if s, ok := v.Interface().(saver.Transformer); ok {
			if r, err := s.TransformToXML(); err != nil {
				return err
			} else {
				b.Write(r)
			}
		}
		for i := 0; i < v.NumField(); i++ {
			if err := saveValue(v.Field(i).Interface(), b); err != nil {
				return err
			}
		}
	}
	b.CloseTag(t.Name())
	return nil
}

func (s *Savegame) SaveXML(path string) error {
	b := saver.NewBuffer()
	if err := saveValue(s, b); err != nil {
		return err
	}
	return b.ToFile(path)
}
