// Code generated by rimworld-editor. DO NOT EDIT.

package placedefs

import (
	"github.com/cruffinoni/rimworld-editor/internal/xml"
	"github.com/cruffinoni/rimworld-editor/internal/xml/attributes"
	"github.com/cruffinoni/rimworld-editor/internal/xml/types"
)

type RulePack struct {
	Attr           attributes.Attributes
	FieldValidated map[string]bool

	RulesStrings *types.Slice[string] `xml:"rulesStrings"`
}

func (r *RulePack) Assign(*xml.Element) error {
	return nil
}

func (r *RulePack) CountValidatedField() int {
	if r.FieldValidated == nil {
		return 0
	}
	return len(r.FieldValidated)
}

func (r *RulePack) Equal(*RulePack) bool {
	return false
}

func (r *RulePack) GetAttributes() attributes.Attributes {
	return r.Attr
}

func (r *RulePack) GetPath() string {
	return ""
}

func (r *RulePack) Greater(*RulePack) bool {
	return false
}

func (r *RulePack) IsValidField(field string) bool {
	return r.FieldValidated[field]
}

func (r *RulePack) Less(*RulePack) bool {
	return false
}

func (r *RulePack) SetAttributes(attr attributes.Attributes) {
	r.Attr = attr
	return
}

func (r *RulePack) Val() *RulePack {
	return nil
}

func (r *RulePack) ValidateField(field string) {
	if r.FieldValidated == nil {
		r.FieldValidated = make(map[string]bool)
	}
	r.FieldValidated[field] = true
	return
}
