package requests

// Discord API Endpoints
const (
	EndpointBaseURL = "https://discord.com/api/v9/"
	slash           = "/"

	// slugs
	applications = "applications"
	interactions = "interactions"
	webhooks     = "webhooks"
	guilds       = "guilds"
	channels     = "channels"
	invites      = "invites"
	stickers     = "stickers"
	users        = "users"

	// resources
	commands        = "commands"
	permissions     = "permissions"
	callback        = "callback"
	messages        = "messages"
	auditlogs       = "audit-logs"
	crosspost       = "crosspost"
	bulkdelete      = "bulk-delete"
	reactions       = "reactions"
	emoji           = "emoji"
	emojis          = "emojis"
	followers       = "followers"
	typing          = "typing"
	pins            = "pins"
	recipients      = "recipients"
	threads         = "threads"
	members         = "members"
	threadmembers   = "thread-members"
	active          = "active"
	archived        = "archived"
	public          = "public"
	private         = "private"
	scheduledevents = "scheduled-events"
	templates       = "templates"
	preview         = "preview"
	search          = "search"
	nick            = "nick"
	roles           = "roles"
	bans            = "bans"
	prune           = "prune"
	regions         = "regions"
	integrations    = "integrations"
	widget          = "widget"
	widgetjson      = "widget.json"
	widgetpng       = "widget.png"
	vanityurl       = "vanity-url"
	welcomescreen   = "welcome-screen"
	voicestates     = "voice-states"
	member          = "member"
	slack           = "slack"
	github          = "github"

	// parameters
	original = "@original"
	me       = "@me"
)

func EndpointGetGlobalApplicationCommands(applicationid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + commands
}

func EndpointCreateGlobalApplicationCommand(applicationid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + commands
}

func EndpointGetGlobalApplicationCommand(applicationid, commandid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + commands + slash + commandid
}

func EndpointEditGlobalApplicationCommand(applicationid, commandid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + commands + slash + commandid
}

func EndpointDeleteGlobalApplicationCommand(applicationid, commandid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + commands + slash + commandid
}

func EndpointBulkOverwriteGlobalApplicationCommands(applicationid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + commands
}

func EndpointGetGuildApplicationCommands(applicationid, guildid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + guilds + slash + guildid + slash + commands
}

func EndpointCreateGuildApplicationCommand(applicationid, guildid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + guilds + slash + guildid + slash + commands
}

func EndpointGetGuildApplicationCommand(applicationid, guildid, commandid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + guilds + slash + guildid + slash + commands + slash + commandid
}

func EndpointEditGuildApplicationCommand(applicationid, guildid, commandid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + guilds + slash + guildid + slash + commands + slash + commandid
}

func EndpointDeleteGuildApplicationCommand(applicationid, guildid, commandid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + guilds + slash + guildid + slash + commands + slash + commandid
}

func EndpointBulkOverwriteGuildApplicationCommands(applicationid, guildid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + guilds + slash + guildid + slash + commands
}

func EndpointGetGuildApplicationCommandPermissions(applicationid, guildid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + guilds + slash + guildid + slash + commands + slash + permissions
}

func EndpointGetApplicationCommandPermissions(applicationid, guildid, commandid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + guilds + slash + guildid + slash + commands + slash + commandid + slash + permissions
}

func EndpointEditApplicationCommandPermissions(applicationid, guildid, commandid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + guilds + slash + guildid + slash + commands + slash + commandid + slash + permissions
}

func EndpointBatchEditApplicationCommandPermissions(applicationid, guildid string) string {
	return EndpointBaseURL + applications + slash + applicationid + slash + guilds + slash + guildid + slash + commands + slash + permissions
}

func EndpointCreateInteractionResponse(interactionid, interactiontoken string) string {
	return EndpointBaseURL + interactions + slash + interactionid + slash + interactiontoken + slash + callback
}

func EndpointGetOriginalInteractionResponse(applicationid, interactiontoken string) string {
	return EndpointBaseURL + webhooks + slash + applicationid + slash + interactiontoken + slash + messages + slash + original
}

func EndpointEditOriginalInteractionResponse(applicationid, interactiontoken string) string {
	return EndpointBaseURL + webhooks + slash + applicationid + slash + interactiontoken + slash + messages + slash + original
}

func EndpointDeleteOriginalInteractionResponse(applicationid, interactiontoken string) string {
	return EndpointBaseURL + webhooks + slash + applicationid + slash + interactiontoken + slash + messages + slash + original
}

func EndpointCreateFollowupMessage(applicationid, interactiontoken string) string {
	return EndpointBaseURL + webhooks + slash + applicationid + slash + interactiontoken
}

func EndpointGetFollowupMessage(applicationid, interactiontoken, messageid string) string {
	return EndpointBaseURL + webhooks + slash + applicationid + slash + interactiontoken + slash + messages + slash + messageid
}

func EndpointEditFollowupMessage(applicationid, interactiontoken, messageid string) string {
	return EndpointBaseURL + webhooks + slash + applicationid + slash + interactiontoken + slash + messages + slash + messageid
}

func EndpointDeleteFollowupMessage(applicationid, interactiontoken, messageid string) string {
	return EndpointBaseURL + webhooks + slash + applicationid + slash + interactiontoken + slash + messages + slash + messageid
}

func EndpointGetGuildAuditLog(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + auditlogs
}

func EndpointGetChannel(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid
}

func EndpointModifyChannel(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid
}

func EndpointDeleteCloseChannel(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid
}

func EndpointGetChannelMessages(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + messages
}

func EndpointGetChannelMessage(channelid, messageid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + messages + slash + messageid
}

func EndpointCreateMessage(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + messages
}

func EndpointCrosspostMessage(channelid, messageid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + messages + slash + messageid + slash + crosspost
}

func EndpointCreateReaction(channelid, messageid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + messages + slash + messageid + slash + reactions + slash + emoji + slash + me
}

func EndpointDeleteOwnReaction(channelid, messageid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + messages + slash + messageid + slash + reactions + slash + emoji + slash + me
}

func EndpointDeleteUserReaction(channelid, messageid, userid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + messages + slash + messageid + slash + reactions + slash + emoji + slash + userid
}

func EndpointGetReactions(channelid, messageid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + messages + slash + messageid + slash + reactions + slash + emoji
}

func EndpointDeleteAllReactions(channelid, messageid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + messages + slash + messageid + slash + reactions
}

func EndpointDeleteAllReactionsforEmoji(channelid, messageid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + messages + slash + messageid + slash + reactions + slash + emoji
}

func EndpointEditMessage(channelid, messageid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + messages + slash + messageid
}

func EndpointDeleteMessage(channelid, messageid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + messages + slash + messageid
}

func EndpointBulkDeleteMessages(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + messages + slash + bulkdelete
}

func EndpointEditChannelPermissions(channelid, overwriteid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + permissions + slash + overwriteid
}

func EndpointGetChannelInvites(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + invites
}

func EndpointCreateChannelInvite(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + invites
}

func EndpointDeleteChannelPermission(channelid, overwriteid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + permissions + slash + overwriteid
}

func EndpointFollowNewsChannel(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + followers
}

func EndpointTriggerTypingIndicator(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + typing
}

func EndpointGetPinnedMessages(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + pins
}

func EndpointPinMessage(channelid, messageid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + pins + slash + messageid
}

func EndpointUnpinMessage(channelid, messageid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + pins + slash + messageid
}

func EndpointGroupDMAddRecipient(channelid, userid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + recipients + slash + userid
}

func EndpointGroupDMRemoveRecipient(channelid, userid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + recipients + slash + userid
}

func EndpointStartThreadfromMessage(channelid, messageid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + messages + slash + messageid + slash + threads
}

func EndpointStartThreadwithoutMessage(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + threads
}

func EndpointStartThreadinForumChannel(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + threads
}

func EndpointJoinThread(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + threadmembers + slash + me
}

func EndpointAddThreadMember(channelid, userid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + threadmembers + slash + userid
}

func EndpointLeaveThread(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + threadmembers + slash + me
}

func EndpointRemoveThreadMember(channelid, userid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + threadmembers + slash + userid
}

func EndpointGetThreadMember(channelid, userid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + threadmembers + slash + userid
}

func EndpointListThreadMembers(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + threadmembers
}

func EndpointListActiveChannelThreads(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + threads + slash + active
}

func EndpointListPublicArchivedThreads(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + threads + slash + archived + slash + public
}

func EndpointListPrivateArchivedThreads(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + threads + slash + archived + slash + private
}

func EndpointListJoinedPrivateArchivedThreads(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + users + slash + me + slash + threads + slash + archived + slash + private
}

func EndpointListGuildEmojis(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + emojis
}

func EndpointGetGuildEmoji(guildid, emojiid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + emojis + slash + emojiid
}

func EndpointCreateGuildEmoji(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + emojis
}

func EndpointModifyGuildEmoji(guildid, emojiid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + emojis + slash + emojiid
}

func EndpointDeleteGuildEmoji(guildid, emojiid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + emojis + slash + emojiid
}

func EndpointListScheduledEventsforGuild(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + scheduledevents
}

func EndpointCreateGuildScheduledEvent(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + scheduledevents
}

func EndpointGetGuildScheduledEvent(guildid, guild_scheduled_eventid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + scheduledevents + slash + guild_scheduled_eventid
}

func EndpointModifyGuildScheduledEvent(guildid, guild_scheduled_eventid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + scheduledevents + slash + guild_scheduled_eventid
}

func EndpointDeleteGuildScheduledEvent(guildid, guild_scheduled_eventid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + scheduledevents + slash + guild_scheduled_eventid
}

func EndpointGetGuildScheduledEventUsers(guildid, guild_scheduled_eventid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + scheduledevents + slash + guild_scheduled_eventid + slash + users
}

func EndpointGetGuildTemplate(templatecode string) string {
	return EndpointBaseURL + guilds + slash + templates + slash + templatecode
}

func EndpointCreateGuildfromGuildTemplate(templatecode string) string {
	return EndpointBaseURL + guilds + slash + templates + slash + templatecode
}

func EndpointGetGuildTemplates(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + templates
}

func EndpointCreateGuildTemplate(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + templates
}

func EndpointSyncGuildTemplate(guildid, templatecode string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + templates + slash + templatecode
}

func EndpointModifyGuildTemplate(guildid, templatecode string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + templates + slash + templatecode
}

func EndpointDeleteGuildTemplate(guildid, templatecode string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + templates + slash + templatecode
}

func EndpointGetGuild(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid
}

func EndpointGetGuildPreview(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + preview
}

func EndpointModifyGuild(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid
}

func EndpointDeleteGuild(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid
}

func EndpointGetGuildChannels(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + channels
}

func EndpointCreateGuildChannel(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + channels
}

func EndpointModifyGuildChannelPositions(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + channels
}

func EndpointListActiveGuildThreads(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + threads + slash + active
}

func EndpointGetGuildMember(guildid, userid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + members + slash + userid
}

func EndpointListGuildMembers(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + members
}

func EndpointSearchGuildMembers(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + members + slash + search
}

func EndpointAddGuildMember(guildid, userid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + members + slash + userid
}

func EndpointModifyGuildMember(guildid, userid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + members + slash + userid
}

func EndpointModifyCurrentMember(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + members + slash + me
}

func EndpointModifyCurrentUserNick(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + members + slash + me + slash + nick
}

func EndpointAddGuildMemberRole(guildid, userid, roleid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + members + slash + userid + slash + roles + slash + roleid
}

func EndpointRemoveGuildMemberRole(guildid, userid, roleid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + members + slash + userid + slash + roles + slash + roleid
}

func EndpointRemoveGuildMember(guildid, userid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + members + slash + userid
}

func EndpointGetGuildBans(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + bans
}

func EndpointGetGuildBan(guildid, userid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + bans + slash + userid
}

func EndpointCreateGuildBan(guildid, userid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + bans + slash + userid
}

func EndpointRemoveGuildBan(guildid, userid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + bans + slash + userid
}

func EndpointGetGuildRoles(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + roles
}

func EndpointCreateGuildRole(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + roles
}

func EndpointModifyGuildRolePositions(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + roles
}

func EndpointModifyGuildRole(guildid, roleid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + roles + slash + roleid
}

func EndpointDeleteGuildRole(guildid, roleid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + roles + slash + roleid
}

func EndpointGetGuildPruneCount(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + prune
}

func EndpointBeginGuildPrune(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + prune
}

func EndpointGetGuildVoiceRegions(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + regions
}

func EndpointGetGuildInvites(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + invites
}

func EndpointGetGuildIntegrations(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + integrations
}

func EndpointDeleteGuildIntegration(guildid, integrationid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + integrations + slash + integrationid
}

func EndpointGetGuildWidgetSettings(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + widget
}

func EndpointModifyGuildWidget(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + widget
}

func EndpointGetGuildWidget(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + widgetjson
}

func EndpointGetGuildVanityURL(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + vanityurl
}

func EndpointGetGuildWidgetImage(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + widgetpng
}

func EndpointGetGuildWelcomeScreen(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + welcomescreen
}

func EndpointModifyGuildWelcomeScreen(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + welcomescreen
}

func EndpointModifyCurrentUserVoiceState(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + voicestates + slash + me
}

func EndpointModifyUserVoiceState(guildid, userid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + voicestates + slash + userid
}

func EndpointGetInvite(invitecode string) string {
	return EndpointBaseURL + invites + slash + invitecode
}

func EndpointDeleteInvite(invitecode string) string {
	return EndpointBaseURL + invites + slash + invitecode
}

func EndpointGetSticker(stickerid string) string {
	return EndpointBaseURL + stickers + slash + stickerid
}

func EndpointListGuildStickers(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + stickers
}

func EndpointGetGuildSticker(guildid, stickerid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + stickers + slash + stickerid
}

func EndpointCreateGuildSticker(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + stickers
}

func EndpointModifyGuildSticker(guildid, stickerid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + stickers + slash + stickerid
}

func EndpointDeleteGuildSticker(guildid, stickerid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + stickers + slash + stickerid
}

func EndpointGetUser(userid string) string {
	return EndpointBaseURL + users + slash + userid
}

func EndpointGetCurrentUserGuildMember(guildid string) string {
	return EndpointBaseURL + users + slash + me + slash + guilds + slash + guildid + slash + member
}

func EndpointLeaveGuild(guildid string) string {
	return EndpointBaseURL + users + slash + me + slash + guilds + slash + guildid
}

func EndpointCreateWebhook(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + webhooks
}

func EndpointGetChannelWebhooks(channelid string) string {
	return EndpointBaseURL + channels + slash + channelid + slash + webhooks
}

func EndpointGetGuildWebhooks(guildid string) string {
	return EndpointBaseURL + guilds + slash + guildid + slash + webhooks
}

func EndpointGetWebhook(webhookid string) string {
	return EndpointBaseURL + webhooks + slash + webhookid
}

func EndpointGetWebhookwithToken(webhookid, webhooktoken string) string {
	return EndpointBaseURL + webhooks + slash + webhookid + slash + webhooktoken
}

func EndpointModifyWebhook(webhookid string) string {
	return EndpointBaseURL + webhooks + slash + webhookid
}

func EndpointModifyWebhookwithToken(webhookid, webhooktoken string) string {
	return EndpointBaseURL + webhooks + slash + webhookid + slash + webhooktoken
}

func EndpointDeleteWebhook(webhookid string) string {
	return EndpointBaseURL + webhooks + slash + webhookid
}

func EndpointDeleteWebhookwithToken(webhookid, webhooktoken string) string {
	return EndpointBaseURL + webhooks + slash + webhookid + slash + webhooktoken
}

func EndpointExecuteWebhook(webhookid, webhooktoken string) string {
	return EndpointBaseURL + webhooks + slash + webhookid + slash + webhooktoken
}

func EndpointExecuteSlackCompatibleWebhook(webhookid, webhooktoken string) string {
	return EndpointBaseURL + webhooks + slash + webhookid + slash + webhooktoken + slash + slack
}

func EndpointExecuteGitHubCompatibleWebhook(webhookid, webhooktoken string) string {
	return EndpointBaseURL + webhooks + slash + webhookid + slash + webhooktoken + slash + github
}

func EndpointGetWebhookMessage(webhookid, webhooktoken, messageid string) string {
	return EndpointBaseURL + webhooks + slash + webhookid + slash + webhooktoken + slash + messages + slash + messageid
}

func EndpointEditWebhookMessage(webhookid, webhooktoken, messageid string) string {
	return EndpointBaseURL + webhooks + slash + webhookid + slash + webhooktoken + slash + messages + slash + messageid
}

func EndpointDeleteWebhookMessage(webhookid, webhooktoken, messageid string) string {
	return EndpointBaseURL + webhooks + slash + webhookid + slash + webhooktoken + slash + messages + slash + messageid
}
