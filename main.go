package main

import (
	"log"

	"github.com/mclark4386/dt_benchmark/actions"
)

func main() {
	app := actions.App()
	if err := app.Serve(); err != nil {
		log.Fatal(err)
	}
}
