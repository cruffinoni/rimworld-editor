package main

import (
	"github.com/cruffinoni/rimworld-editor/cmd/app/ui"
	"github.com/cruffinoni/rimworld-editor/cmd/app/ui/terminal"
	"github.com/cruffinoni/rimworld-editor/file"
	"github.com/cruffinoni/rimworld-editor/generated"
	"github.com/cruffinoni/rimworld-editor/generator"
	"github.com/cruffinoni/rimworld-editor/xml/unmarshal"
	"github.com/jawher/mow.cli"
	"log"
	"os"
)

const (
	cliVersion = "0.0.1"
)

const (
	modeConsole = "console"
	modeGUI     = "gui"
)

var modes = []string{modeConsole, modeGUI}

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

/*
	TODO:
		- Console mode: saveXML, generateGoFiles
		- GUI mode: edit data
*/

func CreateApplication() *Application {
	app := &Application{}
	app.Cli = cli.App("rimworld-editor", "Rimworld save game editor")
	app.Version("version", cliVersion)
	app.BoolOptPtr(&app.Verbose, "v verbose", false, "Verbose mode")
	app.BoolOptPtr(&app.Generate, "g generate", false, "Generate go files from xml")
	app.StringOptPtr(&app.Output, "o output", "generated", "Output folder for generated files")
	app.StringOptPtr(&app.Mode, "m mode", modeConsole, "The mode to run the application in")
	app.StringArgPtr(&app.Input, "INPUT", "", "Save game file to explore") // TODO: Later use StringOptPtr and discover the file automatically
	app.Before = app.beforeExecution
	app.Action = func() {
		if app.Mode == modeConsole {
			app.ui = &terminal.Console{}
		} else if app.Mode == modeGUI {
			panic("not implemented")
			//app.ui = app.guiMode
		}
		save := &generated.Save{}
		log.Println("Unmarshalling XML...")
		if err := unmarshal.Element(app.fileOpening.XML.Root, save); err != nil {
			log.Fatal(err)
		}
		log.Println("Initializing UI...")
		app.ui.Init(&app.Options, save)
		log.Println("Running UI...")
		if err := app.ui.Execute(os.Args); err != nil {
			log.Fatal(err)
		}
	}
	return app
}

func (app *Application) beforeExecution() {
	if !isValidMode(app.Mode) {
		log.Printf("->invalid usage: %v", app.Mode)
		app.PrintHelp()
		cli.Exit(1)
	}
	var err error
	if err = app.ReadFile(); err != nil {
		log.Fatal(err)
	}
	if app.Generate {
		root := generator.GenerateGoFiles(app.fileOpening.XML.Root)
		if err = root.WriteGoFile(app.Output); err != nil {
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

func (app *Application) ReadFile() error {
	var err error
	app.fileOpening, err = file.Open(app.Input)
	return err
}
