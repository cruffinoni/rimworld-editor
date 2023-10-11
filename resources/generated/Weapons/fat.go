// Code generated by rimworld-editor. DO NOT EDIT.

package weapons

import (
	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/attributes"
)

type Fat struct {
	Attr           attributes.Attributes
	FieldValidated map[string]bool

	Offset string `xml:"offset"`
	Scale  string `xml:"scale"`
}

func (f *Fat) Assign(*xml.Element) error {
	return nil
}

func (f *Fat) CountValidatedField() int {
	if f.FieldValidated == nil {
		return 0
	}
	return len(f.FieldValidated)
}

func (f *Fat) Equal(*Fat) bool {
	return false
}

func (f *Fat) GetAttributes() attributes.Attributes {
	return f.Attr
}

func (f *Fat) GetPath() string {
	return ""
}

func (f *Fat) Greater(*Fat) bool {
	return false
}

func (f *Fat) IsValidField(field string) bool {
	return f.FieldValidated[field]
}

func (f *Fat) Less(*Fat) bool {
	return false
}

func (f *Fat) SetAttributes(attr attributes.Attributes) {
	f.Attr = attr
	return
}

func (f *Fat) Val() *Fat {
	return nil
}

func (f *Fat) ValidateField(field string) {
	if f.FieldValidated == nil {
		f.FieldValidated = make(map[string]bool)
	}
	f.FieldValidated[field] = true
	return
}