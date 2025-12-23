package application

import (
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/briandowns/spinner"
	"github.com/jawher/mow.cli"
	"github.com/tcnksm/go-input"

	"github.com/cruffinoni/rimworld-editor/internal/rimworld/discovery"
	"github.com/cruffinoni/rimworld-editor/internal/rimworld/gamedata"
	"github.com/cruffinoni/rimworld-editor/internal/xml/loader"

	"github.com/cruffinoni/rimworld-editor/generated"
	"github.com/cruffinoni/rimworld-editor/internal/application/terminal"
	"github.com/cruffinoni/rimworld-editor/internal/application/userinterface"
	"github.com/cruffinoni/rimworld-editor/internal/codegen"
	"github.com/cruffinoni/rimworld-editor/internal/codegen/writer"
	"github.com/cruffinoni/rimworld-editor/internal/xml/binder"
	"github.com/cruffinoni/rimworld-editor/internal/xml/encoder/reflection"
	"github.com/cruffinoni/rimworld-editor/pkg/logging"
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
	userinterface.Options
	*cli.Cli

	fileOpening *loader.Opening
	ui          userinterface.Mode
	logger      logging.Logger
}

func CreateApplication(logger logging.Logger) *Application {
	app := &Application{
		logger: logger,
	}
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
			app.ui = &terminal.Console{}
		} else if app.Mode == modeGUI {
			panic("not implemented")
			// app.ui = app.guiMode
		}
		app.ui.SetLogger(app.logger)
		structInit := &generated.GeneratedStructStarter0{}
		app.logger.Debug("Unmarshalling XML")
		s.FinalMSG = "XML file unmarshalled successfully\n"
		s.Start()
		if err := binder.Element(app.logger, app.fileOpening.XML.Root, structInit); err != nil {
			app.logger.WithError(err).Fatal("Failed to unmarshal XML")
		}
		s.Stop()
		structInit.ValidateField("Savegame")
		app.logger.Debug("Initializing UI")
		app.ui.Init(&app.Options, structInit.Savegame)
		app.logger.Debug("Running UI")
		if err := app.ui.Execute(os.Args); err != nil {
			app.logger.WithError(err).Error("UI execution failed")
			return
		}
		if app.Save {
			app.logger.Debug("End of execution, generating new file")
			if err := app.SaveGameFile(structInit.Savegame); err != nil {
				app.logger.WithError(err).Error("Failed to save game file")
			}
		}
	}
	return app
}

func (app *Application) SaveGameFile(sg *generated.Savegame) error {
	buffer, err := reflection.SaveWithBuffer(app.logger, sg)
	if err != nil {
		app.logger.WithError(err).Panic("Failed to save XML buffer")
	}
	p, err := discovery.GetSavegamePath(app.OperatingSystem)
	if err != nil {
		return err
	}
	path := p + "/" + "Generated_" + strconv.FormatInt(time.Now().Unix(), 10) + ".rws"
	app.logger.WithField("path", path).Debug("Saving file")
	if err := buffer.ToFile(path); err != nil {
		return err
	}
	return nil
}

func (app *Application) beforeExecution() {
	if !isValidMode(app.Mode) {
		app.logger.WithField("mode", app.Mode).Error("Invalid mode")
		app.PrintHelp()
		cli.Exit(1)
	}
	gameData := gamedata.NewGameData(app.logger)
	app.logger.Debug("Discovering game data")
	err := gameData.DiscoverGameData(app.OperatingSystem)
	if err != nil {
		app.logger.WithError(err).Fatal("Failed to discover game data")
	}
	// gameData.PrintThemes()
	// e, err := gameData.FindElement("", "Scavenger22")
	// printer.Debugf("E: %v & Err %v", e.XMLPath(), err)
	app.logger.Debug("Generating Go files from game data")
	if err := gameData.GenerateGoFiles(); err != nil {
		app.logger.WithError(err).Fatal("Failed to generate game files")
	}
	//if err := gameData.ReadGameFiles(); err != nil {
	//	log.Fatal(err)
	//}
	//os.Exit(0)
	if err = app.ReadSaveGame(); err != nil {
		app.logger.WithError(err).Fatal("Failed to read save game")
	}
	if app.Generate {
		s := spinner.New(spinner.CharSets[35], 100*time.Millisecond)
		s.FinalMSG = "Generating Go files successfully\n"
		s.Start()
		root := codegen.GenerateGoFiles(app.logger, app.fileOpening.XML.Root, true)
		s.Stop()
		gw := writer.NewGoWriter(app.logger, nil, true, "")
		if err = gw.WriteGoFile(app.Output, root); err != nil {
			app.logger.WithError(err).Fatal("Failed to write generated files")
		}
		if err = app.fileOpening.ReOpen(); err != nil {
			app.logger.WithError(err).Fatal("Failed to reopen savegame")
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
	savegamesDirPath, err := discovery.GetSavegamePath(app.OperatingSystem)
	if err != nil {
		return err
	}
	saves, err := discovery.GetLatestSavegameFiles(app.MaxSaveGameFileDiscover, savegamesDirPath)
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
	app.fileOpening, err = loader.Open(savegame)
	if err != nil {
		app.logger.WithError(err).Error("Failed to open savegame")
	} else {
		app.logger.WithField("path", savegame).Debug("Savegame found")
	}
	return err
}
