package main

import "fmt"

// TODO: note that Group DM requests are not being used

var (
	// endpoints represents a dependency graph of endpoints (map[endpoint][]dependencies).
	endpoints = map[string][]string{
		"GetGlobalApplicationCommands":           {"CreateGlobalApplicationCommand"},
		"CreateGlobalApplicationCommand":         {},
		"GetGlobalApplicationCommand":            {"CreateGlobalApplicationCommand"},
		"EditGlobalApplicationCommand":           {"CreateGlobalApplicationCommand"},
		"DeleteGlobalApplicationCommand":         {"CreateGlobalApplicationCommand"},
		"BulkOverwriteGlobalApplicationCommands": {"CreateGlobalApplicationCommand"},
		"GetGuildApplicationCommands":            {},
		"CreateGuildApplicationCommand":          {},
		"GetGuildApplicationCommand":             {"CreateGuildApplicationCommand"},
		"EditGuildApplicationCommand":            {"CreateGuildApplicationCommand"},
		"DeleteGuildApplicationCommand":          {"CreateGuildApplicationCommand"},
		"BulkOverwriteGuildApplicationCommands":  {"CreateGuildApplicationCommand"},
		"GetGuildApplicationCommandPermissions":  {},
		"GetApplicationCommandPermissions":       {},
		"EditApplicationCommandPermissions":      {},
		"BatchEditApplicationCommandPermissions": {},
		"CreateInteractionResponse":              {},
		"GetOriginalInteractionResponse":         {"CreateInteractionResponse"},
		"EditOriginalInteractionResponse":        {"CreateInteractionResponse"},
		"DeleteOriginalInteractionResponse":      {"CreateInteractionResponse"},
		"CreateFollowupMessage":                  {},
		"GetFollowupMessage":                     {"CreateFollowupMessage"},
		"EditFollowupMessage":                    {"CreateFollowupMessage"},
		"DeleteFollowupMessage":                  {"CreateFollowupMessage"},
		"GetGuildAuditLog":                       {"CreateGuild"},
		"ListAutoModerationRulesForGuild":        {"CreateAutoModerationRule"},
		"GetAutoModerationRule":                  {"CreateAutoModerationRule"},
		"CreateAutoModerationRule":               {"CreateGuild"},
		"ModifyAutoModerationRule":               {"CreateAutoModerationRule"},
		"DeleteAutoModerationRule":               {"CreateAutoModerationRule"},
		"GetChannel":                             {"CreateGuild"},
		"ModifyChannel":                          {"CreateGuild"},
		"ModifyChannelGuild":                     {"CreateGuild"},
		"ModifyChannelThread":                    {},
		"DeleteCloseChannel":                     {},
		"GetChannelMessages":                     {},
		"GetChannelMessage":                      {},
		"CreateMessage":                          {},
		"CrosspostMessage":                       {},
		"CreateReaction":                         {},
		"DeleteOwnReaction":                      {"CreateReaction"},
		"DeleteUserReaction":                     {"CreateReaction"},
		"GetReactions":                           {"CreateReaction"},
		"DeleteAllReactions":                     {"CreateReaction"},
		"DeleteAllReactionsforEmoji":             {},
		"EditMessage":                            {"CreateMessage"},
		"DeleteMessage":                          {"CreateMessage"},
		"BulkDeleteMessages":                     {"CreateMessage"},
		"EditChannelPermissions":                 {},
		"GetChannelInvites":                      {},
		"CreateChannelInvite":                    {},
		"DeleteChannelPermission":                {},
		"FollowNewsChannel":                      {},
		"TriggerTypingIndicator":                 {},
		"GetPinnedMessages":                      {},
		"PinMessage":                             {"CreateMessage"},
		"UnpinMessage":                           {"CreateMessage"},
		"StartThreadfromMessage":                 {"CreateMessage"},
		"StartThreadwithoutMessage":              {},
		"StartThreadinForumChannel":              {},
		"JoinThread":                             {"StartThreadfromMessage"},
		"AddThreadMember":                        {"StartThreadfromMessage"},
		"LeaveThread":                            {"StartThreadfromMessage"},
		"RemoveThreadMember":                     {"StartThreadfromMessage"},
		"GetThreadMember":                        {"StartThreadfromMessage"},
		"ListThreadMembers":                      {"StartThreadfromMessage"},
		"ListPublicArchivedThreads":              {"StartThreadfromMessage"},
		"ListPrivateArchivedThreads":             {"StartThreadfromMessage"},
		"ListJoinedPrivateArchivedThreads":       {"StartThreadfromMessage"},
		"ListGuildEmojis":                        {},
		"GetGuildEmoji":                          {},
		"CreateGuildEmoji":                       {},
		"ModifyGuildEmoji":                       {},
		"DeleteGuildEmoji":                       {},
		"CreateGuild":                            {},
		"GetGuild":                               {"CreateGuild"},
		"GetGuildPreview":                        {"CreateGuild"},
		"ModifyGuild":                            {"CreateGuild"},
		"DeleteGuild":                            {"CreateGuild"},
		"GetGuildChannels":                       {"CreateGuild"},
		"CreateGuildChannel":                     {"CreateGuild"},
		"ModifyGuildChannelPositions":            {},
		"ListActiveGuildThreads":                 {"CreateGuild"},
		"GetGuildMember":                         {"CreateGuild"},
		"ListGuildMembers":                       {"CreateGuild"},
		"SearchGuildMembers":                     {"CreateGuild"},
		"AddGuildMember":                         {"CreateGuild"},
		"ModifyGuildMember":                      {"CreateGuild"},
		"ModifyCurrentMember":                    {"AddGuildMember"},
		"AddGuildMemberRole":                     {"AddGuildMember"},
		"RemoveGuildMemberRole":                  {"AddGuildMember"},
		"RemoveGuildMember":                      {"AddGuildMember"},
		"GetGuildBans":                           {"CreateGuildBan"},
		"GetGuildBan":                            {"CreateGuildBan"},
		"CreateGuildBan":                         {"CreateGuild"},
		"RemoveGuildBan":                         {"CreateGuildBan"},
		"GetGuildRoles":                          {"CreateGuild"},
		"CreateGuildRole":                        {"CreateGuild"},
		"ModifyGuildRolePositions":               {"CreateGuild"},
		"ModifyGuildRole":                        {"CreateGuildRole"},
		"DeleteGuildRole":                        {"CreateGuildRole"},
		"GetGuildPruneCount":                     {"CreateGuild"},
		"BeginGuildPrune":                        {"CreateGuild", "AddGuildMember"},
		"GetGuildVoiceRegions":                   {"CreateGuild"},
		"GetGuildInvites":                        {"CreateGuild"},
		"GetGuildIntegrations":                   {"CreateGuild"},
		"DeleteGuildIntegration":                 {"CreateGuild"},
		"GetGuildWidgetSettings":                 {"CreateGuild"},
		"ModifyGuildWidget":                      {"CreateGuild"},
		"GetGuildWidget":                         {"CreateGuild"},
		"GetGuildVanityURL":                      {"CreateGuild"},
		"GetGuildWidgetImage":                    {"CreateGuild"},
		"GetGuildWelcomeScreen":                  {"CreateGuild"},
		"ModifyGuildWelcomeScreen":               {"CreateGuild"},
		"ModifyCurrentUserVoiceState":            {},
		"ModifyUserVoiceState":                   {},
		"ListScheduledEventsforGuild":            {"CreateGuildScheduledEvent"},
		"CreateGuildScheduledEvent":              {"CreateGuild"},
		"GetGuildScheduledEvent":                 {"CreateGuildScheduledEvent"},
		"ModifyGuildScheduledEvent":              {"CreateGuildScheduledEvent"},
		"DeleteGuildScheduledEvent":              {"CreateGuildScheduledEvent"},
		"GetGuildScheduledEventUsers":            {"CreateGuildScheduledEvent"},
		"GetGuildTemplate":                       {"CreateGuildTemplate"},
		"CreateGuildfromGuildTemplate":           {},
		"GetGuildTemplates":                      {},
		"CreateGuildTemplate":                    {},
		"SyncGuildTemplate":                      {"CreateGuildTemplate"},
		"ModifyGuildTemplate":                    {"CreateGuildTemplate"},
		"DeleteGuildTemplate":                    {"CreateGuildTemplate"},
		"GetInvite":                              {},
		"DeleteInvite":                           {},
		"CreateStageInstance":                    {},
		"GetStageInstance":                       {"CreateStageInstance"},
		"ModifyStageInstance":                    {"CreateStageInstance"},
		"DeleteStageInstance":                    {"CreateStageInstance"},
		"GetSticker":                             {},
		"ListNitroStickerPacks":                  {},
		"ListGuildStickers":                      {},
		"GetGuildSticker":                        {"CreateGuildSticker"},
		"CreateGuildSticker":                     {},
		"ModifyGuildSticker":                     {"CreateGuildSticker"},
		"DeleteGuildSticker":                     {"CreateGuildSticker"},
		"GetCurrentUser":                         {},
		"GetUser":                                {},
		"ModifyCurrentUser":                      {},
		"GetCurrentUserGuilds":                   {},
		"GetCurrentUserGuildMember":              {},
		"LeaveGuild":                             {"CreateGuild"},
		"GetUserConnections":                     {"GetUser"},
		"ListVoiceRegions":                       {},
		"CreateWebhook":                          {},
		"GetChannelWebhooks":                     {},
		"GetGuildWebhooks":                       {"CreateGuild"},
		"GetWebhook":                             {"CreateWebhook"},
		"GetWebhookwithToken":                    {"CreateWebhook"},
		"ModifyWebhook":                          {"CreateWebhook"},
		"ModifyWebhookwithToken":                 {"CreateWebhook"},
		"DeleteWebhook":                          {"CreateWebhook"},
		"DeleteWebhookwithToken":                 {"CreateWebhook"},
		"ExecuteWebhook":                         {"CreateWebhook"},
		"ExecuteSlackCompatibleWebhook":          {"CreateWebhook"},
		"ExecuteGitHubCompatibleWebhook":         {"CreateWebhook"},
		"GetWebhookMessage":                      {"CreateWebhook"},
		"EditWebhookMessage":                     {"CreateWebhook"},
		"DeleteWebhookMessage":                   {"CreateWebhook"},
		"GetGateway":                             {},
		"GetGatewayBot":                          {},
		"GetCurrentBotApplicationInformation":    {},
		"GetCurrentAuthorizationInformation":     {},
	}
)

// findOrder finds the order of endpoints with the least amount of dependencies
// to the most amount of dependencies.
func findOrder(endpoints map[string][]string) []string {
	numEndpoints := len(endpoints)

	// indegree represents the amount of endpoints required for an endpoint.
	indegree := make(map[string]int, numEndpoints)

	// calculate the indegrees of the endpoints.
	for endpoint, dependencies := range endpoints {
		indegree[endpoint] = len(dependencies)
	}

	// queue represents a first-in first-out data structure.
	//
	// queue can't have more entries than the number of endpoints,
	// so initialize a map of length 0 with capacity = numEndpoints.
	queue := make([]string, 0, numEndpoints)

	// fill the queue with endpoints that have no dependencies.
	for endpoint, numDependencies := range indegree {
		if numDependencies == 0 {
			queue = append(queue, endpoint)
		}
	}

	// output represents the returned result.
	output := make([]string, 0, numEndpoints)

	// add the endpoints with no dependencies to the output.
	for len(queue) > 0 {
		// select the first entry in the queue (of endpoints with no dependencies).
		current := queue[0]

		// remove the first entry out the queue (of endpoints with no dependencies).
		queue = queue[1:]

		// add the entry to the output.
		output = append(output, current)
		numEndpoints--

		// add endpoints that no longer contain dependencies to the queue.
		for _, dependency := range endpoints[current] {
			indegree[dependency]--

			// when the endpoints's dependency count is 0, add it to the queue.
			if indegree[dependency] == 0 {
				queue = append(queue, dependency)
			}
		}
	}

	// numEndpoints is not 0 when a cycle occurs (i.e [a: b],[b: a])
	if numEndpoints != 0 {
		fmt.Println("WARNING: cycle occurred")

		return []string{}
	}

	return output
}

func main() {
	for i, endpoint := range findOrder(endpoints) {
		fmt.Printf(fmt.Sprintf("%d. %v", i, endpoint))
	}
}
