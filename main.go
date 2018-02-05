package main

import (
	"log"

	"cpsg-git.mattclark.guru/highlands/dt_benchmark/actions"
)

func main() {
	app := actions.App()
	if err := app.Serve(); err != nil {
		log.Fatal(err)
	}
}
