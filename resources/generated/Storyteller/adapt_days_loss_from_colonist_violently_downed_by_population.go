// Code generated by rimworld-editor. DO NOT EDIT.

package storyteller

import (
	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/attributes"
	"github.com/cruffinoni/rimworld-editor/xml/types"
)

type AdaptDaysLossFromColonistViolentlyDownedByPopulation struct {
	Attr           attributes.Attributes
	FieldValidated map[string]bool

	Points *types.Slice[string] `xml:"points"`
}

func (a *AdaptDaysLossFromColonistViolentlyDownedByPopulation) Assign(*xml.Element) error {
	return nil
}

func (a *AdaptDaysLossFromColonistViolentlyDownedByPopulation) CountValidatedField() int {
	if a.FieldValidated == nil {
		return 0
	}
	return len(a.FieldValidated)
}

func (a *AdaptDaysLossFromColonistViolentlyDownedByPopulation) Equal(*AdaptDaysLossFromColonistViolentlyDownedByPopulation) bool {
	return false
}

func (a *AdaptDaysLossFromColonistViolentlyDownedByPopulation) GetAttributes() attributes.Attributes {
	return a.Attr
}

func (a *AdaptDaysLossFromColonistViolentlyDownedByPopulation) GetPath() string {
	return ""
}

func (a *AdaptDaysLossFromColonistViolentlyDownedByPopulation) Greater(*AdaptDaysLossFromColonistViolentlyDownedByPopulation) bool {
	return false
}

func (a *AdaptDaysLossFromColonistViolentlyDownedByPopulation) IsValidField(field string) bool {
	return a.FieldValidated[field]
}

func (a *AdaptDaysLossFromColonistViolentlyDownedByPopulation) Less(*AdaptDaysLossFromColonistViolentlyDownedByPopulation) bool {
	return false
}

func (a *AdaptDaysLossFromColonistViolentlyDownedByPopulation) SetAttributes(attr attributes.Attributes) {
	a.Attr = attr
	return
}

func (a *AdaptDaysLossFromColonistViolentlyDownedByPopulation) Val() *AdaptDaysLossFromColonistViolentlyDownedByPopulation {
	return nil
}

func (a *AdaptDaysLossFromColonistViolentlyDownedByPopulation) ValidateField(field string) {
	if a.FieldValidated == nil {
		a.FieldValidated = make(map[string]bool)
	}
	a.FieldValidated[field] = true
	return
}
