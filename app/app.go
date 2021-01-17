package app

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
	"github.com/tobiaskohlbau/unturned-admin/store"
	bolt "go.etcd.io/bbolt"
)

type appServer struct {
	router  *chi.Mux
	devMode bool

	privateKey ed25519.PrivateKey
	userStore  *store.UserStore
}

func New(devMode bool) http.Handler {
	db, err := bolt.Open("database.db", 0600, nil)
	if err != nil {
		log.Error().Err(err).Msg("failed to open database")
		return nil
	}

	seed := []byte(os.Getenv("SIGNATURE_SEED"))
	srv := &appServer{
		devMode:    devMode,
		router:     chi.NewRouter(),
		privateKey: ed25519.NewKeyFromSeed(seed),
		userStore:  store.NewUserStore(db),
	}

	// insert or update default users
	f, err := os.Open("users.json")
	var pathErr fs.PathError
	if err != nil && errors.As(err, &pathErr) {
		log.Error().Err(err).Msg("failed to open default users")
	}
	users := []store.User{}
	if err := json.NewDecoder(f).Decode(&users); err != nil {
		log.Error().Err(err).Msg("failed to decode default users")
		return nil
	}

	for _, user := range users {
		if err := srv.userStore.Add(user); err != nil {
			log.Error().Err(err).Str("username", user.Username).Msg("failed to insert default user")
		}
	}

	srv.routes()

	return srv.router
}

type Token struct {
	Username    string   `json:"username"`
	Permissions []string `json:"permissions"`
	Activated   bool     `json:"activated"`
}

type ctxKeyUser int

const UserKey ctxKeyUser = 0

func GetUser(ctx context.Context) (store.User, error) {
	if ctx == nil {
		return store.User{}, errors.New("invalid context")
	}
	if user, ok := ctx.Value(UserKey).(store.User); ok {
		return user, nil
	}
	return store.User{}, errors.New("no token in context")
}

func (s *appServer) userMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tidCookie, err := r.Cookie("tid")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		sidCookie, err := r.Cookie("sid")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		signature, err := base64.RawStdEncoding.DecodeString(sidCookie.Value)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		token, err := base64.RawStdEncoding.DecodeString(tidCookie.Value)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		if !ed25519.Verify(s.privateKey.Public().(ed25519.PublicKey), token, signature) {
			next.ServeHTTP(w, r)
			return
		}

		tok := Token{}
		if err = json.Unmarshal(token, &tok); err != nil {
			next.ServeHTTP(w, r)
			return
		}

		needsUpdate := false
		if tidCookie.Expires.Sub(time.Now()).Minutes() < 30 {
			needsUpdate = true
		}

		usr, ok := s.userStore.GetByUsername(tok.Username)
		if ok {
			if usr.Activated != tok.Activated {
				log.Info().Msg("user acitvation changed")
				tok.Activated = usr.Activated
				needsUpdate = true
			}

			sort.Strings(usr.Permissions)
			sort.Strings(tok.Permissions)

			if len(tok.Permissions) != len(usr.Permissions) {
				log.Info().Msg("user permissions changed")
				tok.Permissions = usr.Permissions
				needsUpdate = true
			} else {
				for i, p := range usr.Permissions {
					if tok.Permissions[i] != p {
						log.Info().Msg("user permissions changed")
						tok.Permissions = usr.Permissions
						needsUpdate = true
						break
					}
				}
			}
		}

		if needsUpdate {
			s.setAuthCookies(w, r, &tok)
		}

		ctx := context.WithValue(r.Context(), UserKey, usr)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *appServer) setAuthCookies(w http.ResponseWriter, r *http.Request, token *Token) error {
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
		user, err := GetUser(r.Context())
		if err != nil {
			log.Error().Err(err).Msg("no user present")
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		hasPermission := false
		for _, p := range user.Permissions {
			if p == permission {
				hasPermission = true
				break
			}
		}
		if !hasPermission {
			log.Info().Str("username", user.Username).Str("need", permission).Strs("has", user.Permissions).Msg("user does not have permission")
			http.Error(w, "insufficient permissions", http.StatusUnauthorized)
			return
		}

		fn(w, r)
	}
}

func (s *appServer) requireAuthenticated(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := GetUser(r.Context())
		if err != nil {
			log.Error().Err(err).Msg("no user present")
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		fn(w, r)
	}
}
