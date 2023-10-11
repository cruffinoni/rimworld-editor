// Code generated by rimworld-editor. DO NOT EDIT.

package factiondefs

import (
	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/attributes"
	"github.com/cruffinoni/rimworld-editor/xml/types"
)

type RaidLootValueFromPointsCurve struct {
	Attr           attributes.Attributes
	FieldValidated map[string]bool

	Points *types.Slice[string] `xml:"points"`
}

func (r *RaidLootValueFromPointsCurve) Assign(*xml.Element) error {
	return nil
}

func (r *RaidLootValueFromPointsCurve) CountValidatedField() int {
	if r.FieldValidated == nil {
		return 0
	}
	return len(r.FieldValidated)
}

func (r *RaidLootValueFromPointsCurve) Equal(*RaidLootValueFromPointsCurve) bool {
	return false
}

func (r *RaidLootValueFromPointsCurve) GetAttributes() attributes.Attributes {
	return r.Attr
}

func (r *RaidLootValueFromPointsCurve) GetPath() string {
	return ""
}

func (r *RaidLootValueFromPointsCurve) Greater(*RaidLootValueFromPointsCurve) bool {
	return false
}

func (r *RaidLootValueFromPointsCurve) IsValidField(field string) bool {
	return r.FieldValidated[field]
}

func (r *RaidLootValueFromPointsCurve) Less(*RaidLootValueFromPointsCurve) bool {
	return false
}

func (r *RaidLootValueFromPointsCurve) SetAttributes(attr attributes.Attributes) {
	r.Attr = attr
	return
}

func (r *RaidLootValueFromPointsCurve) Val() *RaidLootValueFromPointsCurve {
	return nil
}

func (r *RaidLootValueFromPointsCurve) ValidateField(field string) {
	if r.FieldValidated == nil {
		r.FieldValidated = make(map[string]bool)
	}
	r.FieldValidated[field] = true
	return
}