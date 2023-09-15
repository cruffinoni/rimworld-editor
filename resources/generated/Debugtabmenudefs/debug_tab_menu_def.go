// Code generated by rimworld-editor. DO NOT EDIT.

package debugtabmenudefs

import (
	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/attributes"
)

type DebugTabMenuDef struct {
	Attr           attributes.Attributes
	FieldValidated map[string]bool

	DefName      string `xml:"defName"`
	Label        string `xml:"label"`
	MenuClass    string `xml:"menuClass"`
	DisplayOrder int64  `xml:"displayOrder"`
}

func (d *DebugTabMenuDef) Assign(*xml.Element) error {
	return nil
}

func (d *DebugTabMenuDef) CountValidatedField() int {
	if d.FieldValidated == nil {
		return 0
	}
	return len(d.FieldValidated)
}

func (d *DebugTabMenuDef) Equal(*DebugTabMenuDef) bool {
	return false
}

func (d *DebugTabMenuDef) GetAttributes() attributes.Attributes {
	return d.Attr
}

func (d *DebugTabMenuDef) GetPath() string {
	return ""
}

func (d *DebugTabMenuDef) Greater(*DebugTabMenuDef) bool {
	return false
}

func (d *DebugTabMenuDef) IsValidField(field string) bool {
	return d.FieldValidated[field]
}

func (d *DebugTabMenuDef) Less(*DebugTabMenuDef) bool {
	return false
}

func (d *DebugTabMenuDef) SetAttributes(attr attributes.Attributes) {
	d.Attr = attr
	return
}

func (d *DebugTabMenuDef) Val() *DebugTabMenuDef {
	return nil
}

func (d *DebugTabMenuDef) ValidateField(field string) {
	if d.FieldValidated == nil {
		d.FieldValidated = make(map[string]bool)
	}
	d.FieldValidated[field] = true
	return
}
