package service

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/bluesky-social/indigo/atproto/auth/oauth"
	"github.com/bluesky-social/indigo/atproto/identity"
	"github.com/gorilla/sessions"

	"github.com/jr-dragon/pds-oauth/internal/data"
	"github.com/jr-dragon/pds-oauth/internal/lib/libhttp"
	"github.com/jr-dragon/pds-oauth/internal/lib/liblogs"
)

type OAuth struct {
	client *oauth.ClientApp
	store  sessions.Store
}

func NewOAuth(cfg *data.Config, client *oauth.ClientApp) *OAuth {
	return &OAuth{
		client: client,
		store:  sessions.NewCookieStore(cfg.App.Key),
	}
}

func (svc *OAuth) ClientConfig(w http.ResponseWriter, r *http.Request) {
	libhttp.WriteJSON(w, http.StatusOK, svc.client.Config.ClientMetadata())
}

// Start handles GET /oauth/start?handle=<handle>.
//
// Resolves the user's PDS, initiates a PAR request, and returns the
// authorization redirect URL for the SPA to navigate the user to.
func (svc *OAuth) Start(w http.ResponseWriter, r *http.Request) {
	handle := strings.TrimSpace(r.URL.Query().Get("handle"))
	if handle == "" {
		libhttp.WriteError(w, http.StatusBadRequest, "missing required query parameter: handle")
		return
	}

	redirectURL, err := svc.client.StartAuthFlow(r.Context(), handle)
	if err != nil {
		slog.ErrorContext(r.Context(), "failed to start auth flow", liblogs.ErrAttr(err))
		libhttp.WriteError(w, http.StatusBadRequest, "failed to start auth flow")
		return
	}

	http.Redirect(w, r, redirectURL, http.StatusFound)
}

// Callback handles GET /oauth/callback.
//
// The PDS redirects the user's browser here with ?code=...&state=... after
// authorization. API verifies the state (via indigo), performs the token
// exchange, creates a server-side session, sets the httpOnly cookie, and
// redirects the browser to the home page.
func (svc *OAuth) Callback(w http.ResponseWriter, r *http.Request) {
	atSession, err := svc.client.ProcessCallback(r.Context(), r.URL.Query())
	if err != nil {
		slog.WarnContext(r.Context(), "failed to get OAuth callback", liblogs.ErrAttr(err))
		http.Redirect(w, r, "/?error="+err.Error(), http.StatusFound)
		return
	}

	atIdent, err := identity.DefaultDirectory().LookupDID(r.Context(), atSession.AccountDID)
	if err != nil {
		slog.WarnContext(r.Context(), "failed to lookup DID", liblogs.ErrAttr(err), slog.Any("at_session", atSession))
	}

	sess, err := svc.store.Get(r, "atproto")
	if err != nil {
		slog.ErrorContext(r.Context(), "failed to get oauth session", liblogs.ErrAttr(err))
		http.Redirect(w, r, "/?error=failed to process oauth session", http.StatusFound)
		return
	}

	if sess.Values["Session"], err = json.Marshal(atSession); err != nil {
		slog.ErrorContext(r.Context(), "failed to marshal json", liblogs.ErrAttr(err))
		libhttp.WriteError(w, http.StatusInternalServerError, "")
		return
	}
	if sess.Values["Identity"], err = json.Marshal(atIdent); err != nil {
		slog.ErrorContext(r.Context(), "failed to marshal json", liblogs.ErrAttr(err))
		libhttp.WriteError(w, http.StatusInternalServerError, "")
	}
	if err := sess.Save(r, w); err != nil {
		slog.ErrorContext(r.Context(), "failed to save oauth session", liblogs.ErrAttr(err))
		http.Redirect(w, r, "/?error=failed to process oauth session", http.StatusFound)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
