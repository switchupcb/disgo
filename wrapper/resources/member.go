// Thread Member Object
package resource

import "time"

// Thread Member Object
// https://discord.com/developers/docs/resources/channel#thread-member-object
type ThreadMember struct {
	ThreadID      int64     `json:"id,omitempty"`
	UserID        int64     `json:"user_id,omitempty"`
	JoinTimestamp time.Time `json:"join_timestamp"`
	Flags         int       `json:"flags"`
}

// Guild Member Object
// https://discord.com/developers/docs/resources/guild#guild-member-object
type GuildMember struct {
	User                       *User      `json:"user"`
	Nick                       string     `json:"nick,omitempty"`
	Avatar                     string     `json:"avatar"`
	Roles                      []int64    `json:"roles"`
	GuildID                    int64      `json:"guild_id,omitempty"`
	JoinedAt                   time.Time  `json:"joined_at,omitempty"`
	PremiumSince               time.Time  `json:"premium_since,omitempty"`
	Deaf                       bool       `json:"deaf"`
	Mute                       bool       `json:"mute"`
	Pending                    bool       `json:"pending"`
	CommunicationDisabledUntil *time.Time `json:"communication_disabled_until"`
	Permissions                int64      `json:"permissions,string"`
}
