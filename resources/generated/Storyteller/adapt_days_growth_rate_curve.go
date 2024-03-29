// Code generated by rimworld-editor. DO NOT EDIT.

package storyteller

import (
	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/attributes"
	"github.com/cruffinoni/rimworld-editor/xml/types"
)

type AdaptDaysGrowthRateCurve struct {
	Attr           attributes.Attributes
	FieldValidated map[string]bool

	Points *types.Slice[string] `xml:"points"`
}

func (a *AdaptDaysGrowthRateCurve) Assign(*xml.Element) error {
	return nil
}

func (a *AdaptDaysGrowthRateCurve) CountValidatedField() int {
	if a.FieldValidated == nil {
		return 0
	}
	return len(a.FieldValidated)
}

func (a *AdaptDaysGrowthRateCurve) Equal(*AdaptDaysGrowthRateCurve) bool {
	return false
}

func (a *AdaptDaysGrowthRateCurve) GetAttributes() attributes.Attributes {
	return a.Attr
}

func (a *AdaptDaysGrowthRateCurve) GetPath() string {
	return ""
}

func (a *AdaptDaysGrowthRateCurve) Greater(*AdaptDaysGrowthRateCurve) bool {
	return false
}

func (a *AdaptDaysGrowthRateCurve) IsValidField(field string) bool {
	return a.FieldValidated[field]
}

func (a *AdaptDaysGrowthRateCurve) Less(*AdaptDaysGrowthRateCurve) bool {
	return false
}

func (a *AdaptDaysGrowthRateCurve) SetAttributes(attr attributes.Attributes) {
	a.Attr = attr
	return
}

func (a *AdaptDaysGrowthRateCurve) Val() *AdaptDaysGrowthRateCurve {
	return nil
}

func (a *AdaptDaysGrowthRateCurve) ValidateField(field string) {
	if a.FieldValidated == nil {
		a.FieldValidated = make(map[string]bool)
	}
	a.FieldValidated[field] = true
	return
}
