package wrapper

import (
	"time"

	"github.com/valyala/fasthttp"
)

// Default Configuration Values.
const (
	module           = "github.com/switchupcb/disgo"
	defaultUserAgent = "DiscordBot (https://" + module + ", v" + VersionDiscordAPI + ")"
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
	// Request holds configuration variables that pertain to the Discord HTTP API.
	Request Request

	// Gateway holds configuration variables that pertain to the Discord Gateway.
	Gateway Gateway
}

// DefaultConfig returns a default client configuration.
func DefaultConfig() *Config {
	c := new(Config)
	c.Request = DefaultRequest()
	c.Gateway = DefaultGateway(false)

	return c
}

// Request represents Discord Request parameters used to perform various actions by the client.
type Request struct {
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

const (
	// defaultRequestTimeout represents the default amount of time to wait on a request.
	defaultRequestTimeout = time.Second

	// totalRoutes represents the total amount of Discord HTTP Routes (174) + the Global Route (1).
	totalRoutes = 175
)

// DefaultRequest returns a default Request configuration.
func DefaultRequest() Request {
	// configure the client.
	client := new(fasthttp.Client)
	client.Name = defaultUserAgent

	// configure the rate limiter.
	ratelimiter := &RateLimit{ //nolint:exhaustruct
		ids:     make(map[uint16]string, totalRoutes),
		buckets: make(map[string]*Bucket, totalRoutes),
		entries: make(map[string]int, totalRoutes),
	}

	// https://discord.com/developers/docs/topics/rate-limits#global-rate-limit
	ratelimiter.SetBucket(
		0, &Bucket{ //nolint:exhaustruct
			Limit:     FlagGlobalRateLimitRequest,
			Remaining: FlagGlobalRateLimitRequest,
		},
	)

	return Request{
		Client:      client,
		Timeout:     defaultRequestTimeout,
		Retries:     1,
		RateLimiter: ratelimiter,
	}
}

// Gateway represents Discord Gateway parameters used to perform various actions by the client.
type Gateway struct {
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

	// GatewayPresenceUpdate represents the presence or status update of a bot.
	//
	// GatewayPresenceUpdate is used when the bot connects to a session.
	//
	// https://discord.com/developers/docs/topics/gateway#update-presence
	GatewayPresenceUpdate *GatewayPresenceUpdate

	// RateLimiter represents an object that provides rate limit functionality.
	RateLimiter RateLimiter
}

const (
	// totalIntents represents the total amount of Discord Intents.
	totalIntents = 19

	// totalGatewayBuckets represents the total amount of Discord Gateway Rate Limits.
	totalGatewayBuckets = 2
)

// DefaultGateway returns a default Gateway configuration.
//
// When privileged intents are enabled, the MESSAGE_CONTENT intent will be included.
//
// MESSAGE_CONTENT is required to receive message content fields
// (content, attachments, embeds, and components).
//
// https://discord.com/developers/docs/topics/gateway#privileged-intents
func DefaultGateway(privileged bool) Gateway {
	// configure the rate limiter.
	ratelimiter := &RateLimit{ //nolint:exhaustruct
		ids:     make(map[uint16]string, totalGatewayBuckets),
		buckets: make(map[string]*Bucket, totalGatewayBuckets),
	}

	// https://discord.com/developers/docs/topics/gateway#rate-limiting
	ratelimiter.SetBucket(
		0, &Bucket{ //nolint:exhaustruct
			Limit:     FlagGlobalRateLimitGateway,
			Remaining: FlagGlobalRateLimitGateway,
			Expiry:    time.Now().Add(FlagGlobalRateLimitGatewayInterval),
		},
	)

	if privileged {
		is := make(map[BitFlag]bool, totalIntents)
		is[FlagIntentMESSAGE_CONTENT] = true

		return Gateway{
			Intents:               FlagIntentMESSAGE_CONTENT,
			IntentSet:             is,
			GatewayPresenceUpdate: new(GatewayPresenceUpdate),
			RateLimiter:           ratelimiter,
		}
	} else {
		return Gateway{
			Intents:               0,
			IntentSet:             make(map[BitFlag]bool, totalIntents),
			GatewayPresenceUpdate: new(GatewayPresenceUpdate),
			RateLimiter:           ratelimiter,
		}
	}
}
