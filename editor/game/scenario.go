package game

import "github.com/cruffinoni/rimworld-editor/xml"

type Scenario struct {
	Name        string       `xml:"name"`
	Summary     string       `xml:"summary"`
	Description *xml.Element `xml:"description"`
	//PlayerFaction []*xml.Element `xml:"playerFaction"`
	//Parts         []*xml.Element `xml:"parts"`
}

func (s *Scenario) Assign(e *xml.Element) error {
	return nil
}

func (s *Scenario) GetPath() string {
	return ""
}
