package wrapper

import "github.com/valyala/fasthttp"

// Client represents a Discord Application.
type Client struct {
	Authorization *Authorization

	ApplicationID string
	client        *fasthttp.Client
}

// Authorization represents authorization parameters required to authorize a client.
type Authorization struct {
	// ClientID represents the application's client_id.
	ClientID string

	// ClientSecret represents the application's client_secret.
	ClientSecret string

	// Scopes represents a list of OAuth2 scopes.
	Scopes []string

	// RedirectURI represents the registered url-encoded URL of the application.
	RedirectURI string

	// state represents the state parameter used to prevent CSRF and Clickjacking.
	// https://discord.com/developers/docs/topics/oauth2#state-and-security
	State string

	// prompt controls how the authorization flow handles existing authorizations.
	Prompt string
}
