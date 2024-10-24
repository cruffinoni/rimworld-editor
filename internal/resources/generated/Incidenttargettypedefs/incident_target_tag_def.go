// Code generated by rimworld-editor. DO NOT EDIT.

package incidenttargettypedefs

import (
	"github.com/cruffinoni/rimworld-editor/internal/xml"
	"github.com/cruffinoni/rimworld-editor/internal/xml/attributes"
)

type IncidentTargetTagDef struct {
	Attr           attributes.Attributes
	FieldValidated map[string]bool

	DefName string `xml:"defName"`
}

func (i *IncidentTargetTagDef) Assign(*xml.Element) error {
	return nil
}

func (i *IncidentTargetTagDef) CountValidatedField() int {
	if i.FieldValidated == nil {
		return 0
	}
	return len(i.FieldValidated)
}

func (i *IncidentTargetTagDef) Equal(*IncidentTargetTagDef) bool {
	return false
}

func (i *IncidentTargetTagDef) GetAttributes() attributes.Attributes {
	return i.Attr
}

func (i *IncidentTargetTagDef) GetPath() string {
	return ""
}

func (i *IncidentTargetTagDef) Greater(*IncidentTargetTagDef) bool {
	return false
}

func (i *IncidentTargetTagDef) IsValidField(field string) bool {
	return i.FieldValidated[field]
}

func (i *IncidentTargetTagDef) Less(*IncidentTargetTagDef) bool {
	return false
}

func (i *IncidentTargetTagDef) SetAttributes(attr attributes.Attributes) {
	i.Attr = attr
	return
}

func (i *IncidentTargetTagDef) Val() *IncidentTargetTagDef {
	return nil
}

func (i *IncidentTargetTagDef) ValidateField(field string) {
	if i.FieldValidated == nil {
		i.FieldValidated = make(map[string]bool)
	}
	i.FieldValidated[field] = true
	return
}
