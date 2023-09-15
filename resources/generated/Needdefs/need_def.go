// Code generated by rimworld-editor. DO NOT EDIT.

package needdefs

import (
	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/attributes"
	"github.com/cruffinoni/rimworld-editor/xml/types"
)

type NeedDef struct {
	Attr           attributes.Attributes
	FieldValidated map[string]bool

	DefName                  string               `xml:"defName"`
	NeedClass                string               `xml:"needClass"`
	Label                    string               `xml:"label"`
	Description              string               `xml:"description"`
	ShowOnNeedList           bool                 `xml:"showOnNeedList"`
	MinIntelligence          string               `xml:"minIntelligence"`
	BaseLevel                float64              `xml:"baseLevel"`
	SeekerRisePerHour        float64              `xml:"seekerRisePerHour"`
	SeekerFallPerHour        float64              `xml:"seekerFallPerHour"`
	ListPriority             int64                `xml:"listPriority"`
	Major                    bool                 `xml:"major"`
	FreezeWhileSleeping      bool                 `xml:"freezeWhileSleeping"`
	ShowForCaravanMembers    bool                 `xml:"showForCaravanMembers"`
	DevelopmentalStageFilter string               `xml:"developmentalStageFilter"`
	ShowUnitTicks            bool                 `xml:"showUnitTicks"`
	ColonistsOnly            bool                 `xml:"colonistsOnly"`
	NeverOnPrisoner          bool                 `xml:"neverOnPrisoner"`
	NeverOnSlave             bool                 `xml:"neverOnSlave"`
	FreezeInMentalState      bool                 `xml:"freezeInMentalState"`
	ColonistAndPrisonersOnly bool                 `xml:"colonistAndPrisonersOnly"`
	NullifyingPrecepts       *types.Slice[string] `xml:"nullifyingPrecepts"`
	FallPerDay               float64              `xml:"fallPerDay"`
	SlavesOnly               bool                 `xml:"slavesOnly"`
}

func (n *NeedDef) Assign(*xml.Element) error {
	return nil
}

func (n *NeedDef) CountValidatedField() int {
	if n.FieldValidated == nil {
		return 0
	}
	return len(n.FieldValidated)
}

func (n *NeedDef) Equal(*NeedDef) bool {
	return false
}

func (n *NeedDef) GetAttributes() attributes.Attributes {
	return n.Attr
}

func (n *NeedDef) GetPath() string {
	return ""
}

func (n *NeedDef) Greater(*NeedDef) bool {
	return false
}

func (n *NeedDef) IsValidField(field string) bool {
	return n.FieldValidated[field]
}

func (n *NeedDef) Less(*NeedDef) bool {
	return false
}

func (n *NeedDef) SetAttributes(attr attributes.Attributes) {
	n.Attr = attr
	return
}

func (n *NeedDef) Val() *NeedDef {
	return nil
}

func (n *NeedDef) ValidateField(field string) {
	if n.FieldValidated == nil {
		n.FieldValidated = make(map[string]bool)
	}
	n.FieldValidated[field] = true
	return
}
