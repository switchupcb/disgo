package resources

import "encoding/json"

// Gateway Payload Structure
// https://discord.com/developers/docs/topics/gateway#payloads-gateway-payload-structure
type GatewayPayload struct {
	Op             uint8           `json:"op,omitempty"`
	Data           json.RawMessage `json:"d,omitempty"`
	SequenceNumber uint32          `json:"s,omitempty"`
	EventName      string          `json:"t,omitempty"`
}

// Gateway URL Query String Params
// https://discord.com/developers/docs/topics/gateway#connecting-gateway-url-query-string-params
// TODO

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
	FlagListofIntentsGUILDS = 1 << 0

	// GUILD_MEMBER_ADD
	// GUILD_MEMBER_UPDATE
	// GUILD_MEMBER_REMOVE
	// THREAD_MEMBERS_UPDATE *
	FlagListofIntentsGUILD_MEMBERS = 1 << 1

	// GUILD_BAN_ADD
	// GUILD_BAN_REMOVE
	FlagListofIntentsGUILD_BANS = 1 << 2

	// GUILD_EMOJIS_UPDATE
	// GUILD_STICKERS_UPDATE
	FlagListofIntentsGUILD_EMOJIS_AND_STICKERS = 1 << 3

	// GUILD_INTEGRATIONS_UPDATE
	// INTEGRATION_CREATE
	// INTEGRATION_UPDATE
	// INTEGRATION_DELETE
	FlagListofIntentsGUILD_INTEGRATIONS = 1 << 4

	// WEBHOOKS_UPDATE
	FlagListofIntentsGUILD_WEBHOOKS = 1 << 5

	// INVITE_CREATE
	// INVITE_DELETE
	FlagListofIntentsGUILD_INVITES = 1 << 6

	// VOICE_STATE_UPDATE
	FlagListofIntentsGUILD_VOICE_STATES = 1 << 7

	// PRESENCE_UPDATE
	FlagListofIntentsGUILD_PRESENCES = 1 << 8

	// MESSAGE_CREATE
	// MESSAGE_UPDATE
	// MESSAGE_DELETE
	// MESSAGE_DELETE_BULK
	FlagListofIntentsGUILD_MESSAGES = 1 << 9

	// MESSAGE_REACTION_ADD
	// MESSAGE_REACTION_REMOVE
	// MESSAGE_REACTION_REMOVE_ALL
	// MESSAGE_REACTION_REMOVE_EMOJI
	FlagListofIntentsGUILD_MESSAGE_REACTIONS = 1 << 10

	// TYPING_START

	FlagListofIntentsGUILD_MESSAGE_TYPING = 1 << 11

	// MESSAGE_CREATE
	// MESSAGE_UPDATE
	// MESSAGE_DELETE
	// CHANNEL_PINS_UPDATE
	FlagListofIntentsDIRECT_MESSAGES = 1 << 12

	// MESSAGE_REACTION_ADD
	// MESSAGE_REACTION_REMOVE
	// MESSAGE_REACTION_REMOVE_ALL
	// MESSAGE_REACTION_REMOVE_EMOJI
	FlagListofIntentsDIRECT_MESSAGE_REACTIONS = 1 << 13

	// TYPING_START
	FlagListofIntentsDIRECT_MESSAGE_TYPING = 1 << 14

	// GUILD_SCHEDULED_EVENT_CREATE
	// GUILD_SCHEDULED_EVENT_UPDATE
	// GUILD_SCHEDULED_EVENT_DELETE
	// GUILD_SCHEDULED_EVENT_USER_ADD
	// GUILD_SCHEDULED_EVENT_USER_REMOVE
	FlagListofIntentsGUILD_SCHEDULED_EVENTS = 1 << 16
)

// Gateway Commands
// https://discord.com/developers/docs/topics/gateway#commands-and-events-gateway-commands
// TODO

// Gateway Events
// https://discord.com/developers/docs/topics/gateway#commands-and-events-gateway-events
// TODO

// Identify Structure
// https://discord.com/developers/docs/topics/gateway#identify-identify-structure
// TODO

// Identify Connection Properties
// https://discord.com/developers/docs/topics/gateway#identify-identify-connection-properties
// TODO

// Resume Structure
// https://discord.com/developers/docs/topics/gateway#resume-resume-structure
// TODO

// Guild Request Members Structure
// https://discord.com/developers/docs/topics/gateway#request-guild-members-guild-request-members-structure
// TODO

// Gateway Voice State Update Structure
// https://discord.com/developers/docs/topics/gateway#update-voice-state-gateway-voice-state-update-structure
// TODO

// Gateway Presence Update Structure
// https://discord.com/developers/docs/topics/gateway#update-presence-gateway-presence-update-structure
// TODO

// Status Types
// https://discord.com/developers/docs/topics/gateway#update-presence-status-types
const (
	FlagStatusTypesOnline       = "online"
	FlagStatusTypesDoNotDisturb = "dnd"
	FlagStatusTypesAFK          = "idle"
	FlagStatusTypesInvisible    = "invisible"
	FlagStatusTypesOffline      = "offline"
)

// Presence Update Event Fields
// https://discord.com/developers/docs/topics/gateway#presence-update-presence-update-event-fields
type PresenceUpdate struct {
	User         *User        `json:"user,omitempty"`
	GuildID      Snowflake    `json:"guild_id,omitempty"`
	Status       string       `json:"status,omitempty"`
	Activities   []*Activity  `json:"activities,omitempty"`
	ClientStatus ClientStatus `json:"client_status,omitempty"`
}

// Client Status Object
// https://discord.com/developers/docs/topics/gateway#client-status-object
type ClientStatus struct {
	Desktop string `json:"desktop,omitempty"`
	Mobile  string `json:"mobile,omitempty"`
	Web     string `json:"web,omitempty"`
}
