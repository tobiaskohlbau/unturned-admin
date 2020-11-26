package app

import (
	"crypto/ed25519"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/tobiaskohlbau/unturned-admin/web"
)

func (s *appServer) routes() {
	api := &apiServer{
		appServer: s,

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

	seed := []byte(os.Getenv("SIGNATURE_SEED"))
	s.privateKey = ed25519.NewKeyFromSeed(seed)

	apiRouter := chi.NewRouter()

	apiRouter.Post("/login", api.handleLogin())
	apiRouter.Get("/rcon", s.requirePermission("ADMIN", api.handleRCON(rconConnector)))
	apiRouter.Get("/backup", s.requirePermission("ADMIN", api.handleBackup()))
	apiRouter.Get("/update", s.requirePermission("ADMIN", api.handleUpdate(rconConnector)))
	apiRouter.Delete("/update", s.requirePermission("ADMIN", api.handleUpdateCancel()))
	apiRouter.Route("/files", func(r chi.Router) {
		r.Get("/", s.requirePermission("ADMIN", api.handleFiles()))
		r.Put("/", s.requirePermission("ADMIN", api.handleFilesSave()))
		r.Delete("/", s.requirePermission("ADMIN", api.handleFilesDelete()))
	})

	s.router.Mount("/api", apiRouter)

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
