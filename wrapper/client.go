package wrapper

import (
	"time"

	"github.com/valyala/fasthttp"
)

// Client represents a Discord Application.
type Client struct {
	ApplicationID string

	// Authentication contains parameters required to authenticate the bot.
	Authentication *Authentication

	// Authorization contains parameters required to authorize a client's access to resources.
	Authorization *Authorization

	// Config represents parameters used to perform various actions by the client.
	Config *Config
}

// Authentication represents authentication parameters required to authenticate the bot.
// https://discord.com/developers/docs/reference#authentication
type Authentication struct {
	// Header represents a Token Authorization Header.
	Header string
}

// BotToken generates a Bot Token Authorization Header.
func BotToken(token string) *Authentication {
	return &Authentication{
		Header: "Bot " + token,
	}
}

// BearerToken generates a Bearer Token Authorization Header.
func BearerToken(token string) *Authentication {
	return &Authentication{
		"Bearer" + token,
	}
}

// Authorization represents authorization parameters required to authorize a client's access to resources.
type Authorization struct {
	// ClientID represents the application's client_id.
	ClientID string

	// ClientSecret represents the application's client_secret.
	ClientSecret string

	// Scopes represents a list of OAuth2 scopes.
	Scopes []string

	// RedirectURI represents the registered URL of the application.
	//
	// The URL should be non-url-encoded (i.e "https://localhost"),
	// NOT url-encoded (i.e "https%3A%2F%2Flocalhost").
	RedirectURI string

	// state represents the state parameter used to prevent CSRF and Clickjacking.
	// https://discord.com/developers/docs/topics/oauth2#state-and-security
	State string

	// prompt controls how the authorization flow handles existing authorizations.
	Prompt string
}

// Config represents parameters used to perform various actions by the client.
type Config struct {
	// Client is used to send requests.
	//
	// Use Client to set a custom User-Agent in the HTTP Request Header.
	// https://discord.com/developers/docs/reference#user-agent
	//
	// https://pkg.go.dev/github.com/valyala/fasthttp#Client
	Client *fasthttp.Client

	// Timeout represents the amount of time a request will wait for a response.
	Timeout time.Duration

	// Retries represents the amount of time a request will retry a bad gateway.
	Retries int

	// RateLimiter represents an object that provides rate limit functionality.
	RateLimiter RateLimiter
}

// Default Configuration Values.
const (
	defaultUserAgent      = "DiscordBot (https://github.com/switchupcb/disgo, " + "v" + VersionDiscordAPI + ")"
	defaultRequestTimeout = time.Second * 3
)

// DefaultConfig returns a default client configuration.
func DefaultConfig() *Config {
	c := new(Config)
	c.Client = new(fasthttp.Client)
	c.Client.Name = defaultUserAgent
	c.Timeout = defaultRequestTimeout
	c.Retries = 0
	c.RateLimiter = RateLimit{}

	return c
}
