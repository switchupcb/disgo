package resource

import "encoding/json"

// Gateway Payload Structure
// https://discord.com/developers/docs/topics/gateway#payloads-gateway-payload-structure
type GatewayPayload struct {
	Op             uint8           `json:"op"`
	Data           json.RawMessage `json:"d"`
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
	GUILDS = 1 << 0

	// GUILD_MEMBER_ADD
	// GUILD_MEMBER_UPDATE
	// GUILD_MEMBER_REMOVE
	// THREAD_MEMBERS_UPDATE *
	GUILD_MEMBERS = 1 << 1

	// GUILD_BAN_ADD
	// GUILD_BAN_REMOVE
	GUILD_BANS = 1 << 2

	// GUILD_EMOJIS_UPDATE
	// GUILD_STICKERS_UPDATE
	GUILD_EMOJIS_AND_STICKERS = 1 << 3

	// GUILD_INTEGRATIONS_UPDATE
	// INTEGRATION_CREATE
	// INTEGRATION_UPDATE
	// INTEGRATION_DELETE
	GUILD_INTEGRATIONS = 1 << 4

	// WEBHOOKS_UPDATE
	GUILD_WEBHOOKS = 1 << 5

	// INVITE_CREATE
	// INVITE_DELETE
	GUILD_INVITES = 1 << 6

	// VOICE_STATE_UPDATE
	GUILD_VOICE_STATES = 1 << 7

	// PRESENCE_UPDATE
	GUILD_PRESENCES = 1 << 8

	// MESSAGE_CREATE
	// MESSAGE_UPDATE
	// MESSAGE_DELETE
	// MESSAGE_DELETE_BULK
	GUILD_MESSAGES = 1 << 9

	// MESSAGE_REACTION_ADD
	// MESSAGE_REACTION_REMOVE
	// MESSAGE_REACTION_REMOVE_ALL
	// MESSAGE_REACTION_REMOVE_EMOJI
	GUILD_MESSAGE_REACTIONS = 1 << 10

	// TYPING_START

	GUILD_MESSAGE_TYPING = 1 << 11

	// MESSAGE_CREATE
	// MESSAGE_UPDATE
	// MESSAGE_DELETE
	// CHANNEL_PINS_UPDATE
	DIRECT_MESSAGES = 1 << 12

	// MESSAGE_REACTION_ADD
	// MESSAGE_REACTION_REMOVE
	// MESSAGE_REACTION_REMOVE_ALL
	// MESSAGE_REACTION_REMOVE_EMOJI
	DIRECT_MESSAGE_REACTIONS = 1 << 13

	// TYPING_START
	DIRECT_MESSAGE_TYPING = 1 << 14

	// GUILD_SCHEDULED_EVENT_CREATE
	// GUILD_SCHEDULED_EVENT_UPDATE
	// GUILD_SCHEDULED_EVENT_DELETE
	// GUILD_SCHEDULED_EVENT_USER_ADD
	// GUILD_SCHEDULED_EVENT_USER_REMOVE
	GUILD_SCHEDULED_EVENTS = 1 << 16
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
	Online       = "online"
	DoNotDisturb = "dnd"
	AFK          = "idle"
	Invisible    = "invisible"
	Offline      = "offline"
)

// Presence Update Event Fields
// https://discord.com/developers/docs/topics/gateway#presence-update-presence-update-event-fields
type PresenceUpdate struct {
	User         *User        `json:"user"`
	GuildID      int64        `json:"guild_id"`
	Status       string       `json:"status"`
	Activities   []*Activity  `json:"activities"`
	ClientStatus ClientStatus `json:"client_status"`
}

// Client Status Object
// https://discord.com/developers/docs/topics/gateway#client-status-object
type ClientStatus struct {
	Desktop string `json:"desktop"`
	Mobile  string `json:"mobile"`
	Web     string `json:"web"`
}

// Activity Object
// https://discord.com/developers/docs/topics/gateway#activity-object-activity-structure
type Activity struct {
	Name          string      `json:"name"`
	Type          uint8       `json:"type"`
	URL           string      `json:"url,omitempty"`
	CreatedAt     int         `json:"created_at"`
	Timestamps    *Timestamps `json:"timestamps,omitempty"`
	ApplicationID int64       `json:"application_id,omitempty"`
	Details       string      `json:"details,omitempty"`
	State         string      `json:"state,omitempty"`
	Emoji         *Emoji      `json:"emoji,omitempty"`
	Party         *Party      `json:"party,omitempty"`
	Assets        *Assets     `json:"assets,omitempty"`
	Secrets       *Secrets    `json:"secrets,omitempty"`
	Instance      bool        `json:"instance,omitempty"`
	Flags         uint8       `json:"flags,omitempty"`
	Buttons       []Button    `json:"buttons,omitempty"`
}
