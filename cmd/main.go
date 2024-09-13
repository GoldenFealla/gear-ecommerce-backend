package main

import (
	"log"

	"github.com/goldenfealla/gear-manager/app"
)

func main() {
	server := app.New()

	err := server.Start()

	if err != nil {
		log.Fatalln(err)
	}
}
