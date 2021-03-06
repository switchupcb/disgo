package wrapper

import (
	"encoding/json"
	"time"
)

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
	Code        int
	Description string
	Explanation string
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
	Code        int
	Description string
	Explanation string
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
		40060:  "Interaction has already been acknowledged",
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
		50041:  "Invalid API version provided",
		50045:  "File uploaded exceeds the maximum size",
		50046:  "Invalid file uploaded",
		50054:  "Cannot self-redeem this gift",
		50055:  "Invalid Guild",
		50068:  "Invalid message type",
		50070:  "Payment source required to redeem gift",
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
		220003: "Webhooks can only create threads in forum channels",
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

// Snowflake represents a Discord API Snowflake.
type Snowflake uint64

// Flag represents an (unused) alias for a Discord API Flag ranging from 0 - 255.
type Flag uint8

// BitFlag represents an alias for a Discord API Bitwise Flag denoted by 1 << x.
type BitFlag uint

// CodeFlag represents an alias for a Discord API Code ranging from 0 - 65535.
type CodeFlag uint16

// Gateway Events
// https://discord.com/developers/docs/topics/gateway#commands-and-events-gateway-events
type Event interface{}

// Gateway Event Names
// https://discord.com/developers/docs/topics/gateway#event-names
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
// https://discord.com/developers/docs/topics/gateway#hello-hello-structure
type Hello struct {
	HeartbeatInterval int `json:"heartbeat_interval"`
}

// Ready Event Fields
// https://discord.com/developers/docs/topics/gateway#ready-ready-event-fields
type Ready struct {
	Version     int          `json:"v"`
	User        *User        `json:"user"`
	Guilds      []*Guild     `json:"guilds"`
	SessionID   string       `json:"session_id"`
	Shard       *[2]int      `json:"shard,omitempty"`
	Application *Application `json:"application"`
}

// Resumed
// https://discord.com/developers/docs/topics/gateway#resumed
type Resumed struct{}

// Reconnect
// https://discord.com/developers/docs/topics/gateway#reconnect
type Reconnect struct{}

// Invalid Session
// https://discord.com/developers/docs/topics/gateway#invalid-session
type InvalidSession struct {
	Data bool `json:"d"`
}

// Application Command Permissions Update
// https://discord.com/developers/docs/topics/gateway#application-command-permissions-update
type ApplicationCommandPermissionsUpdate struct {
	*GuildApplicationCommandPermissions
}

// Auto Moderation Rule Create
// https://discord.com/developers/docs/topics/gateway#auto-moderation-rule-create
type AutoModerationRuleCreate struct {
	*AutoModerationRule
}

// Auto Moderation Rule Update
// https://discord.com/developers/docs/topics/gateway#auto-moderation-rule-update
type AutoModerationRuleUpdate struct {
	*AutoModerationRule
}

// Auto Moderation Rule Delete
// https://discord.com/developers/docs/topics/gateway#auto-moderation-rule-delete
type AutoModerationRuleDelete struct {
	*AutoModerationRule
}

// Auto Moderation Action Execution
// https://discord.com/developers/docs/topics/gateway#auto-moderation-action-execution
type AutoModerationActionExecution struct {
	GuildID              string               `json:"guild_id"`
	Action               AutoModerationAction `json:"action"`
	RuleID               string               `json:"rule_id"`
	RuleTriggerType      Flag                 `json:"rule_trigger_type"`
	UserID               string               `json:"user_id"`
	ChannelID            string               `json:"channel_id"`
	MessageID            string               `json:"message_id"`
	AlertSystemMessageID string               `json:"alert_system_message_id"`
	Content              string               `json:"content"`
	MatchedKeyword       *string              `json:"matched_keyword"`
	MatchedContent       *string              `json:"matched_content"`
}

// Channel Create
// https://discord.com/developers/docs/topics/gateway#channel-create
type ChannelCreate struct {
	*Channel
}

// Channel Update
// https://discord.com/developers/docs/topics/gateway#channel-update
type ChannelUpdate struct {
	*Channel
}

// Channel Delete
// https://discord.com/developers/docs/topics/gateway#channel-delete
type ChannelDelete struct {
	*Channel
}

// Thread Create
// https://discord.com/developers/docs/topics/gateway#thread-create
type ThreadCreate struct {
	*Channel
	NewlyCreated bool `json:"newly_created,omitempty"`
}

// Thread Update
// https://discord.com/developers/docs/topics/gateway#thread-update
type ThreadUpdate struct {
	*Channel
}

// Thread Delete
// https://discord.com/developers/docs/topics/gateway#thread-delete
type ThreadDelete struct {
	*Channel
}

// Thread List Sync Event Fields
// https://discord.com/developers/docs/topics/gateway#thread-list-sync
type ThreadListSync struct {
	GuildID    string          `json:"guild_id"`
	ChannelIDs []string        `json:"channel_ids,omitempty"`
	Threads    []*Channel      `json:"threads"`
	Members    []*ThreadMember `json:"members"`
}

// Thread Member Update
// https://discord.com/developers/docs/topics/gateway#thread-member-update
type ThreadMemberUpdate struct {
	*ThreadMember
	GuildID string `json:"guild_id"`
}

// Thread Members Update
// https://discord.com/developers/docs/topics/gateway#thread-members-update
type ThreadMembersUpdate struct {
	ID             string          `json:"id"`
	GuildID        string          `json:"guild_id"`
	MemberCount    int             `json:"member_count"`
	AddedMembers   []*ThreadMember `json:"added_members,omitempty"`
	RemovedMembers []string        `json:"removed_member_ids,omitempty"`
}

// Channel Pins Update
// https://discord.com/developers/docs/topics/gateway#channel-pins-update
type ChannelPinsUpdate struct {
	GuildID          string    `json:"guild_id,omitempty"`
	ChannelID        string    `json:"channel_id"`
	LastPinTimestamp time.Time `json:"last_pin_timestamp,omitempty"`
}

// Guild Create
// https://discord.com/developers/docs/topics/gateway#guild-create
type GuildCreate struct {
	*Guild

	// https://discord.com/developers/docs/topics/threads#gateway-events
	Threads []*Channel `json:"threads,omitempty"`
}

// Guild Update
// https://discord.com/developers/docs/topics/gateway#guild-update
type GuildUpdate struct {
	*Guild
}

// Guild Delete
// https://discord.com/developers/docs/topics/gateway#guild-delete
type GuildDelete struct {
	*Guild
}

// Guild Ban Add
// https://discord.com/developers/docs/topics/gateway#guild-ban-add
type GuildBanAdd struct {
	GuildID string `json:"guild_id"`
	User    *User  `json:"user"`
}

// Guild Ban Remove
// https://discord.com/developers/docs/topics/gateway#guild-ban-remove
type GuildBanRemove struct {
	GuildID string `json:"guild_id"`
	User    *User  `json:"user"`
}

// Guild Emojis Update
// https://discord.com/developers/docs/topics/gateway#guild-emojis-update
type GuildEmojisUpdate struct {
	GuildID string   `json:"guild_id"`
	Emojis  []*Emoji `json:"emojis"`
}

// Guild Stickers Update
// https://discord.com/developers/docs/topics/gateway#guild-stickers-update
type GuildStickersUpdate struct {
	GuildID  string     `json:"guild_id"`
	Stickers []*Sticker `json:"stickers"`
}

// Guild Integrations Update
// https://discord.com/developers/docs/topics/gateway#guild-integrations-update
type GuildIntegrationsUpdate struct {
	GuildID string `json:"guild_id"`
}

// Guild Member Add
// https://discord.com/developers/docs/topics/gateway#guild-member-add
type GuildMemberAdd struct {
	GuildID string `json:"guild_id"`
	*GuildMember
}

// Guild Member Remove
// https://discord.com/developers/docs/topics/gateway#guild-member-remove
type GuildMemberRemove struct {
	GuildID string `json:"guild_id"`
	User    *User  `json:"user"`
}

// Guild Member Update
// https://discord.com/developers/docs/topics/gateway#guild-member-update
type GuildMemberUpdate struct {
	*GuildMember
}

// Guild Members Chunk
// https://discord.com/developers/docs/topics/gateway#guild-members-chunk
type GuildMembersChunk struct {
	GuildID    string            `json:"guild_id"`
	Members    []*GuildMember    `json:"members"`
	ChunkIndex int               `json:"chunk_index"`
	ChunkCount int               `json:"chunk_count"`
	Presences  []*PresenceUpdate `json:"presences,omitempty"`
	NotFound   []string          `json:"not_found,omitempty"`
	Nonce      *string           `json:"nonce,omitempty"`
}

// Guild Role Create
// https://discord.com/developers/docs/topics/gateway#guild-role-create
type GuildRoleCreate struct {
	GuildID string `json:"guild_id"`
	Role    *Role  `json:"role"`
}

// Guild Role Update
// https://discord.com/developers/docs/topics/gateway#guild-role-update
type GuildRoleUpdate struct {
	GuildID string `json:"guild_id"`
	Role    *Role  `json:"role"`
}

// Guild Role Delete
// https://discord.com/developers/docs/topics/gateway#guild-role-delete
type GuildRoleDelete struct {
	GuildID string `json:"guild_id"`
	RoleID  string `json:"role_id"`
}

// Guild Scheduled Event Create
// https://discord.com/developers/docs/topics/gateway#guild-scheduled-event-create
type GuildScheduledEventCreate struct {
	*GuildScheduledEvent
}

// Guild Scheduled Event Update
// https://discord.com/developers/docs/topics/gateway#guild-scheduled-event-update
type GuildScheduledEventUpdate struct {
	*GuildScheduledEvent
}

// Guild Scheduled Event Delete
// https://discord.com/developers/docs/topics/gateway#guild-scheduled-event-delete
type GuildScheduledEventDelete struct {
	*GuildScheduledEvent
}

// Guild Scheduled Event User Add
// https://discord.com/developers/docs/topics/gateway#guild-scheduled-event-user-add
type GuildScheduledEventUserAdd struct {
	GuildScheduledEventID string `json:"guild_scheduled_event_id"`
	UserID                string `json:"user_id"`
	GuildID               string `json:"guild_id"`
}

// Guild Scheduled Event User Remove
// https://discord.com/developers/docs/topics/gateway#guild-scheduled-event-user-remove
type GuildScheduledEventUserRemove struct {
	GuildScheduledEventID string `json:"guild_scheduled_event_id"`
	UserID                string `json:"user_id"`
	GuildID               string `json:"guild_id"`
}

// Integration Create
// https://discord.com/developers/docs/topics/gateway#integration-create
type IntegrationCreate struct {
	*Integration
	GuildID string `json:"guild_id"`
}

// Integration Update
// https://discord.com/developers/docs/topics/gateway#integration-update
type IntegrationUpdate struct {
	*Integration
	GuildID string `json:"guild_id"`
}

// Integration Delete
// https://discord.com/developers/docs/topics/gateway#integration-delete
type IntegrationDelete struct {
	IntegrationID string `json:"id"`
	GuildID       string `json:"guild_id"`
	ApplicationID string `json:"application_id,omitempty"`
}

// Interaction Create
// https://discord.com/developers/docs/topics/gateway#interaction-create
type InteractionCreate struct {
	*Interaction
}

// Invite Create
// https://discord.com/developers/docs/topics/gateway#invite-create
type InviteCreate struct {
	ChannelID         string       `json:"channel_id"`
	Code              string       `json:"code"`
	CreatedAt         time.Time    `json:"created_at"`
	GuildID           string       `json:"guild_id,omitempty"`
	Inviter           *User        `json:"inviter,omitempty"`
	MaxAge            int          `json:"max_age"`
	MaxUses           int          `json:"max_uses"`
	TargetType        *int         `json:"target_user_type,omitempty"`
	TargetUser        *User        `json:"target_user,omitempty"`
	TargetApplication *Application `json:"target_application,omitempty"`
	Temporary         bool         `json:"temporary"`
	Uses              int          `json:"uses"`
}

// Invite Delete
// https://discord.com/developers/docs/topics/gateway#invite-delete
type InviteDelete struct {
	ChannelID string `json:"channel_id"`
	GuildID   string `json:"guild_id,omitempty"`
	Code      string `json:"code"`
}

// Message Create
// https://discord.com/developers/docs/topics/gateway#message-create
type MessageCreate struct {
	*Message
}

// Message Update
// https://discord.com/developers/docs/topics/gateway#message-update
type MessageUpdate struct {
	Message *Message
}

// Message Delete
// https://discord.com/developers/docs/topics/gateway#message-delete
type MessageDelete struct {
	MessageID string `json:"id"`
	ChannelID string `json:"channel_id"`
	GuildID   string `json:"guild_id,omitempty"`
}

// Message Delete Bulk
// https://discord.com/developers/docs/topics/gateway#message-delete-bulk
type MessageDeleteBulk struct {
	MessageIDs []string `json:"ids"`
	ChannelID  string   `json:"channel_id"`
	GuildID    string   `json:"guild_id,omitempty"`
}

// Message Reaction Add
// https://discord.com/developers/docs/topics/gateway#message-reaction-add
type MessageReactionAdd struct {
	UserID    string       `json:"user_id"`
	ChannelID string       `json:"channel_id"`
	MessageID string       `json:"message_id"`
	GuildID   string       `json:"guild_id,omitempty"`
	Member    *GuildMember `json:"member,omitempty"`
	Emoji     *Emoji       `json:"emoji"`
}

// Message Reaction Remove
// https://discord.com/developers/docs/topics/gateway#message-reaction-remove
type MessageReactionRemove struct {
	UserID    string `json:"user_id"`
	ChannelID string `json:"channel_id"`
	MessageID string `json:"message_id"`
	GuildID   string `json:"guild_id,omitempty"`
	Emoji     *Emoji `json:"emoji"`
}

// Message Reaction Remove All
// https://discord.com/developers/docs/topics/gateway#message-reaction-remove-all
type MessageReactionRemoveAll struct {
	ChannelID string `json:"channel_id"`
	MessageID string `json:"message_id"`
	GuildID   string `json:"guild_id,omitempty"`
}

// Message Reaction Remove Emoji
// https://discord.com/developers/docs/topics/gateway#message-reaction-remove-emoji
type MessageReactionRemoveEmoji struct {
	ChannelID string `json:"channel_id"`
	GuildID   string `json:"guild_id,omitempty"`
	MessageID string `json:"message_id"`
	Emoji     *Emoji `json:"emoji"`
}

// Presence Update Event Fields
// https://discord.com/developers/docs/topics/gateway#presence-update-presence-update-event-fields
type PresenceUpdate struct {
	User         *User         `json:"user"`
	GuildID      string        `json:"guild_id"`
	Status       string        `json:"status"`
	Activities   []*Activity   `json:"activities"`
	ClientStatus *ClientStatus `json:"client_status"`
}

// Stage Instance Create
// https://discord.com/developers/docs/topics/gateway#stage-instance-create
type StageInstanceCreate struct {
	*StageInstance
}

// Stage Instance Update
// https://discord.com/developers/docs/topics/gateway#stage-instance-update
type StageInstanceUpdate struct {
	*StageInstance
}

// Stage Instance Delete
// https://discord.com/developers/docs/topics/gateway#stage-instance-delete
type StageInstanceDelete struct {
	*StageInstance
}

// Typing Start
// https://discord.com/developers/docs/topics/gateway#typing-start
type TypingStart struct {
	ChannelID string       `json:"channel_id"`
	GuildID   string       `json:"guild_id,omitempty"`
	UserID    string       `json:"user_id"`
	Member    *GuildMember `json:"member,omitempty"`
	Timestamp int          `json:"timestamp"`
}

// User Update
// https://discord.com/developers/docs/topics/gateway#user-update
type UserUpdate struct {
	*User
}

// Voice State Update
// https://discord.com/developers/docs/topics/gateway#voice-state-update
type VoiceStateUpdate struct {
	*VoiceState
}

// Voice Server Update
// https://discord.com/developers/docs/topics/gateway#voice-server-update
type VoiceServerUpdate struct {
	Token    string `json:"token"`
	GuildID  string `json:"guild_id"`
	Endpoint string `json:"endpoint"`
}

// Webhooks Update
// https://discord.com/developers/docs/topics/gateway#webhooks-update
type WebhooksUpdate struct {
	GuildID   string `json:"guild_id"`
	ChannelID string `json:"channel_id"`
}

// Gateway Payload Structure
// https://discord.com/developers/docs/topics/gateway#payloads-gateway-payload-structure
type GatewayPayload struct {
	Op             int             `json:"op"`
	Data           json.RawMessage `json:"d"`
	SequenceNumber *int64          `json:"s"`
	EventName      *string         `json:"t"`
}

// Gateway URL Query String Params
// https://discord.com/developers/docs/topics/gateway#connecting-gateway-url-query-string-params
type GatewayURLQueryString struct {
	V        int    `url:"v"`
	Encoding string `url:"encoding"`
	Compress string `url:"compress,omitempty"`
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
	FlagIntentGUILDS = 1 << 0

	// GUILD_MEMBER_ADD
	// GUILD_MEMBER_UPDATE
	// GUILD_MEMBER_REMOVE
	// THREAD_MEMBERS_UPDATE *
	FlagIntentGUILD_MEMBERS = 1 << 1

	// GUILD_BAN_ADD
	// GUILD_BAN_REMOVE
	FlagIntentGUILD_BANS = 1 << 2

	// GUILD_EMOJIS_UPDATE
	// GUILD_STICKERS_UPDATE
	FlagIntentGUILD_EMOJIS_AND_STICKERS = 1 << 3

	// GUILD_INTEGRATIONS_UPDATE
	// INTEGRATION_CREATE
	// INTEGRATION_UPDATE
	// INTEGRATION_DELETE
	FlagIntentGUILD_INTEGRATIONS = 1 << 4

	// WEBHOOKS_UPDATE
	FlagIntentGUILD_WEBHOOKS = 1 << 5

	// INVITE_CREATE
	// INVITE_DELETE
	FlagIntentGUILD_INVITES = 1 << 6

	// VOICE_STATE_UPDATE
	FlagIntentGUILD_VOICE_STATES = 1 << 7

	// PRESENCE_UPDATE
	FlagIntentGUILD_PRESENCES = 1 << 8

	// MESSAGE_CREATE
	// MESSAGE_UPDATE
	// MESSAGE_DELETE
	// MESSAGE_DELETE_BULK
	FlagIntentGUILD_MESSAGES = 1 << 9

	// MESSAGE_REACTION_ADD
	// MESSAGE_REACTION_REMOVE
	// MESSAGE_REACTION_REMOVE_ALL
	// MESSAGE_REACTION_REMOVE_EMOJI
	FlagIntentGUILD_MESSAGE_REACTIONS = 1 << 10

	// TYPING_START
	FlagIntentGUILD_MESSAGE_TYPING  = 1 << 11
	FlagIntentDIRECT_MESSAGE_TYPING = 1 << 14

	// MESSAGE_CREATE
	// MESSAGE_UPDATE
	// MESSAGE_DELETE
	// CHANNEL_PINS_UPDATE
	FlagIntentDIRECT_MESSAGES = 1 << 12

	// MESSAGE_REACTION_ADD
	// MESSAGE_REACTION_REMOVE
	// MESSAGE_REACTION_REMOVE_ALL
	// MESSAGE_REACTION_REMOVE_EMOJI
	FlagIntentDIRECT_MESSAGE_REACTIONS = 1 << 13

	FlagIntentMESSAGE_CONTENT = 1 << 15

	// GUILD_SCHEDULED_EVENT_CREATE
	// GUILD_SCHEDULED_EVENT_UPDATE
	// GUILD_SCHEDULED_EVENT_DELETE
	// GUILD_SCHEDULED_EVENT_USER_ADD
	// GUILD_SCHEDULED_EVENT_USER_REMOVE
	FlagIntentGUILD_SCHEDULED_EVENTS = 1 << 16

	// AUTO_MODERATION_RULE_CREATE
	// AUTO_MODERATION_RULE_UPDATE
	// AUTO_MODERATION_RULE_DELETE
	FlagIntentAUTO_MODERATION_CONFIGURATION = 1 << 20

	// AUTO_MODERATION_ACTION_EXECUTION
	FlagIntentAUTO_MODERATION_EXECUTION = 1 << 21
)

// Gateway Commands
// https://discord.com/developers/docs/topics/gateway#commands-and-events
type Command interface{}

// Gateway Command Names
// https://discord.com/developers/docs/topics/gateway#commands-and-events-gateway-commands
const (
	FlagGatewayCommandNameHeartbeat           = "Heartbeat"
	FlagGatewayCommandNameIdentify            = "Identify"
	FlagGatewayCommandNamePresenceUpdate      = "PresenceUpdate"
	FlagGatewayCommandNameVoiceStateUpdate    = "VoiceStateUpdate"
	FlagGatewayCommandNameResume              = "Resume"
	FlagGatewayCommandNameRequestGuildMembers = "RequestGuildMembers"
)

// Identify Structure
// https://discord.com/developers/docs/topics/gateway#identify-identify-structure
type Identify struct {
	Token          string                       `json:"token"`
	Properties     IdentifyConnectionProperties `json:"properties"`
	Compress       bool                         `json:"compress"`
	LargeThreshold int                          `json:"large_threshold,omitempty"`
	Shard          *[2]int                      `json:"shard,omitempty"`
	Presence       GatewayPresenceUpdate        `json:"presence,omitempty"`
	Intents        BitFlag                      `json:"intents"`
}

// Identify Connection Properties
// https://discord.com/developers/docs/topics/gateway#identify-identify-connection-properties
type IdentifyConnectionProperties struct {
	OS      string `json:"os"`
	Browser string `json:"browser"`
	Device  string `json:"device"`
}

// Resume Structure
// https://discord.com/developers/docs/topics/gateway#resume-resume-structure
type Resume struct {
	Token     string `json:"token"`
	SessionID string `json:"session_id"`
	Seq       int64  `json:"seq"`
}

// Heartbeat
// https://discord.com/developers/docs/topics/gateway#heartbeat
type Heartbeat struct {
	Data int64 `json:"d"`
}

// Request Guild Members Structure
// https://discord.com/developers/docs/topics/gateway#request-guild-members-guild-request-members-structure
type RequestGuildMembers struct {
	GuildID   string   `json:"guild_id"`
	Query     *string  `json:"query,omitempty"`
	Limit     int      `json:"limit"`
	Presences *bool    `json:"presences,omitempty"`
	UserIDs   []string `json:"user_ids,omitempty"`
	Nonce     *string  `json:"nonce,omitempty"`
}

// Gateway Voice State Update Structure
// https://discord.com/developers/docs/topics/gateway#update-voice-state-gateway-voice-state-update-structure
type GatewayVoiceStateUpdate struct {
	GuildID   string `json:"guild_id"`
	ChannelID string `json:"channel_id"`
	SelfMute  bool   `json:"self_mute"`
	SelfDeaf  bool   `json:"self_deaf"`
}

// Gateway Presence Update Structure
// https://discord.com/developers/docs/topics/gateway#update-presence-gateway-presence-update-structure
type GatewayPresenceUpdate struct {
	Since  int         `json:"since"`
	Game   []*Activity `json:"game"`
	Status string      `json:"status"`
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
	Limit      int     `http:"X-RateLimit-Limit,omitempty"`
	Remaining  int     `http:"X-RateLimit-Remaining,omitempty"`
	Reset      float64 `http:"X-RateLimit-Reset,omitempty"`
	ResetAfter float64 `http:"X-RateLimit-Reset-After,omitempty"`
	Bucket     string  `http:"X-RateLimit-Bucket,omitempty"`
	Global     bool    `http:"X-RateLimit-Global,omitempty"`
	Scope      string  `http:"X-RateLimit-Scope,omitempty"`
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
	WithLocalizations bool `url:"with_localizations,omitempty"`
}

// Create Global Application Command
// POST /applications/{application.id}/commands
// https://discord.com/developers/docs/interactions/application-commands#create-global-application-command
type CreateGlobalApplicationCommand struct {
	Name                     string                      `json:"name,omitempty"`
	NameLocalizations        map[Flag]string             `json:"name_localizations,omitempty"`
	Description              string                      `json:"description,omitempty"`
	DescriptionLocalizations map[Flag]string             `json:"description_localizations,omitempty"`
	Options                  []*ApplicationCommandOption `json:"options,omitempty"`
	DefaultMemberPermissions string                      `json:"default_member_permissions,omitempty"`
	DMPermission             bool                        `json:"dm_permission,omitempty"`
	Type                     Flag                        `json:"type,omitempty"`
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
	CommandID                string
	Name                     string                      `json:"name,omitempty"`
	NameLocalizations        map[string]string           `json:"name_localizations"`
	Description              string                      `json:"description,omitempty"`
	DescriptionLocalizations map[string]string           `json:"description_localizations"`
	Options                  []*ApplicationCommandOption `json:"options,omitempty"`
	DefaultMemberPermissions *string                     `json:"default_member_permissions,omitempty"`
	DMPermission             *bool                       `json:"dm_permission,omitempty"`
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
	ApplicationCommands []*ApplicationCommand `json:"commands,omitempty"`
}

// Get Guild Application Commands
// GET /applications/{application.id}/guilds/{guild.id}/commands
// https://discord.com/developers/docs/interactions/application-commands#get-guild-application-commands
type GetGuildApplicationCommands struct {
	GuildID           string
	WithLocalizations bool `url:"with_localizations,omitempty"`
}

// Create Guild Application Command
// POST /applications/{application.id}/guilds/{guild.id}/commands
// https://discord.com/developers/docs/interactions/application-commands#create-guild-application-command
type CreateGuildApplicationCommand struct {
	GuildID                  string
	Name                     string                      `json:"name"`
	NameLocalizations        map[string]string           `json:"name_localization"`
	Description              string                      `json:"description"`
	DescriptionLocalizations map[string]string           `json:"description_localizations"`
	Options                  []*ApplicationCommandOption `json:"options,omitempty"`
	DefaultMemberPermissions *string                     `json:"default_member_permissions,omitempty"`
	DMPermission             *bool                       `json:"dm_permission,omitempty"`
	Type                     *Flag                       `json:"type,omitempty"`
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
	GuildID                  string
	CommandID                string
	Name                     string                      `json:"name,omitempty"`
	NameLocalizations        map[string]string           `json:"name_localizations"`
	Description              string                      `json:"description,omitempty"`
	DescriptionLocalizations map[string]string           `json:"description_localizations"`
	Options                  []*ApplicationCommandOption `json:"options,omitempty"`
	DefaultMemberPermissions *string                     `json:"default_member_permissions,omitempty"`
	DMPermission             *bool                       `json:"dm_permission,omitempty"`
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
	GuildID                  string
	ID                       string                      `json:"id,omitempty"`
	Name                     string                      `json:"name"`
	NameLocalizations        map[string]string           `json:"name_localizations"`
	Description              string                      `json:"description"`
	DescriptionLocalizations map[string]string           `json:"description_localizations"`
	Options                  []*ApplicationCommandOption `json:"options,omitempty"`
	DefaultMemberPermissions *string                     `json:"default_member_permissions,omitempty"`
	DMPermission             *bool                       `json:"dm_permission,omitempty"`
	Type                     *Flag                       `json:"type,omitempty"`
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
	GuildID     string
	CommandID   string
	Permissions []*ApplicationCommandPermissions `json:"permissions,omitempty"`
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
	InteractionID    string
	InteractionToken string
}

// Get Original Interaction Response
// GET /webhooks/{application.id}/{interaction.token}/messages/@original
// https://discord.com/developers/docs/interactions/receiving-and-responding#get-original-interaction-response
type GetOriginalInteractionResponse struct {
	InteractionToken string
}

// Edit Original Interaction Response
// PATCH /webhooks/{application.id}/{interaction.token}/messages/@original
// https://discord.com/developers/docs/interactions/receiving-and-responding#edit-original-interaction-response
type EditOriginalInteractionResponse struct {
	InteractionToken string
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
	InteractionToken string
}

// Get Followup Message
// GET /webhooks/{application.id}/{interaction.token}/messages/{message.id}
// https://discord.com/developers/docs/interactions/receiving-and-responding#get-followup-message
type GetFollowupMessage struct {
	InteractionToken string
	MessageID        string
}

// Edit Followup Message
// PATCH /webhooks/{application.id}/{interaction.token}/messages/{message.id}
// https://discord.com/developers/docs/interactions/receiving-and-responding#edit-followup-message
type EditFollowupMessage struct {
	InteractionToken string
	MessageID        string
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
	GuildID    string
	UserID     string `url:"user_id"`
	ActionType Flag   `url:"action_type"`
	Before     string `url:"before"`
	Limit      int    `url:"limit"`
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
	GuildID         string
	Name            string                  `json:"name"`
	EventType       Flag                    `json:"event_type"`
	TriggerType     Flag                    `json:"trigger_type"`
	TriggerMetadata TriggerMetadata         `json:"trigger_metadata"`
	Actions         []*AutoModerationAction `json:"actions"`
	Enabled         bool                    `json:"enabled"`
	ExemptRoles     []string                `json:"exempt_roles,omitempty"`
	ExemptChannels  []string                `json:"exempt_channels,omitempty"`
}

// Modify Auto Moderation Rule
// PATCH /guilds/{guild.id}/auto-moderation/rules/{auto_moderation_rule.id}
// https://discord.com/developers/docs/resources/auto-moderation#modify-auto-moderation-rule
type ModifyAutoModerationRule struct {
	GuildID              string
	AutoModerationRuleID string
	Name                 string                  `json:"name"`
	EventType            Flag                    `json:"event_type"`
	TriggerType          Flag                    `json:"trigger_type"`
	TriggerMetadata      TriggerMetadata         `json:"trigger_metadata"`
	Actions              []*AutoModerationAction `json:"actions"`
	Enabled              bool                    `json:"enabled"`
	ExemptRoles          []string                `json:"exempt_roles"`
	ExemptChannels       []string                `json:"exempt_channels"`
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
	ChannelID string
	Name      string `json:"name"`
	Icon      int    `json:"icon"`
}

// Modify Channel Guild
// PATCH /channels/{channel.id}
// https://discord.com/developers/docs/resources/channel#modify-channel-json-params-guild-channel
type ModifyChannelGuild struct {
	ChannelID                  string
	Name                       *string                `json:"name"`
	Type                       *Flag                  `json:"type"`
	Position                   *int                   `json:"position"`
	Topic                      *string                `json:"topic"`
	NSFW                       bool                   `json:"nsfw"`
	RateLimitPerUser           *int                   `json:"rate_limit_per_user"`
	Bitrate                    *int                   `json:"bitrate"`
	UserLimit                  *int                   `json:"user_limit"`
	PermissionOverwrites       *[]PermissionOverwrite `json:"permission_overwrites"`
	ParentID                   *string                `json:"parent_id"`
	RTCRegion                  *string                `json:"rtc_region"`
	VideoQualityMode           *Flag                  `json:"video_quality_mode"`
	DefaultAutoArchiveDuration *int                   `json:"default_auto_archive_duration"`
}

// Modify Channel
// PATCH /channels/{channel.id}
// https://discord.com/developers/docs/resources/channel#modify-channel-json-params-thread
type ModifyChannelThread struct {
	ChannelID           string
	Name                string  `json:"name"`
	Archived            bool    `json:"archived"`
	AutoArchiveDuration int     `json:"auto_archive_duration"`
	Locked              bool    `json:"locked"`
	Invitable           bool    `json:"invitable"`
	RateLimitPerUser    *int    `json:"rate_limit_per_user"`
	Flags               BitFlag `json:"flags,omitempty"`
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
	ChannelID string
	Around    *string `url:"around,omitempty"`
	Before    *string `url:"before,omitempty"`
	After     *string `url:"after,omitempty"`
	Limit     Flag    `url:"limit,omitempty"`
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
	ChannelID        string
	Content          *string           `json:"content,omitempty"`
	TTS              *bool             `json:"tts,omitempty"`
	Embeds           []*Embed          `json:"embeds,omitempty"`
	AllowedMentions  *AllowedMentions  `json:"allowed_mentions,omitempty"`
	MessageReference *MessageReference `json:"message_reference,omitempty"`
	Components       []*Component      `json:"components,omitempty"`
	StickerIDS       []*string         `json:"sticker_ids,omitempty"`
	Files            []byte            `dasgo:"files,omitempty"`
	PayloadJSON      *string           `json:"payload_json,omitempty"`
	Attachments      []*Attachment     `json:"attachments,omitempty"`
	Flags            *BitFlag          `json:"flags,omitempty"`
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
	ChannelID string
	MessageID string
	Emoji     string
	After     string `url:"after,omitempty"`
	Limit     int    `url:"limit,omitempty"`
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
	ChannelID       string
	MessageID       string
	Content         *string          `json:"content"`
	Embeds          []*Embed         `json:"embeds"`
	Flags           *BitFlag         `json:"flags"`
	AllowedMentions *AllowedMentions `json:"allowed_mentions"`
	Components      []*Component     `json:"components"`
	Files           []byte           `dasgo:"files"`
	PayloadJSON     *string          `json:"payload_json"`
	Attachments     []*Attachment    `json:"attachments"`
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
	ChannelID string
	Messages  []*string `json:"messages"`
}

// Edit Channel Permissions
// PUT /channels/{channel.id}/permissions/{overwrite.id}
// https://discord.com/developers/docs/resources/channel#edit-channel-permissions
type EditChannelPermissions struct {
	ChannelID   string
	OverwriteID string
	Allow       *string `json:"allow,omitempty"`
	Deny        *string `json:"deny,omitempty"`
	Type        *Flag   `json:"type"`
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
	ChannelID           string
	MaxAge              *int   `json:"max_age"`
	MaxUses             *int   `json:"max_uses"`
	Temporary           bool   `json:"temporary"`
	Unique              bool   `json:"unique"`
	TargetType          Flag   `json:"target_type"`
	TargetUserID        string `json:"target_user_id"`
	TargetApplicationID string `json:"target_application_id"`
}

// Delete Channel Permission
// DELETE /channels/{channel.id}/permissions/{overwrite.id}
// https://discord.com/developers/docs/resources/channel#delete-channel-permission
type DeleteChannelPermission struct {
	ChannelID   string
	OverwriteID string
}

// Follow News Channel
// POST /channels/{channel.id}/followers
// https://discord.com/developers/docs/resources/channel#follow-news-channel
type FollowNewsChannel struct {
	ChannelID        string
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
	ChannelID   string
	UserID      string
	AccessToken string  `json:"access_token"`
	Nickname    *string `json:"nick"`
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
	ChannelID           string
	MessageID           string
	Name                string `json:"name"`
	AutoArchiveDuration int    `json:"auto_archive_duration,omitempty"`
	RateLimitPerUser    int    `json:"rate_limit_per_user,omitempty"`
}

// Start Thread without Message
// POST /channels/{channel.id}/threads
// https://discord.com/developers/docs/resources/channel#start-thread-without-message
type StartThreadwithoutMessage struct {
	ChannelID           string
	Name                string `json:"name"`
	AutoArchiveDuration int    `json:"auto_archive_duration,omitempty"`
	Type                Flag   `json:"type,omitempty"`
	Invitable           bool   `json:"invitable,omitempty"`
	RateLimitPerUser    int    `json:"rate_limit_per_user,omitempty"`
}

// Start Thread in Forum Channel
// POST /channels/{channel.id}/threads
// https://discord.com/developers/docs/resources/channel#start-thread-in-forum-channel
type StartThreadinForumChannel struct {
	ChannelID           string
	Name                string                    `json:"name"`
	AutoArchiveDuration int                       `json:"auto_archive_duration,omitempty"`
	RateLimitPerUser    int                       `json:"rate_limit_per_user,omitempty"`
	Message             *ForumThreadMessageParams `json:"message"`
}

// Forum Thread Message Params Object
// https://discord.com/developers/docs/resources/channel#start-thread-in-forum-channel-forum-thread-message-params-object
type ForumThreadMessageParams struct {
	Content         *string          `json:"content,omitempty"`
	Embeds          []*Embed         `json:"embeds,omitempty"`
	AllowedMentions *AllowedMentions `json:"allowed_mentions,omitempty"`
	Components      []*Component     `json:"components,omitempty"`
	StickerIDS      []*string        `json:"sticker_ids,omitempty"`
	Attachments     []*Attachment    `json:"attachments,omitempty"`
	Files           []byte           `dasgo:"files"`
	PayloadJSON     string           `json:"payload_json,omitempty"`
	Flags           BitFlag          `json:"flags,omitempty"`
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

// List Active Channel Threads
// GET /channels/{channel.id}/threads/active
// https://discord.com/developers/docs/resources/channel#list-active-threads
type ListActiveChannelThreads struct {
	ChannelID string
	Before    string `json:"before,omitempty"`
	Limit     int    `json:"limit,omitempty"`
}

// List Public Archived Threads
// GET /channels/{channel.id}/threads/archived/public
// https://discord.com/developers/docs/resources/channel#list-public-archived-threads
type ListPublicArchivedThreads struct {
	ChannelID string
	Before    string `url:"before,omitempty"`
	Limit     int    `url:"limit,omitempty"`
}

// List Private Archived Threads
// GET /channels/{channel.id}/threads/archived/private
// https://discord.com/developers/docs/resources/channel#list-private-archived-threads
type ListPrivateArchivedThreads struct {
	ChannelID string
	Before    string `url:"before,omitempty"`
	Limit     int    `url:"limit,omitempty"`
}

// List Joined Private Archived Threads
// GET /channels/{channel.id}/users/@me/threads/archived/private
// https://discord.com/developers/docs/resources/channel#list-joined-private-archived-threads
type ListJoinedPrivateArchivedThreads struct {
	ChannelID string
	Before    string `url:"before,omitempty"`
	Limit     int    `url:"limit,omitempty"`
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
	GuildID string
	Name    string    `json:"name"`
	Image   string    `json:"image"`
	Roles   []*string `json:"roles"`
}

// Modify Guild Emoji
// PATCH /guilds/{guild.id}/emojis/{emoji.id}
// https://discord.com/developers/docs/resources/emoji#modify-guild-emoji
type ModifyGuildEmoji struct {
	GuildID string
	EmojiID string
	Name    string    `json:"name"`
	Roles   []*string `json:"roles"`
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
	Name                        string     `json:"name"`
	Region                      *string    `json:"region,omitempty"`
	Icon                        *string    `json:"icon,omitempty"`
	VerificationLevel           *Flag      `json:"verification_level,omitempty"`
	DefaultMessageNotifications *Flag      `json:"default_message_notifications,omitempty"`
	ExplicitContentFilter       *Flag      `json:"explicit_content_filter,omitempty"`
	Roles                       []*Role    `json:"roles,omitempty"`
	Channels                    []*Channel `json:"channels,omitempty"`
	AfkChannelID                string     `json:"afk_channel_id,omitempty"`
	AfkTimeout                  int        `json:"afk_timeout,omitempty"`
	SystemChannelID             string     `json:"system_channel_id,omitempty"`
	SystemChannelFlags          BitFlag    `json:"system_channel_flags,omitempty"`
}

// Get Guild
// GET /guilds/{guild.id}
// https://discord.com/developers/docs/resources/guild#get-guild
type GetGuild struct {
	GuildID    string
	WithCounts *bool `url:"with_counts,omitempty"`
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
	GuildID                     string
	Name                        string    `json:"name"`
	Region                      string    `json:"region"`
	VerificationLevel           *Flag     `json:"verification_lvl"`
	DefaultMessageNotifications *Flag     `json:"default_message_notifications"`
	ExplicitContentFilter       *Flag     `json:"explicit_content_filter"`
	AFKChannelID                string    `json:"afk_channel_id"`
	AfkTimeout                  int       `json:"afk_timeout"`
	Icon                        *string   `json:"icon"`
	OwnerID                     string    `json:"owner_id"`
	Splash                      *string   `json:"splash"`
	DiscoverySplash             *string   `json:"discovery_splash"`
	Banner                      *string   `json:"banner"`
	SystemChannelID             *string   `json:"system_channel_id"`
	SystemChannelFlags          BitFlag   `json:"system_channel_flags"`
	RulesChannelID              *string   `json:"rules_channel_id"`
	PublicUpdatesChannelID      *string   `json:"public_updates_channel_id"`
	PreferredLocale             *string   `json:"preferred_locale"`
	Features                    []*string `json:"features"`
	Description                 *string   `json:"description"`
	PremiumProgressBarEnabled   bool      `json:"premium_progress_bar_enabled"`
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
	GuildID                    string
	Name                       string                 `json:"name"`
	Type                       *Flag                  `json:"type"`
	Topic                      *string                `json:"topic"`
	Bitrate                    *int                   `json:"bitrate"`
	UserLimit                  *int                   `json:"user_limit"`
	RateLimitPerUser           *int                   `json:"rate_limit_per_user"`
	Position                   *int                   `json:"position"`
	PermissionOverwrites       []*PermissionOverwrite `json:"permission_overwrites"`
	ParentID                   *string                `json:"parent_id"`
	NSFW                       *bool                  `json:"nsfw"`
	RTCRegion                  string                 `json:"rtc_region"`
	VideoQualityMode           *Flag                  `json:"video_quality_mode"`
	DefaultAutoArchiveDuration int                    `json:"default_auto_archive_duration"`
}

// Modify Guild Channel Positions
// PATCH /guilds/{guild.id}/channels
// https://discord.com/developers/docs/resources/guild#modify-guild-channel-positions
type ModifyGuildChannelPositions struct {
	GuildID         string
	ID              string  `json:"id"`
	Position        int     `json:"position"`
	LockPermissions bool    `json:"lock_permissions"`
	ParentID        *string `json:"parent_id"`
}

// List Active Guild Threads
// GET /guilds/{guild.id}/threads/active
// https://discord.com/developers/docs/resources/guild#list-active-guild-threads
type ListActiveGuildThreads struct {
	GuildID string
	Threads []*Channel      `json:"threads"`
	Members []*ThreadMember `json:"members"`
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
	GuildID string
	Limit   *int    `url:"limit"`
	After   *string `url:"after"`
}

// Search Guild Members
// GET /guilds/{guild.id}/members/search
// https://discord.com/developers/docs/resources/guild#search-guild-members
type SearchGuildMembers struct {
	GuildID string
	Query   *string `url:"query"`
	Limit   *int    `url:"limit"`
}

// Add Guild Member
// PUT /guilds/{guild.id}/members/{user.id}
// https://discord.com/developers/docs/resources/guild#add-guild-member
type AddGuildMember struct {
	GuildID     string
	UserID      string
	AccessToken string   `json:"access_token"`
	Nick        string   `json:"nick"`
	Roles       []string `json:"roles"`
	Mute        bool     `json:"mute"`
	Deaf        bool     `json:"deaf"`
}

// Modify Guild Member
// PATCH /guilds/{guild.id}/members/{user.id}
// https://discord.com/developers/docs/resources/guild#modify-guild-member
type ModifyGuildMember struct {
	GuildID                    string
	UserID                     string
	Nick                       string     `json:"nick"`
	Roles                      []string   `json:"roles"`
	Mute                       bool       `json:"mute"`
	Deaf                       bool       `json:"deaf"`
	ChannelID                  string     `json:"channel_id"`
	CommunicationDisabledUntil *time.Time `json:"communication_disabled_until"`
}

// Modify Current Member
// PATCH /guilds/{guild.id}/members/@me
// https://discord.com/developers/docs/resources/guild#modify-current-member
type ModifyCurrentMember struct {
	GuildID string
	Nick    *string `json:"nick,omitempty"`
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
	GuildID string
	Limit   *int    `url:"limit,omitempty"`
	Before  *string `url:"before,omitempty"`
	After   *string `url:"after,omitempty"`
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
	GuildID           string
	UserID            string
	DeleteMessageDays *int    `json:"delete_message_days,omitempty"`
	Reason            *string `json:"reason,omitempty"`
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
	GuildID      string
	Name         string  `json:"name"`
	Permissions  string  `json:"permissions"`
	Color        *int    `json:"color"`
	Hoist        bool    `json:"hoist"`
	Icon         *string `json:"icon"`
	UnicodeEmoji *string `json:"unicode_emoji"`
	Mentionable  bool    `json:"mentionable"`
}

// Modify Guild Role Positions
// PATCH /guilds/{guild.id}/roles
// https://discord.com/developers/docs/resources/guild#modify-guild-role-positions
type ModifyGuildRolePositions struct {
	GuildID  string
	ID       string `json:"id"`
	Position *int   `json:"position,omitempty"`
}

// Modify Guild Role
// PATCH /guilds/{guild.id}/roles/{role.id}
// https://discord.com/developers/docs/resources/guild#modify-guild-role
type ModifyGuildRole struct {
	GuildID      string
	RoleID       string
	Name         string  `json:"name"`
	Permissions  BitFlag `json:"permissions"`
	Color        *int    `json:"color"`
	Hoist        bool    `json:"hoist"`
	Icon         *string `json:"icon"`
	UnicodeEmoji *string `json:"unicode_emoji"`
	Mentionable  bool    `json:"mentionable"`
}

// Modify Guild MFA Level
// POST /guilds/{guild.id}/mfa
// https://discord.com/developers/docs/resources/guild#modify-guild-mfa-level
type ModifyGuildMFALevel struct {
	GuildID string
	Level   Flag `json:"level"`
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
	GuildID      string
	Days         int      `url:"days"`
	IncludeRoles []string `url:"include_roles"`
}

// Begin Guild Prune
// POST /guilds/{guild.id}/prune
// https://discord.com/developers/docs/resources/guild#begin-guild-prune
type BeginGuildPrune struct {
	GuildID           string
	Days              int      `json:"days"`
	ComputePruneCount bool     `json:"compute_prune_count"`
	IncludeRoles      []string `json:"include_roles"`
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
	GuildID string
	Code    string `json:"code,omitempty"`
	Uses    int    `json:"uses,omitempty"`
}

// Get Guild Widget Image
// GET /guilds/{guild.id}/widget.png
// https://discord.com/developers/docs/resources/guild#get-guild-widget-image
type GetGuildWidgetImage struct {
	GuildID string

	// Widget Style Options
	// https://discord.com/developers/docs/resources/guild#get-guild-widget-image-widget-style-options
	Style string `url:"style,omitempty"`
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
	GuildID         string
	Enabled         bool                    `json:"enabled"`
	WelcomeChannels []*WelcomeScreenChannel `json:"welcome_channels"`
	Description     *string                 `json:"description"`
}

// Modify Current User Voice State
// PATCH /guilds/{guild.id}/voice-states/@me
// https://discord.com/developers/docs/resources/guild#modify-current-user-voice-state
type ModifyCurrentUserVoiceState struct {
	GuildID                 string
	Suppress                bool       `json:"suppress,omitempty"`
	RequestToSpeakTimestamp *time.Time `json:"request_to_speak_timestamp,omitempty"`
}

// Modify User Voice State
// PATCH /guilds/{guild.id}/voice-states/{user.id}
// https://discord.com/developers/docs/resources/guild#modify-user-voice-state
type ModifyUserVoiceState struct {
	GuildID  string
	UserID   string
	Suppress *bool `json:"suppress,omitempty"`
}

// List Scheduled Events for Guild
// GET /guilds/{guild.id}/scheduled-events
// https://discord.com/developers/docs/resources/guild-scheduled-event#list-scheduled-events-for-guild
type ListScheduledEventsforGuild struct {
	GuildID       string
	WithUserCount *bool `url:"with_user_count,omitempty"`
}

// Create Guild Scheduled Event
// POST /guilds/{guild.id}/scheduled-events
// https://discord.com/developers/docs/resources/guild-scheduled-event#create-guild-scheduled-event
type CreateGuildScheduledEvent struct {
	GuildID            string
	ChannelID          *string                            `json:"channel_id,omitempty"`
	EntityMetadata     *GuildScheduledEventEntityMetadata `json:"entity_metadata,omitempty"`
	Name               *string                            `json:"name"`
	PrivacyLevel       Flag                               `json:"privacy_level"`
	ScheduledStartTime string                             `json:"scheduled_start_time"`
	ScheduledEndTime   string                             `json:"scheduled_end_time,omitempty"`
	Description        *string                            `json:"description,omitempty"`
	EntityType         Flag                               `json:"entity_type"`
	Image              *string                            `json:"image,omitempty"`
}

// Get Guild Scheduled Event
// GET /guilds/{guild.id}/scheduled-events/{guild_scheduled_event.id}
// https://discord.com/developers/docs/resources/guild-scheduled-event#get-guild-scheduled-event
type GetGuildScheduledEvent struct {
	GuildID               string
	GuildScheduledEventID string
	WithUserCount         *bool `url:"with_user_count,omitempty"`
}

// Modify Guild Scheduled Event
// PATCH /guilds/{guild.id}/scheduled-events/{guild_scheduled_event.id}
// https://discord.com/developers/docs/resources/guild-scheduled-event#modify-guild-scheduled-event
type ModifyGuildScheduledEvent struct {
	GuildID               string
	GuildScheduledEventID string
	ChannelID             *string                            `json:"channel_id,omitempty"`
	EntityMetadata        *GuildScheduledEventEntityMetadata `json:"entity_metadata,omitempty"`
	Name                  *string                            `json:"name,omitempty"`
	PrivacyLevel          Flag                               `json:"privacy_level,omitempty"`
	ScheduledStartTime    string                             `json:"scheduled_start_time,omitempty"`
	ScheduledEndTime      string                             `json:"scheduled_end_time,omitempty"`
	Description           *string                            `json:"description,omitempty"`
	EntityType            *Flag                              `json:"entity_type,omitempty"`
	Status                Flag                               `json:"status,omitempty"`
	Image                 *string                            `json:"image,omitempty"`
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
	GuildID               string
	GuildScheduledEventID string
	Limit                 *int    `url:"limit,omitempty"`
	WithMember            *bool   `url:"with_member,omitempty"`
	Before                *string `url:"before,omitempty"`
	After                 *string `url:"after,omitempty"`
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
	TemplateCode string
	Name         string `json:"name"`
	Icon         string `json:"icon,omitempty"`
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
	GuildID     string
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
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
	GuildID      string
	TemplateCode string
	Name         string  `json:"name,omitempty"`
	Description  *string `json:"description,omitempty"`
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
	InviteCode            string
	GuildScheduledEventID string `url:"guild_scheduled_event_id,omitempty"`
	WithCounts            *bool  `url:"with_counts,omitempty"`
	WithExpiration        *bool  `url:"with_expiration,omitempty"`
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
	ChannelID             string `json:"channel_id"`
	Topic                 string `json:"topic"`
	PrivacyLevel          Flag   `json:"privacy_level,omitempty"`
	SendStartNotification bool   `json:"send_start_notification,omitempty"`
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
	ChannelID    string
	Topic        string `json:"topic,omitempty"`
	PrivacyLevel Flag   `json:"privacy_level,omitempty"`
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
type ListNitroStickerPacks struct {
	StickerPacks []*StickerPack `json:"sticker_packs"`
}

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
	GuildID     string
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Tags        *string `json:"tags"`
	File        []byte  `dasgo:"file"`
}

// Modify Guild Sticker
// PATCH /guilds/{guild.id}/stickers/{sticker.id}
// https://discord.com/developers/docs/resources/sticker#modify-guild-sticker
type ModifyGuildSticker struct {
	GuildID     string
	StickerID   string
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	Tags        *string `json:"tags,omitempty"`
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
	Username string  `json:"username,omitempty"`
	Avatar   *string `json:"avatar,omitempty"`
}

// Get Current User Guilds
// GET /users/@me/guilds
// https://discord.com/developers/docs/resources/user#get-current-user-guilds
type GetCurrentUserGuilds struct {
	Before *string `json:"before"`
	After  *string `json:"after"`
	Limit  *int    `json:"limit"`
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
	AccessTokens []*string         `json:"access_tokens"`
	Nicks        map[string]string `json:"nicks"`
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
	ChannelID string
	Name      string `json:"name"`
	Avatar    string `json:"avatar,omitempty"`
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
	WebhookID string
	Name      string `json:"name"`
	Avatar    string `json:"avatar"`
	ChannelID string `json:"channel_id"`
}

// Modify Webhook with Token
// PATCH /webhooks/{webhook.id}/{webhook.token}
// https://discord.com/developers/docs/resources/webhook#modify-webhook-with-token
type ModifyWebhookwithToken struct {
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
	WebhookID       string
	WebhookToken    string
	Wait            bool             `url:"wait"`
	ThreadID        string           `url:"thread_id"`
	Content         string           `json:"content"`
	Username        string           `json:"username,omitempty"`
	AvatarURL       string           `json:"avatar_url,omitempty"`
	TTS             bool             `json:"tts"`
	Embeds          []*Embed         `json:"embeds"`
	AllowedMentions *AllowedMentions `json:"allowed_mentions,omitempty"`
	Components      []Component      `json:"components,omitempty"`
	Files           []byte           `dasgo:"files"`
	PayloadJSON     string           `json:"payload_json"`
	Attachments     []*Attachment    `json:"attachments,omitempty"`
	Flags           BitFlag          `json:"flags,omitempty"`
	ThreadName      string           `json:"thread_name,omitempty"`
}

// Execute Slack-Compatible Webhook
// POST /webhooks/{webhook.id}/{webhook.token}/slack
// https://discord.com/developers/docs/resources/webhook#execute-slackcompatible-webhook
type ExecuteSlackCompatibleWebhook struct {
	WebhookID    string
	WebhookToken string
	ThreadID     string `url:"thread_id"`
	Wait         bool   `url:"wait"`
}

// Execute GitHub-Compatible Webhook
// POST /webhooks/{webhook.id}/{webhook.token}/github
// https://discord.com/developers/docs/resources/webhook#execute-githubcompatible-webhook
type ExecuteGitHubCompatibleWebhook struct {
	WebhookID    string
	WebhookToken string
	ThreadID     string `url:"thread_id"`
	Wait         bool   `url:"wait"`
}

// Get Webhook Message
// GET /webhooks/{webhook.id}/{webhook.token}/messages/{message.id}
// https://discord.com/developers/docs/resources/webhook#get-webhook-message
type GetWebhookMessage struct {
	WebhookID    string
	WebhookToken string
	MessageID    string
	ThreadID     string `url:"thread_id"`
}

// Edit Webhook Message
// PATCH /webhooks/{webhook.id}/{webhook.token}/messages/{message.id}
// https://discord.com/developers/docs/resources/webhook#edit-webhook-message
type EditWebhookMessage struct {
	WebhookID       string
	WebhookToken    string
	MessageID       string
	ThreadID        string           `url:"thread_id"`
	Content         *string          `json:"content"`
	Embeds          []*Embed         `json:"embeds"`
	Components      []*Component     `json:"components"`
	Files           []byte           `dasgo:"files"`
	AllowedMentions *AllowedMentions `json:"allowed_mentions"`
	PayloadJSON     string           `json:"payload_json"`
	Attachments     []*Attachment    `json:"attachments"`
}

// Delete Webhook Message
// DELETE /webhooks/{webhook.id}/{webhook.token}/messages/{message.id}
// https://discord.com/developers/docs/resources/webhook#delete-webhook-message
type DeleteWebhookMessage struct {
	WebhookID    string
	WebhookToken string
	MessageID    string
	ThreadID     *string `url:"thread_id"`
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
	Permissions        BitFlag `url:"permissions"`
	GuildID            string  `url:"guild_id"`
	DisableGuildSelect bool    `url:"disable_guild_select"`
}

// Application Command Structure
// https://discord.com/developers/docs/interactions/application-commands#application-command-object-application-command-structure
type ApplicationCommand struct {
	ID                       string                      `json:"id"`
	Type                     Flag                        `json:"type,omitempty"`
	ApplicationID            string                      `json:"application_id"`
	GuildID                  string                      `json:"guild_id,omitempty"`
	Name                     string                      `json:"name"`
	NameLocalizations        map[Flag]string             `json:"name_localizations"`
	Description              string                      `json:"description"`
	DescriptionLocalizations map[Flag]string             `json:"description_localizations"`
	Options                  []*ApplicationCommandOption `json:"options,omitempty"`
	DefaultMemberPermissions *string                     `json:"default_member_permissions"`
	DMPermission             *bool                       `json:"dm_permission,omitempty"`
	Version                  string                      `json:"version,omitempty"`
}

// Application Command Types
// https://discord.com/developers/docs/interactions/application-commands#application-command-object-application-command-types
const (
	FlagApplicationCommandTypeCHAT_INPUT = 1
	FlagApplicationCommandTypeUSER       = 2
	FlagApplicationCommandTypeMESSAGE    = 3
)

// Application Command Option Structure
// https://discord.com/developers/docs/interactions/application-commands#application-command-object-application-command-option-structure
type ApplicationCommandOption struct {
	Type                     Flag                              `json:"type"`
	Name                     string                            `json:"name"`
	NameLocalizations        map[Flag]string                   `json:"name_localizations"`
	Description              string                            `json:"description"`
	DescriptionLocalizations map[Flag]string                   `json:"description_localizations"`
	Required                 *bool                             `json:"required,omitempty"`
	Choices                  []*ApplicationCommandOptionChoice `json:"choices,omitempty"`
	Options                  []*ApplicationCommandOption       `json:"options,omitempty"`
	ChannelTypes             []*Flag                           `json:"channel_types,omitempty"`
	MinValue                 *float64                          `json:"min_value,omitempty"`
	MaxValue                 *float64                          `json:"max_value,omitempty"`
	Autocomplete             *bool                             `json:"autocomplete,omitempty"`
}

// Application Command Option Type
// https://discord.com/developers/docs/interactions/application-commands#application-command-object-application-command-option-type
const (
	FlagApplicationCommandOptionTypeSUB_COMMAND       = 1
	FlagApplicationCommandOptionTypeSUB_COMMAND_GROUP = 2
	FlagApplicationCommandOptionTypeSTRING            = 3
	FlagApplicationCommandOptionTypeINTEGER           = 4
	FlagApplicationCommandOptionTypeBOOLEAN           = 5
	FlagApplicationCommandOptionTypeUSER              = 6
	FlagApplicationCommandOptionTypeCHANNEL           = 7
	FlagApplicationCommandOptionTypeROLE              = 8
	FlagApplicationCommandOptionTypeMENTIONABLE       = 9
	FlagApplicationCommandOptionTypeNUMBER            = 10
	FlagApplicationCommandOptionTypeATTACHMENT        = 11
)

// Application Command Option Choice
// https://discord.com/developers/docs/interactions/application-commands#application-command-object-application-command-option-choice-structure
type ApplicationCommandOptionChoice struct {
	Name              string          `json:"name"`
	NameLocalizations map[Flag]string `json:"name_localizations"`
	Value             interface{}     `json:"value"`
}

// Application Command Interaction Data Option Structure
// https://discord.com/developers/docs/interactions/application-commands#application-command-object-application-command-interaction-data-option-structure
type ApplicationCommandInteractionDataOption struct {
	Name    string                                     `json:"name"`
	Type    Flag                                       `json:"type"`
	Value   *interface{}                               `json:"value,omitempty"`
	Options []*ApplicationCommandInteractionDataOption `json:"options,omitempty"`
	Focused *bool                                      `json:"focused,omitempty"`
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
	FlagApplicationCommandPermissionTypeROLE = 1
	FlagApplicationCommandPermissionTypeUSER = 2
)

// Component Object
type Component interface {
	ComponentType() Flag
}

// Component Types
// https://discord.com/developers/docs/interactions/message-components#component-object-component-types
const (
	FlagComponentTypeActionRow  = 1
	FlagComponentTypeButton     = 2
	FlagComponentTypeSelectMenu = 3
	FlagComponentTypeTextInput  = 4
)

// https://discord.com/developers/docs/interactions/message-components#component-object
type ActionsRow struct {
	Components []Component `json:"components"`
}

// Button Object
// https://discord.com/developers/docs/interactions/message-components#button-object
type Button struct {
	Style    Flag    `json:"style"`
	Label    *string `json:"label,omitempty"`
	Emoji    *Emoji  `json:"emoji,omitempty"`
	CustomID *string `json:"custom_id,omitempty"`
	URL      *string `json:"url,omitempty"`
	Disabled *bool   `json:"disabled,omitempty"`
}

// Button Styles
// https://discord.com/developers/docs/interactions/message-components#button-object-button-styles
const (
	FlagButtonStylePRIMARY   = 1
	FlagButtonStyleBLURPLE   = 1
	FlagButtonStyleSecondary = 2
	FlagButtonStyleGREY      = 2
	FlagButtonStyleSuccess   = 3
	FlagButtonStyleGREEN     = 3
	FlagButtonStyleDanger    = 4
	FlagButtonStyleRED       = 4
	FlagButtonStyleLINK      = 5
)

// Select Menu Structure
// https://discord.com/developers/docs/interactions/message-components#select-menu-object-select-menu-structure
type SelectMenu struct {
	CustomID    string             `json:"custom_id"`
	Options     []SelectMenuOption `json:"options"`
	Placeholder *string            `json:"placeholder,omitempty"`
	MinValues   *Flag              `json:"min_values,omitempty"`
	MaxValues   *Flag              `json:"max_values,omitempty"`
	Disabled    *bool              `json:"disabled,omitempty"`
}

// Select Menu Option Structure
// https://discord.com/developers/docs/interactions/message-components#select-menu-object-select-option-structure
type SelectMenuOption struct {
	Label       string  `json:"label"`
	Value       string  `json:"value"`
	Description *string `json:"description,omitempty"`
	Emoji       *Emoji  `json:"emoji,omitempty"`
	Default     *bool   `json:"default,omitempty"`
}

// Text Input Structure
// https://discord.com/developers/docs/interactions/message-components#text-inputs-text-input-structure
type TextInput struct {
	CustomID    string  `json:"custom_id"`
	Style       Flag    `json:"style"`
	Label       *string `json:"label"`
	MinLength   *int    `json:"min_length,omitempty"`
	MaxLength   *int    `json:"max_length,omitempty"`
	Required    *bool   `json:"required,omitempty"`
	Value       *string `json:"value,omitempty"`
	Placeholder *string `json:"placeholder,omitempty"`
}

// Text Input Styles
// https://discord.com/developers/docs/interactions/message-components#text-inputs-text-input-styles
const (
	FlagTextInputStyleShort     = 1
	FlagTextInputStyleParagraph = 2
)

// Interaction Object
// https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-object-interaction-structure
type Interaction struct {
	ID             string          `json:"id"`
	ApplicationID  string          `json:"application_id"`
	Type           Flag            `json:"type"`
	Data           InteractionData `json:"data,omitempty"`
	GuildID        string          `json:"guild_id,omitempty"`
	ChannelID      string          `json:"channel_id,omitempty"`
	Member         *GuildMember    `json:"member,omitempty"`
	User           *User           `json:"user,omitempty"`
	Token          string          `json:"token"`
	Version        Flag            `json:"version"`
	Message        *Message        `json:"message,omitempty"`
	AppPermissions *BitFlag        `json:"app_permissions,omitempty"`
	Locale         *string         `json:"locale,omitempty"`
	GuildLocale    *string         `json:"guild_locale,omitempty"`
}

// Interaction Type
// https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-object-interaction-type
const (
	FlagInteractionTypePING                             = 1
	FlagInteractionTypeAPPLICATION_COMMAND              = 2
	FlagInteractionTypeMESSAGE_COMPONENT                = 3
	FlagInteractionTypeAPPLICATION_COMMAND_AUTOCOMPLETE = 4
	FlagInteractionTypeMODAL_SUBMIT                     = 5
)

// Interaction Data
type InteractionData interface {
	InteractionDataType() Flag
}

// Application Command Data Structure
// https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-object-application-command-data-structure
type ApplicationCommandData struct {
	ID       string                                     `json:"id"`
	Name     string                                     `json:"name"`
	Type     Flag                                       `json:"type"`
	Resolved *ResolvedData                              `json:"resolved,omitempty"`
	Options  []*ApplicationCommandInteractionDataOption `json:"options,omitempty"`
	GuildID  string                                     `json:"guild_id,omitempty"`
	TargetID string                                     `json:"target_id,omitempty"`
}

// Message Component Data Structure
// https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-object-message-component-data-structure
type MessageComponentData struct {
	CustomID      string              `json:"custom_id"`
	ComponentType Flag                `json:"component_type"`
	Values        []*SelectMenuOption `json:"values,omitempty"`
}

// Modal Submit Data Structure
// https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-object-modal-submit-data-structure
type ModalSubmitData struct {
	CustomID   string       `json:"custom_id"`
	Components []*Component `json:"components"`
}

// Resolved Data Structure
// https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-object-resolved-data-structure
type ResolvedData struct {
	Users       map[string]*User        `json:"users"`
	Members     map[string]*GuildMember `json:"members"`
	Roles       map[string]*Role        `json:"roles"`
	Channels    map[string]*Channel     `json:"channels"`
	Messages    map[string]*Message     `json:"messages"`
	Attachments map[string]*Attachment  `json:"attachments"`
}

// Message Interaction Structure
// https://discord.com/developers/docs/interactions/receiving-and-responding#message-interaction-object-message-interaction-structure
type MessageInteraction struct {
	ID     string       `json:"id"`
	Type   Flag         `json:"type"`
	Name   string       `json:"name"`
	User   *User        `json:"user"`
	Member *GuildMember `json:"member,omitempty"`
}

// Interaction Response Structure
// https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-response-object-interaction-response-structure
type InteractionResponse struct {
	Type Flag                     `json:"type"`
	Data *InteractionCallbackData `json:"data,omitempty"`
}

// Interaction Callback Type
// https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-response-object-interaction-callback-type
const (
	FlagInteractionCallbackTypePONG                                    = 1
	FlagInteractionCallbackTypeCHANNEL_MESSAGE_WITH_SOURCE             = 4
	FlagInteractionCallbackTypeDEFERRED_CHANNEL_MESSAGE_WITH_SOURCE    = 5
	FlagInteractionCallbackTypeDEFERRED_UPDATE_MESSAGE                 = 6
	FlagInteractionCallbackTypeUPDATE_MESSAGE                          = 7
	FlagInteractionCallbackTypeAPPLICATION_COMMAND_AUTOCOMPLETE_RESULT = 8
	FlagInteractionCallbackTypeMODAL                                   = 9
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
	CustomID   *string     `json:"custom_id"`
	Title      string      `json:"title"`
	Components []Component `json:"components"`
}

// Application Object
// https://discord.com/developers/docs/resources/application
type Application struct {
	ID                  string         `json:"id"`
	Name                string         `json:"name"`
	Icon                string         `json:"icon"`
	Description         string         `json:"description"`
	RPCOrigins          []string       `json:"rpc_origins,omitempty"`
	BotPublic           bool           `json:"bot_public"`
	BotRequireCodeGrant bool           `json:"bot_require_code_grant"`
	TermsOfServiceURL   string         `json:"terms_of_service_url,omitempty"`
	PrivacyProxyURL     string         `json:"privacy_policy_url,omitempty"`
	Owner               *User          `json:"owner,omitempty"`
	VerifyKey           string         `json:"verify_key"`
	Team                *Team          `json:"team"`
	GuildID             string         `json:"guild_id,omitempty"`
	PrimarySKUID        string         `json:"primary_sku_id,omitempty"`
	Slug                *string        `json:"slug,omitempty"`
	CoverImage          *string        `json:"cover_image,omitempty"`
	Flags               *Flag          `json:"flags,omitempty"`
	Tags                []string       `json:"tags,omitempty"`
	InstallParams       *InstallParams `json:"install_params,omitempty"`
	CustomInstallURL    string         `json:"custom_install_url,omitempty"`
}

// Application Flags
// https://discord.com/developers/docs/resources/application#application-object-application-flags
const (
	FlagApplicationGATEWAY_PRESENCE                 = 1 << 12
	FlagApplicationGATEWAY_PRESENCE_LIMITED         = 1 << 13
	FlagApplicationGATEWAY_GUILD_MEMBERS            = 1 << 14
	FlagApplicationGATEWAY_GUILD_MEMBERS_LIMITED    = 1 << 15
	FlagApplicationVERIFICATION_PENDING_GUILD_LIMIT = 1 << 16
	FlagApplicationEMBEDDED                         = 1 << 17
	FlagApplicationGATEWAY_MESSAGE_CONTENT          = 1 << 18
	FlagApplicationGATEWAY_MESSAGE_CONTENT_LIMITED  = 1 << 19
)

// Install Params Object
// https://discord.com/developers/docs/resources/application#install-params-object
type InstallParams struct {
	Scopes      []string `json:"scopes"`
	Permissions string   `json:"permissions"`
}

// Audit Log Object
// https://discord.com/developers/docs/resources/audit-log
type AuditLog struct {
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
	TargetID   string            `json:"target_id"`
	Changes    []*AuditLogChange `json:"changes,omitempty"`
	UserID     string            `json:"user_id"`
	ID         string            `json:"id"`
	ActionType Flag              `json:"action_type"`
	Options    *AuditLogOptions  `json:"options,omitempty"`
	Reason     *string           `json:"reason,omitempty"`
}

// Audit Log Events
// https://discord.com/developers/docs/resources/audit-log#audit-log-entry-object-audit-log-events
const (
	FlagAuditLogEventGUILD_UPDATE                          = 1
	FlagAuditLogEventCHANNEL_CREATE                        = 10
	FlagAuditLogEventCHANNEL_UPDATE                        = 11
	FlagAuditLogEventCHANNEL_DELETE                        = 12
	FlagAuditLogEventCHANNEL_OVERWRITE_CREATE              = 13
	FlagAuditLogEventCHANNEL_OVERWRITE_UPDATE              = 14
	FlagAuditLogEventCHANNEL_OVERWRITE_DELETE              = 15
	FlagAuditLogEventMEMBER_KICK                           = 20
	FlagAuditLogEventMEMBER_PRUNE                          = 21
	FlagAuditLogEventMEMBER_BAN_ADD                        = 22
	FlagAuditLogEventMEMBER_BAN_REMOVE                     = 23
	FlagAuditLogEventMEMBER_UPDATE                         = 24
	FlagAuditLogEventMEMBER_ROLE_UPDATE                    = 25
	FlagAuditLogEventMEMBER_MOVE                           = 26
	FlagAuditLogEventMEMBER_DISCONNECT                     = 27
	FlagAuditLogEventBOT_ADD                               = 28
	FlagAuditLogEventROLE_CREATE                           = 30
	FlagAuditLogEventROLE_UPDATE                           = 31
	FlagAuditLogEventROLE_DELETE                           = 32
	FlagAuditLogEventINVITE_CREATE                         = 40
	FlagAuditLogEventINVITE_UPDATE                         = 41
	FlagAuditLogEventINVITE_DELETE                         = 42
	FlagAuditLogEventWEBHOOK_CREATE                        = 50
	FlagAuditLogEventWEBHOOK_UPDATE                        = 51
	FlagAuditLogEventWEBHOOK_DELETE                        = 52
	FlagAuditLogEventEMOJI_CREATE                          = 60
	FlagAuditLogEventEMOJI_UPDATE                          = 61
	FlagAuditLogEventEMOJI_DELETE                          = 62
	FlagAuditLogEventMESSAGE_DELETE                        = 72
	FlagAuditLogEventMESSAGE_BULK_DELETE                   = 73
	FlagAuditLogEventMESSAGE_PIN                           = 74
	FlagAuditLogEventMESSAGE_UNPIN                         = 75
	FlagAuditLogEventINTEGRATION_CREATE                    = 80
	FlagAuditLogEventINTEGRATION_UPDATE                    = 81
	FlagAuditLogEventINTEGRATION_DELETE                    = 82
	FlagAuditLogEventSTAGE_INSTANCE_CREATE                 = 83
	FlagAuditLogEventSTAGE_INSTANCE_UPDATE                 = 84
	FlagAuditLogEventSTAGE_INSTANCE_DELETE                 = 85
	FlagAuditLogEventSTICKER_CREATE                        = 90
	FlagAuditLogEventSTICKER_UPDATE                        = 91
	FlagAuditLogEventSTICKER_DELETE                        = 92
	FlagAuditLogEventGUILD_SCHEDULED_EVENT_CREATE          = 100
	FlagAuditLogEventGUILD_SCHEDULED_EVENT_UPDATE          = 101
	FlagAuditLogEventGUILD_SCHEDULED_EVENT_DELETE          = 102
	FlagAuditLogEventTHREAD_CREATE                         = 110
	FlagAuditLogEventTHREAD_UPDATE                         = 111
	FlagAuditLogEventTHREAD_DELETE                         = 112
	FlagAuditLogEventAPPLICATION_COMMAND_PERMISSION_UPDATE = 121
)

// Optional Audit Entry Info
// https://discord.com/developers/docs/resources/audit-log#audit-log-entry-object-optional-audit-entry-info
type AuditLogOptions struct {
	ApplicationID    string `json:"application_id"`
	ChannelID        string `json:"channel_id"`
	Count            string `json:"count"`
	DeleteMemberDays string `json:"delete_member_days"`
	ID               string `json:"id"`
	MembersRemoved   string `json:"members_removed"`
	MessageID        string `json:"message_id"`
	RoleName         string `json:"role_name"`
	Type             string `json:"type"`
}

// Audit Log Change Object
// https://discord.com/developers/docs/resources/audit-log#audit-log-change-object
type AuditLogChange struct {
	NewValue interface{} `json:"new_value,omitempty"`
	OldValue interface{} `json:"old_value,omitempty"`
	Key      string      `json:"key"`
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
	EventType       Flag                    `json:"event_type"`
	TriggerType     Flag                    `json:"trigger_type"`
	TriggerMetadata TriggerMetadata         `json:"trigger_metadata"`
	Actions         []*AutoModerationAction `json:"actions"`
	Enabled         bool                    `json:"enabled"`
	ExemptRoles     []string                `json:"exempt_roles"`
	ExemptChannels  []string                `json:"exempt_channels"`
}

// Trigger Types
// https://discord.com/developers/docs/resources/auto-moderation#auto-moderation-rule-object-trigger-types
const (
	FlagTriggerTypeKEYWORD        = 1
	FlagTriggerTypeHARMFUL_LINK   = 2
	FlagTriggerTypeSPAM           = 3
	FlagTriggerTypeKEYWORD_PRESET = 4
)

// Trigger Metadata
// https://discord.com/developers/docs/resources/auto-moderation#auto-moderation-rule-object-trigger-metadata
type TriggerMetadata struct {
	// https://discord.com/developers/docs/resources/auto-moderation#auto-moderation-rule-object-keyword-matching-strategies
	KeywordFilter []string `json:"keyword_filter"`
	Presets       []Flag   `json:"presets"`
	AllowList     []string `json:"allow_list"`
}

// Keyword Preset Types
// https://discord.com/developers/docs/resources/auto-moderation#auto-moderation-rule-object-keyword-preset-types
const (
	FlagKeywordPresetTypePROFANITY      = 1
	FlagKeywordPresetTypeSEXUAL_CONTENT = 2
	FlagKeywordPresetTypeSLURS          = 3
)

// Event Types
// https://discord.com/developers/docs/resources/auto-moderation#auto-moderation-rule-object-event-types
const (
	FlagEventTypeMESSAGE_SEND = 1
)

// Auto Moderation Action Structure
// https://discord.com/developers/docs/resources/auto-moderation#auto-moderation-action-object
type AutoModerationAction struct {
	Type     Flag            `json:"type"`
	Metadata *ActionMetadata `json:"metadata,omitempty"`
}

// Action Types
// https://discord.com/developers/docs/resources/auto-moderation#auto-moderation-action-object-action-types
const (
	FlagActionTypeBLOCK_MESSAGE      = 1
	FlagActionTypeSEND_ALERT_MESSAGE = 2
	FlagActionTypeTIMEOUT            = 3
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
	ID                         string                `json:"id"`
	Type                       *Flag                 `json:"type"`
	GuildID                    string                `json:"guild_id,omitempty"`
	Position                   *int                  `json:"position,omitempty"`
	PermissionOverwrites       []PermissionOverwrite `json:"permission_overwrites,omitempty"`
	Name                       string                `json:"name,omitempty"`
	Topic                      *string               `json:"topic,omitempty"`
	NSFW                       *bool                 `json:"nsfw,omitempty"`
	LastMessageID              *string               `json:"last_message_id"`
	Bitrate                    *Flag                 `json:"bitrate,omitempty"`
	UserLimit                  *Flag                 `json:"user_limit,omitempty"`
	RateLimitPerUser           *int                  `json:"rate_limit_per_user,omitempty"`
	Recipients                 []*User               `json:"recipients,omitempty"`
	Icon                       *string               `json:"icon"`
	OwnerID                    string                `json:"owner_id,omitempty"`
	ApplicationID              string                `json:"application_id,omitempty"`
	ParentID                   *string               `json:"parent_id"`
	LastPinTimestamp           time.Time             `json:"last_pin_timestamp"`
	RTCRegion                  *string               `json:"rtc_region"`
	VideoQualityMode           Flag                  `json:"video_quality_mode,omitempty"`
	MessageCount               *int                  `json:"message_count,omitempty"`
	MemberCount                *int                  `json:"member_count,omitempty"`
	ThreadMetadata             *ThreadMetadata       `json:"thread_metadata,omitempty"`
	Member                     *ThreadMember         `json:"member,omitempty"`
	DefaultAutoArchiveDuration int                   `json:"default_auto_archive_duration,omitempty"`
	Permissions                *string               `json:"permissions,omitempty"`
	Flags                      BitFlag               `json:"flags,omitempty"`
	TotalMessageSent           int                   `json:"total_message_sent,omitempty"`
}

// Channel Types
// https://discord.com/developers/docs/resources/channel#channel-object-channel-types
const (
	FlagChannelTypeGUILD_TEXT           = 0
	FlagChannelTypeDM                   = 1
	FlagChannelTypeGUILD_VOICE          = 2
	FlagChannelTypeGROUP_DM             = 3
	FlagChannelTypeGUILD_CATEGORY       = 4
	FlagChannelTypeGUILD_NEWS           = 5
	FlagChannelTypeGUILD_NEWS_THREAD    = 10
	FlagChannelTypeGUILD_PUBLIC_THREAD  = 11
	FlagChannelTypeGUILD_PRIVATE_THREAD = 12
	FlagChannelTypeGUILD_STAGE_VOICE    = 13
	FlagChannelTypeGUILD_DIRECTORY      = 14
	FlagChannelTypeGUILD_FORUM          = 15
)

// Video Quality Modes
// https://discord.com/developers/docs/resources/channel#channel-object-video-quality-modes
const (
	FlagVideoQualityModeAUTO = 1
	FlagVideoQualityModeFULL = 2
)

// Channel Flags
// https://discord.com/developers/docs/resources/channel#channel-object-channel-flags
const (
	FlagChannelPINNED = 1 << 1
)

// Message Object
// https://discord.com/developers/docs/resources/channel#message-object
type Message struct {
	ID                string            `json:"id"`
	ChannelID         string            `json:"channel_id"`
	Author            *User             `json:"author"`
	Content           string            `json:"content"`
	Timestamp         time.Time         `json:"timestamp"`
	EditedTimestamp   *time.Time        `json:"edited_timestamp"`
	TTS               bool              `json:"tts"`
	MentionEveryone   bool              `json:"mention_everyone"`
	Mentions          []*User           `json:"mentions"`
	MentionRoles      []*string         `json:"mention_roles"`
	MentionChannels   []*ChannelMention `json:"mention_channels,omitempty"`
	Attachments       []*Attachment     `json:"attachments"`
	Embeds            []*Embed          `json:"embeds"`
	Reactions         []*Reaction       `json:"reactions,omitempty"`
	Nonce             interface{}       `json:"nonce,omitempty"`
	Pinned            bool              `json:"pinned"`
	WebhookID         string            `json:"webhook_id,omitempty"`
	Type              Flag              `json:"type"`
	Activity          *MessageActivity  `json:"activity,omitempty"`
	Application       *Application      `json:"application,omitempty"`
	ApplicationID     string            `json:"application_id,omitempty"`
	MessageReference  *MessageReference `json:"message_reference,omitempty"`
	Flags             *BitFlag          `json:"flags,omitempty"`
	ReferencedMessage *Message          `json:"referenced_message"`
	Interaction       *Interaction      `json:"interaction"`
	Thread            *Channel          `json:"thread"`
	Components        []*Component      `json:"components"`
	StickerItems      []*StickerItem    `json:"sticker_items"`
	Stickers          []*Sticker        `json:"stickers"`
	Position          int               `json:"position,omitempty"`

	// MessageCreate Event Extra Fields
	// https://discord.com/developers/docs/topics/gateway#message-create-message-create-extra-fields
	GuildID string       `json:"guild_id,omitempty"`
	Member  *GuildMember `json:"member,omitempty"`
}

// Message Types
// https://discord.com/developers/docs/resources/channel#message-object-message-types
const (
	FlagMessageTypeDEFAULT                                      = 0
	FlagMessageTypeRECIPIENT_ADD                                = 1
	FlagMessageTypeRECIPIENT_REMOVE                             = 2
	FlagMessageTypeCALL                                         = 3
	FlagMessageTypeCHANNEL_NAME_CHANGE                          = 4
	FlagMessageTypeCHANNEL_ICON_CHANGE                          = 5
	FlagMessageTypeCHANNEL_PINNED_MESSAGE                       = 6
	FlagMessageTypeUSER_JOIN                                    = 7
	FlagMessageTypeGUILD_BOOST                                  = 8
	FlagMessageTypeGUILD_BOOST_TIER_1                           = 9
	FlagMessageTypeGUILD_BOOST_TIER_2                           = 10
	FlagMessageTypeGUILD_BOOST_TIER_3                           = 11
	FlagMessageTypeCHANNEL_FOLLOW_ADD                           = 12
	FlagMessageTypeGUILD_DISCOVERY_DISQUALIFIED                 = 14
	FlagMessageTypeGUILD_DISCOVERY_REQUALIFIED                  = 15
	FlagMessageTypeGUILD_DISCOVERY_GRACE_PERIOD_INITIAL_WARNING = 16
	FlagMessageTypeGUILD_DISCOVERY_GRACE_PERIOD_FINAL_WARNING   = 17
	FlagMessageTypeTHREAD_CREATED                               = 18
	FlagMessageTypeREPLY                                        = 19
	FlagMessageTypeCHAT_INPUT_COMMAND                           = 20
	FlagMessageTypeTHREAD_STARTER_MESSAGE                       = 21
	FlagMessageTypeGUILD_INVITE_REMINDER                        = 22
	FlagMessageTypeCONTEXT_MENU_COMMAND                         = 23
)

// Message Activity Structure
// https://discord.com/developers/docs/resources/channel#message-object-message-activity-structure
type MessageActivity struct {
	Type    int     `json:"type"`
	PartyID *string `json:"party_id,omitempty"`
}

// Message Activity Types
// https://discord.com/developers/docs/resources/channel#message-object-message-activity-types
const (
	FlagMessageActivityTypeJOIN         = 1
	FlagMessageActivityTypeSPECTATE     = 2
	FlagMessageActivityTypeLISTEN       = 3
	FlagMessageActivityTypeJOIN_REQUEST = 5
)

// Message Flags
// https://discord.com/developers/docs/resources/channel#message-object-message-flags
const (
	FlagMessageCROSSPOSTED                            = 1 << 0
	FlagMessageIS_CROSSPOST                           = 1 << 1
	FlagMessageSUPPRESS_EMBEDS                        = 1 << 2
	FlagMessageSOURCE_MESSAGE_DELETED                 = 1 << 3
	FlagMessageURGENT                                 = 1 << 4
	FlagMessageHAS_THREAD                             = 1 << 5
	FlagMessageEPHEMERAL                              = 1 << 6
	FlagMessageLOADING                                = 1 << 7
	FlagMessageFAILED_TO_MENTION_SOME_ROLES_IN_THREAD = 1 << 8
)

// Message Reference Object
// https://discord.com/developers/docs/resources/channel#message-reference-object
type MessageReference struct {
	MessageID       string  `json:"message_id,omitempty"`
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
	Count int    `json:"count"`
	Me    bool   `json:"me"`
	Emoji *Emoji `json:"emoji"`
}

// Overwrite Object
// https://discord.com/developers/docs/resources/channel#overwrite-object
type PermissionOverwrite struct {
	ID    string `json:"id"`
	Type  *Flag  `json:"type"`
	Deny  string `json:"deny"`
	Allow string `json:"allow"`
}

// Thread Metadata Object
// https://discord.com/developers/docs/resources/channel#thread-metadata-object
type ThreadMetadata struct {
	Archived            bool       `json:"archived"`
	AutoArchiveDuration int        `json:"auto_archive_duration"`
	Locked              bool       `json:"locked"`
	Invitable           *bool      `json:"invitable,omitempty"`
	CreateTimestamp     *time.Time `json:"create_timestamp"`
}

// Thread Member Object
// https://discord.com/developers/docs/resources/channel#thread-member-object
type ThreadMember struct {
	ThreadID      string    `json:"id,omitempty"`
	UserID        string    `json:"user_id,omitempty"`
	JoinTimestamp time.Time `json:"join_timestamp"`
	Flags         Flag      `json:"flags"`
}

// Embed Object
// https://discord.com/developers/docs/resources/channel#embed-object
type Embed struct {
	Title       *string         `json:"title,omitempty"`
	Type        *string         `json:"type,omitempty"`
	Description *string         `json:"description,omitempty"`
	URL         *string         `json:"url,omitempty"`
	Timestamp   time.Time       `json:"timestamp,omitempty"`
	Color       int             `json:"color,omitempty"`
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
	URL      string  `json:"url"`
	ProxyURL *string `json:"proxy_url,omitempty"`
	Height   *int    `json:"height,omitempty"`
	Width    *int    `json:"width,omitempty"`
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
	URL      string  `json:"url"`
	ProxyURL *string `json:"proxy_url,omitempty"`
	Height   *int    `json:"height,omitempty"`
	Width    *int    `json:"width,omitempty"`
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
	Name         string  `json:"name"`
	URL          *string `json:"url,omitempty"`
	IconURL      *string `json:"icon_url,omitempty"`
	ProxyIconURL *string `json:"proxy_icon_url,omitempty"`
}

// Embed Footer Structure
// https://discord.com/developers/docs/resources/channel#embed-object-embed-footer-structure
type EmbedFooter struct {
	Text         string  `json:"text"`
	IconURL      *string `json:"icon_url,omitempty"`
	ProxyIconURL *string `json:"proxy_icon_url,omitempty"`
}

// Embed Field Structure
// https://discord.com/developers/docs/resources/channel#embed-object-embed-field-structure
type EmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline *bool  `json:"inline,omitempty"`
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
	ID          string  `json:"id"`
	Filename    string  `json:"filename"`
	Description *string `json:"description,omitempty"`
	ContentType *string `json:"content_type,omitempty"`
	Size        int     `json:"size"`
	URL         string  `json:"url"`
	ProxyURL    *string `json:"proxy_url"`
	Height      *int    `json:"height"`
	Width       *int    `json:"width"`
	Emphemeral  *bool   `json:"ephemeral,omitempty"`
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
	ID            string   `json:"id"`
	Name          *string  `json:"name"`
	Roles         []string `json:"roles,omitempty"`
	User          *User    `json:"user,omitempty"`
	RequireColons *bool    `json:"require_colons,omitempty"`
	Managed       *bool    `json:"managed,omitempty"`
	Animated      *bool    `json:"animated,omitempty"`
	Available     *bool    `json:"available,omitempty"`
}

// Guild Object
// https://discord.com/developers/docs/resources/guild#guild-object
type Guild struct {
	ID                          string         `json:"id"`
	Name                        string         `json:"name"`
	Icon                        string         `json:"icon"`
	IconHash                    *string        `json:"icon_hash"`
	Splash                      string         `json:"splash"`
	DiscoverySplash             string         `json:"discovery_splash"`
	Owner                       *bool          `json:"owner,omitempty"`
	OwnerID                     string         `json:"owner_id"`
	Permissions                 *string        `json:"permissions,omitempty"`
	Region                      *string        `json:"region"`
	AfkChannelID                string         `json:"afk_channel_id"`
	AfkTimeout                  int            `json:"afk_timeout"`
	WidgetEnabled               *bool          `json:"widget_enabled,omitempty"`
	WidgetChannelID             *string        `json:"widget_channel_id"`
	VerificationLevel           *Flag          `json:"verification_level"`
	DefaultMessageNotifications *Flag          `json:"default_message_notifications"`
	ExplicitContentFilter       *Flag          `json:"explicit_content_filter"`
	Roles                       []*Role        `json:"roles"`
	Emojis                      []*Emoji       `json:"emojis"`
	Features                    []*string      `json:"features"`
	MFALevel                    *Flag          `json:"mfa_level"`
	ApplicationID               string         `json:"application_id"`
	SystemChannelID             string         `json:"system_channel_id"`
	SystemChannelFlags          BitFlag        `json:"system_channel_flags"`
	RulesChannelID              string         `json:"rules_channel_id"`
	MaxPresences                *int           `json:"max_presences"`
	MaxMembers                  int            `json:"max_members,omitempty"`
	VanityUrl                   *string        `json:"vanity_url_code"`
	Description                 *string        `json:"description"`
	Banner                      *string        `json:"banner"`
	PremiumTier                 *Flag          `json:"premium_tier"`
	PremiumSubscriptionCount    *int           `json:"premium_subscription_count,omitempty"`
	PreferredLocale             string         `json:"preferred_locale"`
	PublicUpdatesChannelID      string         `json:"public_updates_channel_id"`
	MaxVideoChannelUsers        *int           `json:"max_video_channel_users,omitempty"`
	ApproximateMemberCount      *int           `json:"approximate_member_count,omitempty"`
	ApproximatePresenceCount    *int           `json:"approximate_presence_count,omitempty"`
	WelcomeScreen               *WelcomeScreen `json:"welcome_screen,omitempty"`
	NSFWLevel                   *Flag          `json:"nsfw_level"`
	Stickers                    []*Sticker     `json:"stickers,omitempty"`
	PremiumProgressBarEnabled   bool           `json:"premium_progress_bar_enabled"`

	// Unavailable Guild Object
	// https://discord.com/developers/docs/resources/guild#unavailable-guild-object
	Unavailable bool `json:"unavailable,omitempty"`
}

// Default Message Notification Level
// https://discord.com/developers/docs/resources/guild#guild-object-default-message-notification-level
const (
	FlagDefaultMessageNotificationLevelALL_MESSAGES  = 0
	FlagDefaultMessageNotificationLevelONLY_MENTIONS = 1
)

// Explicit Content Filter Level
// https://discord.com/developers/docs/resources/guild#guild-object-explicit-content-filter-level
const (
	FlagExplicitContentFilterLevelDISABLED              = 0
	FlagExplicitContentFilterLevelMEMBERS_WITHOUT_ROLES = 1
	FlagExplicitContentFilterLevelALL_MEMBERS           = 2
)

// MFA Level
// https://discord.com/developers/docs/resources/guild#guild-object-mfa-level
const (
	FlagMFALevelNONE     = 0
	FlagMFALevelELEVATED = 1
)

// Verification Level
// https://discord.com/developers/docs/resources/guild#guild-object-verification-level
const (
	FlagVerificationLevelNONE      = 0
	FlagVerificationLevelLOW       = 1
	FlagVerificationLevelMEDIUM    = 2
	FlagVerificationLevelHIGH      = 3
	FlagVerificationLevelVERY_HIGH = 4
)

// Guild NSFW Level
// https://discord.com/developers/docs/resources/guild#guild-object-guild-nsfw-level
const (
	FlagGuildNSFWLevelDEFAULT        = 0
	FlagGuildNSFWLevelEXPLICIT       = 1
	FlagGuildNSFWLevelSAFE           = 2
	FlagGuildNSFWLevelAGE_RESTRICTED = 3
)

// Premium Tier
// https://discord.com/developers/docs/resources/guild#guild-object-premium-tier
const (
	FlagPremiumTierNONE  = 0
	FlagPremiumTierONE   = 1
	FlagPremiumTierTWO   = 2
	FlagPremiumTierTHREE = 3
)

// System Channel Flags
// https://discord.com/developers/docs/resources/guild#guild-object-system-channel-flags
const (
	FlagSystemChannelSUPPRESS_JOIN_NOTIFICATIONS           = 1 << 0
	FlagSystemChannelSUPPRESS_PREMIUM_SUBSCRIPTIONS        = 1 << 1
	FlagSystemChannelSUPPRESS_GUILD_REMINDER_NOTIFICATIONS = 1 << 2
	FlagSystemChannelSUPPRESS_JOIN_NOTIFICATION_REPLIES    = 1 << 3
)

// Guild Features
// https://discord.com/developers/docs/resources/guild#guild-object-guild-features
const (
	FlagGuildFeatureANIMATED_BANNER                  = "ANIMATED_BANNER"
	FlagGuildFeatureANIMATED_ICON                    = "ANIMATED_ICON"
	FlagGuildFeatureBANNER                           = "BANNER"
	FlagGuildFeatureCOMMERCE                         = "COMMERCE"
	FlagGuildFeatureCOMMUNITY                        = "COMMUNITY"
	FlagGuildFeatureDISCOVERABLE                     = "DISCOVERABLE"
	FlagGuildFeatureFEATURABLE                       = "FEATURABLE"
	FlagGuildFeatureINVITE_SPLASH                    = "INVITE_SPLASH"
	FlagGuildFeatureMEMBER_VERIFICATION_GATE_ENABLED = "MEMBER_VERIFICATION_GATE_ENABLED"
	FlagGuildFeatureMONETIZATION_ENABLED             = "MONETIZATION_ENABLED"
	FlagGuildFeatureMORE_STICKERS                    = "MORE_STICKERS"
	FlagGuildFeatureNEWS                             = "NEWS"
	FlagGuildFeaturePARTNERED                        = "PARTNERED"
	FlagGuildFeaturePREVIEW_ENABLED                  = "PREVIEW_ENABLED"
	FlagGuildFeaturePRIVATE_THREADS                  = "PRIVATE_THREADS"
	FlagGuildFeatureROLE_ICONS                       = "ROLE_ICONS"
	FlagGuildFeatureSEVEN_DAY_THREAD_ARCHIVE         = "SEVEN_DAY_THREAD_ARCHIVE"
	FlagGuildFeatureTHREE_DAY_THREAD_ARCHIVE         = "THREE_DAY_THREAD_ARCHIVE"
	FlagGuildFeatureTICKETED_EVENTS_ENABLED          = "TICKETED_EVENTS_ENABLED"
	FlagGuildFeatureVANITY_URL                       = "VANITY_URL"
	FlagGuildFeatureVERIFIED                         = "VERIFIED"
	FlagGuildFeatureVIP_REGIONS                      = "VIP_REGIONS"
	FlagGuildFeatureWELCOME_SCREEN_ENABLED           = "WELCOME_SCREEN_ENABLED"
)

// Guild Preview Object
// https://discord.com/developers/docs/resources/guild#guild-preview-object-guild-preview-structure
type GuildPreview struct {
	ID                       string     `json:"id"`
	Name                     string     `json:"name"`
	Icon                     string     `json:"icon"`
	Splash                   string     `json:"splash"`
	DiscoverySplash          string     `json:"discovery_splash"`
	Emojis                   []*Emoji   `json:"emojis"`
	Features                 []*string  `json:"features"`
	ApproximateMemberCount   int        `json:"approximate_member_count"`
	ApproximatePresenceCount int        `json:"approximate_presence_count"`
	Description              *string    `json:"description"`
	Stickers                 []*Sticker `json:"stickers"`
}

// Guild Widget Settings Object
// https://discord.com/developers/docs/resources/guild#guild-widget-settings-object
type GuildWidgetSettings struct {
	Enabled   bool   `json:"enabled"`
	ChannelID string `json:"channel_id"`
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
	User                       *User      `json:"user,omitempty"`
	Nick                       *string    `json:"nick"`
	Avatar                     *string    `json:"avatar"`
	Roles                      []*string  `json:"roles"`
	GuildID                    string     `json:"guild_id,omitempty"`
	JoinedAt                   time.Time  `json:"joined_at"`
	PremiumSince               *time.Time `json:"premium_since"`
	Deaf                       bool       `json:"deaf"`
	Mute                       bool       `json:"mute"`
	Pending                    *bool      `json:"pending,omitempty"`
	Permissions                *string    `json:"permissions,omitempty"`
	CommunicationDisabledUntil time.Time  `json:"communication_disabled_until"`
}

// Integration Object
// https://discord.com/developers/docs/resources/guild#integration-object
type Integration struct {
	ID                string             `json:"id"`
	Name              string             `json:"name"`
	Type              string             `json:"type"`
	Enabled           *bool              `json:"enabled,omitempty"`
	Syncing           *bool              `json:"syncing,omitempty"`
	RoleID            string             `json:"role_id,omitempty"`
	EnableEmoticons   *bool              `json:"enable_emoticons,omitempty"`
	ExpireBehavior    *Flag              `json:"expire_behavior,omitempty"`
	ExpireGracePeriod *int               `json:"expire_grace_period,omitempty"`
	User              *User              `json:"user,omitempty"`
	Account           IntegrationAccount `json:"account"`
	SyncedAt          *time.Time         `json:"synced_at,omitempty"`
	SubscriberCount   *int               `json:"subscriber_count,omitempty"`
	Revoked           *bool              `json:"revoked,omitempty"`
	Application       *Application       `json:"application,omitempty"`
}

// Integration Expire Behaviors
// https://discord.com/developers/docs/resources/guild#integration-object-integration-expire-behaviors
const (
	FlagIntegrationExpireBehaviorREMOVEROLE = 0
	FlagIntegrationExpireBehaviorKICK       = 1
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
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Icon        *string `json:"icon"`
	Description *string `json:"description"`
	Bot         *User   `json:"bot,omitempty"`
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
	ChannelID   string  `json:"channel_id"`
	Description *string `json:"description"`
	EmojiID     *string `json:"emoji_id"`
	EmojiName   *string `json:"emoji_name"`
}

// Guild Scheduled Event Object
// https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event-object-guild-scheduled-event-structure
type GuildScheduledEvent struct {
	ID                 string                             `json:"id"`
	GuildID            string                             `json:"guild_id"`
	ChannelID          *string                            `json:"channel_id"`
	CreatorID          *string                            `json:"creator_id"`
	Name               string                             `json:"name"`
	Description        *string                            `json:"description"`
	ScheduledStartTime time.Time                          `json:"scheduled_start_time"`
	ScheduledEndTime   time.Time                          `json:"scheduled_end_time"`
	PrivacyLevel       Flag                               `json:"privacy_level"`
	Status             Flag                               `json:"status"`
	EntityType         Flag                               `json:"entity_type"`
	EntityID           *string                            `json:"entity_id"`
	EntityMetadata     *GuildScheduledEventEntityMetadata `json:"entity_metadata"`
	Creator            *User                              `json:"creator,omitempty"`
	UserCount          int                                `json:"user_count,omitempty"`
	Image              string                             `json:"image,omitempty"`
}

// Guild Scheduled Event Privacy Level
// https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event-object-guild-scheduled-event-privacy-level
const (
	FlagGuildScheduledEventPrivacyLevelGUILD_ONLY = 2
)

// Guild Scheduled Event Entity Types
// https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event-object-guild-scheduled-event-entity-types
const (
	FlagGuildScheduledEventEntityTypeSTAGE_INSTANCE = 1
	FlagGuildScheduledEventEntityTypeVOICE          = 2
	FlagGuildScheduledEventEntityTypeEXTERNAL       = 3
)

// Guild Scheduled Event Status
// https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event-object-guild-scheduled-event-status
const (
	FlagGuildScheduledEventStatusSCHEDULED = 1
	FlagGuildScheduledEventStatusACTIVE    = 2
	FlagGuildScheduledEventStatusCOMPLETED = 3
	FlagGuildScheduledEventStatusCANCELED  = 4
)

// Guild Scheduled Event Entity Metadata
// https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event-object-guild-scheduled-event-entity-metadata
type GuildScheduledEventEntityMetadata struct {
	Location string `json:"location,omitempty"`
}

// Guild Scheduled Event User Object
// https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event-user-object-guild-scheduled-event-user-structure
type GuildScheduledEventUser struct {
	GuildScheduledEventID string       `json:"guild_scheduled_event_id"`
	User                  *User        `json:"user"`
	Member                *GuildMember `json:"member,omitempty"`
}

// Guild Template Object
// https://discord.com/developers/docs/resources/guild-template#guild-template-object
type GuildTemplate struct {
	Code                  string    `json:"code"`
	Name                  string    `json:"name"`
	Description           *string   `json:"description"`
	UsageCount            *int      `json:"usage_count"`
	CreatorID             string    `json:"creator_id"`
	Creator               *User     `json:"creator"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
	SourceGuildID         string    `json:"source_guild_id"`
	SerializedSourceGuild *Guild    `json:"serialized_source_guild"`
	IsDirty               *bool     `json:"is_dirty"`
}

// Invite Object
// https://discord.com/developers/docs/resources/invite#invite-object
type Invite struct {
	Code                     string               `json:"code"`
	Guild                    *Guild               `json:"guild,omitempty"`
	Channel                  *Channel             `json:"channel"`
	Inviter                  *User                `json:"inviter,omitempty"`
	TargetType               Flag                 `json:"target_type,omitempty"`
	TargetUser               *User                `json:"target_user,omitempty"`
	TargetApplication        *Application         `json:"target_application,omitempty"`
	ApproximatePresenceCount int                  `json:"approximate_presence_count,omitempty"`
	ApproximateMemberCount   int                  `json:"approximate_member_count,omitempty"`
	ExpiresAt                *time.Time           `json:"expires_at"`
	StageInstance            StageInstance        `json:"stage_instance,omitempty"`
	GuildScheduledEvent      *GuildScheduledEvent `json:"guild_scheduled_event,omitempty"`
}

// Invite Target Types
// https://discord.com/developers/docs/resources/invite#invite-object-invite-target-types
const (
	FlagInviteTargetTypeSTREAM               = 1
	FlagInviteTargetTypeEMBEDDED_APPLICATION = 2
)

// Invite Metadata Object
// https://discord.com/developers/docs/resources/invite#invite-metadata-object-invite-metadata-structure
type InviteMetadata struct {
	Uses      *int      `json:"uses"`
	MaxUses   *int      `json:"max_uses"`
	MaxAge    int       `json:"max_age"`
	Temporary bool      `json:"temporary"`
	CreatedAt time.Time `json:"created_at"`
}

// Stage Instance Object
// https://discord.com/developers/docs/resources/stage-instance#stage-instance-object
type StageInstance struct {
	ID                    string  `json:"id"`
	GuildID               *string `json:"guild_id"`
	ChannelID             *string `json:"channel_id"`
	Topic                 string  `json:"topic"`
	PrivacyLevel          Flag    `json:"privacy_level"`
	DiscoverableDisabled  bool    `json:"discoverable_disabled"`
	GuildScheduledEventID *string `json:"guild_scheduled_event_id"`
}

// Privacy Level
// https://discord.com/developers/docs/resources/stage-instance#stage-instance-object-privacy-level
const (
	FlagPrivacyLevelPUBLIC     = 1
	FlagPrivacyLevelGUILD_ONLY = 2
)

// Sticker Structure
// https://discord.com/developers/docs/resources/sticker#sticker-object-sticker-structure
type Sticker struct {
	ID          string  `json:"id"`
	PackID      string  `json:"pack_id,omitempty"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Tags        *string `json:"tags"`
	Asset       *string `json:"asset"`
	Type        Flag    `json:"type"`
	FormatType  Flag    `json:"format_type"`
	Available   *bool   `json:"available,omitempty"`
	GuildID     *string `json:"guild_id,omitempty"`
	User        *User   `json:"user,omitempty"`
	SortValue   *int    `json:"sort_value,omitempty"`
}

// Sticker Types
// https://discord.com/developers/docs/resources/sticker#sticker-object-sticker-types
const (
	FlagStickerTypeSTANDARD = 1
	FlagStickerTypeGUILD    = 2
)

// Sticker Format Types
// https://discord.com/developers/docs/resources/sticker#sticker-object-sticker-format-types
const (
	FlagStickerFormatTypePNG    = 1
	FlagStickerFormatTypeAPNG   = 2
	FlagStickerFormatTypeLOTTIE = 3
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
	ID             string     `json:"id"`
	Stickers       []*Sticker `json:"stickers"`
	Name           string     `json:"name"`
	SKU_ID         string     `json:"sku_id"`
	CoverStickerID string     `json:"cover_sticker_id,omitempty"`
	Description    string     `json:"description"`
	BannerAssetID  string     `json:"banner_asset_id,omitempty"`
}

// User Object
// https://discord.com/developers/docs/resources/user#user-object
type User struct {
	ID            string   `json:"id"`
	Username      string   `json:"username"`
	Discriminator string   `json:"discriminator"`
	Avatar        *string  `json:"avatar"`
	Bot           *bool    `json:"bot,omitempty"`
	System        *bool    `json:"system,omitempty"`
	MFAEnabled    *bool    `json:"mfa_enabled,omitempty"`
	Banner        *string  `json:"banner"`
	AccentColor   *int     `json:"accent_color"`
	Locale        string   `json:"locale,omitempty"`
	Verified      bool     `json:"verified,omitempty"`
	Email         *string  `json:"email"`
	Flags         *BitFlag `json:"flag,omitempty"`
	PremiumType   *Flag    `json:"premium_type,omitempty"`
	PublicFlags   BitFlag  `json:"public_flag,omitempty"`
}

// User Flags
// https://discord.com/developers/docs/resources/user#user-object-user-flags
const (
	FlagUserNONE                         = 0
	FlagUserSTAFF                        = 1 << 0
	FlagUserPARTNER                      = 1 << 1
	FlagUserHYPESQUAD                    = 1 << 2
	FlagUserBUG_HUNTER_LEVEL_1           = 1 << 3
	FlagUserHYPESQUAD_ONLINE_HOUSE_ONE   = 1 << 6
	FlagUserHYPESQUAD_ONLINE_HOUSE_TWO   = 1 << 7
	FlagUserHYPESQUAD_ONLINE_HOUSE_THREE = 1 << 8
	FlagUserPREMIUM_EARLY_SUPPORTER      = 1 << 9
	FlagUserTEAM_PSEUDO_USER             = 1 << 10
	FlagUserBUG_HUNTER_LEVEL_2           = 1 << 14
	FlagUserVERIFIED_BOT                 = 1 << 16
	FlagUserVERIFIED_DEVELOPER           = 1 << 17
	FlagUserCERTIFIED_MODERATOR          = 1 << 18
	FlagUserBOT_HTTP_INTERACTIONS        = 1 << 19
)

// Premium Types
// https://discord.com/developers/docs/resources/user#user-object-premium-types
const (
	FlagPremiumTypeNONE         = 0
	FlagPremiumTypeNITROCLASSIC = 1
	FlagPremiumTypeNITRO        = 2
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
	Visibility   Flag           `json:"visibility"`
}

// Visibility Types
// https://discord.com/developers/docs/resources/user#connection-object-visibility-types
const (
	FlagVisibilityTypeNONE     = 0
	FlagVisibilityTypeEVERYONE = 1
)

// Voice State Object
// https://discord.com/developers/docs/resources/voice#voice-state-object-voice-state-structure
type VoiceState struct {
	GuildID                 string       `json:"guild_id,omitempty"`
	ChannelID               *string      `json:"channel_id"`
	UserID                  string       `json:"user_id"`
	Member                  *GuildMember `json:"member,omitempty"`
	SessionID               string       `json:"session_id"`
	Deaf                    bool         `json:"deaf"`
	Mute                    bool         `json:"mute"`
	SelfDeaf                bool         `json:"self_deaf"`
	SelfMute                bool         `json:"self_mute"`
	SelfStream              *bool        `json:"self_stream,omitempty"`
	SelfVideo               bool         `json:"self_video"`
	Suppress                bool         `json:"suppress"`
	RequestToSpeakTimestamp *time.Time   `json:"request_to_speak_timestamp"`
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
	ID            string   `json:"id"`
	Type          Flag     `json:"type"`
	GuildID       *string  `json:"guild_id"`
	ChannelID     *string  `json:"channel_id"`
	User          *User    `json:"user,omitempty"`
	Name          *string  `json:"name"`
	Avatar        *string  `json:"avatar"`
	Token         string   `json:"token,omitempty"`
	ApplicationID *string  `json:"application_id"`
	SourceGuild   *Guild   `json:"source_guild,omitempty"`
	SourceChannel *Channel `json:"source_channel,omitempty"`
	URL           *string  `json:"url,omitempty"`
}

// Webhook Types
// https://discord.com/developers/docs/resources/webhook#webhook-object-webhook-types
const (
	FlagWebhookTypeINCOMING        = 1
	FlagWebhookTypeCHANNELFOLLOWER = 2
	FlagWebhookTypeAPPLICATION     = 3
)

// Bitwise Permission Flags
// https://discord.com/developers/docs/topics/permissions#permissions-bitwise-permission-flags
const (
	FlagBitwisePermissionCREATE_INSTANT_INVITE      = 1 << 0
	FlagBitwisePermissionKICK_MEMBERS               = 1 << 1
	FlagBitwisePermissionBAN_MEMBERS                = 1 << 2
	FlagBitwisePermissionADMINISTRATOR              = 1 << 3
	FlagBitwisePermissionMANAGE_CHANNELS            = 1 << 4
	FlagBitwisePermissionMANAGE_GUILD               = 1 << 5
	FlagBitwisePermissionADD_REACTIONS              = 1 << 6
	FlagBitwisePermissionVIEW_AUDIT_LOG             = 1 << 7
	FlagBitwisePermissionPRIORITY_SPEAKER           = 1 << 8
	FlagBitwisePermissionSTREAM                     = 1 << 9
	FlagBitwisePermissionVIEW_CHANNEL               = 1 << 10
	FlagBitwisePermissionSEND_MESSAGES              = 1 << 11
	FlagBitwisePermissionSEND_TTS_MESSAGES          = 1 << 12
	FlagBitwisePermissionMANAGE_MESSAGES            = 1 << 13
	FlagBitwisePermissionEMBED_LINKS                = 1 << 14
	FlagBitwisePermissionATTACH_FILES               = 1 << 15
	FlagBitwisePermissionREAD_MESSAGE_HISTORY       = 1 << 16
	FlagBitwisePermissionMENTION_EVERYONE           = 1 << 17
	FlagBitwisePermissionUSE_EXTERNAL_EMOJIS        = 1 << 18
	FlagBitwisePermissionVIEW_GUILD_INSIGHTS        = 1 << 19
	FlagBitwisePermissionCONNECT                    = 1 << 20
	FlagBitwisePermissionSPEAK                      = 1 << 21
	FlagBitwisePermissionMUTE_MEMBERS               = 1 << 22
	FlagBitwisePermissionDEAFEN_MEMBERS             = 1 << 23
	FlagBitwisePermissionMOVE_MEMBERS               = 1 << 24
	FlagBitwisePermissionUSE_VAD                    = 1 << 25
	FlagBitwisePermissionCHANGE_NICKNAME            = 1 << 26
	FlagBitwisePermissionMANAGE_NICKNAMES           = 1 << 27
	FlagBitwisePermissionMANAGE_ROLES               = 1 << 28
	FlagBitwisePermissionMANAGE_WEBHOOKS            = 1 << 29
	FlagBitwisePermissionMANAGE_EMOJIS_AND_STICKERS = 1 << 30
	FlagBitwisePermissionUSE_APPLICATION_COMMANDS   = 1 << 31
	FlagBitwisePermissionREQUEST_TO_SPEAK           = 1 << 32
	FlagBitwisePermissionMANAGE_EVENTS              = 1 << 33
	FlagBitwisePermissionMANAGE_THREADS             = 1 << 34
	FlagBitwisePermissionCREATE_PUBLIC_THREADS      = 1 << 35
	FlagBitwisePermissionCREATE_PRIVATE_THREADS     = 1 << 36
	FlagBitwisePermissionUSE_EXTERNAL_STICKERS      = 1 << 37
	FlagBitwisePermissionSEND_MESSAGES_IN_THREADS   = 1 << 38
	FlagBitwisePermissionUSE_EMBEDDED_ACTIVITIES    = 1 << 39
	FlagBitwisePermissionMODERATE_MEMBERS           = 1 << 40
)

// Permission Overwrite Types
const (
	FlagPermissionOverwriteTypeRole   = 0
	FlagPermissionOverwriteTypeMember = 1
)

// Role Object
// https://discord.com/developers/docs/topics/permissions#role-object
type Role struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Color        int       `json:"color"`
	Hoist        bool      `json:"hoist"`
	Icon         *string   `json:"icon,omitempty"`
	UnicodeEmoji *string   `json:"unicode_emoji,omitempty"`
	Position     int       `json:"position"`
	Permissions  string    `json:"permissions"`
	Managed      bool      `json:"managed"`
	Mentionable  bool      `json:"mentionable"`
	Tags         *RoleTags `json:"tags,omitempty"`
}

// Role Tags Structure
// https://discord.com/developers/docs/topics/permissions#role-object-role-tags-structure
type RoleTags struct {
	BotID             string `json:"bot_id,omitempty"`
	IntegrationID     string `json:"integration_id,omitempty"`
	PremiumSubscriber bool   `json:"premium_subscriber,omitempty"`
}

// Team Object
// https://discord.com/developers/docs/topics/teams#data-models-team-object
type Team struct {
	Icon        *string       `json:"icon"`
	ID          string        `json:"id"`
	Members     []*TeamMember `json:"members"`
	Name        string        `json:"name"`
	Description *string       `json:"description"`
	OwnerUserID string        `json:"owner_user_id"`
}

// Team Member Object
// https://discord.com/developers/docs/topics/teams#data-models-team-member-object
type TeamMember struct {
	MembershipState Flag     `json:"membership_state"`
	Permissions     []string `json:"permissions"`
	TeamID          string   `json:"team_id"`
	User            *User    `json:"user"`
}

// Membership State Enum
// https://discord.com/developers/docs/topics/teams#data-models-membership-state-enum
const (
	FlagMembershipStateEnumINVITED  = 1
	FlagMembershipStateEnumACCEPTED = 2
)

// Client Status Object
// https://discord.com/developers/docs/topics/gateway#client-status-object
type ClientStatus struct {
	Desktop *string `json:"desktop,omitempty"`
	Mobile  *string `json:"mobile,omitempty"`
	Web     *string `json:"web,omitempty"`
}

// Activity Object
// https://discord.com/developers/docs/topics/gateway#activity-object-activity-structure
type Activity struct {
	Name          string              `json:"name"`
	Type          *Flag               `json:"type"`
	URL           *string             `json:"url"`
	CreatedAt     int                 `json:"created_at"`
	Timestamps    *ActivityTimestamps `json:"timestamps,omitempty"`
	ApplicationID string              `json:"application_id,omitempty"`
	Details       *string             `json:"details"`
	State         *string             `json:"state"`
	Emoji         *Emoji              `json:"emoji"`
	Party         *ActivityParty      `json:"party,omitempty"`
	Assets        *ActivityAssets     `json:"assets,omitempty"`
	Secrets       *ActivitySecrets    `json:"secrets,omitempty"`
	Instance      *bool               `json:"instance,omitempty"`
	Flags         BitFlag             `json:"flags,omitempty"`
	Buttons       []Button            `json:"buttons,omitempty"`
}

// Activity Types
// https://discord.com/developers/docs/topics/gateway#activity-object-activity-types
const (
	FlagActivityTypePlaying   = 0
	FlagActivityTypeStreaming = 1
	FlagActivityTypeListening = 2
	FlagActivityTypeWatching  = 3
	FlagActivityTypeCustom    = 4
	FlagActivityTypeCompeting = 5
)

// Activity Timestamps Struct
// htthttps://discord.com/developers/docs/topics/gateway#activity-object-activity-timestamps
type ActivityTimestamps struct {
	Start *int `json:"start,omitempty"`
	End   *int `json:"end,omitempty"`
}

// Activity Emoji
// https://discord.com/developers/docs/topics/gateway#activity-object-activity-emoji
type ActivityEmoji struct {
	Name     string `json:"name"`
	ID       string `json:"id,omitempty"`
	Animated *bool  `json:"animated,omitempty"`
}

// Activity Party Struct
// https://discord.com/developers/docs/topics/gateway#activity-object-activity-party
type ActivityParty struct {
	ID   *string `json:"id,omitempty"`
	Size *[2]int `json:"size,omitempty"`
}

// Activity Assets Struct
// https://discord.com/developers/docs/topics/gateway#activity-object-activity-assets
type ActivityAssets struct {
	LargeImage *string `json:"large_image,omitempty"`
	LargeText  *string `json:"large_text,omitempty"`
	SmallImage *string `json:"small_image,omitempty"`
	SmallText  *string `json:"small_text,omitempty"`
}

// Activity Asset Image
// https://discord.com/developers/docs/topics/gateway#activity-object-activity-asset-image
type ActivityAssetImage struct {
	ApplicationAsset string `json:"application_asset_id"`
	MediaProxyImage  string `json:"image_id"`
}

// Activity Secrets Struct
// https://discord.com/developers/docs/topics/gateway#activity-object-activity-secrets
type ActivitySecrets struct {
	Join     *string `json:"join,omitempty"`
	Spectate *string `json:"spectate,omitempty"`
	Match    *string `json:"match,omitempty"`
}

// Activity Flags
// https://discord.com/developers/docs/topics/gateway#activity-object-activity-flags
const (
	FlagActivityINSTANCE                    = 1 << 0
	FlagActivityJOIN                        = 1 << 1
	FlagActivitySPECTATE                    = 1 << 2
	FlagActivityJOIN_REQUEST                = 1 << 3
	FlagActivitySYNC                        = 1 << 4
	FlagActivityPLAY                        = 1 << 5
	FlagActivityPARTY_PRIVACY_FRIENDS       = 1 << 6
	FlagActivityPARTY_PRIVACY_VOICE_CHANNEL = 1 << 7
	FlagActivityEMBEDDED                    = 1 << 8
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

// Current Authorization Information Response Structure
// https://discord.com/developers/docs/topics/oauth2#get-current-authorization-information
type CurrentAuthorizationInformationResponse struct {
	Application *Application `json:"application"`
	Scopes      []*int       `json:"scopes"`
	Expires     *time.Time   `json:"expires"`
	User        *User        `json:"user,omitempty"`
}

// Get Gateway Response
// https://discord.com/developers/docs/topics/gateway#get-gateway-example-response
type GetGatewayResponse struct {
	URL string `json:"url,omitempty"`
}

// Get Gateway Bot Response
// https://discord.com/developers/docs/topics/gateway#get-gateway-example-response
type GetGatewayBotResponse struct {
	URL               string            `json:"url"`
	Shards            *int              `json:"shards"`
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
	ExpiresIn    time.Duration `json:"expires_in,omitempty"`
	RefreshToken string        `json:"refresh_token,omitempty"`
	Scope        string        `json:"scope,omitempty"`
}

// Redirect URI
// https://discord.com/developers/docs/topics/oauth2#implicit-grant-redirect-url-example
type RedirectURI struct {
	AccessToken string        `url:"access_token,omitempty"`
	TokenType   string        `url:"token_type,omitempty"`
	ExpiresIn   time.Duration `url:"expires_in,omitempty"`
	Scope       string        `url:"scope,omitempty"`
	State       string        `url:"state,omitempty"`
}

// Client Credentials Access Token Response
// https://discord.com/developers/docs/topics/oauth2#client-credentials-grant-client-credentials-access-token-response
type ClientCredentialsAccessTokenResponse struct {
	AccessToken string        `json:"access_token,omitempty"`
	TokenType   string        `json:"token_type,omitempty"`
	ExpiresIn   time.Duration `json:"expires_in,omitempty"`
	Scope       string        `json:"scope,omitempty"`
}

// Webhook Token Response
// https://discord.com/developers/docs/topics/oauth2#webhooks-webhook-token-response-example
type WebhookTokenResponse struct {
	TokenType    string        `json:"token_type,omitempty"`
	AccessToken  string        `json:"access_token,omitempty"`
	Scope        string        `json:"scope,omitempty"`
	ExpiresIn    time.Duration `json:"expires_in,omitempty"`
	RefreshToken string        `json:"refresh_token,omitempty"`
	Webhook      *Webhook      `json:"webhook,omitempty"`
}

// Extended Bot Authorization Access Token Response
// https://discord.com/developers/docs/topics/oauth2#authorization-code-grant-access-token-response
type ExtendedBotAuthorizationAccessTokenResponse struct {
	TokenType    string        `json:"token_type,omitempty"`
	Guild        *Guild        `json:"guild,omitempty"`
	AccessToken  string        `json:"access_token,omitempty"`
	Scope        string        `json:"scope,omitempty"`
	ExpiresIn    time.Duration `json:"expires_in,omitempty"`
	RefreshToken string        `json:"refresh_token,omitempty"`
}

// Pointer returns a pointer to the given value.
func Pointer[T any](v T) *T {
	return &v
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
