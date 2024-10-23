// Code generated by rimworld-editor. DO NOT EDIT.

package designations

import (
	"github.com/cruffinoni/rimworld-editor/internal/xml"
	"github.com/cruffinoni/rimworld-editor/internal/xml/attributes"
	"github.com/cruffinoni/rimworld-editor/internal/xml/types"
)

type DesignationCategoryDef struct {
	Attr           attributes.Attributes
	FieldValidated map[string]bool

	DefName                   string               `xml:"defName"`
	Label                     string               `xml:"label"`
	Order                     int64                `xml:"order"`
	SpecialDesignatorClasses  *types.Slice[string] `xml:"specialDesignatorClasses"`
	ShowPowerGrid             bool                 `xml:"showPowerGrid"`
	ResearchPrerequisites     *types.Slice[string] `xml:"researchPrerequisites"`
	TexturePath               string               `xml:"texturePath"`
	TargetType                string               `xml:"targetType"`
	ShouldBatchDraw           bool                 `xml:"shouldBatchDraw"`
	RemoveIfBuildingDespawned bool                 `xml:"removeIfBuildingDespawned"`
	DesignateCancelable       bool                 `xml:"designateCancelable"`
}

func (d *DesignationCategoryDef) Assign(*xml.Element) error {
	return nil
}

func (d *DesignationCategoryDef) CountValidatedField() int {
	if d.FieldValidated == nil {
		return 0
	}
	return len(d.FieldValidated)
}

func (d *DesignationCategoryDef) Equal(*DesignationCategoryDef) bool {
	return false
}

func (d *DesignationCategoryDef) GetAttributes() attributes.Attributes {
	return d.Attr
}

func (d *DesignationCategoryDef) GetPath() string {
	return ""
}

func (d *DesignationCategoryDef) Greater(*DesignationCategoryDef) bool {
	return false
}

func (d *DesignationCategoryDef) IsValidField(field string) bool {
	return d.FieldValidated[field]
}

func (d *DesignationCategoryDef) Less(*DesignationCategoryDef) bool {
	return false
}

func (d *DesignationCategoryDef) SetAttributes(attr attributes.Attributes) {
	d.Attr = attr
	return
}

func (d *DesignationCategoryDef) Val() *DesignationCategoryDef {
	return nil
}

func (d *DesignationCategoryDef) ValidateField(field string) {
	if d.FieldValidated == nil {
		d.FieldValidated = make(map[string]bool)
	}
	d.FieldValidated[field] = true
	return
}