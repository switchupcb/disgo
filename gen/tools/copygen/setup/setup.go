package wrapper

import (
	disgo "github.com/switchupcb/disgo/wrapper"
)

// Copygen defines the functions that will be generated.
type Copygen interface {
	SendGetGlobalApplicationCommands(*disgo.GetGlobalApplicationCommands) ([]*disgo.ApplicationCommand, error)
	SendCreateGlobalApplicationCommand(*disgo.CreateGlobalApplicationCommand) (*disgo.ApplicationCommand, error)
	SendGetGlobalApplicationCommand(*disgo.GetGlobalApplicationCommand) (*disgo.ApplicationCommand, error)
	SendEditGlobalApplicationCommand(*disgo.EditGlobalApplicationCommand) (*disgo.ApplicationCommand, error)
	SendDeleteGlobalApplicationCommand(*disgo.DeleteGlobalApplicationCommand) error
	SendBulkOverwriteGlobalApplicationCommands(*disgo.BulkOverwriteGlobalApplicationCommands) ([]*disgo.ApplicationCommand, error)
	SendGetGuildApplicationCommands(*disgo.GetGuildApplicationCommands) ([]*disgo.ApplicationCommand, error)
	SendCreateGuildApplicationCommand(*disgo.CreateGuildApplicationCommand) (*disgo.ApplicationCommand, error)
	SendGetGuildApplicationCommand(*disgo.GetGuildApplicationCommand) (*disgo.ApplicationCommand, error)
	SendEditGuildApplicationCommand(*disgo.EditGuildApplicationCommand) (*disgo.ApplicationCommand, error)
	SendDeleteGuildApplicationCommand(*disgo.DeleteGuildApplicationCommand) error
	SendBulkOverwriteGuildApplicationCommands(*disgo.BulkOverwriteGuildApplicationCommands) ([]*disgo.ApplicationCommand, error)
	SendGetGuildApplicationCommandPermissions(*disgo.GetGuildApplicationCommandPermissions) (*disgo.GuildApplicationCommandPermissions, error)
	SendGetApplicationCommandPermissions(*disgo.GetApplicationCommandPermissions) (*disgo.GuildApplicationCommandPermissions, error)
	SendEditApplicationCommandPermissions(*disgo.EditApplicationCommandPermissions) (*disgo.GuildApplicationCommandPermissions, error)
	SendBatchEditApplicationCommandPermissions(*disgo.BatchEditApplicationCommandPermissions) (*disgo.GuildApplicationCommandPermissions, error)
	SendCreateInteractionResponse(*disgo.CreateInteractionResponse) error
	SendGetOriginalInteractionResponse(*disgo.GetOriginalInteractionResponse) error
	SendEditOriginalInteractionResponse(*disgo.EditOriginalInteractionResponse) (*disgo.Message, error)
	SendDeleteOriginalInteractionResponse(*disgo.DeleteOriginalInteractionResponse) error
	SendCreateFollowupMessage(*disgo.CreateFollowupMessage) (*disgo.Message, error)
	SendGetFollowupMessage(*disgo.GetFollowupMessage) (*disgo.Message, error)
	SendEditFollowupMessage(*disgo.EditFollowupMessage) (*disgo.Message, error)
	SendDeleteFollowupMessage(*disgo.DeleteFollowupMessage) error
	SendGetGuildAuditLog(*disgo.GetGuildAuditLog) (*disgo.AuditLog, error)
	SendGetChannel(*disgo.GetChannel) (*disgo.Channel, error)
	SendModifyChannel(*disgo.ModifyChannel) (*disgo.Channel, error)
	SendModifyChannelGroupDM(*disgo.ModifyChannelGroupDM) (*disgo.Channel, error)
	SendModifyChannelGuild(*disgo.ModifyChannelGuild) (*disgo.Channel, error)
	SendModifyChannelThread(*disgo.ModifyChannelThread) (*disgo.Channel, error)
	SendDeleteCloseChannel(*disgo.DeleteCloseChannel) (*disgo.Channel, error)
	SendGetChannelMessages(*disgo.GetChannelMessages) ([]*disgo.Message, error)
	SendGetChannelMessage(*disgo.GetChannelMessage) (*disgo.Message, error)
	SendCreateMessage(*disgo.CreateMessage) (*disgo.Message, error)
	SendCrosspostMessage(*disgo.CrosspostMessage) (*disgo.Message, error)
	SendCreateReaction(*disgo.CreateReaction) error
	SendDeleteOwnReaction(*disgo.DeleteOwnReaction) error
	SendDeleteUserReaction(*disgo.DeleteUserReaction) error
	SendGetReactions(*disgo.GetReactions) ([]*disgo.User, error)
	SendDeleteAllReactions(*disgo.DeleteAllReactions) error
	SendDeleteAllReactionsforEmoji(*disgo.DeleteAllReactionsforEmoji) error
	SendEditMessage(*disgo.EditMessage) (*disgo.Message, error)
	SendDeleteMessage(*disgo.DeleteMessage) error
	SendBulkDeleteMessages(*disgo.BulkDeleteMessages) error
	SendEditChannelPermissions(*disgo.EditChannelPermissions) error
	SendGetChannelInvites(*disgo.GetChannelInvites) ([]*disgo.Invite, error)
	SendCreateChannelInvite(*disgo.CreateChannelInvite) (*disgo.Invite, error)
	SendDeleteChannelPermission(*disgo.DeleteChannelPermission) error
	// SendFollowNewsChannel(*disgo.FollowNewsChannel) (*disgo.FollowedChannel, error)
	SendTriggerTypingIndicator(*disgo.TriggerTypingIndicator) error
	SendGetPinnedMessages(*disgo.GetPinnedMessages) ([]*disgo.Message, error)
	SendPinMessage(*disgo.PinMessage) error
	SendUnpinMessage(*disgo.UnpinMessage) error
	SendGroupDMAddRecipient(*disgo.GroupDMAddRecipient) error
	SendGroupDMRemoveRecipient(*disgo.GroupDMRemoveRecipient) error
	SendStartThreadfromMessage(*disgo.StartThreadfromMessage) (*disgo.Channel, error)
	SendStartThreadwithoutMessage(*disgo.StartThreadwithoutMessage) (*disgo.Channel, error)
	SendStartThreadinForumChannel(*disgo.StartThreadinForumChannel) (*disgo.Channel, error)
	SendStartThreadinForumChannelMessage(*disgo.StartThreadinForumChannelMessage) (*disgo.Channel, error)
	SendJoinThread(*disgo.JoinThread) error
	SendAddThreadMember(*disgo.AddThreadMember) error
	SendLeaveThread(*disgo.LeaveThread) error
	SendRemoveThreadMember(*disgo.RemoveThreadMember) error
	SendGetThreadMember(*disgo.GetThreadMember) (*disgo.ThreadMember, error)
	SendListThreadMembers(*disgo.ListThreadMembers) ([]*disgo.ThreadMember, error)
	SendListPublicArchivedThreads(*disgo.ListPublicArchivedThreads) (*disgo.ListPublicArchivedThreadsResponse, error)
	SendListPrivateArchivedThreads(*disgo.ListPrivateArchivedThreads) (*disgo.ListPrivateArchivedThreadsResponse, error)
	SendListJoinedPrivateArchivedThreads(*disgo.ListJoinedPrivateArchivedThreads) (*disgo.ListJoinedPrivateArchivedThreadsResponse, error)
	SendListGuildEmojis(*disgo.ListGuildEmojis) ([]*disgo.Emoji, error)
	SendGetGuildEmoji(*disgo.GetGuildEmoji) (*disgo.Emoji, error)
	SendCreateGuildEmoji(*disgo.CreateGuildEmoji) (*disgo.Emoji, error)
	SendModifyGuildEmoji(*disgo.ModifyGuildEmoji) (*disgo.Emoji, error)
	SendDeleteGuildEmoji(*disgo.DeleteGuildEmoji) error
	SendListScheduledEventsforGuild(*disgo.ListScheduledEventsforGuild) ([]*disgo.GuildScheduledEvent, error)
	SendCreateGuildScheduledEvent(*disgo.CreateGuildScheduledEvent) (*disgo.GuildScheduledEvent, error)
	SendGetGuildScheduledEvent(*disgo.GetGuildScheduledEvent) (*disgo.GuildScheduledEvent, error)
	SendModifyGuildScheduledEvent(*disgo.ModifyGuildScheduledEvent) (*disgo.GuildScheduledEvent, error)
	SendDeleteGuildScheduledEvent(*disgo.DeleteGuildScheduledEvent) error
	SendGetGuildScheduledEventUsers(*disgo.GetGuildScheduledEventUsers) ([]*disgo.GuildScheduledEventUser, error)
	SendGetGuildTemplate(*disgo.GetGuildTemplate) (*disgo.GuildTemplate, error)
	SendCreateGuildfromGuildTemplate(*disgo.CreateGuildfromGuildTemplate) ([]*disgo.GuildTemplate, error)
	SendGetGuildTemplates(*disgo.GetGuildTemplates) ([]*disgo.GuildTemplate, error)
	SendCreateGuildTemplate(*disgo.CreateGuildTemplate) (*disgo.GuildTemplate, error)
	SendSyncGuildTemplate(*disgo.SyncGuildTemplate) (*disgo.GuildTemplate, error)
	SendModifyGuildTemplate(*disgo.ModifyGuildTemplate) (*disgo.GuildTemplate, error)
	SendDeleteGuildTemplate(*disgo.DeleteGuildTemplate) (*disgo.GuildTemplate, error)
	SendCreateGuild(*disgo.CreateGuild) (*disgo.Guild, error)
	SendGetGuild(*disgo.GetGuild) (*disgo.Guild, error)
	SendGetGuildPreview(*disgo.GetGuildPreview) (*disgo.GuildPreview, error)
	SendModifyGuild(*disgo.ModifyGuild) (*disgo.Guild, error)
	SendDeleteGuild(*disgo.DeleteGuild) error
	SendGetGuildChannels(*disgo.GetGuildChannels) ([]*disgo.Channel, error)
	SendCreateGuildChannel(*disgo.CreateGuildChannel) (*disgo.Channel, error)
	SendModifyGuildChannelPositions(*disgo.ModifyGuildChannelPositions) error
	SendListActiveGuildThreads(*disgo.ListActiveGuildThreads) (*disgo.ListActiveThreadsResponse, error)
	SendGetGuildMember(*disgo.GetGuildMember) (*disgo.GuildMember, error)
	SendListGuildMembers(*disgo.ListGuildMembers) ([]*disgo.GuildMember, error)
	SendSearchGuildMembers(*disgo.SearchGuildMembers) ([]*disgo.GuildMember, error)
	SendAddGuildMember(*disgo.AddGuildMember) (*disgo.GuildMember, error)
	SendModifyGuildMember(*disgo.ModifyGuildMember) (*disgo.GuildMember, error)
	SendModifyCurrentMember(*disgo.ModifyCurrentMember) (*disgo.GuildMember, error)
	SendModifyCurrentUserNick(*disgo.ModifyCurrentUserNick) (*disgo.ModifyCurrentUserNick, error)
	SendAddGuildMemberRole(*disgo.AddGuildMemberRole) error
	SendRemoveGuildMemberRole(*disgo.RemoveGuildMemberRole) error
	SendRemoveGuildMember(*disgo.RemoveGuildMember) error
	SendGetGuildBans(*disgo.GetGuildBans) ([]*disgo.Ban, error)
	SendGetGuildBan(*disgo.GetGuildBan) (*disgo.Ban, error)
	SendCreateGuildBan(*disgo.CreateGuildBan) error
	SendRemoveGuildBan(*disgo.RemoveGuildBan) error
	SendGetGuildRoles(*disgo.GetGuildRoles) ([]*disgo.Role, error)
	SendCreateGuildRole(*disgo.CreateGuildRole) (*disgo.Role, error)
	SendModifyGuildRolePositions(*disgo.ModifyGuildRolePositions) ([]*disgo.Role, error)
	SendModifyGuildRole(*disgo.ModifyGuildRole) (*disgo.Role, error)
	SendDeleteGuildRole(*disgo.DeleteGuildRole) error
	SendGetGuildPruneCount(*disgo.GetGuildPruneCount) error
	SendBeginGuildPrune(*disgo.BeginGuildPrune) error
	SendGetGuildVoiceRegions(*disgo.GetGuildVoiceRegions) (*disgo.VoiceRegion, error)
	SendGetGuildInvites(*disgo.GetGuildInvites) ([]*disgo.Invite, error)
	SendGetGuildIntegrations(*disgo.GetGuildIntegrations) ([]*disgo.Integration, error)
	SendDeleteGuildIntegration(*disgo.DeleteGuildIntegration) error
	SendGetGuildWidgetSettings(*disgo.GetGuildWidgetSettings) (*disgo.GuildWidget, error)
	SendModifyGuildWidget(*disgo.ModifyGuildWidget) (*disgo.GuildWidget, error)
	SendGetGuildWidget(*disgo.GetGuildWidget) (*disgo.GuildWidget, error)
	SendGetGuildVanityURL(*disgo.GetGuildVanityURL) (*disgo.Invite, error)
	SendGetGuildWidgetImage(*disgo.GetGuildWidgetImage) (*disgo.EmbedImage, error)
	SendGetGuildWelcomeScreen(*disgo.GetGuildWelcomeScreen) (*disgo.WelcomeScreen, error)
	SendModifyGuildWelcomeScreen(*disgo.ModifyGuildWelcomeScreen) (*disgo.WelcomeScreen, error)
	SendModifyCurrentUserVoiceState(*disgo.ModifyCurrentUserVoiceState) error
	SendModifyUserVoiceState(*disgo.ModifyUserVoiceState) error
	SendGetInvite(*disgo.GetInvite) (*disgo.Invite, error)
	SendDeleteInvite(*disgo.DeleteInvite) (*disgo.Invite, error)
	SendCreateStageInstance(*disgo.CreateStageInstance) (*disgo.StageInstance, error)
	SendGetStageInstance(*disgo.GetStageInstance) error
	SendModifyStageInstance(*disgo.ModifyStageInstance) (*disgo.StageInstance, error)
	SendDeleteStageInstance(*disgo.DeleteStageInstance) error
	SendGetSticker(*disgo.GetSticker) (*disgo.Sticker, error)
	SendListNitroStickerPacks(*disgo.ListNitroStickerPacks) ([]*disgo.StickerPack, error)
	SendListGuildStickers(*disgo.ListGuildStickers) ([]*disgo.Sticker, error)
	SendGetGuildSticker(*disgo.GetGuildSticker) (*disgo.Sticker, error)
	SendCreateGuildSticker(*disgo.CreateGuildSticker) (*disgo.Sticker, error)
	SendModifyGuildSticker(*disgo.ModifyGuildSticker) (*disgo.Sticker, error)
	SendDeleteGuildSticker(*disgo.DeleteGuildSticker) error
	SendModifyCurrentUser(*disgo.ModifyCurrentUser) (*disgo.User, error)
	SendGetCurrentUserGuilds(*disgo.GetCurrentUserGuilds) ([]*disgo.Guild, error)
	SendGetCurrentUserGuildMember(*disgo.GetCurrentUserGuildMember) (*disgo.GuildMember, error)
	SendLeaveGuild(*disgo.LeaveGuild) error
	SendCreateGroupDM(*disgo.CreateGroupDM) (*disgo.Channel, error)
	SendGetUserConnections(*disgo.GetUserConnections) ([]*disgo.Connection, error)
	SendListVoiceRegions(*disgo.ListVoiceRegions) ([]*disgo.VoiceRegion, error)
	SendCreateWebhook(*disgo.CreateWebhook) (*disgo.Webhook, error)
	SendGetChannelWebhooks(*disgo.GetChannelWebhooks) ([]*disgo.Webhook, error)
	SendGetGuildWebhooks(*disgo.GetGuildWebhooks) ([]*disgo.Webhook, error)
	SendGetWebhook(*disgo.GetWebhook) (*disgo.Webhook, error)
	SendGetWebhookwithToken(*disgo.GetWebhookwithToken) (*disgo.Webhook, error)
	SendModifyWebhook(*disgo.ModifyWebhook) (*disgo.Webhook, error)
	SendModifyWebhookwithToken(*disgo.ModifyWebhookwithToken) (*disgo.Webhook, error)
	SendDeleteWebhook(*disgo.DeleteWebhook) error
	SendDeleteWebhookwithToken(*disgo.DeleteWebhookwithToken) error
	SendExecuteWebhook(*disgo.ExecuteWebhook) error
	SendExecuteSlackCompatibleWebhook(*disgo.ExecuteSlackCompatibleWebhook) error
	SendExecuteGitHubCompatibleWebhook(*disgo.ExecuteGitHubCompatibleWebhook) error
	SendGetWebhookMessage(*disgo.GetWebhookMessage) (*disgo.Message, error)
	SendEditWebhookMessage(*disgo.EditWebhookMessage) (*disgo.Message, error)
	SendDeleteWebhookMessage(*disgo.DeleteWebhookMessage) error
	SendGetGateway(*disgo.GetGateway) (*disgo.GetGateway, error)
	SendGetGatewayBot(*disgo.GetGatewayBot) (*disgo.GetGatewayBot, error)
	SendGetCurrentBotApplicationInformation(*disgo.GetCurrentBotApplicationInformation) (*disgo.Application, error)
	SendGetCurrentAuthorizationInformation(*disgo.GetCurrentAuthorizationInformation) (*disgo.CurrentAuthorizationInformationResponse, error)
}
