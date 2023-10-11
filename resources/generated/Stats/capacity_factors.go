// Code generated by rimworld-editor. DO NOT EDIT.

package stats

import (
	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/attributes"
)

type CapacityFactors struct {
	Attr           attributes.Attributes
	FieldValidated map[string]bool

	Capacity                 string               `xml:"capacity"`
	Weight                   float64              `xml:"weight"`
	Max                      float64              `xml:"max"`
	AllowedDefect            float64              `xml:"allowedDefect"`
	FactorFromGlowCurve      *FactorFromGlowCurve `xml:"factorFromGlowCurve"`
	IgnoreIfIncapableOfSight bool                 `xml:"ignoreIfIncapableOfSight"`
	IgnoreIfPrefersDarkness  bool                 `xml:"ignoreIfPrefersDarkness"`
	HumanlikeOnly            bool                 `xml:"humanlikeOnly"`
}

func (c *CapacityFactors) Assign(*xml.Element) error {
	return nil
}

func (c *CapacityFactors) CountValidatedField() int {
	if c.FieldValidated == nil {
		return 0
	}
	return len(c.FieldValidated)
}

func (c *CapacityFactors) Equal(*CapacityFactors) bool {
	return false
}

func (c *CapacityFactors) GetAttributes() attributes.Attributes {
	return c.Attr
}

func (c *CapacityFactors) GetPath() string {
	return ""
}

func (c *CapacityFactors) Greater(*CapacityFactors) bool {
	return false
}

func (c *CapacityFactors) IsValidField(field string) bool {
	return c.FieldValidated[field]
}

func (c *CapacityFactors) Less(*CapacityFactors) bool {
	return false
}

func (c *CapacityFactors) SetAttributes(attr attributes.Attributes) {
	c.Attr = attr
	return
}

func (c *CapacityFactors) Val() *CapacityFactors {
	return nil
}

func (c *CapacityFactors) ValidateField(field string) {
	if c.FieldValidated == nil {
		c.FieldValidated = make(map[string]bool)
	}
	c.FieldValidated[field] = true
	return
}