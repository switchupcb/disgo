package resource

import "time"

// Invite Object
// https://discord.com/developers/docs/resources/invite#invite-object
type Invite struct {
	Code                     string              `json:"code"`
	Guild                    *Guild              `json:"guild"`
	Channel                  *Channel            `json:"channel"`
	Inviter                  *User               `json:"inviter"`
	TargetUser               *User               `json:"target_user"`
	TargetType               uint8               `json:"target_type"`
	TargetApplication        *Application        `json:"target_application"`
	ApproximatePresenceCount int                 `json:"approximate_presence_count"`
	ApproximateMemberCount   int                 `json:"approximate_member_count"`
	ExpiresAt                time.Time           `json:"expires_at"`
	StageInstance            StageInstance       `json:"stage_instance"`
	GuildScheduledEvent      GuildScheduledEvent `json:"guild_scheduled_event"`
}

// Invite Target Types
// https://discord.com/developers/docs/resources/invite#invite-object-invite-target-types
const (
	InviteTargetTypesSTREAM = 1
	EMBEDDED_APPLICATION    = 2
)

// Invite Metadata Object
// https://discord.com/developers/docs/resources/invite#invite-metadata-object-invite-metadata-structure
type InviteMetadata struct {
	Uses      int       `json:"uses"`
	MaxUses   int       `json:"max_uses"`
	MaxAge    int       `json:"max_age"`
	Temporary bool      `json:"temporary"`
	CreatedAt time.Time `json:"created_at"`
}
