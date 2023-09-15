// Code generated by rimworld-editor. DO NOT EDIT.

package lifestagedefs

import (
	"github.com/cruffinoni/rimworld-editor/xml"
	"github.com/cruffinoni/rimworld-editor/xml/attributes"
	"github.com/cruffinoni/rimworld-editor/xml/types/embedded"
)

type LifeStageDef struct {
	Attr           attributes.Attributes
	FieldValidated map[string]bool

	DefName                         string                           `xml:"defName"`
	Label                           string                           `xml:"label"`
	Visible                         bool                             `xml:"visible"`
	HairStyleFilter                 *HairStyleFilter                 `xml:"hairStyleFilter"`
	BodySizeFactor                  float64                          `xml:"bodySizeFactor"`
	BodyWidth                       float64                          `xml:"bodyWidth"`
	BodyDrawOffset                  string                           `xml:"bodyDrawOffset"`
	VoxPitch                        float64                          `xml:"voxPitch"`
	HeadSizeFactor                  float64                          `xml:"headSizeFactor"`
	EyeSizeFactor                   float64                          `xml:"eyeSizeFactor"`
	FoodMaxFactor                   float64                          `xml:"foodMaxFactor"`
	HungerRateFactor                float64                          `xml:"hungerRateFactor"`
	HealthScaleFactor               float64                          `xml:"healthScaleFactor"`
	MarketValueFactor               float64                          `xml:"marketValueFactor"`
	MeleeDamageFactor               float64                          `xml:"meleeDamageFactor"`
	DevelopmentalStage              string                           `xml:"developmentalStage"`
	FallAsleepMaxThresholdOverride  float64                          `xml:"fallAsleepMaxThresholdOverride"`
	NaturalWakeThresholdOverride    float64                          `xml:"naturalWakeThresholdOverride"`
	AlwaysDowned                    bool                             `xml:"alwaysDowned"`
	Claimable                       bool                             `xml:"claimable"`
	InvoluntarySleepIsNegativeEvent bool                             `xml:"involuntarySleepIsNegativeEvent"`
	ThinkTreeMainOverride           *embedded.Type[string]           `xml:"thinkTreeMainOverride"`
	ThinkTreeConstantOverride       *embedded.Type[string]           `xml:"thinkTreeConstantOverride"`
	CanDoRandomMentalBreaks         bool                             `xml:"canDoRandomMentalBreaks"`
	CanSleepWhileHeld               bool                             `xml:"canSleepWhileHeld"`
	CanVoluntarilySleep             bool                             `xml:"canVoluntarilySleep"`
	CanSleepWhenStarving            bool                             `xml:"canSleepWhenStarving"`
	CanInitiateSocialInteraction    bool                             `xml:"canInitiateSocialInteraction"`
	CustomMoodTipString             string                           `xml:"customMoodTipString"`
	StatFactors                     *StatFactors                     `xml:"statFactors"`
	StatOffsets                     *StatOffsets                     `xml:"statOffsets"`
	InvoluntarySleepMtbdaysFromRest *InvoluntarySleepMtbdaysFromRest `xml:"involuntarySleepMTBDaysFromRest"`
	WorkerClass                     *embedded.Type[string]           `xml:"workerClass"`
	EquipmentDrawDistanceFactor     float64                          `xml:"equipmentDrawDistanceFactor"`
	SittingOffset                   float64                          `xml:"sittingOffset"`
	Adjective                       string                           `xml:"adjective"`
	Reproductive                    bool                             `xml:"reproductive"`
	VoxVolume                       float64                          `xml:"voxVolume"`
	Milkable                        bool                             `xml:"milkable"`
	Shearable                       bool                             `xml:"shearable"`
	CaravanRideable                 bool                             `xml:"caravanRideable"`
}

func (l *LifeStageDef) Assign(*xml.Element) error {
	return nil
}

func (l *LifeStageDef) CountValidatedField() int {
	if l.FieldValidated == nil {
		return 0
	}
	return len(l.FieldValidated)
}

func (l *LifeStageDef) Equal(*LifeStageDef) bool {
	return false
}

func (l *LifeStageDef) GetAttributes() attributes.Attributes {
	return l.Attr
}

func (l *LifeStageDef) GetPath() string {
	return ""
}

func (l *LifeStageDef) Greater(*LifeStageDef) bool {
	return false
}

func (l *LifeStageDef) IsValidField(field string) bool {
	return l.FieldValidated[field]
}

func (l *LifeStageDef) Less(*LifeStageDef) bool {
	return false
}

func (l *LifeStageDef) SetAttributes(attr attributes.Attributes) {
	l.Attr = attr
	return
}

func (l *LifeStageDef) Val() *LifeStageDef {
	return nil
}

func (l *LifeStageDef) ValidateField(field string) {
	if l.FieldValidated == nil {
		l.FieldValidated = make(map[string]bool)
	}
	l.FieldValidated[field] = true
	return
}
