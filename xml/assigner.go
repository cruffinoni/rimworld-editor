package xml

// AttributeAssigner is an interface for objects that can be assigned XML attributes.
type AttributeAssigner interface {
	SetAttributes(attributes Attributes)
	GetAttributes() Attributes
}

type Assigner interface {
	Assign(e *Element) error
	GetPath() string
	AttributeAssigner
}
