// Code generated by rimworld-editor. DO NOT EDIT.

package factiondefs

import (
	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/attributes"
	"github.com/cruffinoni/rimworld-editor/xml/types/embedded"
)

type XenotypeChances struct {
	Attr           attributes.Attributes
	FieldValidated map[string]bool

	Neanderthal *embedded.Type[float64] `xml:"Neanderthal"`
	Hussar      *embedded.Type[float64] `xml:"Hussar"`
	Genie       *embedded.Type[float64] `xml:"Genie"`
	Highmate    *embedded.Type[float64] `xml:"Highmate"`
	Pigskin     *embedded.Type[float64] `xml:"Pigskin"`
	Yttakin     *embedded.Type[float64] `xml:"Yttakin"`
	Dirtmole    *embedded.Type[float64] `xml:"Dirtmole"`
	Impid       *embedded.Type[float64] `xml:"Impid"`
	Waster      *embedded.Type[float64] `xml:"Waster"`
}

func (x *XenotypeChances) Assign(*xml.Element) error {
	return nil
}

func (x *XenotypeChances) CountValidatedField() int {
	if x.FieldValidated == nil {
		return 0
	}
	return len(x.FieldValidated)
}

func (x *XenotypeChances) Equal(*XenotypeChances) bool {
	return false
}

func (x *XenotypeChances) GetAttributes() attributes.Attributes {
	return x.Attr
}

func (x *XenotypeChances) GetPath() string {
	return ""
}

func (x *XenotypeChances) Greater(*XenotypeChances) bool {
	return false
}

func (x *XenotypeChances) IsValidField(field string) bool {
	return x.FieldValidated[field]
}

func (x *XenotypeChances) Less(*XenotypeChances) bool {
	return false
}

func (x *XenotypeChances) SetAttributes(attr attributes.Attributes) {
	x.Attr = attr
	return
}

func (x *XenotypeChances) Val() *XenotypeChances {
	return nil
}

func (x *XenotypeChances) ValidateField(field string) {
	if x.FieldValidated == nil {
		x.FieldValidated = make(map[string]bool)
	}
	x.FieldValidated[field] = true
	return
}