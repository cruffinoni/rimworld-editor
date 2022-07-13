package world

import "github.com/cruffinoni/rimworld-editor/xml"

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

func (w *World) SetAttributes(_ xml.Attributes) {
	// No attributes need to be set.
}

func (w *World) GetAttributes() xml.Attributes {
	return nil
}
