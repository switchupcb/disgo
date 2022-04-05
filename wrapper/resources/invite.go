package resources

import "time"

// Invite Object
// https://discord.com/developers/docs/resources/invite#invite-object
type Invite struct {
	Code                     string              `json:"code,omitempty"`
	Guild                    *Guild              `json:"guild,omitempty"`
	Channel                  *Channel            `json:"channel,omitempty"`
	Inviter                  *User               `json:"inviter,omitempty"`
	TargetUser               *User               `json:"target_user,omitempty"`
	TargetType               Flag                `json:"target_type,omitempty"`
	TargetApplication        *Application        `json:"target_application,omitempty"`
	ApproximatePresenceCount int                 `json:"approximate_presence_count,omitempty"`
	ApproximateMemberCount   int                 `json:"approximate_member_count,omitempty"`
	ExpiresAt                time.Time           `json:"expires_at,omitempty"`
	StageInstance            StageInstance       `json:"stage_instance,omitempty"`
	GuildScheduledEvent      GuildScheduledEvent `json:"guild_scheduled_event,omitempty"`
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
	Uses      int       `json:"uses,omitempty"`
	MaxUses   int       `json:"max_uses,omitempty"`
	MaxAge    int       `json:"max_age,omitempty"`
	Temporary bool      `json:"temporary,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}
