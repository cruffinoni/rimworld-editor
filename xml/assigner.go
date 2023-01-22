// Package xml provides types and functions for working with XML elements and attributes.
package xml

import "github.com/cruffinoni/rimworld-editor/xml/attributes"

// AttributeAssigner is an interface for types that can set and get XML attributes.
type AttributeAssigner interface {
	// SetAttributes sets the XML attributes for the implementing type.
	SetAttributes(attributes attributes.Attributes)
	// GetAttributes gets the XML attributes for the implementing type.
	GetAttributes() attributes.Attributes
}

// Assigner is an interface for types that can assign values to XML elements and have a path and attributes.
type Assigner interface {
	// Assign assigns values to the XML element.
	Assign(e *Element) error
	GetPath() string
	// AttributeAssigner is an interface for types that can set and get XML attributes.
	AttributeAssigner
}
