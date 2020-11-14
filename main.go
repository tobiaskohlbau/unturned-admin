package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/tobiaskohlbau/unturned-admin/app"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	devMode := flag.Bool("dev", false, "development mode")
	flag.Parse()

	srv := app.New(*devMode)
	return http.ListenAndServe(":8080", srv)
}
