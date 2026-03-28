package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type UserHandler interface {
	Me(w http.ResponseWriter, r *http.Request)
}

func NewHandler(oauthsvc OAuthHandler, usersvc UserHandler) http.Handler {
	mux := chi.NewMux()

	mux.Group(oauthHandlers(oauthsvc))
	mux.HandleFunc("GET /api/me", usersvc.Me)

	return mux
}
