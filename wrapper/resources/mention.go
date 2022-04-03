package resources

// Channel Mention Object
// https://discord.com/developers/docs/resources/channel#channel-mention-object
type ChannelMention struct {
	ID      int64  `json:"id"`
	GuildID int64  `json:"guild_id"`
	Type    int    `json:"type"`
	Name    string `json:"name"`
}

// Allowed Mention Types
// https://discord.com/developers/docs/resources/channel#allowed-mentions-object-allowed-mentions-structure
const (
	AllowedMentionRoles     = "roles"
	AllowedMentionsUsers    = "users"
	AllowedMentionsEveryone = "everyone"
)

// Allowed Mentions Structure
// https://discord.com/developers/docs/resources/channel#allowed-mentions-object-allowed-mentions-structure
type AllowedMentions struct {
	Parse       []string `json:"parse"`
	Roles       []string `json:"roles,omitempty"`
	Users       []string `json:"users,omitempty"`
	RepliedUser bool     `json:"replied_user,omitempty"`
}
