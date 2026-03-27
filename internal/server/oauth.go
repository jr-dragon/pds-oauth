package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type OAuthHandler interface {
	ClientConfig(http.ResponseWriter, *http.Request)
	Start(http.ResponseWriter, *http.Request)
	Callback(http.ResponseWriter, *http.Request)
}

func oauthHandlers(h OAuthHandler) func(r chi.Router) {
	return func(r chi.Router) {
		r.HandleFunc("GET /oauth/client-metadata.json", h.ClientConfig)
		r.HandleFunc("GET /oauth/start", h.Start)
		r.HandleFunc("GET /oauth/callback", h.Callback)
	}
}
