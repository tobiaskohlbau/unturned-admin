package app

import (
	"fmt"
	"log"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/go-chi/chi"
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

	s.router.Use(s.userMiddleware)

	apiRouter := chi.NewRouter()

	apiRouter.Post("/login", api.handleLogin())
	apiRouter.Get("/login/steam", api.handleSteamLogin(fmt.Sprintf("%s/api/login/steam/callback", os.Getenv("BASE_URL"))))
	apiRouter.Get("/login/steam/callback", api.handleSteamCallback(os.Getenv("STEAM_APIKEY")))
	apiRouter.Get("/users", s.requirePermission("ADMIN", api.handleUsers()))
	apiRouter.Post("/users/{username}", s.requirePermission("ADMIN", api.handleUserSave()))
	apiRouter.Get("/players", s.requireAuthenticated(api.handlePlayers(os.Getenv("UNTURNEDADMIN_ENDPOINT"))))
	apiRouter.Get("/rcon", s.requirePermission("MODERATOR", api.handleRCON(rconConnector)))
	apiRouter.Get("/backup", s.requirePermission("MODERATOR", api.handleBackup()))
	apiRouter.Get("/update", s.requirePermission("MODERATOR", api.handleUpdate(rconConnector)))
	apiRouter.Delete("/update", s.requirePermission("MODERATOR", api.handleUpdateCancel()))
	apiRouter.Route("/files", func(r chi.Router) {
		r.Get("/", s.requirePermission("MODERATOR", api.handleFiles()))
		r.Put("/", s.requirePermission("MODERATOR", api.handleFilesSave()))
		r.Delete("/", s.requirePermission("MODERATOR", api.handleFilesDelete()))
	})
	apiRouter.Get("/user", api.handleUser())

	s.router.Mount("/api", apiRouter)

	if s.devMode {
		target, err := url.Parse("http://localhost:3000")
		if err != nil {
			log.Panic("failed to create dev proxy: %w", err)
		}
		s.router.Get("/*", httputil.NewSingleHostReverseProxy(target).ServeHTTP)
	}
}
