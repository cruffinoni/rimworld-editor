// Code generated by rimworld-editor. DO NOT EDIT.

package pawnkinddefs_humanlikes

import (
	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/attributes"
	"github.com/cruffinoni/rimworld-editor/xml/types"
)

type InventoryOptions struct {
	Attr           attributes.Attributes
	FieldValidated map[string]bool

	SkipChance          float64                            `xml:"skipChance"`
	SubOptionsChooseOne *types.Slice[*SubOptionsChooseOne] `xml:"subOptionsChooseOne"`
}

func (i *InventoryOptions) Assign(*xml.Element) error {
	return nil
}

func (i *InventoryOptions) CountValidatedField() int {
	if i.FieldValidated == nil {
		return 0
	}
	return len(i.FieldValidated)
}

func (i *InventoryOptions) Equal(*InventoryOptions) bool {
	return false
}

func (i *InventoryOptions) GetAttributes() attributes.Attributes {
	return i.Attr
}

func (i *InventoryOptions) GetPath() string {
	return ""
}

func (i *InventoryOptions) Greater(*InventoryOptions) bool {
	return false
}

func (i *InventoryOptions) IsValidField(field string) bool {
	return i.FieldValidated[field]
}

func (i *InventoryOptions) Less(*InventoryOptions) bool {
	return false
}

func (i *InventoryOptions) SetAttributes(attr attributes.Attributes) {
	i.Attr = attr
	return
}

func (i *InventoryOptions) Val() *InventoryOptions {
	return nil
}

func (i *InventoryOptions) ValidateField(field string) {
	if i.FieldValidated == nil {
		i.FieldValidated = make(map[string]bool)
	}
	i.FieldValidated[field] = true
	return
}
