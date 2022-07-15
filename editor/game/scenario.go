package game

import (
	"fmt"
	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/attributes"
	"github.com/cruffinoni/rimworld-editor/xml/types"
	"github.com/cruffinoni/rimworld-editor/xml/types/primary"
	"log"
)

type Parts struct {
	fmt.Stringer
	attr        attributes.Attributes
	Def         string `xml:"def"`
	OtherFields map[string]string
}

func (p *Parts) Assign(e *xml.Element) error {
	log.Print("Assign from parts is called")
	return nil
}

func (p *Parts) GetPath() string {
	return ""
}

func (p *Parts) SetAttributes(attributes attributes.Attributes) {
	log.Printf("attributes: %+#v", attributes)
	p.attr = attributes
}

func (p *Parts) GetAttributes() attributes.Attributes {
	return p.attr
}

func (p Parts) String() string {
	log.Printf("p: %+#v", p)
	return fmt.Sprintf("<%v>", p.Def)
}

type Scenario struct {
	Name          primary.EmbeddedType[string] `xml:"name"`
	Summary       string                       `xml:"summary"`
	Description   string                       `xml:"description"`
	PlayerFaction xml.Elements                 `xml:"playerFaction"`
	Parts         types.Slice[Parts]           `xml:"parts"`
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
