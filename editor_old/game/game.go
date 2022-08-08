package game

import (
	"github.com/cruffinoni/rimworld-editor/editor_old/game/world"
	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/attributes"
)

type Game struct {
	CurrentMapIndex int64            `xml:"currentMapIndex"`
	Info            string           `xml:"info"`
	Rules           string           `xml:"rules"`
	Scenario        *Scenario        `xml:"scenario"`
	TickManager     *TickManager     `xml:"tickManager"`
	ResearchManager *ResearchManager `xml:"researchManager"`
	World           *world.World     `xml:"world"`
}

func (g *Game) Assign(e *xml.Element) error {
	return nil
}

func (g *Game) GetPath() string {
	return ""
}

func (g *Game) SetAttributes(_ attributes.Attributes) {
	// No attributes need to be set.
}

func (g *Game) GetAttributes() attributes.Attributes {
	return nil
}

//
//func (g *Game) GetCurrentMapIndex() int64 {
//	return g.CurrentMapIndex
//}
//
//func (g *Game) GetInfo() string {
//	return g.Info
//}
