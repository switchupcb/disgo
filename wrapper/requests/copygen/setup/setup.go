// Package requests contains the setup information for copygen generated code.
package requests

import (
	"github.com/switchupcb/disgo/wrapper/requests"
	"github.com/switchupcb/disgo/wrapper/requests/responses"
	"github.com/switchupcb/disgo/wrapper/resources"
)

// Copygen defines the functions that will be generated.
type Copygen interface {
	SendGetGlobalApplicationCommands(*requests.GetGlobalApplicationCommands) ([]*resources.ApplicationCommand, error)
	SendCreateGlobalApplicationCommand(*requests.CreateGlobalApplicationCommand) (*resources.ApplicationCommand, error)
	SendGetGlobalApplicationCommand(*requests.GetGlobalApplicationCommand) (*resources.ApplicationCommand, error)
	SendEditGlobalApplicationCommand(*requests.EditGlobalApplicationCommand) (*resources.ApplicationCommand, error)
	SendDeleteGlobalApplicationCommand(*requests.DeleteGlobalApplicationCommand) error
	SendBulkOverwriteGlobalApplicationCommands(*requests.BulkOverwriteGlobalApplicationCommands) ([]*resources.ApplicationCommand, error)
	SendGetGuildApplicationCommands(*requests.GetGuildApplicationCommands) ([]*resources.ApplicationCommand, error)
	SendCreateGuildApplicationCommand(*requests.CreateGuildApplicationCommand) (*resources.ApplicationCommand, error)
	SendGetGuildApplicationCommand(*requests.GetGuildApplicationCommand) (*resources.ApplicationCommand, error)
	SendEditGuildApplicationCommand(*requests.EditGuildApplicationCommand) (*resources.ApplicationCommand, error)
	SendDeleteGuildApplicationCommand(*requests.DeleteGuildApplicationCommand) error
	SendBulkOverwriteGuildApplicationCommands(*requests.BulkOverwriteGuildApplicationCommands) ([]*resources.ApplicationCommand, error)
	SendGetGuildApplicationCommandPermissions(*requests.GetGuildApplicationCommandPermissions) (*resources.GuildApplicationCommandPermissions, error)
	SendGetApplicationCommandPermissions(*requests.GetApplicationCommandPermissions) (*resources.GuildApplicationCommandPermissions, error)
	SendEditApplicationCommandPermissions(*requests.EditApplicationCommandPermissions) (*resources.GuildApplicationCommandPermissions, error)
	SendBatchEditApplicationCommandPermissions(*requests.BatchEditApplicationCommandPermissions) (*resources.GuildApplicationCommandPermissions, error)
	SendCreateInteractionResponse(*requests.CreateInteractionResponse) error
	SendGetOriginalInteractionResponse(*requests.GetOriginalInteractionResponse) error
	SendEditOriginalInteractionResponse(*requests.EditOriginalInteractionResponse) (*resources.Message, error)
	SendDeleteOriginalInteractionResponse(*requests.DeleteOriginalInteractionResponse) error
	SendCreateFollowupMessage(*requests.CreateFollowupMessage) (*resources.Message, error)
	SendGetFollowupMessage(*requests.GetFollowupMessage) (*resources.Message, error)
	SendEditFollowupMessage(*requests.EditFollowupMessage) (*resources.Message, error)
	SendDeleteFollowupMessage(*requests.DeleteFollowupMessage) error
	SendGetGuildAuditLog(*requests.GetGuildAuditLog) (*resources.AuditLog, error)
	SendGetChannel(*requests.GetChannel) (*resources.Channel, error)
	SendModifyChannel(*requests.ModifyChannel) (*resources.Channel, error)
	SendModifyChannelGroupDM(*requests.ModifyChannelGroupDM) (*resources.Channel, error)
	SendModifyChannelGuild(*requests.ModifyChannelGuild) (*resources.Channel, error)
	SendModifyChannelThread(*requests.ModifyChannelThread) (*resources.Channel, error)
	SendDeleteCloseChannel(*requests.DeleteCloseChannel) (*resources.Channel, error)
	SendGetChannelMessages(*requests.GetChannelMessages) ([]*resources.Message, error)
	SendGetChannelMessage(*requests.GetChannelMessage) (*resources.Message, error)
	SendCreateMessage(*requests.CreateMessage) (*resources.Message, error)
	SendCrosspostMessage(*requests.CrosspostMessage) (*resources.Message, error)
	SendCreateReaction(*requests.CreateReaction) error
	SendDeleteOwnReaction(*requests.DeleteOwnReaction) error
	SendDeleteUserReaction(*requests.DeleteUserReaction) error
	SendGetReactions(*requests.GetReactions) ([]*resources.User, error)
	SendDeleteAllReactions(*requests.DeleteAllReactions) error
	SendDeleteAllReactionsforEmoji(*requests.DeleteAllReactionsforEmoji) error
	SendEditMessage(*requests.EditMessage) (*resources.Message, error)
	SendDeleteMessage(*requests.DeleteMessage) error
	SendBulkDeleteMessages(*requests.BulkDeleteMessages) error
	SendEditChannelPermissions(*requests.EditChannelPermissions) error
	SendGetChannelInvites(*requests.GetChannelInvites) ([]*resources.Invite, error)
	SendCreateChannelInvite(*requests.CreateChannelInvite) (*resources.Invite, error)
	SendDeleteChannelPermission(*requests.DeleteChannelPermission) error
	SendFollowNewsChannel(*requests.FollowNewsChannel) (*resources.FollowedChannel, error)
	SendTriggerTypingIndicator(*requests.TriggerTypingIndicator) error
	SendGetPinnedMessages(*requests.GetPinnedMessages) ([]*resources.Message, error)
	SendPinMessage(*requests.PinMessage) error
	SendUnpinMessage(*requests.UnpinMessage) error
	SendGroupDMAddRecipient(*requests.GroupDMAddRecipient) error
	SendGroupDMRemoveRecipient(*requests.GroupDMRemoveRecipient) error
	SendStartThreadfromMessage(*requests.StartThreadfromMessage) (*resources.Channel, error)
	SendStartThreadwithoutMessage(*requests.StartThreadwithoutMessage) (*resources.Channel, error)
	SendStartThreadinForumChannel(*requests.StartThreadinForumChannel) (*resources.Channel, error)
	SendStartThreadinForumChannelMessage(*requests.StartThreadinForumChannelMessage) (*resources.Channel, error)
	SendJoinThread(*requests.JoinThread) error
	SendAddThreadMember(*requests.AddThreadMember) error
	SendLeaveThread(*requests.LeaveThread) error
	SendRemoveThreadMember(*requests.RemoveThreadMember) error
	SendGetThreadMember(*requests.GetThreadMember) (*resources.ThreadMember, error)
	SendListThreadMembers(*requests.ListThreadMembers) ([]*resources.ThreadMember, error)
	SendListActiveChannelThreads(*requests.ListActiveChannelThreads) (*responses.ListActiveThreadsResponse, error)
	SendListPublicArchivedThreads(*requests.ListPublicArchivedThreads) (*responses.ListPublicArchivedThreadsResponse, error)
	SendListPrivateArchivedThreads(*requests.ListPrivateArchivedThreads) (*responses.ListPrivateArchivedThreadsResponse, error)
	SendListJoinedPrivateArchivedThreads(*requests.ListJoinedPrivateArchivedThreads) (*responses.ListJoinedPrivateArchivedThreadsResponse, error)
	SendListGuildEmojis(*requests.ListGuildEmojis) ([]*resources.Emoji, error)
	SendGetGuildEmoji(*requests.GetGuildEmoji) (*resources.Emoji, error)
	SendCreateGuildEmoji(*requests.CreateGuildEmoji) (*resources.Emoji, error)
	SendModifyGuildEmoji(*requests.ModifyGuildEmoji) (*resources.Emoji, error)
	SendDeleteGuildEmoji(*requests.DeleteGuildEmoji) error
	SendListScheduledEventsforGuild(*requests.ListScheduledEventsforGuild) ([]*resources.GuildScheduledEvent, error)
	SendCreateGuildScheduledEvent(*requests.CreateGuildScheduledEvent) (*resources.GuildScheduledEvent, error)
	SendGetGuildScheduledEvent(*requests.GetGuildScheduledEvent) (*resources.GuildScheduledEvent, error)
	SendModifyGuildScheduledEvent(*requests.ModifyGuildScheduledEvent) (*resources.GuildScheduledEvent, error)
	SendDeleteGuildScheduledEvent(*requests.DeleteGuildScheduledEvent) error
	SendGetGuildScheduledEventUsers(*requests.GetGuildScheduledEventUsers) ([]*resources.GuildScheduledEventUser, error)
	SendGetGuildTemplate(*requests.GetGuildTemplate) (*resources.GuildTemplate, error)
	SendCreateGuildfromGuildTemplate(*requests.CreateGuildfromGuildTemplate) ([]*resources.GuildTemplate, error)
	SendGetGuildTemplates(*requests.GetGuildTemplates) ([]*resources.GuildTemplate, error)
	SendCreateGuildTemplate(*requests.CreateGuildTemplate) (*resources.GuildTemplate, error)
	SendSyncGuildTemplate(*requests.SyncGuildTemplate) (*resources.GuildTemplate, error)
	SendModifyGuildTemplate(*requests.ModifyGuildTemplate) (*resources.GuildTemplate, error)
	SendDeleteGuildTemplate(*requests.DeleteGuildTemplate) (*resources.GuildTemplate, error)
	SendCreateGuild(*requests.CreateGuild) (*resources.Guild, error)
	SendGetGuild(*requests.GetGuild) (*resources.Guild, error)
	SendGetGuildPreview(*requests.GetGuildPreview) (*resources.GuildPreview, error)
	SendModifyGuild(*requests.ModifyGuild) (*resources.Guild, error)
	SendDeleteGuild(*requests.DeleteGuild) error
	SendGetGuildChannels(*requests.GetGuildChannels) ([]*resources.Channel, error)
	SendCreateGuildChannel(*requests.CreateGuildChannel) (*resources.Channel, error)
	SendModifyGuildChannelPositions(*requests.ModifyGuildChannelPositions) error
	SendListActiveGuildThreads(*requests.ListActiveGuildThreads) (*responses.ListActiveThreadsResponse, error)
	SendGetGuildMember(*requests.GetGuildMember) (*resources.GuildMember, error)
	SendListGuildMembers(*requests.ListGuildMembers) ([]*resources.GuildMember, error)
	SendSearchGuildMembers(*requests.SearchGuildMembers) ([]*resources.GuildMember, error)
	SendAddGuildMember(*requests.AddGuildMember) (*resources.GuildMember, error)
	SendModifyGuildMember(*requests.ModifyGuildMember) (*resources.GuildMember, error)
	SendModifyCurrentMember(*requests.ModifyCurrentMember) (*resources.GuildMember, error)
	SendModifyCurrentUserNick(*requests.ModifyCurrentUserNick) (*responses.ModifyCurrentUserNick, error)
	SendAddGuildMemberRole(*requests.AddGuildMemberRole) error
	SendRemoveGuildMemberRole(*requests.RemoveGuildMemberRole) error
	SendRemoveGuildMember(*requests.RemoveGuildMember) error
	SendGetGuildBans(*requests.GetGuildBans) ([]*resources.Ban, error)
	SendGetGuildBan(*requests.GetGuildBan) (*resources.Ban, error)
	SendCreateGuildBan(*requests.CreateGuildBan) error
	SendRemoveGuildBan(*requests.RemoveGuildBan) error
	SendGetGuildRoles(*requests.GetGuildRoles) ([]*resources.Role, error)
	SendCreateGuildRole(*requests.CreateGuildRole) (*resources.Role, error)
	SendModifyGuildRolePositions(*requests.ModifyGuildRolePositions) ([]*resources.Role, error)
	SendModifyGuildRole(*requests.ModifyGuildRole) (*resources.Role, error)
	SendDeleteGuildRole(*requests.DeleteGuildRole) error
	//SendGetGuildPruneCount(*requests.GetGuildPruneCount) error  https://discord.com/developers/docs/resources/guild#get-guild-prune-count
	//SendBeginGuildPrune(*requests.BeginGuildPrune) error        https://discord.com/developers/docs/resources/guild#get-guild-prune-count
	SendGetGuildVoiceRegions(*requests.GetGuildVoiceRegions) (*resources.VoiceRegion, error)
	SendGetGuildInvites(*requests.GetGuildInvites) ([]*resources.Invite, error)
	SendGetGuildIntegrations(*requests.GetGuildIntegrations) ([]*resources.Integration, error)
	SendDeleteGuildIntegration(*requests.DeleteGuildIntegration) error
	SendGetGuildWidgetSettings(*requests.GetGuildWidgetSettings) (*resources.GuildWidget, error)
	SendModifyGuildWidget(*requests.ModifyGuildWidget) (*resources.GuildWidget, error)
	SendGetGuildWidget(*requests.GetGuildWidget) (*resources.GuildWidget, error)
	SendGetGuildVanityURL(*requests.GetGuildVanityURL) (*resources.Invite, error)
	SendGetGuildWidgetImage(*requests.GetGuildWidgetImage) (*resources.EmbedImage, error)
	SendGetGuildWelcomeScreen(*requests.GetGuildWelcomeScreen) (*resources.WelcomeScreen, error)
	SendModifyGuildWelcomeScreen(*requests.ModifyGuildWelcomeScreen) (*resources.WelcomeScreen, error)
	SendModifyCurrentUserVoiceState(*requests.ModifyCurrentUserVoiceState) error
	SendModifyUserVoiceState(*requests.ModifyUserVoiceState) error
	SendGetInvite(*requests.GetInvite) (*resources.Invite, error)
	SendDeleteInvite(*requests.DeleteInvite) (*resources.Invite, error)
	SendCreateStageInstance(*requests.CreateStageInstance) (*resources.StageInstance, error)
	SendGetStageInstance(*requests.GetStageInstance) error
	SendModifyStageInstance(*requests.ModifyStageInstance) (*resources.StageInstance, error)
	SendDeleteStageInstance(*requests.DeleteStageInstance) error
	SendGetSticker(*requests.GetSticker) (*resources.Sticker, error)
	SendListNitroStickerPacks(*requests.ListNitroStickerPacks) ([]*resources.StickerPack, error)
	SendListGuildStickers(*requests.ListGuildStickers) ([]*resources.Sticker, error)
	SendGetGuildSticker(*requests.GetGuildSticker) (*resources.Sticker, error)
	SendCreateGuildSticker(*requests.CreateGuildSticker) (*resources.Sticker, error)
	SendModifyGuildSticker(*requests.ModifyGuildSticker) (*resources.Sticker, error)
	SendDeleteGuildSticker(*requests.DeleteGuildSticker) error
	SendModifyCurrentUser(*requests.ModifyCurrentUser) (*resources.User, error)
	SendGetCurrentUserGuilds(*requests.GetCurrentUserGuilds) ([]*resources.Guild, error)
	SendGetCurrentUserGuildMember(*requests.GetCurrentUserGuildMember) (*resources.GuildMember, error)
	SendLeaveGuild(*requests.LeaveGuild) error
	SendCreateGroupDM(*requests.CreateGroupDM) (*resources.Channel, error)
	SendGetUserConnections(*requests.GetUserConnections) ([]*resources.Connection, error)
	SendListVoiceRegions(*requests.ListVoiceRegions) ([]*resources.VoiceRegion, error)
	SendCreateWebhook(*requests.CreateWebhook) (*resources.Webhook, error)
	SendGetChannelWebhooks(*requests.GetChannelWebhooks) ([]*resources.Webhook, error)
	SendGetGuildWebhooks(*requests.GetGuildWebhooks) ([]*resources.Webhook, error)
	SendGetWebhook(*requests.GetWebhook) (*resources.Webhook, error)
	SendGetWebhookwithToken(*requests.GetWebhookwithToken) (*resources.Webhook, error)
	SendModifyWebhook(*requests.ModifyWebhook) (*resources.Webhook, error)
	SendModifyWebhookwithToken(*requests.ModifyWebhookwithToken) (*resources.Webhook, error)
	SendDeleteWebhook(*requests.DeleteWebhook) error
	SendDeleteWebhookwithToken(*requests.DeleteWebhookwithToken) error
	SendExecuteWebhook(*requests.ExecuteWebhook) error
	SendExecuteSlackCompatibleWebhook(*requests.ExecuteSlackCompatibleWebhook) error
	SendExecuteGitHubCompatibleWebhook(*requests.ExecuteGitHubCompatibleWebhook) error
	SendGetWebhookMessage(*requests.GetWebhookMessage) (*resources.Message, error)
	SendEditWebhookMessage(*requests.EditWebhookMessage) (*resources.Message, error)
	SendDeleteWebhookMessage(*requests.DeleteWebhookMessage) error
	SendGetGateway(*requests.GetGateway) (*responses.GetGateway, error)
	SendGetGatewayBot(*requests.GetGatewayBot) (*responses.GetGatewayBot, error)
	SendGetCurrentBotApplicationInformation(*requests.GetCurrentBotApplicationInformation) (*resources.Application, error)
	SendGetCurrentAuthorizationInformation(*requests.GetCurrentAuthorizationInformation) (*responses.CurrentAuthorizationInformation, error)
}
