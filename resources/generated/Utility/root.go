// Code generated by rimworld-editor. DO NOT EDIT.

package utility

import (
	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/attributes"
	"github.com/cruffinoni/rimworld-editor/xml/types"
)

type Root struct {
	Attr           attributes.Attributes
	FieldValidated map[string]bool

	Nodes                  *types.Slice[*Nodes] `xml:"nodes"`
	CustomLetterLabel      string               `xml:"customLetterLabel"`
	CustomLetterText       string               `xml:"customLetterText"`
	CustomLetterLabelRules string               `xml:"customLetterLabelRules"`
	CustomLetterTextRules  string               `xml:"customLetterTextRules"`
}

func (r *Root) Assign(*xml.Element) error {
	return nil
}

func (r *Root) CountValidatedField() int {
	if r.FieldValidated == nil {
		return 0
	}
	return len(r.FieldValidated)
}

func (r *Root) Equal(*Root) bool {
	return false
}

func (r *Root) GetAttributes() attributes.Attributes {
	return r.Attr
}

func (r *Root) GetPath() string {
	return ""
}

func (r *Root) Greater(*Root) bool {
	return false
}

func (r *Root) IsValidField(field string) bool {
	return r.FieldValidated[field]
}

func (r *Root) Less(*Root) bool {
	return false
}

func (r *Root) SetAttributes(attr attributes.Attributes) {
	r.Attr = attr
	return
}

func (r *Root) Val() *Root {
	return nil
}

func (r *Root) ValidateField(field string) {
	if r.FieldValidated == nil {
		r.FieldValidated = make(map[string]bool)
	}
	r.FieldValidated[field] = true
	return
}
