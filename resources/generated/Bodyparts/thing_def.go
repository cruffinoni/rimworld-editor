// Code generated by rimworld-editor. DO NOT EDIT.

package bodyparts

import (
	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/attributes"
	"github.com/cruffinoni/rimworld-editor/xml/types"
)

type ThingDef struct {
	Attr           attributes.Attributes
	FieldValidated map[string]bool

	TechLevel                  string                     `xml:"techLevel"`
	ThingCategories            *types.Slice[string]       `xml:"thingCategories"`
	GraphicData                *GraphicData               `xml:"graphicData"`
	StatBases                  *StatBases                 `xml:"statBases"`
	ThingSetMakerTags          *types.Slice[string]       `xml:"thingSetMakerTags"`
	DefName                    string                     `xml:"defName"`
	Label                      string                     `xml:"label"`
	Description                string                     `xml:"description"`
	DescriptionHyperlinks      *DescriptionHyperlinks     `xml:"descriptionHyperlinks"`
	Category                   string                     `xml:"category"`
	UseHitPoints               bool                       `xml:"useHitPoints"`
	Selectable                 bool                       `xml:"selectable"`
	IsTechHediff               bool                       `xml:"isTechHediff"`
	PathCost                   int64                      `xml:"pathCost"`
	ThingClass                 string                     `xml:"thingClass"`
	DrawerType                 string                     `xml:"drawerType"`
	AltitudeLayer              string                     `xml:"altitudeLayer"`
	TickerType                 string                     `xml:"tickerType"`
	AlwaysHaulable             bool                       `xml:"alwaysHaulable"`
	AllowedArchonexusCount     int64                      `xml:"allowedArchonexusCount"`
	Comps                      *types.Slice[*Comps]       `xml:"comps"`
	TradeTags                  *types.Slice[string]       `xml:"tradeTags"`
	RecipeMaker                *RecipeMaker               `xml:"recipeMaker"`
	CostList                   *CostList                  `xml:"costList"`
	TechHediffsTags            *types.Slice[string]       `xml:"techHediffsTags"`
	SpawnThingOnRemoved        string                     `xml:"spawnThingOnRemoved"`
	AddedPartProps             *AddedPartProps            `xml:"addedPartProps"`
	LabelNoun                  string                     `xml:"labelNoun"`
	JobString                  string                     `xml:"jobString"`
	SkillRequirements          *SkillRequirements         `xml:"skillRequirements"`
	FixedIngredientFilter      *FixedIngredientFilter     `xml:"fixedIngredientFilter"`
	AppliedOnFixedBodyParts    *types.Slice[string]       `xml:"appliedOnFixedBodyParts"`
	Stages                     *types.Slice[*Stages]      `xml:"stages"`
	RemovesHediff              string                     `xml:"removesHediff"`
	DeathOnFailedSurgeryChance float64                    `xml:"deathOnFailedSurgeryChance"`
	Ingredients                *types.Slice[*Ingredients] `xml:"ingredients"`
	AddsHediff                 string                     `xml:"addsHediff"`
}

func (t *ThingDef) Assign(*xml.Element) error {
	return nil
}

func (t *ThingDef) CountValidatedField() int {
	if t.FieldValidated == nil {
		return 0
	}
	return len(t.FieldValidated)
}

func (t *ThingDef) Equal(*ThingDef) bool {
	return false
}

func (t *ThingDef) GetAttributes() attributes.Attributes {
	return t.Attr
}

func (t *ThingDef) GetPath() string {
	return ""
}

func (t *ThingDef) Greater(*ThingDef) bool {
	return false
}

func (t *ThingDef) IsValidField(field string) bool {
	return t.FieldValidated[field]
}

func (t *ThingDef) Less(*ThingDef) bool {
	return false
}

func (t *ThingDef) SetAttributes(attr attributes.Attributes) {
	t.Attr = attr
	return
}

func (t *ThingDef) Val() *ThingDef {
	return nil
}

func (t *ThingDef) ValidateField(field string) {
	if t.FieldValidated == nil {
		t.FieldValidated = make(map[string]bool)
	}
	t.FieldValidated[field] = true
	return
}
