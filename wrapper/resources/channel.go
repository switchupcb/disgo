package resources

import (
	"time"
)

// Channel Object
// https://discord.com/developers/docs/resources/channel
type Channel struct {
	ID                         Snowflake             `json:"id,omitempty"`
	Type                       Flag                  `json:"type,omitempty"`
	GuildID                    Snowflake             `json:"guild_id,omitempty"`
	Position                   int                   `json:"position,omitempty"`
	PermissionOverwrites       []PermissionOverwrite `json:"permission_overwrites,omitempty"`
	Name                       string                `json:"name,omitempty"`
	Topic                      string                `json:"topic,omitempty"`
	NSFW                       bool                  `json:"nsfw,omitempty"`
	LastMessageID              Snowflake             `json:"last_message_id,omitempty"`
	Bitrate                    Flag                  `json:"bitrate,omitempty"`
	UserLimit                  Flag                  `json:"user_limit,omitempty"`
	RateLimitPerUser           CodeFlag              `json:"rate_limit_per_user,omitempty"`
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
	Permissions                string                `json:"permissions,omitempty"`
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
