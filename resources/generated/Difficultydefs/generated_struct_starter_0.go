// Code generated by rimworld-editor. DO NOT EDIT.

package difficultydefs

import (
	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/attributes"
	"github.com/cruffinoni/rimworld-editor/xml/types"
)

type GeneratedStructStarter0 struct {
	Attr           attributes.Attributes
	FieldValidated map[string]bool

	Defs *types.Slice[*DifficultyDef] `xml:"Defs"`
}

func (g *GeneratedStructStarter0) Assign(*xml.Element) error {
	return nil
}

func (g *GeneratedStructStarter0) CountValidatedField() int {
	if g.FieldValidated == nil {
		return 0
	}
	return len(g.FieldValidated)
}

func (g *GeneratedStructStarter0) Equal(*GeneratedStructStarter0) bool {
	return false
}

func (g *GeneratedStructStarter0) GetAttributes() attributes.Attributes {
	return g.Attr
}

func (g *GeneratedStructStarter0) GetPath() string {
	return ""
}

func (g *GeneratedStructStarter0) Greater(*GeneratedStructStarter0) bool {
	return false
}

func (g *GeneratedStructStarter0) IsValidField(field string) bool {
	return g.FieldValidated[field]
}

func (g *GeneratedStructStarter0) Less(*GeneratedStructStarter0) bool {
	return false
}

func (g *GeneratedStructStarter0) SetAttributes(attr attributes.Attributes) {
	g.Attr = attr
	return
}

func (g *GeneratedStructStarter0) Val() *GeneratedStructStarter0 {
	return nil
}

func (g *GeneratedStructStarter0) ValidateField(field string) {
	if g.FieldValidated == nil {
		g.FieldValidated = make(map[string]bool)
	}
	g.FieldValidated[field] = true
	return
}
