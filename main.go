package main

import (
	"github.com/cruffinoni/rimworld-editor/editor"
	"github.com/cruffinoni/rimworld-editor/generator"
	"log"
	"time"
)

func main() {
	log.Println("Starting...")
	t := time.Now()
	f, err := editor.Open("test/huge.rws")
	if err != nil {
		log.Fatal(err)
	}

	//log.Println("Parsing...")
	//var m editor.Savegame
	//err = unmarshal.Element(f.XML.Root, &m)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.Println("Saving...")
	//if err = m.SaveXML("test/part_save.xml"); err != nil {
	//	log.Fatal(err)
	//}
	root := generator.GenerateGoFiles(f.XML.Root)
	if err = root.WriteGoFile("generated"); err != nil {
		log.Fatal(err)
	}
	log.Printf("Done in %s", time.Since(t))
}
