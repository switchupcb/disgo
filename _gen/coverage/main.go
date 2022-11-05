package main

import "fmt"

var (
	// endpoints represents a dependency graph of endpoints (map[dependency][]dependents).
	endpoints = map[string][]string{
		"GetGlobalApplicationCommands":           {"CreateGlobalApplicationCommand"},
		"CreateGlobalApplicationCommand":         {},
		"GetGlobalApplicationCommand":            {"CreateGlobalApplicationCommand"},
		"EditGlobalApplicationCommand":           {"CreateGlobalApplicationCommand"},
		"DeleteGlobalApplicationCommand":         {"CreateGlobalApplicationCommand"},
		"BulkOverwriteGlobalApplicationCommands": {"CreateGlobalApplicationCommand"},
		"GetGuildApplicationCommands":            {"CreateGuildApplicationCommand"},
		"CreateGuildApplicationCommand":          {"CreateGuild"},
		"GetGuildApplicationCommand":             {"CreateGuildApplicationCommand"},
		"EditGuildApplicationCommand":            {"CreateGuildApplicationCommand"},
		"DeleteGuildApplicationCommand":          {"CreateGuildApplicationCommand"},
		"BulkOverwriteGuildApplicationCommands":  {"CreateGuildApplicationCommand"},
		"GetGuildApplicationCommandPermissions":  {"CreateGuildApplicationCommand"},
		"GetApplicationCommandPermissions":       {"CreateGlobalApplicationCommand"},
		"EditApplicationCommandPermissions":      {"CreateGlobalApplicationCommand"},
		"BatchEditApplicationCommandPermissions": {"CreateGlobalApplicationCommand"},
		"CreateInteractionResponse":              {}, // Interaction Create Gateway Event
		"GetOriginalInteractionResponse":         {"CreateInteractionResponse"},
		"EditOriginalInteractionResponse":        {"CreateInteractionResponse"},
		"DeleteOriginalInteractionResponse":      {"CreateInteractionResponse"},
		"CreateFollowupMessage":                  {}, // Interaction Create Gateway Event
		"GetFollowupMessage":                     {"CreateFollowupMessage"},
		"EditFollowupMessage":                    {"CreateFollowupMessage"},
		"DeleteFollowupMessage":                  {"CreateFollowupMessage"},
		"GetGuildAuditLog":                       {"CreateGuild"},
		"ListAutoModerationRulesForGuild":        {"CreateGuild"},
		"GetAutoModerationRule":                  {"CreateAutoModerationRule"},
		"CreateAutoModerationRule":               {"CreateGuild"},
		"ModifyAutoModerationRule":               {"CreateAutoModerationRule"},
		"DeleteAutoModerationRule":               {"CreateAutoModerationRule"},
		"GetChannel":                             {"CreateGuildChannel"},
		"ModifyChannel":                          {"CreateGuildChannel"},
		"ModifyChannelGroupDM":                   {},
		"ModifyChannelGuild":                     {"CreateGuildChannel"},
		"ModifyChannelThread":                    {"StartThreadwithoutMessage"},
		"DeleteCloseChannel":                     {"CreateGuildChannel"},
		"GetChannelMessages":                     {"CreateGuildChannel"},
		"GetChannelMessage":                      {"CreateMessage"},
		"CreateMessage":                          {"CreateGuildChannel"},
		"CrosspostMessage":                       {"CreateMessage"},
		"CreateReaction":                         {"CreateMessage"},
		"DeleteOwnReaction":                      {"CreateReaction"},
		"DeleteUserReaction":                     {"CreateReaction"},
		"GetReactions":                           {"CreateMessage"},
		"DeleteAllReactions":                     {"CreateMessage"},
		"DeleteAllReactionsforEmoji":             {"CreateReaction"},
		"EditMessage":                            {"CreateMessage"},
		"DeleteMessage":                          {"CreateMessage"},
		"BulkDeleteMessages":                     {"CreateMessage"},
		"EditChannelPermissions":                 {"CreateGuildChannel"},
		"GetChannelInvites":                      {"CreateGuildChannel"},
		"CreateChannelInvite":                    {"CreateGuildChannel"},
		"DeleteChannelPermission":                {"CreateGuildChannel"},
		"FollowAnnouncementChannel":              {"CreateGuildChannel"},
		"TriggerTypingIndicator":                 {"CreateGuildChannel"},
		"GetPinnedMessages":                      {"CreateGuildChannel"},
		"PinMessage":                             {"CreateMessage"},
		"UnpinMessage":                           {"PinMessage"},
		"GroupDMAddRecipient":                    {},
		"GroupDMRemoveRecipient":                 {},
		"StartThreadfromMessage":                 {"CreateGuildChannel"},
		"StartThreadwithoutMessage":              {"CreateGuildChannel"},
		"StartThreadinForumChannel":              {"CreateGuildChannel"},
		"JoinThread":                             {"StartThreadwithoutMessage"},
		"AddThreadMember":                        {"StartThreadwithoutMessage"},
		"LeaveThread":                            {"JoinThread"},
		"RemoveThreadMember":                     {"AddThreadMember"},
		"GetThreadMember":                        {"AddThreadMember"},
		"ListThreadMembers":                      {"StartThreadwithoutMessage"},
		"ListPublicArchivedThreads":              {"CreateGuildChannel"},
		"ListPrivateArchivedThreads":             {"CreateGuildChannel"},
		"ListJoinedPrivateArchivedThreads":       {"CreateGuildChannel"},
		"ListGuildEmojis":                        {"CreateGuildChannel"},
		"GetGuildEmoji":                          {"CreateGuildEmoji"},
		"CreateGuildEmoji":                       {},
		"ModifyGuildEmoji":                       {"CreateGuildEmoji"},
		"DeleteGuildEmoji":                       {"CreateGuildEmoji"},
		"CreateGuild":                            {},
		"GetGuild":                               {"CreateGuild"},
		"GetGuildPreview":                        {"CreateGuild"},
		"ModifyGuild":                            {"CreateGuild"},
		"DeleteGuild":                            {"CreateGuild"},
		"GetGuildChannels":                       {"CreateGuild"},
		"CreateGuildChannel":                     {"CreateGuild"},
		"ModifyGuildChannelPositions":            {"CreateGuild", "CreateGuildChannel"},
		"ListActiveGuildThreads":                 {"CreateGuild"},
		"GetGuildMember":                         {"AddGuildMember"},
		"ListGuildMembers":                       {"CreateGuild"},
		"SearchGuildMembers":                     {"CreateGuild"},
		"AddGuildMember":                         {"CreateGuild"},
		"ModifyGuildMember":                      {"AddGuildMember"},
		"ModifyCurrentMember":                    {"CreateGuild"},
		"AddGuildMemberRole":                     {"AddGuildMember", "CreateGuildRole"},
		"RemoveGuildMemberRole":                  {"AddGuildMemberRole"},
		"RemoveGuildMember":                      {"AddGuildMember"},
		"GetGuildBans":                           {"CreateGuild"},
		"GetGuildBan":                            {"CreateGuildBan"},
		"CreateGuildBan":                         {"AddGuildMember"},
		"RemoveGuildBan":                         {"CreateGuildBan"},
		"GetGuildRoles":                          {"CreateGuild"},
		"CreateGuildRole":                        {"CreateGuild"},
		"ModifyGuildRolePositions":               {"CreateGuild", "CreateGuildRole"},
		"ModifyGuildRole":                        {"CreateGuildRole"},
		"DeleteGuildRole":                        {"CreateGuildRole"},
		"GetGuildPruneCount":                     {"CreateGuild"},
		"BeginGuildPrune":                        {"CreateGuild"},
		"GetGuildVoiceRegions":                   {"CreateGuild"},
		"GetGuildInvites":                        {"CreateGuild"},
		"GetGuildIntegrations":                   {"CreateGuild"},
		"DeleteGuildIntegration":                 {}, // Client Required
		"GetGuildWidgetSettings":                 {"CreateGuild"},
		"ModifyGuildWidget":                      {"CreateGuild"},
		"GetGuildWidget":                         {"CreateGuild"},
		"GetGuildVanityURL":                      {"CreateGuild"},
		"GetGuildWidgetImage":                    {"CreateGuild"},
		"GetGuildWelcomeScreen":                  {"CreateGuild"},
		"ModifyGuildWelcomeScreen":               {"CreateGuild"},
		"ModifyCurrentUserVoiceState":            {"CreateGuild"},
		"ModifyUserVoiceState":                   {"CreateGuild", "AddGuildMember"},
		"ListScheduledEventsforGuild":            {"CreateGuild"},
		"CreateGuildScheduledEvent":              {"CreateGuild"},
		"GetGuildScheduledEvent":                 {"CreateGuildScheduledEvent"},
		"ModifyGuildScheduledEvent":              {"CreateGuildScheduledEvent"},
		"DeleteGuildScheduledEvent":              {"CreateGuildScheduledEvent"},
		"GetGuildScheduledEventUsers":            {"CreateGuildScheduledEvent"},
		"GetGuildTemplate":                       {"CreateGuild"},
		"CreateGuildfromGuildTemplate":           {"CreateGuildTemplate"},
		"GetGuildTemplates":                      {"CreateGuild"},
		"CreateGuildTemplate":                    {"CreateGuild"},
		"SyncGuildTemplate":                      {"CreateGuildTemplate"},
		"ModifyGuildTemplate":                    {"CreateGuildTemplate"},
		"DeleteGuildTemplate":                    {"CreateGuildTemplate"},
		"GetInvite":                              {"GetGuildInvites"},
		"DeleteInvite":                           {"GetInvite"},
		"CreateStageInstance":                    {"CreateGuild"},
		"GetStageInstance":                       {"CreateStageInstance"},
		"ModifyStageInstance":                    {"CreateStageInstance"},
		"DeleteStageInstance":                    {"CreateStageInstance"},
		"GetSticker":                             {"CreateGuildSticker"},
		"ListNitroStickerPacks":                  {},
		"ListGuildStickers":                      {"CreateGuildSticker"},
		"GetGuildSticker":                        {"CreateGuildSticker"},
		"CreateGuildSticker":                     {"CreateGuildSticker"},
		"ModifyGuildSticker":                     {"CreateGuildSticker"},
		"DeleteGuildSticker":                     {"CreateGuildSticker"},
		"GetCurrentUser":                         {},
		"GetUser":                                {},
		"ModifyCurrentUser":                      {"GetUser"},
		"GetCurrentUserGuilds":                   {"GetUser"},
		"GetCurrentUserGuildMember":              {"CreateGuild"},
		"LeaveGuild":                             {"CreateGuild"},
		"CreateGroupDM":                          {},
		"GetUserConnections":                     {},
		"ListVoiceRegions":                       {},
		"CreateWebhook":                          {"CreateGuildChannel"},
		"GetChannelWebhooks":                     {"CreateGuildChannel"},
		"GetGuildWebhooks":                       {"CreateGuild"},
		"GetWebhook":                             {"GetWebhook"},
		"GetWebhookwithToken":                    {"CreateWebhook"},
		"ModifyWebhook":                          {"CreateWebhook"},
		"ModifyWebhookwithToken":                 {"GetWebhook"},
		"DeleteWebhook":                          {"CreateWebhook"},
		"DeleteWebhookwithToken":                 {"GetWebhook"},
		"ExecuteWebhook":                         {"CreateWebhook"},
		"ExecuteSlackCompatibleWebhook":          {"CreateWebhook"},
		"ExecuteGitHubCompatibleWebhook":         {"CreateWebhook"},
		"GetWebhookMessage":                      {"CreateWebhook"},
		"EditWebhookMessage":                     {"CreateWebhook", "GetWebhookMessage"},
		"DeleteWebhookMessage":                   {"CreateWebhook", "GetWebhookMessage"},
		"GetGateway":                             {},
		"GetGatewayBot":                          {},
		"GetCurrentBotApplicationInformation":    {},
		"GetCurrentAuthorizationInformation":     {},
	}

	// unused represents a map of unused endpoints.
	unused = map[string]bool{
		"BulkOverwriteGlobalApplicationCommands": true,
		"BulkOverwriteGuildApplicationCommands":  true,
		"BatchEditApplicationCommandPermissions": true,
		"ModifyChannelGroupDM":                   true,
		"BulkDeleteMessages":                     true,
		"TriggerTypingIndicator":                 true,
		"GroupDMAddRecipient":                    true,
		"GroupDMRemoveRecipient":                 true,
		"DeleteGuildIntegration":                 true,
		"CreateGroupDM":                          true,

		// Webhooks
		"CreateWebhook":                  true,
		"GetChannelWebhooks":             true,
		"GetGuildWebhooks":               true,
		"GetWebhook":                     true,
		"GetWebhookwithToken":            true,
		"ModifyWebhook":                  true,
		"ModifyWebhookwithToken":         true,
		"DeleteWebhook":                  true,
		"DeleteWebhookwithToken":         true,
		"ExecuteWebhook":                 true,
		"ExecuteSlackCompatibleWebhook":  true,
		"ExecuteGitHubCompatibleWebhook": true,
		"GetWebhookMessage":              true,
		"EditWebhookMessage":             true,
		"DeleteWebhookMessage":           true,

		// Session Test
		"GetGateway":    true,
		"GetGatewayBot": true,
	}
)

// filterEndpoints removes unused endpoints from the endpoint map.
func filterEndpoints(endpoints map[string][]string) {
	for dependency := range unused {
		delete(endpoints, dependency)
	}
}

// findOrder finds the order of endpoints with the least amount of dependencies
// to the most amount of dependencies.
func findOrder(endpoints map[string][]string) []string {
	filterEndpoints(endpoints)

	// TODO

	return []string{}
}

// contains returns whether the slice s contains the string x.
func contains(s []string, x string) bool {
	for _, item := range s {
		if x == item {
			return true
		}
	}

	return false
}

func main() {
	for i, endpoint := range findOrder(endpoints) {
		fmt.Printf("%d. %v\n", i, endpoint)
	}
}
