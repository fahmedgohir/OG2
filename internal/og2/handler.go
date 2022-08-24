package og2

import (
	"bytes"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"hunter.io/og2/internal/og2/game"
	"net/http"
)

type Handler struct {
	sessions Sessions
}

func NewHandler(sessions Sessions) *Handler {
	return &Handler{
		sessions: sessions,
	}
}

func (h *Handler) Route(router *chi.Mux) {
	router.Post("/user", h.HandleUser())
	router.Get("/user", h.HandleSession())
}

func (h *Handler) HandleUser() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var user game.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		err := h.sessions.Create(user)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		session, err := h.sessions.Get(user)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		render.JSON(rw, r, session)
	}
}

func (h *Handler) HandleSession() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var user game.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		session, err := h.sessions.Get(user)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		b, err := json.Marshal(session)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		render.JSON(rw, r, bytes.NewReader(b))
	}
}
