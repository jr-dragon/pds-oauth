package main

//go:generate go tool kessoku $GOFILE

import (
	"net/http"
	"path"

	"github.com/bluesky-social/indigo/atproto/auth/oauth"
	"github.com/mazrean/kessoku"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/jr-dragon/pds-oauth/internal/data"
	"github.com/jr-dragon/pds-oauth/internal/server"
	"github.com/jr-dragon/pds-oauth/internal/service"
)

// newServer is the kessoku-generated DI initializer.
// Inputs:  *data.Config	(created in main with special lifecycle handling)
// Output:  *server.Server	(fully wired HTTP server)
//
//nolint:unused
var _ = kessoku.Inject[*http.Server](
	"newServer",
	kessoku.Provide(newOAuthClientApp),
	kessoku.Provide(service.NewOAuth),
	kessoku.Bind[server.OAuthHandler](kessoku.Provide(func(svc *service.OAuth) server.OAuthHandler { return svc })),
	kessoku.Provide(server.NewHandler),
	kessoku.Provide(func(cfg *data.Config, handler http.Handler) *http.Server {
		return &http.Server{
			Addr:    cfg.Server.Addr,
			Handler: h2c.NewHandler(handler, &http2.Server{}),
		}
	}),
)

func newOAuthClientApp(cfg *data.Config) *oauth.ClientApp {
	scopes := []string{
		"atproto",
	}

	var clientCfg oauth.ClientConfig
	if cfg.App.IsLocal() {
		clientCfg = oauth.NewLocalhostConfig(cfg.App.URL, scopes)
	} else {
		clientCfg = oauth.NewPublicConfig(
			path.Join(cfg.App.URL, "oauth/client-metadata.json"),
			cfg.App.URL,
			scopes,
		)
	}

	return oauth.NewClientApp(&clientCfg, oauth.NewMemStore())
}
