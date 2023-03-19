package main

import (
	"flag"
	"log"
	"time"

	"github.com/briandowns/spinner"

	"github.com/cruffinoni/rimworld-editor/file"
	"github.com/cruffinoni/rimworld-editor/generator"
	"github.com/cruffinoni/rimworld-editor/generator/files"
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
	s := spinner.New(spinner.CharSets[35], 100*time.Millisecond)
	log.Printf("Opening and decoding XML file from %s", path)
	s.FinalMSG = "XML file decoded successfully\n"
	s.Start()
	fo, err = file.Open(path)
	if err != nil {
		log.Fatal(err)
		return
	}
	s.Stop()
	//s.Prefix = "Generating go files to './generated'... "
	s.FinalMSG = "Go files successfully generated\n"
	//s.Start()
	root := generator.GenerateGoFiles(fo.XML.Root, true)
	if err = files.WriteGoFile("./generated", root, true, nil); err != nil {
		log.Fatal(err)
	}
	//s.Stop()
}
