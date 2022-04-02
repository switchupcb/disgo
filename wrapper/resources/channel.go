package resource

import (
	"time"
)

// Channel Object
// https://discord.com/developers/docs/resources/channel
type Channel struct {
	ID                         int64                 `json:"id"`
	Type                       uint8                 `json:"type"`
	GuildID                    int64                 `json:"guild_id,omitempty"`
	Position                   int                   `json:"position,omitempty"` // can be less than 0
	PermissionOverwrites       []PermissionOverwrite `json:"permission_overwrites,omitempty"`
	Name                       string                `json:"name,omitempty"`
	Topic                      string                `json:"topic,omitempty"`
	NSFW                       bool                  `json:"nsfw,omitempty"`
	LastMessageID              int64                 `json:"last_message_id,omitempty"`
	Bitrate                    uint                  `json:"bitrate,omitempty"`
	UserLimit                  uint                  `json:"user_limit,omitempty"`
	RateLimitPerUser           uint                  `json:"rate_limit_per_user,omitempty"`
	Recipients                 []*User               `json:"recipients,omitempty"` // empty if not DM/GroupDM
	Icon                       string                `json:"icon,omitempty"`
	OwnerID                    int64                 `json:"owner_id,omitempty"`
	ApplicationID              int64                 `json:"application_id,omitempty"`
	ParentID                   int64                 `json:"parent_id,omitempty"`
	LastPinTimestamp           time.Time             `json:"last_pin_timestamp,omitempty"`
	RTCRegion                  string                `json:"rtc_region,omitempty"`
	MessageCount               int                   `json:"message_count,omitempty"`
	MemberCount                int                   `json:"member_count,omitempty"`
	ThreadMetadata             *ThreadMetadata       `json:"thread_metadata,omitempty"`
	Member                     *ThreadMember         `json:"member,omitempty"`
	DefaultAutoArchiveDuration int                   `json:"default_auto_archive_duration,omitempty"`
	Permissions                string                `json:"permissions,omitempty"`
}

// Channel Types
// https://discord.com/developers/docs/resources/channel#channel-object-channel-types
const (
	GUILD_TEXT           = 0
	DM                   = 1
	GUILD_VOICE          = 2
	GROUP_DM             = 3
	GUILD_CATEGORY       = 4
	GUILD_NEWS           = 5
	GUILD_NEWS_THREAD    = 10
	GUILD_PUBLIC_THREAD  = 11
	GUILD_PRIVATE_THREAD = 12
	GUILD_STAGE_VOICE    = 13
	GUILD_DIRECTORY      = 14
)

// Video Quality Modes
// https://discord.com/developers/docs/resources/channel#channel-object-video-quality-modes
const (
	AUTO = 1
	FULL = 2
)

// Thread Metadata Object
// https://discord.com/developers/docs/resources/channel#thread-metadata-object
type ThreadMetadata struct {
	Archived            bool      `json:"archived,omitempty"`
	AutoArchiveDuration int       `json:"auto_archive_duration,omitempty"`
	Locked              bool      `json:"locked,bool"`
	Invitable           bool      `json:"invitable,omitempty"`
	CreateTimestamp     time.Time `json:"create_timestamp,omitempty"`
}
