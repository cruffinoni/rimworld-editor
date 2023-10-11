// Code generated by rimworld-editor. DO NOT EDIT.

package effects

import (
	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/attributes"
)

type Grains struct {
	Attr           attributes.Attributes
	FieldValidated map[string]bool

	ClipPath       string `xml:"clipPath"`
	ClipFolderPath string `xml:"clipFolderPath"`
}

func (g *Grains) Assign(*xml.Element) error {
	return nil
}

func (g *Grains) CountValidatedField() int {
	if g.FieldValidated == nil {
		return 0
	}
	return len(g.FieldValidated)
}

func (g *Grains) Equal(*Grains) bool {
	return false
}

func (g *Grains) GetAttributes() attributes.Attributes {
	return g.Attr
}

func (g *Grains) GetPath() string {
	return ""
}

func (g *Grains) Greater(*Grains) bool {
	return false
}

func (g *Grains) IsValidField(field string) bool {
	return g.FieldValidated[field]
}

func (g *Grains) Less(*Grains) bool {
	return false
}

func (g *Grains) SetAttributes(attr attributes.Attributes) {
	g.Attr = attr
	return
}

func (g *Grains) Val() *Grains {
	return nil
}

func (g *Grains) ValidateField(field string) {
	if g.FieldValidated == nil {
		g.FieldValidated = make(map[string]bool)
	}
	g.FieldValidated[field] = true
	return
}