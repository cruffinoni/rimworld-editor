// Code generated by rimworld-editor. DO NOT EDIT.

package stats

import (
	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/attributes"
	"github.com/cruffinoni/rimworld-editor/xml/types"
)

type SkillNeedFactors struct {
	Attr           attributes.Attributes
	FieldValidated map[string]bool

	Skill          string                `xml:"skill"`
	BaseValue      float64               `xml:"baseValue"`
	BonusPerLevel  float64               `xml:"bonusPerLevel"`
	Capacity       string                `xml:"capacity"`
	Weight         float64               `xml:"weight"`
	Max            float64               `xml:"max"`
	ValuesPerLevel *types.Slice[float64] `xml:"valuesPerLevel"`
	Required       bool                  `xml:"required"`
	AllowedDefect  float64               `xml:"allowedDefect"`
}

func (s *SkillNeedFactors) Assign(*xml.Element) error {
	return nil
}

func (s *SkillNeedFactors) CountValidatedField() int {
	if s.FieldValidated == nil {
		return 0
	}
	return len(s.FieldValidated)
}

func (s *SkillNeedFactors) Equal(*SkillNeedFactors) bool {
	return false
}

func (s *SkillNeedFactors) GetAttributes() attributes.Attributes {
	return s.Attr
}

func (s *SkillNeedFactors) GetPath() string {
	return ""
}

func (s *SkillNeedFactors) Greater(*SkillNeedFactors) bool {
	return false
}

func (s *SkillNeedFactors) IsValidField(field string) bool {
	return s.FieldValidated[field]
}

func (s *SkillNeedFactors) Less(*SkillNeedFactors) bool {
	return false
}

func (s *SkillNeedFactors) SetAttributes(attr attributes.Attributes) {
	s.Attr = attr
	return
}

func (s *SkillNeedFactors) Val() *SkillNeedFactors {
	return nil
}

func (s *SkillNeedFactors) ValidateField(field string) {
	if s.FieldValidated == nil {
		s.FieldValidated = make(map[string]bool)
	}
	s.FieldValidated[field] = true
	return
}
