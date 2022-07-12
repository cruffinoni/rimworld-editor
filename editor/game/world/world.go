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
