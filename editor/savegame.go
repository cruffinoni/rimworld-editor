package editor

import (
	"github.com/cruffinoni/rimworld-editor/editor/game"
	"github.com/cruffinoni/rimworld-editor/editor/meta"
)

type Savegame struct {
	//Xml      string `xml:"savegame"`
	//Version  string `xml:"version,attr"`
	//Encoding string `xml:"encoding,attr"`
	Meta *meta.Meta `xml:"meta"`
	Game *game.Game `xml:"game"`
}
