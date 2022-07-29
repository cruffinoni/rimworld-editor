package main

import (
	"github.com/cruffinoni/rimworld-editor/editor"
	"github.com/cruffinoni/rimworld-editor/generator"
	"github.com/jawher/mow.cli"
	"log"
	"os"
)

const (
	cliVersion = "0.0.1"
)

type modeType string

const (
	modeConsole = modeType("console")
	modeGUI     = modeType("gui")
)

func isValidMode(mode string) bool {
	return mode == string(modeConsole) || mode == string(modeGUI)
}

// options is a struct that contains the options for the application.
type options struct {
	verbose  bool
	generate bool
	output   string
	input    string
	mode     string
}

// Application is the main application.
type Application struct {
	options
	*cli.Cli

	fileOpening *editor.FileOpening
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
	app.BoolOptPtr(&app.verbose, "v verbose", false, "Verbose mode")
	app.BoolOptPtr(&app.generate, "g generate", false, "Generate go files from xml")
	app.StringOptPtr(&app.output, "o output", "generated", "Output file")
	app.StringArgPtr(&app.mode, "MODE", string(modeConsole), "The mode to run the application in")
	app.StringArgPtr(&app.input, "INPUT", "", "Save game file to explore") // TODO: Later use StringOptPtr and discover the file automatically
	app.Before = app.beforeExecution
	app.Action = app.doAction
	return app
}

func (app *Application) beforeExecution() {
	if !isValidMode(app.mode) {
		app.PrintHelp()
		cli.Exit(1)
	}
	var err error
	if err = app.ReadFile(); err != nil {
		log.Fatal(err)
	}
	if app.generate {
		root := generator.GenerateGoFiles(app.fileOpening.XML.Root)
		if err = root.WriteGoFile(app.output); err != nil {
			log.Fatal(err)
		}
	}
	//if app.mode == string(modeConsole) {
	//	app.Action = app.consoleMode
	//} else if app.mode == string(modeGUI) {
	//	app.Action = app.guiMode
	//}
}

func (app *Application) doAction() {

}

func (app *Application) Run() error {
	return app.Cli.Run(os.Args)
}

func (app *Application) RunWithArgs(args []string) error {
	return app.Cli.Run(args)
}

func (app *Application) ReadFile() error {
	var err error
	app.fileOpening, err = editor.Open(app.input)
	return err
}
