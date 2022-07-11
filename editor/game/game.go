package game

import "github.com/cruffinoni/rimworld-editor/xml"

type Game struct {
	CurrentMapIndex int              `xml:"currentMapIndex"`
	Info            string           `xml:"info"`
	Rules           string           `xml:"rules"`
	Scenario        *Scenario        `xml:"scenario"`
	TickManager     *TickManager     `xml:"tickManager"`
	ResearchManager *ResearchManager `xml:"researchManager"`
}

func (g *Game) Assign(e *xml.Element) error {
	return nil
}

func (g *Game) GetPath() string {
	return ""
}

func (g *Game) GetCurrentMapIndex() int {
	return g.CurrentMapIndex
}

func (g *Game) GetInfo() string {
	return g.Info
}
