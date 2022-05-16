package disgo

// Discord API Endpoints
const (
	EndpointBaseURL = "https://discord.com/api/v9/"
	guilds          = "guilds"
	archived        = "archived"
	stageinstances  = "stage-instances"
	voice           = "voice"
	voicestates     = "voice-states"
	applications    = "applications"
	bulkdelete      = "bulk-delete"
	pins            = "pins"
	preview         = "preview"
	members         = "members"
	interactions    = "interactions"
	messages        = "messages"
	templates       = "templates"
	threads         = "threads"
	users           = "users"
	welcomescreen   = "welcome-screen"
	stickerpacks    = "sticker-packs"
	recipients      = "recipients"
	widget          = "widget"
	bot             = "bot"
	commands        = "commands"
	callback        = "callback"
	public          = "public"
	private         = "private"
	connections     = "connections"
	permissions     = "permissions"
	original        = "@original"
	scheduledevents = "scheduled-events"
	vanityurl       = "vanity-url"
	stickers        = "stickers"
	auditlogs       = "audit-logs"
	reactions       = "reactions"
	prune           = "prune"
	regions         = "regions"
	crosspost       = "crosspost"
	invites         = "invites"
	nick            = "nick"
	widgetpng       = "widget.png"
	member          = "member"
	channels        = "channels"
	bans            = "bans"
	slack           = "slack"
	gateway         = "gateway"
	webhooks        = "webhooks"
	typing          = "typing"
	emojis          = "emojis"
	integrations    = "integrations"
	threadmembers   = "thread-members"
	active          = "active"
	roles           = "roles"
	me              = "@me"
	widgetjson      = "widget.json"
	github          = "github"
	oauth           = "oauth2"
	followers       = "followers"
	search          = "search"
)

// EndpointGetGlobalApplicationCommandsbuilds a query for an HTTP request.
func EndpointGetGlobalApplicationCommands(applications, applicationid, commands string) string {
	return applications + applicationid + commands
}

// EndpointCreateGlobalApplicationCommandbuilds a query for an HTTP request.
func EndpointCreateGlobalApplicationCommand(applications, applicationid, commands string) string {
	return applications + applicationid + commands
}

// EndpointGetGlobalApplicationCommandbuilds a query for an HTTP request.
func EndpointGetGlobalApplicationCommand(applications, applicationid, commands, commandid string) string {
	return applications + applicationid + commands + commandid
}

// EndpointEditGlobalApplicationCommandbuilds a query for an HTTP request.
func EndpointEditGlobalApplicationCommand(applications, applicationid, commands, commandid string) string {
	return applications + applicationid + commands + commandid
}

// EndpointDeleteGlobalApplicationCommandbuilds a query for an HTTP request.
func EndpointDeleteGlobalApplicationCommand(applications, applicationid, commands, commandid string) string {
	return applications + applicationid + commands + commandid
}

// EndpointBulkOverwriteGlobalApplicationCommandsbuilds a query for an HTTP request.
func EndpointBulkOverwriteGlobalApplicationCommands(applications, applicationid, commands string) string {
	return applications + applicationid + commands
}

// EndpointGetGuildApplicationCommandsbuilds a query for an HTTP request.
func EndpointGetGuildApplicationCommands(applications, applicationid, guilds, guildid, commands string) string {
	return applications + applicationid + guilds + guildid + commands
}

// EndpointCreateGuildApplicationCommandbuilds a query for an HTTP request.
func EndpointCreateGuildApplicationCommand(applications, applicationid, guilds, guildid, commands string) string {
	return applications + applicationid + guilds + guildid + commands
}

// EndpointGetGuildApplicationCommandbuilds a query for an HTTP request.
func EndpointGetGuildApplicationCommand(applications, applicationid, guilds, guildid, commands, commandid string) string {
	return applications + applicationid + guilds + guildid + commands + commandid
}

// EndpointEditGuildApplicationCommandbuilds a query for an HTTP request.
func EndpointEditGuildApplicationCommand(applications, applicationid, guilds, guildid, commands, commandid string) string {
	return applications + applicationid + guilds + guildid + commands + commandid
}

// EndpointDeleteGuildApplicationCommandbuilds a query for an HTTP request.
func EndpointDeleteGuildApplicationCommand(applications, applicationid, guilds, guildid, commands, commandid string) string {
	return applications + applicationid + guilds + guildid + commands + commandid
}

// EndpointBulkOverwriteGuildApplicationCommandsbuilds a query for an HTTP request.
func EndpointBulkOverwriteGuildApplicationCommands(applications, applicationid, guilds, guildid, commands string) string {
	return applications + applicationid + guilds + guildid + commands
}

// EndpointGetGuildApplicationCommandPermissionsbuilds a query for an HTTP request.
func EndpointGetGuildApplicationCommandPermissions(applications, applicationid, guilds, guildid, commands, permissions string) string {
	return applications + applicationid + guilds + guildid + commands + permissions
}

// EndpointGetApplicationCommandPermissionsbuilds a query for an HTTP request.
func EndpointGetApplicationCommandPermissions(applications, applicationid, guilds, guildid, commands, commandid, permissions string) string {
	return applications + applicationid + guilds + guildid + commands + commandid + permissions
}

// EndpointEditApplicationCommandPermissionsbuilds a query for an HTTP request.
func EndpointEditApplicationCommandPermissions(applications, applicationid, guilds, guildid, commands, commandid, permissions string) string {
	return applications + applicationid + guilds + guildid + commands + commandid + permissions
}

// EndpointBatchEditApplicationCommandPermissionsbuilds a query for an HTTP request.
func EndpointBatchEditApplicationCommandPermissions(applications, applicationid, guilds, guildid, commands, permissions string) string {
	return applications + applicationid + guilds + guildid + commands + permissions
}

// EndpointCreateInteractionResponsebuilds a query for an HTTP request.
func EndpointCreateInteractionResponse(interactions, interactionid, interactiontoken, callback string) string {
	return interactions + interactionid + interactiontoken + callback
}

// EndpointGetOriginalInteractionResponsebuilds a query for an HTTP request.
func EndpointGetOriginalInteractionResponse(webhooks, applicationid, interactiontoken, messages, original string) string {
	return webhooks + applicationid + interactiontoken + messages + original
}

// EndpointEditOriginalInteractionResponsebuilds a query for an HTTP request.
func EndpointEditOriginalInteractionResponse(webhooks, applicationid, interactiontoken, messages, original string) string {
	return webhooks + applicationid + interactiontoken + messages + original
}

// EndpointDeleteOriginalInteractionResponsebuilds a query for an HTTP request.
func EndpointDeleteOriginalInteractionResponse(webhooks, applicationid, interactiontoken, messages, original string) string {
	return webhooks + applicationid + interactiontoken + messages + original
}

// EndpointCreateFollowupMessagebuilds a query for an HTTP request.
func EndpointCreateFollowupMessage(webhooks, applicationid, interactiontoken string) string {
	return webhooks + applicationid + interactiontoken
}

// EndpointGetFollowupMessagebuilds a query for an HTTP request.
func EndpointGetFollowupMessage(webhooks, applicationid, interactiontoken, messages, messageid string) string {
	return webhooks + applicationid + interactiontoken + messages + messageid
}

// EndpointEditFollowupMessagebuilds a query for an HTTP request.
func EndpointEditFollowupMessage(webhooks, applicationid, interactiontoken, messages, messageid string) string {
	return webhooks + applicationid + interactiontoken + messages + messageid
}

// EndpointDeleteFollowupMessagebuilds a query for an HTTP request.
func EndpointDeleteFollowupMessage(webhooks, applicationid, interactiontoken, messages, messageid string) string {
	return webhooks + applicationid + interactiontoken + messages + messageid
}

// EndpointGetGuildAuditLogbuilds a query for an HTTP request.
func EndpointGetGuildAuditLog(guilds, guildid, auditlogs string) string {
	return guilds + guildid + auditlogs
}

// EndpointGetChannelbuilds a query for an HTTP request.
func EndpointGetChannel(channels, channelid string) string {
	return channels + channelid
}

// EndpointModifyChannelbuilds a query for an HTTP request.
func EndpointModifyChannel(channels, channelid string) string {
	return channels + channelid
}

// EndpointDeleteCloseChannelbuilds a query for an HTTP request.
func EndpointDeleteCloseChannel(channels, channelid string) string {
	return channels + channelid
}

// EndpointGetChannelMessagesbuilds a query for an HTTP request.
func EndpointGetChannelMessages(channels, channelid, messages string) string {
	return channels + channelid + messages
}

// EndpointGetChannelMessagebuilds a query for an HTTP request.
func EndpointGetChannelMessage(channels, channelid, messages, messageid string) string {
	return channels + channelid + messages + messageid
}

// EndpointCreateMessagebuilds a query for an HTTP request.
func EndpointCreateMessage(channels, channelid, messages string) string {
	return channels + channelid + messages
}

// EndpointCrosspostMessagebuilds a query for an HTTP request.
func EndpointCrosspostMessage(channels, channelid, messages, messageid, crosspost string) string {
	return channels + channelid + messages + messageid + crosspost
}

// EndpointCreateReactionbuilds a query for an HTTP request.
func EndpointCreateReaction(channels, channelid, messages, messageid, reactions, emoji, me string) string {
	return channels + channelid + messages + messageid + reactions + emoji + me
}

// EndpointDeleteOwnReactionbuilds a query for an HTTP request.
func EndpointDeleteOwnReaction(channels, channelid, messages, messageid, reactions, emoji, me string) string {
	return channels + channelid + messages + messageid + reactions + emoji + me
}

// EndpointDeleteUserReactionbuilds a query for an HTTP request.
func EndpointDeleteUserReaction(channels, channelid, messages, messageid, reactions, emoji, userid string) string {
	return channels + channelid + messages + messageid + reactions + emoji + userid
}

// EndpointGetReactionsbuilds a query for an HTTP request.
func EndpointGetReactions(channels, channelid, messages, messageid, reactions, emoji string) string {
	return channels + channelid + messages + messageid + reactions + emoji
}

// EndpointDeleteAllReactionsbuilds a query for an HTTP request.
func EndpointDeleteAllReactions(channels, channelid, messages, messageid, reactions string) string {
	return channels + channelid + messages + messageid + reactions
}

// EndpointDeleteAllReactionsforEmojibuilds a query for an HTTP request.
func EndpointDeleteAllReactionsforEmoji(channels, channelid, messages, messageid, reactions, emoji string) string {
	return channels + channelid + messages + messageid + reactions + emoji
}

// EndpointEditMessagebuilds a query for an HTTP request.
func EndpointEditMessage(channels, channelid, messages, messageid string) string {
	return channels + channelid + messages + messageid
}

// EndpointDeleteMessagebuilds a query for an HTTP request.
func EndpointDeleteMessage(channels, channelid, messages, messageid string) string {
	return channels + channelid + messages + messageid
}

// EndpointBulkDeleteMessagesbuilds a query for an HTTP request.
func EndpointBulkDeleteMessages(channels, channelid, messages, bulkdelete string) string {
	return channels + channelid + messages + bulkdelete
}

// EndpointEditChannelPermissionsbuilds a query for an HTTP request.
func EndpointEditChannelPermissions(channels, channelid, permissions, overwriteid string) string {
	return channels + channelid + permissions + overwriteid
}

// EndpointGetChannelInvitesbuilds a query for an HTTP request.
func EndpointGetChannelInvites(channels, channelid, invites string) string {
	return channels + channelid + invites
}

// EndpointCreateChannelInvitebuilds a query for an HTTP request.
func EndpointCreateChannelInvite(channels, channelid, invites string) string {
	return channels + channelid + invites
}

// EndpointDeleteChannelPermissionbuilds a query for an HTTP request.
func EndpointDeleteChannelPermission(channels, channelid, permissions, overwriteid string) string {
	return channels + channelid + permissions + overwriteid
}

// EndpointFollowNewsChannelbuilds a query for an HTTP request.
func EndpointFollowNewsChannel(channels, channelid, followers string) string {
	return channels + channelid + followers
}

// EndpointTriggerTypingIndicatorbuilds a query for an HTTP request.
func EndpointTriggerTypingIndicator(channels, channelid, typing string) string {
	return channels + channelid + typing
}

// EndpointGetPinnedMessagesbuilds a query for an HTTP request.
func EndpointGetPinnedMessages(channels, channelid, pins string) string {
	return channels + channelid + pins
}

// EndpointPinMessagebuilds a query for an HTTP request.
func EndpointPinMessage(channels, channelid, pins, messageid string) string {
	return channels + channelid + pins + messageid
}

// EndpointUnpinMessagebuilds a query for an HTTP request.
func EndpointUnpinMessage(channels, channelid, pins, messageid string) string {
	return channels + channelid + pins + messageid
}

// EndpointGroupDMAddRecipientbuilds a query for an HTTP request.
func EndpointGroupDMAddRecipient(channels, channelid, recipients, userid string) string {
	return channels + channelid + recipients + userid
}

// EndpointGroupDMRemoveRecipientbuilds a query for an HTTP request.
func EndpointGroupDMRemoveRecipient(channels, channelid, recipients, userid string) string {
	return channels + channelid + recipients + userid
}

// EndpointStartThreadfromMessagebuilds a query for an HTTP request.
func EndpointStartThreadfromMessage(channels, channelid, messages, messageid, threads string) string {
	return channels + channelid + messages + messageid + threads
}

// EndpointStartThreadwithoutMessagebuilds a query for an HTTP request.
func EndpointStartThreadwithoutMessage(channels, channelid, threads string) string {
	return channels + channelid + threads
}

// EndpointStartThreadinForumChannelbuilds a query for an HTTP request.
func EndpointStartThreadinForumChannel(channels, channelid, threads string) string {
	return channels + channelid + threads
}

// EndpointJoinThreadbuilds a query for an HTTP request.
func EndpointJoinThread(channels, channelid, threadmembers, me string) string {
	return channels + channelid + threadmembers + me
}

// EndpointAddThreadMemberbuilds a query for an HTTP request.
func EndpointAddThreadMember(channels, channelid, threadmembers, userid string) string {
	return channels + channelid + threadmembers + userid
}

// EndpointLeaveThreadbuilds a query for an HTTP request.
func EndpointLeaveThread(channels, channelid, threadmembers, me string) string {
	return channels + channelid + threadmembers + me
}

// EndpointRemoveThreadMemberbuilds a query for an HTTP request.
func EndpointRemoveThreadMember(channels, channelid, threadmembers, userid string) string {
	return channels + channelid + threadmembers + userid
}

// EndpointGetThreadMemberbuilds a query for an HTTP request.
func EndpointGetThreadMember(channels, channelid, threadmembers, userid string) string {
	return channels + channelid + threadmembers + userid
}

// EndpointListThreadMembersbuilds a query for an HTTP request.
func EndpointListThreadMembers(channels, channelid, threadmembers string) string {
	return channels + channelid + threadmembers
}

// EndpointListActiveThreadsbuilds a query for an HTTP request.
func EndpointListActiveThreads(channels, channelid, threads, active string) string {
	return channels + channelid + threads + active
}

// EndpointListPublicArchivedThreadsbuilds a query for an HTTP request.
func EndpointListPublicArchivedThreads(channels, channelid, threads, archived, public string) string {
	return channels + channelid + threads + archived + public
}

// EndpointListPrivateArchivedThreadsbuilds a query for an HTTP request.
func EndpointListPrivateArchivedThreads(channels, channelid, threads, archived, private string) string {
	return channels + channelid + threads + archived + private
}

// EndpointListJoinedPrivateArchivedThreadsbuilds a query for an HTTP request.
func EndpointListJoinedPrivateArchivedThreads(channels, channelid, users, me, threads, archived, private string) string {
	return channels + channelid + users + me + threads + archived + private
}

// EndpointListGuildEmojisbuilds a query for an HTTP request.
func EndpointListGuildEmojis(guilds, guildid, emojis string) string {
	return guilds + guildid + emojis
}

// EndpointGetGuildEmojibuilds a query for an HTTP request.
func EndpointGetGuildEmoji(guilds, guildid, emojis, emojiid string) string {
	return guilds + guildid + emojis + emojiid
}

// EndpointCreateGuildEmojibuilds a query for an HTTP request.
func EndpointCreateGuildEmoji(guilds, guildid, emojis string) string {
	return guilds + guildid + emojis
}

// EndpointModifyGuildEmojibuilds a query for an HTTP request.
func EndpointModifyGuildEmoji(guilds, guildid, emojis, emojiid string) string {
	return guilds + guildid + emojis + emojiid
}

// EndpointDeleteGuildEmojibuilds a query for an HTTP request.
func EndpointDeleteGuildEmoji(guilds, guildid, emojis, emojiid string) string {
	return guilds + guildid + emojis + emojiid
}

// EndpointListScheduledEventsforGuildbuilds a query for an HTTP request.
func EndpointListScheduledEventsforGuild(guilds, guildid, scheduledevents string) string {
	return guilds + guildid + scheduledevents
}

// EndpointCreateGuildScheduledEventbuilds a query for an HTTP request.
func EndpointCreateGuildScheduledEvent(guilds, guildid, scheduledevents string) string {
	return guilds + guildid + scheduledevents
}

// EndpointGetGuildScheduledEventbuilds a query for an HTTP request.
func EndpointGetGuildScheduledEvent(guilds, guildid, scheduledevents, guildscheduledeventid string) string {
	return guilds + guildid + scheduledevents + guildscheduledeventid
}

// EndpointModifyGuildScheduledEventbuilds a query for an HTTP request.
func EndpointModifyGuildScheduledEvent(guilds, guildid, scheduledevents, guildscheduledeventid string) string {
	return guilds + guildid + scheduledevents + guildscheduledeventid
}

// EndpointDeleteGuildScheduledEventbuilds a query for an HTTP request.
func EndpointDeleteGuildScheduledEvent(guilds, guildid, scheduledevents, guildscheduledeventid string) string {
	return guilds + guildid + scheduledevents + guildscheduledeventid
}

// EndpointGetGuildScheduledEventUsersbuilds a query for an HTTP request.
func EndpointGetGuildScheduledEventUsers(guilds, guildid, scheduledevents, guildscheduledeventid, users string) string {
	return guilds + guildid + scheduledevents + guildscheduledeventid + users
}

// EndpointGetGuildTemplatebuilds a query for an HTTP request.
func EndpointGetGuildTemplate(guilds, templates, templatecode string) string {
	return guilds + templates + templatecode
}

// EndpointCreateGuildfromGuildTemplatebuilds a query for an HTTP request.
func EndpointCreateGuildfromGuildTemplate(guilds, templates, templatecode string) string {
	return guilds + templates + templatecode
}

// EndpointGetGuildTemplatesbuilds a query for an HTTP request.
func EndpointGetGuildTemplates(guilds, guildid, templates string) string {
	return guilds + guildid + templates
}

// EndpointCreateGuildTemplatebuilds a query for an HTTP request.
func EndpointCreateGuildTemplate(guilds, guildid, templates string) string {
	return guilds + guildid + templates
}

// EndpointSyncGuildTemplatebuilds a query for an HTTP request.
func EndpointSyncGuildTemplate(guilds, guildid, templates, templatecode string) string {
	return guilds + guildid + templates + templatecode
}

// EndpointModifyGuildTemplatebuilds a query for an HTTP request.
func EndpointModifyGuildTemplate(guilds, guildid, templates, templatecode string) string {
	return guilds + guildid + templates + templatecode
}

// EndpointDeleteGuildTemplatebuilds a query for an HTTP request.
func EndpointDeleteGuildTemplate(guilds, guildid, templates, templatecode string) string {
	return guilds + guildid + templates + templatecode
}

// EndpointCreateGuildbuilds a query for an HTTP request.
func EndpointCreateGuild(guilds string) string {
	return guilds
}

// EndpointGetGuildbuilds a query for an HTTP request.
func EndpointGetGuild(guilds, guildid string) string {
	return guilds + guildid
}

// EndpointGetGuildPreviewbuilds a query for an HTTP request.
func EndpointGetGuildPreview(guilds, guildid, preview string) string {
	return guilds + guildid + preview
}

// EndpointModifyGuildbuilds a query for an HTTP request.
func EndpointModifyGuild(guilds, guildid string) string {
	return guilds + guildid
}

// EndpointDeleteGuildbuilds a query for an HTTP request.
func EndpointDeleteGuild(guilds, guildid string) string {
	return guilds + guildid
}

// EndpointGetGuildChannelsbuilds a query for an HTTP request.
func EndpointGetGuildChannels(guilds, guildid, channels string) string {
	return guilds + guildid + channels
}

// EndpointCreateGuildChannelbuilds a query for an HTTP request.
func EndpointCreateGuildChannel(guilds, guildid, channels string) string {
	return guilds + guildid + channels
}

// EndpointModifyGuildChannelPositionsbuilds a query for an HTTP request.
func EndpointModifyGuildChannelPositions(guilds, guildid, channels string) string {
	return guilds + guildid + channels
}

// EndpointListActiveGuildThreadsbuilds a query for an HTTP request.
func EndpointListActiveGuildThreads(guilds, guildid, threads, active string) string {
	return guilds + guildid + threads + active
}

// EndpointGetGuildMemberbuilds a query for an HTTP request.
func EndpointGetGuildMember(guilds, guildid, members, userid string) string {
	return guilds + guildid + members + userid
}

// EndpointListGuildMembersbuilds a query for an HTTP request.
func EndpointListGuildMembers(guilds, guildid, members string) string {
	return guilds + guildid + members
}

// EndpointSearchGuildMembersbuilds a query for an HTTP request.
func EndpointSearchGuildMembers(guilds, guildid, members, search string) string {
	return guilds + guildid + members + search
}

// EndpointAddGuildMemberbuilds a query for an HTTP request.
func EndpointAddGuildMember(guilds, guildid, members, userid string) string {
	return guilds + guildid + members + userid
}

// EndpointModifyGuildMemberbuilds a query for an HTTP request.
func EndpointModifyGuildMember(guilds, guildid, members, userid string) string {
	return guilds + guildid + members + userid
}

// EndpointModifyCurrentMemberbuilds a query for an HTTP request.
func EndpointModifyCurrentMember(guilds, guildid, members, me string) string {
	return guilds + guildid + members + me
}

// EndpointModifyCurrentUserNickbuilds a query for an HTTP request.
func EndpointModifyCurrentUserNick(guilds, guildid, members, me, nick string) string {
	return guilds + guildid + members + me + nick
}

// EndpointAddGuildMemberRolebuilds a query for an HTTP request.
func EndpointAddGuildMemberRole(guilds, guildid, members, userid, roles, roleid string) string {
	return guilds + guildid + members + userid + roles + roleid
}

// EndpointRemoveGuildMemberRolebuilds a query for an HTTP request.
func EndpointRemoveGuildMemberRole(guilds, guildid, members, userid, roles, roleid string) string {
	return guilds + guildid + members + userid + roles + roleid
}

// EndpointRemoveGuildMemberbuilds a query for an HTTP request.
func EndpointRemoveGuildMember(guilds, guildid, members, userid string) string {
	return guilds + guildid + members + userid
}

// EndpointGetGuildBansbuilds a query for an HTTP request.
func EndpointGetGuildBans(guilds, guildid, bans string) string {
	return guilds + guildid + bans
}

// EndpointGetGuildBanbuilds a query for an HTTP request.
func EndpointGetGuildBan(guilds, guildid, bans, userid string) string {
	return guilds + guildid + bans + userid
}

// EndpointCreateGuildBanbuilds a query for an HTTP request.
func EndpointCreateGuildBan(guilds, guildid, bans, userid string) string {
	return guilds + guildid + bans + userid
}

// EndpointRemoveGuildBanbuilds a query for an HTTP request.
func EndpointRemoveGuildBan(guilds, guildid, bans, userid string) string {
	return guilds + guildid + bans + userid
}

// EndpointGetGuildRolesbuilds a query for an HTTP request.
func EndpointGetGuildRoles(guilds, guildid, roles string) string {
	return guilds + guildid + roles
}

// EndpointCreateGuildRolebuilds a query for an HTTP request.
func EndpointCreateGuildRole(guilds, guildid, roles string) string {
	return guilds + guildid + roles
}

// EndpointModifyGuildRolePositionsbuilds a query for an HTTP request.
func EndpointModifyGuildRolePositions(guilds, guildid, roles string) string {
	return guilds + guildid + roles
}

// EndpointModifyGuildRolebuilds a query for an HTTP request.
func EndpointModifyGuildRole(guilds, guildid, roles, roleid string) string {
	return guilds + guildid + roles + roleid
}

// EndpointDeleteGuildRolebuilds a query for an HTTP request.
func EndpointDeleteGuildRole(guilds, guildid, roles, roleid string) string {
	return guilds + guildid + roles + roleid
}

// EndpointGetGuildPruneCountbuilds a query for an HTTP request.
func EndpointGetGuildPruneCount(guilds, guildid, prune string) string {
	return guilds + guildid + prune
}

// EndpointBeginGuildPrunebuilds a query for an HTTP request.
func EndpointBeginGuildPrune(guilds, guildid, prune string) string {
	return guilds + guildid + prune
}

// EndpointGetGuildVoiceRegionsbuilds a query for an HTTP request.
func EndpointGetGuildVoiceRegions(guilds, guildid, regions string) string {
	return guilds + guildid + regions
}

// EndpointGetGuildInvitesbuilds a query for an HTTP request.
func EndpointGetGuildInvites(guilds, guildid, invites string) string {
	return guilds + guildid + invites
}

// EndpointGetGuildIntegrationsbuilds a query for an HTTP request.
func EndpointGetGuildIntegrations(guilds, guildid, integrations string) string {
	return guilds + guildid + integrations
}

// EndpointDeleteGuildIntegrationbuilds a query for an HTTP request.
func EndpointDeleteGuildIntegration(guilds, guildid, integrations, integrationid string) string {
	return guilds + guildid + integrations + integrationid
}

// EndpointGetGuildWidgetSettingsbuilds a query for an HTTP request.
func EndpointGetGuildWidgetSettings(guilds, guildid, widget string) string {
	return guilds + guildid + widget
}

// EndpointModifyGuildWidgetbuilds a query for an HTTP request.
func EndpointModifyGuildWidget(guilds, guildid, widget string) string {
	return guilds + guildid + widget
}

// EndpointGetGuildWidgetbuilds a query for an HTTP request.
func EndpointGetGuildWidget(guilds, guildid, widgetjson string) string {
	return guilds + guildid + widgetjson
}

// EndpointGetGuildVanityURLbuilds a query for an HTTP request.
func EndpointGetGuildVanityURL(guilds, guildid, vanityurl string) string {
	return guilds + guildid + vanityurl
}

// EndpointGetGuildWidgetImagebuilds a query for an HTTP request.
func EndpointGetGuildWidgetImage(guilds, guildid, widgetpng string) string {
	return guilds + guildid + widgetpng
}

// EndpointGetGuildWelcomeScreenbuilds a query for an HTTP request.
func EndpointGetGuildWelcomeScreen(guilds, guildid, welcomescreen string) string {
	return guilds + guildid + welcomescreen
}

// EndpointModifyGuildWelcomeScreenbuilds a query for an HTTP request.
func EndpointModifyGuildWelcomeScreen(guilds, guildid, welcomescreen string) string {
	return guilds + guildid + welcomescreen
}

// EndpointModifyCurrentUserVoiceStatebuilds a query for an HTTP request.
func EndpointModifyCurrentUserVoiceState(guilds, guildid, voicestates, me string) string {
	return guilds + guildid + voicestates + me
}

// EndpointModifyUserVoiceStatebuilds a query for an HTTP request.
func EndpointModifyUserVoiceState(guilds, guildid, voicestates, userid string) string {
	return guilds + guildid + voicestates + userid
}

// EndpointGetInvitebuilds a query for an HTTP request.
func EndpointGetInvite(invites, invitecode string) string {
	return invites + invitecode
}

// EndpointDeleteInvitebuilds a query for an HTTP request.
func EndpointDeleteInvite(invites, invitecode string) string {
	return invites + invitecode
}

// EndpointCreateStageInstancebuilds a query for an HTTP request.
func EndpointCreateStageInstance(stageinstances string) string {
	return stageinstances
}

// EndpointGetStageInstancebuilds a query for an HTTP request.
func EndpointGetStageInstance(stageinstances, channelid string) string {
	return stageinstances + channelid
}

// EndpointModifyStageInstancebuilds a query for an HTTP request.
func EndpointModifyStageInstance(stageinstances, channelid string) string {
	return stageinstances + channelid
}

// EndpointDeleteStageInstancebuilds a query for an HTTP request.
func EndpointDeleteStageInstance(stageinstances, channelid string) string {
	return stageinstances + channelid
}

// EndpointGetStickerbuilds a query for an HTTP request.
func EndpointGetSticker(stickers, stickerid string) string {
	return stickers + stickerid
}

// EndpointListNitroStickerPacksbuilds a query for an HTTP request.
func EndpointListNitroStickerPacks(stickerpacks string) string {
	return stickerpacks
}

// EndpointListGuildStickersbuilds a query for an HTTP request.
func EndpointListGuildStickers(guilds, guildid, stickers string) string {
	return guilds + guildid + stickers
}

// EndpointGetGuildStickerbuilds a query for an HTTP request.
func EndpointGetGuildSticker(guilds, guildid, stickers, stickerid string) string {
	return guilds + guildid + stickers + stickerid
}

// EndpointCreateGuildStickerbuilds a query for an HTTP request.
func EndpointCreateGuildSticker(guilds, guildid, stickers string) string {
	return guilds + guildid + stickers
}

// EndpointModifyGuildStickerbuilds a query for an HTTP request.
func EndpointModifyGuildSticker(guilds, guildid, stickers, stickerid string) string {
	return guilds + guildid + stickers + stickerid
}

// EndpointDeleteGuildStickerbuilds a query for an HTTP request.
func EndpointDeleteGuildSticker(guilds, guildid, stickers, stickerid string) string {
	return guilds + guildid + stickers + stickerid
}

// EndpointGetCurrentUserbuilds a query for an HTTP request.
func EndpointGetCurrentUser(users, me string) string {
	return users + me
}

// EndpointGetUserbuilds a query for an HTTP request.
func EndpointGetUser(users, userid string) string {
	return users + userid
}

// EndpointModifyCurrentUserbuilds a query for an HTTP request.
func EndpointModifyCurrentUser(users, me string) string {
	return users + me
}

// EndpointGetCurrentUserGuildsbuilds a query for an HTTP request.
func EndpointGetCurrentUserGuilds(users, me, guilds string) string {
	return users + me + guilds
}

// EndpointGetCurrentUserGuildMemberbuilds a query for an HTTP request.
func EndpointGetCurrentUserGuildMember(users, me, guilds, guildid, member string) string {
	return users + me + guilds + guildid + member
}

// EndpointLeaveGuildbuilds a query for an HTTP request.
func EndpointLeaveGuild(users, me, guilds, guildid string) string {
	return users + me + guilds + guildid
}

// EndpointCreateDMbuilds a query for an HTTP request.
func EndpointCreateDM(users, me, channels string) string {
	return users + me + channels
}

// EndpointCreateGroupDMbuilds a query for an HTTP request.
func EndpointCreateGroupDM(users, me, channels string) string {
	return users + me + channels
}

// EndpointGetUserConnectionsbuilds a query for an HTTP request.
func EndpointGetUserConnections(users, me, connections string) string {
	return users + me + connections
}

// EndpointListVoiceRegionsbuilds a query for an HTTP request.
func EndpointListVoiceRegions(voice, regions string) string {
	return voice + regions
}

// EndpointCreateWebhookbuilds a query for an HTTP request.
func EndpointCreateWebhook(channels, channelid, webhooks string) string {
	return channels + channelid + webhooks
}

// EndpointGetChannelWebhooksbuilds a query for an HTTP request.
func EndpointGetChannelWebhooks(channels, channelid, webhooks string) string {
	return channels + channelid + webhooks
}

// EndpointGetGuildWebhooksbuilds a query for an HTTP request.
func EndpointGetGuildWebhooks(guilds, guildid, webhooks string) string {
	return guilds + guildid + webhooks
}

// EndpointGetWebhookbuilds a query for an HTTP request.
func EndpointGetWebhook(webhooks, webhookid string) string {
	return webhooks + webhookid
}

// EndpointGetWebhookwithTokenbuilds a query for an HTTP request.
func EndpointGetWebhookwithToken(webhooks, webhookid, webhooktoken string) string {
	return webhooks + webhookid + webhooktoken
}

// EndpointModifyWebhookbuilds a query for an HTTP request.
func EndpointModifyWebhook(webhooks, webhookid string) string {
	return webhooks + webhookid
}

// EndpointModifyWebhookwithTokenbuilds a query for an HTTP request.
func EndpointModifyWebhookwithToken(webhooks, webhookid, webhooktoken string) string {
	return webhooks + webhookid + webhooktoken
}

// EndpointDeleteWebhookbuilds a query for an HTTP request.
func EndpointDeleteWebhook(webhooks, webhookid string) string {
	return webhooks + webhookid
}

// EndpointDeleteWebhookwithTokenbuilds a query for an HTTP request.
func EndpointDeleteWebhookwithToken(webhooks, webhookid, webhooktoken string) string {
	return webhooks + webhookid + webhooktoken
}

// EndpointExecuteWebhookbuilds a query for an HTTP request.
func EndpointExecuteWebhook(webhooks, webhookid, webhooktoken string) string {
	return webhooks + webhookid + webhooktoken
}

// EndpointExecuteSlackCompatibleWebhookbuilds a query for an HTTP request.
func EndpointExecuteSlackCompatibleWebhook(webhooks, webhookid, webhooktoken, slack string) string {
	return webhooks + webhookid + webhooktoken + slack
}

// EndpointExecuteGitHubCompatibleWebhookbuilds a query for an HTTP request.
func EndpointExecuteGitHubCompatibleWebhook(webhooks, webhookid, webhooktoken, github string) string {
	return webhooks + webhookid + webhooktoken + github
}

// EndpointGetWebhookMessagebuilds a query for an HTTP request.
func EndpointGetWebhookMessage(webhooks, webhookid, webhooktoken, messages, messageid string) string {
	return webhooks + webhookid + webhooktoken + messages + messageid
}

// EndpointEditWebhookMessagebuilds a query for an HTTP request.
func EndpointEditWebhookMessage(webhooks, webhookid, webhooktoken, messages, messageid string) string {
	return webhooks + webhookid + webhooktoken + messages + messageid
}

// EndpointDeleteWebhookMessagebuilds a query for an HTTP request.
func EndpointDeleteWebhookMessage(webhooks, webhookid, webhooktoken, messages, messageid string) string {
	return webhooks + webhookid + webhooktoken + messages + messageid
}

// EndpointGetGatewaybuilds a query for an HTTP request.
func EndpointGetGateway(gateway string) string {
	return gateway
}

// EndpointGetGatewayBotbuilds a query for an HTTP request.
func EndpointGetGatewayBot(gateway, bot string) string {
	return gateway + bot
}

// EndpointGetCurrentBotApplicationInformationbuilds a query for an HTTP request.
func EndpointGetCurrentBotApplicationInformation(oauth, applications, me string) string {
	return oauth + applications + me
}

// EndpointGetCurrentAuthorizationInformationbuilds a query for an HTTP request.
func EndpointGetCurrentAuthorizationInformation(oauth, me string) string {
	return oauth + me
}
