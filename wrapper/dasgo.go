package disgo

import (
	"encoding/json"
	"time"
)

// Gateway Opcodes
// https://discord.com/developers/docs/topics/opcodes-and-status-codes#gateway-gateway-opcodes
const (
	FlagOpcodesGatewayDispatch  = 0
	FlagOpcodesGatewayHeartbeat = 1
	FlagOpcodesGatewayIdentify  = 2
	FlagOpcodesGatewayPresence  = 3
	FlagOpcodesGatewayVoice     = 4
	FlagOpcodesGatewayResume    = 6
	FlagOpcodesGatewayReconnect = 7
	FlagOpcodesGatewayRequest   = 8
	FlagOpcodesGatewayInvalid   = 9
	FlagOpcodesGatewayHello     = 10
	FlagOpcodesGatewayAck       = 11
)

// Gateway Close Event Codes
// https://discord.com/developers/docs/topics/opcodes-and-status-codes#gateway-gateway-close-event-codes
const (
	FlagCodesEventCloseGatewayUnknownError         = 4000
	FlagCodesEventCloseGatewayUnknownOpcode        = 4001
	FlagCodesEventCloseGatewayDecodeError          = 4002
	FlagCodesEventCloseGatewayNotAuthenticated     = 4003
	FlagCodesEventCloseGatewayAuthenticationFailed = 4004
	FlagCodesEventCloseGatewayAlreadyAuthenticated = 4005
	FlagCodesEventCloseGatewayInvalidSeq           = 4007
	FlagCodesEventCloseGatewayRateLimited          = 4008
	FlagCodesEventCloseGatewaySessionTimed         = 4009
	FlagCodesEventCloseGatewayInvalidShard         = 4010
	FlagCodesEventCloseGatewayShardingRequired     = 4011
	FlagCodesEventCloseGatewayInvalidAPIVersion    = 4012
	FlagCodesEventCloseGatewayInvalidIntent        = 4013
	FlagCodesEventCloseGatewayDisallowedIntent     = 4014
)

// Voice Opcodes
// https://discord.com/developers/docs/topics/opcodes-and-status-codes#voice-voice-opcodes
const (
	FlagOpcodesVoiceIdentify           = 0
	FlagOpcodesVoiceSelectProtocol     = 1
	FlagOpcodesVoiceReadyServer        = 2
	FlagOpcodesVoiceHeartbeat          = 3
	FlagOpcodesVoiceSessionDescription = 4
	FlagOpcodesVoiceSpeaking           = 5
	FlagOpcodesVoiceHeartbeatACK       = 6
	FlagOpcodesVoiceResume             = 7
	FlagOpcodesVoiceHello              = 8
	FlagOpcodesVoiceResumed            = 9
	FlagOpcodesVoiceClientDisconnect   = 13
)

// Voice Close Event Codes
// https://discord.com/developers/docs/topics/opcodes-and-status-codes#voice-voice-close-event-codes
const (
	FlagCodesEventCloseVoiceUnknownOpcode         = 4001
	FlagCodesEventCloseVoiceFailedDecode          = 4002
	FlagCodesEventCloseVoiceNotAuthenticated      = 4003
	FlagCodesEventCloseVoiceAuthenticationFailed  = 4004
	FlagCodesEventCloseVoiceAlreadyAuthenticated  = 4005
	FlagCodesEventCloseVoiceInvalidSession        = 4006
	FlagCodesEventCloseVoiceSessionTimeout        = 4009
	FlagCodesEventCloseVoiceServerNotFound        = 4011
	FlagCodesEventCloseVoiceUnknownProtocol       = 4012
	FlagCodesEventCloseVoiceDisconnectedChannel   = 4014
	FlagCodesEventCloseVoiceVoiceServerCrash      = 4015
	FlagCodesEventCloseVoiceUnknownEncryptionMode = 4016
)

// HTTP Response Codes
// https://discord.com/developers/docs/topics/opcodes-and-status-codes#http-http-response-codes
const (
	FlagCodesResponseHTTPOK                 = 200
	FlagCodesResponseHTTPCREATED            = 201
	FlagCodesResponseHTTPNOCONTENT          = 204
	FlagCodesResponseHTTPNOTMODIFIED        = 304
	FlagCodesResponseHTTPBADREQUEST         = 400
	FlagCodesResponseHTTPUNAUTHORIZED       = 401
	FlagCodesResponseHTTPFORBIDDEN          = 403
	FlagCodesResponseHTTPNOTFOUND           = 404
	FlagCodesResponseHTTPMETHODNOTALLOWED   = 405
	FlagCodesResponseHTTPTOOMANYREQUESTS    = 429
	FlagCodesResponseHTTPGATEWAYUNAVAILABLE = 502
	FlagCodesResponseHTTPSERVERERROR        = 504 // 5xx (504 Not Guaranteed)
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
		20001:  "Bots cannot use this endpoint",
		20002:  "Only bots can use this endpoint",
		20009:  "Explicit content cannot be sent to the desired recipient(s)",
		20012:  "You are not authorized to perform this action on this application",
		20016:  "This action cannot be performed due to slowmode rate limit",
		20018:  "Only the owner of this account can perform this action",
		20022:  "This message cannot be edited due to announcement rate limits",
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
		30033:  "Max number of thread participants has been reached (1000)",
		30035:  "Maximum number of bans for non-guild members have been exceeded",
		30037:  "Maximum number of bans fetches has been reached",
		30038:  "Maximum number of uncompleted guild scheduled events reached (100)",
		30039:  "Maximum number of stickers reached",
		30040:  "Maximum number of prune requests has been reached. Try again later",
		30042:  "Maximum number of guild widget settings updates has been reached. Try again later",
		30046:  "Maximum number of edits to messages older than 1 hour reached. Try again later",
		40001:  "Unauthorized. Provide a valid token and try again",
		40002:  "You need to verify your account in order to perform this action",
		40003:  "You are opening direct messages too fast",
		40004:  "Send messages has been temporarily disabled",
		40005:  "Request entity too large. Try sending something smaller in size",
		40006:  "This feature has been temporarily disabled server-side",
		40007:  "The user is banned from this guild",
		40032:  "Target user is not connected to voice",
		40033:  "This message has already been crossposted",
		40041:  "An application command with that name already exists",
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
		60003:  "Two factor is required for this operation",
		80004:  "No users with DiscordTag exist",
		90001:  "Reaction was blocked",
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
	}
)

// RPC Error Codes
// https://discord.com/developers/docs/topics/opcodes-and-status-codes#rpc-rpc-error-codes
const (
	FlagCodesErrorRPCUnknownError                    = 1000
	FlagCodesErrorRPCInvalidPayload                  = 4000
	FlagCodesErrorRPCInvalidCommand                  = 4002
	FlagCodesErrorRPCInvalidGuild                    = 4003
	FlagCodesErrorRPCInvalidEvent                    = 4004
	FlagCodesErrorRPCInvalidChannel                  = 4005
	FlagCodesErrorRPCInvalidPermissions              = 4006
	FlagCodesErrorRPCInvalidClientID                 = 4007
	FlagCodesErrorRPCInvalidOrigin                   = 4008
	FlagCodesErrorRPCInvalidToken                    = 4009
	FlagCodesErrorRPCInvalidUser                     = 4010
	FlagCodesErrorRPCOAuth2Error                     = 5000
	FlagCodesErrorRPCSelectChannelTimedOut           = 5001
	FlagCodesErrorRPCGET_GUILDTimedOut               = 5002
	FlagCodesErrorRPCSelectVoiceForceRequired        = 5003
	FlagCodesErrorRPCCaptureShortcutAlreadyListening = 5004
)

// RPC Close Event Codes
// https://discord.com/developers/docs/topics/opcodes-and-status-codes#rpc-rpc-close-event-codes
const (
	FlagCodesEventCloseRPCInvalidClientID = 4000
	FlagCodesEventCloseRPCInvalidOrigin   = 4001
	FlagCodesEventCloseRPCRateLimited     = 4002
	FlagCodesEventCloseRPCTokenRevoked    = 4003
	FlagCodesEventCloseRPCInvalidVersion  = 4004
	FlagCodesEventCloseRPCInvalidEncoding = 4005
)

// Gateway Commands
// https://discord.com/developers/docs/topics/gateway#commands-and-events-gateway-commands
type Command interface{}

// Identify Structure
// https://discord.com/developers/docs/topics/gateway#identify-identify-structure
type Identify struct {
	Token          string                       `json:"token,omitempty"`
	Properties     IdentifyConnectionProperties `json:"properties,omitempty"`
	Compress       bool                         `json:"compress,omitempty"`
	LargeThreshold int                          `json:"large_threshold,omitempty"`
	Shard          *[2]int                      `json:"shard,omitempty"`
	Presence       GatewayPresenceUpdate        `json:"presence,omitempty"`
	Intents        BitFlag                      `json:"intents,omitempty"`
}

// Identify Connection Properties
// https://discord.com/developers/docs/topics/gateway#identify-identify-connection-properties
type IdentifyConnectionProperties struct {
	OS      string `json:"$os,omitempty"`
	Browser string `json:"$browser,omitempty"`
	Device  string `json:"$device,omitempty"`
}

// Resume Structure
// https://discord.com/developers/docs/topics/gateway#resume-resume-structure
type Resume struct {
	Token     string `json:"token,omitempty"`
	SessionID string `json:"session_id,omitempty"`
	Seq       uint32 `json:"seq,omitempty"`
}

// Heartbeat
// https://discord.com/developers/docs/topics/gateway#heartbeat
type Heartbeat struct {
	Op   int   `json:"op,omitempty"`
	Data int64 `json:"d,omitempty"`
}

// Guild Request Members Structure
// https://discord.com/developers/docs/topics/gateway#request-guild-members-guild-request-members-structure
type GuildRequestMembers struct {
	GuildID   Snowflake   `json:"guild_id,omitempty"`
	Query     string      `json:"query,omitempty"`
	Limit     uint        `json:"limit,omitempty"`
	Presences bool        `json:"presences,omitempty"`
	UserIDs   []Snowflake `json:"user_ids,omitempty"`
	Nonce     string      `json:"nonce,omitempty"`
}

// Gateway Voice State Update Structure
// https://discord.com/developers/docs/topics/gateway#update-voice-state-gateway-voice-state-update-structure
type GatewayVoiceStateUpdate struct {
	GuildID   Snowflake `json:"guild_id,omitempty"`
	ChannelID Snowflake `json:"channel_id,omitempty"`
	SelfMute  bool      `json:"self_mute,omitempty"`
	SelfDeaf  bool      `json:"self_deaf,omitempty"`
}

// Gateway Presence Update Structure
// https://discord.com/developers/docs/topics/gateway#update-presence-gateway-presence-update-structure
type GatewayPresenceUpdate struct {
	Since  int         `json:"since,omitempty"`
	Game   []*Activity `json:"game,omitempty"`
	Status string      `json:"status,omitempty"`
	AFK    bool        `json:"afk,omitempty"`
}

// Status Types
// https://discord.com/developers/docs/topics/gateway#update-presence-status-types
const (
	FlagTypesStatusOnline       = "online"
	FlagTypesStatusDoNotDisturb = "dnd"
	FlagTypesStatusAFK          = "idle"
	FlagTypesStatusInvisible    = "invisible"
	FlagTypesStatusOffline      = "offline"
)

// Snowflake represents a Discord API Snowflake.
type Snowflake uint64

// Flag represents an (unused) alias for a Discord API Flag ranging from 0 - 255.
type Flag uint8

// BitFlag represents an alias for a Discord API Bitwise Flag denoted by 1 << x.
type BitFlag uint

// CodeFlag represents an alias for a Discord API code ranging from 0 - 65535.
type CodeFlag uint16

// Gateway Events
// https://discord.com/developers/docs/topics/gateway#commands-and-events-gateway-events
type Event interface{}

// Hello Structure
// https://discord.com/developers/docs/topics/gateway#hello-hello-structure
type Hello struct {
	HeartbeatInterval time.Duration `json:"heartbeat_interval,omitempty"`
}

// Ready Event Fields
// https://discord.com/developers/docs/topics/gateway#ready-ready-event-fields
type Ready struct {
	Version     int          `json:"v,omitempty"`
	User        *User        `json:"user,omitempty"`
	Guilds      []*Guild     `json:"guilds,omitempty"`
	SessionID   string       `json:"session_id,omitempty"`
	Shard       *[2]int      `json:"shard,omitempty"`
	Application *Application `json:"application,omitempty"`
}

// Resumed
// https://discord.com/developers/docs/topics/gateway#resumed
type Resumed struct {
	Op int `json:"op,omitempty"`
}

// Reconnect
// https://discord.com/developers/docs/topics/gateway#reconnect
type Reconnect struct {
	Op int `json:"op,omitempty"`
}

// Invalid Session
// https://discord.com/developers/docs/topics/gateway#invalid-session
type InvalidSession struct {
	Op   int  `json:"op,omitempty"`
	Data bool `json:"d,omitempty"`
}

// Application Command Permissions Update
// https://discord.com/developers/docs/topics/gateway#application-command-permissions-update
type ApplicationCommandPermissionsUpdate struct {
	*GuildApplicationCommandPermissions
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
	GuildID    Snowflake       `json:"guild_id,omitempty"`
	ChannelIDs []Snowflake     `json:"channel_ids,omitempty"`
	Threads    []*Channel      `json:"threads,omitempty"`
	Members    []*ThreadMember `json:"members,omitempty"`
}

// Thread Member Update
// https://discord.com/developers/docs/topics/gateway#thread-member-update
type ThreadMemberUpdate struct {
	*ThreadMember
	GuildID Snowflake `json:"guild_id,omitempty"`
}

// Thread Members Update
// https://discord.com/developers/docs/topics/gateway#thread-members-update
type ThreadMembersUpdate struct {
	ID             Snowflake       `json:"id,omitempty"`
	GuildID        Snowflake       `json:"guild_id,omitempty"`
	MemberCount    int             `json:"member_count,omitempty"`
	AddedMembers   []*ThreadMember `json:"added_members,omitempty"`
	RemovedMembers []Snowflake     `json:"removed_member_ids,omitempty"`
}

// Channel Pins Update
// https://discord.com/developers/docs/topics/gateway#channel-pins-update
type ChannelPinsUpdate struct {
	GuildID          Snowflake `json:"guild_id,omitempty"`
	ChannelID        Snowflake `json:"channel_id,omitempty"`
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
	GuildID string `json:"guild_id,omitempty"`
	User    *User  `json:"user,omitempty"`
}

// Guild Ban Remove
// https://discord.com/developers/docs/topics/gateway#guild-ban-remove
type GuildBanRemove struct {
	GuildID Snowflake `json:"guild_id,omitempty"`
	User    *User     `json:"user,omitempty"`
}

// Guild Emojis Update
// https://discord.com/developers/docs/topics/gateway#guild-emojis-update
type GuildEmojisUpdate struct {
	GuildID Snowflake `json:"guild_id,omitempty"`
	Emojis  []*Emoji  `json:"emojis,omitempty"`
}

// Guild Stickers Update
// https://discord.com/developers/docs/topics/gateway#guild-stickers-update
type GuildStickersUpdate struct {
	GuildID  Snowflake  `json:"guild_id,omitempty"`
	Stickers []*Sticker `json:"stickers,omitempty"`
}

// Guild Integrations Update
// https://discord.com/developers/docs/topics/gateway#guild-integrations-update
type GuildIntegrationsUpdate struct {
	GuildID Snowflake `json:"guild_id,omitempty"`
}

// Guild Member Add
// https://discord.com/developers/docs/topics/gateway#guild-member-add
type GuildMemberAdd struct {
	GuildID Snowflake `json:"guild_id,omitempty"`
	*GuildMember
}

// Guild Member Remove
// https://discord.com/developers/docs/topics/gateway#guild-member-remove
type GuildMemberRemove struct {
	GuildID Snowflake `json:"guild_id,omitempty"`
	User    *User     `json:"user,omitempty"`
}

// Guild Member Update
// https://discord.com/developers/docs/topics/gateway#guild-member-update
type GuildMemberUpdate struct {
	*GuildMember
}

// Guild Members Chunk
// https://discord.com/developers/docs/topics/gateway#guild-members-chunk
type GuildMembersChunk struct {
	GuildID    Snowflake         `json:"guild_id,omitempty"`
	Members    []*GuildMember    `json:"members,omitempty"`
	ChunkIndex int               `json:"chunk_index,omitempty"`
	ChunkCount int               `json:"chunk_count,omitempty"`
	Presences  []*PresenceUpdate `json:"presences,omitempty"`
	NotFound   []Snowflake       `json:"not_found,omitempty"`
	Nonce      string            `json:"nonce,omitempty"`
}

// Guild Role Create
// https://discord.com/developers/docs/topics/gateway#guild-role-create
type GuildRoleCreate struct {
	GuildID Snowflake `json:"guild_id,omitempty"`
	Role    *Role     `json:"role,omitempty"`
}

// Guild Role Update
// https://discord.com/developers/docs/topics/gateway#guild-role-update
type GuildRoleUpdate struct {
	GuildID Snowflake `json:"guild_id,omitempty"`
	Role    *Role     `json:"role,omitempty"`
}

// Guild Role Delete
// https://discord.com/developers/docs/topics/gateway#guild-role-delete
type GuildRoleDelete struct {
	GuildID Snowflake `json:"guild_id,omitempty"`
	RoleID  Snowflake `json:"role_id,omitempty"`
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
	GuildScheduledEventID Snowflake `json:"guild_scheduled_event_id,omitempty"`
	UserID                Snowflake `json:"user_id,omitempty"`
	GuildID               Snowflake `json:"guild_id,omitempty"`
}

// Guild Scheduled Event User Remove
// https://discord.com/developers/docs/topics/gateway#guild-scheduled-event-user-remove
type GuildScheduledEventUserRemove struct {
	GuildScheduledEventID Snowflake `json:"guild_scheduled_event_id,omitempty"`
	UserID                Snowflake `json:"user_id,omitempty"`
	GuildID               Snowflake `json:"guild_id,omitempty"`
}

// Integration Create
// https://discord.com/developers/docs/topics/gateway#integration-create
type IntegrationCreate struct {
	*Integration
	GuildID Snowflake `json:"guild_id,omitempty"`
}

// Integration Update
// https://discord.com/developers/docs/topics/gateway#integration-update
type IntegrationUpdate struct {
	*Integration
	GuildID Snowflake `json:"guild_id,omitempty"`
}

// Integration Delete
// https://discord.com/developers/docs/topics/gateway#integration-delete
type IntegrationDelete struct {
	IntegrationID Snowflake `json:"id,omitempty"`
	GuildID       Snowflake `json:"guild_id,omitempty"`
	ApplicationID Snowflake `json:"application_id,omitempty"`
}

// Interaction Create
// https://discord.com/developers/docs/topics/gateway#interaction-create
type InteractionCreate struct {
	*Interaction
}

// Invite Create
// https://discord.com/developers/docs/topics/gateway#invite-create
type InviteCreate struct {
	ChannelID         Snowflake    `json:"channel_id,omitempty"`
	Code              string       `json:"code,omitempty"`
	CreatedAt         time.Time    `json:"created_at,omitempty"`
	GuildID           Snowflake    `json:"guild_id,omitempty"`
	Inviter           *User        `json:"inviter,omitempty"`
	MaxAge            int          `json:"max_age,omitempty"`
	MaxUses           int          `json:"max_uses,omitempty"`
	TargetType        int          `json:"target_user_type,omitempty"`
	TargetUser        *User        `json:"target_user,omitempty"`
	TargetApplication *Application `json:"target_application,omitempty"`
	Temporary         bool         `json:"temporary,omitempty"`
	Uses              int          `json:"uses,omitempty"`
}

// Invite Delete
// https://discord.com/developers/docs/topics/gateway#invite-delete
type InviteDelete struct {
	ChannelID Snowflake `json:"channel_id,omitempty"`
	GuildID   Snowflake `json:"guild_id,omitempty"`
	Code      Snowflake `json:"code,omitempty"`
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
	MessageID Snowflake `json:"id,omitempty"`
	ChannelID Snowflake `json:"channel_id,omitempty"`
	GuildID   Snowflake `json:"guild_id,omitempty"`
}

// Message Delete Bulk
// https://discord.com/developers/docs/topics/gateway#message-delete-bulk
type MessageDeleteBulk struct {
	MessageIDs []Snowflake `json:"ids,omitempty"`
	ChannelID  Snowflake   `json:"channel_id,omitempty"`
	GuildID    Snowflake   `json:"guild_id,omitempty"`
}

// Message Reaction Add
// https://discord.com/developers/docs/topics/gateway#message-reaction-add
type MessageReactionAdd struct {
	UserID    Snowflake    `json:"user_id,omitempty"`
	ChannelID Snowflake    `json:"channel_id,omitempty"`
	MessageID Snowflake    `json:"message_id,omitempty"`
	GuildID   Snowflake    `json:"guild_id,omitempty"`
	Member    *GuildMember `json:"member,omitempty"`
	Emoji     *Emoji       `json:"emoji,omitempty"`
}

// Message Reaction Remove
// https://discord.com/developers/docs/topics/gateway#message-reaction-remove
type MessageReactionRemove struct {
	UserID    Snowflake    `json:"user_id,omitempty"`
	ChannelID Snowflake    `json:"channel_id,omitempty"`
	MessageID Snowflake    `json:"message_id,omitempty"`
	GuildID   Snowflake    `json:"guild_id,omitempty"`
	Member    *GuildMember `json:"member,omitempty"`
	Emoji     *Emoji       `json:"emoji,omitempty"`
}

// Message Reaction Remove All
// https://discord.com/developers/docs/topics/gateway#message-reaction-remove-all
type MessageReactionRemoveAll struct {
	ChannelID Snowflake `json:"channel_id,omitempty"`
	MessageID Snowflake `json:"message_id,omitempty"`
	GuildID   Snowflake `json:"guild_id,omitempty"`
}

// Message Reaction Remove Emoji
// https://discord.com/developers/docs/topics/gateway#message-reaction-remove-emoji
type MessageReactionRemoveEmoji struct {
	ChannelID Snowflake `json:"channel_id,omitempty"`
	GuildID   Snowflake `json:"guild_id,omitempty"`
	MessageID Snowflake `json:"message_id,omitempty"`
	Emoji     *Emoji    `json:"emoji,omitempty"`
}

// Presence Update Event Fields
// https://discord.com/developers/docs/topics/gateway#presence-update-presence-update-event-fields
type PresenceUpdate struct {
	User         *User         `json:"user,omitempty"`
	GuildID      Snowflake     `json:"guild_id,omitempty"`
	Status       string        `json:"status,omitempty"`
	Activities   []*Activity   `json:"activities,omitempty"`
	ClientStatus *ClientStatus `json:"client_status,omitempty"`
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
	ChannelID Snowflake    `json:"channel_id,omitempty"`
	GuildID   Snowflake    `json:"guild_id,omitempty"`
	UserID    Snowflake    `json:"user_id,omitempty"`
	Member    *GuildMember `json:"member,omitempty"`
	Timestamp int          `json:"timestamp,omitempty"`
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
	Token    string    `json:"token,omitempty"`
	GuildID  Snowflake `json:"guild_id,omitempty"`
	Endpoint string    `json:"endpoint,omitempty"`
}

// Webhooks Update
// https://discord.com/developers/docs/topics/gateway#webhooks-update
type WebhooksUpdate struct {
	GuildID   Snowflake `json:"guild_id,omitempty"`
	ChannelID Snowflake `json:"channel_id,omitempty"`
}

// Gateway Payload Structure
// https://discord.com/developers/docs/topics/gateway#payloads-gateway-payload-structure
type GatewayPayload struct {
	Op             *Flag           `json:"op,omitempty"`
	Data           json.RawMessage `json:"d,omitempty"`
	SequenceNumber uint32          `json:"s,omitempty"`
	EventName      string          `json:"t,omitempty"`
}

// Gateway URL Query String Params
// https://discord.com/developers/docs/topics/gateway#connecting-gateway-url-query-string-params
type GatewayURLQueryString struct {
	V        int    `json:"v,omitempty"`
	Encoding string `json:"encoding,omitempty"`
	Compress string `json:"compress,omitempty"`
}

// Session Start Limit Structure
// https://discord.com/developers/docs/topics/gateway#session-start-limit-object-session-start-limit-structure
type SessionStartLimit struct {
	Total          int `json:"total,omitempty"`
	Remaining      int `json:"remaining,omitempty"`
	ResetAfter     int `json:"reset_after,omitempty"`
	MaxConcurrency int `json:"max_concurrency,omitempty"`
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
	FlagIntentsofListGUILDS = 1 << 0

	// GUILD_MEMBER_ADD
	// GUILD_MEMBER_UPDATE
	// GUILD_MEMBER_REMOVE
	// THREAD_MEMBERS_UPDATE *
	FlagIntentsofListGUILD_MEMBERS = 1 << 1

	// GUILD_BAN_ADD
	// GUILD_BAN_REMOVE
	FlagIntentsofListGUILD_BANS = 1 << 2

	// GUILD_EMOJIS_UPDATE
	// GUILD_STICKERS_UPDATE
	FlagIntentsofListGUILD_EMOJIS_AND_STICKERS = 1 << 3

	// GUILD_INTEGRATIONS_UPDATE
	// INTEGRATION_CREATE
	// INTEGRATION_UPDATE
	// INTEGRATION_DELETE
	FlagIntentsofListGUILD_INTEGRATIONS = 1 << 4

	// WEBHOOKS_UPDATE
	FlagIntentsofListGUILD_WEBHOOKS = 1 << 5

	// INVITE_CREATE
	// INVITE_DELETE
	FlagIntentsofListGUILD_INVITES = 1 << 6

	// VOICE_STATE_UPDATE
	FlagIntentsofListGUILD_VOICE_STATES = 1 << 7

	// PRESENCE_UPDATE
	FlagIntentsofListGUILD_PRESENCES = 1 << 8

	// MESSAGE_CREATE
	// MESSAGE_UPDATE
	// MESSAGE_DELETE
	// MESSAGE_DELETE_BULK
	FlagIntentsofListGUILD_MESSAGES = 1 << 9

	// MESSAGE_REACTION_ADD
	// MESSAGE_REACTION_REMOVE
	// MESSAGE_REACTION_REMOVE_ALL
	// MESSAGE_REACTION_REMOVE_EMOJI
	FlagIntentsofListGUILD_MESSAGE_REACTIONS = 1 << 10

	// TYPING_START

	FlagIntentsofListGUILD_MESSAGE_TYPING = 1 << 11

	// MESSAGE_CREATE
	// MESSAGE_UPDATE
	// MESSAGE_DELETE
	// CHANNEL_PINS_UPDATE
	FlagIntentsofListDIRECT_MESSAGES = 1 << 12

	// MESSAGE_REACTION_ADD
	// MESSAGE_REACTION_REMOVE
	// MESSAGE_REACTION_REMOVE_ALL
	// MESSAGE_REACTION_REMOVE_EMOJI
	FlagIntentsofListDIRECT_MESSAGE_REACTIONS = 1 << 13

	// TYPING_START
	FlagIntentsofListDIRECT_MESSAGE_TYPING = 1 << 14

	// GUILD_SCHEDULED_EVENT_CREATE
	// GUILD_SCHEDULED_EVENT_UPDATE
	// GUILD_SCHEDULED_EVENT_DELETE
	// GUILD_SCHEDULED_EVENT_USER_ADD
	// GUILD_SCHEDULED_EVENT_USER_REMOVE
	FlagIntentsofListGUILD_SCHEDULED_EVENTS = 1 << 16
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
// GET/applications/{application.id}/commands
// https://discord.com/developers/docs/interactions/application-commands#get-global-application-commands
type GetGlobalApplicationCommands struct {
	WithLocalizations bool `json:"with_localizations,omitempty"`
}

// Create Global Application Command
// POST/applications/{application.id}/commands
// https://discord.com/developers/docs/interactions/application-commands#create-global-application-command
type CreateGlobalApplicationCommand struct {
	Name                     string                      `json:"name,omitempty"`
	NameLocalizations        map[Flag]string             `json:"name_localizations,omitempty"`
	Description              string                      `json:"description,omitempty"`
	DescriptionLocalizations map[Flag]string             `json:"description_localizations,omitempty"`
	Options                  []*ApplicationCommandOption `json:"options,omitempty"`
	DefaultPermission        bool                        `json:"default_permission,omitempty"`
	Type                     Flag                        `json:"type,omitempty"`
}

// Get Global Application Command
// GET/applications/{application.id}/commands/{command.id}
// https://discord.com/developers/docs/interactions/application-commands#get-global-application-command
type GetGlobalApplicationCommand struct {
	CommandID Snowflake
}

// Edit Global Application Command
// PATCH/applications/{application.id}/commands/{command.id}
// https://discord.com/developers/docs/interactions/application-commands#edit-global-application-command
type EditGlobalApplicationCommand struct {
	CommandID                Snowflake
	Name                     string                      `json:"name,omitempty"`
	NameLocalizations        map[Flag]string             `json:"name_localizations,omitempty"`
	Description              string                      `json:"description,omitempty"`
	DescriptionLocalizations map[Flag]string             `json:"description_localizations,omitempty"`
	Options                  []*ApplicationCommandOption `json:"options,omitempty"`
	DefaultPermission        bool                        `json:"default_permission,omitempty"`
}

// Delete Global Application Command
// DELETE /applications/{application.id}/commands/{command.id}
// https://discord.com/developers/docs/interactions/application-commands#delete-global-application-command
type DeleteGlobalApplicationCommand struct {
	CommandID Snowflake
}

// Bulk Overwrite Global Application Commands
// PUT /applications/{application.id}/commands
// https://discord.com/developers/docs/interactions/application-commands#bulk-overwrite-global-application-commands
type BulkOverwriteGlobalApplicationCommands struct {
	ApplicationCommands []*ApplicationCommand
}

// Get Guild Application Commands
// GET /applications/{application.id}/guilds/{guild.id}/commands
// https://discord.com/developers/docs/interactions/application-commands#get-guild-application-commands
type GetGuildApplicationCommands struct {
	GuildID           Snowflake
	WithLocalizations bool `json:"with_localizations,omitempty"`
}

// Create Guild Application Command
// POST /applications/{application.id}/guilds/{guild.id}/commands
// https://discord.com/developers/docs/interactions/application-commands#create-guild-application-command
type CreateGuildApplicationCommand struct {
	GuildID                  Snowflake
	Name                     string                      `json:"name,omitempty"`
	NameLocalizations        map[Flag]string             `json:"name_localizations,omitempty"`
	Description              string                      `json:"description,omitempty"`
	DescriptionLocalizations map[Flag]string             `json:"description_localizations,omitempty"`
	Options                  []*ApplicationCommandOption `json:"options,omitempty"`
	DefaultPermission        bool                        `json:"default_permission,omitempty"`
	Type                     Flag                        `json:"type,omitempty"`
}

// Get Guild Application Command
// GET /applications/{application.id}/guilds/{guild.id}/commands/{command.id}
// https://discord.com/developers/docs/interactions/application-commands#get-guild-application-command
type GetGuildApplicationCommand struct {
	GuildID Snowflake
}

// Edit Guild Application Command
// PATCH /applications/{application.id}/guilds/{guild.id}/commands/{command.id}
// https://discord.com/developers/docs/interactions/application-commands#edit-guild-application-command
type EditGuildApplicationCommand struct {
	Name                     string                      `json:"name,omitempty"`
	NameLocalizations        map[Flag]string             `json:"name_localizations,omitempty"`
	Description              string                      `json:"description,omitempty"`
	DescriptionLocalizations map[Flag]string             `json:"description_localizations,omitempty"`
	Options                  []*ApplicationCommandOption `json:"options,omitempty"`
	DefaultPermission        bool                        `json:"default_permission,omitempty"`
}

// Delete Guild Application Command
// DELETE /applications/{application.id}/guilds/{guild.id}/commands/{command.id}
// https://discord.com/developers/docs/interactions/application-commands#delete-guild-application-command
type DeleteGuildApplicationCommand struct {
	GuildID Snowflake
}

// Bulk Overwrite Guild Application Commands
// PUT /applications/{application.id}/guilds/{guild.id}/commands
// https://discord.com/developers/docs/interactions/application-commands#bulk-overwrite-guild-application-commands
type BulkOverwriteGuildApplicationCommands struct {
	CommandID                Snowflake
	Name                     string                      `json:"name,omitempty"`
	NameLocalizations        map[Flag]string             `json:"name_localizations,omitempty"`
	Description              string                      `json:"description,omitempty"`
	DescriptionLocalizations map[Flag]string             `json:"description_localizations,omitempty"`
	Options                  []*ApplicationCommandOption `json:"options,omitempty"`
	DefaultPermission        bool                        `json:"default_permission,omitempty"`
	Type                     Flag                        `json:"type,omitempty"`
}

// Get Guild Application Command Permissions
// GET /applications/{application.id}/guilds/{guild.id}/commands/permissions
// https://discord.com/developers/docs/interactions/application-commands#get-guild-application-command-permissions
type GetGuildApplicationCommandPermissions struct {
	GuildID Snowflake
}

// Get Application Command Permissions
// GET /applications/{application.id}/guilds/{guild.id}/commands/{command.id}/permissions
// https://discord.com/developers/docs/interactions/application-commands#get-application-command-permissions
type GetApplicationCommandPermissions struct {
	GuildID Snowflake
}

// Edit Application Command Permissions
// PUT /applications/{application.id}/guilds/{guild.id}/commands/{command.id}/permissions
// https://discord.com/developers/docs/interactions/application-commands#edit-application-command-permissions
type EditApplicationCommandPermissions struct {
	Permissions []*ApplicationCommandPermissions `json:"permissions,omitempty"`
}

// Batch Edit Application Command Permissions
// PUT /applications/{application.id}/guilds/{guild.id}/commands/permissions
// https://discord.com/developers/docs/interactions/application-commands#batch-edit-application-command-permissions
type BatchEditApplicationCommandPermissions struct {
	GuildID Snowflake
}

// Get Guild Audit Log
// GET /guilds/{guild.id}/audit-logs
// https://discord.com/developers/docs/resources/audit-log#get-guild-audit-log
type GetGuildAuditLog struct {
	UserID     Snowflake `json:"user_id"`
	ActionType Flag      `json:"action_type"`
	Before     Snowflake `json:"before,omitempty"`
	Limit      Flag      `json:"limit,omitempty"`
}

// Get Channel
// GET /channels/{channel.id}
// https://discord.com/developers/docs/resources/channel#get-channel
type GetChannel struct {
	ChannelID Snowflake
}

// Modify Channel
// PATCH /channels/{channel.id}
// https://discord.com/developers/docs/resources/channel#modify-channel
type ModifyChannel struct {
	ChannelID Snowflake
}

// Modify Channel
// PATCH /channels/{channel.id}
// https://discord.com/developers/docs/resources/channel#modify-channel-json-params-group-dm
type ModifyChannelGroupDM struct {
	Name string `json:"name,omitempty"`
	Icon int    `json:"icon,omitempty"`
}

// Modify Channel
// PATCH /channels/{channel.id}
// https://discord.com/developers/docs/resources/channel#modify-channel-json-params-guild-channel
type ModifyChannelGuild struct {
	Name                       *string                `json:"name,omitempty"`
	Type                       *Flag                  `json:"type,omitempty"`
	Position                   *uint                  `json:"position,omitempty"`
	Topic                      *string                `json:"topic,omitempty"`
	NSFW                       bool                   `json:"nsfw,omitempty"`
	RateLimitPerUser           *CodeFlag              `json:"rate_limit_per_user,omitempty"`
	Bitrate                    *uint                  `json:"bitrate,omitempty"`
	UserLimit                  *Flag                  `json:"user_limit,omitempty"`
	PermissionOverwrites       *[]PermissionOverwrite `json:"permission_overwrites,omitempty"`
	ParentID                   *Snowflake             `json:"parent_id,omitempty"`
	RTCRegion                  *string                `json:"rtc_region,omitempty"`
	VideoQualityMode           Flag                   `json:"video_quality_mode,omitempty"`
	DefaultAutoArchiveDuration *uint                  `json:"default_auto_archive_duration,omitempty"`
}

// Modify Channel
// PATCH /channels/{channel.id}
// https://discord.com/developers/docs/resources/channel#modify-channel-json-params-thread
type ModifyChannelThread struct {
	Name                string    `json:"name,omitempty"`
	Archived            bool      `json:"archived,omitempty"`
	AutoArchiveDuration CodeFlag  `json:"auto_archive_duration,omitempty"`
	Locked              bool      `json:"locked,omitempty"`
	Invitable           bool      `json:"invitable,omitempty"`
	RateLimitPerUser    *CodeFlag `json:"rate_limit_per_user,omitempty"`
}

// Delete/Close Channel
// DELETE /channels/{channel.id}
// https://discord.com/developers/docs/resources/channel#deleteclose-channel
type DeleteCloseChannel struct {
	ChannelID Snowflake
}

// Get Channel Messages
// GET /channels/{channel.id}/messages
// https://discord.com/developers/docs/resources/channel#get-channel-messages
type GetChannelMessages struct {
	Around *Snowflake `json:"around,omitempty"`
	Before *Snowflake `json:"before,omitempty"`
	After  *Snowflake `json:"after,omitempty"`
	Limit  Flag       `json:"limit,omitempty"`
}

// Get Channel Message
// GET /channels/{channel.id}/messages/{message.id}
// https://discord.com/developers/docs/resources/channel#get-channel-message
type GetChannelMessage struct {
	MessageID Snowflake
}

// Create Message
// POST /channels/{channel.id}/messages
// https://discord.com/developers/docs/resources/channel#create-message
type CreateMessage struct {
	Content         string            `json:"content,omitempty"`
	TTS             bool              `json:"tts,omitempty"`
	Embeds          []*Embed          `json:"embeds,omitempty"`
	Embed           *Embed            `json:"embed,omitempty"`
	AllowedMentions *AllowedMentions  `json:"allowed_mentions,omitempty"`
	Reference       *MessageReference `json:"message_reference,omitempty"`
	StickerID       []*Snowflake      `json:"sticker_ids,omitempty"`
	Components      []*Component      `json:"components,omitempty"`
	Files           []byte            `dasgo:"files"`
	PayloadJSON     *string           `json:"payload_json,omitempty"`
	Attachments     []*Attachment     `json:"attachments,omitempty"`
	Flags           BitFlag           `json:"flags,omitempty"`
}

// Crosspost Message
// POST /channels/{channel.id}/messages/{message.id}/crosspost
// https://discord.com/developers/docs/resources/channel#crosspost-message
type CrosspostMessage struct {
	MessageID *Snowflake
}

// Create Reaction
// PUT /channels/{channel.id}/messages/{message.id}/reactions/{emoji}/@me
// https://discord.com/developers/docs/resources/channel#create-reaction
type CreateReaction struct {
	MessageID *Snowflake
}

// Delete Own Reaction
// DELETE /channels/{channel.id}/messages/{message.id}/reactions/{emoji}/@me
// https://discord.com/developers/docs/resources/channel#delete-own-reaction
type DeleteOwnReaction struct {
	MessageID *Snowflake
}

// Delete User Reaction
// DELETE /channels/{channel.id}/messages/{message.id}/reactions/{emoji}/{user.id}
// https://discord.com/developers/docs/resources/channel#delete-user-reaction
type DeleteUserReaction struct {
	MessageID *Snowflake
}

// Get Reactions
// GET /channels/{channel.id}/messages/{message.id}/reactions/{emoji}
// https://discord.com/developers/docs/resources/channel#get-reactions
type GetReactions struct {
	After Snowflake `json:"after,omitempty"`
	Limit Flag      `json:"limit,omitempty"` // 1 is default. even if 0 is supplied.
}

// Delete All Reactions
// DELETE /channels/{channel.id}/messages/{message.id}/reactions
// https://discord.com/developers/docs/resources/channel#delete-all-reactions
type DeleteAllReactions struct {
	MessageID *Snowflake
}

// Delete All Reactions for Emoji
// DELETE /channels/{channel.id}/messages/{message.id}/reactions/{emoji}
// https://discord.com/developers/docs/resources/channel#delete-all-reactions-for-emoji
type DeleteAllReactionsforEmoji struct {
	MessageID *Snowflake
}

// Edit Message
// PATCH /channels/{channel.id}/messages/{message.id}
// https://discord.com/developers/docs/resources/channel#edit-message
type EditMessage struct {
	Content         *string          `json:"content,omitempty"`
	Embeds          []*Embed         `json:"embeds,omitempty"`
	Embed           *Embed           `json:"embed,omitempty"`
	Flags           *BitFlag         `json:"flags,omitempty"`
	AllowedMentions *AllowedMentions `json:"allowed_mentions,omitempty"`
	Components      []*Component     `json:"components,omitempty"`
	Files           []byte           `dasgo:"files"`
	PayloadJSON     *string          `json:"payload_json,omitempty"`
	Attachments     []*Attachment    `json:"attachments,omitempty"`
}

// Delete Message
// DELETE /channels/{channel.id}/messages/{message.id}
// https://discord.com/developers/docs/resources/channel#delete-message
type DeleteMessage struct {
	MessageID *Snowflake
}

// Bulk Delete Messages
// POST /channels/{channel.id}/messages/bulk-delete
// https://discord.com/developers/docs/resources/channel#bulk-delete-messages
type BulkDeleteMessages struct {
	Messages []*Snowflake `json:"messages,omitempty"`
}

// Edit Channel Permissions
// PUT /channels/{channel.id}/permissions/{overwrite.id}
// https://discord.com/developers/docs/resources/channel#edit-channel-permissions
type EditChannelPermissions struct {
	Allow string `json:"allow,omitempty"`
	Deny  string `json:"deny,omitempty"`
	Type  *Flag  `json:"type,omitempty"`
}

// Get Channel Invites
// GET /channels/{channel.id}/invites
// https://discord.com/developers/docs/resources/channel#get-channel-invites
type GetChannelInvites struct {
	ChannelID Snowflake
}

// Create Channel Invite
// POST /channels/{channel.id}/invites
// https://discord.com/developers/docs/resources/channel#create-channel-invite
type CreateChannelInvite struct {
	MaxAge              *int      `json:"max_age,omitempty"`
	MaxUses             *Flag     `json:"max_uses,omitempty"`
	Temporary           bool      `json:"temporary,omitempty"`
	Unique              bool      `json:"unique,omitempty"`
	TargetType          Flag      `json:"target_type,omitempty"`
	TargetUserID        Snowflake `json:"target_user_id,omitempty"`
	TargetApplicationID Snowflake `json:"target_application_id,omitempty"`
}

// Delete Channel Permission
// DELETE /channels/{channel.id}/permissions/{overwrite.id}
// https://discord.com/developers/docs/resources/channel#delete-channel-permission
type DeleteChannelPermission struct {
	OverwriteID Snowflake
}

// Follow News Channel
// POST /channels/{channel.id}/followers
// https://discord.com/developers/docs/resources/channel#follow-news-channel
type FollowNewsChannel struct {
	WebhookChannelID Snowflake
}

// Trigger Typing Indicator
// POST /channels/{channel.id}/typing
// https://discord.com/developers/docs/resources/channel#trigger-typing-indicator
type TriggerTypingIndicator struct {
	ChannelID Snowflake
}

// Get Pinned Messages
// GET /channels/{channel.id}/pins
// https://discord.com/developers/docs/resources/channel#get-pinned-messages
type GetPinnedMessages struct {
	ChannelID Snowflake
}

// Pin Message
// PUT /channels/{channel.id}/pins/{message.id}
// https://discord.com/developers/docs/resources/channel#pin-message
type PinMessage struct {
	MessageID Snowflake
}

// Unpin Message
// DELETE /channels/{channel.id}/pins/{message.id}
// https://discord.com/developers/docs/resources/channel#unpin-message
type UnpinMessage struct {
	MessageID Snowflake
}

// Group DM Add Recipient
// PUT /channels/{channel.id}/recipients/{user.id}
// https://discord.com/developers/docs/resources/channel#group-dm-add-recipient
type GroupDMAddRecipient struct {
	AccessToken string  `json:"access_token,omitempty"`
	Nickname    *string `json:"nick,omitempty"`
}

// Group DM Remove Recipient
// DELETE /channels/{channel.id}/recipients/{user.id}
// https://discord.com/developers/docs/resources/channel#group-dm-remove-recipient
type GroupDMRemoveRecipient struct {
	UserID Snowflake
}

// Start Thread from Message
// POST /channels/{channel.id}/messages/{message.id}/threads
// https://discord.com/developers/docs/resources/channel#start-thread-from-message
type StartThreadfromMessage struct {
	Name                string `json:"name,omitempty"`
	RateLimitPerUser    uint   `json:"rate_limit_per_user,omitempty"`
	AutoArchiveDuration *int   `json:"auto_archive_duration,omitempty"`
}

// Start Thread without Message
// POST /channels/{channel.id}/threads
// https://discord.com/developers/docs/resources/channel#start-thread-without-message
type StartThreadwithoutMessage struct {
	Name                string    `json:"name,omitempty"`
	AutoArchiveDuration CodeFlag  `json:"auto_archive_duration,omitempty"`
	Type                *Flag     `json:"type,omitempty"`
	Invitable           bool      `json:"invitable,omitempty"`
	RateLimitPerUser    *CodeFlag `json:"rate_limit_per_user,omitempty"`
}

// Start Thread in Forum Channel
// POST /channels/{channel.id}/threads
// https://discord.com/developers/docs/resources/channel#start-thread-in-forum-channel
type StartThreadinForumChannel struct {
	ChannelID Snowflake
}

// Start Thread in Forum Channel
// POST /channels/{channel.id}/threads
// https://discord.com/developers/docs/resources/channel#start-thread-in-forum-channel-json-params-for-the-thread
type StartThreadinForumChannelThread struct {
	Name                string    `json:"name,omitempty"`
	RateLimitPerUser    uint      `json:"rate_limit_per_user,omitempty"`
	AutoArchiveDuration *CodeFlag `json:"auto_archive_duration,omitempty"`
}

// Start Thread in Forum Channel
// POST /channels/{channel.id}/threads
// https://discord.com/developers/docs/resources/channel#start-thread-in-forum-channel-json-params-for-the-message
type StartThreadinForumChannelMessage struct {
	Content         *string          `json:"content,omitempty"`
	Embeds          []*Embed         `json:"embeds,omitempty"`
	AllowedMentions *AllowedMentions `json:"allowed_mentions,omitempty"`
	Components      []*Component     `json:"components,omitempty"`
	StickerIDS      []*Snowflake     `json:"sticker_ids,omitempty"`
	Attachments     []*Attachment    `json:"attachments,omitempty"`
	Files           []byte           `dasgo:"files"`
	PayloadJSON     string           `json:"payload_json,omitempty"`
	Flags           BitFlag          `json:"flags,omitempty"`
}

// Join Thread
// PUT /channels/{channel.id}/thread-members/@me
// https://discord.com/developers/docs/resources/channel#join-thread
type JoinThread struct {
	ChannelID Snowflake
}

// Add Thread Member
// PUT /channels/{channel.id}/thread-members/{user.id}
// https://discord.com/developers/docs/resources/channel#add-thread-member
type AddThreadMember struct {
	UserID Snowflake
}

// Leave Thread
// DELETE /channels/{channel.id}/thread-members/@me
// https://discord.com/developers/docs/resources/channel#leave-thread
type LeaveThread struct {
	ChannelID Snowflake
}

// Remove Thread Member
// DELETE /channels/{channel.id}/thread-members/{user.id}
// https://discord.com/developers/docs/resources/channel#remove-thread-member
type RemoveThreadMember struct {
	UserID Snowflake
}

// Get Thread Member
// GET /channels/{channel.id}/thread-members/{user.id}
// https://discord.com/developers/docs/resources/channel#get-thread-member
type GetThreadMember struct {
	UserID Snowflake
}

// List Thread Members
// GET /channels/{channel.id}/thread-members
// https://discord.com/developers/docs/resources/channel#list-thread-members
type ListThreadMembers struct {
	ChannelID Snowflake
}

// List Active Channel Threads
// GET /channels/{channel.id}/threads/active
// https://discord.com/developers/docs/resources/channel#list-active-threads
type ListActiveChannelThreads struct {
	Before Snowflake `json:"before,omitempty"`
	Limit  int       `json:"limit,omitempty"`
}

// List Public Archived Threads
// GET /channels/{channel.id}/threads/archived/public
// https://discord.com/developers/docs/resources/channel#list-public-archived-threads
type ListPublicArchivedThreads struct {
	Before Snowflake `json:"before,omitempty"`
	Limit  int       `json:"limit,omitempty"`
}

// List Private Archived Threads
// GET /channels/{channel.id}/threads/archived/private
// https://discord.com/developers/docs/resources/channel#list-private-archived-threads
type ListPrivateArchivedThreads struct {
	Before Snowflake `json:"before,omitempty"`
	Limit  int       `json:"limit,omitempty"`
}

// List Joined Private Archived Threads
// GET /channels/{channel.id}/users/@me/threads/archived/private
// https://discord.com/developers/docs/resources/channel#list-joined-private-archived-threads
type ListJoinedPrivateArchivedThreads struct {
	Before Snowflake `json:"before,omitempty"`
	Limit  int       `json:"limit,omitempty"`
}

// List Guild Emojis
// GET /guilds/{guild.id}/emojis
// https://discord.com/developers/docs/resources/emoji#list-guild-emojis
type ListGuildEmojis struct {
	GuildID Snowflake
}

// Get Guild Emoji
// GET /guilds/{guild.id}/emojis/{emoji.id}
// https://discord.com/developers/docs/resources/emoji#get-guild-emoji
type GetGuildEmoji struct {
	EmojiID Snowflake
}

// Create Guild Emoji
// POST /guilds/{guild.id}/emojis
// https://discord.com/developers/docs/resources/emoji#create-guild-emoji
type CreateGuildEmoji struct {
	GuildID Snowflake
	Name    string       `json:"name,omitempty"`
	Image   string       `json:"image,omitempty"`
	Roles   []*Snowflake `json:"roles,omitempty"`
}

// Modify Guild Emoji
// PATCH /guilds/{guild.id}/emojis/{emoji.id}
// https://discord.com/developers/docs/resources/emoji#modify-guild-emoji
type ModifyGuildEmoji struct {
	EmojiID Snowflake
	Name    string       `json:"name,omitempty"`
	Roles   []*Snowflake `json:"roles,omitempty"`
}

// Delete Guild Emoji
// DELETE /guilds/{guild.id}/emojis/{emoji.id}
// https://discord.com/developers/docs/resources/emoji#delete-guild-emoji
type DeleteGuildEmoji struct {
	EmojiID Snowflake
}

// Get Gateway
// GET /gateway
// https://discord.com/developers/docs/topics/gateway#get-gateway
type GetGateway struct{}

// Get Gateway Bot
// GET /gateway/bot
// https://discord.com/developers/docs/topics/gateway#get-gateway-bot
type GetGatewayBot struct {
	URL               string `json:"url,omitempty"`
	Shards            int    `json:"shards,omitempty"`
	SessionStartLimit int    `json:"session_start_limit,omitempty"`
}

// Create Guild
// POST /guilds
// https://discord.com/developers/docs/resources/guild#create-guild
type CreateGuild struct {
	Name                        string     `json:"name,omitempty"`
	Region                      string     `json:"region,omitempty"`
	Icon                        string     `json:"icon,omitempty"`
	VerificationLevel           *Flag      `json:"verification_level,omitempty"`
	DefaultMessageNotifications *Flag      `json:"default_message_notifications,omitempty"`
	AfkChannelID                string     `json:"afk_channel_id,omitempty"`
	AfkTimeout                  int        `json:"afk_timeout,omitempty"`
	OwnerID                     string     `json:"owner_id,omitempty"`
	Splash                      string     `json:"splash,omitempty"`
	Banner                      string     `json:"banner,omitempty"`
	Roles                       []*Role    `json:"roles,omitempty"`
	Channels                    []*Channel `json:"channels,omitempty"`
	SystemChannelID             Snowflake  `json:"system_channel_id,omitempty"`
	SystemChannelFlags          BitFlag    `json:"system_channel_flags,omitempty"`
	ExplicitContentFilter       *Flag      `json:"explicit_content_filter,omitempty"`
}

// Get Guild
// GET /guilds/{guild.id}
// https://discord.com/developers/docs/resources/guild#get-guild
type GetGuild struct {
	GuildID    *Snowflake
	WithCounts bool `json:"with_counts,omitempty"`
}

// Get Guild Preview
// GET /guilds/{guild.id}/preview
// https://discord.com/developers/docs/resources/guild#get-guild-preview
type GetGuildPreview struct {
	GuildID *Snowflake
}

// Modify Guild
// PATCH /guilds/{guild.id}
// https://discord.com/developers/docs/resources/guild#modify-guild
type ModifyGuild struct {
	GuildID                     Snowflake
	Name                        string     `json:"name,omitempty"`
	Region                      string     `json:"region,omitempty"`
	VerificationLevel           *Flag      `json:"verification_lvl,omitempty"`
	DefaultMessageNotifications *Flag      `json:"default_message_notifications,omitempty"`
	ExplicitContentFilter       *Flag      `json:"explicit_content_filter,omitempty"`
	AFKChannelID                Snowflake  `json:"afk_channel_id,omitempty"`
	Icon                        *string    `json:"icon,omitempty"`
	OwnerID                     Snowflake  `json:"owner_id,omitempty"`
	Splash                      *string    `json:"splash,omitempty"`
	DiscoverySplash             *string    `json:"discovery_splash,omitempty"`
	Banner                      *string    `json:"banner,omitempty"`
	SystemChannelID             Snowflake  `json:"system_channel_id,omitempty"`
	SystemChannelFlags          BitFlag    `json:"system_channel_flags,omitempty"`
	RulesChannelID              Snowflake  `json:"rules_channel_id,omitempty"`
	PublicUpdatesChannelID      *Snowflake `json:"public_updates_channel_id,omitempty"`
	PreferredLocale             *string    `json:"preferred_locale,omitempty"`
	Features                    []*string  `json:"features,omitempty"`
	Description                 *string    `json:"description,omitempty"`
	PremiumProgressBarEnabled   bool       `json:"premium_progress_bar_enabled,omitempty"`
}

// Delete Guild
// DELETE /guilds/{guild.id}
// https://discord.com/developers/docs/resources/guild#delete-guild
type DeleteGuild struct {
	GuildID *Snowflake
}

// Get Guild Channels
// GET /guilds/{guild.id}/channels
// https://discord.com/developers/docs/resources/guild#get-guild-channels
type GetGuildChannels struct {
	GuildID *Snowflake
}

// Create Guild Channel
// POST /guilds/{guild.id}/channels
// https://discord.com/developers/docs/resources/guild#create-guild-channel
type CreateGuildChannel struct {
	Name                       string                 `json:"name,omitempty"`
	Type                       *Flag                  `json:"type,omitempty"`
	Topic                      *string                `json:"topic,omitempty"`
	NSFW                       bool                   `json:"nsfw,omitempty"`
	Position                   int                    `json:"position,omitempty"`
	Bitrate                    int                    `json:"bitrate,omitempty"`
	UserLimit                  int                    `json:"user_limit,omitempty"`
	PermissionOverwrites       []*PermissionOverwrite `json:"permission_overwrites,omitempty"`
	ParentID                   *Snowflake             `json:"parent_id,omitempty"`
	RateLimitPerUser           *CodeFlag              `json:"rate_limit_per_user,omitempty"`
	DefaultAutoArchiveDuration int                    `json:"default_auto_archive_duration,omitempty"`
}

// Modify Guild Channel Positions
// PATCH /guilds/{guild.id}/channels
// https://discord.com/developers/docs/resources/guild#modify-guild-channel-positions
type ModifyGuildChannelPositions struct {
	ID              Snowflake  `json:"id,omitempty"`
	Position        int        `json:"position,omitempty"`
	LockPermissions bool       `json:"lock_permissions,omitempty"`
	ParentID        *Snowflake `json:"parent_id,omitempty"`
}

// List Active Guild Threads
// GET /guilds/{guild.id}/threads/active
// https://discord.com/developers/docs/resources/guild#list-active-threads
type ListActiveGuildThreads struct {
	GuildID *Snowflake
	Threads []*Channel      `json:"threads,omitempty"`
	Members []*ThreadMember `json:"members,omitempty"`
}

// Get Guild Member
// GET /guilds/{guild.id}/members/{user.id}
// https://discord.com/developers/docs/resources/guild#get-guild-member
type GetGuildMember struct {
	UserID *Snowflake
}

// List Guild Members
// GET /guilds/{guild.id}/members
// https://discord.com/developers/docs/resources/guild#list-guild-members
type ListGuildMembers struct {
	After *Snowflake `json:"after,omitempty"`
	Limit *CodeFlag  `json:"limit,omitempty"`
}

// Search Guild Members
// GET /guilds/{guild.id}/members/search
// https://discord.com/developers/docs/resources/guild#search-guild-members
type SearchGuildMembers struct {
	GuildID *Snowflake
	Query   string    `json:"query,omitempty"`
	Limit   *CodeFlag `json:"limit,omitempty"`
}

// Add Guild Member
// PUT /guilds/{guild.id}/members/{user.id}
// https://discord.com/developers/docs/resources/guild#add-guild-member
type AddGuildMember struct {
	UserID      *Snowflake
	AccessToken string       `json:"access_token,omitempty"`
	Nick        string       `json:"nick,omitempty"`
	Roles       []*Snowflake `json:"roles,omitempty"`
	Mute        bool         `json:"mute,omitempty"`
	Deaf        bool         `json:"deaf,omitempty"`
}

// Modify Guild Member
// PATCH /guilds/{guild.id}/members/{user.id}
// https://discord.com/developers/docs/resources/guild#modify-guild-member
type ModifyGuildMember struct {
	UserID                     *Snowflake
	Nick                       string       `json:"nick,omitempty"`
	Roles                      []*Snowflake `json:"roles,omitempty"`
	Mute                       bool         `json:"mute,omitempty"`
	Deaf                       bool         `json:"deaf,omitempty"`
	ChannelID                  Snowflake    `json:"channel_id,omitempty"`
	CommunicationDisabledUntil *time.Time   `json:"communication_disabled_until,omitempty"`
}

// Modify Current Member
// PATCH /guilds/{guild.id}/members/@me
// https://discord.com/developers/docs/resources/guild#modify-current-member
type ModifyCurrentMember struct {
	GuildID Snowflake
	Nick    string `json:"nick,omitempty"`
}

// Modify Current User Nick
// PATCH /guilds/{guild.id}/members/@me/nick
// https://discord.com/developers/docs/resources/guild#modify-current-user-nick
type ModifyCurrentUserNick struct {
	GuildID Snowflake
	Nick    string `json:"nick,omitempty"`
}

// Add Guild Member Role
// PUT /guilds/{guild.id}/members/{user.id}/roles/{role.id}
// https://discord.com/developers/docs/resources/guild#add-guild-member-role
type AddGuildMemberRole struct {
	RoleID Snowflake
}

// Remove Guild Member Role
// DELETE /guilds/{guild.id}/members/{user.id}/roles/{role.id}
// https://discord.com/developers/docs/resources/guild#remove-guild-member-role
type RemoveGuildMemberRole struct {
	RoleID Snowflake
}

// Remove Guild Member
// DELETE /guilds/{guild.id}/members/{user.id}
// https://discord.com/developers/docs/resources/guild#remove-guild-member
type RemoveGuildMember struct {
	UserID *Snowflake
}

// Get Guild Bans
// GET /guilds/{guild.id}/bans
// https://discord.com/developers/docs/resources/guild#get-guild-bans
type GetGuildBans struct {
	GuildID Snowflake
	Before  *Snowflake `json:"before,omitempty"`
	After   *Snowflake `json:"after,omitempty"`
	Limit   *CodeFlag  `json:"limit,omitempty"`
}

// Get Guild Ban
// GET /guilds/{guild.id}/bans/{user.id}
// https://discord.com/developers/docs/resources/guild#get-guild-ban
type GetGuildBan struct {
	UserID *Snowflake
}

// Create Guild Ban
// PUT /guilds/{guild.id}/bans/{user.id}
// https://discord.com/developers/docs/resources/guild#create-guild-ban
type CreateGuildBan struct {
	UserID            *Snowflake
	DeleteMessageDays *Flag   `json:"delete_message_days,omitempty"`
	Reason            *string `json:"reason,omitempty"`
}

// Remove Guild Ban
// DELETE /guilds/{guild.id}/bans/{user.id}
// https://discord.com/developers/docs/resources/guild#remove-guild-ban
type RemoveGuildBan struct {
	UserID *Snowflake
}

// Get Guild Roles
// GET /guilds/{guild.id}/roles
// https://discord.com/developers/docs/resources/guild#get-guild-roles
type GetGuildRoles struct {
	GuildID Snowflake
}

// Create Guild Role
// POST /guilds/{guild.id}/roles
// https://discord.com/developers/docs/resources/guild#create-guild-role
type CreateGuildRole struct {
	GuildID      Snowflake
	Name         string  `json:"name,omitempty"`
	Permissions  string  `json:"permissions,omitempty"`
	Color        *int    `json:"color,omitempty"`
	Hoist        bool    `json:"hoist,omitempty"`
	Icon         *int    `json:"icon,omitempty"`
	UnicodeEmoji *string `json:"unicode_emoji,omitempty"`
	Mentionable  bool    `json:"mentionable,omitempty"`
}

// Modify Guild Role Positions
// PATCH /guilds/{guild.id}/roles
// https://discord.com/developers/docs/resources/guild#modify-guild-role-positions
type ModifyGuildRolePositions struct {
	ID       Snowflake `json:"id,omitempty"`
	Position int       `json:"position,omitempty"`
}

// Modify Guild Role
// PATCH /guilds/{guild.id}/roles/{role.id}
// https://discord.com/developers/docs/resources/guild#modify-guild-role
type ModifyGuildRole struct {
	RoleID       Snowflake
	Name         string  `json:"name,omitempty"`
	Permissions  int64   `json:"permissions,string,omitempty"`
	Color        *int    `json:"color,omitempty"`
	Hoist        bool    `json:"hoist,omitempty"`
	Icon         *int    `json:"icon,omitempty"`
	UnicodeEmoji *string `json:"unicode_emoji,omitempty"`
	Mentionable  bool    `json:"mentionable,omitempty"`
}

// Delete Guild Role
// DELETE /guilds/{guild.id}/roles/{role.id}
// https://discord.com/developers/docs/resources/guild#delete-guild-role
type DeleteGuildRole struct {
	RoleID Snowflake
}

// Get Guild Prune Count
// GET /guilds/{guild.id}/prune
// https://discord.com/developers/docs/resources/guild#get-guild-prune-count
type GetGuildPruneCount struct {
	GuildID      Snowflake
	Days         Flag         `json:"days,omitempty"`
	IncludeRoles []*Snowflake `json:"include_roles,omitempty"`
}

// Begin Guild Prune
// POST /guilds/{guild.id}/prune
// https://discord.com/developers/docs/resources/guild#begin-guild-prune
type BeginGuildPrune struct {
	GuildID           Snowflake
	Days              Flag         `json:"days,omitempty"`
	ComputePruneCount bool         `json:"compute_prune_count,omitempty"`
	IncludeRoles      []*Snowflake `json:"include_roles,omitempty"`
	Reason            *string      `json:"reason,omitempty"`
}

// Get Guild Voice Regions
// GET /guilds/{guild.id}/regions
// https://discord.com/developers/docs/resources/guild#get-guild-voice-regions
type GetGuildVoiceRegions struct {
	GuildID Snowflake
}

// Get Guild Invites
// GET /guilds/{guild.id}/invites
// https://discord.com/developers/docs/resources/guild#get-guild-invites
type GetGuildInvites struct {
	GuildID Snowflake
}

// Get Guild Integrations
// GET /guilds/{guild.id}/integrations
// https://discord.com/developers/docs/resources/guild#get-guild-integrations
type GetGuildIntegrations struct {
	GuildID Snowflake
}

// Delete Guild Integration
// DELETE /guilds/{guild.id}/integrations/{integration.id}
// https://discord.com/developers/docs/resources/guild#delete-guild-integration
type DeleteGuildIntegration struct {
	IntegrationID Snowflake
}

// Get Guild Widget Settings
// GET /guilds/{guild.id}/widget
// https://discord.com/developers/docs/resources/guild#get-guild-widget-settings
type GetGuildWidgetSettings struct {
	GuildID Snowflake
}

// Modify Guild Widget
// PATCH /guilds/{guild.id}/widget
// https://discord.com/developers/docs/resources/guild#modify-guild-widget
type ModifyGuildWidget struct {
	GuildID Snowflake
}

// Get Guild Widget
// GET /guilds/{guild.id}/widget.json
// https://discord.com/developers/docs/resources/guild#get-guild-widget
type GetGuildWidget struct {
	GuildID Snowflake
}

// Get Guild Vanity URL
// GET /guilds/{guild.id}/vanity-url
// https://discord.com/developers/docs/resources/guild#get-guild-vanity-url
type GetGuildVanityURL struct {
	GuildID Snowflake
}

// Get Guild Widget Image
// GET /guilds/{guild.id}/widget.png
// https://discord.com/developers/docs/resources/guild#get-guild-widget-image
type GetGuildWidgetImage struct {
	GuildID Snowflake
	Style   string `json:"style,omitempty"`
}

// Get Guild Welcome Screen
// GET /guilds/{guild.id}/welcome-screen
// https://discord.com/developers/docs/resources/guild#get-guild-welcome-screen
type GetGuildWelcomeScreen struct {
	GuildID Snowflake
}

// Modify Guild Welcome Screen
// PATCH /guilds/{guild.id}/welcome-screen
// https://discord.com/developers/docs/resources/guild#modify-guild-welcome-screen
type ModifyGuildWelcomeScreen struct {
	GuildID         Snowflake
	Enabled         bool                    `json:"enabled,omitempty"`
	WelcomeChannels []*WelcomeScreenChannel `json:"welcome_channels,omitempty"`
	Description     *string                 `json:"description,omitempty"`
}

// Modify Current User Voice State
// PATCH /guilds/{guild.id}/voice-states/@me
// https://discord.com/developers/docs/resources/guild#modify-current-user-voice-state
type ModifyCurrentUserVoiceState struct {
	ChannelID               Snowflake
	Suppress                bool       `json:"suppress,omitempty"`
	RequestToSpeakTimestamp *time.Time `json:"request_to_speak_timestamp,omitempty"`
}

// Modify User Voice State
// PATCH /guilds/{guild.id}/voice-states/{user.id}
// https://discord.com/developers/docs/resources/guild#modify-user-voice-state
type ModifyUserVoiceState struct {
	ChannelID Snowflake
	Suppress  bool `json:"suppress,omitempty"`
}

// List Scheduled Events for Guild
// GET /guilds/{guild.id}/scheduled-events
// https://discord.com/developers/docs/resources/guild-scheduled-event#list-scheduled-events-for-guild
type ListScheduledEventsforGuild struct {
	WithUserCount bool `json:"with_user_count,omitempty"`
}

// Create Guild Scheduled Event
// POST /guilds/{guild.id}/scheduled-events
// https://discord.com/developers/docs/resources/guild-scheduled-event#create-guild-scheduled-event
type CreateGuildScheduledEvent struct {
	ChannelID          *Snowflake                         `json:"channel_id,omitempty"`
	EntityMetadata     *GuildScheduledEventEntityMetadata `json:"entity_metadata,omitempty"`
	Name               *string                            `json:"name,omitempty"`
	PrivacyLevel       Flag                               `json:"privacy_level,omitempty"`
	ScheduledStartTime Snowflake                          `json:"scheduled_start_time,omitempty"`
	ScheduledEndTime   Snowflake                          `json:"scheduled_end_time,omitempty"`
	Description        *string                            `json:"description,omitempty"`
	EntityType         *Flag                              `json:"entity_type,omitempty"`
	Image              *string                            `json:"image,omitempty"`
}

// Get Guild Scheduled Event
// GET /guilds/{guild.id}/scheduled-events/{guild_scheduled_event.id}
// https://discord.com/developers/docs/resources/guild-scheduled-event#get-guild-scheduled-event
type GetGuildScheduledEvent struct {
	WithUserCount bool `json:"with_user_count,omitempty"`
}

// Modify Guild Scheduled Event
// PATCH /guilds/{guild.id}/scheduled-events/{guild_scheduled_event.id}
// https://discord.com/developers/docs/resources/guild-scheduled-event#modify-guild-scheduled-event
type ModifyGuildScheduledEvent struct {
	ChannelID          *Snowflake                         `json:"channel_id,omitempty"`
	EntityMetadata     *GuildScheduledEventEntityMetadata `json:"entity_metadata,omitempty"`
	Name               *string                            `json:"name,omitempty"`
	PrivacyLevel       Flag                               `json:"privacy_level,omitempty"`
	ScheduledStartTime Snowflake                          `json:"scheduled_start_time,omitempty"`
	ScheduledEndTime   Snowflake                          `json:"scheduled_end_time,omitempty"`
	Description        *string                            `json:"description,omitempty"`
	EntityType         *Flag                              `json:"entity_type,omitempty"`
	Image              *string                            `json:"image,omitempty"`
	Status             Flag                               `json:"status,omitempty"`
}

// Delete Guild Scheduled Event
// DELETE /guilds/{guild.id}/scheduled-events/{guild_scheduled_event.id}
// https://discord.com/developers/docs/resources/guild-scheduled-event#delete-guild-scheduled-event
type DeleteGuildScheduledEvent struct {
	GuildScheduledEventID Snowflake
}

// Get Guild Scheduled Event Users
// GET /guilds/{guild.id}/scheduled-events/{guild_scheduled_event.id}/users
// https://discord.com/developers/docs/resources/guild-scheduled-event#get-guild-scheduled-event-users
type GetGuildScheduledEventUsers struct {
	Limit      Flag       `json:"limit,omitempty"`
	WithMember bool       `json:"with_member,omitempty"`
	Before     *Snowflake `json:"before,omitempty"`
	After      *Snowflake `json:"after,omitempty"`
}

// Get Guild Template
// GET /guilds/templates/{template.code}
// https://discord.com/developers/docs/resources/guild-template#get-guild-template
type GetGuildTemplate struct {
	TemplateCode string `json:"code,omitempty"`
}

// Create Guild from Guild Template
// POST /guilds/templates/{template.code}
// https://discord.com/developers/docs/resources/guild-template#create-guild-from-guild-template
type CreateGuildfromGuildTemplate struct {
	Name string `json:"name,omitempty"`
	Icon string `json:"icon,omitempty"`
}

// Get Guild Templates
// GET /guilds/{guild.id}/templates
// https://discord.com/developers/docs/resources/guild-template#get-guild-templates
type GetGuildTemplates struct {
	GuildID Snowflake
}

// Create Guild Template
// POST /guilds/{guild.id}/templates
// https://discord.com/developers/docs/resources/guild-template#create-guild-template
type CreateGuildTemplate struct {
	Name        string  `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

// Sync Guild Template
// PUT /guilds/{guild.id}/templates/{template.code}
// https://discord.com/developers/docs/resources/guild-template#sync-guild-template
type SyncGuildTemplate struct {
	TemplateCode string `json:"code,omitempty"`
}

// Modify Guild Template
// PATCH /guilds/{guild.id}/templates/{template.code}
// https://discord.com/developers/docs/resources/guild-template#modify-guild-template
type ModifyGuildTemplate struct {
	TemplateCode string  `json:"code,omitempty"`
	Name         string  `json:"name,omitempty"`
	Description  *string `json:"description,omitempty"`
}

// Delete Guild Template
// DELETE /guilds/{guild.id}/templates/{template.code}
// https://discord.com/developers/docs/resources/guild-template#delete-guild-template
type DeleteGuildTemplate struct {
	TemplateCode string `json:"code,omitempty"`
}

// Create Interaction Response
// POST /interactions/{interaction.id}/{interaction.token}/callback
// https://discord.com/developers/docs/interactions/receiving-and-responding#create-interaction-response
type CreateInteractionResponse struct {
	InteractionToken string `json:"token,omitempty"`
}

// Get Original Interaction Response
// GET /webhooks/{application.id}/{interaction.token}/messages/@original
// https://discord.com/developers/docs/interactions/receiving-and-responding#get-original-interaction-response
type GetOriginalInteractionResponse struct {
	InteractionToken string `json:"token,omitempty"`
}

// Edit Original Interaction Response
// PATCH /webhooks/{application.id}/{interaction.token}/messages/@original
// https://discord.com/developers/docs/interactions/receiving-and-responding#edit-original-interaction-response
type EditOriginalInteractionResponse struct {
	InteractionToken string `json:"token,omitempty"`
}

// Delete Original Interaction Response
// DELETE /webhooks/{application.id}/{interaction.token}/messages/@original
// https://discord.com/developers/docs/interactions/receiving-and-responding#delete-original-interaction-response
type DeleteOriginalInteractionResponse struct {
	InteractionToken string `json:"token,omitempty"`
}

// Create Followup Message
// POST /webhooks/{application.id}/{interaction.token}
// https://discord.com/developers/docs/interactions/receiving-and-responding#create-followup-message
type CreateFollowupMessage struct {
	InteractionToken string `json:"token,omitempty"`
}

// Get Followup Message
// GET /webhooks/{application.id}/{interaction.token}/messages/{message.id}
// https://discord.com/developers/docs/interactions/receiving-and-responding#get-followup-message
type GetFollowupMessage struct {
	MessageID Snowflake
}

// Edit Followup Message
// PATCH /webhooks/{application.id}/{interaction.token}/messages/{message.id}
// https://discord.com/developers/docs/interactions/receiving-and-responding#edit-followup-message
type EditFollowupMessage struct {
	MessageID Snowflake
}

// Delete Followup Message
// DELETE /webhooks/{application.id}/{interaction.token}/messages/{message.id}
// https://discord.com/developers/docs/interactions/receiving-and-responding#delete-followup-message
type DeleteFollowupMessage struct {
	MessageID Snowflake
}

// Get Invite
// GET /invites/{invite.code}
// https://discord.com/developers/docs/resources/invite#get-invite
type GetInvite struct {
	GuildScheduledEventID Snowflake `json:"guild_scheduled_event_id,omitempty"`
	WithCounts            bool      `json:"with_counts,omitempty"`
	WithExpiration        bool      `json:"with_expiration,omitempty"`
}

// Delete Invite
// DELETE /invites/{invite.code}
// https://discord.com/developers/docs/resources/invite#delete-invite
type DeleteInvite struct {
	InviteCode string `json:"code,omitempty"`
}

// Get Current Bot Application Information
// GET /oauth2/applications/@me
// https://discord.com/developers/docs/topics/oauth2#get-current-bot-application-information
type GetCurrentBotApplicationInformation struct{}

// Get Current Authorization Information
// GET /oauth2/@me
// https://discord.com/developers/docs/topics/oauth2#get-current-authorization-information
type GetCurrentAuthorizationInformation struct{}

// Create Stage Instance
// POST /stage-instances
// https://discord.com/developers/docs/resources/stage-instance#create-stage-instance
type CreateStageInstance struct {
	ChannelID             Snowflake
	Topic                 string `json:"topic,omitempty"`
	PrivacyLevel          Flag   `json:"privacy_level,omitempty"`
	SendStartNotification bool   `json:"send_start_notification,omitempty"`
}

// Get Stage Instance
// GET /stage-instances/{channel.id}
// https://discord.com/developers/docs/resources/stage-instance#get-stage-instance
type GetStageInstance struct {
	ChannelID Snowflake
}

// Modify Stage Instance
// PATCH /stage-instances/{channel.id}
// https://discord.com/developers/docs/resources/stage-instance#modify-stage-instance
type ModifyStageInstance struct {
	ChannelID    Snowflake
	Topic        string `json:"topic,omitempty"`
	PrivacyLevel Flag   `json:"privacy_level,omitempty"`
}

// Delete Stage Instance
// DELETE /stage-instances/{channel.id}
// https://discord.com/developers/docs/resources/stage-instance#delete-stage-instance
type DeleteStageInstance struct {
	ChannelID Snowflake
}

// Get Sticker
// GET /stickers/{sticker.id}
// https://discord.com/developers/docs/resources/sticker#get-sticker
type GetSticker struct {
	StickerID Snowflake
}

// List Nitro Sticker Packs
// GET /sticker-packs
// https://discord.com/developers/docs/resources/sticker#list-nitro-sticker-packs
type ListNitroStickerPacks struct {
	StickerPacks []*StickerPack `json:"sticker_packs,omitempty"`
}

// List Guild Stickers
// GET /guilds/{guild.id}/stickers
// https://discord.com/developers/docs/resources/sticker#list-guild-stickers
type ListGuildStickers struct {
	GuildID Snowflake
}

// Get Guild Sticker
// GET /guilds/{guild.id}/stickers/{sticker.id}
// https://discord.com/developers/docs/resources/sticker#get-guild-sticker
type GetGuildSticker struct {
	StickerID Snowflake
}

// Create Guild Sticker
// POST /guilds/{guild.id}/stickers
// https://discord.com/developers/docs/resources/sticker#create-guild-sticker
type CreateGuildSticker struct {
	GuildID     Snowflake
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	Tags        *string `json:"tags,omitempty"`
	Files       []byte  `dasgo:"files"`
}

// Modify Guild Sticker
// PATCH /guilds/{guild.id}/stickers/{sticker.id}
// https://discord.com/developers/docs/resources/sticker#modify-guild-sticker
type ModifyGuildSticker struct {
	StickerID   Snowflake
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	Tags        *string `json:"tags,omitempty"`
}

// Delete Guild Sticker
// DELETE /guilds/{guild.id}/stickers/{sticker.id}
// https://discord.com/developers/docs/resources/sticker#delete-guild-sticker
type DeleteGuildSticker struct {
	StickerID Snowflake
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
	Before *Snowflake `json:"before,omitempty"`
	After  *Snowflake `json:"after,omitempty"`
	Limit  Flag       `json:"limit,omitempty"`
}

// Get Current User Guild Member
// GET /users/@me/guilds/{guild.id}/member
// https://discord.com/developers/docs/resources/user#get-current-user-guild-member
type GetCurrentUserGuildMember struct {
	GuildID Snowflake
}

// Leave Guild
// DELETE /users/@me/guilds/{guild.id}
// https://discord.com/developers/docs/resources/user#leave-guild
type LeaveGuild struct {
	GuildID Snowflake
}

// Create DM
// POST /users/@me/channels
// https://discord.com/developers/docs/resources/user#create-dm
type CreateDM struct {
	RecipientID Snowflake `json:"recipient_id,omitempty"`
}

// Create Group DM
// POST /users/@me/channels
// https://discord.com/developers/docs/resources/user#create-group-dm
type CreateGroupDM struct {
	AccessTokens []*string            `json:"access_tokens,omitempty"`
	Nicks        map[Snowflake]string `json:"nicks,omitempty"`
}

// Get User Connections
// GET /users/@me/connections
// https://discord.com/developers/docs/resources/user#get-user-connections
type GetUserConnections struct {
	RecipientID Snowflake `json:"recipient_id,omitempty"`
}

// List Voice Regions
// GET /voice/regions
// https://discord.com/developers/docs/resources/voice#list-voice-regions
type ListVoiceRegions struct{}

// Create Webhook
// POST /channels/{channel.id}/webhooks
// https://discord.com/developers/docs/resources/webhook#create-webhook
type CreateWebhook struct {
	Name   string `json:"name,omitempty"`
	Avatar string `json:"avatar,omitempty"`
}

// Get Channel Webhooks
// GET /channels/{channel.id}/webhooks
// https://discord.com/developers/docs/resources/webhook#get-channel-webhooks
type GetChannelWebhooks struct {
	ChannelID Snowflake
}

// Get Guild Webhooks
// GET /guilds/{guild.id}/webhooks
// https://discord.com/developers/docs/resources/webhook#get-guild-webhooks
type GetGuildWebhooks struct {
	GuildID Snowflake
}

// Get Webhook
// GET /webhooks/{webhook.id}
// https://discord.com/developers/docs/resources/webhook#get-webhook
type GetWebhook struct {
	WebhookID Snowflake
}

// Get Webhook with Token
// GET /webhooks/{webhook.id}/{webhook.token}
// https://discord.com/developers/docs/resources/webhook#get-webhook-with-token
type GetWebhookwithToken struct {
	WebhookID Snowflake
}

// Modify Webhook
// PATCH /webhooks/{webhook.id}
// https://discord.com/developers/docs/resources/webhook#modify-webhook
type ModifyWebhook struct {
	WebhookID Snowflake
	Name      string    `json:"name,omitempty"`
	Avatar    string    `json:"avatar,omitempty"`
	ChannelID Snowflake `json:"channel_id,omitempty"`
}

// Modify Webhook with Token
// PATCH /webhooks/{webhook.id}/{webhook.token}
// https://discord.com/developers/docs/resources/webhook#modify-webhook-with-token
type ModifyWebhookwithToken struct {
	WebhookID Snowflake
}

// Delete Webhook
// DELETE /webhooks/{webhook.id}
// https://discord.com/developers/docs/resources/webhook#delete-webhook
type DeleteWebhook struct {
	WebhookID Snowflake
}

// Delete Webhook with Token
// DELETE /webhooks/{webhook.id}/{webhook.token}
// https://discord.com/developers/docs/resources/webhook#delete-webhook-with-token
type DeleteWebhookwithToken struct {
	WebhookID Snowflake
}

// Execute Webhook
// POST /webhooks/{webhook.id}/{webhook.token}
// https://discord.com/developers/docs/resources/webhook#execute-webhook
type ExecuteWebhook struct {
	Wait            bool             `json:"wait,omitempty"`
	ThreadID        Snowflake        `json:"thread_id,omitempty"`
	Content         string           `json:"content,omitempty"`
	Username        string           `json:"username,omitempty"`
	AvatarURL       string           `json:"avatar_url,omitempty"`
	TTS             bool             `json:"tts,omitempty"`
	Files           []byte           `dasgo:"files"`
	Components      []Component      `json:"components,omitempty"`
	Embeds          []*Embed         `json:"embeds,omitempty"`
	AllowedMentions *AllowedMentions `json:"allowed_mentions,omitempty"`
	PayloadJSON     string           `json:"payload_json,omitempty"`
	Attachments     []*Attachment    `json:"attachments,omitempty"`
}

// Execute Slack-Compatible Webhook
// POST /webhooks/{webhook.id}/{webhook.token}/slack
// https://discord.com/developers/docs/resources/webhook#execute-slackcompatible-webhook
type ExecuteSlackCompatibleWebhook struct {
	ThreadID Snowflake
	Wait     bool `json:"wait,omitempty"`
}

// Execute GitHub-Compatible Webhook
// POST /webhooks/{webhook.id}/{webhook.token}/github
// https://discord.com/developers/docs/resources/webhook#execute-githubcompatible-webhook
type ExecuteGitHubCompatibleWebhook struct {
	ThreadID Snowflake
	Wait     bool `json:"wait,omitempty"`
}

// Get Webhook Message
// GET /webhooks/{webhook.id}/{webhook.token}/messages/{message.id}
// https://discord.com/developers/docs/resources/webhook#get-webhook-message
type GetWebhookMessage struct {
	ThreadID Snowflake
}

// Edit Webhook Message
// PATCH /webhooks/{webhook.id}/{webhook.token}/messages/{message.id}
// https://discord.com/developers/docs/resources/webhook#edit-webhook-message
type EditWebhookMessage struct {
	WebhookID       Snowflake
	ThreadID        Snowflake        `json:"thread_id,omitempty"`
	Content         *string          `json:"content,omitempty"`
	Components      []*Component     `json:"components,omitempty"`
	Embeds          []*Embed         `json:"embeds,omitempty"`
	Files           []byte           `dasgo:"files"`
	AllowedMentions *AllowedMentions `json:"allowed_mentions,omitempty"`
	PayloadJSON     string           `json:"payload_json,omitempty"`
	Attachments     []*Attachment    `json:"attachments,omitempty"`
}

// Delete Webhook Message
// DELETE /webhooks/{webhook.id}/{webhook.token}/messages/{message.id}
// https://discord.com/developers/docs/resources/webhook#delete-webhook-message
type DeleteWebhookMessage struct {
	ThreadID Snowflake
}

// Application Object
// https://discord.com/developers/docs/resources/application
type Application struct {
	ID                  Snowflake      `json:"id,omitempty"`
	Name                string         `json:"name,omitempty"`
	Icon                string         `json:"icon,omitempty"`
	Description         string         `json:"description,omitempty"`
	RPCOrigins          []string       `json:"rpc_origins,omitempty"`
	BotPublic           bool           `json:"bot_public,omitempty"`
	BotRequireCodeGrant bool           `json:"bot_require_code_grant,omitempty"`
	TermsOfServiceURL   string         `json:"terms_of_service_url,omitempty"`
	PrivacyProxyURL     string         `json:"privacy_policy_url,omitempty"`
	Owner               *User          `json:"owner,omitempty"`
	VerifyKey           string         `json:"verify_key,omitempty"`
	Team                *Team          `json:"team,omitempty"`
	GuildID             Snowflake      `json:"guild_id,omitempty"`
	PrimarySKUID        Snowflake      `json:"primary_sku_id,omitempty"`
	Slug                *string        `json:"slug,omitempty"`
	CoverImage          string         `json:"cover_image,omitempty"`
	Flags               Flag           `json:"flags,omitempty"`
	Summary             string         `json:"summary,omitempty"`
	InstallParams       *InstallParams `json:"install_params,omitempty"`
	CustomInstallURL    string         `json:"custom_install_url,omitempty"`
}

// Application Flags
// https://discord.com/developers/docs/resources/application#application-object-application-flags
const (
	FlagFlagsApplicationGATEWAY_PRESENCE                 = 1 << 12
	FlagFlagsApplicationGATEWAY_PRESENCE_LIMITED         = 1 << 13
	FlagFlagsApplicationGATEWAY_GUILD_MEMBERS            = 1 << 14
	FlagFlagsApplicationGATEWAY_GUILD_MEMBERS_LIMITED    = 1 << 15
	FlagFlagsApplicationVERIFICATION_PENDING_GUILD_LIMIT = 1 << 16
	FlagFlagsApplicationEMBEDDED                         = 1 << 17
	FlagFlagsApplicationGATEWAY_MESSAGE_CONTENT          = 1 << 18
	FlagFlagsApplicationGATEWAY_MESSAGE_CONTENT_LIMITED  = 1 << 19
)

// Install Params Object
// https://discord.com/developers/docs/resources/application#install-params-object
type InstallParams struct {
	Scopes      []string `json:"scopes,omitempty"`
	Permissions string   `json:"permissions,omitempty"`
}

// Application Command Structure
// https://discord.com/developers/docs/interactions/application-commands#application-command-object-application-command-structure
type ApplicationCommand struct {
	ID                       Snowflake                   `json:"id,omitempty"`
	Type                     Flag                        `json:"type,omitempty"`
	ApplicationID            Snowflake                   `json:"application_id,omitempty"`
	GuildID                  Snowflake                   `json:"guild_id,omitempty"`
	Name                     string                      `json:"name,omitempty"`
	NameLocalizations        map[Flag]string             `json:"name_localizations,omitempty"`
	Description              string                      `json:"description,omitempty"`
	DescriptionLocalizations map[Flag]string             `json:"description_localizations,omitempty"`
	Options                  []*ApplicationCommandOption `json:"options,omitempty"`
	DefaultPermission        bool                        `json:"default_permission,omitempty"`
	Version                  Snowflake                   `json:"version,omitempty"`
}

// Application Command Types
// https://discord.com/developers/docs/interactions/application-commands#application-command-object-application-command-types
const (
	FlagTypesCommandApplicationCHAT_INPUT = 1
	FlagTypesCommandApplicationUSER       = 2
	FlagTypesCommandApplicationMESSAGE    = 3
)

// Application Command Option Structure
// https://discord.com/developers/docs/interactions/application-commands#application-command-object-application-command-option-structure
type ApplicationCommandOption struct {
	Type                     Flag                              `json:"type,omitempty"`
	Name                     string                            `json:"name,omitempty"`
	NameLocalizations        map[Flag]string                   `json:"name_localizations,omitempty"`
	Description              string                            `json:"description,omitempty"`
	DescriptionLocalizations map[Flag]string                   `json:"description_localizations,omitempty"`
	Required                 bool                              `json:"required,omitempty"`
	Choices                  []*ApplicationCommandOptionChoice `json:"choices,omitempty"`
	Options                  []*ApplicationCommandOption       `json:"options,omitempty"`
	ChannelTypes             []*Flag                           `json:"channel_types,omitempty"`
	MinValue                 float64                           `json:"min_value,omitempty"`
	MaxValue                 float64                           `json:"max_value,omitempty"`
	Autocomplete             bool                              `json:"autocomplete,omitempty"`
}

// Application Command Option Type
// https://discord.com/developers/docs/interactions/application-commands#application-command-object-application-command-option-type
const (
	FlagTypeOptionCommandApplicationSUB_COMMAND       = 1
	FlagTypeOptionCommandApplicationSUB_COMMAND_GROUP = 2
	FlagTypeOptionCommandApplicationSTRING            = 3
	FlagTypeOptionCommandApplicationINTEGER           = 4
	FlagTypeOptionCommandApplicationBOOLEAN           = 5
	FlagTypeOptionCommandApplicationUSER              = 6
	FlagTypeOptionCommandApplicationCHANNEL           = 7
	FlagTypeOptionCommandApplicationROLE              = 8
	FlagTypeOptionCommandApplicationMENTIONABLE       = 9
	FlagTypeOptionCommandApplicationNUMBER            = 10
	FlagTypeOptionCommandApplicationATTACHMENT        = 11
)

// Application Command Option Choice
// https://discord.com/developers/docs/interactions/application-commands#application-command-object-application-command-option-choice-structure
type ApplicationCommandOptionChoice struct {
	Name              string          `json:"name,omitempty"`
	NameLocalizations map[Flag]string `json:"name_localizations,omitempty"`
	Value             interface{}     `json:"value,omitempty"`
}

// Application Command Interaction Data Option Structure
// https://discord.com/developers/docs/interactions/application-commands#application-command-object-application-command-interaction-data-option-structure
type ApplicationCommandInteractionDataOption struct {
	Name    string                                     `json:"name,omitempty"`
	Type    Flag                                       `json:"type,omitempty"`
	Value   interface{}                                `json:"value,omitempty"`
	Options []*ApplicationCommandInteractionDataOption `json:"options,omitempty"`
	Focused bool                                       `json:"focused,omitempty"`
}

// Guild Application Command Permissions Object
// https://discord.com/developers/docs/interactions/application-commands#application-command-permissions-object-guild-application-command-permissions-structure
type GuildApplicationCommandPermissions struct {
	ID            Snowflake                        `json:"id,omitempty"`
	ApplicationID Snowflake                        `json:"application_id,omitempty"`
	GuildID       Snowflake                        `json:"guild_id,omitempty"`
	Permissions   []*ApplicationCommandPermissions `json:"permissions,omitempty"`
}

// Application Command Permissions Structure
// https://discord.com/developers/docs/interactions/application-commands#application-command-permissions-object-application-command-permissions-structure
type ApplicationCommandPermissions struct {
	ID         Snowflake `json:"id,omitempty"`
	Type       Flag      `json:"type,omitempty"`
	Permission bool      `json:"permission,omitempty"`
}

// Application Command Permission Type
// https://discord.com/developers/docs/interactions/application-commands#application-command-permissions-object-application-command-permission-type
const (
	FlagTypePermissionCommandApplicationROLE = 1
	FlagTypePermissionCommandApplicationUSER = 2
)

// Audit Log Object
// https://discord.com/developers/docs/resources/audit-log
type AuditLog struct {
	AuditLogEntries      []*AuditLogEntry       `json:"audit_log_entries,omitempty"`
	GuildScheduledEvents []*GuildScheduledEvent `json:"guild_scheduled_events,omitempty"`
	Integration          []*Integration         `json:"integrations,omitempty"`
	Threads              []*Channel             `json:"threads,omitempty"`
	Users                []*User                `json:"users,omitempty"`
	Webhooks             []*Webhook             `json:"webhooks,omitempty"`
}

// Audit Log Events
// https://discord.com/developers/docs/resources/audit-log#audit-log-entry-object-audit-log-events
const (
	FlagEventsLogAuditGUILD_UPDATE                 = 1
	FlagEventsLogAuditCHANNEL_CREATE               = 10
	FlagEventsLogAuditCHANNEL_UPDATE               = 11
	FlagEventsLogAuditCHANNEL_DELETE               = 12
	FlagEventsLogAuditCHANNEL_OVERWRITE_CREATE     = 13
	FlagEventsLogAuditCHANNEL_OVERWRITE_UPDATE     = 14
	FlagEventsLogAuditCHANNEL_OVERWRITE_DELETE     = 15
	FlagEventsLogAuditMEMBER_KICK                  = 20
	FlagEventsLogAuditMEMBER_PRUNE                 = 21
	FlagEventsLogAuditMEMBER_BAN_ADD               = 22
	FlagEventsLogAuditMEMBER_BAN_REMOVE            = 23
	FlagEventsLogAuditMEMBER_UPDATE                = 24
	FlagEventsLogAuditMEMBER_ROLE_UPDATE           = 25
	FlagEventsLogAuditMEMBER_MOVE                  = 26
	FlagEventsLogAuditMEMBER_DISCONNECT            = 27
	FlagEventsLogAuditBOT_ADD                      = 28
	FlagEventsLogAuditROLE_CREATE                  = 30
	FlagEventsLogAuditROLE_UPDATE                  = 31
	FlagEventsLogAuditROLE_DELETE                  = 32
	FlagEventsLogAuditINVITE_CREATE                = 40
	FlagEventsLogAuditINVITE_UPDATE                = 41
	FlagEventsLogAuditINVITE_DELETE                = 42
	FlagEventsLogAuditWEBHOOK_CREATE               = 50
	FlagEventsLogAuditWEBHOOK_UPDATE               = 51
	FlagEventsLogAuditWEBHOOK_DELETE               = 52
	FlagEventsLogAuditEMOJI_CREATE                 = 60
	FlagEventsLogAuditEMOJI_UPDATE                 = 61
	FlagEventsLogAuditEMOJI_DELETE                 = 62
	FlagEventsLogAuditMESSAGE_DELETE               = 72
	FlagEventsLogAuditMESSAGE_BULK_DELETE          = 73
	FlagEventsLogAuditMESSAGE_PIN                  = 74
	FlagEventsLogAuditMESSAGE_UNPIN                = 75
	FlagEventsLogAuditINTEGRATION_CREATE           = 80
	FlagEventsLogAuditINTEGRATION_UPDATE           = 81
	FlagEventsLogAuditINTEGRATION_DELETE           = 82
	FlagEventsLogAuditSTAGE_INSTANCE_CREATE        = 83
	FlagEventsLogAuditSTAGE_INSTANCE_UPDATE        = 84
	FlagEventsLogAuditSTAGE_INSTANCE_DELETE        = 85
	FlagEventsLogAuditSTICKER_CREATE               = 90
	FlagEventsLogAuditSTICKER_UPDATE               = 91
	FlagEventsLogAuditSTICKER_DELETE               = 92
	FlagEventsLogAuditGUILD_SCHEDULED_EVENT_CREATE = 100
	FlagEventsLogAuditGUILD_SCHEDULED_EVENT_UPDATE = 101
	FlagEventsLogAuditGUILD_SCHEDULED_EVENT_DELETE = 102
	FlagEventsLogAuditTHREAD_CREATE                = 110
	FlagEventsLogAuditTHREAD_UPDATE                = 111
	FlagEventsLogAuditTHREAD_DELETE                = 112
)

// Audit Log Entry Object
// https://discord.com/developers/docs/resources/audit-log#audit-log-object-audit-log-structure
type AuditLogEntry struct {
	TargetID   string            `json:"target_id,omitempty"`
	Changes    []*AuditLogChange `json:"changes,omitempty"`
	UserID     Snowflake         `json:"user_id,omitempty"`
	ID         Snowflake         `json:"id,omitempty"`
	ActionType Flag              `json:"action_type,omitempty"`
	Options    *AuditLogOptions  `json:"options,omitempty"`
	Reason     *string           `json:"reason,omitempty"`
}

// Optional Audit Entry Info
// https://discord.com/developers/docs/resources/audit-log#audit-log-entry-object-optional-audit-entry-info
type AuditLogOptions struct {
	ChannelID        Snowflake `json:"channel_id,omitempty"`
	Count            string    `json:"count,omitempty"`
	DeleteMemberDays string    `json:"delete_member_days,omitempty"`
	ID               Snowflake `json:"id,omitempty"`
	MembersRemoved   string    `json:"members_removed,omitempty"`
	MessageID        Snowflake `json:"message_id,omitempty"`
	RoleName         string    `json:"role_name,omitempty"`
}

// Audit Log Change Object
// https://discord.com/developers/docs/resources/audit-log#audit-log-change-object
type AuditLogChange struct {
	NewValue interface{} `json:"new_value,omitempty"`
	OldValue interface{} `json:"old_value,omitempty"`
	Key      string      `json:"key,omitempty"`
}

// Audit Log Change Key
// https://discord.com/developers/docs/resources/audit-log#audit-log-change-object-audit-log-change-key
const (
	FlagKeyChangeLogAuditafk_channel_id                = "afk_channel_id"
	FlagKeyChangeLogAuditafk_timeout                   = "afk_timeout"
	FlagKeyChangeLogAuditallow                         = "allow"
	FlagKeyChangeLogAuditapplication_id                = "application_id"
	FlagKeyChangeLogAuditarchived                      = "archived"
	FlagKeyChangeLogAuditasset                         = "asset"
	FlagKeyChangeLogAuditauto_archive_duration         = "auto_archive_duration"
	FlagKeyChangeLogAuditavailable                     = "available"
	FlagKeyChangeLogAuditavatar_hash                   = "avatar_hash"
	FlagKeyChangeLogAuditbanner_hash                   = "banner_hash"
	FlagKeyChangeLogAuditbitrate                       = "bitrate"
	FlagKeyChangeLogAuditchannel_id                    = "channel_id"
	FlagKeyChangeLogAuditcode                          = "code"
	FlagKeyChangeLogAuditcolor                         = "color"
	FlagKeyChangeLogAuditcommunication_disabled_until  = "communication_disabled_until"
	FlagKeyChangeLogAuditdeaf                          = "deaf"
	FlagKeyChangeLogAuditdefault_auto_archive_duration = "default_auto_archive_duration"
	FlagKeyChangeLogAuditdefault_message_notifications = "default_message_notifications"
	FlagKeyChangeLogAuditdeny                          = "deny"
	FlagKeyChangeLogAuditdescription                   = "description"
	FlagKeyChangeLogAuditdiscovery_splash_hash         = "discovery_splash_hash"
	FlagKeyChangeLogAuditenable_emoticons              = "enable_emoticons"
	FlagKeyChangeLogAuditentity_type                   = "entity_type"
	FlagKeyChangeLogAuditexpire_behavior               = "expire_behavior"
	FlagKeyChangeLogAuditexpire_grace_period           = "expire_grace_period"
	FlagKeyChangeLogAuditexplicit_content_filter       = "explicit_content_filter"
	FlagKeyChangeLogAuditformat_type                   = "format_type"
	FlagKeyChangeLogAuditguild_id                      = "guild_id"
	FlagKeyChangeLogAudithoist                         = "hoist"
	FlagKeyChangeLogAuditicon_hash                     = "icon_hash"
	FlagKeyChangeLogAuditimage_hash                    = "image_hash"
	FlagKeyChangeLogAuditid                            = "id"
	FlagKeyChangeLogAuditinvitable                     = "invitable"
	FlagKeyChangeLogAuditinviter_id                    = "inviter_id"
	FlagKeyChangeLogAuditlocation                      = "location"
	FlagKeyChangeLogAuditlocked                        = "locked"
	FlagKeyChangeLogAuditmax_age                       = "max_age"
	FlagKeyChangeLogAuditmax_uses                      = "max_uses"
	FlagKeyChangeLogAuditmentionable                   = "mentionable"
	FlagKeyChangeLogAuditmfa_level                     = "mfa_level"
	FlagKeyChangeLogAuditmute                          = "mute"
	FlagKeyChangeLogAuditname                          = "name"
	FlagKeyChangeLogAuditnick                          = "nick"
	FlagKeyChangeLogAuditnsfw                          = "nsfw"
	FlagKeyChangeLogAuditowner_id                      = "owner_id"
	FlagKeyChangeLogAuditpermission_overwrites         = "permission_overwrites"
	FlagKeyChangeLogAuditpermissions                   = "permissions"
	FlagKeyChangeLogAuditposition                      = "position"
	FlagKeyChangeLogAuditpreferred_locale              = "preferred_locale"
	FlagKeyChangeLogAuditprivacy_level                 = "privacy_level"
	FlagKeyChangeLogAuditprune_delete_days             = "prune_delete_days"
	FlagKeyChangeLogAuditpublic_updates_channel_id     = "public_updates_channel_id"
	FlagKeyChangeLogAuditrate_limit_per_user           = "rate_limit_per_user"
	FlagKeyChangeLogAuditregion                        = "region"
	FlagKeyChangeLogAuditrules_channel_id              = "rules_channel_id"
	FlagKeyChangeLogAuditsplash_hash                   = "splash_hash"
	FlagKeyChangeLogAuditstatus                        = "status"
	FlagKeyChangeLogAuditsystem_channel_id             = "system_channel_id"
	FlagKeyChangeLogAudittags                          = "tags"
	FlagKeyChangeLogAudittemporary                     = "temporary"
	FlagKeyChangeLogAudittopic                         = "topic"
	FlagKeyChangeLogAudittype                          = "type"
	FlagKeyChangeLogAuditunicode_emoji                 = "unicode_emoji"
	FlagKeyChangeLogAudituser_limit                    = "user_limit"
	FlagKeyChangeLogAudituses                          = "uses"
	FlagKeyChangeLogAuditvanity_url_code               = "vanity_url_code"
	FlagKeyChangeLogAuditverification_level            = "verification_level"
	FlagKeyChangeLogAuditwidget_channel_id             = "widget_channel_id"
	FlagKeyChangeLogAuditwidget_enabled                = "widget_enabled"
	FlagKeyChangeLogAuditadd                           = "add"
	FlagKeyChangeLogAuditremove                        = "remove"
)

// Channel Object
// https://discord.com/developers/docs/resources/channel
type Channel struct {
	ID                         Snowflake             `json:"id,omitempty"`
	Type                       *Flag                 `json:"type,omitempty"`
	GuildID                    Snowflake             `json:"guild_id,omitempty"`
	Position                   int                   `json:"position,omitempty"`
	PermissionOverwrites       []PermissionOverwrite `json:"permission_overwrites,omitempty"`
	Name                       string                `json:"name,omitempty"`
	Topic                      *string               `json:"topic,omitempty"`
	NSFW                       bool                  `json:"nsfw,omitempty"`
	LastMessageID              Snowflake             `json:"last_message_id,omitempty"`
	Bitrate                    Flag                  `json:"bitrate,omitempty"`
	UserLimit                  Flag                  `json:"user_limit,omitempty"`
	RateLimitPerUser           *CodeFlag             `json:"rate_limit_per_user,omitempty"`
	Recipients                 []*User               `json:"recipients,omitempty"`
	Icon                       string                `json:"icon,omitempty"`
	OwnerID                    Snowflake             `json:"owner_id,omitempty"`
	ApplicationID              Snowflake             `json:"application_id,omitempty"`
	ParentID                   Snowflake             `json:"parent_id,omitempty"`
	LastPinTimestamp           time.Time             `json:"last_pin_timestamp,omitempty"`
	RTCRegion                  string                `json:"rtc_region,omitempty"`
	MessageCount               Flag                  `json:"message_count,omitempty"`
	MemberCount                Flag                  `json:"member_count,omitempty"`
	ThreadMetadata             *ThreadMetadata       `json:"thread_metadata,omitempty"`
	Member                     *ThreadMember         `json:"member,omitempty"`
	DefaultAutoArchiveDuration CodeFlag              `json:"default_auto_archive_duration,omitempty"`
	Permissions                *string               `json:"permissions,omitempty"`
}

// Channel Types
// https://discord.com/developers/docs/resources/channel#channel-object-channel-types
const (
	FlagTypesChannelGUILD_TEXT           = 0
	FlagTypesChannelDM                   = 1
	FlagTypesChannelGUILD_VOICE          = 2
	FlagTypesChannelGROUP_DM             = 3
	FlagTypesChannelGUILD_CATEGORY       = 4
	FlagTypesChannelGUILD_NEWS           = 5
	FlagTypesChannelGUILD_NEWS_THREAD    = 10
	FlagTypesChannelGUILD_PUBLIC_THREAD  = 11
	FlagTypesChannelGUILD_PRIVATE_THREAD = 12
	FlagTypesChannelGUILD_STAGE_VOICE    = 13
	FlagTypesChannelGUILD_DIRECTORY      = 14
)

// Video Quality Modes
// https://discord.com/developers/docs/resources/channel#channel-object-video-quality-modes
const (
	FlagModesQualityVideoAUTO = 1
	FlagModesQualityVideoFULL = 2
)

// Thread Metadata Object
// https://discord.com/developers/docs/resources/channel#thread-metadata-object
type ThreadMetadata struct {
	Archived            bool      `json:"archived,omitempty"`
	AutoArchiveDuration CodeFlag  `json:"auto_archive_duration,omitempty"`
	Locked              bool      `json:"locked,omitempty"`
	Invitable           bool      `json:"invitable,omitempty"`
	CreateTimestamp     time.Time `json:"create_timestamp,omitempty"`
}

// Thread Member Object
// https://discord.com/developers/docs/resources/channel#thread-member-object
type ThreadMember struct {
	ThreadID      Snowflake `json:"id,omitempty"`
	UserID        Snowflake `json:"user_id,omitempty"`
	JoinTimestamp time.Time `json:"join_timestamp,omitempty"`
	Flags         CodeFlag  `json:"flags,omitempty"`
}

// Component Types
// https://discord.com/developers/docs/interactions/message-components#component-object-component-types
const (
	FlagTypesComponentActionRow  = 1
	FlagTypesComponentButton     = 2
	FlagTypesComponentSelectMenu = 3
	FlagTypesComponentTextInput  = 4
)

// Component Object
type Component interface{}

// https://discord.com/developers/docs/interactions/message-components#component-object
type ActionsRow struct {
	Components []Component `json:"components,omitempty"`
}

// Button Object
// https://discord.com/developers/docs/interactions/message-components#button-object
type Button struct {
	Style    Flag    `json:"style,omitempty"`
	Label    *string `json:"label,omitempty"`
	Emoji    *Emoji  `json:"emoji,omitempty"`
	CustomID string  `json:"custom_id,omitempty"`
	URL      string  `json:"url,omitempty"`
	Disabled bool    `json:"disabled,omitempty"`
}

// Button Styles
// https://discord.com/developers/docs/interactions/message-components#button-object-button-styles
const (
	FlagStylesbuttonPRIMARY   = 1
	FlagStylesbuttonBLURPLE   = 1
	FlagStylesbuttonSecondary = 2
	FlagStylesbuttonGREY      = 2
	FlagStylesbuttonSuccess   = 3
	FlagStylesbuttonGREEN     = 3
	FlagStylesbuttonDanger    = 4
	FlagStylesbuttonRED       = 4
	FlagStylesbuttonLINK      = 5
)

// Select Menu Structure
// https://discord.com/developers/docs/interactions/message-components#select-menu-object-select-menu-structure
type SelectMenu struct {
	CustomID    string             `json:"custom_id,omitempty"`
	Options     []SelectMenuOption `json:"options,omitempty"`
	Placeholder string             `json:"placeholder,omitempty"`
	MinValues   *Flag              `json:"min_values,omitempty"`
	MaxValues   Flag               `json:"max_values,omitempty"`
	Disabled    bool               `json:"disabled,omitempty"`
}

// Select Menu Option Structure
// https://discord.com/developers/docs/interactions/message-components#select-menu-object-select-option-structure
type SelectMenuOption struct {
	Label       *string `json:"label,omitempty"`
	Value       *string `json:"value,omitempty"`
	Description *string `json:"description,omitempty"`
	Emoji       Emoji   `json:"emoji,omitempty"`
	Default     bool    `json:"default,omitempty"`
}

// Text Input Structure
// https://discord.com/developers/docs/interactions/message-components#text-inputs-text-input-structure
type TextInput struct {
	CustomID    string    `json:"custom_id,omitempty"`
	Style       Flag      `json:"style,omitempty"`
	Label       *string   `json:"label,omitempty"`
	MinLength   *CodeFlag `json:"min_length,omitempty"`
	MaxLength   CodeFlag  `json:"max_length,omitempty"`
	Required    bool      `json:"required,omitempty"`
	Value       string    `json:"value,omitempty"`
	Placeholder *string   `json:"placeholder,omitempty"`
}

// TextInputStyle
// https://discord.com/developers/docs/interactions/message-components#text-inputs-text-input-styles
const (
	FlagStyleInputTextShort     = 1
	FlagStyleInputTextParagraph = 2
)

// Embed Object
// https://discord.com/developers/docs/resources/channel#embed-object
type Embed struct {
	Title       string          `json:"title,omitempty"`
	Type        string          `json:"type,omitempty"`
	Description *string         `json:"description,omitempty"`
	URL         string          `json:"url,omitempty"`
	Timestamp   time.Time       `json:"timestamp,omitempty"`
	Color       CodeFlag        `json:"color,omitempty"`
	Footer      *EmbedFooter    `json:"footer,omitempty"`
	Image       *EmbedImage     `json:"image,omitempty"`
	Thumbnail   *EmbedThumbnail `json:"thumbnail,omitempty"`
	Video       *EmbedVideo     `json:"video,omitempty"`
	Provider    *EmbedProvider  `json:"provider,omitempty"`
	Author      *EmbedAuthor    `json:"author,omitempty"`
	Fields      []*EmbedField   `json:"fields,omitempty"`
}

// Embed Types
// https://discord.com/developers/docs/resources/channel#embed-object-embed-types
var (
	EmbedTypes = map[string]string{
		"rich":    "generic embed rendered from embed attributes",
		"image":   "image embed",
		"video":   "video embed",
		"gifv":    "animated gif image embed rendered as a video embed",
		"article": "article embed",
		"link":    "link embed",
	}
)

// Embed Thumbnail Structure
// https://discord.com/developers/docs/resources/channel#embed-object-embed-thumbnail-structure
type EmbedThumbnail struct {
	URL      string  `json:"url,omitempty"`
	ProxyURL *string `json:"proxy_url,omitempty"`
	Height   int     `json:"height,omitempty"`
	Width    int     `json:"width,omitempty"`
}

// Embed Video Structure
// https://discord.com/developers/docs/resources/channel#embed-object-embed-video-structure
type EmbedVideo struct {
	URL      string  `json:"url,omitempty"`
	ProxyURL *string `json:"proxy_url,omitempty"`
	Height   int     `json:"height,omitempty"`
	Width    int     `json:"width,omitempty"`
}

// Embed Image Structure
// https://discord.com/developers/docs/resources/channel#embed-object-embed-image-structure
type EmbedImage struct {
	URL      string  `json:"url,omitempty"`
	ProxyURL *string `json:"proxy_url,omitempty"`
	Height   int     `json:"height,omitempty"`
	Width    int     `json:"width,omitempty"`
}

// Embed Provider Structure
// https://discord.com/developers/docs/resources/channel#embed-object-embed-provider-structure
type EmbedProvider struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

// Embed Author Structure
// https://discord.com/developers/docs/resources/channel#embed-object-embed-author-structure
type EmbedAuthor struct {
	Name         string `json:"name,omitempty"`
	URL          string `json:"url,omitempty"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

// Embed Footer Structure
// https://discord.com/developers/docs/resources/channel#embed-object-embed-footer-structure
type EmbedFooter struct {
	Text         *string `json:"text,omitempty"`
	IconURL      string  `json:"icon_url,omitempty"`
	ProxyIconURL string  `json:"proxy_icon_url,omitempty"`
}

// Embed Field Structure
// https://discord.com/developers/docs/resources/channel#embed-object-embed-field-structure
type EmbedField struct {
	Name   string `json:"name,omitempty"`
	Value  string `json:"value,omitempty"`
	Inline bool   `json:"inline,omitempty"`
}

// Embed Limits
// https://discord.com/developers/docs/resources/channel#embed-object-embed-limits
const (
	FlagLimitsEmbedTitle            = 256
	FlagLimitsEmbedDescription      = 4096
	FlagLimitsEmbedEmbedLimitFields = 25
	FlagLimitsEmbedFieldName        = 256
	FlagLimitsEmbedFieldValue       = 1024
	FlagLimitsEmbedFooterText       = 2048
	FlagLimitsEmbedAuthorName       = 256
)

// Emoji Object
// https://discord.com/developers/docs/resources/emoji#emoji-object-emoji-structure
type Emoji struct {
	ID            Snowflake   `json:"id,omitempty"`
	Name          *string     `json:"name,omitempty"`
	Roles         []Snowflake `json:"roles,omitempty"`
	User          *User       `json:"user,omitempty"`
	RequireColons bool        `json:"require_colons,omitempty"`
	Managed       bool        `json:"managed,omitempty"`
	Animated      bool        `json:"animated,omitempty"`
	Available     bool        `json:"available,omitempty"`
}

// Reaction Object
// https://discord.com/developers/docs/resources/channel#reaction-object
type Reaction struct {
	Count CodeFlag `json:"count,omitempty"`
	Me    bool     `json:"me,omitempty"`
	Emoji *Emoji   `json:"emoji,omitempty"`
}

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
	Name          string              `json:"name,omitempty"`
	Type          *Flag               `json:"type,omitempty"`
	URL           string              `json:"url,omitempty"`
	CreatedAt     int                 `json:"created_at,omitempty"`
	Timestamps    *ActivityTimestamps `json:"timestamps,omitempty"`
	ApplicationID Snowflake           `json:"application_id,omitempty"`
	Details       string              `json:"details,omitempty"`
	State         string              `json:"state,omitempty"`
	Emoji         *Emoji              `json:"emoji,omitempty"`
	Party         *ActivityParty      `json:"party,omitempty"`
	Assets        *ActivityAssets     `json:"assets,omitempty"`
	Secrets       *ActivitySecrets    `json:"secrets,omitempty"`
	Instance      bool                `json:"instance,omitempty"`
	Flags         BitFlag             `json:"flags,omitempty"`
	Buttons       []Button            `json:"buttons,omitempty"`
}

// Activity Types
// https://discord.com/developers/docs/topics/gateway#activity-object-activity-types
const (
	FlagEnumTypeActivityPlaying   = 0
	FlagEnumTypeActivityStreaming = 1
	FlagEnumTypeActivityListening = 2
	FlagEnumTypeActivityWatching  = 3
	FlagEnumTypeActivityCustom    = 4
	FlagEnumTypeActivityCompeting = 5
)

// Activity Timestamps Struct
// htthttps://discord.com/developers/docs/topics/gateway#activity-object-activity-timestamps
type ActivityTimestamps struct {
	Start int `json:"start,omitempty"`
	End   int `json:"end,omitempty"`
}

// Activity Emoji
// https://discord.com/developers/docs/topics/gateway#activity-object-activity-emoji
type ActivityEmoji struct {
	Name     string    `json:"name,omitempty"`
	ID       Snowflake `json:"id,omitempty"`
	Animated bool      `json:"animated,omitempty"`
}

// Activity Party Struct
// https://discord.com/developers/docs/topics/gateway#activity-object-activity-party
type ActivityParty struct {
	ID   string  `json:"id,omitempty"`
	Size *[2]int `json:"size,omitempty"`
}

// Activity Assets Struct
// https://discord.com/developers/docs/topics/gateway#activity-object-activity-assets
type ActivityAssets struct {
	LargeImage string `json:"large_image,omitempty"`
	LargeText  string `json:"large_text,omitempty"`
	SmallImage string `json:"small_image,omitempty"`
	SmallText  string `json:"small_text,omitempty"`
}

// Activity Asset Image
// https://discord.com/developers/docs/topics/gateway#activity-object-activity-asset-image
type ActivityAssetImage struct {
	ApplicationAsset string `json:"application_asset_id,omitempty"`
	MediaProxyImage  string `json:"image_id,omitempty"`
}

// Activity Secrets Struct
// https://discord.com/developers/docs/topics/gateway#activity-object-activity-secrets
type ActivitySecrets struct {
	Join     string `json:"join,omitempty"`
	Spectate string `json:"spectate,omitempty"`
	Match    string `json:"match,omitempty"`
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

// Guild Object
// https://discord.com/developers/docs/resources/guild#guild-object
type Guild struct {
	ID                          Snowflake              `json:"id,omitempty"`
	Name                        string                 `json:"name,omitempty"`
	Icon                        string                 `json:"icon,omitempty"`
	Splash                      string                 `json:"splash,omitempty"`
	DiscoverySplash             string                 `json:"discovery_splash,omitempty"`
	Owner                       bool                   `json:"owner,omitempty"`
	OwnerID                     Snowflake              `json:"owner_id,omitempty"`
	Permissions                 *string                `json:"permissions,omitempty"`
	Region                      string                 `json:"region,omitempty"`
	AfkChannelID                Snowflake              `json:"afk_channel_id,omitempty"`
	AfkTimeout                  *uint                  `json:"afk_timeout,omitempty"`
	WidgetEnabled               bool                   `json:"widget_enabled,omitempty"`
	WidgetChannelID             Snowflake              `json:"widget_channel_id,omitempty"`
	VerificationLevel           *Flag                  `json:"verification_level,omitempty"`
	DefaultMessageNotifications *Flag                  `json:"default_message_notifications,omitempty"`
	ExplicitContentFilter       *Flag                  `json:"explicit_content_filter,omitempty"`
	Roles                       []*Role                `json:"roles,omitempty"`
	Emojis                      []*Emoji               `json:"emojis,omitempty"`
	Features                    []*string              `json:"features,omitempty"`
	MFALevel                    *Flag                  `json:"mfa_level,omitempty"`
	ApplicationID               Snowflake              `json:"application_id,omitempty"`
	SystemChannelID             Snowflake              `json:"system_channel_id,omitempty"`
	SystemChannelFlags          BitFlag                `json:"system_channel_flags,omitempty"`
	RulesChannelID              Snowflake              `json:"rules_channel_id,omitempty"`
	JoinedAt                    time.Time              `json:"joined_at,omitempty"`
	Large                       bool                   `json:"large,omitempty"`
	Unavailable                 bool                   `json:"unavailable,omitempty"`
	MemberCount                 uint                   `json:"member_count,omitempty"`
	VoiceStates                 []*VoiceState          `json:"voice_states,omitempty"`
	Members                     []*GuildMember         `json:"members,omitempty"`
	Channels                    []*Channel             `json:"channels,omitempty"`
	Threads                     []*Channel             `json:"threads,omitempty"`
	Presences                   []*PresenceUpdate      `json:"presences,omitempty"`
	MaxPresences                CodeFlag               `json:"max_presences,omitempty"`
	MaxMembers                  int                    `json:"max_members,omitempty"`
	VanityUrl                   *string                `json:"vanity_url_code,omitempty"`
	Description                 *string                `json:"description,omitempty"`
	Banner                      string                 `json:"banner,omitempty"`
	PremiumTier                 *Flag                  `json:"premium_tier,omitempty"`
	PremiumSubscriptionCount    *CodeFlag              `json:"premium_subscription_count,omitempty"`
	PreferredLocale             string                 `json:"preferred_locale,omitempty"`
	PublicUpdatesChannelID      Snowflake              `json:"public_updates_channel_id,omitempty"`
	MaxVideoChannelUsers        int                    `json:"max_video_channel_users,omitempty"`
	ApproximateMemberCount      int                    `json:"approximate_member_count,omitempty"`
	ApproximatePresenceCount    int                    `json:"approximate_presence_count,omitempty"`
	WelcomeScreen               *WelcomeScreen         `json:"welcome_screen,omitempty"`
	NSFWLevel                   *Flag                  `json:"nsfw_level,omitempty"`
	StageInstances              []*StageInstance       `json:"stage_instances,omitempty"`
	Stickers                    []*Sticker             `json:"stickers,omitempty"`
	GuildScheduledEvents        []*GuildScheduledEvent `json:"guild_scheduled_events,omitempty"`
	PremiumProgressBarEnabled   bool                   `json:"premium_progress_bar_enabled,omitempty"`
}

// Default Message Notification Level
// https://discord.com/developers/docs/resources/guild#guild-object-default-message-notification-level
const (
	FlagLevelNotificationMessageDefaultALL_MESSAGES  = 0
	FlagLevelNotificationMessageDefaultONLY_MENTIONS = 1
)

// Explicit Content Filter Level
// https://discord.com/developers/docs/resources/guild#guild-object-explicit-content-filter-level
const (
	FlagLevelFilterContentExplicitDISABLED              = 0
	FlagLevelFilterContentExplicitMEMBERS_WITHOUT_ROLES = 1
	FlagLevelFilterContentExplicitALL_MEMBERS           = 2
)

// MFA Level
// https://discord.com/developers/docs/resources/guild#guild-object-mfa-level
const (
	FlagLevelMFANONE     = 0
	FlagLevelMFAELEVATED = 1
)

// Verification Level
// https://discord.com/developers/docs/resources/guild#guild-object-verification-level
const (
	FlagLevelVerificationNONE      = 0
	FlagLevelVerificationLOW       = 1
	FlagLevelVerificationMEDIUM    = 2
	FlagLevelVerificationHIGH      = 3
	FlagLevelVerificationVERY_HIGH = 4
)

// Guild NSFW Level
// https://discord.com/developers/docs/resources/guild#guild-object-guild-nsfw-level
const (
	FlagLevelNSFWGuildDEFAULT        = 0
	FlagLevelNSFWGuildEXPLICIT       = 1
	FlagLevelNSFWGuildSAFE           = 2
	FlagLevelNSFWGuildAGE_RESTRICTED = 3
)

// Premium Tier
// https://discord.com/developers/docs/resources/guild#guild-object-premium-tier
const (
	FlagTierPremiumNONE  = 0
	FlagTierPremiumONE   = 1
	FlagTierPremiumTWO   = 2
	FlagTierPremiumTHREE = 3
)

// System Channel Flags
// https://discord.com/developers/docs/resources/guild#guild-object-system-channel-flags
const (
	FlagFlagsChannelSystemSUPPRESS_JOIN_NOTIFICATIONS           = 1 << 0
	FlagFlagsChannelSystemSUPPRESS_PREMIUM_SUBSCRIPTIONS        = 1 << 1
	FlagFlagsChannelSystemSUPPRESS_GUILD_REMINDER_NOTIFICATIONS = 1 << 2
	FlagFlagsChannelSystemSUPPRESS_JOIN_NOTIFICATION_REPLIES    = 1 << 3
)

// Guild Features
// https://discord.com/developers/docs/resources/guild#guild-object-guild-features
var (
	GuildFeatures = map[string]string{
		"ANIMATED_BANNER":                  "guild has access to set an animated guild banner image",
		"ANIMATED_ICON":                    "guild has access to set an animated guild icon",
		"BANNER":                           "guild has access to set a guild banner image",
		"COMMERCE":                         "guild has access to use commerce features (i.e. create store channels)",
		"COMMUNITY":                        "guild can enable welcome screen, Membership Screening, stage channels and discovery, and receives community updates",
		"DISCOVERABLE":                     "guild is able to be discovered in the directory",
		"FEATURABLE":                       "guild is able to be featured in the directory",
		"INVITE_SPLASH":                    "guild has access to set an invite splash background",
		"MEMBER_VERIFICATION_GATE_ENABLED": "guild has enabled Membership Screening",
		"MONETIZATION_ENABLED":             "guild has enabled monetization",
		"MORE_STICKERS":                    "guild has increased custom sticker slots",
		"NEWS":                             "guild has access to create news channels",
		"PARTNERED":                        "guild is partnered",
		"PREVIEW_ENABLED":                  "guild can be previewed before joining via Membership Screening or the directory",
		"PRIVATE_THREADS":                  "guild has access to create private threads",
		"ROLE_ICONS":                       "guild is able to set role icons",
		"SEVEN_DAY_THREAD_ARCHIVE":         "guild has access to the seven day archive time for threads",
		"THREE_DAY_THREAD_ARCHIVE":         "guild has access to the three day archive time for threads",
		"TICKETED_EVENTS_ENABLED":          "guild has enabled ticketed events",
		"VANITY_URL":                       "guild has access to set a vanity URL",
		"VERIFIED":                         "guild is verified",
		"VIP_REGIONS":                      "guild has access to set 384kbps bitrate in voice (previously VIP voice servers)",
		"WELCOME_SCREEN_ENABLED":           "guild has enabled the welcome screen",
	}
)

// Guild Preview Object
// https://discord.com/developers/docs/resources/guild#guild-preview-object-guild-preview-structure
type GuildPreview struct {
	ID                       string     `json:"id,omitempty"`
	Name                     string     `json:"name,omitempty"`
	Icon                     string     `json:"icon,omitempty"`
	Splash                   string     `json:"splash,omitempty"`
	DiscoverySplash          string     `json:"discovery_splash,omitempty"`
	Emojis                   []*Emoji   `json:"emojis,omitempty"`
	Features                 []*string  `json:"features,omitempty"`
	ApproximateMemberCount   int        `json:"approximate_member_count,omitempty"`
	ApproximatePresenceCount int        `json:"approximate_presence_count,omitempty"`
	Description              *string    `json:"description,omitempty"`
	Stickers                 []*Sticker `json:"stickers,omitempty"`
}

// Guild Widget Settings Object
// https://discord.com/developers/docs/resources/guild#guild-widget-settings-object
type GuildWidget struct {
	Enabled   bool      `json:"enabled,omitempty"`
	ChannelID Snowflake `json:"channel_id,omitempty"`
}

// Guild Member Object
// https://discord.com/developers/docs/resources/guild#guild-member-object
type GuildMember struct {
	User                       *User        `json:"user,omitempty"`
	Nick                       *string      `json:"nick,omitempty"`
	Avatar                     string       `json:"avatar,omitempty"`
	Roles                      []*Snowflake `json:"roles,omitempty"`
	GuildID                    Snowflake    `json:"guild_id,omitempty"`
	JoinedAt                   time.Time    `json:"joined_at,omitempty"`
	PremiumSince               time.Time    `json:"premium_since,omitempty"`
	Deaf                       bool         `json:"deaf,omitempty"`
	Mute                       bool         `json:"mute,omitempty"`
	Pending                    bool         `json:"pending,omitempty"`
	CommunicationDisabledUntil time.Time    `json:"communication_disabled_until,omitempty"`
	Permissions                *string      `json:"permissions,omitempty"`
}

// Guild Ban Object
// https://discord.com/developers/docs/resources/guild#ban-object
type Ban struct {
	Reason *string `json:"reason,omitempty"`
	User   *User   `json:"user,omitempty"`
}

// Guild Scheduled Event Object
// https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event-object-guild-scheduled-event-structure
type GuildScheduledEvent struct {
	ID                 Snowflake                         `json:"id,omitempty"`
	GuildID            Snowflake                         `json:"guild_id,omitempty"`
	ChannelID          Snowflake                         `json:"channel_id,omitempty"`
	CreatorID          Snowflake                         `json:"creator_id,omitempty"`
	Name               string                            `json:"name,omitempty"`
	Description        string                            `json:"description,omitempty"`
	ScheduledStartTime time.Time                         `json:"scheduled_start_time,omitempty"`
	ScheduledEndTime   time.Time                         `json:"scheduled_end_time,omitempty"`
	PrivacyLevel       Flag                              `json:"privacy_level,omitempty"`
	Status             Flag                              `json:"status,omitempty"`
	EntityType         Flag                              `json:"entity_type,omitempty"`
	EntityID           Snowflake                         `json:"entity_id,omitempty"`
	EntityMetadata     GuildScheduledEventEntityMetadata `json:"entity_metadata,omitempty"`
	Creator            *User                             `json:"creator,omitempty"`
	UserCount          CodeFlag                          `json:"user_count,omitempty"`
	Image              string                            `json:"image,omitempty"`
}

// Guild Scheduled Event Privacy Level
// https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event-object-guild-scheduled-event-privacy-level
const (
	FlagGuildScheduledEventPrivacyLevelGUILD_ONLY = 2
)

// Guild Scheduled Event Entity Types
// https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event-object-guild-scheduled-event-entity-types
const (
	FlagTypesEntityEventScheduledGuildSTAGE_INSTANCE = 1
	FlagTypesEntityEventScheduledGuildVOICE          = 2
	FlagTypesEntityEventScheduledGuildEXTERNAL       = 3
)

// Guild Scheduled Event Status
// https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event-object-guild-scheduled-event-status
const (
	FlagStatusEventScheduledGuildSCHEDULED = 1
	FlagStatusEventScheduledGuildACTIVE    = 2
	FlagStatusEventScheduledGuildCOMPLETED = 3
	FlagStatusEventScheduledGuildCANCELED  = 4
)

// Guild Scheduled Event Entity Metadata
// https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event-object-guild-scheduled-event-entity-metadata
type GuildScheduledEventEntityMetadata struct {
	// location of the event (1-100 characters)
	// required for events with 'entity_type': EXTERNAL
	Location string `json:"location,omitempty"`
}

// Guild Scheduled Event User Object
// https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event-user-object-guild-scheduled-event-user-structure
type GuildScheduledEventUser struct {
	GuildScheduledEventID Snowflake    `json:"guild_scheduled_event_id,omitempty"`
	User                  *User        `json:"user,omitempty"`
	Member                *GuildMember `json:"member,omitempty"`
}

// Guild Template Object
// https://discord.com/developers/docs/resources/guild-template#guild-template-object
type GuildTemplate struct {
	Code                  string    `json:"code,omitempty"`
	Name                  string    `json:"name,omitempty"`
	Description           *string   `json:"description,omitempty"`
	UsageCount            *int      `json:"usage_count,omitempty"`
	CreatorID             Snowflake `json:"creator_id,omitempty"`
	Creator               *User     `json:"creator,omitempty"`
	CreatedAt             time.Time `json:"created_at,omitempty"`
	UpdatedAt             time.Time `json:"updated_at,omitempty"`
	SourceGuildID         Snowflake `json:"source_guild_id,omitempty"`
	SerializedSourceGuild *Guild    `json:"serialized_source_guild,omitempty"`
	IsDirty               bool      `json:"is_dirty,omitempty"`
}

// Welcome Screen Object
// https://discord.com/developers/docs/resources/guild#welcome-screen-object-welcome-screen-structure
type WelcomeScreen struct {
	Description           *string                 `json:"description,omitempty"`
	WelcomeScreenChannels []*WelcomeScreenChannel `json:"welcome_channels,omitempty"`
}

// Welcome Screen Channel Structure
// https://discord.com/developers/docs/resources/guild#welcome-screen-object-welcome-screen-channel-structure
type WelcomeScreenChannel struct {
	ChannelID   Snowflake  `json:"channel_id,omitempty"`
	Description *string    `json:"description,omitempty"`
	EmojiID     *Snowflake `json:"emoji_id,omitempty"`
	EmojiName   *string    `json:"emoji_name,omitempty"`
}

// Integration Object
// https://discord.com/developers/docs/resources/guild#integration-object
type Integration struct {
	ID                Snowflake          `json:"id,omitempty"`
	Name              string             `json:"name,omitempty"`
	Type              string             `json:"type,omitempty"`
	Enabled           bool               `json:"enabled,omitempty"`
	Syncing           bool               `json:"syncing,omitempty"`
	RoleID            Snowflake          `json:"role_id,omitempty"`
	EnableEmoticons   bool               `json:"enable_emoticons,omitempty"`
	ExpireBehavior    *Flag              `json:"expire_behavior,omitempty"`
	ExpireGracePeriod *int               `json:"expire_grace_period,omitempty"`
	User              *User              `json:"user,omitempty"`
	Account           IntegrationAccount `json:"account,omitempty"`
	SyncedAt          time.Time          `json:"synced_at,omitempty"`
	SubscriberCount   *int               `json:"subscriber_count,omitempty"`
	Revoked           bool               `json:"revoked,omitempty"`
	Application       *Application       `json:"application,omitempty"`
}

// Integration Expire Behaviors
// https://discord.com/developers/docs/resources/guild#integration-object-integration-expire-behaviors
const (
	FlagBehaviorsExpireIntegrationREMOVEROLE = 0
	FlagBehaviorsExpireIntegrationKICK       = 1
)

// Integration Account Object
// https://discord.com/developers/docs/resources/guild#integration-account-object
type IntegrationAccount struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// Integration Application Object
// https://discord.com/developers/docs/resources/guild#integration-application-object-integration-application-structure
type IntegrationApplication struct {
	ID          Snowflake `json:"id,omitempty"`
	Name        string    `json:"name,omitempty"`
	Icon        string    `json:"icon,omitempty"`
	Description *string   `json:"description,omitempty"`
	Bot         *User     `json:"bot,omitempty"`
}

// Interaction Object
// https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-object-interaction-structure
type Interaction struct {
	ID            Snowflake       `json:"id,omitempty"`
	ApplicationID Snowflake       `json:"application_id,omitempty"`
	Type          Flag            `json:"type,omitempty"`
	Data          InteractionData `json:"data,omitempty"`
	GuildID       Snowflake       `json:"guild_id,omitempty"`
	ChannelID     Snowflake       `json:"channel_id,omitempty"`
	Member        *GuildMember    `json:"member,omitempty"`
	User          *User           `json:"user,omitempty"`
	Token         string          `json:"token,omitempty"`
	Version       Flag            `json:"version,omitempty"`
	Message       *Message        `json:"message,omitempty"`
	Locale        string          `json:"locale,omitempty"`
	GuildLocale   string          `json:"guild_locale,omitempty"`
}

// Interaction Type
// https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-object-interaction-type
const (
	FlagTypeInteractionPING                             = 1
	FlagTypeInteractionAPPLICATION_COMMAND              = 2
	FlagTypeInteractionMESSAGE_COMPONENT                = 3
	FlagTypeInteractionAPPLICATION_COMMAND_AUTOCOMPLETE = 4
	FlagTypeInteractionMODAL_SUBMIT                     = 5
)

// Interaction Data Structure
// https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-object-interaction-data-structure
type InteractionData struct {
	ID            Snowflake                                  `json:"id,omitempty"`
	Name          string                                     `json:"name,omitempty"`
	Type          Flag                                       `json:"type,omitempty"`
	Resolved      *ResolvedData                              `json:"resolved,omitempty"`
	Options       []*ApplicationCommandInteractionDataOption `json:"options,omitempty"`
	CustomID      string                                     `json:"custom_id,omitempty"`
	ComponentType Flag                                       `json:"component_type,omitempty"`
	Values        []*string                                  `json:"values,omitempty"`
	TargetID      Snowflake                                  `json:"target_id,omitempty"`
	Components    []*Component                               `json:"components,omitempty"`
}

// Resolved Data Structure
// https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-object-resolved-data-structure
type ResolvedData struct {
	Users       map[Snowflake]*User        `json:"users,omitempty"`
	Members     map[Snowflake]*GuildMember `json:"members,omitempty"`
	Roles       map[Snowflake]*Role        `json:"roles,omitempty"`
	Channels    map[Snowflake]*Channel     `json:"channels,omitempty"`
	Messages    map[Snowflake]*Message     `json:"messages,omitempty"`
	Attachments map[Snowflake]*Attachment  `json:"attachments,omitempty"`
}

// Message Interaction Structure
// https://discord.com/developers/docs/interactions/receiving-and-responding#message-interaction-object-message-interaction-structure
type MessageInteraction struct {
	ID     Snowflake    `json:"id,omitempty"`
	Type   Flag         `json:"type,omitempty"`
	Name   string       `json:"name,omitempty"`
	User   *User        `json:"user,omitempty"`
	Member *GuildMember `json:"member,omitempty"`
}

// Interaction Response Structure
// https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-response-object-interaction-response-structure
type InteractionResponse struct {
	Type Flag                     `json:"type,omitempty"`
	Data *InteractionCallbackData `json:"data,omitempty"`
}

// Interaction Callback Type
// https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-response-object-interaction-callback-type
const (
	FlagTypeCallbackInteractionPONG                                    = 1
	FlagTypeCallbackInteractionCHANNEL_MESSAGE_WITH_SOURCE             = 4
	FlagTypeCallbackInteractionDEFERRED_CHANNEL_MESSAGE_WITH_SOURCE    = 5
	FlagTypeCallbackInteractionDEFERRED_UPDATE_MESSAGE                 = 6
	FlagTypeCallbackInteractionUPDATE_MESSAGE                          = 7
	FlagTypeCallbackInteractionAPPLICATION_COMMAND_AUTOCOMPLETE_RESULT = 8
	FlagTypeCallbackInteractionMODAL                                   = 9
)

// Interaction Callback Data Structure
// https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-response-object-interaction-callback-data-structure
type InteractionCallbackData interface{}

// Messages
// https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-response-object-messages
type Messages struct {
	TTS             bool             `json:"tts,omitempty"`
	Content         string           `json:"content,omitempty"`
	Embeds          []*Embed         `json:"embeds,omitempty"`
	Components      []Component      `json:"components,omitempty"`
	AllowedMentions *AllowedMentions `json:"allowed_mentions,omitempty"`
	Flags           BitFlag          `json:"flags,omitempty"`
	Attachments     []*Attachment    `json:"attachments,omitempty"`
}

// Autocomplete
// https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-response-object-autocomplete
type Autocomplete struct {
	Choices []*ApplicationCommandOptionChoice `json:"choices,omitempty"`
}

// Modal
// https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-response-object-modal
type ModalSubmitInteractionData struct {
	CustomID   *string     `json:"custom_id,omitempty"`
	Title      string      `json:"title,omitempty"`
	Components []Component `json:"components,omitempty"`
}

// Invite Object
// https://discord.com/developers/docs/resources/invite#invite-object
type Invite struct {
	Code                     string               `json:"code,omitempty"`
	Guild                    *Guild               `json:"guild,omitempty"`
	Channel                  *Channel             `json:"channel,omitempty"`
	Inviter                  *User                `json:"inviter,omitempty"`
	TargetUser               *User                `json:"target_user,omitempty"`
	TargetType               Flag                 `json:"target_type,omitempty"`
	TargetApplication        *Application         `json:"target_application,omitempty"`
	ApproximatePresenceCount int                  `json:"approximate_presence_count,omitempty"`
	ApproximateMemberCount   int                  `json:"approximate_member_count,omitempty"`
	ExpiresAt                time.Time            `json:"expires_at,omitempty"`
	StageInstance            StageInstance        `json:"stage_instance,omitempty"`
	GuildScheduledEvent      *GuildScheduledEvent `json:"guild_scheduled_event,omitempty"`
}

// Invite Target Types
// https://discord.com/developers/docs/resources/invite#invite-object-invite-target-types
const (
	FlagTypesTargetInviteSTREAM               = 1
	FlagTypesTargetInviteEMBEDDED_APPLICATION = 2
)

// Invite Metadata Object
// https://discord.com/developers/docs/resources/invite#invite-metadata-object-invite-metadata-structure
type InviteMetadata struct {
	Uses      *int      `json:"uses,omitempty"`
	MaxUses   *int      `json:"max_uses,omitempty"`
	MaxAge    int       `json:"max_age,omitempty"`
	Temporary bool      `json:"temporary,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

// Channel Mention Object
// https://discord.com/developers/docs/resources/channel#channel-mention-object
type ChannelMention struct {
	ID      Snowflake `json:"id,omitempty"`
	GuildID Snowflake `json:"guild_id,omitempty"`
	Type    *Flag     `json:"type,omitempty"`
	Name    string    `json:"name,omitempty"`
}

// Allowed Mentions Structure
// https://discord.com/developers/docs/resources/channel#allowed-mentions-object-allowed-mentions-structure
type AllowedMentions struct {
	Parse       []*string    `json:"parse,omitempty"`
	Roles       []*Snowflake `json:"roles,omitempty"`
	Users       []*Snowflake `json:"users,omitempty"`
	RepliedUser bool         `json:"replied_user,omitempty"`
}

// Allowed Mention Types
// https://discord.com/developers/docs/resources/channel#allowed-mentions-object-allowed-mentions-structure
const (
	FlagTypesMentionAllowedRoles     = "roles"
	FlagTypesMentionAllowedsUsers    = "users"
	FlagTypesMentionAllowedsEveryone = "everyone"
)

// Message Object
// https://discord.com/developers/docs/resources/channel#message-object
type Message struct {
	ID                Snowflake         `json:"id,omitempty"`
	ChannelID         *Snowflake        `json:"channel_id,omitempty"`
	GuildID           *Snowflake        `json:"guild_id,omitempty"`
	Author            *User             `json:"author,omitempty"`
	Member            *GuildMember      `json:"member,omitempty"`
	Content           string            `json:"content,omitempty"`
	Timestamp         time.Time         `json:"timestamp,omitempty"`
	EditedTimestamp   time.Time         `json:"edited_timestamp,omitempty"`
	TTS               bool              `json:"tts,omitempty"`
	MentionEveryone   bool              `json:"mention_everyone,omitempty"`
	Mentions          []*User           `json:"mentions,omitempty"`
	MentionRoles      []*Snowflake      `json:"mention_roles,omitempty"`
	MentionChannels   []*ChannelMention `json:"mention_channels,omitempty"`
	Attachments       []*Attachment     `json:"attachments,omitempty"`
	Embeds            []*Embed          `json:"embeds,omitempty"`
	Reactions         []*Reaction       `json:"reactions,omitempty"`
	Nonce             interface{}       `json:"nonce,omitempty"`
	Pinned            bool              `json:"pinned,omitempty"`
	WebhookID         *Snowflake        `json:"webhook_id,omitempty"`
	Type              *Flag             `json:"type,omitempty"`
	Activity          MessageActivity   `json:"activity,omitempty"`
	Application       *Application      `json:"application,omitempty"`
	MessageReference  *MessageReference `json:"message_reference,omitempty"`
	Flags             CodeFlag          `json:"flags,omitempty"`
	ReferencedMessage *Message          `json:"referenced_message,omitempty"`
	Interaction       *Interaction      `json:"interaction,omitempty"`
	Thread            *Channel          `json:"thread,omitempty"`
	Components        []*Component      `json:"components,omitempty"`
	StickerItems      []*StickerItem    `json:"sticker_items,omitempty"`
}

// Message Types
// https://discord.com/developers/docs/resources/channel#message-object-message-types
const (
	FlagTypesMessageDEFAULT                                      = 0
	FlagTypesMessageRECIPIENT_ADD                                = 1
	FlagTypesMessageRECIPIENT_REMOVE                             = 2
	FlagTypesMessageCALL                                         = 3
	FlagTypesMessageCHANNEL_NAME_CHANGE                          = 4
	FlagTypesMessageCHANNEL_ICON_CHANGE                          = 5
	FlagTypesMessageCHANNEL_PINNED_MESSAGE                       = 6
	FlagTypesMessageGUILD_MEMBER_JOIN                            = 7
	FlagTypesMessageUSER_PREMIUM_GUILD_SUBSCRIPTION              = 8
	FlagTypesMessageUSER_PREMIUM_GUILD_SUBSCRIPTION_TIER_ONE     = 9
	FlagTypesMessageUSER_PREMIUM_GUILD_SUBSCRIPTION_TIER_TWO     = 10
	FlagTypesMessageUSER_PREMIUM_GUILD_SUBSCRIPTION_TIER_THREE   = 11
	FlagTypesMessageCHANNEL_FOLLOW_ADD                           = 12
	FlagTypesMessageGUILD_DISCOVERY_DISQUALIFIED                 = 14
	FlagTypesMessageGUILD_DISCOVERY_REQUALIFIED                  = 15
	FlagTypesMessageGUILD_DISCOVERY_GRACE_PERIOD_INITIAL_WARNING = 16
	FlagTypesMessageGUILD_DISCOVERY_GRACE_PERIOD_FINAL_WARNING   = 17
	FlagTypesMessageTHREAD_CREATED                               = 18
	FlagTypesMessageREPLY                                        = 19
	FlagTypesMessageCHAT_INPUT_COMMAND                           = 20
	FlagTypesMessageTHREAD_STARTER_MESSAGE                       = 21
	FlagTypesMessageGUILD_INVITE_REMINDER                        = 22
	FlagTypesMessageCONTEXT_MENU_COMMAND                         = 23
)

// Message Activity Structure
// https://discord.com/developers/docs/resources/channel#message-object-message-activity-structure
type MessageActivity struct {
	Type    int     `json:"type,omitempty"`
	PartyID *string `json:"party_id,omitempty"`
}

// Message Activity Types
// https://discord.com/developers/docs/resources/channel#message-object-message-activity-types
const (
	FlagTypesActivityMessageJOIN         = 1
	FlagTypesActivityMessageSPECTATE     = 2
	FlagTypesActivityMessageLISTEN       = 3
	FlagTypesActivityMessageJOIN_REQUEST = 5
)

// Message Flags
// https://discord.com/developers/docs/resources/channel#message-object-message-flags
const (
	FlagFlagsMessageCROSSPOSTED                            = 1 << 0
	FlagFlagsMessageIS_CROSSPOST                           = 1 << 1
	FlagFlagsMessageSUPPRESS_EMBEDS                        = 1 << 2
	FlagFlagsMessageSOURCE_MESSAGE_DELETED                 = 1 << 3
	FlagFlagsMessageURGENT                                 = 1 << 4
	FlagFlagsMessageHAS_THREAD                             = 1 << 5
	FlagFlagsMessageEPHEMERAL                              = 1 << 6
	FlagFlagsMessageLOADING                                = 1 << 7
	FlagFlagsMessageFAILED_TO_MENTION_SOME_ROLES_IN_THREAD = 1 << 8
)

// Message Reference Object
// https://discord.com/developers/docs/resources/channel#message-reference-object
type MessageReference struct {
	MessageID       Snowflake  `json:"message_id,omitempty"`
	ChannelID       *Snowflake `json:"channel_id,omitempty"`
	GuildID         *Snowflake `json:"guild_id,omitempty"`
	FailIfNotExists bool       `json:"fail_if_not_exists,omitempty"`
}

// Message Attachment Object
// https://discord.com/developers/docs/resources/channel#attachment-object
type Attachment struct {
	ID       Snowflake `json:"id,omitempty"`
	Filename string    `json:"filename,omitempty"`
	Size     uint      `json:"size,omitempty"`
	URL      string    `json:"url,omitempty"`
	ProxyURL *string   `json:"proxy_url,omitempty"`
	Height   uint      `json:"height,omitempty"`
	Width    uint      `json:"width,omitempty"`

	SpoilerTag bool `json:"-,omitempty"`
}

// Sticker Structure
// https://discord.com/developers/docs/resources/sticker#sticker-object-sticker-structure
type Sticker struct {
	ID          Snowflake  `json:"id,omitempty"`
	PackID      Snowflake  `json:"pack_id,omitempty"`
	Name        string     `json:"name,omitempty"`
	Description *string    `json:"description,omitempty"`
	Tags        *string    `json:"tags,omitempty"`
	Asset       *string    `json:"asset,omitempty"`
	Type        Flag       `json:"type,omitempty"`
	FormatType  Flag       `json:"format_type,omitempty"`
	Available   bool       `json:"available,omitempty"`
	GuildID     *Snowflake `json:"guild_id,omitempty"`
	User        *User      `json:"user,omitempty"`
	SortValue   int        `json:"sort_value,omitempty"`
}

// Sticker Types
// https://discord.com/developers/docs/resources/sticker#sticker-object-sticker-types
const (
	FlagTypesStickerSTANDARD = 1
	FlagTypesStickerGUILD    = 2
)

// Sticker Format Types
// https://discord.com/developers/docs/resources/sticker#sticker-object-sticker-format-types
const (
	FlagTypesFormatStickerPNG    = 1
	FlagTypesFormatStickerAPNG   = 2
	FlagTypesFormatStickerLOTTIE = 3
)

// Sticker Item Object
// https://discord.com/developers/docs/resources/sticker#sticker-item-object
type StickerItem struct {
	ID         Snowflake `json:"id,omitempty"`
	Name       string    `json:"name,omitempty"`
	FormatType Flag      `json:"format_type,omitempty"`
}

// Sticker Pack Object
// StickerPack represents a pack of standard stickers.
type StickerPack struct {
	ID            Snowflake `json:"id,omitempty"`
	Type          Flag      `json:"type,omitempty"`
	GuildID       Snowflake `json:"guild_id,omitempty"`
	ChannelID     Snowflake `json:"channel_id,omitempty"`
	User          *User     `json:"user,omitempty"`
	Name          string    `json:"name,omitempty"`
	Avatar        string    `json:"avatar,omitempty"`
	Token         string    `json:"token,omitempty"`
	ApplicationID Snowflake `json:"application_id,omitempty"`
	SourceGuild   *Guild    `json:"source_guild,omitempty"`
	SourceChannel *Channel  `json:"source_channel,omitempty"`
	URL           string    `json:"url,omitempty"`
}

// Webhook Object
// https://discord.com/developers/docs/resources/webhook#webhook-object
type Webhook struct {
	ID            Snowflake  `json:"id,omitempty"`
	Type          Flag       `json:"type,omitempty"`
	GuildID       *Snowflake `json:"guild_id,omitempty"`
	ChannelID     *Snowflake `json:"channel_id,omitempty"`
	User          *User      `json:"user,omitempty"`
	Name          string     `json:"name,omitempty"`
	Avatar        string     `json:"avatar,omitempty"`
	Token         string     `json:"token,omitempty"`
	ApplicationID *Snowflake `json:"application_id,omitempty"`
	SourceGuild   *Guild     `json:"source_guild,omitempty"`
	SourceChannel *Channel   `json:"source_channel,omitempty"`
	URL           string     `json:"url,omitempty"`
}

// Webhook Types
// https://discord.com/developers/docs/resources/webhook#webhook-object-webhook-types
const (
	FlagTypesWebhookINCOMING        = 1
	FlagTypesWebhookCHANNELFOLLOWER = 2
	FlagTypesWebhookAPPLICATION     = 3
)

// Bitwise Permission Flags
// https://discord.com/developers/docs/topics/permissions#permissions-bitwise-permission-flags
const (
	FlagFlagsPermissionBitwiseCREATE_INSTANT_INVITE      = 1 << 0
	FlagFlagsPermissionBitwiseKICK_MEMBERS               = 1 << 1
	FlagFlagsPermissionBitwiseBAN_MEMBERS                = 1 << 2
	FlagFlagsPermissionBitwiseADMINISTRATOR              = 1 << 3
	FlagFlagsPermissionBitwiseMANAGE_CHANNELS            = 1 << 4
	FlagFlagsPermissionBitwiseMANAGE_GUILD               = 1 << 5
	FlagFlagsPermissionBitwiseADD_REACTIONS              = 1 << 6
	FlagFlagsPermissionBitwiseVIEW_AUDIT_LOG             = 1 << 7
	FlagFlagsPermissionBitwisePRIORITY_SPEAKER           = 1 << 8
	FlagFlagsPermissionBitwiseSTREAM                     = 1 << 9
	FlagFlagsPermissionBitwiseVIEW_CHANNEL               = 1 << 10
	FlagFlagsPermissionBitwiseSEND_MESSAGES              = 1 << 11
	FlagFlagsPermissionBitwiseSEND_TTS_MESSAGES          = 1 << 12
	FlagFlagsPermissionBitwiseMANAGE_MESSAGES            = 1 << 13
	FlagFlagsPermissionBitwiseEMBED_LINKS                = 1 << 14
	FlagFlagsPermissionBitwiseATTACH_FILES               = 1 << 15
	FlagFlagsPermissionBitwiseREAD_MESSAGE_HISTORY       = 1 << 16
	FlagFlagsPermissionBitwiseMENTION_EVERYONE           = 1 << 17
	FlagFlagsPermissionBitwiseUSE_EXTERNAL_EMOJIS        = 1 << 18
	FlagFlagsPermissionBitwiseVIEW_GUILD_INSIGHTS        = 1 << 19
	FlagFlagsPermissionBitwiseCONNECT                    = 1 << 20
	FlagFlagsPermissionBitwiseSPEAK                      = 1 << 21
	FlagFlagsPermissionBitwiseMUTE_MEMBERS               = 1 << 22
	FlagFlagsPermissionBitwiseDEAFEN_MEMBERS             = 1 << 23
	FlagFlagsPermissionBitwiseMOVE_MEMBERS               = 1 << 24
	FlagFlagsPermissionBitwiseUSE_VAD                    = 1 << 25
	FlagFlagsPermissionBitwiseCHANGE_NICKNAME            = 1 << 26
	FlagFlagsPermissionBitwiseMANAGE_NICKNAMES           = 1 << 27
	FlagFlagsPermissionBitwiseMANAGE_ROLES               = 1 << 28
	FlagFlagsPermissionBitwiseMANAGE_WEBHOOKS            = 1 << 29
	FlagFlagsPermissionBitwiseMANAGE_EMOJIS_AND_STICKERS = 1 << 30
	FlagFlagsPermissionBitwiseUSE_APPLICATION_COMMANDS   = 1 << 31
	FlagFlagsPermissionBitwiseREQUEST_TO_SPEAK           = 1 << 32
	FlagFlagsPermissionBitwiseMANAGE_EVENTS              = 1 << 33
	FlagFlagsPermissionBitwiseMANAGE_THREADS             = 1 << 34
	FlagFlagsPermissionBitwiseCREATE_PUBLIC_THREADS      = 1 << 35
	FlagFlagsPermissionBitwiseCREATE_PRIVATE_THREADS     = 1 << 36
	FlagFlagsPermissionBitwiseUSE_EXTERNAL_STICKERS      = 1 << 37
	FlagFlagsPermissionBitwiseSEND_MESSAGES_IN_THREADS   = 1 << 38
	FlagFlagsPermissionBitwiseUSE_EMBEDDED_ACTIVITIES    = 1 << 39
	FlagFlagsPermissionBitwiseMODERATE_MEMBERS           = 1 << 40
)

// Overwrite Object
// https://discord.com/developers/docs/resources/channel#overwrite-object
type PermissionOverwrite struct {
	ID    Snowflake `json:"id,omitempty"`
	Type  *Flag     `json:"type,omitempty"`
	Deny  string    `json:"deny,omitempty"`
	Allow string    `json:"allow,omitempty"`
}

const (
	FlagPermissionOverwriteTypeRole   = 0
	FlagPermissionOverwriteTypeMember = 1
)

// Role Object
// https://discord.com/developers/docs/topics/permissions#role-object
type Role struct {
	ID           Snowflake `json:"id,omitempty"`
	Name         string    `json:"name,omitempty"`
	Color        uint      `json:"color,omitempty"`
	Hoist        bool      `json:"hoist,omitempty"`
	Icon         string    `json:"icon,omitempty"`
	UnicodeEmoji string    `json:"unicode_emoji,omitempty"`
	Position     int       `json:"position,omitempty"`
	Permissions  string    `json:"permissions,omitempty"`
	Managed      bool      `json:"managed,omitempty"`
	Mentionable  bool      `json:"mentionable,omitempty"`
	Tags         *RoleTags `json:"tags,omitempty"`
}

// Role Tags Structure
// https://discord.com/developers/docs/topics/permissions#role-object-role-tags-structure
type RoleTags struct {
	BotID             Snowflake `json:"bot_id,omitempty"`
	IntegrationID     Snowflake `json:"integration_id,omitempty"`
	PremiumSubscriber bool      `json:"premium_subscriber,omitempty"`
}

// Stage Instance Object
// https://discord.com/developers/docs/resources/stage-instance#stage-instance-object
type StageInstance struct {
	ID                    Snowflake  `json:"id,omitempty"`
	GuildID               *Snowflake `json:"guild_id,omitempty"`
	ChannelID             *Snowflake `json:"channel_id,omitempty"`
	Topic                 string     `json:"topic,omitempty"`
	PrivacyLevel          Flag       `json:"privacy_level,omitempty"`
	DiscoverableDisabled  bool       `json:"discoverable_disabled,omitempty"`
	GuildScheduledEventID Snowflake  `json:"guild_scheduled_event_id,omitempty"`
}

// Privacy Level
// https://discord.com/developers/docs/resources/stage-instance#stage-instance-object-privacy-level
const (
	FlagLevelPrivacyPUBLIC     = 1
	FlagLevelPrivacyGUILD_ONLY = 2
)

// Team Object
// https://discord.com/developers/docs/topics/teams#data-models-team-object
type Team struct {
	Icon        string        `json:"icon,omitempty"`
	ID          Snowflake     `json:"id,omitempty"`
	Members     []*TeamMember `json:"members,omitempty"`
	Name        string        `json:"name,omitempty"`
	Description *string       `json:"description,omitempty"`
	OwnerUserID Snowflake     `json:"owner_user_id,omitempty"`
}

// Team Member Object
// https://discord.com/developers/docs/topics/teams#data-models-team-member-object
type TeamMember struct {
	MembershipState Flag      `json:"membership_state,omitempty"`
	Permissions     []string  `json:"permissions,omitempty"`
	TeamID          Snowflake `json:"team_id,omitempty"`
	User            *User     `json:"user,omitempty"`
}

// Membership State Enum
// https://discord.com/developers/docs/topics/teams#data-models-membership-state-enum
const (
	FlagEnumStateMembershipINVITED  = 1
	FlagEnumStateMembershipACCEPTED = 2
)

// User Object
// https://discord.com/developers/docs/resources/user#user-object
type User struct {
	ID            Snowflake `json:"id,omitempty"`
	Username      string    `json:"username,omitempty"`
	Discriminator string    `json:"discriminator,omitempty"`
	Avatar        string    `json:"avatar,omitempty"`
	Bot           bool      `json:"bot,omitempty"`
	System        bool      `json:"system,omitempty"`
	MFAEnabled    bool      `json:"mfa_enabled,omitempty"`
	Banner        string    `json:"banner,omitempty"`
	AccentColor   int       `json:"accent_color,omitempty"`
	Locale        string    `json:"locale,omitempty"`
	Verified      bool      `json:"verified,omitempty"`
	Email         *string   `json:"email,omitempty"`
	Flags         *BitFlag  `json:"flag,omitempty"`
	PremiumType   *Flag     `json:"premium_type,omitempty"`
	PublicFlags   BitFlag   `json:"public_flag,omitempty"`
}

// User Flags
// https://discord.com/developers/docs/resources/user#user-object-user-flags
const (
	FlagFlagsUserNONE                         = 0
	FlagFlagsUserSTAFF                        = 1 << 0
	FlagFlagsUserPARTNER                      = 1 << 1
	FlagFlagsUserHYPESQUAD                    = 1 << 2
	FlagFlagsUserBUG_HUNTER_LEVEL_1           = 1 << 3
	FlagFlagsUserHYPESQUAD_ONLINE_HOUSE_ONE   = 1 << 6
	FlagFlagsUserHYPESQUAD_ONLINE_HOUSE_TWO   = 1 << 7
	FlagFlagsUserHYPESQUAD_ONLINE_HOUSE_THREE = 1 << 8
	FlagFlagsUserPREMIUM_EARLY_SUPPORTER      = 1 << 9
	FlagFlagsUserTEAM_PSEUDO_USER             = 1 << 10
	FlagFlagsUserBUG_HUNTER_LEVEL_2           = 1 << 14
	FlagFlagsUserVERIFIED_BOT                 = 1 << 16
	FlagFlagsUserVERIFIED_DEVELOPER           = 1 << 17
	FlagFlagsUserCERTIFIED_MODERATOR          = 1 << 18
	FlagFlagsUserBOT_HTTP_INTERACTIONS        = 1 << 19
)

// Premium Types
// https://discord.com/developers/docs/resources/user#user-object-premium-types
const (
	FlagTypesPremiumNONE         = 0
	FlagTypesPremiumNITROCLASSIC = 1
	FlagTypesPremiumNITRO        = 2
)

// User Connection Object
// https://discord.com/developers/docs/resources/user#connection-object-connection-structure
type Connection struct {
	ID           string         `json:"id,omitempty"`
	Name         string         `json:"name,omitempty"`
	Type         string         `json:"type,omitempty"`
	Revoked      bool           `json:"revoked,omitempty"`
	Integrations []*Integration `json:"integrations,omitempty"`
	Verified     bool           `json:"verified,omitempty"`
	FriendSync   bool           `json:"friend_sync,omitempty"`
	ShowActivity bool           `json:"show_activity,omitempty"`
	Visibility   Flag           `json:"visibility,omitempty"`
}

// Visibility Types
// https://discord.com/developers/docs/resources/user#connection-object-visibility-types
const (
	FlagTypesVisibilityNONE     = 0
	FlagTypesVisibilityEVERYONE = 1
)

// Voice State Object
// https://discord.com/developers/docs/resources/voice#voice-state-object-voice-state-structure
type VoiceState struct {
	GuildID                 Snowflake    `json:"guild_id,omitempty"`
	ChannelID               Snowflake    `json:"channel_id,omitempty"`
	UserID                  Snowflake    `json:"user_id,omitempty"`
	Member                  *GuildMember `json:"member,omitempty"`
	SessionID               string       `json:"session_id,omitempty"`
	Deaf                    bool         `json:"deaf,omitempty"`
	Mute                    bool         `json:"mute,omitempty"`
	SelfDeaf                bool         `json:"self_deaf,omitempty"`
	SelfMute                bool         `json:"self_mute,omitempty"`
	SelfStream              bool         `json:"self_stream,omitempty"`
	SelfVideo               bool         `json:"self_video,omitempty"`
	Suppress                bool         `json:"suppress,omitempty"`
	RequestToSpeakTimestamp time.Time    `json:"request_to_speak_timestamp,omitempty"`
}

// Voice Region Object
// https://discord.com/developers/docs/resources/voice#voice-region-object-voice-region-structure
type VoiceRegion struct {
	ID         string `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	Optimal    bool   `json:"optimal,omitempty"`
	Deprecated bool   `json:"deprecated,omitempty"`
	Custom     bool   `json:"custom,omitempty"`
}

// Get Gateway Response
// https://discord.com/developers/docs/topics/gateway#get-gateway-example-response
type GetGatewayResponse struct {
	URL string `json:"url,omitempty"`
}

// Get Gateway Bot Response
// https://discord.com/developers/docs/topics/gateway#get-gateway-example-response
type GetGatewayBotResponse struct {
	URL               string            `json:"url,omitempty"`
	Shards            *int              `json:"shards,omitempty"`
	SessionStartLimit SessionStartLimit `json:"session_start_limit"`
}

// Current Authorization Information Response Structure
// https://discord.com/developers/docs/topics/oauth2#get-current-authorization-information
type CurrentAuthorizationInformationResponse struct {
	Application *Application `json:"application"`
	Scopes      []*int       `json:"scopes"`
	Expires     *time.Time   `json:"expires"`
	User        *User        `json:"user"`
}

// List Active Threads Response Body
// https://discord.com/developers/docs/resources/channel#list-active-threads-response-body
type ListActiveThreadsResponse struct {
	Threads []*Channel      `json:"threads"`
	Members []*ThreadMember `json:"members"`
	HasMore bool            `json:"has_more"`
}

// List Public Archived Threads Response Body
// https://discord.com/developers/docs/resources/channel#list-active-threads-response-body
type ListPublicArchivedThreadsResponse struct {
	Threads []*Channel      `json:"threads"`
	Members []*ThreadMember `json:"members"`
	HasMore bool            `json:"has_more"`
}

// List Private Archived Threads Response Body
// https://discord.com/developers/docs/resources/channel#list-active-threads-response-body
type ListPrivateArchivedThreadsResponse struct {
	Threads []*Channel      `json:"threads"`
	Members []*ThreadMember `json:"members"`
	HasMore bool            `json:"has_more"`
}

// List Joined Private Archived Threads Response Body
// https://discord.com/developers/docs/resources/channel#list-active-threads-response-body
type ListJoinedPrivateArchivedThreadsResponse struct {
	Threads []*Channel      `json:"threads"`
	Members []*ThreadMember `json:"members"`
	HasMore bool            `json:"has_more"`
}

// Modify Current User Nick Response
// https://discord.com/developers/docs/topics/gateway#get-gateway-example-response
type ModifyCurrentUserNickResponse struct {
	Nick *string `json:"nick,omitempty"`
}

