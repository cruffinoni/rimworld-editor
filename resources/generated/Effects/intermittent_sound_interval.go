// Code generated by rimworld-editor. DO NOT EDIT.

package effects

import (
	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/attributes"
)

type IntermittentSoundInterval struct {
	Attr           attributes.Attributes
	FieldValidated map[string]bool

	Min int64 `xml:"min"`
	Max int64 `xml:"max"`
}

func (i *IntermittentSoundInterval) Assign(*xml.Element) error {
	return nil
}

func (i *IntermittentSoundInterval) CountValidatedField() int {
	if i.FieldValidated == nil {
		return 0
	}
	return len(i.FieldValidated)
}

func (i *IntermittentSoundInterval) Equal(*IntermittentSoundInterval) bool {
	return false
}

func (i *IntermittentSoundInterval) GetAttributes() attributes.Attributes {
	return i.Attr
}

func (i *IntermittentSoundInterval) GetPath() string {
	return ""
}

func (i *IntermittentSoundInterval) Greater(*IntermittentSoundInterval) bool {
	return false
}

func (i *IntermittentSoundInterval) IsValidField(field string) bool {
	return i.FieldValidated[field]
}

func (i *IntermittentSoundInterval) Less(*IntermittentSoundInterval) bool {
	return false
}

func (i *IntermittentSoundInterval) SetAttributes(attr attributes.Attributes) {
	i.Attr = attr
	return
}

func (i *IntermittentSoundInterval) Val() *IntermittentSoundInterval {
	return nil
}

func (i *IntermittentSoundInterval) ValidateField(field string) {
	if i.FieldValidated == nil {
		i.FieldValidated = make(map[string]bool)
	}
	i.FieldValidated[field] = true
	return
}
