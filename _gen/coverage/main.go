package main

import "fmt"

var (
	// endpoints represents a dependency graph of endpoints (map[dependent][]dependencies).
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
		"AddThreadMember":                        {"StartThreadwithoutMessage", "GetUser"},
		"LeaveThread":                            {"JoinThread"},
		"RemoveThreadMember":                     {"AddThreadMember"},
		"GetThreadMember":                        {"AddThreadMember"},
		"ListThreadMembers":                      {"StartThreadwithoutMessage"},
		"ListPublicArchivedThreads":              {"CreateGuildChannel"},
		"ListPrivateArchivedThreads":             {"CreateGuildChannel"},
		"ListJoinedPrivateArchivedThreads":       {"CreateGuildChannel"},
		"ListGuildEmojis":                        {"CreateGuildChannel"},
		"GetGuildEmoji":                          {"CreateGuildEmoji"},
		"CreateGuildEmoji":                       {"CreateGuild"},
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
		"CreateGuildSticker":                     {"CreateGuild"},
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
		"GetWebhook":                             {},
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
		"GetGuildIntegration":                    true,
		"DeleteGuildIntegration":                 true,
		"ListNitroStickerPacks":                  true,
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

// filterOutput removes unused endpoints from the endpoint output slice.
func filterOutput(endpoints []string) []string {
	output := make([]string, 0, len(unused))

	for _, endpoint := range endpoints {
		if unused[endpoint] {
			continue
		}

		output = append(output, endpoint)
	}

	return output
}

// findOrder finds the optimal order of endpoints using dependency graph.
func findOrder(endpoints map[string][]string) []string {
	numEndpoints := len(endpoints)

	// dependents represents a map of dependent endpoints to
	// the respective amount of dependencies (map[dependency]numDependencies).
	dependents := make(map[string]int, numEndpoints)

	// calculate the number of dependencies for each dependent.
	for endpoint, dependencies := range endpoints {
		dependents[endpoint] = len(dependencies)
	}

	// queue represents a first-in first-out data structure.
	//
	// queue can't have more entries than the number of endpoints,
	// so initialize a map of length 0 with capacity = numEndpoints.
	queue := make([]string, 0, numEndpoints)

	// fill the queue with dependent endpoints that have no dependencies (i.e `CreateGuild`).
	for endpoint, numDependencies := range dependents {
		if numDependencies == 0 {
			queue = append(queue, endpoint)
		}
	}

	// output represents the returned result.
	output := make([]string, 0, numEndpoints)

	// add dependent endpoints with no dependencies to the output list.
	for len(queue) > 0 {
		// select the first entry in the queue (of dependent endpoints with no dependencies).
		current := queue[0]

		// remove the first entry out the queue (of dependent endpoints with no dependencies).
		queue = queue[1:]

		// add the entry to the output.
		output = append(output, current)

		// decrease the amount of endpoints remaining.
		delete(endpoints, current)
		numEndpoints--

		// The operation above removes any amount of dependencies from the queue.
		//
		// add endpoints with no dependencies to the queue.
		for endpoint, dependencies := range endpoints {
			// when a dependency is added to the output,
			// remove that endpoint (current) from dependent (endpoint)s' dependencies.
			if contains(dependencies, current) {
				dependents[endpoint]--

				// when the endpoint is dependent on no other endpoints, add it to the queue.
				if dependents[endpoint] == 0 {
					queue = append(queue, endpoint)
				}
			}
		}
	}

	if numEndpoints != 0 {
		fmt.Println("WARNING: dependency cycle occurred (i.e [a: b],[b: a]) or necessary endpoint is unused.\n")
		fmt.Println("Examine the following endpoints.")

		for endpoint, dependencies := range endpoints {
			fmt.Println(endpoint, "depends on", dependencies)
		}

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
	for i, endpoint := range filterOutput(findOrder(endpoints)) {
		fmt.Printf("%d. %v\n", i, endpoint)
	}
}
