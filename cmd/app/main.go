package main

import (
	"log"
)

func main() {
	app := CreateApplication()
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
