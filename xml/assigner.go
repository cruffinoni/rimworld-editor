package xml

import "github.com/cruffinoni/rimworld-editor/xml/attributes"

// AttributeAssigner is an interface for objects that can be assigned XML attributes.
type AttributeAssigner interface {
	SetAttributes(attributes attributes.Attributes)
	GetAttributes() attributes.Attributes
}

type Assigner interface {
	Assign(e *Element) error
	GetPath() string
	AttributeAssigner
}
