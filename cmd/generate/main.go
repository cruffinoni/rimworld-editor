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
	flag.StringVar(&path, "path", "./generated", "Path to the save game file")
	flag.Parse()
	if path == "" {
		log.Println("no path specified")
		flag.Usage()
		return
	}
	fo, err = file.Open(path)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(fo.XML.Root.XMLPath(), fo.XML.Root.Child.GetName())
	//log.Printf("Root: %v & Children: %v", fo.XML.Root.GetName(), fo.XML.Root.Child.GetName())
	root := generator.GenerateGoFiles(fo.XML.Root)
	if err = root.WriteGoFile("./generated"); err != nil {
		log.Fatal(err)
	}
}
