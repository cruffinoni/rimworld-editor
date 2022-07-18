package game

import (
	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/attributes"
	"github.com/cruffinoni/rimworld-editor/xml/types"
	"github.com/cruffinoni/rimworld-editor/xml/types/primary"
)

type Parts struct {
	attr            attributes.Attributes
	Def             string `xml:"def"`
	PawnCount       string `xml:"pawnCount"`
	PawnChoiceCount string `xml:"pawnChoiceCount"`
	Method          string `xml:"method"`
	Chance          string `xml:"chance"`
	Context         string `xml:"context"`
	Text            string `xml:"text"`
	CloseSound      string `xml:"closeSound"`
}

func (p *Parts) Assign(e *xml.Element) error {
	return nil
}

func (p *Parts) GetPath() string {
	return ""
}

func (p *Parts) SetAttributes(attributes attributes.Attributes) {
	p.attr = attributes
}

func (p *Parts) GetAttributes() attributes.Attributes {
	return p.attr
}

type Scenario struct {
	Name          *primary.EmbeddedType[string] `xml:"name"`
	Summary       string                        `xml:"summary"`
	Description   string                        `xml:"description"`
	PlayerFaction *xml.Elements                 `xml:"playerFaction"`
	Parts         *types.Slice[*Parts]          `xml:"parts"`
}

func (s *Scenario) SetAttributes(_ attributes.Attributes) {
	// No attributes need to be set.
}

func (s *Scenario) GetAttributes() attributes.Attributes {
	return nil
}

func (s *Scenario) Assign(e *xml.Element) error {
	return nil
}

func (s *Scenario) GetPath() string {
	return ""
}
