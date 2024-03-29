// Code generated by rimworld-editor. DO NOT EDIT.

package bodyparts

import (
	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/attributes"
)

type GraphicData struct {
	Attr           attributes.Attributes
	FieldValidated map[string]bool

	TexPath      string  `xml:"texPath"`
	Color        string  `xml:"color"`
	GraphicClass string  `xml:"graphicClass"`
	DrawSize     float64 `xml:"drawSize"`
}

func (g *GraphicData) Assign(*xml.Element) error {
	return nil
}

func (g *GraphicData) CountValidatedField() int {
	if g.FieldValidated == nil {
		return 0
	}
	return len(g.FieldValidated)
}

func (g *GraphicData) Equal(*GraphicData) bool {
	return false
}

func (g *GraphicData) GetAttributes() attributes.Attributes {
	return g.Attr
}

func (g *GraphicData) GetPath() string {
	return ""
}

func (g *GraphicData) Greater(*GraphicData) bool {
	return false
}

func (g *GraphicData) IsValidField(field string) bool {
	return g.FieldValidated[field]
}

func (g *GraphicData) Less(*GraphicData) bool {
	return false
}

func (g *GraphicData) SetAttributes(attr attributes.Attributes) {
	g.Attr = attr
	return
}

func (g *GraphicData) Val() *GraphicData {
	return nil
}

func (g *GraphicData) ValidateField(field string) {
	if g.FieldValidated == nil {
		g.FieldValidated = make(map[string]bool)
	}
	g.FieldValidated[field] = true
	return
}
