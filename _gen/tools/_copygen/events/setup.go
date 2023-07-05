package wrapper

import (
	disgo "github.com/switchupcb/disgo/wrapper"
)

// Copygen defines the functions that will be generated.
type Copygen interface {
	Hello(*disgo.Hello)
	Ready(*disgo.Ready)
	Resumed(*disgo.Resumed)
	Reconnect(*disgo.Reconnect)
	InvalidSession(*disgo.InvalidSession)
	ApplicationCommandPermissionsUpdate(*disgo.ApplicationCommandPermissionsUpdate)
	// intents FlagIntentAUTO_MODERATION_CONFIGURATION
	AutoModerationRuleCreate(*disgo.AutoModerationRuleCreate)
	// intents FlagIntentAUTO_MODERATION_CONFIGURATION
	AutoModerationRuleUpdate(*disgo.AutoModerationRuleUpdate)
	// intents FlagIntentAUTO_MODERATION_CONFIGURATION
	AutoModerationRuleDelete(*disgo.AutoModerationRuleDelete)
	// intents FlagIntentAUTO_MODERATION_EXECUTION
	AutoModerationActionExecution(*disgo.AutoModerationActionExecution)
	InteractionCreate(*disgo.InteractionCreate)
	VoiceServerUpdate(*disgo.VoiceServerUpdate)
	GuildMembersChunk(*disgo.GuildMembersChunk)
	UserUpdate(*disgo.UserUpdate)
	// intents FlagIntentGUILDS
	ChannelCreate(*disgo.ChannelCreate)
	// intents FlagIntentGUILDS
	ChannelUpdate(*disgo.ChannelUpdate)
	// intents FlagIntentGUILDS
	ChannelDelete(*disgo.ChannelDelete)
	// intents FlagIntentGUILDS FlagIntentDIRECT_MESSAGES
	ChannelPinsUpdate(*disgo.ChannelPinsUpdate)
	// intents FlagIntentGUILDS
	ThreadCreate(*disgo.ThreadCreate)
	// intents FlagIntentGUILDS
	ThreadUpdate(*disgo.ThreadUpdate)
	// intents FlagIntentGUILDS
	ThreadDelete(*disgo.ThreadDelete)
	// intents FlagIntentGUILDS
	ThreadListSync(*disgo.ThreadListSync)
	// intents FlagIntentGUILDS
	ThreadMemberUpdate(*disgo.ThreadMemberUpdate)
	// intents FlagIntentGUILDS FlagIntentGUILD_MEMBERS
	ThreadMembersUpdate(*disgo.ThreadMembersUpdate)
	// intents FlagIntentGUILDS
	GuildCreate(*disgo.GuildCreate)
	// intents FlagIntentGUILDS
	GuildUpdate(*disgo.GuildUpdate)
	// intents FlagIntentGUILDS
	GuildDelete(*disgo.GuildDelete)
	// intents FlagIntentGUILD_MODERATION
	GuildAuditLogEntryCreate(*disgo.GuildAuditLogEntryCreate)
	// intents FlagIntentGUILD_MODERATION
	GuildBanAdd(*disgo.GuildBanAdd)
	// intents FlagIntentGUILD_MODERATION
	GuildBanRemove(*disgo.GuildBanRemove)
	// intents FlagIntentGUILD_EMOJIS_AND_STICKERS
	GuildEmojisUpdate(*disgo.GuildEmojisUpdate)
	// intents FlagIntentGUILD_EMOJIS_AND_STICKERS
	GuildStickersUpdate(*disgo.GuildStickersUpdate)
	// intents FlagIntentGUILD_INTEGRATIONS
	GuildIntegrationsUpdate(*disgo.GuildIntegrationsUpdate)
	// intents FlagIntentGUILD_MEMBERS
	GuildMemberAdd(*disgo.GuildMemberAdd)
	// intents FlagIntentGUILD_MEMBERS
	GuildMemberRemove(*disgo.GuildMemberRemove)
	// intents FlagIntentGUILD_MEMBERS
	GuildMemberUpdate(*disgo.GuildMemberUpdate)
	// intents FlagIntentGUILDS
	GuildRoleCreate(*disgo.GuildRoleCreate)
	// intents FlagIntentGUILDS
	GuildRoleUpdate(*disgo.GuildRoleUpdate)
	// intents FlagIntentGUILDS
	GuildRoleDelete(*disgo.GuildRoleDelete)
	// intents FlagIntentGUILD_SCHEDULED_EVENTS
	GuildScheduledEventCreate(*disgo.GuildScheduledEventCreate)
	// intents FlagIntentGUILD_SCHEDULED_EVENTS
	GuildScheduledEventUpdate(*disgo.GuildScheduledEventUpdate)
	// intents FlagIntentGUILD_SCHEDULED_EVENTS
	GuildScheduledEventDelete(*disgo.GuildScheduledEventDelete)
	// intents FlagIntentGUILD_SCHEDULED_EVENTS
	GuildScheduledEventUserAdd(*disgo.GuildScheduledEventUserAdd)
	// intents FlagIntentGUILD_SCHEDULED_EVENTS
	GuildScheduledEventUserRemove(*disgo.GuildScheduledEventUserRemove)
	// intents FlagIntentGUILD_INTEGRATIONS
	IntegrationCreate(*disgo.IntegrationCreate)
	// intents FlagIntentGUILD_INTEGRATIONS
	IntegrationUpdate(*disgo.IntegrationUpdate)
	// intents FlagIntentGUILD_INTEGRATIONS
	IntegrationDelete(*disgo.IntegrationDelete)
	// intents FlagIntentGUILD_INVITES
	InviteCreate(*disgo.InviteCreate)
	// intents FlagIntentGUILD_INVITES
	InviteDelete(*disgo.InviteDelete)
	// intents FlagIntentGUILD_MESSAGES FlagIntentDIRECT_MESSAGES
	MessageCreate(*disgo.MessageCreate)
	// intents FlagIntentGUILD_MESSAGES FlagIntentDIRECT_MESSAGES
	MessageUpdate(*disgo.MessageUpdate)
	// intents FlagIntentGUILD_MESSAGES FlagIntentDIRECT_MESSAGES
	MessageDelete(*disgo.MessageDelete)
	// intents FlagIntentGUILD_MESSAGES
	MessageDeleteBulk(*disgo.MessageDeleteBulk)
	// intents FlagIntentGUILD_MESSAGE_REACTIONS FlagIntentDIRECT_MESSAGE_REACTIONS
	MessageReactionAdd(*disgo.MessageReactionAdd)
	// intents FlagIntentGUILD_MESSAGE_REACTIONS FlagIntentDIRECT_MESSAGE_REACTIONS
	MessageReactionRemove(*disgo.MessageReactionRemove)
	// intents FlagIntentGUILD_MESSAGE_REACTIONS FlagIntentDIRECT_MESSAGE_REACTIONS
	MessageReactionRemoveAll(*disgo.MessageReactionRemoveAll)
	// intents FlagIntentGUILD_MESSAGE_REACTIONS FlagIntentDIRECT_MESSAGE_REACTIONS
	MessageReactionRemoveEmoji(*disgo.MessageReactionRemoveEmoji)
	// intents FlagIntentGUILD_PRESENCES
	PresenceUpdate(*disgo.PresenceUpdate)
	// intents FlagIntentGUILDS
	StageInstanceCreate(*disgo.StageInstanceCreate)
	// intents FlagIntentGUILDS
	StageInstanceDelete(*disgo.StageInstanceDelete)
	// intents FlagIntentGUILDS
	StageInstanceUpdate(*disgo.StageInstanceUpdate)
	// intents FlagIntentGUILD_MESSAGE_REACTIONS FlagIntentDIRECT_MESSAGE_TYPING
	TypingStart(*disgo.TypingStart)
	// intents FlagIntentGUILD_VOICE_STATES
	VoiceStateUpdate(*disgo.VoiceStateUpdate)
	// intents FlagIntentGUILD_WEBHOOKS
	WebhooksUpdate(*disgo.WebhooksUpdate)
}
