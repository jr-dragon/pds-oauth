package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewHandler(oauthsvc OAuthHandler) http.Handler {
	mux := chi.NewMux()

	mux.Group(oauthHandlers(oauthsvc))

	return mux
}
