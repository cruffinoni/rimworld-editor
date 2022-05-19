package editor

import (
	"github.com/cruffinoni/rimworld-editor/xml"
)

type Savegame struct {
	//Xml      string `xml:"savegame"`
	//Version  string `xml:"version,attr"`
	//Encoding string `xml:"encoding,attr"`
	Meta *xml.Tree `xml:"meta"`
	Game *xml.Tree `xml:"game"`
}
