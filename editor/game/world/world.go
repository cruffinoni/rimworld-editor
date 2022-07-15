package world

import (
	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/attributes"
)

type World struct {
	FactionManager *FactionManager `xml:"factionManager"`
	Info           *Info           `xml:"info"`
}

func (w *World) Assign(e *xml.Element) error {
	return nil
}

func (w *World) GetPath() string {
	return ""
}

func (w *World) SetAttributes(_ attributes.Attributes) {
	// No attributes need to be set.
}

func (w *World) GetAttributes() attributes.Attributes {
	return nil
}
