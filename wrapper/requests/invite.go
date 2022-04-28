package requests

import "github.com/switchupcb/disgo/wrapper/resources"

// Get Invite
// GET /invites/{invite.code}
// https://discord.com/developers/docs/resources/invite#get-invite
type GetInvite struct {
	GuildScheduledEventID resources.Snowflake `json:"guild_scheduled_event_id,omitempty"`
	WithCounts            bool                `json:"with_counts,omitempty"`
	WithExpiration        bool                `json:"with_expiration,omitempty"`
}

// Delete Invite
// DELETE /invites/{invite.code}
// https://discord.com/developers/docs/resources/invite#delete-invite
type DeleteInvite struct {
	InviteCode string `json:"code,omitempty"`
}
