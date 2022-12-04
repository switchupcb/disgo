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

	// Config represents parameters used to perform various actions by the client.
	Config *Config

	// Handlers represents a bot's event handlers.
	Handlers *Handlers

	// Sessions contains sessions a bot uses to interact with the Discord Gateway.
	Sessions []*Session
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

	// Scopes represents a list of OAuth2 scopes.
	Scopes []string
}

// Config represents parameters used to perform various actions by the client.
type Config struct {
	// Gateway holds configuration variables that pertain to the Discord Gateway.
	Gateway Gateway

	// Request holds configuration variables that pertain to the Discord HTTP API.
	Request Request
}

// DefaultConfig returns a default client configuration.
func DefaultConfig() *Config {
	c := new(Config)
	c.Request = DefaultRequest()
	c.Gateway = DefaultGateway()

	return c
}

// Request represents Discord Request parameters used to perform various actions by the client.
type Request struct {
	// RateLimiter represents an object that provides rate limit functionality.
	RateLimiter RateLimiter

	// Client is used to send requests.
	//
	// Use Client to set a custom User-Agent in the HTTP Request Header.
	// https://discord.com/developers/docs/reference#user-agent
	//
	// https://pkg.go.dev/github.com/valyala/fasthttp#Client
	Client *fasthttp.Client

	// Timeout represents the amount of time a request will wait for a response.
	Timeout time.Duration

	// Retries represents the number of times a request may be retried upon failure.
	//
	// A request is ONLY retried when a Bad Gateway or Rate Limit is encountered.
	Retries int

	// RetryShared determines the behavior of a request when
	// a (shared) per-resource rate limit is hit.
	//
	// set RetryShared to true (default) to retry a request (within the per-route rate limit)
	// until it's successful or until it experiences a non-shared 429 status code.
	RetryShared bool
}

const (
	// defaultRequestTimeout represents the default amount of time to wait on a request.
	defaultRequestTimeout = time.Second
)

// DefaultRequest returns a Default Request configuration.
func DefaultRequest() Request {
	// configure the client.
	client := new(fasthttp.Client)
	client.Name = defaultUserAgent

	// configure the rate limiter.
	ratelimiter := &RateLimit{ //nolint:exhaustruct
		ids:     make(map[string]string, len(RouteIDs)),
		buckets: make(map[string]*Bucket, len(RouteIDs)),
		entries: make(map[string]int, len(RouteIDs)),
	}

	ratelimiter.DefaultBucket = &Bucket{ //nolint:exhaustruct
		Limit: 1,
	}

	// https://discord.com/developers/docs/topics/rate-limits#global-rate-limit
	ratelimiter.SetBucket(
		GlobalRateLimitRouteID, &Bucket{ //nolint:exhaustruct
			Limit:     FlagGlobalRateLimitRequest,
			Remaining: FlagGlobalRateLimitRequest,
		},
	)

	return Request{
		RateLimiter: ratelimiter,
		Client:      client,
		Timeout:     defaultRequestTimeout,
		Retries:     1,
		RetryShared: true,
	}
}

// Gateway represents Discord Gateway parameters used to perform various actions by the client.
type Gateway struct {
	// RateLimiter represents an object that provides rate limit functionality.
	RateLimiter RateLimiter

	// Intents represents a Discord Gateway Intent.
	//
	// You must specify a Gateway Intent in order to receive specific information from an event.
	//
	// https://discord.com/developers/docs/topics/gateway#gateway-intents
	Intents BitFlag

	// IntentSet represents a set of Discord Gateway Intents, that a bot needs
	// to receive specific information from an event.
	//
	// IntentSet is used for automatic intent calculation when a user adds an event handler.
	IntentSet map[BitFlag]bool

	// GatewayPresenceUpdate represents the presence or status update of a bot.
	//
	// GatewayPresenceUpdate is used when the bot connects to a session.
	//
	// https://discord.com/developers/docs/topics/gateway#update-presence
	GatewayPresenceUpdate *GatewayPresenceUpdate
}

const (
	// totalIntents represents the total amount of Discord Intents.
	totalIntents = 19

	// totalGatewayBuckets represents the total amount of Discord Gateway Rate Limits.
	totalGatewayBuckets = 2
)

// DefaultGateway returns a default Gateway configuration.
//
// Privileged Intents are disabled by default.
// https://discord.com/developers/docs/topics/gateway#privileged-intents
func DefaultGateway() Gateway {
	// configure the rate limiter.
	ratelimiter := &RateLimit{ //nolint:exhaustruct
		ids:     make(map[string]string, totalGatewayBuckets),
		buckets: make(map[string]*Bucket, totalGatewayBuckets),
	}

	ratelimiter.DefaultBucket = &Bucket{ //nolint:exhaustruct
		Limit: 1,
	}

	// https://discord.com/developers/docs/topics/gateway#rate-limiting
	ratelimiter.SetBucket(
		GlobalRateLimitRouteID, &Bucket{ //nolint:exhaustruct
			Limit:     FlagGlobalRateLimitGateway,
			Remaining: FlagGlobalRateLimitGateway,
			Expiry:    time.Now().Add(FlagGlobalRateLimitGatewayInterval),
		},
	)

	// disable Privileged Intents.
	// https://discord.com/developers/docs/topics/gateway#privileged-intents
	intentSet := make(map[BitFlag]bool, totalIntents)
	for privilegedIntent := range PrivilegedIntents {
		intentSet[privilegedIntent] = true
	}

	return Gateway{
		Intents:               0,
		IntentSet:             intentSet,
		GatewayPresenceUpdate: new(GatewayPresenceUpdate),
		RateLimiter:           ratelimiter,
	}
}

// EnableIntent enables an intent.
//
// This function does NOT check whether the intent is already enabled.
// Use the Gateway.IntentSet to check whether the intent is already enabled.
func (g *Gateway) EnableIntent(intent BitFlag) {
	g.IntentSet[FlagIntentAUTO_MODERATION_CONFIGURATION] = true
	g.Intents |= intent
}

// DisableIntent disables an intent.
//
// Disclaimer: The Bitwise OR operation (used) to add an intent is a DESTRUCTIVE operation.
//
// This means that it can NOT be reversed. As a result, this function will NOT remove
// an intent that is already enabled.
func (g Gateway) DisableIntent(intent BitFlag) {
	g.IntentSet[intent] = true
}
