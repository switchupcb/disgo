package requests

import "github.com/switchupcb/disgo/wrapper/resources"

/// .go fileresources\Channel.md
// Get Channel
// GET /channels/{channel.id}
// https://discord.com/developers/docs/resources/channel#get-channel
type GetChannel struct {
	ChannelID resources.Snowflake
}

// Modify Channel
// PATCH /channels/{channel.id}
// https://discord.com/developers/docs/resources/channel#modify-channel
type ModifyChannel struct {
	ChannelID resources.Snowflake
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
	Name                       *string                          `json:"name,omitempty"`
	Type                       *resources.Flag                  `json:"type,omitempty"`
	Position                   *uint                            `json:"position,omitempty"`
	Topic                      *string                          `json:"topic,omitempty"`
	NSFW                       bool                             `json:"nsfw,omitempty"`
	RateLimitPerUser           *resources.CodeFlag              `json:"rate_limit_per_user,omitempty"`
	Bitrate                    *uint                            `json:"bitrate,omitempty"`
	UserLimit                  *resources.Flag                  `json:"user_limit,omitempty"`
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
	Name                string              `json:"name,omitempty"`
	Archived            bool                `json:"archived,omitempty"`
	AutoArchiveDuration resources.CodeFlag  `json:"auto_archive_duration,omitempty"`
	Locked              bool                `json:"locked,omitempty"`
	Invitable           bool                `json:"invitable,omitempty"`
	RateLimitPerUser    *resources.CodeFlag `json:"rate_limit_per_user,omitempty"`
}

// Delete/Close Channel
// DELETE /channels/{channel.id}
// https://discord.com/developers/docs/resources/channel#deleteclose-channel
type DeleteCloseChannel struct {
	ChannelID resources.Snowflake
}

// Get Channel Messages
// GET /channels/{channel.id}/messages
// https://discord.com/developers/docs/resources/channel#get-channel-messages
type GetChannelMessages struct {
	Around *resources.Snowflake `json:"around,omitempty"`
	Before *resources.Snowflake `json:"before,omitempty"`
	After  *resources.Snowflake `json:"after,omitempty"`
	Limit  resources.Flag       `json:"limit,omitempty"`
}

// Get Channel Message
// GET /channels/{channel.id}/messages/{message.id}
// https://discord.com/developers/docs/resources/channel#get-channel-message
type GetChannelMessage struct {
	MessageID resources.Snowflake
}

// Create Message
// POST /channels/{channel.id}/messages
// https://discord.com/developers/docs/resources/channel#create-message
type CreateMessage struct {
	Content         string                      `json:"content,omitempty"`
	TTS             bool                        `json:"tts,omitempty"`
	Embeds          []*resources.Embed          `json:"embeds,omitempty"`
	Embed           *resources.Embed            `json:"embed,omitempty"`
	AllowedMentions *resources.AllowedMentions  `json:"allowed_mentions,omitempty"`
	Reference       *resources.MessageReference `json:"message_reference,omitempty"`
	StickerID       []*resources.Snowflake      `json:"sticker_ids,omitempty"`
	Components      []*resources.Component      `json:"components,omitempty"`
	Files           []byte                      `disgo:"TODO"`
	PayloadJSON     *string                     `json:"payload_json,omitempty"`
	Attachments     []*resources.Attachment     `json:"attachments,omitempty"`
	Flags           resources.BitFlag           `json:"flags,omitempty"`
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
	After resources.Snowflake `json:"after,omitempty"`
	Limit resources.Flag      `json:"limit,omitempty"` // 1 is default. even if 0 is supplied.
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
	Content         *string                    `json:"content,omitempty"`
	Embeds          []*resources.Embed         `json:"embeds,omitempty"`
	Embed           *resources.Embed           `json:"embed,omitempty"`
	Flags           *resources.BitFlag         `json:"flags,omitempty"`
	AllowedMentions *resources.AllowedMentions `json:"allowed_mentions,omitempty"`
	Components      []*resources.Component     `json:"components,omitempty"`
	Files           []byte                     `disgo:"TODO"`
	PayloadJSON     *string                    `json:"payload_json,omitempty"`
	Attachments     []*resources.Attachment    `json:"attachments,omitempty"`
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
	Messages []*resources.Snowflake `json:"messages,omitempty"`
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
	MaxAge              *int                `json:"max_age,omitempty"`
	MaxUses             *resources.Flag     `json:"max_uses,omitempty"`
	Temporary           bool                `json:"temporary,omitempty"`
	Unique              bool                `json:"unique,omitempty"`
	TargetType          resources.Flag      `json:"target_type,omitempty"`
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
	AccessToken string  `json:"access_token,omitempty"`
	Nickname    *string `json:"nick,omitempty"`
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
	Name                string `json:"name,omitempty"`
	RateLimitPerUser    uint   `json:"rate_limit_per_user,omitempty"`
	AutoArchiveDuration *int   `json:"auto_archive_duration,omitempty"`
}

// Start Thread without Message
// POST /channels/{channel.id}/threads
// https://discord.com/developers/docs/resources/channel#start-thread-without-message
type StartThreadwithoutMessage struct {
	Name                string              `json:"name,omitempty"`
	AutoArchiveDuration resources.CodeFlag  `json:"auto_archive_duration,omitempty"`
	Type                *resources.Flag     `json:"type,omitempty"`
	Invitable           bool                `json:"invitable,omitempty"`
	RateLimitPerUser    *resources.CodeFlag `json:"rate_limit_per_user,omitempty"`
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
	Name                string              `json:"name,omitempty"`
	RateLimitPerUser    uint                `json:"rate_limit_per_user,omitempty"`
	AutoArchiveDuration *resources.CodeFlag `json:"auto_archive_duration,omitempty"`
}

// Start Thread in Forum Channel
// POST /channels/{channel.id}/threads
// https://discord.com/developers/docs/resources/channel#start-thread-in-forum-channel-json-params-for-the-message
type StartThreadinForumChannelMessage struct {
	Content         *string                    `json:"content,omitempty"`
	Embeds          []*resources.Embed         `json:"embeds,omitempty"`
	AllowedMentions *resources.AllowedMentions `json:"allowed_mentions,omitempty"`
	Components      []*resources.Component     `json:"components,omitempty"`
	StickerIDS      []*resources.Snowflake     `json:"sticker_ids,omitempty"`
	Attachments     []*resources.Attachment    `json:"attachments,omitempty"`
	Files           []byte                     `disgo:"TODO"`
	PayloadJSON     string                     `json:"payload_json,omitempty"`
	Flags           resources.BitFlag          `json:"flags,omitempty"`
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
	Before resources.Snowflake `json:"before,omitempty"`
	Limit  int                 `json:"limit,omitempty"`
}

// List Public Archived Threads
// GET /channels/{channel.id}/threads/archived/public
// https://discord.com/developers/docs/resources/channel#list-public-archived-threads
type ListPublicArchivedThreads struct {
	Before resources.Snowflake `json:"before,omitempty"`
	Limit  int                 `json:"limit,omitempty"`
}

// List Private Archived Threads
// GET /channels/{channel.id}/threads/archived/private
// https://discord.com/developers/docs/resources/channel#list-private-archived-threads
type ListPrivateArchivedThreads struct {
	Before resources.Snowflake `json:"before,omitempty"`
	Limit  int                 `json:"limit,omitempty"`
}

// List Joined Private Archived Threads
// GET /channels/{channel.id}/users/@me/threads/archived/private
// https://discord.com/developers/docs/resources/channel#list-joined-private-archived-threads
type ListJoinedPrivateArchivedThreads struct {
	Before resources.Snowflake `json:"before,omitempty"`
	Limit  int                 `json:"limit,omitempty"`
}
