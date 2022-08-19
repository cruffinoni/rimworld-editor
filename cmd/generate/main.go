package main

import (
	"flag"
	"github.com/cruffinoni/rimworld-editor/file"
	"github.com/cruffinoni/rimworld-editor/generator"
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
	fo, err = file.Open(path)
	if err != nil {
		log.Println(err)
		return
	}
	root := generator.GenerateGoFiles(fo.XML.Root)
	if err = root.WriteGoFile("./generated"); err != nil {
		log.Fatal(err)
	}
}
