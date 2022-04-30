package game

import (
	"github.com/cruffinoni/rimworld-editor/utils"
)

type Scenario struct {
	Name        string `xml:"name"`
	Summary     string `xml:"summary"`
	Description string `xml:"description"`
	//PlayerFaction utils.GenericSimpleFormat `xml:"playerFaction"`
	Parts utils.GenericFormatByMap `xml:"parts"` // TODO: parts>li / Implement XML list reader
}

type Game struct {
	CurrentMapIndex int `xml:"currentMapIndex"`
	//Info            utils.GenericSimpleFormat `xml:"info"`
	//Rules           utils.GenericSimpleFormat `xml:"rules"`
	Scenario Scenario `xml:"scenario"`
}
