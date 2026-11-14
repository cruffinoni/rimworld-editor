package application

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"github.com/tcnksm/go-input"

	"github.com/cruffinoni/rimworld-editor/internal/rimworld/discovery"
	"github.com/cruffinoni/rimworld-editor/internal/rimworld/gamedata"
	"github.com/cruffinoni/rimworld-editor/internal/xml/loader"

	"github.com/cruffinoni/rimworld-editor/generated"
	"github.com/cruffinoni/rimworld-editor/internal/application/ui"
	"github.com/cruffinoni/rimworld-editor/internal/codegen"
	"github.com/cruffinoni/rimworld-editor/internal/codegen/writer"
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
	ui.Options
	rootCmd *cobra.Command

	fileOpening *loader.Opening
	ui          ui.Mode
	logger      logging.Logger
}

func CreateApplication(logger logging.Logger) *Application {
	app := &Application{
		logger: logger,
	}
	app.rootCmd = app.newRootCommand()
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

func (app *Application) beforeExecution(cmd *cobra.Command) error {
	if !isValidMode(app.Mode) {
		app.logger.WithField("mode", app.Mode).Error("Invalid mode")
		if cmd != nil {
			_ = cmd.Help()
		}
		return fmt.Errorf("invalid mode: %s", app.Mode)
	}
	gameData := gamedata.NewGameData(app.logger)
	app.logger.Debug("Discovering game data")
	err := gameData.DiscoverGameData(app.OperatingSystem)
	if err != nil {
		return err
	}
	// gameData.PrintThemes()
	// e, err := gameData.FindElement("", "Scavenger22")
	// printer.Debugf("E: %v & Err %v", e.XMLPath(), err)
	app.logger.Debug("Generating Go files from game data")
	if err := gameData.GenerateGoFiles(); err != nil {
		return err
	}
	//if err := gameData.ReadGameFiles(); err != nil {
	//	log.Fatal(err)
	//}
	//os.Exit(0)
	if err = app.ReadSaveGame(); err != nil {
		return err
	}
	if app.Generate {
		s := spinner.New(spinner.CharSets[35], 100*time.Millisecond)
		s.FinalMSG = "Generating Go files successfully\n"
		s.Start()
		root := codegen.GenerateGoFiles(app.logger, app.fileOpening.XML.Root, true)
		s.Stop()
		gw := writer.NewGoWriter(app.logger, nil, true, "")
		if err = gw.WriteGoFile(app.Output, root); err != nil {
			return err
		}
		if err = app.fileOpening.ReOpen(); err != nil {
			return err
		}
	}
	return nil
}

func (app *Application) Run() error {
	app.rootCmd.SetArgs(os.Args[1:])
	return app.rootCmd.Execute()
}

func (app *Application) RunWithArgs(args []string) error {
	app.rootCmd.SetArgs(args)
	return app.rootCmd.Execute()
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
		u := &input.UI{
			Writer: os.Stdout,
			Reader: os.Stdin,
		}
		var joinedFileName []string
		for _, s := range saves {
			joinedFileName = append(joinedFileName, s.Name())
		}

		selected, err := u.Select("Which savegame do you want to select", joinedFileName, &input.Options{
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
