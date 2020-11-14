package app

import (
	"net/http"

	"github.com/go-chi/chi"
)

type appServer struct {
	router  *chi.Mux
	devMode bool
}

func New(devMode bool) http.Handler {
	srv := &appServer{
		devMode: devMode,
		router:  chi.NewRouter(),
	}

	srv.routes()

	return srv.router
}
