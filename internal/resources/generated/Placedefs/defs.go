// Code generated by rimworld-editor. DO NOT EDIT.

package placedefs

import (
	"github.com/cruffinoni/rimworld-editor/internal/xml"
	"github.com/cruffinoni/rimworld-editor/internal/xml/attributes"
)

type Defs struct {
	Attr           attributes.Attributes
	FieldValidated map[string]bool

	RulePackDef *RulePackDef `xml:"RulePackDef"`
	PlaceDef    *PlaceDef    `xml:"PlaceDef"`
}

func (d *Defs) Assign(*xml.Element) error {
	return nil
}

func (d *Defs) CountValidatedField() int {
	if d.FieldValidated == nil {
		return 0
	}
	return len(d.FieldValidated)
}

func (d *Defs) Equal(*Defs) bool {
	return false
}

func (d *Defs) GetAttributes() attributes.Attributes {
	return d.Attr
}

func (d *Defs) GetPath() string {
	return ""
}

func (d *Defs) Greater(*Defs) bool {
	return false
}

func (d *Defs) IsValidField(field string) bool {
	return d.FieldValidated[field]
}

func (d *Defs) Less(*Defs) bool {
	return false
}

func (d *Defs) SetAttributes(attr attributes.Attributes) {
	d.Attr = attr
	return
}

func (d *Defs) Val() *Defs {
	return nil
}

func (d *Defs) ValidateField(field string) {
	if d.FieldValidated == nil {
		d.FieldValidated = make(map[string]bool)
	}
	d.FieldValidated[field] = true
	return
}
