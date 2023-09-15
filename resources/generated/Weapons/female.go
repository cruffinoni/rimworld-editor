// Code generated by rimworld-editor. DO NOT EDIT.

package weapons

import (
	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/attributes"
)

type Female struct {
	Attr           attributes.Attributes
	FieldValidated map[string]bool

	Offset string `xml:"offset"`
	Scale  string `xml:"scale"`
}

func (f *Female) Assign(*xml.Element) error {
	return nil
}

func (f *Female) CountValidatedField() int {
	if f.FieldValidated == nil {
		return 0
	}
	return len(f.FieldValidated)
}

func (f *Female) Equal(*Female) bool {
	return false
}

func (f *Female) GetAttributes() attributes.Attributes {
	return f.Attr
}

func (f *Female) GetPath() string {
	return ""
}

func (f *Female) Greater(*Female) bool {
	return false
}

func (f *Female) IsValidField(field string) bool {
	return f.FieldValidated[field]
}

func (f *Female) Less(*Female) bool {
	return false
}

func (f *Female) SetAttributes(attr attributes.Attributes) {
	f.Attr = attr
	return
}

func (f *Female) Val() *Female {
	return nil
}

func (f *Female) ValidateField(field string) {
	if f.FieldValidated == nil {
		f.FieldValidated = make(map[string]bool)
	}
	f.FieldValidated[field] = true
	return
}
