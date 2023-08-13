package primary

import (
	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/attributes"
)

// Empty is a type that represent an XML tag which is empty but can contain
// attributes.Attributes.
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

func (e *Empty) SetAttributes(attr attributes.Attributes) {
	e.attr = attr
}

func (e *Empty) GetAttributes() attributes.Attributes {
	return e.attr
}

func (e *Empty) String() string {
	return e.name
}
