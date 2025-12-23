package multiple

import (
	"github.com/cruffinoni/rimworld-editor/internal/xml"
	"github.com/cruffinoni/rimworld-editor/internal/xml/attributes"
	"github.com/cruffinoni/rimworld-editor/internal/xml/saver"
	"github.com/cruffinoni/rimworld-editor/pkg/logging"
)

type Data struct {
	Element *xml.Element
	Next    *Data
}

// Type is a linked list of xml.Element that represents multiple types for the same field
// Actually, it's only used in the types.Slice type
type Type struct {
	last  *Data
	first *Data

	logger logging.Logger
}

func (t *Type) Assign(e *xml.Element) error {
	t.logger.WithField("tag", e.GetName()).Debug("multiple.Type.Assign called")
	if t.last == nil {
		t.last = &Data{
			Element: e,
		}
		t.first = t.last
		return nil
	}
	t.last.Next = &Data{
		Element: e,
	}
	t.last = t.last.Next
	return nil
}

func (t *Type) GetPath() string {
	return ""
}

func (t *Type) SetAttributes(_ attributes.Attributes) {
	//printer.Debugf("SetAttributes called on multiple.Type: %v", attributes)
}

func (t *Type) GetAttributes() attributes.Attributes {
	return nil
}

func (t *Type) TransformToXML(buffer *saver.Buffer) error {
	if t.first.Element == nil {
		return nil
	}
	// We are in a list, so don't write twice the same tag
	if t.first.Element.GetName() == "li" {
		buffer.WriteString(t.first.Element.Data.GetString())
		t.first = t.first.Next
		return nil
	} else {
		buffer.WriteString(t.first.Element.ToXML(0))
	}
	t.first = t.first.Next
	return nil
}

func (t *Type) GetXMLTag() []byte {
	t.logger.Debug("multiple.Type.GetXMLTag called")
	return []byte("")
}

func (t *Type) SetLogger(logger logging.Logger) {
	t.logger = logger
}
