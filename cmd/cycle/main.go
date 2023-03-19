package main

import (
	"flag"
	"log"
	"strconv"
	"time"

	"github.com/cruffinoni/rimworld-editor/file"
	"github.com/cruffinoni/rimworld-editor/generated"
	"github.com/cruffinoni/rimworld-editor/generator"
	"github.com/cruffinoni/rimworld-editor/generator/files"
	"github.com/cruffinoni/rimworld-editor/xml/saver/xmlFile"
	"github.com/cruffinoni/rimworld-editor/xml/unmarshal"
)

func main() {
	var (
		fo       *file.Opening
		err      error
		path     string
		fileName string
	)
	flag.StringVar(&path, "path", "", "Path to the save game file")
	flag.StringVar(&fileName, "fileName", "CUSTOM_FILE", "File name for the generated XML")
	flag.Parse()
	if fileName == "CUSTOM_FILE" {
		fileName = "C_" + strconv.FormatInt(time.Now().Unix(), 10)
	}
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
	log.Print("Generating go files to './generated")
	root := generator.GenerateGoFiles(fo.XML.Root, true)
	if err = files.WriteGoFile("./generated", root, true, nil); err != nil {
		log.Fatal(err)
	}
	save := &generated.Savegame{}
	log.Println("Unmarshalling XML...")
	if err := unmarshal.Element(fo.XML.Root.Child, save); err != nil {
		log.Fatal(err)
	}
	save.ValidateField("Savegame")
	log.Print("Generating XML file to folder")
	buffer, err := xmlFile.SaveWithBuffer(save)
	if err != nil {
		log.Panic(err)
	}
	if err := buffer.ToFile("generated/" + fileName + ".rws"); err != nil {
		log.Panic(err)
	}
}
