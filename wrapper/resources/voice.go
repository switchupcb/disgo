package resources

import "time"

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
