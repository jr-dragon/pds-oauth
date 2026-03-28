package service

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/jr-dragon/pds-oauth/internal/data"
	"github.com/jr-dragon/pds-oauth/internal/lib/libhttp"
	"github.com/jr-dragon/pds-oauth/internal/lib/liblogs"
)

type User struct {
	cfg   *data.Config
	store sessions.Store
}

func NewUser(cfg *data.Config) *User {
	return &User{
		cfg:   cfg,
		store: sessions.NewCookieStore(cfg.App.Key),
	}
}

func (svc *User) Me(w http.ResponseWriter, r *http.Request) {
	sess, err := svc.store.Get(r, "atproto")
	if err != nil {
		slog.ErrorContext(r.Context(), "failed to get session", liblogs.ErrAttr(err))
		libhttp.WriteError(w, http.StatusUnauthorized, "failed to authorize")
		return
	}

	slog.InfoContext(r.Context(), "session", slog.Any("session", sess))
	w.WriteHeader(http.StatusNoContent)
}
