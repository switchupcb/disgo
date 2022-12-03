package wrapper

import (
	"github.com/switchupcb/disgo"
)

// Copygen defines the functions that will be generated.
type Copygen interface {
	// http GET
	GetGlobalApplicationCommands(*disgo.GetGlobalApplicationCommands) ([]*disgo.ApplicationCommand, error)
	// http POST
	CreateGlobalApplicationCommand(*disgo.CreateGlobalApplicationCommand) (*disgo.ApplicationCommand, error)
	// http GET
	GetGlobalApplicationCommand(*disgo.GetGlobalApplicationCommand) (*disgo.ApplicationCommand, error)
	// http PATCH
	EditGlobalApplicationCommand(*disgo.EditGlobalApplicationCommand) (*disgo.ApplicationCommand, error)
	// http DELETE
	DeleteGlobalApplicationCommand(*disgo.DeleteGlobalApplicationCommand) error
	// http PUT
	BulkOverwriteGlobalApplicationCommands(*disgo.BulkOverwriteGlobalApplicationCommands) ([]*disgo.ApplicationCommand, error)
	// http GET
	GetGuildApplicationCommands(*disgo.GetGuildApplicationCommands) ([]*disgo.ApplicationCommand, error)
	// http POST
	CreateGuildApplicationCommand(*disgo.CreateGuildApplicationCommand) (*disgo.ApplicationCommand, error)
	// http GET
	GetGuildApplicationCommand(*disgo.GetGuildApplicationCommand) (*disgo.ApplicationCommand, error)
	// http PATCH
	EditGuildApplicationCommand(*disgo.EditGuildApplicationCommand) (*disgo.ApplicationCommand, error)
	// http DELETE
	DeleteGuildApplicationCommand(*disgo.DeleteGuildApplicationCommand) error
	// http PUT
	BulkOverwriteGuildApplicationCommands(*disgo.BulkOverwriteGuildApplicationCommands) ([]*disgo.ApplicationCommand, error)
	// http GET
	GetGuildApplicationCommandPermissions(*disgo.GetGuildApplicationCommandPermissions) (*disgo.GuildApplicationCommandPermissions, error)
	// http GET
	GetApplicationCommandPermissions(*disgo.GetApplicationCommandPermissions) (*disgo.GuildApplicationCommandPermissions, error)
	// http PUT
	EditApplicationCommandPermissions(*disgo.EditApplicationCommandPermissions) (*disgo.GuildApplicationCommandPermissions, error)
	// http PUT
	BatchEditApplicationCommandPermissions(*disgo.BatchEditApplicationCommandPermissions) (*disgo.GuildApplicationCommandPermissions, error)
	// http POST
	CreateInteractionResponse(*disgo.CreateInteractionResponse) error
	// http PATCH
	GetOriginalInteractionResponse(*disgo.GetOriginalInteractionResponse) error
	// http PATCH
	EditOriginalInteractionResponse(*disgo.EditOriginalInteractionResponse) (*disgo.Message, error)
	// http DELETE
	DeleteOriginalInteractionResponse(*disgo.DeleteOriginalInteractionResponse) error
	// http POST
	CreateFollowupMessage(*disgo.CreateFollowupMessage) (*disgo.Message, error)
	// http GET
	GetFollowupMessage(*disgo.GetFollowupMessage) (*disgo.Message, error)
	// http PATCH
	EditFollowupMessage(*disgo.EditFollowupMessage) (*disgo.Message, error)
	// http DELETE
	DeleteFollowupMessage(*disgo.DeleteFollowupMessage) error
	// http GET
	GetGuildAuditLog(*disgo.GetGuildAuditLog) (*disgo.AuditLog, error)
	// http GET
	ListAutoModerationRulesForGuild(*disgo.ListAutoModerationRulesForGuild) ([]*disgo.AutoModerationAction, error)
	// http GET
	GetAutoModerationRule(*disgo.GetAutoModerationRule) (*disgo.AutoModerationRule, error)
	// http POST
	CreateAutoModerationRule(*disgo.CreateAutoModerationRule) (*disgo.AutoModerationRule, error)
	// http PATCH
	ModifyAutoModerationRule(*disgo.ModifyAutoModerationRule) (*disgo.AutoModerationRule, error)
	// http DELETE
	DeleteAutoModerationRule(*disgo.DeleteAutoModerationRule) error
	// http GET
	GetChannel(*disgo.GetChannel) (*disgo.Channel, error)
	// http PATCH
	ModifyChannel(*disgo.ModifyChannel) (*disgo.Channel, error)
	// http PATCH
	ModifyChannelGroupDM(*disgo.ModifyChannelGroupDM) (*disgo.Channel, error)
	// http PATCH
	ModifyChannelGuild(*disgo.ModifyChannelGuild) (*disgo.Channel, error)
	// http PATCH
	ModifyChannelThread(*disgo.ModifyChannelThread) (*disgo.Channel, error)
	// http DELETE
	DeleteCloseChannel(*disgo.DeleteCloseChannel) (*disgo.Channel, error)
	// http GET
	GetChannelMessages(*disgo.GetChannelMessages) ([]*disgo.Message, error)
	// http GET
	GetChannelMessage(*disgo.GetChannelMessage) (*disgo.Message, error)
	// http POST
	CreateMessage(*disgo.CreateMessage) (*disgo.Message, error)
	// http POST
	CrosspostMessage(*disgo.CrosspostMessage) (*disgo.Message, error)
	// http PUT
	CreateReaction(*disgo.CreateReaction) error
	// http DELETE
	DeleteOwnReaction(*disgo.DeleteOwnReaction) error
	// http DELETE
	DeleteUserReaction(*disgo.DeleteUserReaction) error
	// http GET
	GetReactions(*disgo.GetReactions) ([]*disgo.User, error)
	// http DELETE
	DeleteAllReactions(*disgo.DeleteAllReactions) error
	// http DELETE
	DeleteAllReactionsforEmoji(*disgo.DeleteAllReactionsforEmoji) error
	// http PATCH
	EditMessage(*disgo.EditMessage) (*disgo.Message, error)
	// http DELETE
	DeleteMessage(*disgo.DeleteMessage) error
	// http POST
	BulkDeleteMessages(*disgo.BulkDeleteMessages) error
	// http PUT
	EditChannelPermissions(*disgo.EditChannelPermissions) error
	// http GET
	GetChannelInvites(*disgo.GetChannelInvites) ([]*disgo.Invite, error)
	// http POST
	CreateChannelInvite(*disgo.CreateChannelInvite) (*disgo.Invite, error)
	// http DELETE
	DeleteChannelPermission(*disgo.DeleteChannelPermission) error
	// http POST
	FollowAnnouncementChannel(*disgo.FollowAnnouncementChannel) (*disgo.FollowedChannel, error)
	// http POST
	TriggerTypingIndicator(*disgo.TriggerTypingIndicator) error
	// http GET
	GetPinnedMessages(*disgo.GetPinnedMessages) ([]*disgo.Message, error)
	// http PUT
	PinMessage(*disgo.PinMessage) error
	// http DELETE
	UnpinMessage(*disgo.UnpinMessage) error
	// http PUT
	GroupDMAddRecipient(*disgo.GroupDMAddRecipient) error
	// http DELETE
	GroupDMRemoveRecipient(*disgo.GroupDMRemoveRecipient) error
	// http POST
	StartThreadfromMessage(*disgo.StartThreadfromMessage) (*disgo.Channel, error)
	// http POST
	StartThreadwithoutMessage(*disgo.StartThreadwithoutMessage) (*disgo.Channel, error)
	// http POST
	StartThreadinForumChannel(*disgo.StartThreadinForumChannel) (*disgo.Channel, error)
	// http PUT
	JoinThread(*disgo.JoinThread) error
	// http PUT
	AddThreadMember(*disgo.AddThreadMember) error
	// http DELETE
	LeaveThread(*disgo.LeaveThread) error
	// http DELETE
	RemoveThreadMember(*disgo.RemoveThreadMember) error
	// http GET
	GetThreadMember(*disgo.GetThreadMember) (*disgo.ThreadMember, error)
	// http GET
	ListThreadMembers(*disgo.ListThreadMembers) ([]*disgo.ThreadMember, error)
	// http GET
	ListPublicArchivedThreads(*disgo.ListPublicArchivedThreads) (*disgo.ListPublicArchivedThreadsResponse, error)
	// http GET
	ListPrivateArchivedThreads(*disgo.ListPrivateArchivedThreads) (*disgo.ListPrivateArchivedThreadsResponse, error)
	// http GET
	ListJoinedPrivateArchivedThreads(*disgo.ListJoinedPrivateArchivedThreads) (*disgo.ListJoinedPrivateArchivedThreadsResponse, error)
	// http GET
	ListGuildEmojis(*disgo.ListGuildEmojis) ([]*disgo.Emoji, error)
	// http GET
	GetGuildEmoji(*disgo.GetGuildEmoji) (*disgo.Emoji, error)
	// http POST
	CreateGuildEmoji(*disgo.CreateGuildEmoji) (*disgo.Emoji, error)
	// http PATCH
	ModifyGuildEmoji(*disgo.ModifyGuildEmoji) (*disgo.Emoji, error)
	// http DELETE
	DeleteGuildEmoji(*disgo.DeleteGuildEmoji) error
	// http POST
	CreateGuild(*disgo.CreateGuild) (*disgo.Guild, error)
	// http GET
	GetGuild(*disgo.GetGuild) (*disgo.Guild, error)
	// http GET
	GetGuildPreview(*disgo.GetGuildPreview) (*disgo.GuildPreview, error)
	// http PATCH
	ModifyGuild(*disgo.ModifyGuild) (*disgo.Guild, error)
	// http DELETE
	DeleteGuild(*disgo.DeleteGuild) error
	// http GET
	GetGuildChannels(*disgo.GetGuildChannels) ([]*disgo.Channel, error)
	// http POST
	CreateGuildChannel(*disgo.CreateGuildChannel) (*disgo.Channel, error)
	// http PATCH
	ModifyGuildChannelPositions(*disgo.ModifyGuildChannelPositions) error
	// http GET
	ListActiveGuildThreads(*disgo.ListActiveGuildThreads) (*disgo.ListActiveGuildThreadsResponse, error)
	// http GET
	GetGuildMember(*disgo.GetGuildMember) (*disgo.GuildMember, error)
	// http GET
	ListGuildMembers(*disgo.ListGuildMembers) ([]*disgo.GuildMember, error)
	// http GET
	SearchGuildMembers(*disgo.SearchGuildMembers) ([]*disgo.GuildMember, error)
	// http PUT
	AddGuildMember(*disgo.AddGuildMember) (*disgo.GuildMember, error)
	// http PATCH
	ModifyGuildMember(*disgo.ModifyGuildMember) (*disgo.GuildMember, error)
	// http PATCH
	ModifyCurrentMember(*disgo.ModifyCurrentMember) (*disgo.GuildMember, error)
	// http PUT
	AddGuildMemberRole(*disgo.AddGuildMemberRole) error
	// http DELETE
	RemoveGuildMemberRole(*disgo.RemoveGuildMemberRole) error
	// http DELETE
	RemoveGuildMember(*disgo.RemoveGuildMember) error
	// http GET
	GetGuildBans(*disgo.GetGuildBans) ([]*disgo.Ban, error)
	// http GET
	GetGuildBan(*disgo.GetGuildBan) (*disgo.Ban, error)
	// http PUT
	CreateGuildBan(*disgo.CreateGuildBan) error
	// http DELETE
	RemoveGuildBan(*disgo.RemoveGuildBan) error
	// http GET
	GetGuildRoles(*disgo.GetGuildRoles) ([]*disgo.Role, error)
	// http POST
	CreateGuildRole(*disgo.CreateGuildRole) (*disgo.Role, error)
	// http PATCH
	ModifyGuildRolePositions(*disgo.ModifyGuildRolePositions) ([]*disgo.Role, error)
	// http PATCH
	ModifyGuildRole(*disgo.ModifyGuildRole) (*disgo.Role, error)
	// http DELETE
	DeleteGuildRole(*disgo.DeleteGuildRole) error
	// http POST
	ModifyGuildMFALevel(*disgo.ModifyGuildMFALevel) (*disgo.ModifyGuildMFALevelResponse, error)
	// http GET
	GetGuildPruneCount(*disgo.GetGuildPruneCount) (*disgo.GetGuildPruneCountResponse, error)
	// http POST
	BeginGuildPrune(*disgo.BeginGuildPrune) error
	// http GET
	GetGuildVoiceRegions(*disgo.GetGuildVoiceRegions) ([]*disgo.VoiceRegion, error)
	// http GET
	GetGuildInvites(*disgo.GetGuildInvites) ([]*disgo.Invite, error)
	// http GET
	GetGuildIntegrations(*disgo.GetGuildIntegrations) ([]*disgo.Integration, error)
	// http DELETE
	DeleteGuildIntegration(*disgo.DeleteGuildIntegration) error
	// http GET
	GetGuildWidgetSettings(*disgo.GetGuildWidgetSettings) (*disgo.GuildWidget, error)
	// http PATCH
	ModifyGuildWidget(*disgo.ModifyGuildWidget) (*disgo.GuildWidget, error)
	// http GET
	GetGuildWidget(*disgo.GetGuildWidget) (*disgo.GuildWidget, error)
	// http GET
	GetGuildVanityURL(*disgo.GetGuildVanityURL) (*disgo.Invite, error)
	// http GET
	GetGuildWidgetImage(*disgo.GetGuildWidgetImage) (*disgo.EmbedImage, error)
	// http GET
	GetGuildWelcomeScreen(*disgo.GetGuildWelcomeScreen) (*disgo.WelcomeScreen, error)
	// http PATCH
	ModifyGuildWelcomeScreen(*disgo.ModifyGuildWelcomeScreen) (*disgo.WelcomeScreen, error)
	// http PATCH
	ModifyCurrentUserVoiceState(*disgo.ModifyCurrentUserVoiceState) error
	// http PATCH
	ModifyUserVoiceState(*disgo.ModifyUserVoiceState) error
	// http GET
	ListScheduledEventsforGuild(*disgo.ListScheduledEventsforGuild) ([]*disgo.GuildScheduledEvent, error)
	// http POST
	CreateGuildScheduledEvent(*disgo.CreateGuildScheduledEvent) (*disgo.GuildScheduledEvent, error)
	// http GET
	GetGuildScheduledEvent(*disgo.GetGuildScheduledEvent) (*disgo.GuildScheduledEvent, error)
	// http PATCH
	ModifyGuildScheduledEvent(*disgo.ModifyGuildScheduledEvent) (*disgo.GuildScheduledEvent, error)
	// http DELETE
	DeleteGuildScheduledEvent(*disgo.DeleteGuildScheduledEvent) error
	// http GET
	GetGuildScheduledEventUsers(*disgo.GetGuildScheduledEventUsers) ([]*disgo.GuildScheduledEventUser, error)
	// http GET
	GetGuildTemplate(*disgo.GetGuildTemplate) (*disgo.GuildTemplate, error)
	// http POST
	CreateGuildfromGuildTemplate(*disgo.CreateGuildfromGuildTemplate) ([]*disgo.GuildTemplate, error)
	// http GET
	GetGuildTemplates(*disgo.GetGuildTemplates) ([]*disgo.GuildTemplate, error)
	// http POST
	CreateGuildTemplate(*disgo.CreateGuildTemplate) (*disgo.GuildTemplate, error)
	// http PUT
	SyncGuildTemplate(*disgo.SyncGuildTemplate) (*disgo.GuildTemplate, error)
	// http PATCH
	ModifyGuildTemplate(*disgo.ModifyGuildTemplate) (*disgo.GuildTemplate, error)
	// http DELETE
	DeleteGuildTemplate(*disgo.DeleteGuildTemplate) (*disgo.GuildTemplate, error)
	// http GET
	GetInvite(*disgo.GetInvite) (*disgo.Invite, error)
	// http DELETE
	DeleteInvite(*disgo.DeleteInvite) (*disgo.Invite, error)
	// http POST
	CreateStageInstance(*disgo.CreateStageInstance) (*disgo.StageInstance, error)
	// http GET
	GetStageInstance(*disgo.GetStageInstance) (*disgo.StageInstance, error)
	// http PATCH
	ModifyStageInstance(*disgo.ModifyStageInstance) (*disgo.StageInstance, error)
	// http DELETE
	DeleteStageInstance(*disgo.DeleteStageInstance) error
	// http GET
	GetSticker(*disgo.GetSticker) (*disgo.Sticker, error)
	// http GET
	ListNitroStickerPacks(*disgo.ListNitroStickerPacks) (*disgo.ListNitroStickerPacksResponse, error)
	// http GET
	ListGuildStickers(*disgo.ListGuildStickers) ([]*disgo.Sticker, error)
	// http GET
	GetGuildSticker(*disgo.GetGuildSticker) (*disgo.Sticker, error)
	// http POST
	CreateGuildSticker(*disgo.CreateGuildSticker) (*disgo.Sticker, error)
	// http PATCH
	ModifyGuildSticker(*disgo.ModifyGuildSticker) (*disgo.Sticker, error)
	// http DELETE
	DeleteGuildSticker(*disgo.DeleteGuildSticker) error
	// http GET
	GetCurrentUser(*disgo.GetCurrentUser) (*disgo.User, error)
	// http GET
	GetUser(*disgo.GetUser) (*disgo.User, error)
	// http PATCH
	ModifyCurrentUser(*disgo.ModifyCurrentUser) (*disgo.User, error)
	// http GET
	GetCurrentUserGuilds(*disgo.GetCurrentUserGuilds) ([]*disgo.Guild, error)
	// http GET
	GetCurrentUserGuildMember(*disgo.GetCurrentUserGuildMember) (*disgo.GuildMember, error)
	// http DELETE
	LeaveGuild(*disgo.LeaveGuild) error
	// http POST
	CreateDM(*disgo.CreateDM) (*disgo.Channel, error)
	// http POST
	CreateGroupDM(*disgo.CreateGroupDM) (*disgo.Channel, error)
	// http GET
	GetUserConnections(*disgo.GetUserConnections) ([]*disgo.Connection, error)
	// http GET
	ListVoiceRegions(*disgo.ListVoiceRegions) ([]*disgo.VoiceRegion, error)
	// http POST
	CreateWebhook(*disgo.CreateWebhook) (*disgo.Webhook, error)
	// http GET
	GetChannelWebhooks(*disgo.GetChannelWebhooks) ([]*disgo.Webhook, error)
	// http GET
	GetGuildWebhooks(*disgo.GetGuildWebhooks) ([]*disgo.Webhook, error)
	// http GET
	GetWebhook(*disgo.GetWebhook) (*disgo.Webhook, error)
	// http GET
	GetWebhookwithToken(*disgo.GetWebhookwithToken) (*disgo.Webhook, error)
	// http PATCH
	ModifyWebhook(*disgo.ModifyWebhook) (*disgo.Webhook, error)
	// http PATCH
	ModifyWebhookwithToken(*disgo.ModifyWebhookwithToken) (*disgo.Webhook, error)
	// http DELETE
	DeleteWebhook(*disgo.DeleteWebhook) error
	// http DELETE
	DeleteWebhookwithToken(*disgo.DeleteWebhookwithToken) error
	// http POST
	ExecuteWebhook(*disgo.ExecuteWebhook) error
	// http POST
	ExecuteSlackCompatibleWebhook(*disgo.ExecuteSlackCompatibleWebhook) error
	// http POST
	ExecuteGitHubCompatibleWebhook(*disgo.ExecuteGitHubCompatibleWebhook) error
	// http GET
	GetWebhookMessage(*disgo.GetWebhookMessage) (*disgo.Message, error)
	// http PATCH
	EditWebhookMessage(*disgo.EditWebhookMessage) (*disgo.Message, error)
	// http DELETE
	DeleteWebhookMessage(*disgo.DeleteWebhookMessage) error
	// http GET
	GetGateway(*disgo.GetGateway) (*disgo.GetGatewayBotResponse, error)
	// http GET
	GetGatewayBot(*disgo.GetGatewayBot) (*disgo.GetGatewayBotResponse, error)
	// http GET
	GetCurrentBotApplicationInformation(*disgo.GetCurrentBotApplicationInformation) (*disgo.Application, error)
	// http GET
	GetCurrentAuthorizationInformation(*disgo.GetCurrentAuthorizationInformation) (*disgo.CurrentAuthorizationInformationResponse, error)
}
