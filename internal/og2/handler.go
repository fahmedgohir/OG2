package og2

import (
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
	router.Get("/dashboard", h.HandleSession())
	router.Post("/upgrade", h.HandleUpgrade())
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

		render.JSON(rw, r, session)
	}
}

func (h *Handler) HandleUpgrade() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var upgrade game.Upgrade
		if err := json.NewDecoder(r.Body).Decode(&upgrade); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		session, err := h.sessions.Get(upgrade.User)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		session, err = session.Upgrade(upgrade.Factory)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := h.sessions.Set(session); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		render.JSON(rw, r, session)
	}
}
