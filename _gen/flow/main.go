package main

import "fmt"

// TODO: note that Group DM requests are not being used

var (
	// endpoints represents a dependency graph of endpoints (map[dependency][]endpoints).
	endpoints = map[string][]string{
		"GetGlobalApplicationCommands":           {"CreateGlobalApplicationCommand"},
		"CreateGlobalApplicationCommand":         {},
		"GetGlobalApplicationCommand":            {"CreateGlobalApplicationCommand"},
		"EditGlobalApplicationCommand":           {"CreateGlobalApplicationCommand"},
		"DeleteGlobalApplicationCommand":         {"CreateGlobalApplicationCommand"},
		"BulkOverwriteGlobalApplicationCommands": {"CreateGlobalApplicationCommand"},
		"GetGuildApplicationCommands":            {"CreateGlobalApplicationCommand"},
		"CreateGuildApplicationCommand":          {},
		"GetGuildApplicationCommand":             {"CreateGuildApplicationCommand"},
		"EditGuildApplicationCommand":            {"CreateGuildApplicationCommand"},
		"DeleteGuildApplicationCommand":          {"CreateGuildApplicationCommand"},
		"BulkOverwriteGuildApplicationCommands":  {"CreateGuildApplicationCommand"},
		"GetGuildApplicationCommandPermissions":  {"CreateGuildApplicationCommand"},
		"GetApplicationCommandPermissions":       {"CreateGuildApplicationCommand"},
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
		"BeginGuildPrune":                        {"CreateGuild"},
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

	// dependents represents a map of the amount of endpoints that
	// an endpoint depends on (map[endpoint]numDependents).
	dependents := make(map[string]int, numEndpoints)

	// calculate the dependents of the endpoints.
	for endpoint, dependencies := range endpoints {
		dependents[endpoint] = len(dependencies)
	}

	// queue represents a first-in first-out data structure.
	//
	// queue can't have more entries than the number of endpoints,
	// so initialize a map of length 0 with capacity = numEndpoints.
	queue := make([]string, 0, numEndpoints)

	// fill the queue with endpoints that have no dependencies.
	for endpoint, numDependencies := range dependents {
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

		// decrease the amount of endpoints remaining.
		numEndpoints--

		// add endpoints that no longer depend on any other endpoints to the queue.
		//
		// in other words, add endpoints with no dependencies.
		for endpoint, dependencies := range endpoints {
			// when an endpoint depends on the current endpoint (that is now accounted for),
			// decrement the amount of dependents the endpoint has.
			if contains(dependencies, current) {
				dependents[endpoint]--

				// when the endpoint is dependent on no other endpoints, add it to the queue.
				if dependents[endpoint] == 0 {
					queue = append(queue, endpoint)

					if numEndpoints < 0 {
						fmt.Println("WARNING: cycle occurred")

						return []string{}
					}
				}
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

	for _, endpoint := range findOrder(endpoints) {
		// fmt.Println(fmt.Sprintf("%d. %v", i, endpoint))
		fmt.Println(endpoint)
	}

}
