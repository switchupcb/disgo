package wrapper

import (
	disgo "github.com/switchupcb/disgo/wrapper"
)

// Copygen defines the functions that will be generated.
type Copygen interface {
	CreateGuild(*disgo.CreateGuild) (*disgo.Guild, error)
	CreateGlobalApplicationCommand(*disgo.CreateGlobalApplicationCommand) (*disgo.ApplicationCommand, error)
	CreateGuildApplicationCommand(*disgo.CreateGuildApplicationCommand) (*disgo.ApplicationCommand, error)
	CreateInteractionResponse(*disgo.CreateInteractionResponse) error
	CreateFollowupMessage(*disgo.CreateFollowupMessage) (*disgo.Message, error)
	CreateAutoModerationRule(*disgo.CreateAutoModerationRule) (*disgo.AutoModerationRule, error)
	CreateMessage(*disgo.CreateMessage) (*disgo.Message, error)
	CreateReaction(*disgo.CreateReaction) error
	CreateChannelInvite(*disgo.CreateChannelInvite) (*disgo.Invite, error)
	CreateGuildEmoji(*disgo.CreateGuildEmoji) (*disgo.Emoji, error)
	CreateGuildChannel(*disgo.CreateGuildChannel) (*disgo.Channel, error)
	CreateGuildBan(*disgo.CreateGuildBan) error
	CreateGuildRole(*disgo.CreateGuildRole) (*disgo.Role, error)
	CreateGuildScheduledEvent(*disgo.CreateGuildScheduledEvent) (*disgo.GuildScheduledEvent, error)
	CreateGuildfromGuildTemplate(*disgo.CreateGuildfromGuildTemplate) ([]*disgo.GuildTemplate, error)
	CreateGuildTemplate(*disgo.CreateGuildTemplate) (*disgo.GuildTemplate, error)
	CreateStageInstance(*disgo.CreateStageInstance) (*disgo.StageInstance, error)
	CreateGuildSticker(*disgo.CreateGuildSticker) (*disgo.Sticker, error)
	CreateGroupDM(*disgo.CreateGroupDM) (*disgo.Channel, error)
	CreateWebhook(*disgo.CreateWebhook) (*disgo.Webhook, error)

	GetGlobalApplicationCommands(*disgo.GetGlobalApplicationCommands) ([]*disgo.ApplicationCommand, error)
	GetGlobalApplicationCommand(*disgo.GetGlobalApplicationCommand) (*disgo.ApplicationCommand, error)
	GetGuildApplicationCommands(*disgo.GetGuildApplicationCommands) ([]*disgo.ApplicationCommand, error)
	GetGuildApplicationCommand(*disgo.GetGuildApplicationCommand) (*disgo.ApplicationCommand, error)
	GetGuildApplicationCommandPermissions(*disgo.GetGuildApplicationCommandPermissions) (*disgo.GuildApplicationCommandPermissions, error)
	GetApplicationCommandPermissions(*disgo.GetApplicationCommandPermissions) (*disgo.GuildApplicationCommandPermissions, error)
	GetOriginalInteractionResponse(*disgo.GetOriginalInteractionResponse) error
	GetFollowupMessage(*disgo.GetFollowupMessage) (*disgo.Message, error)
	GetGuildAuditLog(*disgo.GetGuildAuditLog) (*disgo.AuditLog, error)
	GetAutoModerationRule(*disgo.GetAutoModerationRule) (*disgo.AutoModerationRule, error)
	GetChannel(*disgo.GetChannel) (*disgo.Channel, error)
	GetChannelMessages(*disgo.GetChannelMessages) ([]*disgo.Message, error)
	GetChannelMessage(*disgo.GetChannelMessage) (*disgo.Message, error)
	GetReactions(*disgo.GetReactions) ([]*disgo.User, error)
	GetChannelInvites(*disgo.GetChannelInvites) ([]*disgo.Invite, error)
	GetPinnedMessages(*disgo.GetPinnedMessages) ([]*disgo.Message, error)
	GetThreadMember(*disgo.GetThreadMember) (*disgo.ThreadMember, error)
	GetGuildEmoji(*disgo.GetGuildEmoji) (*disgo.Emoji, error)
	GetGuild(*disgo.GetGuild) (*disgo.Guild, error)
	GetGuildPreview(*disgo.GetGuildPreview) (*disgo.GuildPreview, error)
	GetGuildChannels(*disgo.GetGuildChannels) ([]*disgo.Channel, error)
	GetGuildMember(*disgo.GetGuildMember) (*disgo.GuildMember, error)
	GetGuildBans(*disgo.GetGuildBans) ([]*disgo.Ban, error)
	GetGuildBan(*disgo.GetGuildBan) (*disgo.Ban, error)
	GetGuildRoles(*disgo.GetGuildRoles) ([]*disgo.Role, error)
	GetGuildPruneCount(*disgo.GetGuildPruneCount) error
	GetGuildVoiceRegions(*disgo.GetGuildVoiceRegions) (*disgo.VoiceRegion, error)
	GetGuildInvites(*disgo.GetGuildInvites) ([]*disgo.Invite, error)
	GetGuildIntegrations(*disgo.GetGuildIntegrations) ([]*disgo.Integration, error)
	GetGuildWidgetSettings(*disgo.GetGuildWidgetSettings) (*disgo.GuildWidget, error)
	GetGuildWidget(*disgo.GetGuildWidget) (*disgo.GuildWidget, error)
	GetGuildVanityURL(*disgo.GetGuildVanityURL) (*disgo.Invite, error)
	GetGuildWidgetImage(*disgo.GetGuildWidgetImage) (*disgo.EmbedImage, error)
	GetGuildWelcomeScreen(*disgo.GetGuildWelcomeScreen) (*disgo.WelcomeScreen, error)
	GetGuildScheduledEvent(*disgo.GetGuildScheduledEvent) (*disgo.GuildScheduledEvent, error)
	GetGuildScheduledEventUsers(*disgo.GetGuildScheduledEventUsers) ([]*disgo.GuildScheduledEventUser, error)
	GetGuildTemplate(*disgo.GetGuildTemplate) (*disgo.GuildTemplate, error)
	GetGuildTemplates(*disgo.GetGuildTemplates) ([]*disgo.GuildTemplate, error)
	GetInvite(*disgo.GetInvite) (*disgo.Invite, error)
	GetStageInstance(*disgo.GetStageInstance) error
	GetSticker(*disgo.GetSticker) (*disgo.Sticker, error)
	GetGuildSticker(*disgo.GetGuildSticker) (*disgo.Sticker, error)
	GetCurrentUser(*disgo.GetCurrentUser) (*disgo.User, error)
	GetUser(*disgo.GetUser) (*disgo.User, error)
	GetCurrentUserGuilds(*disgo.GetCurrentUserGuilds) ([]*disgo.Guild, error)
	GetCurrentUserGuildMember(*disgo.GetCurrentUserGuildMember) (*disgo.GuildMember, error)
	GetUserConnections(*disgo.GetUserConnections) ([]*disgo.Connection, error)
	GetChannelWebhooks(*disgo.GetChannelWebhooks) ([]*disgo.Webhook, error)
	GetGuildWebhooks(*disgo.GetGuildWebhooks) ([]*disgo.Webhook, error)
	GetWebhook(*disgo.GetWebhook) (*disgo.Webhook, error)
	GetWebhookwithToken(*disgo.GetWebhookwithToken) (*disgo.Webhook, error)
	GetWebhookMessage(*disgo.GetWebhookMessage) (*disgo.Message, error)
	GetGateway(*disgo.GetGateway) (*disgo.GetGatewayBotResponse, error)
	GetGatewayBot(*disgo.GetGatewayBot) (*disgo.GetGatewayBotResponse, error)
	GetCurrentBotApplicationInformation(*disgo.GetCurrentBotApplicationInformation) (*disgo.Application, error)
	GetCurrentAuthorizationInformation(*disgo.GetCurrentAuthorizationInformation) (*disgo.CurrentAuthorizationInformationResponse, error)

	EditGlobalApplicationCommand(*disgo.EditGlobalApplicationCommand) (*disgo.ApplicationCommand, error)
	EditGuildApplicationCommand(*disgo.EditGuildApplicationCommand) (*disgo.ApplicationCommand, error)
	EditApplicationCommandPermissions(*disgo.EditApplicationCommandPermissions) (*disgo.GuildApplicationCommandPermissions, error)
	EditOriginalInteractionResponse(*disgo.EditOriginalInteractionResponse) (*disgo.Message, error)
	EditFollowupMessage(*disgo.EditFollowupMessage) (*disgo.Message, error)
	EditMessage(*disgo.EditMessage) (*disgo.Message, error)
	EditChannelPermissions(*disgo.EditChannelPermissions) error
	EditWebhookMessage(*disgo.EditWebhookMessage) (*disgo.Message, error)
	BatchEditApplicationCommandPermissions(*disgo.BatchEditApplicationCommandPermissions) (*disgo.GuildApplicationCommandPermissions, error)

	ModifyGuildWidget(*disgo.ModifyGuildWidget) (*disgo.GuildWidget, error)
	ModifyAutoModerationRule(*disgo.ModifyAutoModerationRule) (*disgo.AutoModerationRule, error)
	ModifyChannel(*disgo.ModifyChannel) (*disgo.Channel, error)
	ModifyChannelGroupDM(*disgo.ModifyChannelGroupDM) (*disgo.Channel, error)
	ModifyChannelGuild(*disgo.ModifyChannelGuild) (*disgo.Channel, error)
	ModifyChannelThread(*disgo.ModifyChannelThread) (*disgo.Channel, error)
	ModifyGuildEmoji(*disgo.ModifyGuildEmoji) (*disgo.Emoji, error)
	ModifyGuild(*disgo.ModifyGuild) (*disgo.Guild, error)
	ModifyGuildChannelPositions(*disgo.ModifyGuildChannelPositions) error
	ModifyGuildMember(*disgo.ModifyGuildMember) (*disgo.GuildMember, error)
	ModifyCurrentMember(*disgo.ModifyCurrentMember) (*disgo.GuildMember, error)
	ModifyGuildRolePositions(*disgo.ModifyGuildRolePositions) ([]*disgo.Role, error)
	ModifyGuildRole(*disgo.ModifyGuildRole) (*disgo.Role, error)
	ModifyGuildWelcomeScreen(*disgo.ModifyGuildWelcomeScreen) (*disgo.WelcomeScreen, error)
	ModifyCurrentUserVoiceState(*disgo.ModifyCurrentUserVoiceState) error
	ModifyUserVoiceState(*disgo.ModifyUserVoiceState) error
	ModifyGuildScheduledEvent(*disgo.ModifyGuildScheduledEvent) (*disgo.GuildScheduledEvent, error)
	ModifyGuildTemplate(*disgo.ModifyGuildTemplate) (*disgo.GuildTemplate, error)
	ModifyStageInstance(*disgo.ModifyStageInstance) (*disgo.StageInstance, error)
	ModifyGuildSticker(*disgo.ModifyGuildSticker) (*disgo.Sticker, error)
	ModifyCurrentUser(*disgo.ModifyCurrentUser) (*disgo.User, error)
	ModifyWebhook(*disgo.ModifyWebhook) (*disgo.Webhook, error)
	ModifyWebhookwithToken(*disgo.ModifyWebhookwithToken) (*disgo.Webhook, error)

	ListAutoModerationRulesForGuild(*disgo.ListAutoModerationRulesForGuild) ([]*disgo.AutoModerationAction, error)
	ListThreadMembers(*disgo.ListThreadMembers) ([]*disgo.ThreadMember, error)
	ListPublicArchivedThreads(*disgo.ListPublicArchivedThreads) (*disgo.ListPublicArchivedThreadsResponse, error)
	ListPrivateArchivedThreads(*disgo.ListPrivateArchivedThreads) (*disgo.ListPrivateArchivedThreadsResponse, error)
	ListJoinedPrivateArchivedThreads(*disgo.ListJoinedPrivateArchivedThreads) (*disgo.ListJoinedPrivateArchivedThreadsResponse, error)
	ListGuildEmojis(*disgo.ListGuildEmojis) ([]*disgo.Emoji, error)
	ListActiveGuildThreads(*disgo.ListActiveGuildThreads) (*disgo.ListActiveGuildThreadsResponse, error)
	ListGuildMembers(*disgo.ListGuildMembers) ([]*disgo.GuildMember, error)
	ListScheduledEventsforGuild(*disgo.ListScheduledEventsforGuild) ([]*disgo.GuildScheduledEvent, error)
	ListNitroStickerPacks(*disgo.ListNitroStickerPacks) ([]*disgo.StickerPack, error)
	ListGuildStickers(*disgo.ListGuildStickers) ([]*disgo.Sticker, error)
	ListVoiceRegions(*disgo.ListVoiceRegions) ([]*disgo.VoiceRegion, error)

	GroupDMAddRecipient(*disgo.GroupDMAddRecipient) error
	AddThreadMember(*disgo.AddThreadMember) error
	AddGuildMember(*disgo.AddGuildMember) (*disgo.GuildMember, error)
	AddGuildMemberRole(*disgo.AddGuildMemberRole) error

	StartThreadfromMessage(*disgo.StartThreadfromMessage) (*disgo.Channel, error)
	StartThreadwithoutMessage(*disgo.StartThreadwithoutMessage) (*disgo.Channel, error)
	StartThreadinForumChannel(*disgo.StartThreadinForumChannel) (*disgo.Channel, error)

	BulkOverwriteGlobalApplicationCommands(*disgo.BulkOverwriteGlobalApplicationCommands) ([]*disgo.ApplicationCommand, error)
	BulkOverwriteGuildApplicationCommands(*disgo.BulkOverwriteGuildApplicationCommands) ([]*disgo.ApplicationCommand, error)

	CrosspostMessage(*disgo.CrosspostMessage) (*disgo.Message, error)

	FollowNewsChannel(*disgo.FollowNewsChannel) (*disgo.FollowedChannel, error)
	TriggerTypingIndicator(*disgo.TriggerTypingIndicator) error

	PinMessage(*disgo.PinMessage) error
	UnpinMessage(*disgo.UnpinMessage) error

	JoinThread(*disgo.JoinThread) error

	SearchGuildMembers(*disgo.SearchGuildMembers) ([]*disgo.GuildMember, error)

	BeginGuildPrune(*disgo.BeginGuildPrune) error

	SyncGuildTemplate(*disgo.SyncGuildTemplate) (*disgo.GuildTemplate, error)

	ExecuteWebhook(*disgo.ExecuteWebhook) error
	ExecuteSlackCompatibleWebhook(*disgo.ExecuteSlackCompatibleWebhook) error
	ExecuteGitHubCompatibleWebhook(*disgo.ExecuteGitHubCompatibleWebhook) error

	GroupDMRemoveRecipient(*disgo.GroupDMRemoveRecipient) error
	RemoveThreadMember(*disgo.RemoveThreadMember) error
	RemoveGuildMemberRole(*disgo.RemoveGuildMemberRole) error
	RemoveGuildMember(*disgo.RemoveGuildMember) error
	RemoveGuildBan(*disgo.RemoveGuildBan) error

	DeleteGlobalApplicationCommand(*disgo.DeleteGlobalApplicationCommand) error
	DeleteGuildApplicationCommand(*disgo.DeleteGuildApplicationCommand) error
	DeleteOriginalInteractionResponse(*disgo.DeleteOriginalInteractionResponse) error
	DeleteFollowupMessage(*disgo.DeleteFollowupMessage) error
	DeleteAutoModerationRule(*disgo.DeleteAutoModerationRule) error
	DeleteCloseChannel(*disgo.DeleteCloseChannel) (*disgo.Channel, error)
	DeleteOwnReaction(*disgo.DeleteOwnReaction) error
	DeleteUserReaction(*disgo.DeleteUserReaction) error
	DeleteAllReactions(*disgo.DeleteAllReactions) error
	DeleteAllReactionsforEmoji(*disgo.DeleteAllReactionsforEmoji) error
	DeleteMessage(*disgo.DeleteMessage) error
	BulkDeleteMessages(*disgo.BulkDeleteMessages) error
	DeleteChannelPermission(*disgo.DeleteChannelPermission) error
	DeleteGuildEmoji(*disgo.DeleteGuildEmoji) error
	DeleteGuild(*disgo.DeleteGuild) error
	DeleteGuildRole(*disgo.DeleteGuildRole) error
	DeleteGuildIntegration(*disgo.DeleteGuildIntegration) error
	DeleteGuildScheduledEvent(*disgo.DeleteGuildScheduledEvent) error
	DeleteGuildTemplate(*disgo.DeleteGuildTemplate) (*disgo.GuildTemplate, error)
	DeleteInvite(*disgo.DeleteInvite) (*disgo.Invite, error)
	DeleteStageInstance(*disgo.DeleteStageInstance) error
	DeleteGuildSticker(*disgo.DeleteGuildSticker) error
	DeleteWebhook(*disgo.DeleteWebhook) error
	DeleteWebhookwithToken(*disgo.DeleteWebhookwithToken) error
	DeleteWebhookMessage(*disgo.DeleteWebhookMessage) error

	LeaveThread(*disgo.LeaveThread) error
	LeaveGuild(*disgo.LeaveGuild) error
}
