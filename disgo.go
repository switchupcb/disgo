//go:generate bundle -o disgo.go -dst . -pkg disgo -prefix "" ./wrapper

package disgo

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	json "github.com/goccy/go-json"
	"github.com/gorilla/schema"
	"github.com/rs/xid"
	"github.com/rs/zerolog"
	"github.com/switchupcb/disgo/wrapper/socket"
	"github.com/switchupcb/websocket"
	"github.com/valyala/fasthttp"
	"golang.org/x/sync/errgroup"
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
	ClientID     string
	ClientSecret string
	RedirectURI  string
	State        string
	Prompt       string
	Scopes       []string
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
	c.Gateway = DefaultGateway()

	return c
}

// Request represents Discord Request parameters used to perform various actions by the client.
type Request struct {
	RateLimiter RateLimiter
	Client      *fasthttp.Client
	Timeout     time.Duration
	Retries     int
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
		Client:      client,
		Timeout:     defaultRequestTimeout,
		Retries:     1,
		RetryShared: true,
		RateLimiter: ratelimiter,
	}
}

// Gateway represents Discord Gateway parameters used to perform various actions by the client.
type Gateway struct {
	RateLimiter           RateLimiter
	IntentSet             map[BitFlag]bool
	GatewayPresenceUpdate *GatewayPresenceUpdate
	Intents               BitFlag
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

// Gateway Opcodes
// https://discord.com/developers/docs/topics/opcodes-and-status-codes#gateway-gateway-opcodes
const (
	FlagGatewayOpcodeDispatch            = 0
	FlagGatewayOpcodeHeartbeat           = 1
	FlagGatewayOpcodeIdentify            = 2
	FlagGatewayOpcodePresenceUpdate      = 3
	FlagGatewayOpcodeVoiceStateUpdate    = 4
	FlagGatewayOpcodeResume              = 6
	FlagGatewayOpcodeReconnect           = 7
	FlagGatewayOpcodeRequestGuildMembers = 8
	FlagGatewayOpcodeInvalidSession      = 9
	FlagGatewayOpcodeHello               = 10
	FlagGatewayOpcodeHeartbeatACK        = 11
)

// Gateway Close Event Codes
// https://discord.com/developers/docs/topics/opcodes-and-status-codes#gateway-gateway-close-event-codes
type GatewayCloseEventCode struct {
	Description string
	Explanation string
	Code        int
	Reconnect   bool
}

var (
	FlagGatewayCloseEventCodeUnknownError = GatewayCloseEventCode{
		Code:        4000,
		Description: "Unknown error",
		Explanation: "We're not sure what went wrong. Try reconnecting?",
		Reconnect:   true,
	}

	FlagGatewayCloseEventCodeUnknownOpcode = GatewayCloseEventCode{
		Code:        4001,
		Description: "Unknown opcode",
		Explanation: "You sent an invalid Gateway opcode or an invalid payload for an opcode. Don't do that!",
		Reconnect:   true,
	}

	FlagGatewayCloseEventCodeDecodeError = GatewayCloseEventCode{
		Code:        4002,
		Description: "Decode error",
		Explanation: "You sent an invalid payload to us. Don't do that!",
		Reconnect:   true,
	}

	FlagGatewayCloseEventCodeNotAuthenticated = GatewayCloseEventCode{
		Code:        4003,
		Description: "Not authenticated",
		Explanation: "You sent us a payload prior to identifying.",
		Reconnect:   true,
	}

	FlagGatewayCloseEventCodeAuthenticationFailed = GatewayCloseEventCode{
		Code:        4004,
		Description: "Authentication failed",
		Explanation: "The account token sent with your identify payload is incorrect.",
		Reconnect:   false,
	}

	FlagGatewayCloseEventCodeAlreadyAuthenticated = GatewayCloseEventCode{
		Code:        4005,
		Description: "Already authenticated",
		Explanation: "You sent more than one identify payload. Don't do that!",
		Reconnect:   true,
	}

	FlagGatewayCloseEventCodeInvalidSeq = GatewayCloseEventCode{
		Code:        4007,
		Description: "Invalid seq",
		Explanation: "The sequence sent when resuming the session was invalid. Reconnect and start a new session.",
		Reconnect:   true,
	}

	FlagGatewayCloseEventCodeRateLimited = GatewayCloseEventCode{
		Code:        4008,
		Description: "Rate limited.",
		Explanation: "You're sending payloads to us too quickly. Slow it down! You will be disconnected on receiving this.",
		Reconnect:   true,
	}

	FlagGatewayCloseEventCodeSessionTimed = GatewayCloseEventCode{
		Code:        4009,
		Description: "Session timed out",
		Explanation: "Your session timed out. Reconnect and start a new one.",
		Reconnect:   true,
	}

	FlagGatewayCloseEventCodeInvalidShard = GatewayCloseEventCode{
		Code:        4010,
		Description: "Invalid shard",
		Explanation: "You sent us an invalid shard when identifying.",
		Reconnect:   false,
	}

	FlagGatewayCloseEventCodeShardingRequired = GatewayCloseEventCode{
		Code:        4011,
		Description: "Sharding required",
		Explanation: "The session would have handled too many guilds - you are required to shard your connection in order to connect.",
		Reconnect:   false,
	}

	FlagGatewayCloseEventCodeInvalidAPIVersion = GatewayCloseEventCode{
		Code:        4012,
		Description: "Invalid API version",
		Explanation: "You sent an invalid version for the gateway.",
		Reconnect:   false,
	}

	FlagGatewayCloseEventCodeInvalidIntent = GatewayCloseEventCode{
		Code:        4013,
		Description: "Invalid intent(s)",
		Explanation: "You sent an invalid intent for a Gateway Intent. You may have incorrectly calculated the bitwise value.",
		Reconnect:   false,
	}

	FlagGatewayCloseEventCodeDisallowedIntent = GatewayCloseEventCode{
		Code:        4014,
		Description: "Disallowed intent(s)",
		Explanation: "You sent a disallowed intent for a Gateway Intent. You may have tried to specify an intent that you have not enabled or are not approved for.",
		Reconnect:   false,
	}

	GatewayCloseEventCodes = map[int]*GatewayCloseEventCode{
		FlagGatewayCloseEventCodeUnknownError.Code:         &FlagGatewayCloseEventCodeUnknownError,
		FlagGatewayCloseEventCodeUnknownOpcode.Code:        &FlagGatewayCloseEventCodeUnknownOpcode,
		FlagGatewayCloseEventCodeDecodeError.Code:          &FlagGatewayCloseEventCodeDecodeError,
		FlagGatewayCloseEventCodeNotAuthenticated.Code:     &FlagGatewayCloseEventCodeNotAuthenticated,
		FlagGatewayCloseEventCodeAuthenticationFailed.Code: &FlagGatewayCloseEventCodeAuthenticationFailed,
		FlagGatewayCloseEventCodeAlreadyAuthenticated.Code: &FlagGatewayCloseEventCodeAlreadyAuthenticated,
		FlagGatewayCloseEventCodeInvalidSeq.Code:           &FlagGatewayCloseEventCodeInvalidSeq,
		FlagGatewayCloseEventCodeRateLimited.Code:          &FlagGatewayCloseEventCodeRateLimited,
		FlagGatewayCloseEventCodeSessionTimed.Code:         &FlagGatewayCloseEventCodeSessionTimed,
		FlagGatewayCloseEventCodeInvalidShard.Code:         &FlagGatewayCloseEventCodeInvalidShard,
		FlagGatewayCloseEventCodeInvalidAPIVersion.Code:    &FlagGatewayCloseEventCodeInvalidAPIVersion,
		FlagGatewayCloseEventCodeInvalidIntent.Code:        &FlagGatewayCloseEventCodeInvalidIntent,
		FlagGatewayCloseEventCodeDisallowedIntent.Code:     &FlagGatewayCloseEventCodeDisallowedIntent,
	}
)

// Client Close Event Codes
// https://discord.com/developers/docs/topics/gateway#disconnections
var (
	FlagClientCloseEventCodeNormal = 1000
	FlagClientCloseEventCodeAway   = 1001

	// https://www.rfc-editor.org/rfc/rfc6455#section-7.4.1
	FlagClientCloseEventCodeReconnect = 3000
)

// Voice Opcodes
// https://discord.com/developers/docs/topics/opcodes-and-status-codes#voice-voice-opcodes
const (
	FlagVoiceOpcodeIdentify           = 0
	FlagVoiceOpcodeSelectProtocol     = 1
	FlagVoiceOpcodeReadyServer        = 2
	FlagVoiceOpcodeHeartbeat          = 3
	FlagVoiceOpcodeSessionDescription = 4
	FlagVoiceOpcodeSpeaking           = 5
	FlagVoiceOpcodeHeartbeatACK       = 6
	FlagVoiceOpcodeResume             = 7
	FlagVoiceOpcodeHello              = 8
	FlagVoiceOpcodeResumed            = 9
	FlagVoiceOpcodeClientDisconnect   = 13
)

// Voice Close Event Codes
// https://discord.com/developers/docs/topics/opcodes-and-status-codes#voice-voice-close-event-codes
type VoiceCloseEventCode struct {
	Description string
	Explanation string
	Code        int
}

var (
	FlagVoiceCloseEventCodeUnknownOpcode = VoiceCloseEventCode{
		Code:        4001,
		Description: "Unknown opcode",
		Explanation: "You sent an invalid opcode.",
	}

	FlagVoiceCloseEventCodeFailedDecode = VoiceCloseEventCode{
		Code:        4002,
		Description: "Failed to decode payload",
		Explanation: "You sent a invalid payload in your identifying to the Gateway.",
	}

	FlagVoiceCloseEventCodeNotAuthenticated = VoiceCloseEventCode{
		Code:        4003,
		Description: "Not authenticated",
		Explanation: "You sent a payload before identifying with the Gateway.",
	}

	FlagVoiceCloseEventCodeAuthenticationFailed = VoiceCloseEventCode{
		Code:        4004,
		Description: "Authentication failed",
		Explanation: "The token you sent in your identify payload is incorrect.",
	}

	FlagVoiceCloseEventCodeAlreadyAuthenticated = VoiceCloseEventCode{
		Code:        4005,
		Description: "Already authenticated",
		Explanation: "You sent more than one identify payload. Stahp.",
	}

	FlagVoiceCloseEventCodeInvalidSession = VoiceCloseEventCode{
		Code:        4006,
		Description: "Session no longer valid",
		Explanation: "Your session is no longer valid.",
	}

	FlagVoiceCloseEventCodeSessionTimeout = VoiceCloseEventCode{
		Code:        4009,
		Description: "Session timeout",
		Explanation: "Your session has timed out.",
	}

	FlagVoiceCloseEventCodeServerNotFound = VoiceCloseEventCode{
		Code:        4011,
		Description: "Server not found",
		Explanation: "We can't find the server you're trying to connect to.",
	}

	FlagVoiceCloseEventCodeUnknownProtocol = VoiceCloseEventCode{
		Code:        4012,
		Description: "Unknown protocol",
		Explanation: "We didn't recognize the protocol you sent.",
	}

	FlagVoiceCloseEventCodeDisconnectedChannel = VoiceCloseEventCode{
		Code:        4014,
		Description: "Disconnected",
		Explanation: "Channel was deleted, you were kicked, voice server changed, or the main gateway session was dropped. Don't reconnect.",
	}

	FlagVoiceCloseEventCodeVoiceServerCrash = VoiceCloseEventCode{
		Code:        4015,
		Description: "Voice server crashed",
		Explanation: "The server crashed. Our bad! Try resuming.",
	}

	FlagVoiceCloseEventCodeUnknownEncryptionMode = VoiceCloseEventCode{
		Code:        4016,
		Description: "Unknown encryption mode",
		Explanation: "We didn't recognize your encryption.",
	}

	VoiceCloseEventCodes = map[int]*VoiceCloseEventCode{
		FlagVoiceCloseEventCodeUnknownOpcode.Code:         &FlagVoiceCloseEventCodeUnknownOpcode,
		FlagVoiceCloseEventCodeFailedDecode.Code:          &FlagVoiceCloseEventCodeFailedDecode,
		FlagVoiceCloseEventCodeNotAuthenticated.Code:      &FlagVoiceCloseEventCodeNotAuthenticated,
		FlagVoiceCloseEventCodeAuthenticationFailed.Code:  &FlagVoiceCloseEventCodeAuthenticationFailed,
		FlagVoiceCloseEventCodeAlreadyAuthenticated.Code:  &FlagVoiceCloseEventCodeAlreadyAuthenticated,
		FlagVoiceCloseEventCodeInvalidSession.Code:        &FlagVoiceCloseEventCodeInvalidSession,
		FlagVoiceCloseEventCodeSessionTimeout.Code:        &FlagVoiceCloseEventCodeSessionTimeout,
		FlagVoiceCloseEventCodeServerNotFound.Code:        &FlagVoiceCloseEventCodeServerNotFound,
		FlagVoiceCloseEventCodeUnknownProtocol.Code:       &FlagVoiceCloseEventCodeUnknownProtocol,
		FlagVoiceCloseEventCodeDisconnectedChannel.Code:   &FlagVoiceCloseEventCodeDisconnectedChannel,
		FlagVoiceCloseEventCodeVoiceServerCrash.Code:      &FlagVoiceCloseEventCodeVoiceServerCrash,
		FlagVoiceCloseEventCodeUnknownEncryptionMode.Code: &FlagVoiceCloseEventCodeUnknownEncryptionMode,
	}
)

// HTTP Response Codes
// https://discord.com/developers/docs/topics/opcodes-and-status-codes#http-http-response-codes
const (
	FlagHTTPResponseCodeOK                 = 200
	FlagHTTPResponseCodeCREATED            = 201
	FlagHTTPResponseCodeNOCONTENT          = 204
	FlagHTTPResponseCodeNOTMODIFIED        = 304
	FlagHTTPResponseCodeBADREQUEST         = 400
	FlagHTTPResponseCodeUNAUTHORIZED       = 401
	FlagHTTPResponseCodeFORBIDDEN          = 403
	FlagHTTPResponseCodeNOTFOUND           = 404
	FlagHTTPResponseCodeMETHODNOTALLOWED   = 405
	FlagHTTPResponseCodeTOOMANYREQUESTS    = 429
	FlagHTTPResponseCodeGATEWAYUNAVAILABLE = 502
	FlagHTTPResponseCodeSERVERERROR        = 500 // 5xx (500 Not Guaranteed)
)

var (
	HTTPResponseCodes = map[int]string{
		FlagHTTPResponseCodeOK:                 "The request completed successfully.",
		FlagHTTPResponseCodeCREATED:            "The entity was created successfully.",
		FlagHTTPResponseCodeNOCONTENT:          "The request completed successfully but returned no content.",
		FlagHTTPResponseCodeNOTMODIFIED:        "The entity was not modified (no action was taken).",
		FlagHTTPResponseCodeBADREQUEST:         "The request was improperly formatted, or the server couldn't understand it.",
		FlagHTTPResponseCodeUNAUTHORIZED:       "The Authorization header was missing or invalid.",
		FlagHTTPResponseCodeFORBIDDEN:          "The Authorization token you passed did not have permission to the resource.",
		FlagHTTPResponseCodeNOTFOUND:           "The resource at the location specified doesn't exist.",
		FlagHTTPResponseCodeMETHODNOTALLOWED:   "The HTTP method used is not valid for the location specified.",
		FlagHTTPResponseCodeTOOMANYREQUESTS:    "You are being rate limited, see Rate Limits.",
		FlagHTTPResponseCodeGATEWAYUNAVAILABLE: "There was not a gateway available to process your request. Wait a bit and retry.",
		FlagHTTPResponseCodeSERVERERROR:        "The server had an error processing your request (these are rare).",
	}
)

// JSON Error Codes
// https://discord.com/developers/docs/topics/opcodes-and-status-codes#json-json-error-codes
var (
	JSONErrorCodes = map[int]string{
		0:      "General error (such as a malformed request body, amongst other things)",
		10001:  "Unknown account",
		10002:  "Unknown application",
		10003:  "Unknown channel",
		10004:  "Unknown guild",
		10005:  "Unknown integration",
		10006:  "Unknown invite",
		10007:  "Unknown member",
		10008:  "Unknown message",
		10009:  "Unknown permission overwrite",
		10010:  "Unknown provider",
		10011:  "Unknown role",
		10012:  "Unknown token",
		10013:  "Unknown user",
		10014:  "Unknown emoji",
		10015:  "Unknown webhook",
		10016:  "Unknown webhook service",
		10020:  "Unknown session",
		10026:  "Unknown ban",
		10027:  "Unknown SKU",
		10028:  "Unknown Store Listing",
		10029:  "Unknown entitlement",
		10030:  "Unknown build",
		10031:  "Unknown lobby",
		10032:  "Unknown branch",
		10033:  "Unknown store directory layout",
		10036:  "Unknown redistributable",
		10038:  "Unknown gift code",
		10049:  "Unknown stream",
		10050:  "Unknown premium server subscribe cooldown",
		10057:  "Unknown guild template",
		10059:  "Unknown discoverable server category",
		10060:  "Unknown sticker",
		10062:  "Unknown interaction",
		10063:  "Unknown application command",
		10065:  "Unknown voice state",
		10066:  "Unknown application command permissions",
		10067:  "Unknown Stage Instance",
		10068:  "Unknown Guild Member Verification Form",
		10069:  "Unknown Guild Welcome Screen",
		10070:  "Unknown Guild Scheduled Event",
		10071:  "Unknown Guild Scheduled Event User",
		10087:  "Unknown Tag",
		20001:  "Bots cannot use this endpoint",
		20002:  "Only bots can use this endpoint",
		20009:  "Explicit content cannot be sent to the desired recipient(s)",
		20012:  "You are not authorized to perform this action on this application",
		20016:  "This action cannot be performed due to slowmode rate limit",
		20018:  "Only the owner of this account can perform this action",
		20022:  "This message cannot be edited due to announcement rate limits",
		20024:  "The owner of this account is not old enough to join an NSFW server.",
		20028:  "The channel you are writing has hit the write rate limit",
		20029:  "The write action you are performing on the server has hit the write rate limit",
		20031:  "Your Stage topic, server name, server description, or channel names contain words that are not allowed",
		20035:  "Guild premium subscription level too low",
		30001:  "Maximum number of guilds reached (100)",
		30002:  "Maximum number of friends reached (1000)",
		30003:  "Maximum number of pins reached for the channel (50)",
		30004:  "Maximum number of recipients reached (10)",
		30005:  "Maximum number of guild roles reached (250)",
		30007:  "Maximum number of webhooks reached (10)",
		30008:  "Maximum number of emojis reached",
		30010:  "Maximum number of reactions reached (20)",
		30013:  "Maximum number of guild channels reached (500)",
		30015:  "Maximum number of attachments in a message reached (10)",
		30016:  "Maximum number of invites reached (1000)",
		30018:  "Maximum number of animated emojis reached",
		30019:  "Maximum number of server members reached",
		30030:  "Maximum number of server categories has been reached (5)",
		30031:  "Guild already has a template",
		30032:  "Maximum number of application commands reached",
		30033:  "Max number of thread participants has been reached (1000)",
		30034:  "Max number of daily application command creates has been reached (200)",
		30035:  "Maximum number of bans for non-guild members have been exceeded",
		30037:  "Maximum number of bans fetches has been reached",
		30038:  "Maximum number of uncompleted guild scheduled events reached (100)",
		30039:  "Maximum number of stickers reached",
		30040:  "Maximum number of prune requests has been reached. Try again later",
		30042:  "Maximum number of guild widget settings updates has been reached. Try again later",
		30046:  "Maximum number of edits to messages older than 1 hour reached. Try again later",
		30047:  "Maximum number of pinned threads in a forum channel has been reached",
		30048:  "Maximum number of tags in a forum channel has been reached",
		30052:  "Bitrate is too high for channel of this type",
		40001:  "Unauthorized. Provide a valid token and try again",
		40002:  "You need to verify your account in order to perform this action",
		40003:  "You are opening direct messages too fast",
		40004:  "Send messages has been temporarily disabled",
		40005:  "Request entity too large. Try sending something smaller in size",
		40006:  "This feature has been temporarily disabled server-side",
		40007:  "The user is banned from this guild",
		40012:  "Connection has been revoked",
		40032:  "Target user is not connected to voice",
		40033:  "This message has already been crossposted",
		40041:  "An application command with that name already exists",
		40043:  "Application interaction failed to send",
		40058:  "Cannot send a message in a forum channel",
		40060:  "Interaction has already been acknowledged",
		40061:  "Tag names must be unique",
		40066:  "There are no tags available that can be set by non-moderators",
		40067:  "A tag is required to create a forum post in this channel",
		50001:  "Missing access",
		50002:  "Invalid account type",
		50003:  "Cannot execute action on a DM channel",
		50004:  "Guild widget disabled",
		50005:  "Cannot edit a message authored by another user",
		50006:  "Cannot send an empty message",
		50007:  "Cannot send messages to this user",
		50008:  "Cannot send messages in a non-text channel",
		50009:  "Channel verification level is too high for you to gain access",
		50010:  "OAuth2 application does not have a bot",
		50011:  "OAuth2 application limit reached",
		50012:  "Invalid OAuth2 state",
		50013:  "You lack permissions to perform that action",
		50014:  "Invalid authentication token provided",
		50015:  "Note was too long",
		50016:  "Provided too few or too many messages to delete. Must provide at least 2 and fewer than 100 messages to delete",
		50017:  "Invalid MFA Level",
		50019:  "A message can only be pinned to the channel it was sent in",
		50020:  "Invite code was either invalid or taken",
		50021:  "Cannot execute action on a system message",
		50024:  "Cannot execute action on this channel type",
		50025:  "Invalid OAuth2 access token provided",
		50026:  "Missing required OAuth2 scope",
		50027:  "Invalid webhook token provided",
		50028:  "Invalid role",
		50033:  "Invalid Recipient(s)",
		50034:  "A message provided was too old to bulk delete",
		50035:  "Invalid form body (returned for both application/json and multipart/form-data bodies), or invalid Content-Type provided",
		50036:  "An invite was accepted to a guild the application's bot is not in",
		50039:  "Invalid Activity Action",
		50041:  "Invalid API version provided",
		50045:  "File uploaded exceeds the maximum size",
		50046:  "Invalid file uploaded",
		50054:  "Cannot self-redeem this gift",
		50055:  "Invalid Guild",
		50068:  "Invalid message type",
		50070:  "Payment source required to redeem gift",
		50073:  "Cannot modify a system webhook",
		50074:  "Cannot delete a channel required for Community guilds",
		50081:  "Invalid sticker sent",
		50083:  "Tried to perform an operation on an archived thread, such as editing a message or adding a user to the thread",
		50084:  "Invalid thread notification settings",
		50085:  "before value is earlier than the thread creation date",
		50086:  "Community server channels must be text channels",
		50095:  "This server is not available in your location",
		50097:  "This server needs monetization enabled in order to perform this action",
		50101:  "This server needs more boosts to perform this action",
		50109:  "The request body contains invalid JSON.",
		50132:  "Ownership cannot be transferred to a bot user",
		50138:  "Failed to resize asset below the maximum size: 262144",
		50146:  "Uploaded file not found.",
		50600:  "You do not have permission to send this sticker.",
		60003:  "Two factor is required for this operation",
		80004:  "No users with DiscordTag exist",
		90001:  "Reaction was blocked",
		110001: "Application not yet available. Try again later",
		130000: "API resource is currently overloaded. Try again a little later",
		150006: "The Stage is already open",
		160002: "Cannot reply without permission to read message history",
		160004: "A thread has already been created for this message",
		160005: "Thread is locked",
		160006: "Maximum number of active threads reached",
		160007: "Maximum number of active announcement threads reached",
		170001: "Invalid JSON for uploaded Lottie file",
		170002: "Uploaded Lotties cannot contain rasterized images such as PNG or JPEG",
		170003: "Sticker maximum framerate exceeded",
		170004: "Sticker frame count exceeds maximum of 1000 frames",
		170005: "Lottie animation maximum dimensions exceeded",
		170006: "Sticker frame rate is either too small or too large",
		170007: "Sticker animation duration exceeds maximum of 5 seconds",
		180000: "Cannot update a finished event",
		180002: "Failed to create stage needed for stage event",
		200000: "Message was blocked by automatic moderation",
		200001: "Title was blocked by automatic moderation",
		220001: "Webhooks posted to forum channels must have a thread_name or thread_id",
		220002: " Webhooks posted to forum channels cannot have both a thread_name and thread_id",
		220003: "Webhooks can only create threads in forum channels",
		240000: "Message blocked by harmful links filter",
	}
)

// RPC Error Codes
// https://discord.com/developers/docs/topics/opcodes-and-status-codes#rpc-rpc-error-codes
const (
	FlagRPCErrorCodeUnknownError                    = 1000
	FlagRPCErrorCodeInvalidPayload                  = 4000
	FlagRPCErrorCodeInvalidCommand                  = 4002
	FlagRPCErrorCodeInvalidGuild                    = 4003
	FlagRPCErrorCodeInvalidEvent                    = 4004
	FlagRPCErrorCodeInvalidChannel                  = 4005
	FlagRPCErrorCodeInvalidPermissions              = 4006
	FlagRPCErrorCodeInvalidClientID                 = 4007
	FlagRPCErrorCodeInvalidOrigin                   = 4008
	FlagRPCErrorCodeInvalidToken                    = 4009
	FlagRPCErrorCodeInvalidUser                     = 4010
	FlagRPCErrorCodeOAuth2Error                     = 5000
	FlagRPCErrorCodeSelectChannelTimedOut           = 5001
	FlagRPCErrorCodeGET_GUILDTimedOut               = 5002
	FlagRPCErrorCodeSelectVoiceForceRequired        = 5003
	FlagRPCErrorCodeCaptureShortcutAlreadyListening = 5004
)

// RPC Close Event Codes
// https://discord.com/developers/docs/topics/opcodes-and-status-codes#rpc-rpc-close-event-codes
const (
	FlagRPCCloseEventCodeInvalidClientID = 4000
	FlagRPCCloseEventCodeInvalidOrigin   = 4001
	FlagRPCCloseEventCodeRateLimited     = 4002
	FlagRPCCloseEventCodeTokenRevoked    = 4003
	FlagRPCCloseEventCodeInvalidVersion  = 4004
	FlagRPCCloseEventCodeInvalidEncoding = 4005
)

// Flag represents an (unused) alias for a Discord API Flag ranging from 0 - 255.
type Flag uint8

// BitFlag represents an alias for a Discord API Bitwise Flag denoted by 1 << x.
type BitFlag uint64

// File represents a file attachment.
type File struct {
	Name        string
	ContentType string
	Data        []byte
}

// Nonce represents a Discord nonce (integer or string).
type Nonce string

// Value represents a value (string, integer, or double).
type Value string

// PointerIndicator represents a Dasgo double pointer value indicator.
type PointerIndicator uint8

const (
	// IsValueNothing indicates that the field was not provided.
	//
	// The double pointer is nil.
	IsValueNothing PointerIndicator = 0

	// IsValueNull indicates the the field was provided with a null value.
	//
	// The double pointer points to a nil pointer.
	IsValueNull PointerIndicator = 1

	// IsValueValid indicates that the field is a valid value.
	//
	// The double pointer points to a pointer that points to a value.
	IsValueValid PointerIndicator = 2
)

// Gateway Events
// https://discord.com/developers/docs/topics/gateway#gateway-events
type Event interface{}

// Gateway Event Names
// https://discord.com/developers/docs/topics/gateway-events
const (
	FlagGatewayEventNameHello                               = "HELLO"
	FlagGatewayEventNameReady                               = "READY"
	FlagGatewayEventNameResumed                             = "RESUMED"
	FlagGatewayEventNameReconnect                           = "RECONNECT"
	FlagGatewayEventNameInvalidSession                      = "INVALID_SESSION"
	FlagGatewayEventNameApplicationCommandPermissionsUpdate = "APPLICATION_COMMAND_PERMISSIONS_UPDATE"
	FlagGatewayEventNameAutoModerationRuleCreate            = "AUTO_MODERATION_RULE_CREATE"
	FlagGatewayEventNameAutoModerationRuleUpdate            = "AUTO_MODERATION_RULE_UPDATE"
	FlagGatewayEventNameAutoModerationRuleDelete            = "AUTO_MODERATION_RULE_DELETE"
	FlagGatewayEventNameAutoModerationActionExecution       = "AUTO_MODERATION_ACTION_EXECUTION"
	FlagGatewayEventNameChannelCreate                       = "CHANNEL_CREATE"
	FlagGatewayEventNameChannelUpdate                       = "CHANNEL_UPDATE"
	FlagGatewayEventNameChannelDelete                       = "CHANNEL_DELETE"
	FlagGatewayEventNameChannelPinsUpdate                   = "CHANNEL_PINS_UPDATE"
	FlagGatewayEventNameThreadCreate                        = "THREAD_CREATE"
	FlagGatewayEventNameThreadUpdate                        = "THREAD_UPDATE"
	FlagGatewayEventNameThreadDelete                        = "THREAD_DELETE"
	FlagGatewayEventNameThreadListSync                      = "THREAD_LIST_SYNC"
	FlagGatewayEventNameThreadMemberUpdate                  = "THREAD_MEMBER_UPDATE"
	FlagGatewayEventNameThreadMembersUpdate                 = "THREAD_MEMBERS_UPDATE"
	FlagGatewayEventNameGuildCreate                         = "GUILD_CREATE"
	FlagGatewayEventNameGuildUpdate                         = "GUILD_UPDATE"
	FlagGatewayEventNameGuildDelete                         = "GUILD_DELETE"
	FlagGatewayEventNameGuildBanAdd                         = "GUILD_BAN_ADD"
	FlagGatewayEventNameGuildBanRemove                      = "GUILD_BAN_REMOVE"
	FlagGatewayEventNameGuildEmojisUpdate                   = "GUILD_EMOJIS_UPDATE"
	FlagGatewayEventNameGuildStickersUpdate                 = "GUILD_STICKERS_UPDATE"
	FlagGatewayEventNameGuildIntegrationsUpdate             = "GUILD_INTEGRATIONS_UPDATE"
	FlagGatewayEventNameGuildMemberAdd                      = "GUILD_MEMBER_ADD"
	FlagGatewayEventNameGuildMemberRemove                   = "GUILD_MEMBER_REMOVE"
	FlagGatewayEventNameGuildMemberUpdate                   = "GUILD_MEMBER_UPDATE"
	FlagGatewayEventNameGuildMembersChunk                   = "GUILD_MEMBERS_CHUNK"
	FlagGatewayEventNameGuildRoleCreate                     = "GUILD_ROLE_CREATE"
	FlagGatewayEventNameGuildRoleUpdate                     = "GUILD_ROLE_UPDATE"
	FlagGatewayEventNameGuildRoleDelete                     = "GUILD_ROLE_DELETE"
	FlagGatewayEventNameGuildScheduledEventCreate           = "GUILD_SCHEDULED_EVENT_CREATE"
	FlagGatewayEventNameGuildScheduledEventUpdate           = "GUILD_SCHEDULED_EVENT_UPDATE"
	FlagGatewayEventNameGuildScheduledEventDelete           = "GUILD_SCHEDULED_EVENT_DELETE"
	FlagGatewayEventNameGuildScheduledEventUserAdd          = "GUILD_SCHEDULED_EVENT_USER_ADD"
	FlagGatewayEventNameGuildScheduledEventUserRemove       = "GUILD_SCHEDULED_EVENT_USER_REMOVE"
	FlagGatewayEventNameIntegrationCreate                   = "INTEGRATION_CREATE"
	FlagGatewayEventNameIntegrationUpdate                   = "INTEGRATION_UPDATE"
	FlagGatewayEventNameIntegrationDelete                   = "INTEGRATION_DELETE"
	FlagGatewayEventNameInteractionCreate                   = "INTERACTION_CREATE"
	FlagGatewayEventNameInviteCreate                        = "INVITE_CREATE"
	FlagGatewayEventNameInviteDelete                        = "INVITE_DELETE"
	FlagGatewayEventNameMessageCreate                       = "MESSAGE_CREATE"
	FlagGatewayEventNameMessageUpdate                       = "MESSAGE_UPDATE"
	FlagGatewayEventNameMessageDelete                       = "MESSAGE_DELETE"
	FlagGatewayEventNameMessageDeleteBulk                   = "MESSAGE_DELETE_BULK"
	FlagGatewayEventNameMessageReactionAdd                  = "MESSAGE_REACTION_ADD"
	FlagGatewayEventNameMessageReactionRemove               = "MESSAGE_REACTION_REMOVE"
	FlagGatewayEventNameMessageReactionRemoveAll            = "MESSAGE_REACTION_REMOVE_ALL"
	FlagGatewayEventNameMessageReactionRemoveEmoji          = "MESSAGE_REACTION_REMOVE_EMOJI"
	FlagGatewayEventNamePresenceUpdate                      = "PRESENCE_UPDATE"
	FlagGatewayEventNameStageInstanceCreate                 = "STAGE_INSTANCE_CREATE"
	FlagGatewayEventNameStageInstanceDelete                 = "STAGE_INSTANCE_DELETE"
	FlagGatewayEventNameStageInstanceUpdate                 = "STAGE_INSTANCE_UPDATE"
	FlagGatewayEventNameTypingStart                         = "TYPING_START"
	FlagGatewayEventNameUserUpdate                          = "USER_UPDATE"
	FlagGatewayEventNameVoiceStateUpdate                    = "VOICE_STATE_UPDATE"
	FlagGatewayEventNameVoiceServerUpdate                   = "VOICE_SERVER_UPDATE"
	FlagGatewayEventNameWebhooksUpdate                      = "WEBHOOKS_UPDATE"
)

// Hello Structure
// https://discord.com/developers/docs/topics/gateway-events#hello-hello-structure
type Hello struct {
	HeartbeatInterval int `json:"heartbeat_interval"`
}

// Ready Event Fields
// https://discord.com/developers/docs/topics/gateway-events#ready-ready-event-fields
type Ready struct {
	Application      *Application `json:"application"`
	User             *User        `json:"user"`
	Shard            *[2]int      `json:"shard,omitempty"`
	SessionID        string       `json:"session_id"`
	ResumeGatewayURL string       `json:"resume_gateway_url"`
	Guilds           []*Guild     `json:"guilds"`
	Version          int          `json:"v"`
}

// Resumed
// https://discord.com/developers/docs/topics/gateway-events#resumed
type Resumed struct{}

// Reconnect
// https://discord.com/developers/docs/topics/gateway-events#reconnect
type Reconnect struct{}

// Invalid Session
// https://discord.com/developers/docs/topics/gateway-events#invalid-session
type InvalidSession struct {
	Data bool `json:"d"`
}

// Application Command Permissions Update
// https://discord.com/developers/docs/topics/gateway-events#application-command-permissions-update
type ApplicationCommandPermissionsUpdate struct {
	*GuildApplicationCommandPermissions
}

// Auto Moderation Rule Create
// https://discord.com/developers/docs/topics/gateway-events#auto-moderation-rule-create
type AutoModerationRuleCreate struct {
	*AutoModerationRule
}

// Auto Moderation Rule Update
// https://discord.com/developers/docs/topics/gateway-events#auto-moderation-rule-update
type AutoModerationRuleUpdate struct {
	*AutoModerationRule
}

// Auto Moderation Rule Delete
// https://discord.com/developers/docs/topics/gateway-events#auto-moderation-rule-delete
type AutoModerationRuleDelete struct {
	*AutoModerationRule
}

// Auto Moderation Action Execution
// https://discord.com/developers/docs/topics/gateway-events#auto-moderation-action-execution
type AutoModerationActionExecution struct {
	Action               AutoModerationAction `json:"action"`
	MessageID            *string              `json:"message_id,omitempty"`
	MatchedKeyword       *string              `json:"matched_keyword"`
	MatchedContent       *string              `json:"matched_content"`
	ChannelID            *string              `json:"channel_id,omitempty"`
	AlertSystemMessageID *string              `json:"alert_system_message_id,omitempty"`
	RuleID               string               `json:"rule_id"`
	GuildID              string               `json:"guild_id"`
	Content              string               `json:"content"`
	UserID               string               `json:"user_id"`
	RuleTriggerType      Flag                 `json:"rule_trigger_type"`
}

// Channel Create
// https://discord.com/developers/docs/topics/gateway-events#channel-create
type ChannelCreate struct {
	*Channel
}

// Channel Update
// https://discord.com/developers/docs/topics/gateway-events#channel-update
type ChannelUpdate struct {
	*Channel
}

// Channel Delete
// https://discord.com/developers/docs/topics/gateway-events#channel-delete
type ChannelDelete struct {
	*Channel
}

// Thread Create
// https://discord.com/developers/docs/topics/gateway-events#thread-create
type ThreadCreate struct {
	*Channel
	NewlyCreated *bool `json:"newly_created,omitempty"`
}

// Thread Update
// https://discord.com/developers/docs/topics/gateway-events#thread-update
type ThreadUpdate struct {
	*Channel
}

// Thread Delete
// https://discord.com/developers/docs/topics/gateway-events#thread-delete
type ThreadDelete struct {
	*Channel
}

// Thread List Sync Event Fields
// https://discord.com/developers/docs/topics/gateway-events#thread-list-sync
type ThreadListSync struct {
	GuildID    string          `json:"guild_id"`
	ChannelIDs []string        `json:"channel_ids,omitempty"`
	Threads    []*Channel      `json:"threads"`
	Members    []*ThreadMember `json:"members"`
}

// Thread Member Update
// https://discord.com/developers/docs/topics/gateway-events#thread-member-update
type ThreadMemberUpdate struct {
	*ThreadMember
	GuildID string `json:"guild_id"`
}

// Thread Members Update
// https://discord.com/developers/docs/topics/gateway-events#thread-members-update
type ThreadMembersUpdate struct {
	ID             string          `json:"id"`
	GuildID        string          `json:"guild_id"`
	AddedMembers   []*ThreadMember `json:"added_members,omitempty"`
	RemovedMembers []string        `json:"removed_member_ids,omitempty"`
	MemberCount    int             `json:"member_count"`
}

// Channel Pins Update
// https://discord.com/developers/docs/topics/gateway-events#channel-pins-update
type ChannelPinsUpdate struct {
	LastPinTimestamp **time.Time `json:"last_pin_timestamp,omitempty"`
	GuildID          string      `json:"guild_id,omitempty"`
	ChannelID        string      `json:"channel_id"`
}

// Guild Create
// https://discord.com/developers/docs/topics/gateway-events#guild-create
type GuildCreate struct {
	*Guild

	// https://discord.com/developers/docs/topics/threads#gateway-events
	Threads []*Channel `json:"threads,omitempty"`
}

// Guild Update
// https://discord.com/developers/docs/topics/gateway-events#guild-update
type GuildUpdate struct {
	*Guild
}

// Guild Delete
// https://discord.com/developers/docs/topics/gateway-events#guild-delete
type GuildDelete struct {
	*Guild
}

// Guild Ban Add
// https://discord.com/developers/docs/topics/gateway-events#guild-ban-add
type GuildBanAdd struct {
	User    *User  `json:"user"`
	GuildID string `json:"guild_id"`
}

// Guild Ban Remove
// https://discord.com/developers/docs/topics/gateway-events#guild-ban-remove
type GuildBanRemove struct {
	User    *User  `json:"user"`
	GuildID string `json:"guild_id"`
}

// Guild Emojis Update
// https://discord.com/developers/docs/topics/gateway-events#guild-emojis-update
type GuildEmojisUpdate struct {
	GuildID string   `json:"guild_id"`
	Emojis  []*Emoji `json:"emojis"`
}

// Guild Stickers Update
// https://discord.com/developers/docs/topics/gateway-events#guild-stickers-update
type GuildStickersUpdate struct {
	GuildID  string     `json:"guild_id"`
	Stickers []*Sticker `json:"stickers"`
}

// Guild Integrations Update
// https://discord.com/developers/docs/topics/gateway-events#guild-integrations-update
type GuildIntegrationsUpdate struct {
	GuildID string `json:"guild_id"`
}

// Guild Member Add
// https://discord.com/developers/docs/topics/gateway-events#guild-member-add
type GuildMemberAdd struct {
	*GuildMember
	GuildID string `json:"guild_id"`
}

// Guild Member Remove
// https://discord.com/developers/docs/topics/gateway-events#guild-member-remove
type GuildMemberRemove struct {
	User    *User  `json:"user"`
	GuildID string `json:"guild_id"`
}

// Guild Member Update
// https://discord.com/developers/docs/topics/gateway-events#guild-member-update
type GuildMemberUpdate struct {
	*GuildMember
}

// Guild Members Chunk
// https://discord.com/developers/docs/topics/gateway-events#guild-members-chunk
type GuildMembersChunk struct {
	Nonce      *string           `json:"nonce,omitempty"`
	GuildID    string            `json:"guild_id"`
	Members    []*GuildMember    `json:"members"`
	Presences  []*PresenceUpdate `json:"presences,omitempty"`
	NotFound   []string          `json:"not_found,omitempty"`
	ChunkIndex int               `json:"chunk_index"`
	ChunkCount int               `json:"chunk_count"`
}

// Guild Role Create
// https://discord.com/developers/docs/topics/gateway-events#guild-role-create
type GuildRoleCreate struct {
	Role    *Role  `json:"role"`
	GuildID string `json:"guild_id"`
}

// Guild Role Update
// https://discord.com/developers/docs/topics/gateway-events#guild-role-update
type GuildRoleUpdate struct {
	Role    *Role  `json:"role"`
	GuildID string `json:"guild_id"`
}

// Guild Role Delete
// https://discord.com/developers/docs/topics/gateway-events#guild-role-delete
type GuildRoleDelete struct {
	GuildID string `json:"guild_id"`
	RoleID  string `json:"role_id"`
}

// Guild Scheduled Event Create
// https://discord.com/developers/docs/topics/gateway-events#guild-scheduled-event-create
type GuildScheduledEventCreate struct {
	*GuildScheduledEvent
}

// Guild Scheduled Event Update
// https://discord.com/developers/docs/topics/gateway-events#guild-scheduled-event-update
type GuildScheduledEventUpdate struct {
	*GuildScheduledEvent
}

// Guild Scheduled Event Delete
// https://discord.com/developers/docs/topics/gateway-events#guild-scheduled-event-delete
type GuildScheduledEventDelete struct {
	*GuildScheduledEvent
}

// Guild Scheduled Event User Add
// https://discord.com/developers/docs/topics/gateway-events#guild-scheduled-event-user-add
type GuildScheduledEventUserAdd struct {
	GuildScheduledEventID string `json:"guild_scheduled_event_id"`
	UserID                string `json:"user_id"`
	GuildID               string `json:"guild_id"`
}

// Guild Scheduled Event User Remove
// https://discord.com/developers/docs/topics/gateway-events#guild-scheduled-event-user-remove
type GuildScheduledEventUserRemove struct {
	GuildScheduledEventID string `json:"guild_scheduled_event_id"`
	UserID                string `json:"user_id"`
	GuildID               string `json:"guild_id"`
}

// Integration Create
// https://discord.com/developers/docs/topics/gateway-events#integration-create
type IntegrationCreate struct {
	*Integration
	GuildID string `json:"guild_id"`
}

// Integration Update
// https://discord.com/developers/docs/topics/gateway-events#integration-update
type IntegrationUpdate struct {
	*Integration
	GuildID string `json:"guild_id"`
}

// Integration Delete
// https://discord.com/developers/docs/topics/gateway-events#integration-delete
type IntegrationDelete struct {
	ApplicationID *string `json:"application_id,omitempty"`
	IntegrationID string  `json:"id"`
	GuildID       string  `json:"guild_id"`
}

// Interaction Create
// https://discord.com/developers/docs/topics/gateway-events#interaction-create
type InteractionCreate struct {
	*Interaction
}

// Invite Create
// https://discord.com/developers/docs/topics/gateway-events#invite-create
type InviteCreate struct {
	CreatedAt         time.Time    `json:"created_at"`
	TargetType        *int         `json:"target_user_type,omitempty"`
	GuildID           *string      `json:"guild_id,omitempty"`
	Inviter           *User        `json:"inviter,omitempty"`
	TargetUser        *User        `json:"target_user,omitempty"`
	TargetApplication *Application `json:"target_application,omitempty"`
	ChannelID         string       `json:"channel_id"`
	Code              string       `json:"code"`
	MaxAge            int          `json:"max_age"`
	MaxUses           int          `json:"max_uses"`
	Uses              int          `json:"uses"`
	Temporary         bool         `json:"temporary"`
}

// Invite Delete
// https://discord.com/developers/docs/topics/gateway-events#invite-delete
type InviteDelete struct {
	ChannelID string  `json:"channel_id"`
	GuildID   *string `json:"guild_id,omitempty"`
	Code      string  `json:"code"`
}

// Message Create
// https://discord.com/developers/docs/topics/gateway-events#message-create
type MessageCreate struct {
	*Message
}

// Message Update
// https://discord.com/developers/docs/topics/gateway-events#message-update
type MessageUpdate struct {
	*Message
}

// Message Delete
// https://discord.com/developers/docs/topics/gateway-events#message-delete
type MessageDelete struct {
	GuildID   *string `json:"guild_id,omitempty"`
	MessageID string  `json:"id"`
	ChannelID string  `json:"channel_id"`
}

// Message Delete Bulk
// https://discord.com/developers/docs/topics/gateway-events#message-delete-bulk
type MessageDeleteBulk struct {
	GuildID    *string  `json:"guild_id,omitempty"`
	ChannelID  string   `json:"channel_id"`
	MessageIDs []string `json:"ids"`
}

// Message Reaction Add
// https://discord.com/developers/docs/topics/gateway-events#message-reaction-add
type MessageReactionAdd struct {
	GuildID   *string      `json:"guild_id,omitempty"`
	Member    *GuildMember `json:"member,omitempty"`
	Emoji     *Emoji       `json:"emoji"`
	UserID    string       `json:"user_id"`
	ChannelID string       `json:"channel_id"`
	MessageID string       `json:"message_id"`
}

// Message Reaction Remove
// https://discord.com/developers/docs/topics/gateway-events#message-reaction-remove
type MessageReactionRemove struct {
	GuildID   *string `json:"guild_id,omitempty"`
	Emoji     *Emoji  `json:"emoji"`
	UserID    string  `json:"user_id"`
	ChannelID string  `json:"channel_id"`
	MessageID string  `json:"message_id"`
}

// Message Reaction Remove All
// https://discord.com/developers/docs/topics/gateway-events#message-reaction-remove-all
type MessageReactionRemoveAll struct {
	GuildID   *string `json:"guild_id,omitempty"`
	ChannelID string  `json:"channel_id"`
	MessageID string  `json:"message_id"`
}

// Message Reaction Remove Emoji
// https://discord.com/developers/docs/topics/gateway-events#message-reaction-remove-emoji
type MessageReactionRemoveEmoji struct {
	GuildID   *string `json:"guild_id,omitempty"`
	Emoji     *Emoji  `json:"emoji"`
	ChannelID string  `json:"channel_id"`
	MessageID string  `json:"message_id"`
}

// Presence Update Event Fields
// https://discord.com/developers/docs/topics/gateway-events#presence-update-presence-update-event-fields
type PresenceUpdate struct {
	User         *User         `json:"user"`
	ClientStatus *ClientStatus `json:"client_status"`
	GuildID      string        `json:"guild_id"`
	Status       string        `json:"status"`
	Activities   []*Activity   `json:"activities"`
}

// Stage Instance Create
// https://discord.com/developers/docs/topics/gateway-events#stage-instance-create
type StageInstanceCreate struct {
	*StageInstance
}

// Stage Instance Update
// https://discord.com/developers/docs/topics/gateway-events#stage-instance-update
type StageInstanceUpdate struct {
	*StageInstance
}

// Stage Instance Delete
// https://discord.com/developers/docs/topics/gateway-events#stage-instance-delete
type StageInstanceDelete struct {
	*StageInstance
}

// Typing Start
// https://discord.com/developers/docs/topics/gateway-events#typing-start
type TypingStart struct {
	Timestamp time.Time    `json:"timestamp"`
	GuildID   *string      `json:"guild_id,omitempty"`
	Member    *GuildMember `json:"member,omitempty"`
	ChannelID string       `json:"channel_id"`
	UserID    string       `json:"user_id"`
}

// User Update
// https://discord.com/developers/docs/topics/gateway-events#user-update
type UserUpdate struct {
	*User
}

// Voice State Update
// https://discord.com/developers/docs/topics/gateway-events#voice-state-update
type VoiceStateUpdate struct {
	*VoiceState
}

// Voice Server Update
// https://discord.com/developers/docs/topics/gateway-events#voice-server-update
type VoiceServerUpdate struct {
	Endpoint *string `json:"endpoint"`
	Token    string  `json:"token"`
	GuildID  string  `json:"guild_id"`
}

// Webhooks Update
// https://discord.com/developers/docs/topics/gateway-events#webhooks-update
type WebhooksUpdate struct {
	GuildID   string `json:"guild_id"`
	ChannelID string `json:"channel_id"`
}

// Gateway Payload Structure
// https://discord.com/developers/docs/topics/gateway-events#payload-structure
type GatewayPayload struct {
	SequenceNumber *int64          `json:"s"`
	EventName      *string         `json:"t"`
	Data           json.RawMessage `json:"d"`
	Op             int             `json:"op"`
}

// Gateway URL Query String Params
// https://discord.com/developers/docs/topics/gateway#connecting-gateway-url-query-string-params
type GatewayURLQueryString struct {
	Compress *string `url:"compress,omitempty"`
	Encoding string  `url:"encoding"`
	V        int     `url:"v"`
}

// Session Start Limit Structure
// https://discord.com/developers/docs/topics/gateway#session-start-limit-object-session-start-limit-structure
type SessionStartLimit struct {
	Total          int `json:"total"`
	Remaining      int `json:"remaining"`
	ResetAfter     int `json:"reset_after"`
	MaxConcurrency int `json:"max_concurrency"`
}

// List of Intents
// https://discord.com/developers/docs/topics/gateway#list-of-intents
const (
	// GUILD_CREATE
	// GUILD_UPDATE
	// GUILD_DELETE
	// GUILD_ROLE_CREATE
	// GUILD_ROLE_UPDATE
	// GUILD_ROLE_DELETE
	// CHANNEL_CREATE
	// CHANNEL_UPDATE
	// CHANNEL_DELETE
	// CHANNEL_PINS_UPDATE
	// THREAD_CREATE
	// THREAD_UPDATE
	// THREAD_DELETE
	// THREAD_LIST_SYNC
	// THREAD_MEMBER_UPDATE
	// THREAD_MEMBERS_UPDATE *
	// STAGE_INSTANCE_CREATE
	// STAGE_INSTANCE_UPDATE
	// STAGE_INSTANCE_DELETE
	FlagIntentGUILDS BitFlag = 1 << 0

	// GUILD_MEMBER_ADD
	// GUILD_MEMBER_UPDATE
	// GUILD_MEMBER_REMOVE
	// THREAD_MEMBERS_UPDATE *
	FlagIntentGUILD_MEMBERS BitFlag = 1 << 1

	// GUILD_BAN_ADD
	// GUILD_BAN_REMOVE
	FlagIntentGUILD_BANS BitFlag = 1 << 2

	// GUILD_EMOJIS_UPDATE
	// GUILD_STICKERS_UPDATE
	FlagIntentGUILD_EMOJIS_AND_STICKERS BitFlag = 1 << 3

	// GUILD_INTEGRATIONS_UPDATE
	// INTEGRATION_CREATE
	// INTEGRATION_UPDATE
	// INTEGRATION_DELETE
	FlagIntentGUILD_INTEGRATIONS BitFlag = 1 << 4

	// WEBHOOKS_UPDATE
	FlagIntentGUILD_WEBHOOKS BitFlag = 1 << 5

	// INVITE_CREATE
	// INVITE_DELETE
	FlagIntentGUILD_INVITES BitFlag = 1 << 6

	// VOICE_STATE_UPDATE
	FlagIntentGUILD_VOICE_STATES BitFlag = 1 << 7

	// PRESENCE_UPDATE
	FlagIntentGUILD_PRESENCES BitFlag = 1 << 8

	// MESSAGE_CREATE
	// MESSAGE_UPDATE
	// MESSAGE_DELETE
	// MESSAGE_DELETE_BULK
	FlagIntentGUILD_MESSAGES BitFlag = 1 << 9

	// MESSAGE_REACTION_ADD
	// MESSAGE_REACTION_REMOVE
	// MESSAGE_REACTION_REMOVE_ALL
	// MESSAGE_REACTION_REMOVE_EMOJI
	FlagIntentGUILD_MESSAGE_REACTIONS BitFlag = 1 << 10

	// TYPING_START
	FlagIntentGUILD_MESSAGE_TYPING  BitFlag = 1 << 11
	FlagIntentDIRECT_MESSAGE_TYPING BitFlag = 1 << 14

	// MESSAGE_CREATE
	// MESSAGE_UPDATE
	// MESSAGE_DELETE
	// CHANNEL_PINS_UPDATE
	FlagIntentDIRECT_MESSAGES BitFlag = 1 << 12

	// MESSAGE_REACTION_ADD
	// MESSAGE_REACTION_REMOVE
	// MESSAGE_REACTION_REMOVE_ALL
	// MESSAGE_REACTION_REMOVE_EMOJI
	FlagIntentDIRECT_MESSAGE_REACTIONS BitFlag = 1 << 13

	FlagIntentMESSAGE_CONTENT BitFlag = 1 << 15

	// GUILD_SCHEDULED_EVENT_CREATE
	// GUILD_SCHEDULED_EVENT_UPDATE
	// GUILD_SCHEDULED_EVENT_DELETE
	// GUILD_SCHEDULED_EVENT_USER_ADD
	// GUILD_SCHEDULED_EVENT_USER_REMOVE
	FlagIntentGUILD_SCHEDULED_EVENTS BitFlag = 1 << 16

	// AUTO_MODERATION_RULE_CREATE
	// AUTO_MODERATION_RULE_UPDATE
	// AUTO_MODERATION_RULE_DELETE
	FlagIntentAUTO_MODERATION_CONFIGURATION BitFlag = 1 << 20

	// AUTO_MODERATION_ACTION_EXECUTION
	FlagIntentAUTO_MODERATION_EXECUTION BitFlag = 1 << 21
)

// Privileged Intents
// https://discord.com/developers/docs/topics/gateway#privileged-intents
var (
	PrivilegedIntents = map[BitFlag]bool{
		FlagIntentGUILD_PRESENCES: true,
		FlagIntentGUILD_MEMBERS:   true,
		FlagIntentMESSAGE_CONTENT: true,
	}
)

// Gateway SendEvent
// https://discord.com/developers/docs/topics/gateway-events#send-events
type SendEvent interface{}

// Gateway SendEvent Names
// https://discord.com/developers/docs/topics/gateway-events#send-events
const (
	FlagGatewaySendEventNameHeartbeat           = "Heartbeat"
	FlagGatewaySendEventNameIdentify            = "Identify"
	FlagGatewaySendEventNameUpdatePresence      = "UpdatePresence"
	FlagGatewaySendEventNameUpdateVoiceState    = "UpdateVoiceState "
	FlagGatewaySendEventNameResume              = "Resume"
	FlagGatewaySendEventNameRequestGuildMembers = "RequestGuildMembers"
)

// Identify Structure
// https://discord.com/developers/docs/topics/gateway-events#identify-identify-structure
type Identify struct {
	Compress       *bool                        `json:"compress,omitempty"`
	LargeThreshold *int                         `json:"large_threshold,omitempty"`
	Shard          *[2]int                      `json:"shard,omitempty"`
	Presence       *GatewayPresenceUpdate       `json:"presence,omitempty"`
	Properties     IdentifyConnectionProperties `json:"properties"`
	Token          string                       `json:"token"`
	Intents        BitFlag                      `json:"intents"`
}

// Identify Connection Properties
// https://discord.com/developers/docs/topics/gateway-events#identify-identify-connection-properties
type IdentifyConnectionProperties struct {
	OS      string `json:"os"`
	Browser string `json:"browser"`
	Device  string `json:"device"`
}

// Resume Structure
// https://discord.com/developers/docs/topics/gateway-events#resume-resume-structure
type Resume struct {
	Token     string `json:"token"`
	SessionID string `json:"session_id"`
	Seq       int64  `json:"seq"`
}

// Heartbeat
// https://discord.com/developers/docs/topics/gateway-events#heartbeat
type Heartbeat struct {
	Data int64 `json:"d"`
}

// Request Guild Members Structure
// https://discord.com/developers/docs/topics/gateway-events#request-guild-members-guild-request-members-structure
type RequestGuildMembers struct {
	Query     *string  `json:"query,omitempty"`
	Presences *bool    `json:"presences,omitempty"`
	Nonce     *string  `json:"nonce,omitempty"`
	GuildID   string   `json:"guild_id"`
	UserIDs   []string `json:"user_ids,omitempty"`
	Limit     int      `json:"limit"`
}

// Gateway Voice State Update Structure
// https://discord.com/developers/docs/topics/gateway-events#update-voice-state-gateway-voice-state-update-structure
type GatewayVoiceStateUpdate struct {
	ChannelID *string `json:"channel_id"`
	GuildID   string  `json:"guild_id"`
	SelfMute  bool    `json:"self_mute"`
	SelfDeaf  bool    `json:"self_deaf"`
}

// Gateway Presence Update Structure
// https://discord.com/developers/docs/topics/gateway-events#update-presence-gateway-presence-update-structure
type GatewayPresenceUpdate struct {
	Since  *int        `json:"since"`
	Status string      `json:"status"`
	Game   []*Activity `json:"game"`
	AFK    bool        `json:"afk"`
}

// Status Types
// https://discord.com/developers/docs/topics/gateway#update-presence-status-types
const (
	FlagStatusTypeOnline       = "online"
	FlagStatusTypeDoNotDisturb = "dnd"
	FlagStatusTypeAFK          = "idle"
	FlagStatusTypeInvisible    = "invisible"
	FlagStatusTypeOffline      = "offline"
)

// Rate Limit Headers
// https://discord.com/developers/docs/topics/rate-limits#header-format-rate-limit-header-examples
const (
	FlagRateLimitHeaderDate       = "Date"
	FlagRateLimitHeaderLimit      = "X-RateLimit-Limit"
	FlagRateLimitHeaderRemaining  = "X-RateLimit-Remaining"
	FlagRateLimitHeaderReset      = "X-RateLimit-Reset"
	FlagRateLimitHeaderResetAfter = "X-RateLimit-Reset-After"
	FlagRateLimitHeaderBucket     = "X-RateLimit-Bucket"
	FlagRateLimitHeaderGlobal     = "X-RateLimit-Global"
	FlagRateLimitHeaderScope      = "X-RateLimit-Scope"
	FlagRateLimitHeaderRetryAfter = "Retry-After"
)

// Rate Limit Header
// https://discord.com/developers/docs/topics/rate-limits#header-format
type RateLimitHeader struct {
	Scope      string  `http:"X-RateLimit-Scope,omitempty"`
	Bucket     string  `http:"X-RateLimit-Bucket,omitempty"`
	Remaining  int     `http:"X-RateLimit-Remaining,omitempty"`
	Reset      float64 `http:"X-RateLimit-Reset,omitempty"`
	ResetAfter float64 `http:"X-RateLimit-Reset-After,omitempty"`
	Limit      int     `http:"X-RateLimit-Limit,omitempty"`
	Global     bool    `http:"X-RateLimit-Global,omitempty"`
}

// Rate Limit Scope Values
// https://discord.com/developers/docs/topics/rate-limits#header-format-rate-limit-header-examples
const (
	RateLimitScopeValueUser   = "user"
	RateLimitScopeValueGlobal = "global"
	RateLimitScopeValueShared = "shared"
)

// Rate Limit Response Structure
// https://discord.com/developers/docs/topics/rate-limits#exceeding-a-rate-limit-rate-limit-response-structure
type RateLimitResponse struct {
	Message    string  `json:"message"`
	RetryAfter float64 `json:"retry_after"`
	Global     bool    `json:"global"`
}

// Global Rate Limits
// https://discord.com/developers/docs/topics/rate-limits#global-rate-limit
const (
	// Global Rate Limit (Requests): 50 requests per second.
	FlagGlobalRateLimitRequest = 50

	// Global Rate Limit (Gateway): 120 commands per minute.
	FlagGlobalRateLimitGateway         = 120
	FlagGlobalRateLimitGatewayInterval = time.Minute

	// Global Rate Limit (Identify Command): Get Gateway Bot `max_concurrency + 1` per 5 Seconds.
	FlagGlobalRateLimitIdentifyInterval = time.Second * 5

	// Global Rate Limit (Identify Command): 1000 per day.
	FlagGlobalRateLimitIdentifyDaily         = 1000
	FlagGlobalRateLimitIdentifyDailyInterval = time.Hour * 24
)

// Invalid Request Limit (CloudFlare Bans)
// https://discord.com/developers/docs/topics/rate-limits#invalid-request-limit-aka-cloudflare-bans
const (
	// 10,000 requests per 10 minutes.
	FlagInvalidRequestRateLimit = 10000
)

var (
	InvalidRateLimitRequests = map[int]string{
		FlagHTTPResponseCodeUNAUTHORIZED:    HTTPResponseCodes[FlagHTTPResponseCodeUNAUTHORIZED],
		FlagHTTPResponseCodeFORBIDDEN:       HTTPResponseCodes[FlagHTTPResponseCodeFORBIDDEN],
		FlagHTTPResponseCodeTOOMANYREQUESTS: HTTPResponseCodes[FlagHTTPResponseCodeTOOMANYREQUESTS],
	}
)

// Version
// https://discord.com/developers/docs/reference#api-versioning
const (
	VersionDiscordAPI = "10"
)

// time.Time Format
// https://discord.com/developers/docs/reference#iso8601-datetime
const (
	TimestampFormatISO8601 = time.RFC3339
)

// Image Formats
// https://discord.com/developers/docs/reference#image-formatting-image-formats
const (
	ImageFormatJPEG   = "JPEG"
	ImageFormatPNG    = "PNG"
	ImageFormatWebP   = "WebP"
	ImageFormatGIF    = "GIF"
	ImageFormatLottie = "Lottie"
)

// CDN Endpoint Exceptions
// https://discord.com/developers/docs/reference#image-formatting-cdn-endpoints
const (
	CDNEndpointAnimatedHashPrefix       = "a_"
	CDNEndpointUserDiscriminatorDivisor = 5
)

// Locales
// https://discord.com/developers/docs/reference#locales
const (
	FlagLocalesDanish              = "da"
	FlagLocalesGerman              = "de"
	FlagLocalesEnglishUK           = "en-GB"
	FlagLocalesEnglishUS           = "en-US"
	FlagLocalesSpanish             = "es-ES"
	FlagLocalesFrench              = "fr"
	FlagLocalesCroatian            = "hr"
	FlagLocalesItalian             = "it"
	FlagLocalesLithuanian          = "lt"
	FlagLocalesHungarian           = "hu"
	FlagLocalesDutch               = "nl"
	FlagLocalesNorwegian           = "no"
	FlagLocalesPolish              = "pl"
	FlagLocalesPortugueseBrazilian = "pt-BR"
	FlagLocalesRomanian            = "ro"
	FlagLocalesFinnish             = "fi"
	FlagLocalesSwedish             = "sv-SE"
	FlagLocalesVietnamese          = "vi"
	FlagLocalesTurkish             = "tr"
	FlagLocalesCzech               = "cs"
	FlagLocalesGreek               = "el"
	FlagLocalesBulgarian           = "bg"
	FlagLocalesRussian             = "ru"
	FlagLocalesUkrainian           = "uk"
	FlagLocalesHindi               = "hi"
	FlagLocalesThai                = "th"
	FlagLocalesChineseChina        = "zh-CN"
	FlagLocalesJapanese            = "ja"
	FlagLocalesChineseTaiwan       = "zh-TW"
	FlagLocalesKorean              = "ko"
)

// Get Global Application Commands
// GET /applications/{application.id}/commands
// https://discord.com/developers/docs/interactions/application-commands#get-global-application-commands
type GetGlobalApplicationCommands struct {
	WithLocalizations *bool `url:"with_localizations,omitempty"`
}

// Create Global Application Command
// POST /applications/{application.id}/commands
// https://discord.com/developers/docs/interactions/application-commands#create-global-application-command
type CreateGlobalApplicationCommand struct {
	DMPermission             **bool                      `json:"dm_permission,omitempty"`
	NameLocalizations        *map[string]string          `json:"name_localizations,omitempty"`
	Description              *string                     `json:"description,omitempty"`
	DescriptionLocalizations *map[string]string          `json:"description_localizations,omitempty"`
	DefaultMemberPermissions **string                    `json:"default_member_permissions,omitempty"`
	Type                     *Flag                       `json:"type,omitempty"`
	Name                     string                      `json:"name,omitempty"`
	Options                  []*ApplicationCommandOption `json:"options,omitempty"`
}

// Get Global Application Command
// GET /applications/{application.id}/commands/{command.id}
// https://discord.com/developers/docs/interactions/application-commands#get-global-application-command
type GetGlobalApplicationCommand struct {
	CommandID string
}

// Edit Global Application Command
// PATCH /applications/{application.id}/commands/{command.id}
// https://discord.com/developers/docs/interactions/application-commands#edit-global-application-command
type EditGlobalApplicationCommand struct {
	DefaultMemberPermissions **string                    `json:"default_member_permissions,omitempty"`
	Name                     *string                     `json:"name,omitempty"`
	NameLocalizations        *map[string]string          `json:"name_localizations,omitempty"`
	Description              *string                     `json:"description,omitempty"`
	DescriptionLocalizations *map[string]string          `json:"description_localizations,omitempty"`
	DMPermission             **bool                      `json:"dm_permission,omitempty"`
	CommandID                string                      `json:"-"`
	Options                  []*ApplicationCommandOption `json:"options,omitempty"`
}

// Delete Global Application Command
// DELETE /applications/{application.id}/commands/{command.id}
// https://discord.com/developers/docs/interactions/application-commands#delete-global-application-command
type DeleteGlobalApplicationCommand struct {
	CommandID string
}

// Bulk Overwrite Global Application Commands
// PUT /applications/{application.id}/commands
// https://discord.com/developers/docs/interactions/application-commands#bulk-overwrite-global-application-commands
type BulkOverwriteGlobalApplicationCommands struct {
	ApplicationCommands []*ApplicationCommand `json:"commands"`
}

// Get Guild Application Commands
// GET /applications/{application.id}/guilds/{guild.id}/commands
// https://discord.com/developers/docs/interactions/application-commands#get-guild-application-commands
type GetGuildApplicationCommands struct {
	WithLocalizations *bool  `url:"with_localizations,omitempty"`
	GuildID           string `url:"-"`
}

// Create Guild Application Command
// POST /applications/{application.id}/guilds/{guild.id}/commands
// https://discord.com/developers/docs/interactions/application-commands#create-guild-application-command
type CreateGuildApplicationCommand struct {
	DefaultMemberPermissions **string                    `json:"default_member_permissions,omitempty"`
	Type                     *Flag                       `json:"type,omitempty"`
	NameLocalizations        *map[string]string          `json:"name_localizations,omitempty"`
	Description              *string                     `json:"description,omitempty"`
	DescriptionLocalizations *map[string]string          `json:"description_localizations,omitempty"`
	GuildID                  string                      `json:"-"`
	Name                     string                      `json:"name"`
	Options                  []*ApplicationCommandOption `json:"options,omitempty"`
}

// Get Guild Application Command
// GET /applications/{application.id}/guilds/{guild.id}/commands/{command.id}
// https://discord.com/developers/docs/interactions/application-commands#get-guild-application-command
type GetGuildApplicationCommand struct {
	GuildID   string
	CommandID string
}

// Edit Guild Application Command
// PATCH /applications/{application.id}/guilds/{guild.id}/commands/{command.id}
// https://discord.com/developers/docs/interactions/application-commands#edit-guild-application-command
type EditGuildApplicationCommand struct {
	DefaultMemberPermissions **string                    `json:"default_member_permissions,omitempty"`
	Name                     *string                     `json:"name,omitempty"`
	NameLocalizations        *map[string]string          `json:"name_localizations,omitempty"`
	Description              *string                     `json:"description,omitempty"`
	DescriptionLocalizations *map[string]string          `json:"description_localizations,omitempty"`
	GuildID                  string                      `json:"-"`
	CommandID                string                      `json:"-"`
	Options                  []*ApplicationCommandOption `json:"options,omitempty"`
}

// Delete Guild Application Command
// DELETE /applications/{application.id}/guilds/{guild.id}/commands/{command.id}
// https://discord.com/developers/docs/interactions/application-commands#delete-guild-application-command
type DeleteGuildApplicationCommand struct {
	GuildID   string
	CommandID string
}

// Bulk Overwrite Guild Application Commands
// PUT /applications/{application.id}/guilds/{guild.id}/commands
// https://discord.com/developers/docs/interactions/application-commands#bulk-overwrite-guild-application-commands
type BulkOverwriteGuildApplicationCommands struct {
	GuildID             string                `json:"-"`
	ApplicationCommands []*ApplicationCommand `json:"commands"`
}

// Get Guild Application Command Permissions
// GET /applications/{application.id}/guilds/{guild.id}/commands/permissions
// https://discord.com/developers/docs/interactions/application-commands#get-guild-application-command-permissions
type GetGuildApplicationCommandPermissions struct {
	GuildID string
}

// Get Application Command Permissions
// GET /applications/{application.id}/guilds/{guild.id}/commands/{command.id}/permissions
// https://discord.com/developers/docs/interactions/application-commands#get-application-command-permissions
type GetApplicationCommandPermissions struct {
	GuildID   string
	CommandID string
}

// Edit Application Command Permissions
// PUT /applications/{application.id}/guilds/{guild.id}/commands/{command.id}/permissions
// https://discord.com/developers/docs/interactions/application-commands#edit-application-command-permissions
type EditApplicationCommandPermissions struct {
	GuildID     string                           `json:"-"`
	CommandID   string                           `json:"-"`
	Permissions []*ApplicationCommandPermissions `json:"permissions"`
}

// Batch Edit Application Command Permissions
// PUT /applications/{application.id}/guilds/{guild.id}/commands/permissions
// https://discord.com/developers/docs/interactions/application-commands#batch-edit-application-command-permissions
type BatchEditApplicationCommandPermissions struct {
	GuildID string
}

// Create Interaction Response
// POST /interactions/{interaction.id}/{interaction.token}/callback
// https://discord.com/developers/docs/interactions/receiving-and-responding#create-interaction-response
type CreateInteractionResponse struct {
	*InteractionResponse
	InteractionID    string
	InteractionToken string
}

// Get Original Interaction Response
// GET /webhooks/{application.id}/{interaction.token}/messages/@original
// https://discord.com/developers/docs/interactions/receiving-and-responding#get-original-interaction-response
type GetOriginalInteractionResponse struct {
	ThreadID         *string `url:"thread_id,omitempty"`
	InteractionToken string  `url:"-"`
}

// Edit Original Interaction Response
// PATCH /webhooks/{application.id}/{interaction.token}/messages/@original
// https://discord.com/developers/docs/interactions/receiving-and-responding#edit-original-interaction-response
type EditOriginalInteractionResponse struct {
	AllowedMentions  **AllowedMentions `json:"allowed_mentions,omitempty" url:"-"`
	ThreadID         *string           `json:"-" url:"thread_id,omitempty"`
	Content          **string          `json:"content,omitempty" url:"-"`
	Embeds           *[]*Embed         `json:"embeds,omitempty" url:"-"`
	Components       *[]Component      `json:"components,omitempty" url:"-"`
	Attachments      *[]*Attachment    `json:"attachments,omitempty" url:"-"`
	ApplicationID    string            `json:"-" url:"-"`
	InteractionToken string            `json:"-" url:"-"`
	Files            []*File           `json:"-" url:"-" dasgo:"files"`
}

// Delete Original Interaction Response
// DELETE /webhooks/{application.id}/{interaction.token}/messages/@original
// https://discord.com/developers/docs/interactions/receiving-and-responding#delete-original-interaction-response
type DeleteOriginalInteractionResponse struct {
	InteractionToken string
}

// Create Followup Message
// POST /webhooks/{application.id}/{interaction.token}
// https://discord.com/developers/docs/interactions/receiving-and-responding#create-followup-message
type CreateFollowupMessage struct {
	AllowedMentions  *AllowedMentions `json:"allowed_mentions,omitempty" url:"-"`
	Flags            *BitFlag         `json:"flags,omitempty" url:"-"`
	ThreadID         *string          `json:"-" url:"thread_id,omitempty"`
	Content          *string          `json:"content,omitempty" url:"-"`
	Username         *string          `json:"username,omitempty" url:"-"`
	AvatarURL        *string          `json:"avatar_url,omitempty" url:"-"`
	TTS              *bool            `json:"tts,omitempty" url:"-"`
	ThreadName       *string          `json:"thread_name,omitempty" url:"-"`
	ApplicationID    string           `json:"-" url:"-"`
	InteractionToken string           `json:"-" url:"-"`
	Files            []*File          `json:"-" url:"-" dasgo:"files"`
	Attachments      []*Attachment    `json:"attachments,omitempty" url:"-"`
	Components       []Component      `json:"components,omitempty" url:"-"`
	Embeds           []*Embed         `json:"embeds,omitempty" url:"-"`
}

// Get Followup Message
// GET /webhooks/{application.id}/{interaction.token}/messages/{message.id}
// https://discord.com/developers/docs/interactions/receiving-and-responding#get-followup-message
type GetFollowupMessage struct {
	ThreadID         *string `url:"thread_id,omitempty"`
	InteractionToken string  `url:"-"`
	MessageID        string  `url:"-"`
}

// Edit Followup Message
// PATCH /webhooks/{application.id}/{interaction.token}/messages/{message.id}
// https://discord.com/developers/docs/interactions/receiving-and-responding#edit-followup-message
type EditFollowupMessage struct {
	Components       *[]Component      `json:"components,omitempty" url:"-"`
	AllowedMentions  **AllowedMentions `json:"allowed_mentions,omitempty" url:"-"`
	ThreadID         *string           `json:"-" url:"thread_id,omitempty"`
	Content          **string          `json:"content,omitempty" url:"-"`
	Embeds           *[]*Embed         `json:"embeds,omitempty" url:"-"`
	Attachments      *[]*Attachment    `json:"attachments,omitempty" url:"-"`
	InteractionToken string            `json:"-" url:"-"`
	ApplicationID    string            `json:"-" url:"-"`
	MessageID        string            `json:"-" url:"-"`
	Files            []*File           `json:"-" url:"-" dasgo:"files"`
}

// Delete Followup Message
// DELETE /webhooks/{application.id}/{interaction.token}/messages/{message.id}
// https://discord.com/developers/docs/interactions/receiving-and-responding#delete-followup-message
type DeleteFollowupMessage struct {
	InteractionToken string
	MessageID        string
}

// Get Guild Audit Log
// GET /guilds/{guild.id}/audit-logs
// https://discord.com/developers/docs/resources/audit-log#get-guild-audit-log
type GetGuildAuditLog struct {
	GuildID    string `url:"-"`
	UserID     string `url:"user_id,omitempty"`
	Before     string `url:"before,omitempty"`
	Limit      int    `url:"limit,omitempty"`
	ActionType Flag   `url:"action_type,omitempty"`
}

// List Auto Moderation Rules for Guild
// GET /guilds/{guild.id}/auto-moderation/rules
// https://discord.com/developers/docs/resources/auto-moderation#list-auto-moderation-rules-for-guild
type ListAutoModerationRulesForGuild struct {
	GuildID string
}

// Get Auto Moderation Rule
// GET /guilds/{guild.id}/auto-moderation/rules/{auto_moderation_rule.id}
// https://discord.com/developers/docs/resources/auto-moderation#get-auto-moderation-rule
type GetAutoModerationRule struct {
	GuildID              string
	AutoModerationRuleID string
}

// Create Auto Moderation Rule
// POST /guilds/{guild.id}/auto-moderation/rules
// https://discord.com/developers/docs/resources/auto-moderation#create-auto-moderation-rule
type CreateAutoModerationRule struct {
	Enabled         *bool                   `json:"enabled,omitempty"`
	TriggerMetadata *TriggerMetadata        `json:"trigger_metadata,omitempty"`
	Name            string                  `json:"name"`
	GuildID         string                  `json:"-"`
	ExemptChannels  []string                `json:"exempt_channels,omitempty"`
	Actions         []*AutoModerationAction `json:"actions"`
	ExemptRoles     []string                `json:"exempt_roles,omitempty"`
	TriggerType     Flag                    `json:"trigger_type"`
	EventType       Flag                    `json:"event_type"`
}

// Modify Auto Moderation Rule
// PATCH /guilds/{guild.id}/auto-moderation/rules/{auto_moderation_rule.id}
// https://discord.com/developers/docs/resources/auto-moderation#modify-auto-moderation-rule
type ModifyAutoModerationRule struct {
	GuildID              string                  `json:"-"`
	AutoModerationRuleID string                  `json:"-"`
	Name                 *string                 `json:"name,omitempty"`
	EventType            *Flag                   `json:"event_type,omitempty"`
	TriggerType          *Flag                   `json:"trigger_type,omitempty"`
	TriggerMetadata      *TriggerMetadata        `json:"trigger_metadata,omitempty"`
	Actions              []*AutoModerationAction `json:"actions,omitempty"`
	Enabled              *bool                   `json:"enabled,omitempty"`
	ExemptRoles          []string                `json:"exempt_roles,omitempty"`
	ExemptChannels       []string                `json:"exempt_channels,omitempty"`
}

// Delete Auto Moderation Rule
// DELETE /guilds/{guild.id}/auto-moderation/rules/{auto_moderation_rule.id}
// https://discord.com/developers/docs/resources/auto-moderation#delete-auto-moderation-rule
type DeleteAutoModerationRule struct {
	GuildID              string
	AutoModerationRuleID string
}

// Get Channel
// GET /channels/{channel.id}
// https://discord.com/developers/docs/resources/channel#get-channel
type GetChannel struct {
	ChannelID string
}

// Modify Channel
// PATCH /channels/{channel.id}
// https://discord.com/developers/docs/resources/channel#modify-channel
type ModifyChannel struct {
	ChannelID string
}

// Modify Channel Group DM
// PATCH /channels/{channel.id}
// https://discord.com/developers/docs/resources/channel#modify-channel-json-params-group-dm
type ModifyChannelGroupDM struct {
	Name      *string `json:"name,omitempty"`
	Icon      *string `json:"icon,omitempty"`
	ChannelID string  `json:"-"`
}

// Modify Channel Guild
// PATCH /channels/{channel.id}
// https://discord.com/developers/docs/resources/channel#modify-channel-json-params-guild-channel
type ModifyChannelGuild struct {
	DefaultSortOrder              **int                   `json:"default_sort_order,omitempty"`
	Name                          *string                 `json:"name,omitempty"`
	Type                          *Flag                   `json:"type,omitempty"`
	Position                      **int                   `json:"position,omitempty"`
	Topic                         **string                `json:"topic,omitempty"`
	NSFW                          **bool                  `json:"nsfw,omitempty"`
	RateLimitPerUser              **int                   `json:"rate_limit_per_user,omitempty"`
	Bitrate                       **int                   `json:"bitrate,omitempty"`
	UserLimit                     **int                   `json:"user_limit,omitempty"`
	PermissionOverwrites          *[]*PermissionOverwrite `json:"permission_overwrites,omitempty"`
	ParentID                      **string                `json:"parent_id,omitempty"`
	RTCRegion                     **string                `json:"rtc_region,omitempty"`
	VideoQualityMode              **Flag                  `json:"video_quality_mode,omitempty"`
	DefaultAutoArchiveDuration    **int                   `json:"default_auto_archive_duration,omitempty"`
	Flags                         *BitFlag                `json:"flags,omitempty"`
	DefaultThreadRateLimitPerUser *int                    `json:"default_thread_rate_limit_per_user,omitempty"`
	DefaultReactionEmoji          **DefaultReaction       `json:"default_reaction_emoji,omitempty"`
	ChannelID                     string                  `json:"-"`
	AvailableTags                 []*ForumTag             `json:"available_tags,omitempty"`
}

// Modify Channel
// PATCH /channels/{channel.id}
// https://discord.com/developers/docs/resources/channel#modify-channel-json-params-thread
type ModifyChannelThread struct {
	ChannelID           string   `json:"-"`
	Name                *string  `json:"name,omitempty"`
	Archived            *bool    `json:"archived,omitempty"`
	AutoArchiveDuration *int     `json:"auto_archive_duration,omitempty"`
	Locked              *bool    `json:"locked,omitempty"`
	Invitable           *bool    `json:"invitable,omitempty"`
	RateLimitPerUser    **int    `json:"rate_limit_per_user,omitempty"`
	Flags               *BitFlag `json:"flags,omitempty"`
	AppliedTags         []string `json:"applied_tags,omitempty"`
}

// Delete/Close Channel
// DELETE /channels/{channel.id}
// https://discord.com/developers/docs/resources/channel#deleteclose-channel
type DeleteCloseChannel struct {
	ChannelID string
}

// Get Channel Messages
// GET /channels/{channel.id}/messages
// https://discord.com/developers/docs/resources/channel#get-channel-messages
type GetChannelMessages struct {
	Around    *string `url:"around,omitempty"`
	Before    *string `url:"before,omitempty"`
	After     *string `url:"after,omitempty"`
	Limit     *Flag   `url:"limit,omitempty"`
	ChannelID string  `url:"-"`
}

// Get Channel Message
// GET /channels/{channel.id}/messages/{message.id}
// https://discord.com/developers/docs/resources/channel#get-channel-message
type GetChannelMessage struct {
	ChannelID string
	MessageID string
}

// Create Message
// POST /channels/{channel.id}/messages
// https://discord.com/developers/docs/resources/channel#create-message
type CreateMessage struct {
	MessageReference *MessageReference `json:"message_reference,omitempty"`
	Content          *string           `json:"content,omitempty"`
	Nonce            *Nonce            `json:"nonce,omitempty"`
	TTS              *bool             `json:"tts,omitempty"`
	AllowedMentions  *AllowedMentions  `json:"allowed_mentions,omitempty"`
	Flags            *BitFlag          `json:"flags,omitempty"`
	ChannelID        string            `json:"-"`
	Embeds           []*Embed          `json:"embeds,omitempty"`
	Components       []Component       `json:"components,omitempty"`
	StickerIDS       []*string         `json:"sticker_ids,omitempty"`
	Files            []*File           `json:"-" dasgo:"files,omitempty"`
	Attachments      []*Attachment     `json:"attachments,omitempty"`
}

// Crosspost Message
// POST /channels/{channel.id}/messages/{message.id}/crosspost
// https://discord.com/developers/docs/resources/channel#crosspost-message
type CrosspostMessage struct {
	ChannelID string
	MessageID string
}

// Create Reaction
// PUT /channels/{channel.id}/messages/{message.id}/reactions/{emoji}/@me
// https://discord.com/developers/docs/resources/channel#create-reaction
type CreateReaction struct {
	ChannelID string
	MessageID string
	Emoji     string
}

// Delete Own Reaction
// DELETE /channels/{channel.id}/messages/{message.id}/reactions/{emoji}/@me
// https://discord.com/developers/docs/resources/channel#delete-own-reaction
type DeleteOwnReaction struct {
	ChannelID string
	MessageID string
	Emoji     string
}

// Delete User Reaction
// DELETE /channels/{channel.id}/messages/{message.id}/reactions/{emoji}/{user.id}
// https://discord.com/developers/docs/resources/channel#delete-user-reaction
type DeleteUserReaction struct {
	ChannelID string
	MessageID string
	Emoji     string
	UserID    string
}

// Get Reactions
// GET /channels/{channel.id}/messages/{message.id}/reactions/{emoji}
// https://discord.com/developers/docs/resources/channel#get-reactions
type GetReactions struct {
	After     *string `url:"after,omitempty"`
	Limit     *int    `url:"limit,omitempty"`
	ChannelID string  `url:"-"`
	MessageID string  `url:"-"`
	Emoji     string  `url:"-"`
}

// Delete All Reactions
// DELETE /channels/{channel.id}/messages/{message.id}/reactions
// https://discord.com/developers/docs/resources/channel#delete-all-reactions
type DeleteAllReactions struct {
	ChannelID string
	MessageID string
}

// Delete All Reactions for Emoji
// DELETE /channels/{channel.id}/messages/{message.id}/reactions/{emoji}
// https://discord.com/developers/docs/resources/channel#delete-all-reactions-for-emoji
type DeleteAllReactionsforEmoji struct {
	ChannelID string
	MessageID string
	Emoji     string
}

// Edit Message
// PATCH /channels/{channel.id}/messages/{message.id}
// https://discord.com/developers/docs/resources/channel#edit-message
type EditMessage struct {
	Components      *[]Component      `json:"components,omitempty"`
	Content         **string          `json:"content,omitempty"`
	Embeds          *[]*Embed         `json:"embeds,omitempty"`
	Flags           **BitFlag         `json:"flags,omitempty"`
	AllowedMentions **AllowedMentions `json:"allowed_mentions,omitempty"`
	Attachments     *[]*Attachment    `json:"attachments,omitempty"`
	MessageID       string            `json:"-"`
	ChannelID       string            `json:"-"`
	Files           []*File           `json:"-" dasgo:"files"`
}

// Delete Message
// DELETE /channels/{channel.id}/messages/{message.id}
// https://discord.com/developers/docs/resources/channel#delete-message
type DeleteMessage struct {
	ChannelID string
	MessageID string
}

// Bulk Delete Messages
// POST /channels/{channel.id}/messages/bulk-delete
// https://discord.com/developers/docs/resources/channel#bulk-delete-messages
type BulkDeleteMessages struct {
	ChannelID string    `json:"-"`
	Messages  []*string `json:"messages"`
}

// Edit Channel Permissions
// PUT /channels/{channel.id}/permissions/{overwrite.id}
// https://discord.com/developers/docs/resources/channel#edit-channel-permissions
type EditChannelPermissions struct {
	Allow       **string `json:"allow,omitempty"`
	Deny        **string `json:"deny,omitempty"`
	ChannelID   string   `json:"-"`
	OverwriteID string   `json:"-"`
	Type        Flag     `json:"type"`
}

// Get Channel Invites
// GET /channels/{channel.id}/invites
// https://discord.com/developers/docs/resources/channel#get-channel-invites
type GetChannelInvites struct {
	ChannelID string
}

// Create Channel Invite
// POST /channels/{channel.id}/invites
// https://discord.com/developers/docs/resources/channel#create-channel-invite
type CreateChannelInvite struct {
	MaxAge              *int   `json:"max_age"`
	MaxUses             *int   `json:"max_uses"`
	ChannelID           string `json:"-"`
	TargetUserID        string `json:"target_user_id"`
	TargetApplicationID string `json:"target_application_id"`
	Temporary           bool   `json:"temporary"`
	Unique              bool   `json:"unique"`
	TargetType          Flag   `json:"target_type"`
}

// Delete Channel Permission
// DELETE /channels/{channel.id}/permissions/{overwrite.id}
// https://discord.com/developers/docs/resources/channel#delete-channel-permission
type DeleteChannelPermission struct {
	ChannelID   string
	OverwriteID string
}

// Follow Announcement Channel
// POST /channels/{channel.id}/followers
// https://discord.com/developers/docs/resources/channel#follow-announcement-channel
type FollowAnnouncementChannel struct {
	ChannelID        string `json:"-"`
	WebhookChannelID string `json:"webhook_channel_id"`
}

// Trigger Typing Indicator
// POST /channels/{channel.id}/typing
// https://discord.com/developers/docs/resources/channel#trigger-typing-indicator
type TriggerTypingIndicator struct {
	ChannelID string
}

// Get Pinned Messages
// GET /channels/{channel.id}/pins
// https://discord.com/developers/docs/resources/channel#get-pinned-messages
type GetPinnedMessages struct {
	ChannelID string
}

// Pin Message
// PUT /channels/{channel.id}/pins/{message.id}
// https://discord.com/developers/docs/resources/channel#pin-message
type PinMessage struct {
	ChannelID string
	MessageID string
}

// Unpin Message
// DELETE /channels/{channel.id}/pins/{message.id}
// https://discord.com/developers/docs/resources/channel#unpin-message
type UnpinMessage struct {
	ChannelID string
	MessageID string
}

// Group DM Add Recipient
// PUT /channels/{channel.id}/recipients/{user.id}
// https://discord.com/developers/docs/resources/channel#group-dm-add-recipient
type GroupDMAddRecipient struct {
	Nickname    *string `json:"nick"`
	ChannelID   string  `json:"-"`
	UserID      string  `json:"-"`
	AccessToken string  `json:"access_token"`
}

// Group DM Remove Recipient
// DELETE /channels/{channel.id}/recipients/{user.id}
// https://discord.com/developers/docs/resources/channel#group-dm-remove-recipient
type GroupDMRemoveRecipient struct {
	ChannelID string
	UserID    string
}

// Start Thread from Message
// POST /channels/{channel.id}/messages/{message.id}/threads
// https://discord.com/developers/docs/resources/channel#start-thread-from-message
type StartThreadfromMessage struct {
	AutoArchiveDuration *int   `json:"auto_archive_duration,omitempty"`
	RateLimitPerUser    **int  `json:"rate_limit_per_user,omitempty"`
	ChannelID           string `json:"-"`
	MessageID           string `json:"-"`
	Name                string `json:"name"`
}

// Start Thread without Message
// POST /channels/{channel.id}/threads
// https://discord.com/developers/docs/resources/channel#start-thread-without-message
type StartThreadwithoutMessage struct {
	AutoArchiveDuration *int   `json:"auto_archive_duration,omitempty"`
	Type                *Flag  `json:"type,omitempty"`
	Invitable           *bool  `json:"invitable,omitempty"`
	RateLimitPerUser    **int  `json:"rate_limit_per_user,omitempty"`
	ChannelID           string `json:"-"`
	Name                string `json:"name"`
}

// Start Thread in Forum Channel
// POST /channels/{channel.id}/threads
// https://discord.com/developers/docs/resources/channel#start-thread-in-forum-channel
type StartThreadinForumChannel struct {
	ChannelID           string                    `json:"-"`
	Name                string                    `json:"name"`
	AutoArchiveDuration *int                      `json:"auto_archive_duration,omitempty"`
	RateLimitPerUser    **int                     `json:"rate_limit_per_user,omitempty"`
	Message             *ForumThreadMessageParams `json:"message"`
	AppliedTags         []string                  `json:"applied_tags,omitempty"`
}

// Forum Thread Message Params Object
// https://discord.com/developers/docs/resources/channel#start-thread-in-forum-channel-forum-thread-message-params-object
type ForumThreadMessageParams struct {
	Content         *string          `json:"content,omitempty"`
	Flags           *BitFlag         `json:"flags,omitempty"`
	AllowedMentions *AllowedMentions `json:"allowed_mentions,omitempty"`
	Components      []Component      `json:"components,omitempty"`
	StickerIDS      []*string        `json:"sticker_ids,omitempty"`
	Attachments     []*Attachment    `json:"attachments,omitempty"`
	Files           []*File          `json:"-" dasgo:"files"`
	Embeds          []*Embed         `json:"embeds,omitempty"`
}

// Join Thread
// PUT /channels/{channel.id}/thread-members/@me
// https://discord.com/developers/docs/resources/channel#join-thread
type JoinThread struct {
	ChannelID string
}

// Add Thread Member
// PUT /channels/{channel.id}/thread-members/{user.id}
// https://discord.com/developers/docs/resources/channel#add-thread-member
type AddThreadMember struct {
	ChannelID string
	UserID    string
}

// Leave Thread
// DELETE /channels/{channel.id}/thread-members/@me
// https://discord.com/developers/docs/resources/channel#leave-thread
type LeaveThread struct {
	ChannelID string
}

// Remove Thread Member
// DELETE /channels/{channel.id}/thread-members/{user.id}
// https://discord.com/developers/docs/resources/channel#remove-thread-member
type RemoveThreadMember struct {
	ChannelID string
	UserID    string
}

// Get Thread Member
// GET /channels/{channel.id}/thread-members/{user.id}
// https://discord.com/developers/docs/resources/channel#get-thread-member
type GetThreadMember struct {
	ChannelID string
	UserID    string
}

// List Thread Members
// GET /channels/{channel.id}/thread-members
// https://discord.com/developers/docs/resources/channel#list-thread-members
type ListThreadMembers struct {
	ChannelID string
}

// List Public Archived Threads
// GET /channels/{channel.id}/threads/archived/public
// https://discord.com/developers/docs/resources/channel#list-public-archived-threads
type ListPublicArchivedThreads struct {
	Before    *time.Time `url:"before,omitempty"`
	Limit     *int       `url:"limit,omitempty"`
	ChannelID string     `url:"-"`
}

// List Private Archived Threads
// GET /channels/{channel.id}/threads/archived/private
// https://discord.com/developers/docs/resources/channel#list-private-archived-threads
type ListPrivateArchivedThreads struct {
	Before    *time.Time `url:"before,omitempty"`
	Limit     *int       `url:"limit,omitempty"`
	ChannelID string     `url:"-"`
}

// List Joined Private Archived Threads
// GET /channels/{channel.id}/users/@me/threads/archived/private
// https://discord.com/developers/docs/resources/channel#list-joined-private-archived-threads
type ListJoinedPrivateArchivedThreads struct {
	Before    *time.Time `url:"before,omitempty"`
	Limit     *int       `url:"limit,omitempty"`
	ChannelID string     `url:"-"`
}

// List Guild Emojis
// GET /guilds/{guild.id}/emojis
// https://discord.com/developers/docs/resources/emoji#list-guild-emojis
type ListGuildEmojis struct {
	GuildID string
}

// Get Guild Emoji
// GET /guilds/{guild.id}/emojis/{emoji.id}
// https://discord.com/developers/docs/resources/emoji#get-guild-emoji
type GetGuildEmoji struct {
	GuildID string
	EmojiID string
}

// Create Guild Emoji
// POST /guilds/{guild.id}/emojis
// https://discord.com/developers/docs/resources/emoji#create-guild-emoji
type CreateGuildEmoji struct {
	GuildID string    `json:"-"`
	Name    string    `json:"name"`
	Image   string    `json:"image"`
	Roles   []*string `json:"roles"`
}

// Modify Guild Emoji
// PATCH /guilds/{guild.id}/emojis/{emoji.id}
// https://discord.com/developers/docs/resources/emoji#modify-guild-emoji
type ModifyGuildEmoji struct {
	Name    *string    `json:"name,omitempty"`
	Roles   *[]*string `json:"roles,omitempty"`
	GuildID string     `json:"-"`
	EmojiID string     `json:"-"`
}

// Delete Guild Emoji
// DELETE /guilds/{guild.id}/emojis/{emoji.id}
// https://discord.com/developers/docs/resources/emoji#delete-guild-emoji
type DeleteGuildEmoji struct {
	GuildID string
	EmojiID string
}

// Create Guild
// POST /guilds
// https://discord.com/developers/docs/resources/guild#create-guild
type CreateGuild struct {
	Region                      **string   `json:"region,omitempty"`
	Icon                        *string    `json:"icon,omitempty"`
	VerificationLevel           *Flag      `json:"verification_level,omitempty"`
	DefaultMessageNotifications *Flag      `json:"default_message_notifications,omitempty"`
	ExplicitContentFilter       *Flag      `json:"explicit_content_filter,omitempty"`
	AfkChannelID                *string    `json:"afk_channel_id,omitempty"`
	AfkTimeout                  *int       `json:"afk_timeout,omitempty"`
	SystemChannelID             *string    `json:"system_channel_id,omitempty"`
	SystemChannelFlags          *BitFlag   `json:"system_channel_flags,omitempty"`
	Name                        string     `json:"name"`
	Roles                       []*Role    `json:"roles,omitempty"`
	Channels                    []*Channel `json:"channels,omitempty"`
}

// Get Guild
// GET /guilds/{guild.id}
// https://discord.com/developers/docs/resources/guild#get-guild
type GetGuild struct {
	WithCounts *bool  `url:"with_counts,omitempty"`
	GuildID    string `url:"-"`
}

// Get Guild Preview
// GET /guilds/{guild.id}/preview
// https://discord.com/developers/docs/resources/guild#get-guild-preview
type GetGuildPreview struct {
	GuildID string
}

// Modify Guild
// PATCH /guilds/{guild.id}
// https://discord.com/developers/docs/resources/guild#modify-guild
type ModifyGuild struct {
	PremiumProgressBarEnabled   *bool     `json:"premium_progress_bar_enabled,omitempty"`
	Name                        *string   `json:"name,omitempty"`
	VerificationLevel           **Flag    `json:"verification_level,omitempty"`
	DefaultMessageNotifications **Flag    `json:"default_message_notifications,omitempty"`
	ExplicitContentFilter       **Flag    `json:"explicit_content_filter,omitempty"`
	AFKChannelID                **string  `json:"afk_channel_id,omitempty"`
	Description                 **string  `json:"description,omitempty"`
	Icon                        **string  `json:"icon,omitempty"`
	PreferredLocale             **string  `json:"preferred_locale,omitempty"`
	Splash                      **string  `json:"splash,omitempty"`
	DiscoverySplash             **string  `json:"discovery_splash,omitempty"`
	Banner                      **string  `json:"banner,omitempty"`
	SystemChannelID             **string  `json:"system_channel_id,omitempty"`
	SystemChannelFlags          *BitFlag  `json:"system_channel_flags,omitempty"`
	RulesChannelID              **string  `json:"rules_channel_id,omitempty"`
	PublicUpdatesChannelID      **string  `json:"public_updates_channel_id,omitempty"`
	OwnerID                     string    `json:"owner_id,omitempty"`
	GuildID                     string    `json:"-"`
	Features                    []*string `json:"features,omitempty"`
	AfkTimeout                  int       `json:"afk_timeout,omitempty"`
}

// Delete Guild
// DELETE /guilds/{guild.id}
// https://discord.com/developers/docs/resources/guild#delete-guild
type DeleteGuild struct {
	GuildID string
}

// Get Guild Channels
// GET /guilds/{guild.id}/channels
// https://discord.com/developers/docs/resources/guild#get-guild-channels
type GetGuildChannels struct {
	GuildID string
}

// Create Guild Channel
// POST /guilds/{guild.id}/channels
// https://discord.com/developers/docs/resources/guild#create-guild-channel
type CreateGuildChannel struct {
	DefaultSortOrder           **int                   `json:"default_sort_order,omitempty"`
	AvailableTags              *[]*ForumTag            `json:"available_tags,omitempty"`
	Type                       **Flag                  `json:"type,omitempty"`
	Topic                      **string                `json:"topic,omitempty"`
	Bitrate                    **int                   `json:"bitrate,omitempty"`
	UserLimit                  **int                   `json:"user_limit,omitempty"`
	RateLimitPerUser           **int                   `json:"rate_limit_per_user,omitempty"`
	Position                   **int                   `json:"position,omitempty"`
	PermissionOverwrites       *[]*PermissionOverwrite `json:"permission_overwrites,omitempty"`
	ParentID                   **string                `json:"parent_id,omitempty"`
	NSFW                       **bool                  `json:"nsfw,omitempty"`
	RTCRegion                  **string                `json:"rtc_region,omitempty"`
	VideoQualityMode           **Flag                  `json:"video_quality_mode,omitempty"`
	DefaultAutoArchiveDuration **int                   `json:"default_auto_archive_duration,omitempty"`
	DefaultReactionEmoji       **DefaultReaction       `json:"default_reaction_emoji,omitempty"`
	Name                       string                  `json:"name"`
	GuildID                    string                  `json:"-"`
}

// Modify Guild Channel Positions
// PATCH /guilds/{guild.id}/channels
// https://discord.com/developers/docs/resources/guild#modify-guild-channel-positions
type ModifyGuildChannelPositions struct {
	GuildID    string                                  `json:"-"`
	Parameters []*ModifyGuildChannelPositionParameters `json:"parameters"`
}

// Modify Guild Channel Position Parameters
// https://discord.com/developers/docs/resources/guild#modify-guild-channel-positions-json-params
type ModifyGuildChannelPositionParameters struct {
	Position        *int    `json:"position"`
	LockPermissions *bool   `json:"lock_permissions"`
	ParentID        *string `json:"parent_id"`
	ID              string  `json:"id"`
}

// List Active Guild Threads
// GET /guilds/{guild.id}/threads/active
// https://discord.com/developers/docs/resources/guild#list-active-guild-threads
type ListActiveGuildThreads struct {
	GuildID string `json:"-"`
}

// Get Guild Member
// GET /guilds/{guild.id}/members/{user.id}
// https://discord.com/developers/docs/resources/guild#get-guild-member
type GetGuildMember struct {
	GuildID string
	UserID  string
}

// List Guild Members
// GET /guilds/{guild.id}/members
// https://discord.com/developers/docs/resources/guild#list-guild-members
type ListGuildMembers struct {
	Limit   *int    `url:"limit,omitempty"`
	After   *string `url:"after,omitempty"`
	GuildID string  `url:"-"`
}

// Search Guild Members
// GET /guilds/{guild.id}/members/search
// https://discord.com/developers/docs/resources/guild#search-guild-members
type SearchGuildMembers struct {
	Limit   *int   `url:"limit,omitempty"`
	GuildID string `url:"-"`
	Query   string `url:"query"`
}

// Add Guild Member
// PUT /guilds/{guild.id}/members/{user.id}
// https://discord.com/developers/docs/resources/guild#add-guild-member
type AddGuildMember struct {
	Deaf        *bool    `json:"deaf,omitempty"`
	Nick        *string  `json:"nick,omitempty"`
	Mute        *bool    `json:"mute,omitempty"`
	UserID      string   `json:"-"`
	AccessToken string   `json:"access_token"`
	GuildID     string   `json:"-"`
	Roles       []string `json:"roles,omitempty"`
}

// Modify Guild Member
// PATCH /guilds/{guild.id}/members/{user.id}
// https://discord.com/developers/docs/resources/guild#modify-guild-member
type ModifyGuildMember struct {
	ChannelID                  **string    `json:"channel_id,omitempty"`
	CommunicationDisabledUntil **time.Time `json:"communication_disabled_until,omitempty"`
	Nick                       **string    `json:"nick,omitempty"`
	Roles                      *[]string   `json:"roles,omitempty"`
	Mute                       **bool      `json:"mute,omitempty"`
	Deaf                       **bool      `json:"deaf,omitempty"`
	GuildID                    string      `json:"-"`
	UserID                     string      `json:"-"`
}

// Modify Current Member
// PATCH /guilds/{guild.id}/members/@me
// https://discord.com/developers/docs/resources/guild#modify-current-member
type ModifyCurrentMember struct {
	Nick    **string `json:"nick,omitempty"`
	GuildID string   `json:"-"`
}

// Add Guild Member Role
// PUT /guilds/{guild.id}/members/{user.id}/roles/{role.id}
// https://discord.com/developers/docs/resources/guild#add-guild-member-role
type AddGuildMemberRole struct {
	GuildID string
	UserID  string
	RoleID  string
}

// Remove Guild Member Role
// DELETE /guilds/{guild.id}/members/{user.id}/roles/{role.id}
// https://discord.com/developers/docs/resources/guild#remove-guild-member-role
type RemoveGuildMemberRole struct {
	GuildID string
	UserID  string
	RoleID  string
}

// Remove Guild Member
// DELETE /guilds/{guild.id}/members/{user.id}
// https://discord.com/developers/docs/resources/guild#remove-guild-member
type RemoveGuildMember struct {
	GuildID string
	UserID  string
}

// Get Guild Bans
// GET /guilds/{guild.id}/bans
// https://discord.com/developers/docs/resources/guild#get-guild-bans
type GetGuildBans struct {
	Limit   *int    `url:"limit,omitempty"`
	Before  *string `url:"before,omitempty"`
	After   *string `url:"after,omitempty"`
	GuildID string  `url:"-"`
}

// Get Guild Ban
// GET /guilds/{guild.id}/bans/{user.id}
// https://discord.com/developers/docs/resources/guild#get-guild-ban
type GetGuildBan struct {
	GuildID string
	UserID  string
}

// Create Guild Ban
// PUT /guilds/{guild.id}/bans/{user.id}
// https://discord.com/developers/docs/resources/guild#create-guild-ban
type CreateGuildBan struct {
	DeleteMessageSeconds *int   `json:"delete_message_seconds,omitempty"`
	GuildID              string `json:"-"`
	UserID               string `json:"-"`
}

// Remove Guild Ban
// DELETE /guilds/{guild.id}/bans/{user.id}
// https://discord.com/developers/docs/resources/guild#remove-guild-ban
type RemoveGuildBan struct {
	GuildID string
	UserID  string
}

// Get Guild Roles
// GET /guilds/{guild.id}/roles
// https://discord.com/developers/docs/resources/guild#get-guild-roles
type GetGuildRoles struct {
	GuildID string
}

// Create Guild Role
// POST /guilds/{guild.id}/roles
// https://discord.com/developers/docs/resources/guild#create-guild-role
type CreateGuildRole struct {
	UnicodeEmoji **string `json:"unicode_emoji,omitempty"`
	Name         *string  `json:"name,omitempty"`
	Permissions  *string  `json:"permissions,omitempty"`
	Color        *int     `json:"color,omitempty"`
	Hoist        *bool    `json:"hoist,omitempty"`
	Icon         **string `json:"icon,omitempty"`
	Mentionable  *bool    `json:"mentionable,omitempty"`
	GuildID      string   `json:"-"`
}

// Modify Guild Role Positions
// PATCH /guilds/{guild.id}/roles
// https://discord.com/developers/docs/resources/guild#modify-guild-role-positions
type ModifyGuildRolePositions struct {
	GuildID    string                               `json:"-"`
	Parameters []*ModifyGuildRolePositionParameters `json:"parameters"`
}

// Modify Guild Role Position Parameters
// https://discord.com/developers/docs/resources/guild#create-guild-role-json-params
type ModifyGuildRolePositionParameters struct {
	Position **int  `json:"position,omitempty"`
	ID       string `json:"id"`
}

// Modify Guild Role
// PATCH /guilds/{guild.id}/roles/{role.id}
// https://discord.com/developers/docs/resources/guild#modify-guild-role
type ModifyGuildRole struct {
	Icon         **string `json:"icon,omitempty"`
	UnicodeEmoji **string `json:"unicode_emoji,omitempty"`
	Name         **string `json:"name,omitempty"`
	Permissions  **string `json:"permissions,omitempty"`
	Color        **int    `json:"color,omitempty"`
	Hoist        **bool   `json:"hoist,omitempty"`
	Mentionable  **bool   `json:"mentionable,omitempty"`
	GuildID      string   `json:"-"`
	RoleID       string   `json:"-"`
}

// Modify Guild MFA Level
// POST /guilds/{guild.id}/mfa
// https://discord.com/developers/docs/resources/guild#modify-guild-mfa-level
type ModifyGuildMFALevel struct {
	GuildID string `json:"-"`
	Level   Flag   `json:"level"`
}

// Delete Guild Role
// DELETE /guilds/{guild.id}/roles/{role.id}
// https://discord.com/developers/docs/resources/guild#delete-guild-role
type DeleteGuildRole struct {
	GuildID string
	RoleID  string
}

// Get Guild Prune Count
// GET /guilds/{guild.id}/prune
// https://discord.com/developers/docs/resources/guild#get-guild-prune-count
type GetGuildPruneCount struct {
	GuildID      string   `url:"-"`
	IncludeRoles []string `url:"include_roles,omitempty"`
	Days         int      `url:"days,omitempty"`
}

// Begin Guild Prune
// POST /guilds/{guild.id}/prune
// https://discord.com/developers/docs/resources/guild#begin-guild-prune
type BeginGuildPrune struct {
	GuildID           string   `json:"-"`
	IncludeRoles      []string `json:"include_roles"`
	Days              int      `json:"days"`
	ComputePruneCount bool     `json:"compute_prune_count"`
}

// Get Guild Voice Regions
// GET /guilds/{guild.id}/regions
// https://discord.com/developers/docs/resources/guild#get-guild-voice-regions
type GetGuildVoiceRegions struct {
	GuildID string
}

// Get Guild Invites
// GET /guilds/{guild.id}/invites
// https://discord.com/developers/docs/resources/guild#get-guild-invites
type GetGuildInvites struct {
	GuildID string
}

// Get Guild Integrations
// GET /guilds/{guild.id}/integrations
// https://discord.com/developers/docs/resources/guild#get-guild-integrations
type GetGuildIntegrations struct {
	GuildID string
}

// Delete Guild Integration
// DELETE /guilds/{guild.id}/integrations/{integration.id}
// https://discord.com/developers/docs/resources/guild#delete-guild-integration
type DeleteGuildIntegration struct {
	GuildID       string
	IntegrationID string
}

// Get Guild Widget Settings
// GET /guilds/{guild.id}/widget
// https://discord.com/developers/docs/resources/guild#get-guild-widget-settings
type GetGuildWidgetSettings struct {
	GuildID string
}

// Modify Guild Widget
// PATCH /guilds/{guild.id}/widget
// https://discord.com/developers/docs/resources/guild#modify-guild-widget
type ModifyGuildWidget struct {
	GuildID string
}

// Get Guild Widget
// GET /guilds/{guild.id}/widget.json
// https://discord.com/developers/docs/resources/guild#get-guild-widget
type GetGuildWidget struct {
	GuildID string
}

// Get Guild Vanity URL
// GET /guilds/{guild.id}/vanity-url
// https://discord.com/developers/docs/resources/guild#get-guild-vanity-url
type GetGuildVanityURL struct {
	Code    *string `json:"code"`
	GuildID string  `json:"-"`
	Uses    int     `json:"uses,omitempty"`
}

// Get Guild Widget Image
// GET /guilds/{guild.id}/widget.png
// https://discord.com/developers/docs/resources/guild#get-guild-widget-image
type GetGuildWidgetImage struct {
	Style   *string `url:"style,omitempty"`
	GuildID string  `url:"-"`
}

// Widget Style Options
// https://discord.com/developers/docs/resources/guild#get-guild-widget-image-widget-style-options
const (
	FlagWidgetStyleOptionShield  = "shield"
	FlagWidgetStyleOptionBanner1 = "banner1"
	FlagWidgetStyleOptionBanner2 = "banner2"
	FlagWidgetStyleOptionBanner3 = "banner3"
	FlagWidgetStyleOptionBanner4 = "banner4"
)

// Get Guild Welcome Screen
// GET /guilds/{guild.id}/welcome-screen
// https://discord.com/developers/docs/resources/guild#get-guild-welcome-screen
type GetGuildWelcomeScreen struct {
	GuildID string
}

// Modify Guild Welcome Screen
// PATCH /guilds/{guild.id}/welcome-screen
// https://discord.com/developers/docs/resources/guild#modify-guild-welcome-screen
type ModifyGuildWelcomeScreen struct {
	Enabled         **bool                   `json:"enabled,omitempty"`
	WelcomeChannels *[]*WelcomeScreenChannel `json:"welcome_channels,omitempty"`
	Description     **string                 `json:"description,omitempty"`
	GuildID         string                   `json:"-"`
}

// Modify Current User Voice State
// PATCH /guilds/{guild.id}/voice-states/@me
// https://discord.com/developers/docs/resources/guild#modify-current-user-voice-state
type ModifyCurrentUserVoiceState struct {
	ChannelID               *string     `json:"channel_id,omitempty"`
	Suppress                *bool       `json:"suppress,omitempty"`
	RequestToSpeakTimestamp **time.Time `json:"request_to_speak_timestamp,omitempty"`
	GuildID                 string      `json:"-"`
}

// Modify User Voice State
// PATCH /guilds/{guild.id}/voice-states/{user.id}
// https://discord.com/developers/docs/resources/guild#modify-user-voice-state
type ModifyUserVoiceState struct {
	Suppress  *bool  `json:"suppress,omitempty"`
	GuildID   string `json:"-"`
	UserID    string `json:"-"`
	ChannelID string `json:"channel_id"`
}

// List Scheduled Events for Guild
// GET /guilds/{guild.id}/scheduled-events
// https://discord.com/developers/docs/resources/guild-scheduled-event#list-scheduled-events-for-guild
type ListScheduledEventsforGuild struct {
	WithUserCount *bool  `url:"with_user_count,omitempty"`
	GuildID       string `url:"-"`
}

// Create Guild Scheduled Event
// POST /guilds/{guild.id}/scheduled-events
// https://discord.com/developers/docs/resources/guild-scheduled-event#create-guild-scheduled-event
type CreateGuildScheduledEvent struct {
	ScheduledStartTime time.Time                          `json:"scheduled_start_time"`
	ScheduledEndTime   *time.Time                         `json:"scheduled_end_time,omitempty"`
	ChannelID          *string                            `json:"channel_id,omitempty"`
	EntityMetadata     *GuildScheduledEventEntityMetadata `json:"entity_metadata,omitempty"`
	Image              *string                            `json:"image,omitempty"`
	Description        *string                            `json:"description,omitempty"`
	EntityType         *Flag                              `json:"entity_type,omitempty"`
	GuildID            string                             `json:"-"`
	Name               string                             `json:"name"`
	PrivacyLevel       Flag                               `json:"privacy_level"`
}

// Get Guild Scheduled Event
// GET /guilds/{guild.id}/scheduled-events/{guild_scheduled_event.id}
// https://discord.com/developers/docs/resources/guild-scheduled-event#get-guild-scheduled-event
type GetGuildScheduledEvent struct {
	WithUserCount         *bool  `url:"with_user_count,omitempty"`
	GuildID               string `url:"-"`
	GuildScheduledEventID string `url:"-"`
}

// Modify Guild Scheduled Event
// PATCH /guilds/{guild.id}/scheduled-events/{guild_scheduled_event.id}
// https://discord.com/developers/docs/resources/guild-scheduled-event#modify-guild-scheduled-event
type ModifyGuildScheduledEvent struct {
	ScheduledStartTime    *time.Time                          `json:"scheduled_start_time,omitempty"`
	ScheduledEndTime      *time.Time                          `json:"scheduled_end_time,omitempty"`
	ChannelID             *string                             `json:"channel_id,omitempty"`
	EntityMetadata        **GuildScheduledEventEntityMetadata `json:"entity_metadata,omitempty"`
	Name                  *string                             `json:"name,omitempty"`
	PrivacyLevel          *Flag                               `json:"privacy_level,omitempty"`
	Description           **string                            `json:"description,omitempty"`
	EntityType            *Flag                               `json:"entity_type,omitempty"`
	Status                *Flag                               `json:"status,omitempty"`
	Image                 *string                             `json:"image,omitempty"`
	GuildID               string                              `json:"-"`
	GuildScheduledEventID string                              `json:"-"`
}

// Delete Guild Scheduled Event
// DELETE /guilds/{guild.id}/scheduled-events/{guild_scheduled_event.id}
// https://discord.com/developers/docs/resources/guild-scheduled-event#delete-guild-scheduled-event
type DeleteGuildScheduledEvent struct {
	GuildID               string
	GuildScheduledEventID string
}

// Get Guild Scheduled Event Users
// GET /guilds/{guild.id}/scheduled-events/{guild_scheduled_event.id}/users
// https://discord.com/developers/docs/resources/guild-scheduled-event#get-guild-scheduled-event-users
type GetGuildScheduledEventUsers struct {
	Limit                 *int    `url:"limit,omitempty"`
	WithMember            *bool   `url:"with_member,omitempty"`
	Before                *string `url:"before,omitempty"`
	After                 *string `url:"after,omitempty"`
	GuildID               string  `url:"-"`
	GuildScheduledEventID string  `url:"-"`
}

// Get Guild Template
// GET /guilds/templates/{template.code}
// https://discord.com/developers/docs/resources/guild-template#get-guild-template
type GetGuildTemplate struct {
	TemplateCode string
}

// Create Guild from Guild Template
// POST /guilds/templates/{template.code}
// https://discord.com/developers/docs/resources/guild-template#create-guild-from-guild-template
type CreateGuildfromGuildTemplate struct {
	Icon         *string `json:"icon,omitempty"`
	TemplateCode string  `json:"-"`
	Name         string  `json:"name"`
}

// Get Guild Templates
// GET /guilds/{guild.id}/templates
// https://discord.com/developers/docs/resources/guild-template#get-guild-templates
type GetGuildTemplates struct {
	GuildID string
}

// Create Guild Template
// POST /guilds/{guild.id}/templates
// https://discord.com/developers/docs/resources/guild-template#create-guild-template
type CreateGuildTemplate struct {
	Description **string `json:"description,omitempty"`
	GuildID     string   `json:"-"`
	Name        string   `json:"name"`
}

// Sync Guild Template
// PUT /guilds/{guild.id}/templates/{template.code}
// https://discord.com/developers/docs/resources/guild-template#sync-guild-template
type SyncGuildTemplate struct {
	GuildID      string
	TemplateCode string
}

// Modify Guild Template
// PATCH /guilds/{guild.id}/templates/{template.code}
// https://discord.com/developers/docs/resources/guild-template#modify-guild-template
type ModifyGuildTemplate struct {
	Name         *string  `json:"name,omitempty"`
	Description  **string `json:"description,omitempty"`
	GuildID      string
	TemplateCode string `json:"-"`
}

// Delete Guild Template
// DELETE /guilds/{guild.id}/templates/{template.code}
// https://discord.com/developers/docs/resources/guild-template#delete-guild-template
type DeleteGuildTemplate struct {
	GuildID      string
	TemplateCode string
}

// Get Invite
// GET /invites/{invite.code}
// https://discord.com/developers/docs/resources/invite#get-invite
type GetInvite struct {
	WithCounts            *bool   `url:"with_counts,omitempty"`
	WithExpiration        *bool   `url:"with_expiration,omitempty"`
	GuildScheduledEventID *string `url:"guild_scheduled_event_id,omitempty"`
	InviteCode            string  `url:"-"`
}

// Delete Invite
// DELETE /invites/{invite.code}
// https://discord.com/developers/docs/resources/invite#delete-invite
type DeleteInvite struct {
	InviteCode string
}

// Create Stage Instance
// POST /stage-instances
// https://discord.com/developers/docs/resources/stage-instance#create-stage-instance
type CreateStageInstance struct {
	PrivacyLevel          *Flag  `json:"privacy_level,omitempty"`
	SendStartNotification *bool  `json:"send_start_notification,omitempty"`
	ChannelID             string `json:"channel_id"`
	Topic                 string `json:"topic"`
}

// Get Stage Instance
// GET /stage-instances/{channel.id}
// https://discord.com/developers/docs/resources/stage-instance#get-stage-instance
type GetStageInstance struct {
	ChannelID string
}

// Modify Stage Instance
// PATCH /stage-instances/{channel.id}
// https://discord.com/developers/docs/resources/stage-instance#modify-stage-instance
type ModifyStageInstance struct {
	Topic        *string `json:"topic,omitempty"`
	PrivacyLevel *Flag   `json:"privacy_level,omitempty"`
	ChannelID    string  `json:"-"`
}

// Delete Stage Instance
// DELETE /stage-instances/{channel.id}
// https://discord.com/developers/docs/resources/stage-instance#delete-stage-instance
type DeleteStageInstance struct {
	ChannelID string
}

// Get Sticker
// GET /stickers/{sticker.id}
// https://discord.com/developers/docs/resources/sticker#get-sticker
type GetSticker struct {
	StickerID string
}

// List Nitro Sticker Packs
// GET /sticker-packs
// https://discord.com/developers/docs/resources/sticker#list-nitro-sticker-packs
type ListNitroStickerPacks struct{}

// List Guild Stickers
// GET /guilds/{guild.id}/stickers
// https://discord.com/developers/docs/resources/sticker#list-guild-stickers
type ListGuildStickers struct {
	GuildID string
}

// Get Guild Sticker
// GET /guilds/{guild.id}/stickers/{sticker.id}
// https://discord.com/developers/docs/resources/sticker#get-guild-sticker
type GetGuildSticker struct {
	GuildID   string
	StickerID string
}

// Create Guild Sticker
// POST /guilds/{guild.id}/stickers
// https://discord.com/developers/docs/resources/sticker#create-guild-sticker
type CreateGuildSticker struct {
	GuildID     string  `json:"-"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Tags        *string `json:"tags"`
	File        File    `json:"-" dasgo:"file"`
}

// Modify Guild Sticker
// PATCH /guilds/{guild.id}/stickers/{sticker.id}
// https://discord.com/developers/docs/resources/sticker#modify-guild-sticker
type ModifyGuildSticker struct {
	Name        *string  `json:"name,omitempty"`
	Description **string `json:"description,omitempty"`
	Tags        *string  `json:"tags,omitempty"`
	GuildID     string   `json:"-"`
	StickerID   string   `json:"-"`
}

// Delete Guild Sticker
// DELETE /guilds/{guild.id}/stickers/{sticker.id}
// https://discord.com/developers/docs/resources/sticker#delete-guild-sticker
type DeleteGuildSticker struct {
	GuildID   string
	StickerID string
}

// Get Current User
// GET/users/@me
// https://discord.com/developers/docs/resources/user#get-current-user
type GetCurrentUser struct{}

// Get User
// GET/users/{user.id}
// https://discord.com/developers/docs/resources/user#get-user
type GetUser struct {
	UserID string
}

// Modify Current User
// PATCH /users/@me
// https://discord.com/developers/docs/resources/user#modify-current-user
type ModifyCurrentUser struct {
	Username *string `json:"username,omitempty"`
	Avatar   *string `json:"avatar,omitempty"`
}

// Get Current User Guilds
// GET /users/@me/guilds
// https://discord.com/developers/docs/resources/user#get-current-user-guilds
type GetCurrentUserGuilds struct {
	Before *string `json:"before,omitempty"`
	After  *string `json:"after,omitempty"`
	Limit  *int    `json:"limit,omitempty"`
}

// Get Current User Guild Member
// GET /users/@me/guilds/{guild.id}/member
// https://discord.com/developers/docs/resources/user#get-current-user-guild-member
type GetCurrentUserGuildMember struct {
	GuildID string
}

// Leave Guild
// DELETE /users/@me/guilds/{guild.id}
// https://discord.com/developers/docs/resources/user#leave-guild
type LeaveGuild struct {
	GuildID string
}

// Create DM
// POST /users/@me/channels
// https://discord.com/developers/docs/resources/user#create-dm
type CreateDM struct {
	RecipientID string `json:"recipient_id"`
}

// Create Group DM
// POST /users/@me/channels
// https://discord.com/developers/docs/resources/user#create-group-dm
type CreateGroupDM struct {
	Nicks        map[string]string `json:"nicks"`
	AccessTokens []*string         `json:"access_tokens"`
}

// Get User Connections
// GET /users/@me/connections
// https://discord.com/developers/docs/resources/user#get-user-connections
type GetUserConnections struct{}

// List Voice Regions
// GET /voice/regions
// https://discord.com/developers/docs/resources/voice#list-voice-regions
type ListVoiceRegions struct{}

// Create Webhook
// POST /channels/{channel.id}/webhooks
// https://discord.com/developers/docs/resources/webhook#create-webhook
type CreateWebhook struct {
	Avatar    **string `json:"avatar,omitempty"`
	ChannelID string   `json:"-"`
	Name      string   `json:"name"`
}

// Get Channel Webhooks
// GET /channels/{channel.id}/webhooks
// https://discord.com/developers/docs/resources/webhook#get-channel-webhooks
type GetChannelWebhooks struct {
	ChannelID string
}

// Get Guild Webhooks
// GET /guilds/{guild.id}/webhooks
// https://discord.com/developers/docs/resources/webhook#get-guild-webhooks
type GetGuildWebhooks struct {
	GuildID string
}

// Get Webhook
// GET /webhooks/{webhook.id}
// https://discord.com/developers/docs/resources/webhook#get-webhook
type GetWebhook struct {
	WebhookID string
}

// Get Webhook with Token
// GET /webhooks/{webhook.id}/{webhook.token}
// https://discord.com/developers/docs/resources/webhook#get-webhook-with-token
type GetWebhookwithToken struct {
	WebhookID    string
	WebhookToken string
}

// Modify Webhook
// PATCH /webhooks/{webhook.id}
// https://discord.com/developers/docs/resources/webhook#modify-webhook
type ModifyWebhook struct {
	Name      *string  `json:"name,omitempty"`
	Avatar    **string `json:"avatar,omitempty"`
	ChannelID *string  `json:"channel_id,omitempty"`
	WebhookID string   `json:"-"`
}

// Modify Webhook with Token
// PATCH /webhooks/{webhook.id}/{webhook.token}
// https://discord.com/developers/docs/resources/webhook#modify-webhook-with-token
type ModifyWebhookwithToken struct {
	Name         *string  `json:"name,omitempty"`
	Avatar       **string `json:"avatar,omitempty"`
	WebhookID    string
	WebhookToken string
}

// Delete Webhook
// DELETE /webhooks/{webhook.id}
// https://discord.com/developers/docs/resources/webhook#delete-webhook
type DeleteWebhook struct {
	WebhookID string
}

// Delete Webhook with Token
// DELETE /webhooks/{webhook.id}/{webhook.token}
// https://discord.com/developers/docs/resources/webhook#delete-webhook-with-token
type DeleteWebhookwithToken struct {
	WebhookID    string
	WebhookToken string
}

// Execute Webhook
// POST /webhooks/{webhook.id}/{webhook.token}
// https://discord.com/developers/docs/resources/webhook#execute-webhook
type ExecuteWebhook struct {
	AllowedMentions *AllowedMentions `json:"allowed_mentions,omitempty" url:"-"`
	Flags           *BitFlag         `json:"flags,omitempty" url:"-"`
	Wait            *bool            `json:"-" url:"wait,omitempty"`
	ThreadID        *string          `json:"-" url:"thread_id,omitempty"`
	Content         *string          `json:"content,omitempty" url:"-"`
	Username        *string          `json:"username,omitempty" url:"-"`
	AvatarURL       *string          `json:"avatar_url,omitempty" url:"-"`
	TTS             *bool            `json:"tts,omitempty" url:"-"`
	ThreadName      *string          `json:"thread_name,omitempty" url:"-"`
	WebhookToken    string           `json:"-" url:"-"`
	WebhookID       string           `json:"-" url:"-"`
	Embeds          []*Embed         `json:"embeds,omitempty" url:"-"`
	Components      []Component      `json:"components,omitempty" url:"-"`
	Files           []*File          `json:"-" url:"-" dasgo:"files"`
	Attachments     []*Attachment    `json:"attachments,omitempty" url:"-"`
}

// Execute Slack-Compatible Webhook
// POST /webhooks/{webhook.id}/{webhook.token}/slack
// https://discord.com/developers/docs/resources/webhook#execute-slackcompatible-webhook
type ExecuteSlackCompatibleWebhook struct {
	ThreadID     *string `url:"thread_id,omitempty"`
	Wait         *bool   `url:"wait,omitempty"`
	WebhookID    string  `url:"-"`
	WebhookToken string  `url:"-"`
}

// Execute GitHub-Compatible Webhook
// POST /webhooks/{webhook.id}/{webhook.token}/github
// https://discord.com/developers/docs/resources/webhook#execute-githubcompatible-webhook
type ExecuteGitHubCompatibleWebhook struct {
	ThreadID     *string `url:"thread_id,omitempty"`
	Wait         *bool   `url:"wait,omitempty"`
	WebhookID    string  `url:"-"`
	WebhookToken string  `url:"-"`
}

// Get Webhook Message
// GET /webhooks/{webhook.id}/{webhook.token}/messages/{message.id}
// https://discord.com/developers/docs/resources/webhook#get-webhook-message
type GetWebhookMessage struct {
	ThreadID     *string `url:"thread_id,omitempty"`
	WebhookID    string  `url:"-"`
	WebhookToken string  `url:"-"`
	MessageID    string  `url:"-"`
}

// Edit Webhook Message
// PATCH /webhooks/{webhook.id}/{webhook.token}/messages/{message.id}
// https://discord.com/developers/docs/resources/webhook#edit-webhook-message
type EditWebhookMessage struct {
	Components      *[]Component      `json:"components,omitempty" url:"-"`
	AllowedMentions **AllowedMentions `json:"allowed_mentions,omitempty" url:"-"`
	ThreadID        *string           `url:"thread_id,omitempty"`
	Content         **string          `json:"content,omitempty" url:"-"`
	Embeds          *[]*Embed         `json:"embeds,omitempty" url:"-"`
	Attachments     *[]*Attachment    `json:"attachments,omitempty" url:"-"`
	WebhookToken    string            `json:"-" url:"-"`
	WebhookID       string            `json:"-" url:"-"`
	MessageID       string            `json:"-" url:"-"`
	Files           []*File           `json:"-" url:"-" dasgo:"files"`
}

// Delete Webhook Message
// DELETE /webhooks/{webhook.id}/{webhook.token}/messages/{message.id}
// https://discord.com/developers/docs/resources/webhook#delete-webhook-message
type DeleteWebhookMessage struct {
	ThreadID     *string `url:"thread_id,omitempty"`
	WebhookID    string  `url:"-"`
	WebhookToken string  `url:"-"`
	MessageID    string  `url:"-"`
}

// Get Current Bot Application Information
// GET /oauth2/applications/@me
// https://discord.com/developers/docs/topics/oauth2#get-current-bot-application-information
type GetCurrentBotApplicationInformation struct{}

// Get Current Authorization Information
// GET /oauth2/@me
// https://discord.com/developers/docs/topics/oauth2#get-current-authorization-information
type GetCurrentAuthorizationInformation struct{}

// Get Gateway
// GET /gateway
// https://discord.com/developers/docs/topics/gateway#get-gateway
type GetGateway struct{}

// Get Gateway Bot
// GET /gateway/bot
// https://discord.com/developers/docs/topics/gateway#get-gateway-bot
type GetGatewayBot struct{}

// Authorization URL
// GET /oauth2/authorize
// https://discord.com/developers/docs/topics/oauth2#authorization-code-grant-authorization-url-example
type AuthorizationURL struct {
	ResponseType string `url:"response_type,omitempty"`
	ClientID     string `url:"client_id,omitempty"`
	Scope        string `url:"scope,omitempty"`
	State        string `url:"state,omitempty"`
	RedirectURI  string `url:"redirect_uri,omitempty"`
	Prompt       string `url:"prompt,omitempty"`
}

// Access Token Exchange
// POST /oauth2/token
// https://discord.com/developers/docs/topics/oauth2#authorization-code-grant-access-token-exchange-example
type AccessTokenExchange struct {
	ClientID     string `url:"client_id,omitempty"`
	ClientSecret string `url:"client_secret,omitempty"`
	GrantType    string `url:"grant_type,omitempty"`
	Code         string `url:"code,omitempty"`
	RedirectURI  string `url:"redirect_uri,omitempty"`
}

// Refresh Token Exchange
// POST /oauth2/token
// https://discord.com/developers/docs/topics/oauth2#authorization-code-grant-refresh-token-exchange-example
type RefreshTokenExchange struct {
	ClientID     string `url:"client_id,omitempty"`
	ClientSecret string `url:"client_secret,omitempty"`
	GrantType    string `url:"grant_type,omitempty"`
	RefreshToken string `url:"refresh_token,omitempty"`
}

// Client Credentials Token Request
// POST /oauth2/token
// https://discord.com/developers/docs/topics/oauth2#client-credentials-grant-client-credentials-token-request-example
type ClientCredentialsTokenRequest struct {
	GrantType string `url:"grant_type,omitempty"`
	Scope     string `url:"scope,omitempty"`
}

// Bot Auth Parameters
// GET /oauth2/authorize
// https://discord.com/developers/docs/topics/oauth2#bot-authorization-flow-bot-auth-parameters
type BotAuth struct {
	ClientID           string  `url:"client_id"`
	Scope              string  `url:"scope"`
	GuildID            string  `url:"guild_id"`
	Permissions        BitFlag `url:"permissions"`
	DisableGuildSelect bool    `url:"disable_guild_select"`
}

// Application Command Structure
// https://discord.com/developers/docs/interactions/application-commands#application-command-object-application-command-structure
type ApplicationCommand struct {
	Type                     *Flag                       `json:"type,omitempty"`
	GuildID                  *string                     `json:"guild_id,omitempty"`
	DMPermission             *bool                       `json:"dm_permission,omitempty"`
	NameLocalizations        *map[string]string          `json:"name_localizations,omitempty"`
	DescriptionLocalizations *map[string]string          `json:"description_localizations,omitempty"`
	DefaultMemberPermissions *string                     `json:"default_member_permissions"`
	ID                       string                      `json:"id"`
	ApplicationID            string                      `json:"application_id"`
	Description              string                      `json:"description"`
	Name                     string                      `json:"name"`
	Version                  string                      `json:"version,omitempty"`
	Options                  []*ApplicationCommandOption `json:"options,omitempty"`
}

// Application Command Types
// https://discord.com/developers/docs/interactions/application-commands#application-command-object-application-command-types
const (
	FlagApplicationCommandTypeCHAT_INPUT Flag = 1
	FlagApplicationCommandTypeUSER       Flag = 2
	FlagApplicationCommandTypeMESSAGE    Flag = 3
)

// Application Command Option Structure
// https://discord.com/developers/docs/interactions/application-commands#application-command-object-application-command-option-structure
type ApplicationCommandOption struct {
	MaxValue                 *float64                          `json:"max_value,omitempty"`
	Autocomplete             *bool                             `json:"autocomplete,omitempty"`
	NameLocalizations        *map[string]string                `json:"name_localizations,omitempty"`
	MinValue                 *float64                          `json:"min_value,omitempty"`
	DescriptionLocalizations *map[string]string                `json:"description_localizations,omitempty"`
	Required                 *bool                             `json:"required,omitempty"`
	MaxLength                *int                              `json:"max_length,omitempty"`
	MinLength                *int                              `json:"min_length,omitempty"`
	Name                     string                            `json:"name"`
	Description              string                            `json:"description"`
	Options                  []*ApplicationCommandOption       `json:"options,omitempty"`
	ChannelTypes             []Flag                            `json:"channel_types,omitempty"`
	Choices                  []*ApplicationCommandOptionChoice `json:"choices,omitempty"`
	Type                     Flag                              `json:"type"`
}

// Application Command Option Type
// https://discord.com/developers/docs/interactions/application-commands#application-command-object-application-command-option-type
const (
	FlagApplicationCommandOptionTypeSUB_COMMAND       Flag = 1
	FlagApplicationCommandOptionTypeSUB_COMMAND_GROUP Flag = 2
	FlagApplicationCommandOptionTypeSTRING            Flag = 3
	FlagApplicationCommandOptionTypeINTEGER           Flag = 4
	FlagApplicationCommandOptionTypeBOOLEAN           Flag = 5
	FlagApplicationCommandOptionTypeUSER              Flag = 6
	FlagApplicationCommandOptionTypeCHANNEL           Flag = 7
	FlagApplicationCommandOptionTypeROLE              Flag = 8
	FlagApplicationCommandOptionTypeMENTIONABLE       Flag = 9
	FlagApplicationCommandOptionTypeNUMBER            Flag = 10
	FlagApplicationCommandOptionTypeATTACHMENT        Flag = 11
)

// Application Command Option Choice
// https://discord.com/developers/docs/interactions/application-commands#application-command-object-application-command-option-choice-structure
type ApplicationCommandOptionChoice struct {
	Name              string             `json:"name"`
	NameLocalizations *map[string]string `json:"name_localizations,omitempty"`
	Value             Value              `json:"value"`
}

// Guild Application Command Permissions Object
// https://discord.com/developers/docs/interactions/application-commands#application-command-permissions-object-guild-application-command-permissions-structure
type GuildApplicationCommandPermissions struct {
	ID            string                           `json:"id"`
	ApplicationID string                           `json:"application_id"`
	GuildID       string                           `json:"guild_id"`
	Permissions   []*ApplicationCommandPermissions `json:"permissions"`
}

// Application Command Permissions Structure
// https://discord.com/developers/docs/interactions/application-commands#application-command-permissions-object-application-command-permissions-structure
type ApplicationCommandPermissions struct {
	ID         string `json:"id"`
	Type       Flag   `json:"type"`
	Permission bool   `json:"permission"`
}

// Application Command Permission Type
// https://discord.com/developers/docs/interactions/application-commands#application-command-permissions-object-application-command-permission-type
const (
	FlagApplicationCommandPermissionTypeROLE    Flag = 1
	FlagApplicationCommandPermissionTypeUSER    Flag = 2
	FlagApplicationCommandPermissionTypeCHANNEL Flag = 3
)

// Component Object
// https://discord.com/developers/docs/interactions/message-components#component-object
type Component interface {
	ComponentType() Flag
}

// Component Types
// https://discord.com/developers/docs/interactions/message-components#component-object-component-types
const (
	FlagComponentTypeActionRow         Flag = 1
	FlagComponentTypeButton            Flag = 2
	FlagComponentTypeSelectMenu        Flag = 3
	FlagComponentTypeTextInput         Flag = 4
	FlagComponentTypeUserSelect        Flag = 5
	FlagComponentTypeRoleSelect        Flag = 6
	FlagComponentTypeMentionableSelect Flag = 7
	FlagComponentTypeChannelSelect     Flag = 8
)

// https://discord.com/developers/docs/interactions/message-components#component-object
type ActionsRow struct {
	Components []Component `json:"components"`
}

// Button Object
// https://discord.com/developers/docs/interactions/message-components#button-object
type Button struct {
	Label    *string `json:"label,omitempty"`
	Emoji    *Emoji  `json:"emoji,omitempty"`
	CustomID *string `json:"custom_id,omitempty"`
	URL      *string `json:"url,omitempty"`
	Disabled *bool   `json:"disabled,omitempty"`
	Style    Flag    `json:"style"`
}

// Button Styles
// https://discord.com/developers/docs/interactions/message-components#button-object-button-styles
const (
	FlagButtonStylePRIMARY   Flag = 1
	FlagButtonStyleBLURPLE   Flag = 1
	FlagButtonStyleSecondary Flag = 2
	FlagButtonStyleGREY      Flag = 2
	FlagButtonStyleSuccess   Flag = 3
	FlagButtonStyleGREEN     Flag = 3
	FlagButtonStyleDanger    Flag = 4
	FlagButtonStyleRED       Flag = 4
	FlagButtonStyleLINK      Flag = 5
)

// Select Menu Structure
// https://discord.com/developers/docs/interactions/message-components#select-menu-object-select-menu-structure
type SelectMenu struct {
	MaxValues    *Flag              `json:"max_values,omitempty"`
	Disabled     *bool              `json:"disabled,omitempty"`
	Placeholder  *string            `json:"placeholder,omitempty"`
	MinValues    *Flag              `json:"min_values,omitempty"`
	CustomID     string             `json:"custom_id"`
	Options      []SelectMenuOption `json:"options"`
	ChannelTypes []Flag             `json:"channel_types,omitempty"`
	Type         int                `json:"type"`
}

// Select Menu Option Structure
// https://discord.com/developers/docs/interactions/message-components#select-menu-object-select-option-structure
type SelectMenuOption struct {
	Description *string `json:"description,omitempty"`
	Emoji       *Emoji  `json:"emoji,omitempty"`
	Default     *bool   `json:"default,omitempty"`
	Label       string  `json:"label"`
	Value       string  `json:"value"`
}

// Text Input Structure
// https://discord.com/developers/docs/interactions/message-components#text-inputs-text-input-structure
type TextInput struct {
	Value       *string `json:"value,omitempty"`
	Placeholder *string `json:"placeholder,omitempty"`
	Label       *string `json:"label"`
	MinLength   *int    `json:"min_length,omitempty"`
	MaxLength   *int    `json:"max_length,omitempty"`
	Required    *bool   `json:"required,omitempty"`
	CustomID    string  `json:"custom_id"`
	Style       Flag    `json:"style"`
}

// Text Input Styles
// https://discord.com/developers/docs/interactions/message-components#text-inputs-text-input-styles
const (
	FlagTextInputStyleShort     Flag = 1
	FlagTextInputStyleParagraph Flag = 2
)

// Interaction Object
// https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-object-interaction-structure
type Interaction struct {
	Data           InteractionData `json:"data,omitempty"`
	Message        *Message        `json:"message,omitempty"`
	Locale         *string         `json:"locale,omitempty"`
	AppPermissions *BitFlag        `json:"app_permissions,omitempty,string"`
	GuildID        *string         `json:"guild_id,omitempty"`
	ChannelID      *string         `json:"channel_id,omitempty"`
	Member         *GuildMember    `json:"member,omitempty"`
	User           *User           `json:"user,omitempty"`
	GuildLocale    *string         `json:"guild_locale,omitempty"`
	Token          string          `json:"token"`
	ApplicationID  string          `json:"application_id"`
	ID             string          `json:"id"`
	Version        int             `json:"version,omitempty"`
	Type           Flag            `json:"type"`
}

// Interaction Type
// https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-object-interaction-type
const (
	FlagInteractionTypePING                             Flag = 1
	FlagInteractionTypeAPPLICATION_COMMAND              Flag = 2
	FlagInteractionTypeMESSAGE_COMPONENT                Flag = 3
	FlagInteractionTypeAPPLICATION_COMMAND_AUTOCOMPLETE Flag = 4
	FlagInteractionTypeMODAL_SUBMIT                     Flag = 5
)

// Interaction Data
// https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-object-interaction-data
type InteractionData interface {
	InteractionDataType() Flag
}

// Application Command Data Structure
// https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-object-application-command-data-structure
type ApplicationCommandData struct {
	TargetID *string                                    `json:"target_id,omitempty"`
	Resolved *ResolvedData                              `json:"resolved,omitempty"`
	GuildID  *string                                    `json:"guild_id,omitempty"`
	Name     string                                     `json:"name"`
	ID       string                                     `json:"id"`
	Options  []*ApplicationCommandInteractionDataOption `json:"options,omitempty"`
	Type     Flag                                       `json:"type"`
}

// Message Component Data Structure
// https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-object-message-component-data-structure
type MessageComponentData struct {
	CustomID      string              `json:"custom_id"`
	Values        []*SelectMenuOption `json:"values,omitempty"`
	ComponentType Flag                `json:"component_type"`
}

// Modal Submit Data Structure
// https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-object-modal-submit-data-structure
type ModalSubmitData struct {
	CustomID   string      `json:"custom_id"`
	Components []Component `json:"components"`
}

// Resolved Data Structure
// https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-object-resolved-data-structure
type ResolvedData struct {
	Users       map[string]*User        `json:"users,omitempty"`
	Members     map[string]*GuildMember `json:"members,omitempty"`
	Roles       map[string]*Role        `json:"roles,omitempty"`
	Channels    map[string]*Channel     `json:"channels,omitempty"`
	Messages    map[string]*Message     `json:"messages,omitempty"`
	Attachments map[string]*Attachment  `json:"attachments,omitempty"`
}

// Application Command Interaction Data Option Structure
// https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-object-application-command-interaction-data-option-structure
type ApplicationCommandInteractionDataOption struct {
	Value   *Value                                     `json:"value,omitempty"`
	Focused *bool                                      `json:"focused,omitempty"`
	Name    string                                     `json:"name"`
	Options []*ApplicationCommandInteractionDataOption `json:"options,omitempty"`
	Type    Flag                                       `json:"type"`
}

// Message Interaction Structure
// https://discord.com/developers/docs/interactions/receiving-and-responding#message-interaction-object-message-interaction-structure
type MessageInteraction struct {
	User   *User        `json:"user"`
	Member *GuildMember `json:"member,omitempty"`
	ID     string       `json:"id"`
	Name   string       `json:"name"`
	Type   Flag         `json:"type"`
}

// Interaction Response Structure
// https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-response-object-interaction-response-structure
type InteractionResponse struct {
	Data InteractionCallbackData `json:"data,omitempty"`
	Type Flag                    `json:"type"`
}

// Interaction Callback Type
// https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-response-object-interaction-callback-type
const (
	FlagInteractionCallbackTypePONG                                    Flag = 1
	FlagInteractionCallbackTypeCHANNEL_MESSAGE_WITH_SOURCE             Flag = 4
	FlagInteractionCallbackTypeDEFERRED_CHANNEL_MESSAGE_WITH_SOURCE    Flag = 5
	FlagInteractionCallbackTypeDEFERRED_UPDATE_MESSAGE                 Flag = 6
	FlagInteractionCallbackTypeUPDATE_MESSAGE                          Flag = 7
	FlagInteractionCallbackTypeAPPLICATION_COMMAND_AUTOCOMPLETE_RESULT Flag = 8
	FlagInteractionCallbackTypeMODAL                                   Flag = 9
)

// Interaction Callback Data Structure
// https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-response-object-interaction-callback-data-structure
type InteractionCallbackData interface {
	InteractionCallbackDataType() Flag
}

// Messages
// https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-response-object-messages
type Messages struct {
	TTS             *bool            `json:"tts,omitempty"`
	Content         *string          `json:"content,omitempty"`
	Embeds          []*Embed         `json:"embeds,omitempty"`
	AllowedMentions *AllowedMentions `json:"allowed_mentions,omitempty"`
	Flags           *BitFlag         `json:"flags,omitempty"`
	Components      []Component      `json:"components,omitempty"`
	Attachments     []*Attachment    `json:"attachments,omitempty"`
}

// Autocomplete
// https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-response-object-autocomplete
type Autocomplete struct {
	Choices []*ApplicationCommandOptionChoice `json:"choices"`
}

// Modal
// https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-response-object-modal
type Modal struct {
	CustomID   string      `json:"custom_id"`
	Title      string      `json:"title"`
	Components []Component `json:"components"`
}

// Application Object
// https://discord.com/developers/docs/resources/application
type Application struct {
	Slug                *string        `json:"slug,omitempty"`
	GuildID             *string        `json:"guild_id,omitempty"`
	Icon                *string        `json:"icon"`
	Team                *Team          `json:"team"`
	InstallParams       *InstallParams `json:"install_params,omitempty"`
	Flags               *BitFlag       `json:"flags,omitempty"`
	CoverImage          *string        `json:"cover_image,omitempty"`
	TermsOfServiceURL   *string        `json:"terms_of_service_url,omitempty"`
	PrivacyProxyURL     *string        `json:"privacy_policy_url,omitempty"`
	Owner               *User          `json:"owner,omitempty"`
	CustomInstallURL    *string        `json:"custom_install_url,omitempty"`
	PrimarySKUID        *string        `json:"primary_sku_id,omitempty"`
	Description         string         `json:"description"`
	Name                string         `json:"name"`
	VerifyKey           string         `json:"verify_key"`
	ID                  string         `json:"id"`
	Tags                []string       `json:"tags,omitempty"`
	RPCOrigins          []string       `json:"rpc_origins,omitempty"`
	BotRequireCodeGrant bool           `json:"bot_require_code_grant"`
	BotPublic           bool           `json:"bot_public"`
}

// Application Flags
// https://discord.com/developers/docs/resources/application#application-object-application-flags
const (
	FlagApplicationGATEWAY_PRESENCE                 BitFlag = 1 << 12
	FlagApplicationGATEWAY_PRESENCE_LIMITED         BitFlag = 1 << 13
	FlagApplicationGATEWAY_GUILD_MEMBERS            BitFlag = 1 << 14
	FlagApplicationGATEWAY_GUILD_MEMBERS_LIMITED    BitFlag = 1 << 15
	FlagApplicationVERIFICATION_PENDING_GUILD_LIMIT BitFlag = 1 << 16
	FlagApplicationEMBEDDED                         BitFlag = 1 << 17
	FlagApplicationGATEWAY_MESSAGE_CONTENT          BitFlag = 1 << 18
	FlagApplicationGATEWAY_MESSAGE_CONTENT_LIMITED  BitFlag = 1 << 19
	FlagApplicationAPPLICATION_COMMAND_BADGE        BitFlag = 1 << 23
)

// Install Params Object
// https://discord.com/developers/docs/resources/application#install-params-object
type InstallParams struct {
	Permissions string   `json:"permissions"`
	Scopes      []string `json:"scopes"`
}

// Audit Log Object
// https://discord.com/developers/docs/resources/audit-log
type AuditLog struct {
	ApplicationCommands  []*ApplicationCommand  `json:"application_commands"`
	AuditLogEntries      []*AuditLogEntry       `json:"audit_log_entries"`
	GuildScheduledEvents []*GuildScheduledEvent `json:"guild_scheduled_events"`
	Integration          []*Integration         `json:"integrations"`
	Threads              []*Channel             `json:"threads"`
	Users                []*User                `json:"users"`
	Webhooks             []*Webhook             `json:"webhooks"`
}

// Audit Log Entry Object
// https://discord.com/developers/docs/resources/audit-log#audit-log-entry-object-audit-log-entry-structure
type AuditLogEntry struct {
	TargetID   *string           `json:"target_id"`
	UserID     *string           `json:"user_id"`
	Options    *AuditLogOptions  `json:"options,omitempty"`
	Reason     *string           `json:"reason,omitempty"`
	ID         string            `json:"id"`
	Changes    []*AuditLogChange `json:"changes,omitempty"`
	ActionType Flag              `json:"action_type"`
}

// Audit Log Events
// https://discord.com/developers/docs/resources/audit-log#audit-log-entry-object-audit-log-events
const (
	FlagAuditLogEventGUILD_UPDATE                                Flag = 1
	FlagAuditLogEventCHANNEL_CREATE                              Flag = 10
	FlagAuditLogEventCHANNEL_UPDATE                              Flag = 11
	FlagAuditLogEventCHANNEL_DELETE                              Flag = 12
	FlagAuditLogEventCHANNEL_OVERWRITE_CREATE                    Flag = 13
	FlagAuditLogEventCHANNEL_OVERWRITE_UPDATE                    Flag = 14
	FlagAuditLogEventCHANNEL_OVERWRITE_DELETE                    Flag = 15
	FlagAuditLogEventMEMBER_KICK                                 Flag = 20
	FlagAuditLogEventMEMBER_PRUNE                                Flag = 21
	FlagAuditLogEventMEMBER_BAN_ADD                              Flag = 22
	FlagAuditLogEventMEMBER_BAN_REMOVE                           Flag = 23
	FlagAuditLogEventMEMBER_UPDATE                               Flag = 24
	FlagAuditLogEventMEMBER_ROLE_UPDATE                          Flag = 25
	FlagAuditLogEventMEMBER_MOVE                                 Flag = 26
	FlagAuditLogEventMEMBER_DISCONNECT                           Flag = 27
	FlagAuditLogEventBOT_ADD                                     Flag = 28
	FlagAuditLogEventROLE_CREATE                                 Flag = 30
	FlagAuditLogEventROLE_UPDATE                                 Flag = 31
	FlagAuditLogEventROLE_DELETE                                 Flag = 32
	FlagAuditLogEventINVITE_CREATE                               Flag = 40
	FlagAuditLogEventINVITE_UPDATE                               Flag = 41
	FlagAuditLogEventINVITE_DELETE                               Flag = 42
	FlagAuditLogEventWEBHOOK_CREATE                              Flag = 50
	FlagAuditLogEventWEBHOOK_UPDATE                              Flag = 51
	FlagAuditLogEventWEBHOOK_DELETE                              Flag = 52
	FlagAuditLogEventEMOJI_CREATE                                Flag = 60
	FlagAuditLogEventEMOJI_UPDATE                                Flag = 61
	FlagAuditLogEventEMOJI_DELETE                                Flag = 62
	FlagAuditLogEventMESSAGE_DELETE                              Flag = 72
	FlagAuditLogEventMESSAGE_BULK_DELETE                         Flag = 73
	FlagAuditLogEventMESSAGE_PIN                                 Flag = 74
	FlagAuditLogEventMESSAGE_UNPIN                               Flag = 75
	FlagAuditLogEventINTEGRATION_CREATE                          Flag = 80
	FlagAuditLogEventINTEGRATION_UPDATE                          Flag = 81
	FlagAuditLogEventINTEGRATION_DELETE                          Flag = 82
	FlagAuditLogEventSTAGE_INSTANCE_CREATE                       Flag = 83
	FlagAuditLogEventSTAGE_INSTANCE_UPDATE                       Flag = 84
	FlagAuditLogEventSTAGE_INSTANCE_DELETE                       Flag = 85
	FlagAuditLogEventSTICKER_CREATE                              Flag = 90
	FlagAuditLogEventSTICKER_UPDATE                              Flag = 91
	FlagAuditLogEventSTICKER_DELETE                              Flag = 92
	FlagAuditLogEventGUILD_SCHEDULED_EVENT_CREATE                Flag = 100
	FlagAuditLogEventGUILD_SCHEDULED_EVENT_UPDATE                Flag = 101
	FlagAuditLogEventGUILD_SCHEDULED_EVENT_DELETE                Flag = 102
	FlagAuditLogEventTHREAD_CREATE                               Flag = 110
	FlagAuditLogEventTHREAD_UPDATE                               Flag = 111
	FlagAuditLogEventTHREAD_DELETE                               Flag = 112
	FlagAuditLogEventAPPLICATION_COMMAND_PERMISSION_UPDATE       Flag = 121
	FlagAuditLogEventAUTO_MODERATION_RULE_CREATE                 Flag = 140
	FlagAuditLogEventAUTO_MODERATION_RULE_UPDATE                 Flag = 141
	FlagAuditLogEventAUTO_MODERATION_RULE_DELETE                 Flag = 142
	FlagAuditLogEventAUTO_MODERATION_BLOCK_MESSAGE               Flag = 143
	FlagAuditLogEventAUTO_MODERATION_FLAG_TO_CHANNEL             Flag = 144
	FlagAuditLogEventAUTO_MODERATION_USER_COMMUNICATION_DISABLED Flag = 145
)

// Optional Audit Entry Info
// https://discord.com/developers/docs/resources/audit-log#audit-log-entry-object-optional-audit-entry-info
type AuditLogOptions struct {
	ApplicationID                 string `json:"application_id"`
	AutoModerationRuleName        string `json:"auto_moderation_rule_name"`
	AutoModerationRuleTriggerType string `json:"auto_moderation_rule_trigger_type"`
	ChannelID                     string `json:"channel_id"`
	Count                         string `json:"count"`
	DeleteMemberDays              string `json:"delete_member_days"`
	ID                            string `json:"id"`
	MembersRemoved                string `json:"members_removed"`
	MessageID                     string `json:"message_id"`
	RoleName                      string `json:"role_name"`
	Type                          string `json:"type"`
}

// Audit Log Change Object
// https://discord.com/developers/docs/resources/audit-log#audit-log-change-object
type AuditLogChange struct {
	Key      string          `json:"key"`
	NewValue json.RawMessage `json:"new_value,omitempty"`
	OldValue json.RawMessage `json:"old_value,omitempty"`
}

// Audit Log Change Exceptions
// https://discord.com/developers/docs/resources/audit-log#audit-log-change-object-audit-log-change-exceptions

// Auto Moderation Rule Structure
// https://discord.com/developers/docs/resources/auto-moderation#auto-moderation-rule-object-auto-moderation-rule-structure
type AutoModerationRule struct {
	ID              string                  `json:"id"`
	GuildID         string                  `json:"guild_id"`
	Name            string                  `json:"name"`
	CreatorID       string                  `json:"creator_id"`
	ExemptChannels  []string                `json:"exempt_channels"`
	Actions         []*AutoModerationAction `json:"actions"`
	ExemptRoles     []string                `json:"exempt_roles"`
	TriggerMetadata TriggerMetadata         `json:"trigger_metadata"`
	TriggerType     Flag                    `json:"trigger_type"`
	Enabled         bool                    `json:"enabled"`
	EventType       Flag                    `json:"event_type"`
}

// Trigger Types
// https://discord.com/developers/docs/resources/auto-moderation#auto-moderation-rule-object-trigger-types
const (
	FlagTriggerTypeKEYWORD        Flag = 1
	FlagTriggerTypeHARMFUL_LINK   Flag = 2
	FlagTriggerTypeSPAM           Flag = 3
	FlagTriggerTypeKEYWORD_PRESET Flag = 4
	FlagTriggerTypeMENTION_SPAM   Flag = 5
)

// Trigger Metadata
// https://discord.com/developers/docs/resources/auto-moderation#auto-moderation-rule-object-trigger-metadata
type TriggerMetadata struct {
	// https://discord.com/developers/docs/resources/auto-moderation#auto-moderation-rule-object-keyword-matching-strategies
	KeywordFilter     []string `json:"keyword_filter"`
	RegexPatterns     []Flag   `json:"regex_patterns"`
	Presets           []Flag   `json:"presets"`
	AllowList         []string `json:"allow_list"`
	MentionTotalLimit int      `json:"mention_total_limit"`
}

// Keyword Preset Types
// https://discord.com/developers/docs/resources/auto-moderation#auto-moderation-rule-object-keyword-preset-types
const (
	FlagKeywordPresetTypePROFANITY      Flag = 1
	FlagKeywordPresetTypeSEXUAL_CONTENT Flag = 2
	FlagKeywordPresetTypeSLURS          Flag = 3
)

// Event Types
// https://discord.com/developers/docs/resources/auto-moderation#auto-moderation-rule-object-event-types
const (
	FlagEventTypeMESSAGE_SEND Flag = 1
)

// Auto Moderation Action Structure
// https://discord.com/developers/docs/resources/auto-moderation#auto-moderation-action-object
type AutoModerationAction struct {
	Metadata *ActionMetadata `json:"metadata,omitempty"`
	Type     Flag            `json:"type"`
}

// Action Types
// https://discord.com/developers/docs/resources/auto-moderation#auto-moderation-action-object-action-types
const (
	FlagActionTypeBLOCK_MESSAGE      Flag = 1
	FlagActionTypeSEND_ALERT_MESSAGE Flag = 2
	FlagActionTypeTIMEOUT            Flag = 3
)

// Action Metadata
// https://discord.com/developers/docs/resources/auto-moderation#auto-moderation-action-object-action-metadata
type ActionMetadata struct {
	ChannelID       string `json:"channel_id"`
	DurationSeconds int    `json:"duration_seconds"`
}

// Channel Object
// https://discord.com/developers/docs/resources/channel
type Channel struct {
	DefaultSortOrder              **int                  `json:"default_sort_order,omitempty"`
	Type                          *Flag                  `json:"type"`
	GuildID                       *string                `json:"guild_id,omitempty"`
	Position                      *int                   `json:"position,omitempty"`
	DefaultThreadRateLimitPerUser *int                   `json:"default_thread_rate_limit_per_user,omitempty"`
	Name                          **string               `json:"name,omitempty"`
	Topic                         **string               `json:"topic,omitempty"`
	NSFW                          *bool                  `json:"nsfw,omitempty"`
	LastMessageID                 **string               `json:"last_message_id,omitempty"`
	Bitrate                       *int                   `json:"bitrate,omitempty"`
	UserLimit                     *int                   `json:"user_limit,omitempty"`
	RateLimitPerUser              *int                   `json:"rate_limit_per_user,omitempty"`
	DefaultReactionEmoji          *DefaultReaction       `json:"default_reaction_emoji"`
	Icon                          **string               `json:"icon,omitempty"`
	OwnerID                       *string                `json:"owner_id,omitempty"`
	ApplicationID                 *string                `json:"application_id,omitempty"`
	Flags                         *BitFlag               `json:"flags,omitempty"`
	LastPinTimestamp              **time.Time            `json:"last_pin_timestamp,omitempty"`
	RTCRegion                     **string               `json:"rtc_region,omitempty"`
	VideoQualityMode              *Flag                  `json:"video_quality_mode,omitempty"`
	MessageCount                  *int                   `json:"message_count,omitempty"`
	MemberCount                   *int                   `json:"member_count,omitempty"`
	ThreadMetadata                *ThreadMetadata        `json:"thread_metadata,omitempty"`
	Member                        *ThreadMember          `json:"member,omitempty"`
	DefaultAutoArchiveDuration    *int                   `json:"default_auto_archive_duration,omitempty"`
	Permissions                   *string                `json:"permissions,omitempty"`
	ParentID                      **string               `json:"parent_id,omitempty"`
	TotalMessageSent              *int                   `json:"total_message_sent,omitempty"`
	ID                            string                 `json:"id"`
	AvailableTags                 []*ForumTag            `json:"available_tags,omitempty"`
	AppliedTags                   []string               `json:"applied_tags,omitempty"`
	Recipients                    []*User                `json:"recipients,omitempty"`
	PermissionOverwrites          []*PermissionOverwrite `json:"permission_overwrites,omitempty"`
}

// Channel Types
// https://discord.com/developers/docs/resources/channel#channel-object-channel-types
const (
	FlagChannelTypeGUILD_TEXT          Flag = 0
	FlagChannelTypeDM                  Flag = 1
	FlagChannelTypeGUILD_VOICE         Flag = 2
	FlagChannelTypeGROUP_DM            Flag = 3
	FlagChannelTypeGUILD_CATEGORY      Flag = 4
	FlagChannelTypeGUILD_ANNOUNCEMENT  Flag = 5
	FlagChannelTypeANNOUNCEMENT_THREAD Flag = 10
	FlagChannelTypePUBLIC_THREAD       Flag = 11
	FlagChannelTypePRIVATE_THREAD      Flag = 12
	FlagChannelTypeGUILD_STAGE_VOICE   Flag = 13
	FlagChannelTypeGUILD_DIRECTORY     Flag = 14
	FlagChannelTypeGUILD_FORUM         Flag = 15
)

// Video Quality Modes
// https://discord.com/developers/docs/resources/channel#channel-object-video-quality-modes
const (
	FlagVideoQualityModeAUTO Flag = 1
	FlagVideoQualityModeFULL Flag = 2
)

// Channel Flags
// https://discord.com/developers/docs/resources/channel#channel-object-channel-flags
const (
	FlagChannelPINNED      BitFlag = 1 << 1
	FlagChannelREQUIRE_TAG BitFlag = 1 << 4
)

// Sort Order Types
// https://discord.com/developers/docs/resources/channel#channel-object-sort-order-types
const (
	FlagSortOrderTypeLATEST_ACTIVITY Flag = 0
	FlagSortOrderTypeCREATION_DATE   Flag = 1
)

// Message Object
// https://discord.com/developers/docs/resources/channel#message-object
type Message struct {
	Timestamp         time.Time         `json:"timestamp"`
	WebhookID         *string           `json:"webhook_id,omitempty"`
	Author            *User             `json:"author"`
	Position          *int              `json:"position,omitempty"`
	GuildID           *string           `json:"guild_id,omitempty"`
	EditedTimestamp   *time.Time        `json:"edited_timestamp"`
	Thread            *Channel          `json:"thread"`
	Interaction       *Interaction      `json:"interaction"`
	ReferencedMessage **Message         `json:"referenced_message,omitempty"`
	Flags             *BitFlag          `json:"flags,omitempty"`
	MessageReference  *MessageReference `json:"message_reference,omitempty"`
	ApplicationID     *string           `json:"application_id,omitempty"`
	Application       *Application      `json:"application,omitempty"`
	Activity          *MessageActivity  `json:"activity,omitempty"`
	Nonce             *Nonce            `json:"nonce,omitempty"`
	Member            *GuildMember      `json:"member,omitempty"`
	ChannelID         string            `json:"channel_id"`
	ID                string            `json:"id"`
	Content           string            `json:"content"`
	Stickers          []*Sticker        `json:"stickers"`
	Attachments       []*Attachment     `json:"attachments"`
	MentionChannels   []*ChannelMention `json:"mention_channels,omitempty"`
	MentionRoles      []*string         `json:"mention_roles"`
	Mentions          []*User           `json:"mentions"`
	Reactions         []*Reaction       `json:"reactions,omitempty"`
	Embeds            []*Embed          `json:"embeds"`
	Components        []Component       `json:"components"`
	StickerItems      []*StickerItem    `json:"sticker_items"`
	MentionEveryone   bool              `json:"mention_everyone"`
	TTS               bool              `json:"tts"`
	Type              Flag              `json:"type"`
	Pinned            bool              `json:"pinned"`
}

// Message Types
// https://discord.com/developers/docs/resources/channel#message-object-message-types
const (
	FlagMessageTypeDEFAULT                                      Flag = 0
	FlagMessageTypeRECIPIENT_ADD                                Flag = 1
	FlagMessageTypeRECIPIENT_REMOVE                             Flag = 2
	FlagMessageTypeCALL                                         Flag = 3
	FlagMessageTypeCHANNEL_NAME_CHANGE                          Flag = 4
	FlagMessageTypeCHANNEL_ICON_CHANGE                          Flag = 5
	FlagMessageTypeCHANNEL_PINNED_MESSAGE                       Flag = 6
	FlagMessageTypeUSER_JOIN                                    Flag = 7
	FlagMessageTypeGUILD_BOOST                                  Flag = 8
	FlagMessageTypeGUILD_BOOST_TIER_1                           Flag = 9
	FlagMessageTypeGUILD_BOOST_TIER_2                           Flag = 10
	FlagMessageTypeGUILD_BOOST_TIER_3                           Flag = 11
	FlagMessageTypeCHANNEL_FOLLOW_ADD                           Flag = 12
	FlagMessageTypeGUILD_DISCOVERY_DISQUALIFIED                 Flag = 14
	FlagMessageTypeGUILD_DISCOVERY_REQUALIFIED                  Flag = 15
	FlagMessageTypeGUILD_DISCOVERY_GRACE_PERIOD_INITIAL_WARNING Flag = 16
	FlagMessageTypeGUILD_DISCOVERY_GRACE_PERIOD_FINAL_WARNING   Flag = 17
	FlagMessageTypeTHREAD_CREATED                               Flag = 18
	FlagMessageTypeREPLY                                        Flag = 19
	FlagMessageTypeCHAT_INPUT_COMMAND                           Flag = 20
	FlagMessageTypeTHREAD_STARTER_MESSAGE                       Flag = 21
	FlagMessageTypeGUILD_INVITE_REMINDER                        Flag = 22
	FlagMessageTypeCONTEXT_MENU_COMMAND                         Flag = 23
	FlagMessageTypeAUTO_MODERATION_ACTION                       Flag = 24
)

// Message Activity Structure
// https://discord.com/developers/docs/resources/channel#message-object-message-activity-structure
type MessageActivity struct {
	PartyID *string `json:"party_id,omitempty"`
	Type    int     `json:"type"`
}

// Message Activity Types
// https://discord.com/developers/docs/resources/channel#message-object-message-activity-types
const (
	FlagMessageActivityTypeJOIN         Flag = 1
	FlagMessageActivityTypeSPECTATE     Flag = 2
	FlagMessageActivityTypeLISTEN       Flag = 3
	FlagMessageActivityTypeJOIN_REQUEST Flag = 5
)

// Message Flags
// https://discord.com/developers/docs/resources/channel#message-object-message-flags
const (
	FlagMessageCROSSPOSTED                            BitFlag = 1 << 0
	FlagMessageIS_CROSSPOST                           BitFlag = 1 << 1
	FlagMessageSUPPRESS_EMBEDS                        BitFlag = 1 << 2
	FlagMessageSOURCE_MESSAGE_DELETED                 BitFlag = 1 << 3
	FlagMessageURGENT                                 BitFlag = 1 << 4
	FlagMessageHAS_THREAD                             BitFlag = 1 << 5
	FlagMessageEPHEMERAL                              BitFlag = 1 << 6
	FlagMessageLOADING                                BitFlag = 1 << 7
	FlagMessageFAILED_TO_MENTION_SOME_ROLES_IN_THREAD BitFlag = 1 << 8
)

// Message Reference Object
// https://discord.com/developers/docs/resources/channel#message-reference-object
type MessageReference struct {
	MessageID       *string `json:"message_id,omitempty"`
	ChannelID       *string `json:"channel_id,omitempty"`
	GuildID         *string `json:"guild_id,omitempty"`
	FailIfNotExists *bool   `json:"fail_if_not_exists,omitempty"`
}

// Followed Channel Structure
// https://discord.com/developers/docs/resources/channel#followed-channel-object-followed-channel-structure
type FollowedChannel struct {
	ChannelID string `json:"channel_id"`
	WebhookID string `json:"webhook_id"`
}

// Reaction Object
// https://discord.com/developers/docs/resources/channel#reaction-object
type Reaction struct {
	Emoji *Emoji `json:"emoji"`
	Count int    `json:"count"`
	Me    bool   `json:"me"`
}

// Overwrite Object
// https://discord.com/developers/docs/resources/channel#overwrite-object
type PermissionOverwrite struct {
	ID    string `json:"id"`
	Deny  string `json:"deny"`
	Allow string `json:"allow"`
	Type  Flag   `json:"type"`
}

// Thread Metadata Object
// https://discord.com/developers/docs/resources/channel#thread-metadata-object
type ThreadMetadata struct {
	ArchiveTimestamp    time.Time   `json:"archive_timestamp"`
	Invitable           *bool       `json:"invitable,omitempty"`
	CreateTimestamp     **time.Time `json:"create_timestamp"`
	AutoArchiveDuration int         `json:"auto_archive_duration"`
	Archived            bool        `json:"archived"`
	Locked              bool        `json:"locked"`
}

// Thread Member Object
// https://discord.com/developers/docs/resources/channel#thread-member-object
type ThreadMember struct {
	ThreadID      *string   `json:"id,omitempty"`
	UserID        *string   `json:"user_id,omitempty"`
	JoinTimestamp time.Time `json:"join_timestamp"`
	Flags         Flag      `json:"flags"`
}

// Default Reaction Structure
// https://discord.com/developers/docs/resources/channel#default-reaction-object-default-reaction-structure
type DefaultReaction struct {
	EmojiID   *string `json:"emoji_id"`
	EmojiName *string `json:"emoji_name"`
}

// Forum Tag Structure
// https://discord.com/developers/docs/resources/channel#forum-tag-object-forum-tag-structure
type ForumTag struct {
	EmojiName *string `json:"emoji_name"`
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	EmojiID   string  `json:"emoji_id"`
	Moderated bool    `json:"moderated"`
}

// Embed Object
// https://discord.com/developers/docs/resources/channel#embed-object
type Embed struct {
	Title       *string         `json:"title,omitempty"`
	Type        *string         `json:"type,omitempty"`
	Description *string         `json:"description,omitempty"`
	URL         *string         `json:"url,omitempty"`
	Timestamp   *time.Time      `json:"timestamp,omitempty"`
	Color       *int            `json:"color,omitempty"`
	Footer      *EmbedFooter    `json:"footer,omitempty"`
	Image       *EmbedImage     `json:"image,omitempty"`
	Thumbnail   *EmbedThumbnail `json:"thumbnail,omitempty"`
	Video       *EmbedVideo     `json:"video,omitempty"`
	Provider    *EmbedProvider  `json:"provider,omitempty"`
	Author      *EmbedAuthor    `json:"author,omitempty"`
	Fields      []*EmbedField   `json:"fields,omitempty"`
}

// Embed Thumbnail Structure
// https://discord.com/developers/docs/resources/channel#embed-object-embed-thumbnail-structure
type EmbedThumbnail struct {
	ProxyURL *string `json:"proxy_url,omitempty"`
	Height   *int    `json:"height,omitempty"`
	Width    *int    `json:"width,omitempty"`
	URL      string  `json:"url"`
}

// Embed Video Structure
// https://discord.com/developers/docs/resources/channel#embed-object-embed-video-structure
type EmbedVideo struct {
	URL      *string `json:"url,omitempty"`
	ProxyURL *string `json:"proxy_url,omitempty"`
	Height   *int    `json:"height,omitempty"`
	Width    *int    `json:"width,omitempty"`
}

// Embed Image Structure
// https://discord.com/developers/docs/resources/channel#embed-object-embed-image-structure
type EmbedImage struct {
	ProxyURL *string `json:"proxy_url,omitempty"`
	Height   *int    `json:"height,omitempty"`
	Width    *int    `json:"width,omitempty"`
	URL      string  `json:"url"`
}

// Embed Provider Structure
// https://discord.com/developers/docs/resources/channel#embed-object-embed-provider-structure
type EmbedProvider struct {
	Name *string `json:"name,omitempty"`
	URL  *string `json:"url,omitempty"`
}

// Embed Author Structure
// https://discord.com/developers/docs/resources/channel#embed-object-embed-author-structure
type EmbedAuthor struct {
	URL          *string `json:"url,omitempty"`
	IconURL      *string `json:"icon_url,omitempty"`
	ProxyIconURL *string `json:"proxy_icon_url,omitempty"`
	Name         string  `json:"name"`
}

// Embed Footer Structure
// https://discord.com/developers/docs/resources/channel#embed-object-embed-footer-structure
type EmbedFooter struct {
	IconURL      *string `json:"icon_url,omitempty"`
	ProxyIconURL *string `json:"proxy_icon_url,omitempty"`
	Text         string  `json:"text"`
}

// Embed Field Structure
// https://discord.com/developers/docs/resources/channel#embed-object-embed-field-structure
type EmbedField struct {
	Inline *bool  `json:"inline,omitempty"`
	Name   string `json:"name"`
	Value  string `json:"value"`
}

// Embed Limits
// https://discord.com/developers/docs/resources/channel#embed-object-embed-limits
const (
	FlagEmbedLimitTitle       = 256
	FlagEmbedLimitDescription = 4096
	FlagEmbedLimitFields      = 25
	FlagEmbedLimitFieldName   = 256
	FlagEmbedLimitFieldValue  = 1024
	FlagEmbedLimitFooterText  = 2048
	FlagEmbedLimitAuthorName  = 256
)

// Message Attachment Object
// https://discord.com/developers/docs/resources/channel#attachment-object-attachment-structure
type Attachment struct {
	Height      **int   `json:"height,omitempty"`
	Description *string `json:"description,omitempty"`
	ContentType *string `json:"content_type,omitempty"`
	Width       **int   `json:"width,omitempty"`
	Emphemeral  *bool   `json:"ephemeral,omitempty"`
	ID          string  `json:"id"`
	URL         string  `json:"url"`
	ProxyURL    string  `json:"proxy_url"`
	Filename    string  `json:"filename"`
	Size        int     `json:"size"`
}

// Channel Mention Object
// https://discord.com/developers/docs/resources/channel#channel-mention-object
type ChannelMention struct {
	ID      string `json:"id"`
	GuildID string `json:"guild_id"`
	Type    *Flag  `json:"type"`
	Name    string `json:"name"`
}

// Allowed Mentions Structure
// https://discord.com/developers/docs/resources/channel#allowed-mentions-object-allowed-mentions-structure
type AllowedMentions struct {
	Parse       []*string `json:"parse"`
	Roles       []*string `json:"roles"`
	Users       []*string `json:"users"`
	RepliedUser bool      `json:"replied_user"`
}

// Allowed Mention Types
// https://discord.com/developers/docs/resources/channel#allowed-mentions-object-allowed-mentions-structure
const (
	FlagAllowedMentionTypeRoles    = "roles"
	FlagAllowedMentionTypeUsers    = "users"
	FlagAllowedMentionTypeEveryone = "everyone"
)

// Emoji Object
// https://discord.com/developers/docs/resources/emoji#emoji-object-emoji-structure
type Emoji struct {
	ID            *string  `json:"id"`
	Name          *string  `json:"name,omitempty"`
	User          *User    `json:"user,omitempty"`
	RequireColons *bool    `json:"require_colons,omitempty"`
	Managed       *bool    `json:"managed,omitempty"`
	Animated      *bool    `json:"animated,omitempty"`
	Available     *bool    `json:"available,omitempty"`
	Roles         []string `json:"roles,omitempty"`
}

// Guild Object
// https://discord.com/developers/docs/resources/guild#guild-object
type Guild struct {
	RulesChannelID              *string        `json:"rules_channel_id"`
	WelcomeScreen               *WelcomeScreen `json:"welcome_screen,omitempty"`
	Icon                        *string        `json:"icon"`
	IconHash                    **string       `json:"icon_hash,omitempty"`
	Splash                      *string        `json:"splash"`
	DiscoverySplash             *string        `json:"discovery_splash"`
	Owner                       *bool          `json:"owner,omitempty"`
	ApproximatePresenceCount    *int           `json:"approximate_presence_count,omitempty"`
	Permissions                 *string        `json:"permissions,omitempty"`
	Region                      **string       `json:"region,omitempty"`
	ApproximateMemberCount      *int           `json:"approximate_member_count,omitempty"`
	MaxVideoChannelUsers        *int           `json:"max_video_channel_users,omitempty"`
	WidgetEnabled               *bool          `json:"widget_enabled,omitempty"`
	WidgetChannelID             **string       `json:"widget_channel_id,omitempty"`
	PublicUpdatesChannelID      *string        `json:"public_updates_channel_id"`
	PremiumSubscriptionCount    *int           `json:"premium_subscription_count,omitempty"`
	ApplicationID               *string        `json:"application_id"`
	Banner                      *string        `json:"banner"`
	Description                 *string        `json:"description"`
	VanityUrl                   *string        `json:"vanity_url_code"`
	MaxPresences                **int          `json:"max_presences,omitempty"`
	MaxMembers                  *int           `json:"max_members,omitempty"`
	SystemChannelID             *string        `json:"system_channel_id"`
	AfkChannelID                *string        `json:"afk_channel_id"`
	Unavailable                 *bool          `json:"unavailable,omitempty"`
	OwnerID                     string         `json:"owner_id"`
	ID                          string         `json:"id"`
	Name                        string         `json:"name"`
	PreferredLocale             string         `json:"preferred_locale"`
	Emojis                      []*Emoji       `json:"emojis"`
	Roles                       []*Role        `json:"roles"`
	Stickers                    []*Sticker     `json:"stickers,omitempty"`
	Features                    []*string      `json:"features"`
	AfkTimeout                  int            `json:"afk_timeout"`
	SystemChannelFlags          BitFlag        `json:"system_channel_flags"`
	DefaultMessageNotifications Flag           `json:"default_message_notifications"`
	PremiumTier                 Flag           `json:"premium_tier"`
	ExplicitContentFilter       Flag           `json:"explicit_content_filter"`
	NSFWLevel                   Flag           `json:"nsfw_level"`
	MFALevel                    Flag           `json:"mfa_level"`
	PremiumProgressBarEnabled   bool           `json:"premium_progress_bar_enabled"`
	VerificationLevel           Flag           `json:"verification_level"`
}

// Default Message Notification Level
// https://discord.com/developers/docs/resources/guild#guild-object-default-message-notification-level
const (
	FlagDefaultMessageNotificationLevelALL_MESSAGES  Flag = 0
	FlagDefaultMessageNotificationLevelONLY_MENTIONS Flag = 1
)

// Explicit Content Filter Level
// https://discord.com/developers/docs/resources/guild#guild-object-explicit-content-filter-level
const (
	FlagExplicitContentFilterLevelDISABLED              Flag = 0
	FlagExplicitContentFilterLevelMEMBERS_WITHOUT_ROLES Flag = 1
	FlagExplicitContentFilterLevelALL_MEMBERS           Flag = 2
)

// MFA Level
// https://discord.com/developers/docs/resources/guild#guild-object-mfa-level
const (
	FlagMFALevelNONE     Flag = 0
	FlagMFALevelELEVATED Flag = 1
)

// Verification Level
// https://discord.com/developers/docs/resources/guild#guild-object-verification-level
const (
	FlagVerificationLevelNONE      Flag = 0
	FlagVerificationLevelLOW       Flag = 1
	FlagVerificationLevelMEDIUM    Flag = 2
	FlagVerificationLevelHIGH      Flag = 3
	FlagVerificationLevelVERY_HIGH Flag = 4
)

// Guild NSFW Level
// https://discord.com/developers/docs/resources/guild#guild-object-guild-nsfw-level
const (
	FlagGuildNSFWLevelDEFAULT        Flag = 0
	FlagGuildNSFWLevelEXPLICIT       Flag = 1
	FlagGuildNSFWLevelSAFE           Flag = 2
	FlagGuildNSFWLevelAGE_RESTRICTED Flag = 3
)

// Premium Tier
// https://discord.com/developers/docs/resources/guild#guild-object-premium-tier
const (
	FlagPremiumTierNONE  Flag = 0
	FlagPremiumTierONE   Flag = 1
	FlagPremiumTierTWO   Flag = 2
	FlagPremiumTierTHREE Flag = 3
)

// System Channel Flags
// https://discord.com/developers/docs/resources/guild#guild-object-system-channel-flags
const (
	FlagSystemChannelSUPPRESS_JOIN_NOTIFICATIONS           BitFlag = 1 << 0
	FlagSystemChannelSUPPRESS_PREMIUM_SUBSCRIPTIONS        BitFlag = 1 << 1
	FlagSystemChannelSUPPRESS_GUILD_REMINDER_NOTIFICATIONS BitFlag = 1 << 2
	FlagSystemChannelSUPPRESS_JOIN_NOTIFICATION_REPLIES    BitFlag = 1 << 3
)

// Guild Features
// https://discord.com/developers/docs/resources/guild#guild-object-guild-features
const (
	FlagGuildFeatureANIMATED_BANNER                    = "ANIMATED_BANNER"
	FlagGuildFeatureANIMATED_ICON                      = "ANIMATED_ICON"
	FlagGuildFeatureAPPLICATION_COMMAND_PERMISSIONS_V2 = "APPLICATION_COMMAND_PERMISSIONS_V2"
	FlagGuildFeatureAUTO_MODERATION                    = "AUTO_MODERATION"
	FlagGuildFeatureBANNER                             = "BANNER"
	FlagGuildFeatureCOMMUNITY                          = "COMMUNITY"
	FlagGuildFeatureDEVELOPER_SUPPORT_SERVER           = "DEVELOPER_SUPPORT_SERVER"
	FlagGuildFeatureDISCOVERABLE                       = "DISCOVERABLE"
	FlagGuildFeatureFEATURABLE                         = "FEATURABLE"
	FlagGuildFeatureINVITES_DISABLED                   = "INVITES_DISABLED"
	FlagGuildFeatureINVITE_SPLASH                      = "INVITE_SPLASH"
	FlagGuildFeatureMEMBER_VERIFICATION_GATE_ENABLED   = "MEMBER_VERIFICATION_GATE_ENABLED"
	FlagGuildFeatureMONETIZATION_ENABLED               = "MONETIZATION_ENABLED"
	FlagGuildFeatureMORE_STICKERS                      = "MORE_STICKERS"
	FlagGuildFeatureNEWS                               = "NEWS"
	FlagGuildFeaturePARTNERED                          = "PARTNERED"
	FlagGuildFeaturePREVIEW_ENABLED                    = "PREVIEW_ENABLED"
	FlagGuildFeatureROLE_ICONS                         = "ROLE_ICONS"
	FlagGuildFeatureTICKETED_EVENTS_ENABLED            = "TICKETED_EVENTS_ENABLED"
	FlagGuildFeatureVANITY_URL                         = "VANITY_URL"
	FlagGuildFeatureVERIFIED                           = "VERIFIED"
	FlagGuildFeatureVIP_REGIONS                        = "VIP_REGIONS"
	FlagGuildFeatureWELCOME_SCREEN_ENABLED             = "WELCOME_SCREEN_ENABLED"
)

// Mutable Guild Features
// https://discord.com/developers/docs/resources/guild#guild-object-mutable-guild-features
var (
	MutableGuildFeatures = map[string]bool{
		FlagGuildFeatureCOMMUNITY:        true,
		FlagGuildFeatureINVITES_DISABLED: true,
		FlagGuildFeatureDISCOVERABLE:     true,
	}
)

// Guild Preview Object
// https://discord.com/developers/docs/resources/guild#guild-preview-object-guild-preview-structure
type GuildPreview struct {
	Icon                     *string    `json:"icon"`
	Splash                   *string    `json:"splash"`
	DiscoverySplash          *string    `json:"discovery_splash"`
	Description              *string    `json:"description"`
	ID                       string     `json:"id"`
	Name                     string     `json:"name"`
	Emojis                   []*Emoji   `json:"emojis"`
	Features                 []*string  `json:"features"`
	Stickers                 []*Sticker `json:"stickers"`
	ApproximateMemberCount   int        `json:"approximate_member_count"`
	ApproximatePresenceCount int        `json:"approximate_presence_count"`
}

// Guild Widget Settings Object
// https://discord.com/developers/docs/resources/guild#guild-widget-settings-object
type GuildWidgetSettings struct {
	ChannelID *string `json:"channel_id"`
	Enabled   bool    `json:"enabled"`
}

// Guild Widget Object
// https://discord.com/developers/docs/resources/guild#et-gguild-widget-object-get-guild-widget-structure*
type GuildWidget struct {
	ID            string     `json:"id"`
	Name          string     `json:"name"`
	InstantInvite *string    `json:"instant_invite"`
	Channels      []*Channel `json:"channels"`
	Members       []*User    `json:"members"`
	PresenceCount int        `json:"presence_count"`
}

// Guild Member Object
// https://discord.com/developers/docs/resources/guild#guild-member-object
type GuildMember struct {
	JoinedAt                   time.Time   `json:"joined_at"`
	User                       *User       `json:"user,omitempty"`
	Nick                       **string    `json:"nick,omitempty"`
	Avatar                     **string    `json:"avatar,omitempty"`
	Permissions                *string     `json:"permissions,omitempty"`
	PremiumSince               **time.Time `json:"premium_since,omitempty"`
	Pending                    *bool       `json:"pending,omitempty"`
	CommunicationDisabledUntil **time.Time `json:"communication_disabled_until,omitempty"`
	Roles                      []*string   `json:"roles"`
	Deaf                       bool        `json:"deaf"`
	Mute                       bool        `json:"mute"`
}

// Integration Object
// https://discord.com/developers/docs/resources/guild#integration-object
type Integration struct {
	ID                string             `json:"id"`
	Name              string             `json:"name"`
	Type              string             `json:"type"`
	Enabled           *bool              `json:"enabled,omitempty"`
	Syncing           *bool              `json:"syncing,omitempty"`
	RoleID            *string            `json:"role_id,omitempty"`
	EnableEmoticons   *bool              `json:"enable_emoticons,omitempty"`
	ExpireBehavior    *Flag              `json:"expire_behavior,omitempty"`
	ExpireGracePeriod *int               `json:"expire_grace_period,omitempty"`
	User              *User              `json:"user,omitempty"`
	Account           IntegrationAccount `json:"account"`
	SyncedAt          *time.Time         `json:"synced_at,omitempty"`
	SubscriberCount   *int               `json:"subscriber_count,omitempty"`
	Revoked           *bool              `json:"revoked,omitempty"`
	Application       *Application       `json:"application,omitempty"`
	Scopes            []string           `json:"scopes,omitempty"`
}

// Integration Expire Behaviors
// https://discord.com/developers/docs/resources/guild#integration-object-integration-expire-behaviors
const (
	FlagIntegrationExpireBehaviorREMOVEROLE Flag = 0
	FlagIntegrationExpireBehaviorKICK       Flag = 1
)

// Integration Account Object
// https://discord.com/developers/docs/resources/guild#integration-account-object
type IntegrationAccount struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Integration Application Object
// https://discord.com/developers/docs/resources/guild#integration-application-object-integration-application-structure
type IntegrationApplication struct {
	Icon        *string `json:"icon"`
	Bot         *User   `json:"bot,omitempty"`
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
}

// Guild Ban Object
// https://discord.com/developers/docs/resources/guild#ban-object
type Ban struct {
	Reason *string `json:"reason"`
	User   *User   `json:"user"`
}

// Welcome Screen Object
// https://discord.com/developers/docs/resources/guild#welcome-screen-object-welcome-screen-structure
type WelcomeScreen struct {
	Description           *string                 `json:"description"`
	WelcomeScreenChannels []*WelcomeScreenChannel `json:"welcome_channels"`
}

// Welcome Screen Channel Structure
// https://discord.com/developers/docs/resources/guild#welcome-screen-object-welcome-screen-channel-structure
type WelcomeScreenChannel struct {
	Description *string `json:"description"`
	EmojiID     *string `json:"emoji_id"`
	EmojiName   *string `json:"emoji_name"`
	ChannelID   string  `json:"channel_id"`
}

// Guild Scheduled Event Object
// https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event-object-guild-scheduled-event-structure
type GuildScheduledEvent struct {
	ScheduledStartTime time.Time                          `json:"scheduled_start_time"`
	EntityMetadata     *GuildScheduledEventEntityMetadata `json:"entity_metadata"`
	ChannelID          *string                            `json:"channel_id"`
	CreatorID          **string                           `json:"creator_id,omitempty"`
	EntityID           *string                            `json:"entity_id"`
	Description        **string                           `json:"description,omitempty"`
	Creator            *User                              `json:"creator,omitempty"`
	ScheduledEndTime   *time.Time                         `json:"scheduled_end_time"`
	Image              **string                           `json:"image,omitempty"`
	UserCount          *int                               `json:"user_count,omitempty"`
	ID                 string                             `json:"id"`
	Name               string                             `json:"name"`
	GuildID            string                             `json:"guild_id"`
	EntityType         Flag                               `json:"entity_type"`
	Status             Flag                               `json:"status"`
	PrivacyLevel       Flag                               `json:"privacy_level"`
}

// Guild Scheduled Event Privacy Level
// https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event-object-guild-scheduled-event-privacy-level
const (
	FlagGuildScheduledEventPrivacyLevelGUILD_ONLY Flag = 2
)

// Guild Scheduled Event Entity Types
// https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event-object-guild-scheduled-event-entity-types
const (
	FlagGuildScheduledEventEntityTypeSTAGE_INSTANCE Flag = 1
	FlagGuildScheduledEventEntityTypeVOICE          Flag = 2
	FlagGuildScheduledEventEntityTypeEXTERNAL       Flag = 3
)

// Guild Scheduled Event Status
// https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event-object-guild-scheduled-event-status
const (
	FlagGuildScheduledEventStatusSCHEDULED Flag = 1
	FlagGuildScheduledEventStatusACTIVE    Flag = 2
	FlagGuildScheduledEventStatusCOMPLETED Flag = 3
	FlagGuildScheduledEventStatusCANCELED  Flag = 4
)

// Guild Scheduled Event Entity Metadata
// https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event-object-guild-scheduled-event-entity-metadata
type GuildScheduledEventEntityMetadata struct {
	Location string `json:"location,omitempty"`
}

// Guild Scheduled Event User Object
// https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event-user-object-guild-scheduled-event-user-structure
type GuildScheduledEventUser struct {
	User                  *User        `json:"user"`
	Member                *GuildMember `json:"member,omitempty"`
	GuildScheduledEventID string       `json:"guild_scheduled_event_id"`
}

// Guild Template Object
// https://discord.com/developers/docs/resources/guild-template#guild-template-object
type GuildTemplate struct {
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
	Description           *string   `json:"description"`
	SerializedSourceGuild *Guild    `json:"serialized_source_guild"`
	IsDirty               *bool     `json:"is_dirty"`
	Creator               *User     `json:"creator"`
	Code                  string    `json:"code"`
	Name                  string    `json:"name"`
	SourceGuildID         string    `json:"source_guild_id"`
	CreatorID             string    `json:"creator_id"`
	UsageCount            int       `json:"usage_count"`
}

// Invite Object
// https://discord.com/developers/docs/resources/invite#invite-object
type Invite struct {
	TargetApplication        *Application         `json:"target_application,omitempty"`
	Guild                    *Guild               `json:"guild,omitempty"`
	Channel                  *Channel             `json:"channel"`
	Inviter                  *User                `json:"inviter,omitempty"`
	TargetType               *Flag                `json:"target_type,omitempty"`
	TargetUser               *User                `json:"target_user,omitempty"`
	ApproximatePresenceCount *int                 `json:"approximate_presence_count,omitempty"`
	ApproximateMemberCount   *int                 `json:"approximate_member_count,omitempty"`
	ExpiresAt                **time.Time          `json:"expires_at,omitempty"`
	GuildScheduledEvent      *GuildScheduledEvent `json:"guild_scheduled_event,omitempty"`
	Code                     string               `json:"code"`
}

// Invite Target Types
// https://discord.com/developers/docs/resources/invite#invite-object-invite-target-types
const (
	FlagInviteTargetTypeSTREAM               Flag = 1
	FlagInviteTargetTypeEMBEDDED_APPLICATION Flag = 2
)

// Invite Metadata Object
// https://discord.com/developers/docs/resources/invite#invite-metadata-object-invite-metadata-structure
type InviteMetadata struct {
	CreatedAt time.Time `json:"created_at"`
	Uses      int       `json:"uses"`
	MaxUses   int       `json:"max_uses"`
	MaxAge    int       `json:"max_age"`
	Temporary bool      `json:"temporary"`
}

// Stage Instance Object
// https://discord.com/developers/docs/resources/stage-instance#stage-instance-object
type StageInstance struct {
	GuildScheduledEventID *string `json:"guild_scheduled_event_id"`
	GuildID               string  `json:"guild_id"`
	ChannelID             string  `json:"channel_id"`
	Topic                 string  `json:"topic"`
	ID                    string  `json:"id"`
	PrivacyLevel          Flag    `json:"privacy_level"`
	DiscoverableDisabled  bool    `json:"discoverable_disabled"`
}

// Stage Instance Privacy Level
// https://discord.com/developers/docs/resources/stage-instance#stage-instance-object-privacy-level
const (
	FlagStageInstancePrivacyLevelGUILD_ONLY Flag = 2
)

// Sticker Structure
// https://discord.com/developers/docs/resources/sticker#sticker-object-sticker-structure
type Sticker struct {
	PackID      *string `json:"pack_id,omitempty"`
	Available   *bool   `json:"available,omitempty"`
	Description *string `json:"description"`
	User        *User   `json:"user,omitempty"`
	Asset       *string `json:"asset,omitempty"`
	GuildID     *string `json:"guild_id,omitempty"`
	SortValue   *int    `json:"sort_value,omitempty"`
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Tags        string  `json:"tags"`
	Type        Flag    `json:"type"`
	FormatType  Flag    `json:"format_type"`
}

// Sticker Types
// https://discord.com/developers/docs/resources/sticker#sticker-object-sticker-types
const (
	FlagStickerTypeSTANDARD Flag = 1
	FlagStickerTypeGUILD    Flag = 2
)

// Sticker Format Types
// https://discord.com/developers/docs/resources/sticker#sticker-object-sticker-format-types
const (
	FlagStickerFormatTypePNG    Flag = 1
	FlagStickerFormatTypeAPNG   Flag = 2
	FlagStickerFormatTypeLOTTIE Flag = 3
)

// Sticker Item Object
// https://discord.com/developers/docs/resources/sticker#sticker-item-object
type StickerItem struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	FormatType Flag   `json:"format_type"`
}

// Sticker Pack Object
// https://discord.com/developers/docs/resources/sticker#sticker-pack-object-sticker-pack-structure
type StickerPack struct {
	BannerAssetID  *string    `json:"banner_asset_id,omitempty"`
	CoverStickerID *string    `json:"cover_sticker_id,omitempty"`
	Name           string     `json:"name"`
	SKU_ID         string     `json:"sku_id"`
	Description    string     `json:"description"`
	ID             string     `json:"id"`
	Stickers       []*Sticker `json:"stickers"`
}

// User Object
// https://discord.com/developers/docs/resources/user#user-object
type User struct {
	PublicFlags   *BitFlag `json:"public_flag,omitempty"`
	PremiumType   *Flag    `json:"premium_type,omitempty"`
	Flags         *BitFlag `json:"flag,omitempty"`
	Avatar        *string  `json:"avatar"`
	Bot           *bool    `json:"bot,omitempty"`
	System        *bool    `json:"system,omitempty"`
	MFAEnabled    *bool    `json:"mfa_enabled,omitempty"`
	Banner        **string `json:"banner,omitempty"`
	AccentColor   **int    `json:"accent_color,omitempty"`
	Locale        *string  `json:"locale,omitempty"`
	Verified      *bool    `json:"verified,omitempty"`
	Email         **string `json:"email,omitempty"`
	Discriminator string   `json:"discriminator"`
	Username      string   `json:"username"`
	ID            string   `json:"id"`
}

// User Flags
// https://discord.com/developers/docs/resources/user#user-object-user-flags
const (
	FlagUserNONE                         BitFlag = 0
	FlagUserSTAFF                        BitFlag = 1 << 0
	FlagUserPARTNER                      BitFlag = 1 << 1
	FlagUserHYPESQUAD                    BitFlag = 1 << 2
	FlagUserBUG_HUNTER_LEVEL_1           BitFlag = 1 << 3
	FlagUserHYPESQUAD_ONLINE_HOUSE_ONE   BitFlag = 1 << 6
	FlagUserHYPESQUAD_ONLINE_HOUSE_TWO   BitFlag = 1 << 7
	FlagUserHYPESQUAD_ONLINE_HOUSE_THREE BitFlag = 1 << 8
	FlagUserPREMIUM_EARLY_SUPPORTER      BitFlag = 1 << 9
	FlagUserTEAM_PSEUDO_USER             BitFlag = 1 << 10
	FlagUserBUG_HUNTER_LEVEL_2           BitFlag = 1 << 14
	FlagUserVERIFIED_BOT                 BitFlag = 1 << 16
	FlagUserVERIFIED_DEVELOPER           BitFlag = 1 << 17
	FlagUserCERTIFIED_MODERATOR          BitFlag = 1 << 18
	FlagUserBOT_HTTP_INTERACTIONS        BitFlag = 1 << 19
	FlagUserACTIVE_DEVELOPER             BitFlag = 1 << 22
)

// Premium Types
// https://discord.com/developers/docs/resources/user#user-object-premium-types
const (
	FlagPremiumTypeNONE         Flag = 0
	FlagPremiumTypeNITROCLASSIC Flag = 1
	FlagPremiumTypeNITRO        Flag = 2
	FlagPremiumTypeNITROBASIC   Flag = 3
)

// User Connection Object
// https://discord.com/developers/docs/resources/user#connection-object-connection-structure
type Connection struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	Type         string         `json:"type"`
	Revoked      *bool          `json:"revoked,omitempty"`
	Integrations []*Integration `json:"integrations,omitempty"`
	Verified     bool           `json:"verified"`
	FriendSync   bool           `json:"friend_sync"`
	ShowActivity bool           `json:"show_activity"`
	TwoWayLink   bool           `json:"two_way_link"`
	Visibility   Flag           `json:"visibility"`
}

// Visibility Types
// https://discord.com/developers/docs/resources/user#connection-object-visibility-types
const (
	FlagVisibilityTypeNONE     Flag = 0
	FlagVisibilityTypeEVERYONE Flag = 1
)

// Voice State Object
// https://discord.com/developers/docs/resources/voice#voice-state-object-voice-state-structure
type VoiceState struct {
	GuildID                 *string      `json:"guild_id,omitempty"`
	ChannelID               *string      `json:"channel_id"`
	SelfStream              *bool        `json:"self_stream,omitempty"`
	Member                  *GuildMember `json:"member,omitempty"`
	RequestToSpeakTimestamp *time.Time   `json:"request_to_speak_timestamp"`
	SessionID               string       `json:"session_id"`
	UserID                  string       `json:"user_id"`
	Deaf                    bool         `json:"deaf"`
	SelfDeaf                bool         `json:"self_deaf"`
	SelfMute                bool         `json:"self_mute"`
	SelfVideo               bool         `json:"self_video"`
	Suppress                bool         `json:"suppress"`
	Mute                    bool         `json:"mute"`
}

// Voice Region Object
// https://discord.com/developers/docs/resources/voice#voice-region-object-voice-region-structure
type VoiceRegion struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Optimal    bool   `json:"optimal"`
	Deprecated bool   `json:"deprecated"`
	Custom     bool   `json:"custom"`
}

// Webhook Object
// https://discord.com/developers/docs/resources/webhook#webhook-object
type Webhook struct {
	Avatar        *string  `json:"avatar"`
	Token         *string  `json:"token,omitempty"`
	GuildID       **string `json:"guild_id,omitempty"`
	ChannelID     *string  `json:"channel_id"`
	User          *User    `json:"user,omitempty"`
	Name          *string  `json:"name"`
	ApplicationID *string  `json:"application_id"`
	SourceGuild   *Guild   `json:"source_guild,omitempty"`
	SourceChannel *Channel `json:"source_channel,omitempty"`
	URL           *string  `json:"url,omitempty"`
	ID            string   `json:"id"`
	Type          Flag     `json:"type"`
}

// Webhook Types
// https://discord.com/developers/docs/resources/webhook#webhook-object-webhook-types
const (
	FlagWebhookTypeINCOMING        Flag = 1
	FlagWebhookTypeCHANNELFOLLOWER Flag = 2
	FlagWebhookTypeAPPLICATION     Flag = 3
)

// Bitwise Permission Flags
// https://discord.com/developers/docs/topics/permissions#permissions-bitwise-permission-flags
const (
	FlagBitwisePermissionCREATE_INSTANT_INVITE      BitFlag = 1 << 0
	FlagBitwisePermissionKICK_MEMBERS               BitFlag = 1 << 1
	FlagBitwisePermissionBAN_MEMBERS                BitFlag = 1 << 2
	FlagBitwisePermissionADMINISTRATOR              BitFlag = 1 << 3
	FlagBitwisePermissionMANAGE_CHANNELS            BitFlag = 1 << 4
	FlagBitwisePermissionMANAGE_GUILD               BitFlag = 1 << 5
	FlagBitwisePermissionADD_REACTIONS              BitFlag = 1 << 6
	FlagBitwisePermissionVIEW_AUDIT_LOG             BitFlag = 1 << 7
	FlagBitwisePermissionPRIORITY_SPEAKER           BitFlag = 1 << 8
	FlagBitwisePermissionSTREAM                     BitFlag = 1 << 9
	FlagBitwisePermissionVIEW_CHANNEL               BitFlag = 1 << 10
	FlagBitwisePermissionSEND_MESSAGES              BitFlag = 1 << 11
	FlagBitwisePermissionSEND_TTS_MESSAGES          BitFlag = 1 << 12
	FlagBitwisePermissionMANAGE_MESSAGES            BitFlag = 1 << 13
	FlagBitwisePermissionEMBED_LINKS                BitFlag = 1 << 14
	FlagBitwisePermissionATTACH_FILES               BitFlag = 1 << 15
	FlagBitwisePermissionREAD_MESSAGE_HISTORY       BitFlag = 1 << 16
	FlagBitwisePermissionMENTION_EVERYONE           BitFlag = 1 << 17
	FlagBitwisePermissionUSE_EXTERNAL_EMOJIS        BitFlag = 1 << 18
	FlagBitwisePermissionVIEW_GUILD_INSIGHTS        BitFlag = 1 << 19
	FlagBitwisePermissionCONNECT                    BitFlag = 1 << 20
	FlagBitwisePermissionSPEAK                      BitFlag = 1 << 21
	FlagBitwisePermissionMUTE_MEMBERS               BitFlag = 1 << 22
	FlagBitwisePermissionDEAFEN_MEMBERS             BitFlag = 1 << 23
	FlagBitwisePermissionMOVE_MEMBERS               BitFlag = 1 << 24
	FlagBitwisePermissionUSE_VAD                    BitFlag = 1 << 25
	FlagBitwisePermissionCHANGE_NICKNAME            BitFlag = 1 << 26
	FlagBitwisePermissionMANAGE_NICKNAMES           BitFlag = 1 << 27
	FlagBitwisePermissionMANAGE_ROLES               BitFlag = 1 << 28
	FlagBitwisePermissionMANAGE_WEBHOOKS            BitFlag = 1 << 29
	FlagBitwisePermissionMANAGE_EMOJIS_AND_STICKERS BitFlag = 1 << 30
	FlagBitwisePermissionUSE_APPLICATION_COMMANDS   BitFlag = 1 << 31
	FlagBitwisePermissionREQUEST_TO_SPEAK           BitFlag = 1 << 32
	FlagBitwisePermissionMANAGE_EVENTS              BitFlag = 1 << 33
	FlagBitwisePermissionMANAGE_THREADS             BitFlag = 1 << 34
	FlagBitwisePermissionCREATE_PUBLIC_THREADS      BitFlag = 1 << 35
	FlagBitwisePermissionCREATE_PRIVATE_THREADS     BitFlag = 1 << 36
	FlagBitwisePermissionUSE_EXTERNAL_STICKERS      BitFlag = 1 << 37
	FlagBitwisePermissionSEND_MESSAGES_IN_THREADS   BitFlag = 1 << 38
	FlagBitwisePermissionUSE_EMBEDDED_ACTIVITIES    BitFlag = 1 << 39
	FlagBitwisePermissionMODERATE_MEMBERS           BitFlag = 1 << 40
)

// Permission Overwrite Types
const (
	FlagPermissionOverwriteTypeRole   Flag = 0
	FlagPermissionOverwriteTypeMember Flag = 1
)

// Role Object
// https://discord.com/developers/docs/topics/permissions#role-object
type Role struct {
	Icon         **string  `json:"icon,omitempty"`
	UnicodeEmoji **string  `json:"unicode_emoji,omitempty"`
	Tags         *RoleTags `json:"tags,omitempty"`
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Permissions  string    `json:"permissions"`
	Color        int       `json:"color"`
	Position     int       `json:"position"`
	Hoist        bool      `json:"hoist"`
	Managed      bool      `json:"managed"`
	Mentionable  bool      `json:"mentionable"`
}

// Role Tags Structure
// https://discord.com/developers/docs/topics/permissions#role-object-role-tags-structure
type RoleTags struct {
	BotID             *string `json:"bot_id,omitempty"`
	IntegrationID     *string `json:"integration_id,omitempty"`
	PremiumSubscriber *bool   `json:"premium_subscriber,omitempty"`
}

// Team Object
// https://discord.com/developers/docs/topics/teams#data-models-team-object
type Team struct {
	Icon        *string       `json:"icon"`
	Description *string       `json:"description"`
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	OwnerUserID string        `json:"owner_user_id"`
	Members     []*TeamMember `json:"members"`
}

// Team Member Object
// https://discord.com/developers/docs/topics/teams#data-models-team-member-object
type TeamMember struct {
	User            *User    `json:"user"`
	TeamID          string   `json:"team_id"`
	Permissions     []string `json:"permissions"`
	MembershipState Flag     `json:"membership_state"`
}

// Membership State Enum
// https://discord.com/developers/docs/topics/teams#data-models-membership-state-enum
const (
	FlagMembershipStateEnumINVITED  Flag = 1
	FlagMembershipStateEnumACCEPTED Flag = 2
)

// Client Status Object
// https://discord.com/developers/docs/topics/gateway-events#client-status-object
type ClientStatus struct {
	Desktop *string `json:"desktop,omitempty"`
	Mobile  *string `json:"mobile,omitempty"`
	Web     *string `json:"web,omitempty"`
}

// Activity Object
// https://discord.com/developers/docs/topics/gateway-events#activity-object
type Activity struct {
	Assets        *ActivityAssets     `json:"assets,omitempty"`
	Instance      *bool               `json:"instance,omitempty"`
	URL           **string            `json:"url,omitempty"`
	Secrets       *ActivitySecrets    `json:"secrets,omitempty"`
	Timestamps    *ActivityTimestamps `json:"timestamps,omitempty"`
	ApplicationID *string             `json:"application_id,omitempty"`
	Details       **string            `json:"details,omitempty"`
	State         **string            `json:"state,omitempty"`
	Emoji         **Emoji             `json:"emoji,omitempty"`
	Party         *ActivityParty      `json:"party,omitempty"`
	Name          string              `json:"name"`
	Buttons       []*Button           `json:"buttons,omitempty"`
	CreatedAt     int                 `json:"created_at"`
	Flags         BitFlag             `json:"flags,omitempty"`
	Type          Flag                `json:"type"`
}

// Activity Types
// https://discord.com/developers/docs/topics/gateway-events#activity-object-activity-types
const (
	FlagActivityTypePlaying   Flag = 0
	FlagActivityTypeStreaming Flag = 1
	FlagActivityTypeListening Flag = 2
	FlagActivityTypeWatching  Flag = 3
	FlagActivityTypeCustom    Flag = 4
	FlagActivityTypeCompeting Flag = 5
)

// Activity Timestamps Struct
// https://discord.com/developers/docs/topics/gateway-events#activity-object-activity-timestamps
type ActivityTimestamps struct {
	Start *int `json:"start,omitempty"`
	End   *int `json:"end,omitempty"`
}

// Activity Emoji
// https://discord.com/developers/docs/topics/gateway-events#activity-object-activity-emoji
type ActivityEmoji struct {
	ID       *string `json:"id,omitempty"`
	Animated *bool   `json:"animated,omitempty"`
	Name     string  `json:"name"`
}

// Activity Party Struct
// https://discord.com/developers/docs/topics/gateway-events#activity-object-activity-party
type ActivityParty struct {
	ID   *string `json:"id,omitempty"`
	Size *[2]int `json:"size,omitempty"`
}

// Activity Assets Struct
// https://discord.com/developers/docs/topics/gateway-events#activity-object-activity-assets
type ActivityAssets struct {
	LargeImage *string `json:"large_image,omitempty"`
	LargeText  *string `json:"large_text,omitempty"`
	SmallImage *string `json:"small_image,omitempty"`
	SmallText  *string `json:"small_text,omitempty"`
}

// Activity Asset Image
// https://discord.com/developers/docs/topics/gateway-events#activity-object-activity-asset-image
type ActivityAssetImage struct {
	ApplicationAsset string `json:"application_asset_id"`
	MediaProxyImage  string `json:"image_id"`
}

// Activity Secrets Struct
// https://discord.com/developers/docs/topics/gateway-events#activity-object-activity-secrets
type ActivitySecrets struct {
	Join     *string `json:"join,omitempty"`
	Spectate *string `json:"spectate,omitempty"`
	Match    *string `json:"match,omitempty"`
}

// Activity Flags
// https://discord.com/developers/docs/topics/gateway#activity-object-activity-flags
const (
	FlagActivityINSTANCE                    BitFlag = 1 << 0
	FlagActivityJOIN                        BitFlag = 1 << 1
	FlagActivitySPECTATE                    BitFlag = 1 << 2
	FlagActivityJOIN_REQUEST                BitFlag = 1 << 3
	FlagActivitySYNC                        BitFlag = 1 << 4
	FlagActivityPLAY                        BitFlag = 1 << 5
	FlagActivityPARTY_PRIVACY_FRIENDS       BitFlag = 1 << 6
	FlagActivityPARTY_PRIVACY_VOICE_CHANNEL BitFlag = 1 << 7
	FlagActivityEMBEDDED                    BitFlag = 1 << 8
)

// OAuth2 Scopes
// https://discord.com/developers/docs/topics/oauth2#shared-resources-oauth2-scopes
const (
	FlagOAuth2ScopeActivitiesRead                        = "activities.read"
	FlagOAuth2ScopeActivitiesWrite                       = "activities.write"
	FlagOAuth2ScopeApplicationsBuildsRead                = "applications.builds.read"
	FlagOAuth2ScopeApplicationsBuildsUpload              = "applications.builds.upload"
	FlagOAuth2ScopeApplicationsCommands                  = "applications.commands"
	FlagOAuth2ScopeApplicationsCommandsUpdate            = "applications.commands.update"
	FlagOAuth2ScopeApplicationsCommandsPermissionsUpdate = "applications.commands.permissions.update"
	FlagOAuth2ScopeApplicationsEntitlements              = "applications.entitlements"
	FlagOAuth2ScopeApplicationsStoreUpdate               = "applications.store.update"
	FlagOAuth2ScopeBot                                   = "bot"
	FlagOAuth2ScopeConnections                           = "connections"
	FlagOAuth2ScopeDM_channelsRead                       = "dm_channels.read"
	FlagOAuth2ScopeEmail                                 = "email"
	FlagOAuth2ScopeGDMJoin                               = "gdm.join"
	FlagOAuth2ScopeGuilds                                = "guilds"
	FlagOAuth2ScopeGuildsJoin                            = "guilds.join"
	FlagOAuth2ScopeGuildsMembersRead                     = "guilds.members.read"
	FlagOAuth2ScopeIdentify                              = "identify"
	FlagOAuth2ScopeMessagesRead                          = "messages.read"
	FlagOAuth2ScopeRelationshipsRead                     = "relationships.read"
	FlagOAuth2ScopeRPC                                   = "rpc"
	FlagOAuth2ScopeRPCActivitiesWrite                    = "rpc.activities.write"
	FlagOAuth2ScopeRPCNotificationsRead                  = "rpc.notifications.read"
	FlagOAuth2ScopeRPCVoiceRead                          = "rpc.voice.read"
	FlagOAuth2ScopeRPCVoiceWrite                         = "rpc.voice.write"
	FlagOAuth2ScopeVoice                                 = "voice"
	FlagOAuth2ScopeWebhookIncoming                       = "webhook.incoming"
)

// List Public Archived Threads Response Body
// https://discord.com/developers/docs/resources/channel#list-public-archived-threads-response-body
type ListPublicArchivedThreadsResponse struct {
	Threads []*Channel      `json:"threads"`
	Members []*ThreadMember `json:"members"`
	HasMore bool            `json:"has_more"`
}

// List Private Archived Threads Response Body
// https://discord.com/developers/docs/resources/channel#list-private-archived-threads-response-body
type ListPrivateArchivedThreadsResponse struct {
	Threads []*Channel      `json:"threads"`
	Members []*ThreadMember `json:"members"`
	HasMore bool            `json:"has_more"`
}

// List Joined Private Archived Threads Response Body
// https://discord.com/developers/docs/resources/channel#list-joined-private-archived-threads-response-body
type ListJoinedPrivateArchivedThreadsResponse struct {
	Threads []*Channel      `json:"threads"`
	Members []*ThreadMember `json:"members"`
	HasMore bool            `json:"has_more"`
}

// List Active Guild Threads Response Body
// https://discord.com/developers/docs/resources/guild#list-active-guild-threads-response-body
type ListActiveGuildThreadsResponse struct {
	Threads []*Channel      `json:"threads"`
	Members []*ThreadMember `json:"members"`
}

// Get Guild Prune Count Response Body
// https://discord.com/developers/docs/resources/guild#get-guild-prune-count
type GetGuildPruneCountResponse struct {
	Pruned int `json:"pruned"`
}

// Modify Guild MFA Level Response
// https://discord.com/developers/docs/resources/guild#modify-guild-mfa-level
type ModifyGuildMFALevelResponse struct {
	Level Flag `json:"level"`
}

// List Nitro Sticker Packs Response
// https://discord.com/developers/docs/resources/sticker#list-nitro-sticker-packs
type ListNitroStickerPacksResponse struct {
	StickerPacks []*StickerPack `json:"sticker_packs"`
}

// Current Authorization Information Response Structure
// https://discord.com/developers/docs/topics/oauth2#get-current-authorization-information
type CurrentAuthorizationInformationResponse struct {
	Application *Application `json:"application"`
	Expires     *time.Time   `json:"expires"`
	User        *User        `json:"user,omitempty"`
	Scopes      []int        `json:"scopes"`
}

// Get Gateway Response
// https://discord.com/developers/docs/topics/gateway#get-gateway-example-response
type GetGatewayResponse struct {
	URL string `json:"url,omitempty"`
}

// Get Gateway Bot Response
// https://discord.com/developers/docs/topics/gateway#get-gateway-example-response
type GetGatewayBotResponse struct {
	Shards            *int              `json:"shards"`
	URL               string            `json:"url"`
	SessionStartLimit SessionStartLimit `json:"session_start_limit"`
}

// Redirect URL
// https://discord.com/developers/docs/topics/oauth2#authorization-code-grant-redirect-url-example
type RedirectURL struct {
	Code  string `url:"code,omitempty"`
	State string `url:"state,omitempty"`

	// https://discord.com/developers/docs/topics/oauth2#advanced-bot-authorization
	GuildID     string  `url:"guild_id,omitempty"`
	Permissions BitFlag `url:"permissions,omitempty"`
}

// Access Token Response
// https://discord.com/developers/docs/topics/oauth2#authorization-code-grant-access-token-response
type AccessTokenResponse struct {
	AccessToken  string        `json:"access_token,omitempty"`
	TokenType    string        `json:"token_type,omitempty"`
	RefreshToken string        `json:"refresh_token,omitempty"`
	Scope        string        `json:"scope,omitempty"`
	ExpiresIn    time.Duration `json:"expires_in,omitempty"`
}

// Redirect URI
// https://discord.com/developers/docs/topics/oauth2#implicit-grant-redirect-url-example
type RedirectURI struct {
	AccessToken string        `url:"access_token,omitempty"`
	TokenType   string        `url:"token_type,omitempty"`
	Scope       string        `url:"scope,omitempty"`
	State       string        `url:"state,omitempty"`
	ExpiresIn   time.Duration `url:"expires_in,omitempty"`
}

// Client Credentials Access Token Response
// https://discord.com/developers/docs/topics/oauth2#client-credentials-grant-client-credentials-access-token-response
type ClientCredentialsAccessTokenResponse struct {
	AccessToken string        `json:"access_token,omitempty"`
	TokenType   string        `json:"token_type,omitempty"`
	Scope       string        `json:"scope,omitempty"`
	ExpiresIn   time.Duration `json:"expires_in,omitempty"`
}

// Webhook Token Response
// https://discord.com/developers/docs/topics/oauth2#webhooks-webhook-token-response-example
type WebhookTokenResponse struct {
	Webhook      *Webhook      `json:"webhook,omitempty"`
	TokenType    string        `json:"token_type,omitempty"`
	AccessToken  string        `json:"access_token,omitempty"`
	Scope        string        `json:"scope,omitempty"`
	RefreshToken string        `json:"refresh_token,omitempty"`
	ExpiresIn    time.Duration `json:"expires_in,omitempty"`
}

// Extended Bot Authorization Access Token Response
// https://discord.com/developers/docs/topics/oauth2#authorization-code-grant-access-token-response
type ExtendedBotAuthorizationAccessTokenResponse struct {
	Guild        *Guild        `json:"guild,omitempty"`
	TokenType    string        `json:"token_type,omitempty"`
	AccessToken  string        `json:"access_token,omitempty"`
	Scope        string        `json:"scope,omitempty"`
	RefreshToken string        `json:"refresh_token,omitempty"`
	ExpiresIn    time.Duration `json:"expires_in,omitempty"`
}

// Pointer returns a pointer to the given value.
func Pointer[T any](v T) *T {
	return &v
}

// Pointer2 returns a double pointer to the given value.
//
// set `null` to true in order to point the double pointer to a `nil` pointer.
func Pointer2[T any](v T, null ...bool) **T {
	if len(null) != 0 && null[0] {
		return new(*T)
	}

	pointer := Pointer(v)

	return &pointer
}

// IsValue returns whether the given pointer contains a value.
func IsValue[T any](p *T) bool {
	return p != nil
}

// IsValue2 returns returns whether the given double pointer contains a pointer.
func IsValue2[T any](dp **T) bool {
	return dp != nil
}

// PointerCheck returns whether the given double pointer contains a value.
//
// returns IsValueNothing, IsValueNull, or IsValueValid.
//
// 	IsValueNothing indicates that the field was not provided.
// 	IsValueNull indicates the the field was provided with a null value.
// 	IsValueValid indicates that the field is a valid value.
func PointerCheck[T any](dp **T) PointerIndicator {
	if dp != nil {
		if *dp != nil {
			return IsValueValid
		}

		return IsValueNull
	}

	return IsValueNothing
}

func (c ActionsRow) ComponentType() Flag {
	return FlagComponentTypeActionRow
}

func (c Button) ComponentType() Flag {
	return FlagComponentTypeButton
}

func (c SelectMenu) ComponentType() Flag {
	return FlagComponentTypeSelectMenu
}

func (c TextInput) ComponentType() Flag {
	return FlagComponentTypeTextInput
}

func (d ApplicationCommandData) InteractionDataType() Flag {
	return FlagInteractionTypeAPPLICATION_COMMAND
}

func (d MessageComponentData) InteractionDataType() Flag {
	return FlagInteractionTypeMESSAGE_COMPONENT
}

func (d ModalSubmitData) InteractionDataType() Flag {
	return FlagInteractionTypeMODAL_SUBMIT
}

func (d Messages) InteractionCallbackDataType() Flag {
	return FlagInteractionCallbackTypeCHANNEL_MESSAGE_WITH_SOURCE
}

func (d Autocomplete) InteractionCallbackDataType() Flag {
	return FlagInteractionCallbackTypeAPPLICATION_COMMAND_AUTOCOMPLETE_RESULT
}

func (d Modal) InteractionCallbackDataType() Flag {
	return FlagInteractionCallbackTypeMODAL
}

// Discord API Endpoints
const (
	EndpointBaseURL    = "https://discord.com/api/v" + VersionDiscordAPI + "/"
	CDNEndpointBaseURL = "https://cdn.discordapp.com/"
	achievements       = "achievements"
	active             = "active"
	appassets          = "app-assets"
	appicons           = "app-icons"
	applications       = "applications"
	archived           = "archived"
	auditlogs          = "audit-logs"
	authorize          = "authorize"
	automoderation     = "auto-moderation"
	avatars            = "avatars"
	banners            = "banners"
	bans               = "bans"
	bot                = "bot"
	bulkdelete         = "bulk-delete"
	callback           = "callback"
	channels           = "channels"
	commands           = "commands"
	connections        = "connections"
	crosspost          = "crosspost"
	discoverysplashes  = "discovery-splashes"
	embed              = "embed"
	emojis             = "emojis"
	followers          = "followers"
	gateway            = "gateway"
	github             = "github"
	guildevents        = "guild-events"
	guilds             = "guilds"
	icons              = "icons"
	integrations       = "integrations"
	interactions       = "interactions"
	invites            = "invites"
	me                 = "@me"
	member             = "member"
	members            = "members"
	messages           = "messages"
	mfa                = "mfa"
	nick               = "nick"
	oauth              = "oauth2"
	original           = "@original"
	permissions        = "permissions"
	pins               = "pins"
	preview            = "preview"
	private            = "private"
	prune              = "prune"
	public             = "public"
	reactions          = "reactions"
	recipients         = "recipients"
	regions            = "regions"
	revoke             = "revoke"
	roleicons          = "role-icons"
	roles              = "roles"
	rules              = "rules"
	scheduledevents    = "scheduled-events"
	search             = "search"
	slack              = "slack"
	slash              = "/"
	splashes           = "splashes"
	stageinstances     = "stage-instances"
	stickerpacks       = "sticker-packs"
	stickers           = "stickers"
	store              = "store"
	teamicons          = "team-icons"
	templates          = "templates"
	threadmembers      = "thread-members"
	threads            = "threads"
	token              = "token"
	typing             = "typing"
	users              = "users"
	vanityurl          = "vanity-url"
	voice              = "voice"
	voicestates        = "voice-states"
	webhooks           = "webhooks"
	welcomescreen      = "welcome-screen"
	widget             = "widget"
	widgetjson         = "widget.json"
	widgetpng          = "widget.png"
)

// EndpointGetGlobalApplicationCommands builds a query for an HTTP request.
func EndpointGetGlobalApplicationCommands(applicationid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + commands
}

// EndpointCreateGlobalApplicationCommand builds a query for an HTTP request.
func EndpointCreateGlobalApplicationCommand(applicationid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + commands
}

// EndpointGetGlobalApplicationCommand builds a query for an HTTP request.
func EndpointGetGlobalApplicationCommand(applicationid, commandid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + commands + slash + commandid
}

// EndpointEditGlobalApplicationCommand builds a query for an HTTP request.
func EndpointEditGlobalApplicationCommand(applicationid, commandid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + commands + slash + commandid
}

// EndpointDeleteGlobalApplicationCommand builds a query for an HTTP request.
func EndpointDeleteGlobalApplicationCommand(applicationid, commandid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + commands + slash + commandid
}

// EndpointBulkOverwriteGlobalApplicationCommands builds a query for an HTTP request.
func EndpointBulkOverwriteGlobalApplicationCommands(applicationid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + commands
}

// EndpointGetGuildApplicationCommands builds a query for an HTTP request.
func EndpointGetGuildApplicationCommands(applicationid, guildid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + guilds + slash + guildid + slash + commands
}

// EndpointCreateGuildApplicationCommand builds a query for an HTTP request.
func EndpointCreateGuildApplicationCommand(applicationid, guildid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + guilds + slash + guildid + slash + commands
}

// EndpointGetGuildApplicationCommand builds a query for an HTTP request.
func EndpointGetGuildApplicationCommand(applicationid, guildid, commandid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + guilds + slash + guildid + slash + commands + slash + commandid
}

// EndpointEditGuildApplicationCommand builds a query for an HTTP request.
func EndpointEditGuildApplicationCommand(applicationid, guildid, commandid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + guilds + slash + guildid + slash + commands + slash + commandid
}

// EndpointDeleteGuildApplicationCommand builds a query for an HTTP request.
func EndpointDeleteGuildApplicationCommand(applicationid, guildid, commandid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + guilds + slash + guildid + slash + commands + slash + commandid
}

// EndpointBulkOverwriteGuildApplicationCommands builds a query for an HTTP request.
func EndpointBulkOverwriteGuildApplicationCommands(applicationid, guildid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + guilds + slash + guildid + slash + commands
}

// EndpointGetGuildApplicationCommandPermissions builds a query for an HTTP request.
func EndpointGetGuildApplicationCommandPermissions(applicationid, guildid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + guilds + slash + guildid + slash + commands + slash + permissions
}

// EndpointGetApplicationCommandPermissions builds a query for an HTTP request.
func EndpointGetApplicationCommandPermissions(applicationid, guildid, commandid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + guilds + slash + guildid + slash + commands + slash + commandid + slash + permissions
}

// EndpointEditApplicationCommandPermissions builds a query for an HTTP request.
func EndpointEditApplicationCommandPermissions(applicationid, guildid, commandid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + guilds + slash + guildid + slash + commands + slash + commandid + slash + permissions
}

// EndpointBatchEditApplicationCommandPermissions builds a query for an HTTP request.
func EndpointBatchEditApplicationCommandPermissions(applicationid, guildid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + guilds + slash + guildid + slash + commands + slash + permissions
}

// EndpointCreateInteractionResponse builds a query for an HTTP request.
func EndpointCreateInteractionResponse(interactionid, interactiontoken string) string {
	return EndpointBaseURL + interactions + slash + interactionid + slash + interactiontoken + slash + callback
}

// EndpointGetOriginalInteractionResponse builds a query for an HTTP request.
func EndpointGetOriginalInteractionResponse(applicationid, interactiontoken string) string {
	return EndpointBaseURL + webhooks + slash + applicationid + slash + interactiontoken + slash + messages + slash + original
}

// EndpointEditOriginalInteractionResponse builds a query for an HTTP request.
func EndpointEditOriginalInteractionResponse(applicationid, interactiontoken string) string {
	return EndpointBaseURL + webhooks + slash + applicationid + slash + interactiontoken + slash + messages + slash + original
}

// EndpointDeleteOriginalInteractionResponse builds a query for an HTTP request.
func EndpointDeleteOriginalInteractionResponse(applicationid, interactiontoken string) string {
	return EndpointBaseURL + webhooks + slash + applicationid + slash + interactiontoken + slash + messages + slash + original
}

// EndpointCreateFollowupMessage builds a query for an HTTP request.
func EndpointCreateFollowupMessage(applicationid, interactiontoken string) string {
	return EndpointBaseURL + webhooks + slash + applicationid + slash + interactiontoken
}

// EndpointGetFollowupMessage builds a query for an HTTP request.
func EndpointGetFollowupMessage(applicationid, interactiontoken, messageid string) string {
	return EndpointBaseURL + webhooks + slash + applicationid + slash + interactiontoken + slash + messages + slash + messageid
}

// EndpointEditFollowupMessage builds a query for an HTTP request.
func EndpointEditFollowupMessage(applicationid, interactiontoken, messageid string) string {
	return EndpointBaseURL + webhooks + slash + applicationid + slash + interactiontoken + slash + messages + slash + messageid
}

// EndpointDeleteFollowupMessage builds a query for an HTTP request.
func EndpointDeleteFollowupMessage(applicationid, interactiontoken, messageid string) string {
	return EndpointBaseURL + webhooks + slash + applicationid + slash + interactiontoken + slash + messages + slash + messageid
}

// EndpointGetGuildAuditLog builds a query for an HTTP request.
func EndpointGetGuildAuditLog(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + auditlogs
}

// EndpointListAutoModerationRulesForGuild builds a query for an HTTP request.
func EndpointListAutoModerationRulesForGuild(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + automoderation + slash + rules
}

// EndpointGetAutoModerationRule builds a query for an HTTP request.
func EndpointGetAutoModerationRule(guildid, automoderationruleid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + automoderation + slash + rules + slash + automoderationruleid
}

// EndpointCreateAutoModerationRule builds a query for an HTTP request.
func EndpointCreateAutoModerationRule(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + automoderation + slash + rules
}

// EndpointModifyAutoModerationRule builds a query for an HTTP request.
func EndpointModifyAutoModerationRule(guildid, automoderationruleid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + automoderation + slash + rules + slash + automoderationruleid
}

// EndpointDeleteAutoModerationRule builds a query for an HTTP request.
func EndpointDeleteAutoModerationRule(guildid, automoderationruleid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + automoderation + slash + rules + slash + automoderationruleid
}

// EndpointGetChannel builds a query for an HTTP request.
func EndpointGetChannel(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid
}

// EndpointModifyChannel builds a query for an HTTP request.
func EndpointModifyChannel(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid
}

// EndpointDeleteCloseChannel builds a query for an HTTP request.
func EndpointDeleteCloseChannel(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid
}

// EndpointGetChannelMessages builds a query for an HTTP request.
func EndpointGetChannelMessages(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + messages
}

// EndpointGetChannelMessage builds a query for an HTTP request.
func EndpointGetChannelMessage(channelid, messageid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + messages + slash + messageid
}

// EndpointCreateMessage builds a query for an HTTP request.
func EndpointCreateMessage(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + messages
}

// EndpointCrosspostMessage builds a query for an HTTP request.
func EndpointCrosspostMessage(channelid, messageid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + messages + slash + messageid + slash + crosspost
}

// EndpointCreateReaction builds a query for an HTTP request.
func EndpointCreateReaction(channelid, messageid, emoji string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + messages + slash + messageid + slash + reactions + slash + emoji + slash + me
}

// EndpointDeleteOwnReaction builds a query for an HTTP request.
func EndpointDeleteOwnReaction(channelid, messageid, emoji string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + messages + slash + messageid + slash + reactions + slash + emoji + slash + me
}

// EndpointDeleteUserReaction builds a query for an HTTP request.
func EndpointDeleteUserReaction(channelid, messageid, emoji, userid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + messages + slash + messageid + slash + reactions + slash + emoji + slash + userid
}

// EndpointGetReactions builds a query for an HTTP request.
func EndpointGetReactions(channelid, messageid, emoji string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + messages + slash + messageid + slash + reactions + slash + emoji
}

// EndpointDeleteAllReactions builds a query for an HTTP request.
func EndpointDeleteAllReactions(channelid, messageid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + messages + slash + messageid + slash + reactions
}

// EndpointDeleteAllReactionsforEmoji builds a query for an HTTP request.
func EndpointDeleteAllReactionsforEmoji(channelid, messageid, emoji string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + messages + slash + messageid + slash + reactions + slash + emoji
}

// EndpointEditMessage builds a query for an HTTP request.
func EndpointEditMessage(channelid, messageid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + messages + slash + messageid
}

// EndpointDeleteMessage builds a query for an HTTP request.
func EndpointDeleteMessage(channelid, messageid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + messages + slash + messageid
}

// EndpointBulkDeleteMessages builds a query for an HTTP request.
func EndpointBulkDeleteMessages(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + messages + slash + bulkdelete
}

// EndpointEditChannelPermissions builds a query for an HTTP request.
func EndpointEditChannelPermissions(channelid, overwriteid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + permissions + slash + overwriteid
}

// EndpointGetChannelInvites builds a query for an HTTP request.
func EndpointGetChannelInvites(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + invites
}

// EndpointCreateChannelInvite builds a query for an HTTP request.
func EndpointCreateChannelInvite(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + invites
}

// EndpointDeleteChannelPermission builds a query for an HTTP request.
func EndpointDeleteChannelPermission(channelid, overwriteid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + permissions + slash + overwriteid
}

// EndpointFollowAnnouncementChannel builds a query for an HTTP request.
func EndpointFollowAnnouncementChannel(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + followers
}

// EndpointTriggerTypingIndicator builds a query for an HTTP request.
func EndpointTriggerTypingIndicator(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + typing
}

// EndpointGetPinnedMessages builds a query for an HTTP request.
func EndpointGetPinnedMessages(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + pins
}

// EndpointPinMessage builds a query for an HTTP request.
func EndpointPinMessage(channelid, messageid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + pins + slash + messageid
}

// EndpointUnpinMessage builds a query for an HTTP request.
func EndpointUnpinMessage(channelid, messageid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + pins + slash + messageid
}

// EndpointGroupDMAddRecipient builds a query for an HTTP request.
func EndpointGroupDMAddRecipient(channelid, userid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + recipients + slash + userid
}

// EndpointGroupDMRemoveRecipient builds a query for an HTTP request.
func EndpointGroupDMRemoveRecipient(channelid, userid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + recipients + slash + userid
}

// EndpointStartThreadfromMessage builds a query for an HTTP request.
func EndpointStartThreadfromMessage(channelid, messageid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + messages + slash + messageid + slash + threads
}

// EndpointStartThreadwithoutMessage builds a query for an HTTP request.
func EndpointStartThreadwithoutMessage(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + threads
}

// EndpointStartThreadinForumChannel builds a query for an HTTP request.
func EndpointStartThreadinForumChannel(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + threads
}

// EndpointJoinThread builds a query for an HTTP request.
func EndpointJoinThread(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + threadmembers + slash + me
}

// EndpointAddThreadMember builds a query for an HTTP request.
func EndpointAddThreadMember(channelid, userid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + threadmembers + slash + userid
}

// EndpointLeaveThread builds a query for an HTTP request.
func EndpointLeaveThread(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + threadmembers + slash + me
}

// EndpointRemoveThreadMember builds a query for an HTTP request.
func EndpointRemoveThreadMember(channelid, userid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + threadmembers + slash + userid
}

// EndpointGetThreadMember builds a query for an HTTP request.
func EndpointGetThreadMember(channelid, userid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + threadmembers + slash + userid
}

// EndpointListThreadMembers builds a query for an HTTP request.
func EndpointListThreadMembers(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + threadmembers
}

// EndpointListPublicArchivedThreads builds a query for an HTTP request.
func EndpointListPublicArchivedThreads(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + threads + slash + archived + slash + public
}

// EndpointListPrivateArchivedThreads builds a query for an HTTP request.
func EndpointListPrivateArchivedThreads(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + threads + slash + archived + slash + private
}

// EndpointListJoinedPrivateArchivedThreads builds a query for an HTTP request.
func EndpointListJoinedPrivateArchivedThreads(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + users + slash + me + slash + threads + slash + archived + slash + private
}

// EndpointListGuildEmojis builds a query for an HTTP request.
func EndpointListGuildEmojis(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + emojis
}

// EndpointGetGuildEmoji builds a query for an HTTP request.
func EndpointGetGuildEmoji(guildid, emojiid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + emojis + slash + emojiid
}

// EndpointCreateGuildEmoji builds a query for an HTTP request.
func EndpointCreateGuildEmoji(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + emojis
}

// EndpointModifyGuildEmoji builds a query for an HTTP request.
func EndpointModifyGuildEmoji(guildid, emojiid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + emojis + slash + emojiid
}

// EndpointDeleteGuildEmoji builds a query for an HTTP request.
func EndpointDeleteGuildEmoji(guildid, emojiid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + emojis + slash + emojiid
}

// EndpointCreateGuild builds a query for an HTTP request.
func EndpointCreateGuild() string {
	return EndpointBaseURL + guilds
}

// EndpointGetGuild builds a query for an HTTP request.
func EndpointGetGuild(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid
}

// EndpointGetGuildPreview builds a query for an HTTP request.
func EndpointGetGuildPreview(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + preview
}

// EndpointModifyGuild builds a query for an HTTP request.
func EndpointModifyGuild(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid
}

// EndpointDeleteGuild builds a query for an HTTP request.
func EndpointDeleteGuild(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid
}

// EndpointGetGuildChannels builds a query for an HTTP request.
func EndpointGetGuildChannels(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + channels
}

// EndpointCreateGuildChannel builds a query for an HTTP request.
func EndpointCreateGuildChannel(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + channels
}

// EndpointModifyGuildChannelPositions builds a query for an HTTP request.
func EndpointModifyGuildChannelPositions(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + channels
}

// EndpointListActiveGuildThreads builds a query for an HTTP request.
func EndpointListActiveGuildThreads(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + threads + slash + active
}

// EndpointGetGuildMember builds a query for an HTTP request.
func EndpointGetGuildMember(guildid, userid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + members + slash + userid
}

// EndpointListGuildMembers builds a query for an HTTP request.
func EndpointListGuildMembers(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + members
}

// EndpointSearchGuildMembers builds a query for an HTTP request.
func EndpointSearchGuildMembers(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + members + slash + search
}

// EndpointAddGuildMember builds a query for an HTTP request.
func EndpointAddGuildMember(guildid, userid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + members + slash + userid
}

// EndpointModifyGuildMember builds a query for an HTTP request.
func EndpointModifyGuildMember(guildid, userid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + members + slash + userid
}

// EndpointModifyCurrentMember builds a query for an HTTP request.
func EndpointModifyCurrentMember(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + members + slash + me
}

// EndpointModifyCurrentUserNick builds a query for an HTTP request.
func EndpointModifyCurrentUserNick(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + members + slash + me + slash + nick
}

// EndpointAddGuildMemberRole builds a query for an HTTP request.
func EndpointAddGuildMemberRole(guildid, userid, roleid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + members + slash + userid + slash + roles + slash + roleid
}

// EndpointRemoveGuildMemberRole builds a query for an HTTP request.
func EndpointRemoveGuildMemberRole(guildid, userid, roleid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + members + slash + userid + slash + roles + slash + roleid
}

// EndpointRemoveGuildMember builds a query for an HTTP request.
func EndpointRemoveGuildMember(guildid, userid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + members + slash + userid
}

// EndpointGetGuildBans builds a query for an HTTP request.
func EndpointGetGuildBans(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + bans
}

// EndpointGetGuildBan builds a query for an HTTP request.
func EndpointGetGuildBan(guildid, userid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + bans + slash + userid
}

// EndpointCreateGuildBan builds a query for an HTTP request.
func EndpointCreateGuildBan(guildid, userid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + bans + slash + userid
}

// EndpointRemoveGuildBan builds a query for an HTTP request.
func EndpointRemoveGuildBan(guildid, userid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + bans + slash + userid
}

// EndpointGetGuildRoles builds a query for an HTTP request.
func EndpointGetGuildRoles(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + roles
}

// EndpointCreateGuildRole builds a query for an HTTP request.
func EndpointCreateGuildRole(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + roles
}

// EndpointModifyGuildRolePositions builds a query for an HTTP request.
func EndpointModifyGuildRolePositions(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + roles
}

// EndpointModifyGuildRole builds a query for an HTTP request.
func EndpointModifyGuildRole(guildid, roleid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + roles + slash + roleid
}

// EndpointModifyGuildMFALevel builds a query for an HTTP request.
func EndpointModifyGuildMFALevel(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + mfa
}

// EndpointDeleteGuildRole builds a query for an HTTP request.
func EndpointDeleteGuildRole(guildid, roleid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + roles + slash + roleid
}

// EndpointGetGuildPruneCount builds a query for an HTTP request.
func EndpointGetGuildPruneCount(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + prune
}

// EndpointBeginGuildPrune builds a query for an HTTP request.
func EndpointBeginGuildPrune(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + prune
}

// EndpointGetGuildVoiceRegions builds a query for an HTTP request.
func EndpointGetGuildVoiceRegions(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + regions
}

// EndpointGetGuildInvites builds a query for an HTTP request.
func EndpointGetGuildInvites(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + invites
}

// EndpointGetGuildIntegrations builds a query for an HTTP request.
func EndpointGetGuildIntegrations(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + integrations
}

// EndpointDeleteGuildIntegration builds a query for an HTTP request.
func EndpointDeleteGuildIntegration(guildid, integrationid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + integrations + slash + integrationid
}

// EndpointGetGuildWidgetSettings builds a query for an HTTP request.
func EndpointGetGuildWidgetSettings(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + widget
}

// EndpointModifyGuildWidget builds a query for an HTTP request.
func EndpointModifyGuildWidget(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + widget
}

// EndpointGetGuildWidget builds a query for an HTTP request.
func EndpointGetGuildWidget(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + widgetjson
}

// EndpointGetGuildVanityURL builds a query for an HTTP request.
func EndpointGetGuildVanityURL(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + vanityurl
}

// EndpointGetGuildWidgetImage builds a query for an HTTP request.
func EndpointGetGuildWidgetImage(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + widgetpng
}

// EndpointGetGuildWelcomeScreen builds a query for an HTTP request.
func EndpointGetGuildWelcomeScreen(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + welcomescreen
}

// EndpointModifyGuildWelcomeScreen builds a query for an HTTP request.
func EndpointModifyGuildWelcomeScreen(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + welcomescreen
}

// EndpointModifyCurrentUserVoiceState builds a query for an HTTP request.
func EndpointModifyCurrentUserVoiceState(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + voicestates + slash + me
}

// EndpointModifyUserVoiceState builds a query for an HTTP request.
func EndpointModifyUserVoiceState(guildid, userid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + voicestates + slash + userid
}

// EndpointListScheduledEventsforGuild builds a query for an HTTP request.
func EndpointListScheduledEventsforGuild(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + scheduledevents
}

// EndpointCreateGuildScheduledEvent builds a query for an HTTP request.
func EndpointCreateGuildScheduledEvent(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + scheduledevents
}

// EndpointGetGuildScheduledEvent builds a query for an HTTP request.
func EndpointGetGuildScheduledEvent(guildid, guildscheduledeventid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + scheduledevents + slash + guildscheduledeventid
}

// EndpointModifyGuildScheduledEvent builds a query for an HTTP request.
func EndpointModifyGuildScheduledEvent(guildid, guildscheduledeventid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + scheduledevents + slash + guildscheduledeventid
}

// EndpointDeleteGuildScheduledEvent builds a query for an HTTP request.
func EndpointDeleteGuildScheduledEvent(guildid, guildscheduledeventid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + scheduledevents + slash + guildscheduledeventid
}

// EndpointGetGuildScheduledEventUsers builds a query for an HTTP request.
func EndpointGetGuildScheduledEventUsers(guildid, guildscheduledeventid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + scheduledevents + slash + guildscheduledeventid + slash + users
}

// EndpointGetGuildTemplate builds a query for an HTTP request.
func EndpointGetGuildTemplate(templatecode string) string {
	return EndpointBaseURL + guilds + slash + templates + slash + templatecode
}

// EndpointCreateGuildfromGuildTemplate builds a query for an HTTP request.
func EndpointCreateGuildfromGuildTemplate(templatecode string) string {
	return EndpointBaseURL + guilds + slash + templates + slash + templatecode
}

// EndpointGetGuildTemplates builds a query for an HTTP request.
func EndpointGetGuildTemplates(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + templates
}

// EndpointCreateGuildTemplate builds a query for an HTTP request.
func EndpointCreateGuildTemplate(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + templates
}

// EndpointSyncGuildTemplate builds a query for an HTTP request.
func EndpointSyncGuildTemplate(guildid, templatecode string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + templates + slash + templatecode
}

// EndpointModifyGuildTemplate builds a query for an HTTP request.
func EndpointModifyGuildTemplate(guildid, templatecode string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + templates + slash + templatecode
}

// EndpointDeleteGuildTemplate builds a query for an HTTP request.
func EndpointDeleteGuildTemplate(guildid, templatecode string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + templates + slash + templatecode
}

// EndpointGetInvite builds a query for an HTTP request.
func EndpointGetInvite(invitecode string) string {
	return EndpointBaseURL + invites + slash + invitecode
}

// EndpointDeleteInvite builds a query for an HTTP request.
func EndpointDeleteInvite(invitecode string) string {
	return EndpointBaseURL + invites + slash + invitecode
}

// EndpointCreateStageInstance builds a query for an HTTP request.
func EndpointCreateStageInstance() string {
	return EndpointBaseURL + stageinstances
}

// EndpointGetStageInstance builds a query for an HTTP request.
func EndpointGetStageInstance(channelid string) string {
	return EndpointBaseURL + stageinstances + slash + channelid
}

// EndpointModifyStageInstance builds a query for an HTTP request.
func EndpointModifyStageInstance(channelid string) string {
	return EndpointBaseURL + stageinstances + slash + channelid
}

// EndpointDeleteStageInstance builds a query for an HTTP request.
func EndpointDeleteStageInstance(channelid string) string {
	return EndpointBaseURL + stageinstances + slash + channelid
}

// EndpointGetSticker builds a query for an HTTP request.
func EndpointGetSticker(stickerid string) string {
	return EndpointBaseURL + stickers + slash + stickerid
}

// EndpointListNitroStickerPacks builds a query for an HTTP request.
func EndpointListNitroStickerPacks() string {
	return EndpointBaseURL + stickerpacks
}

// EndpointListGuildStickers builds a query for an HTTP request.
func EndpointListGuildStickers(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + stickers
}

// EndpointGetGuildSticker builds a query for an HTTP request.
func EndpointGetGuildSticker(guildid, stickerid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + stickers + slash + stickerid
}

// EndpointCreateGuildSticker builds a query for an HTTP request.
func EndpointCreateGuildSticker(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + stickers
}

// EndpointModifyGuildSticker builds a query for an HTTP request.
func EndpointModifyGuildSticker(guildid, stickerid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + stickers + slash + stickerid
}

// EndpointDeleteGuildSticker builds a query for an HTTP request.
func EndpointDeleteGuildSticker(guildid, stickerid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + stickers + slash + stickerid
}

// EndpointGetCurrentUser builds a query for an HTTP request.
func EndpointGetCurrentUser() string {
	return EndpointBaseURL + users + slash + me
}

// EndpointGetUser builds a query for an HTTP request.
func EndpointGetUser(userid string) string {
	return EndpointBaseURL + users + slash + userid
}

// EndpointModifyCurrentUser builds a query for an HTTP request.
func EndpointModifyCurrentUser() string {
	return EndpointBaseURL + users + slash + me
}

// EndpointGetCurrentUserGuilds builds a query for an HTTP request.
func EndpointGetCurrentUserGuilds() string {
	return EndpointBaseURL + users + slash + me + slash + guilds
}

// EndpointGetCurrentUserGuildMember builds a query for an HTTP request.
func EndpointGetCurrentUserGuildMember(guildid string) string {
	return EndpointBaseURL + users + slash + me + slash + guilds + slash + guildid + slash + member
}

// EndpointLeaveGuild builds a query for an HTTP request.
func EndpointLeaveGuild(guildid string) string {
	return EndpointBaseURL + users + slash + me + slash + guilds + slash + guildid
}

// EndpointCreateDM builds a query for an HTTP request.
func EndpointCreateDM() string {
	return EndpointBaseURL + users + slash + me + slash + channels
}

// EndpointCreateGroupDM builds a query for an HTTP request.
func EndpointCreateGroupDM() string {
	return EndpointBaseURL + users + slash + me + slash + channels
}

// EndpointGetUserConnections builds a query for an HTTP request.
func EndpointGetUserConnections() string {
	return EndpointBaseURL + users + slash + me + slash + connections
}

// EndpointListVoiceRegions builds a query for an HTTP request.
func EndpointListVoiceRegions() string {
	return EndpointBaseURL + voice + slash + regions
}

// EndpointCreateWebhook builds a query for an HTTP request.
func EndpointCreateWebhook(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + webhooks
}

// EndpointGetChannelWebhooks builds a query for an HTTP request.
func EndpointGetChannelWebhooks(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + webhooks
}

// EndpointGetGuildWebhooks builds a query for an HTTP request.
func EndpointGetGuildWebhooks(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + webhooks
}

// EndpointGetWebhook builds a query for an HTTP request.
func EndpointGetWebhook(webhookid string) string {
	return EndpointBaseURL + webhooks + slash + webhookid
}

// EndpointGetWebhookwithToken builds a query for an HTTP request.
func EndpointGetWebhookwithToken(webhookid, webhooktoken string) string {
	return EndpointBaseURL + webhooks + slash + webhookid + slash + webhooktoken
}

// EndpointModifyWebhook builds a query for an HTTP request.
func EndpointModifyWebhook(webhookid string) string {
	return EndpointBaseURL + webhooks + slash + webhookid
}

// EndpointModifyWebhookwithToken builds a query for an HTTP request.
func EndpointModifyWebhookwithToken(webhookid, webhooktoken string) string {
	return EndpointBaseURL + webhooks + slash + webhookid + slash + webhooktoken
}

// EndpointDeleteWebhook builds a query for an HTTP request.
func EndpointDeleteWebhook(webhookid string) string {
	return EndpointBaseURL + webhooks + slash + webhookid
}

// EndpointDeleteWebhookwithToken builds a query for an HTTP request.
func EndpointDeleteWebhookwithToken(webhookid, webhooktoken string) string {
	return EndpointBaseURL + webhooks + slash + webhookid + slash + webhooktoken
}

// EndpointExecuteWebhook builds a query for an HTTP request.
func EndpointExecuteWebhook(webhookid, webhooktoken string) string {
	return EndpointBaseURL + webhooks + slash + webhookid + slash + webhooktoken
}

// EndpointExecuteSlackCompatibleWebhook builds a query for an HTTP request.
func EndpointExecuteSlackCompatibleWebhook(webhookid, webhooktoken string) string {
	return EndpointBaseURL + webhooks + slash + webhookid + slash + webhooktoken + slash + slack
}

// EndpointExecuteGitHubCompatibleWebhook builds a query for an HTTP request.
func EndpointExecuteGitHubCompatibleWebhook(webhookid, webhooktoken string) string {
	return EndpointBaseURL + webhooks + slash + webhookid + slash + webhooktoken + slash + github
}

// EndpointGetWebhookMessage builds a query for an HTTP request.
func EndpointGetWebhookMessage(webhookid, webhooktoken, messageid string) string {
	return EndpointBaseURL + webhooks + slash + webhookid + slash + webhooktoken + slash + messages + slash + messageid
}

// EndpointEditWebhookMessage builds a query for an HTTP request.
func EndpointEditWebhookMessage(webhookid, webhooktoken, messageid string) string {
	return EndpointBaseURL + webhooks + slash + webhookid + slash + webhooktoken + slash + messages + slash + messageid
}

// EndpointDeleteWebhookMessage builds a query for an HTTP request.
func EndpointDeleteWebhookMessage(webhookid, webhooktoken, messageid string) string {
	return EndpointBaseURL + webhooks + slash + webhookid + slash + webhooktoken + slash + messages + slash + messageid
}

// EndpointGetGateway builds a query for an HTTP request.
func EndpointGetGateway() string {
	return EndpointBaseURL + gateway
}

// EndpointGetGatewayBot builds a query for an HTTP request.
func EndpointGetGatewayBot() string {
	return EndpointBaseURL + gateway + slash + bot
}

// EndpointAuthorizationURL builds a query for an HTTP request.
func EndpointAuthorizationURL() string {
	return EndpointBaseURL + oauth + slash + authorize
}

// EndpointTokenURL builds a query for an HTTP request.
func EndpointTokenURL() string {
	return EndpointBaseURL + oauth + slash + token
}

// EndpointTokenRevocationURL builds a query for an HTTP request.
func EndpointTokenRevocationURL() string {
	return EndpointBaseURL + oauth + slash + token + slash + revoke
}

// EndpointGetCurrentBotApplicationInformation builds a query for an HTTP request.
func EndpointGetCurrentBotApplicationInformation() string {
	return EndpointBaseURL + oauth + slash + applications + slash + me
}

// EndpointGetCurrentAuthorizationInformation builds a query for an HTTP request.
func EndpointGetCurrentAuthorizationInformation() string {
	return EndpointBaseURL + oauth + slash + me
}

// CDNEndpointCustomEmoji builds a query for an HTTP request.
func CDNEndpointCustomEmoji(emojiid string) string {
	return CDNEndpointBaseURL + emojis + slash + emojiid
}

// CDNEndpointGuildIcon builds a query for an HTTP request.
func CDNEndpointGuildIcon(guildid, guildicon string) string {
	return CDNEndpointBaseURL + icons + slash + guildid + slash + guildicon
}

// CDNEndpointGuildSplash builds a query for an HTTP request.
func CDNEndpointGuildSplash(guildid, guildsplash string) string {
	return CDNEndpointBaseURL + splashes + slash + guildid + slash + guildsplash
}

// CDNEndpointGuildDiscoverySplash builds a query for an HTTP request.
func CDNEndpointGuildDiscoverySplash(guildid, guilddiscoverysplash string) string {
	return CDNEndpointBaseURL + discoverysplashes + slash + guildid + slash + guilddiscoverysplash
}

// CDNEndpointGuildBanner builds a query for an HTTP request.
func CDNEndpointGuildBanner(guildid, guildbanner string) string {
	return CDNEndpointBaseURL + banners + slash + guildid + slash + guildbanner
}

// CDNEndpointUserBanner builds a query for an HTTP request.
func CDNEndpointUserBanner(userid, userbanner string) string {
	return CDNEndpointBaseURL + banners + slash + userid + slash + userbanner
}

// CDNEndpointDefaultUserAvatar builds a query for an HTTP request.
func CDNEndpointDefaultUserAvatar(userdiscriminator string) string {
	return CDNEndpointBaseURL + embed + slash + avatars + slash + userdiscriminator
}

// CDNEndpointUserAvatar builds a query for an HTTP request.
func CDNEndpointUserAvatar(userid, useravatar string) string {
	return CDNEndpointBaseURL + avatars + slash + userid + slash + useravatar
}

// CDNEndpointGuildMemberAvatar builds a query for an HTTP request.
func CDNEndpointGuildMemberAvatar(guildid, userid, memberavatar string) string {
	return CDNEndpointBaseURL + guilds + slash + guildid + slash + users + slash + userid + slash + avatars + slash + memberavatar
}

// CDNEndpointApplicationIcon builds a query for an HTTP request.
func CDNEndpointApplicationIcon(applicationid, icon string) string {
	return CDNEndpointBaseURL + appicons + slash + applicationid + slash + icon
}

// CDNEndpointApplicationCover builds a query for an HTTP request.
func CDNEndpointApplicationCover(applicationid, coverimage string) string {
	return CDNEndpointBaseURL + appicons + slash + applicationid + slash + coverimage
}

// CDNEndpointApplicationAsset builds a query for an HTTP request.
func CDNEndpointApplicationAsset(applicationid, assetid string) string {
	return CDNEndpointBaseURL + appassets + slash + applicationid + slash + assetid
}

// CDNEndpointAchievementIcon builds a query for an HTTP request.
func CDNEndpointAchievementIcon(applicationid, achievementid, iconhash string) string {
	return CDNEndpointBaseURL + appassets + slash + applicationid + slash + achievements + slash + achievementid + slash + icons + slash + iconhash
}

// CDNEndpointStickerPackBanner builds a query for an HTTP request.
func CDNEndpointStickerPackBanner(applicationid, stickerpackbannerassetid string) string {
	return CDNEndpointBaseURL + appassets + slash + applicationid + slash + store + slash + stickerpackbannerassetid
}

// CDNEndpointTeamIcon builds a query for an HTTP request.
func CDNEndpointTeamIcon(teamid, teamicon string) string {
	return CDNEndpointBaseURL + teamicons + slash + teamid + slash + teamicon
}

// CDNEndpointSticker builds a query for an HTTP request.
func CDNEndpointSticker(stickerid string) string {
	return CDNEndpointBaseURL + stickers + slash + stickerid
}

// CDNEndpointRoleIcon builds a query for an HTTP request.
func CDNEndpointRoleIcon(roleid, roleicon string) string {
	return CDNEndpointBaseURL + roleicons + slash + roleid + slash + roleicon
}

// CDNEndpointGuildScheduledEventCover builds a query for an HTTP request.
func CDNEndpointGuildScheduledEventCover(scheduledeventid, scheduledeventcoverimage string) string {
	return CDNEndpointBaseURL + guildevents + slash + scheduledeventid + slash + scheduledeventcoverimage
}

// CDNEndpointGuildMemberBanner builds a query for an HTTP request.
func CDNEndpointGuildMemberBanner(guildid, userid, memberbanner string) string {
	return CDNEndpointBaseURL + guilds + slash + guildid + slash + users + slash + userid + slash + banners + slash + memberbanner
}

var (
	EndpointModifyChannelGroupDM = EndpointModifyChannel
	EndpointModifyChannelGuild   = EndpointModifyChannel
	EndpointModifyChannelThread  = EndpointModifyChannel
)

// Send Request Error Messages.
const (
	errSendMarshal = "marshalling an HTTP body: %w"
	errUnmarshal   = "error unmarshalling into %T: %w"
	errRateLimit   = "error converting the HTTP rate limit header %q: %w"
)

// ErrorRequest represents an HTTP Request error that occurs when an attempt to send a request fails.
type ErrorRequest struct {
	Err           error
	ClientID      string
	CorrelationID string
	RouteID       string
	ResourceID    string
	Endpoint      string
}

func (e ErrorRequest) Error() string {
	return fmt.Errorf("REQUEST ERROR: client %q: x %q: route: %q: resource: %q: endpoint: %q: error: %w",
		e.ClientID, e.CorrelationID, e.RouteID, e.ResourceID, e.Endpoint, e.Err).Error()
}

// Status Code Error Messages.
const (
	errStatusCodeKnown   = "Status Code %d: %v"
	errStatusCodeUnknown = "Status Code %d: Unknown status code error from Discord"
)

// StatusCodeError handles a Discord API HTTP Status Code and returns the relevant error message.
func StatusCodeError(status int) error {
	if msg, ok := HTTPResponseCodes[status]; ok {
		return fmt.Errorf(errStatusCodeKnown, status, msg)
	}

	return fmt.Errorf(errStatusCodeUnknown, status)
}

// JSON Error Code Messages.
const (
	errJSONErrorCodeKnown = "JSON Error Code %d: %v"
	errJSONErrorUnknown   = "JSON Error Code %d: Unknown JSON Error Code from Discord"
)

// JSONCodeError handles a Discord API JSON Error Code and returns the relevant error message.
func JSONCodeError(status int) error {
	if msg, ok := JSONErrorCodes[status]; ok {
		return fmt.Errorf(errJSONErrorCodeKnown, status, msg)
	}

	return fmt.Errorf(errJSONErrorUnknown, status)
}

// Event Handler Error Messages.
const (
	errHandleNotRemoved   = "event handler was not added"
	errRemoveInvalidIndex = "event handler cannot be removed since there is no event handler at index %d"
)

// ErrorEventHandler represents an Event Handler error that occurs when an attempt to
// add or remove an event handler fails.
type ErrorEventHandler struct {
	Err      error
	ClientID string
	Event    string
}

func (e ErrorEventHandler) Error() string {
	return fmt.Errorf("EVENT HANDLER ERROR: client %q: event %q: error: %w",
		e.ClientID, e.Event, e.Err).Error()
}

const (
	ErrorEventActionUnmarshal = "unmarshalling"
	ErrorEventActionMarshal   = "marshalling"
	ErrorEventActionRead      = "reading"
	ErrorEventActionWrite     = "writing"
)

// ErrorEvent represents a WebSocket error that occurs when an attempt to {action} an event fails.
type ErrorEvent struct {
	// ClientID represents the Application ID of the event handler caller.
	ClientID string

	// Event represents the name of the event involved in this error.
	Event string

	// Err represents the error that occurred while performing the action.
	Err error

	// Action represents the action that prompted the error.
	//
	// ErrorEventAction's can be one of four values:
	// ErrorEventActionUnmarshal: an error occurred while unmarshalling the Event from a JSON.
	// ErrorEventActionMarshal:   an error occurred while marshalling the Event to a JSON.
	// ErrorEventActionRead:      an error occurred while reading the Event from a Websocket Connection.
	// ErrorEventActionWrite:     an error occurred while writing the Event to a Websocket Connection.
	Action string
}

func (e ErrorEvent) Error() string {
	return fmt.Errorf("EVENT ERROR: client %q: event %q: action: %q, error: %w",
		e.ClientID, e.Event, e.Action, e.Err).Error()
}

// ErrorSession represents a WebSocket Session error that occurs during an active session.
type ErrorSession struct {
	Err       error
	SessionID string
}

func (e ErrorSession) Error() string {
	return fmt.Errorf("SESSION ERROR: session %q: error: %w", e.SessionID, e.Err).Error()
}

const (
	ErrConnectionSession = "Discord Gateway"
)

// ErrorDisconnect represents a disconnection error that occurs when
// an attempt to gracefully disconnect from a connection fails.
type ErrorDisconnect struct {
	Action     error
	Err        error
	Connection string
}

func (e ErrorDisconnect) Error() string {
	return fmt.Errorf("error disconnecting from %q\n"+
		"\tDisconnect(): %v\n"+
		"\treason: %w\n",
		e.Connection, e.Err, e.Action).Error()
}

// Handlers represents a bot's event handlers.
type Handlers struct {
	Hello                               []func(*Hello)
	Ready                               []func(*Ready)
	Resumed                             []func(*Resumed)
	Reconnect                           []func(*Reconnect)
	InvalidSession                      []func(*InvalidSession)
	ApplicationCommandPermissionsUpdate []func(*ApplicationCommandPermissionsUpdate)
	AutoModerationRuleCreate            []func(*AutoModerationRuleCreate)
	AutoModerationRuleUpdate            []func(*AutoModerationRuleUpdate)
	AutoModerationRuleDelete            []func(*AutoModerationRuleDelete)
	AutoModerationActionExecution       []func(*AutoModerationActionExecution)
	InteractionCreate                   []func(*InteractionCreate)
	VoiceServerUpdate                   []func(*VoiceServerUpdate)
	GuildMembersChunk                   []func(*GuildMembersChunk)
	UserUpdate                          []func(*UserUpdate)
	ChannelCreate                       []func(*ChannelCreate)
	ChannelUpdate                       []func(*ChannelUpdate)
	ChannelDelete                       []func(*ChannelDelete)
	ChannelPinsUpdate                   []func(*ChannelPinsUpdate)
	ThreadCreate                        []func(*ThreadCreate)
	ThreadUpdate                        []func(*ThreadUpdate)
	ThreadDelete                        []func(*ThreadDelete)
	ThreadListSync                      []func(*ThreadListSync)
	ThreadMemberUpdate                  []func(*ThreadMemberUpdate)
	ThreadMembersUpdate                 []func(*ThreadMembersUpdate)
	GuildCreate                         []func(*GuildCreate)
	GuildUpdate                         []func(*GuildUpdate)
	GuildDelete                         []func(*GuildDelete)
	GuildBanAdd                         []func(*GuildBanAdd)
	GuildBanRemove                      []func(*GuildBanRemove)
	GuildEmojisUpdate                   []func(*GuildEmojisUpdate)
	GuildStickersUpdate                 []func(*GuildStickersUpdate)
	GuildIntegrationsUpdate             []func(*GuildIntegrationsUpdate)
	GuildMemberAdd                      []func(*GuildMemberAdd)
	GuildMemberRemove                   []func(*GuildMemberRemove)
	GuildMemberUpdate                   []func(*GuildMemberUpdate)
	GuildRoleCreate                     []func(*GuildRoleCreate)
	GuildRoleUpdate                     []func(*GuildRoleUpdate)
	GuildRoleDelete                     []func(*GuildRoleDelete)
	GuildScheduledEventCreate           []func(*GuildScheduledEventCreate)
	GuildScheduledEventUpdate           []func(*GuildScheduledEventUpdate)
	GuildScheduledEventDelete           []func(*GuildScheduledEventDelete)
	GuildScheduledEventUserAdd          []func(*GuildScheduledEventUserAdd)
	GuildScheduledEventUserRemove       []func(*GuildScheduledEventUserRemove)
	IntegrationCreate                   []func(*IntegrationCreate)
	IntegrationUpdate                   []func(*IntegrationUpdate)
	IntegrationDelete                   []func(*IntegrationDelete)
	InviteCreate                        []func(*InviteCreate)
	InviteDelete                        []func(*InviteDelete)
	MessageCreate                       []func(*MessageCreate)
	MessageUpdate                       []func(*MessageUpdate)
	MessageDelete                       []func(*MessageDelete)
	MessageDeleteBulk                   []func(*MessageDeleteBulk)
	MessageReactionAdd                  []func(*MessageReactionAdd)
	MessageReactionRemove               []func(*MessageReactionRemove)
	MessageReactionRemoveAll            []func(*MessageReactionRemoveAll)
	MessageReactionRemoveEmoji          []func(*MessageReactionRemoveEmoji)
	PresenceUpdate                      []func(*PresenceUpdate)
	StageInstanceCreate                 []func(*StageInstanceCreate)
	StageInstanceDelete                 []func(*StageInstanceDelete)
	StageInstanceUpdate                 []func(*StageInstanceUpdate)
	TypingStart                         []func(*TypingStart)
	VoiceStateUpdate                    []func(*VoiceStateUpdate)
	WebhooksUpdate                      []func(*WebhooksUpdate)
	mu                                  sync.RWMutex
}

// Handle adds an event handler for the given event to the bot.
func (bot *Client) Handle(eventname string, function interface{}) error {
	bot.Handlers.mu.Lock()
	defer bot.Handlers.mu.Unlock()

	switch eventname {
	case FlagGatewayEventNameHello:
		if f, ok := function.(func(*Hello)); ok {
			bot.Handlers.Hello = append(bot.Handlers.Hello, f)
			LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("added event handler")
			return nil
		}

	case FlagGatewayEventNameReady:
		if f, ok := function.(func(*Ready)); ok {
			bot.Handlers.Ready = append(bot.Handlers.Ready, f)
			LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("added event handler")
			return nil
		}

	case FlagGatewayEventNameResumed:
		if f, ok := function.(func(*Resumed)); ok {
			bot.Handlers.Resumed = append(bot.Handlers.Resumed, f)
			LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("added event handler")
			return nil
		}

	case FlagGatewayEventNameReconnect:
		if f, ok := function.(func(*Reconnect)); ok {
			bot.Handlers.Reconnect = append(bot.Handlers.Reconnect, f)
			LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("added event handler")
			return nil
		}

	case FlagGatewayEventNameInvalidSession:
		if f, ok := function.(func(*InvalidSession)); ok {
			bot.Handlers.InvalidSession = append(bot.Handlers.InvalidSession, f)
			LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("added event handler")
			return nil
		}

	case FlagGatewayEventNameApplicationCommandPermissionsUpdate:
		if f, ok := function.(func(*ApplicationCommandPermissionsUpdate)); ok {
			bot.Handlers.ApplicationCommandPermissionsUpdate = append(bot.Handlers.ApplicationCommandPermissionsUpdate, f)
			LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("added event handler")
			return nil
		}

	case FlagGatewayEventNameAutoModerationRuleCreate:
		if !bot.Config.Gateway.IntentSet[FlagIntentAUTO_MODERATION_CONFIGURATION] {
			bot.Config.Gateway.IntentSet[FlagIntentAUTO_MODERATION_CONFIGURATION] = true
			bot.Config.Gateway.Intents |= FlagIntentAUTO_MODERATION_CONFIGURATION
		}

		if f, ok := function.(func(*AutoModerationRuleCreate)); ok {
			bot.Handlers.AutoModerationRuleCreate = append(bot.Handlers.AutoModerationRuleCreate, f)
			LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("added event handler")
			return nil
		}

	case FlagGatewayEventNameAutoModerationRuleUpdate:
		if !bot.Config.Gateway.IntentSet[FlagIntentAUTO_MODERATION_CONFIGURATION] {
			bot.Config.Gateway.IntentSet[FlagIntentAUTO_MODERATION_CONFIGURATION] = true
			bot.Config.Gateway.Intents |= FlagIntentAUTO_MODERATION_CONFIGURATION
		}

		if f, ok := function.(func(*AutoModerationRuleUpdate)); ok {
			bot.Handlers.AutoModerationRuleUpdate = append(bot.Handlers.AutoModerationRuleUpdate, f)
			LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("added event handler")
			return nil
		}

	case FlagGatewayEventNameAutoModerationRuleDelete:
		if !bot.Config.Gateway.IntentSet[FlagIntentAUTO_MODERATION_CONFIGURATION] {
			bot.Config.Gateway.IntentSet[FlagIntentAUTO_MODERATION_CONFIGURATION] = true
			bot.Config.Gateway.Intents |= FlagIntentAUTO_MODERATION_CONFIGURATION
		}

		if f, ok := function.(func(*AutoModerationRuleDelete)); ok {
			bot.Handlers.AutoModerationRuleDelete = append(bot.Handlers.AutoModerationRuleDelete, f)
			LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("added event handler")
			return nil
		}

	case FlagGatewayEventNameAutoModerationActionExecution:
		if !bot.Config.Gateway.IntentSet[FlagIntentAUTO_MODERATION_EXECUTION] {
			bot.Config.Gateway.IntentSet[FlagIntentAUTO_MODERATION_EXECUTION] = true
			bot.Config.Gateway.Intents |= FlagIntentAUTO_MODERATION_EXECUTION
		}

		if f, ok := function.(func(*AutoModerationActionExecution)); ok {
			bot.Handlers.AutoModerationActionExecution = append(bot.Handlers.AutoModerationActionExecution, f)
			LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("added event handler")
			return nil
		}

	case FlagGatewayEventNameInteractionCreate:
		if f, ok := function.(func(*InteractionCreate)); ok {
			bot.Handlers.InteractionCreate = append(bot.Handlers.InteractionCreate, f)
			LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("added event handler")
			return nil
		}

	case FlagGatewayEventNameVoiceServerUpdate:
		if f, ok := function.(func(*VoiceServerUpdate)); ok {
			bot.Handlers.VoiceServerUpdate = append(bot.Handlers.VoiceServerUpdate, f)
			LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("added event handler")
			return nil
		}

	case FlagGatewayEventNameGuildMembersChunk:
		if f, ok := function.(func(*GuildMembersChunk)); ok {
			bot.Handlers.GuildMembersChunk = append(bot.Handlers.GuildMembersChunk, f)
			LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("added event handler")
			return nil
		}

	case FlagGatewayEventNameUserUpdate:
		if f, ok := function.(func(*UserUpdate)); ok {
			bot.Handlers.UserUpdate = append(bot.Handlers.UserUpdate, f)
			LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("added event handler")
			return nil
		}

	case FlagGatewayEventNameChannelCreate:
		if !bot.Config.Gateway.IntentSet[FlagIntentGUILDS] {
			bot.Config.Gateway.IntentSet[FlagIntentGUILDS] = true
			bot.Config.Gateway.Intents |= FlagIntentGUILDS
		}

		if f, ok := function.(func(*ChannelCreate)); ok {
			bot.Handlers.ChannelCreate = append(bot.Handlers.ChannelCreate, f)
			LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("added event handler")
			return nil
		}

	case FlagGatewayEventNameChannelUpdate:
		if !bot.Config.Gateway.IntentSet[FlagIntentGUILDS] {
			bot.Config.Gateway.IntentSet[FlagIntentGUILDS] = true
			bot.Config.Gateway.Intents |= FlagIntentGUILDS
		}

		if f, ok := function.(func(*ChannelUpdate)); ok {
			bot.Handlers.ChannelUpdate = append(bot.Handlers.ChannelUpdate, f)
			LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("added event handler")
			return nil
		}

	case FlagGatewayEventNameChannelDelete:
		if !bot.Config.Gateway.IntentSet[FlagIntentGUILDS] {
			bot.Config.Gateway.IntentSet[FlagIntentGUILDS] = true
			bot.Config.Gateway.Intents |= FlagIntentGUILDS
		}

		if f, ok := function.(func(*ChannelDelete)); ok {
			bot.Handlers.ChannelDelete = append(bot.Handlers.ChannelDelete, f)
			LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("added event handler")
			return nil
		}

	case FlagGatewayEventNameChannelPinsUpdate:
		if !bot.Config.Gateway.IntentSet[FlagIntentDIRECT_MESSAGES] {
			bot.Config.Gateway.IntentSet[FlagIntentDIRECT_MESSAGES] = true
			bot.Config.Gateway.Intents |= FlagIntentDIRECT_MESSAGES
		}

		if !bot.Config.Gateway.IntentSet[FlagIntentGUILDS] {
			bot.Config.Gateway.IntentSet[FlagIntentGUILDS] = true
			bot.Config.Gateway.Intents |= FlagIntentGUILDS
		}

		if f, ok := function.(func(*ChannelPinsUpdate)); ok {
			bot.Handlers.ChannelPinsUpdate = append(bot.Handlers.ChannelPinsUpdate, f)
			LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("added event handler")
			return nil
		}

	case FlagGatewayEventNameThreadCreate:
		if !bot.Config.Gateway.IntentSet[FlagIntentGUILDS] {
			bot.Config.Gateway.IntentSet[FlagIntentGUILDS] = true
			bot.Config.Gateway.Intents |= FlagIntentGUILDS
		}

		if f, ok := function.(func(*ThreadCreate)); ok {
			bot.Handlers.ThreadCreate = append(bot.Handlers.ThreadCreate, f)
			LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("added event handler")
			return nil
		}

	case FlagGatewayEventNameThreadUpdate:
		if !bot.Config.Gateway.IntentSet[FlagIntentGUILDS] {
			bot.Config.Gateway.IntentSet[FlagIntentGUILDS] = true
			bot.Config.Gateway.Intents |= FlagIntentGUILDS
		}

		if f, ok := function.(func(*ThreadUpdate)); ok {
			bot.Handlers.ThreadUpdate = append(bot.Handlers.ThreadUpdate, f)
			LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("added event handler")
			return nil
		}

	case FlagGatewayEventNameThreadDelete:
		if !bot.Config.Gateway.IntentSet[FlagIntentGUILDS] {
			bot.Config.Gateway.IntentSet[FlagIntentGUILDS] = true
			bot.Config.Gateway.Intents |= FlagIntentGUILDS
		}

		if f, ok := function.(func(*ThreadDelete)); ok {
			bot.Handlers.ThreadDelete = append(bot.Handlers.ThreadDelete, f)
			LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("added event handler")
			return nil
		}

	case FlagGatewayEventNameThreadListSync:
		if !bot.Config.Gateway.IntentSet[FlagIntentGUILDS] {
			bot.Config.Gateway.IntentSet[FlagIntentGUILDS] = true
			bot.Config.Gateway.Intents |= FlagIntentGUILDS
		}

		if f, ok := function.(func(*ThreadListSync)); ok {
			bot.Handlers.ThreadListSync = append(bot.Handlers.ThreadListSync, f)
			LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("added event handler")
			return nil
		}

	case FlagGatewayEventNameThreadMemberUpdate:
		if !bot.Config.Gateway.IntentSet[FlagIntentGUILDS] {
			bot.Config.Gateway.IntentSet[FlagIntentGUILDS] = true
			bot.Config.Gateway.Intents |= FlagIntentGUILDS
		}

		if f, ok := function.(func(*ThreadMemberUpdate)); ok {
			bot.Handlers.ThreadMemberUpdate = append(bot.Handlers.ThreadMemberUpdate, f)
			LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("added event handler")
			return nil
		}

	case FlagGatewayEventNameThreadMembersUpdate:
		if !bot.Config.Gateway.IntentSet[FlagIntentGUILDS] {
			bot.Config.Gateway.IntentSet[FlagIntentGUILDS] = true
			bot.Config.Gateway.Intents |= FlagIntentGUILDS
		}

		if !bot.Config.Gateway.IntentSet[FlagIntentGUILD_MEMBERS] {
			bot.Config.Gateway.IntentSet[FlagIntentGUILD_MEMBERS] = true
			bot.Config.Gateway.Intents |= FlagIntentGUILD_MEMBERS
		}

		if f, ok := function.(func(*ThreadMembersUpdate)); ok {
			bot.Handlers.ThreadMembersUpdate = append(bot.Handlers.ThreadMembersUpdate, f)
			LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("added event handler")
			return nil
		}

	case FlagGatewayEventNameGuildCreate:
		if !bot.Config.Gateway.IntentSet[FlagIntentGUILDS] {
			bot.Config.Gateway.IntentSet[FlagIntentGUILDS] = true
			bot.Config.Gateway.Intents |= FlagIntentGUILDS
		}

		if f, ok := function.(func(*GuildCreate)); ok {
			bot.Handlers.GuildCreate = append(bot.Handlers.GuildCreate, f)
			LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("added event handler")
			return nil
		}

	case FlagGatewayEventNameGuildUpdate:
		if !bot.Config.Gateway.IntentSet[FlagIntentGUILDS] {
			bot.Config.Gateway.IntentSet[FlagIntentGUILDS] = true
			bot.Config.Gateway.Intents |= FlagIntentGUILDS
		}

		if f, ok := function.(func(*GuildUpdate)); ok {
			bot.Handlers.GuildUpdate = append(bot.Handlers.GuildUpdate, f)
			LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("added event handler")
			return nil
		}

	case FlagGatewayEventNameGuildDelete:
		if !bot.Config.Gateway.IntentSet[FlagIntentGUILDS] {
			bot.Config.Gateway.IntentSet[FlagIntentGUILDS] = true
			bot.Config.Gateway.Intents |= FlagIntentGUILDS
		}

		if f, ok := function.(func(*GuildDelete)); ok {
			bot.Handlers.GuildDelete = append(bot.Handlers.GuildDelete, f)
			LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("added event handler")
			return nil
		}

	case FlagGatewayEventNameGuildBanAdd:
		if !bot.Config.Gateway.IntentSet[FlagIntentGUILD_BANS] {
			bot.Config.Gateway.IntentSet[FlagIntentGUILD_BANS] = true
			bot.Config.Gateway.Intents |= FlagIntentGUILD_BANS
		}

		if f, ok := function.(func(*GuildBanAdd)); ok {
			bot.Handlers.GuildBanAdd = append(bot.Handlers.GuildBanAdd, f)
			LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("added event handler")
			return nil
		}

	case FlagGatewayEventNameGuildBanRemove:
		if !bot.Config.Gateway.IntentSet[FlagIntentGUILD_BANS] {
			bot.Config.Gateway.IntentSet[FlagIntentGUILD_BANS] = true
			bot.Config.Gateway.Intents |= FlagIntentGUILD_BANS
		}

		if f, ok := function.(func(*GuildBanRemove)); ok {
			bot.Handlers.GuildBanRemove = append(bot.Handlers.GuildBanRemove, f)
			LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("added event handler")
			return nil
		}

	case FlagGatewayEventNameGuildEmojisUpdate:
		if !bot.Config.Gateway.IntentSet[FlagIntentGUILD_EMOJIS_AND_STICKERS] {
			bot.Config.Gateway.IntentSet[FlagIntentGUILD_EMOJIS_AND_STICKERS] = true
			bot.Config.Gateway.Intents |= FlagIntentGUILD_EMOJIS_AND_STICKERS
		}

		if f, ok := function.(func(*GuildEmojisUpdate)); ok {
			bot.Handlers.GuildEmojisUpdate = append(bot.Handlers.GuildEmojisUpdate, f)
			LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("added event handler")
			return nil
		}

	case FlagGatewayEventNameGuildStickersUpdate:
		if !bot.Config.Gateway.IntentSet[FlagIntentGUILD_EMOJIS_AND_STICKERS] {
			bot.Config.Gateway.IntentSet[FlagIntentGUILD_EMOJIS_AND_STICKERS] = true
			bot.Config.Gateway.Intents |= FlagIntentGUILD_EMOJIS_AND_STICKERS
		}

		if f, ok := function.(func(*GuildStickersUpdate)); ok {
			bot.Handlers.GuildStickersUpdate = append(bot.Handlers.GuildStickersUpdate, f)
			LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("added event handler")
			return nil
		}

	case FlagGatewayEventNameGuildIntegrationsUpdate:
		if !bot.Config.Gateway.IntentSet[FlagIntentGUILD_INTEGRATIONS] {
			bot.Config.Gateway.IntentSet[FlagIntentGUILD_INTEGRATIONS] = true
			bot.Config.Gateway.Intents |= FlagIntentGUILD_INTEGRATIONS
		}

		if f, ok := function.(func(*GuildIntegrationsUpdate)); ok {
			bot.Handlers.GuildIntegrationsUpdate = append(bot.Handlers.GuildIntegrationsUpdate, f)
			LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("added event handler")
			return nil
		}

	case FlagGatewayEventNameGuildMemberAdd:
		if !bot.Config.Gateway.IntentSet[FlagIntentGUILD_MEMBERS] {
			bot.Config.Gateway.IntentSet[FlagIntentGUILD_MEMBERS] = true
			bot.Config.Gateway.Intents |= FlagIntentGUILD_MEMBERS
		}

		if f, ok := function.(func(*GuildMemberAdd)); ok {
			bot.Handlers.GuildMemberAdd = append(bot.Handlers.GuildMemberAdd, f)
			LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("added event handler")
			return nil
		}

	case FlagGatewayEventNameGuildMemberRemove:
		if !bot.Config.Gateway.IntentSet[FlagIntentGUILD_MEMBERS] {
			bot.Config.Gateway.IntentSet[FlagIntentGUILD_MEMBERS] = true
			bot.Config.Gateway.Intents |= FlagIntentGUILD_MEMBERS
		}

		if f, ok := function.(func(*GuildMemberRemove)); ok {
			bot.Handlers.GuildMemberRemove = append(bot.Handlers.GuildMemberRemove, f)
			LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("added event handler")
			return nil
		}

	case FlagGatewayEventNameGuildMemberUpdate:
		if !bot.Config.Gateway.IntentSet[FlagIntentGUILD_MEMBERS] {
			bot.Config.Gateway.IntentSet[FlagIntentGUILD_MEMBERS] = true
			bot.Config.Gateway.Intents |= FlagIntentGUILD_MEMBERS
		}

		if f, ok := function.(func(*GuildMemberUpdate)); ok {
			bot.Handlers.GuildMemberUpdate = append(bot.Handlers.GuildMemberUpdate, f)
			LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("added event handler")
			return nil
		}

	case FlagGatewayEventNameGuildRoleCreate:
		if !bot.Config.Gateway.IntentSet[FlagIntentGUILDS] {
			bot.Config.Gateway.IntentSet[FlagIntentGUILDS] = true
			bot.Config.Gateway.Intents |= FlagIntentGUILDS
		}

		if f, ok := function.(func(*GuildRoleCreate)); ok {
			bot.Handlers.GuildRoleCreate = append(bot.Handlers.GuildRoleCreate, f)
			LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("added event handler")
			return nil
		}

	case FlagGatewayEventNameGuildRoleUpdate:
		if !bot.Config.Gateway.IntentSet[FlagIntentGUILDS] {
			bot.Config.Gateway.IntentSet[FlagIntentGUILDS] = true
			bot.Config.Gateway.Intents |= FlagIntentGUILDS
		}

		if f, ok := function.(func(*GuildRoleUpdate)); ok {
			bot.Handlers.GuildRoleUpdate = append(bot.Handlers.GuildRoleUpdate, f)
			LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("added event handler")
			return nil
		}

	case FlagGatewayEventNameGuildRoleDelete:
		if !bot.Config.Gateway.IntentSet[FlagIntentGUILDS] {
			bot.Config.Gateway.IntentSet[FlagIntentGUILDS] = true
			bot.Config.Gateway.Intents |= FlagIntentGUILDS
		}

		if f, ok := function.(func(*GuildRoleDelete)); ok {
			bot.Handlers.GuildRoleDelete = append(bot.Handlers.GuildRoleDelete, f)
			LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("added event handler")
			return nil
		}

	case FlagGatewayEventNameGuildScheduledEventCreate:
		if !bot.Config.Gateway.IntentSet[FlagIntentGUILD_SCHEDULED_EVENTS] {
			bot.Config.Gateway.IntentSet[FlagIntentGUILD_SCHEDULED_EVENTS] = true
			bot.Config.Gateway.Intents |= FlagIntentGUILD_SCHEDULED_EVENTS
		}

		if f, ok := function.(func(*GuildScheduledEventCreate)); ok {
			bot.Handlers.GuildScheduledEventCreate = append(bot.Handlers.GuildScheduledEventCreate, f)
			LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("added event handler")
			return nil
		}

	case FlagGatewayEventNameGuildScheduledEventUpdate:
		if !bot.Config.Gateway.IntentSet[FlagIntentGUILD_SCHEDULED_EVENTS] {
			bot.Config.Gateway.IntentSet[FlagIntentGUILD_SCHEDULED_EVENTS] = true
			bot.Config.Gateway.Intents |= FlagIntentGUILD_SCHEDULED_EVENTS
		}

		if f, ok := function.(func(*GuildScheduledEventUpdate)); ok {
			bot.Handlers.GuildScheduledEventUpdate = append(bot.Handlers.GuildScheduledEventUpdate, f)
			LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("added event handler")
			return nil
		}

	case FlagGatewayEventNameGuildScheduledEventDelete:
		if !bot.Config.Gateway.IntentSet[FlagIntentGUILD_SCHEDULED_EVENTS] {
			bot.Config.Gateway.IntentSet[FlagIntentGUILD_SCHEDULED_EVENTS] = true
			bot.Config.Gateway.Intents |= FlagIntentGUILD_SCHEDULED_EVENTS
		}

		if f, ok := function.(func(*GuildScheduledEventDelete)); ok {
			bot.Handlers.GuildScheduledEventDelete = append(bot.Handlers.GuildScheduledEventDelete, f)
			LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("added event handler")
			return nil
		}

	case FlagGatewayEventNameGuildScheduledEventUserAdd:
		if !bot.Config.Gateway.IntentSet[FlagIntentGUILD_SCHEDULED_EVENTS] {
			bot.Config.Gateway.IntentSet[FlagIntentGUILD_SCHEDULED_EVENTS] = true
			bot.Config.Gateway.Intents |= FlagIntentGUILD_SCHEDULED_EVENTS
		}

		if f, ok := function.(func(*GuildScheduledEventUserAdd)); ok {
			bot.Handlers.GuildScheduledEventUserAdd = append(bot.Handlers.GuildScheduledEventUserAdd, f)
			LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("added event handler")
			return nil
		}

	case FlagGatewayEventNameGuildScheduledEventUserRemove:
		if !bot.Config.Gateway.IntentSet[FlagIntentGUILD_SCHEDULED_EVENTS] {
			bot.Config.Gateway.IntentSet[FlagIntentGUILD_SCHEDULED_EVENTS] = true
			bot.Config.Gateway.Intents |= FlagIntentGUILD_SCHEDULED_EVENTS
		}

		if f, ok := function.(func(*GuildScheduledEventUserRemove)); ok {
			bot.Handlers.GuildScheduledEventUserRemove = append(bot.Handlers.GuildScheduledEventUserRemove, f)
			LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("added event handler")
			return nil
		}

	case FlagGatewayEventNameIntegrationCreate:
		if !bot.Config.Gateway.IntentSet[FlagIntentGUILD_INTEGRATIONS] {
			bot.Config.Gateway.IntentSet[FlagIntentGUILD_INTEGRATIONS] = true
			bot.Config.Gateway.Intents |= FlagIntentGUILD_INTEGRATIONS
		}

		if f, ok := function.(func(*IntegrationCreate)); ok {
			bot.Handlers.IntegrationCreate = append(bot.Handlers.IntegrationCreate, f)
			LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("added event handler")
			return nil
		}

	case FlagGatewayEventNameIntegrationUpdate:
		if !bot.Config.Gateway.IntentSet[FlagIntentGUILD_INTEGRATIONS] {
			bot.Config.Gateway.IntentSet[FlagIntentGUILD_INTEGRATIONS] = true
			bot.Config.Gateway.Intents |= FlagIntentGUILD_INTEGRATIONS
		}

		if f, ok := function.(func(*IntegrationUpdate)); ok {
			bot.Handlers.IntegrationUpdate = append(bot.Handlers.IntegrationUpdate, f)
			LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("added event handler")
			return nil
		}

	case FlagGatewayEventNameIntegrationDelete:
		if !bot.Config.Gateway.IntentSet[FlagIntentGUILD_INTEGRATIONS] {
			bot.Config.Gateway.IntentSet[FlagIntentGUILD_INTEGRATIONS] = true
			bot.Config.Gateway.Intents |= FlagIntentGUILD_INTEGRATIONS
		}

		if f, ok := function.(func(*IntegrationDelete)); ok {
			bot.Handlers.IntegrationDelete = append(bot.Handlers.IntegrationDelete, f)
			LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("added event handler")
			return nil
		}

	case FlagGatewayEventNameInviteCreate:
		if !bot.Config.Gateway.IntentSet[FlagIntentGUILD_INVITES] {
			bot.Config.Gateway.IntentSet[FlagIntentGUILD_INVITES] = true
			bot.Config.Gateway.Intents |= FlagIntentGUILD_INVITES
		}

		if f, ok := function.(func(*InviteCreate)); ok {
			bot.Handlers.InviteCreate = append(bot.Handlers.InviteCreate, f)
			LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("added event handler")
			return nil
		}

	case FlagGatewayEventNameInviteDelete:
		if !bot.Config.Gateway.IntentSet[FlagIntentGUILD_INVITES] {
			bot.Config.Gateway.IntentSet[FlagIntentGUILD_INVITES] = true
			bot.Config.Gateway.Intents |= FlagIntentGUILD_INVITES
		}

		if f, ok := function.(func(*InviteDelete)); ok {
			bot.Handlers.InviteDelete = append(bot.Handlers.InviteDelete, f)
			LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("added event handler")
			return nil
		}

	case FlagGatewayEventNameMessageCreate:
		if !bot.Config.Gateway.IntentSet[FlagIntentDIRECT_MESSAGES] {
			bot.Config.Gateway.IntentSet[FlagIntentDIRECT_MESSAGES] = true
			bot.Config.Gateway.Intents |= FlagIntentDIRECT_MESSAGES
		}

		if !bot.Config.Gateway.IntentSet[FlagIntentGUILD_MESSAGES] {
			bot.Config.Gateway.IntentSet[FlagIntentGUILD_MESSAGES] = true
			bot.Config.Gateway.Intents |= FlagIntentGUILD_MESSAGES
		}

		if f, ok := function.(func(*MessageCreate)); ok {
			bot.Handlers.MessageCreate = append(bot.Handlers.MessageCreate, f)
			LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("added event handler")
			return nil
		}

	case FlagGatewayEventNameMessageUpdate:
		if !bot.Config.Gateway.IntentSet[FlagIntentDIRECT_MESSAGES] {
			bot.Config.Gateway.IntentSet[FlagIntentDIRECT_MESSAGES] = true
			bot.Config.Gateway.Intents |= FlagIntentDIRECT_MESSAGES
		}

		if !bot.Config.Gateway.IntentSet[FlagIntentGUILD_MESSAGES] {
			bot.Config.Gateway.IntentSet[FlagIntentGUILD_MESSAGES] = true
			bot.Config.Gateway.Intents |= FlagIntentGUILD_MESSAGES
		}

		if f, ok := function.(func(*MessageUpdate)); ok {
			bot.Handlers.MessageUpdate = append(bot.Handlers.MessageUpdate, f)
			LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("added event handler")
			return nil
		}

	case FlagGatewayEventNameMessageDelete:
		if !bot.Config.Gateway.IntentSet[FlagIntentDIRECT_MESSAGES] {
			bot.Config.Gateway.IntentSet[FlagIntentDIRECT_MESSAGES] = true
			bot.Config.Gateway.Intents |= FlagIntentDIRECT_MESSAGES
		}

		if !bot.Config.Gateway.IntentSet[FlagIntentGUILD_MESSAGES] {
			bot.Config.Gateway.IntentSet[FlagIntentGUILD_MESSAGES] = true
			bot.Config.Gateway.Intents |= FlagIntentGUILD_MESSAGES
		}

		if f, ok := function.(func(*MessageDelete)); ok {
			bot.Handlers.MessageDelete = append(bot.Handlers.MessageDelete, f)
			LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("added event handler")
			return nil
		}

	case FlagGatewayEventNameMessageDeleteBulk:
		if !bot.Config.Gateway.IntentSet[FlagIntentGUILD_MESSAGES] {
			bot.Config.Gateway.IntentSet[FlagIntentGUILD_MESSAGES] = true
			bot.Config.Gateway.Intents |= FlagIntentGUILD_MESSAGES
		}

		if f, ok := function.(func(*MessageDeleteBulk)); ok {
			bot.Handlers.MessageDeleteBulk = append(bot.Handlers.MessageDeleteBulk, f)
			LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("added event handler")
			return nil
		}

	case FlagGatewayEventNameMessageReactionAdd:
		if !bot.Config.Gateway.IntentSet[FlagIntentDIRECT_MESSAGE_REACTIONS] {
			bot.Config.Gateway.IntentSet[FlagIntentDIRECT_MESSAGE_REACTIONS] = true
			bot.Config.Gateway.Intents |= FlagIntentDIRECT_MESSAGE_REACTIONS
		}

		if !bot.Config.Gateway.IntentSet[FlagIntentGUILD_MESSAGE_REACTIONS] {
			bot.Config.Gateway.IntentSet[FlagIntentGUILD_MESSAGE_REACTIONS] = true
			bot.Config.Gateway.Intents |= FlagIntentGUILD_MESSAGE_REACTIONS
		}

		if f, ok := function.(func(*MessageReactionAdd)); ok {
			bot.Handlers.MessageReactionAdd = append(bot.Handlers.MessageReactionAdd, f)
			LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("added event handler")
			return nil
		}

	case FlagGatewayEventNameMessageReactionRemove:
		if !bot.Config.Gateway.IntentSet[FlagIntentDIRECT_MESSAGE_REACTIONS] {
			bot.Config.Gateway.IntentSet[FlagIntentDIRECT_MESSAGE_REACTIONS] = true
			bot.Config.Gateway.Intents |= FlagIntentDIRECT_MESSAGE_REACTIONS
		}

		if !bot.Config.Gateway.IntentSet[FlagIntentGUILD_MESSAGE_REACTIONS] {
			bot.Config.Gateway.IntentSet[FlagIntentGUILD_MESSAGE_REACTIONS] = true
			bot.Config.Gateway.Intents |= FlagIntentGUILD_MESSAGE_REACTIONS
		}

		if f, ok := function.(func(*MessageReactionRemove)); ok {
			bot.Handlers.MessageReactionRemove = append(bot.Handlers.MessageReactionRemove, f)
			LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("added event handler")
			return nil
		}

	case FlagGatewayEventNameMessageReactionRemoveAll:
		if !bot.Config.Gateway.IntentSet[FlagIntentDIRECT_MESSAGE_REACTIONS] {
			bot.Config.Gateway.IntentSet[FlagIntentDIRECT_MESSAGE_REACTIONS] = true
			bot.Config.Gateway.Intents |= FlagIntentDIRECT_MESSAGE_REACTIONS
		}

		if !bot.Config.Gateway.IntentSet[FlagIntentGUILD_MESSAGE_REACTIONS] {
			bot.Config.Gateway.IntentSet[FlagIntentGUILD_MESSAGE_REACTIONS] = true
			bot.Config.Gateway.Intents |= FlagIntentGUILD_MESSAGE_REACTIONS
		}

		if f, ok := function.(func(*MessageReactionRemoveAll)); ok {
			bot.Handlers.MessageReactionRemoveAll = append(bot.Handlers.MessageReactionRemoveAll, f)
			LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("added event handler")
			return nil
		}

	case FlagGatewayEventNameMessageReactionRemoveEmoji:
		if !bot.Config.Gateway.IntentSet[FlagIntentDIRECT_MESSAGE_REACTIONS] {
			bot.Config.Gateway.IntentSet[FlagIntentDIRECT_MESSAGE_REACTIONS] = true
			bot.Config.Gateway.Intents |= FlagIntentDIRECT_MESSAGE_REACTIONS
		}

		if !bot.Config.Gateway.IntentSet[FlagIntentGUILD_MESSAGE_REACTIONS] {
			bot.Config.Gateway.IntentSet[FlagIntentGUILD_MESSAGE_REACTIONS] = true
			bot.Config.Gateway.Intents |= FlagIntentGUILD_MESSAGE_REACTIONS
		}

		if f, ok := function.(func(*MessageReactionRemoveEmoji)); ok {
			bot.Handlers.MessageReactionRemoveEmoji = append(bot.Handlers.MessageReactionRemoveEmoji, f)
			LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("added event handler")
			return nil
		}

	case FlagGatewayEventNamePresenceUpdate:
		if !bot.Config.Gateway.IntentSet[FlagIntentGUILD_PRESENCES] {
			bot.Config.Gateway.IntentSet[FlagIntentGUILD_PRESENCES] = true
			bot.Config.Gateway.Intents |= FlagIntentGUILD_PRESENCES
		}

		if f, ok := function.(func(*PresenceUpdate)); ok {
			bot.Handlers.PresenceUpdate = append(bot.Handlers.PresenceUpdate, f)
			LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("added event handler")
			return nil
		}

	case FlagGatewayEventNameStageInstanceCreate:
		if !bot.Config.Gateway.IntentSet[FlagIntentGUILDS] {
			bot.Config.Gateway.IntentSet[FlagIntentGUILDS] = true
			bot.Config.Gateway.Intents |= FlagIntentGUILDS
		}

		if f, ok := function.(func(*StageInstanceCreate)); ok {
			bot.Handlers.StageInstanceCreate = append(bot.Handlers.StageInstanceCreate, f)
			LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("added event handler")
			return nil
		}

	case FlagGatewayEventNameStageInstanceDelete:
		if !bot.Config.Gateway.IntentSet[FlagIntentGUILDS] {
			bot.Config.Gateway.IntentSet[FlagIntentGUILDS] = true
			bot.Config.Gateway.Intents |= FlagIntentGUILDS
		}

		if f, ok := function.(func(*StageInstanceDelete)); ok {
			bot.Handlers.StageInstanceDelete = append(bot.Handlers.StageInstanceDelete, f)
			LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("added event handler")
			return nil
		}

	case FlagGatewayEventNameStageInstanceUpdate:
		if !bot.Config.Gateway.IntentSet[FlagIntentGUILDS] {
			bot.Config.Gateway.IntentSet[FlagIntentGUILDS] = true
			bot.Config.Gateway.Intents |= FlagIntentGUILDS
		}

		if f, ok := function.(func(*StageInstanceUpdate)); ok {
			bot.Handlers.StageInstanceUpdate = append(bot.Handlers.StageInstanceUpdate, f)
			LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("added event handler")
			return nil
		}

	case FlagGatewayEventNameTypingStart:
		if !bot.Config.Gateway.IntentSet[FlagIntentDIRECT_MESSAGE_TYPING] {
			bot.Config.Gateway.IntentSet[FlagIntentDIRECT_MESSAGE_TYPING] = true
			bot.Config.Gateway.Intents |= FlagIntentDIRECT_MESSAGE_TYPING
		}

		if !bot.Config.Gateway.IntentSet[FlagIntentGUILD_MESSAGE_REACTIONS] {
			bot.Config.Gateway.IntentSet[FlagIntentGUILD_MESSAGE_REACTIONS] = true
			bot.Config.Gateway.Intents |= FlagIntentGUILD_MESSAGE_REACTIONS
		}

		if f, ok := function.(func(*TypingStart)); ok {
			bot.Handlers.TypingStart = append(bot.Handlers.TypingStart, f)
			LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("added event handler")
			return nil
		}

	case FlagGatewayEventNameVoiceStateUpdate:
		if !bot.Config.Gateway.IntentSet[FlagIntentGUILD_VOICE_STATES] {
			bot.Config.Gateway.IntentSet[FlagIntentGUILD_VOICE_STATES] = true
			bot.Config.Gateway.Intents |= FlagIntentGUILD_VOICE_STATES
		}

		if f, ok := function.(func(*VoiceStateUpdate)); ok {
			bot.Handlers.VoiceStateUpdate = append(bot.Handlers.VoiceStateUpdate, f)
			LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("added event handler")
			return nil
		}

	case FlagGatewayEventNameWebhooksUpdate:
		if !bot.Config.Gateway.IntentSet[FlagIntentGUILD_WEBHOOKS] {
			bot.Config.Gateway.IntentSet[FlagIntentGUILD_WEBHOOKS] = true
			bot.Config.Gateway.Intents |= FlagIntentGUILD_WEBHOOKS
		}

		if f, ok := function.(func(*WebhooksUpdate)); ok {
			bot.Handlers.WebhooksUpdate = append(bot.Handlers.WebhooksUpdate, f)
			LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("added event handler")
			return nil
		}
	}

	err := ErrorEventHandler{
		ClientID: bot.ApplicationID,
		Event:    eventname,
		Err:      fmt.Errorf("%s", errHandleNotRemoved),
	}
	LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")

	return err
}

// Remove removes the event handler at the given index from the bot.
// This function does NOT remove intents automatically.
func (bot *Client) Remove(eventname string, index int) error {
	bot.Handlers.mu.Lock()
	defer bot.Handlers.mu.Unlock()

	switch eventname {
	case FlagGatewayEventNameHello:
		if len(bot.Handlers.Hello) <= index {
			err := ErrorEventHandler{
				ClientID: bot.ApplicationID,
				Event:    eventname,
				Err:      fmt.Errorf(errRemoveInvalidIndex, index),
			}
			LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")
			return err
		}

		bot.Handlers.Hello = append(bot.Handlers.Hello[:index], bot.Handlers.Hello[index+1:]...)

	case FlagGatewayEventNameReady:
		if len(bot.Handlers.Ready) <= index {
			err := ErrorEventHandler{
				ClientID: bot.ApplicationID,
				Event:    eventname,
				Err:      fmt.Errorf(errRemoveInvalidIndex, index),
			}
			LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")
			return err
		}

		bot.Handlers.Ready = append(bot.Handlers.Ready[:index], bot.Handlers.Ready[index+1:]...)

	case FlagGatewayEventNameResumed:
		if len(bot.Handlers.Resumed) <= index {
			err := ErrorEventHandler{
				ClientID: bot.ApplicationID,
				Event:    eventname,
				Err:      fmt.Errorf(errRemoveInvalidIndex, index),
			}
			LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")
			return err
		}

		bot.Handlers.Resumed = append(bot.Handlers.Resumed[:index], bot.Handlers.Resumed[index+1:]...)

	case FlagGatewayEventNameReconnect:
		if len(bot.Handlers.Reconnect) <= index {
			err := ErrorEventHandler{
				ClientID: bot.ApplicationID,
				Event:    eventname,
				Err:      fmt.Errorf(errRemoveInvalidIndex, index),
			}
			LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")
			return err
		}

		bot.Handlers.Reconnect = append(bot.Handlers.Reconnect[:index], bot.Handlers.Reconnect[index+1:]...)

	case FlagGatewayEventNameInvalidSession:
		if len(bot.Handlers.InvalidSession) <= index {
			err := ErrorEventHandler{
				ClientID: bot.ApplicationID,
				Event:    eventname,
				Err:      fmt.Errorf(errRemoveInvalidIndex, index),
			}
			LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")
			return err
		}

		bot.Handlers.InvalidSession = append(bot.Handlers.InvalidSession[:index], bot.Handlers.InvalidSession[index+1:]...)

	case FlagGatewayEventNameApplicationCommandPermissionsUpdate:
		if len(bot.Handlers.ApplicationCommandPermissionsUpdate) <= index {
			err := ErrorEventHandler{
				ClientID: bot.ApplicationID,
				Event:    eventname,
				Err:      fmt.Errorf(errRemoveInvalidIndex, index),
			}
			LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")
			return err
		}

		bot.Handlers.ApplicationCommandPermissionsUpdate = append(bot.Handlers.ApplicationCommandPermissionsUpdate[:index], bot.Handlers.ApplicationCommandPermissionsUpdate[index+1:]...)

	case FlagGatewayEventNameAutoModerationRuleCreate:
		if len(bot.Handlers.AutoModerationRuleCreate) <= index {
			err := ErrorEventHandler{
				ClientID: bot.ApplicationID,
				Event:    eventname,
				Err:      fmt.Errorf(errRemoveInvalidIndex, index),
			}
			LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")
			return err
		}

		bot.Handlers.AutoModerationRuleCreate = append(bot.Handlers.AutoModerationRuleCreate[:index], bot.Handlers.AutoModerationRuleCreate[index+1:]...)

	case FlagGatewayEventNameAutoModerationRuleUpdate:
		if len(bot.Handlers.AutoModerationRuleUpdate) <= index {
			err := ErrorEventHandler{
				ClientID: bot.ApplicationID,
				Event:    eventname,
				Err:      fmt.Errorf(errRemoveInvalidIndex, index),
			}
			LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")
			return err
		}

		bot.Handlers.AutoModerationRuleUpdate = append(bot.Handlers.AutoModerationRuleUpdate[:index], bot.Handlers.AutoModerationRuleUpdate[index+1:]...)

	case FlagGatewayEventNameAutoModerationRuleDelete:
		if len(bot.Handlers.AutoModerationRuleDelete) <= index {
			err := ErrorEventHandler{
				ClientID: bot.ApplicationID,
				Event:    eventname,
				Err:      fmt.Errorf(errRemoveInvalidIndex, index),
			}
			LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")
			return err
		}

		bot.Handlers.AutoModerationRuleDelete = append(bot.Handlers.AutoModerationRuleDelete[:index], bot.Handlers.AutoModerationRuleDelete[index+1:]...)

	case FlagGatewayEventNameAutoModerationActionExecution:
		if len(bot.Handlers.AutoModerationActionExecution) <= index {
			err := ErrorEventHandler{
				ClientID: bot.ApplicationID,
				Event:    eventname,
				Err:      fmt.Errorf(errRemoveInvalidIndex, index),
			}
			LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")
			return err
		}

		bot.Handlers.AutoModerationActionExecution = append(bot.Handlers.AutoModerationActionExecution[:index], bot.Handlers.AutoModerationActionExecution[index+1:]...)

	case FlagGatewayEventNameInteractionCreate:
		if len(bot.Handlers.InteractionCreate) <= index {
			err := ErrorEventHandler{
				ClientID: bot.ApplicationID,
				Event:    eventname,
				Err:      fmt.Errorf(errRemoveInvalidIndex, index),
			}
			LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")
			return err
		}

		bot.Handlers.InteractionCreate = append(bot.Handlers.InteractionCreate[:index], bot.Handlers.InteractionCreate[index+1:]...)

	case FlagGatewayEventNameVoiceServerUpdate:
		if len(bot.Handlers.VoiceServerUpdate) <= index {
			err := ErrorEventHandler{
				ClientID: bot.ApplicationID,
				Event:    eventname,
				Err:      fmt.Errorf(errRemoveInvalidIndex, index),
			}
			LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")
			return err
		}

		bot.Handlers.VoiceServerUpdate = append(bot.Handlers.VoiceServerUpdate[:index], bot.Handlers.VoiceServerUpdate[index+1:]...)

	case FlagGatewayEventNameGuildMembersChunk:
		if len(bot.Handlers.GuildMembersChunk) <= index {
			err := ErrorEventHandler{
				ClientID: bot.ApplicationID,
				Event:    eventname,
				Err:      fmt.Errorf(errRemoveInvalidIndex, index),
			}
			LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")
			return err
		}

		bot.Handlers.GuildMembersChunk = append(bot.Handlers.GuildMembersChunk[:index], bot.Handlers.GuildMembersChunk[index+1:]...)

	case FlagGatewayEventNameUserUpdate:
		if len(bot.Handlers.UserUpdate) <= index {
			err := ErrorEventHandler{
				ClientID: bot.ApplicationID,
				Event:    eventname,
				Err:      fmt.Errorf(errRemoveInvalidIndex, index),
			}
			LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")
			return err
		}

		bot.Handlers.UserUpdate = append(bot.Handlers.UserUpdate[:index], bot.Handlers.UserUpdate[index+1:]...)

	case FlagGatewayEventNameChannelCreate:
		if len(bot.Handlers.ChannelCreate) <= index {
			err := ErrorEventHandler{
				ClientID: bot.ApplicationID,
				Event:    eventname,
				Err:      fmt.Errorf(errRemoveInvalidIndex, index),
			}
			LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")
			return err
		}

		bot.Handlers.ChannelCreate = append(bot.Handlers.ChannelCreate[:index], bot.Handlers.ChannelCreate[index+1:]...)

	case FlagGatewayEventNameChannelUpdate:
		if len(bot.Handlers.ChannelUpdate) <= index {
			err := ErrorEventHandler{
				ClientID: bot.ApplicationID,
				Event:    eventname,
				Err:      fmt.Errorf(errRemoveInvalidIndex, index),
			}
			LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")
			return err
		}

		bot.Handlers.ChannelUpdate = append(bot.Handlers.ChannelUpdate[:index], bot.Handlers.ChannelUpdate[index+1:]...)

	case FlagGatewayEventNameChannelDelete:
		if len(bot.Handlers.ChannelDelete) <= index {
			err := ErrorEventHandler{
				ClientID: bot.ApplicationID,
				Event:    eventname,
				Err:      fmt.Errorf(errRemoveInvalidIndex, index),
			}
			LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")
			return err
		}

		bot.Handlers.ChannelDelete = append(bot.Handlers.ChannelDelete[:index], bot.Handlers.ChannelDelete[index+1:]...)

	case FlagGatewayEventNameChannelPinsUpdate:
		if len(bot.Handlers.ChannelPinsUpdate) <= index {
			err := ErrorEventHandler{
				ClientID: bot.ApplicationID,
				Event:    eventname,
				Err:      fmt.Errorf(errRemoveInvalidIndex, index),
			}
			LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")
			return err
		}

		bot.Handlers.ChannelPinsUpdate = append(bot.Handlers.ChannelPinsUpdate[:index], bot.Handlers.ChannelPinsUpdate[index+1:]...)

	case FlagGatewayEventNameThreadCreate:
		if len(bot.Handlers.ThreadCreate) <= index {
			err := ErrorEventHandler{
				ClientID: bot.ApplicationID,
				Event:    eventname,
				Err:      fmt.Errorf(errRemoveInvalidIndex, index),
			}
			LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")
			return err
		}

		bot.Handlers.ThreadCreate = append(bot.Handlers.ThreadCreate[:index], bot.Handlers.ThreadCreate[index+1:]...)

	case FlagGatewayEventNameThreadUpdate:
		if len(bot.Handlers.ThreadUpdate) <= index {
			err := ErrorEventHandler{
				ClientID: bot.ApplicationID,
				Event:    eventname,
				Err:      fmt.Errorf(errRemoveInvalidIndex, index),
			}
			LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")
			return err
		}

		bot.Handlers.ThreadUpdate = append(bot.Handlers.ThreadUpdate[:index], bot.Handlers.ThreadUpdate[index+1:]...)

	case FlagGatewayEventNameThreadDelete:
		if len(bot.Handlers.ThreadDelete) <= index {
			err := ErrorEventHandler{
				ClientID: bot.ApplicationID,
				Event:    eventname,
				Err:      fmt.Errorf(errRemoveInvalidIndex, index),
			}
			LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")
			return err
		}

		bot.Handlers.ThreadDelete = append(bot.Handlers.ThreadDelete[:index], bot.Handlers.ThreadDelete[index+1:]...)

	case FlagGatewayEventNameThreadListSync:
		if len(bot.Handlers.ThreadListSync) <= index {
			err := ErrorEventHandler{
				ClientID: bot.ApplicationID,
				Event:    eventname,
				Err:      fmt.Errorf(errRemoveInvalidIndex, index),
			}
			LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")
			return err
		}

		bot.Handlers.ThreadListSync = append(bot.Handlers.ThreadListSync[:index], bot.Handlers.ThreadListSync[index+1:]...)

	case FlagGatewayEventNameThreadMemberUpdate:
		if len(bot.Handlers.ThreadMemberUpdate) <= index {
			err := ErrorEventHandler{
				ClientID: bot.ApplicationID,
				Event:    eventname,
				Err:      fmt.Errorf(errRemoveInvalidIndex, index),
			}
			LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")
			return err
		}

		bot.Handlers.ThreadMemberUpdate = append(bot.Handlers.ThreadMemberUpdate[:index], bot.Handlers.ThreadMemberUpdate[index+1:]...)

	case FlagGatewayEventNameThreadMembersUpdate:
		if len(bot.Handlers.ThreadMembersUpdate) <= index {
			err := ErrorEventHandler{
				ClientID: bot.ApplicationID,
				Event:    eventname,
				Err:      fmt.Errorf(errRemoveInvalidIndex, index),
			}
			LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")
			return err
		}

		bot.Handlers.ThreadMembersUpdate = append(bot.Handlers.ThreadMembersUpdate[:index], bot.Handlers.ThreadMembersUpdate[index+1:]...)

	case FlagGatewayEventNameGuildCreate:
		if len(bot.Handlers.GuildCreate) <= index {
			err := ErrorEventHandler{
				ClientID: bot.ApplicationID,
				Event:    eventname,
				Err:      fmt.Errorf(errRemoveInvalidIndex, index),
			}
			LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")
			return err
		}

		bot.Handlers.GuildCreate = append(bot.Handlers.GuildCreate[:index], bot.Handlers.GuildCreate[index+1:]...)

	case FlagGatewayEventNameGuildUpdate:
		if len(bot.Handlers.GuildUpdate) <= index {
			err := ErrorEventHandler{
				ClientID: bot.ApplicationID,
				Event:    eventname,
				Err:      fmt.Errorf(errRemoveInvalidIndex, index),
			}
			LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")
			return err
		}

		bot.Handlers.GuildUpdate = append(bot.Handlers.GuildUpdate[:index], bot.Handlers.GuildUpdate[index+1:]...)

	case FlagGatewayEventNameGuildDelete:
		if len(bot.Handlers.GuildDelete) <= index {
			err := ErrorEventHandler{
				ClientID: bot.ApplicationID,
				Event:    eventname,
				Err:      fmt.Errorf(errRemoveInvalidIndex, index),
			}
			LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")
			return err
		}

		bot.Handlers.GuildDelete = append(bot.Handlers.GuildDelete[:index], bot.Handlers.GuildDelete[index+1:]...)

	case FlagGatewayEventNameGuildBanAdd:
		if len(bot.Handlers.GuildBanAdd) <= index {
			err := ErrorEventHandler{
				ClientID: bot.ApplicationID,
				Event:    eventname,
				Err:      fmt.Errorf(errRemoveInvalidIndex, index),
			}
			LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")
			return err
		}

		bot.Handlers.GuildBanAdd = append(bot.Handlers.GuildBanAdd[:index], bot.Handlers.GuildBanAdd[index+1:]...)

	case FlagGatewayEventNameGuildBanRemove:
		if len(bot.Handlers.GuildBanRemove) <= index {
			err := ErrorEventHandler{
				ClientID: bot.ApplicationID,
				Event:    eventname,
				Err:      fmt.Errorf(errRemoveInvalidIndex, index),
			}
			LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")
			return err
		}

		bot.Handlers.GuildBanRemove = append(bot.Handlers.GuildBanRemove[:index], bot.Handlers.GuildBanRemove[index+1:]...)

	case FlagGatewayEventNameGuildEmojisUpdate:
		if len(bot.Handlers.GuildEmojisUpdate) <= index {
			err := ErrorEventHandler{
				ClientID: bot.ApplicationID,
				Event:    eventname,
				Err:      fmt.Errorf(errRemoveInvalidIndex, index),
			}
			LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")
			return err
		}

		bot.Handlers.GuildEmojisUpdate = append(bot.Handlers.GuildEmojisUpdate[:index], bot.Handlers.GuildEmojisUpdate[index+1:]...)

	case FlagGatewayEventNameGuildStickersUpdate:
		if len(bot.Handlers.GuildStickersUpdate) <= index {
			err := ErrorEventHandler{
				ClientID: bot.ApplicationID,
				Event:    eventname,
				Err:      fmt.Errorf(errRemoveInvalidIndex, index),
			}
			LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")
			return err
		}

		bot.Handlers.GuildStickersUpdate = append(bot.Handlers.GuildStickersUpdate[:index], bot.Handlers.GuildStickersUpdate[index+1:]...)

	case FlagGatewayEventNameGuildIntegrationsUpdate:
		if len(bot.Handlers.GuildIntegrationsUpdate) <= index {
			err := ErrorEventHandler{
				ClientID: bot.ApplicationID,
				Event:    eventname,
				Err:      fmt.Errorf(errRemoveInvalidIndex, index),
			}
			LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")
			return err
		}

		bot.Handlers.GuildIntegrationsUpdate = append(bot.Handlers.GuildIntegrationsUpdate[:index], bot.Handlers.GuildIntegrationsUpdate[index+1:]...)

	case FlagGatewayEventNameGuildMemberAdd:
		if len(bot.Handlers.GuildMemberAdd) <= index {
			err := ErrorEventHandler{
				ClientID: bot.ApplicationID,
				Event:    eventname,
				Err:      fmt.Errorf(errRemoveInvalidIndex, index),
			}
			LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")
			return err
		}

		bot.Handlers.GuildMemberAdd = append(bot.Handlers.GuildMemberAdd[:index], bot.Handlers.GuildMemberAdd[index+1:]...)

	case FlagGatewayEventNameGuildMemberRemove:
		if len(bot.Handlers.GuildMemberRemove) <= index {
			err := ErrorEventHandler{
				ClientID: bot.ApplicationID,
				Event:    eventname,
				Err:      fmt.Errorf(errRemoveInvalidIndex, index),
			}
			LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")
			return err
		}

		bot.Handlers.GuildMemberRemove = append(bot.Handlers.GuildMemberRemove[:index], bot.Handlers.GuildMemberRemove[index+1:]...)

	case FlagGatewayEventNameGuildMemberUpdate:
		if len(bot.Handlers.GuildMemberUpdate) <= index {
			err := ErrorEventHandler{
				ClientID: bot.ApplicationID,
				Event:    eventname,
				Err:      fmt.Errorf(errRemoveInvalidIndex, index),
			}
			LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")
			return err
		}

		bot.Handlers.GuildMemberUpdate = append(bot.Handlers.GuildMemberUpdate[:index], bot.Handlers.GuildMemberUpdate[index+1:]...)

	case FlagGatewayEventNameGuildRoleCreate:
		if len(bot.Handlers.GuildRoleCreate) <= index {
			err := ErrorEventHandler{
				ClientID: bot.ApplicationID,
				Event:    eventname,
				Err:      fmt.Errorf(errRemoveInvalidIndex, index),
			}
			LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")
			return err
		}

		bot.Handlers.GuildRoleCreate = append(bot.Handlers.GuildRoleCreate[:index], bot.Handlers.GuildRoleCreate[index+1:]...)

	case FlagGatewayEventNameGuildRoleUpdate:
		if len(bot.Handlers.GuildRoleUpdate) <= index {
			err := ErrorEventHandler{
				ClientID: bot.ApplicationID,
				Event:    eventname,
				Err:      fmt.Errorf(errRemoveInvalidIndex, index),
			}
			LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")
			return err
		}

		bot.Handlers.GuildRoleUpdate = append(bot.Handlers.GuildRoleUpdate[:index], bot.Handlers.GuildRoleUpdate[index+1:]...)

	case FlagGatewayEventNameGuildRoleDelete:
		if len(bot.Handlers.GuildRoleDelete) <= index {
			err := ErrorEventHandler{
				ClientID: bot.ApplicationID,
				Event:    eventname,
				Err:      fmt.Errorf(errRemoveInvalidIndex, index),
			}
			LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")
			return err
		}

		bot.Handlers.GuildRoleDelete = append(bot.Handlers.GuildRoleDelete[:index], bot.Handlers.GuildRoleDelete[index+1:]...)

	case FlagGatewayEventNameGuildScheduledEventCreate:
		if len(bot.Handlers.GuildScheduledEventCreate) <= index {
			err := ErrorEventHandler{
				ClientID: bot.ApplicationID,
				Event:    eventname,
				Err:      fmt.Errorf(errRemoveInvalidIndex, index),
			}
			LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")
			return err
		}

		bot.Handlers.GuildScheduledEventCreate = append(bot.Handlers.GuildScheduledEventCreate[:index], bot.Handlers.GuildScheduledEventCreate[index+1:]...)

	case FlagGatewayEventNameGuildScheduledEventUpdate:
		if len(bot.Handlers.GuildScheduledEventUpdate) <= index {
			err := ErrorEventHandler{
				ClientID: bot.ApplicationID,
				Event:    eventname,
				Err:      fmt.Errorf(errRemoveInvalidIndex, index),
			}
			LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")
			return err
		}

		bot.Handlers.GuildScheduledEventUpdate = append(bot.Handlers.GuildScheduledEventUpdate[:index], bot.Handlers.GuildScheduledEventUpdate[index+1:]...)

	case FlagGatewayEventNameGuildScheduledEventDelete:
		if len(bot.Handlers.GuildScheduledEventDelete) <= index {
			err := ErrorEventHandler{
				ClientID: bot.ApplicationID,
				Event:    eventname,
				Err:      fmt.Errorf(errRemoveInvalidIndex, index),
			}
			LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")
			return err
		}

		bot.Handlers.GuildScheduledEventDelete = append(bot.Handlers.GuildScheduledEventDelete[:index], bot.Handlers.GuildScheduledEventDelete[index+1:]...)

	case FlagGatewayEventNameGuildScheduledEventUserAdd:
		if len(bot.Handlers.GuildScheduledEventUserAdd) <= index {
			err := ErrorEventHandler{
				ClientID: bot.ApplicationID,
				Event:    eventname,
				Err:      fmt.Errorf(errRemoveInvalidIndex, index),
			}
			LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")
			return err
		}

		bot.Handlers.GuildScheduledEventUserAdd = append(bot.Handlers.GuildScheduledEventUserAdd[:index], bot.Handlers.GuildScheduledEventUserAdd[index+1:]...)

	case FlagGatewayEventNameGuildScheduledEventUserRemove:
		if len(bot.Handlers.GuildScheduledEventUserRemove) <= index {
			err := ErrorEventHandler{
				ClientID: bot.ApplicationID,
				Event:    eventname,
				Err:      fmt.Errorf(errRemoveInvalidIndex, index),
			}
			LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")
			return err
		}

		bot.Handlers.GuildScheduledEventUserRemove = append(bot.Handlers.GuildScheduledEventUserRemove[:index], bot.Handlers.GuildScheduledEventUserRemove[index+1:]...)

	case FlagGatewayEventNameIntegrationCreate:
		if len(bot.Handlers.IntegrationCreate) <= index {
			err := ErrorEventHandler{
				ClientID: bot.ApplicationID,
				Event:    eventname,
				Err:      fmt.Errorf(errRemoveInvalidIndex, index),
			}
			LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")
			return err
		}

		bot.Handlers.IntegrationCreate = append(bot.Handlers.IntegrationCreate[:index], bot.Handlers.IntegrationCreate[index+1:]...)

	case FlagGatewayEventNameIntegrationUpdate:
		if len(bot.Handlers.IntegrationUpdate) <= index {
			err := ErrorEventHandler{
				ClientID: bot.ApplicationID,
				Event:    eventname,
				Err:      fmt.Errorf(errRemoveInvalidIndex, index),
			}
			LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")
			return err
		}

		bot.Handlers.IntegrationUpdate = append(bot.Handlers.IntegrationUpdate[:index], bot.Handlers.IntegrationUpdate[index+1:]...)

	case FlagGatewayEventNameIntegrationDelete:
		if len(bot.Handlers.IntegrationDelete) <= index {
			err := ErrorEventHandler{
				ClientID: bot.ApplicationID,
				Event:    eventname,
				Err:      fmt.Errorf(errRemoveInvalidIndex, index),
			}
			LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")
			return err
		}

		bot.Handlers.IntegrationDelete = append(bot.Handlers.IntegrationDelete[:index], bot.Handlers.IntegrationDelete[index+1:]...)

	case FlagGatewayEventNameInviteCreate:
		if len(bot.Handlers.InviteCreate) <= index {
			err := ErrorEventHandler{
				ClientID: bot.ApplicationID,
				Event:    eventname,
				Err:      fmt.Errorf(errRemoveInvalidIndex, index),
			}
			LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")
			return err
		}

		bot.Handlers.InviteCreate = append(bot.Handlers.InviteCreate[:index], bot.Handlers.InviteCreate[index+1:]...)

	case FlagGatewayEventNameInviteDelete:
		if len(bot.Handlers.InviteDelete) <= index {
			err := ErrorEventHandler{
				ClientID: bot.ApplicationID,
				Event:    eventname,
				Err:      fmt.Errorf(errRemoveInvalidIndex, index),
			}
			LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")
			return err
		}

		bot.Handlers.InviteDelete = append(bot.Handlers.InviteDelete[:index], bot.Handlers.InviteDelete[index+1:]...)

	case FlagGatewayEventNameMessageCreate:
		if len(bot.Handlers.MessageCreate) <= index {
			err := ErrorEventHandler{
				ClientID: bot.ApplicationID,
				Event:    eventname,
				Err:      fmt.Errorf(errRemoveInvalidIndex, index),
			}
			LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")
			return err
		}

		bot.Handlers.MessageCreate = append(bot.Handlers.MessageCreate[:index], bot.Handlers.MessageCreate[index+1:]...)

	case FlagGatewayEventNameMessageUpdate:
		if len(bot.Handlers.MessageUpdate) <= index {
			err := ErrorEventHandler{
				ClientID: bot.ApplicationID,
				Event:    eventname,
				Err:      fmt.Errorf(errRemoveInvalidIndex, index),
			}
			LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")
			return err
		}

		bot.Handlers.MessageUpdate = append(bot.Handlers.MessageUpdate[:index], bot.Handlers.MessageUpdate[index+1:]...)

	case FlagGatewayEventNameMessageDelete:
		if len(bot.Handlers.MessageDelete) <= index {
			err := ErrorEventHandler{
				ClientID: bot.ApplicationID,
				Event:    eventname,
				Err:      fmt.Errorf(errRemoveInvalidIndex, index),
			}
			LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")
			return err
		}

		bot.Handlers.MessageDelete = append(bot.Handlers.MessageDelete[:index], bot.Handlers.MessageDelete[index+1:]...)

	case FlagGatewayEventNameMessageDeleteBulk:
		if len(bot.Handlers.MessageDeleteBulk) <= index {
			err := ErrorEventHandler{
				ClientID: bot.ApplicationID,
				Event:    eventname,
				Err:      fmt.Errorf(errRemoveInvalidIndex, index),
			}
			LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")
			return err
		}

		bot.Handlers.MessageDeleteBulk = append(bot.Handlers.MessageDeleteBulk[:index], bot.Handlers.MessageDeleteBulk[index+1:]...)

	case FlagGatewayEventNameMessageReactionAdd:
		if len(bot.Handlers.MessageReactionAdd) <= index {
			err := ErrorEventHandler{
				ClientID: bot.ApplicationID,
				Event:    eventname,
				Err:      fmt.Errorf(errRemoveInvalidIndex, index),
			}
			LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")
			return err
		}

		bot.Handlers.MessageReactionAdd = append(bot.Handlers.MessageReactionAdd[:index], bot.Handlers.MessageReactionAdd[index+1:]...)

	case FlagGatewayEventNameMessageReactionRemove:
		if len(bot.Handlers.MessageReactionRemove) <= index {
			err := ErrorEventHandler{
				ClientID: bot.ApplicationID,
				Event:    eventname,
				Err:      fmt.Errorf(errRemoveInvalidIndex, index),
			}
			LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")
			return err
		}

		bot.Handlers.MessageReactionRemove = append(bot.Handlers.MessageReactionRemove[:index], bot.Handlers.MessageReactionRemove[index+1:]...)

	case FlagGatewayEventNameMessageReactionRemoveAll:
		if len(bot.Handlers.MessageReactionRemoveAll) <= index {
			err := ErrorEventHandler{
				ClientID: bot.ApplicationID,
				Event:    eventname,
				Err:      fmt.Errorf(errRemoveInvalidIndex, index),
			}
			LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")
			return err
		}

		bot.Handlers.MessageReactionRemoveAll = append(bot.Handlers.MessageReactionRemoveAll[:index], bot.Handlers.MessageReactionRemoveAll[index+1:]...)

	case FlagGatewayEventNameMessageReactionRemoveEmoji:
		if len(bot.Handlers.MessageReactionRemoveEmoji) <= index {
			err := ErrorEventHandler{
				ClientID: bot.ApplicationID,
				Event:    eventname,
				Err:      fmt.Errorf(errRemoveInvalidIndex, index),
			}
			LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")
			return err
		}

		bot.Handlers.MessageReactionRemoveEmoji = append(bot.Handlers.MessageReactionRemoveEmoji[:index], bot.Handlers.MessageReactionRemoveEmoji[index+1:]...)

	case FlagGatewayEventNamePresenceUpdate:
		if len(bot.Handlers.PresenceUpdate) <= index {
			err := ErrorEventHandler{
				ClientID: bot.ApplicationID,
				Event:    eventname,
				Err:      fmt.Errorf(errRemoveInvalidIndex, index),
			}
			LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")
			return err
		}

		bot.Handlers.PresenceUpdate = append(bot.Handlers.PresenceUpdate[:index], bot.Handlers.PresenceUpdate[index+1:]...)

	case FlagGatewayEventNameStageInstanceCreate:
		if len(bot.Handlers.StageInstanceCreate) <= index {
			err := ErrorEventHandler{
				ClientID: bot.ApplicationID,
				Event:    eventname,
				Err:      fmt.Errorf(errRemoveInvalidIndex, index),
			}
			LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")
			return err
		}

		bot.Handlers.StageInstanceCreate = append(bot.Handlers.StageInstanceCreate[:index], bot.Handlers.StageInstanceCreate[index+1:]...)

	case FlagGatewayEventNameStageInstanceDelete:
		if len(bot.Handlers.StageInstanceDelete) <= index {
			err := ErrorEventHandler{
				ClientID: bot.ApplicationID,
				Event:    eventname,
				Err:      fmt.Errorf(errRemoveInvalidIndex, index),
			}
			LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")
			return err
		}

		bot.Handlers.StageInstanceDelete = append(bot.Handlers.StageInstanceDelete[:index], bot.Handlers.StageInstanceDelete[index+1:]...)

	case FlagGatewayEventNameStageInstanceUpdate:
		if len(bot.Handlers.StageInstanceUpdate) <= index {
			err := ErrorEventHandler{
				ClientID: bot.ApplicationID,
				Event:    eventname,
				Err:      fmt.Errorf(errRemoveInvalidIndex, index),
			}
			LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")
			return err
		}

		bot.Handlers.StageInstanceUpdate = append(bot.Handlers.StageInstanceUpdate[:index], bot.Handlers.StageInstanceUpdate[index+1:]...)

	case FlagGatewayEventNameTypingStart:
		if len(bot.Handlers.TypingStart) <= index {
			err := ErrorEventHandler{
				ClientID: bot.ApplicationID,
				Event:    eventname,
				Err:      fmt.Errorf(errRemoveInvalidIndex, index),
			}
			LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")
			return err
		}

		bot.Handlers.TypingStart = append(bot.Handlers.TypingStart[:index], bot.Handlers.TypingStart[index+1:]...)

	case FlagGatewayEventNameVoiceStateUpdate:
		if len(bot.Handlers.VoiceStateUpdate) <= index {
			err := ErrorEventHandler{
				ClientID: bot.ApplicationID,
				Event:    eventname,
				Err:      fmt.Errorf(errRemoveInvalidIndex, index),
			}
			LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")
			return err
		}

		bot.Handlers.VoiceStateUpdate = append(bot.Handlers.VoiceStateUpdate[:index], bot.Handlers.VoiceStateUpdate[index+1:]...)

	case FlagGatewayEventNameWebhooksUpdate:
		if len(bot.Handlers.WebhooksUpdate) <= index {
			err := ErrorEventHandler{
				ClientID: bot.ApplicationID,
				Event:    eventname,
				Err:      fmt.Errorf(errRemoveInvalidIndex, index),
			}
			LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(err).Msg("")
			return err
		}

		bot.Handlers.WebhooksUpdate = append(bot.Handlers.WebhooksUpdate[:index], bot.Handlers.WebhooksUpdate[index+1:]...)
	}

	LogEventHandler(Logger.Info(), bot.ApplicationID, eventname).Msg("removed event handler")

	return nil
}

// handle handles an event using its name and data.
func (bot *Client) handle(eventname string, data json.RawMessage) {
	bot.Handlers.mu.RLock()
	defer bot.Handlers.mu.RUnlock()

	switch eventname {
	case FlagGatewayEventNameHello:
		if len(bot.Handlers.Hello) != 0 {
			event := new(Hello)
			if err := json.Unmarshal(data, event); err != nil {
				LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(ErrorEvent{ClientID: bot.ApplicationID, Event: FlagGatewayEventNameHello, Err: err, Action: ErrorEventActionUnmarshal}).Msg("")
				return
			}

			for _, handler := range bot.Handlers.Hello {
				go handler(event)
			}
		}

	case FlagGatewayEventNameReady:
		if len(bot.Handlers.Ready) != 0 {
			event := new(Ready)
			if err := json.Unmarshal(data, event); err != nil {
				LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(ErrorEvent{ClientID: bot.ApplicationID, Event: FlagGatewayEventNameReady, Err: err, Action: ErrorEventActionUnmarshal}).Msg("")
				return
			}

			for _, handler := range bot.Handlers.Ready {
				go handler(event)
			}
		}

	case FlagGatewayEventNameResumed:
		if len(bot.Handlers.Resumed) != 0 {
			event := new(Resumed)
			if err := json.Unmarshal(data, event); err != nil {
				LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(ErrorEvent{ClientID: bot.ApplicationID, Event: FlagGatewayEventNameResumed, Err: err, Action: ErrorEventActionUnmarshal}).Msg("")
				return
			}

			for _, handler := range bot.Handlers.Resumed {
				go handler(event)
			}
		}

	case FlagGatewayEventNameReconnect:
		if len(bot.Handlers.Reconnect) != 0 {
			event := new(Reconnect)
			if err := json.Unmarshal(data, event); err != nil {
				LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(ErrorEvent{ClientID: bot.ApplicationID, Event: FlagGatewayEventNameReconnect, Err: err, Action: ErrorEventActionUnmarshal}).Msg("")
				return
			}

			for _, handler := range bot.Handlers.Reconnect {
				go handler(event)
			}
		}

	case FlagGatewayEventNameInvalidSession:
		if len(bot.Handlers.InvalidSession) != 0 {
			event := new(InvalidSession)
			if err := json.Unmarshal(data, event); err != nil {
				LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(ErrorEvent{ClientID: bot.ApplicationID, Event: FlagGatewayEventNameInvalidSession, Err: err, Action: ErrorEventActionUnmarshal}).Msg("")
				return
			}

			for _, handler := range bot.Handlers.InvalidSession {
				go handler(event)
			}
		}

	case FlagGatewayEventNameApplicationCommandPermissionsUpdate:
		if len(bot.Handlers.ApplicationCommandPermissionsUpdate) != 0 {
			event := new(ApplicationCommandPermissionsUpdate)
			if err := json.Unmarshal(data, event); err != nil {
				LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(ErrorEvent{ClientID: bot.ApplicationID, Event: FlagGatewayEventNameApplicationCommandPermissionsUpdate, Err: err, Action: ErrorEventActionUnmarshal}).Msg("")
				return
			}

			for _, handler := range bot.Handlers.ApplicationCommandPermissionsUpdate {
				go handler(event)
			}
		}

	case FlagGatewayEventNameAutoModerationRuleCreate:
		if len(bot.Handlers.AutoModerationRuleCreate) != 0 {
			event := new(AutoModerationRuleCreate)
			if err := json.Unmarshal(data, event); err != nil {
				LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(ErrorEvent{ClientID: bot.ApplicationID, Event: FlagGatewayEventNameAutoModerationRuleCreate, Err: err, Action: ErrorEventActionUnmarshal}).Msg("")
				return
			}

			for _, handler := range bot.Handlers.AutoModerationRuleCreate {
				go handler(event)
			}
		}

	case FlagGatewayEventNameAutoModerationRuleUpdate:
		if len(bot.Handlers.AutoModerationRuleUpdate) != 0 {
			event := new(AutoModerationRuleUpdate)
			if err := json.Unmarshal(data, event); err != nil {
				LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(ErrorEvent{ClientID: bot.ApplicationID, Event: FlagGatewayEventNameAutoModerationRuleUpdate, Err: err, Action: ErrorEventActionUnmarshal}).Msg("")
				return
			}

			for _, handler := range bot.Handlers.AutoModerationRuleUpdate {
				go handler(event)
			}
		}

	case FlagGatewayEventNameAutoModerationRuleDelete:
		if len(bot.Handlers.AutoModerationRuleDelete) != 0 {
			event := new(AutoModerationRuleDelete)
			if err := json.Unmarshal(data, event); err != nil {
				LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(ErrorEvent{ClientID: bot.ApplicationID, Event: FlagGatewayEventNameAutoModerationRuleDelete, Err: err, Action: ErrorEventActionUnmarshal}).Msg("")
				return
			}

			for _, handler := range bot.Handlers.AutoModerationRuleDelete {
				go handler(event)
			}
		}

	case FlagGatewayEventNameAutoModerationActionExecution:
		if len(bot.Handlers.AutoModerationActionExecution) != 0 {
			event := new(AutoModerationActionExecution)
			if err := json.Unmarshal(data, event); err != nil {
				LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(ErrorEvent{ClientID: bot.ApplicationID, Event: FlagGatewayEventNameAutoModerationActionExecution, Err: err, Action: ErrorEventActionUnmarshal}).Msg("")
				return
			}

			for _, handler := range bot.Handlers.AutoModerationActionExecution {
				go handler(event)
			}
		}

	case FlagGatewayEventNameInteractionCreate:
		if len(bot.Handlers.InteractionCreate) != 0 {
			event := new(InteractionCreate)
			if err := json.Unmarshal(data, event); err != nil {
				LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(ErrorEvent{ClientID: bot.ApplicationID, Event: FlagGatewayEventNameInteractionCreate, Err: err, Action: ErrorEventActionUnmarshal}).Msg("")
				return
			}

			for _, handler := range bot.Handlers.InteractionCreate {
				go handler(event)
			}
		}

	case FlagGatewayEventNameVoiceServerUpdate:
		if len(bot.Handlers.VoiceServerUpdate) != 0 {
			event := new(VoiceServerUpdate)
			if err := json.Unmarshal(data, event); err != nil {
				LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(ErrorEvent{ClientID: bot.ApplicationID, Event: FlagGatewayEventNameVoiceServerUpdate, Err: err, Action: ErrorEventActionUnmarshal}).Msg("")
				return
			}

			for _, handler := range bot.Handlers.VoiceServerUpdate {
				go handler(event)
			}
		}

	case FlagGatewayEventNameGuildMembersChunk:
		if len(bot.Handlers.GuildMembersChunk) != 0 {
			event := new(GuildMembersChunk)
			if err := json.Unmarshal(data, event); err != nil {
				LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(ErrorEvent{ClientID: bot.ApplicationID, Event: FlagGatewayEventNameGuildMembersChunk, Err: err, Action: ErrorEventActionUnmarshal}).Msg("")
				return
			}

			for _, handler := range bot.Handlers.GuildMembersChunk {
				go handler(event)
			}
		}

	case FlagGatewayEventNameUserUpdate:
		if len(bot.Handlers.UserUpdate) != 0 {
			event := new(UserUpdate)
			if err := json.Unmarshal(data, event); err != nil {
				LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(ErrorEvent{ClientID: bot.ApplicationID, Event: FlagGatewayEventNameUserUpdate, Err: err, Action: ErrorEventActionUnmarshal}).Msg("")
				return
			}

			for _, handler := range bot.Handlers.UserUpdate {
				go handler(event)
			}
		}

	case FlagGatewayEventNameChannelCreate:
		if len(bot.Handlers.ChannelCreate) != 0 {
			event := new(ChannelCreate)
			if err := json.Unmarshal(data, event); err != nil {
				LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(ErrorEvent{ClientID: bot.ApplicationID, Event: FlagGatewayEventNameChannelCreate, Err: err, Action: ErrorEventActionUnmarshal}).Msg("")
				return
			}

			for _, handler := range bot.Handlers.ChannelCreate {
				go handler(event)
			}
		}

	case FlagGatewayEventNameChannelUpdate:
		if len(bot.Handlers.ChannelUpdate) != 0 {
			event := new(ChannelUpdate)
			if err := json.Unmarshal(data, event); err != nil {
				LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(ErrorEvent{ClientID: bot.ApplicationID, Event: FlagGatewayEventNameChannelUpdate, Err: err, Action: ErrorEventActionUnmarshal}).Msg("")
				return
			}

			for _, handler := range bot.Handlers.ChannelUpdate {
				go handler(event)
			}
		}

	case FlagGatewayEventNameChannelDelete:
		if len(bot.Handlers.ChannelDelete) != 0 {
			event := new(ChannelDelete)
			if err := json.Unmarshal(data, event); err != nil {
				LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(ErrorEvent{ClientID: bot.ApplicationID, Event: FlagGatewayEventNameChannelDelete, Err: err, Action: ErrorEventActionUnmarshal}).Msg("")
				return
			}

			for _, handler := range bot.Handlers.ChannelDelete {
				go handler(event)
			}
		}

	case FlagGatewayEventNameChannelPinsUpdate:
		if len(bot.Handlers.ChannelPinsUpdate) != 0 {
			event := new(ChannelPinsUpdate)
			if err := json.Unmarshal(data, event); err != nil {
				LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(ErrorEvent{ClientID: bot.ApplicationID, Event: FlagGatewayEventNameChannelPinsUpdate, Err: err, Action: ErrorEventActionUnmarshal}).Msg("")
				return
			}

			for _, handler := range bot.Handlers.ChannelPinsUpdate {
				go handler(event)
			}
		}

	case FlagGatewayEventNameThreadCreate:
		if len(bot.Handlers.ThreadCreate) != 0 {
			event := new(ThreadCreate)
			if err := json.Unmarshal(data, event); err != nil {
				LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(ErrorEvent{ClientID: bot.ApplicationID, Event: FlagGatewayEventNameThreadCreate, Err: err, Action: ErrorEventActionUnmarshal}).Msg("")
				return
			}

			for _, handler := range bot.Handlers.ThreadCreate {
				go handler(event)
			}
		}

	case FlagGatewayEventNameThreadUpdate:
		if len(bot.Handlers.ThreadUpdate) != 0 {
			event := new(ThreadUpdate)
			if err := json.Unmarshal(data, event); err != nil {
				LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(ErrorEvent{ClientID: bot.ApplicationID, Event: FlagGatewayEventNameThreadUpdate, Err: err, Action: ErrorEventActionUnmarshal}).Msg("")
				return
			}

			for _, handler := range bot.Handlers.ThreadUpdate {
				go handler(event)
			}
		}

	case FlagGatewayEventNameThreadDelete:
		if len(bot.Handlers.ThreadDelete) != 0 {
			event := new(ThreadDelete)
			if err := json.Unmarshal(data, event); err != nil {
				LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(ErrorEvent{ClientID: bot.ApplicationID, Event: FlagGatewayEventNameThreadDelete, Err: err, Action: ErrorEventActionUnmarshal}).Msg("")
				return
			}

			for _, handler := range bot.Handlers.ThreadDelete {
				go handler(event)
			}
		}

	case FlagGatewayEventNameThreadListSync:
		if len(bot.Handlers.ThreadListSync) != 0 {
			event := new(ThreadListSync)
			if err := json.Unmarshal(data, event); err != nil {
				LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(ErrorEvent{ClientID: bot.ApplicationID, Event: FlagGatewayEventNameThreadListSync, Err: err, Action: ErrorEventActionUnmarshal}).Msg("")
				return
			}

			for _, handler := range bot.Handlers.ThreadListSync {
				go handler(event)
			}
		}

	case FlagGatewayEventNameThreadMemberUpdate:
		if len(bot.Handlers.ThreadMemberUpdate) != 0 {
			event := new(ThreadMemberUpdate)
			if err := json.Unmarshal(data, event); err != nil {
				LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(ErrorEvent{ClientID: bot.ApplicationID, Event: FlagGatewayEventNameThreadMemberUpdate, Err: err, Action: ErrorEventActionUnmarshal}).Msg("")
				return
			}

			for _, handler := range bot.Handlers.ThreadMemberUpdate {
				go handler(event)
			}
		}

	case FlagGatewayEventNameThreadMembersUpdate:
		if len(bot.Handlers.ThreadMembersUpdate) != 0 {
			event := new(ThreadMembersUpdate)
			if err := json.Unmarshal(data, event); err != nil {
				LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(ErrorEvent{ClientID: bot.ApplicationID, Event: FlagGatewayEventNameThreadMembersUpdate, Err: err, Action: ErrorEventActionUnmarshal}).Msg("")
				return
			}

			for _, handler := range bot.Handlers.ThreadMembersUpdate {
				go handler(event)
			}
		}

	case FlagGatewayEventNameGuildCreate:
		if len(bot.Handlers.GuildCreate) != 0 {
			event := new(GuildCreate)
			if err := json.Unmarshal(data, event); err != nil {
				LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(ErrorEvent{ClientID: bot.ApplicationID, Event: FlagGatewayEventNameGuildCreate, Err: err, Action: ErrorEventActionUnmarshal}).Msg("")
				return
			}

			for _, handler := range bot.Handlers.GuildCreate {
				go handler(event)
			}
		}

	case FlagGatewayEventNameGuildUpdate:
		if len(bot.Handlers.GuildUpdate) != 0 {
			event := new(GuildUpdate)
			if err := json.Unmarshal(data, event); err != nil {
				LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(ErrorEvent{ClientID: bot.ApplicationID, Event: FlagGatewayEventNameGuildUpdate, Err: err, Action: ErrorEventActionUnmarshal}).Msg("")
				return
			}

			for _, handler := range bot.Handlers.GuildUpdate {
				go handler(event)
			}
		}

	case FlagGatewayEventNameGuildDelete:
		if len(bot.Handlers.GuildDelete) != 0 {
			event := new(GuildDelete)
			if err := json.Unmarshal(data, event); err != nil {
				LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(ErrorEvent{ClientID: bot.ApplicationID, Event: FlagGatewayEventNameGuildDelete, Err: err, Action: ErrorEventActionUnmarshal}).Msg("")
				return
			}

			for _, handler := range bot.Handlers.GuildDelete {
				go handler(event)
			}
		}

	case FlagGatewayEventNameGuildBanAdd:
		if len(bot.Handlers.GuildBanAdd) != 0 {
			event := new(GuildBanAdd)
			if err := json.Unmarshal(data, event); err != nil {
				LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(ErrorEvent{ClientID: bot.ApplicationID, Event: FlagGatewayEventNameGuildBanAdd, Err: err, Action: ErrorEventActionUnmarshal}).Msg("")
				return
			}

			for _, handler := range bot.Handlers.GuildBanAdd {
				go handler(event)
			}
		}

	case FlagGatewayEventNameGuildBanRemove:
		if len(bot.Handlers.GuildBanRemove) != 0 {
			event := new(GuildBanRemove)
			if err := json.Unmarshal(data, event); err != nil {
				LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(ErrorEvent{ClientID: bot.ApplicationID, Event: FlagGatewayEventNameGuildBanRemove, Err: err, Action: ErrorEventActionUnmarshal}).Msg("")
				return
			}

			for _, handler := range bot.Handlers.GuildBanRemove {
				go handler(event)
			}
		}

	case FlagGatewayEventNameGuildEmojisUpdate:
		if len(bot.Handlers.GuildEmojisUpdate) != 0 {
			event := new(GuildEmojisUpdate)
			if err := json.Unmarshal(data, event); err != nil {
				LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(ErrorEvent{ClientID: bot.ApplicationID, Event: FlagGatewayEventNameGuildEmojisUpdate, Err: err, Action: ErrorEventActionUnmarshal}).Msg("")
				return
			}

			for _, handler := range bot.Handlers.GuildEmojisUpdate {
				go handler(event)
			}
		}

	case FlagGatewayEventNameGuildStickersUpdate:
		if len(bot.Handlers.GuildStickersUpdate) != 0 {
			event := new(GuildStickersUpdate)
			if err := json.Unmarshal(data, event); err != nil {
				LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(ErrorEvent{ClientID: bot.ApplicationID, Event: FlagGatewayEventNameGuildStickersUpdate, Err: err, Action: ErrorEventActionUnmarshal}).Msg("")
				return
			}

			for _, handler := range bot.Handlers.GuildStickersUpdate {
				go handler(event)
			}
		}

	case FlagGatewayEventNameGuildIntegrationsUpdate:
		if len(bot.Handlers.GuildIntegrationsUpdate) != 0 {
			event := new(GuildIntegrationsUpdate)
			if err := json.Unmarshal(data, event); err != nil {
				LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(ErrorEvent{ClientID: bot.ApplicationID, Event: FlagGatewayEventNameGuildIntegrationsUpdate, Err: err, Action: ErrorEventActionUnmarshal}).Msg("")
				return
			}

			for _, handler := range bot.Handlers.GuildIntegrationsUpdate {
				go handler(event)
			}
		}

	case FlagGatewayEventNameGuildMemberAdd:
		if len(bot.Handlers.GuildMemberAdd) != 0 {
			event := new(GuildMemberAdd)
			if err := json.Unmarshal(data, event); err != nil {
				LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(ErrorEvent{ClientID: bot.ApplicationID, Event: FlagGatewayEventNameGuildMemberAdd, Err: err, Action: ErrorEventActionUnmarshal}).Msg("")
				return
			}

			for _, handler := range bot.Handlers.GuildMemberAdd {
				go handler(event)
			}
		}

	case FlagGatewayEventNameGuildMemberRemove:
		if len(bot.Handlers.GuildMemberRemove) != 0 {
			event := new(GuildMemberRemove)
			if err := json.Unmarshal(data, event); err != nil {
				LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(ErrorEvent{ClientID: bot.ApplicationID, Event: FlagGatewayEventNameGuildMemberRemove, Err: err, Action: ErrorEventActionUnmarshal}).Msg("")
				return
			}

			for _, handler := range bot.Handlers.GuildMemberRemove {
				go handler(event)
			}
		}

	case FlagGatewayEventNameGuildMemberUpdate:
		if len(bot.Handlers.GuildMemberUpdate) != 0 {
			event := new(GuildMemberUpdate)
			if err := json.Unmarshal(data, event); err != nil {
				LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(ErrorEvent{ClientID: bot.ApplicationID, Event: FlagGatewayEventNameGuildMemberUpdate, Err: err, Action: ErrorEventActionUnmarshal}).Msg("")
				return
			}

			for _, handler := range bot.Handlers.GuildMemberUpdate {
				go handler(event)
			}
		}

	case FlagGatewayEventNameGuildRoleCreate:
		if len(bot.Handlers.GuildRoleCreate) != 0 {
			event := new(GuildRoleCreate)
			if err := json.Unmarshal(data, event); err != nil {
				LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(ErrorEvent{ClientID: bot.ApplicationID, Event: FlagGatewayEventNameGuildRoleCreate, Err: err, Action: ErrorEventActionUnmarshal}).Msg("")
				return
			}

			for _, handler := range bot.Handlers.GuildRoleCreate {
				go handler(event)
			}
		}

	case FlagGatewayEventNameGuildRoleUpdate:
		if len(bot.Handlers.GuildRoleUpdate) != 0 {
			event := new(GuildRoleUpdate)
			if err := json.Unmarshal(data, event); err != nil {
				LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(ErrorEvent{ClientID: bot.ApplicationID, Event: FlagGatewayEventNameGuildRoleUpdate, Err: err, Action: ErrorEventActionUnmarshal}).Msg("")
				return
			}

			for _, handler := range bot.Handlers.GuildRoleUpdate {
				go handler(event)
			}
		}

	case FlagGatewayEventNameGuildRoleDelete:
		if len(bot.Handlers.GuildRoleDelete) != 0 {
			event := new(GuildRoleDelete)
			if err := json.Unmarshal(data, event); err != nil {
				LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(ErrorEvent{ClientID: bot.ApplicationID, Event: FlagGatewayEventNameGuildRoleDelete, Err: err, Action: ErrorEventActionUnmarshal}).Msg("")
				return
			}

			for _, handler := range bot.Handlers.GuildRoleDelete {
				go handler(event)
			}
		}

	case FlagGatewayEventNameGuildScheduledEventCreate:
		if len(bot.Handlers.GuildScheduledEventCreate) != 0 {
			event := new(GuildScheduledEventCreate)
			if err := json.Unmarshal(data, event); err != nil {
				LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(ErrorEvent{ClientID: bot.ApplicationID, Event: FlagGatewayEventNameGuildScheduledEventCreate, Err: err, Action: ErrorEventActionUnmarshal}).Msg("")
				return
			}

			for _, handler := range bot.Handlers.GuildScheduledEventCreate {
				go handler(event)
			}
		}

	case FlagGatewayEventNameGuildScheduledEventUpdate:
		if len(bot.Handlers.GuildScheduledEventUpdate) != 0 {
			event := new(GuildScheduledEventUpdate)
			if err := json.Unmarshal(data, event); err != nil {
				LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(ErrorEvent{ClientID: bot.ApplicationID, Event: FlagGatewayEventNameGuildScheduledEventUpdate, Err: err, Action: ErrorEventActionUnmarshal}).Msg("")
				return
			}

			for _, handler := range bot.Handlers.GuildScheduledEventUpdate {
				go handler(event)
			}
		}

	case FlagGatewayEventNameGuildScheduledEventDelete:
		if len(bot.Handlers.GuildScheduledEventDelete) != 0 {
			event := new(GuildScheduledEventDelete)
			if err := json.Unmarshal(data, event); err != nil {
				LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(ErrorEvent{ClientID: bot.ApplicationID, Event: FlagGatewayEventNameGuildScheduledEventDelete, Err: err, Action: ErrorEventActionUnmarshal}).Msg("")
				return
			}

			for _, handler := range bot.Handlers.GuildScheduledEventDelete {
				go handler(event)
			}
		}

	case FlagGatewayEventNameGuildScheduledEventUserAdd:
		if len(bot.Handlers.GuildScheduledEventUserAdd) != 0 {
			event := new(GuildScheduledEventUserAdd)
			if err := json.Unmarshal(data, event); err != nil {
				LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(ErrorEvent{ClientID: bot.ApplicationID, Event: FlagGatewayEventNameGuildScheduledEventUserAdd, Err: err, Action: ErrorEventActionUnmarshal}).Msg("")
				return
			}

			for _, handler := range bot.Handlers.GuildScheduledEventUserAdd {
				go handler(event)
			}
		}

	case FlagGatewayEventNameGuildScheduledEventUserRemove:
		if len(bot.Handlers.GuildScheduledEventUserRemove) != 0 {
			event := new(GuildScheduledEventUserRemove)
			if err := json.Unmarshal(data, event); err != nil {
				LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(ErrorEvent{ClientID: bot.ApplicationID, Event: FlagGatewayEventNameGuildScheduledEventUserRemove, Err: err, Action: ErrorEventActionUnmarshal}).Msg("")
				return
			}

			for _, handler := range bot.Handlers.GuildScheduledEventUserRemove {
				go handler(event)
			}
		}

	case FlagGatewayEventNameIntegrationCreate:
		if len(bot.Handlers.IntegrationCreate) != 0 {
			event := new(IntegrationCreate)
			if err := json.Unmarshal(data, event); err != nil {
				LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(ErrorEvent{ClientID: bot.ApplicationID, Event: FlagGatewayEventNameIntegrationCreate, Err: err, Action: ErrorEventActionUnmarshal}).Msg("")
				return
			}

			for _, handler := range bot.Handlers.IntegrationCreate {
				go handler(event)
			}
		}

	case FlagGatewayEventNameIntegrationUpdate:
		if len(bot.Handlers.IntegrationUpdate) != 0 {
			event := new(IntegrationUpdate)
			if err := json.Unmarshal(data, event); err != nil {
				LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(ErrorEvent{ClientID: bot.ApplicationID, Event: FlagGatewayEventNameIntegrationUpdate, Err: err, Action: ErrorEventActionUnmarshal}).Msg("")
				return
			}

			for _, handler := range bot.Handlers.IntegrationUpdate {
				go handler(event)
			}
		}

	case FlagGatewayEventNameIntegrationDelete:
		if len(bot.Handlers.IntegrationDelete) != 0 {
			event := new(IntegrationDelete)
			if err := json.Unmarshal(data, event); err != nil {
				LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(ErrorEvent{ClientID: bot.ApplicationID, Event: FlagGatewayEventNameIntegrationDelete, Err: err, Action: ErrorEventActionUnmarshal}).Msg("")
				return
			}

			for _, handler := range bot.Handlers.IntegrationDelete {
				go handler(event)
			}
		}

	case FlagGatewayEventNameInviteCreate:
		if len(bot.Handlers.InviteCreate) != 0 {
			event := new(InviteCreate)
			if err := json.Unmarshal(data, event); err != nil {
				LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(ErrorEvent{ClientID: bot.ApplicationID, Event: FlagGatewayEventNameInviteCreate, Err: err, Action: ErrorEventActionUnmarshal}).Msg("")
				return
			}

			for _, handler := range bot.Handlers.InviteCreate {
				go handler(event)
			}
		}

	case FlagGatewayEventNameInviteDelete:
		if len(bot.Handlers.InviteDelete) != 0 {
			event := new(InviteDelete)
			if err := json.Unmarshal(data, event); err != nil {
				LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(ErrorEvent{ClientID: bot.ApplicationID, Event: FlagGatewayEventNameInviteDelete, Err: err, Action: ErrorEventActionUnmarshal}).Msg("")
				return
			}

			for _, handler := range bot.Handlers.InviteDelete {
				go handler(event)
			}
		}

	case FlagGatewayEventNameMessageCreate:
		if len(bot.Handlers.MessageCreate) != 0 {
			event := new(MessageCreate)
			if err := json.Unmarshal(data, event); err != nil {
				LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(ErrorEvent{ClientID: bot.ApplicationID, Event: FlagGatewayEventNameMessageCreate, Err: err, Action: ErrorEventActionUnmarshal}).Msg("")
				return
			}

			for _, handler := range bot.Handlers.MessageCreate {
				go handler(event)
			}
		}

	case FlagGatewayEventNameMessageUpdate:
		if len(bot.Handlers.MessageUpdate) != 0 {
			event := new(MessageUpdate)
			if err := json.Unmarshal(data, event); err != nil {
				LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(ErrorEvent{ClientID: bot.ApplicationID, Event: FlagGatewayEventNameMessageUpdate, Err: err, Action: ErrorEventActionUnmarshal}).Msg("")
				return
			}

			for _, handler := range bot.Handlers.MessageUpdate {
				go handler(event)
			}
		}

	case FlagGatewayEventNameMessageDelete:
		if len(bot.Handlers.MessageDelete) != 0 {
			event := new(MessageDelete)
			if err := json.Unmarshal(data, event); err != nil {
				LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(ErrorEvent{ClientID: bot.ApplicationID, Event: FlagGatewayEventNameMessageDelete, Err: err, Action: ErrorEventActionUnmarshal}).Msg("")
				return
			}

			for _, handler := range bot.Handlers.MessageDelete {
				go handler(event)
			}
		}

	case FlagGatewayEventNameMessageDeleteBulk:
		if len(bot.Handlers.MessageDeleteBulk) != 0 {
			event := new(MessageDeleteBulk)
			if err := json.Unmarshal(data, event); err != nil {
				LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(ErrorEvent{ClientID: bot.ApplicationID, Event: FlagGatewayEventNameMessageDeleteBulk, Err: err, Action: ErrorEventActionUnmarshal}).Msg("")
				return
			}

			for _, handler := range bot.Handlers.MessageDeleteBulk {
				go handler(event)
			}
		}

	case FlagGatewayEventNameMessageReactionAdd:
		if len(bot.Handlers.MessageReactionAdd) != 0 {
			event := new(MessageReactionAdd)
			if err := json.Unmarshal(data, event); err != nil {
				LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(ErrorEvent{ClientID: bot.ApplicationID, Event: FlagGatewayEventNameMessageReactionAdd, Err: err, Action: ErrorEventActionUnmarshal}).Msg("")
				return
			}

			for _, handler := range bot.Handlers.MessageReactionAdd {
				go handler(event)
			}
		}

	case FlagGatewayEventNameMessageReactionRemove:
		if len(bot.Handlers.MessageReactionRemove) != 0 {
			event := new(MessageReactionRemove)
			if err := json.Unmarshal(data, event); err != nil {
				LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(ErrorEvent{ClientID: bot.ApplicationID, Event: FlagGatewayEventNameMessageReactionRemove, Err: err, Action: ErrorEventActionUnmarshal}).Msg("")
				return
			}

			for _, handler := range bot.Handlers.MessageReactionRemove {
				go handler(event)
			}
		}

	case FlagGatewayEventNameMessageReactionRemoveAll:
		if len(bot.Handlers.MessageReactionRemoveAll) != 0 {
			event := new(MessageReactionRemoveAll)
			if err := json.Unmarshal(data, event); err != nil {
				LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(ErrorEvent{ClientID: bot.ApplicationID, Event: FlagGatewayEventNameMessageReactionRemoveAll, Err: err, Action: ErrorEventActionUnmarshal}).Msg("")
				return
			}

			for _, handler := range bot.Handlers.MessageReactionRemoveAll {
				go handler(event)
			}
		}

	case FlagGatewayEventNameMessageReactionRemoveEmoji:
		if len(bot.Handlers.MessageReactionRemoveEmoji) != 0 {
			event := new(MessageReactionRemoveEmoji)
			if err := json.Unmarshal(data, event); err != nil {
				LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(ErrorEvent{ClientID: bot.ApplicationID, Event: FlagGatewayEventNameMessageReactionRemoveEmoji, Err: err, Action: ErrorEventActionUnmarshal}).Msg("")
				return
			}

			for _, handler := range bot.Handlers.MessageReactionRemoveEmoji {
				go handler(event)
			}
		}

	case FlagGatewayEventNamePresenceUpdate:
		if len(bot.Handlers.PresenceUpdate) != 0 {
			event := new(PresenceUpdate)
			if err := json.Unmarshal(data, event); err != nil {
				LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(ErrorEvent{ClientID: bot.ApplicationID, Event: FlagGatewayEventNamePresenceUpdate, Err: err, Action: ErrorEventActionUnmarshal}).Msg("")
				return
			}

			for _, handler := range bot.Handlers.PresenceUpdate {
				go handler(event)
			}
		}

	case FlagGatewayEventNameStageInstanceCreate:
		if len(bot.Handlers.StageInstanceCreate) != 0 {
			event := new(StageInstanceCreate)
			if err := json.Unmarshal(data, event); err != nil {
				LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(ErrorEvent{ClientID: bot.ApplicationID, Event: FlagGatewayEventNameStageInstanceCreate, Err: err, Action: ErrorEventActionUnmarshal}).Msg("")
				return
			}

			for _, handler := range bot.Handlers.StageInstanceCreate {
				go handler(event)
			}
		}

	case FlagGatewayEventNameStageInstanceDelete:
		if len(bot.Handlers.StageInstanceDelete) != 0 {
			event := new(StageInstanceDelete)
			if err := json.Unmarshal(data, event); err != nil {
				LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(ErrorEvent{ClientID: bot.ApplicationID, Event: FlagGatewayEventNameStageInstanceDelete, Err: err, Action: ErrorEventActionUnmarshal}).Msg("")
				return
			}

			for _, handler := range bot.Handlers.StageInstanceDelete {
				go handler(event)
			}
		}

	case FlagGatewayEventNameStageInstanceUpdate:
		if len(bot.Handlers.StageInstanceUpdate) != 0 {
			event := new(StageInstanceUpdate)
			if err := json.Unmarshal(data, event); err != nil {
				LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(ErrorEvent{ClientID: bot.ApplicationID, Event: FlagGatewayEventNameStageInstanceUpdate, Err: err, Action: ErrorEventActionUnmarshal}).Msg("")
				return
			}

			for _, handler := range bot.Handlers.StageInstanceUpdate {
				go handler(event)
			}
		}

	case FlagGatewayEventNameTypingStart:
		if len(bot.Handlers.TypingStart) != 0 {
			event := new(TypingStart)
			if err := json.Unmarshal(data, event); err != nil {
				LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(ErrorEvent{ClientID: bot.ApplicationID, Event: FlagGatewayEventNameTypingStart, Err: err, Action: ErrorEventActionUnmarshal}).Msg("")
				return
			}

			for _, handler := range bot.Handlers.TypingStart {
				go handler(event)
			}
		}

	case FlagGatewayEventNameVoiceStateUpdate:
		if len(bot.Handlers.VoiceStateUpdate) != 0 {
			event := new(VoiceStateUpdate)
			if err := json.Unmarshal(data, event); err != nil {
				LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(ErrorEvent{ClientID: bot.ApplicationID, Event: FlagGatewayEventNameVoiceStateUpdate, Err: err, Action: ErrorEventActionUnmarshal}).Msg("")
				return
			}

			for _, handler := range bot.Handlers.VoiceStateUpdate {
				go handler(event)
			}
		}

	case FlagGatewayEventNameWebhooksUpdate:
		if len(bot.Handlers.WebhooksUpdate) != 0 {
			event := new(WebhooksUpdate)
			if err := json.Unmarshal(data, event); err != nil {
				LogEventHandler(Logger.Error(), bot.ApplicationID, eventname).Err(ErrorEvent{ClientID: bot.ApplicationID, Event: FlagGatewayEventNameWebhooksUpdate, Err: err, Action: ErrorEventActionUnmarshal}).Msg("")
				return
			}

			for _, handler := range bot.Handlers.WebhooksUpdate {
				go handler(event)
			}
		}
	}
}

func (r *BulkOverwriteGlobalApplicationCommands) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.ApplicationCommands) // nolint:wrapcheck
}

func (r *BulkOverwriteGuildApplicationCommands) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.ApplicationCommands) // nolint:wrapcheck
}

func (r *ModifyGuildChannelPositions) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Parameters) // nolint:wrapcheck
}

func (r *ModifyGuildRolePositions) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Parameters) // nolint:wrapcheck
}

/**unmarshal.go contains custom UnmarshalJSON() functions.

This enables json.Unmarshal() to unmarshal into types that contain fields that are interfaces.

In addition, structs that contain an embedded field - that implements UnmarshalJSON() - will
use the embedded field's implementation of UnmarshalJSON(). As a result, these structs must
also implement UnmarshalJSON() to prevent null pointer dereferences. */

/** Unused: Command, Event */

/** Nonce
Includes: CreateMessage, Message */

func (v *Nonce) UnmarshalJSON(b []byte) error {
	var x interface{}

	if err := json.Unmarshal(b, &x); err != nil {
		return fmt.Errorf(errUnmarshal, x, err)
	}

	switch xValue := x.(type) {
	case string:
		*v = Nonce(xValue)

	case int:
		*v = Nonce(strconv.Itoa(xValue))

	default:
		return fmt.Errorf(errUnmarshal, v, fmt.Errorf("value is type %T", x))
	}

	return nil
}

/** Value
Includes: ApplicationCommandOptionChoice, ApplicationCommandInteractionDataOption */

func (v *Value) UnmarshalJSON(b []byte) error {
	var x interface{}

	if err := json.Unmarshal(b, &x); err != nil {
		return fmt.Errorf(errUnmarshal, x, err)
	}

	switch xValue := x.(type) {
	case string:
		*v = Value(xValue)

	case int:
		*v = Value(strconv.Itoa(xValue))

	case float64:
		*v = Value(strconv.FormatFloat(xValue, 'f', -1, bit64))

	default:
		return fmt.Errorf(errUnmarshal, v, fmt.Errorf("value is type %T", x))
	}

	return nil
}

/** Component */

// unmarshalComponents unmarshals a JSON component array into a slice of Go Interface Components (with underlying structs).
func unmarshalComponents(b []byte) ([]Component, error) {
	if len(b) == 0 {
		return nil, nil
	}

	type unmarshalComponent struct {

		// https://discord.com/developers/docs/interactions/message-components#component-object-example-component
		Type Flag `json:"type"`
	}

	// Components are always provided in a JSON array.
	// Create a variable (of type []unmarshalComponent) that can read all of the Component Types.
	var unmarshalledComponents []unmarshalComponent

	// unmarshal the JSON (components.{component.Type}) into unmarshalledComponents.
	if err := json.Unmarshal(b, &unmarshalledComponents); err != nil {
		return nil, fmt.Errorf(errUnmarshal, unmarshalledComponents, err)
	}

	// use the known component types to return a slice of Go Interface Components with underlying structs.
	components := make([]Component, len(unmarshalledComponents))
	for i, unmarshalledComponent := range unmarshalledComponents {
		// set the component (interface) to an underlying type.
		switch unmarshalledComponent.Type {
		case FlagComponentTypeActionRow:
			components[i] = &ActionsRow{} //nolint:exhaustruct

		case FlagComponentTypeButton:
			components[i] = &Button{} //nolint:exhaustruct

		case FlagComponentTypeSelectMenu,
			FlagComponentTypeUserSelect,
			FlagComponentTypeRoleSelect,
			FlagComponentTypeMentionableSelect,
			FlagComponentTypeChannelSelect:
			components[i] = &SelectMenu{} //nolint:exhaustruct

		case FlagComponentTypeTextInput:
			components[i] = &TextInput{} //nolint:exhaustruct

		default:
			return nil, fmt.Errorf(
				"attempt to unmarshal into unknown component type (%d)",
				unmarshalledComponent.Type,
			)
		}
	}

	return components, nil
}

func (r *EditOriginalInteractionResponse) UnmarshalJSON(b []byte) error {
	// The following pattern is present throughout this file
	// in order to prevent a stack overflow (of r.UnmarshalJSON()).
	type alias EditOriginalInteractionResponse

	var unmarshalled struct {
		alias
		Components json.RawMessage `json:"components"`
	}

	if err := json.Unmarshal(b, &unmarshalled); err != nil {
		return fmt.Errorf(errUnmarshal, r, err)
	}

	components, err := unmarshalComponents(unmarshalled.Components)
	if err != nil {
		return fmt.Errorf(errUnmarshal, r, err)
	}

	unmarshalled.alias.Components = &components

	if r == nil {
		r = new(EditOriginalInteractionResponse)
	}

	*r = EditOriginalInteractionResponse(unmarshalled.alias)

	return nil
}

func (r *CreateFollowupMessage) UnmarshalJSON(b []byte) error {
	type alias CreateFollowupMessage

	var unmarshalled struct {
		alias
		Components json.RawMessage `json:"components"`
	}

	var err error
	if err = json.Unmarshal(b, &unmarshalled); err != nil {
		return fmt.Errorf(errUnmarshal, r, err)
	}

	if unmarshalled.alias.Components, err = unmarshalComponents(unmarshalled.Components); err != nil {
		return fmt.Errorf(errUnmarshal, r, err)
	}

	if r == nil {
		r = new(CreateFollowupMessage)
	}

	*r = CreateFollowupMessage(unmarshalled.alias)

	return nil
}

func (r *EditFollowupMessage) UnmarshalJSON(b []byte) error {
	type alias EditFollowupMessage

	var unmarshalled struct {
		alias
		Components json.RawMessage `json:"components"`
	}

	if err := json.Unmarshal(b, &unmarshalled); err != nil {
		return fmt.Errorf(errUnmarshal, r, err)
	}

	components, err := unmarshalComponents(unmarshalled.Components)
	if err != nil {
		return fmt.Errorf(errUnmarshal, r, err)
	}

	unmarshalled.alias.Components = &components

	if r == nil {
		r = new(EditFollowupMessage)
	}

	*r = EditFollowupMessage(unmarshalled.alias)

	return nil
}

func (r *EditMessage) UnmarshalJSON(b []byte) error {
	type alias EditMessage

	var unmarshalled struct {
		alias
		Components json.RawMessage `json:"components"`
	}

	if err := json.Unmarshal(b, &unmarshalled); err != nil {
		return fmt.Errorf(errUnmarshal, r, err)
	}

	components, err := unmarshalComponents(unmarshalled.Components)
	if err != nil {
		return fmt.Errorf(errUnmarshal, r, err)
	}

	unmarshalled.alias.Components = &components

	if r == nil {
		r = new(EditMessage)
	}

	*r = EditMessage(unmarshalled.alias)

	return nil
}

func (r *ForumThreadMessageParams) UnmarshalJSON(b []byte) error {
	type alias ForumThreadMessageParams

	var unmarshalled struct {
		alias
		Components json.RawMessage `json:"components"`
	}

	var err error
	if err = json.Unmarshal(b, &unmarshalled); err != nil {
		return fmt.Errorf(errUnmarshal, r, err)
	}

	if unmarshalled.alias.Components, err = unmarshalComponents(unmarshalled.Components); err != nil {
		return fmt.Errorf(errUnmarshal, r, err)
	}

	if r == nil {
		r = new(ForumThreadMessageParams)
	}

	*r = ForumThreadMessageParams(unmarshalled.alias)

	return nil
}

func (r *ExecuteWebhook) UnmarshalJSON(b []byte) error {
	type alias ExecuteWebhook

	var unmarshalled struct {
		alias
		Components json.RawMessage `json:"components"`
	}

	var err error
	if err = json.Unmarshal(b, &unmarshalled); err != nil {
		return fmt.Errorf(errUnmarshal, r, err)
	}

	if unmarshalled.alias.Components, err = unmarshalComponents(unmarshalled.Components); err != nil {
		return fmt.Errorf(errUnmarshal, r, err)
	}

	if r == nil {
		r = new(ExecuteWebhook)
	}

	*r = ExecuteWebhook(unmarshalled.alias)

	return nil
}

func (r *EditWebhookMessage) UnmarshalJSON(b []byte) error {
	type alias EditWebhookMessage

	var unmarshalled struct {
		alias
		Components json.RawMessage `json:"components"`
	}

	if err := json.Unmarshal(b, &unmarshalled); err != nil {
		return fmt.Errorf(errUnmarshal, r, err)
	}

	components, err := unmarshalComponents(unmarshalled.Components)
	if err != nil {
		return fmt.Errorf(errUnmarshal, r, err)
	}

	unmarshalled.alias.Components = &components

	if r == nil {
		r = new(EditWebhookMessage)
	}

	*r = EditWebhookMessage(unmarshalled.alias)

	return nil
}

func (r *ActionsRow) UnmarshalJSON(b []byte) error {
	type alias ActionsRow

	var unmarshalled struct {
		alias
		Components json.RawMessage `json:"components"`
	}

	var err error
	if err = json.Unmarshal(b, &unmarshalled); err != nil {
		return fmt.Errorf(errUnmarshal, r, err)
	}

	if unmarshalled.alias.Components, err = unmarshalComponents(unmarshalled.Components); err != nil {
		return fmt.Errorf(errUnmarshal, r, err)
	}

	if r == nil {
		r = new(ActionsRow)
	}

	*r = ActionsRow(unmarshalled.alias)

	return nil
}

func (r *ModalSubmitData) UnmarshalJSON(b []byte) error {
	type alias ModalSubmitData

	var unmarshalled struct {
		alias
		Components json.RawMessage `json:"components"`
	}

	var err error
	if err = json.Unmarshal(b, &unmarshalled); err != nil {
		return fmt.Errorf(errUnmarshal, r, err)
	}

	if unmarshalled.alias.Components, err = unmarshalComponents(unmarshalled.Components); err != nil {
		return fmt.Errorf(errUnmarshal, r, err)
	}

	if r == nil {
		r = new(ModalSubmitData)
	}

	*r = ModalSubmitData(unmarshalled.alias)

	return nil
}

func (r *Messages) UnmarshalJSON(b []byte) error {
	type alias Messages

	var unmarshalled struct {
		alias
		Components json.RawMessage `json:"components"`
	}

	var err error
	if err = json.Unmarshal(b, &unmarshalled); err != nil {
		return fmt.Errorf(errUnmarshal, r, err)
	}

	if unmarshalled.alias.Components, err = unmarshalComponents(unmarshalled.Components); err != nil {
		return fmt.Errorf(errUnmarshal, r, err)
	}

	if r == nil {
		r = new(Messages)
	}

	*r = Messages(unmarshalled.alias)

	return nil
}

func (r *Modal) UnmarshalJSON(b []byte) error {
	type alias Modal

	var unmarshalled struct {
		alias
		Components json.RawMessage `json:"components"`
	}

	var err error
	if err = json.Unmarshal(b, &unmarshalled); err != nil {
		return fmt.Errorf(errUnmarshal, r, err)
	}

	if unmarshalled.alias.Components, err = unmarshalComponents(unmarshalled.Components); err != nil {
		return fmt.Errorf(errUnmarshal, r, err)
	}

	if r == nil {
		r = new(Modal)
	}

	*r = Modal(unmarshalled.alias)

	return nil
}

func (r *Message) UnmarshalJSON(b []byte) error {
	type alias Message

	var unmarshalled struct {
		alias
		Components json.RawMessage `json:"components"`
	}

	var err error
	if err = json.Unmarshal(b, &unmarshalled); err != nil {
		return fmt.Errorf(errUnmarshal, r, err)
	}

	if unmarshalled.alias.Components, err = unmarshalComponents(unmarshalled.Components); err != nil {
		return fmt.Errorf(errUnmarshal, r, err)
	}

	if r == nil {
		r = new(Message)
	}

	*r = Message(unmarshalled.alias)

	return nil
}

/** InteractionData */

// unmarshalInteractionData unmarshals a JSON InteractionData object into
// a Go Interface InteractionData (with an underlying struct).
func unmarshalInteractionData(b json.RawMessage, x Flag) (InteractionData, error) {
	if len(b) == 0 {
		return nil, nil
	}

	var interactionData InteractionData

	// use the known Interaction Data type to return
	// a Go Interface InteractionData with an underlying struct.
	switch x {
	case FlagInteractionTypePING:
		return nil, nil

	case FlagInteractionTypeAPPLICATION_COMMAND,
		FlagInteractionTypeAPPLICATION_COMMAND_AUTOCOMPLETE:
		interactionData = &ApplicationCommandData{} //nolint:exhaustruct

	case FlagInteractionTypeMESSAGE_COMPONENT:
		interactionData = &MessageComponentData{} //nolint:exhaustruct

	case FlagInteractionTypeMODAL_SUBMIT:
		interactionData = &ModalSubmitData{} //nolint:exhaustruct
	}

	if interactionData == nil {
		return nil, fmt.Errorf("attempt to unmarshal into unknown interaction data type (%d)", x)
	}

	// unmarshal into the underlying struct.
	if err := json.Unmarshal(b, interactionData); err != nil {
		return nil, fmt.Errorf(errUnmarshal, interactionData, err)
	}

	return interactionData, nil
}

func (r *Interaction) UnmarshalJSON(b []byte) error {
	type alias Interaction

	var unmarshalledInteraction struct {
		alias
		Data json.RawMessage `json:"data,omitempty"`
	}

	var err error
	if err = json.Unmarshal(b, &unmarshalledInteraction); err != nil {
		return fmt.Errorf(errUnmarshal, r, err)
	}

	if unmarshalledInteraction.alias.Data, err =
		unmarshalInteractionData(unmarshalledInteraction.Data, unmarshalledInteraction.Type); err != nil {
		return fmt.Errorf(errUnmarshal, r, err)
	}

	if r == nil {
		r = new(Interaction)
	}

	*r = Interaction(unmarshalledInteraction.alias)

	return nil
}

/** InteractionCallbackData */

// unmarshalInteractionCallbackData unmarshals a JSON InteractionCallbackData object into
// a Go Interface InteractionCallbackData (with an underlying struct).
func unmarshalInteractionCallbackData(b []byte, x Flag) (InteractionCallbackData, error) {
	if len(b) == 0 {
		return nil, nil
	}

	var interactionCallbackData InteractionCallbackData

	// use the known Interaction Callback Data type to return
	// a Go Interface InteractionCallbackData with an underlying struct.
	switch x {
	case FlagInteractionCallbackTypePONG:
		return nil, nil // Ping

	case FlagInteractionCallbackTypeCHANNEL_MESSAGE_WITH_SOURCE,
		FlagInteractionCallbackTypeUPDATE_MESSAGE:
		interactionCallbackData = &Messages{} //nolint:exhaustruct

	case FlagInteractionCallbackTypeDEFERRED_CHANNEL_MESSAGE_WITH_SOURCE:
		return nil, nil // Edit a followup response later.

	case FlagInteractionCallbackTypeDEFERRED_UPDATE_MESSAGE:
		return nil, nil // Edit the original response later.

	case FlagInteractionCallbackTypeAPPLICATION_COMMAND_AUTOCOMPLETE_RESULT:
		interactionCallbackData = &Autocomplete{} //nolint:exhaustruct

	case FlagInteractionCallbackTypeMODAL:
		interactionCallbackData = &Modal{} //nolint:exhaustruct
	}

	if interactionCallbackData == nil {
		return nil, fmt.Errorf(
			"attempt to unmarshal into unknown interaction callback data type (%d)",
			x)
	}

	// unmarshal into the underlying struct.
	if err := json.Unmarshal(b, interactionCallbackData); err != nil {
		return nil, fmt.Errorf(errUnmarshal, interactionCallbackData, err)
	}

	return interactionCallbackData, nil
}

func (r *InteractionResponse) UnmarshalJSON(b []byte) error {
	type alias InteractionResponse

	var unmarshalledInteractionResponse struct {
		alias
		Data json.RawMessage `json:"data,omitempty"`
	}

	var err error
	if err = json.Unmarshal(b, &unmarshalledInteractionResponse); err != nil {
		return fmt.Errorf(errUnmarshal, r, err)
	}

	if unmarshalledInteractionResponse.alias.Data, err =
		unmarshalInteractionCallbackData(
			unmarshalledInteractionResponse.Data, unmarshalledInteractionResponse.Type); err != nil {
		return fmt.Errorf(errUnmarshal, r, err)
	}

	if r == nil {
		r = new(InteractionResponse)
	}

	*r = InteractionResponse(unmarshalledInteractionResponse.alias)

	return nil
}

/** Structs that contain embedded fields that implement UnmarshalJSON() */

func (e *MessageCreate) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &e.Message); err != nil {
		return fmt.Errorf(errUnmarshal, e, err)
	}

	return nil
}

func (e *MessageUpdate) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &e.Message); err != nil {
		return fmt.Errorf(errUnmarshal, e, err)
	}

	return nil
}

func (e *InteractionCreate) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &e.Interaction); err != nil {
		return fmt.Errorf(errUnmarshal, e, err)
	}

	return nil
}

func (e *CreateInteractionResponse) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &e.InteractionResponse); err != nil {
		return fmt.Errorf(errUnmarshal, e, err)
	}

	return nil
}

/**unmarshal_convert.go contains type conversion functions for interfaces.

This enables users (developers) to easily type convert interfaces. */

const (
	errTypeConvert = "attempted to type convert InteractionData of type %v to type %s"
)

/* Nonce */

func (n Nonce) String() string {
	return string(n)
}

func (n Nonce) Int64() (int64, error) {
	return strconv.ParseInt(string(n), base10, bit64) //nolint:wrapcheck
}

/* Value */

func (n Value) String() string {
	return string(n)
}

func (n Value) Float64() (float64, error) {
	return strconv.ParseFloat(string(n), bit64) //nolint:wrapcheck
}

func (n Value) Int64() (int64, error) {
	return strconv.ParseInt(string(n), base10, bit64) //nolint:wrapcheck
}

/* InteractionData */

// ApplicationCommand type converts an InteractionData field into an ApplicationCommandData struct.
func (i *Interaction) ApplicationCommand() *ApplicationCommandData {
	switch i.Data.InteractionDataType() {
	case FlagInteractionTypeAPPLICATION_COMMAND,
		FlagInteractionTypeAPPLICATION_COMMAND_AUTOCOMPLETE:
		return i.Data.(*ApplicationCommandData) //nolint:forcetypeassert

	case FlagInteractionTypePING:
		panic(fmt.Sprintf(errTypeConvert, "Ping", "ApplicationCommandData"))

	case FlagInteractionTypeMESSAGE_COMPONENT:
		panic(fmt.Sprintf(errTypeConvert, "MessageComponentData", "ApplicationCommandData"))

	case FlagInteractionTypeMODAL_SUBMIT:
		panic(fmt.Sprintf(errTypeConvert, "ModalSubmitData", "ApplicationCommandData"))
	}

	panic(fmt.Sprintf(errTypeConvert, i.Data.InteractionDataType(), "ApplicationCommandData"))
}

// MessageComponent type converts an InteractionData field into a MessageComponentData struct.
func (i *Interaction) MessageComponent() *MessageComponentData {
	switch i.Data.InteractionDataType() {
	case FlagInteractionTypeMESSAGE_COMPONENT:
		return i.Data.(*MessageComponentData) //nolint:forcetypeassert

	case FlagInteractionTypePING:
		panic(fmt.Sprintf(errTypeConvert, "Ping", "MessageComponentData"))

	case FlagInteractionTypeAPPLICATION_COMMAND,
		FlagInteractionTypeAPPLICATION_COMMAND_AUTOCOMPLETE:
		panic(fmt.Sprintf(errTypeConvert, "ApplicationCommandData", "MessageComponentData"))

	case FlagInteractionTypeMODAL_SUBMIT:
		panic(fmt.Sprintf(errTypeConvert, "ModalSubmitData", "MessageComponentData"))
	}

	panic(fmt.Sprintf(errTypeConvert, i.Data.InteractionDataType(), "MessageComponentData"))
}

// ModalSubmit type converts an InteractionData field into a ModalSubmitData struct.
func (i *Interaction) ModalSubmit() *ModalSubmitData {
	switch i.Data.InteractionDataType() {
	case FlagInteractionTypeMODAL_SUBMIT:
		return i.Data.(*ModalSubmitData) //nolint:forcetypeassert

	case FlagInteractionTypePING:
		panic(fmt.Sprintf(errTypeConvert, "Ping", "ModalSubmitData"))

	case FlagInteractionTypeAPPLICATION_COMMAND,
		FlagInteractionTypeAPPLICATION_COMMAND_AUTOCOMPLETE:
		panic(fmt.Sprintf(errTypeConvert, "ApplicationCommandData", "ModalSubmitData"))

	case FlagInteractionTypeMESSAGE_COMPONENT:
		panic(fmt.Sprintf(errTypeConvert, "MessageComponentData", "ModalSubmitData"))
	}

	panic(fmt.Sprintf(errTypeConvert, i.Data.InteractionDataType(), "ModalSubmitData"))
}

// init is called at the start of the application.
func init() {
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.SetGlobalLevel(zerolog.Disabled)
}

var (
	// Logger represents the Disgo Logger used to log information.
	Logger = zerolog.New(os.Stdout)
)

// Logger Contexts
const (
	// LogCtxClient represents the log key for a Client (Bot) Application ID.
	LogCtxClient = "client"

	// LogCtxCorrelation represents the log key for a Correlation ID.
	LogCtxCorrelation = "xid"

	// LogCtxRequest represents the log key for a Request ID.
	LogCtxRequest = "request"

	// LogCtxRoute represents the log key for a Route ID.
	LogCtxRoute = "route"

	// LogCtxResource represents the log key for a Resource ID.
	LogCtxResource = "resource"

	// LogCtxEndpoint represents the log key for an Endpoint.
	LogCtxEndpoint = "endpoint"

	// LogCtxRequestBody represents the log key for an HTTP Request Body.
	LogCtxRequestBody = "body"

	// LogCtxBucket represents the log key for a Rate Limit Bucket ID.
	LogCtxBucket = "bucket"

	// LogCtxReset represents the log key for a Discord Bucket reset time.
	LogCtxReset = "reset"

	// LogCtxResponse represents the log key for an HTTP Request Response.
	LogCtxResponse = "response"

	// LogCtxResponseHeader represents the log key for an HTTP Request Response header.
	LogCtxResponseHeader = "header"

	// LogCtxResponseBody represents the log key for an HTTP Request Response body.
	LogCtxResponseBody = "body"

	// LogCtxSession represents the log key for a Discord Session ID.
	LogCtxSession = "session"

	// LogCtxPayload represents the log key for a Discord Gateway Payload.
	LogCtxPayload = "payload"

	// LogCtxPayloadOpcode represents the log key for a Discord Gateway Payload opcode.
	LogCtxPayloadOpcode = "opcode"

	// LogCtxPayloadData represents the log key for Discord Gateway Payload data.
	LogCtxPayloadData = "data"

	// LogCtxEvent represents the log key for a Discord Gateway Event.
	LogCtxEvent = "event"

	// LogCtxCommand represents the log key for a Discord Gateway command.
	LogCtxCommand = "command"

	// LogCtxCommandOpcode represents the log key for a Discord Gateway command opcode.
	LogCtxCommandOpcode = "opcode"

	// LogCtxCommandName represents the log key for a Discord Gateway command name.
	LogCtxCommandName = "name"
)

// LogRequest logs a request.
func LogRequest(log *zerolog.Event, clientid, xid, routeid, resourceid, endpoint string) *zerolog.Event {
	return log.Timestamp().
		Str(LogCtxClient, clientid).
		Dict(LogCtxRequest, zerolog.Dict().
			Str(LogCtxCorrelation, xid).
			Str(LogCtxRoute, routeid).
			Str(LogCtxResource, resourceid).
			Str(LogCtxEndpoint, endpoint),
		)
}

// LogRequestBody logs a request with its body.
func LogRequestBody(log *zerolog.Event, clientid, xid, routeid, resourceid, endpoint, body string) *zerolog.Event {
	return log.Timestamp().
		Str(LogCtxClient, clientid).
		Dict(LogCtxRequest, zerolog.Dict().
			Str(LogCtxCorrelation, xid).
			Str(LogCtxRoute, routeid).
			Str(LogCtxResource, resourceid).
			Str(LogCtxEndpoint, endpoint).
			Str(LogCtxRequestBody, body),
		)
}

// LogResponse logs a response (typically using LogRequest).
func LogResponse(log *zerolog.Event, header, body string) *zerolog.Event {
	return log.Dict(LogCtxResponse, zerolog.Dict().
		Str(LogCtxResponseHeader, header).
		Str(LogCtxResponseBody, body),
	)
}

// LogEventHandler logs an event handler action.
func LogEventHandler(log *zerolog.Event, clientid, event string) *zerolog.Event {
	return log.Timestamp().
		Str(LogCtxClient, clientid).
		Str(LogCtxEvent, event)
}

// LogSession logs a session.
func LogSession(log *zerolog.Event, sessionid string) *zerolog.Event {
	return log.Timestamp().
		Str(LogCtxSession, sessionid)
}

// LogPayload logs a Discord Gateway Payload (typically using LogSession).
func LogPayload(log *zerolog.Event, op int, data json.RawMessage) *zerolog.Event {
	return log.Dict(LogCtxPayload, zerolog.Dict().
		Int(LogCtxPayloadOpcode, op).
		Bytes(LogCtxPayloadData, data),
	)
}

// LogCommand logs a Gateway Command (typically using a LogSession).
func LogCommand(log *zerolog.Event, clientid string, op int, command string) *zerolog.Event {
	return log.Str(LogCtxClient, clientid).
		Dict(LogCtxCommand, zerolog.Dict().
			Int(LogCtxCommandOpcode, op).
			Str(LogCtxCommandName, command),
		)
}

const (
	grantTypeAuthorizationCodeGrant = "authorization_code"
	grantTypeRefreshToken           = "refresh_token"
	grantTypeClientCredentials      = "client_credentials"
	amountAuthURLParams             = 6
	amountAuthURLParamsBot          = 3
)

// GenerateAuthorizationURL generates an authorization URL from a given client and response type.
func GenerateAuthorizationURL(bot *Client, response string) string {
	params := make([]string, 0, amountAuthURLParams)

	// response_type is the type of response the redirect will return.
	if response != "" {
		params = append(params, "responsetype="+response)
	}

	// client_id is the application client id.
	params = append(params, "client_id="+bot.Authorization.ClientID)

	// scope is a list of OAuth2 scopes separated by url encoded spaces (%20).
	scope := urlQueryStringScope(bot.Authorization.Scopes)
	if scope != "" {
		params = append(params, scope)
	}

	// redirect_uri is the URL registered while creating the application.
	if bot.Authorization.RedirectURI != "" {
		params = append(params, "redirect_uri="+url.QueryEscape(bot.Authorization.RedirectURI))
	}

	// state is the unique string mentioned in State and Security.
	if bot.Authorization.State != "" {
		params = append(params, "state="+bot.Authorization.State)
	}

	// prompt controls how the authorization flow handles existing authorizations.
	if bot.Authorization.Prompt != "" {
		params = append(params, "prompt="+bot.Authorization.Prompt)
	}

	return EndpointAuthorizationURL() + "?" + strings.Join(params, "&")
}

// BotAuthParams represents parameters used to generate a bot authorization URL.
type BotAuthParams struct {
	Bot                *Client
	GuildID            string
	ResponseType       string
	Permissions        BitFlag
	DisableGuildSelect bool
}

// GenerateBotAuthorizationURL generates a bot authorization URL using the given BotAuthParams.
//
// Bot.Scopes must include "bot" to enable the OAuth2 Bot Flow.
func GenerateBotAuthorizationURL(p BotAuthParams) string {
	params := make([]string, 0, amountAuthURLParamsBot)

	// permissions is permissions the bot is requesting.
	params = append(params, "permissions="+strconv.FormatUint(uint64(p.Permissions), base10))

	// guild_id is the Guild ID of the guild that is pre-selected in the authorization prompt.
	if p.GuildID != "" {
		params = append(params, "guild_id="+p.GuildID)
	}

	// disable_guild_select determines whether the user will be allowed to select a guild
	// other than the guild_id.
	params = append(params, "disable_guild_select="+strconv.FormatBool(p.DisableGuildSelect))

	return GenerateAuthorizationURL(p.Bot, p.ResponseType) + "&" + strings.Join(params, "&")
}

// AuthorizationCodeGrant performs an OAuth2 authorization code grant.
//
// Send the user a valid Authorization URL, which can be generated using
// GenerateAuthorizationURL(bot, "code").
//
// When the user visits the Authorization URL, they will be prompted for authorization.
// If the user accepts the prompt, they will be redirected to the `redirect_uri`.
// This issues a GET request to the `redirect_uri` web server which YOU MUST HANDLE
// by parsing the request's URL Query String into a disgo.RedirectURL object.
//
// Retrieve the user's access token by calling THIS FUNCTION (with the disgo.RedirectURL parameter),
// which performs an Access Token Exchange.
//
// Refresh the token by using RefreshAuthorizationCodeGrant(bot, token).
//
// For more information, read https://discord.com/developers/docs/topics/oauth2#authorization-code-grant
func AuthorizationCodeGrant(bot *Client, ru *RedirectURL) (*AccessTokenResponse, error) {
	exchange := &AccessTokenExchange{
		ClientID:     bot.Authorization.ClientID,
		ClientSecret: bot.Authorization.ClientSecret,
		GrantType:    grantTypeAuthorizationCodeGrant,
		Code:         ru.Code,
		RedirectURI:  bot.Authorization.RedirectURI,
	}

	return exchange.Send(bot)
}

// RefreshAuthorizationCodeGrant refreshes an Access Token from an OAuth2 authorization code grant.
func RefreshAuthorizationCodeGrant(bot *Client, token *AccessTokenResponse) (*AccessTokenResponse, error) {
	exchange := &RefreshTokenExchange{
		ClientID:     bot.Authorization.ClientID,
		ClientSecret: bot.Authorization.ClientSecret,
		GrantType:    grantTypeRefreshToken,
		RefreshToken: token.RefreshToken,
	}

	return exchange.Send(bot)
}

// ImplicitGrant converts a RedirectURI (from a simplified OAuth2 grant) to an AccessTokenResponse.
//
// Send the user a valid Authorization URL, which can be generated using
// GenerateAuthorizationURL(bot, "token").
//
// When the user visits the Authorization URL, they will be prompted for authorization.
// If the user accepts the prompt, they will be redirected to the `redirect_uri`.
// This issues a GET request to the `redirect_uri` web server which YOU MUST HANDLE
// by parsing the request's URI Fragments into a disgo.RedirectURI object.
//
// A disgo.RedirectURI object is equivalent to a disgo.AccessTokenResponse,
// but it does NOT contain a refresh token.
//
// For more information, read https://discord.com/developers/docs/topics/oauth2#implicit-grant
func ImplicitGrant(ru *RedirectURI) *AccessTokenResponse {
	return &AccessTokenResponse{
		AccessToken:  ru.AccessToken,
		TokenType:    ru.TokenType,
		ExpiresIn:    ru.ExpiresIn,
		RefreshToken: "",
		Scope:        ru.Scope,
	}
}

// ClientCredentialsGrant performs a client credential OAuth2 grant for TESTING PURPOSES.
//
// The bot client's Authentication Header will be set to a Basic Authentication Header that
// uses the bot's ClientID as a username and ClientSecret as a password.
//
// A request will be made for a Client Credential grant which returns a disgo.AccessTokenResponse
// that does NOT contain a refresh token.
//
// For more information, read https://discord.com/developers/docs/topics/oauth2#client-credentials-grant
func ClientCredentialsGrant(bot *Client) (*AccessTokenResponse, error) {
	bot.Authentication.Header = "Basic " +
		base64.StdEncoding.EncodeToString([]byte(bot.Authorization.ClientID+":"+bot.Authorization.ClientSecret))

	grant := &ClientCredentialsTokenRequest{
		GrantType: grantTypeClientCredentials,
		Scope:     urlQueryStringScope(bot.Authorization.Scopes),
	}

	return grant.Send(bot)
}

// BotAuthorization performs a specialized OAuth2 flow for users to add bots to guilds.
//
// Send the user a valid Bot Authorization URL, which can be generated using
// GenerateBotAuthorizationURL(disgo.BotAuthParams{...}).
//
// When the user visits the Bot Authorization URL, they will be prompted for authorization.
// If the user accepts the prompt (with a guild), the bot will be added to the selected guild.
//
// For more information read, https://discord.com/developers/docs/topics/oauth2#bot-authorization-flow
func BotAuthorization() {}

// AdvancedBotAuthorization performs a specialized OAuth2 flow for users to add bots to guilds.
//
// Send the user a valid Bot Authorization URL, which can be generated using
// GenerateBotAuthorizationURL(disgo.BotAuthParams{...}).
//
// If the user accepts the prompt (with a guild), they will be redirected to the `redirect_uri`.
// This issues a GET request to the `redirect_uri` web server which YOU MUST HANDLE
// by parsing the request's URL Query String into a disgo.RedirectURL object.
//
// Retrieve the user's access token by calling THIS FUNCTION (with the disgo.RedirectURL parameter),
// which performs an Access Token Exchange.
//
// Refresh the token by using RefreshAuthorizationCodeGrant(bot, token).
//
// For more information read, https://discord.com/developers/docs/topics/oauth2#advanced-bot-authorization
func AdvancedBotAuthorization(bot *Client, ru *RedirectURL) (*AccessTokenResponse, error) {
	return AuthorizationCodeGrant(bot, ru)
}

// WebhookAuthorization performs a specialized OAuth2 authorization code grant.
//
// Send the user a valid Authorization URL, which can be generated using
// GenerateAuthorizationURL(bot, "code") when bot.Scopes is set to `webhook.incoming`.
//
// When the user visits the Authorization URL, they will be prompted for authorization.
// If the user accepts the prompt (with a channel), they will be redirected to the `redirect_uri`.
// This issues a GET request to the `redirect_uri` web server which YOU MUST HANDLE
// by parsing the request's URL Query String into a disgo.RedirectURL object.
//
// Retrieve the user's access token by calling THIS FUNCTION (with the disgo.RedirectURL parameter),
// which performs an Access Token Exchange.
//
// Refresh the token by using RefreshAuthorizationCodeGrant(bot, token).
//
// For more information read, https://discord.com/developers/docs/topics/oauth2#webhooks
func WebhookAuthorization(bot *Client, ru *RedirectURL) (*AccessTokenResponse, *Webhook, error) {
	exchange := &AccessTokenExchange{
		ClientID:     bot.Authorization.ClientID,
		ClientSecret: bot.Authorization.ClientSecret,
		GrantType:    grantTypeAuthorizationCodeGrant,
		Code:         ru.Code,
		RedirectURI:  bot.Authorization.RedirectURI,
	}

	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[1]("1")
	query, err := EndpointQueryString(exchange)
	if err != nil {
		return nil, nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      "",
			Err:           err,
		}
	}
	endpoint := EndpointTokenURL() + "?" + query

	result := new(WebhookTokenResponse)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPost, endpoint, ContentTypeURLQueryString, nil, result)
	if err != nil {
		return nil, nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	// convert the webhook token response to an access token response (and webhook).
	token := &AccessTokenResponse{
		AccessToken:  result.AccessToken,
		TokenType:    result.TokenType,
		ExpiresIn:    result.ExpiresIn,
		RefreshToken: result.RefreshToken,
		Scope:        result.Scope,
	}

	return token, result.Webhook, nil
}

// Send sends an AccessTokenExchange request to Discord and returns an AccessTokenResponse.
func (r *AccessTokenExchange) Send(bot *Client) (*AccessTokenResponse, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[1]("1")
	query, err := EndpointQueryString(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      "",
			Err:           err,
		}
	}
	endpoint := EndpointTokenURL() + "?" + query

	result := new(AccessTokenResponse)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPost, endpoint, ContentTypeURLQueryString, nil, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a RefreshTokenExchange request to Discord and returns an AccessTokenResponse.
//
// Uses the RefreshTokenExchange ClientID and ClientSecret.
func (r *RefreshTokenExchange) Send(bot *Client) (*AccessTokenResponse, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[1]("1")
	query, err := EndpointQueryString(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      "",
			Err:           err,
		}
	}
	endpoint := EndpointTokenURL() + "?" + query

	result := new(AccessTokenResponse)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPost, endpoint, ContentTypeURLQueryString, nil, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a ClientCredentialsTokenRequest to Discord and returns a ClientCredentialsTokenRequest.
func (r *ClientCredentialsTokenRequest) Send(bot *Client) (*AccessTokenResponse, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[1]("1")
	query, err := EndpointQueryString(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      "",
			Err:           err,
		}
	}
	endpoint := EndpointTokenURL() + "?" + query

	result := new(AccessTokenResponse)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPost, endpoint, ContentTypeURLQueryString, nil, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// urlQueryStringScope parses a given slice of scopes to generate a valid URL Query String.
func urlQueryStringScope(scopes []string) string {
	if len(scopes) > 0 {
		var scope strings.Builder
		scope.WriteString("scope=")

		for i, s := range scopes {
			if i > 0 {
				scope.WriteString("%20")
			}

			scope.WriteString(s)
		}

		return scope.String()
	}

	return ""
}

// rlbpool represents a synchronized Rate Limit Bucket pool.
var rlbpool sync.Pool

// getBucket gets a Bucket from a pool.
func getBucket() *Bucket {
	if b := rlbpool.Get(); b != nil {
		return b.(*Bucket) //nolint:forcetypeassert
	}

	return new(Bucket)
}

// putBucket puts a Rate Limit Bucket into the pool.
func putBucket(b *Bucket) {
	b.ID = ""
	b.Limit = 0
	b.Remaining = 0
	b.Pending = 0
	b.Date = time.Time{}
	b.Expiry = time.Time{}

	rlbpool.Put(b)
}

// spool represents a synchronized Session pool.
var spool sync.Pool

// NewSession gets a Session from a pool.
func NewSession() *Session {
	if s := spool.Get(); s != nil {
		return s.(*Session) //nolint:forcetypeassert
	}

	return new(Session)
}

// putSession puts a Session into the pool.
func putSession(s *Session) {
	s.Lock()
	defer s.Unlock()

	// reset the Session.
	s.ID = ""
	s.Seq = 0
	s.Endpoint = ""
	s.Context = nil
	s.Conn = nil
	s.heartbeat = nil
	s.manager = nil

	spool.Put(s)
}

// gpool represents a synchronized Gateway Payload pool.
var gpool sync.Pool

// getPayload gets a Gateway Payload from the pool.
func getPayload() *GatewayPayload {
	if g := gpool.Get(); g != nil {
		return g.(*GatewayPayload) //nolint:forcetypeassert
	}

	return new(GatewayPayload)
}

// putPayload puts a Gateway Payload into the pool.
func putPayload(g *GatewayPayload) {
	// reset the Gateway Payload.
	g.Op = 0
	g.Data = nil
	g.SequenceNumber = nil
	g.EventName = nil

	gpool.Put(g)
}

const (
	nilRouteBucket = "NIL"
)

var (
	// IgnoreGlobalRateLimitRouteIDs represents a set of Route IDs that do NOT adhere to the Global Rate Limit.
	//
	// Interaction endpoints are not bound to the bot's Global Rate Limit.
	// https://discord.com/developers/docs/interactions/receiving-and-responding#endpoints
	IgnoreGlobalRateLimitRouteIDs = map[string]bool{
		"18": true, "19": true, "20": true, "21": true, "22": true, "23": true, "24": true, "25": true,
	}
)

// RateLimit provides concurrency-safe rate limit functionality by implementing the RateLimiter interface.
type RateLimit struct {
	ids           map[string]string
	buckets       map[string]*Bucket
	entries       map[string]int
	DefaultBucket *Bucket
	muQueue       sync.Mutex
	muTx          sync.Mutex
}

func (r *RateLimit) SetBucketID(routeid string, bucketid string) {
	currentBucketID := r.ids[routeid]

	// when the current Bucket ID is not the same as the new Bucket ID.
	if currentBucketID != bucketid {
		// update the entries for the current Bucket ID.
		if currentBucketID != "" {
			r.entries[currentBucketID]--

			// when the current Bucket ID is no longer referenced by a Route,
			// delete the respective Bucket (and recycle it).
			if r.entries[currentBucketID] <= 0 {
				if currentBucket := r.buckets[currentBucketID]; currentBucket != nil {
					putBucket(currentBucket)
				}
				delete(r.entries, currentBucketID)
				delete(r.buckets, currentBucketID)

				Logger.Info().Timestamp().Str(LogCtxRequest, routeid).Str(LogCtxBucket, currentBucketID).Msg("deleted bucket")
			}
		}

		// set the Route ID to the new Bucket ID.
		r.ids[routeid] = bucketid

		// update the entries for the new Bucket ID.
		r.entries[bucketid]++

		Logger.Info().Timestamp().Str(LogCtxRequest, routeid).Str(LogCtxBucket, bucketid).Msg("set route to bucket")
	}
}

func (r *RateLimit) GetBucketID(routeid string) string {
	return r.ids[routeid]
}

func (r *RateLimit) SetBucketFromID(bucketid string, bucket *Bucket) {
	r.buckets[bucketid] = bucket

	Logger.Info().Timestamp().Str(LogCtxBucket, bucketid).Msgf("set bucket to object %p", bucket)
}

func (r *RateLimit) GetBucketFromID(bucketid string) *Bucket {
	return r.buckets[bucketid]
}

func (r *RateLimit) SetBucket(routeid string, bucket *Bucket) {
	r.buckets[r.ids[routeid]] = bucket
}

func (r *RateLimit) GetBucket(routeid string, resourceid string) *Bucket {
	requestid := routeid + resourceid

	// ID 0 is used as a Global Rate Limit Bucket (or nil).
	if routeid != GlobalRateLimitRouteID {
		switch r.ids[requestid] {
		// when a non-global route is initialized and (BucketID == "NIL"), NO rate limit applies.
		case nilRouteBucket:
			return nil

		// This rate limiter implementation points the Route ID 0 (which is reserved for a
		// Global Rate Limit) to Bucket ID "".
		//
		// As a result (of the Default Bucket mechanism), non-0 Route IDs must be handled accordingly.
		case "":
			// when a non-global route is uninitialized, set it to the Default Bucket.
			//
			// While GetBucket() can be called multiple times BEFORE a request is sent,
			// this case is only true the FIRST time a GetBucket() call is made (for that request),
			// As a result, the Bucket that is allocated from this call will ALWAYS be
			// immediately used.
			//
			// The Route's Bucket ID (Hash) is set to an artificial value while
			// the Route ID is pointed to a Default Bucket. This results in
			// subsequent calls to the Route's Default Bucket to return to
			// the same initialized bucket.
			//
			// When a Default Bucket is exhausted, it will never expire.
			// As a result, the Remaining field-value will remain at 0 until
			// the pending request (and its actual Bucket) is confirmed.
			//
			// requestID = routeid + resourceid
			// temporaryBucketID = requestID
			r.SetBucketID(requestid, requestid)

			// DefaultBucket (Per-Route) = RateLimit.DefaultBucket
			if "" == resourceid {
				if r.DefaultBucket == nil {
					return nil
				}

				b := getBucket()
				b.Limit = r.DefaultBucket.Limit
				b.Remaining = r.DefaultBucket.Limit
				r.SetBucketFromID(requestid, b)

				return b
			}

			// DefaultBucket (Per-Resource) = GetBucket(routeid, "")
			defaultBucket := r.GetBucket(routeid, "")
			if defaultBucket == nil {
				return nil
			}

			b := getBucket()
			b.Limit = defaultBucket.Limit
			b.Remaining = defaultBucket.Limit
			r.SetBucketFromID(requestid, b)

			return b
		}
	}

	return r.buckets[r.ids[requestid]]
}

func (r *RateLimit) SetDefaultBucket(bucket *Bucket) {
	r.DefaultBucket = bucket
}

func (r *RateLimit) Lock() {
	r.muQueue.Lock()
}

func (r *RateLimit) Unlock() {
	r.muQueue.Unlock()
}

func (r *RateLimit) StartTx() {
	r.muTx.Lock()
}

func (r *RateLimit) EndTx() {
	r.muTx.Unlock()
}

type hash func(routeid string, parameters ...string) (string, string)

// Hash returns a hashing function which hashes a request using its routeID.
//
// n represents the degree of resources used to hash the request.
// 	Per-Route: n = 0
//	Per-Resource: n = 1
//  ...
func Hash(n int) hash {
	return func(r string, p ...string) (string, string) {
		return r, strings.Join(p[:n], "")
	}
}

var (
	// HashPerRoute hashes a request that uses a per-route rate limit algorithm.
	HashPerRoute = Hash(0)
)

var (
	// RateLimitHashFuncs represents a map of routes to respective rate limit algorithms.
	//
	// used to determine the hashing function for routes during runtime (map[routeID]algorithm).
	RateLimitHashFuncs = map[uint8]hash{
		RouteIDs["OAuth"]:                                  HashPerRoute,
		RouteIDs["GetGlobalApplicationCommands"]:           HashPerRoute,
		RouteIDs["CreateGlobalApplicationCommand"]:         HashPerRoute,
		RouteIDs["GetGlobalApplicationCommand"]:            HashPerRoute,
		RouteIDs["EditGlobalApplicationCommand"]:           HashPerRoute,
		RouteIDs["DeleteGlobalApplicationCommand"]:         HashPerRoute,
		RouteIDs["BulkOverwriteGlobalApplicationCommands"]: HashPerRoute,
		RouteIDs["GetGuildApplicationCommands"]:            HashPerRoute,
		RouteIDs["CreateGuildApplicationCommand"]:          HashPerRoute,
		RouteIDs["GetGuildApplicationCommand"]:             HashPerRoute,
		RouteIDs["EditGuildApplicationCommand"]:            HashPerRoute,
		RouteIDs["DeleteGuildApplicationCommand"]:          HashPerRoute,
		RouteIDs["BulkOverwriteGuildApplicationCommands"]:  HashPerRoute,
		RouteIDs["GetGuildApplicationCommandPermissions"]:  HashPerRoute,
		RouteIDs["GetApplicationCommandPermissions"]:       HashPerRoute,
		RouteIDs["EditApplicationCommandPermissions"]:      HashPerRoute,
		RouteIDs["BatchEditApplicationCommandPermissions"]: HashPerRoute,
		RouteIDs["CreateInteractionResponse"]:              HashPerRoute,
		RouteIDs["GetOriginalInteractionResponse"]:         HashPerRoute,
		RouteIDs["EditOriginalInteractionResponse"]:        HashPerRoute,
		RouteIDs["DeleteOriginalInteractionResponse"]:      HashPerRoute,
		RouteIDs["CreateFollowupMessage"]:                  HashPerRoute,
		RouteIDs["GetFollowupMessage"]:                     HashPerRoute,
		RouteIDs["EditFollowupMessage"]:                    HashPerRoute,
		RouteIDs["DeleteFollowupMessage"]:                  HashPerRoute,
		RouteIDs["GetGuildAuditLog"]:                       HashPerRoute,
		RouteIDs["ListAutoModerationRulesForGuild"]:        HashPerRoute,
		RouteIDs["GetAutoModerationRule"]:                  HashPerRoute,
		RouteIDs["CreateAutoModerationRule"]:               HashPerRoute,
		RouteIDs["ModifyAutoModerationRule"]:               HashPerRoute,
		RouteIDs["DeleteAutoModerationRule"]:               HashPerRoute,
		RouteIDs["GetChannel"]:                             HashPerRoute,
		RouteIDs["ModifyChannel"]:                          HashPerRoute,
		RouteIDs["ModifyChannelGroupDM"]:                   HashPerRoute,
		RouteIDs["ModifyChannelGuild"]:                     HashPerRoute,
		RouteIDs["ModifyChannelThread"]:                    HashPerRoute,
		RouteIDs["DeleteCloseChannel"]:                     HashPerRoute,
		RouteIDs["GetChannelMessages"]:                     HashPerRoute,
		RouteIDs["GetChannelMessage"]:                      HashPerRoute,
		RouteIDs["CreateMessage"]:                          HashPerRoute,
		RouteIDs["CrosspostMessage"]:                       HashPerRoute,
		RouteIDs["CreateReaction"]:                         HashPerRoute,
		RouteIDs["DeleteOwnReaction"]:                      HashPerRoute,
		RouteIDs["DeleteUserReaction"]:                     HashPerRoute,
		RouteIDs["GetReactions"]:                           HashPerRoute,
		RouteIDs["DeleteAllReactions"]:                     HashPerRoute,
		RouteIDs["DeleteAllReactionsforEmoji"]:             HashPerRoute,
		RouteIDs["EditMessage"]:                            HashPerRoute,
		RouteIDs["DeleteMessage"]:                          HashPerRoute,
		RouteIDs["BulkDeleteMessages"]:                     HashPerRoute,
		RouteIDs["EditChannelPermissions"]:                 HashPerRoute,
		RouteIDs["GetChannelInvites"]:                      HashPerRoute,
		RouteIDs["CreateChannelInvite"]:                    HashPerRoute,
		RouteIDs["DeleteChannelPermission"]:                HashPerRoute,
		RouteIDs["FollowAnnouncementChannel"]:              HashPerRoute,
		RouteIDs["TriggerTypingIndicator"]:                 HashPerRoute,
		RouteIDs["GetPinnedMessages"]:                      HashPerRoute,
		RouteIDs["PinMessage"]:                             HashPerRoute,
		RouteIDs["UnpinMessage"]:                           HashPerRoute,
		RouteIDs["GroupDMAddRecipient"]:                    HashPerRoute,
		RouteIDs["GroupDMRemoveRecipient"]:                 HashPerRoute,
		RouteIDs["StartThreadfromMessage"]:                 HashPerRoute,
		RouteIDs["StartThreadwithoutMessage"]:              HashPerRoute,
		RouteIDs["StartThreadinForumChannel"]:              HashPerRoute,
		RouteIDs["JoinThread"]:                             HashPerRoute,
		RouteIDs["AddThreadMember"]:                        HashPerRoute,
		RouteIDs["LeaveThread"]:                            HashPerRoute,
		RouteIDs["RemoveThreadMember"]:                     HashPerRoute,
		RouteIDs["GetThreadMember"]:                        HashPerRoute,
		RouteIDs["ListThreadMembers"]:                      HashPerRoute,
		RouteIDs["ListPublicArchivedThreads"]:              HashPerRoute,
		RouteIDs["ListPrivateArchivedThreads"]:             HashPerRoute,
		RouteIDs["ListJoinedPrivateArchivedThreads"]:       HashPerRoute,
		RouteIDs["ListGuildEmojis"]:                        HashPerRoute,
		RouteIDs["GetGuildEmoji"]:                          HashPerRoute,
		RouteIDs["CreateGuildEmoji"]:                       HashPerRoute,
		RouteIDs["ModifyGuildEmoji"]:                       HashPerRoute,
		RouteIDs["DeleteGuildEmoji"]:                       HashPerRoute,
		RouteIDs["CreateGuild"]:                            HashPerRoute,
		RouteIDs["GetGuild"]:                               HashPerRoute,
		RouteIDs["GetGuildPreview"]:                        HashPerRoute,
		RouteIDs["ModifyGuild"]:                            HashPerRoute,
		RouteIDs["DeleteGuild"]:                            HashPerRoute,
		RouteIDs["CreateDM"]:                               HashPerRoute,
		RouteIDs["GetGuildChannels"]:                       HashPerRoute,
		RouteIDs["CreateGuildChannel"]:                     HashPerRoute,
		RouteIDs["ModifyGuildChannelPositions"]:            HashPerRoute,
		RouteIDs["ListActiveGuildThreads"]:                 HashPerRoute,
		RouteIDs["GetGuildMember"]:                         HashPerRoute,
		RouteIDs["ListGuildMembers"]:                       HashPerRoute,
		RouteIDs["SearchGuildMembers"]:                     HashPerRoute,
		RouteIDs["AddGuildMember"]:                         HashPerRoute,
		RouteIDs["ModifyGuildMember"]:                      HashPerRoute,
		RouteIDs["ModifyCurrentMember"]:                    HashPerRoute,
		RouteIDs["AddGuildMemberRole"]:                     HashPerRoute,
		RouteIDs["RemoveGuildMemberRole"]:                  HashPerRoute,
		RouteIDs["RemoveGuildMember"]:                      HashPerRoute,
		RouteIDs["GetGuildBans"]:                           HashPerRoute,
		RouteIDs["GetGuildBan"]:                            HashPerRoute,
		RouteIDs["CreateGuildBan"]:                         HashPerRoute,
		RouteIDs["RemoveGuildBan"]:                         HashPerRoute,
		RouteIDs["GetGuildRoles"]:                          HashPerRoute,
		RouteIDs["CreateGuildRole"]:                        HashPerRoute,
		RouteIDs["ModifyGuildRolePositions"]:               HashPerRoute,
		RouteIDs["ModifyGuildRole"]:                        HashPerRoute,
		RouteIDs["DeleteGuildRole"]:                        HashPerRoute,
		RouteIDs["ModifyGuildMFALevel"]:                    HashPerRoute,
		RouteIDs["GetGuildPruneCount"]:                     HashPerRoute,
		RouteIDs["BeginGuildPrune"]:                        HashPerRoute,
		RouteIDs["GetGuildVoiceRegions"]:                   HashPerRoute,
		RouteIDs["GetGuildInvites"]:                        HashPerRoute,
		RouteIDs["GetGuildIntegrations"]:                   HashPerRoute,
		RouteIDs["DeleteGuildIntegration"]:                 HashPerRoute,
		RouteIDs["GetGuildWidgetSettings"]:                 HashPerRoute,
		RouteIDs["ModifyGuildWidget"]:                      HashPerRoute,
		RouteIDs["GetGuildWidget"]:                         HashPerRoute,
		RouteIDs["GetGuildVanityURL"]:                      HashPerRoute,
		RouteIDs["GetGuildWidgetImage"]:                    HashPerRoute,
		RouteIDs["GetGuildWelcomeScreen"]:                  HashPerRoute,
		RouteIDs["ModifyGuildWelcomeScreen"]:               HashPerRoute,
		RouteIDs["ModifyCurrentUserVoiceState"]:            HashPerRoute,
		RouteIDs["ModifyUserVoiceState"]:                   HashPerRoute,
		RouteIDs["ListScheduledEventsforGuild"]:            HashPerRoute,
		RouteIDs["CreateGuildScheduledEvent"]:              HashPerRoute,
		RouteIDs["GetGuildScheduledEvent"]:                 HashPerRoute,
		RouteIDs["ModifyGuildScheduledEvent"]:              HashPerRoute,
		RouteIDs["DeleteGuildScheduledEvent"]:              HashPerRoute,
		RouteIDs["GetGuildScheduledEventUsers"]:            HashPerRoute,
		RouteIDs["GetGuildTemplate"]:                       HashPerRoute,
		RouteIDs["CreateGuildfromGuildTemplate"]:           HashPerRoute,
		RouteIDs["GetGuildTemplates"]:                      HashPerRoute,
		RouteIDs["CreateGuildTemplate"]:                    HashPerRoute,
		RouteIDs["SyncGuildTemplate"]:                      HashPerRoute,
		RouteIDs["ModifyGuildTemplate"]:                    HashPerRoute,
		RouteIDs["DeleteGuildTemplate"]:                    HashPerRoute,
		RouteIDs["GetInvite"]:                              HashPerRoute,
		RouteIDs["DeleteInvite"]:                           HashPerRoute,
		RouteIDs["CreateStageInstance"]:                    HashPerRoute,
		RouteIDs["GetStageInstance"]:                       HashPerRoute,
		RouteIDs["ModifyStageInstance"]:                    HashPerRoute,
		RouteIDs["DeleteStageInstance"]:                    HashPerRoute,
		RouteIDs["GetSticker"]:                             HashPerRoute,
		RouteIDs["ListNitroStickerPacks"]:                  HashPerRoute,
		RouteIDs["ListGuildStickers"]:                      HashPerRoute,
		RouteIDs["GetGuildSticker"]:                        HashPerRoute,
		RouteIDs["CreateGuildSticker"]:                     HashPerRoute,
		RouteIDs["ModifyGuildSticker"]:                     HashPerRoute,
		RouteIDs["DeleteGuildSticker"]:                     HashPerRoute,
		RouteIDs["GetCurrentUser"]:                         HashPerRoute,
		RouteIDs["GetUser"]:                                HashPerRoute,
		RouteIDs["ModifyCurrentUser"]:                      HashPerRoute,
		RouteIDs["GetCurrentUserGuilds"]:                   HashPerRoute,
		RouteIDs["GetCurrentUserGuildMember"]:              HashPerRoute,
		RouteIDs["LeaveGuild"]:                             HashPerRoute,
		RouteIDs["CreateGroupDM"]:                          HashPerRoute,
		RouteIDs["GetUserConnections"]:                     HashPerRoute,
		RouteIDs["ListVoiceRegions"]:                       HashPerRoute,
		RouteIDs["CreateWebhook"]:                          HashPerRoute,
		RouteIDs["GetChannelWebhooks"]:                     HashPerRoute,
		RouteIDs["GetGuildWebhooks"]:                       HashPerRoute,
		RouteIDs["GetWebhook"]:                             HashPerRoute,
		RouteIDs["GetWebhookwithToken"]:                    HashPerRoute,
		RouteIDs["ModifyWebhook"]:                          HashPerRoute,
		RouteIDs["ModifyWebhookwithToken"]:                 HashPerRoute,
		RouteIDs["DeleteWebhook"]:                          HashPerRoute,
		RouteIDs["DeleteWebhookwithToken"]:                 HashPerRoute,
		RouteIDs["ExecuteWebhook"]:                         HashPerRoute,
		RouteIDs["ExecuteSlackCompatibleWebhook"]:          HashPerRoute,
		RouteIDs["ExecuteGitHubCompatibleWebhook"]:         HashPerRoute,
		RouteIDs["GetWebhookMessage"]:                      HashPerRoute,
		RouteIDs["EditWebhookMessage"]:                     HashPerRoute,
		RouteIDs["DeleteWebhookMessage"]:                   HashPerRoute,
		RouteIDs["GetGateway"]:                             HashPerRoute,
		RouteIDs["GetGatewayBot"]:                          HashPerRoute,
		RouteIDs["GetCurrentBotApplicationInformation"]:    HashPerRoute,
		RouteIDs["GetCurrentAuthorizationInformation"]:     HashPerRoute,
	}
)

const (
	GlobalRateLimitRouteID = "0"
)

// RateLimiter represents an interface for rate limits.
//
// RateLimiter is an interface which allows developers to use multi-application architectures,
// which run multiple applications on separate processes or servers.
type RateLimiter interface {
	// SetBucketID maps a Route ID to a Rate Limit Bucket ID (Discord Hash).
	//
	// ID 0 is reserved for a Global Rate Limit Bucket or nil.
	SetBucketID(routeid string, bucketid string)

	// GetBucketID gets a Rate Limit Bucket ID (Discord Hash) using a Route ID.
	GetBucketID(routeid string) string

	// SetBucketFromID maps a Bucket ID to a Rate Limit Bucket.
	SetBucketFromID(bucketid string, bucket *Bucket)

	// GetBucketFromID gets a Rate Limit Bucket using the given Bucket ID.
	GetBucketFromID(bucketid string) *Bucket

	// SetBucket maps a Route ID to a Rate Limit Bucket.
	//
	// ID 0 is reserved for a Global Rate Limit Bucket or nil.
	SetBucket(routeid string, bucket *Bucket)

	// GetBucket gets a Rate Limit Bucket using the given Route ID + Resource ID.
	//
	// Implements the Default Bucket mechanism by assigning the GetBucketID(routeid) when applicable.
	GetBucket(routeid string, resourceid string) *Bucket

	// SetDefaultBucket sets the Default Bucket for per-route rate limits.
	SetDefaultBucket(bucket *Bucket)

	// Lock locks the rate limiter.
	//
	// If the lock is already in use, the calling goroutine blocks until the rate limiter is available.
	//
	// This prevents multiple requests from being PROCESSED at once, which prevents race conditions.
	// In other words, a single request is PROCESSED from a rate limiter when Lock is implemented and called.
	//
	// This does NOT prevent multiple requests from being SENT at a time.
	Lock()

	// Unlock unlocks the rate limiter.
	//
	// If the rate limiter holds multiple locks, unlocking will unblock another goroutine,
	// which allows another request to be processed.
	Unlock()

	// StartTx starts a transaction with the rate limiter.
	//
	// If a transaction is already started, the calling goroutine blocks until the rate limiter is available.
	//
	// This prevents the transaction (of Rate Limit Bucket reads and writes) from concurrent manipulation.
	StartTx()

	// EndTx ends a transaction with the rate limiter.
	//
	// If the rate limiter holds multiple transactions, ending one will unblock another goroutine,
	// which allows another transaction to start.
	EndTx()
}

// Bucket represents a Discord API Rate Limit Bucket.
type Bucket struct {
	Date      time.Time
	Expiry    time.Time
	ID        string
	Limit     int16
	Remaining int16
	Pending   int16
}

// Reset resets a Discord API Rate Limit Bucket and sets its expiry.
func (b *Bucket) Reset(expiry time.Time) {
	b.Expiry = expiry

	// Remaining = Limit - Pending
	b.Remaining = b.Limit - b.Pending
}

// Use uses the given amount of tokens for a Discord API Rate Limit Bucket.
func (b *Bucket) Use(amount int16) {
	b.Remaining -= amount
	b.Pending += amount
}

// ConfirmDate confirms the usage of a given amount of tokens for a Discord API Rate Limit Bucket,
// using the bucket's current expiry and given (Discord Header) Date time.
//
// Used for the Global Rate Limit Bucket.
func (b *Bucket) ConfirmDate(amount int16, date time.Time) {
	b.Pending -= amount

	switch {
	// Date is zero when a request has never been sent to Discord.
	//
	// set the Date of the current Bucket to the date of the current Discord Bucket.
	case b.Date.IsZero():
		b.Date = date

		// The EXACT reset period of Discord's Global Rate Limit Bucket will always occur
		// BEFORE the current Bucket resets (due to this implementation).
		//
		// reset the current Bucket with an expiry that occurs [0, 1) seconds
		// AFTER the Discord Global Rate Limit Bucket will be reset.
		//
		// This results in a Bucket's expiry that is eventually consistent with
		// Discord's Bucket expiry over time (once determined between requests).
		b.Expiry = time.Now().Add(time.Second)

	// Date is EQUAL to the Discord Bucket's Date when the request applies to the current Bucket.
	case b.Date.Equal(date):

	// Date occurs BEFORE a Discord Bucket's Date when the request applies to the next Bucket.
	//
	// update the current Bucket to the next Bucket.
	case b.Date.Before(date):
		b.Date = date

		// align the current Bucket's expiry to Discord's Bucket expiry.
		b.Expiry = time.Now().Add(time.Second)

	// Date occurs AFTER a Discord Bucket's Date when the request applied to a previous Bucket.
	case b.Date.After(date):
		b.Remaining += amount
	}
}

// ConfirmHeader confirms the usage of a given amount of tokens for a Discord API Rate Limit Bucket,
// using a given Route ID and respective Discord Rate Limit Header.
//
// Used for Route Rate Limits.
func (b *Bucket) ConfirmHeader(amount int16, header RateLimitHeader) {
	b.Pending -= amount

	// determine the reset time.
	//
	// Discord recommends to rely on the `Retry-After` header.
	// https://discord.com/developers/docs/topics/rate-limits#exceeding-a-rate-limit
	reset := time.Now().Add(time.Millisecond*time.Duration(header.ResetAfter*msPerSecond) + time.Millisecond)

	// Expiry is zero when a request from the Route ID has never been sent to Discord.
	//
	// set the current Bucket to the current Discord Bucket.
	if b.Expiry.IsZero() {
		b.Limit = int16(header.Limit)
		b.Remaining = int16(header.Remaining) - b.Pending
		b.Expiry = reset

		return
	}

	switch {
	// Expiry is EQUAL to the Discord Bucket's Reset when the request applies to the current Bucket.
	case b.Expiry == reset:

	// Expiry occurs BEFORE a Discord Bucket's Reset when the request applies to the next Bucket.
	//
	// update the current Bucket to the next Bucket.
	case b.Expiry.Before(reset):
		b.Limit = int16(header.Limit)
		b.Expiry = reset

	// Expiry occurs AFTER a Discord Bucket's Reset when the request applied to a previous Bucket.
	case b.Expiry.After(reset):
		b.Remaining += amount
	}
}

// Conversion Constants.
const (
	base10              = 10
	bit64               = 64
	msPerSecond float64 = 1000
)

// Content Types
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Type
var (
	// ContentTypeURLQueryString is an HTTP Header Content Type that indicates
	// a payload with an encoded URL Query String.
	ContentTypeURLQueryString = []byte("application/x-www-form-urlencoded")

	// ContentTypeJSON is an HTTP Header Content Type that indicates a payload with a JSON body.
	ContentTypeJSON = []byte("application/json")

	// ContentTypeMultipartForm is an HTTP Header Content Type that indicates
	// a payload with multiple content types.
	ContentTypeMultipartForm = []byte("multipart/form-data")

	// ContentTypeJPEG is an HTTP Header Content Type that indicates a payload with a JPEG image.
	ContentTypeJPEG = []byte("image/jpeg")

	// ContentTypePNG is an HTTP Header Content Type that indicates a payload with a PNG image.
	ContentTypePNG = []byte("image/png")

	// ContentTypeWebP is an HTTP Header Content Type that indicates a payload with a WebP image.
	ContentTypeWebP = []byte("image/webp")

	// ContentTypeGIF is an HTTP Header Content Type that indicates a payload with a GIF animated image.
	ContentTypeGIF = []byte("image/gif")
)

// HTTP Header Variables.
const (
	// headerAuthorizationKey represents the key for an "Authorization" HTTP Header.
	headerAuthorizationKey = "Authorization"
)

// HTTP Header Rate Limit Variables.
var (
	// headerDate represents a byte representation of "Date" for HTTP Header functionality.
	headerDate = []byte(FlagRateLimitHeaderDate)

	// headerRetryAfter represents a byte representation of "Retry-After" for HTTP Header functionality.
	headerRateLimitRetryAfter = []byte(FlagRateLimitHeaderRetryAfter)

	// headerRateLimit represents a byte representation of "X-RateLimit-Limit" for HTTP Header functionality.
	headerRateLimit = []byte(FlagRateLimitHeaderLimit)

	// headerRateLimitRemaining represents a byte representation of "X-RateLimit-Remaining" for HTTP Header functionality.
	headerRateLimitRemaining = []byte(FlagRateLimitHeaderRemaining)

	// headerRateLimitReset represents a byte representation of "X-RateLimit-Reset" for HTTP Header functionality.
	headerRateLimitReset = []byte(FlagRateLimitHeaderReset)

	// headerRateLimitResetAfter represents a byte representation of "X-RateLimit-Reset-After" for HTTP Header functionality.
	headerRateLimitResetAfter = []byte(FlagRateLimitHeaderResetAfter)

	// headerRateLimitBucket represents a byte representation of "X-RateLimit-Bucket" for HTTP Header functionality.
	headerRateLimitBucket = []byte(FlagRateLimitHeaderBucket)

	// headerRateLimitGlobal represents a byte representation of "X-RateLimit-Global" for HTTP Header functionality.
	headerRateLimitGlobal = []byte(FlagRateLimitHeaderGlobal)

	// headerRateLimitScope represents a byte representation of "X-RateLimit-Scope" for HTTP Header functionality.
	headerRateLimitScope = []byte(FlagRateLimitHeaderScope)
)

// peekDate peeks an HTTP Header for the Date.
func peekDate(r *fasthttp.Response) (time.Time, error) {
	date, err := time.Parse(time.RFC1123, string(r.Header.PeekBytes(headerDate)))
	if err != nil {
		return time.Time{}, fmt.Errorf("error occurred parsing the \"Date\" HTTP Header: %w", err)
	}

	return date, nil
}

// peekHeaderRetryAfter peeks an HTTP Header for the Rate Limit Header "Retry-After".
func peekHeaderRetryAfter(r *fasthttp.Response) (float64, error) {
	retryafter, err := strconv.ParseFloat(string(r.Header.PeekBytes(headerRateLimitRetryAfter)), bit64)
	if err != nil {
		return 0, fmt.Errorf(errRateLimit, string(headerRateLimitRetryAfter), err)
	}

	return retryafter, nil
}

// peekHeaderRateLimit peeks an HTTP Header for Discord Rate Limit Header values.
func peekHeaderRateLimit(r *fasthttp.Response) RateLimitHeader {
	limit, _ := strconv.Atoi(string(r.Header.PeekBytes(headerRateLimit)))
	remaining, _ := strconv.Atoi(string(r.Header.PeekBytes(headerRateLimitRemaining)))
	reset, _ := strconv.ParseFloat(string(r.Header.PeekBytes(headerRateLimitReset)), bit64)
	resetafter, _ := strconv.ParseFloat(string(r.Header.PeekBytes(headerRateLimitResetAfter)), bit64)
	global, _ := strconv.ParseBool(string(r.Header.PeekBytes(headerRateLimitGlobal)))

	return RateLimitHeader{
		Limit:      limit,
		Remaining:  remaining,
		Reset:      reset,
		ResetAfter: resetafter,
		Bucket:     string(r.Header.PeekBytes(headerRateLimitBucket)),
		Global:     global,
		Scope:      string(r.Header.PeekBytes(headerRateLimitScope)),
	}
}

// SendRequest sends a fasthttp.Request using the given route ID, HTTP method, URI, content type and body,
// then parses the response into dst.
func SendRequest(bot *Client, xid, routeid, resourceid, method, uri string, content, body []byte, dst any) error { //nolint:gocyclo,maintidx
	retries := 0
	requestid := routeid + resourceid
	request := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(request)
	request.Header.SetMethod(method)
	request.Header.SetContentTypeBytes(content)
	request.Header.Set(headerAuthorizationKey, bot.Authentication.Header)
	request.SetRequestURI(uri)
	request.SetBodyRaw(body)
	response := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(response)

	// Certain endpoints are not bound to the bot's Global Rate Limit.
	if IgnoreGlobalRateLimitRouteIDs[requestid] {
		goto SEND
	}

RATELIMIT:
	// a single request or response is PROCESSED at any point in time.
	bot.Config.Request.RateLimiter.Lock()

	LogRequest(Logger.Trace(), bot.ApplicationID, xid, routeid, resourceid, uri).Msg("processing request")
	if Logger.GetLevel() == zerolog.TraceLevel {
		LogRequestBody(Logger.Trace(), bot.ApplicationID, xid, routeid, resourceid, uri, string(body)).Msg("")
	}

	// check Global and Route Rate Limit Buckets prior to sending the current request.
	for {
		bot.Config.Request.RateLimiter.StartTx()

		globalBucket := bot.Config.Request.RateLimiter.GetBucket(GlobalRateLimitRouteID, "")

		// stop waiting when the Global Rate Limit Bucket is NOT empty.
		if isNotEmpty(globalBucket) {
			routeBucket := bot.Config.Request.RateLimiter.GetBucket(routeid, resourceid)

			if isNotEmpty(routeBucket) {
				break
			}

			if isExpired(routeBucket) {
				// When a Route Bucket expires, its new expiry becomes unknown.
				// As a result, it will never reset (again) until a pending request's
				// response sets a new expiry.
				routeBucket.Reset(time.Time{})
			}

			var wait time.Time
			if routeBucket != nil {
				wait = routeBucket.Expiry
			}

			// do NOT block other requests due to a Route Rate Limit.
			bot.Config.Request.RateLimiter.EndTx()
			bot.Config.Request.RateLimiter.Unlock()

			// reduce CPU usage by blocking the current goroutine
			// until it's eligible for action.
			if routeBucket != nil {
				<-time.After(time.Until(wait))
			}

			goto RATELIMIT
		}

		// reset the Global Rate Limit Bucket when the current Bucket has passed its expiry.
		if isExpired(globalBucket) {
			globalBucket.Reset(time.Now().Add(time.Second))
		}

		bot.Config.Request.RateLimiter.EndTx()
	}

	if globalBucket := bot.Config.Request.RateLimiter.GetBucket(GlobalRateLimitRouteID, ""); globalBucket != nil {
		globalBucket.Use(1)
	}

	if routeBucket := bot.Config.Request.RateLimiter.GetBucket(routeid, resourceid); routeBucket != nil {
		routeBucket.Use(1)
	}

	bot.Config.Request.RateLimiter.EndTx()
	bot.Config.Request.RateLimiter.Unlock()

SEND:
	LogRequest(Logger.Trace(), bot.ApplicationID, xid, routeid, resourceid, uri).Msg("sending request")

	// send the request.
	if err := bot.Config.Request.Client.DoTimeout(request, response, bot.Config.Request.Timeout); err != nil {
		return fmt.Errorf("%w", err)
	}

	LogResponse(LogRequest(Logger.Info(), bot.ApplicationID, xid, routeid, resourceid, uri),
		response.Header.String(), string(response.Body()),
	).Msg("")

	var header RateLimitHeader

	// confirm the response with the rate limiter.
	//
	// Certain endpoints are not bound to the bot's Global Rate Limit.
	if !IgnoreGlobalRateLimitRouteIDs[requestid] {
		// parse the Rate Limit Header for per-route rate limit functionality.
		header = peekHeaderRateLimit(response)

		// parse the Date header for Global Rate Limit functionality.
		date, err := peekDate(response)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		bot.Config.Request.RateLimiter.StartTx()

		// confirm the Global Rate Limit Bucket.
		globalBucket := bot.Config.Request.RateLimiter.GetBucket(GlobalRateLimitRouteID, "")
		if globalBucket != nil {
			globalBucket.ConfirmDate(1, date)
		}

		// confirm the Route Rate Limit Bucket (if applicable).
		routeBucket := bot.Config.Request.RateLimiter.GetBucket(routeid, resourceid)
		switch {
		// when there is no Discord Bucket, remove the route's mapping to a rate limit Bucket.
		case header.Bucket == "":
			if routeBucket != nil {
				bot.Config.Request.RateLimiter.SetBucketID(requestid, nilRouteBucket)
				routeBucket = nil
			}

		// when the route's Bucket ID does NOT match the Discord Bucket, update it.
		case routeBucket == nil && header.Bucket != "" || routeBucket.ID != header.Bucket:
			var pending int16
			if routeBucket != nil {
				pending = routeBucket.Pending
			}

			// update the route ID mapping to a rate limit Bucket ID.
			bot.Config.Request.RateLimiter.SetBucketID(requestid, header.Bucket)

			// map the Bucket ID to the updated Rate Limit Bucket.
			if bucket := bot.Config.Request.RateLimiter.GetBucketFromID(header.Bucket); bucket != nil {
				routeBucket = bucket
			} else {
				routeBucket = getBucket()
				bot.Config.Request.RateLimiter.SetBucketFromID(header.Bucket, routeBucket)
			}

			routeBucket.Pending += pending
			routeBucket.ID = header.Bucket
		}

		if routeBucket != nil {
			routeBucket.ConfirmHeader(1, header)
		}

		if response.StatusCode() != fasthttp.StatusTooManyRequests {
			bot.Config.Request.RateLimiter.EndTx()
		}
	}

	// handle the response.
	switch response.StatusCode() {
	case fasthttp.StatusOK, fasthttp.StatusCreated:
		// parse the response data.
		if err := json.Unmarshal(response.Body(), dst); err != nil {
			return fmt.Errorf(errUnmarshal, dst, err)
		}

		return nil

	case fasthttp.StatusNoContent:
		return nil

	// process the rate limit.
	case fasthttp.StatusTooManyRequests:
		retry := retries < bot.Config.Request.Retries
		retries++

		switch header.Scope { //nolint:gocritic
		// Discord per-resource (shared) rate limit headers include the per-route (user) bucket.
		//
		// when a per-resource rate limit is encountered, send another request or return.
		case RateLimitScopeValueShared:
			bot.Config.Request.RateLimiter.EndTx()

			if retry || bot.Config.Request.RetryShared {
				goto RATELIMIT
			}

			return StatusCodeError(response.StatusCode())
		}

		// parse the rate limit response data for `retry_after`.
		var data RateLimitResponse
		if err := json.Unmarshal(response.Body(), &data); err != nil {
			return fmt.Errorf("%w", err)
		}

		// determine the reset time.
		var reset time.Time
		if data.RetryAfter == 0 {
			// when the 429 is from Discord, use the `retry_after` value (s).
			reset = time.Now().Add(time.Millisecond * time.Duration(data.RetryAfter*msPerSecond))
		} else {
			// when the 429 is from a Cloudflare Ban, use the `"Retry-After"` value (s).
			retryafter, err := peekHeaderRetryAfter(response)
			if err != nil {
				return fmt.Errorf("%w", err)
			}

			reset = time.Now().Add(time.Millisecond * time.Duration(retryafter*msPerSecond))
		}

		LogRequest(Logger.Debug(), bot.ApplicationID, xid, routeid, resourceid, uri).
			Time(LogCtxReset, reset).Msg("")

		switch header.Global {
		// when the global request rate limit is encountered.
		case true:
			// when the current time is BEFORE the reset time,
			// all requests must wait until the 429 expires.
			if time.Now().Before(reset) {
				if globalBucket := bot.Config.Request.RateLimiter.GetBucket(GlobalRateLimitRouteID, ""); globalBucket != nil {
					globalBucket.Remaining = 0
					globalBucket.Expiry = reset.Add(time.Millisecond)
				}
			}

			bot.Config.Request.RateLimiter.EndTx()

		// when a per-route (user) rate limit is encountered.
		case false:
			// do NOT block other requests while waiting for a Route Rate Limit.
			bot.Config.Request.RateLimiter.EndTx()

			// when the current time is BEFORE the reset time,
			// requests with the same Rate Limit Bucket must wait until the 429 expires.
			if time.Now().Before(reset) {
				if routeBucket := bot.Config.Request.RateLimiter.GetBucket(routeid, resourceid); routeBucket != nil {
					routeBucket.Remaining = 0
					routeBucket.Expiry = reset.Add(time.Millisecond)
				}
			}
		}

		if retry {
			goto RATELIMIT
		}

		return StatusCodeError(fasthttp.StatusTooManyRequests)

	// retry the request on a bad gateway server error.
	case fasthttp.StatusBadGateway:
		if retries < bot.Config.Request.Retries {
			retries++

			goto RATELIMIT
		}

		return StatusCodeError(fasthttp.StatusBadGateway)

	default:
		return StatusCodeError(response.StatusCode())
	}
}

// isExpired determines whether a rate limit Bucket is expired.
func isExpired(b *Bucket) bool {
	// a rate limit bucket is expired when
	// 1. the bucket exists AND
	// 2. the current time occurs after the non-zero expiry time.
	return b != nil && !b.Expiry.IsZero() && time.Now().After(b.Expiry)
}

// isNotEmpty determines whether a rate limit Bucket is NOT empty.
func isNotEmpty(b *Bucket) bool {
	// a rate limit bucket is NOT empty when
	// 1. the bucket does not exist OR
	// 2. there is one or more remaining request token(s).
	return b == nil || b.Remaining > 0
}

// boundary represents the boundary that is used in every multipart form.
var boundary = randomBoundary()

// randomBoundary generates a random value to be used as a boundary in multipart forms.
//
// Implemented from the mime/multipart package.
func randomBoundary() string {
	var buf [30]byte
	_, err := io.ReadFull(rand.Reader, buf[:])
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%x", buf[:])
}

// quoteEscaper escapes quotes and backslashes in a multipart form.
var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

// createMultipartForm creates a multipart/form for Discord using a given JSON body and files.
func createMultipartForm(json []byte, files ...*File) ([]byte, []byte, error) {
	form := bytes.NewBuffer(nil)

	// set the boundary.
	multipartWriter := multipart.NewWriter(form)
	err := multipartWriter.SetBoundary(boundary)
	if err != nil {
		return nil, nil, fmt.Errorf("error setting multipart form boundary: %w", err)
	}

	// add the `payload_json` JSON to the form.
	multipartPayloadJSONPart, err := createPayloadJSONForm(multipartWriter)
	if err != nil {
		return nil, nil, fmt.Errorf("error adding JSON payload header to multipart form: %w", err)
	}

	if _, err := multipartPayloadJSONPart.Write(json); err != nil {
		return nil, nil, fmt.Errorf("error writing JSON payload data to multipart form: %w", err)
	}

	// add the remaining files to the form.
	for i, file := range files {
		name := strings.Join([]string{"files[", strconv.Itoa(i), "]"}, "")
		multipartFilePart, err := createFormFile(multipartWriter, name, file.Name, file.ContentType)
		if err != nil {
			return nil, nil, fmt.Errorf("error adding a file %q to a multipart form: %w", file.Name, err)
		}

		if _, err := multipartFilePart.Write(file.Data); err != nil {
			return nil, nil, fmt.Errorf("error writing file %q data to multipart form: %w", file.Name, err)
		}
	}

	// write the trailing boundary.
	if err := multipartWriter.Close(); err != nil {
		return nil, nil, fmt.Errorf("error closing the multipart form: %w", err)
	}

	return []byte(multipartWriter.FormDataContentType()), form.Bytes(), nil
}

var (
	// contentTypeJSONString is an HTTP Header Content Type that indicates a payload with a JSON body.
	contentTypeJSONString = "application/json"

	// contentTypeOctetStreamString is an HTTP Header Content Type that indicates a payload with binary data.
	contentTypeOctetStreamString = "application/octet-stream"
)

// createPayloadJSONForm creates a form-data header for the `payload_json` file in a multipart form.
func createPayloadJSONForm(m *multipart.Writer) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="payload_json"`)
	h.Set("Content-Type", contentTypeJSONString)

	return m.CreatePart(h) //nolint:wrapcheck
}

// createFormFile creates a form-data header for file attachments in a multipart form.
func createFormFile(m *multipart.Writer, name, filename, contentType string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`, name, quoteEscaper.Replace(filename)))

	if "" == contentType {
		contentType = contentTypeOctetStreamString
	}

	h.Set("Content-Type", contentType)

	return m.CreatePart(h) //nolint:wrapcheck
}

var (
	// qsEncoder is used to create URL Query Strings from objects.
	qsEncoder = schema.NewEncoder()
)

// init runs at the start of the program.
func init() {
	// use `url` tags for the URL Query String encoder and decoder.
	qsEncoder.SetAliasTag("url")
}

// EndpointQueryString returns a URL Query String from a given object.
func EndpointQueryString(dst any) (string, error) {
	params := url.Values{}
	err := qsEncoder.Encode(dst, params)
	if err != nil {
		return "", err //nolint:wrapcheck
	}

	return params.Encode(), nil
}

var (
	// RouteIDs represents a map of Routes to Route IDs (map[string]uint8).
	RouteIDs = map[string]uint8{
		"":                                       0,
		"OAuth":                                  1,
		"GetGlobalApplicationCommands":           2,
		"CreateGlobalApplicationCommand":         3,
		"GetGlobalApplicationCommand":            4,
		"EditGlobalApplicationCommand":           5,
		"DeleteGlobalApplicationCommand":         6,
		"BulkOverwriteGlobalApplicationCommands": 7,
		"GetGuildApplicationCommands":            8,
		"CreateGuildApplicationCommand":          9,
		"GetGuildApplicationCommand":             10,
		"EditGuildApplicationCommand":            11,
		"DeleteGuildApplicationCommand":          12,
		"BulkOverwriteGuildApplicationCommands":  13,
		"GetGuildApplicationCommandPermissions":  14,
		"GetApplicationCommandPermissions":       15,
		"EditApplicationCommandPermissions":      16,
		"BatchEditApplicationCommandPermissions": 17,
		"CreateInteractionResponse":              18,
		"GetOriginalInteractionResponse":         19,
		"EditOriginalInteractionResponse":        20,
		"DeleteOriginalInteractionResponse":      21,
		"CreateFollowupMessage":                  22,
		"GetFollowupMessage":                     23,
		"EditFollowupMessage":                    24,
		"DeleteFollowupMessage":                  25,
		"GetGuildAuditLog":                       26,
		"ListAutoModerationRulesForGuild":        27,
		"GetAutoModerationRule":                  28,
		"CreateAutoModerationRule":               29,
		"ModifyAutoModerationRule":               30,
		"DeleteAutoModerationRule":               31,
		"GetChannel":                             32,
		"ModifyChannel":                          33,
		"ModifyChannelGroupDM":                   34,
		"ModifyChannelGuild":                     35,
		"ModifyChannelThread":                    36,
		"DeleteCloseChannel":                     37,
		"GetChannelMessages":                     38,
		"GetChannelMessage":                      39,
		"CreateMessage":                          40,
		"CrosspostMessage":                       41,
		"CreateReaction":                         42,
		"DeleteOwnReaction":                      43,
		"DeleteUserReaction":                     44,
		"GetReactions":                           45,
		"DeleteAllReactions":                     46,
		"DeleteAllReactionsforEmoji":             47,
		"EditMessage":                            48,
		"DeleteMessage":                          49,
		"BulkDeleteMessages":                     50,
		"EditChannelPermissions":                 51,
		"GetChannelInvites":                      52,
		"CreateChannelInvite":                    53,
		"DeleteChannelPermission":                54,
		"FollowAnnouncementChannel":              55,
		"TriggerTypingIndicator":                 56,
		"GetPinnedMessages":                      57,
		"PinMessage":                             58,
		"UnpinMessage":                           59,
		"GroupDMAddRecipient":                    60,
		"GroupDMRemoveRecipient":                 61,
		"StartThreadfromMessage":                 62,
		"StartThreadwithoutMessage":              63,
		"StartThreadinForumChannel":              64,
		"JoinThread":                             65,
		"AddThreadMember":                        66,
		"LeaveThread":                            67,
		"RemoveThreadMember":                     68,
		"GetThreadMember":                        69,
		"ListThreadMembers":                      70,
		"ListPublicArchivedThreads":              71,
		"ListPrivateArchivedThreads":             72,
		"ListJoinedPrivateArchivedThreads":       73,
		"ListGuildEmojis":                        74,
		"GetGuildEmoji":                          75,
		"CreateGuildEmoji":                       76,
		"ModifyGuildEmoji":                       77,
		"DeleteGuildEmoji":                       78,
		"CreateGuild":                            79,
		"GetGuild":                               80,
		"GetGuildPreview":                        81,
		"ModifyGuild":                            82,
		"DeleteGuild":                            83,
		"GetGuildChannels":                       84,
		"CreateGuildChannel":                     85,
		"ModifyGuildChannelPositions":            86,
		"ListActiveGuildThreads":                 87,
		"GetGuildMember":                         88,
		"ListGuildMembers":                       89,
		"SearchGuildMembers":                     90,
		"AddGuildMember":                         91,
		"ModifyGuildMember":                      92,
		"ModifyCurrentMember":                    93,
		"AddGuildMemberRole":                     94,
		"RemoveGuildMemberRole":                  95,
		"RemoveGuildMember":                      96,
		"GetGuildBans":                           97,
		"GetGuildBan":                            98,
		"CreateGuildBan":                         99,
		"RemoveGuildBan":                         100,
		"GetGuildRoles":                          101,
		"CreateGuildRole":                        102,
		"ModifyGuildRolePositions":               103,
		"ModifyGuildRole":                        104,
		"DeleteGuildRole":                        105,
		"ModifyGuildMFALevel":                    106,
		"GetGuildPruneCount":                     107,
		"BeginGuildPrune":                        108,
		"GetGuildVoiceRegions":                   109,
		"GetGuildInvites":                        110,
		"GetGuildIntegrations":                   111,
		"DeleteGuildIntegration":                 112,
		"GetGuildWidgetSettings":                 113,
		"ModifyGuildWidget":                      114,
		"GetGuildWidget":                         115,
		"GetGuildVanityURL":                      116,
		"GetGuildWidgetImage":                    117,
		"GetGuildWelcomeScreen":                  118,
		"ModifyGuildWelcomeScreen":               119,
		"ModifyCurrentUserVoiceState":            120,
		"ModifyUserVoiceState":                   121,
		"ListScheduledEventsforGuild":            122,
		"CreateGuildScheduledEvent":              123,
		"GetGuildScheduledEvent":                 124,
		"ModifyGuildScheduledEvent":              125,
		"DeleteGuildScheduledEvent":              126,
		"GetGuildScheduledEventUsers":            127,
		"GetGuildTemplate":                       128,
		"CreateGuildfromGuildTemplate":           129,
		"GetGuildTemplates":                      130,
		"CreateGuildTemplate":                    131,
		"SyncGuildTemplate":                      132,
		"ModifyGuildTemplate":                    133,
		"DeleteGuildTemplate":                    134,
		"GetInvite":                              135,
		"DeleteInvite":                           136,
		"CreateStageInstance":                    137,
		"GetStageInstance":                       138,
		"ModifyStageInstance":                    139,
		"DeleteStageInstance":                    140,
		"GetSticker":                             141,
		"ListNitroStickerPacks":                  142,
		"ListGuildStickers":                      143,
		"GetGuildSticker":                        144,
		"CreateGuildSticker":                     145,
		"ModifyGuildSticker":                     146,
		"DeleteGuildSticker":                     147,
		"GetCurrentUser":                         148,
		"GetUser":                                149,
		"ModifyCurrentUser":                      150,
		"GetCurrentUserGuilds":                   151,
		"GetCurrentUserGuildMember":              152,
		"LeaveGuild":                             153,
		"CreateDM":                               154,
		"CreateGroupDM":                          155,
		"GetUserConnections":                     156,
		"ListVoiceRegions":                       157,
		"CreateWebhook":                          158,
		"GetChannelWebhooks":                     159,
		"GetGuildWebhooks":                       160,
		"GetWebhook":                             161,
		"GetWebhookwithToken":                    162,
		"ModifyWebhook":                          163,
		"ModifyWebhookwithToken":                 164,
		"DeleteWebhook":                          165,
		"DeleteWebhookwithToken":                 166,
		"ExecuteWebhook":                         167,
		"ExecuteSlackCompatibleWebhook":          168,
		"ExecuteGitHubCompatibleWebhook":         169,
		"GetWebhookMessage":                      170,
		"EditWebhookMessage":                     171,
		"DeleteWebhookMessage":                   172,
		"GetGateway":                             173,
		"GetGatewayBot":                          174,
		"GetCurrentBotApplicationInformation":    175,
		"GetCurrentAuthorizationInformation":     176,
	}
)

// Send sends a GetGlobalApplicationCommands request to Discord and returns a []*ApplicationCommand.
func (r *GetGlobalApplicationCommands) Send(bot *Client) ([]*ApplicationCommand, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[2]("2")
	query, err := EndpointQueryString(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      "",
			Err:           err,
		}
	}
	endpoint := EndpointGetGlobalApplicationCommands(bot.ApplicationID) + "?" + query

	result := make([]*ApplicationCommand, 0)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, ContentTypeURLQueryString, nil, &result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a CreateGlobalApplicationCommand request to Discord and returns a ApplicationCommand.
func (r *CreateGlobalApplicationCommand) Send(bot *Client) (*ApplicationCommand, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[3]("3")
	endpoint := EndpointCreateGlobalApplicationCommand(bot.ApplicationID)

	body, err := json.Marshal(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           fmt.Errorf(errSendMarshal, err),
		}
	}

	result := new(ApplicationCommand)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPost, endpoint, ContentTypeJSON, body, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a GetGlobalApplicationCommand request to Discord and returns a ApplicationCommand.
func (r *GetGlobalApplicationCommand) Send(bot *Client) (*ApplicationCommand, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[4]("4", "297ffb1f"+r.CommandID)
	endpoint := EndpointGetGlobalApplicationCommand(bot.ApplicationID, r.CommandID)

	result := new(ApplicationCommand)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, nil, nil, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a EditGlobalApplicationCommand request to Discord and returns a ApplicationCommand.
func (r *EditGlobalApplicationCommand) Send(bot *Client) (*ApplicationCommand, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[5]("5", "297ffb1f"+r.CommandID)
	endpoint := EndpointEditGlobalApplicationCommand(bot.ApplicationID, r.CommandID)

	body, err := json.Marshal(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           fmt.Errorf(errSendMarshal, err),
		}
	}

	result := new(ApplicationCommand)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPatch, endpoint, ContentTypeJSON, body, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a DeleteGlobalApplicationCommand request to Discord and returns a error.
func (r *DeleteGlobalApplicationCommand) Send(bot *Client) error {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[6]("6", "297ffb1f"+r.CommandID)
	endpoint := EndpointDeleteGlobalApplicationCommand(bot.ApplicationID, r.CommandID)

	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodDelete, endpoint, nil, nil, nil)
	if err != nil {
		return ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return nil
}

// Send sends a BulkOverwriteGlobalApplicationCommands request to Discord and returns a []*ApplicationCommand.
func (r *BulkOverwriteGlobalApplicationCommands) Send(bot *Client) ([]*ApplicationCommand, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[7]("7")
	endpoint := EndpointBulkOverwriteGlobalApplicationCommands(bot.ApplicationID)

	body, err := json.Marshal(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           fmt.Errorf(errSendMarshal, err),
		}
	}

	result := make([]*ApplicationCommand, 0)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPut, endpoint, ContentTypeJSON, body, &result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a GetGuildApplicationCommands request to Discord and returns a []*ApplicationCommand.
func (r *GetGuildApplicationCommands) Send(bot *Client) ([]*ApplicationCommand, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[8]("8", "45892a5d"+r.GuildID)
	query, err := EndpointQueryString(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      "",
			Err:           err,
		}
	}
	endpoint := EndpointGetGuildApplicationCommands(bot.ApplicationID, r.GuildID) + "?" + query

	result := make([]*ApplicationCommand, 0)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, ContentTypeURLQueryString, nil, &result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a CreateGuildApplicationCommand request to Discord and returns a ApplicationCommand.
func (r *CreateGuildApplicationCommand) Send(bot *Client) (*ApplicationCommand, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[9]("9", "45892a5d"+r.GuildID)
	endpoint := EndpointCreateGuildApplicationCommand(bot.ApplicationID, r.GuildID)

	body, err := json.Marshal(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           fmt.Errorf(errSendMarshal, err),
		}
	}

	result := new(ApplicationCommand)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPost, endpoint, ContentTypeJSON, body, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a GetGuildApplicationCommand request to Discord and returns a ApplicationCommand.
func (r *GetGuildApplicationCommand) Send(bot *Client) (*ApplicationCommand, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[10]("10", "45892a5d"+r.GuildID, "297ffb1f"+r.CommandID)
	endpoint := EndpointGetGuildApplicationCommand(bot.ApplicationID, r.GuildID, r.CommandID)

	result := new(ApplicationCommand)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, nil, nil, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a EditGuildApplicationCommand request to Discord and returns a ApplicationCommand.
func (r *EditGuildApplicationCommand) Send(bot *Client) (*ApplicationCommand, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[11]("11", "45892a5d"+r.GuildID, "297ffb1f"+r.CommandID)
	endpoint := EndpointEditGuildApplicationCommand(bot.ApplicationID, r.GuildID, r.CommandID)

	body, err := json.Marshal(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           fmt.Errorf(errSendMarshal, err),
		}
	}

	result := new(ApplicationCommand)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPatch, endpoint, ContentTypeJSON, body, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a DeleteGuildApplicationCommand request to Discord and returns a error.
func (r *DeleteGuildApplicationCommand) Send(bot *Client) error {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[12]("12", "45892a5d"+r.GuildID, "297ffb1f"+r.CommandID)
	endpoint := EndpointDeleteGuildApplicationCommand(bot.ApplicationID, r.GuildID, r.CommandID)

	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodDelete, endpoint, nil, nil, nil)
	if err != nil {
		return ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return nil
}

// Send sends a BulkOverwriteGuildApplicationCommands request to Discord and returns a []*ApplicationCommand.
func (r *BulkOverwriteGuildApplicationCommands) Send(bot *Client) ([]*ApplicationCommand, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[13]("13", "45892a5d"+r.GuildID)
	endpoint := EndpointBulkOverwriteGuildApplicationCommands(bot.ApplicationID, r.GuildID)

	body, err := json.Marshal(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           fmt.Errorf(errSendMarshal, err),
		}
	}

	result := make([]*ApplicationCommand, 0)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPut, endpoint, ContentTypeJSON, body, &result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a GetGuildApplicationCommandPermissions request to Discord and returns a GuildApplicationCommandPermissions.
func (r *GetGuildApplicationCommandPermissions) Send(bot *Client) (*GuildApplicationCommandPermissions, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[14]("14", "45892a5d"+r.GuildID)
	endpoint := EndpointGetGuildApplicationCommandPermissions(bot.ApplicationID, r.GuildID)

	result := new(GuildApplicationCommandPermissions)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, nil, nil, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a GetApplicationCommandPermissions request to Discord and returns a GuildApplicationCommandPermissions.
func (r *GetApplicationCommandPermissions) Send(bot *Client) (*GuildApplicationCommandPermissions, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[15]("15", "45892a5d"+r.GuildID, "297ffb1f"+r.CommandID)
	endpoint := EndpointGetApplicationCommandPermissions(bot.ApplicationID, r.GuildID, r.CommandID)

	result := new(GuildApplicationCommandPermissions)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, nil, nil, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a EditApplicationCommandPermissions request to Discord and returns a GuildApplicationCommandPermissions.
func (r *EditApplicationCommandPermissions) Send(bot *Client) (*GuildApplicationCommandPermissions, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[16]("16", "45892a5d"+r.GuildID, "297ffb1f"+r.CommandID)
	endpoint := EndpointEditApplicationCommandPermissions(bot.ApplicationID, r.GuildID, r.CommandID)

	body, err := json.Marshal(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           fmt.Errorf(errSendMarshal, err),
		}
	}

	result := new(GuildApplicationCommandPermissions)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPut, endpoint, ContentTypeJSON, body, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a BatchEditApplicationCommandPermissions request to Discord and returns a GuildApplicationCommandPermissions.
func (r *BatchEditApplicationCommandPermissions) Send(bot *Client) (*GuildApplicationCommandPermissions, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[17]("17", "45892a5d"+r.GuildID)
	endpoint := EndpointBatchEditApplicationCommandPermissions(bot.ApplicationID, r.GuildID)

	result := new(GuildApplicationCommandPermissions)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPut, endpoint, nil, nil, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a CreateInteractionResponse request to Discord and returns a error.
func (r *CreateInteractionResponse) Send(bot *Client) error {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[18]("18", "beb3d0e6"+r.InteractionID, "cb69bb28"+r.InteractionToken)
	endpoint := EndpointCreateInteractionResponse(r.InteractionID, r.InteractionToken)

	body, err := json.Marshal(r)
	if err != nil {
		return ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           fmt.Errorf(errSendMarshal, err),
		}
	}

	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPost, endpoint, ContentTypeJSON, body, nil)
	if err != nil {
		return ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return nil
}

// Send sends a GetOriginalInteractionResponse request to Discord and returns a error.
func (r *GetOriginalInteractionResponse) Send(bot *Client) error {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[19]("19", "cb69bb28"+r.InteractionToken)
	query, err := EndpointQueryString(r)
	if err != nil {
		return ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      "",
			Err:           err,
		}
	}
	endpoint := EndpointGetOriginalInteractionResponse(bot.ApplicationID, r.InteractionToken) + "?" + query

	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPatch, endpoint, ContentTypeURLQueryString, nil, nil)
	if err != nil {
		return ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return nil
}

// Send sends a EditOriginalInteractionResponse request to Discord and returns a Message.
func (r *EditOriginalInteractionResponse) Send(bot *Client) (*Message, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[20]("20", "cb69bb28"+r.InteractionToken)
	query, err := EndpointQueryString(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      "",
			Err:           err,
		}
	}
	endpoint := EndpointEditOriginalInteractionResponse(bot.ApplicationID, r.InteractionToken) + "?" + query

	body, err := json.Marshal(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           fmt.Errorf(errSendMarshal, err),
		}
	}

	contentType := ContentTypeJSON
	if len(r.Files) != 0 {
		var multipartErr error
		if contentType, body, multipartErr = createMultipartForm(body, r.Files...); multipartErr != nil {
			return nil, ErrorRequest{
				ClientID:      bot.ApplicationID,
				CorrelationID: xid,
				RouteID:       routeid,
				ResourceID:    resourceid,
				Endpoint:      "",
				Err:           err,
			}
		}
	}

	result := new(Message)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPatch, endpoint, contentType, body, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a DeleteOriginalInteractionResponse request to Discord and returns a error.
func (r *DeleteOriginalInteractionResponse) Send(bot *Client) error {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[21]("21", "cb69bb28"+r.InteractionToken)
	endpoint := EndpointDeleteOriginalInteractionResponse(bot.ApplicationID, r.InteractionToken)

	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodDelete, endpoint, nil, nil, nil)
	if err != nil {
		return ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return nil
}

// Send sends a CreateFollowupMessage request to Discord and returns a Message.
func (r *CreateFollowupMessage) Send(bot *Client) (*Message, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[22]("22", "cb69bb28"+r.InteractionToken)
	query, err := EndpointQueryString(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      "",
			Err:           err,
		}
	}
	endpoint := EndpointCreateFollowupMessage(bot.ApplicationID, r.InteractionToken) + "?" + query

	body, err := json.Marshal(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           fmt.Errorf(errSendMarshal, err),
		}
	}

	contentType := ContentTypeJSON
	if len(r.Files) != 0 {
		var multipartErr error
		if contentType, body, multipartErr = createMultipartForm(body, r.Files...); multipartErr != nil {
			return nil, ErrorRequest{
				ClientID:      bot.ApplicationID,
				CorrelationID: xid,
				RouteID:       routeid,
				ResourceID:    resourceid,
				Endpoint:      "",
				Err:           err,
			}
		}
	}

	result := new(Message)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPost, endpoint, contentType, body, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a GetFollowupMessage request to Discord and returns a Message.
func (r *GetFollowupMessage) Send(bot *Client) (*Message, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[23]("23", "cb69bb28"+r.InteractionToken, "d57d6589"+r.MessageID)
	query, err := EndpointQueryString(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      "",
			Err:           err,
		}
	}
	endpoint := EndpointGetFollowupMessage(bot.ApplicationID, r.InteractionToken, r.MessageID) + "?" + query

	result := new(Message)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, ContentTypeURLQueryString, nil, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a EditFollowupMessage request to Discord and returns a Message.
func (r *EditFollowupMessage) Send(bot *Client) (*Message, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[24]("24", "cb69bb28"+r.InteractionToken, "d57d6589"+r.MessageID)
	query, err := EndpointQueryString(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      "",
			Err:           err,
		}
	}
	endpoint := EndpointEditFollowupMessage(bot.ApplicationID, r.InteractionToken, r.MessageID) + "?" + query

	body, err := json.Marshal(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           fmt.Errorf(errSendMarshal, err),
		}
	}

	contentType := ContentTypeJSON
	if len(r.Files) != 0 {
		var multipartErr error
		if contentType, body, multipartErr = createMultipartForm(body, r.Files...); multipartErr != nil {
			return nil, ErrorRequest{
				ClientID:      bot.ApplicationID,
				CorrelationID: xid,
				RouteID:       routeid,
				ResourceID:    resourceid,
				Endpoint:      "",
				Err:           err,
			}
		}
	}

	result := new(Message)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPatch, endpoint, contentType, body, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a DeleteFollowupMessage request to Discord and returns a error.
func (r *DeleteFollowupMessage) Send(bot *Client) error {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[25]("25", "cb69bb28"+r.InteractionToken, "d57d6589"+r.MessageID)
	endpoint := EndpointDeleteFollowupMessage(bot.ApplicationID, r.InteractionToken, r.MessageID)

	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodDelete, endpoint, nil, nil, nil)
	if err != nil {
		return ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return nil
}

// Send sends a GetGuildAuditLog request to Discord and returns a AuditLog.
func (r *GetGuildAuditLog) Send(bot *Client) (*AuditLog, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[26]("26", "45892a5d"+r.GuildID)
	query, err := EndpointQueryString(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      "",
			Err:           err,
		}
	}
	endpoint := EndpointGetGuildAuditLog(r.GuildID) + "?" + query

	result := new(AuditLog)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, ContentTypeURLQueryString, nil, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a ListAutoModerationRulesForGuild request to Discord and returns a []*AutoModerationAction.
func (r *ListAutoModerationRulesForGuild) Send(bot *Client) ([]*AutoModerationAction, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[27]("27", "45892a5d"+r.GuildID)
	endpoint := EndpointListAutoModerationRulesForGuild(r.GuildID)

	result := make([]*AutoModerationAction, 0)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, nil, nil, &result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a GetAutoModerationRule request to Discord and returns a AutoModerationRule.
func (r *GetAutoModerationRule) Send(bot *Client) (*AutoModerationRule, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[28]("28", "45892a5d"+r.GuildID, "1b7efe5d"+r.AutoModerationRuleID)
	endpoint := EndpointGetAutoModerationRule(r.GuildID, r.AutoModerationRuleID)

	result := new(AutoModerationRule)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, nil, nil, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a CreateAutoModerationRule request to Discord and returns a AutoModerationRule.
func (r *CreateAutoModerationRule) Send(bot *Client) (*AutoModerationRule, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[29]("29", "45892a5d"+r.GuildID)
	endpoint := EndpointCreateAutoModerationRule(r.GuildID)

	body, err := json.Marshal(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           fmt.Errorf(errSendMarshal, err),
		}
	}

	result := new(AutoModerationRule)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPost, endpoint, ContentTypeJSON, body, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a ModifyAutoModerationRule request to Discord and returns a AutoModerationRule.
func (r *ModifyAutoModerationRule) Send(bot *Client) (*AutoModerationRule, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[30]("30", "45892a5d"+r.GuildID, "1b7efe5d"+r.AutoModerationRuleID)
	endpoint := EndpointModifyAutoModerationRule(r.GuildID, r.AutoModerationRuleID)

	body, err := json.Marshal(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           fmt.Errorf(errSendMarshal, err),
		}
	}

	result := new(AutoModerationRule)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPatch, endpoint, ContentTypeJSON, body, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a DeleteAutoModerationRule request to Discord and returns a error.
func (r *DeleteAutoModerationRule) Send(bot *Client) error {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[31]("31", "45892a5d"+r.GuildID, "1b7efe5d"+r.AutoModerationRuleID)
	endpoint := EndpointDeleteAutoModerationRule(r.GuildID, r.AutoModerationRuleID)

	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodDelete, endpoint, nil, nil, nil)
	if err != nil {
		return ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return nil
}

// Send sends a GetChannel request to Discord and returns a Channel.
func (r *GetChannel) Send(bot *Client) (*Channel, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[32]("32", "e5416649"+r.ChannelID)
	endpoint := EndpointGetChannel(r.ChannelID)

	result := new(Channel)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, nil, nil, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a ModifyChannel request to Discord and returns a Channel.
func (r *ModifyChannel) Send(bot *Client) (*Channel, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[33]("33", "e5416649"+r.ChannelID)
	endpoint := EndpointModifyChannel(r.ChannelID)

	result := new(Channel)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPatch, endpoint, nil, nil, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a ModifyChannelGroupDM request to Discord and returns a Channel.
func (r *ModifyChannelGroupDM) Send(bot *Client) (*Channel, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[34]("34", "e5416649"+r.ChannelID)
	endpoint := EndpointModifyChannelGroupDM(r.ChannelID)

	body, err := json.Marshal(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           fmt.Errorf(errSendMarshal, err),
		}
	}

	result := new(Channel)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPatch, endpoint, ContentTypeJSON, body, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a ModifyChannelGuild request to Discord and returns a Channel.
func (r *ModifyChannelGuild) Send(bot *Client) (*Channel, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[35]("35", "e5416649"+r.ChannelID)
	endpoint := EndpointModifyChannelGuild(r.ChannelID)

	body, err := json.Marshal(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           fmt.Errorf(errSendMarshal, err),
		}
	}

	result := new(Channel)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPatch, endpoint, ContentTypeJSON, body, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a ModifyChannelThread request to Discord and returns a Channel.
func (r *ModifyChannelThread) Send(bot *Client) (*Channel, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[36]("36", "e5416649"+r.ChannelID)
	endpoint := EndpointModifyChannelThread(r.ChannelID)

	body, err := json.Marshal(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           fmt.Errorf(errSendMarshal, err),
		}
	}

	result := new(Channel)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPatch, endpoint, ContentTypeJSON, body, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a DeleteCloseChannel request to Discord and returns a Channel.
func (r *DeleteCloseChannel) Send(bot *Client) (*Channel, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[37]("37", "e5416649"+r.ChannelID)
	endpoint := EndpointDeleteCloseChannel(r.ChannelID)

	result := new(Channel)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodDelete, endpoint, nil, nil, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a GetChannelMessages request to Discord and returns a []*Message.
func (r *GetChannelMessages) Send(bot *Client) ([]*Message, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[38]("38", "e5416649"+r.ChannelID)
	query, err := EndpointQueryString(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      "",
			Err:           err,
		}
	}
	endpoint := EndpointGetChannelMessages(r.ChannelID) + "?" + query

	result := make([]*Message, 0)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, ContentTypeURLQueryString, nil, &result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a GetChannelMessage request to Discord and returns a Message.
func (r *GetChannelMessage) Send(bot *Client) (*Message, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[39]("39", "e5416649"+r.ChannelID, "d57d6589"+r.MessageID)
	endpoint := EndpointGetChannelMessage(r.ChannelID, r.MessageID)

	result := new(Message)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, nil, nil, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a CreateMessage request to Discord and returns a Message.
func (r *CreateMessage) Send(bot *Client) (*Message, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[40]("40", "e5416649"+r.ChannelID)
	endpoint := EndpointCreateMessage(r.ChannelID)

	body, err := json.Marshal(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           fmt.Errorf(errSendMarshal, err),
		}
	}

	contentType := ContentTypeJSON
	if len(r.Files) != 0 {
		var multipartErr error
		if contentType, body, multipartErr = createMultipartForm(body, r.Files...); multipartErr != nil {
			return nil, ErrorRequest{
				ClientID:      bot.ApplicationID,
				CorrelationID: xid,
				RouteID:       routeid,
				ResourceID:    resourceid,
				Endpoint:      "",
				Err:           err,
			}
		}
	}

	result := new(Message)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPost, endpoint, contentType, body, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a CrosspostMessage request to Discord and returns a Message.
func (r *CrosspostMessage) Send(bot *Client) (*Message, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[41]("41", "e5416649"+r.ChannelID, "d57d6589"+r.MessageID)
	endpoint := EndpointCrosspostMessage(r.ChannelID, r.MessageID)

	result := new(Message)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPost, endpoint, nil, nil, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a CreateReaction request to Discord and returns a error.
func (r *CreateReaction) Send(bot *Client) error {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[42]("42", "e5416649"+r.ChannelID, "d57d6589"+r.MessageID, "033ebcdd"+r.Emoji)
	endpoint := EndpointCreateReaction(r.ChannelID, r.MessageID, r.Emoji)

	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPut, endpoint, nil, nil, nil)
	if err != nil {
		return ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return nil
}

// Send sends a DeleteOwnReaction request to Discord and returns a error.
func (r *DeleteOwnReaction) Send(bot *Client) error {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[43]("43", "e5416649"+r.ChannelID, "d57d6589"+r.MessageID, "033ebcdd"+r.Emoji)
	endpoint := EndpointDeleteOwnReaction(r.ChannelID, r.MessageID, r.Emoji)

	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodDelete, endpoint, nil, nil, nil)
	if err != nil {
		return ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return nil
}

// Send sends a DeleteUserReaction request to Discord and returns a error.
func (r *DeleteUserReaction) Send(bot *Client) error {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[44]("44", "e5416649"+r.ChannelID, "d57d6589"+r.MessageID, "033ebcdd"+r.Emoji, "209c92df"+r.UserID)
	endpoint := EndpointDeleteUserReaction(r.ChannelID, r.MessageID, r.Emoji, r.UserID)

	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodDelete, endpoint, nil, nil, nil)
	if err != nil {
		return ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return nil
}

// Send sends a GetReactions request to Discord and returns a []*User.
func (r *GetReactions) Send(bot *Client) ([]*User, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[45]("45", "e5416649"+r.ChannelID, "d57d6589"+r.MessageID, "033ebcdd"+r.Emoji)
	query, err := EndpointQueryString(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      "",
			Err:           err,
		}
	}
	endpoint := EndpointGetReactions(r.ChannelID, r.MessageID, r.Emoji) + "?" + query

	result := make([]*User, 0)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, ContentTypeURLQueryString, nil, &result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a DeleteAllReactions request to Discord and returns a error.
func (r *DeleteAllReactions) Send(bot *Client) error {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[46]("46", "e5416649"+r.ChannelID, "d57d6589"+r.MessageID)
	endpoint := EndpointDeleteAllReactions(r.ChannelID, r.MessageID)

	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodDelete, endpoint, nil, nil, nil)
	if err != nil {
		return ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return nil
}

// Send sends a DeleteAllReactionsforEmoji request to Discord and returns a error.
func (r *DeleteAllReactionsforEmoji) Send(bot *Client) error {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[47]("47", "e5416649"+r.ChannelID, "d57d6589"+r.MessageID, "033ebcdd"+r.Emoji)
	endpoint := EndpointDeleteAllReactionsforEmoji(r.ChannelID, r.MessageID, r.Emoji)

	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodDelete, endpoint, nil, nil, nil)
	if err != nil {
		return ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return nil
}

// Send sends a EditMessage request to Discord and returns a Message.
func (r *EditMessage) Send(bot *Client) (*Message, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[48]("48", "e5416649"+r.ChannelID, "d57d6589"+r.MessageID)
	endpoint := EndpointEditMessage(r.ChannelID, r.MessageID)

	body, err := json.Marshal(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           fmt.Errorf(errSendMarshal, err),
		}
	}

	contentType := ContentTypeJSON
	if len(r.Files) != 0 {
		var multipartErr error
		if contentType, body, multipartErr = createMultipartForm(body, r.Files...); multipartErr != nil {
			return nil, ErrorRequest{
				ClientID:      bot.ApplicationID,
				CorrelationID: xid,
				RouteID:       routeid,
				ResourceID:    resourceid,
				Endpoint:      "",
				Err:           err,
			}
		}
	}

	result := new(Message)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPatch, endpoint, contentType, body, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a DeleteMessage request to Discord and returns a error.
func (r *DeleteMessage) Send(bot *Client) error {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[49]("49", "e5416649"+r.ChannelID, "d57d6589"+r.MessageID)
	endpoint := EndpointDeleteMessage(r.ChannelID, r.MessageID)

	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodDelete, endpoint, nil, nil, nil)
	if err != nil {
		return ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return nil
}

// Send sends a BulkDeleteMessages request to Discord and returns a error.
func (r *BulkDeleteMessages) Send(bot *Client) error {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[50]("50", "e5416649"+r.ChannelID)
	endpoint := EndpointBulkDeleteMessages(r.ChannelID)

	body, err := json.Marshal(r)
	if err != nil {
		return ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           fmt.Errorf(errSendMarshal, err),
		}
	}

	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPost, endpoint, ContentTypeJSON, body, nil)
	if err != nil {
		return ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return nil
}

// Send sends a EditChannelPermissions request to Discord and returns a error.
func (r *EditChannelPermissions) Send(bot *Client) error {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[51]("51", "e5416649"+r.ChannelID, "9167175f"+r.OverwriteID)
	endpoint := EndpointEditChannelPermissions(r.ChannelID, r.OverwriteID)

	body, err := json.Marshal(r)
	if err != nil {
		return ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           fmt.Errorf(errSendMarshal, err),
		}
	}

	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPut, endpoint, ContentTypeJSON, body, nil)
	if err != nil {
		return ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return nil
}

// Send sends a GetChannelInvites request to Discord and returns a []*Invite.
func (r *GetChannelInvites) Send(bot *Client) ([]*Invite, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[52]("52", "e5416649"+r.ChannelID)
	endpoint := EndpointGetChannelInvites(r.ChannelID)

	result := make([]*Invite, 0)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, nil, nil, &result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a CreateChannelInvite request to Discord and returns a Invite.
func (r *CreateChannelInvite) Send(bot *Client) (*Invite, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[53]("53", "e5416649"+r.ChannelID)
	endpoint := EndpointCreateChannelInvite(r.ChannelID)

	body, err := json.Marshal(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           fmt.Errorf(errSendMarshal, err),
		}
	}

	result := new(Invite)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPost, endpoint, ContentTypeJSON, body, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a DeleteChannelPermission request to Discord and returns a error.
func (r *DeleteChannelPermission) Send(bot *Client) error {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[54]("54", "e5416649"+r.ChannelID, "9167175f"+r.OverwriteID)
	endpoint := EndpointDeleteChannelPermission(r.ChannelID, r.OverwriteID)

	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodDelete, endpoint, nil, nil, nil)
	if err != nil {
		return ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return nil
}

// Send sends a FollowAnnouncementChannel request to Discord and returns a FollowedChannel.
func (r *FollowAnnouncementChannel) Send(bot *Client) (*FollowedChannel, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[55]("55", "e5416649"+r.ChannelID)
	endpoint := EndpointFollowAnnouncementChannel(r.ChannelID)

	body, err := json.Marshal(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           fmt.Errorf(errSendMarshal, err),
		}
	}

	result := new(FollowedChannel)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPost, endpoint, ContentTypeJSON, body, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a TriggerTypingIndicator request to Discord and returns a error.
func (r *TriggerTypingIndicator) Send(bot *Client) error {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[56]("56", "e5416649"+r.ChannelID)
	endpoint := EndpointTriggerTypingIndicator(r.ChannelID)

	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPost, endpoint, nil, nil, nil)
	if err != nil {
		return ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return nil
}

// Send sends a GetPinnedMessages request to Discord and returns a []*Message.
func (r *GetPinnedMessages) Send(bot *Client) ([]*Message, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[57]("57", "e5416649"+r.ChannelID)
	endpoint := EndpointGetPinnedMessages(r.ChannelID)

	result := make([]*Message, 0)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, nil, nil, &result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a PinMessage request to Discord and returns a error.
func (r *PinMessage) Send(bot *Client) error {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[58]("58", "e5416649"+r.ChannelID, "d57d6589"+r.MessageID)
	endpoint := EndpointPinMessage(r.ChannelID, r.MessageID)

	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPut, endpoint, nil, nil, nil)
	if err != nil {
		return ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return nil
}

// Send sends a UnpinMessage request to Discord and returns a error.
func (r *UnpinMessage) Send(bot *Client) error {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[59]("59", "e5416649"+r.ChannelID, "d57d6589"+r.MessageID)
	endpoint := EndpointUnpinMessage(r.ChannelID, r.MessageID)

	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodDelete, endpoint, nil, nil, nil)
	if err != nil {
		return ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return nil
}

// Send sends a GroupDMAddRecipient request to Discord and returns a error.
func (r *GroupDMAddRecipient) Send(bot *Client) error {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[60]("60", "e5416649"+r.ChannelID, "209c92df"+r.UserID)
	endpoint := EndpointGroupDMAddRecipient(r.ChannelID, r.UserID)

	body, err := json.Marshal(r)
	if err != nil {
		return ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           fmt.Errorf(errSendMarshal, err),
		}
	}

	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPut, endpoint, ContentTypeJSON, body, nil)
	if err != nil {
		return ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return nil
}

// Send sends a GroupDMRemoveRecipient request to Discord and returns a error.
func (r *GroupDMRemoveRecipient) Send(bot *Client) error {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[61]("61", "e5416649"+r.ChannelID, "209c92df"+r.UserID)
	endpoint := EndpointGroupDMRemoveRecipient(r.ChannelID, r.UserID)

	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodDelete, endpoint, nil, nil, nil)
	if err != nil {
		return ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return nil
}

// Send sends a StartThreadfromMessage request to Discord and returns a Channel.
func (r *StartThreadfromMessage) Send(bot *Client) (*Channel, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[62]("62", "e5416649"+r.ChannelID, "d57d6589"+r.MessageID)
	endpoint := EndpointStartThreadfromMessage(r.ChannelID, r.MessageID)

	body, err := json.Marshal(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           fmt.Errorf(errSendMarshal, err),
		}
	}

	result := new(Channel)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPost, endpoint, ContentTypeJSON, body, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a StartThreadwithoutMessage request to Discord and returns a Channel.
func (r *StartThreadwithoutMessage) Send(bot *Client) (*Channel, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[63]("63", "e5416649"+r.ChannelID)
	endpoint := EndpointStartThreadwithoutMessage(r.ChannelID)

	body, err := json.Marshal(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           fmt.Errorf(errSendMarshal, err),
		}
	}

	result := new(Channel)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPost, endpoint, ContentTypeJSON, body, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a StartThreadinForumChannel request to Discord and returns a Channel.
func (r *StartThreadinForumChannel) Send(bot *Client) (*Channel, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[64]("64", "e5416649"+r.ChannelID)
	endpoint := EndpointStartThreadinForumChannel(r.ChannelID)

	body, err := json.Marshal(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           fmt.Errorf(errSendMarshal, err),
		}
	}

	result := new(Channel)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPost, endpoint, ContentTypeJSON, body, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a JoinThread request to Discord and returns a error.
func (r *JoinThread) Send(bot *Client) error {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[65]("65", "e5416649"+r.ChannelID)
	endpoint := EndpointJoinThread(r.ChannelID)

	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPut, endpoint, nil, nil, nil)
	if err != nil {
		return ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return nil
}

// Send sends a AddThreadMember request to Discord and returns a error.
func (r *AddThreadMember) Send(bot *Client) error {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[66]("66", "e5416649"+r.ChannelID, "209c92df"+r.UserID)
	endpoint := EndpointAddThreadMember(r.ChannelID, r.UserID)

	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPut, endpoint, nil, nil, nil)
	if err != nil {
		return ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return nil
}

// Send sends a LeaveThread request to Discord and returns a error.
func (r *LeaveThread) Send(bot *Client) error {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[67]("67", "e5416649"+r.ChannelID)
	endpoint := EndpointLeaveThread(r.ChannelID)

	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodDelete, endpoint, nil, nil, nil)
	if err != nil {
		return ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return nil
}

// Send sends a RemoveThreadMember request to Discord and returns a error.
func (r *RemoveThreadMember) Send(bot *Client) error {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[68]("68", "e5416649"+r.ChannelID, "209c92df"+r.UserID)
	endpoint := EndpointRemoveThreadMember(r.ChannelID, r.UserID)

	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodDelete, endpoint, nil, nil, nil)
	if err != nil {
		return ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return nil
}

// Send sends a GetThreadMember request to Discord and returns a ThreadMember.
func (r *GetThreadMember) Send(bot *Client) (*ThreadMember, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[69]("69", "e5416649"+r.ChannelID, "209c92df"+r.UserID)
	endpoint := EndpointGetThreadMember(r.ChannelID, r.UserID)

	result := new(ThreadMember)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, nil, nil, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a ListThreadMembers request to Discord and returns a []*ThreadMember.
func (r *ListThreadMembers) Send(bot *Client) ([]*ThreadMember, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[70]("70", "e5416649"+r.ChannelID)
	endpoint := EndpointListThreadMembers(r.ChannelID)

	result := make([]*ThreadMember, 0)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, nil, nil, &result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a ListPublicArchivedThreads request to Discord and returns a ListPublicArchivedThreadsResponse.
func (r *ListPublicArchivedThreads) Send(bot *Client) (*ListPublicArchivedThreadsResponse, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[71]("71", "e5416649"+r.ChannelID)
	query, err := EndpointQueryString(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      "",
			Err:           err,
		}
	}
	endpoint := EndpointListPublicArchivedThreads(r.ChannelID) + "?" + query

	result := new(ListPublicArchivedThreadsResponse)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, ContentTypeURLQueryString, nil, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a ListPrivateArchivedThreads request to Discord and returns a ListPrivateArchivedThreadsResponse.
func (r *ListPrivateArchivedThreads) Send(bot *Client) (*ListPrivateArchivedThreadsResponse, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[72]("72", "e5416649"+r.ChannelID)
	query, err := EndpointQueryString(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      "",
			Err:           err,
		}
	}
	endpoint := EndpointListPrivateArchivedThreads(r.ChannelID) + "?" + query

	result := new(ListPrivateArchivedThreadsResponse)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, ContentTypeURLQueryString, nil, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a ListJoinedPrivateArchivedThreads request to Discord and returns a ListJoinedPrivateArchivedThreadsResponse.
func (r *ListJoinedPrivateArchivedThreads) Send(bot *Client) (*ListJoinedPrivateArchivedThreadsResponse, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[73]("73", "e5416649"+r.ChannelID)
	query, err := EndpointQueryString(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      "",
			Err:           err,
		}
	}
	endpoint := EndpointListJoinedPrivateArchivedThreads(r.ChannelID) + "?" + query

	result := new(ListJoinedPrivateArchivedThreadsResponse)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, ContentTypeURLQueryString, nil, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a ListGuildEmojis request to Discord and returns a []*Emoji.
func (r *ListGuildEmojis) Send(bot *Client) ([]*Emoji, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[74]("74", "45892a5d"+r.GuildID)
	endpoint := EndpointListGuildEmojis(r.GuildID)

	result := make([]*Emoji, 0)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, nil, nil, &result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a GetGuildEmoji request to Discord and returns a Emoji.
func (r *GetGuildEmoji) Send(bot *Client) (*Emoji, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[75]("75", "45892a5d"+r.GuildID, "67c175a8"+r.EmojiID)
	endpoint := EndpointGetGuildEmoji(r.GuildID, r.EmojiID)

	result := new(Emoji)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, nil, nil, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a CreateGuildEmoji request to Discord and returns a Emoji.
func (r *CreateGuildEmoji) Send(bot *Client) (*Emoji, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[76]("76", "45892a5d"+r.GuildID)
	endpoint := EndpointCreateGuildEmoji(r.GuildID)

	body, err := json.Marshal(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           fmt.Errorf(errSendMarshal, err),
		}
	}

	result := new(Emoji)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPost, endpoint, ContentTypeJSON, body, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a ModifyGuildEmoji request to Discord and returns a Emoji.
func (r *ModifyGuildEmoji) Send(bot *Client) (*Emoji, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[77]("77", "45892a5d"+r.GuildID, "67c175a8"+r.EmojiID)
	endpoint := EndpointModifyGuildEmoji(r.GuildID, r.EmojiID)

	body, err := json.Marshal(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           fmt.Errorf(errSendMarshal, err),
		}
	}

	result := new(Emoji)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPatch, endpoint, ContentTypeJSON, body, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a DeleteGuildEmoji request to Discord and returns a error.
func (r *DeleteGuildEmoji) Send(bot *Client) error {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[78]("78", "45892a5d"+r.GuildID, "67c175a8"+r.EmojiID)
	endpoint := EndpointDeleteGuildEmoji(r.GuildID, r.EmojiID)

	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodDelete, endpoint, nil, nil, nil)
	if err != nil {
		return ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return nil
}

// Send sends a CreateGuild request to Discord and returns a Guild.
func (r *CreateGuild) Send(bot *Client) (*Guild, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[79]("79")
	endpoint := EndpointCreateGuild()

	body, err := json.Marshal(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           fmt.Errorf(errSendMarshal, err),
		}
	}

	result := new(Guild)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPost, endpoint, ContentTypeJSON, body, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a GetGuild request to Discord and returns a Guild.
func (r *GetGuild) Send(bot *Client) (*Guild, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[80]("80", "45892a5d"+r.GuildID)
	query, err := EndpointQueryString(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      "",
			Err:           err,
		}
	}
	endpoint := EndpointGetGuild(r.GuildID) + "?" + query

	result := new(Guild)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, ContentTypeURLQueryString, nil, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a GetGuildPreview request to Discord and returns a GuildPreview.
func (r *GetGuildPreview) Send(bot *Client) (*GuildPreview, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[81]("81", "45892a5d"+r.GuildID)
	endpoint := EndpointGetGuildPreview(r.GuildID)

	result := new(GuildPreview)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, nil, nil, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a ModifyGuild request to Discord and returns a Guild.
func (r *ModifyGuild) Send(bot *Client) (*Guild, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[82]("82", "45892a5d"+r.GuildID)
	endpoint := EndpointModifyGuild(r.GuildID)

	body, err := json.Marshal(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           fmt.Errorf(errSendMarshal, err),
		}
	}

	result := new(Guild)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPatch, endpoint, ContentTypeJSON, body, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a DeleteGuild request to Discord and returns a error.
func (r *DeleteGuild) Send(bot *Client) error {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[83]("83", "45892a5d"+r.GuildID)
	endpoint := EndpointDeleteGuild(r.GuildID)

	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodDelete, endpoint, nil, nil, nil)
	if err != nil {
		return ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return nil
}

// Send sends a GetGuildChannels request to Discord and returns a []*Channel.
func (r *GetGuildChannels) Send(bot *Client) ([]*Channel, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[84]("84", "45892a5d"+r.GuildID)
	endpoint := EndpointGetGuildChannels(r.GuildID)

	result := make([]*Channel, 0)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, nil, nil, &result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a CreateGuildChannel request to Discord and returns a Channel.
func (r *CreateGuildChannel) Send(bot *Client) (*Channel, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[85]("85", "45892a5d"+r.GuildID)
	endpoint := EndpointCreateGuildChannel(r.GuildID)

	body, err := json.Marshal(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           fmt.Errorf(errSendMarshal, err),
		}
	}

	result := new(Channel)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPost, endpoint, ContentTypeJSON, body, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a ModifyGuildChannelPositions request to Discord and returns a error.
func (r *ModifyGuildChannelPositions) Send(bot *Client) error {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[86]("86", "45892a5d"+r.GuildID)
	endpoint := EndpointModifyGuildChannelPositions(r.GuildID)

	body, err := json.Marshal(r)
	if err != nil {
		return ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           fmt.Errorf(errSendMarshal, err),
		}
	}

	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPatch, endpoint, ContentTypeJSON, body, nil)
	if err != nil {
		return ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return nil
}

// Send sends a ListActiveGuildThreads request to Discord and returns a ListActiveGuildThreadsResponse.
func (r *ListActiveGuildThreads) Send(bot *Client) (*ListActiveGuildThreadsResponse, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[87]("87", "45892a5d"+r.GuildID)
	endpoint := EndpointListActiveGuildThreads(r.GuildID)

	result := new(ListActiveGuildThreadsResponse)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, nil, nil, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a GetGuildMember request to Discord and returns a GuildMember.
func (r *GetGuildMember) Send(bot *Client) (*GuildMember, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[88]("88", "45892a5d"+r.GuildID, "209c92df"+r.UserID)
	endpoint := EndpointGetGuildMember(r.GuildID, r.UserID)

	result := new(GuildMember)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, nil, nil, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a ListGuildMembers request to Discord and returns a []*GuildMember.
func (r *ListGuildMembers) Send(bot *Client) ([]*GuildMember, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[89]("89", "45892a5d"+r.GuildID)
	query, err := EndpointQueryString(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      "",
			Err:           err,
		}
	}
	endpoint := EndpointListGuildMembers(r.GuildID) + "?" + query

	result := make([]*GuildMember, 0)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, ContentTypeURLQueryString, nil, &result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a SearchGuildMembers request to Discord and returns a []*GuildMember.
func (r *SearchGuildMembers) Send(bot *Client) ([]*GuildMember, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[90]("90", "45892a5d"+r.GuildID)
	query, err := EndpointQueryString(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      "",
			Err:           err,
		}
	}
	endpoint := EndpointSearchGuildMembers(r.GuildID) + "?" + query

	result := make([]*GuildMember, 0)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, ContentTypeURLQueryString, nil, &result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a AddGuildMember request to Discord and returns a GuildMember.
func (r *AddGuildMember) Send(bot *Client) (*GuildMember, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[91]("91", "45892a5d"+r.GuildID, "209c92df"+r.UserID)
	endpoint := EndpointAddGuildMember(r.GuildID, r.UserID)

	body, err := json.Marshal(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           fmt.Errorf(errSendMarshal, err),
		}
	}

	result := new(GuildMember)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPut, endpoint, ContentTypeJSON, body, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a ModifyGuildMember request to Discord and returns a GuildMember.
func (r *ModifyGuildMember) Send(bot *Client) (*GuildMember, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[92]("92", "45892a5d"+r.GuildID, "209c92df"+r.UserID)
	endpoint := EndpointModifyGuildMember(r.GuildID, r.UserID)

	body, err := json.Marshal(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           fmt.Errorf(errSendMarshal, err),
		}
	}

	result := new(GuildMember)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPatch, endpoint, ContentTypeJSON, body, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a ModifyCurrentMember request to Discord and returns a GuildMember.
func (r *ModifyCurrentMember) Send(bot *Client) (*GuildMember, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[93]("93", "45892a5d"+r.GuildID)
	endpoint := EndpointModifyCurrentMember(r.GuildID)

	body, err := json.Marshal(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           fmt.Errorf(errSendMarshal, err),
		}
	}

	result := new(GuildMember)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPatch, endpoint, ContentTypeJSON, body, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a AddGuildMemberRole request to Discord and returns a error.
func (r *AddGuildMemberRole) Send(bot *Client) error {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[94]("94", "45892a5d"+r.GuildID, "209c92df"+r.UserID, "3cf7dd7c"+r.RoleID)
	endpoint := EndpointAddGuildMemberRole(r.GuildID, r.UserID, r.RoleID)

	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPut, endpoint, nil, nil, nil)
	if err != nil {
		return ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return nil
}

// Send sends a RemoveGuildMemberRole request to Discord and returns a error.
func (r *RemoveGuildMemberRole) Send(bot *Client) error {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[95]("95", "45892a5d"+r.GuildID, "209c92df"+r.UserID, "3cf7dd7c"+r.RoleID)
	endpoint := EndpointRemoveGuildMemberRole(r.GuildID, r.UserID, r.RoleID)

	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodDelete, endpoint, nil, nil, nil)
	if err != nil {
		return ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return nil
}

// Send sends a RemoveGuildMember request to Discord and returns a error.
func (r *RemoveGuildMember) Send(bot *Client) error {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[96]("96", "45892a5d"+r.GuildID, "209c92df"+r.UserID)
	endpoint := EndpointRemoveGuildMember(r.GuildID, r.UserID)

	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodDelete, endpoint, nil, nil, nil)
	if err != nil {
		return ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return nil
}

// Send sends a GetGuildBans request to Discord and returns a []*Ban.
func (r *GetGuildBans) Send(bot *Client) ([]*Ban, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[97]("97", "45892a5d"+r.GuildID)
	query, err := EndpointQueryString(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      "",
			Err:           err,
		}
	}
	endpoint := EndpointGetGuildBans(r.GuildID) + "?" + query

	result := make([]*Ban, 0)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, ContentTypeURLQueryString, nil, &result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a GetGuildBan request to Discord and returns a Ban.
func (r *GetGuildBan) Send(bot *Client) (*Ban, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[98]("98", "45892a5d"+r.GuildID, "209c92df"+r.UserID)
	endpoint := EndpointGetGuildBan(r.GuildID, r.UserID)

	result := new(Ban)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, nil, nil, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a CreateGuildBan request to Discord and returns a error.
func (r *CreateGuildBan) Send(bot *Client) error {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[99]("99", "45892a5d"+r.GuildID, "209c92df"+r.UserID)
	endpoint := EndpointCreateGuildBan(r.GuildID, r.UserID)

	body, err := json.Marshal(r)
	if err != nil {
		return ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           fmt.Errorf(errSendMarshal, err),
		}
	}

	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPut, endpoint, ContentTypeJSON, body, nil)
	if err != nil {
		return ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return nil
}

// Send sends a RemoveGuildBan request to Discord and returns a error.
func (r *RemoveGuildBan) Send(bot *Client) error {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[100]("100", "45892a5d"+r.GuildID, "209c92df"+r.UserID)
	endpoint := EndpointRemoveGuildBan(r.GuildID, r.UserID)

	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodDelete, endpoint, nil, nil, nil)
	if err != nil {
		return ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return nil
}

// Send sends a GetGuildRoles request to Discord and returns a []*Role.
func (r *GetGuildRoles) Send(bot *Client) ([]*Role, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[101]("101", "45892a5d"+r.GuildID)
	endpoint := EndpointGetGuildRoles(r.GuildID)

	result := make([]*Role, 0)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, nil, nil, &result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a CreateGuildRole request to Discord and returns a Role.
func (r *CreateGuildRole) Send(bot *Client) (*Role, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[102]("102", "45892a5d"+r.GuildID)
	endpoint := EndpointCreateGuildRole(r.GuildID)

	body, err := json.Marshal(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           fmt.Errorf(errSendMarshal, err),
		}
	}

	result := new(Role)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPost, endpoint, ContentTypeJSON, body, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a ModifyGuildRolePositions request to Discord and returns a []*Role.
func (r *ModifyGuildRolePositions) Send(bot *Client) ([]*Role, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[103]("103", "45892a5d"+r.GuildID)
	endpoint := EndpointModifyGuildRolePositions(r.GuildID)

	body, err := json.Marshal(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           fmt.Errorf(errSendMarshal, err),
		}
	}

	result := make([]*Role, 0)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPatch, endpoint, ContentTypeJSON, body, &result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a ModifyGuildRole request to Discord and returns a Role.
func (r *ModifyGuildRole) Send(bot *Client) (*Role, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[104]("104", "45892a5d"+r.GuildID, "3cf7dd7c"+r.RoleID)
	endpoint := EndpointModifyGuildRole(r.GuildID, r.RoleID)

	body, err := json.Marshal(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           fmt.Errorf(errSendMarshal, err),
		}
	}

	result := new(Role)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPatch, endpoint, ContentTypeJSON, body, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a DeleteGuildRole request to Discord and returns a error.
func (r *DeleteGuildRole) Send(bot *Client) error {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[105]("105", "45892a5d"+r.GuildID, "3cf7dd7c"+r.RoleID)
	endpoint := EndpointDeleteGuildRole(r.GuildID, r.RoleID)

	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodDelete, endpoint, nil, nil, nil)
	if err != nil {
		return ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return nil
}

// Send sends a ModifyGuildMFALevel request to Discord and returns a ModifyGuildMFALevelResponse.
func (r *ModifyGuildMFALevel) Send(bot *Client) (*ModifyGuildMFALevelResponse, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[106]("106", "45892a5d"+r.GuildID)
	endpoint := EndpointModifyGuildMFALevel(r.GuildID)

	body, err := json.Marshal(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           fmt.Errorf(errSendMarshal, err),
		}
	}

	result := new(ModifyGuildMFALevelResponse)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPost, endpoint, ContentTypeJSON, body, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a GetGuildPruneCount request to Discord and returns a GetGuildPruneCountResponse.
func (r *GetGuildPruneCount) Send(bot *Client) (*GetGuildPruneCountResponse, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[107]("107", "45892a5d"+r.GuildID)
	query, err := EndpointQueryString(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      "",
			Err:           err,
		}
	}
	endpoint := EndpointGetGuildPruneCount(r.GuildID) + "?" + query

	result := new(GetGuildPruneCountResponse)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, ContentTypeURLQueryString, nil, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a BeginGuildPrune request to Discord and returns a error.
func (r *BeginGuildPrune) Send(bot *Client) error {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[108]("108", "45892a5d"+r.GuildID)
	endpoint := EndpointBeginGuildPrune(r.GuildID)

	body, err := json.Marshal(r)
	if err != nil {
		return ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           fmt.Errorf(errSendMarshal, err),
		}
	}

	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPost, endpoint, ContentTypeJSON, body, nil)
	if err != nil {
		return ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return nil
}

// Send sends a GetGuildVoiceRegions request to Discord and returns a []*VoiceRegion.
func (r *GetGuildVoiceRegions) Send(bot *Client) ([]*VoiceRegion, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[109]("109", "45892a5d"+r.GuildID)
	endpoint := EndpointGetGuildVoiceRegions(r.GuildID)

	result := make([]*VoiceRegion, 0)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, nil, nil, &result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a GetGuildInvites request to Discord and returns a []*Invite.
func (r *GetGuildInvites) Send(bot *Client) ([]*Invite, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[110]("110", "45892a5d"+r.GuildID)
	endpoint := EndpointGetGuildInvites(r.GuildID)

	result := make([]*Invite, 0)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, nil, nil, &result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a GetGuildIntegrations request to Discord and returns a []*Integration.
func (r *GetGuildIntegrations) Send(bot *Client) ([]*Integration, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[111]("111", "45892a5d"+r.GuildID)
	endpoint := EndpointGetGuildIntegrations(r.GuildID)

	result := make([]*Integration, 0)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, nil, nil, &result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a DeleteGuildIntegration request to Discord and returns a error.
func (r *DeleteGuildIntegration) Send(bot *Client) error {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[112]("112", "45892a5d"+r.GuildID, "cb4479f8"+r.IntegrationID)
	endpoint := EndpointDeleteGuildIntegration(r.GuildID, r.IntegrationID)

	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodDelete, endpoint, nil, nil, nil)
	if err != nil {
		return ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return nil
}

// Send sends a GetGuildWidgetSettings request to Discord and returns a GuildWidget.
func (r *GetGuildWidgetSettings) Send(bot *Client) (*GuildWidget, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[113]("113", "45892a5d"+r.GuildID)
	endpoint := EndpointGetGuildWidgetSettings(r.GuildID)

	result := new(GuildWidget)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, nil, nil, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a ModifyGuildWidget request to Discord and returns a GuildWidget.
func (r *ModifyGuildWidget) Send(bot *Client) (*GuildWidget, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[114]("114", "45892a5d"+r.GuildID)
	endpoint := EndpointModifyGuildWidget(r.GuildID)

	result := new(GuildWidget)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPatch, endpoint, nil, nil, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a GetGuildWidget request to Discord and returns a GuildWidget.
func (r *GetGuildWidget) Send(bot *Client) (*GuildWidget, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[115]("115", "45892a5d"+r.GuildID)
	endpoint := EndpointGetGuildWidget(r.GuildID)

	result := new(GuildWidget)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, nil, nil, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a GetGuildVanityURL request to Discord and returns a Invite.
func (r *GetGuildVanityURL) Send(bot *Client) (*Invite, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[116]("116", "45892a5d"+r.GuildID)
	endpoint := EndpointGetGuildVanityURL(r.GuildID)

	body, err := json.Marshal(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           fmt.Errorf(errSendMarshal, err),
		}
	}

	result := new(Invite)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, ContentTypeJSON, body, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a GetGuildWidgetImage request to Discord and returns a EmbedImage.
func (r *GetGuildWidgetImage) Send(bot *Client) (*EmbedImage, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[117]("117", "45892a5d"+r.GuildID)
	query, err := EndpointQueryString(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      "",
			Err:           err,
		}
	}
	endpoint := EndpointGetGuildWidgetImage(r.GuildID) + "?" + query

	result := new(EmbedImage)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, ContentTypeURLQueryString, nil, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a GetGuildWelcomeScreen request to Discord and returns a WelcomeScreen.
func (r *GetGuildWelcomeScreen) Send(bot *Client) (*WelcomeScreen, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[118]("118", "45892a5d"+r.GuildID)
	endpoint := EndpointGetGuildWelcomeScreen(r.GuildID)

	result := new(WelcomeScreen)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, nil, nil, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a ModifyGuildWelcomeScreen request to Discord and returns a WelcomeScreen.
func (r *ModifyGuildWelcomeScreen) Send(bot *Client) (*WelcomeScreen, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[119]("119", "45892a5d"+r.GuildID)
	endpoint := EndpointModifyGuildWelcomeScreen(r.GuildID)

	body, err := json.Marshal(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           fmt.Errorf(errSendMarshal, err),
		}
	}

	result := new(WelcomeScreen)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPatch, endpoint, ContentTypeJSON, body, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a ModifyCurrentUserVoiceState request to Discord and returns a error.
func (r *ModifyCurrentUserVoiceState) Send(bot *Client) error {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[120]("120", "45892a5d"+r.GuildID)
	endpoint := EndpointModifyCurrentUserVoiceState(r.GuildID)

	body, err := json.Marshal(r)
	if err != nil {
		return ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           fmt.Errorf(errSendMarshal, err),
		}
	}

	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPatch, endpoint, ContentTypeJSON, body, nil)
	if err != nil {
		return ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return nil
}

// Send sends a ModifyUserVoiceState request to Discord and returns a error.
func (r *ModifyUserVoiceState) Send(bot *Client) error {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[121]("121", "45892a5d"+r.GuildID, "209c92df"+r.UserID)
	endpoint := EndpointModifyUserVoiceState(r.GuildID, r.UserID)

	body, err := json.Marshal(r)
	if err != nil {
		return ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           fmt.Errorf(errSendMarshal, err),
		}
	}

	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPatch, endpoint, ContentTypeJSON, body, nil)
	if err != nil {
		return ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return nil
}

// Send sends a ListScheduledEventsforGuild request to Discord and returns a []*GuildScheduledEvent.
func (r *ListScheduledEventsforGuild) Send(bot *Client) ([]*GuildScheduledEvent, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[122]("122", "45892a5d"+r.GuildID)
	query, err := EndpointQueryString(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      "",
			Err:           err,
		}
	}
	endpoint := EndpointListScheduledEventsforGuild(r.GuildID) + "?" + query

	result := make([]*GuildScheduledEvent, 0)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, ContentTypeURLQueryString, nil, &result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a CreateGuildScheduledEvent request to Discord and returns a GuildScheduledEvent.
func (r *CreateGuildScheduledEvent) Send(bot *Client) (*GuildScheduledEvent, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[123]("123", "45892a5d"+r.GuildID)
	endpoint := EndpointCreateGuildScheduledEvent(r.GuildID)

	body, err := json.Marshal(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           fmt.Errorf(errSendMarshal, err),
		}
	}

	result := new(GuildScheduledEvent)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPost, endpoint, ContentTypeJSON, body, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a GetGuildScheduledEvent request to Discord and returns a GuildScheduledEvent.
func (r *GetGuildScheduledEvent) Send(bot *Client) (*GuildScheduledEvent, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[124]("124", "45892a5d"+r.GuildID, "522412fc"+r.GuildScheduledEventID)
	query, err := EndpointQueryString(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      "",
			Err:           err,
		}
	}
	endpoint := EndpointGetGuildScheduledEvent(r.GuildID, r.GuildScheduledEventID) + "?" + query

	result := new(GuildScheduledEvent)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, ContentTypeURLQueryString, nil, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a ModifyGuildScheduledEvent request to Discord and returns a GuildScheduledEvent.
func (r *ModifyGuildScheduledEvent) Send(bot *Client) (*GuildScheduledEvent, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[125]("125", "45892a5d"+r.GuildID, "522412fc"+r.GuildScheduledEventID)
	endpoint := EndpointModifyGuildScheduledEvent(r.GuildID, r.GuildScheduledEventID)

	body, err := json.Marshal(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           fmt.Errorf(errSendMarshal, err),
		}
	}

	result := new(GuildScheduledEvent)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPatch, endpoint, ContentTypeJSON, body, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a DeleteGuildScheduledEvent request to Discord and returns a error.
func (r *DeleteGuildScheduledEvent) Send(bot *Client) error {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[126]("126", "45892a5d"+r.GuildID, "522412fc"+r.GuildScheduledEventID)
	endpoint := EndpointDeleteGuildScheduledEvent(r.GuildID, r.GuildScheduledEventID)

	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodDelete, endpoint, nil, nil, nil)
	if err != nil {
		return ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return nil
}

// Send sends a GetGuildScheduledEventUsers request to Discord and returns a []*GuildScheduledEventUser.
func (r *GetGuildScheduledEventUsers) Send(bot *Client) ([]*GuildScheduledEventUser, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[127]("127", "45892a5d"+r.GuildID, "522412fc"+r.GuildScheduledEventID)
	query, err := EndpointQueryString(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      "",
			Err:           err,
		}
	}
	endpoint := EndpointGetGuildScheduledEventUsers(r.GuildID, r.GuildScheduledEventID) + "?" + query

	result := make([]*GuildScheduledEventUser, 0)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, ContentTypeURLQueryString, nil, &result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a GetGuildTemplate request to Discord and returns a GuildTemplate.
func (r *GetGuildTemplate) Send(bot *Client) (*GuildTemplate, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[128]("128", "61437152"+r.TemplateCode)
	endpoint := EndpointGetGuildTemplate(r.TemplateCode)

	result := new(GuildTemplate)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, nil, nil, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a CreateGuildfromGuildTemplate request to Discord and returns a []*GuildTemplate.
func (r *CreateGuildfromGuildTemplate) Send(bot *Client) ([]*GuildTemplate, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[129]("129", "61437152"+r.TemplateCode)
	endpoint := EndpointCreateGuildfromGuildTemplate(r.TemplateCode)

	body, err := json.Marshal(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           fmt.Errorf(errSendMarshal, err),
		}
	}

	result := make([]*GuildTemplate, 0)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPost, endpoint, ContentTypeJSON, body, &result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a GetGuildTemplates request to Discord and returns a []*GuildTemplate.
func (r *GetGuildTemplates) Send(bot *Client) ([]*GuildTemplate, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[130]("130", "45892a5d"+r.GuildID)
	endpoint := EndpointGetGuildTemplates(r.GuildID)

	result := make([]*GuildTemplate, 0)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, nil, nil, &result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a CreateGuildTemplate request to Discord and returns a GuildTemplate.
func (r *CreateGuildTemplate) Send(bot *Client) (*GuildTemplate, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[131]("131", "45892a5d"+r.GuildID)
	endpoint := EndpointCreateGuildTemplate(r.GuildID)

	body, err := json.Marshal(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           fmt.Errorf(errSendMarshal, err),
		}
	}

	result := new(GuildTemplate)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPost, endpoint, ContentTypeJSON, body, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a SyncGuildTemplate request to Discord and returns a GuildTemplate.
func (r *SyncGuildTemplate) Send(bot *Client) (*GuildTemplate, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[132]("132", "45892a5d"+r.GuildID, "61437152"+r.TemplateCode)
	endpoint := EndpointSyncGuildTemplate(r.GuildID, r.TemplateCode)

	result := new(GuildTemplate)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPut, endpoint, nil, nil, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a ModifyGuildTemplate request to Discord and returns a GuildTemplate.
func (r *ModifyGuildTemplate) Send(bot *Client) (*GuildTemplate, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[133]("133", "45892a5d"+r.GuildID, "61437152"+r.TemplateCode)
	endpoint := EndpointModifyGuildTemplate(r.GuildID, r.TemplateCode)

	body, err := json.Marshal(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           fmt.Errorf(errSendMarshal, err),
		}
	}

	result := new(GuildTemplate)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPatch, endpoint, ContentTypeJSON, body, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a DeleteGuildTemplate request to Discord and returns a GuildTemplate.
func (r *DeleteGuildTemplate) Send(bot *Client) (*GuildTemplate, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[134]("134", "45892a5d"+r.GuildID, "61437152"+r.TemplateCode)
	endpoint := EndpointDeleteGuildTemplate(r.GuildID, r.TemplateCode)

	result := new(GuildTemplate)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodDelete, endpoint, nil, nil, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a GetInvite request to Discord and returns a Invite.
func (r *GetInvite) Send(bot *Client) (*Invite, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[135]("135", "781d4865"+r.InviteCode)
	query, err := EndpointQueryString(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      "",
			Err:           err,
		}
	}
	endpoint := EndpointGetInvite(r.InviteCode) + "?" + query

	result := new(Invite)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, ContentTypeURLQueryString, nil, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a DeleteInvite request to Discord and returns a Invite.
func (r *DeleteInvite) Send(bot *Client) (*Invite, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[136]("136", "781d4865"+r.InviteCode)
	endpoint := EndpointDeleteInvite(r.InviteCode)

	result := new(Invite)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodDelete, endpoint, nil, nil, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a CreateStageInstance request to Discord and returns a StageInstance.
func (r *CreateStageInstance) Send(bot *Client) (*StageInstance, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[137]("137")
	endpoint := EndpointCreateStageInstance()

	body, err := json.Marshal(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           fmt.Errorf(errSendMarshal, err),
		}
	}

	result := new(StageInstance)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPost, endpoint, ContentTypeJSON, body, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a GetStageInstance request to Discord and returns a StageInstance.
func (r *GetStageInstance) Send(bot *Client) (*StageInstance, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[138]("138", "e5416649"+r.ChannelID)
	endpoint := EndpointGetStageInstance(r.ChannelID)

	result := new(StageInstance)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, nil, nil, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a ModifyStageInstance request to Discord and returns a StageInstance.
func (r *ModifyStageInstance) Send(bot *Client) (*StageInstance, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[139]("139", "e5416649"+r.ChannelID)
	endpoint := EndpointModifyStageInstance(r.ChannelID)

	body, err := json.Marshal(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           fmt.Errorf(errSendMarshal, err),
		}
	}

	result := new(StageInstance)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPatch, endpoint, ContentTypeJSON, body, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a DeleteStageInstance request to Discord and returns a error.
func (r *DeleteStageInstance) Send(bot *Client) error {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[140]("140", "e5416649"+r.ChannelID)
	endpoint := EndpointDeleteStageInstance(r.ChannelID)

	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodDelete, endpoint, nil, nil, nil)
	if err != nil {
		return ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return nil
}

// Send sends a GetSticker request to Discord and returns a Sticker.
func (r *GetSticker) Send(bot *Client) (*Sticker, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[141]("141", "6eeeabf1"+r.StickerID)
	endpoint := EndpointGetSticker(r.StickerID)

	result := new(Sticker)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, nil, nil, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a ListNitroStickerPacks request to Discord and returns a ListNitroStickerPacksResponse.
func (r *ListNitroStickerPacks) Send(bot *Client) (*ListNitroStickerPacksResponse, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[142]("142")
	endpoint := EndpointListNitroStickerPacks()

	result := new(ListNitroStickerPacksResponse)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, nil, nil, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a ListGuildStickers request to Discord and returns a []*Sticker.
func (r *ListGuildStickers) Send(bot *Client) ([]*Sticker, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[143]("143", "45892a5d"+r.GuildID)
	endpoint := EndpointListGuildStickers(r.GuildID)

	result := make([]*Sticker, 0)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, nil, nil, &result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a GetGuildSticker request to Discord and returns a Sticker.
func (r *GetGuildSticker) Send(bot *Client) (*Sticker, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[144]("144", "45892a5d"+r.GuildID, "6eeeabf1"+r.StickerID)
	endpoint := EndpointGetGuildSticker(r.GuildID, r.StickerID)

	result := new(Sticker)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, nil, nil, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a CreateGuildSticker request to Discord and returns a Sticker.
func (r *CreateGuildSticker) Send(bot *Client) (*Sticker, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[145]("145", "45892a5d"+r.GuildID)
	endpoint := EndpointCreateGuildSticker(r.GuildID)

	body, err := json.Marshal(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           fmt.Errorf(errSendMarshal, err),
		}
	}

	var contentType []byte
	var multipartErr error
	if contentType, body, multipartErr = createMultipartForm(body, &r.File); multipartErr != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      "",
			Err:           err,
		}
	}

	result := new(Sticker)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPost, endpoint, contentType, body, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a ModifyGuildSticker request to Discord and returns a Sticker.
func (r *ModifyGuildSticker) Send(bot *Client) (*Sticker, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[146]("146", "45892a5d"+r.GuildID, "6eeeabf1"+r.StickerID)
	endpoint := EndpointModifyGuildSticker(r.GuildID, r.StickerID)

	body, err := json.Marshal(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           fmt.Errorf(errSendMarshal, err),
		}
	}

	result := new(Sticker)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPatch, endpoint, ContentTypeJSON, body, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a DeleteGuildSticker request to Discord and returns a error.
func (r *DeleteGuildSticker) Send(bot *Client) error {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[147]("147", "45892a5d"+r.GuildID, "6eeeabf1"+r.StickerID)
	endpoint := EndpointDeleteGuildSticker(r.GuildID, r.StickerID)

	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodDelete, endpoint, nil, nil, nil)
	if err != nil {
		return ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return nil
}

// Send sends a GetCurrentUser request to Discord and returns a User.
func (r *GetCurrentUser) Send(bot *Client) (*User, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[148]("148")
	endpoint := EndpointGetCurrentUser()

	result := new(User)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, nil, nil, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a GetUser request to Discord and returns a User.
func (r *GetUser) Send(bot *Client) (*User, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[149]("149", "209c92df"+r.UserID)
	endpoint := EndpointGetUser(r.UserID)

	result := new(User)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, nil, nil, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a ModifyCurrentUser request to Discord and returns a User.
func (r *ModifyCurrentUser) Send(bot *Client) (*User, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[150]("150")
	endpoint := EndpointModifyCurrentUser()

	body, err := json.Marshal(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           fmt.Errorf(errSendMarshal, err),
		}
	}

	result := new(User)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPatch, endpoint, ContentTypeJSON, body, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a GetCurrentUserGuilds request to Discord and returns a []*Guild.
func (r *GetCurrentUserGuilds) Send(bot *Client) ([]*Guild, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[151]("151")
	endpoint := EndpointGetCurrentUserGuilds()

	body, err := json.Marshal(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           fmt.Errorf(errSendMarshal, err),
		}
	}

	result := make([]*Guild, 0)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, ContentTypeJSON, body, &result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a GetCurrentUserGuildMember request to Discord and returns a GuildMember.
func (r *GetCurrentUserGuildMember) Send(bot *Client) (*GuildMember, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[152]("152", "45892a5d"+r.GuildID)
	endpoint := EndpointGetCurrentUserGuildMember(r.GuildID)

	result := new(GuildMember)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, nil, nil, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a LeaveGuild request to Discord and returns a error.
func (r *LeaveGuild) Send(bot *Client) error {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[153]("153", "45892a5d"+r.GuildID)
	endpoint := EndpointLeaveGuild(r.GuildID)

	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodDelete, endpoint, nil, nil, nil)
	if err != nil {
		return ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return nil
}

// Send sends a CreateDM request to Discord and returns a Channel.
func (r *CreateDM) Send(bot *Client) (*Channel, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[154]("154")
	endpoint := EndpointCreateDM()

	body, err := json.Marshal(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           fmt.Errorf(errSendMarshal, err),
		}
	}

	result := new(Channel)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPost, endpoint, ContentTypeJSON, body, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a CreateGroupDM request to Discord and returns a Channel.
func (r *CreateGroupDM) Send(bot *Client) (*Channel, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[155]("155")
	endpoint := EndpointCreateGroupDM()

	body, err := json.Marshal(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           fmt.Errorf(errSendMarshal, err),
		}
	}

	result := new(Channel)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPost, endpoint, ContentTypeJSON, body, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a GetUserConnections request to Discord and returns a []*Connection.
func (r *GetUserConnections) Send(bot *Client) ([]*Connection, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[156]("156")
	endpoint := EndpointGetUserConnections()

	result := make([]*Connection, 0)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, nil, nil, &result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a ListVoiceRegions request to Discord and returns a []*VoiceRegion.
func (r *ListVoiceRegions) Send(bot *Client) ([]*VoiceRegion, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[157]("157")
	endpoint := EndpointListVoiceRegions()

	result := make([]*VoiceRegion, 0)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, nil, nil, &result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a CreateWebhook request to Discord and returns a Webhook.
func (r *CreateWebhook) Send(bot *Client) (*Webhook, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[158]("158", "e5416649"+r.ChannelID)
	endpoint := EndpointCreateWebhook(r.ChannelID)

	body, err := json.Marshal(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           fmt.Errorf(errSendMarshal, err),
		}
	}

	result := new(Webhook)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPost, endpoint, ContentTypeJSON, body, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a GetChannelWebhooks request to Discord and returns a []*Webhook.
func (r *GetChannelWebhooks) Send(bot *Client) ([]*Webhook, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[159]("159", "e5416649"+r.ChannelID)
	endpoint := EndpointGetChannelWebhooks(r.ChannelID)

	result := make([]*Webhook, 0)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, nil, nil, &result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a GetGuildWebhooks request to Discord and returns a []*Webhook.
func (r *GetGuildWebhooks) Send(bot *Client) ([]*Webhook, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[160]("160", "45892a5d"+r.GuildID)
	endpoint := EndpointGetGuildWebhooks(r.GuildID)

	result := make([]*Webhook, 0)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, nil, nil, &result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a GetWebhook request to Discord and returns a Webhook.
func (r *GetWebhook) Send(bot *Client) (*Webhook, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[161]("161", "6d62b21b"+r.WebhookID)
	endpoint := EndpointGetWebhook(r.WebhookID)

	result := new(Webhook)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, nil, nil, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a GetWebhookwithToken request to Discord and returns a Webhook.
func (r *GetWebhookwithToken) Send(bot *Client) (*Webhook, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[162]("162", "6d62b21b"+r.WebhookID, "8954ac33"+r.WebhookToken)
	endpoint := EndpointGetWebhookwithToken(r.WebhookID, r.WebhookToken)

	result := new(Webhook)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, nil, nil, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a ModifyWebhook request to Discord and returns a Webhook.
func (r *ModifyWebhook) Send(bot *Client) (*Webhook, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[163]("163", "6d62b21b"+r.WebhookID)
	endpoint := EndpointModifyWebhook(r.WebhookID)

	body, err := json.Marshal(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           fmt.Errorf(errSendMarshal, err),
		}
	}

	result := new(Webhook)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPatch, endpoint, ContentTypeJSON, body, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a ModifyWebhookwithToken request to Discord and returns a Webhook.
func (r *ModifyWebhookwithToken) Send(bot *Client) (*Webhook, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[164]("164", "6d62b21b"+r.WebhookID, "8954ac33"+r.WebhookToken)
	endpoint := EndpointModifyWebhookwithToken(r.WebhookID, r.WebhookToken)

	body, err := json.Marshal(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           fmt.Errorf(errSendMarshal, err),
		}
	}

	result := new(Webhook)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPatch, endpoint, ContentTypeJSON, body, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a DeleteWebhook request to Discord and returns a error.
func (r *DeleteWebhook) Send(bot *Client) error {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[165]("165", "6d62b21b"+r.WebhookID)
	endpoint := EndpointDeleteWebhook(r.WebhookID)

	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodDelete, endpoint, nil, nil, nil)
	if err != nil {
		return ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return nil
}

// Send sends a DeleteWebhookwithToken request to Discord and returns a error.
func (r *DeleteWebhookwithToken) Send(bot *Client) error {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[166]("166", "6d62b21b"+r.WebhookID, "8954ac33"+r.WebhookToken)
	endpoint := EndpointDeleteWebhookwithToken(r.WebhookID, r.WebhookToken)

	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodDelete, endpoint, nil, nil, nil)
	if err != nil {
		return ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return nil
}

// Send sends a ExecuteWebhook request to Discord and returns a error.
func (r *ExecuteWebhook) Send(bot *Client) error {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[167]("167", "6d62b21b"+r.WebhookID, "8954ac33"+r.WebhookToken)
	query, err := EndpointQueryString(r)
	if err != nil {
		return ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      "",
			Err:           err,
		}
	}
	endpoint := EndpointExecuteWebhook(r.WebhookID, r.WebhookToken) + "?" + query

	body, err := json.Marshal(r)
	if err != nil {
		return ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           fmt.Errorf(errSendMarshal, err),
		}
	}

	contentType := ContentTypeJSON
	if len(r.Files) != 0 {
		var multipartErr error
		if contentType, body, multipartErr = createMultipartForm(body, r.Files...); multipartErr != nil {
			return ErrorRequest{
				ClientID:      bot.ApplicationID,
				CorrelationID: xid,
				RouteID:       routeid,
				ResourceID:    resourceid,
				Endpoint:      "",
				Err:           err,
			}
		}
	}

	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPost, endpoint, contentType, body, nil)
	if err != nil {
		return ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return nil
}

// Send sends a ExecuteSlackCompatibleWebhook request to Discord and returns a error.
func (r *ExecuteSlackCompatibleWebhook) Send(bot *Client) error {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[168]("168", "6d62b21b"+r.WebhookID, "8954ac33"+r.WebhookToken)
	query, err := EndpointQueryString(r)
	if err != nil {
		return ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      "",
			Err:           err,
		}
	}
	endpoint := EndpointExecuteSlackCompatibleWebhook(r.WebhookID, r.WebhookToken) + "?" + query

	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPost, endpoint, ContentTypeURLQueryString, nil, nil)
	if err != nil {
		return ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return nil
}

// Send sends a ExecuteGitHubCompatibleWebhook request to Discord and returns a error.
func (r *ExecuteGitHubCompatibleWebhook) Send(bot *Client) error {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[169]("169", "6d62b21b"+r.WebhookID, "8954ac33"+r.WebhookToken)
	query, err := EndpointQueryString(r)
	if err != nil {
		return ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      "",
			Err:           err,
		}
	}
	endpoint := EndpointExecuteGitHubCompatibleWebhook(r.WebhookID, r.WebhookToken) + "?" + query

	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPost, endpoint, ContentTypeURLQueryString, nil, nil)
	if err != nil {
		return ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return nil
}

// Send sends a GetWebhookMessage request to Discord and returns a Message.
func (r *GetWebhookMessage) Send(bot *Client) (*Message, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[170]("170", "6d62b21b"+r.WebhookID, "8954ac33"+r.WebhookToken, "d57d6589"+r.MessageID)
	query, err := EndpointQueryString(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      "",
			Err:           err,
		}
	}
	endpoint := EndpointGetWebhookMessage(r.WebhookID, r.WebhookToken, r.MessageID) + "?" + query

	result := new(Message)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, ContentTypeURLQueryString, nil, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a EditWebhookMessage request to Discord and returns a Message.
func (r *EditWebhookMessage) Send(bot *Client) (*Message, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[171]("171", "6d62b21b"+r.WebhookID, "8954ac33"+r.WebhookToken, "d57d6589"+r.MessageID)
	query, err := EndpointQueryString(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      "",
			Err:           err,
		}
	}
	endpoint := EndpointEditWebhookMessage(r.WebhookID, r.WebhookToken, r.MessageID) + "?" + query

	body, err := json.Marshal(r)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           fmt.Errorf(errSendMarshal, err),
		}
	}

	contentType := ContentTypeJSON
	if len(r.Files) != 0 {
		var multipartErr error
		if contentType, body, multipartErr = createMultipartForm(body, r.Files...); multipartErr != nil {
			return nil, ErrorRequest{
				ClientID:      bot.ApplicationID,
				CorrelationID: xid,
				RouteID:       routeid,
				ResourceID:    resourceid,
				Endpoint:      "",
				Err:           err,
			}
		}
	}

	result := new(Message)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodPatch, endpoint, contentType, body, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a DeleteWebhookMessage request to Discord and returns a error.
func (r *DeleteWebhookMessage) Send(bot *Client) error {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[172]("172", "6d62b21b"+r.WebhookID, "8954ac33"+r.WebhookToken, "d57d6589"+r.MessageID)
	query, err := EndpointQueryString(r)
	if err != nil {
		return ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      "",
			Err:           err,
		}
	}
	endpoint := EndpointDeleteWebhookMessage(r.WebhookID, r.WebhookToken, r.MessageID) + "?" + query

	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodDelete, endpoint, ContentTypeURLQueryString, nil, nil)
	if err != nil {
		return ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return nil
}

// Send sends a GetGateway request to Discord and returns a GetGatewayBotResponse.
func (r *GetGateway) Send(bot *Client) (*GetGatewayBotResponse, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[173]("173")
	endpoint := EndpointGetGateway()

	result := new(GetGatewayBotResponse)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, nil, nil, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a GetGatewayBot request to Discord and returns a GetGatewayBotResponse.
func (r *GetGatewayBot) Send(bot *Client) (*GetGatewayBotResponse, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[174]("174")
	endpoint := EndpointGetGatewayBot()

	result := new(GetGatewayBotResponse)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, nil, nil, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a GetCurrentBotApplicationInformation request to Discord and returns a Application.
func (r *GetCurrentBotApplicationInformation) Send(bot *Client) (*Application, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[175]("175")
	endpoint := EndpointGetCurrentBotApplicationInformation()

	result := new(Application)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, nil, nil, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

// Send sends a GetCurrentAuthorizationInformation request to Discord and returns a CurrentAuthorizationInformationResponse.
func (r *GetCurrentAuthorizationInformation) Send(bot *Client) (*CurrentAuthorizationInformationResponse, error) {
	var err error
	xid := xid.New().String()
	routeid, resourceid := RateLimitHashFuncs[176]("176")
	endpoint := EndpointGetCurrentAuthorizationInformation()

	result := new(CurrentAuthorizationInformationResponse)
	err = SendRequest(bot, xid, routeid, resourceid, fasthttp.MethodGet, endpoint, nil, nil, result)
	if err != nil {
		return nil, ErrorRequest{
			ClientID:      bot.ApplicationID,
			CorrelationID: xid,
			RouteID:       routeid,
			ResourceID:    resourceid,
			Endpoint:      endpoint,
			Err:           err,
		}
	}

	return result, nil
}

const (
	gatewayEndpointParams     = "?v=" + VersionDiscordAPI + "&encoding=json"
	invalidSessionWaitTime    = 1 * time.Second
	maxIdentifyLargeThreshold = 250
)

// Session represents a Discord Gateway WebSocket Session.
type Session struct {
	Context   context.Context
	manager   *manager
	Conn      *websocket.Conn
	heartbeat *heartbeat
	Endpoint  string
	ID        string
	Seq       int64
	sync.RWMutex
}

// isConnected returns whether the session is connected.
func (s *Session) isConnected() bool {
	if s.Context == nil {
		return false
	}

	select {
	case <-s.Context.Done():
		return false
	default:
		return true
	}
}

// canReconnect determines whether the session is in a valid state to reconnect.
func (s *Session) canReconnect() bool {
	return s.ID != "" && s.Endpoint != "" && atomic.LoadInt64(&s.Seq) != 0
}

// Connect connects a session to the Discord Gateway (WebSocket Connection).
func (s *Session) Connect(bot *Client) error {
	s.Lock()
	defer s.Unlock()

	LogSession(Logger.Info(), s.ID).Str(LogCtxClient, bot.ApplicationID).Msg("connecting session")

	return s.connect(bot)
}

// connect connects a session to a WebSocket Connection.
func (s *Session) connect(bot *Client) error {
	if s.isConnected() {
		return fmt.Errorf("session %q is already connected", s.ID)
	}

	// request a valid Gateway URL endpoint from the Discord API.
	gatewayEndpoint := s.Endpoint
	if gatewayEndpoint == "" || !s.canReconnect() {
		gateway := GetGateway{}
		response, err := gateway.Send(bot)
		if err != nil {
			return fmt.Errorf("error getting the Gateway API Endpoint: %w", err)
		}

		gatewayEndpoint = response.URL + gatewayEndpointParams

		// set the maximum allowed (Identify) concurrency rate limit.
		//
		// https://discord.com/developers/docs/topics/gateway#rate-limiting
		bot.Config.Gateway.RateLimiter.StartTx()

		identifyBucket := bot.Config.Gateway.RateLimiter.GetBucketFromID(FlagGatewaySendEventNameIdentify)
		if identifyBucket == nil {
			identifyBucket = getBucket()
			bot.Config.Gateway.RateLimiter.SetBucketFromID(FlagGatewaySendEventNameIdentify, identifyBucket)
		}

		identifyBucket.Limit = int16(response.SessionStartLimit.MaxConcurrency) + 1
		identifyBucket.Remaining = identifyBucket.Limit
		identifyBucket.Expiry = time.Now().Add(FlagGlobalRateLimitIdentifyInterval)

		bot.Config.Gateway.RateLimiter.EndTx()
	}

	var err error

	// connect to the Discord Gateway Websocket.
	s.manager = new(manager)
	s.Context, s.manager.cancel = context.WithCancel(context.Background())
	if s.Conn, _, err = websocket.Dial(s.Context, gatewayEndpoint, nil); err != nil {
		return fmt.Errorf("error connecting to the Discord Gateway: %w", err)
	}

	// handle the incoming Hello event upon connecting to the Gateway.
	hello := new(Hello)
	if err := readEvent(s, FlagGatewayEventNameHello, hello); err != nil {
		err = fmt.Errorf("error reading initial Hello event: %w", err)
		sessionErr := ErrorSession{SessionID: s.ID, Err: err}
		if disconnectErr := s.disconnect(FlagClientCloseEventCodeNormal); disconnectErr != nil {
			sessionErr.Err = ErrorDisconnect{
				Connection: ErrConnectionSession,
				Action:     err,
				Err:        disconnectErr,
			}
		}

		return sessionErr
	}

	for _, handler := range bot.Handlers.Hello {
		go handler(hello)
	}

	// begin sending heartbeat payloads every heartbeat_interval ms.
	ms := time.Millisecond * time.Duration(hello.HeartbeatInterval)
	s.heartbeat = &heartbeat{
		interval: ms,
		ticker:   time.NewTicker(ms),
		send:     make(chan Heartbeat),

		// add a HeartbeatACK to the HeartbeatACK channel to prevent
		// the length of the HeartbeatACK channel from being 0 immediately,
		// which results in an attempt to reconnect.
		acks: 1,
	}

	// create a goroutine group for the Session.
	s.manager.Group, s.manager.signal = errgroup.WithContext(s.Context)
	s.manager.err = make(chan error, 1)

	// spawn the heartbeat pulse goroutine.
	s.manager.routines.Add(1)
	atomic.AddInt32(&s.manager.pulses, 1)
	s.manager.Go(func() error {
		s.pulse()
		return nil
	})

	// spawn the heartbeat beat goroutine.
	s.manager.routines.Add(1)
	s.manager.Go(func() error {
		if err := s.beat(bot); err != nil {
			return ErrorSession{
				SessionID: s.ID,
				Err:       fmt.Errorf("heartbeat: %w", err),
			}
		}

		return nil
	})

	// send the initial Identify or Resumed packet.
	if err := s.initial(bot, 0); err != nil {
		sessionErr := ErrorSession{SessionID: s.ID, Err: err}
		if disconnectErr := s.disconnect(FlagClientCloseEventCodeNormal); disconnectErr != nil {
			sessionErr.Err = ErrorDisconnect{
				Connection: ErrConnectionSession,
				Action:     err,
				Err:        disconnectErr,
			}
		}

		return sessionErr
	}

	// spawn the event listener listen goroutine.
	s.manager.routines.Add(1)
	s.manager.Go(func() error {
		if err := s.listen(bot); err != nil {
			return ErrorSession{
				SessionID: s.ID,
				Err:       fmt.Errorf("listen: %w", err),
			}
		}

		return nil
	})

	// spawn the manager goroutine.
	s.manager.routines.Add(1)
	go s.manage(bot)

	// ensure that the Session's goroutines are spawned.
	s.manager.routines.Wait()

	return nil
}

// initial sends the initial Identify or Resume packet required to connect to the Gateway,
// then handles the incoming Ready or Resumed packet that indicates a successful connection.
func (s *Session) initial(bot *Client, attempt int) error {
	if !s.canReconnect() {
		// send an Opcode 2 Identify to the Discord Gateway.
		identify := Identify{
			Token: bot.Authentication.Token,
			Properties: IdentifyConnectionProperties{
				OS:      runtime.GOOS,
				Browser: module,
				Device:  module,
			},
			Compress:       Pointer(true),
			LargeThreshold: Pointer(maxIdentifyLargeThreshold),
			Shard:          nil, // SHARD: set shard information using s.Shard.
			Presence:       bot.Config.Gateway.GatewayPresenceUpdate,
			Intents:        bot.Config.Gateway.Intents,
		}

		if err := identify.SendEvent(bot, s); err != nil {
			return err
		}
	} else {
		// send an Opcode 6 Resume to the Discord Gateway to reconnect the session.
		resume := Resume{
			Token:     bot.Authentication.Token,
			SessionID: s.ID,
			Seq:       atomic.LoadInt64(&s.Seq),
		}

		if err := resume.SendEvent(bot, s); err != nil {
			return err
		}
	}

	// handle the incoming Ready, Resumed or Replayed event (or Opcode 9 Invalid Session).
	payload := new(GatewayPayload)
	if err := socket.Read(s.Context, s.Conn, payload); err != nil {
		return fmt.Errorf("error reading initial payload: %w", err)
	}

	LogPayload(LogSession(Logger.Info(), s.ID), payload.Op, payload.Data).Msg("received initial payload")

	switch payload.Op {
	case FlagGatewayOpcodeDispatch:
		switch {
		// When a connection is successful, the Discord Gateway will respond with a Ready event.
		case *payload.EventName == FlagGatewayEventNameReady:
			ready := new(Ready)
			if err := json.Unmarshal(payload.Data, ready); err != nil {
				return fmt.Errorf("error reading ready event: %w", err)
			}

			s.ID = ready.SessionID
			s.Seq = 0
			s.Endpoint = ready.ResumeGatewayURL
			// SHARD: set shard information using r.Shard
			bot.ApplicationID = ready.Application.ID

			LogSession(Logger.Info(), s.ID).Msg("received Ready event")

			for _, handler := range bot.Handlers.Ready {
				go handler(ready)
			}

		// When a reconnection is successful, the Discord Gateway will respond
		// by replaying all missed events in order, finalized by a Resumed event.
		case *payload.EventName == FlagGatewayEventNameResumed:
			LogSession(Logger.Info(), s.ID).Msg("received Resumed event")

			for _, handler := range bot.Handlers.Resumed {
				go handler(&Resumed{})
			}

		// When a reconnection is successful, the Discord Gateway will respond
		// by replaying all missed events in order, finalized by a Resumed event.
		default:
			// handle the initial payload(s) until a Resumed event is encountered.
			go bot.handle(*payload.EventName, payload.Data)

			for {
				replayed := new(GatewayPayload)
				if err := socket.Read(s.Context, s.Conn, replayed); err != nil {
					return fmt.Errorf("error replaying events: %w", err)
				}

				if replayed.Op == FlagGatewayOpcodeDispatch && *replayed.EventName == FlagGatewayEventNameResumed {
					LogSession(Logger.Info(), s.ID).Msg("received Resumed event")

					for _, handler := range bot.Handlers.Resumed {
						go handler(&Resumed{})
					}

					return nil
				}

				go bot.handle(*payload.EventName, payload.Data)
			}
		}

	// When the maximum concurrency limit has been reached while connecting, or when
	// the session does NOT reconnect in time, the Discord Gateway send an Opcode 9 Invalid Session.
	case FlagGatewayOpcodeInvalidSession:
		if attempt < 1 {
			// wait for Discord to close the session, then complete a fresh connect.
			<-time.NewTimer(invalidSessionWaitTime).C

			s.ID = ""
			s.Seq = 0
			if err := s.initial(bot, attempt+1); err != nil {
				return err
			}

			return nil
		}

		return fmt.Errorf("session %q couldn't connect to the Discord Gateway or has invalidated an active session", s.ID)
	default:
		return fmt.Errorf("session %q received payload %d during connection which is unexpected", s.ID, payload.Op)
	}

	return nil
}

// Disconnect disconnects a session from the Discord Gateway using the given status code.
func (s *Session) Disconnect() error {
	s.Lock()

	if !s.isConnected() {
		s.Unlock()

		return fmt.Errorf("session %q is already disconnected", s.ID)
	}

	id := s.ID
	LogSession(Logger.Info(), id).Msgf("disconnecting session with code %d", FlagClientCloseEventCodeNormal)

	s.manager.signal = context.WithValue(s.manager.signal, keySignal, signalDisconnect)

	if err := s.disconnect(FlagClientCloseEventCodeNormal); err != nil {
		s.Unlock()

		return ErrorDisconnect{
			Connection: ErrConnectionSession,
			Action:     nil,
			Err:        err,
		}
	}

	s.Unlock()

	if err := <-s.manager.err; err != nil {
		return err
	}

	putSession(s)

	LogSession(Logger.Info(), id).Msgf("disconnected session with code %d", FlagClientCloseEventCodeNormal)

	return nil
}

// disconnect disconnects a session from a WebSocket Connection using the given status code.
func (s *Session) disconnect(code int) error {
	// cancel the context to kill the goroutines of the Session.
	defer s.manager.cancel()

	if err := s.Conn.Close(websocket.StatusCode(code), ""); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

// Reconnect reconnects an already connected session to the Discord Gateway
// by disconnecting the session, then connecting again.
func (s *Session) Reconnect(bot *Client) error {
	s.reconnect("reconnecting")

	if err := <-s.manager.err; err != nil {
		return err
	}

	// connect to the Discord Gateway again.
	if err := s.Connect(bot); err != nil {
		return fmt.Errorf("error reconnecting session %q: %w", s.ID, err)
	}

	return nil
}

// readEvent is a helper function for reading events from the WebSocket Session.
func readEvent(s *Session, name string, dst any) error {
	payload := new(GatewayPayload)
	if err := socket.Read(s.Context, s.Conn, payload); err != nil {
		return fmt.Errorf("readEvent: %w", err)
	}

	if err := json.Unmarshal(payload.Data, dst); err != nil {
		return fmt.Errorf("readEvent: %w", err)
	}

	return nil
}

// writeEvent is a helper function for writing events to the WebSocket Session.
func writeEvent(bot *Client, s *Session, op int, name string, dst any) error {
RATELIMIT:
	// a single command is PROCESSED at any point in time.
	bot.Config.Gateway.RateLimiter.Lock()

	LogCommand(LogSession(Logger.Trace(), s.ID), bot.ApplicationID, op, name).Msg("processing gateway command")

	for {
		bot.Config.Gateway.RateLimiter.StartTx()

		globalBucket := bot.Config.Gateway.RateLimiter.GetBucket(GlobalRateLimitRouteID, "")

		// stop waiting when the Global Rate Limit Bucket is NOT empty.
		if isNotEmpty(globalBucket) {
			switch op {
			case FlagGatewayOpcodeIdentify:
				identifyBucket := bot.Config.Gateway.RateLimiter.GetBucketFromID(FlagGatewaySendEventNameIdentify)

				if isNotEmpty(identifyBucket) {
					if globalBucket != nil {
						globalBucket.Remaining--
					}

					if identifyBucket != nil {
						identifyBucket.Remaining--
					}

					bot.Config.Gateway.RateLimiter.EndTx()

					goto SEND
				}

				if isExpired(identifyBucket) {
					if globalBucket != nil {
						globalBucket.Remaining--
					}

					if identifyBucket != nil {
						identifyBucket.Reset(time.Now().Add(FlagGlobalRateLimitIdentifyInterval))
						identifyBucket.Remaining--
					}

					bot.Config.Gateway.RateLimiter.EndTx()

					goto SEND
				}

				var wait time.Time
				if identifyBucket != nil {
					wait = identifyBucket.Expiry
				}

				// do NOT block other requests due to a Command Rate Limit.
				bot.Config.Gateway.RateLimiter.EndTx()
				bot.Config.Gateway.RateLimiter.Unlock()

				// reduce CPU usage by blocking the current goroutine
				// until it's eligible for action.
				if identifyBucket != nil {
					<-time.After(time.Until(wait))
				}

				goto RATELIMIT

			default:
				if globalBucket != nil {
					globalBucket.Remaining--
				}

				bot.Config.Gateway.RateLimiter.EndTx()

				goto SEND
			}
		}

		// reset the Global Rate Limit Bucket when the current Bucket has passed its expiry.
		if isExpired(globalBucket) {
			globalBucket.Reset(time.Now().Add(time.Minute))
		}

		bot.Config.Gateway.RateLimiter.EndTx()
	}

SEND:
	bot.Config.Gateway.RateLimiter.Unlock()

	LogCommand(LogSession(Logger.Trace(), s.ID), bot.ApplicationID, op, name).Msg("sending gateway command")

	// write the event to the WebSocket Connection.
	event, err := json.Marshal(dst)
	if err != nil {
		return fmt.Errorf("writeEvent: %w", err)
	}

	if err = socket.Write(s.Context, s.Conn, websocket.MessageBinary,
		GatewayPayload{ //nolint:exhaustruct
			Op:   op,
			Data: event,
		}); err != nil {
		return fmt.Errorf("writeEvent: %w", err)
	}

	LogCommand(LogSession(Logger.Trace(), s.ID), bot.ApplicationID, op, name).Msg("sent gateway command")

	return nil
}

// SendEvent sends an Opcode 1 Heartbeat event to the Discord Gateway.
func (c *Heartbeat) SendEvent(bot *Client, session *Session) error {
	if err := writeEvent(bot, session, FlagGatewayOpcodeHeartbeat, FlagGatewaySendEventNameHeartbeat, c); err != nil {
		return err
	}

	return nil
}

// SendEvent sends an Opcode 2 Identify event to the Discord Gateway.
func (c *Identify) SendEvent(bot *Client, session *Session) error {
	if err := writeEvent(bot, session, FlagGatewayOpcodeIdentify, FlagGatewaySendEventNameIdentify, c); err != nil {
		return err
	}

	return nil
}

// SendEvent sends an Opcode 3 UpdatePresence event to the Discord Gateway.
func (c *GatewayPresenceUpdate) SendEvent(bot *Client, session *Session) error {
	if err := writeEvent(bot, session, FlagGatewayOpcodePresenceUpdate, FlagGatewaySendEventNameUpdatePresence, c); err != nil {
		return err
	}

	return nil
}

// SendEvent sends an Opcode 4 UpdateVoiceState event to the Discord Gateway.
func (c *VoiceStateUpdate) SendEvent(bot *Client, session *Session) error {
	if err := writeEvent(bot, session, FlagGatewayOpcodeVoiceStateUpdate, FlagGatewaySendEventNameUpdateVoiceState, c); err != nil {
		return err
	}

	return nil
}

// SendEvent sends an Opcode 6 Resume event to the Discord Gateway.
func (c *Resume) SendEvent(bot *Client, session *Session) error {
	if err := writeEvent(bot, session, FlagGatewayOpcodeResume, FlagGatewaySendEventNameResume, c); err != nil {
		return err
	}

	return nil
}

// SendEvent sends an Opcode 8 RequestGuildMembers event to the Discord Gateway.
func (c *RequestGuildMembers) SendEvent(bot *Client, session *Session) error {
	if err := writeEvent(bot, session, FlagGatewayOpcodeRequestGuildMembers, FlagGatewaySendEventNameRequestGuildMembers, c); err != nil {
		return err
	}

	return nil
}

// heartbeat represents the heartbeat mechanism for a Session.
type heartbeat struct {
	ticker   *time.Ticker
	send     chan Heartbeat
	interval time.Duration
	acks     uint32
}

// Monitor returns the current amount of HeartbeatACKs for a Session's heartbeat.
func (s *Session) Monitor() uint32 {
	s.Lock()
	acks := atomic.LoadUint32(&s.heartbeat.acks)
	s.Unlock()

	return acks
}

// beat listens for pulses to send Opcode 1 Heartbeats to the Discord Gateway (to verify the connection is alive).
func (s *Session) beat(bot *Client) error {
	s.manager.routines.Done()

	// ensure that all pulse routines are closed prior to closing.
	defer func() {
		for {
			select {
			case <-s.heartbeat.send:
			case <-s.Context.Done():
				if atomic.LoadInt32(&s.manager.pulses) != 0 {
					break
				}

				s.logClose("heartbeat")

				return
			}
		}
	}()

	for {
		select {
		case hb := <-s.heartbeat.send:
			s.Lock()

			// close the connection if the last sent Heartbeat never received a HeartbeatACK.
			if atomic.LoadUint32(&s.heartbeat.acks) == 0 {
				s.Unlock()

				s.reconnect("attempting to reconnect session due to no HeartbeatACK")

				return nil
			}

			// prevent two Heartbeat Payloads being sent to the Discord Gateway consecutively within nanoseconds,
			// when the ticker queues a Heartbeat while the listen thread (onPayload) queues a Heartbeat
			// (in response to the Discord Gateway).
			//
			// clear queued (outdated) heartbeats.
			for len(s.heartbeat.send) > 0 {
				// ensure the latest sequence is sent.
				if h := <-s.heartbeat.send; h.Data > hb.Data {
					hb.Data = h.Data
				}
			}

			// send a Heartbeat to the Discord Gateway (WebSocket Connection).
			if err := hb.SendEvent(bot, s); err != nil {
				s.Unlock()

				return err
			}

			// reset the ticker (and empty existing ticks).
			s.heartbeat.ticker.Reset(s.heartbeat.interval)
			for len(s.heartbeat.ticker.C) > 0 {
				<-s.heartbeat.ticker.C
			}

			// reset the amount of HeartbeatACKs since the last heartbeat.
			atomic.StoreUint32(&s.heartbeat.acks, 0)

			LogSession(Logger.Info(), s.ID).Msg("sent heartbeat")

			s.Unlock()

		case <-s.Context.Done():
			return nil
		}
	}
}

// pulse generates Opcode 1 Heartbeats for a Session's heartbeat channel.
func (s *Session) pulse() {
	s.manager.routines.Done()
	defer s.decrementPulses()

	// send an Opcode 1 Heartbeat payload after heartbeat_interval * jitter milliseconds
	// (where jitter is a random value between 0 and 1).
	s.Lock()
	s.heartbeat.send <- Heartbeat{Data: s.Seq}
	LogSession(Logger.Info(), s.ID).Msg("queued jitter heartbeat")
	s.Unlock()

	for {
		select {
		// every Heartbeat Interval...
		case <-s.heartbeat.ticker.C:
			s.Lock()

			// queue a heartbeat.
			s.heartbeat.send <- Heartbeat{Data: s.Seq}

			LogSession(Logger.Info(), s.ID).Msg("queued heartbeat")

			s.Unlock()

		case <-s.Context.Done():
			s.Lock()
			s.logClose("pulse")
			s.Unlock()

			return
		}
	}
}

// respond responds to Opcode 1 Heartbeats from the Discord Gateway.
func (s *Session) respond(data json.RawMessage) error {
	defer s.decrementPulses()

	var heartbeat Heartbeat
	if err := json.Unmarshal(data, &heartbeat); err != nil {
		return fmt.Errorf("error unmarshalling incoming Heartbeat: %w", err)
	}

	atomic.StoreInt64(&s.Seq, heartbeat.Data)

	s.Lock()

	// ensure that the heartbeat routine has not been closed.
	if atomic.LoadInt32(&s.manager.pulses) <= 1 {
		s.Unlock()

		return nil
	}

	// heartbeat() checks for the amount of HeartbeatACKs received since the last Heartbeat.
	// There is a possibility for this value to be 0 due to latency rather than a dead connection.
	// For example, when a Heartbeat is queued, sent, responded, and sent.
	//
	// Prevent this possibility by treating this response from Discord as an indication that the
	// connection is still alive.
	atomic.AddUint32(&s.heartbeat.acks, 1)

	// send an Opcode 1 Heartbeat without waiting the remainder of the current interval.
	s.heartbeat.send <- heartbeat

	LogSession(Logger.Info(), s.ID).Msg("responded to heartbeat")

	s.Unlock()

	return nil
}

// listen listens to the connection for payloads from the Discord Gateway.
func (s *Session) listen(bot *Client) error {
	s.manager.routines.Done()

	var err error

	for {
		payload := getPayload()
		if err = socket.Read(s.Context, s.Conn, payload); err != nil {
			break
		}

		LogPayload(LogSession(Logger.Info(), s.ID), payload.Op, payload.Data).Msg("received payload")

		if err = s.onPayload(bot, *payload); err != nil {
			break
		}
	}

	s.Lock()
	defer s.Unlock()
	defer s.logClose("listen")

	select {
	case <-s.Context.Done():
		return nil

	default:
		return err
	}
}

// onPayload handles an Discord Gateway Payload.
func (s *Session) onPayload(bot *Client, payload GatewayPayload) error {
	defer putPayload(&payload)

	// https://discord.com/developers/docs/topics/opcodes-and-status-codes#gateway-gateway-opcodes
	switch payload.Op {
	// run the bot's event handlers.
	case FlagGatewayOpcodeDispatch:
		atomic.StoreInt64(&s.Seq, *payload.SequenceNumber)
		go bot.handle(*payload.EventName, payload.Data)

	// send an Opcode 1 Heartbeat to the Discord Gateway.
	case FlagGatewayOpcodeHeartbeat:
		s.Lock()
		atomic.AddInt32(&s.manager.pulses, 1)
		s.Unlock()

		s.manager.Go(func() error {
			if err := s.respond(payload.Data); err != nil {
				return fmt.Errorf("respond: %w", err)
			}

			return nil
		})

	// handle the successful acknowledgement of the client's last heartbeat.
	case FlagGatewayOpcodeHeartbeatACK:
		s.Lock()
		atomic.AddUint32(&s.heartbeat.acks, 1)
		s.Unlock()

	// occurs when the Discord Gateway is shutting down the connection, while signalling the client to reconnect.
	case FlagGatewayOpcodeReconnect:
		s.reconnect("reconnecting session due to Opcode 7 Reconnect")

		return nil

	// in the context of onPayload, an Invalid Session occurs when an active session is invalidated.
	case FlagGatewayOpcodeInvalidSession:
		// wait for Discord to close the session, then complete a fresh connect.
		<-time.NewTimer(invalidSessionWaitTime).C

		s.Lock()
		defer s.Unlock()

		if err := s.initial(bot, 0); err != nil {
			return err
		}
	}

	return nil
}

// signal represents a manager Context Signal.
type signal string

// manager Context Signals.
const (
	// keySignal represents the Context key for a manager's signals.
	keySignal = signal("signal")

	// keyReason represents the Context key for a manager's reason for disconnection.
	keyReason = signal("reason")

	// signalDisconnect indicates that a Disconnection was called purposefully.
	signalDisconnect = 1

	// signalReconnect signals the manager to reconnect upon a successful disconnection.
	signalReconnect = 2
)

// manager represents a manager of a Session's goroutines.
type manager struct {
	signal context.Context
	cancel context.CancelFunc
	err    chan error
	*errgroup.Group
	routines sync.WaitGroup
	pulses   int32
}

// decrementPulses safely decrements the pulses counter of a Session manager.
func (s *Session) decrementPulses() {
	s.Lock()
	defer s.Unlock()

	atomic.AddInt32(&s.manager.pulses, -1)
}

// logClose safely logs the close of a Session's goroutine.
func (s *Session) logClose(routine string) {
	LogSession(Logger.Info(), s.ID).Msgf("closed %s routine", routine)
}

// reconnect spawns a goroutine for reconnection which prompts the manager
// to reconnect upon a disconnection.
func (s *Session) reconnect(reason string) {
	s.manager.Go(func() error {
		s.Lock()
		defer s.logClose("reconnect")
		defer s.Unlock()

		LogSession(Logger.Info(), s.ID).Msg(reason)

		s.manager.signal = context.WithValue(s.manager.signal, keySignal, signalReconnect)
		if err := s.disconnect(FlagClientCloseEventCodeReconnect); err != nil {
			return fmt.Errorf("reconnect: %w", err)
		}

		return nil
	})
}

// manage manages a Session's goroutines.
func (s *Session) manage(bot *Client) {
	s.manager.routines.Done()
	defer func() {
		s.Lock()
		s.logClose("manager")
		s.Unlock()
	}()

	// wait until all of a Session's goroutines are closed.
	err := s.manager.Wait()

	// log the reason for disconnection (if applicable).
	if reason := s.manager.signal.Value(keyReason); reason != nil {
		LogSession(Logger.Info(), s.ID).Msgf("%v", reason)
	}

	// when a signal is provided, it indicates that the disconnection was purposeful.
	signal := s.manager.signal.Value(keySignal)
	switch signal {
	case signalDisconnect:
		LogSession(Logger.Info(), s.ID).Msg("successfully disconnected")

		s.manager.err <- nil

		return

	case signalReconnect:
		LogSession(Logger.Info(), s.ID).Msg("successfully disconnected (while reconnecting)")

		// allow Discord to close the session.
		<-time.After(time.Second)

		s.manager.err <- nil

		return
	}

	// when an error caused goroutines to close, manage the state of disconnection.
	if err != nil {
		disconnectErr := new(ErrorDisconnect)
		closeErr := new(websocket.CloseError)
		switch {
		// when an error occurs from a purposeful disconnection.
		case errors.As(err, disconnectErr):
			s.manager.err <- err

		// when an error occurs from a WebSocket Close Error.
		case errors.As(err, closeErr):
			s.manager.err <- s.handleGatewayCloseError(bot, closeErr)

		default:
			if cErr := s.Conn.Close(websocket.StatusCode(FlagClientCloseEventCodeAway), ""); cErr != nil {
				s.manager.err <- ErrorDisconnect{
					Connection: ErrConnectionSession,
					Err:        cErr,
					Action:     err,
				}

				return
			}

			s.manager.err <- err
		}

		return
	}

	s.manager.err <- nil
}

// handleGatewayCloseError handles a WebSocket CloseError.
func (s *Session) handleGatewayCloseError(bot *Client, closeErr *websocket.CloseError) error {
	code, ok := GatewayCloseEventCodes[int(closeErr.Code)]
	switch ok {
	// Gateway Close Event Code is known.
	case true:
		LogSession(Logger.Info(), s.ID).
			Msgf("received Gateway Close Event Code %d %s: %s",
				code.Code, code.Description, code.Explanation,
			)

		if code.Reconnect {
			s.reconnect(fmt.Sprintf("reconnecting due to Gateway Close Event Code %d", code.Code))

			return nil
		}

		return closeErr

	// Gateway Close Event Code is unknown.
	default:

		// when another goroutine calls disconnect(),
		// s.Conn.Close is called before s.cancel which will result in
		// a CloseError with the close code that Disgo uses to reconnect.
		if closeErr.Code == websocket.StatusCode(FlagClientCloseEventCodeReconnect) {
			return nil
		}

		LogSession(Logger.Info(), s.ID).
			Msgf("received unknown Gateway Close Event Code %d with reason %q",
				closeErr.Code, closeErr.Reason,
			)

		return closeErr
	}
}
