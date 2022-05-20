package editor

import (
	"github.com/cruffinoni/rimworld-editor/xml"
)

type Savegame struct {
	//Xml      string `XML:"savegame"`
	//Version  string `XML:"version,attr"`
	//Encoding string `XML:"encoding,attr"`
	Meta *xml.Tree `XML:"meta"`
	Game *xml.Tree `XML:"game"`
}
