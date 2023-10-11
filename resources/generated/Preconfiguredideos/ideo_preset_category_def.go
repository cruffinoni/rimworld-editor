// Code generated by rimworld-editor. DO NOT EDIT.

package preconfiguredideos

import (
	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/attributes"
	"github.com/cruffinoni/rimworld-editor/xml/types"
)

type IdeoPresetCategoryDef struct {
	Attr           attributes.Attributes
	FieldValidated map[string]bool

	DefName     string               `xml:"defName"`
	Label       string               `xml:"label"`
	Description string               `xml:"description"`
	GroupLabel  string               `xml:"groupLabel"`
	CategoryDef string               `xml:"categoryDef"`
	IconPath    string               `xml:"iconPath"`
	ClassicPlus bool                 `xml:"classicPlus"`
	Memes       *types.Slice[string] `xml:"memes"`
}

func (i *IdeoPresetCategoryDef) Assign(*xml.Element) error {
	return nil
}

func (i *IdeoPresetCategoryDef) CountValidatedField() int {
	if i.FieldValidated == nil {
		return 0
	}
	return len(i.FieldValidated)
}

func (i *IdeoPresetCategoryDef) Equal(*IdeoPresetCategoryDef) bool {
	return false
}

func (i *IdeoPresetCategoryDef) GetAttributes() attributes.Attributes {
	return i.Attr
}

func (i *IdeoPresetCategoryDef) GetPath() string {
	return ""
}

func (i *IdeoPresetCategoryDef) Greater(*IdeoPresetCategoryDef) bool {
	return false
}

func (i *IdeoPresetCategoryDef) IsValidField(field string) bool {
	return i.FieldValidated[field]
}

func (i *IdeoPresetCategoryDef) Less(*IdeoPresetCategoryDef) bool {
	return false
}

func (i *IdeoPresetCategoryDef) SetAttributes(attr attributes.Attributes) {
	i.Attr = attr
	return
}

func (i *IdeoPresetCategoryDef) Val() *IdeoPresetCategoryDef {
	return nil
}

func (i *IdeoPresetCategoryDef) ValidateField(field string) {
	if i.FieldValidated == nil {
		i.FieldValidated = make(map[string]bool)
	}
	i.FieldValidated[field] = true
	return
}