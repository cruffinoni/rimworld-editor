// Code generated by rimworld-editor. DO NOT EDIT.

package bearddefs

import (
	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/attributes"
	"github.com/cruffinoni/rimworld-editor/xml/types"
)

type BeardDef struct {
	Attr           attributes.Attributes
	FieldValidated map[string]bool

	DefName           string               `xml:"defName"`
	Label             string               `xml:"label"`
	NoGraphic         bool                 `xml:"noGraphic"`
	Category          string               `xml:"category"`
	StyleTags         *types.Slice[string] `xml:"styleTags"`
	TexPath           string               `xml:"texPath"`
	StyleGender       string               `xml:"styleGender"`
	OffsetNarrowEast  string               `xml:"offsetNarrowEast"`
	OffsetNarrowSouth string               `xml:"offsetNarrowSouth"`
}

func (b *BeardDef) Assign(*xml.Element) error {
	return nil
}

func (b *BeardDef) CountValidatedField() int {
	if b.FieldValidated == nil {
		return 0
	}
	return len(b.FieldValidated)
}

func (b *BeardDef) Equal(*BeardDef) bool {
	return false
}

func (b *BeardDef) GetAttributes() attributes.Attributes {
	return b.Attr
}

func (b *BeardDef) GetPath() string {
	return ""
}

func (b *BeardDef) Greater(*BeardDef) bool {
	return false
}

func (b *BeardDef) IsValidField(field string) bool {
	return b.FieldValidated[field]
}

func (b *BeardDef) Less(*BeardDef) bool {
	return false
}

func (b *BeardDef) SetAttributes(attr attributes.Attributes) {
	b.Attr = attr
	return
}

func (b *BeardDef) Val() *BeardDef {
	return nil
}

func (b *BeardDef) ValidateField(field string) {
	if b.FieldValidated == nil {
		b.FieldValidated = make(map[string]bool)
	}
	b.FieldValidated[field] = true
	return
}
