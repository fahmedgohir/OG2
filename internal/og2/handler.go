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

type UserRequest struct {
	User game.User `json:"user"`
}

func (h *Handler) HandleUser() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var req UserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		err := h.sessions.Create(req.User)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		session, err := h.sessions.Get(req.User)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		render.JSON(rw, r, session)
	}
}

type SessionRequest struct {
	User game.User `json:"user"`
}

func (h *Handler) HandleSession() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var req SessionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		session, err := h.sessions.Get(req.User)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		render.JSON(rw, r, session)
	}
}

type UpgradeRequest struct {
	User    game.User     `json:"user"`
	Factory game.Resource `json:"factory"`
}

func (h *Handler) HandleUpgrade() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var req UpgradeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		session, err := h.sessions.Get(req.User)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := session.Upgrade(req.Factory); err != nil {
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
