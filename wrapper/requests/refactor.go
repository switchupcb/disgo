package requests

/// .go filegame_sdk\Achievements.md
// Get Achievements
// GET /applications/{application.id#DOCS_GAME_SDK_SDK_STARTER_GUIDE/get-set-up}/achievements
// TODO
type GetAchievements struct {
	// TODO
}

// Get Achievement
// GET /applications/{application.id#DOCS_GAME_SDK_SDK_STARTER_GUIDE/get-set-up}/achievements/{achievement.id#DOCS_GAME_SDK_ACHIEVEMENTS/data-models-achievement-struct}
// TODO
type GetAchievement struct {
	// TODO
}

// Create Achievement
// POST /applications/{application.id#DOCS_GAME_SDK_SDK_STARTER_GUIDE/get-set-up}/achievements
// TODO
type CreateAchievement struct {
	// TODO
}

// Update Achievement
// PATCH /applications/{application.id#DOCS_GAME_SDK_SDK_STARTER_GUIDE/get-set-up}/achievements/{achievement.id#DOCS_GAME_SDK_ACHIEVEMENTS/data-models-achievement-struct}
// TODO
type UpdateAchievement struct {
	// TODO
}

// Delete Achievement
// DELETE /applications/{application.id#DOCS_GAME_SDK_SDK_STARTER_GUIDE/get-set-up}/achievements/{achievement.id#DOCS_GAME_SDK_ACHIEVEMENTS/data-models-achievement-struct}
// TODO
type DeleteAchievement struct {
	// TODO
}

// Update User Achievement
// PUT /applications/{application.id#DOCS_GAME_SDK_SDK_STARTER_GUIDE/get-set-up}/achievements/{achievement.id#DOCS_GAME_SDK_ACHIEVEMENTS/data-models-achievement-struct}
// TODO
type UpdateUserAchievement struct {
	// TODO
}

// Get User Achievements
// GET /users/@me/applications/{application.id#DOCS_GAME_SDK_SDK_STARTER_GUIDE/get-set-up}/achievements
// TODO
type GetUserAchievements struct {
	// TODO
}

/// .go filegame_sdk\Lobbies.md
// Create Lobby
// POST /lobbies
// TODO
type CreateLobby struct {
	// TODO
}

// Update Lobby
// PATCH /lobbies/{lobby.id#DOCS_LOBBIES/data-models-lobby-struct}
// TODO
type UpdateLobby struct {
	// TODO
}

// Delete Lobby
// DELETE /lobbies/{lobby.id#DOCS_GAME_SDK_LOBBIES/data-models-lobby-struct}
// TODO
type DeleteLobby struct {
	// TODO
}

// Update Lobby Member
// PATCH /lobbies/{lobby.id#DOCS_GAME_SDK_LOBBIES/data-models-lobby-struct}/members/{user.id#DOCS_RESOURCES_USER/user-object}
// TODO
type UpdateLobbyMember struct {
	// TODO
}

// Create Lobby Search
// POST /lobbies/search
// TODO
type CreateLobbySearch struct {
	// TODO
}

// Send Lobby Data
// POST /lobbies/{lobby.id#DOCS_GAME_SDK_LOBBIES/data-models-lobby-struct}/send
// TODO
type SendLobbyData struct {
	// TODO
}

/// .go filegame_sdk\Store.md
// Get Entitlements
// GET /applications/{application.id#DOCS_GAME_SDK_SDK_STARTER_GUIDE/get-set-up}/entitlements
// TODO
type GetEntitlements struct {
	// TODO
}

// Get Entitlement
// GET /applications/{application.id#DOCS_GAME_SDK_SDK_STARTER_GUIDE/get-set-up}/entitlements/{entitlement.id#DOCS_GAME_SDK_STORE/data-models-entitlement-struct}
// TODO
type GetEntitlement struct {
	// TODO
}

// Get SKUs
// GET /applications/{application.id#DOCS_GAME_SDK_SDK_STARTER_GUIDE/get-set-up}/skus
// TODO
type GetSKUs struct {
	// TODO
}

// Consume SKU
// POST /applications/{application.id#DOCS_GAME_SDK_SDK_STARTER_GUIDE/get-set-up}/entitlements/{entitlement.id#DOCS_GAME_SDK_STORE/data-models-entitlement-struct}/consume
// TODO
type ConsumeSKU struct {
	// TODO
}

// Delete Test Entitlement
// DELETE /applications/{application.id#DOCS_GAME_SDK_SDK_STARTER_GUIDE/get-set-up}/entitlements/{entitlement.id#DOCS_GAME_SDK_STORE/data-models-entitlement-struct}
// TODO
type DeleteTestEntitlement struct {
	// TODO
}

// Create Purchase Discount
// PUT /store/skus/{sku.id#DOCS_GAME_SDK_STORE/data-models-sku-struct}/discounts/{user.id#DOCS_RESOURCES_USER/user-object}
// TODO
type CreatePurchaseDiscount struct {
	// TODO
}

// Delete Purchase Discount
// DELETE /store/skus/{sku.id#DOCS_GAME_SDK_STORE/data-models-sku-struct}/discounts/{user.id#DOCS_RESOURCES_USER/user-object}
// TODO
type DeletePurchaseDiscount struct {
	// TODO
}

// Delete Global Application Command
// DELETE /applications/{application.id#DOCS_RESOURCES_APPLICATION/application-object}/commands/{command.id#DOCS_INTERACTIONS_APPLICATION_COMMANDS/application-command-object}
// TODO
type DeleteGlobalApplicationCommand struct {
	// TODO
}

// Bulk Overwrite Global Application Commands
// PUT /applications/{application.id#DOCS_RESOURCES_APPLICATION/application-object}/commands
// TODO
type BulkOverwriteGlobalApplicationCommands struct {
	// TODO
}

// Get Guild Application Commands
// GET /applications/{application.id#DOCS_RESOURCES_APPLICATION/application-object}/guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/commands
// TODO
type GetGuildApplicationCommands struct {
	// TODO
}

// Create Guild Application Command
// POST /applications/{application.id#DOCS_RESOURCES_APPLICATION/application-object}/guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/commands
// TODO
type CreateGuildApplicationCommand struct {
	// TODO
}

// Get Guild Application Command
// GET /applications/{application.id#DOCS_RESOURCES_APPLICATION/application-object}/guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/commands/{command.id#DOCS_INTERACTIONS_APPLICATION_COMMANDS/application-command-object}
// TODO
type GetGuildApplicationCommand struct {
	// TODO
}

// Edit Guild Application Command
// PATCH /applications/{application.id#DOCS_RESOURCES_APPLICATION/application-object}/guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/commands/{command.id#DOCS_INTERACTIONS_APPLICATION_COMMANDS/application-command-object}
// TODO
type EditGuildApplicationCommand struct {
	// TODO
}

// Delete Guild Application Command
// DELETE /applications/{application.id#DOCS_RESOURCES_APPLICATION/application-object}/guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/commands/{command.id#DOCS_INTERACTIONS_APPLICATION_COMMANDS/application-command-object}
// TODO
type DeleteGuildApplicationCommand struct {
	// TODO
}

// Bulk Overwrite Guild Application Commands
// PUT /applications/{application.id#DOCS_RESOURCES_APPLICATION/application-object}/guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/commands
// TODO
type BulkOverwriteGuildApplicationCommands struct {
	// TODO
}

// Get Guild Application Command Permissions
// GET /applications/{application.id#DOCS_RESOURCES_APPLICATION/application-object}/guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/commands/permissions
// TODO
type GetGuildApplicationCommandPermissions struct {
	// TODO
}

// Get Application Command Permissions
// GET /applications/{application.id#DOCS_RESOURCES_APPLICATION/application-object}/guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/commands/{command.id#DOCS_INTERACTIONS_APPLICATION_COMMANDS/application-command-object}/permissions
// TODO
type GetApplicationCommandPermissions struct {
	// TODO
}

// Edit Application Command Permissions
// PUT /applications/{application.id#DOCS_RESOURCES_APPLICATION/application-object}/guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/commands/{command.id#DOCS_INTERACTIONS_APPLICATION_COMMANDS/application-command-object}/permissions
// TODO
type EditApplicationCommandPermissions struct {
	// TODO
}

// Batch Edit Application Command Permissions
// PUT /applications/{application.id#DOCS_RESOURCES_APPLICATION/application-object}/guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/commands/permissions
// TODO
type BatchEditApplicationCommandPermissions struct {
	// TODO
}

/// .go fileinteractions\Receiving_and_Responding.md
// Create Interaction Response
// POST /interactions/{interaction.id#DOCS_INTERACTIONS_RECEIVING_AND_RESPONDING/interaction}/{interaction.token#DOCS_INTERACTIONS_RECEIVING_AND_RESPONDING/interaction-object}/callback
// TODO
type CreateInteractionResponse struct {
	// TODO
}

// Get Original Interaction Response
// GET /webhooks/{application.id#DOCS_RESOURCES_APPLICATION/application-object}/{interaction.token#DOCS_INTERACTIONS_RECEIVING_AND_RESPONDING/interaction-object}/messages/@original
// TODO
type GetOriginalInteractionResponse struct {
	// TODO
}

// Edit Original Interaction Response
// PATCH /webhooks/{application.id#DOCS_RESOURCES_APPLICATION/application-object}/{interaction.token#DOCS_INTERACTIONS_RECEIVING_AND_RESPONDING/interaction-object}/messages/@original
// TODO
type EditOriginalInteractionResponse struct {
	// TODO
}

// Delete Original Interaction Response
// DELETE /webhooks/{application.id#DOCS_RESOURCES_APPLICATION/application-object}/{interaction.token#DOCS_INTERACTIONS_RECEIVING_AND_RESPONDING/interaction-object}/messages/@original
// TODO
type DeleteOriginalInteractionResponse struct {
	// TODO
}

// Create Followup Message
// POST /webhooks/{application.id#DOCS_RESOURCES_APPLICATION/application-object}/{interaction.token#DOCS_INTERACTIONS_RECEIVING_AND_RESPONDING/interaction-object}
// TODO
type CreateFollowupMessage struct {
	// TODO
}

// Get Followup Message
// GET /webhooks/{application.id#DOCS_RESOURCES_APPLICATION/application-object}/{interaction.token#DOCS_INTERACTIONS_RECEIVING_AND_RESPONDING/interaction-object}/messages/{message.id#DOCS_RESOURCES_CHANNEL/message-object}
// TODO
type GetFollowupMessage struct {
	// TODO
}

// Edit Followup Message
// PATCH /webhooks/{application.id#DOCS_RESOURCES_APPLICATION/application-object}/{interaction.token#DOCS_INTERACTIONS_RECEIVING_AND_RESPONDING/interaction-object}/messages/{message.id#DOCS_RESOURCES_CHANNEL/message-object}
// TODO
type EditFollowupMessage struct {
	// TODO
}

// Delete Followup Message
// DELETE /webhooks/{application.id#DOCS_RESOURCES_APPLICATION/application-object}/{interaction.token#DOCS_INTERACTIONS_RECEIVING_AND_RESPONDING/interaction-object}/messages/{message.id#DOCS_RESOURCES_CHANNEL/message-object}
// TODO
type DeleteFollowupMessage struct {
	// TODO
}

/// .go fileresources\Audit_Log.md
// Get Guild Audit Log
// GET /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/audit-logs
// TODO
type GetGuildAuditLog struct {
	// TODO
}

/// .go fileresources\Channel.md
// Get Channel
// GET /channels/{channel.id#DOCS_RESOURCES_CHANNEL/channel-object}
// TODO
type GetChannel struct {
	// TODO
}

// Modify Channel
// PATCH /channels/{channel.id#DOCS_RESOURCES_CHANNEL/channel-object}
// TODO
type ModifyChannel struct {
	// TODO
}

// Delete/Close Channel
// DELETE /channels/{channel.id#DOCS_RESOURCES_CHANNEL/channel-object}
// TODO
type DeleteCloseChannel struct {
	// TODO
}

// Get Channel Messages
// GET /channels/{channel.id#DOCS_RESOURCES_CHANNEL/channel-object}/messages
// TODO
type GetChannelMessages struct {
	// TODO
}

// Get Channel Message
// GET /channels/{channel.id#DOCS_RESOURCES_CHANNEL/channel-object}/messages/{message.id#DOCS_RESOURCES_CHANNEL/message-object}
// TODO
type GetChannelMessage struct {
	// TODO
}

// Create Message
// POST /channels/{channel.id#DOCS_RESOURCES_CHANNEL/channel-object}/messages
// TODO
type CreateMessage struct {
	// TODO
}

// Crosspost Message
// POST /channels/{channel.id#DOCS_RESOURCES_CHANNEL/channel-object}/messages/{message.id#DOCS_RESOURCES_CHANNEL/message-object}/crosspost
// TODO
type CrosspostMessage struct {
	// TODO
}

// Create Reaction
// PUT /channels/{channel.id#DOCS_RESOURCES_CHANNEL/channel-object}/messages/{message.id#DOCS_RESOURCES_CHANNEL/message-object}/reactions/{emoji#DOCS_RESOURCES_EMOJI/emoji-object}/@me
// TODO
type CreateReaction struct {
	// TODO
}

// Delete Own Reaction
// DELETE /channels/{channel.id#DOCS_RESOURCES_CHANNEL/channel-object}/messages/{message.id#DOCS_RESOURCES_CHANNEL/message-object}/reactions/{emoji#DOCS_RESOURCES_EMOJI/emoji-object}/@me
// TODO
type DeleteOwnReaction struct {
	// TODO
}

// Delete User Reaction
// DELETE /channels/{channel.id#DOCS_RESOURCES_CHANNEL/channel-object}/messages/{message.id#DOCS_RESOURCES_CHANNEL/message-object}/reactions/{emoji#DOCS_RESOURCES_EMOJI/emoji-object}/{user.id#DOCS_RESOURCES_USER/user-object}
// TODO
type DeleteUserReaction struct {
	// TODO
}

// Get Reactions
// GET /channels/{channel.id#DOCS_RESOURCES_CHANNEL/channel-object}/messages/{message.id#DOCS_RESOURCES_CHANNEL/message-object}/reactions/{emoji#DOCS_RESOURCES_EMOJI/emoji-object}
// TODO
type GetReactions struct {
	// TODO
}

// Delete All Reactions
// DELETE /channels/{channel.id#DOCS_RESOURCES_CHANNEL/channel-object}/messages/{message.id#DOCS_RESOURCES_CHANNEL/message-object}/reactions
// TODO
type DeleteAllReactions struct {
	// TODO
}

// Delete All Reactions for Emoji
// DELETE /channels/{channel.id#DOCS_RESOURCES_CHANNEL/channel-object}/messages/{message.id#DOCS_RESOURCES_CHANNEL/message-object}/reactions/{emoji#DOCS_RESOURCES_EMOJI/emoji-object}
// TODO
type DeleteAllReactionsforEmoji struct {
	// TODO
}

// Edit Message
// PATCH /channels/{channel.id#DOCS_RESOURCES_CHANNEL/channel-object}/messages/{message.id#DOCS_RESOURCES_CHANNEL/message-object}
// TODO
type EditMessage struct {
	// TODO
}

// Delete Message
// DELETE /channels/{channel.id#DOCS_RESOURCES_CHANNEL/channel-object}/messages/{message.id#DOCS_RESOURCES_CHANNEL/message-object}
// TODO
type DeleteMessage struct {
	// TODO
}

// Bulk Delete Messages
// POST /channels/{channel.id#DOCS_RESOURCES_CHANNEL/channel-object}/messages/bulk-delete
// TODO
type BulkDeleteMessages struct {
	// TODO
}

// Edit Channel Permissions
// PUT /channels/{channel.id#DOCS_RESOURCES_CHANNEL/channel-object}/permissions/{overwrite.id#DOCS_RESOURCES_CHANNEL/overwrite-object}
// TODO
type EditChannelPermissions struct {
	// TODO
}

// Get Channel Invites
// GET /channels/{channel.id#DOCS_RESOURCES_CHANNEL/channel-object}/invites
// TODO
type GetChannelInvites struct {
	// TODO
}

// Create Channel Invite
// POST /channels/{channel.id#DOCS_RESOURCES_CHANNEL/channel-object}/invites
// TODO
type CreateChannelInvite struct {
	// TODO
}

// Delete Channel Permission
// DELETE /channels/{channel.id#DOCS_RESOURCES_CHANNEL/channel-object}/permissions/{overwrite.id#DOCS_RESOURCES_CHANNEL/overwrite-object}
// TODO
type DeleteChannelPermission struct {
	// TODO
}

// Follow News Channel
// POST /channels/{channel.id#DOCS_RESOURCES_CHANNEL/channel-object}/followers
// TODO
type FollowNewsChannel struct {
	// TODO
}

// Trigger Typing Indicator
// POST /channels/{channel.id#DOCS_RESOURCES_CHANNEL/channel-object}/typing
// TODO
type TriggerTypingIndicator struct {
	// TODO
}

// Get Pinned Messages
// GET /channels/{channel.id#DOCS_RESOURCES_CHANNEL/channel-object}/pins
// TODO
type GetPinnedMessages struct {
	// TODO
}

// Pin Message
// PUT /channels/{channel.id#DOCS_RESOURCES_CHANNEL/channel-object}/pins/{message.id#DOCS_RESOURCES_CHANNEL/message-object}
// TODO
type PinMessage struct {
	// TODO
}

// Unpin Message
// DELETE /channels/{channel.id#DOCS_RESOURCES_CHANNEL/channel-object}/pins/{message.id#DOCS_RESOURCES_CHANNEL/message-object}
// TODO
type UnpinMessage struct {
	// TODO
}

// Group DM Add Recipient
// PUT /channels/{channel.id#DOCS_RESOURCES_CHANNEL/channel-object}/recipients/{user.id#DOCS_RESOURCES_USER/user-object}
// TODO
type GroupDMAddRecipient struct {
	// TODO
}

// Group DM Remove Recipient
// DELETE /channels/{channel.id#DOCS_RESOURCES_CHANNEL/channel-object}/recipients/{user.id#DOCS_RESOURCES_USER/user-object}
// TODO
type GroupDMRemoveRecipient struct {
	// TODO
}

// Start Thread from Message
// POST /channels/{channel.id#DOCS_RESOURCES_CHANNEL/channel-object}/messages/{message.id#DOCS_RESOURCES_CHANNEL/message-object}/threads
// TODO
type StartThreadfromMessage struct {
	// TODO
}

// Start Thread without Message
// POST /channels/{channel.id#DOCS_RESOURCES_CHANNEL/channel-object}/threads
// TODO
type StartThreadwithoutMessage struct {
	// TODO
}

// Start Thread in Forum Channel
// POST /channels/{channel.id#DOCS_RESOURCES_CHANNEL/channel-object}/threads
// TODO
type StartThreadinForumChannel struct {
	// TODO
}

// Join Thread
// PUT /channels/{channel.id#DOCS_RESOURCES_CHANNEL/channel-object}/thread-members/@me
// TODO
type JoinThread struct {
	// TODO
}

// Add Thread Member
// PUT /channels/{channel.id#DOCS_RESOURCES_CHANNEL/channel-object}/thread-members/{user.id#DOCS_RESOURCES_USER/user-object}
// TODO
type AddThreadMember struct {
	// TODO
}

// Leave Thread
// DELETE /channels/{channel.id#DOCS_RESOURCES_CHANNEL/channel-object}/thread-members/@me
// TODO
type LeaveThread struct {
	// TODO
}

// Remove Thread Member
// DELETE /channels/{channel.id#DOCS_RESOURCES_CHANNEL/channel-object}/thread-members/{user.id#DOCS_RESOURCES_USER/user-object}
// TODO
type RemoveThreadMember struct {
	// TODO
}

// Get Thread Member
// GET /channels/{channel.id#DOCS_RESOURCES_CHANNEL/channel-object}/thread-members/{user.id#DOCS_RESOURCES_USER/user-object}
// TODO
type GetThreadMember struct {
	// TODO
}

// List Thread Members
// GET /channels/{channel.id#DOCS_RESOURCES_CHANNEL/channel-object}/thread-members
// TODO
type ListThreadMembers struct {
	// TODO
}

// List Active Channel Threads
// GET /channels/{channel.id#DOCS_RESOURCES_CHANNEL/channel-object}/threads/active
// TODO
type ListActiveChannelThreads struct {
	// TODO
}

// List Public Archived Threads
// GET /channels/{channel.id#DOCS_RESOURCES_CHANNEL/channel-object}/threads/archived/public
// TODO
type ListPublicArchivedThreads struct {
	// TODO
}

// List Private Archived Threads
// GET /channels/{channel.id#DOCS_RESOURCES_CHANNEL/channel-object}/threads/archived/private
// TODO
type ListPrivateArchivedThreads struct {
	// TODO
}

// List Joined Private Archived Threads
// GET /channels/{channel.id#DOCS_RESOURCES_CHANNEL/channel-object}/users/@me/threads/archived/private
// TODO
type ListJoinedPrivateArchivedThreads struct {
	// TODO
}

/// .go fileresources\Emoji.md
// List Guild Emojis
// GET /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/emojis
// TODO
type ListGuildEmojis struct {
	// TODO
}

// Get Guild Emoji
// GET /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/emojis/{emoji.id#DOCS_RESOURCES_EMOJI/emoji-object}
// TODO
type GetGuildEmoji struct {
	// TODO
}

// Create Guild Emoji
// POST /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/emojis
// TODO
type CreateGuildEmoji struct {
	// TODO
}

// Modify Guild Emoji
// PATCH /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/emojis/{emoji.id#DOCS_RESOURCES_EMOJI/emoji-object}
// TODO
type ModifyGuildEmoji struct {
	// TODO
}

// Delete Guild Emoji
// DELETE /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/emojis/{emoji.id#DOCS_RESOURCES_EMOJI/emoji-object}
// TODO
type DeleteGuildEmoji struct {
	// TODO
}

/// .go fileresources\Guild_Scheduled_Event.md
// List Scheduled Events for Guild
// GET /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/scheduled-events
// TODO
type ListScheduledEventsforGuild struct {
	// TODO
}

// Create Guild Scheduled Event
// POST /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/scheduled-events
// TODO
type CreateGuildScheduledEvent struct {
	// TODO
}

// Get Guild Scheduled Event
// GET /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/scheduled-events/{guild_scheduled_event.id#DOCS_RESOURCES_GUILD_SCHEDULED_EVENT/guild-scheduled-event-object}
// TODO
type GetGuildScheduledEvent struct {
	// TODO
}

// Modify Guild Scheduled Event
// PATCH /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/scheduled-events/{guild_scheduled_event.id#DOCS_RESOURCES_GUILD_SCHEDULED_EVENT/guild-scheduled-event-object}
// TODO
type ModifyGuildScheduledEvent struct {
	// TODO
}

// Delete Guild Scheduled Event
// DELETE /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/scheduled-events/{guild_scheduled_event.id#DOCS_RESOURCES_GUILD_SCHEDULED_EVENT/guild-scheduled-event-object}
// TODO
type DeleteGuildScheduledEvent struct {
	// TODO
}

// Get Guild Scheduled Event Users
// GET /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/scheduled-events/{guild_scheduled_event.id#DOCS_RESOURCES_GUILD_SCHEDULED_EVENT/guild-scheduled-event-object}/users
// TODO
type GetGuildScheduledEventUsers struct {
	// TODO
}

/// .go fileresources\Guild_Template.md
// Get Guild Template
// GET /guilds/templates/{template.code#DOCS_RESOURCES_GUILD_TEMPLATE/guild-template-object}
// TODO
type GetGuildTemplate struct {
	// TODO
}

// Create Guild from Guild Template
// POST /guilds/templates/{template.code#DOCS_RESOURCES_GUILD_TEMPLATE/guild-template-object}
// TODO
type CreateGuildfromGuildTemplate struct {
	// TODO
}

// Get Guild Templates
// GET /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/templates
// TODO
type GetGuildTemplates struct {
	// TODO
}

// Create Guild Template
// POST /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/templates
// TODO
type CreateGuildTemplate struct {
	// TODO
}

// Sync Guild Template
// PUT /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/templates/{template.code#DOCS_RESOURCES_GUILD_TEMPLATE/guild-template-object}
// TODO
type SyncGuildTemplate struct {
	// TODO
}

// Modify Guild Template
// PATCH /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/templates/{template.code#DOCS_RESOURCES_GUILD_TEMPLATE/guild-template-object}
// TODO
type ModifyGuildTemplate struct {
	// TODO
}

// Delete Guild Template
// DELETE /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/templates/{template.code#DOCS_RESOURCES_GUILD_TEMPLATE/guild-template-object}
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
// GET /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}
// TODO
type GetGuild struct {
	// TODO
}

// Get Guild Preview
// GET /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/preview
// TODO
type GetGuildPreview struct {
	// TODO
}

// Modify Guild
// PATCH /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}
// TODO
type ModifyGuild struct {
	// TODO
}

// Delete Guild
// DELETE /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}
// TODO
type DeleteGuild struct {
	// TODO
}

// Get Guild Channels
// GET /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/channels
// TODO
type GetGuildChannels struct {
	// TODO
}

// Create Guild Channel
// POST /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/channels
// TODO
type CreateGuildChannel struct {
	// TODO
}

// Modify Guild Channel Positions
// PATCH /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/channels
// TODO
type ModifyGuildChannelPositions struct {
	// TODO
}

// List Active Guild Threads
// GET /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/threads/active
// TODO
type ListActiveGuildThreads struct {
	// TODO
}

// Get Guild Member
// GET /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/members/{user.id#DOCS_RESOURCES_USER/user-object}
// TODO
type GetGuildMember struct {
	// TODO
}

// List Guild Members
// GET /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/members
// TODO
type ListGuildMembers struct {
	// TODO
}

// Search Guild Members
// GET /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/members/search
// TODO
type SearchGuildMembers struct {
	// TODO
}

// Add Guild Member
// PUT /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/members/{user.id#DOCS_RESOURCES_USER/user-object}
// TODO
type AddGuildMember struct {
	// TODO
}

// Modify Guild Member
// PATCH /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/members/{user.id#DOCS_RESOURCES_USER/user-object}
// TODO
type ModifyGuildMember struct {
	// TODO
}

// Modify Current Member
// PATCH /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/members/@me
// TODO
type ModifyCurrentMember struct {
	// TODO
}

// Modify Current User Nick
// PATCH /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/members/@me/nick
// TODO
type ModifyCurrentUserNick struct {
	// TODO
}

// Add Guild Member Role
// PUT /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/members/{user.id#DOCS_RESOURCES_USER/user-object}/roles/{role.id#DOCS_TOPICS_PERMISSIONS/role-object}
// TODO
type AddGuildMemberRole struct {
	// TODO
}

// Remove Guild Member Role
// DELETE /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/members/{user.id#DOCS_RESOURCES_USER/user-object}/roles/{role.id#DOCS_TOPICS_PERMISSIONS/role-object}
// TODO
type RemoveGuildMemberRole struct {
	// TODO
}

// Remove Guild Member
// DELETE /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/members/{user.id#DOCS_RESOURCES_USER/user-object}
// TODO
type RemoveGuildMember struct {
	// TODO
}

// Get Guild Bans
// GET /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/bans
// TODO
type GetGuildBans struct {
	// TODO
}

// Get Guild Ban
// GET /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/bans/{user.id#DOCS_RESOURCES_USER/user-object}
// TODO
type GetGuildBan struct {
	// TODO
}

// Create Guild Ban
// PUT /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/bans/{user.id#DOCS_RESOURCES_USER/user-object}
// TODO
type CreateGuildBan struct {
	// TODO
}

// Remove Guild Ban
// DELETE /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/bans/{user.id#DOCS_RESOURCES_USER/user-object}
// TODO
type RemoveGuildBan struct {
	// TODO
}

// Get Guild Roles
// GET /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/roles
// TODO
type GetGuildRoles struct {
	// TODO
}

// Create Guild Role
// POST /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/roles
// TODO
type CreateGuildRole struct {
	// TODO
}

// Modify Guild Role Positions
// PATCH /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/roles
// TODO
type ModifyGuildRolePositions struct {
	// TODO
}

// Modify Guild Role
// PATCH /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/roles/{role.id#DOCS_TOPICS_PERMISSIONS/role-object}
// TODO
type ModifyGuildRole struct {
	// TODO
}

// Delete Guild Role
// DELETE /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/roles/{role.id#DOCS_TOPICS_PERMISSIONS/role-object}
// TODO
type DeleteGuildRole struct {
	// TODO
}

// Get Guild Prune Count
// GET /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/prune
// TODO
type GetGuildPruneCount struct {
	// TODO
}

// Begin Guild Prune
// POST /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/prune
// TODO
type BeginGuildPrune struct {
	// TODO
}

// Get Guild Voice Regions
// GET /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/regions
// TODO
type GetGuildVoiceRegions struct {
	// TODO
}

// Get Guild Invites
// GET /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/invites
// TODO
type GetGuildInvites struct {
	// TODO
}

// Get Guild Integrations
// GET /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/integrations
// TODO
type GetGuildIntegrations struct {
	// TODO
}

// Delete Guild Integration
// DELETE /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/integrations/{integration.id#DOCS_RESOURCES_GUILD/integration-object}
// TODO
type DeleteGuildIntegration struct {
	// TODO
}

// Get Guild Widget Settings
// GET /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/widget
// TODO
type GetGuildWidgetSettings struct {
	// TODO
}

// Modify Guild Widget
// PATCH /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/widget
// TODO
type ModifyGuildWidget struct {
	// TODO
}

// Get Guild Widget
// GET /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/widget.json
// TODO
type GetGuildWidget struct {
	// TODO
}

// Get Guild Vanity URL
// GET /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/vanity-url
// TODO
type GetGuildVanityURL struct {
	// TODO
}

// Get Guild Widget Image
// GET /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/widget.png
// TODO
type GetGuildWidgetImage struct {
	// TODO
}

// Get Guild Welcome Screen
// GET /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/welcome-screen
// TODO
type GetGuildWelcomeScreen struct {
	// TODO
}

// Modify Guild Welcome Screen
// PATCH /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/welcome-screen
// TODO
type ModifyGuildWelcomeScreen struct {
	// TODO
}

// Modify Current User Voice State
// PATCH /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/voice-states/@me
// TODO
type ModifyCurrentUserVoiceState struct {
	// TODO
}

// Modify User Voice State
// PATCH /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/voice-states/{user.id#DOCS_RESOURCES_USER/user-object}
// TODO
type ModifyUserVoiceState struct {
	// TODO
}

/// .go fileresources\Invite.md
// Get Invite
// GET /invites/{invite.code#DOCS_RESOURCES_INVITE/invite-object}
// TODO
type GetInvite struct {
	// TODO
}

// Delete Invite
// DELETE /invites/{invite.code#DOCS_RESOURCES_INVITE/invite-object}
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
// GET /stage-instances/{channel.id#DOCS_RESOURCES_CHANNEL/channel-object}
// TODO
type GetStageInstance struct {
	// TODO
}

// Modify Stage Instance
// PATCH /stage-instances/{channel.id#DOCS_RESOURCES_CHANNEL/channel-object}
// TODO
type ModifyStageInstance struct {
	// TODO
}

// Delete Stage Instance
// DELETE /stage-instances/{channel.id#DOCS_RESOURCES_CHANNEL/channel-object}
// TODO
type DeleteStageInstance struct {
	// TODO
}

/// .go fileresources\Sticker.md
// Get Sticker
// GET /stickers/{sticker.id#DOCS_RESOURCES_STICKER/sticker-object}
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
// GET /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/stickers
// TODO
type ListGuildStickers struct {
	// TODO
}

// Get Guild Sticker
// GET /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/stickers/{sticker.id#DOCS_RESOURCES_STICKER/sticker-object}
// TODO
type GetGuildSticker struct {
	// TODO
}

// Create Guild Sticker
// POST /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/stickers
// TODO
type CreateGuildSticker struct {
	// TODO
}

// Modify Guild Sticker
// PATCH /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/stickers/{sticker.id#DOCS_RESOURCES_STICKER/sticker-object}
// TODO
type ModifyGuildSticker struct {
	// TODO
}

// Delete Guild Sticker
// DELETE /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/stickers/{sticker.id#DOCS_RESOURCES_STICKER/sticker-object}
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
// GET /users/{user.id#DOCS_RESOURCES_USER/user-object}
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
// GET /users/@me/guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/member
// TODO
type GetCurrentUserGuildMember struct {
	// TODO
}

// Leave Guild
// DELETE /users/@me/guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}
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
// POST /channels/{channel.id#DOCS_RESOURCES_CHANNEL/channel-object}/webhooks
// TODO
type CreateWebhook struct {
	// TODO
}

// Get Channel Webhooks
// GET /channels/{channel.id#DOCS_RESOURCES_CHANNEL/channel-object}/webhooks
// TODO
type GetChannelWebhooks struct {
	// TODO
}

// Get Guild Webhooks
// GET /guilds/{guild.id#DOCS_RESOURCES_GUILD/guild-object}/webhooks
// TODO
type GetGuildWebhooks struct {
	// TODO
}

// Get Webhook
// GET /webhooks/{webhook.id#DOCS_RESOURCES_WEBHOOK/webhook-object}
// TODO
type GetWebhook struct {
	// TODO
}

// Get Webhook with Token
// GET /webhooks/{webhook.id#DOCS_RESOURCES_WEBHOOK/webhook-object}/{webhook.token#DOCS_RESOURCES_WEBHOOK/webhook-object}
// TODO
type GetWebhookwithToken struct {
	// TODO
}

// Modify Webhook
// PATCH /webhooks/{webhook.id#DOCS_RESOURCES_WEBHOOK/webhook-object}
// TODO
type ModifyWebhook struct {
	// TODO
}

// Modify Webhook with Token
// PATCH /webhooks/{webhook.id#DOCS_RESOURCES_WEBHOOK/webhook-object}/{webhook.token#DOCS_RESOURCES_WEBHOOK/webhook-object}
// TODO
type ModifyWebhookwithToken struct {
	// TODO
}

// Delete Webhook
// DELETE /webhooks/{webhook.id#DOCS_RESOURCES_WEBHOOK/webhook-object}
// TODO
type DeleteWebhook struct {
	// TODO
}

// Delete Webhook with Token
// DELETE /webhooks/{webhook.id#DOCS_RESOURCES_WEBHOOK/webhook-object}/{webhook.token#DOCS_RESOURCES_WEBHOOK/webhook-object}
// TODO
type DeleteWebhookwithToken struct {
	// TODO
}

// Execute Webhook
// POST /webhooks/{webhook.id#DOCS_RESOURCES_WEBHOOK/webhook-object}/{webhook.token#DOCS_RESOURCES_WEBHOOK/webhook-object}
// TODO
type ExecuteWebhook struct {
	// TODO
}

// Execute Slack-Compatible Webhook
// POST /webhooks/{webhook.id#DOCS_RESOURCES_WEBHOOK/webhook-object}/{webhook.token#DOCS_RESOURCES_WEBHOOK/webhook-object}/slack
// TODO
type ExecuteSlackCompatibleWebhook struct {
	// TODO
}

// Execute GitHub-Compatible Webhook
// POST /webhooks/{webhook.id#DOCS_RESOURCES_WEBHOOK/webhook-object}/{webhook.token#DOCS_RESOURCES_WEBHOOK/webhook-object}/github
// TODO
type ExecuteGitHubCompatibleWebhook struct {
	// TODO
}

// Get Webhook Message
// GET /webhooks/{webhook.id#DOCS_RESOURCES_WEBHOOK/webhook-object}/{webhook.token#DOCS_RESOURCES_WEBHOOK/webhook-object}/messages/{message.id#DOCS_RESOURCES_CHANNEL/message-object}
// TODO
type GetWebhookMessage struct {
	// TODO
}

// Edit Webhook Message
// PATCH /webhooks/{webhook.id#DOCS_RESOURCES_WEBHOOK/webhook-object}/{webhook.token#DOCS_RESOURCES_WEBHOOK/webhook-object}/messages/{message.id#DOCS_RESOURCES_CHANNEL/message-object}
// TODO
type EditWebhookMessage struct {
	// TODO
}

// Delete Webhook Message
// DELETE /webhooks/{webhook.id#DOCS_RESOURCES_WEBHOOK/webhook-object}/{webhook.token#DOCS_RESOURCES_WEBHOOK/webhook-object}/messages/{message.id#DOCS_RESOURCES_CHANNEL/message-object}
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
