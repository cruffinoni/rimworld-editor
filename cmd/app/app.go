package main

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/briandowns/spinner"
	"github.com/jawher/mow.cli"
	"github.com/tcnksm/go-input"

	"github.com/cruffinoni/printer"

	"github.com/cruffinoni/rimworld-editor/internal/file"
	"github.com/cruffinoni/rimworld-editor/internal/resources"
	"github.com/cruffinoni/rimworld-editor/internal/resources/discover"

	"github.com/cruffinoni/rimworld-editor/cmd/app/ui"
	"github.com/cruffinoni/rimworld-editor/cmd/app/ui/term"
	"github.com/cruffinoni/rimworld-editor/generated"
	"github.com/cruffinoni/rimworld-editor/internal/generator"
	"github.com/cruffinoni/rimworld-editor/internal/generator/files"
	"github.com/cruffinoni/rimworld-editor/internal/xml/saver/xmlFile"
	"github.com/cruffinoni/rimworld-editor/internal/xml/unmarshal"
)

const (
	cliVersion = "0.0.1"
)

const (
	modeConsole = "console"
	modeGUI     = "gui"
)

func isValidMode(mode string) bool {
	return mode == modeConsole || mode == modeGUI
}

// Application is the main application.
type Application struct {
	ui.Options
	*cli.Cli

	fileOpening *file.Opening
	ui          ui.Mode
}

func CreateApplication() *Application {
	app := &Application{}
	s := spinner.New(spinner.CharSets[21], 100*time.Millisecond)
	app.Cli = cli.App("rimworld-editor", "Rimworld save game editor")
	app.Version("version", cliVersion)
	app.BoolOptPtr(&app.Verbose, "v verbose", false, "Verbose mode")
	app.BoolOptPtr(&app.Generate, "g generate", false, "Generate go files from xml")
	app.BoolOptPtr(&app.Save, "s save", true, "Save your modifications when exiting the application")
	app.StringOptPtr(&app.Output, "o output", "generated", "Output folder for generated files")
	app.StringOptPtr(&app.Mode, "m mode", modeConsole, "The mode to run the application in")
	app.IntOptPtr(&app.MaxSaveGameFileDiscover, "mx maxnb", 10, "Maximum number of save games to discover")
	app.StringOptPtr(&app.Input, "ds defaultsave", "", "Default save game to load from your Rimworld saves game folder")
	app.StringOptPtr(&app.OperatingSystem, "operating-system os", "", "Force a operating system file path finding")
	app.Before = app.beforeExecution
	app.Action = func() {
		if app.Mode == modeConsole {
			app.ui = &term.Console{}
		} else if app.Mode == modeGUI {
			panic("not implemented")
			// app.ui = app.guiMode
		}
		structInit := &generated.GeneratedStructStarter0{}
		printer.Debugf("Unmarshalling XML...")
		s.FinalMSG = "XML file unmarshalled successfully\n"
		s.Start()
		if err := unmarshal.Element(app.fileOpening.XML.Root, structInit); err != nil {
			log.Fatal(err)
		}
		s.Stop()
		structInit.ValidateField("Savegame")
		printer.Debugf("Initializing UI...")
		app.ui.Init(&app.Options, structInit.Savegame)
		printer.Debugf("Running UI...")
		if err := app.ui.Execute(os.Args); err != nil {
			printer.PrintError(err)
			return
		}
		if app.Save {
			printer.Debugf("End of execution, generating new file...")
			if err := app.SaveGameFile(structInit.Savegame); err != nil {
				printer.PrintError(err)
			}
		}
	}
	return app
}

func (app *Application) SaveGameFile(sg *generated.Savegame) error {
	buffer, err := xmlFile.SaveWithBuffer(sg)
	if err != nil {
		log.Panic(err)
	}
	p, err := discover.GetSavegamePath(app.OperatingSystem)
	if err != nil {
		return err
	}
	path := p + "/" + "Generated_" + strconv.FormatInt(time.Now().Unix(), 10) + ".rws"
	printer.Debugf("Saving file to '%s'", path)
	if err := buffer.ToFile(path); err != nil {
		return err
	}
	return nil
}

func (app *Application) beforeExecution() {
	if !isValidMode(app.Mode) {
		printer.PrintErrorSf("invalid mode: %v", app.Mode)
		app.PrintHelp()
		cli.Exit(1)
	}
	gameData := resources.NewGameData()
	printer.Debugf("Discovering game data...")
	err := gameData.DiscoverGameData(app.OperatingSystem)
	if err != nil {
		log.Fatal(err)
	}
	// gameData.PrintThemes()
	// e, err := gameData.FindElement("", "Scavenger22")
	// printer.Debugf("E: %v & Err %v", e.XMLPath(), err)
	printer.Debugf("Generating Go files from game data...")
	if err := gameData.GenerateGoFiles(); err != nil {
		log.Fatal(err)
	}
	//if err := gameData.ReadGameFiles(); err != nil {
	//	log.Fatal(err)
	//}
	//os.Exit(0)
	if err = app.ReadSaveGame(); err != nil {
		log.Fatal(err)
	}
	if app.Generate {
		s := spinner.New(spinner.CharSets[35], 100*time.Millisecond)
		s.FinalMSG = "Generating Go files successfully\n"
		s.Start()
		root := generator.GenerateGoFiles(app.fileOpening.XML.Root, true)
		s.Stop()
		if err = files.DefaultGoWriter.WriteGoFile(app.Output, root); err != nil {
			log.Fatal(err)
		}
		if err = app.fileOpening.ReOpen(); err != nil {
			log.Fatal(err)
		}
	}
}

func (app *Application) Run() error {
	return app.Cli.Run(os.Args)
}

func (app *Application) RunWithArgs(args []string) error {
	return app.Cli.Run(args)
}

func (app *Application) ReadSaveGame() error {
	var savegame string
	savegamesDirPath, err := discover.GetSavegamePath(app.OperatingSystem)
	if err != nil {
		return err
	}
	saves, err := discover.GetLatestSavegameFiles(app.MaxSaveGameFileDiscover, savegamesDirPath)
	if err != nil {
		return err
	}
	if app.Input != "" {
		for _, s := range saves {
			if s.Name() == app.Input {
				savegame = filepath.Join(savegamesDirPath, s.Name())
				break
			}
		}
	} else {
		ui := &input.UI{
			Writer: os.Stdout,
			Reader: os.Stdin,
		}
		var joinedFileName []string
		for _, s := range saves {
			joinedFileName = append(joinedFileName, s.Name())
		}

		selected, err := ui.Select("What savegame do you want to select", joinedFileName, &input.Options{
			Required: true,
			Loop:     true,
		})
		if err != nil {
			return err
		}
		savegame = filepath.Join(savegamesDirPath, selected)
	}
	app.fileOpening, err = file.Open(savegame)
	if err != nil {
		printer.PrintError(err)
	} else {
		printer.Debugf("Savegame found at %v", savegame)
	}
	return err
}
