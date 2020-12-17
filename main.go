package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/tobiaskohlbau/unturned-admin/app"
	"github.com/tobiaskohlbau/unturned-admin/mock"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	devMode := flag.Bool("dev", false, "development mode")
	flag.Parse()

	if *devMode {
		go func() {
			mock := mock.New()
			http.ListenAndServe(":8000", mock)
		}()
	}

	srv := app.New(*devMode)
	return http.ListenAndServe(":8080", srv)
}
