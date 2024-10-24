// Code generated by rimworld-editor. DO NOT EDIT.

package genedefs

import (
	"github.com/cruffinoni/rimworld-editor/internal/xml"
	"github.com/cruffinoni/rimworld-editor/internal/xml/attributes"
	"github.com/cruffinoni/rimworld-editor/internal/xml/types"
)

type GeneCategoryDef struct {
	Attr           attributes.Attributes
	FieldValidated map[string]bool

	DefName                       string               `xml:"defName"`
	Label                         string               `xml:"label"`
	DisplayPriorityInXenotype     int64                `xml:"displayPriorityInXenotype"`
	DisplayPriorityInGenepack     int64                `xml:"displayPriorityInGenepack"`
	Description                   string               `xml:"description"`
	BiostatCpx                    int64                `xml:"biostatCpx"`
	EndogeneCategory              string               `xml:"endogeneCategory"`
	IconPath                      string               `xml:"iconPath"`
	CanGenerateInGeneSet          bool                 `xml:"canGenerateInGeneSet"`
	DisplayCategory               string               `xml:"displayCategory"`
	RandomBrightnessFactor        float64              `xml:"randomBrightnessFactor"`
	PassOnDirectly                bool                 `xml:"passOnDirectly"`
	ExclusionTags                 *types.Slice[string] `xml:"exclusionTags"`
	HairColorOverride             string               `xml:"hairColorOverride"`
	SelectionWeight               float64              `xml:"selectionWeight"`
	DisplayOrderInCategory        int64                `xml:"displayOrderInCategory"`
	SelectionWeightFactorDarkSkin int64                `xml:"selectionWeightFactorDarkSkin"`
	SkinColorBase                 string               `xml:"skinColorBase"`
	MinMelanin                    float64              `xml:"minMelanin"`
}

func (g *GeneCategoryDef) Assign(*xml.Element) error {
	return nil
}

func (g *GeneCategoryDef) CountValidatedField() int {
	if g.FieldValidated == nil {
		return 0
	}
	return len(g.FieldValidated)
}

func (g *GeneCategoryDef) Equal(*GeneCategoryDef) bool {
	return false
}

func (g *GeneCategoryDef) GetAttributes() attributes.Attributes {
	return g.Attr
}

func (g *GeneCategoryDef) GetPath() string {
	return ""
}

func (g *GeneCategoryDef) Greater(*GeneCategoryDef) bool {
	return false
}

func (g *GeneCategoryDef) IsValidField(field string) bool {
	return g.FieldValidated[field]
}

func (g *GeneCategoryDef) Less(*GeneCategoryDef) bool {
	return false
}

func (g *GeneCategoryDef) SetAttributes(attr attributes.Attributes) {
	g.Attr = attr
	return
}

func (g *GeneCategoryDef) Val() *GeneCategoryDef {
	return nil
}

func (g *GeneCategoryDef) ValidateField(field string) {
	if g.FieldValidated == nil {
		g.FieldValidated = make(map[string]bool)
	}
	g.FieldValidated[field] = true
	return
}
