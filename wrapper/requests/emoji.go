package requests

import "github.com/switchupcb/disgo/wrapper/resources"

// List Guild Emojis
// GET /guilds/{guild.id}/emojis
// https://discord.com/developers/docs/resources/emoji#list-guild-emojis
type ListGuildEmojis struct {
	GuildID resources.Snowflake
}

// Get Guild Emoji
// GET /guilds/{guild.id}/emojis/{emoji.id}
// https://discord.com/developers/docs/resources/emoji#get-guild-emoji
type GetGuildEmoji struct {
	GuildID resources.Snowflake
	EmojiID resources.Snowflake
}

// Create Guild Emoji
// POST /guilds/{guild.id}/emojis
// https://discord.com/developers/docs/resources/emoji#create-guild-emoji
type CreateGuildEmoji struct {
	GuildID resources.Snowflake
	Name    string                 `json:"name,omitempty"`
	Image   string                 `json:"image,omitempty"`
	Roles   []*resources.Snowflake `json:"roles,omitempty"`
}

// Modify Guild Emoji
// PATCH /guilds/{guild.id}/emojis/{emoji.id}
// https://discord.com/developers/docs/resources/emoji#modify-guild-emoji
type ModifyGuildEmoji struct {
	GuildID resources.Snowflake
	EmojiID resources.Snowflake
	Name    string                 `json:"name,omitempty"`
	Roles   []*resources.Snowflake `json:"roles,omitempty"`
}

// Delete Guild Emoji
// DELETE /guilds/{guild.id}/emojis/{emoji.id}
// https://discord.com/developers/docs/resources/emoji#delete-guild-emoji
type DeleteGuildEmoji struct {
	GuildID resources.Snowflake
	EmojiID resources.Snowflake
}
