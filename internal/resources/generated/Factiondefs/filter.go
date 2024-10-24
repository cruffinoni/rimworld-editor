// Code generated by rimworld-editor. DO NOT EDIT.

package factiondefs

import (
	"github.com/cruffinoni/rimworld-editor/internal/xml"
	"github.com/cruffinoni/rimworld-editor/internal/xml/attributes"
	"github.com/cruffinoni/rimworld-editor/internal/xml/types"
)

type Filter struct {
	Attr           attributes.Attributes
	FieldValidated map[string]bool

	ThingDefs *types.Slice[string] `xml:"thingDefs"`
}

func (f *Filter) Assign(*xml.Element) error {
	return nil
}

func (f *Filter) CountValidatedField() int {
	if f.FieldValidated == nil {
		return 0
	}
	return len(f.FieldValidated)
}

func (f *Filter) Equal(*Filter) bool {
	return false
}

func (f *Filter) GetAttributes() attributes.Attributes {
	return f.Attr
}

func (f *Filter) GetPath() string {
	return ""
}

func (f *Filter) Greater(*Filter) bool {
	return false
}

func (f *Filter) IsValidField(field string) bool {
	return f.FieldValidated[field]
}

func (f *Filter) Less(*Filter) bool {
	return false
}

func (f *Filter) SetAttributes(attr attributes.Attributes) {
	f.Attr = attr
	return
}

func (f *Filter) Val() *Filter {
	return nil
}

func (f *Filter) ValidateField(field string) {
	if f.FieldValidated == nil {
		f.FieldValidated = make(map[string]bool)
	}
	f.FieldValidated[field] = true
	return
}
