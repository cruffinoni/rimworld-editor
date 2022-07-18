package main

import (
	"github.com/cruffinoni/rimworld-editor/editor"
	"github.com/cruffinoni/rimworld-editor/xml/unmarshal"
	"log"
)

func main() {
	log.Println("Starting...")
	f, err := editor.Open("test/part.rws")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Parsing...")

	var m editor.Savegame
	err = unmarshal.Element(f.XML.Root, &m)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Saving...")
	if err = m.SaveXML("test/part_save.xml"); err != nil {
		log.Fatal(err)
	}
	log.Println("Done!")
}
