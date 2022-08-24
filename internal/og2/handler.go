package og2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
)

type Handler struct {
	sessions *Sessions
}

func NewHandler(sessions *Sessions) *Handler {
	return &Handler{
		sessions: sessions,
	}
}

func (h *Handler) Route(router *chi.Mux) {
	router.Post("/user", h.HandleUser())
}

func (h *Handler) HandleUser() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)

		var user User
		if err := decoder.Decode(&user); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		session := h.sessions.Create(user)
		if session == nil {
			http.Error(rw, "could not create session", http.StatusInternalServerError)
			return
		}

		b, err := json.Marshal(session)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Println(string(b))

		render.JSON(rw, r, bytes.NewBuffer(b))
	}
}
