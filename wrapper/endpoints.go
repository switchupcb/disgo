package wrapper

// Discord API Endpoints
const (
	EndpointBaseURL = "https://discord.com/api/v9/"
	roles           = "roles"
	voicestates     = "voice-states"
	slack           = "slack"
	bot             = "bot"
	commands        = "commands"
	active          = "active"
	prune           = "prune"
	applications    = "applications"
	me              = "@me"
	archived        = "archived"
	emojis          = "emojis"
	connections     = "connections"
	slash           = "/"
	original        = "@original"
	invites         = "invites"
	interactions    = "interactions"
	crosspost       = "crosspost"
	followers       = "followers"
	threads         = "threads"
	widgetjson      = "widget.json"
	voice           = "voice"
	github          = "github"
	typing          = "typing"
	public          = "public"
	templates       = "templates"
	integrations    = "integrations"
	stickers        = "stickers"
	stickerpacks    = "sticker-packs"
	permissions     = "permissions"
	messages        = "messages"
	threadmembers   = "thread-members"
	preview         = "preview"
	guilds          = "guilds"
	recipients      = "recipients"
	reactions       = "reactions"
	pins            = "pins"
	scheduledevents = "scheduled-events"
	bulkdelete      = "bulk-delete"
	stageinstances  = "stage-instances"
	members         = "members"
	bans            = "bans"
	widgetpng       = "widget.png"
	vanityurl       = "vanity-url"
	member          = "member"
	auditlogs       = "audit-logs"
	channels        = "channels"
	search          = "search"
	gateway         = "gateway"
	webhooks        = "webhooks"
	nick            = "nick"
	welcomescreen   = "welcome-screen"
	private         = "private"
	users           = "users"
	regions         = "regions"
	widget          = "widget"
	oauth           = "oauth2"
	callback        = "callback"
)

// EndpointGetGlobalApplicationCommands builds a query for an HTTP request.
func EndpointGetGlobalApplicationCommands(applicationid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + commands
}

// EndpointCreateGlobalApplicationCommand builds a query for an HTTP request.
func EndpointCreateGlobalApplicationCommand(applicationid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + commands
}

// EndpointGetGlobalApplicationCommand builds a query for an HTTP request.
func EndpointGetGlobalApplicationCommand(applicationid, commandid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + commands + slash + commandid
}

// EndpointEditGlobalApplicationCommand builds a query for an HTTP request.
func EndpointEditGlobalApplicationCommand(applicationid, commandid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + commands + slash + commandid
}

// EndpointDeleteGlobalApplicationCommand builds a query for an HTTP request.
func EndpointDeleteGlobalApplicationCommand(applicationid, commandid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + commands + slash + commandid
}

// EndpointBulkOverwriteGlobalApplicationCommands builds a query for an HTTP request.
func EndpointBulkOverwriteGlobalApplicationCommands(applicationid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + commands
}

// EndpointGetGuildApplicationCommands builds a query for an HTTP request.
func EndpointGetGuildApplicationCommands(applicationid, guildid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + guilds + slash + guildid + slash + commands
}

// EndpointCreateGuildApplicationCommand builds a query for an HTTP request.
func EndpointCreateGuildApplicationCommand(applicationid, guildid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + guilds + slash + guildid + slash + commands
}

// EndpointGetGuildApplicationCommand builds a query for an HTTP request.
func EndpointGetGuildApplicationCommand(applicationid, guildid, commandid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + guilds + slash + guildid + slash + commands + slash + commandid
}

// EndpointEditGuildApplicationCommand builds a query for an HTTP request.
func EndpointEditGuildApplicationCommand(applicationid, guildid, commandid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + guilds + slash + guildid + slash + commands + slash + commandid
}

// EndpointDeleteGuildApplicationCommand builds a query for an HTTP request.
func EndpointDeleteGuildApplicationCommand(applicationid, guildid, commandid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + guilds + slash + guildid + slash + commands + slash + commandid
}

// EndpointBulkOverwriteGuildApplicationCommands builds a query for an HTTP request.
func EndpointBulkOverwriteGuildApplicationCommands(applicationid, guildid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + guilds + slash + guildid + slash + commands
}

// EndpointGetGuildApplicationCommandPermissions builds a query for an HTTP request.
func EndpointGetGuildApplicationCommandPermissions(applicationid, guildid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + guilds + slash + guildid + slash + commands + slash + permissions
}

// EndpointGetApplicationCommandPermissions builds a query for an HTTP request.
func EndpointGetApplicationCommandPermissions(applicationid, guildid, commandid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + guilds + slash + guildid + slash + commands + slash + commandid + slash + permissions
}

// EndpointEditApplicationCommandPermissions builds a query for an HTTP request.
func EndpointEditApplicationCommandPermissions(applicationid, guildid, commandid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + guilds + slash + guildid + slash + commands + slash + commandid + slash + permissions
}

// EndpointBatchEditApplicationCommandPermissions builds a query for an HTTP request.
func EndpointBatchEditApplicationCommandPermissions(applicationid, guildid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + guilds + slash + guildid + slash + commands + slash + permissions
}

// EndpointCreateInteractionResponse builds a query for an HTTP request.
func EndpointCreateInteractionResponse(interactionid, interactiontoken string) string {
	return EndpointBaseURL + interactions + slash + interactionid + slash + interactiontoken + slash + callback
}

// EndpointGetOriginalInteractionResponse builds a query for an HTTP request.
func EndpointGetOriginalInteractionResponse(applicationid, interactiontoken string) string {
	return EndpointBaseURL + webhooks + slash + applicationid + slash + interactiontoken + slash + messages + slash + original
}

// EndpointEditOriginalInteractionResponse builds a query for an HTTP request.
func EndpointEditOriginalInteractionResponse(applicationid, interactiontoken string) string {
	return EndpointBaseURL + webhooks + slash + applicationid + slash + interactiontoken + slash + messages + slash + original
}

// EndpointDeleteOriginalInteractionResponse builds a query for an HTTP request.
func EndpointDeleteOriginalInteractionResponse(applicationid, interactiontoken string) string {
	return EndpointBaseURL + webhooks + slash + applicationid + slash + interactiontoken + slash + messages + slash + original
}

// EndpointCreateFollowupMessage builds a query for an HTTP request.
func EndpointCreateFollowupMessage(applicationid, interactiontoken string) string {
	return EndpointBaseURL + webhooks + slash + applicationid + slash + interactiontoken
}

// EndpointGetFollowupMessage builds a query for an HTTP request.
func EndpointGetFollowupMessage(applicationid, interactiontoken, messageid string) string {
	return EndpointBaseURL + webhooks + slash + applicationid + slash + interactiontoken + slash + messages + slash + messageid
}

// EndpointEditFollowupMessage builds a query for an HTTP request.
func EndpointEditFollowupMessage(applicationid, interactiontoken, messageid string) string {
	return EndpointBaseURL + webhooks + slash + applicationid + slash + interactiontoken + slash + messages + slash + messageid
}

// EndpointDeleteFollowupMessage builds a query for an HTTP request.
func EndpointDeleteFollowupMessage(applicationid, interactiontoken, messageid string) string {
	return EndpointBaseURL + webhooks + slash + applicationid + slash + interactiontoken + slash + messages + slash + messageid
}

// EndpointGetGuildAuditLog builds a query for an HTTP request.
func EndpointGetGuildAuditLog(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + auditlogs
}

// EndpointGetChannel builds a query for an HTTP request.
func EndpointGetChannel(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid
}

// EndpointModifyChannel builds a query for an HTTP request.
func EndpointModifyChannel(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid
}

// EndpointDeleteCloseChannel builds a query for an HTTP request.
func EndpointDeleteCloseChannel(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid
}

// EndpointGetChannelMessages builds a query for an HTTP request.
func EndpointGetChannelMessages(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + messages
}

// EndpointGetChannelMessage builds a query for an HTTP request.
func EndpointGetChannelMessage(channelid, messageid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + messages + slash + messageid
}

// EndpointCreateMessage builds a query for an HTTP request.
func EndpointCreateMessage(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + messages
}

// EndpointCrosspostMessage builds a query for an HTTP request.
func EndpointCrosspostMessage(channelid, messageid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + messages + slash + messageid + slash + crosspost
}

// EndpointCreateReaction builds a query for an HTTP request.
func EndpointCreateReaction(channelid, messageid, emoji string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + messages + slash + messageid + slash + reactions + slash + emoji + slash + me
}

// EndpointDeleteOwnReaction builds a query for an HTTP request.
func EndpointDeleteOwnReaction(channelid, messageid, emoji string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + messages + slash + messageid + slash + reactions + slash + emoji + slash + me
}

// EndpointDeleteUserReaction builds a query for an HTTP request.
func EndpointDeleteUserReaction(channelid, messageid, emoji, userid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + messages + slash + messageid + slash + reactions + slash + emoji + slash + userid
}

// EndpointGetReactions builds a query for an HTTP request.
func EndpointGetReactions(channelid, messageid, emoji string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + messages + slash + messageid + slash + reactions + slash + emoji
}

// EndpointDeleteAllReactions builds a query for an HTTP request.
func EndpointDeleteAllReactions(channelid, messageid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + messages + slash + messageid + slash + reactions
}

// EndpointDeleteAllReactionsforEmoji builds a query for an HTTP request.
func EndpointDeleteAllReactionsforEmoji(channelid, messageid, emoji string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + messages + slash + messageid + slash + reactions + slash + emoji
}

// EndpointEditMessage builds a query for an HTTP request.
func EndpointEditMessage(channelid, messageid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + messages + slash + messageid
}

// EndpointDeleteMessage builds a query for an HTTP request.
func EndpointDeleteMessage(channelid, messageid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + messages + slash + messageid
}

// EndpointBulkDeleteMessages builds a query for an HTTP request.
func EndpointBulkDeleteMessages(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + messages + slash + bulkdelete
}

// EndpointEditChannelPermissions builds a query for an HTTP request.
func EndpointEditChannelPermissions(channelid, overwriteid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + permissions + slash + overwriteid
}

// EndpointGetChannelInvites builds a query for an HTTP request.
func EndpointGetChannelInvites(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + invites
}

// EndpointCreateChannelInvite builds a query for an HTTP request.
func EndpointCreateChannelInvite(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + invites
}

// EndpointDeleteChannelPermission builds a query for an HTTP request.
func EndpointDeleteChannelPermission(channelid, overwriteid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + permissions + slash + overwriteid
}

// EndpointFollowNewsChannel builds a query for an HTTP request.
func EndpointFollowNewsChannel(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + followers
}

// EndpointTriggerTypingIndicator builds a query for an HTTP request.
func EndpointTriggerTypingIndicator(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + typing
}

// EndpointGetPinnedMessages builds a query for an HTTP request.
func EndpointGetPinnedMessages(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + pins
}

// EndpointPinMessage builds a query for an HTTP request.
func EndpointPinMessage(channelid, messageid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + pins + slash + messageid
}

// EndpointUnpinMessage builds a query for an HTTP request.
func EndpointUnpinMessage(channelid, messageid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + pins + slash + messageid
}

// EndpointGroupDMAddRecipient builds a query for an HTTP request.
func EndpointGroupDMAddRecipient(channelid, userid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + recipients + slash + userid
}

// EndpointGroupDMRemoveRecipient builds a query for an HTTP request.
func EndpointGroupDMRemoveRecipient(channelid, userid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + recipients + slash + userid
}

// EndpointStartThreadfromMessage builds a query for an HTTP request.
func EndpointStartThreadfromMessage(channelid, messageid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + messages + slash + messageid + slash + threads
}

// EndpointStartThreadwithoutMessage builds a query for an HTTP request.
func EndpointStartThreadwithoutMessage(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + threads
}

// EndpointStartThreadinForumChannel builds a query for an HTTP request.
func EndpointStartThreadinForumChannel(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + threads
}

// EndpointJoinThread builds a query for an HTTP request.
func EndpointJoinThread(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + threadmembers + slash + me
}

// EndpointAddThreadMember builds a query for an HTTP request.
func EndpointAddThreadMember(channelid, userid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + threadmembers + slash + userid
}

// EndpointLeaveThread builds a query for an HTTP request.
func EndpointLeaveThread(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + threadmembers + slash + me
}

// EndpointRemoveThreadMember builds a query for an HTTP request.
func EndpointRemoveThreadMember(channelid, userid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + threadmembers + slash + userid
}

// EndpointGetThreadMember builds a query for an HTTP request.
func EndpointGetThreadMember(channelid, userid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + threadmembers + slash + userid
}

// EndpointListThreadMembers builds a query for an HTTP request.
func EndpointListThreadMembers(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + threadmembers
}

// EndpointListActiveThreads builds a query for an HTTP request.
func EndpointListActiveThreads(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + threads + slash + active
}

// EndpointListPublicArchivedThreads builds a query for an HTTP request.
func EndpointListPublicArchivedThreads(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + threads + slash + archived + slash + public
}

// EndpointListPrivateArchivedThreads builds a query for an HTTP request.
func EndpointListPrivateArchivedThreads(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + threads + slash + archived + slash + private
}

// EndpointListJoinedPrivateArchivedThreads builds a query for an HTTP request.
func EndpointListJoinedPrivateArchivedThreads(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + users + slash + me + slash + threads + slash + archived + slash + private
}

// EndpointListGuildEmojis builds a query for an HTTP request.
func EndpointListGuildEmojis(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + emojis
}

// EndpointGetGuildEmoji builds a query for an HTTP request.
func EndpointGetGuildEmoji(guildid, emojiid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + emojis + slash + emojiid
}

// EndpointCreateGuildEmoji builds a query for an HTTP request.
func EndpointCreateGuildEmoji(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + emojis
}

// EndpointModifyGuildEmoji builds a query for an HTTP request.
func EndpointModifyGuildEmoji(guildid, emojiid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + emojis + slash + emojiid
}

// EndpointDeleteGuildEmoji builds a query for an HTTP request.
func EndpointDeleteGuildEmoji(guildid, emojiid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + emojis + slash + emojiid
}

// EndpointListScheduledEventsforGuild builds a query for an HTTP request.
func EndpointListScheduledEventsforGuild(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + scheduledevents
}

// EndpointCreateGuildScheduledEvent builds a query for an HTTP request.
func EndpointCreateGuildScheduledEvent(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + scheduledevents
}

// EndpointGetGuildScheduledEvent builds a query for an HTTP request.
func EndpointGetGuildScheduledEvent(guildid, guildscheduledeventid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + scheduledevents + slash + guildscheduledeventid
}

// EndpointModifyGuildScheduledEvent builds a query for an HTTP request.
func EndpointModifyGuildScheduledEvent(guildid, guildscheduledeventid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + scheduledevents + slash + guildscheduledeventid
}

// EndpointDeleteGuildScheduledEvent builds a query for an HTTP request.
func EndpointDeleteGuildScheduledEvent(guildid, guildscheduledeventid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + scheduledevents + slash + guildscheduledeventid
}

// EndpointGetGuildScheduledEventUsers builds a query for an HTTP request.
func EndpointGetGuildScheduledEventUsers(guildid, guildscheduledeventid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + scheduledevents + slash + guildscheduledeventid + slash + users
}

// EndpointGetGuildTemplate builds a query for an HTTP request.
func EndpointGetGuildTemplate(templatecode string) string {
	return EndpointBaseURL + guilds + slash + templates + slash + templatecode
}

// EndpointCreateGuildfromGuildTemplate builds a query for an HTTP request.
func EndpointCreateGuildfromGuildTemplate(templatecode string) string {
	return EndpointBaseURL + guilds + slash + templates + slash + templatecode
}

// EndpointGetGuildTemplates builds a query for an HTTP request.
func EndpointGetGuildTemplates(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + templates
}

// EndpointCreateGuildTemplate builds a query for an HTTP request.
func EndpointCreateGuildTemplate(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + templates
}

// EndpointSyncGuildTemplate builds a query for an HTTP request.
func EndpointSyncGuildTemplate(guildid, templatecode string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + templates + slash + templatecode
}

// EndpointModifyGuildTemplate builds a query for an HTTP request.
func EndpointModifyGuildTemplate(guildid, templatecode string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + templates + slash + templatecode
}

// EndpointDeleteGuildTemplate builds a query for an HTTP request.
func EndpointDeleteGuildTemplate(guildid, templatecode string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + templates + slash + templatecode
}

// EndpointCreateGuild builds a query for an HTTP request.
func EndpointCreateGuild() string {
	return EndpointBaseURL + guilds
}

// EndpointGetGuild builds a query for an HTTP request.
func EndpointGetGuild(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid
}

// EndpointGetGuildPreview builds a query for an HTTP request.
func EndpointGetGuildPreview(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + preview
}

// EndpointModifyGuild builds a query for an HTTP request.
func EndpointModifyGuild(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid
}

// EndpointDeleteGuild builds a query for an HTTP request.
func EndpointDeleteGuild(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid
}

// EndpointGetGuildChannels builds a query for an HTTP request.
func EndpointGetGuildChannels(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + channels
}

// EndpointCreateGuildChannel builds a query for an HTTP request.
func EndpointCreateGuildChannel(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + channels
}

// EndpointModifyGuildChannelPositions builds a query for an HTTP request.
func EndpointModifyGuildChannelPositions(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + channels
}

// EndpointGetGuildMember builds a query for an HTTP request.
func EndpointGetGuildMember(guildid, userid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + members + slash + userid
}

// EndpointListGuildMembers builds a query for an HTTP request.
func EndpointListGuildMembers(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + members
}

// EndpointSearchGuildMembers builds a query for an HTTP request.
func EndpointSearchGuildMembers(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + members + slash + search
}

// EndpointAddGuildMember builds a query for an HTTP request.
func EndpointAddGuildMember(guildid, userid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + members + slash + userid
}

// EndpointModifyGuildMember builds a query for an HTTP request.
func EndpointModifyGuildMember(guildid, userid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + members + slash + userid
}

// EndpointModifyCurrentMember builds a query for an HTTP request.
func EndpointModifyCurrentMember(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + members + slash + me
}

// EndpointModifyCurrentUserNick builds a query for an HTTP request.
func EndpointModifyCurrentUserNick(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + members + slash + me + slash + nick
}

// EndpointAddGuildMemberRole builds a query for an HTTP request.
func EndpointAddGuildMemberRole(guildid, userid, roleid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + members + slash + userid + slash + roles + slash + roleid
}

// EndpointRemoveGuildMemberRole builds a query for an HTTP request.
func EndpointRemoveGuildMemberRole(guildid, userid, roleid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + members + slash + userid + slash + roles + slash + roleid
}

// EndpointRemoveGuildMember builds a query for an HTTP request.
func EndpointRemoveGuildMember(guildid, userid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + members + slash + userid
}

// EndpointGetGuildBans builds a query for an HTTP request.
func EndpointGetGuildBans(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + bans
}

// EndpointGetGuildBan builds a query for an HTTP request.
func EndpointGetGuildBan(guildid, userid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + bans + slash + userid
}

// EndpointCreateGuildBan builds a query for an HTTP request.
func EndpointCreateGuildBan(guildid, userid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + bans + slash + userid
}

// EndpointRemoveGuildBan builds a query for an HTTP request.
func EndpointRemoveGuildBan(guildid, userid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + bans + slash + userid
}

// EndpointGetGuildRoles builds a query for an HTTP request.
func EndpointGetGuildRoles(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + roles
}

// EndpointCreateGuildRole builds a query for an HTTP request.
func EndpointCreateGuildRole(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + roles
}

// EndpointModifyGuildRolePositions builds a query for an HTTP request.
func EndpointModifyGuildRolePositions(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + roles
}

// EndpointModifyGuildRole builds a query for an HTTP request.
func EndpointModifyGuildRole(guildid, roleid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + roles + slash + roleid
}

// EndpointDeleteGuildRole builds a query for an HTTP request.
func EndpointDeleteGuildRole(guildid, roleid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + roles + slash + roleid
}

// EndpointGetGuildPruneCount builds a query for an HTTP request.
func EndpointGetGuildPruneCount(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + prune
}

// EndpointBeginGuildPrune builds a query for an HTTP request.
func EndpointBeginGuildPrune(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + prune
}

// EndpointGetGuildVoiceRegions builds a query for an HTTP request.
func EndpointGetGuildVoiceRegions(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + regions
}

// EndpointGetGuildInvites builds a query for an HTTP request.
func EndpointGetGuildInvites(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + invites
}

// EndpointGetGuildIntegrations builds a query for an HTTP request.
func EndpointGetGuildIntegrations(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + integrations
}

// EndpointDeleteGuildIntegration builds a query for an HTTP request.
func EndpointDeleteGuildIntegration(guildid, integrationid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + integrations + slash + integrationid
}

// EndpointGetGuildWidgetSettings builds a query for an HTTP request.
func EndpointGetGuildWidgetSettings(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + widget
}

// EndpointModifyGuildWidget builds a query for an HTTP request.
func EndpointModifyGuildWidget(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + widget
}

// EndpointGetGuildWidget builds a query for an HTTP request.
func EndpointGetGuildWidget(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + widgetjson
}

// EndpointGetGuildVanityURL builds a query for an HTTP request.
func EndpointGetGuildVanityURL(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + vanityurl
}

// EndpointGetGuildWidgetImage builds a query for an HTTP request.
func EndpointGetGuildWidgetImage(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + widgetpng
}

// EndpointGetGuildWelcomeScreen builds a query for an HTTP request.
func EndpointGetGuildWelcomeScreen(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + welcomescreen
}

// EndpointModifyGuildWelcomeScreen builds a query for an HTTP request.
func EndpointModifyGuildWelcomeScreen(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + welcomescreen
}

// EndpointModifyCurrentUserVoiceState builds a query for an HTTP request.
func EndpointModifyCurrentUserVoiceState(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + voicestates + slash + me
}

// EndpointModifyUserVoiceState builds a query for an HTTP request.
func EndpointModifyUserVoiceState(guildid, userid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + voicestates + slash + userid
}

// EndpointGetInvite builds a query for an HTTP request.
func EndpointGetInvite(invitecode string) string {
	return EndpointBaseURL + invites + slash + invitecode
}

// EndpointDeleteInvite builds a query for an HTTP request.
func EndpointDeleteInvite(invitecode string) string {
	return EndpointBaseURL + invites + slash + invitecode
}

// EndpointCreateStageInstance builds a query for an HTTP request.
func EndpointCreateStageInstance() string {
	return EndpointBaseURL + stageinstances
}

// EndpointGetStageInstance builds a query for an HTTP request.
func EndpointGetStageInstance(channelid string) string {
	return EndpointBaseURL + stageinstances + slash + channelid
}

// EndpointModifyStageInstance builds a query for an HTTP request.
func EndpointModifyStageInstance(channelid string) string {
	return EndpointBaseURL + stageinstances + slash + channelid
}

// EndpointDeleteStageInstance builds a query for an HTTP request.
func EndpointDeleteStageInstance(channelid string) string {
	return EndpointBaseURL + stageinstances + slash + channelid
}

// EndpointGetSticker builds a query for an HTTP request.
func EndpointGetSticker(stickerid string) string {
	return EndpointBaseURL + stickers + slash + stickerid
}

// EndpointListNitroStickerPacks builds a query for an HTTP request.
func EndpointListNitroStickerPacks() string {
	return EndpointBaseURL + stickerpacks
}

// EndpointListGuildStickers builds a query for an HTTP request.
func EndpointListGuildStickers(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + stickers
}

// EndpointGetGuildSticker builds a query for an HTTP request.
func EndpointGetGuildSticker(guildid, stickerid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + stickers + slash + stickerid
}

// EndpointCreateGuildSticker builds a query for an HTTP request.
func EndpointCreateGuildSticker(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + stickers
}

// EndpointModifyGuildSticker builds a query for an HTTP request.
func EndpointModifyGuildSticker(guildid, stickerid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + stickers + slash + stickerid
}

// EndpointDeleteGuildSticker builds a query for an HTTP request.
func EndpointDeleteGuildSticker(guildid, stickerid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + stickers + slash + stickerid
}

// EndpointGetCurrentUser builds a query for an HTTP request.
func EndpointGetCurrentUser() string {
	return EndpointBaseURL + users + slash + me
}

// EndpointGetUser builds a query for an HTTP request.
func EndpointGetUser(userid string) string {
	return EndpointBaseURL + users + slash + userid
}

// EndpointModifyCurrentUser builds a query for an HTTP request.
func EndpointModifyCurrentUser() string {
	return EndpointBaseURL + users + slash + me
}

// EndpointGetCurrentUserGuilds builds a query for an HTTP request.
func EndpointGetCurrentUserGuilds() string {
	return EndpointBaseURL + users + slash + me + slash + guilds
}

// EndpointGetCurrentUserGuildMember builds a query for an HTTP request.
func EndpointGetCurrentUserGuildMember(guildid string) string {
	return EndpointBaseURL + users + slash + me + slash + guilds + slash + guildid + slash + member
}

// EndpointLeaveGuild builds a query for an HTTP request.
func EndpointLeaveGuild(guildid string) string {
	return EndpointBaseURL + users + slash + me + slash + guilds + slash + guildid
}

// EndpointCreateDM builds a query for an HTTP request.
func EndpointCreateDM() string {
	return EndpointBaseURL + users + slash + me + slash + channels
}

// EndpointCreateGroupDM builds a query for an HTTP request.
func EndpointCreateGroupDM() string {
	return EndpointBaseURL + users + slash + me + slash + channels
}

// EndpointGetUserConnections builds a query for an HTTP request.
func EndpointGetUserConnections() string {
	return EndpointBaseURL + users + slash + me + slash + connections
}

// EndpointListVoiceRegions builds a query for an HTTP request.
func EndpointListVoiceRegions() string {
	return EndpointBaseURL + voice + slash + regions
}

// EndpointCreateWebhook builds a query for an HTTP request.
func EndpointCreateWebhook(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + webhooks
}

// EndpointGetChannelWebhooks builds a query for an HTTP request.
func EndpointGetChannelWebhooks(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + webhooks
}

// EndpointGetGuildWebhooks builds a query for an HTTP request.
func EndpointGetGuildWebhooks(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + webhooks
}

// EndpointGetWebhook builds a query for an HTTP request.
func EndpointGetWebhook(webhookid string) string {
	return EndpointBaseURL + webhooks + slash + webhookid
}

// EndpointGetWebhookwithToken builds a query for an HTTP request.
func EndpointGetWebhookwithToken(webhookid, webhooktoken string) string {
	return EndpointBaseURL + webhooks + slash + webhookid + slash + webhooktoken
}

// EndpointModifyWebhook builds a query for an HTTP request.
func EndpointModifyWebhook(webhookid string) string {
	return EndpointBaseURL + webhooks + slash + webhookid
}

// EndpointModifyWebhookwithToken builds a query for an HTTP request.
func EndpointModifyWebhookwithToken(webhookid, webhooktoken string) string {
	return EndpointBaseURL + webhooks + slash + webhookid + slash + webhooktoken
}

// EndpointDeleteWebhook builds a query for an HTTP request.
func EndpointDeleteWebhook(webhookid string) string {
	return EndpointBaseURL + webhooks + slash + webhookid
}

// EndpointDeleteWebhookwithToken builds a query for an HTTP request.
func EndpointDeleteWebhookwithToken(webhookid, webhooktoken string) string {
	return EndpointBaseURL + webhooks + slash + webhookid + slash + webhooktoken
}

// EndpointExecuteWebhook builds a query for an HTTP request.
func EndpointExecuteWebhook(webhookid, webhooktoken string) string {
	return EndpointBaseURL + webhooks + slash + webhookid + slash + webhooktoken
}

// EndpointExecuteSlackCompatibleWebhook builds a query for an HTTP request.
func EndpointExecuteSlackCompatibleWebhook(webhookid, webhooktoken string) string {
	return EndpointBaseURL + webhooks + slash + webhookid + slash + webhooktoken + slash + slack
}

// EndpointExecuteGitHubCompatibleWebhook builds a query for an HTTP request.
func EndpointExecuteGitHubCompatibleWebhook(webhookid, webhooktoken string) string {
	return EndpointBaseURL + webhooks + slash + webhookid + slash + webhooktoken + slash + github
}

// EndpointGetWebhookMessage builds a query for an HTTP request.
func EndpointGetWebhookMessage(webhookid, webhooktoken, messageid string) string {
	return EndpointBaseURL + webhooks + slash + webhookid + slash + webhooktoken + slash + messages + slash + messageid
}

// EndpointEditWebhookMessage builds a query for an HTTP request.
func EndpointEditWebhookMessage(webhookid, webhooktoken, messageid string) string {
	return EndpointBaseURL + webhooks + slash + webhookid + slash + webhooktoken + slash + messages + slash + messageid
}

// EndpointDeleteWebhookMessage builds a query for an HTTP request.
func EndpointDeleteWebhookMessage(webhookid, webhooktoken, messageid string) string {
	return EndpointBaseURL + webhooks + slash + webhookid + slash + webhooktoken + slash + messages + slash + messageid
}

// EndpointGetGateway builds a query for an HTTP request.
func EndpointGetGateway() string {
	return EndpointBaseURL + gateway
}

// EndpointGetGatewayBot builds a query for an HTTP request.
func EndpointGetGatewayBot() string {
	return EndpointBaseURL + gateway + slash + bot
}

// EndpointGetCurrentBotApplicationInformation builds a query for an HTTP request.
func EndpointGetCurrentBotApplicationInformation() string {
	return EndpointBaseURL + oauth + slash + applications + slash + me
}

// EndpointGetCurrentAuthorizationInformation builds a query for an HTTP request.
func EndpointGetCurrentAuthorizationInformation() string {
	return EndpointBaseURL + oauth + slash + me
}
