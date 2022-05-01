package main

import (
	"github.com/cruffinoni/rimworld-editor/editor"
	"log"
)

func main() {
	println("Hello, world.")
	if f, err := editor.Open("TEST.rws"); err != nil {
		panic(err)
	} else {
		//log.Printf("'%+#v'", f)
		for _, v := range f.GetXML().Game.Scenario.Parts.Data["ScenPart_ConfigPage_ConfigureStartingPawns"] {
			log.Printf("=> %+#v", v)
		}
		//log.Printf("-> %v", f.GetXML().Meta.GetModById("ludeon.rimworld"))
	}
}
