package main

import (
	"flag"
	"github.com/cruffinoni/rimworld-editor/file"
	"github.com/cruffinoni/rimworld-editor/generated"
	"github.com/cruffinoni/rimworld-editor/xml/saver/xmlFile"
	"github.com/cruffinoni/rimworld-editor/xml/unmarshal"
	"log"
)

func main() {
	var (
		fo   *file.Opening
		err  error
		path string
	)
	flag.StringVar(&path, "path", "", "Path to the save game file")
	flag.Parse()
	if path == "" {
		log.Println("no path specified")
		flag.Usage()
		return
	}
	log.Printf("Opening and decoding XML file from %s", path)
	fo, err = file.Open(path)
	if err != nil {
		log.Fatal(err)
		return
	}
	save := &generated.Savegame{}
	log.Println("Unmarshalling XML...")
	if err := unmarshal.Element(fo.XML.Root, save); err != nil {
		log.Fatal(err)
	}
	log.Print("Generating XML file to folder")
	buffer, err := xmlFile.SaveWithBuffer(save)
	if err != nil {
		log.Panic(err)
	}
	if err := buffer.ToFile("CUSTOM_FILE.rws"); err != nil {
		log.Panic(err)
	}
}
