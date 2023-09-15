// Code generated by rimworld-editor. DO NOT EDIT.

package factiondefs

import (
	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/attributes"
	"github.com/cruffinoni/rimworld-editor/xml/types"
)

type BackstoryFilters struct {
	Attr           attributes.Attributes
	FieldValidated map[string]bool

	Categories  *types.Slice[string] `xml:"categories"`
	Commonality float64              `xml:"commonality"`
	ThingDefs   *types.Slice[string] `xml:"thingDefs"`
	Li          *types.Slice[string] `xml:"li"`
}

func (b *BackstoryFilters) Assign(*xml.Element) error {
	return nil
}

func (b *BackstoryFilters) CountValidatedField() int {
	if b.FieldValidated == nil {
		return 0
	}
	return len(b.FieldValidated)
}

func (b *BackstoryFilters) Equal(*BackstoryFilters) bool {
	return false
}

func (b *BackstoryFilters) GetAttributes() attributes.Attributes {
	return b.Attr
}

func (b *BackstoryFilters) GetPath() string {
	return ""
}

func (b *BackstoryFilters) Greater(*BackstoryFilters) bool {
	return false
}

func (b *BackstoryFilters) IsValidField(field string) bool {
	return b.FieldValidated[field]
}

func (b *BackstoryFilters) Less(*BackstoryFilters) bool {
	return false
}

func (b *BackstoryFilters) SetAttributes(attr attributes.Attributes) {
	b.Attr = attr
	return
}

func (b *BackstoryFilters) Val() *BackstoryFilters {
	return nil
}

func (b *BackstoryFilters) ValidateField(field string) {
	if b.FieldValidated == nil {
		b.FieldValidated = make(map[string]bool)
	}
	b.FieldValidated[field] = true
	return
}
