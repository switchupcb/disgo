package requests

import "github.com/switchupcb/disgo/wrapper/resources"

// Modify Current User
// PATCH /users/@me
// https://discord.com/developers/docs/resources/user#modify-current-user
type ModifyCurrentUser struct {
	Username string  `json:"username,omitempty"`
	Avatar   *string `json:"avatar,omitempty"`
}

// Get Current User Guilds
// GET /users/@me/guilds
// https://discord.com/developers/docs/resources/user#get-current-user-guilds
type GetCurrentUserGuilds struct {
	Before *resources.Snowflake `json:"before,omitempty"`
	After  *resources.Snowflake `json:"after,omitempty"`
	Limit  resources.Flag       `json:"limit,omitempty"`
}

// Get Current User Guild Member
// GET /users/@me/guilds/{guild.id}/member
// https://discord.com/developers/docs/resources/user#get-current-user-guild-member
type GetCurrentUserGuildMember struct {
	GuildID resources.Snowflake
}

// Leave Guild
// DELETE /users/@me/guilds/{guild.id}
// https://discord.com/developers/docs/resources/user#leave-guild
type LeaveGuild struct {
	GuildID resources.Snowflake
}

// Create DM
// POST /users/@me/channels
// https://discord.com/developers/docs/resources/user#create-dm
type CreateDM struct {
	RecipientID resources.Snowflake `json:"recipient_id,omitempty"`
}

// Create Group DM
// POST /users/@me/channels
// https://discord.com/developers/docs/resources/user#create-group-dm
type CreateGroupDM struct {
	AccessTokens []*string                      `json:"access_tokens,omitempty"`
	Nicks        map[resources.Snowflake]string `json:"nicks,omitempty"`
}

// Get User Connections
// GET /users/@me/connections
// https://discord.com/developers/docs/resources/user#get-user-connections
type GetUserConnections struct {
	RecipientID resources.Snowflake `json:"recipient_id,omitempty"`
}
