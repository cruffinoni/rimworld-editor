package editor

import (
	"github.com/cruffinoni/rimworld-editor/editor/game"
	"github.com/cruffinoni/rimworld-editor/xml"
)

type Savegame struct {
	//Xml      string `xml:"savegame"`
	//Version  string `xml:"version,attr"`
	//Encoding string `xml:"encoding,attr"`
	Meta *xml.Discover `xml:"meta"`
	Game *game.Game    `xml:"game"`
}
