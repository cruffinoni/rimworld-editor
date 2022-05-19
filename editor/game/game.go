package game

import "github.com/cruffinoni/rimworld-editor/xml"

type Scenario struct {
	Name          string         `xml:"name"`
	Summary       string         `xml:"summary"`
	Description   string         `xml:"description"`
	PlayerFaction []*xml.Element `xml:"playerFaction"`
	Parts         []*xml.Element `xml:"parts"`
}

type Game struct {
	CurrentMapIndex int      `xml:"currentMapIndex"`
	Info            string   `xml:"info"`
	Rules           string   `xml:"rules"`
	Scenario        Scenario `xml:"scenario"`
}
