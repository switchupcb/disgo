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

	// Sessions contains sessions a bot uses to interact with the Discord Gateway.
	Sessions []*Session

	// Handlers represents a bot's event handlers.
	Handlers *Handlers

	// Config represents parameters used to perform various actions by the client.
	Config *Config
}

// Authentication represents authentication parameters required to authenticate the bot.
// https://discord.com/developers/docs/reference#authentication
type Authentication struct {
	// Token represents the Authentication Token used to authenticate the bot.
	Token string

	// TokenType represents the type of the Authentication Token.
	TokenType string

	// Header represents a Token Authorization Header.
	Header string
}

// BotToken uses a given token to return a valid Authentication Object for a bot token type.
func BotToken(token string) *Authentication {
	return &Authentication{
		Token:     token,
		TokenType: "Bot",
		Header:    "Bot " + token,
	}
}

// BearerToken uses a given token to return a valid Authentication Object for a bearer token type.
func BearerToken(token string) *Authentication {
	return &Authentication{
		Token:     token,
		TokenType: "Bearer",
		Header:    "Bearer" + token,
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

	// RedirectURI represents the registered url-encoded URL of the application.
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

	// GatewayPresenceUpdate represents the presence or status update of a bot.
	//
	// GatewayPresenceUpdate is used when the bot connects to a session.
	//
	// https://discord.com/developers/docs/topics/gateway#update-presence
	GatewayPresenceUpdate *GatewayPresenceUpdate

	// Intents represents a Discord Gateway Intent.
	//
	// You must specify a Gateway Intent in order to gain access to Events.
	//
	// https://discord.com/developers/docs/topics/gateway#gateway-intents
	Intents BitFlag

	// IntentSet represents a set of Discord Gateway Intents, that a bot needs
	// to gain access to Events (and specific Event Information).
	//
	// IntentSet is used for automatic intent calculation when a user adds an event handler.
	IntentSet map[BitFlag]bool
}

// Default Configuration Values.
const (
	defaultUserAgent      = "DiscordBot (https://" + module + ", v" + VersionDiscordAPI + ")"
	defaultRequestTimeout = time.Second * 3
)

// DefaultConfig returns a default client configuration.
func DefaultConfig() *Config {
	c := new(Config)
	c.Client = new(fasthttp.Client)
	c.Client.Name = defaultUserAgent
	c.Timeout = defaultRequestTimeout
	c.Retries = 1
	c.IntentSet = make(map[BitFlag]bool)
	c.GatewayPresenceUpdate = new(GatewayPresenceUpdate)

	return c
}
