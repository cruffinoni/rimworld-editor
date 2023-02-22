package main

import (
	"flag"
	"log"
	"math"
	"runtime/debug"

	"github.com/cruffinoni/rimworld-editor/file"
	"github.com/cruffinoni/rimworld-editor/generated"
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
	// 1000000000-byte
	// 1000000000-byte limit
	debug.SetMaxStack(1 << 30)
	debug.SetMemoryLimit(int64(5 * math.Pow(2.0, 30)))
	flag.StringVar(&path, "path", "", "Path to the save game file")
	flag.StringVar(&fileName, "fileName", "CUSTOM_FILE", "File name for the generated XML")
	flag.Parse()
	if path == "" {
		log.Println("no path specified")
		flag.Usage()
		return
	}
	//log.Printf("Opening and decoding XML file from %s", path)
	fo, err = file.Open(path)
	if err != nil {
		log.Fatal(err)
		return
	}
	save := &generated.Savegame{}
	//log.Println("Unmarshalling XML...")
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
