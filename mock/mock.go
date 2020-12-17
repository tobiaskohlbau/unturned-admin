package mock

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
)

type mockServer struct {
	router *chi.Mux
}

func New() http.Handler {
	srv := &mockServer{
		router: chi.NewRouter(),
	}

	srv.routes()

	return srv.router
}

func (s *mockServer) handleRoot() http.HandlerFunc {
	type id struct {
		PlayerName string `json:"playerName"`
	}
	type player struct {
		X  float32 `json:"x"`
		Y  float32 `json:"y"`
		ID id      `json:"id"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		players := []player{
			{
				X: 0.0,
				Y: 0.0,
				ID: id{
					PlayerName: "Test",
				},
			},
		}

		if err := json.NewEncoder(w).Encode(players); err != nil {
			http.Error(w, "failed to marshal players", http.StatusInternalServerError)
			return
		}
	}
}
