package main

import (
	"flag"
	"time"

	"github.com/briandowns/spinner"

	"github.com/cruffinoni/rimworld-editor/internal/codegen"
	"github.com/cruffinoni/rimworld-editor/internal/codegen/writer"
	"github.com/cruffinoni/rimworld-editor/internal/config"
	"github.com/cruffinoni/rimworld-editor/internal/xml/loader"
	"github.com/cruffinoni/rimworld-editor/pkg/logging"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fallback := &config.Config{
			Logging: config.LoggingConfig{
				Level:  "info",
				Format: "text",
				Output: "stderr",
			},
		}
		logger := logging.NewLogger("gen_go", &fallback.Logging)
		logger.WithError(err).Fatal("Failed to load config")
	}
	logger := logging.NewLogger("gen_go", &cfg.Logging)
	var (
		fo   *loader.Opening
		path string
	)
	flag.StringVar(&path, "path", "", "Path to the save game file")
	flag.Parse()
	if path == "" {
		logger.Error("No path specified")
		flag.Usage()
		return
	}
	s := spinner.New(spinner.CharSets[35], 100*time.Millisecond)
	logger.WithField("path", path).Debug("Opening and decoding XML file")
	s.FinalMSG = "XML file decoded successfully\n"
	s.Start()
	fo, err = loader.Open(path)
	if err != nil {
		logger.WithError(err).Fatal("Failed to open file")
		return
	}
	s.Stop()
	// s.Prefix = "Generating go files to './generated'... "
	s.FinalMSG = "Go files successfully generated\n"
	// s.Start()
	root := codegen.GenerateGoFiles(logger, fo.XML.Root, true)
	gw := writer.NewGoWriter(logger, nil, true, "")
	if err = gw.WriteGoFile("./generated", root); err != nil {
		logger.WithError(err).Fatal("Failed to write Go files")
	}
	// s.Stop()
}
