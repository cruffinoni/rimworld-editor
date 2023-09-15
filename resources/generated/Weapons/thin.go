// Code generated by rimworld-editor. DO NOT EDIT.

package weapons

import (
	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/attributes"
)

type Thin struct {
	Attr           attributes.Attributes
	FieldValidated map[string]bool

	Offset string `xml:"offset"`
	Scale  string `xml:"scale"`
}

func (t *Thin) Assign(*xml.Element) error {
	return nil
}

func (t *Thin) CountValidatedField() int {
	if t.FieldValidated == nil {
		return 0
	}
	return len(t.FieldValidated)
}

func (t *Thin) Equal(*Thin) bool {
	return false
}

func (t *Thin) GetAttributes() attributes.Attributes {
	return t.Attr
}

func (t *Thin) GetPath() string {
	return ""
}

func (t *Thin) Greater(*Thin) bool {
	return false
}

func (t *Thin) IsValidField(field string) bool {
	return t.FieldValidated[field]
}

func (t *Thin) Less(*Thin) bool {
	return false
}

func (t *Thin) SetAttributes(attr attributes.Attributes) {
	t.Attr = attr
	return
}

func (t *Thin) Val() *Thin {
	return nil
}

func (t *Thin) ValidateField(field string) {
	if t.FieldValidated == nil {
		t.FieldValidated = make(map[string]bool)
	}
	t.FieldValidated[field] = true
	return
}
