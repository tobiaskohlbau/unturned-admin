package app

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
)

type appServer struct {
	router  *chi.Mux
	devMode bool

	privateKey ed25519.PrivateKey
}

func New(devMode bool) http.Handler {
	srv := &appServer{
		devMode: devMode,
		router:  chi.NewRouter(),
	}

	srv.routes()

	return srv.router
}

type Token struct {
	Username    string   `json:"username"`
	Permissions []string `json:"permissions"`
}

func (s *appServer) setAuthCookies(w http.ResponseWriter, r *http.Request, token Token) error {
	data, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("failed to marshal authentication token: %w", err)
	}

	signature := ed25519.Sign(s.privateKey, data)

	http.SetCookie(w, &http.Cookie{
		Name:     "tid",
		Value:    base64.RawStdEncoding.EncodeToString(data),
		SameSite: http.SameSiteDefaultMode,
		Secure:   !s.devMode,
		Expires:  time.Now().Add(24 * time.Hour),
		Path:     "/",
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "sid",
		Value:    base64.RawStdEncoding.EncodeToString(signature),
		SameSite: http.SameSiteDefaultMode,
		HttpOnly: true,
		Secure:   !s.devMode,
		Expires:  time.Now().Add(24 * time.Hour),
		Path:     "/",
	})

	return nil
}

func (s *appServer) requirePermission(permission string, fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tidCookie, err := r.Cookie("tid")
		if err != nil {
			log.Error().Err(err).Msg("tid error")
			http.Error(w, "missing tid", http.StatusUnauthorized)
			return
		}

		sidCookie, err := r.Cookie("sid")
		if err != nil {
			log.Error().Err(err).Msg("sid error")
			http.Error(w, "missing sid", http.StatusUnauthorized)
			return
		}

		signature, err := base64.RawStdEncoding.DecodeString(sidCookie.Value)
		if err != nil {
			http.Error(w, "bad signature value", http.StatusUnauthorized)
			return
		}

		token, err := base64.RawStdEncoding.DecodeString(tidCookie.Value)
		if err != nil {
			http.Error(w, "bad token value", http.StatusUnauthorized)
			return
		}

		if !ed25519.Verify(s.privateKey.Public().(ed25519.PublicKey), token, signature) {
			http.Error(w, "invalid signature", http.StatusUnauthorized)
			return
		}

		var tok Token
		if err = json.Unmarshal(token, &tok); err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		hasPermission := false
		for _, p := range tok.Permissions {
			if p == permission {
				hasPermission = true
				break
			}
		}
		if !hasPermission {
			log.Info().Str("username", tok.Username).Str("need", permission).Strs("has", tok.Permissions).Msg("user does not have permission")
			http.Error(w, "insufficient permissions", http.StatusUnauthorized)
			return
		}

		if tidCookie.Expires.Sub(time.Now()).Minutes() < 30 {
			s.setAuthCookies(w, r, tok)
		}

		fn(w, r)
	}
}
