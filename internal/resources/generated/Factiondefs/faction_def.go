// Code generated by rimworld-editor. DO NOT EDIT.

package factiondefs

import (
	"github.com/cruffinoni/rimworld-editor/internal/xml"
	"github.com/cruffinoni/rimworld-editor/internal/xml/attributes"
	"github.com/cruffinoni/rimworld-editor/internal/xml/types"
)

type FactionDef struct {
	Attr           attributes.Attributes
	FieldValidated map[string]bool

	DefName                             string                               `xml:"defName"`
	Label                               string                               `xml:"label"`
	CategoryTag                         string                               `xml:"categoryTag"`
	Description                         string                               `xml:"description"`
	PawnSingular                        string                               `xml:"pawnSingular"`
	PawnsPlural                         string                               `xml:"pawnsPlural"`
	SettlementGenerationWeight          int64                                `xml:"settlementGenerationWeight"`
	RequiredCountAtGameStart            int64                                `xml:"requiredCountAtGameStart"`
	MaxCountAtGameStart                 int64                                `xml:"maxCountAtGameStart"`
	CanSiege                            bool                                 `xml:"canSiege"`
	CanStageAttacks                     bool                                 `xml:"canStageAttacks"`
	LeaderTitle                         string                               `xml:"leaderTitle"`
	RoyalFavorLabel                     string                               `xml:"royalFavorLabel"`
	RoyalFavorIconPath                  string                               `xml:"royalFavorIconPath"`
	LeaderForceGenerateNewPawn          bool                                 `xml:"leaderForceGenerateNewPawn"`
	FactionIconPath                     string                               `xml:"factionIconPath"`
	FactionNameMaker                    string                               `xml:"factionNameMaker"`
	SettlementNameMaker                 string                               `xml:"settlementNameMaker"`
	AllowedCultures                     *types.Slice[string]                 `xml:"allowedCultures"`
	RequiredMemes                       *types.Slice[string]                 `xml:"requiredMemes"`
	AllowedMemes                        *types.Slice[string]                 `xml:"allowedMemes"`
	DisallowedPrecepts                  *types.Slice[string]                 `xml:"disallowedPrecepts"`
	StructureMemeWeights                *StructureMemeWeights                `xml:"structureMemeWeights"`
	XenotypeSet                         *XenotypeSet                         `xml:"xenotypeSet"`
	PermanentEnemyToEveryoneExcept      *types.Slice[string]                 `xml:"permanentEnemyToEveryoneExcept"`
	TechLevel                           string                               `xml:"techLevel"`
	BackstoryFilters                    *types.Slice[*BackstoryFilters]      `xml:"backstoryFilters"`
	ApparelStuffFilter                  *ApparelStuffFilter                  `xml:"apparelStuffFilter"`
	AllowedArrivalTemperatureRange      string                               `xml:"allowedArrivalTemperatureRange"`
	SettlementTexturePath               string                               `xml:"settlementTexturePath"`
	ColorSpectrum                       *types.Slice[string]                 `xml:"colorSpectrum"`
	FixedLeaderKinds                    *types.Slice[string]                 `xml:"fixedLeaderKinds"`
	RoyalTitleTags                      *types.Slice[string]                 `xml:"royalTitleTags"`
	BaseTraderKinds                     *types.Slice[string]                 `xml:"baseTraderKinds"`
	CaravanTraderKinds                  *types.Slice[string]                 `xml:"caravanTraderKinds"`
	RoyalTitleInheritanceWorkerClass    string                               `xml:"royalTitleInheritanceWorkerClass"`
	RoyalTitleInheritanceRelations      *types.Slice[string]                 `xml:"royalTitleInheritanceRelations"`
	RaidCommonalityFromPointsCurve      *RaidCommonalityFromPointsCurve      `xml:"raidCommonalityFromPointsCurve"`
	RaidLootMaker                       string                               `xml:"raidLootMaker"`
	MaxPawnCostPerTotalPointsCurve      *MaxPawnCostPerTotalPointsCurve      `xml:"maxPawnCostPerTotalPointsCurve"`
	PawnGroupMakers                     *types.Slice[*PawnGroupMakers]       `xml:"pawnGroupMakers"`
	MaxConfigurableAtWorldCreation      int64                                `xml:"maxConfigurableAtWorldCreation"`
	ConfigurationListOrderPriority      int64                                `xml:"configurationListOrderPriority"`
	DisallowedRaidAgeRestrictions       *types.Slice[string]                 `xml:"disallowedRaidAgeRestrictions"`
	BasicMemberKind                     string                               `xml:"basicMemberKind"`
	DisallowedMemes                     *types.Slice[string]                 `xml:"disallowedMemes"`
	Hidden                              bool                                 `xml:"hidden"`
	GenerateNewLeaderFromMapMembersOnly bool                                 `xml:"generateNewLeaderFromMapMembersOnly"`
	MustStartOneEnemy                   bool                                 `xml:"mustStartOneEnemy"`
	CanMakeRandomly                     bool                                 `xml:"canMakeRandomly"`
	RaidLootValueFromPointsCurve        *RaidLootValueFromPointsCurve        `xml:"raidLootValueFromPointsCurve"`
	NaturalEnemy                        bool                                 `xml:"naturalEnemy"`
	ListOrderPriority                   int64                                `xml:"listOrderPriority"`
	PermanentEnemy                      bool                                 `xml:"permanentEnemy"`
	VisitorTraderKinds                  *types.Slice[string]                 `xml:"visitorTraderKinds"`
	ClassicIdeo                         bool                                 `xml:"classicIdeo"`
	FixedName                           string                               `xml:"fixedName"`
	Root                                *Root                                `xml:"root"`
	HumanlikeFaction                    bool                                 `xml:"humanlikeFaction"`
	StartingCountAtWorldCreation        int64                                `xml:"startingCountAtWorldCreation"`
	MinSettlementTemperatureChanceCurve *MinSettlementTemperatureChanceCurve `xml:"minSettlementTemperatureChanceCurve"`
	HostileToFactionlessHumanlikes      bool                                 `xml:"hostileToFactionlessHumanlikes"`
	PlayerInitialSettlementNameMaker    string                               `xml:"playerInitialSettlementNameMaker"`
	AutoFlee                            bool                                 `xml:"autoFlee"`
	EarliestRaidDays                    int64                                `xml:"earliestRaidDays"`
	DropPodActive                       string                               `xml:"dropPodActive"`
	CanUseAvoidGrid                     bool                                 `xml:"canUseAvoidGrid"`
	DropPodIncoming                     string                               `xml:"dropPodIncoming"`
	IsPlayer                            bool                                 `xml:"isPlayer"`
	StartingTechprintsResearchTags      *types.Slice[string]                 `xml:"startingTechprintsResearchTags"`
	ForageabilityFactor                 float64                              `xml:"forageabilityFactor"`
	StartingResearchTags                *types.Slice[string]                 `xml:"startingResearchTags"`
	DisplayInFactionSelection           bool                                 `xml:"displayInFactionSelection"`
	RescueesCanJoin                     bool                                 `xml:"rescueesCanJoin"`
	RecipePrerequisiteTags              *types.Slice[string]                 `xml:"recipePrerequisiteTags"`
}

func (f *FactionDef) Assign(*xml.Element) error {
	return nil
}

func (f *FactionDef) CountValidatedField() int {
	if f.FieldValidated == nil {
		return 0
	}
	return len(f.FieldValidated)
}

func (f *FactionDef) Equal(*FactionDef) bool {
	return false
}

func (f *FactionDef) GetAttributes() attributes.Attributes {
	return f.Attr
}

func (f *FactionDef) GetPath() string {
	return ""
}

func (f *FactionDef) Greater(*FactionDef) bool {
	return false
}

func (f *FactionDef) IsValidField(field string) bool {
	return f.FieldValidated[field]
}

func (f *FactionDef) Less(*FactionDef) bool {
	return false
}

func (f *FactionDef) SetAttributes(attr attributes.Attributes) {
	f.Attr = attr
	return
}

func (f *FactionDef) Val() *FactionDef {
	return nil
}

func (f *FactionDef) ValidateField(field string) {
	if f.FieldValidated == nil {
		f.FieldValidated = make(map[string]bool)
	}
	f.FieldValidated[field] = true
	return
}
