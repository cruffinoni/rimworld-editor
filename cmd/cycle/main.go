package main

import (
	"flag"
	"strconv"
	"time"

	"github.com/cruffinoni/rimworld-editor/generated"
	"github.com/cruffinoni/rimworld-editor/internal/config"
	"github.com/cruffinoni/rimworld-editor/internal/file"
	"github.com/cruffinoni/rimworld-editor/internal/generator"
	"github.com/cruffinoni/rimworld-editor/internal/generator/files"
	"github.com/cruffinoni/rimworld-editor/internal/xml/saver/xmlFile"
	"github.com/cruffinoni/rimworld-editor/internal/xml/unmarshal"
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
		logger := logging.NewLogger("cycle", &fallback.Logging)
		logger.WithError(err).Fatal("Failed to load config")
	}
	logger := logging.NewLogger("cycle", &cfg.Logging)
	var (
		fo       *file.Opening
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
		logger.Error("No path specified")
		flag.Usage()
		return
	}
	logger.WithField("path", path).Debug("Opening and decoding XML file")
	fo, err = file.Open(path)
	if err != nil {
		logger.WithError(err).Fatal("Failed to open file")
		return
	}
	logger.Debug("Generating Go files to ./generated")
	root := generator.GenerateGoFiles(logger, fo.XML.Root, true)
	gw := files.NewGoWriter(logger, nil, true, "")
	if err = gw.WriteGoFile("./generated", root); err != nil {
		logger.WithError(err).Fatal("Failed to write Go files")
	}
	save := &generated.GeneratedStructStarter0{}
	logger.Debug("Unmarshalling XML")
	if err := unmarshal.Element(logger, fo.XML.Root, save); err != nil {
		logger.WithError(err).Fatal("Failed to unmarshal XML")
	}
	save.ValidateField("Savegame")
	logger.Debug("Generating XML file to folder")
	buffer, err := xmlFile.SaveWithBuffer(logger, save.Savegame)
	if err != nil {
		logger.WithError(err).Panic("Failed to save XML buffer")
	}
	if err := buffer.ToFile("generated/" + fileName + ".rws"); err != nil {
		logger.WithError(err).Panic("Failed to write XML file")
	}
}
