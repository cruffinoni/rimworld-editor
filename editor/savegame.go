package editor

import (
	"fmt"
	"github.com/cruffinoni/rimworld-editor/editor/game"
	"github.com/cruffinoni/rimworld-editor/editor/meta"
	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/attributes"
	"github.com/cruffinoni/rimworld-editor/xml/saver"
	"log"
	"reflect"
	"strconv"
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

func saveValue(val any, b *saver.Buffer, tag string) error {
	if val == nil {
		return nil
	}
	t := reflect.TypeOf(val)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	v := reflect.ValueOf(val)
	if v.IsZero() {
		return nil
	}
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	var attr attributes.Attributes
	if attributeAssigner, ok := val.(xml.AttributeAssigner); ok {
		attr = attributeAssigner.GetAttributes()
	}

	if _, ok := val.(xml.Empty); ok {
		b.WriteEmptyTag(tag, attr)
		return nil
	}
	b.OpenTag(tag, attr)
	//b.AddFlag(saver.FlagWriteOpenTag | saver.FlagWriteCloseTag)
	kind := v.Kind()
	//log.Printf("kind: %v", kind)
	//if helper.IsReflectPrimaryType(kind) {
	//	b.Write(v.Bytes())
	//}
	switch kind {
	//case reflect.Slice:
	//	b.Write(v.Bytes())
	case reflect.String:
		b.Write([]byte(v.String()))
	case reflect.Int:
		b.Write([]byte(strconv.FormatInt(v.Int(), 10)))
	case reflect.Float64:
		b.Write([]byte(strconv.FormatFloat(v.Float(), 'f', -1, 64)))
	case reflect.Struct:
		vInterface := v.Interface()
		if vInterface == nil {
			return nil
		}
		switch tInter := vInterface.(type) {
		case fmt.Stringer:
			b.Write([]byte(tInter.String()))
		case saver.Transformer:
			if r, err := tInter.TransformToXML(); err != nil {
				return err
			} else {
				log.Print("writing from transformer")
				b.WriteWithIndent(r)
			}
		}
		for i := 0; i < v.NumField(); i++ {
			f := t.Field(i)
			vTag, ok := f.Tag.Lookup("xml")
			if !ok {
				continue
			}
			vField := v.Field(i)
			// Ignore zero value fields.
			if vField.IsZero() {
				continue
			}
			b.Write([]byte(">\n"))
			if err := saveValue(vField.Interface(), b, vTag); err != nil {
				return err
			}
		}
		b.Write([]byte(">>\n"))
		b.CloseTagWithIndent(tag)
		return nil
	}
	b.CloseTag(tag)
	return nil
}

func (s *Savegame) SaveXML(path string) error {
	b := saver.NewBuffer()
	if err := saveValue(s, b, "savegame"); err != nil {
		return err
	}
	return b.ToFile(path)
}
