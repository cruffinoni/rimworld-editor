// Code generated by rimworld-editor. DO NOT EDIT.

package pawncolumndefs

import (
	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/attributes"
)

type PawnColumnDef struct {
	Attr           attributes.Attributes
	FieldValidated map[string]bool

	DefName                               string `xml:"defName"`
	WorkerClass                           string `xml:"workerClass"`
	WidthPriority                         int64  `xml:"widthPriority"`
	ShowIcon                              bool   `xml:"showIcon"`
	UseLabelShort                         bool   `xml:"useLabelShort"`
	Label                                 string `xml:"label"`
	Sortable                              bool   `xml:"sortable"`
	Width                                 int64  `xml:"width"`
	HeaderTip                             string `xml:"headerTip"`
	HeaderIcon                            string `xml:"headerIcon"`
	Paintable                             bool   `xml:"paintable"`
	Gap                                   int64  `xml:"gap"`
	IgnoreWhenCalculatingOptimalTableSize bool   `xml:"ignoreWhenCalculatingOptimalTableSize"`
}

func (p *PawnColumnDef) Assign(*xml.Element) error {
	return nil
}

func (p *PawnColumnDef) CountValidatedField() int {
	if p.FieldValidated == nil {
		return 0
	}
	return len(p.FieldValidated)
}

func (p *PawnColumnDef) Equal(*PawnColumnDef) bool {
	return false
}

func (p *PawnColumnDef) GetAttributes() attributes.Attributes {
	return p.Attr
}

func (p *PawnColumnDef) GetPath() string {
	return ""
}

func (p *PawnColumnDef) Greater(*PawnColumnDef) bool {
	return false
}

func (p *PawnColumnDef) IsValidField(field string) bool {
	return p.FieldValidated[field]
}

func (p *PawnColumnDef) Less(*PawnColumnDef) bool {
	return false
}

func (p *PawnColumnDef) SetAttributes(attr attributes.Attributes) {
	p.Attr = attr
	return
}

func (p *PawnColumnDef) Val() *PawnColumnDef {
	return nil
}

func (p *PawnColumnDef) ValidateField(field string) {
	if p.FieldValidated == nil {
		p.FieldValidated = make(map[string]bool)
	}
	p.FieldValidated[field] = true
	return
}
