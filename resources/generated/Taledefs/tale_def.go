// Code generated by rimworld-editor. DO NOT EDIT.

package taledefs

import (
	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/attributes"
)

type TaleDef struct {
	Attr           attributes.Attributes
	FieldValidated map[string]bool

	DefName           string    `xml:"defName"`
	Label             string    `xml:"label"`
	TaleClass         string    `xml:"taleClass"`
	Type              string    `xml:"type"`
	BaseInterest      float64   `xml:"baseInterest"`
	RulePack          *RulePack `xml:"rulePack"`
	FirstPawnSymbol   string    `xml:"firstPawnSymbol"`
	SecondPawnSymbol  string    `xml:"secondPawnSymbol"`
	DefType           string    `xml:"defType"`
	DefSymbol         string    `xml:"defSymbol"`
	HistoryGraphColor string    `xml:"historyGraphColor"`
	UsableForArt      bool      `xml:"usableForArt"`
	MaxPerPawn        int64     `xml:"maxPerPawn"`
	IgnoreChance      float64   `xml:"ignoreChance"`
	ColonistOnly      bool      `xml:"colonistOnly"`
	ExpireDays        int64     `xml:"expireDays"`
}

func (t *TaleDef) Assign(*xml.Element) error {
	return nil
}

func (t *TaleDef) CountValidatedField() int {
	if t.FieldValidated == nil {
		return 0
	}
	return len(t.FieldValidated)
}

func (t *TaleDef) Equal(*TaleDef) bool {
	return false
}

func (t *TaleDef) GetAttributes() attributes.Attributes {
	return t.Attr
}

func (t *TaleDef) GetPath() string {
	return ""
}

func (t *TaleDef) Greater(*TaleDef) bool {
	return false
}

func (t *TaleDef) IsValidField(field string) bool {
	return t.FieldValidated[field]
}

func (t *TaleDef) Less(*TaleDef) bool {
	return false
}

func (t *TaleDef) SetAttributes(attr attributes.Attributes) {
	t.Attr = attr
	return
}

func (t *TaleDef) Val() *TaleDef {
	return nil
}

func (t *TaleDef) ValidateField(field string) {
	if t.FieldValidated == nil {
		t.FieldValidated = make(map[string]bool)
	}
	t.FieldValidated[field] = true
	return
}
