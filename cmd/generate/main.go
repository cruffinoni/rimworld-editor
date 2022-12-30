package main

import (
	"flag"
	"log"

	"github.com/cruffinoni/rimworld-editor/file"
	"github.com/cruffinoni/rimworld-editor/generator"
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
	log.Print("Generating go files to './generated")
	root := generator.GenerateGoFiles(fo.XML.Root)
	if err = root.WriteGoFile("./generated"); err != nil {
		log.Fatal(err)
	}
}
