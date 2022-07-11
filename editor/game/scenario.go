package game

import "github.com/cruffinoni/rimworld-editor/xml"

type Scenario struct {
	Name          string       `xml:"name"`
	Summary       string       `xml:"summary"`
	Description   string       `xml:"description"`
	PlayerFaction xml.Elements `xml:"playerFaction"`
	Parts         xml.Elements `xml:"parts"`
}

func (s *Scenario) Assign(e *xml.Element) error {
	return nil
}

func (s *Scenario) GetPath() string {
	return ""
}
