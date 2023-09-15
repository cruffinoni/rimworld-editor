// Code generated by rimworld-editor. DO NOT EDIT.

package bodytypedefs

import (
	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/attributes"
	"github.com/cruffinoni/rimworld-editor/xml/types"
)

type BodyTypeDef struct {
	Attr           attributes.Attributes
	FieldValidated map[string]bool

	DefName                   string                      `xml:"defName"`
	HeadOffset                string                      `xml:"headOffset"`
	BodyNakedGraphicPath      string                      `xml:"bodyNakedGraphicPath"`
	BodyDessicatedGraphicPath string                      `xml:"bodyDessicatedGraphicPath"`
	BodyGraphicScale          string                      `xml:"bodyGraphicScale"`
	WoundAnchors              *types.Slice[*WoundAnchors] `xml:"woundAnchors"`
}

func (b *BodyTypeDef) Assign(*xml.Element) error {
	return nil
}

func (b *BodyTypeDef) CountValidatedField() int {
	if b.FieldValidated == nil {
		return 0
	}
	return len(b.FieldValidated)
}

func (b *BodyTypeDef) Equal(*BodyTypeDef) bool {
	return false
}

func (b *BodyTypeDef) GetAttributes() attributes.Attributes {
	return b.Attr
}

func (b *BodyTypeDef) GetPath() string {
	return ""
}

func (b *BodyTypeDef) Greater(*BodyTypeDef) bool {
	return false
}

func (b *BodyTypeDef) IsValidField(field string) bool {
	return b.FieldValidated[field]
}

func (b *BodyTypeDef) Less(*BodyTypeDef) bool {
	return false
}

func (b *BodyTypeDef) SetAttributes(attr attributes.Attributes) {
	b.Attr = attr
	return
}

func (b *BodyTypeDef) Val() *BodyTypeDef {
	return nil
}

func (b *BodyTypeDef) ValidateField(field string) {
	if b.FieldValidated == nil {
		b.FieldValidated = make(map[string]bool)
	}
	b.FieldValidated[field] = true
	return
}
