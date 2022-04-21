package requests

import (
	"time"

	"github.com/switchupcb/disgo/wrapper/resources"
)

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
	GuildID resources.Snowflake
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
// https://discord.com/developers/docs/interactions/application-commands#delete-guild-application-command
type DeleteGuildApplicationCommand struct {
	GuildID resources.Snowflake
}

// Bulk Overwrite Guild Application Commands
// PUT /applications/{application.id}/guilds/{guild.id}/commands
// https://discord.com/developers/docs/interactions/application-commands#bulk-overwrite-guild-application-commands
type BulkOverwriteGuildApplicationCommands struct {
	CommandID                resources.Snowflake
	Name                     string                                `json:"name,omitempty"`
	NameLocalizations        map[resources.Flag]string             `json:"name_localizations,omitempty"`
	Description              string                                `json:"description,omitempty"`
	DescriptionLocalizations map[resources.Flag]string             `json:"description_localizations,omitempty"`
	Options                  []*resources.ApplicationCommandOption `json:"options,omitempty"`
	DefaultPermission        bool                                  `json:"default_permission,omitempty"`
	Type                     resources.Flag                        `json:"type,omitempty"`
}

// Get Guild Application Command Permissions
// GET /applications/{application.id}/guilds/{guild.id}/commands/permissions
// https://discord.com/developers/docs/interactions/application-commands#get-guild-application-command-permissions
type GetGuildApplicationCommandPermissions struct {
	GuildID resources.Snowflake
}

// Get Application Command Permissions
// GET /applications/{application.id}/guilds/{guild.id}/commands/{command.id}/permissions
// https://discord.com/developers/docs/interactions/application-commands#get-application-command-permissions
type GetApplicationCommandPermissions struct {
	GuildID resources.Snowflake
}

// Edit Application Command Permissions
// PUT /applications/{application.id}/guilds/{guild.id}/commands/{command.id}/permissions
// https://discord.com/developers/docs/interactions/application-commands#edit-application-command-permissions
type EditApplicationCommandPermissions struct {
	Permissions []*resources.ApplicationCommandPermissions `json:"permissions,omitempty"`
}

// Batch Edit Application Command Permissions
// PUT /applications/{application.id}/guilds/{guild.id}/commands/permissions
// https://discord.com/developers/docs/interactions/application-commands#batch-edit-application-command-permissions
type BatchEditApplicationCommandPermissions struct {
	GuildID resources.Snowflake
}

/// .go fileinteractions\Receiving_and_Responding.md
// Create Interaction Response
// POST /interactions/{interaction.id}/{interaction.token}/callback
// https://discord.com/developers/docs/interactions/receiving-and-responding#create-interaction-response
type CreateInteractionResponse struct {
	// TODO
}

// Get Original Interaction Response
// GET /webhooks/{application.id}/{interaction.token}/messages/@original
// https://discord.com/developers/docs/interactions/receiving-and-responding#get-original-interaction-response
type GetOriginalInteractionResponse struct {
	// TODO
}

// Edit Original Interaction Response
// PATCH /webhooks/{application.id}/{interaction.token}/messages/@original
// https://discord.com/developers/docs/interactions/receiving-and-responding#edit-original-interaction-response
type EditOriginalInteractionResponse struct {
	// TODO
}

// Delete Original Interaction Response
// DELETE /webhooks/{application.id}/{interaction.token}/messages/@original
// https://discord.com/developers/docs/interactions/receiving-and-responding#delete-original-interaction-response
type DeleteOriginalInteractionResponse struct {
	// TODO
}

// Create Followup Message
// POST /webhooks/{application.id}/{interaction.token}
// https://discord.com/developers/docs/interactions/receiving-and-responding#create-followup-message
type CreateFollowupMessage struct {
	// TODO
}

// Get Followup Message
// GET /webhooks/{application.id}/{interaction.token}/messages/{message.id}
// https://discord.com/developers/docs/interactions/receiving-and-responding#get-followup-message
type GetFollowupMessage struct {
	// TODO
}

// Edit Followup Message
// PATCH /webhooks/{application.id}/{interaction.token}/messages/{message.id}
// https://discord.com/developers/docs/interactions/receiving-and-responding#edit-followup-message
type EditFollowupMessage struct {
	// TODO
}

// Delete Followup Message
// DELETE /webhooks/{application.id}/{interaction.token}/messages/{message.id}
// https://discord.com/developers/docs/interactions/receiving-and-responding#delete-followup-message
type DeleteFollowupMessage struct {
	// TODO
}

/// .go fileresources\Audit_Log.md
// Get Guild Audit Log
// GET /guilds/{guild.id}/audit-logs
// https://discord.com/developers/docs/resources/audit-log#get-guild-audit-log
type GetGuildAuditLog struct {
	UserID     resources.Snowflake `urlparam:"user_id"`
	ActionType int                 `urlparam:"action_type"`
	Before     resources.Snowflake `urlparam:"before,omitempty"`
	Limit      int                 `urlparam:"limit,omitempty"`
}

/// .go fileresources\Channel.md
// Get Channel
// GET /channels/{channel.id}
// https://discord.com/developers/docs/resources/channel#get-channel
type GetChannel struct {
	// TODO
}

// Modify Channel
// PATCH /channels/{channel.id}
// https://discord.com/developers/docs/resources/channel#modify-channel
type ModifyChannel struct {
	// TODO
}

// Modify Channel
// PATCH /channels/{channel.id}
// https://discord.com/developers/docs/resources/channel#modify-channel-json-params-group-dm
type ModifyChannelGroupDM struct {
	Name string `json:"name,omitempty"`
	Icon int    `json:"icon,omitempty"`
}

// Modify Channel
// PATCH /channels/{channel.id}
// https://discord.com/developers/docs/resources/channel#modify-channel-json-params-guild-channel
type ModifyChannelGuild struct {
	//TODO
	Name                       *string                          `json:"name,omitempty"`
	Type                       *resources.Flag                  `json:"type,omitempty"`
	Position                   *uint                            `json:"position,omitempty"`
	Topic                      *string                          `json:"topic,omitempty"`
	NSFW                       *bool                            `json:"nsfw,omitempty"`
	RateLimitPerUser           *uint                            `json:"rate_limit_per_user,omitempty"`
	Bitrate                    *uint                            `json:"bitrate,omitempty"`
	UserLimit                  *uint                            `json:"user_limit,omitempty"`
	PermissionOverwrites       *[]resources.PermissionOverwrite `json:"permission_overwrites,omitempty"`
	ParentID                   *resources.Snowflake             `json:"parent_id,omitempty"`
	RTCRegion                  *string                          `json:"rtc_region,omitempty"`
	VideoQualityMode           resources.Flag                   `json:"video_quality_mode,omitempty"`
	DefaultAutoArchiveDuration *uint                            `json:"default_auto_archive_duration,omitempty"`
}

// Modify Channel
// PATCH /channels/{channel.id}
// https://discord.com/developers/docs/resources/channel#modify-channel-json-params-thread
type ModifyChannelThread struct {
	Name                string `json:"name,omitempty"`
	Archived            bool   `json:"archived,omitempty"`
	AutoArchiveDuration int    `json:"auto_archive_duration,omitempty"`
	Locked              bool   `json:"locked,omitempty"`
	Invitable           bool   `json:"invitable,omitempty"`
}

// Delete/Close Channel
// DELETE /channels/{channel.id}
// https://discord.com/developers/docs/resources/channel#deleteclose-channel
type DeleteCloseChannel struct {
	// TODO
}

// Get Channel Messages
// GET /channels/{channel.id}/messages
// https://discord.com/developers/docs/resources/channel#get-channel-messages
type GetChannelMessages struct {
	Around *resources.Snowflake `urlparam:"around,omitempty"`
	Before *resources.Snowflake `urlparam:"before,omitempty"`
	After  *resources.Snowflake `urlparam:"after,omitempty"`
	Limit  uint                 `urlparam:"limit,omitempty"`
}

// Get Channel Message
// GET /channels/{channel.id}/messages/{message.id}
// https://discord.com/developers/docs/resources/channel#get-channel-message
type GetChannelMessage struct {
	// TODO
}

// Create Message
// POST /channels/{channel.id}/messages
// https://discord.com/developers/docs/resources/channel#create-message
type CreateMessage struct {
	// TODO: Remove File/Files this when compatibility is not required.
	Content         string                      `json:"content,omitempty"`
	TTS             bool                        `json:"tts,omitempty"`
	Embeds          []*resources.Embed          `json:"embeds,omitempty"`
	Embed           *resources.Embed            `json:"-"`
	AllowedMentions *resources.AllowedMentions  `json:"allowed_mentions,omitempty"`
	Reference       *resources.MessageReference `json:"message_reference,omitempty"`
	StickerID       []*resources.Snowflake      `json:"sticker_ids,omitempty"`
	Components      []resources.Component       `json:"components,omitempty"`
	Files           []byte                      `disgo:"TODO"`
}

// Crosspost Message
// POST /channels/{channel.id}/messages/{message.id}/crosspost
// https://discord.com/developers/docs/resources/channel#crosspost-message
type CrosspostMessage struct {
	MessageID *resources.Snowflake
}

// Create Reaction
// PUT /channels/{channel.id}/messages/{message.id}/reactions/{emoji}/@me
// https://discord.com/developers/docs/resources/channel#create-reaction
type CreateReaction struct {
	MessageID *resources.Snowflake
}

// Delete Own Reaction
// DELETE /channels/{channel.id}/messages/{message.id}/reactions/{emoji}/@me
// https://discord.com/developers/docs/resources/channel#delete-own-reaction
type DeleteOwnReaction struct {
	MessageID *resources.Snowflake
}

// Delete User Reaction
// DELETE /channels/{channel.id}/messages/{message.id}/reactions/{emoji}/{user.id}
// https://discord.com/developers/docs/resources/channel#delete-user-reaction
type DeleteUserReaction struct {
	MessageID *resources.Snowflake
}

// Get Reactions
// GET /channels/{channel.id}/messages/{message.id}/reactions/{emoji}
// https://discord.com/developers/docs/resources/channel#get-reactions
type GetReactions struct {
	After *resources.Snowflake `json:"after,omitempty"`
	Limit int                  `json:"limit,omitempty"` // 1 is default. even if 0 is supplied.
}

// Delete All Reactions
// DELETE /channels/{channel.id}/messages/{message.id}/reactions
// https://discord.com/developers/docs/resources/channel#delete-all-reactions
type DeleteAllReactions struct {
	MessageID *resources.Snowflake
}

// Delete All Reactions for Emoji
// DELETE /channels/{channel.id}/messages/{message.id}/reactions/{emoji}
// https://discord.com/developers/docs/resources/channel#delete-all-reactions-for-emoji
type DeleteAllReactionsforEmoji struct {
	MessageID *resources.Snowflake
}

// Edit Message
// PATCH /channels/{channel.id}/messages/{message.id}
// https://discord.com/developers/docs/resources/channel#edit-message
type EditMessage struct {
	Content *string            `json:"content,omitempty"`
	Embeds  []*resources.Embed `json:"embeds,omitempty"`
	Flags   *resources.BitFlag `json:"flags,omitempty"`
	// PayloadJSON *string `json:"payload_json,omitempty"`
	AllowedMentions *resources.AllowedMentions `json:"allowed_mentions,omitempty"`
	Components      []*resources.Component     `json:"components,omitempty"`
}

// Delete Message
// DELETE /channels/{channel.id}/messages/{message.id}
// https://discord.com/developers/docs/resources/channel#delete-message
type DeleteMessage struct {
	MessageID *resources.Snowflake
}

// Bulk Delete Messages
// POST /channels/{channel.id}/messages/bulk-delete
// https://discord.com/developers/docs/resources/channel#bulk-delete-messages
type BulkDeleteMessages struct {
	Messages resources.Snowflake `json:"messages,omitempty"`
}

// Edit Channel Permissions
// PUT /channels/{channel.id}/permissions/{overwrite.id}
// https://discord.com/developers/docs/resources/channel#edit-channel-permissions
type EditChannelPermissions struct {
	Allow string          `json:"allow,omitempty"`
	Deny  string          `json:"deny,omitempty"`
	Type  *resources.Flag `json:"type,omitempty"`
}

// Get Channel Invites
// GET /channels/{channel.id}/invites
// https://discord.com/developers/docs/resources/channel#get-channel-invites
type GetChannelInvites struct {
	ChannelID resources.Snowflake
}

// Create Channel Invite
// POST /channels/{channel.id}/invites
// https://discord.com/developers/docs/resources/channel#create-channel-invite
type CreateChannelInvite struct {
	MaxAge              int                 `json:"max_age"`
	MaxUses             int                 `json:"max_uses,omitempty"`
	Temporary           bool                `json:"temporary,omitempty"`
	Unique              bool                `json:"unique,omitempty"`
	TargetType          int                 `json:"target_type,omitempty"`
	TargetUserID        resources.Snowflake `json:"target_user_id,omitempty"`
	TargetApplicationID resources.Snowflake `json:"target_application_id,omitempty"`
}

// Delete Channel Permission
// DELETE /channels/{channel.id}/permissions/{overwrite.id}
// https://discord.com/developers/docs/resources/channel#delete-channel-permission
type DeleteChannelPermission struct {
	OverwriteID resources.Snowflake
}

// Follow News Channel
// POST /channels/{channel.id}/followers
// https://discord.com/developers/docs/resources/channel#follow-news-channel
type FollowNewsChannel struct {
	WebhookChannelID resources.Snowflake
}

// Trigger Typing Indicator
// POST /channels/{channel.id}/typing
// https://discord.com/developers/docs/resources/channel#trigger-typing-indicator
type TriggerTypingIndicator struct {
	ChannelID resources.Snowflake
}

// Get Pinned Messages
// GET /channels/{channel.id}/pins
// https://discord.com/developers/docs/resources/channel#get-pinned-messages
type GetPinnedMessages struct {
	ChannelID resources.Snowflake
}

// Pin Message
// PUT /channels/{channel.id}/pins/{message.id}
// https://discord.com/developers/docs/resources/channel#pin-message
type PinMessage struct {
	MessageID resources.Snowflake
}

// Unpin Message
// DELETE /channels/{channel.id}/pins/{message.id}
// https://discord.com/developers/docs/resources/channel#unpin-message
type UnpinMessage struct {
	MessageID resources.Snowflake
}

// Group DM Add Recipient
// PUT /channels/{channel.id}/recipients/{user.id}
// https://discord.com/developers/docs/resources/channel#group-dm-add-recipient
type GroupDMAddRecipient struct {
	AccessToken string `json:"access_token"`
	Nickname    string `json:"nick,omitempty"`
}

// Group DM Remove Recipient
// DELETE /channels/{channel.id}/recipients/{user.id}
// https://discord.com/developers/docs/resources/channel#group-dm-remove-recipient
type GroupDMRemoveRecipient struct {
	UserID resources.Snowflake
}

// Start Thread from Message
// POST /channels/{channel.id}/messages/{message.id}/threads
// https://discord.com/developers/docs/resources/channel#start-thread-from-message
type StartThreadfromMessage struct {
	Name                *string `json:"name,omitempty"`
	RateLimitPerUser    *uint   `json:"rate_limit_per_user,omitempty"`
	AutoArchiveDuration int     `json:"auto_archive_duration,omitempty"`
}

// Start Thread without Message
// POST /channels/{channel.id}/threads
// https://discord.com/developers/docs/resources/channel#start-thread-without-message
type StartThreadwithoutMessage struct {
	Name                *string         `json:"name,omitempty"`
	AutoArchiveDuration int             `json:"auto_archive_duration,omitempty"`
	Type                *resources.Flag `json:"type,omitempty"`
	Invitable           bool            `json:"invitable,omitempty"`
	RateLimitPerUser    *uint           `json:"rate_limit_per_user,omitempty"`
}

// Start Thread in Forum Channel
// POST /channels/{channel.id}/threads
// https://discord.com/developers/docs/resources/channel#start-thread-in-forum-channel
type StartThreadinForumChannel struct {
	ChannelID resources.Snowflake
}

// Start Thread in Forum Channel
// POST /channels/{channel.id}/threads
// https://discord.com/developers/docs/resources/channel#start-thread-in-forum-channel-json-params-for-the-thread
type StartThreadinForumChannelThread struct {
	Name                *string `json:"name,omitempty"`
	RateLimitPerUser    *uint   `json:"rate_limit_per_user,omitempty"`
	AutoArchiveDuration int     `json:"auto_archive_duration,omitempty"`
}

// Start Thread in Forum Channel
// POST /channels/{channel.id}/threads
// https://discord.com/developers/docs/resources/channel#start-thread-in-forum-channel-json-params-for-the-message
type StartThreadinForumChannelMessage struct {
	Content         string                     `json:"content"`
	Embeds          []*resources.Embed         `json:"embeds,omitempty"`
	AllowedMentions *resources.AllowedMentions `json:"allowed_mentions,omitempty"`
	Components      []*resources.Component     `json:"components,omitempty"`
	StickerIDS      []*resources.Snowflake     `json:"sticker_ids,omitempty"`
	Attachments     []*resources.Attachment    `json:"attachments,omitempty"`
	//Files			 []resources. `json:"-"`
	PayloadJSON string            `json:"payload_json, omitempty"`
	Flags       resources.BitFlag `json:"flags,omitempty"`
}

// Join Thread
// PUT /channels/{channel.id}/thread-members/@me
// https://discord.com/developers/docs/resources/channel#join-thread
type JoinThread struct {
	ChannelID resources.Snowflake
}

// Add Thread Member
// PUT /channels/{channel.id}/thread-members/{user.id}
// https://discord.com/developers/docs/resources/channel#add-thread-member
type AddThreadMember struct {
	UserID resources.Snowflake
}

// Leave Thread
// DELETE /channels/{channel.id}/thread-members/@me
// https://discord.com/developers/docs/resources/channel#leave-thread
type LeaveThread struct {
	ChannelID resources.Snowflake
}

// Remove Thread Member
// DELETE /channels/{channel.id}/thread-members/{user.id}
// https://discord.com/developers/docs/resources/channel#remove-thread-member
type RemoveThreadMember struct {
	UserID resources.Snowflake
}

// Get Thread Member
// GET /channels/{channel.id}/thread-members/{user.id}
// https://discord.com/developers/docs/resources/channel#get-thread-member
type GetThreadMember struct {
	UserID resources.Snowflake
}

// List Thread Members
// GET /channels/{channel.id}/thread-members
// https://discord.com/developers/docs/resources/channel#list-thread-members
type ListThreadMembers struct {
	ChannelID resources.Snowflake
}

// List Active Channel Threads
// GET /channels/{channel.id}/threads/active
// https://discord.com/developers/docs/resources/channel#list-active-threads
type ListActiveChannelThreads struct {
	Threads []*resources.Channel      `json:"threads"`
	Members []*resources.ThreadMember `json:"members"`
	HasMore bool                      `json:"has_more"`
	Before  resources.Snowflake       `urlparam:"before,omitempty"`
	Limit   int                       `urlparam:"limit,omitempty"`
}

// List Public Archived Threads
// GET /channels/{channel.id}/threads/archived/public
// https://discord.com/developers/docs/resources/channel#list-public-archived-threads
type ListPublicArchivedThreads struct {
	Threads []*resources.Channel      `json:"threads"`
	Members []*resources.ThreadMember `json:"members"`
	HasMore bool                      `json:"has_more"`
	Before  resources.Snowflake       `urlparam:"before,omitempty"`
	Limit   int                       `urlparam:"limit,omitempty"`
}

// List Private Archived Threads
// GET /channels/{channel.id}/threads/archived/private
// https://discord.com/developers/docs/resources/channel#list-private-archived-threads
type ListPrivateArchivedThreads struct {
	Threads []*resources.Channel      `json:"threads"`
	Members []*resources.ThreadMember `json:"members"`
	HasMore bool                      `json:"has_more"`
	Before  resources.Snowflake       `urlparam:"before,omitempty"`
	Limit   int                       `urlparam:"limit,omitempty"`
}

// List Joined Private Archived Threads
// GET /channels/{channel.id}/users/@me/threads/archived/private
// https://discord.com/developers/docs/resources/channel#list-joined-private-archived-threads
type ListJoinedPrivateArchivedThreads struct {
	Threads []*resources.Channel      `json:"threads"`
	Members []*resources.ThreadMember `json:"members"`
	HasMore bool                      `json:"has_more"`
	Before  resources.Snowflake       `urlparam:"before,omitempty"`
	Limit   int                       `urlparam:"limit,omitempty"`
}

/// .go fileresources\Emoji.md
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
	EmojiID resources.Snowflake
}

// Create Guild Emoji
// POST /guilds/{guild.id}/emojis
// https://discord.com/developers/docs/resources/emoji#create-guild-emoji
type CreateGuildEmoji struct {
	GuildID resources.Snowflake
	Name    string                 `name:"name"`
	Image   string                 `image:"name"`
	Roles   []*resources.Snowflake `roles:"name"`
}

// Modify Guild Emoji
// PATCH /guilds/{guild.id}/emojis/{emoji.id}
// https://discord.com/developers/docs/resources/emoji#modify-guild-emoji
type ModifyGuildEmoji struct {
	EmojiID resources.Snowflake
	Name    string                 `json:"name"`
	Roles   []*resources.Snowflake `json:"roles"`
}

// Delete Guild Emoji
// DELETE /guilds/{guild.id}/emojis/{emoji.id}
// https://discord.com/developers/docs/resources/emoji#delete-guild-emoji
type DeleteGuildEmoji struct {
	EmojiID resources.Snowflake
}

// List Scheduled Events for Guild
// GET /guilds/{guild.id}/scheduled-events
// https://discord.com/developers/docs/resources/guild-scheduled-event#list-scheduled-events-for-guild
type ListScheduledEventsforGuild struct {
	WithUserCount bool `json:"with_user_count"`
}

// Create Guild Scheduled Event
// POST /guilds/{guild.id}/scheduled-events
// https://discord.com/developers/docs/resources/guild-scheduled-event#create-guild-scheduled-event
type CreateGuildScheduledEvent struct {
	ChannelID          *resources.Snowflake                         `json:"channel_id"`
	EntityMetadata     *resources.GuildScheduledEventEntityMetadata `json:"entity_metadata,omitempty"`
	Name               *string                                      `json:"name,omitempty"`
	PrivacyLevel       resources.Flag                               `json:"privacy_level,omitempty"`
	ScheduledStartTime resources.Snowflake                          `json:"scheduled_start_time,omitempty"`
	ScheduledEndTime   resources.Snowflake                          `json:"scheduled_end_time,omitempty"`
	Description        *string                                      `json:"description,omitempty"`
	EntityType         *resources.Flag                              `json:"entity_type,omitempty"`
	Image              string                                       `json:"image"`
}

// Get Guild Scheduled Event
// GET /guilds/{guild.id}/scheduled-events/{guild_scheduled_event.id}
// https://discord.com/developers/docs/resources/guild-scheduled-event#get-guild-scheduled-event
type GetGuildScheduledEvent struct {
	WithUserCount bool `json:"with_user_count"`
}

// Modify Guild Scheduled Event
// PATCH /guilds/{guild.id}/scheduled-events/{guild_scheduled_event.id}
// https://discord.com/developers/docs/resources/guild-scheduled-event#modify-guild-scheduled-event
type ModifyGuildScheduledEvent struct {
	ChannelID          *resources.Snowflake                         `json:"channel_id"`
	EntityMetadata     *resources.GuildScheduledEventEntityMetadata `json:"entity_metadata,omitempty"`
	Name               *string                                      `json:"name,omitempty"`
	PrivacyLevel       resources.Flag                               `json:"privacy_level,omitempty"`
	ScheduledStartTime resources.Snowflake                          `json:"scheduled_start_time,omitempty"`
	ScheduledEndTime   resources.Snowflake                          `json:"scheduled_end_time,omitempty"`
	Description        *string                                      `json:"description,omitempty"`
	EntityType         *resources.Flag                              `json:"entity_type,omitempty"`
	Image              string                                       `json:"image"`
	Status             *resources.Flag                              `json:"status,omitempty"`
}

// Delete Guild Scheduled Event
// DELETE /guilds/{guild.id}/scheduled-events/{guild_scheduled_event.id}
// https://discord.com/developers/docs/resources/guild-scheduled-event#delete-guild-scheduled-event
type DeleteGuildScheduledEvent struct {
	GuildScheduledEventID resources.Snowflake
}

// Get Guild Scheduled Event Users
// GET /guilds/{guild.id}/scheduled-events/{guild_scheduled_event.id}/users
// https://discord.com/developers/docs/resources/guild-scheduled-event#get-guild-scheduled-event-users
type GetGuildScheduledEventUsers struct {
	Limit      resources.Flag      `urlparam:"limit,omitempty"`
	WithMember bool                `urlparam:"with_member,omitempty"`
	Before     resources.Snowflake `urlparam:"before,omitempty"`
	After      resources.Snowflake `urlparam:"after,omitempty"`
}

// Get Guild Template
// GET /guilds/templates/{template.code}
// https://discord.com/developers/docs/resources/guild-template#get-guild-template
type GetGuildTemplate struct {
	// TODO
}

// Create Guild from Guild Template
// POST /guilds/templates/{template.code}
// https://discord.com/developers/docs/resources/guild-template#create-guild-from-guild-template
type CreateGuildfromGuildTemplate struct {
	Name string `json:"name,omitempty"`
	Icon string `json:"icon,omitempty"`
}

// Get Guild Templates
// GET /guilds/{guild.id}/templates
// https://discord.com/developers/docs/resources/guild-template#get-guild-templates
type GetGuildTemplates struct {
	// TODO
}

// Create Guild Template
// POST /guilds/{guild.id}/templates
// https://discord.com/developers/docs/resources/guild-template#create-guild-template
type CreateGuildTemplate struct {
	Name        string  `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

// Sync Guild Template
// PUT /guilds/{guild.id}/templates/{template.code}
// https://discord.com/developers/docs/resources/guild-template#sync-guild-template
type SyncGuildTemplate struct {
	// TODO
}

// Modify Guild Template
// PATCH /guilds/{guild.id}/templates/{template.code}
// https://discord.com/developers/docs/resources/guild-template#modify-guild-template
type ModifyGuildTemplate struct {
	Name        string  `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

// Delete Guild Template
// DELETE /guilds/{guild.id}/templates/{template.code}
// https://discord.com/developers/docs/resources/guild-template#delete-guild-template
type DeleteGuildTemplate struct {
	// TODO
}

/// .go fileresources\Guild.md
// Create Guild
// POST /guilds
// https://discord.com/developers/docs/resources/guild#create-guild
type CreateGuild struct {
	Name   string `json:"name,omitempty"`
	Region string `json:"region,omitempty"`
	Icon   string `json:"icon,omitempty"`
	//VerificationLevel           *resources.VerificationLevel `json:"verification_level,omitempty"`
	DefaultMessageNotifications int                  `json:"default_message_notifications,omitempty"` // TODO: Separate type?
	AfkChannelID                string               `json:"afk_channel_id,omitempty"`
	AfkTimeout                  int                  `json:"afk_timeout,omitempty"`
	OwnerID                     string               `json:"owner_id,omitempty"`
	Splash                      string               `json:"splash,omitempty"`
	Banner                      string               `json:"banner,omitempty"`
	Roles                       []*resources.Role    `json:"roles,omitempty"`
	Channels                    []*resources.Channel `json:"channels,omitempty"`
	SystemChannelID             *resources.Snowflake `json:"system_channel_id,omitempty"`
	SystemChannelFlags          resources.BitFlag    `json:"system_channel_flags,omitempty"`
	//ExplicitContentFilter       *resources.ExplicitContentFilterLevel `json:"explicit_content_filter,omitempty"`
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
	GuildID *resources.Snowflake
	Name    *string `json:"name,omitempty"`
	Region  *string `json:"region,omitempty"`
	//VerificationLvl             *resources.VerificationLevel             `json:"verification_lvl,omitempty"`
	//DefaultMessageNotifications *resources.DefaultMessageNotificationLvl `json:"default_message_notifications,omitempty"`
	//ExplicitContentFilter       *resources.ExplicitContentFilterLevel    `json:"explicit_content_filter,omitempty"`
	AFKChannelID           *resources.Snowflake `json:"afk_channel_id,omitempty"`
	Icon                   *string              `json:"icon,omitempty"`
	OwnerID                *resources.Snowflake `json:"owner_id,omitempty"`
	Splash                 *string              `json:"splash,omitempty"`
	DiscoverySplash        *string              `json:"discovery_splash,omitempty"`
	Banner                 *string              `json:"banner,omitempty"`
	SystemChannelID        *resources.Snowflake `json:"system_channel_id,omitempty"`
	SystemChannelFlags     resources.BitFlag    `json:"system_channel_flags,omitempty"`
	RulesChannelID         *resources.Snowflake `json:"rules_channel_id,omitempty"`
	PublicUpdatesChannelID *resources.Snowflake `json:"public_updates_channel_id,omitempty"`
	PreferredLocale        *string              `json:"preferred_locale,omitempty"`
	Features               *[]string            `json:"features,omitempty"`
	Description            *string              `json:"description,omitempty"`
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
	Position                   int                              `json:"position"`
	Bitrate                    int                              `json:"bitrate,omitempty"`
	UserLimit                  int                              `json:"user_limit,omitempty"`
	PermissionOverwrites       []*resources.PermissionOverwrite `json:"permission_overwrites,omitempty"`
	ParentID                   *resources.Snowflake             `json:"parent_id,omitempty"`
	RateLimitPerUser           *int                             `json:"rate_limit_per_user,omitempty"`
	DefaultAutoArchiveDuration int                              `json:"default_auto_archive_duration,omitempty"`
}

// Modify Guild Channel Positions
// PATCH /guilds/{guild.id}/channels
// https://discord.com/developers/docs/resources/guild#modify-guild-channel-positions
type ModifyGuildChannelPositions struct {
	ID              *resources.Snowflake `json:"id,omitempty"`
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
	After *resources.Snowflake `urlparam:"after,omitempty"`
	Limit int                  `urlparam:"limit,omitempty"`
}

// Search Guild Members
// GET /guilds/{guild.id}/members/search
// https://discord.com/developers/docs/resources/guild#search-guild-members
type SearchGuildMembers struct {
	GuildID *resources.Snowflake
	Query   string `json:"query"`
	Limit   int    `urlparam:"limit,omitempty"`
}

// Add Guild Member
// PUT /guilds/{guild.id}/members/{user.id}
// https://discord.com/developers/docs/resources/guild#add-guild-member
type AddGuildMember struct {
	UserID      *resources.Snowflake
	AccessToken string                 `json:"access_token"`
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
	ChannelID                  resources.Snowflake    `json:"channel_id"`
	CommunicationDisabledUntil *time.Time             `json:"communication_disabled_until"`
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
	Before  *resources.Snowflake `urlparam:"before,omitempty"`
	After   *resources.Snowflake `urlparam:"after,omitempty"`
	Limit   int                  `urlparam:"limit,omitempty"`
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
	DeleteMessageDays *resources.Flag `urlparam:"delete_message_days,omitempty"`
	Reason            *string         `urlparam:"reason,omitempty"`
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
	Name         string `json:"name"`
	Permissions  int64  `json:"permissions,string"`
	Color        int    `json:"color"`
	Hoist        bool   `json:"hoist"`
	Icon         int    `json:"icon"`
	UnicodeEmoji string `json:"unicode_emoji"`
	Mentionable  bool   `json:"mentionable"`
}

// Modify Guild Role Positions
// PATCH /guilds/{guild.id}/roles
// https://discord.com/developers/docs/resources/guild#modify-guild-role-positions
type ModifyGuildRolePositions struct {
	ID       resources.Snowflake `json:"id"`
	Position int                 `json:"position"`
}

// Modify Guild Role
// PATCH /guilds/{guild.id}/roles/{role.id}
// https://discord.com/developers/docs/resources/guild#modify-guild-role
type ModifyGuildRole struct {
	RoleID       resources.Snowflake
	Name         string `json:"name"`
	Permissions  int64  `json:"permissions,string"`
	Color        int    `json:"color"`
	Hoist        bool   `json:"hoist"`
	Icon         int    `json:"icon"`
	UnicodeEmoji string `json:"unicode_emoji"`
	Mentionable  bool   `json:"mentionable"`
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
	Days         int                    `json:"days,omitempty"`
	IncludeRoles []*resources.Snowflake `json:"include_roles,omitempty"`
}

// Begin Guild Prune
// POST /guilds/{guild.id}/prune
// https://discord.com/developers/docs/resources/guild#begin-guild-prune
type BeginGuildPrune struct {
	GuildID           resources.Snowflake
	Days              *int                   `json:"days,omitempty"`
	ComputePruneCount *bool                  `json:"compute_prune_count,omitempty"`
	IncludeRoles      []*resources.Snowflake `json:"include_roles,omitempty"`
	Reason            string                 `json:"reason"`
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
	GuildID         resources.Snowflake
	Enabled         *bool                             `json:"enabled,omitempty"`
	WelcomeChannels []*resources.WelcomeScreenChannel `json:"welcome_channels,omitempty"`
	Description     *string                           `json:"description,omitempty"`
}

// Modify Guild Welcome Screen
// PATCH /guilds/{guild.id}/welcome-screen
// https://discord.com/developers/docs/resources/guild#modify-guild-welcome-screen
type ModifyGuildWelcomeScreen struct {
	GuildID         resources.Snowflake
	Enabled         *bool                             `json:"enabled,omitempty"`
	WelcomeChannels []*resources.WelcomeScreenChannel `json:"welcome_channels,omitempty"`
	Description     *string                           `json:"description,omitempty"`
}

// Modify Current User Voice State
// PATCH /guilds/{guild.id}/voice-states/@me
// https://discord.com/developers/docs/resources/guild#modify-current-user-voice-state
type ModifyCurrentUserVoiceState struct {
	ChannelID               resources.Snowflake
	Suppress                bool       `json:"suppress"`
	RequestToSpeakTimestamp *time.Time `json:"request_to_speak_timestamp"`
}

// Modify User Voice State
// PATCH /guilds/{guild.id}/voice-states/{user.id}
// https://discord.com/developers/docs/resources/guild#modify-user-voice-state
type ModifyUserVoiceState struct {
	ChannelID resources.Snowflake
	Suppress  bool `json:"suppress"`
}

/// .go fileresources\Invite.md
// Get Invite
// GET /invites/{invite.code}
// https://discord.com/developers/docs/resources/invite#get-invite
type GetInvite struct {
	GuildScheduledEventID resources.Snowflake `json:"guild_scheduled_event_id"`
	WithCounts            bool                `json:"with_counts,omitempty"`
	WithExpiration        bool                `json:"with_expiration,omitempty"`
}

// Delete Invite
// DELETE /invites/{invite.code}
// https://discord.com/developers/docs/resources/invite#delete-invite
type DeleteInvite struct {
	// TODO
}

/// .go fileresources\Stage_Instance.md
// Create Stage Instance
// POST /stage-instances
// https://discord.com/developers/docs/resources/stage-instance#create-stage-instance
type CreateStageInstance struct {
	ChannelID             resources.Snowflake
	Topic                 string `json:"topic"`
	PrivacyLevel          int    `json:"privacy_level"`
	SendStartNotification bool   `json:"send_start_notification"`
}

// Get Stage Instance
// GET /stage-instances/{channel.id}
// https://discord.com/developers/docs/resources/stage-instance#get-stage-instance
type GetStageInstance struct {
	ChannelID resources.Snowflake
}

// Modify Stage Instance
// PATCH /stage-instances/{channel.id}
// https://discord.com/developers/docs/resources/stage-instance#modify-stage-instance
type ModifyStageInstance struct {
	ChannelID    resources.Snowflake
	Topic        string `json:"topic"`
	PrivacyLevel int    `json:"privacy_level"`
}

// Delete Stage Instance
// DELETE /stage-instances/{channel.id}
// https://discord.com/developers/docs/resources/stage-instance#delete-stage-instance
type DeleteStageInstance struct {
	ChannelID resources.Snowflake
}

/// .go fileresources\Sticker.md
// Get Sticker
// GET /stickers/{sticker.id}
// https://discord.com/developers/docs/resources/sticker#get-sticker
type GetSticker struct {
	StickerID resources.Snowflake
}

// List Nitro Sticker Packs
// GET /sticker-packs
// https://discord.com/developers/docs/resources/sticker#list-nitro-sticker-packs
type ListNitroStickerPacks struct {
	StickerPacks []*resources.StickerPack `json:"sticker_packs"`
}

// List Guild Stickers
// GET /guilds/{guild.id}/stickers
// https://discord.com/developers/docs/resources/sticker#list-guild-stickers
type ListGuildStickers struct {
	GuildID resources.Snowflake
}

// Get Guild Sticker
// GET /guilds/{guild.id}/stickers/{sticker.id}
// https://discord.com/developers/docs/resources/sticker#get-guild-sticker
type GetGuildSticker struct {
	StickerID resources.Snowflake
}

// Create Guild Sticker
// POST /guilds/{guild.id}/stickers
// https://discord.com/developers/docs/resources/sticker#create-guild-sticker
type CreateGuildSticker struct {
	GuildID     resources.Snowflake
	Name        string `json:"name"`
	Description string `json:"description"`
	Tags        string `json:"tags"`
	////Files           []*File                 `json:"-"`
}

// Modify Guild Sticker
// PATCH /guilds/{guild.id}/stickers/{sticker.id}
// https://discord.com/developers/docs/resources/sticker#modify-guild-sticker
type ModifyGuildSticker struct {
	StickerID   resources.Snowflake
	Name        string `json:"name"`
	Description string `json:"description"`
	Tags        string `json:"tags"`
}

// Delete Guild Sticker
// DELETE /guilds/{guild.id}/stickers/{sticker.id}
// https://discord.com/developers/docs/resources/sticker#delete-guild-sticker
type DeleteGuildSticker struct {
	StickerID resources.Snowflake
}

/// .go fileresources\User.md
// Get Current User
// GET /users/@me
// https://discord.com/developers/docs/game-sdk/users#getcurrentuser
type GetCurrentUser struct {
	UserID resources.Snowflake
}

// Get User
// GET /users/{user.id}
// https://discord.com/developers/docs/game-sdk/users#getuser
type GetUser struct {
	UserID resources.Snowflake
}

// Modify Current User
// PATCH /users/@me
// https://discord.com/developers/docs/resources/user#modify-current-user
type ModifyCurrentUser struct {
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

// Get Current User Guilds
// GET /users/@me/guilds
// https://discord.com/developers/docs/resources/user#get-current-user-guilds
type GetCurrentUserGuilds struct {
	Before resources.Snowflake `urlparam:"before,omitempty"`
	After  resources.Snowflake `urlparam:"after,omitempty"`
	Limit  int                 `urlparam:"limit,omitempty"`
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
	RecipientID resources.Snowflake `json:"recipient_id"`
}

// Create Group DM
// POST /users/@me/channels
// https://discord.com/developers/docs/resources/user#create-group-dm
type CreateGroupDM struct {
	AccessTokens []string                       `json:"access_tokens"`
	Nicks        map[resources.Snowflake]string `json:"nicks"`
}

// Get User Connections
// GET /users/@me/connections
// https://discord.com/developers/docs/resources/user#get-user-connections
type GetUserConnections struct {
	// TODO
}

/// .go fileresources\Voice.md
// List Voice Regions
// GET /voice/regions
// https://discord.com/developers/docs/resources/voice#list-voice-regions
type ListVoiceRegions struct {
	// TODO
}

/// .go fileresources\Webhook.md
// Create Webhook
// POST /channels/{channel.id}/webhooks
// https://discord.com/developers/docs/resources/webhook#create-webhook
type CreateWebhook struct {
	Name   string `json:"name,omitempty"`
	Avatar string `json:"avatar"`
}

// Get Channel Webhooks
// GET /channels/{channel.id}/webhooks
// https://discord.com/developers/docs/resources/webhook#get-channel-webhooks
type GetChannelWebhooks struct {
	ChannelID resources.Snowflake
}

// Get Guild Webhooks
// GET /guilds/{guild.id}/webhooks
// https://discord.com/developers/docs/resources/webhook#get-guild-webhooks
type GetGuildWebhooks struct {
	GuildID resources.Snowflake
}

// Get Webhook
// GET /webhooks/{webhook.id}
// https://discord.com/developers/docs/resources/webhook#get-webhook
type GetWebhook struct {
	WebhookID resources.Snowflake
}

// Get Webhook with Token
// GET /webhooks/{webhook.id}/{webhook.token}
// https://discord.com/developers/docs/resources/webhook#get-webhook-with-token
type GetWebhookwithToken struct {
	WebhookID resources.Snowflake
}

// Modify Webhook
// PATCH /webhooks/{webhook.id}
// https://discord.com/developers/docs/resources/webhook#modify-webhook
type ModifyWebhook struct {
	WebhookID resources.Snowflake
	Name      string              `json:"name,omitempty"`
	Avatar    string              `json:"avatar"`
	ChannelID resources.Snowflake `json:"channel_id"`
}

// Modify Webhook with Token
// PATCH /webhooks/{webhook.id}/{webhook.token}
// https://discord.com/developers/docs/resources/webhook#modify-webhook-with-token
type ModifyWebhookwithToken struct {
	WebhookID resources.Snowflake
}

// Delete Webhook
// DELETE /webhooks/{webhook.id}
// https://discord.com/developers/docs/resources/webhook#delete-webhook
type DeleteWebhook struct {
	WebhookID resources.Snowflake
}

// Delete Webhook with Token
// DELETE /webhooks/{webhook.id}/{webhook.token}
// https://discord.com/developers/docs/resources/webhook#delete-webhook-with-token
type DeleteWebhookwithToken struct {
	WebhookID resources.Snowflake
}

// Execute Webhook
// POST /webhooks/{webhook.id}/{webhook.token}
// https://discord.com/developers/docs/resources/webhook#execute-webhook
type ExecuteWebhook struct {
	Content         string                     `json:"content,omitempty"`
	Username        string                     `json:"username,omitempty"`
	AvatarURL       string                     `json:"avatar_url,omitempty"`
	TTS             bool                       `json:"tts,omitempty"`
	Files           []byte                     `disgo:"TODO"`
	Components      []resources.Component      `json:"components"`
	Embeds          []*resources.Embed         `json:"embeds,omitempty"`
	AllowedMentions *resources.AllowedMentions `json:"allowed_mentions,omitempty"`
	PayloadJSON     string                     `json:"payload_json,omitempty"`
	Attachments     []*resources.Attachment    `json:"attachments,omitempty"`
}

// Execute Slack-Compatible Webhook
// POST /webhooks/{webhook.id}/{webhook.token}/slack
// https://discord.com/developers/docs/resources/webhook#execute-slackcompatible-webhook
type ExecuteSlackCompatibleWebhook struct {
	ThreadID resources.Snowflake
	Wait     bool `json:"wait,omitempty"`
}

// Execute GitHub-Compatible Webhook
// POST /webhooks/{webhook.id}/{webhook.token}/github
// https://discord.com/developers/docs/resources/webhook#execute-githubcompatible-webhook
type ExecuteGitHubCompatibleWebhook struct {
	ThreadID resources.Snowflake
	Wait     bool `json:"wait,omitempty"`
}

// Get Webhook Message
// GET /webhooks/{webhook.id}/{webhook.token}/messages/{message.id}
// https://discord.com/developers/docs/resources/webhook#get-webhook-message
type GetWebhookMessage struct {
	ThreadID resources.Snowflake
}

// Edit Webhook Message
// PATCH /webhooks/{webhook.id}/{webhook.token}/messages/{message.id}
// https://discord.com/developers/docs/resources/webhook#edit-webhook-message
type EditWebhookMessage struct {
	WebhookID       resources.Snowflake
	Content         string                     `json:"content,omitempty"`
	Components      []*resources.Component     `json:"components"`
	Embeds          []*resources.Embed         `json:"embeds,omitempty"`
	Files           []byte                     `json:"-"`
	AllowedMentions *resources.AllowedMentions `json:"allowed_mentions,omitempty"`
	PayloadJSON     string                     `json:"payload_json,omitempty"`
	Attachments     []*resources.Attachment    `json:"attachments,omitempty"`
}

// Delete Webhook Message
// DELETE /webhooks/{webhook.id}/{webhook.token}/messages/{message.id}
// https://discord.com/developers/docs/resources/webhook#delete-webhook-message
type DeleteWebhookMessage struct {
	ThreadID resources.Snowflake
}

/// .go filetopics\Gateway.md
// Get Gateway
// GET /gateway
// https://discord.com/developers/docs/topics/gateway#get-gateway
type GetGateway struct {
	// TODO
}

// Get Gateway Bot
// GET /gateway/bot
// https://discord.com/developers/docs/topics/gateway#get-gateway-bot
type GetGatewayBot struct {
	URL               string `json:"utl"`
	Shards            int    `json:"shards"`
	SessionStartLimit int    `json:"session_start_limit"`
}

/// .go filetopics\OAuth2.md
// Get Current Bot Application Information
// GET /oauth2/applications/@me
// https://discord.com/developers/docs/topics/oauth2#get-current-bot-application-information
type GetCurrentBotApplicationInformation struct {
	// TODO
}

// Get Current Authorization Information
// GET /oauth2/@me
// https://discord.com/developers/docs/topics/oauth2#get-current-authorization-information
type GetCurrentAuthorizationInformation struct {
	// TODO
}
