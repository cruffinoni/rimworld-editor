// Code generated by rimworld-editor. DO NOT EDIT.

package storyteller

import (
	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/attributes"
	"github.com/cruffinoni/rimworld-editor/xml/types"
)

type PopulationIntentFactorFromPopCurve struct {
	Attr           attributes.Attributes
	FieldValidated map[string]bool

	Points *types.Slice[string] `xml:"points"`
}

func (p *PopulationIntentFactorFromPopCurve) Assign(*xml.Element) error {
	return nil
}

func (p *PopulationIntentFactorFromPopCurve) CountValidatedField() int {
	if p.FieldValidated == nil {
		return 0
	}
	return len(p.FieldValidated)
}

func (p *PopulationIntentFactorFromPopCurve) Equal(*PopulationIntentFactorFromPopCurve) bool {
	return false
}

func (p *PopulationIntentFactorFromPopCurve) GetAttributes() attributes.Attributes {
	return p.Attr
}

func (p *PopulationIntentFactorFromPopCurve) GetPath() string {
	return ""
}

func (p *PopulationIntentFactorFromPopCurve) Greater(*PopulationIntentFactorFromPopCurve) bool {
	return false
}

func (p *PopulationIntentFactorFromPopCurve) IsValidField(field string) bool {
	return p.FieldValidated[field]
}

func (p *PopulationIntentFactorFromPopCurve) Less(*PopulationIntentFactorFromPopCurve) bool {
	return false
}

func (p *PopulationIntentFactorFromPopCurve) SetAttributes(attr attributes.Attributes) {
	p.Attr = attr
	return
}

func (p *PopulationIntentFactorFromPopCurve) Val() *PopulationIntentFactorFromPopCurve {
	return nil
}

func (p *PopulationIntentFactorFromPopCurve) ValidateField(field string) {
	if p.FieldValidated == nil {
		p.FieldValidated = make(map[string]bool)
	}
	p.FieldValidated[field] = true
	return
}
