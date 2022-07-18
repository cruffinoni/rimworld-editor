package primary

import (
	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/attributes"
)

type Empty struct {
	name string
	attr attributes.Attributes
}

func (e *Empty) Assign(element *xml.Element) error {
	e.name = element.GetName()
	return nil
}

func (e *Empty) GetPath() string {
	return ""
}

func (e *Empty) SetAttributes(attributes attributes.Attributes) {
	e.attr = attributes
}

func (e *Empty) GetAttributes() attributes.Attributes {
	return e.attr
}

func (e *Empty) String() string {
	return e.name
}
