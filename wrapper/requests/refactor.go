package requests

import "github.com/switchupcb/disgo/wrapper/resources"

// Delete Global Application Command
// DELETE /applications/{application.id}/commands/{command.id}
// https://discord.com/developers/docs/interactions/application-commands#delete-global-application-command
type DeleteGlobalApplicationCommand struct {
	CommandID resources.Snowflake
}

// Bulk Overwrite Global Application Commands
// PUT /applications/{application.id}/commands
// https://discord.com/developers/docs/interactions/application-commands#bulk-overwrite-global-application-commands
type BulkOverwriteGlobalApplicationCommands struct {
	ApplicationCommands []*resources.ApplicationCommand
}

// Get Guild Application Commands
// GET /applications/{application.id}/guilds/{guild.id}/commands
// https://discord.com/developers/docs/interactions/application-commands#get-guild-application-commands
type GetGuildApplicationCommands struct {
	GuildID           resources.Snowflake
	WithLocalizations bool `json:"with_localizations,omitempty"`
}

// Create Guild Application Command
// POST /applications/{application.id}/guilds/{guild.id}/commands
// https://discord.com/developers/docs/interactions/application-commands#create-guild-application-command
type CreateGuildApplicationCommand struct {
	GuildID                  resources.Snowflake
	Name                     string                                `json:"name,omitempty"`
	NameLocalizations        map[resources.Flag]string             `json:"name_localizations,omitempty"`
	Description              string                                `json:"description,omitempty"`
	DescriptionLocalizations map[resources.Flag]string             `json:"description_localizations,omitempty"`
	Options                  []*resources.ApplicationCommandOption `json:"options,omitempty"`
	DefaultPermission        bool                                  `json:"default_permission,omitempty"`
	Type                     resources.Flag                        `json:"type,omitempty"`
}

// Get Guild Application Command
// GET /applications/{application.id}/guilds/{guild.id}/commands/{command.id}
// https://discord.com/developers/docs/interactions/application-commands#get-guild-application-command
type GetGuildApplicationCommand struct {
	// TODO
}

// Edit Guild Application Command
// PATCH /applications/{application.id}/guilds/{guild.id}/commands/{command.id}
// https://discord.com/developers/docs/interactions/application-commands#edit-guild-application-command
type EditGuildApplicationCommand struct {
	Name                     string                                `json:"name,omitempty"`
	NameLocalizations        map[resources.Flag]string             `json:"name_localizations,omitempty"`
	Description              string                                `json:"description,omitempty"`
	DescriptionLocalizations map[resources.Flag]string             `json:"description_localizations,omitempty"`
	Options                  []*resources.ApplicationCommandOption `json:"options,omitempty"`
	DefaultPermission        bool                                  `json:"default_permission,omitempty"`
}

// Delete Guild Application Command
// DELETE /applications/{application.id}/guilds/{guild.id}/commands/{command.id}
// TODO
type DeleteGuildApplicationCommand struct {
	// TODO
}

// Bulk Overwrite Guild Application Commands
// PUT /applications/{application.id}/guilds/{guild.id}/commands
// TODO
type BulkOverwriteGuildApplicationCommands struct {
	// TODO
}

// Get Guild Application Command Permissions
// GET /applications/{application.id}/guilds/{guild.id}/commands/permissions
// TODO
type GetGuildApplicationCommandPermissions struct {
	// TODO
}

// Get Application Command Permissions
// GET /applications/{application.id}/guilds/{guild.id}/commands/{command.id}/permissions
// TODO
type GetApplicationCommandPermissions struct {
	// TODO
}

// Edit Application Command Permissions
// PUT /applications/{application.id}/guilds/{guild.id}/commands/{command.id}/permissions
// TODO
type EditApplicationCommandPermissions struct {
	// TODO
}

// Batch Edit Application Command Permissions
// PUT /applications/{application.id}/guilds/{guild.id}/commands/permissions
// TODO
type BatchEditApplicationCommandPermissions struct {
	// TODO
}

/// .go fileinteractions\Receiving_and_Responding.md
// Create Interaction Response
// POST /interactions/{interaction.id}/{interaction.token}/callback
// TODO
type CreateInteractionResponse struct {
	// TODO
}

// Get Original Interaction Response
// GET /webhooks/{application.id}/{interaction.token}/messages/@original
// TODO
type GetOriginalInteractionResponse struct {
	// TODO
}

// Edit Original Interaction Response
// PATCH /webhooks/{application.id}/{interaction.token}/messages/@original
// TODO
type EditOriginalInteractionResponse struct {
	// TODO
}

// Delete Original Interaction Response
// DELETE /webhooks/{application.id}/{interaction.token}/messages/@original
// TODO
type DeleteOriginalInteractionResponse struct {
	// TODO
}

// Create Followup Message
// POST /webhooks/{application.id}/{interaction.token}
// TODO
type CreateFollowupMessage struct {
	// TODO
}

// Get Followup Message
// GET /webhooks/{application.id}/{interaction.token}/messages/{message.id}
// TODO
type GetFollowupMessage struct {
	// TODO
}

// Edit Followup Message
// PATCH /webhooks/{application.id}/{interaction.token}/messages/{message.id}
// TODO
type EditFollowupMessage struct {
	// TODO
}

// Delete Followup Message
// DELETE /webhooks/{application.id}/{interaction.token}/messages/{message.id}
// TODO
type DeleteFollowupMessage struct {
	// TODO
}

/// .go fileresources\Audit_Log.md
// Get Guild Audit Log
// GET /guilds/{guild.id}/audit-logs
// TODO
type GetGuildAuditLog struct {
	// TODO
}

/// .go fileresources\Channel.md
// Get Channel
// GET /channels/{channel.id}
// TODO
type GetChannel struct {
	// TODO
}

// Modify Channel
// PATCH /channels/{channel.id}
// TODO
type ModifyChannel struct {
	// TODO
}

// Delete/Close Channel
// DELETE /channels/{channel.id}
// TODO
type DeleteCloseChannel struct {
	// TODO
}

// Get Channel Messages
// GET /channels/{channel.id}/messages
// TODO
type GetChannelMessages struct {
	// TODO
}

// Get Channel Message
// GET /channels/{channel.id}/messages/{message.id}
// TODO
type GetChannelMessage struct {
	// TODO
}

// Create Message
// POST /channels/{channel.id}/messages
// TODO
type CreateMessage struct {
	// TODO
}

// Crosspost Message
// POST /channels/{channel.id}/messages/{message.id}/crosspost
// TODO
type CrosspostMessage struct {
	// TODO
}

// Create Reaction
// PUT /channels/{channel.id}/messages/{message.id}/reactions/{emoji}/@me
// TODO
type CreateReaction struct {
	// TODO
}

// Delete Own Reaction
// DELETE /channels/{channel.id}/messages/{message.id}/reactions/{emoji}/@me
// TODO
type DeleteOwnReaction struct {
	// TODO
}

// Delete User Reaction
// DELETE /channels/{channel.id}/messages/{message.id}/reactions/{emoji}/{user.id}
// TODO
type DeleteUserReaction struct {
	// TODO
}

// Get Reactions
// GET /channels/{channel.id}/messages/{message.id}/reactions/{emoji}
// TODO
type GetReactions struct {
	// TODO
}

// Delete All Reactions
// DELETE /channels/{channel.id}/messages/{message.id}/reactions
// TODO
type DeleteAllReactions struct {
	// TODO
}

// Delete All Reactions for Emoji
// DELETE /channels/{channel.id}/messages/{message.id}/reactions/{emoji}
// TODO
type DeleteAllReactionsforEmoji struct {
	// TODO
}

// Edit Message
// PATCH /channels/{channel.id}/messages/{message.id}
// TODO
type EditMessage struct {
	// TODO
}

// Delete Message
// DELETE /channels/{channel.id}/messages/{message.id}
// TODO
type DeleteMessage struct {
	// TODO
}

// Bulk Delete Messages
// POST /channels/{channel.id}/messages/bulk-delete
// TODO
type BulkDeleteMessages struct {
	// TODO
}

// Edit Channel Permissions
// PUT /channels/{channel.id}/permissions/{overwrite.id}
// TODO
type EditChannelPermissions struct {
	// TODO
}

// Get Channel Invites
// GET /channels/{channel.id}/invites
// TODO
type GetChannelInvites struct {
	// TODO
}

// Create Channel Invite
// POST /channels/{channel.id}/invites
// TODO
type CreateChannelInvite struct {
	// TODO
}

// Delete Channel Permission
// DELETE /channels/{channel.id}/permissions/{overwrite.id}
// TODO
type DeleteChannelPermission struct {
	// TODO
}

// Follow News Channel
// POST /channels/{channel.id}/followers
// TODO
type FollowNewsChannel struct {
	// TODO
}

// Trigger Typing Indicator
// POST /channels/{channel.id}/typing
// TODO
type TriggerTypingIndicator struct {
	// TODO
}

// Get Pinned Messages
// GET /channels/{channel.id}/pins
// TODO
type GetPinnedMessages struct {
	// TODO
}

// Pin Message
// PUT /channels/{channel.id}/pins/{message.id}
// TODO
type PinMessage struct {
	// TODO
}

// Unpin Message
// DELETE /channels/{channel.id}/pins/{message.id}
// TODO
type UnpinMessage struct {
	// TODO
}

// Group DM Add Recipient
// PUT /channels/{channel.id}/recipients/{user.id}
// TODO
type GroupDMAddRecipient struct {
	// TODO
}

// Group DM Remove Recipient
// DELETE /channels/{channel.id}/recipients/{user.id}
// TODO
type GroupDMRemoveRecipient struct {
	// TODO
}

// Start Thread from Message
// POST /channels/{channel.id}/messages/{message.id}/threads
// TODO
type StartThreadfromMessage struct {
	// TODO
}

// Start Thread without Message
// POST /channels/{channel.id}/threads
// TODO
type StartThreadwithoutMessage struct {
	// TODO
}

// Start Thread in Forum Channel
// POST /channels/{channel.id}/threads
// TODO
type StartThreadinForumChannel struct {
	// TODO
}

// Join Thread
// PUT /channels/{channel.id}/thread-members/@me
// TODO
type JoinThread struct {
	// TODO
}

// Add Thread Member
// PUT /channels/{channel.id}/thread-members/{user.id}
// TODO
type AddThreadMember struct {
	// TODO
}

// Leave Thread
// DELETE /channels/{channel.id}/thread-members/@me
// TODO
type LeaveThread struct {
	// TODO
}

// Remove Thread Member
// DELETE /channels/{channel.id}/thread-members/{user.id}
// TODO
type RemoveThreadMember struct {
	// TODO
}

// Get Thread Member
// GET /channels/{channel.id}/thread-members/{user.id}
// TODO
type GetThreadMember struct {
	// TODO
}

// List Thread Members
// GET /channels/{channel.id}/thread-members
// TODO
type ListThreadMembers struct {
	// TODO
}

// List Active Channel Threads
// GET /channels/{channel.id}/threads/active
// TODO
type ListActiveChannelThreads struct {
	// TODO
}

// List Public Archived Threads
// GET /channels/{channel.id}/threads/archived/public
// TODO
type ListPublicArchivedThreads struct {
	// TODO
}

// List Private Archived Threads
// GET /channels/{channel.id}/threads/archived/private
// TODO
type ListPrivateArchivedThreads struct {
	// TODO
}

// List Joined Private Archived Threads
// GET /channels/{channel.id}/users/@me/threads/archived/private
// TODO
type ListJoinedPrivateArchivedThreads struct {
	// TODO
}

/// .go fileresources\Emoji.md
// List Guild Emojis
// GET /guilds/{guild.id}/emojis
// TODO
type ListGuildEmojis struct {
	// TODO
}

// Get Guild Emoji
// GET /guilds/{guild.id}/emojis/{emoji.id}
// TODO
type GetGuildEmoji struct {
	// TODO
}

// Create Guild Emoji
// POST /guilds/{guild.id}/emojis
// TODO
type CreateGuildEmoji struct {
	// TODO
}

// Modify Guild Emoji
// PATCH /guilds/{guild.id}/emojis/{emoji.id}
// TODO
type ModifyGuildEmoji struct {
	// TODO
}

// Delete Guild Emoji
// DELETE /guilds/{guild.id}/emojis/{emoji.id}
// TODO
type DeleteGuildEmoji struct {
	// TODO
}

/// .go fileresources\Guild_Scheduled_Event.md
// List Scheduled Events for Guild
// GET /guilds/{guild.id}/scheduled-events
// TODO
type ListScheduledEventsforGuild struct {
	// TODO
}

// Create Guild Scheduled Event
// POST /guilds/{guild.id}/scheduled-events
// TODO
type CreateGuildScheduledEvent struct {
	// TODO
}

// Get Guild Scheduled Event
// GET /guilds/{guild.id}/scheduled-events/{guild_scheduled_event.id}
// TODO
type GetGuildScheduledEvent struct {
	// TODO
}

// Modify Guild Scheduled Event
// PATCH /guilds/{guild.id}/scheduled-events/{guild_scheduled_event.id}
// TODO
type ModifyGuildScheduledEvent struct {
	// TODO
}

// Delete Guild Scheduled Event
// DELETE /guilds/{guild.id}/scheduled-events/{guild_scheduled_event.id}
// TODO
type DeleteGuildScheduledEvent struct {
	// TODO
}

// Get Guild Scheduled Event Users
// GET /guilds/{guild.id}/scheduled-events/{guild_scheduled_event.id}/users
// TODO
type GetGuildScheduledEventUsers struct {
	// TODO
}

/// .go fileresources\Guild_Template.md
// Get Guild Template
// GET /guilds/templates/{template.code}
// TODO
type GetGuildTemplate struct {
	// TODO
}

// Create Guild from Guild Template
// POST /guilds/templates/{template.code}
// TODO
type CreateGuildfromGuildTemplate struct {
	// TODO
}

// Get Guild Templates
// GET /guilds/{guild.id}/templates
// TODO
type GetGuildTemplates struct {
	// TODO
}

// Create Guild Template
// POST /guilds/{guild.id}/templates
// TODO
type CreateGuildTemplate struct {
	// TODO
}

// Sync Guild Template
// PUT /guilds/{guild.id}/templates/{template.code}
// TODO
type SyncGuildTemplate struct {
	// TODO
}

// Modify Guild Template
// PATCH /guilds/{guild.id}/templates/{template.code}
// TODO
type ModifyGuildTemplate struct {
	// TODO
}

// Delete Guild Template
// DELETE /guilds/{guild.id}/templates/{template.code}
// TODO
type DeleteGuildTemplate struct {
	// TODO
}

/// .go fileresources\Guild.md
// Create Guild
// POST /guilds
// TODO
type CreateGuild struct {
	// TODO
}

// Get Guild
// GET /guilds/{guild.id}
// TODO
type GetGuild struct {
	// TODO
}

// Get Guild Preview
// GET /guilds/{guild.id}/preview
// TODO
type GetGuildPreview struct {
	// TODO
}

// Modify Guild
// PATCH /guilds/{guild.id}
// TODO
type ModifyGuild struct {
	// TODO
}

// Delete Guild
// DELETE /guilds/{guild.id}
// TODO
type DeleteGuild struct {
	// TODO
}

// Get Guild Channels
// GET /guilds/{guild.id}/channels
// TODO
type GetGuildChannels struct {
	// TODO
}

// Create Guild Channel
// POST /guilds/{guild.id}/channels
// TODO
type CreateGuildChannel struct {
	// TODO
}

// Modify Guild Channel Positions
// PATCH /guilds/{guild.id}/channels
// TODO
type ModifyGuildChannelPositions struct {
	// TODO
}

// List Active Guild Threads
// GET /guilds/{guild.id}/threads/active
// TODO
type ListActiveGuildThreads struct {
	// TODO
}

// Get Guild Member
// GET /guilds/{guild.id}/members/{user.id}
// TODO
type GetGuildMember struct {
	// TODO
}

// List Guild Members
// GET /guilds/{guild.id}/members
// TODO
type ListGuildMembers struct {
	// TODO
}

// Search Guild Members
// GET /guilds/{guild.id}/members/search
// TODO
type SearchGuildMembers struct {
	// TODO
}

// Add Guild Member
// PUT /guilds/{guild.id}/members/{user.id}
// TODO
type AddGuildMember struct {
	// TODO
}

// Modify Guild Member
// PATCH /guilds/{guild.id}/members/{user.id}
// TODO
type ModifyGuildMember struct {
	// TODO
}

// Modify Current Member
// PATCH /guilds/{guild.id}/members/@me
// TODO
type ModifyCurrentMember struct {
	// TODO
}

// Modify Current User Nick
// PATCH /guilds/{guild.id}/members/@me/nick
// TODO
type ModifyCurrentUserNick struct {
	// TODO
}

// Add Guild Member Role
// PUT /guilds/{guild.id}/members/{user.id}/roles/{role.id}
// TODO
type AddGuildMemberRole struct {
	// TODO
}

// Remove Guild Member Role
// DELETE /guilds/{guild.id}/members/{user.id}/roles/{role.id}
// TODO
type RemoveGuildMemberRole struct {
	// TODO
}

// Remove Guild Member
// DELETE /guilds/{guild.id}/members/{user.id}
// TODO
type RemoveGuildMember struct {
	// TODO
}

// Get Guild Bans
// GET /guilds/{guild.id}/bans
// TODO
type GetGuildBans struct {
	// TODO
}

// Get Guild Ban
// GET /guilds/{guild.id}/bans/{user.id}
// TODO
type GetGuildBan struct {
	// TODO
}

// Create Guild Ban
// PUT /guilds/{guild.id}/bans/{user.id}
// TODO
type CreateGuildBan struct {
	// TODO
}

// Remove Guild Ban
// DELETE /guilds/{guild.id}/bans/{user.id}
// TODO
type RemoveGuildBan struct {
	// TODO
}

// Get Guild Roles
// GET /guilds/{guild.id}/roles
// TODO
type GetGuildRoles struct {
	// TODO
}

// Create Guild Role
// POST /guilds/{guild.id}/roles
// TODO
type CreateGuildRole struct {
	// TODO
}

// Modify Guild Role Positions
// PATCH /guilds/{guild.id}/roles
// TODO
type ModifyGuildRolePositions struct {
	// TODO
}

// Modify Guild Role
// PATCH /guilds/{guild.id}/roles/{role.id}
// TODO
type ModifyGuildRole struct {
	// TODO
}

// Delete Guild Role
// DELETE /guilds/{guild.id}/roles/{role.id}
// TODO
type DeleteGuildRole struct {
	// TODO
}

// Get Guild Prune Count
// GET /guilds/{guild.id}/prune
// TODO
type GetGuildPruneCount struct {
	// TODO
}

// Begin Guild Prune
// POST /guilds/{guild.id}/prune
// TODO
type BeginGuildPrune struct {
	// TODO
}

// Get Guild Voice Regions
// GET /guilds/{guild.id}/regions
// TODO
type GetGuildVoiceRegions struct {
	// TODO
}

// Get Guild Invites
// GET /guilds/{guild.id}/invites
// TODO
type GetGuildInvites struct {
	// TODO
}

// Get Guild Integrations
// GET /guilds/{guild.id}/integrations
// TODO
type GetGuildIntegrations struct {
	// TODO
}

// Delete Guild Integration
// DELETE /guilds/{guild.id}/integrations/{integration.id}
// TODO
type DeleteGuildIntegration struct {
	// TODO
}

// Get Guild Widget Settings
// GET /guilds/{guild.id}/widget
// TODO
type GetGuildWidgetSettings struct {
	// TODO
}

// Modify Guild Widget
// PATCH /guilds/{guild.id}/widget
// TODO
type ModifyGuildWidget struct {
	// TODO
}

// Get Guild Widget
// GET /guilds/{guild.id}/widget.json
// TODO
type GetGuildWidget struct {
	// TODO
}

// Get Guild Vanity URL
// GET /guilds/{guild.id}/vanity-url
// TODO
type GetGuildVanityURL struct {
	// TODO
}

// Get Guild Widget Image
// GET /guilds/{guild.id}/widget.png
// TODO
type GetGuildWidgetImage struct {
	// TODO
}

// Get Guild Welcome Screen
// GET /guilds/{guild.id}/welcome-screen
// TODO
type GetGuildWelcomeScreen struct {
	// TODO
}

// Modify Guild Welcome Screen
// PATCH /guilds/{guild.id}/welcome-screen
// TODO
type ModifyGuildWelcomeScreen struct {
	// TODO
}

// Modify Current User Voice State
// PATCH /guilds/{guild.id}/voice-states/@me
// TODO
type ModifyCurrentUserVoiceState struct {
	// TODO
}

// Modify User Voice State
// PATCH /guilds/{guild.id}/voice-states/{user.id}
// TODO
type ModifyUserVoiceState struct {
	// TODO
}

/// .go fileresources\Invite.md
// Get Invite
// GET /invites/{invite.code}
// TODO
type GetInvite struct {
	// TODO
}

// Delete Invite
// DELETE /invites/{invite.code}
// TODO
type DeleteInvite struct {
	// TODO
}

/// .go fileresources\Stage_Instance.md
// Create Stage Instance
// POST /stage-instances
// TODO
type CreateStageInstance struct {
	// TODO
}

// Get Stage Instance
// GET /stage-instances/{channel.id}
// TODO
type GetStageInstance struct {
	// TODO
}

// Modify Stage Instance
// PATCH /stage-instances/{channel.id}
// TODO
type ModifyStageInstance struct {
	// TODO
}

// Delete Stage Instance
// DELETE /stage-instances/{channel.id}
// TODO
type DeleteStageInstance struct {
	// TODO
}

/// .go fileresources\Sticker.md
// Get Sticker
// GET /stickers/{sticker.id}
// TODO
type GetSticker struct {
	// TODO
}

// List Nitro Sticker Packs
// GET /sticker-packs
// TODO
type ListNitroStickerPacks struct {
	// TODO
}

// List Guild Stickers
// GET /guilds/{guild.id}/stickers
// TODO
type ListGuildStickers struct {
	// TODO
}

// Get Guild Sticker
// GET /guilds/{guild.id}/stickers/{sticker.id}
// TODO
type GetGuildSticker struct {
	// TODO
}

// Create Guild Sticker
// POST /guilds/{guild.id}/stickers
// TODO
type CreateGuildSticker struct {
	// TODO
}

// Modify Guild Sticker
// PATCH /guilds/{guild.id}/stickers/{sticker.id}
// TODO
type ModifyGuildSticker struct {
	// TODO
}

// Delete Guild Sticker
// DELETE /guilds/{guild.id}/stickers/{sticker.id}
// TODO
type DeleteGuildSticker struct {
	// TODO
}

/// .go fileresources\User.md
// Get Current User
// GET /users/@me
// TODO
type GetCurrentUser struct {
	// TODO
}

// Get User
// GET /users/{user.id}
// TODO
type GetUser struct {
	// TODO
}

// Modify Current User
// PATCH /users/@me
// TODO
type ModifyCurrentUser struct {
	// TODO
}

// Get Current User Guilds
// GET /users/@me/guilds
// TODO
type GetCurrentUserGuilds struct {
	// TODO
}

// Get Current User Guild Member
// GET /users/@me/guilds/{guild.id}/member
// TODO
type GetCurrentUserGuildMember struct {
	// TODO
}

// Leave Guild
// DELETE /users/@me/guilds/{guild.id}
// TODO
type LeaveGuild struct {
	// TODO
}

// Create DM
// POST /users/@me/channels
// TODO
type CreateDM struct {
	// TODO
}

// Create Group DM
// POST /users/@me/channels
// TODO
type CreateGroupDM struct {
	// TODO
}

// Get User Connections
// GET /users/@me/connections
// TODO
type GetUserConnections struct {
	// TODO
}

/// .go fileresources\Voice.md
// List Voice Regions
// GET /voice/regions
// TODO
type ListVoiceRegions struct {
	// TODO
}

/// .go fileresources\Webhook.md
// Create Webhook
// POST /channels/{channel.id}/webhooks
// TODO
type CreateWebhook struct {
	// TODO
}

// Get Channel Webhooks
// GET /channels/{channel.id}/webhooks
// TODO
type GetChannelWebhooks struct {
	// TODO
}

// Get Guild Webhooks
// GET /guilds/{guild.id}/webhooks
// TODO
type GetGuildWebhooks struct {
	// TODO
}

// Get Webhook
// GET /webhooks/{webhook.id}
// TODO
type GetWebhook struct {
	// TODO
}

// Get Webhook with Token
// GET /webhooks/{webhook.id}/{webhook.token}
// TODO
type GetWebhookwithToken struct {
	// TODO
}

// Modify Webhook
// PATCH /webhooks/{webhook.id}
// TODO
type ModifyWebhook struct {
	// TODO
}

// Modify Webhook with Token
// PATCH /webhooks/{webhook.id}/{webhook.token}
// TODO
type ModifyWebhookwithToken struct {
	// TODO
}

// Delete Webhook
// DELETE /webhooks/{webhook.id}
// TODO
type DeleteWebhook struct {
	// TODO
}

// Delete Webhook with Token
// DELETE /webhooks/{webhook.id}/{webhook.token}
// TODO
type DeleteWebhookwithToken struct {
	// TODO
}

// Execute Webhook
// POST /webhooks/{webhook.id}/{webhook.token}
// TODO
type ExecuteWebhook struct {
	// TODO
}

// Execute Slack-Compatible Webhook
// POST /webhooks/{webhook.id}/{webhook.token}/slack
// TODO
type ExecuteSlackCompatibleWebhook struct {
	// TODO
}

// Execute GitHub-Compatible Webhook
// POST /webhooks/{webhook.id}/{webhook.token}/github
// TODO
type ExecuteGitHubCompatibleWebhook struct {
	// TODO
}

// Get Webhook Message
// GET /webhooks/{webhook.id}/{webhook.token}/messages/{message.id}
// TODO
type GetWebhookMessage struct {
	// TODO
}

// Edit Webhook Message
// PATCH /webhooks/{webhook.id}/{webhook.token}/messages/{message.id}
// TODO
type EditWebhookMessage struct {
	// TODO
}

// Delete Webhook Message
// DELETE /webhooks/{webhook.id}/{webhook.token}/messages/{message.id}
// TODO
type DeleteWebhookMessage struct {
	// TODO
}

/// .go filetopics\Gateway.md
// Get Gateway
// GET /gateway
// TODO
type GetGateway struct {
	// TODO
}

// Get Gateway Bot
// GET /gateway/bot
// TODO
type GetGatewayBot struct {
	// TODO
}

/// .go filetopics\OAuth2.md
// Get Current Bot Application Information
// GET /oauth2/applications/@me
// TODO
type GetCurrentBotApplicationInformation struct {
	// TODO
}

// Get Current Authorization Information
// GET /oauth2/@me
// TODO
type GetCurrentAuthorizationInformation struct {
	// TODO
}
