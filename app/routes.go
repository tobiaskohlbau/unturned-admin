package app

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/tobiaskohlbau/unturned-admin/web"
)

func (s *appServer) routes() {
	api := &apiServer{
		wsMaxMessageSize:   8192,
		wsWriteWait:        10 * time.Second,
		wsPongWait:         60 * time.Second,
		wsPingPeriod:       (60 * time.Second * 9) / 10,
		wsCloseGracePerdio: 10 * time.Second,
	}

	// handle api endpoints
	rconConnector := establishRCONConnection
	if s.devMode {
		rconConnector = establishMockRCONConnection
	}

	s.router.Get("/api/rcon", api.handleRCON(rconConnector))
	s.router.Get("/api/screenshots", api.handleScreenshots())

	// handle static files
	if s.devMode {
		target, err := url.Parse("http://localhost:3000")
		if err != nil {
			log.Panic("failed to create dev proxy: %w", err)
		}
		s.router.Get("/*", httputil.NewSingleHostReverseProxy(target).ServeHTTP)
	} else {
		s.router.Get("/*", http.FileServer(http.FS(web.Content)).ServeHTTP)
	}
}
