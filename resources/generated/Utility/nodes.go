// Code generated by rimworld-editor. DO NOT EDIT.

package utility

import (
	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/attributes"
	"github.com/cruffinoni/rimworld-editor/xml/types"
)

type Nodes struct {
	Attr           attributes.Attributes
	FieldValidated map[string]bool

	Name     string               `xml:"name"`
	ElseNode *ElseNode            `xml:"elseNode"`
	Nodes    *types.Slice[*Nodes] `xml:"nodes"`
}

func (n *Nodes) Assign(*xml.Element) error {
	return nil
}

func (n *Nodes) CountValidatedField() int {
	if n.FieldValidated == nil {
		return 0
	}
	return len(n.FieldValidated)
}

func (n *Nodes) Equal(*Nodes) bool {
	return false
}

func (n *Nodes) GetAttributes() attributes.Attributes {
	return n.Attr
}

func (n *Nodes) GetPath() string {
	return ""
}

func (n *Nodes) Greater(*Nodes) bool {
	return false
}

func (n *Nodes) IsValidField(field string) bool {
	return n.FieldValidated[field]
}

func (n *Nodes) Less(*Nodes) bool {
	return false
}

func (n *Nodes) SetAttributes(attr attributes.Attributes) {
	n.Attr = attr
	return
}

func (n *Nodes) Val() *Nodes {
	return nil
}

func (n *Nodes) ValidateField(field string) {
	if n.FieldValidated == nil {
		n.FieldValidated = make(map[string]bool)
	}
	n.FieldValidated[field] = true
	return
}