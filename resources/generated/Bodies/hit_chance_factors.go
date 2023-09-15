// Code generated by rimworld-editor. DO NOT EDIT.

package bodies

import (
	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/attributes"
)

type HitChanceFactors struct {
	Attr           attributes.Attributes
	FieldValidated map[string]bool

	Key   string `xml:"key"`
	Value int64  `xml:"value"`
}

func (h *HitChanceFactors) Assign(*xml.Element) error {
	return nil
}

func (h *HitChanceFactors) CountValidatedField() int {
	if h.FieldValidated == nil {
		return 0
	}
	return len(h.FieldValidated)
}

func (h *HitChanceFactors) Equal(*HitChanceFactors) bool {
	return false
}

func (h *HitChanceFactors) GetAttributes() attributes.Attributes {
	return h.Attr
}

func (h *HitChanceFactors) GetPath() string {
	return ""
}

func (h *HitChanceFactors) Greater(*HitChanceFactors) bool {
	return false
}

func (h *HitChanceFactors) IsValidField(field string) bool {
	return h.FieldValidated[field]
}

func (h *HitChanceFactors) Less(*HitChanceFactors) bool {
	return false
}

func (h *HitChanceFactors) SetAttributes(attr attributes.Attributes) {
	h.Attr = attr
	return
}

func (h *HitChanceFactors) Val() *HitChanceFactors {
	return nil
}

func (h *HitChanceFactors) ValidateField(field string) {
	if h.FieldValidated == nil {
		h.FieldValidated = make(map[string]bool)
	}
	h.FieldValidated[field] = true
	return
}
