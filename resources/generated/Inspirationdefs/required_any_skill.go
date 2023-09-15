// Code generated by rimworld-editor. DO NOT EDIT.

package inspirationdefs

import (
	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/attributes"
)

type RequiredAnySkill struct {
	Attr           attributes.Attributes
	FieldValidated map[string]bool

	Construction int64 `xml:"Construction"`
	Artistic     int64 `xml:"Artistic"`
	Crafting     int64 `xml:"Crafting"`
}

func (r *RequiredAnySkill) Assign(*xml.Element) error {
	return nil
}

func (r *RequiredAnySkill) CountValidatedField() int {
	if r.FieldValidated == nil {
		return 0
	}
	return len(r.FieldValidated)
}

func (r *RequiredAnySkill) Equal(*RequiredAnySkill) bool {
	return false
}

func (r *RequiredAnySkill) GetAttributes() attributes.Attributes {
	return r.Attr
}

func (r *RequiredAnySkill) GetPath() string {
	return ""
}

func (r *RequiredAnySkill) Greater(*RequiredAnySkill) bool {
	return false
}

func (r *RequiredAnySkill) IsValidField(field string) bool {
	return r.FieldValidated[field]
}

func (r *RequiredAnySkill) Less(*RequiredAnySkill) bool {
	return false
}

func (r *RequiredAnySkill) SetAttributes(attr attributes.Attributes) {
	r.Attr = attr
	return
}

func (r *RequiredAnySkill) Val() *RequiredAnySkill {
	return nil
}

func (r *RequiredAnySkill) ValidateField(field string) {
	if r.FieldValidated == nil {
		r.FieldValidated = make(map[string]bool)
	}
	r.FieldValidated[field] = true
	return
}
