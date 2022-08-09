package main

import (
	"flag"
	"github.com/cruffinoni/rimworld-editor/editor_old"
	"github.com/cruffinoni/rimworld-editor/generator"
	"log"
)

func main() {
	var (
		file *editor_old.FileOpening
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
	file, err = editor_old.Open(path)
	root := generator.GenerateGoFiles(file.XML.Root)
	if err = root.WriteGoFile("./generated"); err != nil {
		log.Fatal(err)
	}
}
