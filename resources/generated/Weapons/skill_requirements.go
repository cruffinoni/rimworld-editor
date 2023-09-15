// Code generated by rimworld-editor. DO NOT EDIT.

package weapons

import (
	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/attributes"
)

type SkillRequirements struct {
	Attr           attributes.Attributes
	FieldValidated map[string]bool

	Crafting int64 `xml:"Crafting"`
}

func (s *SkillRequirements) Assign(*xml.Element) error {
	return nil
}

func (s *SkillRequirements) CountValidatedField() int {
	if s.FieldValidated == nil {
		return 0
	}
	return len(s.FieldValidated)
}

func (s *SkillRequirements) Equal(*SkillRequirements) bool {
	return false
}

func (s *SkillRequirements) GetAttributes() attributes.Attributes {
	return s.Attr
}

func (s *SkillRequirements) GetPath() string {
	return ""
}

func (s *SkillRequirements) Greater(*SkillRequirements) bool {
	return false
}

func (s *SkillRequirements) IsValidField(field string) bool {
	return s.FieldValidated[field]
}

func (s *SkillRequirements) Less(*SkillRequirements) bool {
	return false
}

func (s *SkillRequirements) SetAttributes(attr attributes.Attributes) {
	s.Attr = attr
	return
}

func (s *SkillRequirements) Val() *SkillRequirements {
	return nil
}

func (s *SkillRequirements) ValidateField(field string) {
	if s.FieldValidated == nil {
		s.FieldValidated = make(map[string]bool)
	}
	s.FieldValidated[field] = true
	return
}
