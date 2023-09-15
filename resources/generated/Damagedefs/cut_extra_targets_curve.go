// Code generated by rimworld-editor. DO NOT EDIT.

package damagedefs

import (
	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/attributes"
	"github.com/cruffinoni/rimworld-editor/xml/types"
)

type CutExtraTargetsCurve struct {
	Attr           attributes.Attributes
	FieldValidated map[string]bool

	Points *types.Slice[string] `xml:"points"`
}

func (c *CutExtraTargetsCurve) Assign(*xml.Element) error {
	return nil
}

func (c *CutExtraTargetsCurve) CountValidatedField() int {
	if c.FieldValidated == nil {
		return 0
	}
	return len(c.FieldValidated)
}

func (c *CutExtraTargetsCurve) Equal(*CutExtraTargetsCurve) bool {
	return false
}

func (c *CutExtraTargetsCurve) GetAttributes() attributes.Attributes {
	return c.Attr
}

func (c *CutExtraTargetsCurve) GetPath() string {
	return ""
}

func (c *CutExtraTargetsCurve) Greater(*CutExtraTargetsCurve) bool {
	return false
}

func (c *CutExtraTargetsCurve) IsValidField(field string) bool {
	return c.FieldValidated[field]
}

func (c *CutExtraTargetsCurve) Less(*CutExtraTargetsCurve) bool {
	return false
}

func (c *CutExtraTargetsCurve) SetAttributes(attr attributes.Attributes) {
	c.Attr = attr
	return
}

func (c *CutExtraTargetsCurve) Val() *CutExtraTargetsCurve {
	return nil
}

func (c *CutExtraTargetsCurve) ValidateField(field string) {
	if c.FieldValidated == nil {
		c.FieldValidated = make(map[string]bool)
	}
	c.FieldValidated[field] = true
	return
}
