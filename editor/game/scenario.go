package game

import "github.com/cruffinoni/rimworld-editor/xml"

type Scenario struct {
	Name          xml.EmbeddedPrimaryType[string] `xml:"name"`
	Summary       string                          `xml:"summary"`
	Description   string                          `xml:"description"`
	PlayerFaction xml.Elements                    `xml:"playerFaction"`
	Parts         xml.Elements                    `xml:"parts"`
}

func (s *Scenario) SetAttributes(_ xml.Attributes) {
	// No attributes need to be set.
}

func (s *Scenario) GetAttributes() xml.Attributes {
	return nil
}

func (s *Scenario) Assign(e *xml.Element) error {
	return nil
}

func (s *Scenario) GetPath() string {
	return ""
}
