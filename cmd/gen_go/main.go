package main

import (
	"flag"
	"log"
	"time"

	"github.com/briandowns/spinner"
	"github.com/cruffinoni/printer"

	"github.com/cruffinoni/rimworld-editor/internal/file"
	"github.com/cruffinoni/rimworld-editor/internal/generator"
	"github.com/cruffinoni/rimworld-editor/internal/generator/files"
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
		printer.Debugf("no path specified")
		flag.Usage()
		return
	}
	s := spinner.New(spinner.CharSets[35], 100*time.Millisecond)
	printer.Debugf("Opening and decoding XML file from %s", path)
	s.FinalMSG = "XML file decoded successfully\n"
	s.Start()
	fo, err = file.Open(path)
	if err != nil {
		log.Fatal(err)
		return
	}
	s.Stop()
	// s.Prefix = "Generating go files to './generated'... "
	s.FinalMSG = "Go files successfully generated\n"
	// s.Start()
	root := generator.GenerateGoFiles(fo.XML.Root, true)
	if err = files.DefaultGoWriter.WriteGoFile("./generated", root); err != nil {
		log.Fatal(err)
	}
	// s.Stop()
}
