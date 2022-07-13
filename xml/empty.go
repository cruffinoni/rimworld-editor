package xml

import "log"

type Empty struct {
	name string
	attr Attributes
}

func (e *Empty) Assign(element *Element) error {
	e.name = element.GetName()
	log.Printf("Assigning %s to %s", element.GetName(), e.name)
	return nil
}

func (e *Empty) GetPath() string {
	return ""
}

func (e *Empty) SetAttributes(attributes Attributes) {
	e.attr = attributes
}

func (e *Empty) GetAttributes() Attributes {
	return e.attr
}

func (e *Empty) String() string {
	return e.name
}
