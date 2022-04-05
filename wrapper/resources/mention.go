package resources

// Channel Mention Object
// https://discord.com/developers/docs/resources/channel#channel-mention-object
type ChannelMention struct {
	ID      Snowflake `json:"id,omitempty"`
	GuildID Snowflake `json:"guild_id,omitempty"`
	Type    Flag      `json:"type,omitempty"`
	Name    string    `json:"name,omitempty"`
}

// Allowed Mention Types
// https://discord.com/developers/docs/resources/channel#allowed-mentions-object-allowed-mentions-structure
const (
	FlagTypesMentionAllowedRoles     = "roles"
	FlagTypesMentionAllowedsUsers    = "users"
	FlagTypesMentionAllowedsEveryone = "everyone"
)

// Allowed Mentions Structure
// https://discord.com/developers/docs/resources/channel#allowed-mentions-object-allowed-mentions-structure
type AllowedMentions struct {
	Parse       []string    `json:"parse,omitempty"`
	Roles       []Snowflake `json:"roles,omitempty"`
	Users       []Snowflake `json:"users,omitempty"`
	RepliedUser bool        `json:"replied_user,omitempty"`
}
