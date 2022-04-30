package game

import (
	"github.com/cruffinoni/rimworld-editor/utils"
)

type Game struct {
	CurrentMapIndex int                         `xml:"currentMapIndex"`
	Info            []utils.GenericSimpleFormat `xml:"info"`
	//Rules           utils.GenericSimpleFormat `xml:"rules"`
}
