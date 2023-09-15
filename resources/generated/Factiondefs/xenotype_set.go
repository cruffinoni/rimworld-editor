// Code generated by rimworld-editor. DO NOT EDIT.

package factiondefs

import (
	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/attributes"
)

type XenotypeSet struct {
	Attr           attributes.Attributes
	FieldValidated map[string]bool

	XenotypeChances *XenotypeChances `xml:"xenotypeChances"`
}

func (x *XenotypeSet) Assign(*xml.Element) error {
	return nil
}

func (x *XenotypeSet) CountValidatedField() int {
	if x.FieldValidated == nil {
		return 0
	}
	return len(x.FieldValidated)
}

func (x *XenotypeSet) Equal(*XenotypeSet) bool {
	return false
}

func (x *XenotypeSet) GetAttributes() attributes.Attributes {
	return x.Attr
}

func (x *XenotypeSet) GetPath() string {
	return ""
}

func (x *XenotypeSet) Greater(*XenotypeSet) bool {
	return false
}

func (x *XenotypeSet) IsValidField(field string) bool {
	return x.FieldValidated[field]
}

func (x *XenotypeSet) Less(*XenotypeSet) bool {
	return false
}

func (x *XenotypeSet) SetAttributes(attr attributes.Attributes) {
	x.Attr = attr
	return
}

func (x *XenotypeSet) Val() *XenotypeSet {
	return nil
}

func (x *XenotypeSet) ValidateField(field string) {
	if x.FieldValidated == nil {
		x.FieldValidated = make(map[string]bool)
	}
	x.FieldValidated[field] = true
	return
}
