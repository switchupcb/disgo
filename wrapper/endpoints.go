package disgo

// Discord API Endpoints
const (
	EndpointBaseURL = "https://discord.com/api/v9/"
	callback        = "callback"
	threads         = "threads"
	public          = "public"
	preview         = "preview"
	widget          = "widget"
	messages        = "messages"
	auditlogs       = "audit-logs"
	followers       = "followers"
	private         = "private"
	me              = "@me"
	bulkdelete      = "bulk-delete"
	threadmembers   = "thread-members"
	roles           = "roles"
	stageinstances  = "stage-instances"
	voice           = "voice"
	archived        = "archived"
	members         = "members"
	integrations    = "integrations"
	voicestates     = "voice-states"
	channels        = "channels"
	crosspost       = "crosspost"
	active          = "active"
	welcomescreen   = "welcome-screen"
	stickerpacks    = "sticker-packs"
	guilds          = "guilds"
	invites         = "invites"
	templates       = "templates"
	prune           = "prune"
	bot             = "bot"
	recipients      = "recipients"
	emojis          = "emojis"
	gateway         = "gateway"
	permissions     = "permissions"
	vanityurl       = "vanity-url"
	stickers        = "stickers"
	slack           = "slack"
	typing          = "typing"
	oauth           = "oauth2"
	applications    = "applications"
	commands        = "commands"
	users           = "users"
	scheduledevents = "scheduled-events"
	webhooks        = "webhooks"
	bans            = "bans"
	widgetjson      = "widget.json"
	github          = "github"
	regions         = "regions"
	widgetpng       = "widget.png"
	interactions    = "interactions"
	search          = "search"
	member          = "member"
	slash           = "/"
	original        = "@original"
	reactions       = "reactions"
	pins            = "pins"
	nick            = "nick"
	connections     = "connections"
)

// EndpointGetGlobalApplicationCommandsbuilds a query for an HTTP request.
func EndpointGetGlobalApplicationCommands(applicationid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + commands
}

// EndpointCreateGlobalApplicationCommandbuilds a query for an HTTP request.
func EndpointCreateGlobalApplicationCommand(applicationid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + commands
}

// EndpointGetGlobalApplicationCommandbuilds a query for an HTTP request.
func EndpointGetGlobalApplicationCommand(applicationid, commandid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + commands + slash + commandid
}

// EndpointEditGlobalApplicationCommandbuilds a query for an HTTP request.
func EndpointEditGlobalApplicationCommand(applicationid, commandid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + commands + slash + commandid
}

// EndpointDeleteGlobalApplicationCommandbuilds a query for an HTTP request.
func EndpointDeleteGlobalApplicationCommand(applicationid, commandid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + commands + slash + commandid
}

// EndpointBulkOverwriteGlobalApplicationCommandsbuilds a query for an HTTP request.
func EndpointBulkOverwriteGlobalApplicationCommands(applicationid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + commands
}

// EndpointGetGuildApplicationCommandsbuilds a query for an HTTP request.
func EndpointGetGuildApplicationCommands(applicationid, guildid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + guilds + slash + guildid + slash + commands
}

// EndpointCreateGuildApplicationCommandbuilds a query for an HTTP request.
func EndpointCreateGuildApplicationCommand(applicationid, guildid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + guilds + slash + guildid + slash + commands
}

// EndpointGetGuildApplicationCommandbuilds a query for an HTTP request.
func EndpointGetGuildApplicationCommand(applicationid, guildid, commandid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + guilds + slash + guildid + slash + commands + slash + commandid
}

// EndpointEditGuildApplicationCommandbuilds a query for an HTTP request.
func EndpointEditGuildApplicationCommand(applicationid, guildid, commandid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + guilds + slash + guildid + slash + commands + slash + commandid
}

// EndpointDeleteGuildApplicationCommandbuilds a query for an HTTP request.
func EndpointDeleteGuildApplicationCommand(applicationid, guildid, commandid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + guilds + slash + guildid + slash + commands + slash + commandid
}

// EndpointBulkOverwriteGuildApplicationCommandsbuilds a query for an HTTP request.
func EndpointBulkOverwriteGuildApplicationCommands(applicationid, guildid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + guilds + slash + guildid + slash + commands
}

// EndpointGetGuildApplicationCommandPermissionsbuilds a query for an HTTP request.
func EndpointGetGuildApplicationCommandPermissions(applicationid, guildid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + guilds + slash + guildid + slash + commands + slash + permissions
}

// EndpointGetApplicationCommandPermissionsbuilds a query for an HTTP request.
func EndpointGetApplicationCommandPermissions(applicationid, guildid, commandid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + guilds + slash + guildid + slash + commands + slash + commandid + slash + permissions
}

// EndpointEditApplicationCommandPermissionsbuilds a query for an HTTP request.
func EndpointEditApplicationCommandPermissions(applicationid, guildid, commandid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + guilds + slash + guildid + slash + commands + slash + commandid + slash + permissions
}

// EndpointBatchEditApplicationCommandPermissionsbuilds a query for an HTTP request.
func EndpointBatchEditApplicationCommandPermissions(applicationid, guildid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + guilds + slash + guildid + slash + commands + slash + permissions
}

// EndpointCreateInteractionResponsebuilds a query for an HTTP request.
func EndpointCreateInteractionResponse(interactionid, interactiontoken string) string {
	return EndpointBaseURL + interactions + slash + interactionid + slash + interactiontoken + slash + callback
}

// EndpointGetOriginalInteractionResponsebuilds a query for an HTTP request.
func EndpointGetOriginalInteractionResponse(applicationid, interactiontoken string) string {
	return EndpointBaseURL + webhooks + slash + applicationid + slash + interactiontoken + slash + messages + slash + original
}

// EndpointEditOriginalInteractionResponsebuilds a query for an HTTP request.
func EndpointEditOriginalInteractionResponse(applicationid, interactiontoken string) string {
	return EndpointBaseURL + webhooks + slash + applicationid + slash + interactiontoken + slash + messages + slash + original
}

// EndpointDeleteOriginalInteractionResponsebuilds a query for an HTTP request.
func EndpointDeleteOriginalInteractionResponse(applicationid, interactiontoken string) string {
	return EndpointBaseURL + webhooks + slash + applicationid + slash + interactiontoken + slash + messages + slash + original
}

// EndpointCreateFollowupMessagebuilds a query for an HTTP request.
func EndpointCreateFollowupMessage(applicationid, interactiontoken string) string {
	return EndpointBaseURL + webhooks + slash + applicationid + slash + interactiontoken
}

// EndpointGetFollowupMessagebuilds a query for an HTTP request.
func EndpointGetFollowupMessage(applicationid, interactiontoken, messageid string) string {
	return EndpointBaseURL + webhooks + slash + applicationid + slash + interactiontoken + slash + messages + slash + messageid
}

// EndpointEditFollowupMessagebuilds a query for an HTTP request.
func EndpointEditFollowupMessage(applicationid, interactiontoken, messageid string) string {
	return EndpointBaseURL + webhooks + slash + applicationid + slash + interactiontoken + slash + messages + slash + messageid
}

// EndpointDeleteFollowupMessagebuilds a query for an HTTP request.
func EndpointDeleteFollowupMessage(applicationid, interactiontoken, messageid string) string {
	return EndpointBaseURL + webhooks + slash + applicationid + slash + interactiontoken + slash + messages + slash + messageid
}

// EndpointGetGuildAuditLogbuilds a query for an HTTP request.
func EndpointGetGuildAuditLog(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + auditlogs
}

// EndpointGetChannelbuilds a query for an HTTP request.
func EndpointGetChannel(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid
}

// EndpointModifyChannelbuilds a query for an HTTP request.
func EndpointModifyChannel(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid
}

// EndpointDeleteCloseChannelbuilds a query for an HTTP request.
func EndpointDeleteCloseChannel(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid
}

// EndpointGetChannelMessagesbuilds a query for an HTTP request.
func EndpointGetChannelMessages(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + messages
}

// EndpointGetChannelMessagebuilds a query for an HTTP request.
func EndpointGetChannelMessage(channelid, messageid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + messages + slash + messageid
}

// EndpointCreateMessagebuilds a query for an HTTP request.
func EndpointCreateMessage(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + messages
}

// EndpointCrosspostMessagebuilds a query for an HTTP request.
func EndpointCrosspostMessage(channelid, messageid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + messages + slash + messageid + slash + crosspost
}

// EndpointCreateReactionbuilds a query for an HTTP request.
func EndpointCreateReaction(channelid, messageid, emoji string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + messages + slash + messageid + slash + reactions + slash + emoji + slash + me
}

// EndpointDeleteOwnReactionbuilds a query for an HTTP request.
func EndpointDeleteOwnReaction(channelid, messageid, emoji string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + messages + slash + messageid + slash + reactions + slash + emoji + slash + me
}

// EndpointDeleteUserReactionbuilds a query for an HTTP request.
func EndpointDeleteUserReaction(channelid, messageid, emoji, userid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + messages + slash + messageid + slash + reactions + slash + emoji + slash + userid
}

// EndpointGetReactionsbuilds a query for an HTTP request.
func EndpointGetReactions(channelid, messageid, emoji string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + messages + slash + messageid + slash + reactions + slash + emoji
}

// EndpointDeleteAllReactionsbuilds a query for an HTTP request.
func EndpointDeleteAllReactions(channelid, messageid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + messages + slash + messageid + slash + reactions
}

// EndpointDeleteAllReactionsforEmojibuilds a query for an HTTP request.
func EndpointDeleteAllReactionsforEmoji(channelid, messageid, emoji string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + messages + slash + messageid + slash + reactions + slash + emoji
}

// EndpointEditMessagebuilds a query for an HTTP request.
func EndpointEditMessage(channelid, messageid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + messages + slash + messageid
}

// EndpointDeleteMessagebuilds a query for an HTTP request.
func EndpointDeleteMessage(channelid, messageid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + messages + slash + messageid
}

// EndpointBulkDeleteMessagesbuilds a query for an HTTP request.
func EndpointBulkDeleteMessages(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + messages + slash + bulkdelete
}

// EndpointEditChannelPermissionsbuilds a query for an HTTP request.
func EndpointEditChannelPermissions(channelid, overwriteid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + permissions + slash + overwriteid
}

// EndpointGetChannelInvitesbuilds a query for an HTTP request.
func EndpointGetChannelInvites(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + invites
}

// EndpointCreateChannelInvitebuilds a query for an HTTP request.
func EndpointCreateChannelInvite(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + invites
}

// EndpointDeleteChannelPermissionbuilds a query for an HTTP request.
func EndpointDeleteChannelPermission(channelid, overwriteid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + permissions + slash + overwriteid
}

// EndpointFollowNewsChannelbuilds a query for an HTTP request.
func EndpointFollowNewsChannel(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + followers
}

// EndpointTriggerTypingIndicatorbuilds a query for an HTTP request.
func EndpointTriggerTypingIndicator(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + typing
}

// EndpointGetPinnedMessagesbuilds a query for an HTTP request.
func EndpointGetPinnedMessages(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + pins
}

// EndpointPinMessagebuilds a query for an HTTP request.
func EndpointPinMessage(channelid, messageid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + pins + slash + messageid
}

// EndpointUnpinMessagebuilds a query for an HTTP request.
func EndpointUnpinMessage(channelid, messageid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + pins + slash + messageid
}

// EndpointGroupDMAddRecipientbuilds a query for an HTTP request.
func EndpointGroupDMAddRecipient(channelid, userid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + recipients + slash + userid
}

// EndpointGroupDMRemoveRecipientbuilds a query for an HTTP request.
func EndpointGroupDMRemoveRecipient(channelid, userid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + recipients + slash + userid
}

// EndpointStartThreadfromMessagebuilds a query for an HTTP request.
func EndpointStartThreadfromMessage(channelid, messageid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + messages + slash + messageid + slash + threads
}

// EndpointStartThreadwithoutMessagebuilds a query for an HTTP request.
func EndpointStartThreadwithoutMessage(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + threads
}

// EndpointStartThreadinForumChannelbuilds a query for an HTTP request.
func EndpointStartThreadinForumChannel(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + threads
}

// EndpointJoinThreadbuilds a query for an HTTP request.
func EndpointJoinThread(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + threadmembers + slash + me
}

// EndpointAddThreadMemberbuilds a query for an HTTP request.
func EndpointAddThreadMember(channelid, userid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + threadmembers + slash + userid
}

// EndpointLeaveThreadbuilds a query for an HTTP request.
func EndpointLeaveThread(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + threadmembers + slash + me
}

// EndpointRemoveThreadMemberbuilds a query for an HTTP request.
func EndpointRemoveThreadMember(channelid, userid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + threadmembers + slash + userid
}

// EndpointGetThreadMemberbuilds a query for an HTTP request.
func EndpointGetThreadMember(channelid, userid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + threadmembers + slash + userid
}

// EndpointListThreadMembersbuilds a query for an HTTP request.
func EndpointListThreadMembers(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + threadmembers
}

// EndpointListActiveThreadsbuilds a query for an HTTP request.
func EndpointListActiveThreads(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + threads + slash + active
}

// EndpointListPublicArchivedThreadsbuilds a query for an HTTP request.
func EndpointListPublicArchivedThreads(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + threads + slash + archived + slash + public
}

// EndpointListPrivateArchivedThreadsbuilds a query for an HTTP request.
func EndpointListPrivateArchivedThreads(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + threads + slash + archived + slash + private
}

// EndpointListJoinedPrivateArchivedThreadsbuilds a query for an HTTP request.
func EndpointListJoinedPrivateArchivedThreads(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + users + slash + me + slash + threads + slash + archived + slash + private
}

// EndpointListGuildEmojisbuilds a query for an HTTP request.
func EndpointListGuildEmojis(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + emojis
}

// EndpointGetGuildEmojibuilds a query for an HTTP request.
func EndpointGetGuildEmoji(guildid, emojiid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + emojis + slash + emojiid
}

// EndpointCreateGuildEmojibuilds a query for an HTTP request.
func EndpointCreateGuildEmoji(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + emojis
}

// EndpointModifyGuildEmojibuilds a query for an HTTP request.
func EndpointModifyGuildEmoji(guildid, emojiid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + emojis + slash + emojiid
}

// EndpointDeleteGuildEmojibuilds a query for an HTTP request.
func EndpointDeleteGuildEmoji(guildid, emojiid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + emojis + slash + emojiid
}

// EndpointListScheduledEventsforGuildbuilds a query for an HTTP request.
func EndpointListScheduledEventsforGuild(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + scheduledevents
}

// EndpointCreateGuildScheduledEventbuilds a query for an HTTP request.
func EndpointCreateGuildScheduledEvent(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + scheduledevents
}

// EndpointGetGuildScheduledEventbuilds a query for an HTTP request.
func EndpointGetGuildScheduledEvent(guildid, guildscheduledeventid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + scheduledevents + slash + guildscheduledeventid
}

// EndpointModifyGuildScheduledEventbuilds a query for an HTTP request.
func EndpointModifyGuildScheduledEvent(guildid, guildscheduledeventid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + scheduledevents + slash + guildscheduledeventid
}

// EndpointDeleteGuildScheduledEventbuilds a query for an HTTP request.
func EndpointDeleteGuildScheduledEvent(guildid, guildscheduledeventid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + scheduledevents + slash + guildscheduledeventid
}

// EndpointGetGuildScheduledEventUsersbuilds a query for an HTTP request.
func EndpointGetGuildScheduledEventUsers(guildid, guildscheduledeventid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + scheduledevents + slash + guildscheduledeventid + slash + users
}

// EndpointGetGuildTemplatebuilds a query for an HTTP request.
func EndpointGetGuildTemplate(templatecode string) string {
	return EndpointBaseURL + guilds + slash + templates + slash + templatecode
}

// EndpointCreateGuildfromGuildTemplatebuilds a query for an HTTP request.
func EndpointCreateGuildfromGuildTemplate(templatecode string) string {
	return EndpointBaseURL + guilds + slash + templates + slash + templatecode
}

// EndpointGetGuildTemplatesbuilds a query for an HTTP request.
func EndpointGetGuildTemplates(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + templates
}

// EndpointCreateGuildTemplatebuilds a query for an HTTP request.
func EndpointCreateGuildTemplate(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + templates
}

// EndpointSyncGuildTemplatebuilds a query for an HTTP request.
func EndpointSyncGuildTemplate(guildid, templatecode string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + templates + slash + templatecode
}

// EndpointModifyGuildTemplatebuilds a query for an HTTP request.
func EndpointModifyGuildTemplate(guildid, templatecode string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + templates + slash + templatecode
}

// EndpointDeleteGuildTemplatebuilds a query for an HTTP request.
func EndpointDeleteGuildTemplate(guildid, templatecode string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + templates + slash + templatecode
}

// EndpointCreateGuildbuilds a query for an HTTP request.
func EndpointCreateGuild() string {
	return EndpointBaseURL + guilds
}

// EndpointGetGuildbuilds a query for an HTTP request.
func EndpointGetGuild(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid
}

// EndpointGetGuildPreviewbuilds a query for an HTTP request.
func EndpointGetGuildPreview(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + preview
}

// EndpointModifyGuildbuilds a query for an HTTP request.
func EndpointModifyGuild(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid
}

// EndpointDeleteGuildbuilds a query for an HTTP request.
func EndpointDeleteGuild(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid
}

// EndpointGetGuildChannelsbuilds a query for an HTTP request.
func EndpointGetGuildChannels(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + channels
}

// EndpointCreateGuildChannelbuilds a query for an HTTP request.
func EndpointCreateGuildChannel(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + channels
}

// EndpointModifyGuildChannelPositionsbuilds a query for an HTTP request.
func EndpointModifyGuildChannelPositions(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + channels
}

// EndpointListActiveGuildThreadsbuilds a query for an HTTP request.
func EndpointListActiveGuildThreads(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + threads + slash + active
}

// EndpointGetGuildMemberbuilds a query for an HTTP request.
func EndpointGetGuildMember(guildid, userid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + members + slash + userid
}

// EndpointListGuildMembersbuilds a query for an HTTP request.
func EndpointListGuildMembers(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + members
}

// EndpointSearchGuildMembersbuilds a query for an HTTP request.
func EndpointSearchGuildMembers(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + members + slash + search
}

// EndpointAddGuildMemberbuilds a query for an HTTP request.
func EndpointAddGuildMember(guildid, userid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + members + slash + userid
}

// EndpointModifyGuildMemberbuilds a query for an HTTP request.
func EndpointModifyGuildMember(guildid, userid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + members + slash + userid
}

// EndpointModifyCurrentMemberbuilds a query for an HTTP request.
func EndpointModifyCurrentMember(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + members + slash + me
}

// EndpointModifyCurrentUserNickbuilds a query for an HTTP request.
func EndpointModifyCurrentUserNick(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + members + slash + me + slash + nick
}

// EndpointAddGuildMemberRolebuilds a query for an HTTP request.
func EndpointAddGuildMemberRole(guildid, userid, roleid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + members + slash + userid + slash + roles + slash + roleid
}

// EndpointRemoveGuildMemberRolebuilds a query for an HTTP request.
func EndpointRemoveGuildMemberRole(guildid, userid, roleid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + members + slash + userid + slash + roles + slash + roleid
}

// EndpointRemoveGuildMemberbuilds a query for an HTTP request.
func EndpointRemoveGuildMember(guildid, userid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + members + slash + userid
}

// EndpointGetGuildBansbuilds a query for an HTTP request.
func EndpointGetGuildBans(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + bans
}

// EndpointGetGuildBanbuilds a query for an HTTP request.
func EndpointGetGuildBan(guildid, userid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + bans + slash + userid
}

// EndpointCreateGuildBanbuilds a query for an HTTP request.
func EndpointCreateGuildBan(guildid, userid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + bans + slash + userid
}

// EndpointRemoveGuildBanbuilds a query for an HTTP request.
func EndpointRemoveGuildBan(guildid, userid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + bans + slash + userid
}

// EndpointGetGuildRolesbuilds a query for an HTTP request.
func EndpointGetGuildRoles(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + roles
}

// EndpointCreateGuildRolebuilds a query for an HTTP request.
func EndpointCreateGuildRole(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + roles
}

// EndpointModifyGuildRolePositionsbuilds a query for an HTTP request.
func EndpointModifyGuildRolePositions(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + roles
}

// EndpointModifyGuildRolebuilds a query for an HTTP request.
func EndpointModifyGuildRole(guildid, roleid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + roles + slash + roleid
}

// EndpointDeleteGuildRolebuilds a query for an HTTP request.
func EndpointDeleteGuildRole(guildid, roleid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + roles + slash + roleid
}

// EndpointGetGuildPruneCountbuilds a query for an HTTP request.
func EndpointGetGuildPruneCount(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + prune
}

// EndpointBeginGuildPrunebuilds a query for an HTTP request.
func EndpointBeginGuildPrune(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + prune
}

// EndpointGetGuildVoiceRegionsbuilds a query for an HTTP request.
func EndpointGetGuildVoiceRegions(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + regions
}

// EndpointGetGuildInvitesbuilds a query for an HTTP request.
func EndpointGetGuildInvites(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + invites
}

// EndpointGetGuildIntegrationsbuilds a query for an HTTP request.
func EndpointGetGuildIntegrations(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + integrations
}

// EndpointDeleteGuildIntegrationbuilds a query for an HTTP request.
func EndpointDeleteGuildIntegration(guildid, integrationid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + integrations + slash + integrationid
}

// EndpointGetGuildWidgetSettingsbuilds a query for an HTTP request.
func EndpointGetGuildWidgetSettings(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + widget
}

// EndpointModifyGuildWidgetbuilds a query for an HTTP request.
func EndpointModifyGuildWidget(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + widget
}

// EndpointGetGuildWidgetbuilds a query for an HTTP request.
func EndpointGetGuildWidget(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + widgetjson
}

// EndpointGetGuildVanityURLbuilds a query for an HTTP request.
func EndpointGetGuildVanityURL(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + vanityurl
}

// EndpointGetGuildWidgetImagebuilds a query for an HTTP request.
func EndpointGetGuildWidgetImage(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + widgetpng
}

// EndpointGetGuildWelcomeScreenbuilds a query for an HTTP request.
func EndpointGetGuildWelcomeScreen(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + welcomescreen
}

// EndpointModifyGuildWelcomeScreenbuilds a query for an HTTP request.
func EndpointModifyGuildWelcomeScreen(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + welcomescreen
}

// EndpointModifyCurrentUserVoiceStatebuilds a query for an HTTP request.
func EndpointModifyCurrentUserVoiceState(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + voicestates + slash + me
}

// EndpointModifyUserVoiceStatebuilds a query for an HTTP request.
func EndpointModifyUserVoiceState(guildid, userid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + voicestates + slash + userid
}

// EndpointGetInvitebuilds a query for an HTTP request.
func EndpointGetInvite(invitecode string) string {
	return EndpointBaseURL + invites + slash + invitecode
}

// EndpointDeleteInvitebuilds a query for an HTTP request.
func EndpointDeleteInvite(invitecode string) string {
	return EndpointBaseURL + invites + slash + invitecode
}

// EndpointCreateStageInstancebuilds a query for an HTTP request.
func EndpointCreateStageInstance() string {
	return EndpointBaseURL + stageinstances
}

// EndpointGetStageInstancebuilds a query for an HTTP request.
func EndpointGetStageInstance(channelid string) string {
	return EndpointBaseURL + stageinstances + slash + channelid
}

// EndpointModifyStageInstancebuilds a query for an HTTP request.
func EndpointModifyStageInstance(channelid string) string {
	return EndpointBaseURL + stageinstances + slash + channelid
}

// EndpointDeleteStageInstancebuilds a query for an HTTP request.
func EndpointDeleteStageInstance(channelid string) string {
	return EndpointBaseURL + stageinstances + slash + channelid
}

// EndpointGetStickerbuilds a query for an HTTP request.
func EndpointGetSticker(stickerid string) string {
	return EndpointBaseURL + stickers + slash + stickerid
}

// EndpointListNitroStickerPacksbuilds a query for an HTTP request.
func EndpointListNitroStickerPacks() string {
	return EndpointBaseURL + stickerpacks
}

// EndpointListGuildStickersbuilds a query for an HTTP request.
func EndpointListGuildStickers(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + stickers
}

// EndpointGetGuildStickerbuilds a query for an HTTP request.
func EndpointGetGuildSticker(guildid, stickerid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + stickers + slash + stickerid
}

// EndpointCreateGuildStickerbuilds a query for an HTTP request.
func EndpointCreateGuildSticker(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + stickers
}

// EndpointModifyGuildStickerbuilds a query for an HTTP request.
func EndpointModifyGuildSticker(guildid, stickerid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + stickers + slash + stickerid
}

// EndpointDeleteGuildStickerbuilds a query for an HTTP request.
func EndpointDeleteGuildSticker(guildid, stickerid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + stickers + slash + stickerid
}

// EndpointGetCurrentUserbuilds a query for an HTTP request.
func EndpointGetCurrentUser() string {
	return EndpointBaseURL + users + slash + me
}

// EndpointGetUserbuilds a query for an HTTP request.
func EndpointGetUser(userid string) string {
	return EndpointBaseURL + users + slash + userid
}

// EndpointModifyCurrentUserbuilds a query for an HTTP request.
func EndpointModifyCurrentUser() string {
	return EndpointBaseURL + users + slash + me
}

// EndpointGetCurrentUserGuildsbuilds a query for an HTTP request.
func EndpointGetCurrentUserGuilds() string {
	return EndpointBaseURL + users + slash + me + slash + guilds
}

// EndpointGetCurrentUserGuildMemberbuilds a query for an HTTP request.
func EndpointGetCurrentUserGuildMember(guildid string) string {
	return EndpointBaseURL + users + slash + me + slash + guilds + slash + guildid + slash + member
}

// EndpointLeaveGuildbuilds a query for an HTTP request.
func EndpointLeaveGuild(guildid string) string {
	return EndpointBaseURL + users + slash + me + slash + guilds + slash + guildid
}

// EndpointCreateDMbuilds a query for an HTTP request.
func EndpointCreateDM() string {
	return EndpointBaseURL + users + slash + me + slash + channels
}

// EndpointCreateGroupDMbuilds a query for an HTTP request.
func EndpointCreateGroupDM() string {
	return EndpointBaseURL + users + slash + me + slash + channels
}

// EndpointGetUserConnectionsbuilds a query for an HTTP request.
func EndpointGetUserConnections() string {
	return EndpointBaseURL + users + slash + me + slash + connections
}

// EndpointListVoiceRegionsbuilds a query for an HTTP request.
func EndpointListVoiceRegions() string {
	return EndpointBaseURL + voice + slash + regions
}

// EndpointCreateWebhookbuilds a query for an HTTP request.
func EndpointCreateWebhook(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + webhooks
}

// EndpointGetChannelWebhooksbuilds a query for an HTTP request.
func EndpointGetChannelWebhooks(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + webhooks
}

// EndpointGetGuildWebhooksbuilds a query for an HTTP request.
func EndpointGetGuildWebhooks(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + webhooks
}

// EndpointGetWebhookbuilds a query for an HTTP request.
func EndpointGetWebhook(webhookid string) string {
	return EndpointBaseURL + webhooks + slash + webhookid
}

// EndpointGetWebhookwithTokenbuilds a query for an HTTP request.
func EndpointGetWebhookwithToken(webhookid, webhooktoken string) string {
	return EndpointBaseURL + webhooks + slash + webhookid + slash + webhooktoken
}

// EndpointModifyWebhookbuilds a query for an HTTP request.
func EndpointModifyWebhook(webhookid string) string {
	return EndpointBaseURL + webhooks + slash + webhookid
}

// EndpointModifyWebhookwithTokenbuilds a query for an HTTP request.
func EndpointModifyWebhookwithToken(webhookid, webhooktoken string) string {
	return EndpointBaseURL + webhooks + slash + webhookid + slash + webhooktoken
}

// EndpointDeleteWebhookbuilds a query for an HTTP request.
func EndpointDeleteWebhook(webhookid string) string {
	return EndpointBaseURL + webhooks + slash + webhookid
}

// EndpointDeleteWebhookwithTokenbuilds a query for an HTTP request.
func EndpointDeleteWebhookwithToken(webhookid, webhooktoken string) string {
	return EndpointBaseURL + webhooks + slash + webhookid + slash + webhooktoken
}

// EndpointExecuteWebhookbuilds a query for an HTTP request.
func EndpointExecuteWebhook(webhookid, webhooktoken string) string {
	return EndpointBaseURL + webhooks + slash + webhookid + slash + webhooktoken
}

// EndpointExecuteSlackCompatibleWebhookbuilds a query for an HTTP request.
func EndpointExecuteSlackCompatibleWebhook(webhookid, webhooktoken string) string {
	return EndpointBaseURL + webhooks + slash + webhookid + slash + webhooktoken + slash + slack
}

// EndpointExecuteGitHubCompatibleWebhookbuilds a query for an HTTP request.
func EndpointExecuteGitHubCompatibleWebhook(webhookid, webhooktoken string) string {
	return EndpointBaseURL + webhooks + slash + webhookid + slash + webhooktoken + slash + github
}

// EndpointGetWebhookMessagebuilds a query for an HTTP request.
func EndpointGetWebhookMessage(webhookid, webhooktoken, messageid string) string {
	return EndpointBaseURL + webhooks + slash + webhookid + slash + webhooktoken + slash + messages + slash + messageid
}

// EndpointEditWebhookMessagebuilds a query for an HTTP request.
func EndpointEditWebhookMessage(webhookid, webhooktoken, messageid string) string {
	return EndpointBaseURL + webhooks + slash + webhookid + slash + webhooktoken + slash + messages + slash + messageid
}

// EndpointDeleteWebhookMessagebuilds a query for an HTTP request.
func EndpointDeleteWebhookMessage(webhookid, webhooktoken, messageid string) string {
	return EndpointBaseURL + webhooks + slash + webhookid + slash + webhooktoken + slash + messages + slash + messageid
}

// EndpointGetGatewaybuilds a query for an HTTP request.
func EndpointGetGateway() string {
	return EndpointBaseURL + gateway
}

// EndpointGetGatewayBotbuilds a query for an HTTP request.
func EndpointGetGatewayBot() string {
	return EndpointBaseURL + gateway + slash + bot
}

// EndpointGetCurrentBotApplicationInformationbuilds a query for an HTTP request.
func EndpointGetCurrentBotApplicationInformation() string {
	return EndpointBaseURL + oauth + slash + applications + slash + me
}

// EndpointGetCurrentAuthorizationInformationbuilds a query for an HTTP request.
func EndpointGetCurrentAuthorizationInformation() string {
	return EndpointBaseURL + oauth + slash + me
}
