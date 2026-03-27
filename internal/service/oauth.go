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

// startResponse is the JSON body returned by GET /oauth/start.
type startResponse struct {
	RedirectURL string `json:"redirect_url"`
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(startResponse{RedirectURL: redirectURL})
}

// Callback handles GET /oauth/callback.
//
// The PDS redirects the user's browser here with ?code=...&state=... after
// authorization. API verifies the state (via indigo), performs the token
// exchange, creates a server-side session, sets the httpOnly cookie, and
// redirects the browser to the home page.
func (svc *OAuth) Callback(w http.ResponseWriter, r *http.Request) {
	sessData, err := svc.client.ProcessCallback(r.Context(), r.URL.Query())
	if err != nil {
		slog.WarnContext(r.Context(), "failed to get OAuth callback", liblogs.ErrAttr(err))
		http.Redirect(w, r, "/?error="+err.Error(), http.StatusFound)
		return
	}

	var handle string
	if ident, err := identity.DefaultDirectory().LookupDID(r.Context(), sessData.AccountDID); err == nil {
		handle = ident.Handle.String()
	} else {
		slog.WarnContext(r.Context(), "failed to lookup DID", liblogs.ErrAttr(err), slog.Any("DID", sessData.AccountDID))
	}

	sess, err := svc.store.Get(r, "oauth")
	if err != nil {
		slog.ErrorContext(r.Context(), "failed to get oauth session", liblogs.ErrAttr(err))
		http.Redirect(w, r, "/?error=failed to process oauth session", http.StatusFound)
		return
	}

	sess.Values["AccountDID"] = sessData.AccountDID
	sess.Values["SessionID"] = sessData.SessionID
	sess.Values["handle"] = handle
	if err := sess.Save(r, w); err != nil {
		slog.ErrorContext(r.Context(), "failed to save oauth session", liblogs.ErrAttr(err))
		http.Redirect(w, r, "/?error=failed to process oauth session", http.StatusFound)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
