// Code generated by rimworld-editor. DO NOT EDIT.

package inspirationdefs

import (
	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/attributes"
	"github.com/cruffinoni/rimworld-editor/xml/types"
)

type InspirationDef struct {
	Attr           attributes.Attributes
	FieldValidated map[string]bool

	DefName                        string               `xml:"defName"`
	Label                          string               `xml:"label"`
	BaseDurationDays               int64                `xml:"baseDurationDays"`
	BeginLetter                    string               `xml:"beginLetter"`
	BeginLetterDef                 string               `xml:"beginLetterDef"`
	EndMessage                     string               `xml:"endMessage"`
	BaseInspectLine                string               `xml:"baseInspectLine"`
	RequiredNonDisabledWorkTags    string               `xml:"requiredNonDisabledWorkTags"`
	MinAge                         int64                `xml:"minAge"`
	StatFactors                    *StatFactors         `xml:"statFactors"`
	AllowedOnDownedPawns           bool                 `xml:"allowedOnDownedPawns"`
	RequiredCapacities             *types.Slice[string] `xml:"requiredCapacities"`
	RequiredNonDisabledStats       *types.Slice[string] `xml:"requiredNonDisabledStats"`
	StatOffsets                    *StatOffsets         `xml:"statOffsets"`
	AssociatedSkills               *types.Slice[string] `xml:"associatedSkills"`
	RequiredSkills                 *RequiredSkills      `xml:"requiredSkills"`
	RequiredNonDisabledWorkTypes   *types.Slice[string] `xml:"requiredNonDisabledWorkTypes"`
	RequiredAnyNonDisabledWorkType *types.Slice[string] `xml:"requiredAnyNonDisabledWorkType"`
	RequiredAnySkill               *RequiredAnySkill    `xml:"requiredAnySkill"`
}

func (i *InspirationDef) Assign(*xml.Element) error {
	return nil
}

func (i *InspirationDef) CountValidatedField() int {
	if i.FieldValidated == nil {
		return 0
	}
	return len(i.FieldValidated)
}

func (i *InspirationDef) Equal(*InspirationDef) bool {
	return false
}

func (i *InspirationDef) GetAttributes() attributes.Attributes {
	return i.Attr
}

func (i *InspirationDef) GetPath() string {
	return ""
}

func (i *InspirationDef) Greater(*InspirationDef) bool {
	return false
}

func (i *InspirationDef) IsValidField(field string) bool {
	return i.FieldValidated[field]
}

func (i *InspirationDef) Less(*InspirationDef) bool {
	return false
}

func (i *InspirationDef) SetAttributes(attr attributes.Attributes) {
	i.Attr = attr
	return
}

func (i *InspirationDef) Val() *InspirationDef {
	return nil
}

func (i *InspirationDef) ValidateField(field string) {
	if i.FieldValidated == nil {
		i.FieldValidated = make(map[string]bool)
	}
	i.FieldValidated[field] = true
	return
}