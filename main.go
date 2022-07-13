package main

import (
	"github.com/cruffinoni/rimworld-editor/editor"
	"github.com/cruffinoni/rimworld-editor/xml/unmarshal"
	"log"
	"os"
)

func main() {
	f, err := editor.Open("test/part.rws")
	if err != nil {
		log.Fatal(err)
	}

	var m editor.Savegame
	err = unmarshal.Element(f.XML.Root, &m)
	if err != nil {
		log.Fatal(err)
	}
	if err = m.SaveXML("test/part_save.xml"); err != nil {
		log.Fatal(err)
	}
}

func outputFileTree() {
	f, err := editor.Open("test/alone.rws")
	if err != nil {
		log.Fatal(err)
	}
	os.WriteFile("test/alone.xml", []byte(f.XML.ToXML()), 0644)
}
