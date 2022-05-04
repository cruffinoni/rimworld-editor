package game

import "github.com/cruffinoni/rimworld-editor/xml"

type Scenario struct {
	Name          string     `xml:"name"`
	Summary       string     `xml:"summary"`
	Description   string     `xml:"description"`
	PlayerFaction xml.Map    `xml:"playerFaction"`
	Parts         xml.Nested `xml:"parts"`
}

type Game struct {
	CurrentMapIndex int `xml:"currentMapIndex"`
	//Info            xml.Map  `xml:"info"`
	//Rules           xml.Map  `xml:"rules"`
	//Scenario        Scenario `xml:"scenario"`
}
