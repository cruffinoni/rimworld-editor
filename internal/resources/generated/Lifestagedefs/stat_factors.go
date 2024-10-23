// Code generated by rimworld-editor. DO NOT EDIT.

package lifestagedefs

import (
	"github.com/cruffinoni/rimworld-editor/internal/xml"
	"github.com/cruffinoni/rimworld-editor/internal/xml/attributes"
)

type StatFactors struct {
	Attr           attributes.Attributes
	FieldValidated map[string]bool

	MoveSpeed float64 `xml:"MoveSpeed"`
}

func (s *StatFactors) Assign(*xml.Element) error {
	return nil
}

func (s *StatFactors) CountValidatedField() int {
	if s.FieldValidated == nil {
		return 0
	}
	return len(s.FieldValidated)
}

func (s *StatFactors) Equal(*StatFactors) bool {
	return false
}

func (s *StatFactors) GetAttributes() attributes.Attributes {
	return s.Attr
}

func (s *StatFactors) GetPath() string {
	return ""
}

func (s *StatFactors) Greater(*StatFactors) bool {
	return false
}

func (s *StatFactors) IsValidField(field string) bool {
	return s.FieldValidated[field]
}

func (s *StatFactors) Less(*StatFactors) bool {
	return false
}

func (s *StatFactors) SetAttributes(attr attributes.Attributes) {
	s.Attr = attr
	return
}

func (s *StatFactors) Val() *StatFactors {
	return nil
}

func (s *StatFactors) ValidateField(field string) {
	if s.FieldValidated == nil {
		s.FieldValidated = make(map[string]bool)
	}
	s.FieldValidated[field] = true
	return
}