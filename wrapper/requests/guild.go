package requests

import (
	"time"

	"github.com/switchupcb/disgo/wrapper/resources"
)

// Create Guild
// POST /guilds
// https://discord.com/developers/docs/resources/guild#create-guild
type CreateGuild struct {
	Name                        string               `json:"name,omitempty"`
	Region                      string               `json:"region,omitempty"`
	Icon                        string               `json:"icon,omitempty"`
	VerificationLevel           *resources.Flag      `json:"verification_level,omitempty"`
	DefaultMessageNotifications *resources.Flag      `json:"default_message_notifications,omitempty"`
	AfkChannelID                string               `json:"afk_channel_id,omitempty"`
	AfkTimeout                  int                  `json:"afk_timeout,omitempty"`
	OwnerID                     string               `json:"owner_id,omitempty"`
	Splash                      string               `json:"splash,omitempty"`
	Banner                      string               `json:"banner,omitempty"`
	Roles                       []*resources.Role    `json:"roles,omitempty"`
	Channels                    []*resources.Channel `json:"channels,omitempty"`
	SystemChannelID             resources.Snowflake  `json:"system_channel_id,omitempty"`
	SystemChannelFlags          resources.BitFlag    `json:"system_channel_flags,omitempty"`
	ExplicitContentFilter       *resources.Flag      `json:"explicit_content_filter,omitempty"`
}

// Get Guild
// GET /guilds/{guild.id}
// https://discord.com/developers/docs/resources/guild#get-guild
type GetGuild struct {
	GuildID    *resources.Snowflake
	WithCounts bool `json:"with_counts,omitempty"`
}

// Get Guild Preview
// GET /guilds/{guild.id}/preview
// https://discord.com/developers/docs/resources/guild#get-guild-preview
type GetGuildPreview struct {
	GuildID *resources.Snowflake
}

// Modify Guild
// PATCH /guilds/{guild.id}
// https://discord.com/developers/docs/resources/guild#modify-guild
type ModifyGuild struct {
	GuildID                     resources.Snowflake
	Name                        string               `json:"name,omitempty"`
	Region                      string               `json:"region,omitempty"`
	VerificationLevel           *resources.Flag      `json:"verification_lvl,omitempty"`
	DefaultMessageNotifications *resources.Flag      `json:"default_message_notifications,omitempty"`
	ExplicitContentFilter       *resources.Flag      `json:"explicit_content_filter,omitempty"`
	AFKChannelID                resources.Snowflake  `json:"afk_channel_id,omitempty"`
	Icon                        *string              `json:"icon,omitempty"`
	OwnerID                     resources.Snowflake  `json:"owner_id,omitempty"`
	Splash                      *string              `json:"splash,omitempty"`
	DiscoverySplash             *string              `json:"discovery_splash,omitempty"`
	Banner                      *string              `json:"banner,omitempty"`
	SystemChannelID             resources.Snowflake  `json:"system_channel_id,omitempty"`
	SystemChannelFlags          resources.BitFlag    `json:"system_channel_flags,omitempty"`
	RulesChannelID              resources.Snowflake  `json:"rules_channel_id,omitempty"`
	PublicUpdatesChannelID      *resources.Snowflake `json:"public_updates_channel_id,omitempty"`
	PreferredLocale             *string              `json:"preferred_locale,omitempty"`
	Features                    []*string            `json:"features,omitempty"`
	Description                 *string              `json:"description,omitempty"`
	PremiumProgressBarEnabled   bool                 `json:"premium_progress_bar_enabled,omitempty"`
}

// Delete Guild
// DELETE /guilds/{guild.id}
// https://discord.com/developers/docs/resources/guild#delete-guild
type DeleteGuild struct {
	GuildID *resources.Snowflake
}

// Get Guild Channels
// GET /guilds/{guild.id}/channels
// https://discord.com/developers/docs/resources/guild#get-guild-channels
type GetGuildChannels struct {
	GuildID *resources.Snowflake
}

// Create Guild Channel
// POST /guilds/{guild.id}/channels
// https://discord.com/developers/docs/resources/guild#create-guild-channel
type CreateGuildChannel struct {
	Name                       string                           `json:"name,omitempty"`
	Type                       *resources.Flag                  `json:"type,omitempty"`
	Topic                      *string                          `json:"topic,omitempty"`
	NSFW                       bool                             `json:"nsfw,omitempty"`
	Position                   int                              `json:"position,omitempty"`
	Bitrate                    int                              `json:"bitrate,omitempty"`
	UserLimit                  int                              `json:"user_limit,omitempty"`
	PermissionOverwrites       []*resources.PermissionOverwrite `json:"permission_overwrites,omitempty"`
	ParentID                   *resources.Snowflake             `json:"parent_id,omitempty"`
	RateLimitPerUser           *resources.CodeFlag              `json:"rate_limit_per_user,omitempty"`
	DefaultAutoArchiveDuration int                              `json:"default_auto_archive_duration,omitempty"`
}

// Modify Guild Channel Positions
// PATCH /guilds/{guild.id}/channels
// https://discord.com/developers/docs/resources/guild#modify-guild-channel-positions
type ModifyGuildChannelPositions struct {
	ID              resources.Snowflake  `json:"id,omitempty"`
	Position        int                  `json:"position,omitempty"`
	LockPermissions bool                 `json:"lock_permissions,omitempty"`
	ParentID        *resources.Snowflake `json:"parent_id,omitempty"`
}

// List Active Guild Threads
// GET /guilds/{guild.id}/threads/active
// https://discord.com/developers/docs/resources/guild#list-active-threads
type ListActiveGuildThreads struct {
	GuildID *resources.Snowflake
	Threads []*resources.Channel      `json:"threads,omitempty"`
	Members []*resources.ThreadMember `json:"members,omitempty"`
}

// Get Guild Member
// GET /guilds/{guild.id}/members/{user.id}
// https://discord.com/developers/docs/resources/guild#get-guild-member
type GetGuildMember struct {
	UserID *resources.Snowflake
}

// List Guild Members
// GET /guilds/{guild.id}/members
// https://discord.com/developers/docs/resources/guild#list-guild-members
type ListGuildMembers struct {
	After *resources.Snowflake `json:"after,omitempty"`
	Limit *resources.CodeFlag  `json:"limit,omitempty"`
}

// Search Guild Members
// GET /guilds/{guild.id}/members/search
// https://discord.com/developers/docs/resources/guild#search-guild-members
type SearchGuildMembers struct {
	GuildID *resources.Snowflake
	Query   string              `json:"query,omitempty"`
	Limit   *resources.CodeFlag `json:"limit,omitempty"`
}

// Add Guild Member
// PUT /guilds/{guild.id}/members/{user.id}
// https://discord.com/developers/docs/resources/guild#add-guild-member
type AddGuildMember struct {
	UserID      *resources.Snowflake
	AccessToken string                 `json:"access_token,omitempty"`
	Nick        string                 `json:"nick,omitempty"`
	Roles       []*resources.Snowflake `json:"roles,omitempty"`
	Mute        bool                   `json:"mute,omitempty"`
	Deaf        bool                   `json:"deaf,omitempty"`
}

// Modify Guild Member
// PATCH /guilds/{guild.id}/members/{user.id}
// https://discord.com/developers/docs/resources/guild#modify-guild-member
type ModifyGuildMember struct {
	UserID                     *resources.Snowflake
	Nick                       string                 `json:"nick,omitempty"`
	Roles                      []*resources.Snowflake `json:"roles,omitempty"`
	Mute                       bool                   `json:"mute,omitempty"`
	Deaf                       bool                   `json:"deaf,omitempty"`
	ChannelID                  resources.Snowflake    `json:"channel_id,omitempty"`
	CommunicationDisabledUntil *time.Time             `json:"communication_disabled_until,omitempty"`
}

// Modify Current Member
// PATCH /guilds/{guild.id}/members/@me
// https://discord.com/developers/docs/resources/guild#modify-current-member
type ModifyCurrentMember struct {
	GuildID resources.Snowflake
	Nick    string `json:"nick,omitempty"`
}

// Modify Current User Nick
// PATCH /guilds/{guild.id}/members/@me/nick
// https://discord.com/developers/docs/resources/guild#modify-current-user-nick
type ModifyCurrentUserNick struct {
	GuildID resources.Snowflake
	Nick    string `json:"nick,omitempty"`
}

// Add Guild Member Role
// PUT /guilds/{guild.id}/members/{user.id}/roles/{role.id}
// https://discord.com/developers/docs/resources/guild#add-guild-member-role
type AddGuildMemberRole struct {
	RoleID resources.Snowflake
}

// Remove Guild Member Role
// DELETE /guilds/{guild.id}/members/{user.id}/roles/{role.id}
// https://discord.com/developers/docs/resources/guild#remove-guild-member-role
type RemoveGuildMemberRole struct {
	RoleID resources.Snowflake
}

// Remove Guild Member
// DELETE /guilds/{guild.id}/members/{user.id}
// https://discord.com/developers/docs/resources/guild#remove-guild-member
type RemoveGuildMember struct {
	UserID *resources.Snowflake
}

// Get Guild Bans
// GET /guilds/{guild.id}/bans
// https://discord.com/developers/docs/resources/guild#get-guild-bans
type GetGuildBans struct {
	GuildID resources.Snowflake
	Before  *resources.Snowflake `json:"before,omitempty"`
	After   *resources.Snowflake `json:"after,omitempty"`
	Limit   *resources.CodeFlag  `json:"limit,omitempty"`
}

// Get Guild Ban
// GET /guilds/{guild.id}/bans/{user.id}
// https://discord.com/developers/docs/resources/guild#get-guild-ban
type GetGuildBan struct {
	UserID *resources.Snowflake
}

// Create Guild Ban
// PUT /guilds/{guild.id}/bans/{user.id}
// https://discord.com/developers/docs/resources/guild#create-guild-ban
type CreateGuildBan struct {
	UserID            *resources.Snowflake
	DeleteMessageDays *resources.Flag `json:"delete_message_days,omitempty"`
	Reason            *string         `json:"reason,omitempty"`
}

// Remove Guild Ban
// DELETE /guilds/{guild.id}/bans/{user.id}
// https://discord.com/developers/docs/resources/guild#remove-guild-ban
type RemoveGuildBan struct {
	UserID *resources.Snowflake
}

// Get Guild Roles
// GET /guilds/{guild.id}/roles
// https://discord.com/developers/docs/resources/guild#get-guild-roles
type GetGuildRoles struct {
	GuildID resources.Snowflake
}

// Create Guild Role
// POST /guilds/{guild.id}/roles
// https://discord.com/developers/docs/resources/guild#create-guild-role
type CreateGuildRole struct {
	GuildID      resources.Snowflake
	Name         string  `json:"name,omitempty"`
	Permissions  string  `json:"permissions,omitempty"`
	Color        *int    `json:"color,omitempty"`
	Hoist        bool    `json:"hoist,omitempty"`
	Icon         *int    `json:"icon,omitempty"`
	UnicodeEmoji *string `json:"unicode_emoji,omitempty"`
	Mentionable  bool    `json:"mentionable,omitempty"`
}

// Modify Guild Role Positions
// PATCH /guilds/{guild.id}/roles
// https://discord.com/developers/docs/resources/guild#modify-guild-role-positions
type ModifyGuildRolePositions struct {
	ID       resources.Snowflake `json:"id,omitempty"`
	Position int                 `json:"position,omitempty"`
}

// Modify Guild Role
// PATCH /guilds/{guild.id}/roles/{role.id}
// https://discord.com/developers/docs/resources/guild#modify-guild-role
type ModifyGuildRole struct {
	RoleID       resources.Snowflake
	Name         string  `json:"name,omitempty"`
	Permissions  int64   `json:"permissions,string,omitempty"`
	Color        *int    `json:"color,omitempty"`
	Hoist        bool    `json:"hoist,omitempty"`
	Icon         *int    `json:"icon,omitempty"`
	UnicodeEmoji *string `json:"unicode_emoji,omitempty"`
	Mentionable  bool    `json:"mentionable,omitempty"`
}

// Delete Guild Role
// DELETE /guilds/{guild.id}/roles/{role.id}
// https://discord.com/developers/docs/resources/guild#delete-guild-role
type DeleteGuildRole struct {
	RoleID resources.Snowflake
}

// Get Guild Prune Count
// GET /guilds/{guild.id}/prune
// https://discord.com/developers/docs/resources/guild#get-guild-prune-count
type GetGuildPruneCount struct {
	GuildID      resources.Snowflake
	Days         resources.Flag         `json:"days,omitempty"`
	IncludeRoles []*resources.Snowflake `json:"include_roles,omitempty"`
}

// Begin Guild Prune
// POST /guilds/{guild.id}/prune
// https://discord.com/developers/docs/resources/guild#begin-guild-prune
type BeginGuildPrune struct {
	GuildID           resources.Snowflake
	Days              resources.Flag         `json:"days,omitempty"`
	ComputePruneCount bool                   `json:"compute_prune_count,omitempty"`
	IncludeRoles      []*resources.Snowflake `json:"include_roles,omitempty"`
	Reason            *string                `json:"reason,omitempty"`
}

// Get Guild Voice Regions
// GET /guilds/{guild.id}/regions
// https://discord.com/developers/docs/resources/guild#get-guild-voice-regions
type GetGuildVoiceRegions struct {
	GuildID resources.Snowflake
}

// Get Guild Invites
// GET /guilds/{guild.id}/invites
// https://discord.com/developers/docs/resources/guild#get-guild-invites
type GetGuildInvites struct {
	GuildID resources.Snowflake
}

// Get Guild Integrations
// GET /guilds/{guild.id}/integrations
// https://discord.com/developers/docs/resources/guild#get-guild-integrations
type GetGuildIntegrations struct {
	GuildID resources.Snowflake
}

// Delete Guild Integration
// DELETE /guilds/{guild.id}/integrations/{integration.id}
// https://discord.com/developers/docs/resources/guild#delete-guild-integration
type DeleteGuildIntegration struct {
	IntegrationID resources.Snowflake
}

// Get Guild Widget Settings
// GET /guilds/{guild.id}/widget
// https://discord.com/developers/docs/resources/guild#get-guild-widget-settings
type GetGuildWidgetSettings struct {
	GuildID resources.Snowflake
}

// Modify Guild Widget
// PATCH /guilds/{guild.id}/widget
// https://discord.com/developers/docs/resources/guild#modify-guild-widget
type ModifyGuildWidget struct {
	GuildID resources.Snowflake
}

// Get Guild Widget
// GET /guilds/{guild.id}/widget.json
// https://discord.com/developers/docs/resources/guild#get-guild-widget
type GetGuildWidget struct {
	GuildID resources.Snowflake
}

// Get Guild Vanity URL
// GET /guilds/{guild.id}/vanity-url
// https://discord.com/developers/docs/resources/guild#get-guild-vanity-url
type GetGuildVanityURL struct {
	GuildID resources.Snowflake
}

// Get Guild Widget Image
// GET /guilds/{guild.id}/widget.png
// https://discord.com/developers/docs/resources/guild#get-guild-widget-image
type GetGuildWidgetImage struct {
	GuildID resources.Snowflake
	Style   string `json:"style,omitempty"`
}

// Get Guild Welcome Screen
// GET /guilds/{guild.id}/welcome-screen
// https://discord.com/developers/docs/resources/guild#get-guild-welcome-screen
type GetGuildWelcomeScreen struct {
	GuildID resources.Snowflake
}

// Modify Guild Welcome Screen
// PATCH /guilds/{guild.id}/welcome-screen
// https://discord.com/developers/docs/resources/guild#modify-guild-welcome-screen
type ModifyGuildWelcomeScreen struct {
	GuildID         resources.Snowflake
	Enabled         bool                              `json:"enabled,omitempty"`
	WelcomeChannels []*resources.WelcomeScreenChannel `json:"welcome_channels,omitempty"`
	Description     *string                           `json:"description,omitempty"`
}

// Modify Current User Voice State
// PATCH /guilds/{guild.id}/voice-states/@me
// https://discord.com/developers/docs/resources/guild#modify-current-user-voice-state
type ModifyCurrentUserVoiceState struct {
	ChannelID               resources.Snowflake
	Suppress                bool       `json:"suppress,omitempty"`
	RequestToSpeakTimestamp *time.Time `json:"request_to_speak_timestamp,omitempty"`
}

// Modify User Voice State
// PATCH /guilds/{guild.id}/voice-states/{user.id}
// https://discord.com/developers/docs/resources/guild#modify-user-voice-state
type ModifyUserVoiceState struct {
	ChannelID resources.Snowflake
	Suppress  bool `json:"suppress,omitempty"`
}
